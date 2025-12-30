package datautils

import "sort"

func RemoveDuplicates(input []string) []string {
	if len(input) == 0 {
		return input
	}

	sort.Strings(input)

	writeIndex := 1
	for readIndex := 1; readIndex < len(input); readIndex++ {
		if input[readIndex] != input[readIndex-1] {
			input[writeIndex] = input[readIndex]
			writeIndex++
		}
	}

	return input[:writeIndex]
}