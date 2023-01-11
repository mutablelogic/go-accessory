package mongodb_test

import (
	"context"
	"testing"

	// Packages
	"github.com/mutablelogic/go-accessory/pkg/mongodb"
	"github.com/stretchr/testify/assert"
)

func Test_Insert_001(t *testing.T) {
	assert := assert.New(t)

	type Doc struct {
		Name string `bson:"name"`
	}

	// Open database connection and register Doc collection
	c, err := mongodb.Open(context.TODO(), uri(t), mongodb.OptDatabase("test"), mongodb.OptCollection(Doc{}, "doc"))
	assert.NoError(err)
	defer c.Close()

	// Insert a single document
	key, err := c.Insert(context.TODO(), Doc{Name: "test"})
	assert.NoError(err)
	assert.NotEmpty(key)
}
