
package main

import "fmt"

func FilterAndDoublePositiveInts(nums []int) []int {
    var result []int
    for _, num := range nums {
        if num > 0 {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{-5, 2, 0, 8, -1, 3}
    output := FilterAndDoublePositiveInts(input)
    fmt.Printf("Input: %v\n", input)
    fmt.Printf("Output: %v\n", output)
}