package meta_test

import (
	"reflect"
	"testing"

	// Packages
	meta "github.com/mutablelogic/go-accessory/pkg/meta"
	assert "github.com/stretchr/testify/assert"
)

const (
	PG_URL = "${PG_URL}"
)

////////////////////////////////////////////////////////////////////////////////

type A struct{}
type B struct{}

func (B) Name() string {
	return "b"
}

type C struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"-"`
}

////////////////////////////////////////////////////////////////////////////////

func Test_Reflect_001(t *testing.T) {
	assert := assert.New(t)
	r := meta.New(reflect.ValueOf(A{}), "")
	assert.NotNil(r)
	assert.Equal(reflect.TypeOf(A{}), r.Type)
}

func Test_Reflect_002(t *testing.T) {
	assert := assert.New(t)
	r := meta.New(reflect.ValueOf(new(A)), "")
	assert.NotNil(r)
	assert.Equal(reflect.TypeOf(A{}), r.Type)
}

func Test_Reflect_003(t *testing.T) {
	assert := assert.New(t)
	r := meta.New(reflect.ValueOf(new(A)), "")
	assert.NotNil(r)
	assert.Equal("A", r.Name)
}

func Test_Reflect_004(t *testing.T) {
	assert := assert.New(t)
	r := meta.New(reflect.ValueOf(B{}), "")
	assert.NotNil(r)
	assert.Equal("b", r.Name)
	t.Log(r)
}

func Test_Reflect_005(t *testing.T) {
	assert := assert.New(t)
	r := meta.New(reflect.ValueOf(C{}), "json")
	assert.NotNil(r)
	assert.Equal("C", r.Name)
	assert.Len(r.Attr, 2)
	assert.Equal("a", r.Attr[0].Name)
	assert.Equal("b", r.Attr[1].Name)
	t.Log(r)
}
