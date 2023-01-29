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

func Test_FindUpdate_001(t *testing.T) {
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
		assert.NotEmpty(doc.Key)
	})

	t.Run("002", func(t *testing.T) {
		filter := c.F()
		filter.Key(doc.Key)
		doc, err := c.Collection(Doc{}).FindUpdate(context.TODO(), Doc{Name: "NewName"}, nil, filter)
		assert.NoError(err)
		assert.NotNil(doc)
		assert.Equal("Test", doc.(*Doc).Name)

		doc2, err := c.Collection(Doc{}).Find(context.TODO(), nil, filter)
		assert.NoError(err)
		assert.NotNil(doc2)
		assert.Equal("NewName", doc2.(*Doc).Name)

		t.Log(doc, doc2)
	})
}
