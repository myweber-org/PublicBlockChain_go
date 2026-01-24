
package datautils

func DeduplicateInts(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, v := range input {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func FilterPositiveInts(input []int) []int {
	result := []int{}
	for _, v := range input {
		if v > 0 {
			result = append(result, v)
		}
	}
	return result
}