package appstoreconnect

import "time"

type Frequency string

type TimeRange struct {
	Start     time.Time
	End       time.Time
	Frequency Frequency
}
