package main

func sumFromCommits(arr []Commit, getValue func(commit Commit) int) int {
	res := 0
	for _, commit := range arr {
		res += getValue(commit)
	}
	return res
}

func Reverse(arr []Commit) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
