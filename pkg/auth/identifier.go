package auth

import "regexp"

var (
	reIdentifier = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-\.]+$`)
)

func isIdentifier(value string) bool {
	return reIdentifier.MatchString(value)
}
