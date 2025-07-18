package timex

import (
	"slices"
	"time"
)

type NightFinder struct {
	nightStart int
	nightEnd   int
	nightRange []int
}

func NewNightFinder(nightStart, nightEnd int) *NightFinder {
	nightRange := make([]int, 0)
	for i := nightStart; ; i++ {
		if i == nightEnd {
			break
		}

		if i == 24 {
			i = 0
		}

		nightRange = append(nightRange, i)
	}

	return &NightFinder{
		nightStart: nightStart,
		nightEnd:   nightEnd,
		nightRange: nightRange,
	}
}

func (finder *NightFinder) IsNight(compare time.Time) bool {
	return slices.Contains(finder.nightRange, compare.Hour())
}

func (finder *NightFinder) IsDay(compare time.Time) bool {
	return !finder.IsNight(compare)
}
