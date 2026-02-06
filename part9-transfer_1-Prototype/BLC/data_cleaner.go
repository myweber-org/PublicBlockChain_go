package datautils

import "sort"

// Deduplicate removes duplicate values from a slice of comparable types
func Deduplicate[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// DeduplicateSorted removes duplicates from a sorted slice more efficiently
func DeduplicateSorted[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	result := make([]T, 0, len(input))
	result = append(result, input[0])

	for i := 1; i < len(input); i++ {
		if input[i] != input[i-1] {
			result = append(result, input[i])
		}
	}

	return result
}

// DeduplicateWithSort sorts and removes duplicates from a slice
func DeduplicateWithSort[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	// Create a copy to avoid modifying original
	sorted := make([]T, len(input))
	copy(sorted, input)

	// Sort the slice
	sort.Slice(sorted, func(i, j int) bool {
		// Use type assertion for comparison
		// This is a simplified comparison - in practice you'd need
		// a more sophisticated approach for generic types
		return false // Placeholder - actual implementation depends on type
	})

	return DeduplicateSorted(sorted)
}