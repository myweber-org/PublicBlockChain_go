
package main

import (
	"fmt"
)

// CalculateMovingAverage computes the moving average of a slice of float64 values.
// It returns a new slice where each element is the average of the preceding 'windowSize' elements.
// If the window size is larger than the number of elements up to the current index,
// it averages the available elements.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 {
		return []float64{}
	}

	result := make([]float64, len(data))
	for i := range data {
		start := i - windowSize + 1
		if start < 0 {
			start = 0
		}
		sum := 0.0
		count := 0
		for j := start; j <= i; j++ {
			sum += data[j]
			count++
		}
		result[i] = sum / float64(count)
	}
	return result
}

func main() {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	averages := CalculateMovingAverage(data, window)
	fmt.Printf("Data: %v\n", data)
	fmt.Printf("Moving Averages (window=%d): %v\n", window, averages)
}