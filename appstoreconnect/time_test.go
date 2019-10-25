package appstoreconnect

import (
	"reflect"
	"testing"
	"time"
)

func TestWeekly(t *testing.T) {
	expected := []time.Time{
		parseTime("2019-09-08"),
		parseTime("2019-09-15"),
		parseTime("2019-09-22"),
		parseTime("2019-09-29"),
		parseTime("2019-10-06"),
		parseTime("2019-10-13"),
		parseTime("2019-10-20"),
		parseTime("2019-10-27"),
	}

	tims := readTimes(NewTimeRange(parseTime("2019-09-10"), parseTime("2019-10-31"), Weekly))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Weekly not what I was expecting")
	}
}

func TestOneWeekly(t *testing.T) {
	expected := []time.Time{
		parseTime("2019-09-08"),
	}
	tims := readTimes(NewTimeRange(parseTime("2019-09-10"), parseTime("2019-09-10"), Weekly))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Weekly not what I was expecting")
	}
}

func TestDaily(t *testing.T) {
	expected := []time.Time{
		parseTime("2018-12-29"),
		parseTime("2018-12-30"),
		parseTime("2018-12-31"),
		parseTime("2019-01-01"),
		parseTime("2019-01-02"),
	}

	tims := readTimes(NewTimeRange(parseTime("2018-12-29"), parseTime("2019-01-02"), Daily))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Daily not what I was expecting")
	}
}

func TestOneDaily(t *testing.T) {
	expected := []time.Time{
		parseTime("2018-12-29"),
	}

	tims := readTimes(NewTimeRange(parseTime("2018-12-29"), parseTime("2018-12-29"), Daily))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Daily not what I was expecting")
	}
}

func TestMonthly(t *testing.T) {
	expected := []time.Time{
		parseTime("2018-12-01"),
		parseTime("2019-01-01"),
		parseTime("2019-02-01"),
		parseTime("2019-03-01"),
		parseTime("2019-04-01"),
	}

	tims := readTimes(NewTimeRange(parseTime("2018-12-29"), parseTime("2019-04-02"), Monthly))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Monthly not what I was expecting: ")
		t.Error(tims)
	}
}

func TestOneMonthly(t *testing.T) {
	expected := []time.Time{
		parseTime("2019-04-01"),
	}

	tims := readTimes(NewTimeRange(parseTime("2019-04-01"), parseTime("2019-04-30"), Monthly))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Monthly not what I was expecting: ")
		t.Error(tims)
	}
}

func TestParseYear(t *testing.T) {
	tm := NewTime("2019")
	if tm.Year() != 2019 || tm.Month() != 01 || tm.Day() != 01 {
		t.Error("unexpected year time parsed: " + tm.String())
	}
}

func TestParseMonth(t *testing.T) {
	tm := NewTime("2018-05")
	if tm.Year() != 2018 || tm.Month() != 05 || tm.Day() != 01 {
		t.Error("unexpected month time parsed: " + tm.String())
	}
}

func TestParseRangeYearly(t *testing.T) {
	tr, err := ParseTimeRange("2019")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Yearly {
		t.Error("frequency should be yearly: " + tr.Frequency.String())
	}

	if !tr.Start.Equal(tr.End) {
		t.Error("start and end should be the same for non-range")
	}

	if tr.Start.Year() != 2019 || tr.Start.Month() != 01 || tr.Start.Day() != 01 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

	tr, err = ParseTimeRange("2018:2019")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Yearly {
		t.Error("frequency should be yearly: " + tr.Frequency.String())
	}

	if tr.Start.Equal(tr.End) {
		t.Error("start and end should not be the same for range")
	}

	if tr.Start.Year() != 2018 || tr.Start.Month() != 01 || tr.Start.Day() != 01 ||
		tr.End.Year() != 2019 || tr.End.Month() != 01 || tr.End.Day() != 01 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

	// and backwards
	tr, err = ParseTimeRange("2019:2018")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Yearly {
		t.Error("frequency should be yearly: " + tr.Frequency.String())
	}

	if tr.Start.Equal(tr.End) {
		t.Error("start and end should not be the same for range")
	}

	if tr.Start.Year() != 2018 || tr.Start.Month() != 01 || tr.Start.Day() != 01 ||
		tr.End.Year() != 2019 || tr.End.Month() != 01 || tr.End.Day() != 01 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}
}

