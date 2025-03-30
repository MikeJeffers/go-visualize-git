package utils

import (
	c "gogitlog/internal/commits"
)

func SumFromCommits(arr []c.Commit, getValue func(commit c.Commit) int) int {
	res := 0
	for _, commit := range arr {
		res += getValue(commit)
	}
	return res
}

func Reverse[T any](arr []T) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
