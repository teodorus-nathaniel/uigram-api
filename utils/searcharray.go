package utils

func SearchArray(array []string, item string) bool {
	for _, el := range array {
		if el == item {
			return true
		}
	}

	return false
}
