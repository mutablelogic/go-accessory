package querybuilder

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func join(v ...any) string {
	return joinsep(" ", v...)
}

func joinsep(sep string, v ...any) string {
	str := ""
	for _, v := range v {
		if part := fmt.Sprint(v); part != "" {
			if str != "" {
				str += sep
			}
			str += part
		}
	}
	return str
}
