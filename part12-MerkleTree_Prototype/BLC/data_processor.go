
package main

import (
	"fmt"
)

// CalculateMovingAverage computes the moving average of a slice of float64 values.
// It takes a slice of data and a window size, returning a slice of averages.
// If window size is invalid (<=0 or > len(data)), it returns nil.
func CalculateMovingAverage(data []float64, window int) []float64 {
	if window <= 0 || window > len(data) {
		return nil
	}

	result := make([]float64, len(data)-window+1)
	for i := 0; i <= len(data)-window; i++ {
		sum := 0.0
		for j := i; j < i+window; j++ {
			sum += data[j]
		}
		result[i] = sum / float64(window)
	}
	return result
}

func main() {
	// Example usage
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	averages := CalculateMovingAverage(data, window)
	fmt.Printf("Data: %v\n", data)
	fmt.Printf("Moving Average (window=%d): %v\n", window, averages)
}