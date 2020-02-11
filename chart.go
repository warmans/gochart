package gochart

import (
	"fmt"
	"image/color"
	"math"

	"github.com/davecgh/go-spew/spew"
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

// MapX takes the given min/max and maps them to the box then returns the value X position within that scale.
// E.g. val:2 min:1 max: 3 of a 100x100 box will return 50
func (b BoundingBox) MapX(min, max, value float64) float64 {
	return normalizeToRange(value, min, max, b.RelX(0), b.W+b.X) //todo: is +X needed?
}

// MapY takes the given min/max and maps them to the box then returns the valu Ye position within that scale.
// E.g. val:2 min:1 max: 3 of a 100x100 box will return 50
func (b BoundingBox) MapY(min, max, value float64) float64 {
	return b.RelY(b.H) - normalizeToRange(value, min, max, 0, b.H)
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

type VisualSeries interface {
	Render(canvas *gg.Context, b BoundingBox) error
	ReplaceSeries(fn func(s *Series) *Series)
	ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale)
	VerticalScale() VerticalScale
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

func (s *Series) AdditiveMerge(add *Series) *Series {
	merged := &Series{
		x: make([]string, len(s.Ys())),
		y: make([]float64, len(s.Ys())),
	}
	for k := range s.Ys() {
		merged.x[k] = s.X(k)
		merged.y[k] = s.Y(k)
	}
	if add != nil {
		for k := range merged.Ys() {
			merged.y[k] += add.Y(k)
		}
	}
	return merged
}

func StackedVisualSeries(vs ...VisualSeries) ([]VisualSeries, VerticalScale) {
	stacked := make([]VisualSeries, len(vs))

	var originalSeries []*Series
	var lastSeries *Series
	for k := range vs {
		vs[k].ReplaceSeries(func(s *Series) *Series {
			originalSeries = append(originalSeries, s)
			merged := s.AdditiveMerge(lastSeries)
			lastSeries = merged
			return merged
		})
		stacked[(len(vs)-1)-k] = vs[k]
	}

	spew.Dump(originalSeries)
	stackedScale := NewStackedVerticalScale(originalSeries...)
	for k := range stacked {
		stacked[k].ReplaceVerticalScale(func(s VerticalScale) VerticalScale {
			return stackedScale
		})
	}
	return stacked, stackedScale
}

func NewPoints(yScale VerticalScale, xScale HorizontalScale, s *Series) VisualSeries {
	return &Points{
		s:         s,
		pointSize: 2,
		yScale:    yScale,
		xScale:    xScale,
	}
}

type Points struct {
	s         *Series
	pointSize float64
	yScale    VerticalScale
	xScale    HorizontalScale
}

func (c *Points) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	for i, v := range c.s.Ys() {
		canvas.SetColor(RandomColor())
		canvas.DrawCircle(
			c.xScale.Position(i, b),
			c.yScale.Position(v, b),
			c.pointSize,
		)
	}
	canvas.Fill()

	return nil
}

func (c *Points) ReplaceSeries(fn func(s *Series) *Series) {
	c.s = fn(c.s)
}

func (c *Points) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	c.yScale = fn(c.VerticalScale())
}

func (c *Points) VerticalScale() VerticalScale {
	return c.yScale
}

func NewLines(yScale VerticalScale, xScale HorizontalScale, s *Series) VisualSeries {
	return &Lines{yScale: yScale, xScale: xScale, s: s, pointSize: 2}
}

type Lines struct {
	yScale    VerticalScale
	xScale    HorizontalScale
	s         *Series
	pointSize float64
}

func (c *Lines) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	canvas.SetColor(RandomColor())

	points := c.s.Ys()

	previousPoint := 0.0
	for i, v := range points {

		// line is drawn from each point to the previous one, so the first one cannot be drawn
		if i < 1 {
			previousPoint = v
			continue
		}

		canvas.DrawLine(
			c.xScale.Position(i, b),
			c.yScale.Position(v, b),
			c.xScale.Position(i-1, b),
			c.yScale.Position(previousPoint, b),
		)
		canvas.Stroke()

		previousPoint = v
	}

	return nil
}

