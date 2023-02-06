package types

import "errors"

func Filter[T any](arr []T, predicate func(T) bool) []T {
	res := make([]T, 0, len(arr))
	for _, val := range arr {
		if predicate(val) {
			res = append(res, val)
		}
	}
	return res
}

func Count[T any](arr []T, predicate func(T) bool) int {
	count := 0
	for _, inv := range arr {
		if predicate(inv) {
			count++
		}
	}
	return count
}

func Remove[T any](arr []T, index int) ([]T, error) {
	N := len(arr)
	if N == 0 {
		return arr, nil
	}
	if index < 0 {
		index = N + index
	}
	if index >= N {
		return arr, errors.New("index out of range")
	}
	arr[index] = arr[N-1]
	return arr[:N-1], nil
}
