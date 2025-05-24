package dateutil

import "time"

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

// FormatDate formats time.Time to YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// FormatDateTime formats time.Time to YYYY-MM-DD HH:MM:SS
func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(date string) (time.Time, error) {
	return time.Parse(DateFormat, date)
}

// ParseDateTime parses a datetime string in YYYY-MM-DD HH:MM:SS format
func ParseDateTime(datetime string) (time.Time, error) {
	return time.Parse(DateTimeFormat, datetime)
}

// StartOfDay returns the start time of the given date (00:00:00)
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end time of the given date (23:59:59)
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}
