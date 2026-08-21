/**
 * Locale-aware date/number formatting (S1 · C5).
 *
 * Formatting follows the effective locale through Intl.* — no custom format
 * templates (VP-007: "首版不暴露任意日期/数字格式模板，随有效 locale").
 * Formatting is fail-safe: invalid inputs render empty, invalid timezones
 * degrade to the locale's default zone instead of throwing.
 */

import { DEFAULT_LOCALE, type Locale } from "./locale";

export interface FormatOptions {
  /** IANA timezone name; omitted = the environment's default zone. */
  timeZone?: string;
}

/** Formats a date value in the given locale. Returns "" for invalid input. */
export function formatDate(
  value: Date | string | number,
  locale: Locale = DEFAULT_LOCALE,
  options: FormatOptions = {},
): string {
  const date = value instanceof Date ? value : new Date(value);
  if (!Number.isFinite(date.getTime())) {
    return "";
  }
  const timeZone =
    options.timeZone !== undefined && options.timeZone !== "" ? options.timeZone : undefined;
  try {
    return new Intl.DateTimeFormat(locale, {
      dateStyle: "medium",
      timeStyle: "short",
      ...(timeZone === undefined ? {} : { timeZone }),
    }).format(date);
  } catch {
    // Invalid IANA name — degrade to the default zone, never throw.
    return new Intl.DateTimeFormat(locale, { dateStyle: "medium", timeStyle: "short" }).format(date);
  }
}

/** Formats a finite number in the given locale. Returns "" for invalid input. */
export function formatNumber(
  value: number,
  locale: Locale = DEFAULT_LOCALE,
  options: Intl.NumberFormatOptions = {},
): string {
  if (typeof value !== "number" || !Number.isFinite(value)) {
    return "";
  }
  return new Intl.NumberFormat(locale, options).format(value);
}
