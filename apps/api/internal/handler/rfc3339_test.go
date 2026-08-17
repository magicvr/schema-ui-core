package handler

import (
	"testing"
	"time"
)

func TestFormatRFC3339Milli(t *testing.T) {
	got := formatRFC3339Milli(time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC))
	if got != "2026-08-17T12:00:00.000Z" {
		t.Fatalf("got %q", got)
	}
}
