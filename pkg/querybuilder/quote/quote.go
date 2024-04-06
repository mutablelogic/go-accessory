package quote

import (
	"regexp"
	"strings"
)

/////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	regexpBareIdentifier = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
)

/////////////////////////////////////////////////////////////////////
// FUNCTIONS

// Single puts single quotes around a string and escapes existing single quotes
func Single(value string) string {
	return "'" + strings.Replace(value, "'", "''", -1) + "'"
}

// Double puts double quotes around a string and escapes existing double quotes
func Double(value string) string {
	return "\"" + strings.Replace(value, "\"", "\"\"", -1) + "\""
}

// Identifier returns a safe version of an identifier
func Identifier(v string) string {
	if IsReservedWord(v) {
		return Double(v)
	} else if isBareIdentifier(v) {
		return v
	} else {
		return Double(v)
	}
}

// Identifiers returns a safe version of a list of identifiers,
// separated by commas
func Identifiers(v ...string) string {
	if len(v) == 0 {
		return ""
	}
	if len(v) == 1 {
		return Identifier(v[0])
	}
	result := make([]string, len(v))
	for i, v_ := range v {
		result[i] = Identifier(v_)
	}
	return strings.Join(result, ",")
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func isBareIdentifier(value string) bool {
	return regexpBareIdentifier.MatchString(value)
}
