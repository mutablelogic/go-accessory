package mongodb_test

import (
	"context"
	"testing"
	"time"

	// Packages
	mongodb "github.com/mutablelogic/go-accessory/pkg/mongodb"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	assert "github.com/stretchr/testify/assert"
	// Namespace imports
	//. "github.com/djthorpe/go-errors"
)

func Test_Find_001(t *testing.T) {
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
		doc, err := c.Collection(Doc{}).Find(context.TODO(), nil, nil)
		assert.NoError(err)
		assert.NotNil(doc)
		t.Log(doc)
	})

	t.Run("003", func(t *testing.T) {
		cursor, err := c.Collection(Doc{}).FindMany(context.TODO(), nil, nil)
		assert.NoError(err)
		assert.NotNil(cursor)
	})
}
