package quote

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Join returns a string that is the concatenation of the arguments, separated by a spaces
// If any argument is an empty string, it is ignored
func Join(v ...any) string {
	return JoinSep(" ", v...)
}

// JoinSep returns a string that is the concatenation of the arguments, separated by
// a separator. If any argument is nil or evaluates to an empty string, it is ignored
func JoinSep(sep string, v ...any) string {
	str := ""
	for _, v := range v {
		if v == nil {
			continue
		}
		part := fmt.Sprint(v)
		if part != "" {
			if str != "" {
				str += sep
			}
			str += part
		}
	}
	if str == "" {
		return ""
	} else if len(sep) == 0 {
		return str
	} else {
		return str[len(sep)-1:]
	}
}
