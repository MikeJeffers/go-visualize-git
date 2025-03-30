package main

import (
	"fmt"
	c "gogitlog/internal/commits"
	d "gogitlog/internal/drawing"
	g "gogitlog/internal/gitlog"
	utils "gogitlog/internal/utils"
	"log"
	"os"
	"strconv"
)

func bucketCommitsByTimeRange(commitArray []c.Commit, days int) [][]c.Commit {
	utils.Reverse(commitArray)
	first := commitArray[0].Date
	last := commitArray[len(commitArray)-1].Date
	var buckets [][]c.Commit
	for date := first; !date.After(last); date = date.AddDate(0, 0, days) {
		next := date.AddDate(0, 0, days)
		var commitsInRange []c.Commit
		for _, commit := range commitArray {
			if commit.Date.After(next) {
				break
			} else if commit.Date.Before(date) {
				continue
			}
			commitsInRange = append(commitsInRange, commit)
		}
		buckets = append(buckets, commitsInRange)
	}
	return buckets
}

func printCommitLog(commits []c.Commit) {
	for _, commit := range commits {
		fmt.Println(commit.String())
	}
}

func parseArgs(args []string) (int, string) {
	if len(args) < 2 {
		log.Fatal("insufficient positional args")
	}
	days := 7
	if len(args) > 2 {
		parsedDays, err := strconv.Atoi(args[2])
		if err == nil {
			days = parsedDays
		}
	}
	repoPath := args[1]
	return days, repoPath
}

func main() {
	days, repoPath := parseArgs(os.Args)

	s := g.GitLog(repoPath)
	commits := c.ParseCommits(s)

	printCommitLog(commits)
	fmt.Println("Total commits in", repoPath, len(commits))

	buckets := bucketCommitsByTimeRange(commits, days)

	d.DrawCommits(buckets, 800, 300)
}
