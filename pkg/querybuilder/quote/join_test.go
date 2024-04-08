package quote_test

import (
	"testing"

	// Packages
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"
	assert "github.com/stretchr/testify/assert"
)

func Test_Join_001(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		In       []any
		Expected any
	}{
		{[]any{}, ``},
		{[]any{1, 2, 3}, `1 2 3`},
		{[]any{"1", 2, "", 3}, `1 2 3`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, quote.Join(test.In...))
	}
}

func Test_Join_002(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		In       []any
		Expected any
	}{
		{[]any{}, ``},
		{[]any{1, 2, 3}, `1,2,3`},
		{[]any{"1", 2, "", 3, ""}, `1,2,3`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, quote.JoinSep(",", test.In...))
	}
}

func Test_Join_003(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		In       []any
		Expected any
	}{
		{[]any{}, ``},
		{[]any{1, 2, 3}, `123`},
		{[]any{"1", 2, "", 3, ""}, `123`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, quote.JoinSep("", test.In...))
	}
}
