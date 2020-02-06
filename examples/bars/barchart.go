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

	series3 := gochart.NewSeries(
		nil,
		gochart.GenTestDataFlat(numPoints, 50),
	)

	canvas := gg.NewContext(640, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	layout := gochart.NewLayout(
		gochart.NewVerticalAxis(series),
		gochart.NewHorizontalAxis(series, 10),
		gochart.NewBars(series, 10),
		gochart.NewPoints(series3, 10),
		gochart.NewLines(series2, 10),
	)

	layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
