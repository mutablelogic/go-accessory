package taskqueue_test

import (
	"context"
	"net/url"
	"os"
	"testing"

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
		t.Log(i, "=>", task)
	}

	// Set the task priority
	//assert.NoError(task.Set(context.TODO(), Tag{TagPriority, 1}))

	assert.NoError(c.Close())
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
