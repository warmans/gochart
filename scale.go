package gochart

import "fmt"

type Label struct {
	Value string
	Tick  int
}

type HorizontalScale interface {
	NumTicks() int
	Labels() []Label
	Position(i int, b BoundingBox) float64
	Offset() float64
}

type VerticalScale interface {
	NumTicks() int
	Labels() []Label
	MinMax() (float64, float64)
	Position(v float64, b BoundingBox) float64
}

func NewHorizontalScale(series Series, offset float64) *StdHorizontalScale {
	return &StdHorizontalScale{series: series, offset: offset}
}

type StdHorizontalScale struct {
	series Series
	offset float64
}

func (s *StdHorizontalScale) NumTicks() int {
	if s.series.Ys() == nil {
		return len(s.series.Ys())
	}
	return len(s.series.Xs())
}

func (s *StdHorizontalScale) Labels() []Label {
	labels := make([]Label, s.NumTicks())
	for i := 0; i < s.NumTicks(); i++ {
		labels[i] = Label{s.series.X(i), i}
	}
	return labels
}

func (s *StdHorizontalScale) Position(i int, b BoundingBox) float64 {

	// todo: this offset acts weirdly with barcharts. More or less data can cause it to shift unevently. It's probably
	// calculated wrong.
	// I think it's because the offset is not taken into account in the tick spacing, just kind of added on the start.
	// Each tick probably needs to be moved by a fraction of the offset

	if i > s.NumTicks() {
		return b.RelX(b.W - s.offset)
	}

	// the actual size available is the total width with the margins removed.
	finalScaleWidth := b.W - (s.offset * 2)

	normalizedPosition := normalizeToRange(float64(i), 0, float64(s.NumTicks()-1), 0, finalScaleWidth)

	return b.RelX(normalizedPosition) + s.offset
}

func (s *StdHorizontalScale) Offset() float64 {
	return s.offset
}

func NewVerticalScale(series ...Series) *StdVerticalScale {
	return &StdVerticalScale{d: series}
}

type StdVerticalScale struct {
	d []Series
}

func (r *StdVerticalScale) MinMax() (float64, float64) {
	return floatsRange(allYData(r.d))
}

func (r *StdVerticalScale) NumTicks() int {
	return 10 //todo: should scale based on the canvas size
}

func (r *StdVerticalScale) Labels() []Label {
	labels := make([]Label, r.NumTicks()+1)
	_, max := r.MinMax()
	for i := 0; i <= r.NumTicks(); i++ {
		labels[i] = Label{fmt.Sprintf("%0.2f", (max/float64(r.NumTicks()))*float64(i)), i}
	}
	return labels
}

func (r *StdVerticalScale) Position(v float64, b BoundingBox) float64 {
	min, max := r.MinMax()
	return b.MapY(min, max, v)
}

func NewStackedVerticalScale(series ...Series) VerticalScale {
	return &StackedVerticalScale{d: series}
}

type StackedVerticalScale struct {
	d []Series
}

func (s *StackedVerticalScale) NumTicks() int {
	return 10
}

func (s *StackedVerticalScale) Labels() []Label {
	labels := make([]Label, s.NumTicks()+1)
	_, max := s.MinMax()
	for i := 0; i <= s.NumTicks(); i++ {
		labels[i] = Label{fmt.Sprintf("%0.2f", (max/float64(s.NumTicks()))*float64(i)), i}
	}
	return labels
}

func (s *StackedVerticalScale) MinMax() (float64, float64) {
	min, _ := floatsRange(allYData(s.d))
	_, max := floatRange(additiveFloatMerge(allYData(s.d)))
	return min, max
}

func (s *StackedVerticalScale) Position(v float64, b BoundingBox) float64 {
	min, max := s.MinMax()
	return b.MapY(min, max, v)
}
