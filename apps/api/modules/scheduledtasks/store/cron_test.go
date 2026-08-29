package store

import (
	"testing"
	"time"
)

func TestParseCron(t *testing.T) {
	for _, expr := range []string{
		"* * * * *", "*/5 * * * *", "0 9 * * 1-5", "30 8,18 * * *", "0 0 1 1 *",
		"15 */2 * * *", "0 12 * * 0",
	} {
		if _, err := ParseCron(expr); err != nil {
			t.Fatalf("ParseCron(%q) = %v, want nil", expr, err)
		}
	}
	for _, expr := range []string{
		"", "a b c d e", "* * * *", "* * * * * *", "60 * * * *", "* 24 * * *",
		"* * 0 * *", "* * * 13 *", "* * * * 7", "*/0 * * * *", "1, * * * *", "*/-1 * * * *",
		// A-003 F-001: a bare scalar with /step would silently drop the step.
		"0/5 * * * *", "5/2 * * * *", "0-59/0 * * * *",
	} {
		if _, err := ParseCron(expr); err == nil {
			t.Fatalf("ParseCron(%q) = nil, want error", expr)
		}
	}
}

func TestCronFieldsNext(t *testing.T) {
	// Every minute at minute 0 of every hour: next is the next hour boundary.
	fields, err := ParseCron("0 * * * *")
	if err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 8, 14, 10, 30, 0, 0, time.UTC)
	next, ok := fields.Next(base)
	if !ok || !next.Equal(time.Date(2026, 8, 14, 11, 0, 0, 0, time.UTC)) {
		t.Fatalf("Next = %v %v, want 11:00", next, ok)
	}

	// Every 5 minutes: next is within 5 minutes.
	fields5, err := ParseCron("*/5 * * * *")
	if err != nil {
		t.Fatal(err)
	}
	next5, ok := fields5.Next(base)
	if !ok || next5.Minute()%5 != 0 || next5.Before(base) {
		t.Fatalf("Next */5 = %v %v", next5, ok)
	}
	// At a non-matching minute the next slot is the following */5 boundary.
	next5b, ok := fields5.Next(time.Date(2026, 8, 14, 10, 31, 0, 0, time.UTC))
	if !ok || next5b.Minute() != 35 {
		t.Fatalf("Next */5 from 10:31 = %v %v, want 10:35", next5b, ok)
	}

	// Daily at 09:15 — inclusive: exactly 09:15 matches the same moment; the
	// scheduler deduplicates per minute slot.
	fieldsDaily, err := ParseCron("15 9 * * *")
	if err != nil {
		t.Fatal(err)
	}
	nextD, ok := fieldsDaily.Next(time.Date(2026, 8, 14, 9, 15, 0, 0, time.UTC))
	if !ok || !nextD.Equal(time.Date(2026, 8, 14, 9, 15, 0, 0, time.UTC)) {
		t.Fatalf("Next daily = %v %v, want 09:15 same moment", nextD, ok)
	}
	nextD2, ok := fieldsDaily.Next(time.Date(2026, 8, 14, 9, 14, 0, 0, time.UTC))
	if !ok || !nextD2.Equal(time.Date(2026, 8, 14, 9, 15, 0, 0, time.UTC)) {
		t.Fatalf("Next daily from 09:14 = %v %v, want 09:15", nextD2, ok)
	}

	// 2026-02-29 does not exist: Feb 29 next matches 2028-02-29.
	fieldsLeap, err := ParseCron("0 0 29 2 *")
	if err != nil {
		t.Fatal(err)
	}
	nextL, ok := fieldsLeap.Next(time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC))
	if !ok || nextL.Year() != 2028 || nextL.Month() != 2 || nextL.Day() != 29 {
		t.Fatalf("Next leap = %v %v", nextL, ok)
	}
}
