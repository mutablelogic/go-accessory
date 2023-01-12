package sqlite

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Convert boolean to integer
func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// Convert integer to boolean
func intToBool(v int) bool {
	if v == 0 {
		return false
	}
	return true
}

