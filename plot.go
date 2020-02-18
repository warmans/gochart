package gochart

import (
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart/pkg/style"
)

type Plot interface {
	Render(canvas *gg.Context, b BoundingBox) error
	ReplaceSeries(fn func(s Series) Series)
	ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale)
	VerticalScale() VerticalScale
	WithStyle(opt ...style.StyleOpt)
}

func StackPlots(vs ...Plot) ([]Plot, VerticalScale) {
	stacked := make([]Plot, len(vs))

	var originalSeries []Series
	var lastSeries Series
	for k := range vs {
		vs[k].ReplaceSeries(func(s Series) Series {
			originalSeries = append(originalSeries, s)
			merged := s.AdditiveMerge(lastSeries)
			lastSeries = merged
			return merged
		})
		stacked[(len(vs)-1)-k] = vs[k]
	}
	stackedScale := NewStackedVerticalScale(originalSeries...)
	for k := range stacked {
		stacked[k].ReplaceVerticalScale(func(s VerticalScale) VerticalScale {
			return stackedScale
		})
	}
	return stacked, stackedScale
}

func NewCompositePlot(plots ...Plot) *CompositePlot {
	return &CompositePlot{plots: plots}
}

type CompositePlot struct {
	plots []Plot
}

func (c *CompositePlot) Render(canvas *gg.Context, container BoundingBox) error {
	for _, p := range c.plots {
		if err := p.Render(canvas, container); err != nil {
			return err
		}
	}
	return nil
}

func NewPointsPlot(yScale VerticalScale, xScale HorizontalScale, s Series) Plot {
	return &PointsPlot{
		s:         s,
		pointSize: 2,
		yScale:    yScale,
		xScale:    xScale,
		styleOpts: style.DefaultPlotOpts,
	}
}

type PointsPlot struct {
	s         Series
	pointSize float64
	yScale    VerticalScale
	xScale    HorizontalScale
	styleOpts style.StyleOpts
}

func (c *PointsPlot) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	c.styleOpts.Apply(canvas)

	for i, v := range c.s.Ys() {
		canvas.DrawCircle(
			c.xScale.Position(i, b),
			c.yScale.Position(v, b),
			c.pointSize,
		)
	}
	canvas.Fill()

	return nil
}

func (c *PointsPlot) ReplaceSeries(fn func(s Series) Series) {
	c.s = fn(c.s)
}

func (c *PointsPlot) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	c.yScale = fn(c.VerticalScale())
}

func (c *PointsPlot) VerticalScale() VerticalScale {
	return c.yScale
}

func (c *PointsPlot) WithStyle(opt ...style.StyleOpt) {
	c.styleOpts = append(c.styleOpts, opt...)
}

func PlotWithStyles(p Plot, opts ...style.StyleOpt) Plot {
	p.WithStyle(opts...)
	return p
}

func NewLinesPlot(yScale VerticalScale, xScale HorizontalScale, s Series) Plot {
	return &LinesPlot{
		yScale:    yScale,
		xScale:    xScale,
		s:         s,
		pointSize: 2,
		styleOpts: style.DefaultPlotOpts,
	}
}

type LinesPlot struct {
	yScale    VerticalScale
	xScale    HorizontalScale
	s         Series
	pointSize float64
	styleOpts style.StyleOpts
}

func (c *LinesPlot) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	c.styleOpts.Apply(canvas)

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

func (c *LinesPlot) ReplaceSeries(fn func(s Series) Series) {
	c.s = fn(c.s)
}

func (c *LinesPlot) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	c.yScale = fn(c.VerticalScale())
}

func (c *LinesPlot) VerticalScale() VerticalScale {
	return c.yScale
}

func (c *LinesPlot) WithStyle(opt ...style.StyleOpt) {
	c.styleOpts = append(c.styleOpts, opt...)
}

func NewBarsPlot(yScale VerticalScale, xScale HorizontalScale, s Series) Plot {
	return &BarsPlot{
		yScale:      yScale,
		xScale:      xScale,
		s:           s,
		maxBarWidth: 20,
		styleOpts:   style.DefaultPlotOpts,
	}
}

type BarsPlot struct {
	yScale      VerticalScale
	xScale      HorizontalScale
	s           Series
	maxBarWidth float64
	styleOpts   style.StyleOpts
}

func (c *BarsPlot) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	c.styleOpts.Apply(canvas)

	maxBarWidth := math.Max(math.Min(b.W/float64(c.xScale.NumTicks()), c.maxBarWidth)-defaultMargin, 1)

	for i, v := range c.s.Ys() {
		canvas.DrawRectangle(
			c.xScale.Position(i, b)-maxBarWidth/2,
			b.RelY(b.H),
			maxBarWidth,
			0-(c.yScale.Position(0, b)-c.yScale.Position(v, b)),
		)
	}
	canvas.Fill()
	canvas.Stroke()

	return nil
}

func (c *BarsPlot) ReplaceSeries(fn func(s Series) Series) {
	c.s = fn(c.s)
}

func (c *BarsPlot) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	c.yScale = fn(c.VerticalScale())
}

func (c *BarsPlot) VerticalScale() VerticalScale {
	return c.yScale
}

func (c *BarsPlot) WithStyle(opt ...style.StyleOpt) {
	c.styleOpts = append(c.styleOpts, opt...)
}

func NewYGrid(yScale VerticalScale) Plot {
	return &YGrid{
		yScale: yScale,
		styleOpts: style.StyleOpts{
			style.Color(color.RGBA{A: 32}),
		},
	}
}

type YGrid struct {
	yScale    VerticalScale
	styleOpts style.StyleOpts
}

func (g *YGrid) Render(canvas *gg.Context, b BoundingBox) error {
	canvas.Push()
	defer canvas.Pop()

	g.styleOpts.Apply(canvas)

	_, max := g.yScale.MinMax()

	for i := range g.yScale.Labels() {

		spacing := max / float64(g.yScale.NumTicks())
		linePos := g.yScale.Position(spacing*float64(i), b)

		canvas.DrawLine(
			b.RelX(0),
			linePos,
			b.RelX(b.W),
			linePos,
		)
	}

	canvas.Stroke()

	return nil
}

func (g *YGrid) ReplaceSeries(fn func(s Series) Series) {
	// no op - grid doesn't need a series
}

func (g *YGrid) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	g.yScale = fn(g.VerticalScale())
}

func (g *YGrid) VerticalScale() VerticalScale {
	return g.yScale
}

func (g *YGrid) WithStyle(opt ...style.StyleOpt) {
	g.styleOpts = append(g.styleOpts, opt...)
}
