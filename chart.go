package gochart

import (
	"fmt"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

const defaultMargin float64 = 8
const defaultTickSize float64 = 4

type Label struct {
	Value string
	Tick  int
}

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
		return len(s.y) + 1
	}
	return len(s.x) + 1
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
	return normalizeToRange(s.Y(i), min, max, b.RelY(0), b.H)
}

func (s *Series) XPos(i int, b BoundingBox, offset float64) float64 {
	if i > s.NumTicksX() {
		return b.RelX(b.W) + offset
	}
	return normalizeToRange(float64(i), 0, float64(s.NumTicksX()-1), b.RelX(0), b.W) + offset
}

func (s *Series) YMinMax() (float64, float64) {
	return floatRange(s.Ys())
}

func (s *Series) YLabels() []Label {
	labels := make([]Label, s.NumTicksY()+1)
	_, max := s.YMinMax()
	for i := 0; i <= s.NumTicksY(); i++ {
		labels[i] = Label{fmt.Sprintf("%0.2f", (max/float64(s.NumTicksY()))*float64(i)), i}
	}
	return labels
}

func (s *Series) XLabels() []Label {
	labels := make([]Label, s.NumTicksX())
	for i := 0; i < s.NumTicksX(); i++ {
		if s.x == nil {
			labels[i] = Label{fmt.Sprintf("%d", i), i}
		} else {
			labels[i] = Label{s.X(i), i}
		}
	}
	return labels
}

func NewPoints(s *Series, xOffset float64) *Points {
	return &Points{s: s, pointSize: 2, xOffset: xOffset}
}

type Points struct {
	s         *Series
	pointSize float64
	xOffset   float64
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
		canvas.SetColor(color.RGBA{0, 0, 0, 255})
		canvas.DrawCircle(
			c.s.XPos(i, b, c.xOffset),
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

func NewLines(s *Series, xOffset float64) *Lines {
	return &Lines{s: s, pointSize: 2, xOffset: xOffset}
}

type Lines struct {
	s         *Series
	pointSize float64
	xOffset   float64
}

func (c *Lines) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	// setup canvas
	canvas.InvertY()

	canvas.SetColor(color.RGBA{0, 0, 255, 255})

	points := c.s.Ys()

	for i := range points {

		nextPoint := minInt(i+1, len(points)-1)

		canvas.DrawLine(
			c.s.XPos(i, b, c.xOffset),
			c.s.YPos(i, b, true),
			c.s.XPos(nextPoint, b, c.xOffset),
			c.s.YPos(nextPoint, b, true),
		)
		canvas.Stroke()
	}

	return nil
}

func (c *Lines) Data() *Series {
	return c.s
}

func NewBars(s *Series, xOffset float64) *Bars {
	return &Bars{s: s, maxBarWidth: 20, xOffset: xOffset}
}

type Bars struct {
	s           *Series
	maxBarWidth float64
	xOffset     float64
}

func (c *Bars) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	// setup canvas
	canvas.InvertY()

	maxBarWidth := math.Max(math.Min(b.W/float64(c.s.NumTicksX()), c.maxBarWidth)-defaultMargin, 1)

	canvas.SetColor(color.RGBA{255, 128, 255, 255})
	for i := range c.s.Ys() {
		canvas.DrawRectangle(
			c.s.XPos(i, b, c.xOffset)-maxBarWidth/2,
			b.RelY(0),
			maxBarWidth,
			c.s.YPos(i, b, true)-b.RelY(0),
		)
	}
	canvas.Fill()

	return nil
}

func (c *Bars) Data() *Series {
	return c.s
}

func NewLayout(yAxis Element, xAxis Element, charts ...Element) Element {
	return &FluidLayout{charts: charts, yAxis: yAxis, xAxis: xAxis}
}

// FluidLayout will resize axis to fit data.
type FluidLayout struct {
	charts []Element
	yAxis  Element
	xAxis  Element
}

func (l *FluidLayout) Data() *Series {
	return l.yAxis.Data()
}

