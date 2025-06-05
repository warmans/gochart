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
	ReplaceYScale(fn func(s YScale) YScale)
	YScale() YScale
	SetStyle(opt ...style.Opt)
}

type PlotOpt func(p Plot)

func PlotStyle(opt ...style.Opt) PlotOpt {
	return func(p Plot) {
		p.SetStyle(opt...)
	}
}

func PlotPointSize(size float64) PlotOpt {
	return func(p Plot) {
		if points, ok := p.(*PointsPlot); ok {
			points.pointSize = size
		}
	}
}

func StackPlots(vs ...Plot) ([]Plot, YScale) {
	stacked := make([]Plot, len(vs))

	var originalSeries []Series
	var lastSeries Series
	var maxYScaleTicks int
	for k := range vs {
		vs[k].ReplaceSeries(func(s Series) Series {
			originalSeries = append(originalSeries, s)
			merged := s.AdditiveMerge(lastSeries)
			lastSeries = merged
			return merged
		})
		stacked[(len(vs)-1)-k] = vs[k]
		if numTicks := vs[k].YScale().NumTicks(); numTicks > maxYScaleTicks {
			maxYScaleTicks = numTicks
		}
	}
	stackedScale := NewStackedYScale(maxYScaleTicks, originalSeries...)
	for k := range stacked {
		stacked[k].ReplaceYScale(func(s YScale) YScale {
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

func NewPointsPlot(yScale YScale, xScale XScale, s Series, opts ...PlotOpt) Plot {
	p := &PointsPlot{
		Styles:    NewStyles(style.DefaultPlotOpts...),
		s:         s,
		pointSize: 2,
		yScale:    yScale,
		xScale:    xScale,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

type PointsPlot struct {
	Styles
	s         Series
	pointSize float64
	yScale    YScale
	xScale    XScale
}

func (c *PointsPlot) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	c.styleOpts.Apply(canvas)

	tickWidth := b.W/float64(len(c.s.Ys())) - defaultMargin

	for i, v := range c.s.Ys() {
		canvas.DrawCircle(
			c.xScale.Position(i, b)+tickWidth/2,
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

func (c *PointsPlot) ReplaceYScale(fn func(s YScale) YScale) {
	c.yScale = fn(c.YScale())
}

func (c *PointsPlot) YScale() YScale {
	return c.yScale
}

func PlotWithStyles(p Plot, opts ...style.Opt) Plot {
	p.SetStyle(opts...)
	return p
}

func NewLinesPlot(yScale YScale, xScale XScale, s Series, opts ...PlotOpt) Plot {
	p := &LinesPlot{
		Styles: NewStyles(style.DefaultPlotOpts...),
		yScale: yScale,
		xScale: xScale,
		s:      s,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

type LinesPlot struct {
	Styles
	yScale YScale
	xScale XScale
	s      Series
}

func (c *LinesPlot) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	c.styleOpts.Apply(canvas)

	points := c.s.Ys()

	tickWidth := b.W/float64(len(points)) - defaultMargin

	previousPoint := 0.0
	for i, v := range points {

		// line is drawn from each point to the previous one, so the first one cannot be drawn
		if i < 1 {
			previousPoint = v
			continue
		}

		canvas.DrawLine(
			c.xScale.Position(i, b)+tickWidth/2,
			c.yScale.Position(v, b),
			c.xScale.Position(i-1, b)+tickWidth/2,
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

func (c *LinesPlot) ReplaceYScale(fn func(s YScale) YScale) {
	c.yScale = fn(c.YScale())
}

func (c *LinesPlot) YScale() YScale {
	return c.yScale
}

func NewBarsPlot(yScale YScale, xScale XScale, s Series, opts ...PlotOpt) *BarsPlot {
	p := &BarsPlot{
		Styles: NewStyles(style.DefaultPlotOpts...),
		yScale: yScale,
		xScale: xScale,
		s:      s,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

type BarsPlot struct {
	Styles
	yScale  YScale
	xScale  XScale
	s       Series
	styleFn func(v float64) style.Opts
}

func (c *BarsPlot) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	c.styleOpts.Apply(canvas)

	maxBarWidth := math.Max(b.W/float64(c.xScale.NumTicks())-defaultMargin, 1)

	for i, v := range c.s.Ys() {
		canvas.Push()
		if c.styleFn != nil {
			c.styleFn(v).Apply(canvas)
		}
		canvas.DrawRectangle(
			c.xScale.Position(i, b),
			b.RelY(b.H),
			maxBarWidth,
			0-(c.yScale.Position(0, b)-c.yScale.Position(v, b)),
		)
		canvas.Fill()
		canvas.Stroke()
		canvas.Pop()
	}

	return nil
}

func (c *BarsPlot) ReplaceSeries(fn func(s Series) Series) {
	c.s = fn(c.s)
}

func (c *BarsPlot) ReplaceYScale(fn func(s YScale) YScale) {
	c.yScale = fn(c.YScale())
}

func (c *BarsPlot) YScale() YScale {
	return c.yScale
}

func (c *BarsPlot) SetStyleFn(f func(v float64) style.Opts) {
	c.styleFn = f
}

func NewYGrid(yScale YScale, opts ...PlotOpt) Plot {
	p := &YGrid{
		Styles: NewStyles(style.Color(color.RGBA{A: 64})),
		yScale: yScale,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

type YGrid struct {
	Styles
	yScale YScale
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

func (g *YGrid) ReplaceYScale(fn func(s YScale) YScale) {
	g.yScale = fn(g.YScale())
}

func (g *YGrid) YScale() YScale {
	return g.yScale
}
