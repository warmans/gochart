package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
	"github.com/warmans/gochart/pkg/style"
)

const numPoints = 128

func main() {

	// setup the canvas size and background color
	canvas := gg.NewContext(882, 88)

	// generate some test data
	series := gochart.NewYSeries(gochart.GenSinWave(numPoints))

	// use a 10 tick vertical scale
	yScale := gochart.NewYScale(10, series)
	xScale := gochart.NewXScale(series, 0)

	plot := gochart.NewLinesPlot(
		yScale,
		xScale,
		series,
		gochart.PlotStyle(style.LineWidth(1), style.Color(color.RGBA{R: 170, G: 57, B: 57, A: 255})),
	)

	plot.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
