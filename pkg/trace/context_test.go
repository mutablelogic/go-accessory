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
	ctx := trace.WithTx(context.TODO())
	ctx = trace.WithUrl(ctx, trace.OpConnect, url, 0)
	t.Log(trace.DumpContextStr(ctx))
}
