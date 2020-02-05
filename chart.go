package gochart

import (
	"fmt"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

const defaultMargin float64 = 4
const defaultTickSize float64 = 4

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

func NewSeries(x []string, y []float64) *Series {
	return &Series{x: x, y: y}
}

type Series struct {
	x []string
	y []float64
}

func (s *Series) X(i int) string {
	if i < len(s.x) {
		return s.x[i]
	}
	return ""
}

func (s *Series) Y(i int) float64 {
	if i < len(s.y) {
		return s.y[i]
	}
	return 0.0
}

func (s *Series) Ys() []float64 {
	return s.y
}

func (s *Series) Xs() []string {
	return s.x
}

func (s *Series) NumTicksX() int {
	if s.x == nil {
		return len(s.y)
	}
	return len(s.x)
}

func (s *Series) NumTicksY() int {
	return 10 //todo: how to normalize this
}

func (s *Series) YPos(i int, b BoundingBox, zeroMin bool) float64 {
	if i > len(s.y) {
		return b.RelY(b.H)
	}
	min, max := floatRange(s.Ys())
	if zeroMin {
		min = 0
	}
	return normalizeToRange(s.Y(i), min, max, b.Y, b.H)
}

func (s *Series) XPos(i int, b BoundingBox) float64 {
	if i > s.NumTicksX() {
		return b.RelX(b.W)
	}
	return normalizeToRange(float64(i), 0, float64(s.NumTicksX()-1), b.X, b.W)
}

func (s *Series) XPosByValue(val string, b BoundingBox) float64 {
	for k, v := range s.XLabels() {
		if val == v {
			return s.XPos(k, b)
		}
	}
	return b.RelX(b.W)
}

func (s *Series) YMinMax() (float64, float64) {
	return floatRange(s.Ys())
}

func (s *Series) YLabels() []string {
	labels := make([]string, s.NumTicksY()+1)
	_, max := s.YMinMax()
	for i := 0; i <= s.NumTicksY(); i++ {
		labels[i] = fmt.Sprintf("%0.8f", (max/float64(s.NumTicksY()))*float64(i))
	}
	return labels
}

func (s *Series) XLabels() []string {
	labels := make([]string, s.NumTicksX())
	for i := 0; i < s.NumTicksX(); i++ {
		if s.x == nil {
			labels[i] = fmt.Sprintf("%d", i)
		} else {
			labels[i] = s.X(i)
		}
	}
	return labels
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
	//canvas.SetColor(color.RGBA{
	//	R: 200,
	//	G: 200,
	//	B: 200,
	//	A: 255,
	//})
	//canvas.DrawRectangle(b.X, b.Y, b.W, b.H)
	//canvas.Fill()

	for i := range c.s.Ys() {
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

func NewLayout(chart Element, yAxis Element, xAxis Element) Element {
	return &Layout{chart: chart, yAxis: yAxis, xAxis: xAxis}
}

type Layout struct {
	chart Element
	yAxis Element
	xAxis Element
}

func (l *Layout) Data() *Series {
	return l.chart.Data()
}

func (l *Layout) Render(canvas *gg.Context, b BoundingBox) error {

	verticalAxisSize := 60.0
	horizontalAxisSize := 30.0

	chartPosition := BoundingBox{
		X: b.X + verticalAxisSize,
		Y: b.Y + horizontalAxisSize,
		W: b.W - verticalAxisSize,
		H: b.H - horizontalAxisSize,
	}
	if err := l.chart.Render(canvas, chartPosition); err != nil {
		return err
	}

	leftAxisPosition := BoundingBox{
		X: b.X,
		Y: b.Y,
		W: verticalAxisSize,
		H: b.H - horizontalAxisSize,
	}
	if err := l.yAxis.Render(canvas, leftAxisPosition); err != nil {
		return err
	}

	bottomAxisPosition := BoundingBox{
		X: b.RelX(0) + verticalAxisSize,
		Y: b.RelY(b.H) - horizontalAxisSize,
		W: b.W - verticalAxisSize,
		H: horizontalAxisSize,
	}
	if err := l.xAxis.Render(canvas, bottomAxisPosition); err != nil {
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

	for i, label := range a.s.YLabels() {

		spacing := b.H / float64(a.s.NumTicksY())
		linePos := b.RelY(b.H - (spacing * float64(i)))

		canvas.DrawStringWrapped(
			truncateStringToMaxSize(canvas, label, b.W-defaultMargin),
			b.RelX(0),
			linePos,
			0,
			0.5,
			b.W-(defaultTickSize+defaultMargin),
			0,
			gg.AlignRight,
		)
		canvas.DrawLine(
			b.RelX(b.W-defaultTickSize),
			linePos,
			b.RelX(b.W),
			linePos,
		)
	}

	canvas.Stroke()

	return nil
}

func NewHorizontalAxis(s *Series) Element {
	return &HorizontalAxis{s: s}
}

type HorizontalAxis struct {
	s *Series
}

func (a *HorizontalAxis) Data() *Series {
	return a.s
}

func (a *HorizontalAxis) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	canvas.SetColor(color.RGBA{0, 0, 0, 255})
	canvas.SetLineWidth(2)

	//debugging
	//canvas.DrawRectangle(b.X, b.Y, b.W, b.H)
	//canvas.Stroke()

	// horizontal line
	canvas.DrawLine(b.RelX(0), b.RelY(0), b.RelX(b.W), b.RelY(0))

	labels := reduceNumStringsToFitSpace(canvas, a.s.XLabels(), b.W)
	spacing := b.W / float64(len(labels))

	for _, label := range labels {

		linePos := a.s.XPosByValue(label, b)

		canvas.DrawLine(
			linePos,
			b.RelY(0),
			linePos,
			b.RelY(0)+defaultTickSize,
		)

		fmt.Printf("%s: %0.2f\n", label, linePos)

		//debugging
		//canvas.DrawRectangle(linePos-spacing/2, b.RelY(0)+defaultTickSize+defaultMargin, spacing, spacing)

		canvas.DrawStringWrapped(
			label,
			linePos,
			b.RelY(0)+defaultTickSize+defaultMargin,
			0.5,
			0.5,
			spacing,
			0,
			gg.AlignCenter,
		)
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
	for i := 0; i < num; i++ {
		values[i] = float64(i) * float64(i)
	}
	return values
}

func truncateStringToMaxSize(canvas *gg.Context, s string, size float64) string {
	for {
		if len([]rune(s)) < 1 {
			return ""
		}
		w, _ := canvas.MeasureString(s)
		if w > size {
			s = s[:len([]rune(s))-1]
		} else {
			return s
		}
	}
}

func reduceNumStringsToFitSpace(canvas *gg.Context, ss []string, size float64) []string {
	for {
		// actually none fit
		if len(ss) == 0 {
			return ss
		}
		if totalStringWidth(canvas, ss, defaultMargin * 2) <= size {
			return ss
		}
		reduced := []string{}
		for k, s := range ss {
			//todo: I think the first and last values should always be in the set
			if k%2 != 0 {
				reduced = append(reduced, s)
			}
		}
		ss = reduced
	}
}

func totalStringWidth(canvas *gg.Context, ss []string, margins float64) float64 {
	total := 0.0
	for _, v := range ss {
		w, _ := canvas.MeasureString(v)
		total += w + margins
	}
	return total
}

func widestStringSize(canvas *gg.Context, ss []string) (w float64, h float64) {
	for _, s := range ss {
		ww, hh := canvas.MeasureString(s)
		if ww > w {
			w = ww
			h = hh
		}
	}
	return
}
