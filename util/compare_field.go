package util

func CompareFields(fieldMap map[string][2]any) []string {
	var result []string

	for msg, pair := range fieldMap {
		if pair[0] == pair[1] {
			result = append(result, msg)
		}
	}

	if len(result) == 0 {
		for msg := range fieldMap {
			result = append(result, msg)
		}
	}

	return result
}
