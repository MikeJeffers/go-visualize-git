package commits

import (
	"strings"
)

func ParseCommits(raw string) []Commit {
	var commits []Commit
	commitRawStrings := strings.Split(raw, "commit ")
	for _, commitRaw := range commitRawStrings {
		commitData := NewCommit()
		commitData.ProcessCommitRaw(commitRaw)
		if commitData.IsValid() {
			commits = append(commits, commitData)
		}
	}
	return commits
}
