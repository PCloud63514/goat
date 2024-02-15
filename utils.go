package goat

func MergeSlicesUnique(a, b []string) []string {
	m := make(map[string]bool)
	var result []string

	for _, item := range a {
		if _, ok := m[item]; !ok {
			m[item] = true
			result = append(result, item)
		}
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			m[item] = true
			result = append(result, item)
		}
	}

	return result
}
