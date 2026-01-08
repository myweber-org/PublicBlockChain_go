
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
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
    cleaned := RemoveDuplicates(data)
    fmt.Printf("Original: %v\n", data)
    fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func CleanData(data []string) []string {
	cleaned := []string{}
	for _, item := range data {
		normalized := NormalizeString(item)
		cleaned = append(cleaned, normalized)
	}
	return RemoveDuplicates(cleaned)
}

func main() {
	rawData := []string{"  Apple", "banana  ", "APPLE", "Banana", "  cherry  "}
	cleanedData := CleanData(rawData)
	fmt.Println("Cleaned data:", cleanedData)
}