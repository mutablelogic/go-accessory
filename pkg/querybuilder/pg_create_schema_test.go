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

func Test_PG_Create_Schema_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{N("a").CreateSchema(), `CREATE SCHEMA a`},
		{N("select").CreateSchema().IfNotExists(), `CREATE SCHEMA IF NOT EXISTS "select"`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}
