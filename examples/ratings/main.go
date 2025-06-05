package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/warmans/gochart/pkg/style"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart"
)

func main() {

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 10})

	canvas := gg.NewContext(1400, 600)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	numPoints := rand.IntN(100)

	series := gochart.NewXYSeries(
		gochart.GenTestEpisodeLabels(numPoints),
		gochart.GenRandomTestData(numPoints, 5),
	)

	yScale := gochart.NewYScale(5, series)
	xScale := gochart.NewXScale(series, 0)

	bars := gochart.NewBarsPlot(yScale, xScale, series, gochart.PlotPointSize(2), gochart.PlotStyle(
		style.Color(color.RGBA{R: 255, A: 255})),
	)

	bars.SetStyleFn(func(v float64) style.Opts {
		if v < 1 {
			return style.Opts{style.Color(color.RGBA{R: 220, A: 255})}
		}
		if v > 1 && v < 3 {
			return style.Opts{style.Color(color.RGBA{R: 245, G: 138, B: 39, A: 255})}
		}
		return style.Opts{style.Color(color.RGBA{R: 23, G: 220, B: 0, A: 255})}

	})

	layout := gochart.NewDynamicLayout(
		gochart.NewStdYAxis(yScale),
		gochart.NewCompactXAxis(series, xScale, gochart.XCompactFontStyles(style.FontFace(face))),
		gochart.NewYGrid(yScale),
		bars,
	)

	if err := layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas)); err != nil {
		panic(err)
	}

	if err := canvas.SavePNG("./example.png"); err != nil {
		panic(err)
	}
}
