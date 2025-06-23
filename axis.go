package gochart

import (
	"github.com/fogleman/gg"
	"github.com/warmans/gochart/pkg/style"
)

type XAxis interface {
	Scale() XScale
	Render(canvas *gg.Context, b BoundingBox) error
	Height(canvas *gg.Context) float64
}

type YAxis interface {
	Scale() YScale
	Render(canvas *gg.Context, b BoundingBox) error
}

func MirrorYStdAxis() YStdAxisOpt {
	return func(ax *YStdAxis) {
		ax.cfg.Mirrored = true
	}
}

func YFontStyles(opt ...style.Opt) YStdAxisOpt {
	return func(ax *YStdAxis) {
		ax.fontStyles.SetStyle(opt...)
	}
}

func YLineStyles(opt ...style.Opt) YStdAxisOpt {
	return func(ax *YStdAxis) {
		ax.lineStyles.SetStyle(opt...)
	}
}

type YStdAxisOpt func(ax *YStdAxis)

type YStdAxisConfig struct {
	Mirrored bool
}

func NewStdYAxis(scale YScale, opts ...YStdAxisOpt) *YStdAxis {
	y := &YStdAxis{
		lineStyles: NewStyles(style.DefaultAxisOpts...),
		fontStyles: NewStyles(style.DefaultAxisOpts...),
		scale:      scale,
		cfg:        &YStdAxisConfig{},
	}
	for _, opt := range opts {
		opt(y)
	}
	return y
}

type YStdAxis struct {
	lineStyles Styles
	fontStyles Styles
	scale      YScale
	cfg        *YStdAxisConfig
}

func (a *YStdAxis) Scale() YScale {
	return a.scale
}

func (a *YStdAxis) Render(canvas *gg.Context, b BoundingBox) error {
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

func NewStdXAxis(s Series, xScale XScale, opts ...XAxisOpt) *XStdAxis {
	x := &XStdAxis{
		lineStyles: NewStyles(style.DefaultAxisOpts...),
		fontStyles: NewStyles(style.DefaultAxisOpts...),
		s:          s,
		xScale:     xScale,
		labelAlign: 0.5,
	}
	for _, o := range opts {
		o(x)
	}
	return x
}

type XAxisOpt func(ax *XStdAxis)

func XFontStyles(opt ...style.Opt) XAxisOpt {
	return func(ax *XStdAxis) {
		ax.fontStyles.SetStyle(opt...)
	}
}

func XLineStyles(opt ...style.Opt) XAxisOpt {
	return func(ax *XStdAxis) {
		ax.lineStyles.SetStyle(opt...)
	}
}

// XLabelAlign aligns the label from left to right.
// 0 = left
// 0.5 = center
// 1 = right
func XLabelAlign(align float64) XAxisOpt {
	return func(ax *XStdAxis) {
		ax.labelAlign = align
	}
}

type XStdAxis struct {
	lineStyles Styles
	fontStyles Styles
	s          Series
	xScale     XScale
	labelAlign float64
}

func (a *XStdAxis) Scale() XScale {
	return a.xScale
}

func (a *XStdAxis) Height(canvas *gg.Context) float64 {
	return canvas.FontHeight() + defaultMargin
}

func (a *XStdAxis) Render(canvas *gg.Context, b BoundingBox) error {

	canvas.Push()
	defer canvas.Pop()

	a.lineStyles.styleOpts.Apply(canvas)

	// horizontal line
	canvas.DrawLine(b.RelX(0), b.RelY(0), b.RelX(b.W), b.RelY(0))

	labels := reduceNumLabelsToFitSpace(canvas, a.xScale.Labels(), b.W)
	totalLabelsWidth := totalLabelsWidth(canvas, labels, defaultMargin*2)
	spacing := totalLabelsWidth / float64(len(labels))

	tickWidth := (b.W / float64(a.xScale.NumTicks())) - defaultMargin

	for _, label := range labels {

		linePos := a.xScale.Position(label.Tick, b) + tickWidth/2

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
			a.labelAlign,
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

type XAxisCompactOpt func(ax *XAxisCompact)

func XCompactFontStyles(opt ...style.Opt) XAxisCompactOpt {
	return func(ax *XAxisCompact) {
		ax.fontStyles.SetStyle(opt...)
	}
}

func NewCompactXAxis(labels []string, xScale XScale, opts ...XAxisCompactOpt) *XAxisCompact {
	x := &XAxisCompact{
		lineStyles: NewStyles(style.DefaultAxisOpts...),
		fontStyles: NewStyles(style.DefaultAxisOpts...),
		xScale:     xScale,
		labelAlign: 0,
		labels:     labels,
	}
	for _, o := range opts {
		o(x)
	}
	return x
}

type XAxisCompact struct {
	lineStyles Styles
	fontStyles Styles
	labels     []string
	xScale     XScale
	labelAlign float64
}

func (a *XAxisCompact) Scale() XScale {
	return a.xScale
}

func (a *XAxisCompact) Height(canvas *gg.Context) float64 {
	canvas.Push()
	defer canvas.Pop()

	// need to apply the font styles to accurately measure the string
	if a.fontStyles.styleOpts != nil {
		a.fontStyles.styleOpts.Apply(canvas)
	}
	longest := 0.0
	for _, v := range a.labels {
		newLen, _ := canvas.MeasureString(v)
		if newLen > longest {
			longest = newLen
		}
	}

	return longest + defaultMargin
}

func (a *XAxisCompact) Render(canvas *gg.Context, b BoundingBox) error {
	canvas.Push()
	defer canvas.Pop()

	a.lineStyles.styleOpts.Apply(canvas)

	// horizontal line
	canvas.DrawLine(b.RelX(0), b.RelY(0), b.RelX(b.W), b.RelY(0))

	labels := a.xScale.Labels()

	tickWidth := (b.W / float64(a.xScale.NumTicks())) - defaultMargin

	for _, label := range labels {

		linePos := a.xScale.Position(label.Tick, b) + tickWidth/2

		canvas.DrawLine(
			linePos,
			b.RelY(0),
			linePos,
			b.RelY(0)+defaultTickSize,
		)

		canvas.Push()
		a.fontStyles.styleOpts.Apply(canvas)

		canvas.RotateAbout(
			45,
			linePos-10,
			b.RelY(0),
		)

		canvas.DrawStringAnchored(
			label.Value,
			linePos,
			b.RelY(0)+defaultTickSize,
			0,
			0,
		)
		canvas.Pop()
	}

	canvas.Stroke()

	return nil
}
