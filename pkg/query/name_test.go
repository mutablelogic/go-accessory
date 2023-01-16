package query_test

import (
	"testing"

	// Packages
	assert "github.com/stretchr/testify/assert"

	// Namespace import
	. "github.com/mutablelogic/go-accessory/pkg/query"
	. "github.com/mutablelogic/go-accessory"
)

func Test_Name_000(t *testing.T) {
	tests := []struct {
		In     Expr
		String string
	}{
		{N("a"), `a`},
		{N("a").WithAlias("b"), `a AS b`},
		{N("a").WithSchema("main"), `main.a`},
		{N("a").WithSchema("main").WithAlias("b"), `main.a AS b`},
		{N("x y").WithSchema("main").WithAlias("b"), `main."x y" AS b`},
		{N("insert").WithSchema("main").WithAlias("b"), `main."insert" AS b`},
		{N("x").WithType("TEXT"), `x TEXT`},
		{N("x").WithDesc(), `x DESC`},
	}

	for _, test := range tests {
		if v := test.In.String(); v != test.String {
			t.Errorf("got %v, wanted %v", v, test.String)
		}
	}
}
