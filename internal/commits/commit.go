package commits

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Date         time.Time
	Insertions   int
	Deletions    int
	FilesChanged int
}

func (c *Commit) IsValid() bool {
	return !c.Date.IsZero()
}

func (c *Commit) ProcessCommitRaw(raw string) {
	for _, line := range strings.Split(raw, "\n") {
		c.processLine(line)
	}
}

func (c *Commit) processLine(line string) {
	if strings.Contains(line, "Date:") {
		c.SetDate(line)
	} else if changesRe.MatchString(line) {
		c.SetValues(line)
	}
}

func (c *Commit) SetDate(line string) {
	cleaned := dateCleanRe.ReplaceAllString(line, "")
	date, err := time.Parse(time.RFC3339, cleaned)
	if err != nil {
		log.Printf("Failed to parse timestamp: %s", cleaned)
		return
	}
	c.Date = date
}

func (c *Commit) SetValues(line string) {
	data := numberRe.FindAllString(line, -1)
	// First result is always the files changed
	if r, e := strconv.Atoi(data[0]); e == nil {
		c.FilesChanged = r
	}
	// All 3 values are found
	if len(data) == 3 {
		if r, e := strconv.Atoi(data[1]); e == nil {
			c.Insertions = r
		}
		if r, e := strconv.Atoi(data[2]); e == nil {
			c.Deletions = r
		}
	} else {
		//Figure out if number is insertion or deletion
		if r, e := strconv.Atoi(data[1]); e == nil {
			if strings.Contains(line, "insertion") {
				c.Insertions = r
			} else {
				c.Deletions = r
			}
		}
	}
}

func (c *Commit) String() string {
	return fmt.Sprintf("Commit: %s \t files: %s \t additions: %s \t deletions: %s",
		c.Date.Format(time.RFC1123Z),
		strconv.Itoa(c.FilesChanged),
		strconv.Itoa(c.Insertions),
		strconv.Itoa(c.Deletions),
	)
}
