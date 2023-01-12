package taskqueue

import "time"

///////////////////////////////////////////////////////////////////////////////
// TYPES

// task is a task in the queue
type task struct {
	Key         string            `bson:"_id,omitempty"`
	Namespace   string            `bson:"namespace,omitempty"`
	Priority    int               `bson:"pri,omitempty"`
	ScheduledAt time.Time         `bson:"scheduled_at,omitempty"`
	ExpiresAt   time.Time         `bson:"expires_at,omitempty"`
	RetryCount  uint              `bson:"retry_count,omitempty"`
	LastError   string            `bson:"last_error,omitempty"`
	Tags        map[string]string `bson:"tags,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewTask creates a new task with the given namespace, priority, and when the task
// should expire. The task is scheduled to run immediately.
func NewTask(namespace string) *task {
	return NewTaskWithPriority(namespace, 0)
}

// NewTaskWithPriority creates a new task with the given namespace and priority
func NewTaskWithPriority(namespace string, priority int) *task {
	return NewTaskWithPriorityAndExpiresAt(namespace, priority, time.Time{})
}

// NewTaskWithPriorityAndExpiresAt creates a new task with the given namespace, priority, and when the task
// should expire.
func NewTaskWithPriorityAndExpiresAt(namespace string, priority int, expires_at time.Time) *task {
	task := new(task)
	task.Namespace = namespace
	task.Priority = priority
	task.ExpiresAt = expires_at
	task.Tags = make(map[string]string)
	return task
}
