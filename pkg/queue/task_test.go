package queue_test

import (
	"testing"
	"time"

	// Packages
	queue "github.com/mutablelogic/go-accessory/pkg/queue"
	assert "github.com/stretchr/testify/assert"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

func Test_Task_001(t *testing.T) {
	assert := assert.New(t)
	task := queue.NewTask("test")
	assert.NotNil(task)
	assert.Equal("", task.Key())
	assert.Equal("test", task.Namespace())
	t.Log(task)
}

func Test_Task_002(t *testing.T) {
	assert := assert.New(t)
	task := queue.NewTaskWithPriority("test", 100)
	assert.NotNil(task)
	assert.Equal(task.Get(TaskPriority).(int), 100)
	t.Log(task)
}

func Test_Task_003(t *testing.T) {
	assert := assert.New(t)
	expiry := time.Now().Add(time.Hour)
	task := queue.NewTaskWithPriorityAndExpiresAt("test", 100, expiry)
	assert.NotNil(task)
	assert.Equal(task.Get(TaskPriority).(int), 100)
	assert.Equal(task.Get(TaskExpiresAt).(time.Time), expiry)
	t.Log(task)
}
