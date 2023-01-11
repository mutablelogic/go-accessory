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
	. "github.com/djthorpe/go-errors"
)

func Test_Insert_001(t *testing.T) {
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
