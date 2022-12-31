package bi

import (
	"errors"
	"time"
)

const SaneDateTimeLayout = "2006-01-02 15:04"

var acceptedTimeFormats = []string{
	SaneDateTimeLayout,
	"2006-01-02",
	"2006/01/02",
	"20060102",
	"20060102 15:04",
	"2006/01/02 15:04",
	"2006-01-02T15:04",
	"2006/01/02T15:04",
	"20060102T15:04",
	time.ANSIC,
	time.Kitchen,
	time.RFC1123,
	time.RubyDate,
	time.UnixDate,
}

func SanitizeDateTime(val string) (string, error) {
	t, err := Parse(val)
	if err != nil {
		return "", err
	}
	return Format(t), nil
}

func Now() string {
	return time.Now().Format(SaneDateTimeLayout)
}

func Format(t time.Time) string {
	return t.Format(SaneDateTimeLayout)
}

func Parse(val string) (time.Time, error) {
	for _, fmt := range acceptedTimeFormats {
		if t, err := time.Parse(fmt, val); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("invalid time format")
}
