package commits

import (
	"strings"
	"testing"
)

func TestCommitParse(t *testing.T) {
	commits := ParseCommits("commit hash\nDate:   2023-12-25T01:11:27-05:00\n2 files changed, 89 insertions(+)")
	if len(commits) < 1 {
		t.Errorf("Failed to yield expected number of commits")
	} else if commits[0].Date.IsZero() {
		t.Errorf("Failed to parse date %s", commits[0].Date.String())
	} else if commits[0].FilesChanged != 2 {
		t.Errorf("Incorrect filesChanged value %d", commits[0].FilesChanged)
	} else if commits[0].Insertions != 89 {
		t.Errorf("Incorrect insertions value %d", commits[0].Insertions)
	} else if commits[0].Deletions != 0 {
		t.Errorf("Incorrect deletions value %d", commits[0].Deletions)
	}
}

func TestToString(t *testing.T) {
	commits := ParseCommits("commit hash\nDate:   2023-12-25T01:11:27-05:00\n2 files changed, 89 insertions(+)")
	str := commits[0].String()
	if !strings.Contains(str, "Dec 2023") || !strings.Contains(str, "additions: 89") {
		t.Errorf("String result unexpected: %s", str)
	}
}

func TestCommitDateParse(t *testing.T) {
	var commit Commit
	commit.SetDate("Date:   2023-12-25T01:11:27-05:00")
	if commit.Date.IsZero() {
		t.Errorf("Failed to parse date %s", commit.Date.String())
	}
}

func TestCommitInvalidDateParse(t *testing.T) {
	var commit Commit
	commit.SetDate("Date:   2023-1BLARG2-25T01:11:ohno27-05:00")
	if !commit.Date.IsZero() {
		t.Errorf("Should have failed to parse date! %s", commit.Date.String())
	}
	if commit.IsValid() {
		t.Errorf("Should not be valid commit data")
	}
}

func TestCommitChangesParse(t *testing.T) {
	var commit Commit
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
	var commit Commit
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
	var commit Commit
	commit.SetValues(" 1 file changed, 1 insertion(+), 1 deletion(-)")
	if commit.FilesChanged != 1 {
		t.Errorf("Incorrect filesChanged value %d", commit.FilesChanged)
	} else if commit.Insertions != 1 {
		t.Errorf("Incorrect insertions value %d", commit.Insertions)
	} else if commit.Deletions != 1 {
		t.Errorf("Incorrect deletions value %d", commit.Deletions)
	}
}
