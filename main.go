package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
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
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	return string(out)
}

func splitCommits(raw string) []Commit {
	var commits []Commit
	re := regexp.MustCompile("[0-9]+")
	commitDetails := strings.Split(raw, "\ncommit ")
	for _, commit := range commitDetails {
		var details Commit
		for _, line := range strings.Split(commit, "\n") {
			if strings.Contains(line, "Date:") {
				cleaned := strings.ReplaceAll(line, "Date:", "")
				cleaned = strings.Trim(cleaned, " ")
				date, err := time.Parse(time.RFC3339, cleaned)
				if err != nil {
					log.Fatalf("Failed to parse timestamp: %s", cleaned)
				}
				details.date = date
			} else if strings.Contains(line, " insertions(+), ") {
				data := re.FindAllString(line, -1)
				if len(data) != 3 {
					log.Fatalf("Unexpected result from regex on commit line: %s", line)
				}
				if r, e := strconv.Atoi(data[0]); e == nil {
					details.filesChanged = r
				}
				if r, e := strconv.Atoi(data[1]); e == nil {
					details.insertions = r
				}
				if r, e := strconv.Atoi(data[2]); e == nil {
					details.deletions = r
				}
			}
		}
		commits = append(commits, details)
	}
	return commits
}

func bucketCommitsByTimeRange(commits []Commit, days int) [][]Commit {
	Reverse(commits)
	first := commits[0].date
	last := commits[len(commits)-1].date
	var buckets [][]Commit
	for d := first; !d.After(last); d = d.AddDate(0, 0, days) {
		next := d.AddDate(0, 0, days)
		var commitsInRange []Commit
		for _, commit := range commits {
			if commit.date.After(next) {
				break
			} else if commit.date.Before(d) {
				continue
			}
			commitsInRange = append(commitsInRange, commit)
		}
		buckets = append(buckets, commitsInRange)
	}
	return buckets
}

func printCommitLog(commits []Commit) {
	for _, commit := range commits {
		fmt.Printf("Commit: %s \t %s \t %s \t %s \n",
			commit.date.Format(time.RFC1123Z),
			strconv.Itoa(commit.filesChanged),
			strconv.Itoa(commit.insertions),
			strconv.Itoa(commit.deletions),
		)
	}
}

func parseArgs(args []string) (int, string) {
	if len(os.Args) < 2 {
		log.Fatal("insufficient positional args")
	}
	days := 7
	if len(os.Args) > 2 {
		d, err := strconv.Atoi(os.Args[2])
		if err == nil {
			days = d
		}
	}
	repoPath := os.Args[1]
	return days, repoPath
}

func main() {
	days, repoPath := parseArgs(os.Args)

	s := getGitLogRaw(repoPath)
	commits := splitCommits(s)

	printCommitLog(commits)
	fmt.Println("Total commits in", repoPath, len(commits))

	buckets := bucketCommitsByTimeRange(commits, days)

	drawCommits(buckets, 800, 300)
}
