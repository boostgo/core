package timex

import (
	"time"

	"github.com/boostgo/core/convert"
	"github.com/rs/zerolog"
)

type Duration struct {
	Nanoseconds  int64 `json:"nanoseconds"`
	Milliseconds int64 `json:"milliseconds"`
	Seconds      int   `json:"seconds"`
	Minutes      int   `json:"minutes"`
	Hours        int   `json:"hours"`
	Days         int   `json:"days"`
}

func NewDuration(duration time.Duration) Duration {
	return Duration{
		Nanoseconds:  duration.Nanoseconds(),
		Milliseconds: duration.Milliseconds(),
		Seconds:      convert.Int(duration.Seconds()),
		Minutes:      convert.Int(duration.Minutes()),
		Hours:        convert.Int(duration.Hours()),
		Days:         convert.Int(duration.Hours() / 24),
	}
}

// Duration converts the Duration struct back to time.Duration
// using the largest non-zero unit available
func (d Duration) Duration() time.Duration {
	if d.Days > 0 {
		return time.Duration(d.Days) * 24 * time.Hour
	}

	if d.Hours > 0 {
		return time.Duration(d.Hours) * time.Hour
	}

	if d.Minutes > 0 {
		return time.Duration(d.Minutes) * time.Minute
	}

	if d.Seconds > 0 {
		return time.Duration(d.Seconds) * time.Second
	}

	if d.Milliseconds > 0 {
		return time.Duration(d.Milliseconds) * time.Millisecond
	}

	if d.Nanoseconds > 0 {
		return time.Duration(d.Nanoseconds) * time.Nanosecond
	}

	// all fields are zero or negative
	return 0
}

func (d Duration) MarshalZerologObject(e *zerolog.Event) {
	e.Int64("nanoseconds", d.Nanoseconds)
	e.Int64("milliseconds", d.Milliseconds)
	e.Int("seconds", d.Seconds)
	e.Int("minutes", d.Minutes)
	e.Int("hours", d.Hours)
	e.Int("days", d.Days)
}
