package drawing

import (
	"gogitlog/internal/commits"
	"testing"
)

func TestComputeBarChart(t *testing.T) {
	empty := ComputeBarChart([][]commits.Commit{}, 100, 200, 10)
	if len(empty) != 0 {
		t.Errorf("Unexpected length %d", len(empty))
	}
	commit := commits.NewCommit()
	commit.Deletions = 2
	commit.Insertions = 4
	bucket := []commits.Commit{commit}
	nonEmpty := ComputeBarChart([][]commits.Commit{bucket}, 100, 200, 10)
	if len(nonEmpty) != 1 {
		t.Errorf("Unexpected length %d", len(nonEmpty))
	}
}
