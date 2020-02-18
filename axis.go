package gochart

import (
	"github.com/fogleman/gg"
	"github.com/warmans/gochart/pkg/style"
)

func MirrorYAxis() YAxisOpt {
	return func(ax *YAxis) {
		ax.cfg.Mirrored = true
	}
}

func YFontStyles(opt ...style.Opt) YAxisOpt {
	return func(ax *YAxis) {
		ax.fontStyles.SetStyle(opt...)
	}
}

func YLineStyles(opt ...style.Opt) YAxisOpt {
	return func(ax *YAxis) {
		ax.lineStyles.SetStyle(opt...)
	}
}

type YAxisOpt func(ax *YAxis)

type YAxisConfig struct {
	Mirrored bool
}

func NewYAxis(scale YScale, opts ...YAxisOpt) *YAxis {
	y := &YAxis{
		lineStyles: NewStyles(style.DefaultAxisOpts...),
		fontStyles: NewStyles(style.DefaultAxisOpts...),
		scale:      scale,
		cfg:        &YAxisConfig{},
	}
	for _, opt := range opts {
		opt(y)
	}
	return y
}

type YAxis struct {
	lineStyles Styles
	fontStyles Styles
	scale      YScale
	cfg        *YAxisConfig
}

func (a *YAxis) Scale() YScale {
	return a.scale
}

func (a *YAxis) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	a.lineStyles.styleOpts.Apply(canvas)

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

		canvas.Push()
		a.fontStyles.styleOpts.Apply(canvas)

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
		canvas.Pop()
	}
	canvas.Stroke()

	return nil
}

func NewXAxis(s Series, xScale XScale, opts ...XAxisOpt) *XAxis {
	x := &XAxis{
		lineStyles: NewStyles(style.DefaultAxisOpts...),
		fontStyles: NewStyles(style.DefaultAxisOpts...),
		s:          s,
		xScale:     xScale,
	}
	for _, o := range opts {
		o(x)
	}
	return x
}

type XAxisOpt func(ax *XAxis)

func XFontStyles(opt ...style.Opt) XAxisOpt {
	return func(ax *XAxis) {
		ax.fontStyles.SetStyle(opt...)
	}
}

func XLineStyles(opt ...style.Opt) XAxisOpt {
	return func(ax *XAxis) {
		ax.lineStyles.SetStyle(opt...)
	}
}

type XAxis struct {
	lineStyles Styles
	fontStyles Styles
	s          Series
	xScale     XScale
}

func (a *XAxis) Scale() XScale {
	return a.xScale
}

func (a *XAxis) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	a.lineStyles.styleOpts.Apply(canvas)

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

		canvas.Push()
		a.fontStyles.styleOpts.Apply(canvas)

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
		canvas.Pop()
	}

	canvas.Stroke()

	return nil
}
