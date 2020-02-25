package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

const numPoints = 23

func main() {

	canvas := gg.NewContext(800, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	series := gochart.NewTimeSeries(gochart.GenTimes(numPoints), gochart.GenTestData(numPoints))

	yScale := gochart.NewYScale(10, series)
	xScale := gochart.NewXScale(series, 0)

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
					El: gochart.NewLinesPlot(
						yScale,
						xScale,
						series,
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
