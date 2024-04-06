package quote_test

import (
	"testing"

	// Import Namespace
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"

	// Import PG
	_ "github.com/mutablelogic/go-accessory/pkg/pg"
)

func Test_Quote_001(t *testing.T) {
	var tests = []struct{ from, to string }{
		{"", `""`},
		{"test", `"test"`},
		{"test\"", `"test"""`},
		{"\"test\"", `"""test"""`},
	}
	for i, test := range tests {
		if quote.Double(test.from) != test.to {
			t.Errorf("%d: Expected %s, got %s", i, test.to, quote.Double(test.from))
		}
	}
}

func Test_Quote_002(t *testing.T) {
	var tests = []struct{ from, to string }{
		{"", `''`},
		{"test", `'test'`},
		{"test'", `'test'''`},
		{"'test'", `'''test'''`},
	}
	for i, test := range tests {
		if quote.Single(test.from) != test.to {
			t.Errorf("%d: Expected %s, got %s", i, test.to, quote.Single(test.from))
		}
	}
}

func Test_Quote_003(t *testing.T) {
	var tests = []struct{ from, to string }{
		{"", `""`},
		{"test", `test`},
		{"order", `"order"`},
		{"some other", `"some other"`},
	}
	for i, test := range tests {
		if quote.Identifier(test.from) != test.to {
			t.Errorf("%d: Expected %s, got %s", i, test.to, quote.Identifier(test.from))
		}
	}
}

func Test_Quote_004(t *testing.T) {
	var tests = []struct {
		from []string
		to   string
	}{
		{[]string{""}, `""`},
		{[]string{"test"}, `test`},
		{[]string{"order"}, `"order"`},
		{[]string{"some other"}, `"some other"`},
		{[]string{"test", "order"}, `test,"order"`},
		{[]string{"order", "test"}, `"order",test`},
		{[]string{"some other", "select"}, `"some other","select"`},
	}
	for i, test := range tests {
		if quote.Identifiers(test.from...) != test.to {
			t.Errorf("%d: Expected %s, got %s", i, test.to, quote.Identifiers(test.from...))
		}
	}
}
