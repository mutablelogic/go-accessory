package auth_test

import (
	"testing"
	"time"

	"github.com/mutablelogic/go-accessory/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func Test_Token_001(t *testing.T) {
	assert := assert.New(t)
	// Create a token with 16 bytes and no expiry
	token := auth.NewByte16(0)
	assert.NotNil(token)
	assert.Equal(32, len(token.Value))
	assert.True(token.IsValid())
	t.Log(token)
}

func Test_Token_002(t *testing.T) {
	assert := assert.New(t)
	// Create a token with 16 bytes and expiry
	token := auth.NewByte16(-time.Second)
	assert.NotNil(token)
	assert.False(token.IsValid())
	assert.Equal(32, len(token.Value))
	t.Log(token)
}

func Test_Token_003(t *testing.T) {
	assert := assert.New(t)
	// Create a token with scope
	token := auth.NewByte16(time.Second, auth.ScopeRead)
	assert.NotNil(token)
	assert.True(token.IsScope(auth.ScopeRead))
	assert.False(token.IsScope(auth.ScopeWrite))
	t.Log(token)
}

func Test_Token_004(t *testing.T) {
	assert := assert.New(t)
	token := auth.NewByte16(-time.Second, auth.ScopeRead)
	assert.NotNil(token)
	assert.False(token.IsScope(auth.ScopeRead))
	t.Log(token)
}
