package quote

import (
	"strings"
	"sync"

	// Packages
	sqlite "github.com/mutablelogic/go-accessory/pkg/sqlite/sys"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	t = `TEXT BLOB FLOAT INTEGER BOOL`
)

var (
	reservedOnce  sync.Once
	reservedWords map[string]bool
	reservedTypes map[string]bool
)

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS

// quoteInit initializes the reserved words and types maps
func quoteInit() {
	reservedOnce.Do(func() {
		reservedWords = make(map[string]bool, sqlite.KeywordCount())
		reservedTypes = make(map[string]bool, 5)
		for i := 0; i < sqlite.KeywordCount(); i++ {
			k := strings.ToUpper(sqlite.KeywordName(i))
			reservedWords[k] = true
		}
		for _, k := range strings.Fields(t) {
			reservedTypes[k] = true
		}
	})
}

// IsReservedWord returns true if the given string is a reserved word
func IsReservedWord(k string) bool {
	quoteInit()
	k = strings.ToUpper(k)
	_, ok := reservedWords[k]
	return ok
}

// IsType returns true if the given string is a sqlite type
func IsType(k string) bool {
	quoteInit()
	_, ok := reservedTypes[k]
	return ok
}

// ReservedWords returns a list of reserved words
func ReservedWords() []string {
	quoteInit()
	result := make([]string, 0, len(reservedWords))
	for k, _ := range reservedWords {
		result = append(result, k)
	}
	return result
}

// Types returns a list of sqlite types
func Types() []string {
	return strings.Fields(t)
}