func TestParseRangeMonthly(t *testing.T) {
	tr, err := ParseTimeRange("2019-03")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Monthly {
		t.Error("frequency should be monthly: " + tr.Frequency.String())
	}

	if !tr.Start.Equal(tr.End) {
		t.Error("start and end should be the same for non-range")
	}

	if tr.Start.Year() != 2019 || tr.Start.Month() != 03 || tr.Start.Day() != 01 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

	tr, err = ParseTimeRange("2018-03:2019-08")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Monthly {
		t.Error("frequency should be monthly: " + tr.Frequency.String())
	}

	if tr.Start.Equal(tr.End) {
		t.Error("start and end should not be the same for range")
	}

	if tr.Start.Year() != 2018 || tr.Start.Month() != 03 || tr.Start.Day() != 01 ||
		tr.End.Year() != 2019 || tr.End.Month() != 8 || tr.End.Day() != 01 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

	// and backwards
	tr, err = ParseTimeRange("2018-03:2019-08")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Monthly {
		t.Error("frequency should be monthly: " + tr.Frequency.String())
	}

	if tr.Start.Equal(tr.End) {
		t.Error("start and end should not be the same for range")
	}

	if tr.Start.Year() != 2018 || tr.Start.Month() != 03 || tr.Start.Day() != 01 ||
		tr.End.Year() != 2019 || tr.End.Month() != 8 || tr.End.Day() != 01 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

}

func TestParseRangeDaily(t *testing.T) {
	tr, err := ParseTimeRange("2019-03-10")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Daily {
		t.Error("frequency should be daily: " + tr.Frequency.String())
	}

	if !tr.Start.Equal(tr.End) {
		t.Error("start and end should be the same for non-range")
	}

	if tr.Start.Year() != 2019 || tr.Start.Month() != 03 || tr.Start.Day() != 10 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

	tr, err = ParseTimeRange("2018-03-12:2019-08-03")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Daily {
		t.Error("frequency should be daily: " + tr.Frequency.String())
	}

	if tr.Start.Equal(tr.End) {
		t.Error("start and end should not be the same for range")
	}

	if tr.Start.Year() != 2018 || tr.Start.Month() != 03 || tr.Start.Day() != 12 ||
		tr.End.Year() != 2019 || tr.End.Month() != 8 || tr.End.Day() != 03 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

	// and backwards
	tr, err = ParseTimeRange("2019-08-03:2018-03-12")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Daily {
		t.Error("frequency should be daily: " + tr.Frequency.String())
	}

	if tr.Start.Equal(tr.End) {
		t.Error("start and end should not be the same for range")
	}

	if tr.Start.Year() != 2018 || tr.Start.Month() != 03 || tr.Start.Day() != 12 ||
		tr.End.Year() != 2019 || tr.End.Month() != 8 || tr.End.Day() != 03 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

}

func TestParseRangeWeekly(t *testing.T) {
	tr, err := ParseTimeRange("2019-03-w1")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Weekly {
		t.Error("frequency should be weekly: " + tr.Frequency.String())
	}

	if !tr.Start.Equal(tr.End) {
		t.Error("start and end should be the same for non-range")
	}

	// first sunday
	if tr.Start.Year() != 2019 || tr.Start.Month() != 3 || tr.Start.Day() != 3 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

	tr, err = ParseTimeRange("2018-03-w1:2019-08-w3")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Weekly {
		t.Error("frequency should be weekly: " + tr.Frequency.String())
	}

	if tr.Start.Equal(tr.End) {
		t.Error("start and end should not be the same for range")
	}

	// first sunday
	if tr.Start.Year() != 2018 || tr.Start.Month() != 3 || tr.Start.Day() != 4 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

	if tr.End.Year() != 2019 || tr.End.Month() != 8 || tr.End.Day() != 18 {
		t.Error("unexpected month or day: " + tr.End.String())
	}
	// and backwards
	tr, err = ParseTimeRange("2019-08-w3:2018-03-w1")
	if err != nil {
		t.Error(err)
	}

	if tr.Frequency != Weekly {
		t.Error("frequency should be weekly: " + tr.Frequency.String())
	}

	if tr.Start.Equal(tr.End) {
		t.Error("start and end should not be the same for range")
	}

	// first sunday
	if tr.Start.Year() != 2018 || tr.Start.Month() != 3 || tr.Start.Day() != 4 {
		t.Error("unexpected month or day: " + tr.Start.String())
	}

	if tr.End.Year() != 2019 || tr.End.Month() != 8 || tr.End.Day() != 18 {
		t.Error("unexpected month or day: " + tr.End.String())
	}

}

func readTimes(itr *TimeRange) []time.Time {
	var tims []time.Time
	for itr.Next() {
		tim := itr.Current()
		tims = append(tims, tim)
	}
	return tims
}
