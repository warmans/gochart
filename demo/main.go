package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/norunners/vue"
	"github.com/warmans/gochart"
	"github.com/warmans/gochart/pkg/style"
)

type Data struct {
	ImgData string

	// cfg
	LineWidth string
}

func main() {

	defaults := &Data{
		LineWidth: "1.0",
	}

	vue.New(
		vue.El("#app"),
		vue.Template(`
			<div style="margin-bottom: 20px; padding: 20px; border: 1px dashed #ccc">
				<img v-bind:src="ImgData" />
			</div>
			<div>
				<div style="margin-bottom: 20px;"><label>Line Width <input v-model="LineWidth"><label></div>
				<div style="margin-bottom: 20px;"><button v-on:click="RenderChart">Render</button></div>
			</div>
		`),
		vue.Data(&Data{
			ImgData:   renderDataURL(defaults),
			LineWidth: defaults.LineWidth,
		}),
		vue.Methods(RenderChart),
	)

	select {}
}

func RenderChart(vctx vue.Context) {
	data := vctx.Data().(*Data)
	data.ImgData = renderDataURL(data)
}

func renderDataURL(cfg *Data) string {
	// setup the canvas size and background color
	canvas := gg.NewContext(882, 88)

	// generate some test data
	series := gochart.NewYSeries(gochart.GenSinWave(64))

	// use a 10 tick vertical scale
	yScale := gochart.NewYScale(10, series)
	xScale := gochart.NewXScale(series, 0)

	plot := gochart.NewLinesPlot(
		yScale,
		xScale,
		series,
		gochart.PlotStyle(
			style.LineWidth(parseFloatOrDefault(cfg.LineWidth, 1.0)),
			style.Color(color.RGBA{R: 170, G: 57, B: 57, A: 255}),
		),
	)
	if err := plot.Render(canvas, gochart.BoundingBoxFromCanvas(canvas)); err != nil {
		panic("failed to render chart: " + err.Error())
	}

	buff := &bytes.Buffer{}
	if err := canvas.EncodePNG(buff); err != nil {
		panic("failed to encode image: " + err.Error())
	}

	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(buff.Bytes()))
}

func parseFloatOrDefault(strVal string, def float64) float64 {
	flVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		return def
	}
	return flVal
}
