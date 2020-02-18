package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
	"github.com/warmans/gochart/pkg/style"
)

func main() {

	numPoints := 10

	series := gochart.NewXYSeries(
		gochart.GenTestTextLabels(numPoints),
		append(gochart.GenTestDataReversed(numPoints/2), gochart.GenTestData(numPoints/2)...),
	)

	series3 := gochart.NewXYSeries(
		nil,
		gochart.GenTestDataFlat(numPoints, 50),
	)

	canvas := gg.NewContext(800, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	xScale := gochart.NewHorizontalScale(series, 10)

	// Stacked Bar Plot
	stackedCharts, stackedScale := gochart.StackPlots(
		gochart.PlotWithStyles(gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series), style.Color(color.RGBA{255, 0, 0, 255}), style.Dash(5)),
		gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series), )

	layout := gochart.NewDynamicLayout(
		gochart.NewVerticalAxis(stackedScale),
		gochart.NewHorizontalAxis(series, xScale),
		append(
			stackedCharts,

			// Background grid lines
			gochart.NewYGrid(stackedScale),

			// Line Plot
			gochart.PlotWithStyles(
				gochart.NewLinesPlot(stackedScale, xScale, series),
				style.Color(color.RGBA{0, 0, 0, 255}),
				style.Dash(5),
			),
			// Points Plot
			gochart.PlotWithStyles(
				gochart.NewPointsPlot(stackedScale, xScale, series3),
				style.Color(color.RGBA{0, 0, 255, 255}),
			),
		)...
	)

	layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
