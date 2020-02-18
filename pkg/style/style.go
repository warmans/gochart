package style

import (
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
)

var DefaultPlotOpts = StyleOpts{
	// set a default random volume for bar fills. This can be overwritten by other options.
	func(canvas *gg.Context) {
		canvas.SetColor(RandomColor())
	},
}

type StyleOpt func(canvas *gg.Context)

type StyleOpts []StyleOpt

func (s StyleOpts) Apply(canvas *gg.Context) {
	for _, o := range s {
		o(canvas)
	}
}

func Color(rgba color.RGBA) StyleOpt {
	return func(canvas *gg.Context) {
		canvas.SetColor(rgba)
	}
}

func Dash(dashes ...float64) StyleOpt {
	return func(canvas *gg.Context) {
		canvas.SetDash(dashes...)
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
