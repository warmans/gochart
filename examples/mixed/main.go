package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
	"github.com/warmans/gochart/pkg/style"
)

const numPoints = 22

func main() {

	series := gochart.NewYSeries(
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

	xScale := gochart.NewXScale(series, 10)

	// Stacked Bar Plot
	stackedCharts, stackedScale := gochart.StackPlots(
		gochart.PlotWithStyles(gochart.NewBarsPlot(gochart.NewYScale(10, series), xScale, series), style.Color(color.RGBA{R: 255, A: 255}), style.Dash(5)),
		gochart.NewBarsPlot(gochart.NewYScale(10, series), xScale, series),
		gochart.NewBarsPlot(gochart.NewYScale(10, series), xScale, series),
		gochart.NewBarsPlot(gochart.NewYScale(10, series), xScale, series),
		gochart.NewBarsPlot(gochart.NewYScale(10, series), xScale, series),
	)

	layout := gochart.NewDynamicLayout(
		gochart.NewYAxis(stackedScale),
		gochart.NewXAxis(series, xScale),
		append(
			// Background grid lines
			[]gochart.Plot{gochart.NewYGrid(stackedScale)},
			append(
				stackedCharts,

				// Line Plot
				gochart.NewLinesPlot(stackedScale, xScale, series, gochart.PlotStyle(
					style.Color(color.RGBA{A: 255}),
					style.Dash(5),
				)),

				// Points Plot
				gochart.NewPointsPlot(stackedScale, xScale, series3, gochart.PlotPointSize(5), gochart.PlotStyle(
					style.Color(color.RGBA{B: 255, A: 255})),
				),
			)...,
		)...
	)

	layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
