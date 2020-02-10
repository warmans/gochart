package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

func main() {

	numPoints := 20

	series := gochart.NewSeries(
		nil,
		append(gochart.GenTestDataReversed(numPoints/2), gochart.GenTestData(numPoints/2)...),
	)

	series2 := gochart.NewSeries(
		nil,
		append(gochart.GenTestData(numPoints/2), gochart.GenTestDataReversed(numPoints/2)...),
	)

	series3 := gochart.NewSeries(
		nil,
		gochart.GenTestDataFlat(numPoints, 48.60),
	)

	canvas := gg.NewContext(640, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	yScale := gochart.NewVerticalScale(series, series2, series3)
	xScale := gochart.NewHorizontalScale(series, 10)

	layout := gochart.NewLayout(
		gochart.NewVerticalAxis(yScale),
		gochart.NewHorizontalAxis(series, xScale),

		gochart.NewBars(yScale, xScale, series),
		gochart.NewPoints(yScale, xScale, series3),
		gochart.NewLines(yScale, xScale, series2),
	)

	layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
