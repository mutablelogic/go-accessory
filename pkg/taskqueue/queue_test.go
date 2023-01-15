package taskqueue_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	// Packages
	mongodb "github.com/mutablelogic/go-accessory/pkg/mongodb"
	taskqueue "github.com/mutablelogic/go-accessory/pkg/taskqueue"
	assert "github.com/stretchr/testify/assert"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

const (
	MONGO_URL = "${MONGO_URL}"
)

func Test_Queue_001(t *testing.T) {
	assert := assert.New(t)
	c, err := mongodb.Open(context.TODO(), uri(t))
	assert.NoError(err)
	assert.NotNil(c)
	defer assert.NoError(c.Close())

	queue := taskqueue.NewQueue(c, "test")
	assert.NotNil(queue)
	t.Log(queue)
}

func Test_Queue_002(t *testing.T) {
	assert := assert.New(t)
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptDatabase("test"))
	assert.NoError(err)
	assert.NotNil(c)

	// Create a queue
	queue := taskqueue.NewQueue(c, "test")
	assert.NotNil(queue)

	// Create N tasks
	for i := 0; i < 10; i++ {
		task, err := queue.New(context.TODO(), Tag{Type: TaskPriority, Value: i})
		assert.NoError(err)
		assert.NotNil(task)
		assert.NotEmpty(task.Key())
		assert.Equal(i, task.Get(TaskPriority).(int))
		t.Log(i, "=>", task)
	}

	// Set the task priority

	assert.NoError(c.Close())
}

func Test_Queue_003(t *testing.T) {
	assert := assert.New(t)
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptDatabase("test"))
	assert.NoError(err)
	assert.NotNil(c)

	// Create a queue
	queue := taskqueue.NewQueue(c, "test")
	assert.NotNil(queue)

	// Read tasks for 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Read tasks until no more tasks
	assert.NoError(queue.Run(ctx, func(ctx context.Context, task Task) error {
		t.Log("TODO: Run task", task)
		return nil
	}))
}

///////////////////////////////////////////////////////////////////////////////
// Return URL or skip test

func uri(t *testing.T) *url.URL {
	if uri := os.ExpandEnv(MONGO_URL); uri == "" {
		t.Skip("no MONGO_URL environment variable, skipping test")
	} else if uri, err := url.Parse(uri); err != nil {
		t.Skip("invalid MONGO_URL environment variable, skipping test")
	} else {
		return uri
	}
	return nil
}
