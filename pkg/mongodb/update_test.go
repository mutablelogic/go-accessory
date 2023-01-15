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

func Test_Update_001(t *testing.T) {
	assert := assert.New(t)

	type Doc struct {
		Key  string `bson:"_id,omitempty"`
		Name string `bson:"name"`
	}

	// List Databases
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptDatabase("test"), mongodb.OptTrace(func(ctx context.Context, delta time.Duration, err error) {
		if err != nil {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", err)
		} else {
			t.Log("TRACE:", trace.DumpContextStr(ctx), "=>", delta)
		}
	}))
	assert.NoError(err)
	defer c.Close()

	doc := Doc{Name: "Test"}
	t.Run("001", func(t *testing.T) {
		err := c.Insert(context.TODO(), &doc)
		assert.NoError(err)
	})

	t.Run("002", func(t *testing.T) {
		filter := c.F()
		filter.Key(doc.Key)
		matched, modified, err := c.Collection(Doc{}).Update(context.TODO(), Doc{Name: "Test2"}, filter)
		assert.NoError(err)
		assert.Equal(int64(1), matched)
		assert.Equal(int64(1), modified)
	})
}
