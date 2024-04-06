package quote

import "strings"

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// IsReservedWord returns true if the given value is a reserved word
func IsReservedWord(k string) bool {
	k = strings.TrimSpace(strings.ToUpper(k))
	_, ok := reservedWords[k]
	return ok
}
