package util

import "time"

func FormatToDate(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}
