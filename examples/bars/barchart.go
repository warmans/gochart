package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

func main() {

	numPoints := 24

	series := gochart.NewSeries(
		nil,
		append(gochart.GenTestDataReversed(numPoints/2), gochart.GenTestData(numPoints/2)...),
	)

	series2 := gochart.NewSeries(
		nil,
		append(gochart.GenTestData(numPoints/2), gochart.GenTestDataReversed(numPoints/2)...),
	)

	// fixme: this doesn't work because the Y scaling gets messed up. Y is scaled based
	// on a different series.
	series3 := gochart.NewSeries(
		nil,
		gochart.GenTestDataFlat(numPoints, 50),
	)

	canvas := gg.NewContext(640, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	scale := gochart.NewAutomaticVerticalScale(series, series2, series3)

	layout := gochart.NewLayout(
		gochart.NewVerticalAxis(scale),
		gochart.NewHorizontalAxis(series, 10),
		gochart.NewBars(scale, series, 10),
		gochart.NewPoints(scale, series3, 10),
		gochart.NewLines(scale, series2, 10),
	)

	layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
