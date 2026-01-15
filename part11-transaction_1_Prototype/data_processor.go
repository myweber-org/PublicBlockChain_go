
package main

import "fmt"

func MovingAverage(data []float64, window int) []float64 {
    if window <= 0 || window > len(data) {
        return nil
    }

    result := make([]float64, len(data)-window+1)
    var sum float64

    for i := 0; i < window; i++ {
        sum += data[i]
    }
    result[0] = sum / float64(window)

    for i := window; i < len(data); i++ {
        sum = sum - data[i-window] + data[i]
        result[i-window+1] = sum / float64(window)
    }

    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    averaged := MovingAverage(sampleData, 3)
    fmt.Printf("Moving average result: %v\n", averaged)
}