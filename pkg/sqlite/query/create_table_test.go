package query_test

import (
	"testing"

	// Packages
	assert "github.com/stretchr/testify/assert"

	// Namespace import
	. "github.com/mutablelogic/go-accessory"
	. "github.com/mutablelogic/go-accessory/pkg/sqlite/query"
)

func Test_CreateTable_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In     Query
		String string
	}{
		{N("a").CreateTable(), `CREATE TABLE a`},
		{N("a").WithSchema("b").CreateTable(), `CREATE TABLE b.a`},
	}
	for _, test := range tests {
		assert.Equal(test.String, test.In.Query())
	}
}
