package taskqueue

import (
	"context"
	"fmt"
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
	Client

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
func NewQueue(client Client, namespace string) TaskQueue {
	return NewQueueWithDelta(client, namespace, 0)
}

// Create a new queue with the given namespace
func NewQueueWithDelta(client Client, namespace string, delta time.Duration) TaskQueue {
	queue := new(queue)
	queue.Client = client
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
	str := "<taskqueue"
	if queue.Client != nil {
		str += fmt.Sprint(" client=", queue.Client)
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
	if task.ScheduledAt_.IsZero() {
		task.ScheduledAt_ = task.CreatedAt_
	}

	// Store in the queue
	if err := queue.Client.Insert(ctx, task); err != nil {
		return nil, err
	}

	// Return success
	return task, nil
}

// Process any tasks in the queue. The function fn is called for each
// task in the queue. The function should return nil if the task is
// successfully processed, or an error if the task should be retried.
func (queue *queue) Run(ctx context.Context, fn TaskFunc) error {
	var result error

	timer := time.NewTimer(0)
	defer timer.Stop()

	// Loop until cancelled
	for {
		select {
		case <-ctx.Done():
			// Return any errors
			return result
		case <-timer.C:
			task, err := queue.Retain(ctx)
			if err != nil {
				result = multierror.Append(result, err)
				continue
			}
			if task == nil {
				// Pass nil to fn to indicate no task
				fn(ctx, nil)
				timer.Reset(queue.delta)
			} else {
				err := fn(ctx, task)
				fmt.Println("task result:", err)
				timer.Reset(0)
			}
		}
	}
}

// Retain the next task to be executed
func (queue *queue) Retain(ctx context.Context) (Task, error) {
	// Sort
	sort := queue.Client.S()
	sort.Desc("priority")
	sort.Asc("scheduled_at")
	sort.Limit(1)

	// Filter
	filter := queue.Client.F()
	if queue.namespace != "" {
		filter.Eq("namespace", queue.namespace)
	}

	// Task collection
	collection := queue.Collection(task{})
	if collection == nil {
		return nil, ErrInternalAppError.With("task collection not found")
	}

	// Find task, and update
	var result Task
	err := queue.Do(ctx, func(ctx context.Context) error {
		item, err := collection.Find(ctx, sort, filter)
		if err != nil {
			return err
		}
		if item == nil {
			return ErrInternalAppError.With("no item returned")
		}
		// Update task to set scheduled_at to nil
		if matched, modified, err := collection.Update(ctx, task{ScheduledAt_: time.Time{}}); err != nil {
			return err
		} else if matched != 1 || modified != 1 {
			return ErrInternalAppError.With("no task updated")
		}
		// Return success
		return nil
	})

	fmt.Println("RETAIN:", result, err)
	return result, err
	// Sort:
	// Where scheduledAt is NOT nill
	// and scheduledAt <= now
	// and ExpiresAt > now
	// Update:
	// set ScheduledAt is nil
}

// Release a task. When the error is nil, the task is released from
// the task queue. When the error is non-nil, the last may be retried.
// The return value may indicate the task will not
// be retried.
func (queue *queue) Release(context.Context, Task, error) error {
	return ErrNotImplemented.With(queue)
}
