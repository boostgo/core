package dt

import "time"

const DateTimeFormat = "2006-01-02 15:04:05"

func Format(dt time.Time) string {
	return dt.Format(DateTimeFormat)
}

func FormatPtr(dt *time.Time) *string {
	if dt == nil {
		return nil
	}

	formatted := Format(*dt)
	return &formatted
}

func Parse(dt string) (time.Time, error) {
	parsed, err := time.Parse(DateTimeFormat, dt)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func ParsePtr(dt *string) (*time.Time, error) {
	if dt == nil {
		return nil, nil
	}

	parsed, err := Parse(*dt)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
