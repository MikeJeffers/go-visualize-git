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

var font *truetype.Font

func init() {
	var err error
	font, err = truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
}

var getTotalChanges = func(commit c.Commit) int {
	return commit.Deletions + commit.Insertions
}
var getInsertions = func(commit c.Commit) int {
	return commit.Insertions
}
var getDeletions = func(commit c.Commit) int {
	return commit.Deletions
}

type Drawable interface {
	Draw(dc *gg.Context)
}

type DrawText struct {
	x     float64
	y     float64
	angle float64
	text  string
}

func (t *DrawText) Draw(dc *gg.Context) {
	face := truetype.NewFace(font, &truetype.Options{Size: 12})
	dc.SetFontFace(face)
	dc.SetRGB(1, 1, 1)
	dc.RotateAbout(t.angle, t.x, t.y)
	dc.DrawStringAnchored(t.text, t.x, t.y, 0, 0)
	dc.RotateAbout(-t.angle, t.x, t.y)
}

type AddDelStackedBar struct {
	x  float64
	y  float64
	y1 float64
	w  float64
	hd float64
	hi float64
}

func (b *AddDelStackedBar) Draw(dc *gg.Context) {
	dc.SetRGB(0, 1, 0)
	dc.DrawRectangle(b.x, b.y, b.w, b.hd+b.hi)
	dc.SetRGB(1, 1, 1)
	dc.SetLineWidth(2)
	dc.Stroke()
	dc.SetRGB(0, 1, 0)
	dc.DrawRectangle(b.x, b.y+b.hd, b.w, b.hi)
	dc.Fill()
	dc.SetRGB(1, 0, 0)
	dc.DrawRectangle(b.x, b.y, b.w, b.hd)
	dc.Fill()
}

func DrawAllOnNewCanvas(drawables []Drawable, width, height int) {
	dc := gg.NewContext(width, height)
	for _, drawable := range drawables {
		drawable.Draw(dc)
	}
	err := dc.SavePNG("out.png")
	if err != nil {
		log.Fatalf("Failed to save output!")
	}
}

func ComputeBarChart(buckets [][]c.Commit, width, height, margin int) []Drawable {
	stepW := float64(width-(margin*2)) / float64(len(buckets))
	maxChanged := 0
	for _, bucket := range buckets {
		value := utils.SumFromCallable(bucket, getTotalChanges)
		if value > maxChanged {
			maxChanged = value
		}
	}
	scaleChanges := float64(height) / float64(maxChanged)
	bars := []Drawable{}
	textStep := max(1, int(len(buckets)/20))
	for i, bucket := range buckets {
		insertions := utils.SumFromCallable(bucket, getInsertions)
		deletions := utils.SumFromCallable(bucket, getDeletions)
		insertHeight := float64(insertions) * scaleChanges
		deleteHeight := float64(deletions) * scaleChanges

		bar := &AddDelStackedBar{
			x:  stepW*float64(i) + float64(margin),
			y:  float64(height) + float64(margin),
			y1: float64(height) - deleteHeight + float64(margin),
			w:  stepW,
			hd: -deleteHeight,
			hi: -insertHeight,
		}
		bars = append(bars, bar)
		if len(bucket) < 1 {
			continue
		} else if i%textStep != 0 {
			continue
		}
		bars = append(bars, &DrawText{
			x:     stepW*float64(i) + float64(margin),
			y:     float64(height) + 10 + float64(margin),
			angle: gg.Radians(60),
			text:  bucket[0].Date.Format("2006/01"),
		})
	}
	return bars
}

func DrawCommits(buckets [][]c.Commit, width, height int) {
	bars := ComputeBarChart(buckets, width, height-100, 25)
	DrawAllOnNewCanvas(bars, width, height)
}

func DrawCommitsByDirChanged(buckets [][]c.Commit, dirs []string, width, height int) {

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

	stepXLegend := width / 10
	stepYLegend := legendH / 4
	xLegend := stepXLegend / 4
	yLegend := height + stepYLegend/2

	face := truetype.NewFace(font, &truetype.Options{Size: 8})
	dc := gg.NewContext(width, height+legendH)
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

	err := dc.SavePNG("out.png")
	if err != nil {
		log.Fatalf("Failed to save output!")
	}
}