func (c *Lines) ReplaceSeries(fn func(s *Series) *Series) {
	c.s = fn(c.s)
}

func (c *Lines) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	c.yScale = fn(c.VerticalScale())
}

func (c *Lines) VerticalScale() VerticalScale {
	return c.yScale
}

func NewBars(yScale VerticalScale, xScale HorizontalScale, s *Series) VisualSeries {
	return &Bars{yScale: yScale, xScale: xScale, s: s, maxBarWidth: 20}
}

type Bars struct {
	yScale      VerticalScale
	xScale      HorizontalScale
	s           *Series
	maxBarWidth float64
}

func (c *Bars) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	maxBarWidth := math.Max(math.Min(b.W/float64(c.xScale.NumTicks()), c.maxBarWidth)-defaultMargin, 1)

	canvas.SetColor(RandomColor())
	for i, v := range c.s.Ys() {
		canvas.DrawRectangle(
			c.xScale.Position(i, b)-maxBarWidth/2,
			b.RelY(b.H),
			maxBarWidth,
			0-(c.yScale.Position(0, b)-c.yScale.Position(v, b)),
		)
	}
	canvas.Fill()

	return nil
}

func (c *Bars) ReplaceSeries(fn func(s *Series) *Series) {
	c.s = fn(c.s)
}

func (c *Bars) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	c.yScale = fn(c.VerticalScale())
}

func (c *Bars) VerticalScale() VerticalScale {
	return c.yScale
}

func NewLayout(yAxis *VerticalAxis, xAxis *HorizontalAxis, charts ...VisualSeries) *FluidLayout {
	return &FluidLayout{charts: charts, yAxis: yAxis, xAxis: xAxis}
}

// FluidLayout will resize axis to fit data.
type FluidLayout struct {
	charts []VisualSeries
	yAxis  *VerticalAxis
	xAxis  *HorizontalAxis
}

