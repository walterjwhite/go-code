package main

import (
	"github.com/dnstap/golang-dnstap"
	"time"
)

type TimeOfDayFilter struct {
	Start TimeOfDay
	End   TimeOfDay
}

type TimeOfDay struct {
	Hour   int
	Minute int
}

/*
const timeLayout = "15:04:05 MST"

func NewTimeOfDayFilter(start string, end string) *TimeOfDayFilter {
	startTime, err := time.Parse(timeLayout, start)
	logging.Panic(err)

	endTime, err := time.Parse(timeLayout, end)
	logging.Panic(err)

	log.Printf("start time: %v\n", startTime)
	log.Printf("end time: %v\n", endTime)

	return &TimeOfDayFilter{Start: startTime, End: endTime}
}
*/

func (f *TimeOfDayFilter) Matches(dnstapRecord *dnstap.Dnstap) bool {
	if dnstapRecord.Message.QueryTimeSec != nil {
		queryTime := time.Unix(int64(*dnstapRecord.Message.QueryTimeSec), 0)
		// check if within range
		//return f.Start.Before(queryTime) && f.End.After(queryTime)
		return f.isAfterStart(queryTime) && f.isBeforeEnd(queryTime)

		// TODO: provide ability to invert
	}

	return false
}

func (f *TimeOfDayFilter) isAfterStart(queryTime time.Time) bool {
	if f.Start.Hour < queryTime.Hour() {
		return true
	}

	if f.Start.Hour == queryTime.Hour() {
		return f.Start.Minute < queryTime.Minute()
	}

	return false
}

func (f *TimeOfDayFilter) isBeforeEnd(queryTime time.Time) bool {
	if f.End.Hour > queryTime.Hour() {
		return true
	}

	if f.End.Hour == queryTime.Hour() {
		return f.End.Minute > queryTime.Minute()
	}

	return false
}
