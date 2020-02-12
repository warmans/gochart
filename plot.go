package gochart

import (
	"github.com/fogleman/gg"
	"math"
)

type Plot interface {
	Render(canvas *gg.Context, b BoundingBox) error
	ReplaceSeries(fn func(s *Series) *Series)
	ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale)
	VerticalScale() VerticalScale
}

func StackPlots(vs ...Plot) ([]Plot, VerticalScale) {
	stacked := make([]Plot, len(vs))

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

func NewPointsPlot(yScale VerticalScale, xScale HorizontalScale, s *Series) Plot {
	return &PointsPlot{
		s:         s,
		pointSize: 2,
		yScale:    yScale,
		xScale:    xScale,
	}
}

type PointsPlot struct {
	s         *Series
	pointSize float64
	yScale    VerticalScale
	xScale    HorizontalScale
}

func (c *PointsPlot) Render(canvas *gg.Context, b BoundingBox) error {

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

func (c *PointsPlot) ReplaceSeries(fn func(s *Series) *Series) {
	c.s = fn(c.s)
}

func (c *PointsPlot) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	c.yScale = fn(c.VerticalScale())
}

func (c *PointsPlot) VerticalScale() VerticalScale {
	return c.yScale
}

func NewLinesPlot(yScale VerticalScale, xScale HorizontalScale, s *Series) Plot {
	return &LinesPlot{yScale: yScale, xScale: xScale, s: s, pointSize: 2}
}

type LinesPlot struct {
	yScale    VerticalScale
	xScale    HorizontalScale
	s         *Series
	pointSize float64
}

func (c *LinesPlot) Render(canvas *gg.Context, b BoundingBox) error {

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

func (c *LinesPlot) ReplaceSeries(fn func(s *Series) *Series) {
	c.s = fn(c.s)
}

func (c *LinesPlot) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	c.yScale = fn(c.VerticalScale())
}

func (c *LinesPlot) VerticalScale() VerticalScale {
	return c.yScale
}

func NewBarsPlot(yScale VerticalScale, xScale HorizontalScale, s *Series) Plot {
	return &BarsPlot{yScale: yScale, xScale: xScale, s: s, maxBarWidth: 20}
}

type BarsPlot struct {
	yScale      VerticalScale
	xScale      HorizontalScale
	s           *Series
	maxBarWidth float64
}

func (c *BarsPlot) Render(canvas *gg.Context, b BoundingBox) error {

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

func (c *BarsPlot) ReplaceSeries(fn func(s *Series) *Series) {
	c.s = fn(c.s)
}

func (c *BarsPlot) ReplaceVerticalScale(fn func(s VerticalScale) VerticalScale) {
	c.yScale = fn(c.VerticalScale())
}

func (c *BarsPlot) VerticalScale() VerticalScale {
	return c.yScale
}
