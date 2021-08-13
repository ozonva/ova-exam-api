package utils

import (
	"errors"
	"math"
	"ova-exam-api/cmd/ova-exam-api/domain/entity/user"
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
func SplitToBulks(entities []user.User, butchSize uint) [][]user.User {
	if len(entities) == 0 {
		panic("Incorrect parameters")
	}

	var floatButchSize = float64(len(entities)) / float64(butchSize)
	var size = int(math.Ceil(floatButchSize))
	result := make([][]user.User, size)
	for i := 0; i < size; i++ {
		currentChunkStart := i * int(butchSize)
		var currentChunkEnd int

		if currentChunkStart+int(butchSize) < len(entities) {
			currentChunkEnd = currentChunkStart + int(butchSize)
		} else {
			currentChunkEnd = len(entities)
		}

		result[i] = entities[currentChunkStart:currentChunkEnd]
	}

	return result
}

func UsersToMap(entities []user.User) (map[uint64]user.User, error) {
	if entities == nil || len(entities) == 0 {
		return nil, errors.New("incorrect parameters")
	}

	result := make(map[uint64]user.User, len(entities))
	for _, value := range entities {
		if _, ok := result[value.UserId]; ok {
			return nil, errors.New("key value is duplicated")
		}

		result[value.UserId] = value
	}

	return result, nil
}