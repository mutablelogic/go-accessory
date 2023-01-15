package queue

import (
	"context"
	"fmt"
	"time"

	// Packages
	multierror "github.com/hashicorp/go-multierror"

	// Namespace imports
	//. "github.com/djthorpe/go-errors"
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

func (queue *queue) Run(ctx context.Context, fn TaskFunc) error {
	<-ctx.Done()
	return nil
}
