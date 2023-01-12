package sqlite_test

import (
	"testing"

	// Packages
	"github.com/mutablelogic/go-accessory/pkg/sqlite/sys"
	"github.com/stretchr/testify/assert"
)

func Test_Sqlite_001(t *testing.T) {
	assert := assert.New(t)

	t.Run("001", func(t *testing.T) {
		version, number, source := sqlite.Version()
		assert.NotEmpty(version)
		assert.NotZero(number)
		assert.NotEmpty(source)
		t.Log(version, number, source)
	})
}
