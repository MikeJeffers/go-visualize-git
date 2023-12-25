package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
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
		fmt.Printf("Failed to execute command: %s", cmd)
		panic(err)
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
					panic(err)
				}
				details.date = date
			} else if strings.Contains(line, " insertions(+), ") {
				stripped := strings.ReplaceAll(line, " ", "")
				stripped = strings.ReplaceAll(stripped, "s", "") //Silly plural hack
				stripped = strings.ReplaceAll(stripped, "filechanged", "")
				stripped = strings.ReplaceAll(stripped, "inertion(+)", "") // Silly consequence of plural hack
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

func bucketCommitsByMonth(commits []Commit) [][]Commit {
	first := commits[0].date
	last := commits[len(commits)-1].date
	var buckets [][]Commit
	for d := first; !d.After(last); d = d.AddDate(0, 0, 28) {
		next := d.AddDate(0, 0, 28)
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

func sumFromCommits(arr []Commit, getValue func(commit Commit) int) int {
	res := 0
	for _, commit := range arr {
		res += getValue(commit)
	}
	return res
}

func main() {
	s := getGitLogRaw(os.Args[1])
	commits := splitCommits(s)
	total := len(commits)
	fmt.Print(total)
	for _, commit := range commits {
		fmt.Printf("Commit: %s \t %s \t %s \t %s \n",
			commit.date.Format(time.RFC1123Z),
			strconv.Itoa(commit.filesChanged),
			strconv.Itoa(commit.insertions),
			strconv.Itoa(commit.deletions),
		)
	}
	// Reverse
	for i, j := 0, len(commits)-1; i < j; i, j = i+1, j-1 {
		commits[i], commits[j] = commits[j], commits[i]
	}

	width := 800
	height := 300

	//stepW := float64(width) / float64(total)

	dc := gg.NewContext(width, height)
	// dc.SetRGB(255, 255, 255)
	// dc.DrawRectangle(0, 0, float64(width), float64(height))
	// dc.Fill()

	buckets := bucketCommitsByMonth(commits)

	stepW := float64(width) / float64(len(buckets))

	var getTotalChanges = func(commit Commit) int {
		return commit.deletions + commit.insertions
	}
	var getInsertions = func(commit Commit) int {
		return commit.insertions
	}
	var getDeletions = func(commit Commit) int {
		return commit.deletions
	}
	var getFilesChanged = func(commit Commit) int {
		return commit.filesChanged
	}
	maxChanged := 0
	maxFChanged := 0
	for _, bucket := range buckets {
		value := sumFromCommits(bucket, getTotalChanges)
		fChanges := sumFromCommits(bucket, getFilesChanged)
		if value > maxChanged {
			maxChanged = value
		}
		if fChanges > maxFChanged {
			maxFChanged = fChanges
		}
	}
	scaleChanges := float64(height) / float64(maxChanged)
	scaleFChanges := float64(height) / float64(maxFChanged)
	prevX, prevY := float64(0.0), float64(height)
	for i, bucket := range buckets {
		totalChanges := sumFromCommits(bucket, getTotalChanges)
		insertions := sumFromCommits(bucket, getInsertions)
		deletions := sumFromCommits(bucket, getDeletions)
		filesChanged := sumFromCommits(bucket, getFilesChanged)
		totalHeight := float64(totalChanges) * scaleChanges
		insertHeight := float64(insertions) * scaleChanges
		deleteHeight := float64(deletions) * scaleChanges
		fHeight := float64(filesChanged) * scaleFChanges * 0.75

		dc.SetRGB(0, 1, 0)
		dc.DrawRectangle(stepW*float64(i), float64(height), stepW, -totalHeight)
		dc.SetRGB(1, 1, 1)
		dc.SetLineWidth(2)
		dc.Stroke()
		dc.SetRGB(0, 1, 0)
		dc.DrawRectangle(stepW*float64(i), float64(height)-deleteHeight, stepW, -insertHeight)
		dc.Fill()
		dc.SetRGB(1, 0, 0)
		dc.DrawRectangle(stepW*float64(i), float64(height), stepW, -deleteHeight)
		dc.Fill()

		dc.DrawLine(prevX, prevY, stepW*float64(i), float64(height)-fHeight)
		prevX, prevY = stepW*float64(i), float64(height)-fHeight
		dc.SetRGB(1, 1, 1)
		dc.SetLineWidth(1)
		dc.Stroke()
	}

	dc.SavePNG("out.png")
}
