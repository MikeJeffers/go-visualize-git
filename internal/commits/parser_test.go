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
