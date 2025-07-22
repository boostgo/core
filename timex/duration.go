package timex

import (
	"time"

	"github.com/boostgo/core/convert"
	"github.com/rs/zerolog"
)

type Duration struct {
	Nanoseconds  int64 `json:"nanoseconds" yaml:"nanoseconds"`
	Milliseconds int64 `json:"milliseconds" yaml:"milliseconds"`
	Seconds      int   `json:"seconds" yaml:"seconds"`
	Minutes      int   `json:"minutes" yaml:"minutes"`
	Hours        int   `json:"hours" yaml:"hours"`
	Days         int   `json:"days" yaml:"days"`
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

func (d Duration) IsZero() bool {
	return d.Nanoseconds == 0 && d.Milliseconds == 0 &&
		d.Seconds == 0 && d.Minutes == 0 && d.Hours == 0 &&
		d.Days == 0
}

func (d Duration) MarshalZerologObject(e *zerolog.Event) {
	e.Int64("nanoseconds", d.Nanoseconds)
	e.Int64("milliseconds", d.Milliseconds)
	e.Int("seconds", d.Seconds)
	e.Int("minutes", d.Minutes)
	e.Int("hours", d.Hours)
	e.Int("days", d.Days)
}
