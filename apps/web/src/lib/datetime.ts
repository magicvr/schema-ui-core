/**
 * Display-time formatting (GOAL-011 table style): ISO-8601 timestamps are
 * rendered in a human-readable local-time form ("2026-08-01 18:02") instead
 * of the raw wire value ("2026-08-01T18:02:44.000Z").
 *
 * Returns null when the value is not a renderable timestamp so callers fall
 * back to the raw text (numbers, booleans, plain strings, arbitrary dates
 * with other shapes stay untouched).
 */

/** Matches the ISO-8601 timestamps shipped by the API (UTC Z or ±hh:mm). */
const ISO_TIMESTAMP_PATTERN =
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})$/;

function pad2(value: number): string {
  return String(value).padStart(2, "0");
}

/**
 * Formats an ISO-8601 timestamp string to local "YYYY-MM-DD HH:mm". Returns
 * null for non-timestamp values, malformed dates, or non-string input.
 */
export function formatDisplayTime(value: unknown): string | null {
  if (typeof value !== "string" || !ISO_TIMESTAMP_PATTERN.test(value)) {
    return null;
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return null;
  }
  return (
    date.getFullYear() +
    "-" +
    pad2(date.getMonth() + 1) +
    "-" +
    pad2(date.getDate()) +
    " " +
    pad2(date.getHours()) +
    ":" +
    pad2(date.getMinutes())
  );
}
