package core

import (
	"log"
	"time"
)

func IsWithDate(date time.Time, start time.Time, end time.Time) bool {
	return (date.Equal(start) || date.After(start)) && (date.Before(end) || date.Equal(end))
}

var now *time.Time

func SetNow(date string) {
	t, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		log.Fatal(err)
	}
	now = &t
}

// gets todays date as time
func GetToday() time.Time {
	return ToStartOfDay(Now())
}

func GetDaysAgo(date time.Time, days int) time.Time {
	return ToStartOfDay(date.AddDate(0, 0, -days))
}

func GetTodayMinusDays(days int) time.Time {
	return GetDaysAgo(Now(), days)
}

func Now() time.Time {
	if now != nil {
		return *now
	}
	return time.Now()
}

func ToStartOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
}

func StringToDate(date string, format string) time.Time {
	t, err := time.ParseInLocation(format, date, time.Local)
	if err != nil {
		log.Fatal(err)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}
