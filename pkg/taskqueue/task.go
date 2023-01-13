package taskqueue

import (
	"fmt"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// task is a task in the queue
type task struct {
	Key_         string          `bson:"_id,omitempty"`
	Namespace_   string          `bson:"namespace,omitempty"`
	Priority_    int             `bson:"pri,omitempty"`
	CreatedAt_   time.Time       `bson:"created_at,omitempty"`
	ScheduledAt_ time.Time       `bson:"scheduled_at,omitempty"`
	ExpiresAt_   time.Time       `bson:"expires_at,omitempty"`
	RetryCount_  uint            `bson:"retry_count,omitempty"`
	LastError_   string          `bson:"last_error,omitempty"`
	Tags_        map[TagType]any `bson:"tags,omitempty"`
}

var _ Task = (*task)(nil)

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
	task.Namespace_ = namespace
	task.Priority_ = priority
	task.CreatedAt_ = time.Now()
	task.ExpiresAt_ = expires_at
	task.Tags_ = make(map[TagType]any)
	return task
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (task *task) String() string {
	str := "<taskqueue.task"
	if task.Key_ != "" {
		str += fmt.Sprintf(" key=%q", task.Key_)
	}
	if task.Namespace_ != "" {
		str += fmt.Sprintf(" namespace=%q", task.Namespace_)
	}
	for _, tag := range task.Tags() {
		str += fmt.Sprintf(" %s=%q", tag.Type, tag.Value)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (task *task) Key() string {
	return task.Key_
}

func (task *task) Namespace() string {
	return task.Namespace_
}

func (task *task) Tags() []Tag {
	var results []Tag
	// Set priority
	if task.Priority_ != 0 {
		results = append(results, Tag{TaskPriority, task.Priority_})
	}
	// Get age
	if !task.CreatedAt_.IsZero() {
		results = append(results, Tag{TaskAge, time.Since(task.CreatedAt_).Truncate(time.Millisecond)})
	}
	// ScheduledAt
	if !task.ScheduledAt_.IsZero() {
		results = append(results, Tag{TaskScheduledAt, task.ScheduledAt_})
	}
	// ExpiresAt
	if !task.ExpiresAt_.IsZero() {
		results = append(results, Tag{TaskExpiresAt, task.ExpiresAt_})
	}
	// RetryCount
	if task.RetryCount_ != 0 {
		results = append(results, Tag{TaskRetryCount, task.RetryCount_})
	}
	// LastError
	if task.LastError_ != "" {
		results = append(results, Tag{TaskLastError, task.LastError_})
	}
	// Other Tags
	for k, v := range task.Tags_ {
		results = append(results, Tag{k, v})
	}
	// Return tags
	return results
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (task *task) set(key TagType, value any) error {
	switch key {
	case TaskPriority:
		if v, ok := value.(int); ok {
			task.Priority_ = v
			return nil
		} else {
			return ErrBadParameter.With("priority must be an integer value")
		}
	case TaskExpiresAt:
		if value == nil {
			task.ExpiresAt_ = time.Time{}
			return nil
		} else if v, ok := value.(time.Time); ok {
			task.ExpiresAt_ = v
			return nil
		} else {
			return ErrBadParameter.With("expires_at must be a time.Time value")
		}
	case TaskLastError:
		if value == nil {
			task.LastError_ = ""
			return nil
		} else if v, ok := value.(string); ok {
			task.LastError_ = v
			return nil
		} else {
			return ErrBadParameter.With("last_error must be a string value")
		}
	case TaskRetryCount, TaskScheduledAt, TaskAge:
		return ErrBadParameter.With("cannot set tag with type: ", key)
	default:
		if value == nil {
			delete(task.Tags_, key)
			return nil
		} else if v, ok := value.(string); ok {
			task.Tags_[key] = v
			return nil
		} else {
			return ErrBadParameter.With("tag must be a string value")
		}
	}
}
