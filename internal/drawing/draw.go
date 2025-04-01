package drawing

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strings"

	c "gogitlog/internal/commits"
	utils "gogitlog/internal/utils"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
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
		value := utils.SumFromCallable(bucket, getTotalChanges)
		fChanges := utils.SumFromCallable(bucket, getFilesChanged)
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
		totalChanges := utils.SumFromCallable(bucket, getTotalChanges)
		insertions := utils.SumFromCallable(bucket, getInsertions)
		deletions := utils.SumFromCallable(bucket, getDeletions)
		filesChanged := utils.SumFromCallable(bucket, getFilesChanged)
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

// TODO
func DrawCommitsByDirChanged(buckets [][]c.Commit, dirs []string, width, height int) {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 8})

	stepW := float64(width) / float64(len(buckets))
	maxChanged := 0
	for _, bucket := range buckets {
		value := utils.SumFromCallable(bucket, getTotalChanges)
		if value > maxChanged {
			maxChanged = value
		}
	}
	colors := []int{}
	for range dirs {
		colors = append(colors, rand.IntN(256*256*256))
	}

	legendH := 120
	dc := gg.NewContext(width, height+legendH)

	stepXLegend := width / 10
	stepYLegend := legendH / 4
	xLegend := stepXLegend / 4
	yLegend := height + stepYLegend/2
	dc.SetFontFace(face)
	for j, dir := range dirs {
		if strings.ContainsAny(dir, ".=>") {
			continue
		}
		dc.SetHexColor(fmt.Sprintf("#%06X", colors[j]))
		xL, yL, wL, hL := float64(xLegend), float64(yLegend), float64(stepXLegend)*0.95, float64(stepYLegend)*0.95
		dc.DrawRectangle(xL, yL, wL, hL)
		dc.Fill()
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(dir, xL+wL/2, yL+hL/2, 0.5, 0.5)
		xLegend += stepXLegend
		if xLegend+stepXLegend >= width {
			xLegend = stepXLegend / 2
			yLegend += stepYLegend
		}
	}

	scaleChanges := float64(height) / float64(maxChanged)
	for i, bucket := range buckets {
		y := float64(height)
		for j, dir := range dirs {
			changesInDir := utils.SumFromCallable(bucket, func(commit c.Commit) int {
				value, ok := commit.ChangesByDir[dir]
				if !ok {
					return 0
				}
				return value
			})
			cHeight := float64(changesInDir) * scaleChanges
			dc.SetHexColor(fmt.Sprintf("#%06X", colors[j]))
			// fmt.Println(fmt.Sprintf("#%06X", colors[j]))
			dc.DrawRectangle(stepW*float64(i), y, stepW, -cHeight)
			dc.Fill()
			y -= cHeight
		}
		totalChanges := utils.SumFromCallable(bucket, getTotalChanges)
		totalHeight := float64(totalChanges) * scaleChanges

		dc.SetRGB(0, 1, 0)
		dc.DrawRectangle(stepW*float64(i), float64(height), stepW, -totalHeight)
		dc.SetRGB(1, 1, 1)
		dc.SetLineWidth(2)
		dc.Stroke()
	}

	err = dc.SavePNG("out.png")
	if err != nil {
		log.Fatalf("Failed to save output!")
	}
}
