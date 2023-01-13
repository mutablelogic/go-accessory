package accessory

import (
	"context"
	"fmt"
)

///////////////////////////////////////////////////////////////////////////////
// TYPE

type TagType string

// Tag represents some task metadata, including Priority and ScheduledAt
type Tag struct {
	Type  TagType
	Value any
}

// TaskFunc is used for executing tasks. When an error is returned, the task
// should be retried at some later date, or released if the error indicates
// the task should not be retried.
type TaskFunc func(context.Context, Task) error

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// TaskQueue represents a set of tasks to be executed in order.
// Create a TaskQueue using:
//
//	queue := taskqueue.NewQueue(client, namespace)
type TaskQueue interface {
	// Schedule a new task to be executed, and return it
	New(context.Context, ...Tag) (Task, error)

	// Set metadata tag values. Delete a tag if value set to nil
	//Set(context.Context, Task, ...Tag) error

	// Run the queue to retain tasks and execute them
	Run(context.Context, TaskFunc) error
}

// Task represents a task
type Task interface {
	Key() string       // A unique identifier for the task
	Namespace() string // Return the namespace of the task
	Tags() []Tag       // Return all metadata tags

	// Get a metadata tag value
	//Get(TagType) any
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	TaskPriority    TagType = "priority"     // int: The priority of the task (higher is more important)
	TaskScheduledAt TagType = "scheduled_at" // time.Time: The time the task is scheduled to be executed
	TaskExpiresAt   TagType = "expires_at"   // time.Time: When the task expires (if not executed before this time)
	TaskAge         TagType = "age"          // time.Duration: The maximum age of the task (how long it has been in the queue)
	TaskRetryCount  TagType = "retry_count"  // int: The number of times the task has been retried
	TaskLastError   TagType = "last_error"   // string: The last task error
)

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t Tag) String() string {
	str := "<tag"
	str += fmt.Sprintf(" type=%q", t.Type)
	if t.Value != nil {
		switch v := t.Value.(type) {
		case string:
			str += fmt.Sprintf(" value=%q", v)
		default:
			str += fmt.Sprint(" value=", t.Value)
		}
	}
	return str + ">"
}
