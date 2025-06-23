package gochart

import "fmt"

type Label struct {
	Value string
	Tick  int
}

type XScale interface {
	NumTicks() int
	Labels() []Label
	Position(i int, b BoundingBox) float64
	Offset() float64
}

type YScale interface {
	NumTicks() int
	Labels() []Label
	MinMax() (float64, float64)
	Position(v float64, b BoundingBox) float64
}

func NewXScaleFromLabels(labels []string) *LabelXScale {
	return &LabelXScale{labels: labels}
}

type LabelXScale struct {
	labels []string
}

func (l *LabelXScale) NumTicks() int {
	return len(l.labels)
}

func (l *LabelXScale) Labels() []Label {
	labels := make([]Label, len(l.labels))
	for k, v := range l.labels {
		labels[k] = Label{Tick: k, Value: v}
	}
	return labels
}

func (l *LabelXScale) Position(i int, b BoundingBox) float64 {
	if i > l.NumTicks() {
		return b.RelX(b.W)
	}

	// the actual size available is the total width with the margins removed.
	finalScaleWidth := b.W

	normalizedPosition := normalizeToRange(float64(i), 0, float64(l.NumTicks()), 0, finalScaleWidth)

	return b.RelX(normalizedPosition)
}

func (l *LabelXScale) Offset() float64 {
	//TODO implement me
	panic("implement me")
}

func NewXScale(series Series, offset float64) *StdXScale {
	return &StdXScale{series: series, offset: offset}
}

type StdXScale struct {
	series Series
	offset float64
}

func (s *StdXScale) NumTicks() int {
	if s.series.Ys() == nil {
		return len(s.series.Ys())
	}
	return len(s.series.Xs())
}

func (s *StdXScale) Labels() []Label {
	labels := make([]Label, s.NumTicks())
	for i := 0; i < s.NumTicks(); i++ {
		labels[i] = Label{s.series.X(i), i}
	}
	return labels
}

func (s *StdXScale) Position(i int, b BoundingBox) float64 {

	if i > s.NumTicks() {
		return b.RelX(b.W - s.offset)
	}

	// the actual size available is the total width with the margins removed.
	finalScaleWidth := b.W - (s.offset * 2)

	normalizedPosition := normalizeToRange(float64(i), 0, float64(s.NumTicks()), 0, finalScaleWidth)

	return b.RelX(normalizedPosition) + s.offset
}

func (s *StdXScale) Offset() float64 {
	return s.offset
}

func NewYScale(numTicks int, series ...Series) *StdYScale {
	return &StdYScale{
		d:        series,
		numTicks: numTicks,
	}
}

type StdYScale struct {
	d        []Series
	numTicks int
}

func (r *StdYScale) MinMax() (float64, float64) {
	return floatsRange(allYData(r.d))
}

func (r *StdYScale) NumTicks() int {
	return r.numTicks //todo: should scale based on the canvas size
}

func (r *StdYScale) Labels() []Label {
	labels := make([]Label, r.NumTicks()+1)
	_, max := r.MinMax()
	for i := 0; i <= r.NumTicks(); i++ {
		labels[i] = Label{fmt.Sprintf("%0.2f", (max/float64(r.NumTicks()))*float64(i)), i}
	}
	return labels
}

func (r *StdYScale) Position(v float64, b BoundingBox) float64 {
	min, max := r.MinMax()
	return b.MapY(min, max, v)
}

func NewStackedYScale(numTicks int, series ...Series) YScale {
	return &StackedYScale{d: series, numTicks: 10}
}

type StackedYScale struct {
	d        []Series
	numTicks int
}

func (s *StackedYScale) NumTicks() int {
	return s.numTicks
}

func (s *StackedYScale) Labels() []Label {
	labels := make([]Label, s.NumTicks()+1)
	_, max := s.MinMax()
	for i := 0; i <= s.NumTicks(); i++ {
		labels[i] = Label{fmt.Sprintf("%0.2f", (max/float64(s.NumTicks()))*float64(i)), i}
	}
	return labels
}

func (s *StackedYScale) MinMax() (float64, float64) {
	min, _ := floatsRange(allYData(s.d))
	_, max := floatRange(additiveFloatMerge(allYData(s.d)))
	return min, max
}

func (s *StackedYScale) Position(v float64, b BoundingBox) float64 {
	min, max := s.MinMax()
	return b.MapY(min, max, v)
}

func NewFixedYScale(numTicks int, maxValue float64) *FixedYScale {
	return &FixedYScale{
		numTicks: numTicks,
		fixedMax: maxValue,
	}
}

type FixedYScale struct {
	numTicks int
	fixedMax float64
}

func (r *FixedYScale) MinMax() (float64, float64) {
	return 0, r.fixedMax
}

func (r *FixedYScale) NumTicks() int {
	return r.numTicks //todo: should scale based on the canvas size
}

func (r *FixedYScale) Labels() []Label {
	labels := make([]Label, r.NumTicks()+1)
	for i := 0; i <= r.NumTicks(); i++ {
		labels[i] = Label{fmt.Sprintf("%0.2f", (r.fixedMax/float64(r.NumTicks()))*float64(i)), i}
	}
	return labels
}

func (r *FixedYScale) Position(v float64, b BoundingBox) float64 {
	min, max := r.MinMax()
	return b.MapY(min, max, v)
}
