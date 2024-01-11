package main

import (
	"testing"
)

func TestCommitParse(t *testing.T) {
	commits := ParseCommits("commit hash\nDate:   2023-12-25T01:11:27-05:00\n2 files changed, 89 insertions(+)")
	if len(commits) < 1 {
		t.Errorf("Failed to yield expected number of commits")
	} else if commits[0].date.IsZero() {
		t.Errorf("Failed to parse date %s", commits[0].date.String())
	} else if commits[0].filesChanged != 2 {
		t.Errorf("Incorrect filesChanged value %d", commits[0].filesChanged)
	} else if commits[0].insertions != 89 {
		t.Errorf("Incorrect insertions value %d", commits[0].insertions)
	} else if commits[0].deletions != 0 {
		t.Errorf("Incorrect deletions value %d", commits[0].deletions)
	}
}

func TestCommitDateParse(t *testing.T) {
	var commit Commit
	commit.setDate("Date:   2023-12-25T01:11:27-05:00")
	if commit.date.IsZero() {
		t.Errorf("Failed to parse date %s", commit.date.String())
	}
}

func TestCommitInvalidDateParse(t *testing.T) {
	var commit Commit
	commit.setDate("Date:   2023-1BLARG2-25T01:11:ohno27-05:00")
	if !commit.date.IsZero() {
		t.Errorf("Should have failed to parse date! %s", commit.date.String())
	}
	if commit.IsValid() {
		t.Errorf("Should not be valid commit data")
	}
}

func TestCommitChangesParse(t *testing.T) {
	var commit Commit
	commit.setValues("2 files changed, 89 insertions(+)")
	if commit.filesChanged != 2 {
		t.Errorf("Incorrect filesChanged value %d", commit.filesChanged)
	} else if commit.insertions != 89 {
		t.Errorf("Incorrect insertions value %d", commit.insertions)
	} else if commit.deletions != 0 {
		t.Errorf("Incorrect deletions value %d", commit.deletions)
	}
}

func TestCommitChangesParse2(t *testing.T) {
	var commit Commit
	commit.setValues("1 file changed, 1 deletion(-)")
	if commit.filesChanged != 1 {
		t.Errorf("Incorrect filesChanged value %d", commit.filesChanged)
	} else if commit.insertions != 0 {
		t.Errorf("Incorrect insertions value %d", commit.insertions)
	} else if commit.deletions != 1 {
		t.Errorf("Incorrect deletions value %d", commit.deletions)
	}
}

func TestCommitChangesParse3(t *testing.T) {
	var commit Commit
	commit.setValues(" 1 file changed, 1 insertion(+), 1 deletion(-)")
	if commit.filesChanged != 1 {
		t.Errorf("Incorrect filesChanged value %d", commit.filesChanged)
	} else if commit.insertions != 1 {
		t.Errorf("Incorrect insertions value %d", commit.insertions)
	} else if commit.deletions != 1 {
		t.Errorf("Incorrect deletions value %d", commit.deletions)
	}
}
