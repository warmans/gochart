package gochart

import (
	"github.com/Pallinder/go-randomdata"
)

func GenTestData(num int) []float64 {
	values := make([]float64, num)
	for i := 0; i < num; i++ {
		values[i] = float64(i) * float64(i)
	}
	return values
}

func GenTestDataFlat(num int, val float64) []float64 {
	values := make([]float64, num)
	for i := 0; i < num; i++ {
		values[i] = val
	}
	return values
}

func GenTestDataReversed(num int) []float64 {
	values := make([]float64, num)
	for i := 0; i < num; i++ {
		values[(num-1)-i] = float64(i) * float64(i)
	}
	return values
}

func GenTestTextLabels(num int) []string {
	labels := make([]string, num)
	for i := 0; i < num; i++ {
		labels[i] = randomdata.City()
	}
	return labels
}

