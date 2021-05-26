package common

// StringInSlice checks if a string is in the slice
func StringInSlice(s []string, str string) bool {
	for i := range s {
		if s[i] == str {
			return true
		}
	}

	return false
}

// StringKeyInMap checks if a string key is in the map
func StringKeyInMap(m map[string]string, str string) bool {
	if _, ok := m[str]; ok {
		return true
	}

	return false
}
