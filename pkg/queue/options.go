package queue

import (
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Option func(*queue) error

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set the queue namespace
func OptNamespace(v string) Option {
	return func(queue *queue) error {
		queue.namespace = v
		return nil
	}
}

// Set the database name
func OptMaxAge(v time.Duration) Option {
	return func(queue *queue) error {
		if v > 0 {
			queue.max_age = v
		} else {
			return ErrBadParameter.With("OptMaxAge")
		}
		return nil
	}
}

// Set the maximum number of retries
func OptMaxRetries(v uint) Option {
	return func(queue *queue) error {
		if v > 0 {
			queue.max_retries = v
		} else {
			return ErrBadParameter.With("OptMaxRetries")
		}
		return nil
	}
}

// Set the maximum number of workers and work deadline
func OptWorkers(n uint, deadline time.Duration) Option {
	return func(queue *queue) error {
		if n > 0 {
			queue.workers = n
		} else {
			return ErrBadParameter.With("OptWorkers")
		}
		if deadline > 0 {
			queue.deadline = deadline
		} else {
			return ErrBadParameter.With("OptWorkers")
		}
		return nil
	}
}

// Set the retry backoff value. The first retry is attempted after
// the backoff value, the second retry is attempted after 2*backoff,
// and so forth, until either the task age is reached or the maximum
// number of retires is reached.
func OptBackoff(v time.Duration) Option {
	return func(queue *queue) error {
		if v > 0 {
			queue.backoff = v
		} else {
			return ErrBadParameter.With("OptBackoff")
		}
		return nil
	}
}
