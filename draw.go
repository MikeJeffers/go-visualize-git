package main

import (
	"github.com/fogleman/gg"
)

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

func drawCommits(buckets [][]Commit, width, height int) {
	stepW := float64(width) / float64(len(buckets))
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

	dc := gg.NewContext(width, height)
	scaleChanges := float64(height) / float64(maxChanged)
	scaleFChanges := float64(height) / float64(maxFChanged)
	prevX, prevY := float64(0), float64(height)
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
