package trace_test

import (
	"context"
	"net/url"
	"testing"

	// Packages
	"github.com/mutablelogic/go-accessory/pkg/trace"
	"github.com/stretchr/testify/assert"
)

func Test_Context_001(t *testing.T) {
	assert := assert.New(t)
	url, err := url.Parse("http://localhost:8080")
	assert.NoError(err)
	assert.NotNil(url)

	parent := context.Background()

	t.Run("001", func(t *testing.T) {
		ctx := trace.WithTx(parent)
		ctx = trace.WithUrl(ctx, trace.OpConnect, url)
		t.Log(trace.DumpContextStr(ctx))
	})
	t.Run("002", func(t *testing.T) {
		ctx := trace.WithTx(parent)
		ctx = trace.WithOp(ctx, trace.OpConnect)
		t.Log(trace.DumpContextStr(ctx))
	})
	t.Run("003", func(t *testing.T) {
		ctx := trace.WithTx(parent)
		ctx = trace.WithCollection(ctx, trace.OpCommit, "db", "collection")
		t.Log(trace.DumpContextStr(ctx))
	})
}
