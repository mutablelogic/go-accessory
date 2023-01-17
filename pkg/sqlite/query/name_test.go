package query_test

import (
	"testing"

	// Packages
	assert "github.com/stretchr/testify/assert"

	// Namespace import
	. "github.com/mutablelogic/go-accessory"
	. "github.com/mutablelogic/go-accessory/pkg/sqlite/query"
)

func Test_Name_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In     Query
		String string
	}{
		{N("a"), `a`},
		{N("a").As("b"), `a AS b`},
		{N("a").WithSchema("main"), `main.a`},
		{N("a").WithSchema("main").As("b"), `main.a AS b`},
		{N("x y").WithSchema("main").As("b"), `main."x y" AS b`},
		{N("insert").WithSchema("main").As("b"), `main."insert" AS b`},
		{N("x").WithType("TEXT"), `x TEXT`},
		{N("x", DESC), `x DESC`},
		{N("x", ASC), `x ASC`},
		{N("x", DESC, ASC), `x ASC`}, // defaults to ASC, not DESC
		{N("x", PRIMARY_KEY, AUTO_INCREMENT), `x PRIMARY KEY AUTOINCREMENT`},
	}
	for _, test := range tests {
		assert.Equal(test.String, test.In.Query())
	}
}
