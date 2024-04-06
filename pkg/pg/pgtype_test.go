package pg_test

import (
	"reflect"
	"testing"
	"time"

	// Packages
	meta "github.com/mutablelogic/go-accessory/pkg/meta"
	assert "github.com/stretchr/testify/assert"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory/pkg/pg"
)

type R struct {
	A string    `bson:"_id"`
	B string    `bson:"b"`
	C int       `bson:"c,omitempty,unique"`
	D time.Time `bson:"d,omitempty"`
	E int       `bson:"e,omitempty,type:numeric"`
	F bool      `bson:"f"`
	G *string   `bson:"g"`
}

func Test_Type_001(t *testing.T) {
	assert := assert.New(t)
	meta := meta.New(reflect.ValueOf(R{}), "bson")
	assert.NotNil(meta)
	for _, field := range meta.Fields {
		col, err := PGColumn(field)
		assert.NoError(err)
		t.Log(field, "=>", col)
	}
}

func Test_Type_002(t *testing.T) {
	assert := assert.New(t)
	meta := meta.New(reflect.ValueOf(R{}), "bson")
	assert.NotNil(meta)
	fields, err := PGColumns(meta)
	assert.NoError(err)
	for _, field := range fields {
		t.Log(field)
	}
}
