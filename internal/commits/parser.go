package commits

import (
	"regexp"
	"strings"
)

var (
	numberRe    = regexp.MustCompile("[0-9]+")
	changesRe   = regexp.MustCompile("(file changed|files changed)")
	dateCleanRe = regexp.MustCompile(`(Date:|\s)`)
)

func ParseCommits(raw string) []Commit {
	var commits []Commit
	commitRawStrings := strings.Split(raw, "commit ")
	for _, commitRaw := range commitRawStrings {
		var commitData Commit
		commitData.ProcessCommitRaw(commitRaw)
		if commitData.IsValid() {
			commits = append(commits, commitData)
		}
	}
	return commits
}
