package querybuilder_test

import (
	"fmt"
	"testing"

	// Packages
	"github.com/stretchr/testify/assert"

	// Import namespaces
	. "github.com/mutablelogic/go-accessory/pkg/querybuilder"

	// Import PG
	_ "github.com/mutablelogic/go-accessory/pkg/pg"
)

func Test_PG_Name_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{N("a"), `a`},
		{N("a").As("b"), `a AS b`},
		{N("a").WithSchema("main"), `main.a`},
		{N("a").WithSchema("main").As("b"), `main.a AS b`},
		{N("x y").WithSchema("main").As("b"), `main."x y" AS b`},
		{N("insert").WithSchema("main").As("b"), `main."insert" AS b`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}