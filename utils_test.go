package main

import (
	"testing"
)

func TestReverse(t *testing.T) {
	var commits []Commit
	var c Commit
	c.setDate("2025-12-01T00:00:00")
	c.filesChanged = 1
	commits = append(commits, c)
	Reverse(commits)
	if len(commits) != 1 || commits[0] != c {
		t.Errorf("Reversing did something unexpected to single item list")
	}
	var d Commit
	d.setDate("2026-12-02T00:00:00")
	d.filesChanged = 2
	commits = append(commits, d)
	if c == d {
		t.Errorf("These two commits should be unequal")
	} else if len(commits) != 2 {
		t.Errorf("Array of unexpected len %d", len(commits))
	} else if commits[0] != c || commits[1] != d {
		t.Errorf("Expected items not in expected order")
	}
	Reverse(commits)
	if commits[1] != c || commits[0] != d {
		t.Errorf("Expected items not in expected order")
	}
}

func TestSumFromCommits(t *testing.T) {
	var commits []Commit
	for i := 0; i < 10; i++ {
		var c Commit
		c.insertions = i
		commits = append(commits, c)
	}
	sum := sumFromCommits(commits, func(commit Commit) int { return commit.insertions })
	if sum != 45 {
		t.Errorf("Unexpected sum %d", sum)
	}
}
