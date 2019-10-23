package appstoreconnect

import (
	"time"
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
	return makeTime(timeString)
}

func NewSingleTimeRange(timeString string, frequency Frequency) *TimeRange {
	return NewTimeRange(
		makeTime(timeString),
		makeTime(timeString),
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

func makeTime(value string) time.Time {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		panic(err)
	}
	return t
}
