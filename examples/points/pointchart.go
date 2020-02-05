package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

func main() {
	series := gochart.NewSeries(
		nil,
		gochart.GenTestData(100),
	)

	canvas := gg.NewContext(640, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	layout := gochart.NewLayout(
		gochart.NewPoints(series),
		gochart.NewVerticalAxis(series),
		gochart.NewHorizontalAxis(series),
	)

	layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
