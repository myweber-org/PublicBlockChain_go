
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
    data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
    cleaned := removeDuplicates(data)
    fmt.Printf("Original: %v\n", data)
    fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	words := strings.Fields(trimmed)
	seen := make(map[string]bool)
	var result []string

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		if !seen[lowerWord] {
			seen[lowerWord] = true
			result = append(result, word)
		}
	}
	return strings.Join(result, " ")
}

func main() {
	testData := "  Apple banana apple   Cherry BANANA   "
	cleaned := CleanString(testData)
	fmt.Printf("Original: '%s'\n", testData)
	fmt.Printf("Cleaned:  '%s'\n", cleaned)
}