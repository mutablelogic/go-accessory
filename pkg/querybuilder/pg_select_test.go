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

func Test_PG_Select_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{N("a").Select(), `SELECT * FROM a`},
		{N("a").Select().Distinct(), `SELECT DISTINCT * FROM a`},
		{N("a").Select("c1", "c2"), `SELECT c1,c2 FROM a`},
		{N("a").Select("c1", "c2").Distinct(), `SELECT DISTINCT c1,c2 FROM a`},
		{N("a").WithSchema("public").Select("c1", "c2"), `SELECT c1,c2 FROM public.a`},
		{N("a").As("d1").Select(N("c1"), N("c2").WithSchema("a")), `SELECT c1,a.c2 FROM a AS d1`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}
