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

func Test_PG_Column_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{N("a").T("text"), `a TEXT`},
		{N("a").T("text").NotNull(), `a TEXT NOT NULL`},
		{N("a").T("text").PrimaryKey(), `a TEXT PRIMARY KEY`},
		{N("a").T("uuid").Default("uuid_generate_v1mc()"), `a UUID DEFAULT uuid_generate_v1mc()`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}
