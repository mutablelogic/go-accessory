package queue

import (
	"context"
	"fmt"
	"io"
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

	// Queue namespace
	namespace string

	// Retry backoff duration
	delta time.Duration
}

var _ TaskQueue = (*queue)(nil)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultRetryDelta = 10 * time.Second
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new queue with the given namespace
func NewQueue(pool Pool, namespace string) TaskQueue {
	return NewQueueWithDelta(pool, namespace, 0)
}

// Create a new queue with the given namespace
func NewQueueWithDelta(pool Pool, namespace string, delta time.Duration) TaskQueue {
	queue := new(queue)
	queue.Pool = pool
	queue.namespace = namespace
	if delta == 0 {
		queue.delta = defaultRetryDelta
	} else {
		queue.delta = delta
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
	str += fmt.Sprint(" delta=", queue.delta)
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

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

	// Save task in the queue
	if err := conn.Insert(ctx, task); err != nil {
		return nil, err
	}

	// Return success
	return task, nil
}

// Perform some operation on up to "limit" tasks
func (queue *queue) Do(ctx context.Context, fn TaskFunc, limit int64, filter ...Filter) error {
	// Get a connection from the pool
	conn := queue.Pool.Get()
	defer queue.Pool.Put(conn)
	if conn == nil {
		return ErrOutOfOrder.With("unable to establish a connection")
	}

	// Sort by priority, then scheduled_at
	sort := conn.S()
	sort.Desc(string(TaskPriority))
	sort.Asc(string(TaskScheduledAt))
	if limit > 0 {
		sort.Limit(limit)
	}
	sort.Limit(limit)

	// Perform operations on tasks in a transaction
	return conn.Do(ctx, func(ctx context.Context) error {
		cursor, err := conn.Collection(task{}).FindMany(ctx, sort)
		if err != nil {
			return err
		}
		for {
			task, err := cursor.Next(ctx)
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			} else if err := fn(ctx, task.(Task)); err != nil {
				return err
			}
		}
		// Return success
		return nil
	})
}

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
