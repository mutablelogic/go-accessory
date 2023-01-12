package taskqueue

import (
	"context"
	"fmt"
	"time"

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

	// Maximum retry count
	retry_count uint

	// Retry backoff duration
	retry_delta time.Duration
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultRetryCount = 6
	defaultRetryDelta = 10 * time.Second
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new queue with the given namespace
func NewQueue(client Client, namespace string) *queue {
	return NewQueueWithRetry(client, namespace, 0, 0)
}

// Create a new queue with the given namespace, retry count, and retry delta
func NewQueueWithRetry(client Client, namespace string, retry_count uint, retry_delta time.Duration) *queue {
	queue := new(queue)
	queue.Client = client
	queue.namespace = namespace
	if retry_count == 0 {
		queue.retry_count = defaultRetryCount
	} else {
		queue.retry_count = retry_count
	}
	if retry_delta <= 0 {
		queue.retry_delta = defaultRetryDelta
	} else {
		queue.retry_delta = retry_delta
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
	str += fmt.Sprint(" retry_count=", queue.retry_count)
	str += fmt.Sprint(" retry_delta=", queue.retry_delta)
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Push a task into the queue
func (queue *queue) Push(ctx context.Context, task *task) error {
	return ErrNotImplemented
}
