package query_test

import (
	"testing"

	// Namespace import
	. "github.com/mutablelogic/go-accessory/pkg/query"
	"github.com/stretchr/testify/assert"
)

func Test_Name_000(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(N("test").Query(), "test")
}
