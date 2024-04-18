package tokenizer_test

import (
	"fmt"
	"strings"
	"testing"

	// Packages
	assert "github.com/stretchr/testify/assert"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory/pkg/querybuilder/tokenizer"
)

///////////////////////////////////////////////////////////////////////////////
// Scanner Tests

func Test_Scanner_001(t *testing.T) {
	// Non-error cases
	tests := []struct {
		in string
	}{
		{""},
		{"      "},
		{"0 1 2 3 4 5 6 7 8 9"},
		{"func(test)"},
		{"   var.test   "},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			assert := assert.New(t)
			scanner := NewScanner(strings.NewReader(test.in), Pos{}, None)
			assert.NotNil(scanner)
			for tok := scanner.Next(); tok.Kind != EOF; tok = scanner.Next() {
				assert.NotNil(tok)
			}
		})
	}
}

func Test_Scanner_002(t *testing.T) {
	// Non-error cases - general
	tests := []struct {
		in     string
		out    string
		values []string
	}{
		{"", "", nil},
		{"    ", "Space", []string{"    "}},
		{" \n\t\t   ", "Space", []string{" \n\t\t   "}},
		{" ; ", "Space SemiColon Space", []string{" ", ";", " "}},
		{"0 1", "NumberInteger Space NumberInteger", []string{"0", " ", "1"}},
		{"func(test)", "Ident OpenParen Ident CloseParen", []string{"func", "(", "test", ")"}},
		{"   var.test   ", "Space Ident Punkt Ident Space", []string{"   ", "var", ".", "test", "   "}},
		{`'test'`, "String", []string{"test"}},
		{`'te"st'`, "String", []string{"te\"st"}},
		{`"e'st"`, "String", []string{"e'st"}},
		{`"te\"st"`, "String", []string{"te\"st"}},
		{`"e\""`, "String", []string{"e\""}},
		{`"t\"'t"`, "String", []string{"t\"'t"}},
		{`'e\'st'`, "String", []string{"e'st"}},
		{`'e"st'`, "String", []string{"e\"st"}},
		{`!!!`, "Not Not Not", []string{"!", "!", "!"}},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			assert := assert.New(t)
			tokens, err := NewScanner(strings.NewReader(test.in), Pos{}, None).Tokens()
			assert.NoError(err)
			assert.NotNil(tokens)
			assert.Equal(test.out, toTokenString(tokens))
			assert.Equal(test.values, toTokenValues(tokens))
		})
	}
}

func Test_Scanner_003(t *testing.T) {
	// Non-error cases - numbers and prefixes
	tests := []struct {
		in  string
		out string
	}{
		{"1267650600228229401496703205376", "NumberInteger"},
		{"-1267650600228229401496703205376", "NumberInteger"},
		{"+1267650600228229401496703205376", "NumberInteger"},
		{"0.00000000000000000000000000000078886090522101180541", "NumberFloat"},
		{".00000000000000000000000000000078886090522101180541", "NumberFloat"},
		{"-.00000000000000000000000000000078886090522101180541", "NumberFloat"},
		{"+.00000000000000000000000000000078886090522101180541", "NumberFloat"},
		{"-0.00000000000000000000000000000078886090522101180541", "NumberFloat"},
		{"+0.00000000000000000000000000000078886090522101180541", "NumberFloat"},
		{"-0.0000000000000000000000000000007E8886090522101180541", "NumberFloat"},
		{"+0.0000000000000000000000000000007e8886090522101180541", "NumberFloat"},
		{"0x12AB", "NumberHex"},
		{"0X12AB", "NumberHex"},
		{"-0x12CD", "NumberHex"},
		{"-0XED", "NumberHex"},
		{"+0xEF12", "NumberHex"},
		{"+0xFE12", "NumberHex"},
		{"012", "NumberOctal"},
		{"-0123445", "NumberOctal"},
		{"+012345", "NumberOctal"},
		{"0b1010", "NumberBinary"},
		{"-0b101011", "NumberBinary"},
		{"+0b1001001", "NumberBinary"},
		{".", "Punkt"},
		{"+", "Plus"},
		{"-", "Minus"},
		{"++", "Plus Plus"},
		{"--", "Minus Minus"},
		{"e", "Ident"},
		{"ee", "Ident"},
		{"e++", "Ident Plus Plus"},
		{"e+-", "Ident Plus Minus"},
		{".+", "Punkt Plus"},
		{".+-", "Punkt Plus Minus"},
		{"-.", "Minus Punkt"},
		{"-..", "Minus Punkt Punkt"},
		{"+.", "Plus Punkt"},
		{"+..", "Plus Punkt Punkt"},
		{"..", "Punkt Punkt"},
		{"..-+", "Punkt Punkt Minus Plus"},
		{"-.e", "Minus Punkt Ident"},
		{"+.e", "Plus Punkt Ident"},
		{"-e", "Minus Ident"},
		{"+e", "Plus Ident"},
		{"+0", "NumberInteger"},
		{"-0", "NumberInteger"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			assert := assert.New(t)
			fmt.Println("Test", t.Name())
			tokens, err := NewScanner(strings.NewReader(test.in), Pos{}, None).Tokens()
			assert.NoError(err)
			assert.NotNil(tokens)
			assert.Equal(test.out, toTokenString(tokens))
			t.Logf("%s => %s", test.in, toTokenString(tokens))
		})
	}
}

