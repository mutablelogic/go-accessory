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

func Test_PG_Insert_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{N("a").Insert(), `INSERT INTO a DEFAULT VALUES`},
		{N("a").As("b").Insert(), `INSERT INTO a AS b DEFAULT VALUES`},
		{N("a").WithSchema("public").Insert(), `INSERT INTO public.a DEFAULT VALUES`},
		{N("a").Insert("a", "b", "c"), `INSERT INTO a ('a', 'b', 'c') VALUES ($1, $2, $3)`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}
