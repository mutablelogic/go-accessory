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

func Test_PG_Create_Table_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{N("a").CreateTable(), `CREATE TABLE a ()`},
		{N("a").CreateTable().IfNotExists(), `CREATE TABLE IF NOT EXISTS a ()`},
		{N("a").WithSchema("public").CreateTable().Temporary(), `CREATE TEMPORARY TABLE public.a ()`},
		{N("a").WithSchema("public").CreateTable("a", "b", "c").Temporary(), `CREATE TEMPORARY TABLE public.a (a TEXT,b TEXT,c TEXT)`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}
