package commits

import (
	"testing"
)

func TestCommitDateParse(t *testing.T) {
	commit := NewCommit()
	commit.SetDate("Date:   2023-12-25T01:11:27-05:00")
	if commit.Date.IsZero() {
		t.Errorf("Failed to parse date %s", commit.Date.String())
	}
}

func TestCommitInvalidDateParse(t *testing.T) {
	commit := NewCommit()
	commit.SetDate("Date:   2023-1BLARG2-25T01:11:ohno27-05:00")
	if !commit.Date.IsZero() {
		t.Errorf("Should have failed to parse date! %s", commit.Date.String())
	}
	if commit.IsValid() {
		t.Errorf("Should not be valid commit data")
	}
}

func TestCommitChangesParse(t *testing.T) {
	commit := NewCommit()
	commit.SetValues("2 files changed, 89 insertions(+)")
	if commit.FilesChanged != 2 {
		t.Errorf("Incorrect filesChanged value %d", commit.FilesChanged)
	} else if commit.Insertions != 89 {
		t.Errorf("Incorrect insertions value %d", commit.Insertions)
	} else if commit.Deletions != 0 {
		t.Errorf("Incorrect deletions value %d", commit.Deletions)
	}
}

func TestCommitChangesParse2(t *testing.T) {
	commit := NewCommit()
	commit.SetValues("1 file changed, 1 deletion(-)")
	if commit.FilesChanged != 1 {
		t.Errorf("Incorrect filesChanged value %d", commit.FilesChanged)
	} else if commit.Insertions != 0 {
		t.Errorf("Incorrect insertions value %d", commit.Insertions)
	} else if commit.Deletions != 1 {
		t.Errorf("Incorrect deletions value %d", commit.Deletions)
	}
}

func TestCommitChangesParse3(t *testing.T) {
	commit := NewCommit()
	commit.SetValues(" 1 file changed, 1 insertion(+), 1 deletion(-)")
	if commit.FilesChanged != 1 {
		t.Errorf("Incorrect filesChanged value %d", commit.FilesChanged)
	} else if commit.Insertions != 1 {
		t.Errorf("Incorrect insertions value %d", commit.Insertions)
	} else if commit.Deletions != 1 {
		t.Errorf("Incorrect deletions value %d", commit.Deletions)
	}
}
