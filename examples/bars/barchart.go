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

	series3 := gochart.NewSeries(
		nil,
		gochart.GenTestDataFlat(numPoints, 300),
	)

	canvas := gg.NewContext(800, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	xScale := gochart.NewHorizontalScale(series, 10)

	stackedCharts, stackedScale := gochart.StackedVisualSeries(
		gochart.NewBars(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBars(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBars(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBars(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBars(gochart.NewVerticalScale(series), xScale, series),
	)

	layout := gochart.NewLayout(
		gochart.NewVerticalAxis(stackedScale),
		gochart.NewHorizontalAxis(series, xScale),
		append(
			stackedCharts,
			gochart.NewLines(stackedScale, xScale, series),
			gochart.NewPoints(stackedScale, xScale, series3),
		)...
	)

	layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
