// 5-field cron support (S-04 · GOAL-010 D-002 `2): self-written validator and
// next-run computation — minute hour dom month dow, supporting *, numbers,
// step (*/n) and lists (a,b). No external dependency (D-001 `5).
package store

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// cronField bounds for the five fields.
var cronRanges = [5][2]int{
	{0, 59}, // minute
	{0, 23}, // hour
	{1, 31}, // dom
	{1, 12}, // month
	{0, 6},  // dow (0 = Sunday)
}

// CronFields is the parsed 5-field schedule.
type CronFields [5]map[int]bool

// ParseCron validates a 5-field cron expression and returns the parsed set.
func ParseCron(expr string) (CronFields, error) {
	parts := strings.Fields(strings.TrimSpace(expr))
	if len(parts) != 5 {
		return CronFields{}, fmt.Errorf("cron must have exactly 5 fields")
	}
	var fields CronFields
	for i, part := range parts {
		values, err := parseCronField(part, cronRanges[i][0], cronRanges[i][1])
		if err != nil {
			return CronFields{}, fmt.Errorf("field %d (%s): %w", i+1, part, err)
		}
		fields[i] = values
	}
	return fields, nil
}

func parseCronField(part string, min, max int) (map[int]bool, error) {
	values := map[int]bool{}
	for _, item := range strings.Split(part, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			return nil, fmt.Errorf("empty list item")
		}
		step := 1
		base := item
		if idx := strings.IndexByte(item, '/'); idx >= 0 {
			base = strings.TrimSpace(item[:idx])
			stepText := strings.TrimSpace(item[idx+1:])
			n, err := strconv.Atoi(stepText)
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("invalid step %q", stepText)
			}
			// A-003 F-001: step is only meaningful on * or a-b ranges (D-002 §2
			// lists */n only). A bare scalar with /step would silently drop the
			// step and match less often than the expression implies — reject it.
			if base != "*" && !strings.Contains(base, "-") {
				return nil, fmt.Errorf("step on a single value %q is not supported (use */n)", item)
			}
			step = n
		}
		if base == "*" {
			for v := min; v <= max; v += step {
				values[v] = true
			}
			continue
		}
		// Range syntax a-b (inclusive); the step applies within the range.
		if idx := strings.IndexByte(base, '-'); idx > 0 {
			loText := strings.TrimSpace(base[:idx])
			hiText := strings.TrimSpace(base[idx+1:])
			lo, errLo := strconv.Atoi(loText)
			hi, errHi := strconv.Atoi(hiText)
			if errLo != nil || errHi != nil || lo < min || hi > max || lo > hi {
				return nil, fmt.Errorf("invalid range %q", base)
			}
			for v := lo; v <= hi; v += step {
				values[v] = true
			}
			continue
		}
		n, err := strconv.Atoi(base)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q", base)
		}
		if n < min || n > max {
			return nil, fmt.Errorf("value %d out of range %d-%d", n, min, max)
		}
		values[n] = true
	}
	return values, nil
}

// Matches reports whether the fields match the given time.
func (f CronFields) Matches(t time.Time) bool {
	if !f[0][t.Minute()] || !f[1][t.Hour()] || !f[2][t.Day()] {
		return false
	}
	if !f[3][int(t.Month())] {
		return false
	}
	return f[4][int(t.Weekday())]
}

// Next returns the next matching time at or after from, searching up to
// five years ahead (a minute-resolution scan; an unbounded search is never
// attempted). ok=false when no match exists in the window. The inclusive
// start lets the scheduler execute the current minute slot exactly once per
// slot (deduplicated in memory — D-002 §3).
func (f CronFields) Next(from time.Time) (time.Time, bool) {
	start := from.Truncate(time.Minute)
	limit := start.AddDate(5, 0, 0)
	for t := start; t.Before(limit); t = t.Add(time.Minute) {
		if f.Matches(t) {
			return t, true
		}
	}
	return time.Time{}, false
}
