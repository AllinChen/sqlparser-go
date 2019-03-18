package main

func ElementExistsInSlice(str string, slice []string) bool {
	for i := range slice {
		if slice[i] == str {
			return true
		}
	}

	return false
}
