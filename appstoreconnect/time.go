package appstoreconnect

import (
	"time"
)

type Frequency string

type TimeIter struct {
	start     time.Time
	end       time.Time
	frequency Frequency
	current   time.Time
	index     int
}

func NewTimeIterator(start time.Time, end time.Time, frequency Frequency) *TimeIter {
	t := TimeIter{
		start:     start,
		end:       end,
		frequency: frequency,
		current:   start,
		index:     0,
	}
	t.start = roundDown(start, frequency)
	t.end = roundDown(end, frequency)
	return &t
}

func (t *TimeIter) Next() bool {
	var next time.Time
	if t.index == 0 {
		next = t.start
	} else {
		// next = t.current.Add(t.interval)
		next = addFrequency(t.current, t.frequency)
	}

	if t.end.Equal(next) || t.end.After(next) {
		t.current = next
		t.index++
		return true
	}

	return false
}

func (t *TimeIter) Current() time.Time {
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
