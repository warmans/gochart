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
)

type Data struct {
	ImgData string

	// cfg
	NumPoints string

	ShowLeftAxis   string
	ShowRightAxis  string
	ShowBottomAxis string
	ShowGrid       string

	ShowPoints string
	ShowLines  string
	ShowBars   string

	ChartWidth  string
	ChartHeight string

	YTicks  string
	XOffset string
}

func main() {

	data := &Data{
		NumPoints: "64",

		ShowLeftAxis:   "true",
		ShowRightAxis:  "true",
		ShowBottomAxis: "true",
		ShowGrid:       "true",

		ShowPoints: "true",
		ShowLines:  "true",
		ShowBars:   "true",

		ChartWidth:  "800",
		ChartHeight: "400",

		YTicks:  "10",
		XOffset: "10",
	}

	data.ImgData = renderDataURL(data)

	vue.New(
		vue.El("#app"),
		vue.Template(`
			<div style="margin-bottom: 20px; padding-bottom: 20px; border-bottom: 1px dashed #ccc">
				<img v-bind:src="ImgData" />
			</div>
			<table>
				<tbody>
					<tr><th style="width: 10rem">Chart Size</td><th><input v-model="ChartWidth" /> x <input v-model="ChartHeight" /></td></tr>
					<tr><th>&nbsp;</th><td></td></tr>

					<tr><th>Num. Datapoints</th><td><input v-model="NumPoints" /></td></tr>
					<tr><th>Num. Y Ticks</th><td><input v-model="YTicks" /></td></tr>
					<tr><th>X Offset</th><td><input v-model="XOffset" /> (compensate for bar width)</td></tr>

					<tr><th>&nbsp;</th><td></td></tr>
					<tr><th>Show Left Axis</th><td><select v-model="ShowLeftAxis"><option value="true">YES</option> with <option value="false">NO</option></select></td></tr>
					<tr><th>Show Right Axis</th><td><select v-model="ShowRightAxis"><option value="true">YES</option><option value="false">NO</option></select></td></tr>
					<tr><th>Show Bottom Axis</th><td><select v-model="ShowBottomAxis"><option value="true">YES</option><option value="false">NO</option></select></td></tr>
					<tr><th>Show Grid</th><td><select v-model="ShowGrid"><option value="true">YES</option><option value="false">NO</option></select></td></tr>

					<tr><th>&nbsp;</th><td></td></tr>
					<tr><th>Show Points</th><td><select v-model="ShowPoints"><option value="true">YES</option><option value="false">NO</option></select></td></tr>
					<tr><th>Show Line</th><td><select v-model="ShowLines"><option value="true">YES</option><option value="false">NO</option></select></td></tr>
					<tr><th>Show Bars</th><td><select v-model="ShowBars"><option value="true">YES</option><option value="false">NO</option></select></td></tr>

					<tr><th>&nbsp;</th><td></td></tr>
					<tr><th></th><td><button v-on:click="RenderChart">Render</button></td><tr>
				</tbody>
			</table>
		`),
		vue.Data(data),
		vue.Methods(RenderChart),
	)

	select {}
}

func RenderChart(vctx vue.Context) {
	data := vctx.Data().(*Data)
	data.ImgData = renderDataURL(data)
}

func renderDataURL(cfg *Data) string {

	numPoints := int(parseFloatOrDefault(cfg.NumPoints, 64))

	series := gochart.NewYSeries(gochart.GenSinWave(numPoints))

	canvas := gg.NewContext(int(parseFloatOrDefault(cfg.ChartWidth, 800)), int(parseFloatOrDefault(cfg.ChartHeight, 400)))
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	xScale := gochart.NewXScale(series, parseFloatOrDefault(cfg.XOffset, 10))

	leftScale := gochart.NewYScale(int(parseFloatOrDefault(cfg.YTicks, 10)), series)

	//plots

	plots := []gochart.Plot{}

	if cfg.ShowGrid == "true" {
		plots = append(plots, gochart.NewYGrid(leftScale))
	}
	if cfg.ShowLines == "true" {
		plots = append(plots, gochart.NewLinesPlot(leftScale, xScale, series))
	}
	if cfg.ShowPoints == "true" {
		plots = append(plots, gochart.NewPointsPlot(leftScale, xScale, series))
	}
	if cfg.ShowBars == "true" {
		plots = append(plots, gochart.NewBarsPlot(leftScale, xScale, series))
	}

	// columns

	topRowCols := []gochart.GridColumn{}
	bottomRowCols := []gochart.GridColumn{}

	if cfg.ShowLeftAxis == "true" {
		topRowCols = append(
			topRowCols,
			gochart.GridColumn{ColSpan: 1, El: gochart.NewStdYAxis(leftScale)},
		)
		//empty column to offset the axis
		bottomRowCols = append(bottomRowCols, gochart.GridColumn{ColSpan: 1})
	}

	topRowCols = append(
		topRowCols,
		gochart.GridColumn{
			ColSpan: 10 + countFalse(cfg.ShowLeftAxis, cfg.ShowRightAxis),
			El:      gochart.NewCompositePlot(plots...),
		},
	)
	bottomRowCols = append(
		bottomRowCols,
		gochart.GridColumn{
			ColSpan: 10 + countFalse(cfg.ShowLeftAxis, cfg.ShowRightAxis),
			El:      gochart.NewStdXAxis(series, xScale),
		},
	)

	if cfg.ShowRightAxis == "true" {
		topRowCols = append(
			topRowCols,
			gochart.GridColumn{ColSpan: 1, El: gochart.NewStdYAxis(leftScale, gochart.MirrorYStdAxis())},
		)
		bottomRowCols = append(bottomRowCols, gochart.GridColumn{ColSpan: 1})
	}

	if cfg.ShowBottomAxis == "false" {
		// clear axis
		bottomRowCols = []gochart.GridColumn{}
	}

	grid := gochart.New12ColGridLayout(
		gochart.GridRow{
			HeightPercent: 0.95,
			Columns:       topRowCols,
		},
		gochart.GridRow{
			HeightPercent: 0.05,
			Columns:       bottomRowCols,
		},
	)

	grid.Render(canvas, gochart.BoundingBoxFromCanvas(canvas))

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

func countFalse(vs ...string) int64 {
	var count int64
	for _, v := range vs {
		if v == "false" {
			count++
		}
	}
	return count
}
