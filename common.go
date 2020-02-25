package gochart

import (
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/warmans/gochart/pkg/style"
)

const defaultMargin float64 = 8
const defaultTickSize float64 = 4

type Renderable interface {
	Render(canvas *gg.Context, container BoundingBox) error
}

func NewStyles(defaults ...style.Opt) Styles {
	return Styles{styleOpts: defaults}
}

type Styles struct {
	styleOpts style.Opts
}

func (a *Styles) SetStyle(opt ...style.Opt) {
	a.styleOpts = append(a.styleOpts, opt...)
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

func timeRange(v []time.Time) (time.Time, time.Time) {
	max := time.Time{}
	min := time.Unix(math.MaxInt32, 0) //2038 :(
	for _, v := range v {
		if v.After(max) {
			max = v
		}
		if v.Before(min) {
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

		// todo: this skews the labels to the left. It needs to center the ticks rather than arrange them
		// from left to right. The problem is the labels are centered on the ticks so the first and
		// last do not take up the expected space.
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

func minInt64(a, b int64) int64 {
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

func allYData(series []Series) [][]float64 {
	all := make([][]float64, 0)
	for _, s := range series {
		all = append(all, s.Ys())
	}
	return all
}

func TimeSeriesDuration(s []time.Time) time.Duration {
	min, max := timeRange(s)
	return max.Sub(min)
}
