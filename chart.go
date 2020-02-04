package gochart

import (
	"fmt"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

type BoundingBox struct {
	X float64
	Y float64
	W float64
	H float64
}

func (b BoundingBox) ScaleX(pos float64) float64 {
	return normalizeToRange(pos, 0, b.W, b.X, b.X+b.W)
}

func (b BoundingBox) ScaleY(pos float64) float64 {
	return normalizeToRange(pos, 0, b.H, b.Y, b.Y+b.H)
}

func (b BoundingBox) RelX(pos float64) float64 {
	return b.X + pos
}

func (b BoundingBox) RelY(pos float64) float64 {
	return b.Y + pos
}

type Element interface {
	Data() *Series
	Render(canvas *gg.Context, b BoundingBox) error
}

type Series struct {
	X []string
	Y []float64
}

func (s *Series) NumTicksX() int {
	if s.X == nil {
		return len(s.Y)
	}
	return len(s.X)
}

func (s *Series) NumTicksY() int {
	return 10 //todo: how to normalize this
}

func (s *Series) YPos(i int, b BoundingBox, zeroMin bool) float64 {
	if i > len(s.Y) {
		return 0
	}
	min, max := floatRange(s.Y)
	if zeroMin {
		min = 0
	}
	return normalizeToRange(s.Y[i], min, max, b.Y, b.H)
}

func (s *Series) XPos(i int, b BoundingBox) float64 {
	if i > s.NumTicksX() {
		return 0
	}
	return normalizeToRange(float64(i), 0, float64(s.NumTicksX()-1), b.X, b.W)
}

func (s *Series) YMinMax() (float64, float64) {
	return floatRange(s.Y)
}

func NewPoints(s *Series) *Points {
	return &Points{s: s, pointSize: 2}
}

type Points struct {
	s         *Series
	pointSize float64
}

func (c *Points) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	// setup canvas
	canvas.InvertY()

	//draw data area
	canvas.SetColor(color.RGBA{
		R: 200,
		G: 200,
		B: 200,
		A: 255,
	})
	canvas.DrawRectangle(b.X, b.Y, b.W, b.H)
	canvas.Fill()

	for i := range c.s.Y {
		canvas.Push()
		canvas.SetColor(color.RGBA{255, 0, 0, 255})
		canvas.DrawCircle(
			c.s.XPos(i, b),
			c.s.YPos(i, b, true),
			c.pointSize,
		)
		canvas.Fill()
		canvas.Pop()
	}

	return nil
}

func (c *Points) Data() *Series {
	return c.s
}

func NewLayout(chart Element, yAxis Element) Element {
	return &Layout{chart: chart, yAxis: yAxis}
}

type Layout struct {
	chart Element
	yAxis Element
}

func (l *Layout) Data() *Series {
	return l.chart.Data()
}

func (l *Layout) Render(canvas *gg.Context, b BoundingBox) error {

	axisWidth := 60.0

	chartPosition := BoundingBox{
		X: b.X + axisWidth,
		Y: b.Y,
		W: b.W - axisWidth,
		H: b.H,
	}
	if err := l.chart.Render(canvas, chartPosition); err != nil {
		return err
	}

	leftAxisPosition := BoundingBox{
		X: b.X,
		Y: b.Y,
		W: axisWidth,
		H: b.H,
	}
	if err := l.yAxis.Render(canvas, leftAxisPosition); err != nil {
		return err
	}

	return nil
}

func NewVerticalAxis(s *Series) Element {
	return &VerticalAxis{s: s}
}

type VerticalAxis struct {
	s *Series
}

func (a *VerticalAxis) Data() *Series {
	return a.s
}

func (a *VerticalAxis) Render(canvas *gg.Context, b BoundingBox) error {
	canvas.Push()
	defer canvas.Pop()

	canvas.SetColor(color.RGBA{0, 0, 0, 255})
	canvas.SetLineWidth(2)

	//debugging
	//canvas.DrawRectangle(b.X, b.Y, b.W, b.H)
	//canvas.Stroke()

	// vertical line
	canvas.DrawLine(b.RelX(b.W), b.RelY(0), b.RelX(b.W), b.RelY(b.H))

	tickWidth := 4.0

	_, max := a.s.YMinMax()
	fmt.Printf("%0.2f, %0.2f, %0.2f, %0.2f\n", b.X, b.Y, b.H, b.W)

	for i := 0; i <= a.s.NumTicksY(); i++ {

		spacing := b.H / float64(a.s.NumTicksY())
		linePos := b.RelY(b.H - (spacing * float64(i)))
		value := (max / float64(a.s.NumTicksY())) * float64(i)

		canvas.DrawString(
			fmt.Sprintf("%0.2f", value),
			b.RelX(0),
			linePos,
		)
		canvas.DrawLine(
			b.RelX(b.W - tickWidth / 2),
			linePos,
			b.RelX(b.W + tickWidth / 2),
			linePos,
		)

		fmt.Printf("pos %0.2f val %0.2f\n", linePos, value)
	}

	canvas.Stroke()

	return nil
}

func floatRange(v []float64) (float64, float64) {
	max := 0.0
	min := math.MaxFloat64
	for _, v := range v {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return min, max
}

func BoundingBoxFromCanvas(ctx *gg.Context) BoundingBox {
	return BoundingBox{
		X: 20,
		Y: 20,
		W: float64(ctx.Width()) - 40,
		H: float64(ctx.Height()) - 40,
	}
}

func normalizeToRange(val, valMin, valMax, scaleMin, scaleMax float64) float64 {
	return (((val - valMin) / valMax) * scaleMax) + scaleMin
}


func GenTestData(num int) []float64 {
	values := make([]float64, num)
	for i := 0; i <  num; i++ {
		values[i] = float64(i) * float64(i)
	}
	return values
}