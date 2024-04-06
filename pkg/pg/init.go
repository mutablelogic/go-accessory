package pg

import (
	"strings"

	// Packages
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"
)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	reservedWords = `SELECT INSERT UPDATE DELETE FROM WHERE ORDER BY GROUP BY JOIN INNER OUTER LEFT RIGHT ON AS AND OR NOT BETWEEN IN LIKE NULL TRUE FALSE IS EXISTS UNIQUE PRIMARY FOREIGN REFERENCES`
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func init() {
	quote.Init(strings.Fields(reservedWords))
}
