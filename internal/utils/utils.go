package utils

import (
	"log"
	"strconv"
)

type StackedBarMode int

const (
	AdditionsDeletions StackedBarMode = iota
	ByDirs
	end
)

func (b StackedBarMode) String() string {
	return [...]string{"AdditionsDeletions", "ByDirs"}[b]
}

func IsValidMode(value int) bool {
	return value < int(end)
}

func ParseArgs(args []string) (int, string, StackedBarMode) {
	if len(args) < 2 {
		log.Fatal("insufficient positional args")
	}
	repoPath := args[1]
	days := 7
	if len(args) > 2 {
		parsedDays, err := strconv.Atoi(args[2])
		if err == nil {
			days = parsedDays
		}
	}
	mode := AdditionsDeletions
	if len(args) > 3 {
		parsedVal, err := strconv.Atoi(args[3])
		if err == nil && IsValidMode(parsedVal) {
			mode = StackedBarMode(parsedVal)
		}
	}
	return days, repoPath, mode
}

func SumFromCallable[T any](arr []T, getValue func(arg T) int) int {
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
