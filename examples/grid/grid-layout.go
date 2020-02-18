package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

const numPoints = 22

func main() {

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

	xScale := gochart.NewXScale(series, 10)

	stackedCharts, stackedScale := gochart.StackPlots(
		gochart.NewBarsPlot(gochart.NewYScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewYScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewYScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewYScale(series), xScale, series),
		gochart.NewBarsPlot(gochart.NewYScale(series), xScale, series),
	)

	rightScale := gochart.NewYScale(series2)

	linePlot := gochart.NewLinesPlot(rightScale, xScale, series2)

	grid := gochart.New12ColGridLayout(
		gochart.GridRow{
			HeightFactor: 0.95,
			Columns: []gochart.GridColumn{
				{ColSpan: 1, El: gochart.NewYAxis(stackedScale)},
				{ColSpan: 10, El: gochart.NewCompositePlot(
					append([]gochart.Plot{gochart.NewYGrid(stackedScale)}, append(stackedCharts, linePlot)...)...),
				},
				{ColSpan: 1, El: gochart.NewYAxis(rightScale, gochart.MirrorYAxis())},
			},
		},
		gochart.GridRow{
			HeightFactor: 0.05,
			Columns: []gochart.GridColumn{
				{ColSpan: 1},
				{ColSpan: 10, El: gochart.NewXAxis(series, xScale)},
				{ColSpan: 1},
			},
		},
	)

	grid.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}

}
