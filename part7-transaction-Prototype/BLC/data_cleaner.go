package datautils

func RemoveDuplicates[T comparable](slice []T) []T {
    seen := make(map[T]bool)
    result := []T{}
    
    for _, item := range slice {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    
    return result
}package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}
package main

import "fmt"

func removeDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
	cleaned := removeDuplicates(input)
	fmt.Println("Original:", input)
	fmt.Println("Cleaned:", cleaned)
}