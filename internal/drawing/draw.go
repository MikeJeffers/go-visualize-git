package drawing

import (
	"log"

	c "gogitlog/internal/commits"
	utils "gogitlog/internal/utils"

	"github.com/fogleman/gg"
)

var getTotalChanges = func(commit c.Commit) int {
	return commit.Deletions + commit.Insertions
}
var getInsertions = func(commit c.Commit) int {
	return commit.Insertions
}
var getDeletions = func(commit c.Commit) int {
	return commit.Deletions
}
var getFilesChanged = func(commit c.Commit) int {
	return commit.FilesChanged
}

func DrawCommits(buckets [][]c.Commit, width, height int) {
	stepW := float64(width) / float64(len(buckets))
	maxChanged := 0
	maxFChanged := 0
	for _, bucket := range buckets {
		value := utils.SumFromCommits(bucket, getTotalChanges)
		fChanges := utils.SumFromCommits(bucket, getFilesChanged)
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
		totalChanges := utils.SumFromCommits(bucket, getTotalChanges)
		insertions := utils.SumFromCommits(bucket, getInsertions)
		deletions := utils.SumFromCommits(bucket, getDeletions)
		filesChanged := utils.SumFromCommits(bucket, getFilesChanged)
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

	err := dc.SavePNG("out.png")
	if err != nil {
		log.Fatalf("Failed to save output!")
	}
}
