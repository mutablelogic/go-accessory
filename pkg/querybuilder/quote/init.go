package quote

import (
	"strings"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	once          sync.Once
	reservedWords map[string]bool
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Init must be called once to initialise the quote package
func Init(words []string) {
	once.Do(func() {
		reservedWords = make(map[string]bool, len(words))
		for _, k := range words {
			v := strings.TrimSpace(strings.ToUpper(k))
			reservedWords[v] = true
		}
	})
}
