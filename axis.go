package gochart

import (
	"github.com/fogleman/gg"
	"image/color"
)

func MirrorVerticalAxis() func(config *VerticalAxisConfig) {
	return func(config *VerticalAxisConfig) {
		config.Mirrored = true
	}
}

type VerticalAxisConfig struct {
	Mirrored bool
}

func NewVerticalAxis(scale VerticalScale, opts ...func(config *VerticalAxisConfig)) *VerticalAxis {
	cfg := &VerticalAxisConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return &VerticalAxis{scale: scale, cfg: cfg}
}

type VerticalAxis struct {
	scale VerticalScale
	cfg   *VerticalAxisConfig
}

func (a *VerticalAxis) Scale() VerticalScale {
	return a.scale
}

func (a *VerticalAxis) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	canvas.SetColor(color.RGBA{0, 0, 0, 255})
	canvas.SetLineWidth(2)

	verticalLinePos := b.RelX(b.W)
	if a.cfg.Mirrored {
		verticalLinePos = b.RelX(0)
	}

	// vertical line
	canvas.DrawLine(verticalLinePos, b.RelY(0), verticalLinePos, b.RelY(b.H))

	_, max := a.scale.MinMax()

	for i, label := range a.scale.Labels() {

		spacing := max / float64(a.scale.NumTicks())
		linePos := a.scale.Position(spacing*float64(i), b)

		// end position of tick line
		tickLinePos := verticalLinePos - defaultTickSize
		if a.cfg.Mirrored {
			tickLinePos = verticalLinePos + defaultTickSize
		}
		canvas.DrawLine(
			tickLinePos,
			linePos,
			verticalLinePos,
			linePos,
		)

		textAlign := gg.AlignRight
		if a.cfg.Mirrored {
			textAlign = gg.AlignLeft
		}

		textStartPos := b.RelX(0)
		if a.cfg.Mirrored {
			textStartPos = b.RelX(0) + defaultTickSize + defaultMargin
		}

		canvas.DrawStringWrapped(
			truncateStringToMaxSize(canvas, label.Value, b.W),
			textStartPos,
			linePos,
			0,
			0.5,
			b.W-(defaultTickSize+defaultMargin),
			0,
			textAlign,
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