func Test_Scanner_004(t *testing.T) {
	// Error cases - numbers
	tests := []struct {
		in string
	}{
		{"12ee"},
		{"-.e16"},
		{"+08888"},
		{"0b2344"},
		{"0x100g100"},
		{"0.0000000e.45"},
		{"0.0000000e-.45"},
		{"0.0000000e+.45"},
		{"-45e++6"},
		{"-45e-.6"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			assert := assert.New(t)
			tokens, err := NewScanner(strings.NewReader(test.in), Pos{}, None).Tokens()
			if err == nil {
				t.Logf("%s => %s", test.in, toTokenString(tokens))
			}
			assert.Error(err)
		})
	}
}

func Test_Scanner_005(t *testing.T) {
	// Non-error cases - comments
	tests := []struct {
		in  string
		out string
	}{
		{"# Hash comment #", "Comment"},
		{"# Hash comment\n # Another comment", "Comment Space Comment"},
		{"// Line comment", "Comment"},
		{"//// Line comment", "Comment"},
		{"/ / // Line comment", "Slash Space Slash Space Comment"},
		{"// Line comment\n// Line comment", "Comment Space Comment"},
		{"test /* Block comment */ test", "Ident Space Comment Space Ident"},
		{"test /* Block \n comment */ test", "Ident Space Comment Space Ident"},
		{"test /*** Block \n comment ***/ test", "Ident Space Comment Space Ident"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			assert := assert.New(t)
			tokens, err := NewScanner(strings.NewReader(test.in), Pos{}, HashComments|LineComments|BlockComments).Tokens()
			assert.NoError(err)
			assert.NotNil(tokens)
			assert.Equal(test.out, toTokenString(tokens))
		})
	}
}

func Test_Scanner_006(t *testing.T) {
	// Non-error cases - unary operators
	tests := []struct {
		in  string
		out string
	}{
		{"!test", "Not Ident"},
		{"!(test)", "Not OpenParen Ident CloseParen"},
		{"!55", "Not NumberInteger"},
		{"(!0b101)", "OpenParen Not NumberBinary CloseParen"},
		{"(!0x101)", "OpenParen Not NumberHex CloseParen"},
		{"(!0101)", "OpenParen Not NumberOctal CloseParen"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			assert := assert.New(t)
			tokens, err := NewScanner(strings.NewReader(test.in), Pos{}, None).Tokens()
			assert.NoError(err)
			assert.NotNil(tokens)
			assert.Equal(test.out, toTokenString(tokens))
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// Private Methods

func toTokenString(tokens []*Token) string {
	var result []string
	for _, token := range tokens {
		result = append(result, token.Kind.String())
	}
	return strings.Join(result, " ")
}

func toTokenValues(tokens []*Token) []string {
	var result []string
	for _, token := range tokens {
		result = append(result, token.Val)
	}
	return result
}
