package util

import (
	"database/sql"
	"time"
)

func SqlNullString(str string) sql.NullString {
	if str == "" {
		return sql.NullString{}
	}
	return sql.NullString{Valid: true, String: str}
}

func SqlNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Valid: true, Time: t}
}

func SqlNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Valid: true, Int64: i}
}

func SqlNullFloat64(f float64) sql.NullFloat64 {
	if f == 0.0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Valid: true, Float64: f}
}

func SqlNullInt16(i int) sql.NullInt16 {
	if i == 0 {
		return sql.NullInt16{}
	}
	return sql.NullInt16{Valid: true, Int16: int16(i)}
}

func SqlNullInt32(i int) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Valid: true, Int32: int32(i)}
}
