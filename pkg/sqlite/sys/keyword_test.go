package sqlite_test

import (
	"testing"

	// Packages
	sqlite "github.com/mutablelogic/go-accessory/pkg/sqlite/sys"
	//assert "github.com/stretchr/testify/assert"
)

func Test_Keyword_001(t *testing.T) {
	for i := 0; i < sqlite.KeywordCount(); i++ {
		name := sqlite.KeywordName(i)
		t.Log("Keyword ", i, "=>", name, "=>", sqlite.KeywordCheck(name))
	}
}
