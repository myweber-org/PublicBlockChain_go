
package main

import "fmt"

func FilterAndDoublePositiveInts(numbers []int) []int {
    var result []int
    for _, num := range numbers {
        if num > 0 {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{-5, 2, 0, 8, -1, 10}
    output := FilterAndDoublePositiveInts(input)
    fmt.Println("Processed slice:", output)
}