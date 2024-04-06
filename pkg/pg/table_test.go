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
		A int `pg:"a,omitempty"`
	}
	meta := meta.New(reflect.ValueOf(X{}), "pg")
	assert.NotNil(meta)
	err = c.(pg.Table).CreateTempTable(context.TODO(), meta)
	assert.NoError(err)
}
