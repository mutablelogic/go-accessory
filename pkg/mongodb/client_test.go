package mongodb_test

import (
	"context"
	"os"
	"testing"
	"time"

	// Packages
	"github.com/mutablelogic/go-accessory/pkg/mongodb"
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
	c, err := mongodb.Open(context.TODO(), uri(t))
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
	c, err := mongodb.Open(context.TODO(), uri(t))
	assert.NoError(err)
	defer c.Close()

	databases, err := c.Databases(context.TODO())
	assert.NoError(err)
	assert.NotNil(databases)
	t.Log(databases)
}

/*
func Test_Client_008(t *testing.T) {
	assert := assert.New(t)

	// List Databases
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptDatabase("admin"))
	assert.NoError(err)
	defer c.Close()

	collections, err := c.Collections(context.TODO())
	assert.NoError(err)
	assert.NotNil(collections)
	t.Log(collections)
}

func Test_Client_009(t *testing.T) {
	assert := assert.New(t)

	type Doc struct {
		Name string `bson:"name"`
	}

	// List Databases
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptDatabase("test"), mongodb.OptCollection(Doc{}, "doc"))
	assert.NoError(err)
	defer c.Close()

	collections, err := c.Collections(context.TODO())
	assert.NoError(err)
	assert.NotNil(collections)
	t.Log(collections)
}
*/

///////////////////////////////////////////////////////////////////////////////
// Return URL or skip test

func uri(t *testing.T) string {
	if uri := os.ExpandEnv(MONGO_URL); uri == "" {
		t.Skip("no MONGO_URL environment variable, skipping test")
		return ""
	} else {
		return uri
	}
}
