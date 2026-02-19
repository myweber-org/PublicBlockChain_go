
package main

import (
    "fmt"
)

// FilterAndDouble filters out even numbers from the input slice,
// doubles the remaining odd numbers, and returns the new slice.
func FilterAndDouble(numbers []int) []int {
    var result []int
    for _, num := range numbers {
        if num%2 != 0 {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := FilterAndDouble(input)
    fmt.Println("Original:", input)
    fmt.Println("Processed:", output)
}