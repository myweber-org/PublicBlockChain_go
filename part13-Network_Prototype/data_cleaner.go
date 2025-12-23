package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Printf("Original: %v\n", numbers)
	fmt.Printf("Unique: %v\n", uniqueNumbers)
}