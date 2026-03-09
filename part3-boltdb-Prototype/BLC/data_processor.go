
package main

import (
	"fmt"
)

// CalculateMovingAverage computes the moving average of a slice of float64 values
// over a specified window size. Returns a slice of averages.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return nil
	}

	averages := make([]float64, 0, len(data)-windowSize+1)

	for i := 0; i <= len(data)-windowSize; i++ {
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			sum += data[i+j]
		}
		averages = append(averages, sum/float64(windowSize))
	}

	return averages
}

func main() {
	// Example usage
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	result := CalculateMovingAverage(data, window)
	fmt.Printf("Moving average (window=%d): %v\n", window, result)
}