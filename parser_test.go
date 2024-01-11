package main

import (
	"testing"
)

func TestCommitParse(t *testing.T) {
	commits := ParseCommits("commit hash\nDate:   2023-12-25T01:11:27-05:00\n2 files changed, 89 insertions(+)")
	if len(commits) < 1 {
		t.Fatalf("Failed to yield expected number of commits")
	} else if commits[0].date.IsZero() {
		t.Fatalf("Failed to parse date %s", commits[0].date.String())
	} else if commits[0].filesChanged != 2 {
		t.Fatalf("Incorrect filesChanged value %d", commits[0].filesChanged)
	} else if commits[0].insertions != 89 {
		t.Fatalf("Incorrect insertions value %d", commits[0].insertions)
	} else if commits[0].deletions != 0 {
		t.Fatalf("Incorrect deletions value %d", commits[0].deletions)
	}
}

func TestCommitDateParse(t *testing.T) {
	var commit Commit
	commit.setDate("Date:   2023-12-25T01:11:27-05:00")
	if commit.date.IsZero() {
		t.Fatalf("Failed to parse date %s", commit.date.String())
	}
}

func TestCommitChangesParse(t *testing.T) {
	var commit Commit
	commit.setValues("2 files changed, 89 insertions(+)")
	if commit.filesChanged != 2 {
		t.Fatalf("Incorrect filesChanged value %d", commit.filesChanged)
	} else if commit.insertions != 89 {
		t.Fatalf("Incorrect insertions value %d", commit.insertions)
	} else if commit.deletions != 0 {
		t.Fatalf("Incorrect deletions value %d", commit.deletions)
	}
}
