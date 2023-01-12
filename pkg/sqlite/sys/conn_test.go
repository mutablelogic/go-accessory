package sqlite_test

import (
	"testing"

	// Packages
	sqlite "github.com/mutablelogic/go-accessory/pkg/sqlite/sys"
	assert "github.com/stretchr/testify/assert"
)

func Test_Conn_001(t *testing.T) {
	assert := assert.New(t)
	db, err := sqlite.OpenPath(sqlite.DefaultMemory, sqlite.SQLITE_OPEN_CREATE, "")
	assert.NoError(err)
	assert.NotNil(db)
	t.Log(db)
	assert.NoError(db.Close())
}
