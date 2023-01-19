package queue

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	// Packages
	multierror "github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type queue struct {
	Pool

	// Parameters
	namespace   string        // Queue namespace
	backoff     time.Duration // Retry backoff duration
	max_age     time.Duration // Maximum task age
	max_retries uint          // Maximum number of retries
	workers     uint          // Maximum number of workers
	deadline    time.Duration // Deadline for any task work
}

var _ TaskQueue = (*queue)(nil)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultBackoff = 10 * time.Second
	defaultRetries = 10
	defaultDelta   = time.Second
)

var (
	defaultWorkers = uint(runtime.NumCPU())
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new queue with the given namespace
func NewQueue(pool Pool, opts ...Option) TaskQueue {
	queue := new(queue)
	if pool == nil {
		return nil
	} else {
		queue.Pool = pool
	}

	// Set some defaults
	queue.backoff = defaultBackoff
	queue.workers = defaultWorkers
	queue.max_retries = defaultRetries

	// Apply options
	for _, opt := range opts {
		if err := opt(queue); err != nil {
			return nil
		}
	}

	return queue
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (queue *queue) String() string {
	str := "<queue"
	if queue.Pool != nil {
		str += fmt.Sprint(" pool=", queue.Pool)
	}
	if queue.namespace != "" {
		str += fmt.Sprintf(" namespace=%q", queue.namespace)
	}
	if queue.backoff > 0 {
		str += fmt.Sprint(" retry_backoff=", queue.backoff)
	}
	if queue.max_age > 0 {
		str += fmt.Sprint(" max_age=", queue.max_age)
	}
	if queue.max_retries > 0 {
		str += fmt.Sprint(" max_retries=", queue.max_retries)
	}
	if queue.workers > 0 {
		str += fmt.Sprint(" workers=", queue.workers)
	}
	if queue.deadline > 0 {
		str += fmt.Sprint(" task_deadline=", queue.deadline)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run the queue
func (queue *queue) Run(ctx context.Context, fn WorkerFunc) error {
	var result error
	var wg sync.WaitGroup
	var ch = make(chan Task, queue.workers)

	// Check parameters
	if fn == nil {
		return ErrBadParameter.With("Run")
	} else if err := ctx.Err(); err != nil {
		return err
	}

	// Spin up workers
	for i := uint(0); i < queue.workers; i++ {
		wg.Add(1)
		go func(i uint, ch <-chan Task, fn WorkerFunc) {
			defer wg.Done()
			queue.run(ctx, ch, fn)
		}(i, ch, fn)
	}

	// Wait for context to be cancelled
	timer := time.NewTimer(0)
	defer timer.Stop()
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		case <-timer.C:
			task, err := queue.Retain(ctx)
			if err != nil {
				fmt.Println("TODO:", err)
			} else if task != nil {
				ch <- task
			}
			timer.Reset(defaultDelta)
		}
	}
	// TODO: Context will retieve tasks from the queue and send them to the channel
	// where workers will pick them up. Wait until context is cancelled, then close
	<-ctx.Done()

	// Close channel
	close(ch)

	// Wait for workers to finish
	wg.Wait()

	// Return any errors
	return result
}

// Schedule a new task to be executed
func (queue *queue) New(ctx context.Context, tag ...Tag) (Task, error) {
	var result error

	// Get a connection from the pool
	conn := queue.Pool.Get()
	defer queue.Pool.Put(conn)

	// Create task
	task := NewTask(queue.namespace)

	// Set tags for the task, report any errors
	for _, tag := range tag {
		if err := task.set(tag.Type, tag.Value); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if result != nil {
		return nil, result
	}

	// Set scheduled_at if not set
	if task.Get(TaskScheduledAt).(time.Time).IsZero() {
		if err := task.set(TaskScheduledAt, task.Get(TaskCreatedAt)); err != nil {
			return nil, err
		}
	}

	// Set expires_at if max_age is set
	if task.Get(TaskExpiresAt).(time.Time).IsZero() && queue.max_age > 0 {
		if err := task.set(TaskExpiresAt, task.Get(TaskCreatedAt).(time.Time).Add(queue.max_age)); err != nil {
			return nil, err
		}
	}

	// Save task in the queue
	if err := conn.Insert(ctx, task); err != nil {
		return nil, err
	}

	// Return success
	return task, nil
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// run is a worker function which performs tasks
func (queue *queue) run(parent context.Context, ch <-chan Task, fn WorkerFunc) {
	// Accept tasks until the channel is closed
	for task := range ch {
		// Create a context with a deadline
		ctx, cancel := taskctx(parent, queue.deadline)

		// Run the task
		if err := fn(ctx, task); err != nil {
			fmt.Println("TODO: Fail ", task, " with error ", err)
		} else {
			fmt.Println("TODO: Release task", task)
		}

		// Call cancel
		cancel()
	}
}

// taskctx can return a context with a deadline if it is set
func taskctx(parent context.Context, deadline time.Duration) (context.Context, context.CancelFunc) {
	if deadline > 0 {
		return context.WithTimeout(parent, deadline)
	} else {
		return parent, func() {}
	}
}

// Retain a task from the queue
func (queue *queue) Retain(ctx context.Context, filter ...Filter) (Task, error) {
	// Get a connection from the pool
	conn := queue.Pool.Get()
	defer queue.Pool.Put(conn)
	if conn == nil {
		return nil, ErrOutOfOrder.With("unable to establish a connection")
	}

	// Sort by priority, then scheduled_at and limit to one task
	sort := conn.S()
	sort.Desc(string(TaskPriority))
	sort.Asc(string(TaskScheduledAt))

	// Filter by task
	task, err := conn.Collection(task{}).Find(ctx, sort, filter...)
	if errors.Is(err, ErrNotFound) || task == nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Return success
	return task.(Task), nil
}

/*
// Release a task with an error or with success. This will delete the task from
// the queue if there is no error, otherwise it will update the task with the
// error and increment the retry count.
func (queue *queue) Release(ctx context.Context, task Task, lastErr error) error {
	// Get a connection from the pool
	conn := queue.Pool.Get()
	defer queue.Pool.Put(conn)
	if conn == nil {
		return ErrOutOfOrder.With("unable to establish a connection")
	}

	// Filter by task
	filter := conn.F()
	if err := filter.Key(task.Key()); err != nil {
		return err
	}

	if lastErr == nil {
		// The case where the lastErr is nil
		if deleted, err := conn.Collection(task).Delete(ctx, filter); err != nil {
			return err
		} else if deleted != 1 {
			return ErrInternalAppError.With("expected to delete one task, got", deleted)
		}
	} else {
		return ErrNotImplemented
	}

	// Return success
	return nil
}

func (queue *queue) Run(ctx context.Context, fn TaskFunc) error {
	<-ctx.Done()
	return nil
}
*/
