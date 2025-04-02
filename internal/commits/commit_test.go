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

func TestEquals(t *testing.T) {
	a := NewCommit()
	a.SetDate("Date:   2023-12-25T01:11:27-05:00")

	b := NewCommit()
	b.SetDate("Date:   2023-12-25T01:11:27-05:00")

	if !b.Equals(&b) {
		t.Errorf("Should be equal")
	}
	if !a.Equals(&b) {
		t.Errorf("Should be equal")
	}
	b.SetValues("2 files changed, 89 insertions(+)")
	if a.Equals(&b) {
		t.Errorf("Should not be equal")
	}
}

func TestCommitParseFileChanges(t *testing.T) {
	commit := NewCommit()
	commit.ProcessCommitRaw(`
	commit abc1234 (HEAD -> master, origin/master)
	Author: commiter <author>
	Date:   2025-03-31T23:41:54-04:00

    message

 	internal/commits/buckets_test.go | 4 ----
 	internal/drawing/draw.go         | 2 --
 	2 files changed, 6 deletions(-)
 `)
	v, ok := commit.ChangesByDir["internal"]
	if !ok {
		t.Errorf("should have changes under dir name %s", "internal")
	} else if v != 6 {
		t.Errorf("Total changes in %s should be 6, got %d", "internal", v)
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

func TestPrint(t *testing.T) {
	commit := NewCommit()
	commit.SetDate("Date:   2023-12-25T01:11:27-05:00")
	PrintCommitLog([]Commit{commit})
}

func TestGetAllTopLevelDirsEmpty(t *testing.T) {
	dirs := GetAllTopLevelDirs([]Commit{})
	if len(dirs) > 0 {
		t.Errorf("Unexpected list content %d", len(dirs))
	}
}

func TestGetAllTopLevelDirs(t *testing.T) {
	commit := NewCommit()
	commit.ProcessCommitRaw(`
	commit abc1234 (HEAD -> master, origin/master)
	Author: commiter <author>
	Date:   2025-03-31T23:41:54-04:00

    message

 	internal/commits/buckets_test.go | 4 ----
 	internal/drawing/draw.go         | 2 --
 	2 files changed, 6 deletions(-)
 `)
	dirs := GetAllTopLevelDirs([]Commit{commit})
	if len(dirs) != 1 {
		t.Errorf("Unexpected len of dirs %d", len(dirs))
	}
}
