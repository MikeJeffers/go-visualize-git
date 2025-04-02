package main

import (
	"fmt"
	c "gogitlog/internal/commits"
	d "gogitlog/internal/drawing"
	g "gogitlog/internal/gitlog"
	utils "gogitlog/internal/utils"
	"os"
)

func main() {
	days, repoPath, mode := utils.ParseArgs(os.Args)

	s := g.GitLog(repoPath)
	commits := c.ParseCommits(s)

	c.PrintCommitLog(commits)
	fmt.Println("Total commits in", repoPath, len(commits))

	buckets := c.BucketCommitsByTimeRange(commits, days)

	switch mode {
	case utils.AdditionsDeletions:
		d.DrawCommits(buckets, 800, 300)
	case utils.ByDirs:
		// TODO: specialized to characterizing work inside of top level directories
		commonDirs := c.GetAllTopLevelDirs(commits)
		d.DrawCommitsByDirChanged(buckets, commonDirs, 800, 300)
	default:
		d.DrawCommits(buckets, 800, 300)
	}
}
