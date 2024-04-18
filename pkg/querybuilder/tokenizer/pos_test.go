package tokenizer_test

import (
	"testing"

	// Packages
	assert "github.com/stretchr/testify/assert"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory/pkg/querybuilder/tokenizer"
)

func Test_Pos_001(t *testing.T) {
	assert := assert.New(t)
	pos := Pos{}
	assert.Equal("pos<1:1>", pos.String())
}

func Test_Pos_002(t *testing.T) {
	assert := assert.New(t)
	path := "test.hcl"
	pos := Pos{Path: &path}
	assert.Equal("pos<test.hcl:1:1>", pos.String())
}
