package mongodb_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	// Packages
	"github.com/mutablelogic/go-accessory/pkg/mongodb"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	"github.com/stretchr/testify/assert"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

const (
	MONGO_URL = "${MONGO_URL}"
)

func Test_Client_001(t *testing.T) {
	assert := assert.New(t)
	c, err := mongodb.Open(context.TODO(), uri(t))
	assert.NoError(err)
	assert.NotNil(c)
	assert.NoError(c.Close())
}

func Test_Client_002(t *testing.T) {
	assert := assert.New(t)
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptTrace(func(ctx context.Context, delta time.Duration) {
		t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", delta)
	}))
	assert.NoError(err)
	defer c.Close()

	// Ping
	assert.NoError(c.Ping(context.TODO()))
}

func Test_Client_003(t *testing.T) {
	assert := assert.New(t)

	// Add default database option
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptDatabase("admin"))
	assert.NoError(err)
	defer c.Close()
	assert.Equal("admin", c.Name())
}

func Test_Client_004(t *testing.T) {
	assert := assert.New(t)

	// Add default timeout option
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptTimeout(5*time.Second))
	assert.NoError(err)
	defer c.Close()

	assert.Equal(5*time.Second, c.Timeout())
}

func Test_Client_005(t *testing.T) {
	assert := assert.New(t)

	// Select specific database
	c, err := mongodb.Open(context.TODO(), uri(t))
	assert.NoError(err)
	defer c.Close()

	db := c.Database("test")
	assert.NotNil(db)
	assert.Equal("test", db.Name())
}

func Test_Client_006(t *testing.T) {
	assert := assert.New(t)

	// Run in a transaction
	c, err := mongodb.Open(context.TODO(), uri(t))
	assert.NoError(err)
	defer c.Close()

	// No error
	assert.NoError(c.Do(context.TODO(), func(ctx context.Context) error {
		return nil
	}))

	// Error
	assert.Error(ErrNotImplemented, c.Do(context.TODO(), func(ctx context.Context) error {
		return ErrNotImplemented
	}))

}

func Test_Client_007(t *testing.T) {
	assert := assert.New(t)

	// List Databases
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptTrace(func(ctx context.Context, delta time.Duration) {
		t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", delta)
	}))
	assert.NoError(err)
	defer c.Close()

	databases, err := c.Databases(context.TODO())
	assert.NoError(err)
	assert.NotNil(databases)
	assert.NotEmpty(databases)
}

func Test_Client_008(t *testing.T) {
	assert := assert.New(t)

	type Doc struct {
		Key  string `bson:"_id,omitempty"`
		Name string `bson:"name"`
	}

	// List Databases
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptDatabase("test"), mongodb.OptTrace(func(ctx context.Context, delta time.Duration) {
		t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", delta)
	}))
	assert.NoError(err)
	defer c.Close()

	t.Run("001", func(t *testing.T) {
		err := c.Insert(context.TODO(), Doc{Name: "Test"})
		assert.NoError(err)
	})

	t.Run("002", func(t *testing.T) {
		err := c.Insert(context.TODO(), Doc{Name: "T1"}, Doc{Name: "T2"}, nil)
		assert.Error(err)
	})

	t.Run("003", func(t *testing.T) {
		a, b := Doc{}, Doc{}
		err := c.Insert(context.TODO(), &a, &b)
		assert.NoError(err)
		assert.NotEmpty(a.Key)
		assert.NotEmpty(b.Key)
	})

	t.Run("004", func(t *testing.T) {
		a, b := Doc{}, Doc{}

		err := c.Do(context.TODO(), func(ctx context.Context) error {
			err := c.Insert(ctx, &a, &b)
			assert.NoError(err)
			assert.NotEmpty(a.Key)
			assert.NotEmpty(b.Key)
			return nil
		})
		assert.NoError(err)

	})

	t.Run("005", func(t *testing.T) {
		a, b := Doc{}, Doc{}

		err := c.Do(context.TODO(), func(ctx context.Context) error {
			err := c.Insert(ctx, &a, &b)
			assert.NoError(err)
			assert.NotEmpty(a.Key)
			assert.NotEmpty(b.Key)
			return ErrBadParameter
		})
		assert.Error(err)

	})
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
