package utils

import (
	c "gogitlog/internal/commits"
	"testing"
)

func TestReverse(t *testing.T) {
	var commits []c.Commit
	var commitA c.Commit
	commitA.SetDate("2025-12-01T00:00:00")
	commitA.FilesChanged = 1
	commits = append(commits, commitA)
	Reverse(commits)
	if len(commits) != 1 || !commits[0].Equals(&commitA) {
		t.Errorf("Reversing did something unexpected to single item list")
	}
	var commitB c.Commit
	commitB.SetDate("2026-12-02T00:00:00")
	commitB.FilesChanged = 2
	commits = append(commits, commitB)
	if commitA.Equals(&commitB) {
		t.Errorf("These two commits should be unequal")
	} else if len(commits) != 2 {
		t.Errorf("Array of unexpected len %d", len(commits))
	} else if !commits[0].Equals(&commitA) || !commits[1].Equals(&commitB) {
		t.Errorf("Expected items not in expected order")
	}
	Reverse(commits)
	if !commits[1].Equals(&commitA) || !commits[0].Equals(&commitB) {
		t.Errorf("Expected items not in expected order")
	}
}

func TestSumFromCommits(t *testing.T) {
	var commits []c.Commit
	for i := 0; i < 10; i++ {
		var commit c.Commit
		commit.Insertions = i
		commits = append(commits, commit)
	}
	sum := SumFromCommits(commits, func(commit c.Commit) int { return commit.Insertions })
	if sum != 45 {
		t.Errorf("Unexpected sum %d", sum)
	}
}
