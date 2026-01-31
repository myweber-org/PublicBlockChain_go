
package main

import (
	"fmt"
	"math"
)

// FilterAndTransform processes a slice of integers, filters out values below threshold,
// and applies a transformation (square root of absolute value).
func FilterAndTransform(numbers []int, threshold int) []float64 {
	var result []float64
	for _, num := range numbers {
		if num > threshold {
			transformed := math.Sqrt(math.Abs(float64(num)))
			result = append(result, transformed)
		}
	}
	return result
}

func main() {
	input := []int{-10, 5, 3, 15, 8, -2, 25}
	threshold := 5
	output := FilterAndTransform(input, threshold)
	fmt.Printf("Processed slice: %v\n", output)
}