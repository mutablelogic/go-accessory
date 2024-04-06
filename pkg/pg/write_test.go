package pg_test

import (
	"context"
	"reflect"
	"testing"

	// Packages
	meta "github.com/mutablelogic/go-accessory/pkg/meta"
	pg "github.com/mutablelogic/go-accessory/pkg/pg"
	assert "github.com/stretchr/testify/assert"
)

func Test_Write_001(t *testing.T) {
	type Token struct {
		Key   string `bson:"_id"`                   // Key
		Value string `bson:"value,default:'hello'"` // Value
	}

	assert := assert.New(t)
	c, err := pg.Open(context.TODO(), uri(t), tracefn(t), pg.OptApplicationName(t.Name()))
	assert.NoError(err)
	if c == nil {
		t.Skip("Skipping, connection error")
	}
	defer c.Close()

	meta := meta.New(reflect.ValueOf(Token{}), "bson")
	assert.NotNil(meta)
	err = c.(pg.Table).CreateTempTable(context.TODO(), meta)
	assert.NoError(err)
	err = c.(pg.Write).WriteInsert(context.TODO(), meta, Token{}, Token{}, Token{})
	assert.NoError(err)
}
