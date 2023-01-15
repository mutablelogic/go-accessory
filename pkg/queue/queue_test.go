package queue_test

import (
	"context"
	"net/url"
	"os"
	"testing"

	// Packages
	pool "github.com/mutablelogic/go-accessory/pkg/pool"
	queue "github.com/mutablelogic/go-accessory/pkg/queue"
	assert "github.com/stretchr/testify/assert"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

const (
	MONGO_URL = "${MONGO_URL}"
)

func Test_Queue_001(t *testing.T) {
	assert := assert.New(t)
	pool := pool.New(context.TODO(), uri(t))
	queue := queue.NewQueue(pool, "test")
	assert.NotNil(queue)
	assert.NoError(pool.Close())
}

func Test_Queue_002(t *testing.T) {
	assert := assert.New(t)
	pool := pool.New(context.TODO(), uri(t))
	queue := queue.NewQueue(pool, "test")
	assert.NotNil(queue)

	// Create N tasks
	for i := 0; i < 10; i++ {
		task, err := queue.New(context.TODO(), Tag{TaskPriority, i})
		assert.NoError(err)
		assert.NotNil(task)
	}

	assert.NoError(pool.Close())
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
