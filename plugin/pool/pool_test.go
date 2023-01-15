package main_test

import (
	"context"
	"testing"

	pool "github.com/mutablelogic/go-accessory/plugin/pool"
	task "github.com/mutablelogic/go-server/pkg/task"
	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Pool_001(t *testing.T) {
	// Create a provider, register dnsregister
	assert := assert.New(t)
	pool := pool.WithLabel(t.Name()).WithUrl("mongodb://cm1")
	provider, err := task.NewProvider(context.Background(), pool)
	assert.NoError(err)
	assert.NotNil(provider)
}
