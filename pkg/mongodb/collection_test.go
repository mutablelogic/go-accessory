package mongodb_test

import (
	"reflect"
	"testing"

	// Packages
	"github.com/djthorpe/go-mongodb/pkg/mongodb"
	"github.com/stretchr/testify/assert"
)

func Test_Collection_001(t *testing.T) {
	type Doc struct {
		Key string `bson:"_id,omitempty"`
		B   int    `bson:"b,unqiue: xxx"`
	}

	t.Run("000", func(t *testing.T) {
		assert := assert.New(t)
		collection := mongodb.NewCollection(reflect.TypeOf(Doc{}), "test")
		assert.NotNil(collection)
		assert.Equal("test", collection.Name)
		assert.Equal(reflect.TypeOf(Doc{}), collection.Type)
	})

	t.Run("001", func(t *testing.T) {
		assert := assert.New(t)
		collection := mongodb.NewCollection(reflect.TypeOf(&Doc{}), "test")
		assert.NotNil(collection)
		assert.Equal("test", collection.Name)
		assert.Equal(reflect.TypeOf(Doc{}), collection.Type)
	})

	t.Run("002", func(t *testing.T) {
		assert := assert.New(t)
		collection := mongodb.NewCollection(reflect.TypeOf(&Doc{}), "test")
		assert.NotNil(collection)
		// Field 0 is the key
		assert.Equal([]int{0}, collection.Key)
	})
}
