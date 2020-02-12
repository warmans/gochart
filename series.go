package gochart

func NewSeries(x []string, y []float64) *Series {
	return &Series{x: x, y: y}
}

type Series struct {
	x []string
	y []float64
}

func (s *Series) X(i int) string {
	if i < len(s.x) {
		return s.x[i]
	}
	return ""
}

func (s *Series) Y(i int) float64 {
	if i < len(s.y) {
		return s.y[i]
	}
	return 0.0
}

func (s *Series) Ys() []float64 {
	return s.y
}

func (s *Series) Xs() []string {
	return s.x
}

func (s *Series) AdditiveMerge(add *Series) *Series {
	merged := &Series{
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