func (l *FluidLayout) Render(canvas *gg.Context, container BoundingBox) error {

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

type HorizontalScale interface {
	NumTicks() int
	Labels() []Label
	Position(i int, b BoundingBox) float64
	Offset() float64
}

type VerticalScale interface {
	NumTicks() int
	Labels() []Label
	MinMax() (float64, float64)
	Position(v float64, b BoundingBox) float64
}

func NewHorizontalScale(series *Series, offset float64) *StdHorizontalScale {
	return &StdHorizontalScale{series: series, offset: offset}
}

type StdHorizontalScale struct {
	series *Series
	offset float64
}

func (s *StdHorizontalScale) NumTicks() int {
	if s.series.x == nil {
		return len(s.series.y) + 1
	}
	return len(s.series.x) + 1
}

func (s *StdHorizontalScale) Labels() []Label {
	labels := make([]Label, s.NumTicks())
	for i := 0; i < s.NumTicks(); i++ {
		if s.series.x == nil {
			labels[i] = Label{fmt.Sprintf("%d", i), i}
		} else {
			labels[i] = Label{s.series.X(i), i}
		}
	}
	return labels
}

func (s *StdHorizontalScale) Position(i int, b BoundingBox) float64 {
	if i > s.NumTicks() {
		return b.RelX(b.W-s.offset) + s.offset
	}
	return normalizeToRange(float64(i), 0, float64(s.NumTicks()-1), b.RelX(0), b.W-s.offset) + s.offset
}

func (s *StdHorizontalScale) Offset() float64 {
	return s.offset
}

func NewVerticalScale(series ...*Series) *StdVerticalScale {
	return &StdVerticalScale{d: series}
}

type StdVerticalScale struct {
	d []*Series
}

func (r *StdVerticalScale) MinMax() (float64, float64) {
	return floatsRange(allYData(r.d))
}

func (r *StdVerticalScale) NumTicks() int {
	return 10 //todo: should scale based on the canvas size
}

func (r *StdVerticalScale) Labels() []Label {
	labels := make([]Label, r.NumTicks()+1)
	_, max := r.MinMax()
	for i := 0; i <= r.NumTicks(); i++ {
		labels[i] = Label{fmt.Sprintf("%0.2f", (max/float64(r.NumTicks()))*float64(i)), i}
	}
	return labels
}

func (r *StdVerticalScale) Position(v float64, b BoundingBox) float64 {
	min, max := r.MinMax()
	return b.MapY(min, max, v)
}

func NewStackedVerticalScale(series ...*Series) VerticalScale {
	return &StackedVerticalScale{d: series}
}

type StackedVerticalScale struct {
	d []*Series
}

func (s *StackedVerticalScale) NumTicks() int {
	return 10
}

func (s *StackedVerticalScale) Labels() []Label {
	labels := make([]Label, s.NumTicks()+1)
	_, max := s.MinMax()
	for i := 0; i <= s.NumTicks(); i++ {
		labels[i] = Label{fmt.Sprintf("%0.2f", (max/float64(s.NumTicks()))*float64(i)), i}
	}
	return labels
}

func (s *StackedVerticalScale) MinMax() (float64, float64) {
	min, _ := floatsRange(allYData(s.d))
	_, max := floatRange(additiveFloatMerge(allYData(s.d)))
	return min, max
}

func (s *StackedVerticalScale) Position(v float64, b BoundingBox) float64 {
	min, max := s.MinMax()
	return b.MapY(min, max, v)
}

func NewVerticalAxis(scale VerticalScale) *VerticalAxis {
	return &VerticalAxis{scale: scale}
}

type VerticalAxis struct {
	scale VerticalScale
}

func (a *VerticalAxis) Scale() VerticalScale {
	return a.scale
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

	_, max := a.scale.MinMax()

	for i, label := range a.scale.Labels() {

		spacing := max / float64(a.scale.NumTicks())
		linePos := a.scale.Position(spacing*float64(i), b)

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

func NewHorizontalAxis(s *Series, xScale HorizontalScale) *HorizontalAxis {
	return &HorizontalAxis{s: s, xScale: xScale}
}

type HorizontalAxis struct {
	s      *Series
	xScale HorizontalScale
}

func (a *HorizontalAxis) Scale() HorizontalScale {
	return a.xScale
}

func (a *HorizontalAxis) Data() *Series {
	return a.s
}

func (a *HorizontalAxis) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	canvas.SetColor(color.RGBA{A: 255})
	canvas.SetLineWidth(2)

	// horizontal line
	canvas.DrawLine(b.RelX(0), b.RelY(0), b.RelX(b.W), b.RelY(0))

	labels := reduceNumLabelsToFitSpace(canvas, a.xScale.Labels(), b.W)
	totalLabelsWidth := totalLabelsWidth(canvas, labels, defaultMargin*2)
	spacing := totalLabelsWidth / float64(len(labels))

	for _, label := range labels {

		linePos := a.xScale.Position(label.Tick, b)

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
			b.RelY(0)+defaultTickSize,
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

// simply find the min and max numbers in the given
// slices.
func floatsRange(vv [][]float64) (float64, float64) {
	overallMax := 0.0
	overallMin := math.MaxFloat64
	for _, v := range vv {
		min, max := floatRange(v)
		if min < overallMin {
			overallMin = min
		}
		if max > overallMax {
			overallMax = max
		}
	}
	return overallMin, overallMax
}

func additiveFloatMerge(slices [][]float64) []float64 {
	res := []float64{}
	for _, sl := range slices {
		for k, v := range sl {
			if len(res)-1 < k {
				res = append(res, v)
			} else {
				res[k] += v
			}
		}
	}
	return res
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

func allYData(series []*Series) [][]float64 {
	all := make([][]float64, 0)
	for _, s := range series {
		all = append(all, s.Ys())
	}
	return all
}
