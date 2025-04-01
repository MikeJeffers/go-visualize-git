package commits

import (
	"sort"
)

func BucketCommitsByTimeRange(commits []Commit, days int) [][]Commit {
	buckets := [][]Commit{}
	if len(commits) < 1 {
		return buckets
	}

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.Before(commits[j].Date)
	})

	first := commits[0].Date
	last := commits[len(commits)-1].Date

	for date := first; !date.After(last); date = date.AddDate(0, 0, days) {
		next := date.AddDate(0, 0, days)
		var commitsInRange []Commit
		for _, commit := range commits {
			if commit.Date.After(next) {
				break
			} else if commit.Date.Before(date) {
				continue
			} else if commit.Date.Equal(next) {
				continue
			}
			commitsInRange = append(commitsInRange, commit)
		}
		buckets = append(buckets, commitsInRange)
	}
	return buckets
}
