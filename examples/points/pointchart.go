package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

func main() {
	series := &gochart.Series{
		Y:    gochart.GenTestData(40), //[]float64{10, 2, 3, 1, 4, 5, 6, 4, 20, 30, 40, 100, 2, 3},
	}

	canvas := gg.NewContext(640, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()


	layout := gochart.NewLayout(
		gochart.NewPoints(series),
		gochart.NewVerticalAxis(series),
	)

	layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
