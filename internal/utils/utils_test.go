package utils

import (
	"testing"
)

func TestArgParse(t *testing.T) {
	args := []string{"cmd", "~/path/to/git/repo", "13", "1"}
	days, repoPath, mode := ParseArgs(args)
	if days != 13 {
		t.Errorf("Arg parsed days: %d", days)
	} else if repoPath != "~/path/to/git/repo" {
		t.Errorf("Arg parsed repoPath: %s", repoPath)
	} else if mode != ByDirs {
		t.Errorf("Arg parsed mode: %d", mode)
	}
}

func TestReverse(t *testing.T) {
	nums := []int{1, 2, 3}
	Reverse(nums)
	if nums[2] != 1 {
		t.Errorf("Last element should be the first")
	}
	if nums[0] != 3 {
		t.Errorf("First element should be the last")
	}
	if len(nums) != 3 {
		t.Errorf("Length changed unexpectedly %d", len(nums))
	}
}
func TestReverseEmpty(t *testing.T) {
	nums := []int{}
	if len(nums) != 0 {
		t.Errorf("Should be empty")
	}
	Reverse(nums)
	if len(nums) != 0 {
		t.Errorf("Should continue to be empty")
	}
}

func TestSumCallable(t *testing.T) {
	nums := []int{1, 2, 3}
	sum := SumFromCallable(nums, func(value int) int {
		return value
	})
	if sum != 1+2+3 {
		t.Errorf("Unexpected sum %d", sum)
	}
}

func TestSumCallableEmpty(t *testing.T) {
	nums := []int{}
	sum := SumFromCallable(nums, func(value int) int {
		return value
	})
	if sum != 0 {
		t.Errorf("Unexpected sum %d", sum)
	}
}
