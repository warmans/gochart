package gochart

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
)

func GenTestData(num int) []float64 {
	values := make([]float64, num)
	for i := 0; i < num; i++ {
		values[i] = float64(i) * float64(i)
	}
	return values
}

func GenRandomTestData(num int, max float64) []float64 {
	values := make([]float64, num)
	for i := 0; i < num; i++ {
		values[i] = rand.Float64() * max
	}
	return values
}

func GenTimes(num int) []time.Time {
	now, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		panic(err)
	}
	values := make([]time.Time, num)
	for i := 0; i < num; i++ {
		values[i] = now.Add(time.Hour * time.Duration(i))
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

func GenTestEpisodeLabels(num int) []string {
	labels := make([]string, num)
	series := 1
	episode := 1
	for i := 0; i < num; i++ {
		labels[i] = fmt.Sprintf("S%02dE%02d", series, episode)
		if episode > 12 {
			series++
			episode = 0
		}
		episode++
	}
	return labels
}

func GenSinWave(num int) []float64 {
	values := make([]float64, num)
	for i := 0; i < num; i++ {
		values[i] = 1 + math.Sin(float64(i))
	}
	return values
}
