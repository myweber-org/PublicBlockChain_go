
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
	input := []int{1, 2, 2, 3, 4, 4, 5}
	output := RemoveDuplicates(input)
	fmt.Println(output)
}
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
    if len(nums) == 0 {
        return nums
    }

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
    cleaned := RemoveDuplicates(input)
    fmt.Printf("Original: %v\n", input)
    fmt.Printf("Cleaned: %v\n", cleaned)
}