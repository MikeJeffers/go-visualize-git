package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	date         time.Time
	insertions   int
	deletions    int
	filesChanged int
}

func getGitLogRaw(repoPath string) string {
	cmd := "cd " + repoPath + " && git log --stat --date iso-strict"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Sprintf("Failed to execute command: %s", cmd)
	}
	return string(out)
}

func splitCommits(raw string) []Commit {
	var commits []Commit
	commitDetails := strings.Split(raw, "\ncommit ")
	for _, commit := range commitDetails {
		var details Commit
		for _, line := range strings.Split(commit, "\n") {
			if strings.Contains(line, "Date:") {
				cleaned := strings.ReplaceAll(line, "Date:", "")
				cleaned = strings.Trim(cleaned, " ")
				date, err := time.Parse(time.RFC3339, cleaned)
				if err != nil {
					fmt.Printf("Failed to parse: %s", cleaned)
					continue
				}
				details.date = date
			} else if strings.Contains(line, " insertions(+), ") {
				stripped := strings.ReplaceAll(line, " ", "")
				stripped = strings.ReplaceAll(stripped, "fileschanged", "")
				stripped = strings.ReplaceAll(stripped, "filechanged", "")
				stripped = strings.ReplaceAll(stripped, "insertions(+)", "")
				stripped = strings.ReplaceAll(stripped, "insertion(+)", "")
				stripped = strings.ReplaceAll(stripped, "deletions(-)", "")
				stripped = strings.ReplaceAll(stripped, "deletion(-)", "")
				commaSep := strings.Split(stripped, ",")
				fChange, err := strconv.Atoi(commaSep[0])
				if err != nil {
					panic(err)
				}
				insertions, err := strconv.Atoi(commaSep[1])
				if err != nil {
					panic(err)
				}
				removals, err := strconv.Atoi(commaSep[2])
				if err != nil {
					panic(err)
				}
				details.deletions = removals
				details.insertions = insertions
				details.filesChanged = fChange
			}
		}
		commits = append(commits, details)
	}
	return commits
}

func main() {
	s := getGitLogRaw(os.Args[1])
	commits := splitCommits(s)
	fmt.Print(len(commits))
	for _, commit := range commits {
		fmt.Printf("Commit: %s \t %s \t %s \t %s \n",
			commit.date.Format(time.RFC1123Z),
			strconv.Itoa(commit.filesChanged),
			strconv.Itoa(commit.insertions),
			strconv.Itoa(commit.deletions),
		)
	}
}
