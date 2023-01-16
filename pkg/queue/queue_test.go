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
	pool := pool.New(context.TODO(), uri(t), pool.OptDatabase("test"))
	queue := queue.NewQueue(pool, "test")
	assert.NotNil(queue)

	// Create N tasks, from lowest to highest priority
	for i := 0; i < 10; i++ {
		task, err := queue.New(context.TODO(), Tag{TaskPriority, i})
		assert.NoError(err)
		assert.NotNil(task)
		assert.NotEmpty(task.Key())
		t.Log("Created task", task)
	}

	// Get N tasks and do work on them within a transaction. The task is retained
	// within the transaction, and the a goroutine will do the work and release the task
	// in another transaction. Need to keep track of how many workers are doing work
	// so we don't create too many workers.
	//
	// TODO: Filter should be:
	//   expires_at is zero or > now
	// and,
	//   scheduled_at is not zero and >= now
	// and,
	//   retry_count < max_retries (where max_retries is always bigger than zero)
	//
	// We need another cycle to expire old tasks:
	//   scheduled_at is not zero
	// and,
	//     expires_at is not zero and >= now (because the task expired or is too old)
	//   or,
	//     retry_count >= max_retries (because the task failed too many times)
	queue.Do(context.TODO(), func(ctx context.Context, task Task) error {
		// TODO: Retain task, do work, then release task, in a worker object
		t.Log("Releasing task", task)
		return queue.Release(ctx, task, nil)
	}, 0)

	// Close connection pool
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
