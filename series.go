package gochart

import (
	"fmt"
	"time"
)

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

type TimeSeriesOpt func(t *TimeSeries)

func TimeFormat(formater func(t time.Time) string) TimeSeriesOpt {
	return func(t *TimeSeries) {
		t.timeFormatter = formater
	}
}

func NewTimeSeries(x []time.Time, y []float64, opts ...TimeSeriesOpt) Series {
	ts := &TimeSeries{x: x, y: y, seriesDuration: TimeSeriesDuration(x)}
	for _, opt := range opts {
		opt(ts)
	}
	return ts
}

type TimeSeries struct {
	x              []time.Time
	y              []float64
	timeFormatter  func(t time.Time) string
	seriesDuration time.Duration
}

func (t *TimeSeries) X(i int) string {
	if i < len(t.x) {
		return t.formatTime(t.x[i])
	}
	return ""
}

func (t *TimeSeries) Y(i int) float64 {
	if i < len(t.y) {
		return t.y[i]
	}
	return 0.0
}

func (t *TimeSeries) Ys() []float64 {
	return t.y
}

func (t *TimeSeries) Xs() []string {
	//if no x-axis is set then just generate numeric values from the Y indexes.
	xs := make([]string, len(t.x))
	for i := 0; i < len(t.x); i++ {
		xs[i] = t.formatTime(t.x[i])
	}
	return xs
}

func (t *TimeSeries) AdditiveMerge(add Series) Series {
	merged := &TimeSeries{
		x: make([]time.Time, len(t.Ys())),
		y: make([]float64, len(t.Ys())),
	}
	for k := range t.Ys() {
		merged.x[k] = t.x[k]
		merged.y[k] = t.Y(k)
	}
	//todo : should this use equality checks on the times?
	if add != nil {
		for k := range merged.Ys() {
			merged.y[k] += add.Y(k)
		}
	}
	return merged
}

func (t *TimeSeries) formatTime(ts time.Time) string {
	if t.timeFormatter == nil {
		switch true {
		case t.seriesDuration < time.Minute:
			return ts.Format("15:04:05")
		case t.seriesDuration < time.Hour*24:
			return ts.Format("15:04")
		default:
			return ts.Format("2006-01-02 15:04:05")
		}
	}
	return t.timeFormatter(ts)
}
