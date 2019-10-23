package appstoreconnect

import (
	"reflect"
	"testing"
	"time"
)

func TestWeekly(t *testing.T) {
	expected := []time.Time{
		makeTime("2019-09-08"),
		makeTime("2019-09-15"),
		makeTime("2019-09-22"),
		makeTime("2019-09-29"),
		makeTime("2019-10-06"),
		makeTime("2019-10-13"),
		makeTime("2019-10-20"),
		makeTime("2019-10-27"),
	}

	tims := readTimes(NewTimeRange(makeTime("2019-09-10"), makeTime("2019-10-31"), Weekly))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Weekly not what I was expecting")
	}
}

func TestOneWeekly(t *testing.T) {
	expected := []time.Time{
		makeTime("2019-09-08"),
	}
	tims := readTimes(NewTimeRange(makeTime("2019-09-10"), makeTime("2019-09-10"), Weekly))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Weekly not what I was expecting")
	}
}

func TestDaily(t *testing.T) {
	expected := []time.Time{
		makeTime("2018-12-29"),
		makeTime("2018-12-30"),
		makeTime("2018-12-31"),
		makeTime("2019-01-01"),
		makeTime("2019-01-02"),
	}

	tims := readTimes(NewTimeRange(makeTime("2018-12-29"), makeTime("2019-01-02"), Daily))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Daily not what I was expecting")
	}
}

func TestOneDaily(t *testing.T) {
	expected := []time.Time{
		makeTime("2018-12-29"),
	}

	tims := readTimes(NewTimeRange(makeTime("2018-12-29"), makeTime("2018-12-29"), Daily))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Daily not what I was expecting")
	}
}

func TestMonthly(t *testing.T) {
	expected := []time.Time{
		makeTime("2018-12-01"),
		makeTime("2019-01-01"),
		makeTime("2019-02-01"),
		makeTime("2019-03-01"),
		makeTime("2019-04-01"),
	}

	tims := readTimes(NewTimeRange(makeTime("2018-12-29"), makeTime("2019-04-02"), Monthly))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Monthly not what I was expecting: ")
		t.Error(tims)
	}
}

func TestOneMonthly(t *testing.T) {
	expected := []time.Time{
		makeTime("2019-04-01"),
	}

	tims := readTimes(NewTimeRange(makeTime("2019-04-01"), makeTime("2019-04-30"), Monthly))
	if !reflect.DeepEqual(expected, tims) {
		t.Error("Monthly not what I was expecting: ")
		t.Error(tims)
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
