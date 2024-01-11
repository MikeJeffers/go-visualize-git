package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var numberRe = regexp.MustCompile("[0-9]+")
var changesRe = regexp.MustCompile("(file changed|files changed)")
var dateCleanRe = regexp.MustCompile(`(Date:|\s)`)

type Commit struct {
	date         time.Time
	insertions   int
	deletions    int
	filesChanged int
}

func (c *Commit) IsValid() bool {
	return !c.date.IsZero()
}

func (c *Commit) ProcessCommitRaw(raw string) {
	for _, line := range strings.Split(raw, "\n") {
		c.processLine(line)
	}
}

func (c *Commit) processLine(line string) {
	if strings.Contains(line, "Date:") {
		c.setDate(line)
	} else if changesRe.MatchString(line) {
		c.setValues(line)
	}
}

func (c *Commit) setDate(line string) {
	cleaned := dateCleanRe.ReplaceAllString(line, "")
	date, err := time.Parse(time.RFC3339, cleaned)
	if err != nil {
		log.Printf("Failed to parse timestamp: %s", cleaned)
		return
	}
	c.date = date
}

func (c *Commit) setValues(line string) {
	data := numberRe.FindAllString(line, -1)
	// First result is always the files changed
	if r, e := strconv.Atoi(data[0]); e == nil {
		c.filesChanged = r
	}
	// All 3 values are found
	if len(data) == 3 {
		if r, e := strconv.Atoi(data[1]); e == nil {
			c.insertions = r
		}
		if r, e := strconv.Atoi(data[2]); e == nil {
			c.deletions = r
		}
	} else {
		//Figure out if number is insertion or deletion
		if r, e := strconv.Atoi(data[1]); e == nil {
			if strings.Contains(line, "insertion") {
				c.insertions = r
			} else {
				c.deletions = r
			}
		}
	}
}

func ParseCommits(raw string) []Commit {
	var commits []Commit
	commitRawStrings := strings.Split(raw, "commit ")
	for _, commitRaw := range commitRawStrings {
		var commitData Commit
		commitData.ProcessCommitRaw(commitRaw)
		if commitData.IsValid() {
			commits = append(commits, commitData)
		}
	}
	return commits
}
