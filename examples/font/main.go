package main

import (
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/warmans/gochart"
	"github.com/warmans/gochart/pkg/style"
	"golang.org/x/image/font/gofont/goregular"
)

const numPoints = 22

func main() {

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 14})

	series := gochart.NewYSeries(gochart.GenTestData(numPoints))

	canvas := gg.NewContext(800, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	yScale := gochart.NewYScale(10, series)
	xScale := gochart.NewXScale(series, 0)

	grid := gochart.New12ColGridLayout(
		gochart.GridRow{
			HeightPercent: 0.95,
			Columns: []gochart.GridColumn{
				{
					ColSpan: 1,
					El: gochart.NewStdYAxis(
						yScale,
						gochart.YFontStyles(
							style.FontFace(face),
							style.Color(color.RGBA{R: 255, A: 255}),
						),
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
					El: gochart.NewStdXAxis(
						series,
						xScale,
						gochart.XFontStyles(
							style.FontFace(face),
							style.Color(color.RGBA{B: 255, A: 255}),
						),
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
