package style

import (
	"golang.org/x/image/font"
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
)

var DefaultPlotOpts = Opts{
	// set a default random volume for bar fills. This can be overwritten by other options.
	func(canvas *gg.Context) {
		canvas.SetColor(RandomColor())
	},
}

var DefaultAxisOpts = Opts{
	Color(color.RGBA{A: 255}),
	LineWidth(2),
}

type Opt func(canvas *gg.Context)

type Opts []Opt

func (s Opts) Apply(canvas *gg.Context) {
	for _, o := range s {
		o(canvas)
	}
}

func Color(rgba color.RGBA) Opt {
	return func(canvas *gg.Context) {
		canvas.SetColor(rgba)
	}
}

func Dash(dashes ...float64) Opt {
	return func(canvas *gg.Context) {
		canvas.SetDash(dashes...)
	}
}

func LineWidth(width float64) Opt {
	return func(canvas *gg.Context) {
		canvas.SetLineWidth(width)
	}
}

func FontFace(fontFace font.Face) Opt {
	return func(canvas *gg.Context) {
		canvas.SetFontFace(fontFace)
	}
}

func RandomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255,
	}
}
