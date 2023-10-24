package tools

import (
	"strings"
)

func ItemInList[T comparable](needle T, haystack []T) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

func StringInListEqualFold(needle string, haystack []string) bool {
	for _, item := range haystack {
		if strings.EqualFold(item, needle) {
			return true
		}
	}

	return false
}

func Map[T any, O any](input []T, mapper func(T) O) []O {
	output := make([]O, len(input))

	for i := range input {
		output[i] = mapper(input[i])
	}

	return output
}
