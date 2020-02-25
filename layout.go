package gochart

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
)

type BoundingBox struct {
	X float64
	Y float64
	W float64
	H float64
}

// MapX takes the given min/max and maps them to the box then returns the value X position within that scale.
// E.g. val:2 min:1 max: 3 of a 100x100 box will return 50
func (b BoundingBox) MapX(min, max, value float64) float64 {
	return normalizeToRange(value, min, max, b.RelX(0), b.W+b.X) //todo: is +X needed?
}

// MapY takes the given min/max and maps them to the box then returns the valu Ye position within that scale.
// E.g. val:2 min:1 max: 3 of a 100x100 box will return 50
func (b BoundingBox) MapY(min, max, value float64) float64 {
	//todo: min value must be zero'd otherwise it pushes the scale off the chart.
	// this is a bug probably in normalizeToRange :/
	return b.RelY(b.H) - normalizeToRange(value, 0, max, 0, b.H)
}

// RelX is the relative position within the canvas i.e. 0 is the far left of the box, not the far left
// of the complete canvas.
func (b BoundingBox) RelX(pos float64) float64 {
	return b.X + pos
}

// RelY is the relative position within the canvas i.e. 0 is the bottom of the box, not the bottom
// of the complete canvas.
func (b BoundingBox) RelY(pos float64) float64 {
	return b.Y + pos
}

func (b BoundingBox) DebugRender(canvas *gg.Context) {
	canvas.Push()
	defer canvas.Pop()
	canvas.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 128})
	canvas.DrawRectangle(b.RelX(0), b.RelY(0), b.W, b.H)
	canvas.DrawString(fmt.Sprintf("x: %0.0f y: %0.0f w: %0.0f h: %0.0f", b.X, b.Y, b.W, b.H), b.X, b.Y+10)
	canvas.Stroke()
}

func NewDynamicLayout(yAxis *YAxis, xAxis *XAxis, charts ...Plot) *DynamicLayout {
	return &DynamicLayout{charts: charts, yAxis: yAxis, xAxis: xAxis}
}

// DynamicLayout will calculate size of axis based on the given data.
type DynamicLayout struct {
	charts []Plot
	yAxis  *YAxis
	xAxis  *XAxis
}

func (l *DynamicLayout) Render(canvas *gg.Context, container BoundingBox) error {

	//container.DebugRender(canvas)

	//todo multiple X axis
	_, maxXLabelH := widestLabelSize(canvas, l.xAxis.Scale().Labels())
	maxYLabelW, _ := widestLabelSize(canvas, l.yAxis.Scale().Labels())

	yAxisWidth := maxYLabelW + defaultMargin
	xAxisHeight := maxXLabelH + defaultMargin

	chartPosition := BoundingBox{
		X: container.RelX(0) + yAxisWidth,
		Y: container.RelY(0),
		W: container.W - yAxisWidth,
		H: container.H - xAxisHeight,
	}

	//chartPosition.DebugRender(canvas)

	for _, ch := range l.charts {
		if err := ch.Render(canvas, chartPosition); err != nil {
			return err
		}
	}

	leftAxisPosition := BoundingBox{
		X: container.RelX(0),
		Y: container.RelY(0),
		W: yAxisWidth,
		H: container.H - xAxisHeight,
	}
	if err := l.yAxis.Render(canvas, leftAxisPosition); err != nil {
		return err
	}

	//leftAxisPosition.DebugRender(canvas)

	bottomAxisPosition := BoundingBox{
		X: container.RelX(0) + yAxisWidth,
		Y: container.RelY(container.H) - xAxisHeight,
		W: container.W - yAxisWidth,
		H: xAxisHeight,
	}
	if err := l.xAxis.Render(canvas, bottomAxisPosition); err != nil {
		return err
	}

	//bottomAxisPosition.DebugRender(canvas)

	return nil
}

func New12ColGridLayout(rows ...GridRow) *GridLayout {
	return &GridLayout{rows: rows, numColumns: 12}
}

type GridRow struct {
	HeightPercent float64 // between 0:1 where 1 is 100% and 0.1 is 10%
	Columns       []GridColumn
}

type GridColumn struct {
	ColSpan int64 // between 1:numColumns
	El      Renderable
}

type GridLayout struct {
	numColumns int64
	rows       []GridRow
}

func (l *GridLayout) Render(canvas *gg.Context, container BoundingBox) error {

	var heightOffset float64
	for _, row := range l.rows {

		rowHeight := container.H * row.HeightPercent

		var widthOffset float64
		var numColumnsRendered int64
		for _, col := range row.Columns {

			colWidth := (container.W / float64(l.numColumns)) * float64(minInt64(col.ColSpan, l.numColumns-numColumnsRendered))
			bb := BoundingBox{
				X: container.RelX(widthOffset),
				Y: container.RelY(heightOffset),
				W: colWidth,
				H: rowHeight,
			}
			widthOffset += colWidth
			numColumnsRendered += col.ColSpan

			if col.El != nil {
				col.El.Render(canvas, bb)
			}
		}
		heightOffset += rowHeight
	}

	return nil
}
