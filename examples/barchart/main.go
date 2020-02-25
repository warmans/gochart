package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

const numPoints = 8

func main() {

	// setup the canvas size and background color
	canvas := gg.NewContext(800, 400)
	canvas.SetColor(color.RGBA{255, 255, 255, 255})
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	// generate some test data
	series := gochart.NewXYSeries(gochart.GenTestTextLabels(numPoints), gochart.GenSinWave(numPoints))

	// use a 10 tick vertical scale
	yScale := gochart.NewYScale(10, series)

	// offset the x scale by 10 to prevent the bars from overlapping the sides of the chart
	// the offset depends on the bar width. So changing it to 50 means the offset needs to
	// be 25.
	xScale := gochart.NewXScale(series, 25)

	grid := gochart.New12ColGridLayout(
		gochart.GridRow{
			HeightPercent: 0.95,
			Columns: []gochart.GridColumn{
				{
					ColSpan: 1,
					El: gochart.NewYAxis(
						yScale,
					),
				},
				{
					ColSpan: 11,
					El: gochart.NewCompositePlot(
						// add a background grid along the same scale as the chart.
						gochart.NewYGrid(yScale),

						// add the bars plot
						gochart.NewBarsPlot(
							yScale,
							xScale,
							series,
							gochart.PlotBarMaxWidth(50),
						),
					),
				},
			},
		},
		gochart.GridRow{
			HeightPercent: 0.05,
			Columns: []gochart.GridColumn{
				{ColSpan: 1},
				{
					ColSpan: 11,
					El: gochart.NewXAxis(
						series,
						xScale,
					),
				},
			},
		},
	)

	grid.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
