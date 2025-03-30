package main

import (
	"fmt"
	c "gogitlog/internal/commits"
	d "gogitlog/internal/drawing"
	g "gogitlog/internal/gitlog"
	utils "gogitlog/internal/utils"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"
)

func getAllTopLevelDirs(commits []c.Commit) []string {
	commonDirsMap := map[string]bool{}
	for _, commit := range commits {
		for dir, _ := range commit.ChangesByDir {
			commonDirsMap[dir] = true
		}
	}
	return slices.Sorted(maps.Keys(commonDirsMap))
}

func bucketCommitsByTimeRange(commits []c.Commit, days int) [][]c.Commit {
	utils.Reverse(commits)
	first := commits[0].Date
	last := commits[len(commits)-1].Date
	var buckets [][]c.Commit
	for date := first; !date.After(last); date = date.AddDate(0, 0, days) {
		next := date.AddDate(0, 0, days)
		var commitsInRange []c.Commit
		for _, commit := range commits {
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
		fmt.Println(commit.ChangesByFile)
		fmt.Println(commit.ChangesByDir)
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

	commonDirs := getAllTopLevelDirs(commits)
	d.DrawCommitsByDirChanged(buckets, commonDirs, 800, 300)

	//d.DrawCommits(buckets, 800, 300)
}
