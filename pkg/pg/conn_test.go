package pg_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	// Packages
	pg "github.com/mutablelogic/go-accessory/pkg/pg"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	assert "github.com/stretchr/testify/assert"
)

const (
	PG_URL = "${PG_URL}"
)

func Test_Client_001(t *testing.T) {
	assert := assert.New(t)
	c, err := pg.Open(context.TODO(), uri(t), tracefn(t), pg.OptApplicationName(t.Name()))
	assert.NoError(err)
	defer func() {
		if c != nil {
			assert.NoError(c.Close())
		}
	}()
}

func Test_Client_002(t *testing.T) {
	assert := assert.New(t)
	c, err := pg.Open(context.TODO(), uri(t), tracefn(t), pg.OptApplicationName(t.Name()))
	assert.NoError(err)
	if c == nil {
		t.Skip("No client")
	}
	defer func() {
		if c != nil {
			assert.NoError(c.Close())
		}
	}()

	// Ping the server
	assert.NoError(c.Ping(context.Background()))
}

func Test_Client_003(t *testing.T) {
	assert := assert.New(t)
	c, err := pg.Open(context.TODO(), uri(t), tracefn(t), pg.OptApplicationName(t.Name()))
	assert.NoError(err)
	if c == nil {
		t.Skip("No client")
	}
	defer func() {
		if c != nil {
			assert.NoError(c.Close())
		}
	}()

	databases, err := c.Databases(context.Background())
	assert.NoError(err)
	assert.NotEmpty(databases)
	t.Log(databases)
}

func Test_Client_004(t *testing.T) {
	assert := assert.New(t)
	c, err := pg.Open(context.TODO(), uri(t), tracefn(t), pg.OptApplicationName(t.Name()))
	assert.NoError(err)
	if c == nil {
		t.Skip("No client")
	}
	defer func() {
		if c != nil {
			assert.NoError(c.Close())
		}
	}()

	assert.NoError(c.Do(context.Background(), func(ctx context.Context) error {
		return nil
	}))
}

///////////////////////////////////////////////////////////////////////////////
// Utility Methods

func uri(t *testing.T) *url.URL {
	if uri := os.ExpandEnv(PG_URL); uri == "" {
		t.Skip("no PG_URL environment variable, skipping test")
	} else if uri, err := url.Parse(uri); err != nil {
		t.Skip("invalid PG_URL environment variable, skipping test")
	} else {
		return uri
	}
	return nil
}

func tracefn(t *testing.T) pg.Opt {
	return pg.OptTrace(func(ctx context.Context, delta time.Duration, err error) {
		if err != nil {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", err)
		} else {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", delta)
		}
	})
}
