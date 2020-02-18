package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

func main() {

	numPoints := 22

	series := gochart.NewYSeries(
		append(gochart.GenTestDataReversed(numPoints/2), gochart.GenTestData(numPoints/2)...),
	)

	series2 := gochart.NewYSeries(
		gochart.GenTestDataReversed(numPoints),
	)

	canvas := gg.NewContext(800, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	xScale := gochart.NewHorizontalScale(series, 10)

	stackedCharts, stackedScale := gochart.StackPlots(
		gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewVerticalScale(series), xScale, series),
	)

	rightScale := gochart.NewVerticalScale(series2)

	linePlot := gochart.NewLinesPlot(rightScale, xScale, series2)

	grid := gochart.New12ColGridLayout(
		gochart.GridRow{
			HeightFactor: 0.95,
			Columns: []gochart.GridColumn{
				{ColSpan: 1, El: gochart.NewVerticalAxis(stackedScale)},
				{ColSpan: 10, El: gochart.NewCompositePlot(append(stackedCharts, gochart.NewYGrid(stackedScale), linePlot)...)},
				{ColSpan: 1, El: gochart.NewVerticalAxis(rightScale, gochart.MirrorVerticalAxis())},
			},
		},
		gochart.GridRow{
			HeightFactor: 0.05,
			Columns: []gochart.GridColumn{
				{ColSpan: 1},
				{ColSpan: 10, El: gochart.NewHorizontalAxis(series, xScale)},
				{ColSpan: 1},
			},
		},
	)

	grid.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}

}
