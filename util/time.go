package util

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/royanqodri/fms-ai-hackathon/config"
)

type TimeFormat string
type CustomDate time.Time
type CustomDateTime time.Time

var LOCATION = "Asia/Makassar"
var LOCATION_TIME *time.Location

const (
	DATE                      string = "2006-01-02"
	DATETIME                  string = "2006-01-02 15:04:05"
	DATETIME_MS               string = "2006-01-02 15:04:05.000"
	DATETIME_MINUTES_COMBINED string = "200601021504"
	DATE_COMBINED             string = "20060102"
	TIME                      string = "15:04:05"
)

func InitLocation() {
	if config.Get().Service.Location != "" {
		LOCATION = config.Get().Service.Location
		LOCATION_TIME, _ = time.LoadLocation(LOCATION)
	}
}

// ToTime converts util.TimeFormat to time.Time
func (tf TimeFormat) ToTime() time.Time {
	return tf.ToTime()
}

func GetTimeNowByLoc() time.Time {
	return time.Now().In(LOCATION_TIME)
}

func GetTimeNowByLocWithoutMs() time.Time {
	return time.Now().In(LOCATION_TIME).Truncate(time.Second)
}

func GetDateNowByLoc() time.Time {
	t := GetTimeNowByLoc()
	year, month, day := t.Date()
	loc := t.Location()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func GetTimeNowMillisInStr() string {
	now := GetTimeNowByLoc().UnixMilli()
	return IntToStr(now)
}

func GetFormattedDateTimeInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATETIME)
}

func GetFormattedDateTimeMsInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATETIME_MS)
}

func GetFormattedTimeInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TIME)
}

func GetFormattedDateInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATE)
}

func GetFormattedDateTimeMinutesCombinedInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATETIME_MINUTES_COMBINED)
}

func GetFormattedDateCombinedInStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATE_COMBINED)
}

func GetFormattedDate(t string) time.Time {
	parsedTime, err := time.ParseInLocation(DATE, t, LOCATION_TIME)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func GetFormattedDateTime(t string) time.Time {
	parsedTime, err := time.ParseInLocation(DATETIME, t, LOCATION_TIME)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func GetFormattedDateTimeNullable(t *string) time.Time {
	if t == nil || *t == "" {
		return time.Time{}
	}

	parsedTime, err := time.ParseInLocation(DATETIME, *t, LOCATION_TIME)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func GetFormattedTime(t string) time.Time {
	parsedTime, err := time.ParseInLocation(TIME, t, LOCATION_TIME)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func ParseAndFormatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}
	return t.Format("2006-01-02")
}

func (ct *CustomDateTime) UnmarshalText(data []byte) error {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", string(data))
	if err != nil {
		return err
	}
	*ct = CustomDateTime(parsedTime)
	return nil
}

func (ct CustomDateTime) ToNullTime() sql.NullTime {
	// Assuming that if the time is zero, it's considered invalid
	t := time.Time(ct)
	valid := !t.IsZero()

	return sql.NullTime{
		Time:  t,
		Valid: valid,
	}
}

func (ct *CustomDate) UnmarshalText(data []byte) error {
	parsedTime, err := time.Parse("2006-01-02", string(data))
	if err != nil {
		return err
	}
	*ct = CustomDate(parsedTime)
	return nil
}

func GetDateNowStrID() map[string]string {
	now := GetTimeNowByLoc()

	// Format date
	date := now.Format("2006-01-02 15:04:05.000")
	id := now.Format("20060102150405000")

	// Create map
	objReturn := map[string]string{
		"date": date,
		"id":   id,
	}

	return objReturn
}
