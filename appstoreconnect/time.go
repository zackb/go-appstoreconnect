package appstoreconnect

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrTimeFormatInvalid bad date string
	ErrTimeFormatInvalid = errors.New("timerange: invalid format")
)

type Frequency string

type TimeRange struct {
	Start     time.Time
	End       time.Time
	Frequency Frequency
	current   time.Time
	index     int
}

func NewTime(timeString string) (time.Time, error) {
	t, _, err := parseTimeAndFrequency(timeString)
	return t, err
}

func ParseTimeRange(value string) (*TimeRange, error) {
	if value == "" {
		return nil, ErrTimeFormatInvalid
	}

	parts := strings.Split(value, ":")
	t1, f, err := parseTimeAndFrequency(parts[0])
	if err != nil {
		return nil, err
	}

	var t2 time.Time
	if len(parts) > 1 {
		t2, _, err = parseTimeAndFrequency(parts[1])
		if err != nil {
			return nil, err
		}
		if t1.After(t2) {
			temp := t1
			t1 = t2
			t2 = temp
		}
	} else {
		t2 = t1
	}

	return NewTimeRange(t1, t2, f), nil
}

func NewSingleTimeRange(timeString string, frequency Frequency) (*TimeRange, error) {
	t, err := NewTime(timeString)
	if err != nil {
		return nil, err
	}
	return NewTimeRange(t, t, frequency), nil
}

func NewTimeRange(start time.Time, end time.Time, frequency Frequency) *TimeRange {
	t := TimeRange{
		Start:     start,
		End:       end,
		Frequency: frequency,
		current:   start,
		index:     0,
	}
	t.Start = roundDown(start, frequency)
	t.End = roundDown(end, frequency)
	return &t
}

func (t *TimeRange) Next() bool {
	var next time.Time
	if t.index == 0 {
		next = t.Start
	} else {
		// next = t.current.Add(t.interval)
		next = addFrequency(t.current, t.Frequency)
	}

	if t.End.Equal(next) || t.End.After(next) {
		t.current = next
		t.index++
		return true
	}

	return false
}

func (t *TimeRange) Current() time.Time {
	return t.current
}

func timeToReportDate(t time.Time, f Frequency) string {
	var format string
	switch f {
	case Daily:
		format = "2006-01-02"
	case Weekly:
		t = t.AddDate(0, 0, -int(t.Weekday()))
		format = "2006-01-02"
	case Monthly:
		format = "2006-01"
	case Yearly:
		format = "2006"
	default:
	}

	return t.Format(format)
}

func addFrequency(t time.Time, f Frequency) time.Time {
	var r time.Time
	switch f {
	case Daily:
		r = t.AddDate(0, 0, 1)
	case Weekly:
		r = t.AddDate(0, 0, -int(t.Weekday()))
		r = r.AddDate(0, 0, 7)
	case Monthly:
		r = t.AddDate(0, 1, 0)
	case Yearly:
		r = t.AddDate(1, 0, 0)
	}
	return r
}

func roundDown(t time.Time, f Frequency) time.Time {
	switch f {
	case Daily:
	case Weekly:
		t = t.AddDate(0, 0, -int(t.Weekday()))
	case Monthly:
		t = t.AddDate(0, 0, -int(t.Day())+1)
	case Yearly:
	}
	return t
}

func parseTimeAndFrequency(value string) (time.Time, Frequency, error) {
	parts := strings.Split(value, "-")
	switch len(parts) {
	case 1:
		t, err := time.Parse("2006", value)
		return t, Yearly, err
	case 2:
		t, err := time.Parse("2006-01", value)
		return t, Monthly, err
	case 3:
		if strings.HasPrefix(parts[2], "w") {
			var year, month, week int
			_, err := fmt.Sscanf(value, "%d-%d-w%d", &year, &month, &week)
			if err != nil {
				return time.Time{}, Weekly, ErrTimeFormatInvalid
			}
			t := time.Date(year, time.Month(month), week, 0, 0, 0, 0, time.UTC)
			t = t.AddDate(0, 0, week*7)
			return t, Weekly, nil
		}
		t, err := time.Parse("2006-01-02", value)
		return t, Daily, err
	}
	return time.Time{}, Daily, ErrTimeFormatInvalid
}

func Yesterday() *TimeRange {
	t := time.Now().AddDate(0, 0, -1)
	return &TimeRange{
		Start:     t,
		End:       t,
		Frequency: Daily,
	}
}
