package main

import (
	"fmt"
	c "gogitlog/internal/commits"
	d "gogitlog/internal/drawing"
	g "gogitlog/internal/gitlog"
	utils "gogitlog/internal/utils"
	"maps"
	"os"
	"slices"
)

func getAllTopLevelDirs(commits []c.Commit) []string {
	commonDirsMap := map[string]bool{}
	for _, commit := range commits {
		for dir := range commit.ChangesByDir {
			commonDirsMap[dir] = true
		}
	}
	return slices.Sorted(maps.Keys(commonDirsMap))
}

func printCommitLog(commits []c.Commit) {
	for _, commit := range commits {
		fmt.Println(commit.String())
		fmt.Println(commit.ChangesByDir)
	}
}

func main() {
	days, repoPath, mode := utils.ParseArgs(os.Args)
	fmt.Printf("%+x %+x %+x\n", days, repoPath, mode)

	s := g.GitLog(repoPath)
	commits := c.ParseCommits(s)

	printCommitLog(commits)
	fmt.Println("Total commits in", repoPath, len(commits))

	buckets := c.BucketCommitsByTimeRange(commits, days)

	commonDirs := getAllTopLevelDirs(commits)
	switch mode {
	case utils.AdditionsDeletions:
		d.DrawCommits(buckets, 800, 300)
	case utils.ByDirs:
		d.DrawCommitsByDirChanged(buckets, commonDirs, 800, 300)
	default:
		d.DrawCommits(buckets, 800, 300)
	}
}
