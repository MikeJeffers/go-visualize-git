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
	for d := first; !d.After(last); d = d.AddDate(0, 1, 0) {
		next := d.AddDate(0, 1, 0)
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

func sumInsertions(arr []Commit) int {
	res := 0
	for i := 0; i < len(arr); i++ {
		res += arr[i].insertions
	}
	return res
}

func sumDeletions(arr []Commit) int {
	res := 0
	for i := 0; i < len(arr); i++ {
		res += arr[i].deletions
	}
	return res
}

func main() {
	s := getGitLogRaw(os.Args[1])
	commits := splitCommits(s)
	total := len(commits)
	fmt.Print(len(commits))
	for _, commit := range commits {
		fmt.Printf("Commit: %s \t %s \t %s \t %s \n",
			commit.date.Format(time.RFC1123Z),
			strconv.Itoa(commit.filesChanged),
			strconv.Itoa(commit.insertions),
			strconv.Itoa(commit.deletions),
		)
	}
	for i, j := 0, len(commits)-1; i < j; i, j = i+1, j-1 {
		commits[i], commits[j] = commits[j], commits[i]
	}

	width := 800
	height := 400

	//stepW := float64(width) / float64(total)

	dc := gg.NewContext(width, height)
	dc.SetRGB(255, 255, 255)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	buckets := bucketCommitsByMonth(commits)

	stepW := float64(width) / float64(len(buckets))
	for i, bucket := range buckets {
		fmt.Println(sumInsertions(bucket))
		dc.SetRGB(0, 200, 0)
		dc.DrawRectangle(stepW*float64(i), float64(height/2), stepW, float64(sumInsertions(bucket))*0.01)
		dc.Fill()
		dc.SetRGB(200, 0, 0)
		dc.DrawRectangle(stepW*float64(i), float64(height/2), stepW, -float64(sumDeletions(bucket))*0.01)
		dc.Fill()
	}

	dc.SavePNG("out.png")
}
