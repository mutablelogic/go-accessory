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

func (B) CollectionName() string {
	return "b"
}

type C struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"-"`
}

type D struct {
	C
	E string `json:"e"`
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
	assert.Len(r.Fields, 2)
	assert.Equal("a", r.Fields[0].Name)
	assert.Equal("b", r.Fields[1].Name)
	t.Log(r)
}

func Test_Reflect_006(t *testing.T) {
	assert := assert.New(t)
	r := meta.New(reflect.ValueOf(D{}), "json")
	assert.NotNil(r)
	assert.Equal("D", r.Name)
	assert.Len(r.Fields, 3)
	assert.Equal("a", r.Fields[0].Name)
	assert.Equal("b", r.Fields[1].Name)
	assert.Equal("e", r.Fields[2].Name)
	t.Log(r)
}
