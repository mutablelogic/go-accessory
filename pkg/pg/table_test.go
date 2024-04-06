package pg_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	// Packages
	meta "github.com/mutablelogic/go-accessory/pkg/meta"
	pg "github.com/mutablelogic/go-accessory/pkg/pg"
	assert "github.com/stretchr/testify/assert"
)

func Test_Table_001(t *testing.T) {
	assert := assert.New(t)
	c, err := pg.Open(context.TODO(), uri(t), tracefn(t), pg.OptApplicationName(t.Name()))
	assert.NoError(err)
	if c == nil {
		t.Skip("Skipping, connection error")
	}
	defer c.Close()

	// Create a table
	type X struct {
		A string `pg:"_id,a"`
	}
	meta := meta.New(reflect.ValueOf(X{}), "pg")
	assert.NotNil(meta)
	err = c.(pg.Table).CreateTempTable(context.TODO(), meta)
	assert.NoError(err)
}

func Test_Table_002(t *testing.T) {
	type TokenType string
	type Token struct {
		Key    string    `bson:"_id"`                  // Key
		Name   string    `bson:"name"`                 // Name of token
		Type   TokenType `bson:"type,type:text"`       // Type of token
		Value  string    `bson:"value"`                // Token value
		Expire time.Time `bson:"expires_at,omitempty"` // Time of expiration for the token
		Time   time.Time `bson:"access_at"`            // Time of last access
		Scope  []string  `bson:"scope"`                // Authorization scopes
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
}
