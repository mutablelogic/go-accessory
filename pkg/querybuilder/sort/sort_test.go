package sort_test

import (
	"fmt"
	"testing"

	// Packages
	sort "github.com/mutablelogic/go-accessory/pkg/querybuilder/sort"
	assert "github.com/stretchr/testify/assert"
)

func Test_Sort_001(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		In       any
		Expected string
	}{
		{sort.Sort(), ``},
		{sort.Sort("a"), `ORDER BY a`},
		{sort.Sort("a", "b"), `ORDER BY a,b`},
		{sort.Sort().Desc("a", "b"), `ORDER BY a DESC,b DESC`},
		{sort.Sort().Asc("a").Desc("b"), `ORDER BY a,b DESC`},
		{sort.Sort().Asc("a").Desc("b", "c"), `ORDER BY a,b DESC,c DESC`},
		{sort.Sort().Limit(1), `LIMIT 1`},
		{sort.Sort().Offset(1), `OFFSET 1`},
		{sort.Sort().Limit(0).Offset(0), ``},
		{sort.Sort().Limit(1).Offset(1), `LIMIT 1 OFFSET 1`},
		{sort.Sort("a").Limit(1).Offset(1), `ORDER BY a LIMIT 1 OFFSET 1`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}