func (l *FluidLayout) Render(canvas *gg.Context, b BoundingBox) error {

	_, maxXLabelH := widestLabelSize(canvas, l.Data().XLabels())
	maxYLabelW, _ := widestLabelSize(canvas, l.Data().YLabels())

	yAxisWidth := maxYLabelW + defaultMargin
	xAxisHeight := maxXLabelH + defaultMargin

	chartPosition := BoundingBox{
		X: b.X + yAxisWidth,
		Y: b.Y + xAxisHeight,
		W: b.W - yAxisWidth,
		H: b.H - xAxisHeight,
	}

	for _, ch := range l.charts {
		if err := ch.Render(canvas, chartPosition); err != nil {
			return err
		}
	}

	leftAxisPosition := BoundingBox{
		X: b.X,
		Y: b.Y,
		W: yAxisWidth,
		H: b.H - xAxisHeight,
	}
	if err := l.yAxis.Render(canvas, leftAxisPosition); err != nil {
		return err
	}

	bottomAxisPosition := BoundingBox{
		X: b.RelX(0) + yAxisWidth,
		Y: b.RelY(b.H) - xAxisHeight,
		W: b.W - yAxisWidth,
		H: xAxisHeight,
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
			truncateStringToMaxSize(canvas, label.Value, b.W),
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

func NewHorizontalAxis(s *Series, xOffset float64) Element {
	return &HorizontalAxis{s: s, xOffset: xOffset}
}

type HorizontalAxis struct {
	s       *Series
	xOffset float64
}

func (a *HorizontalAxis) Data() *Series {
	return a.s
}

func (a *HorizontalAxis) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	canvas.SetColor(color.RGBA{A: 255})
	canvas.SetLineWidth(2)

	//debugging
	//canvas.DrawRectangle(b.RelX(0), b.RelY(0), b.W, b.H)
	//canvas.Stroke()

	// horizontal line
	canvas.DrawLine(b.RelX(0), b.RelY(0), b.RelX(b.W)+a.xOffset, b.RelY(0))

	labels := reduceNumLabelsToFitSpace(canvas, a.s.XLabels(), b.W)
	totalLabelsWidth := totalLabelsWidth(canvas, labels, defaultMargin*2)
	spacing := totalLabelsWidth / float64(len(labels))

	for _, label := range labels {

		linePos := a.s.XPos(label.Tick, b, a.xOffset)

		canvas.DrawLine(
			linePos,
			b.RelY(0),
			linePos,
			b.RelY(0)+defaultTickSize,
		)

		fmt.Printf("%s: %0.2f, %0.2f\n", label.Value, linePos, spacing)

		//debugging
		//canvas.DrawRectangle(linePos-spacing/2, b.RelY(0)+defaultTickSize+defaultMargin, spacing, spacing)

		canvas.DrawStringWrapped(
			label.Value,
			linePos,
			b.RelY(0)+defaultTickSize+defaultMargin,
			0.5,
			0,
			spacing,
			1,
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

func reduceNumLabelsToFitSpace(canvas *gg.Context, ss []Label, size float64) []Label {
	for {
		// actually none fit
		if len(ss) == 0 {
			return ss
		}
		if totalLabelsWidth(canvas, ss, defaultMargin*2) <= size {
			return ss
		}

		reduced := []Label{}
		for k, s := range ss {
			if k%2 == 0 {
				reduced = append(reduced, s)
			}
		}
		ss = reduced
	}
}

func totalLabelsWidth(canvas *gg.Context, ss []Label, margins float64) float64 {
	total := 0.0
	for _, v := range ss {
		w, _ := canvas.MeasureString(v.Value)
		total += w + margins
	}
	return total
}

func widestLabelSize(canvas *gg.Context, ss []Label) (w float64, h float64) {
	for _, s := range ss {
		ww, hh := canvas.MeasureString(s.Value)
		if ww > w {
			w = ww
			h = hh
		}
	}
	return
}

func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}