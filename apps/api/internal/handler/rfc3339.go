package handler

import "time"

const rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

func formatRFC3339Milli(t time.Time) string {
	return t.UTC().Format(rfc3339Milli)
}
