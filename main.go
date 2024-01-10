package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
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
