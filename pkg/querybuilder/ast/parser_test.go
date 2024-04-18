package ast_test

import (
	"strings"
	"testing"

	// Packages
	assert "github.com/stretchr/testify/assert"

	// Namespace imports
	ast "github.com/mutablelogic/go-accessory/pkg/querybuilder/ast"
	tokenizer "github.com/mutablelogic/go-accessory/pkg/querybuilder/tokenizer"
)

///////////////////////////////////////////////////////////////////////////////
// Scanner Tests

func Test_Parser_001(t *testing.T) {
	assert := assert.New(t)
	parser := ast.NewParser(strings.NewReader(""), tokenizer.Pos{})
	assert.NotNil(parser)
}
