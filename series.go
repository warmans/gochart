package gochart

import "fmt"

func NewYSeries(y []float64) Series {
	return &XYSeries{y: y}
}

func NewXYSeries(x []string, y []float64) Series {
	return &XYSeries{x: x, y: y}
}

type Series interface {
	X(i int) string
	Y(i int) float64
	Ys() []float64
	Xs() []string
	AdditiveMerge(add Series) Series
}

type XYSeries struct {
	x []string
	y []float64
}

func (s *XYSeries) X(i int) string {
	x := s.x
	if x == nil {
		x = s.Xs()
	}
	if i < len(x) {
		return x[i]
	}
	return ""
}

func (s *XYSeries) Y(i int) float64 {
	if i < len(s.y) {
		return s.y[i]
	}
	return 0.0
}

func (s *XYSeries) Ys() []float64 {
	return s.y
}

func (s *XYSeries) Xs() []string {
	if s.x == nil {
		//if no x-axis is set then just generate numeric values from the Y indexes.
		xs := make([]string, len(s.y))
		for i := 0; i < len(s.y); i++ {
			xs[i] = fmt.Sprintf("%d", i)
		}
		return xs
	}
	return s.x
}

func (s *XYSeries) AdditiveMerge(add Series) Series {
	merged := &XYSeries{
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
