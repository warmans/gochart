package main

import (
	"github.com/warmans/gochart/pkg/style"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

const numPoints = 100

func main() {

	canvas := gg.NewContext(800, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	series := gochart.NewXYSeries(gochart.GenTestEpisodeLabels(numPoints), gochart.GenRandomTestData(numPoints, 5))
	series2 := gochart.NewXYSeries(gochart.GenTestEpisodeLabels(numPoints), gochart.GenRandomTestData(numPoints, 5))
	series3 := gochart.NewXYSeries(gochart.GenTestEpisodeLabels(numPoints), gochart.GenRandomTestData(numPoints, 5))

	yScale := gochart.NewYScale(5, series)
	xScale := gochart.NewXScale(series, 10)

	layout := gochart.NewDynamicLayout(
		gochart.NewYAxis(yScale),
		gochart.NewXAxis(series, xScale),
		[]gochart.Plot{
			gochart.NewYGrid(yScale),
			gochart.NewPointsPlot(yScale, xScale, series, gochart.PlotPointSize(2), gochart.PlotStyle(
				style.Color(color.RGBA{R: 255, A: 255})),
			),
			gochart.NewPointsPlot(yScale, xScale, series2, gochart.PlotPointSize(2), gochart.PlotStyle(
				style.Color(color.RGBA{R: 255, A: 255})),
			),
			gochart.NewPointsPlot(yScale, xScale, series3, gochart.PlotPointSize(2), gochart.PlotStyle(
				style.Color(color.RGBA{R: 255, A: 255})),
			),
		}...,
	)

	if err := layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas)); err != nil {
		panic(err)
	}

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
