package utils

import (
	"math"
)

func Div(input []int, chunkSize int) [][]int {
	if len(input) == 0 || chunkSize < 1 {
		panic("Incorrect parameters")
	}

	var floatChunkSize = float64(len(input)) / float64(chunkSize)
	var size = int(math.Ceil(floatChunkSize))
	result := make([][]int, size)
	for i := 0; i < size; i++ {
		currentChunkStart := i * chunkSize
		var currentChunkEnd int

		if currentChunkStart+chunkSize < len(input) {
			currentChunkEnd = currentChunkStart + chunkSize
		} else {
			currentChunkEnd = len(input)
		}

		result[i] = input[currentChunkStart:currentChunkEnd]
	}

	return result
}

func Invert(input map[string]string) map[string]string {
	if input == nil || len(input) == 0 {
		panic("Incorrect parameters")
	}

	result := make(map[string]string, len(input))
	for key, value := range input {
		if _, ok := result[value]; ok {
			panic("Key value is duplicated ")
		}
		result[value] = key
	}

	return result
}

func Filter(input []string) []string {
	if len(input) == 0 {
		panic("Incorrect parameters")
	}

	filterElements := map[string]bool{"слово1": true, "слово4": true, "слово5": true}

	result := make([]string, 0)
	for _, inputElement := range input {
		if _, ok := filterElements[inputElement]; ok {
			continue
		}
		result = append(result, inputElement)
	}

	return result
}
