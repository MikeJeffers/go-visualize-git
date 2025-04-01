package commits

import (
	"fmt"
	"testing"
	"time"
)

func TestBucketingEmpty(t *testing.T) {
	buckets := BucketCommitsByTimeRange([]Commit{}, 7)
	if len(buckets) != 0 {
		t.Errorf("Unexpected length of buckets array %d", len(buckets))
	}
}

func TestBucketingSingle(t *testing.T) {
	commit := NewCommit()
	commit.SetDate("Date:   2023-12-25T01:11:27-05:00")
	buckets := BucketCommitsByTimeRange([]Commit{commit}, 7)
	fmt.Printf("%+v\n", buckets)
	if len(buckets) != 1 {
		t.Errorf("Unexpected length of buckets array %d", len(buckets))
	}
	if len(buckets[0]) != 1 {
		t.Errorf("Unexpected length of commits in bucket[0] %d", len(buckets))
	}
}

func TestBucketingOf1PerDay(t *testing.T) {
	commits := []Commit{}
	for i := range 10 {
		commit := NewCommit()
		d := time.Date(2020, 1, i, 1, 1, 1, 1, time.UTC)
		commit.SetDate("Date:   " + d.Format("2006-01-02T15:04:05-07:00"))
		commits = append(commits, commit)
	}
	buckets := BucketCommitsByTimeRange(commits, 1)
	if len(buckets) != 10 {
		t.Errorf("Unexpected length of buckets array %d", len(buckets))
	}
	for i, bucket := range buckets {
		if len(bucket) != 1 {
			t.Errorf("Unexpected length of commits in bucket[%d] %d", i, len(bucket))
		}
	}
}
func TestBucketingOf2Per2Days(t *testing.T) {
	commits := []Commit{}
	for i := range 10 {
		commit := NewCommit()
		d := time.Date(2020, time.April, i, 0, 0, 0, 0, time.UTC)
		commit.SetDate("Date:   " + d.Format("2006-01-02T15:04:05-07:00"))
		commits = append(commits, commit)
	}
	fmt.Println(commits)
	buckets := BucketCommitsByTimeRange(commits, 2)
	fmt.Println(buckets)
	if len(buckets) != 5 {
		t.Errorf("Unexpected length of buckets array %d", len(buckets))
	}
	for i, bucket := range buckets {
		if len(bucket) != 2 {
			t.Errorf("Unexpected length of commits in bucket[%d] %d", i, len(bucket))
		}
	}
}
