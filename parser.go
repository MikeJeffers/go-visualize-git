package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func splitCommits(raw string) []Commit {
	var commits []Commit
	re := regexp.MustCompile("[0-9]+")
	reChangeLine := regexp.MustCompile("(file changed|files changed)")
	commitDetails := strings.Split(raw, "commit ")
	for _, commit := range commitDetails {
		var details Commit
		for _, line := range strings.Split(commit, "\n") {
			fmt.Println(line)
			if strings.Contains(line, "Date:") {
				cleaned := strings.ReplaceAll(line, "Date:", "")
				cleaned = strings.Trim(cleaned, " ")
				date, err := time.Parse(time.RFC3339, cleaned)
				if err != nil {
					log.Fatalf("Failed to parse timestamp: %s", cleaned)
				}
				details.date = date
			} else if reChangeLine.MatchString(line) {
				data := re.FindAllString(line, -1)
				for _, d := range data {
					fmt.Println(d)
				}
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
