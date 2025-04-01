package gitlog

import (
	"log"
	"os/exec"
)

// Calls "git log" on the target repository
// Output is expected in the following format:
//
// commit <COMMIT-SHA> (<BRANCH>)
// Author: <AUTHOR>
// Date:   YYYY-MM-DDTHH:MM:SS-TZ:OFFSET
// EMPTY LINE
//
//	<COMMIT MESSAGE>
//
// EMPTY LINE
//
//	<FILE NAME CHANGED>                | <LINES CHANGED> <ASCII +/- BAR CHART>
//	<NUMBER> files changed, <NUMBER> insertions(+), <NUMBER> deletions(-)
func GitLog(repoPath string) string {
	cmd := "cd " + repoPath + " && git log --stat=10000 --date iso-strict"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	return string(out)
}
