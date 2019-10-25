package appstoreconnect

import (
	"errors"
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

func NewTime(timeString string) time.Time {
	return parseTime(timeString)
}

func ParseTimeRange(value string) (*TimeRange, error) {
	if value == "" {
		return nil, ErrTimeFormatInvalid
	}

	// not a range
	parts := strings.Split(value, ":")
	t1 := parseTime(parts[0])
	var t2 time.Time
	f := parseFrequency(parts[0])
	if len(parts) > 1 {
		t2 = parseTime(parts[1])
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

func NewSingleTimeRange(timeString string, frequency Frequency) *TimeRange {
	return NewTimeRange(
		parseTime(timeString),
		parseTime(timeString),
		frequency,
	)
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

func parseTime(value string) time.Time {
	fmt := "2006-01-02"
	if len(value) == 4 {
		fmt = "2006"
	} else if len(value) == 7 {
		fmt = "2006-01"
	} else if len(value) == 10 && strings.Contains(value, "w") {
		fmt = "2006-01-w2"
	}

	t, err := time.Parse(fmt, value)
	if err != nil {
		panic(err)
	}
	if strings.Contains(value, "w") {
		t = t.AddDate(0, 0, t.Day()*7)
	}
	return t
}

func parseFrequency(value string) Frequency {
	f := Daily
	switch len(value) {
	case 4:
		f = Yearly
	case 7:
		f = Monthly
	case 10:
		if strings.Contains(value, "w") {
			f = Weekly
		} else {
			f = Daily
		}
	}
	return f
}
