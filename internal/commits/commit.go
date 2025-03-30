package commits

import (
	"fmt"
	"log"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	numberRe    = regexp.MustCompile("[0-9]+")
	changesRe   = regexp.MustCompile("(file changed|files changed)")
	dateCleanRe = regexp.MustCompile(`(Date:|\s)`)
	// Regex for "  | 12 +++---" trailing each file change
	fileChangeRe = regexp.MustCompile(`(\s+\|\s+[0-9]+\s+[+-]+)`)
)

type Commit struct {
	Date          time.Time
	Insertions    int
	Deletions     int
	FilesChanged  int
	ChangesByFile map[string]int
	ChangesByDir  map[string]int
}

func NewCommit() Commit {
	return Commit{ChangesByFile: map[string]int{}, ChangesByDir: map[string]int{}}
}

func (c *Commit) IsValid() bool {
	return !c.Date.IsZero()
}

func (c *Commit) ProcessCommitRaw(raw string) {
	for _, line := range strings.Split(raw, "\n") {
		c.processLine(line)
	}
	topLevelPaths := c.GetCommonTopLevelDirs()
	for _, path := range topLevelPaths {
		c.ChangesByDir[path] = c.GetSumChangesByPath(path)
	}
}

func (c *Commit) processLine(line string) {
	if strings.Contains(line, "Date:") {
		c.SetDate(line)
	} else if changesRe.MatchString(line) {
		c.SetValues(line)
	} else if fileChangeRe.MatchString(line) {
		c.trackFileChange(line)
	}
}

func (c *Commit) trackFileChange(line string) {
	trimmed := strings.Trim(line, " +-")
	splits := strings.Split(trimmed, " | ")
	filename := strings.Trim(splits[0], " ")
	numChanges, err := strconv.Atoi(strings.Trim(splits[1], " "))
	if err != nil {
		log.Fatalf("Failed to parse  %s", splits[1])
	}
	_, ok := c.ChangesByFile[filename]
	if !ok {
		c.ChangesByFile[filename] = numChanges
	} else {
		c.ChangesByFile[filename] += numChanges
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

func (c *Commit) GetCommonTopLevelDirs() []string {
	commonDirs := map[string]bool{}
	for key, _ := range c.ChangesByFile {
		str := key
		if strings.Contains(str, "/") {
			str = strings.Split(key, "/")[0]
		}
		commonDirs[str] = true
	}
	return slices.Sorted(maps.Keys(commonDirs))
}

func (c *Commit) GetSumChangesByPath(path string) int {
	sum := 0
	for key, val := range c.ChangesByFile {
		if strings.Contains(key, path) {
			sum += val
		}
	}
	return sum
}

func (c *Commit) String() string {
	return fmt.Sprintf("Commit: %s \t files: %s \t additions: %s \t deletions: %s",
		c.Date.Format(time.RFC1123Z),
		strconv.Itoa(c.FilesChanged),
		strconv.Itoa(c.Insertions),
		strconv.Itoa(c.Deletions),
	)
}

func (c *Commit) Equals(other *Commit) bool {
	return c.Date == other.Date && c.Deletions == other.Deletions && c.FilesChanged == other.FilesChanged && c.Insertions == other.Insertions
}
