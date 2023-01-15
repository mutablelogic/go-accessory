package taskqueue

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// task is a task in the queue
type task struct {
	Id          string          `bson:"_id,omitempty"`
	Name        string          `bson:"namespace,omitempty"`
	Priority    int             `bson:"pri,omitempty"`
	CreatedAt   time.Time       `bson:"created_at,omitempty"`
	ScheduledAt time.Time       `bson:"scheduled_at,omitempty"`
	ExpiresAt   time.Time       `bson:"expires_at,omitempty"`
	RetryCount  uint            `bson:"retry_count,omitempty"`
	LastError   string          `bson:"last_error,omitempty"`
	Tag         map[TagType]any `bson:"tag,omitempty"`
}

var _ Task = (*task)(nil)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	reTagName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]+$`)
)

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
	task.Name = namespace
	task.Priority = priority
	task.CreatedAt = time.Now()
	task.ExpiresAt = expires_at
	task.Tag = make(map[TagType]any)
	return task
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (task *task) String() string {
	str := "<taskqueue.task"
	if task.Id != "" {
		str += fmt.Sprintf(" key=%q", task.Id)
	}
	if task.Name != "" {
		str += fmt.Sprintf(" namespace=%q", task.Name)
	}
	for _, tag := range task.Tags() {
		switch v := tag.Value.(type) {
		case string:
			str += fmt.Sprintf(" %s=%q", tag.Type, v)
		case time.Time:
			if v.IsZero() {
				str += fmt.Sprintf(" %s=nil", tag.Type)
			} else {
				str += fmt.Sprintf(" %s=%q", tag.Type, v.Format(time.RFC3339))
			}
		default:
			str += fmt.Sprintf(" %s=%v", tag.Type, v)
		}
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (task *task) Key() string {
	return task.Id
}

func (task *task) Namespace() string {
	return task.Name
}

func (task *task) Tags() []Tag {
	var results []Tag
	// Set priority
	if task.Priority != 0 {
		results = append(results, Tag{TaskPriority, task.Priority})
	}
	// Get age
	if !task.CreatedAt.IsZero() {
		results = append(results, Tag{TaskAge, time.Since(task.CreatedAt).Truncate(time.Millisecond)})
	}
	// ScheduledAt
	if !task.ScheduledAt.IsZero() {
		results = append(results, Tag{TaskScheduledAt, task.ScheduledAt})
	}
	// ExpiresAt
	if !task.ExpiresAt.IsZero() {
		results = append(results, Tag{TaskExpiresAt, task.ExpiresAt})
	}
	// RetryCount
	if task.RetryCount != 0 {
		results = append(results, Tag{TaskRetryCount, task.RetryCount})
	}
	// LastError
	if task.LastError != "" {
		results = append(results, Tag{TaskLastError, errors.New(task.LastError)})
	}
	// Other Tags
	for k, v := range task.Tag {
		results = append(results, Tag{k, v})
	}
	// Return tags
	return results
}

func (task *task) Get(t TagType) any {
	switch t {
	case TaskPriority:
		return task.Priority_
	case TaskScheduledAt:
		return task.ScheduledAt_
	case TaskExpiresAt:
		return task.ExpiresAt_
	case TaskAge:
		if task.CreatedAt_.IsZero() {
			return nil
		} else {
			return time.Since(task.CreatedAt_).Truncate(time.Millisecond)
		}
	case TaskRetryCount:
		return task.RetryCount_
	case TaskLastError:
		if task.LastError_ == "" {
			return nil
		} else {
			return errors.New(task.LastError_)
		}
	default:
		if value, ok := task.Tags_[t]; ok {
			return value
		} else {
			return nil
		}
	}
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
			return ErrBadParameter.With("priority must be an int value")
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
		if !reTagName.MatchString(string(key)) {
			return ErrBadParameter.Withf("invalid tag name %q", key)
		} else if value == nil {
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
