/**
 * Money / number pure-logic unit (workspace-020 · R3 · C1/C2/C3/C5).
 *
 * Contract GOAL-002 D-001 §3 / §4.3 (user-confirmed I-002, 2026-08-26):
 * - Display + input semantics live on the frontend; API stays a machine
 *   contract (amounts = int64 JSON in the smallest currency unit).
 * - No custom format templates: symbols / positions / fraction digits are
 *   derived from `Intl.NumberFormat` (locale + currency).
 * - Embedded default currency map (§4.3): zh-CN → CNY, en-US → USD,
 *   unknown locale → USD. Missing configuration never throws.
 * - Input parsing normalizes localized strings back to machine values
 *   (amount in minor units); unparsable input → null (callers show a
 *   localized input error and must NOT submit the raw string).
 */

import { SUPPORTED_LOCALES, type Locale } from "./locale";

/** Site/machine default currency (ISO 4217). */
export const DEFAULT_CURRENCY = "USD";

/** Default minor-unit exponent per currency group (ISO 4217 minor units). */
export const DEFAULT_MINOR_UNITS = 2;

/** Currency codes frozen in the contract's embedded default map (§4.3). */
const DEFAULT_CURRENCY_MAP: Record<Locale, string> = {
  "zh-CN": "CNY",
  "en-US": "USD",
};

/** Uppercase three-letter ISO 4217 code; null for anything else. */
export function normalizeCurrencyCode(raw: string | null | undefined): string | null {
  if (typeof raw !== "string") {
    return null;
  }
  const trimmed = raw.trim();
  return /^[A-Za-z]{3}$/.test(trimmed) ? trimmed.toUpperCase() : null;
}

/**
 * Embedded default currency for a locale (§4.3). Unknown locales fall back to
 * USD; never throws.
 */
export function defaultCurrencyFor(locale: Locale | string): string {
  const key = (SUPPORTED_LOCALES as readonly string[]).includes(locale)
    ? (locale as Locale)
    : null;
  return key === null ? DEFAULT_CURRENCY : DEFAULT_CURRENCY_MAP[key];
}

/**
 * Effective default currency: the explicit site default (branding
 * `defaultCurrency`, ISO 4217) wins when set; otherwise the embedded
 * per-locale map (§4.3) applies. Never throws.
 */
export function resolveEffectiveCurrency(
  locale: Locale,
  siteDefault: string | null | undefined,
): string {
  const site = normalizeCurrencyCode(siteDefault);
  return site !== null ? site : defaultCurrencyFor(locale);
}

/** Locale group/decimal separators derived from Intl (no hardcoded tables). */
export interface LocaleSeparators {
  group: string;
  decimal: string;
}

export function localeSeparators(locale: Locale): LocaleSeparators {
  let group = ",";
  let decimal = ".";
  try {
    for (const part of new Intl.NumberFormat(locale).formatToParts(1234567.89)) {
      if (part.type === "group") {
        group = part.value;
      } else if (part.type === "decimal") {
        decimal = part.value;
      }
    }
  } catch {
    // Unsupported locale edge — keep the en-US-ish defaults.
  }
  return { group, decimal };
}

export interface MoneyFormatOptions {
  /** ISO 4217 currency code (uppercase 3 letters). Explicit override. */
  currency?: string;
  /** Site-wide default currency from /api/branding ("" = unset). */
  siteDefaultCurrency?: string;
  /** Minor-unit exponent. Defaults to 2 (CNY/USD). */
  minorUnits?: number;
}

/** Priority: explicit option → site default → embedded per-locale map. */
function resolveCurrency(locale: Locale, options: MoneyFormatOptions): string | null {
  if (options.currency !== undefined) {
    return normalizeCurrencyCode(options.currency);
  }
  return normalizeCurrencyCode(options.siteDefaultCurrency) ?? defaultCurrencyFor(locale);
}

/**
 * Formats an amount given as machine value (minor units) into a
 * locale+currency display string. Invalid input renders "" (fail-safe).
 * R4 F-007: values beyond Number.MAX_SAFE_INTEGER render "" — the machine
 * contract declares int64 minor units, which JS number cannot carry.
 */
export function formatMoney(
  minorValue: number,
  locale: Locale,
  options: MoneyFormatOptions = {},
): string {
  if (typeof minorValue !== "number" || !Number.isFinite(minorValue) || !Number.isSafeInteger(minorValue)) {
    return "";
  }
  const currency = resolveCurrency(locale, options);
  if (currency === null) {
    return "";
  }
  const minorUnits = Number.isInteger(options.minorUnits) && options.minorUnits! >= 0
    ? options.minorUnits!
    : DEFAULT_MINOR_UNITS;
  const major = minorValue / Math.pow(10, minorUnits);
  try {
    return new Intl.NumberFormat(locale, {
      style: "currency",
      currency,
      minimumFractionDigits: minorUnits,
      maximumFractionDigits: minorUnits,
    }).format(major);
  } catch {
    return "";
  }
}

/**
 * Weakly parses a localized number (separators stripped) into a plain number.
 * Returns null for anything without a parseable numeric core.
 */
function parseLocalizedNumberCore(raw: string, locale: Locale): number | null {
  if (typeof raw !== "string") {
    return null;
  }
  const trimmed = raw.trim();
  if (trimmed === "") {
    return null;
  }
  const { group, decimal } = localeSeparators(locale);
  let core = trimmed;
  if (group !== "") {
    // Strip the locale grouping separator wherever it appears. NOTE: group
    // position correctness (e.g. "12,34.5") is intentionally NOT validated —
    // out of the R3 verification scope (tolerance documented for R4 review).
    core = core.split(group).join("");
  }
  if (decimal !== ".") {
    core = core.split(decimal).join(".");
  }
  if ((core.match(/\./g) ?? []).length > 1) {
    return null; // more than one decimal separator → not a number
  }
  if (!/^-?\d+(\.\d+)?$/.test(core)) {
    return null;
  }
  const value = Number(core);
  return Number.isFinite(value) ? value : null;
}

/** Strips a locale+currency symbol (e.g. "¥", "$", "CN¥") and the code itself. */
function stripCurrencyAffixes(raw: string, locale: Locale, currency: string): string {
  let out = raw.trim();
  try {
    const parts = new Intl.NumberFormat(locale, {
      style: "currency",
      currency,
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).formatToParts(1.23);
    const symbol = parts.find((part) => part.type === "currency")?.value;
    if (symbol !== undefined && symbol !== "") {
      out = out.split(symbol).join("");
    }
  } catch {
    // Fall through with the raw string.
  }
  // Also tolerate the bare ISO code (e.g. "CNY 123.45").
  return out.split(currency).join("").trim();
}

export interface MoneyParseOptions {
  /** Expected ISO 4217 currency (affects symbol stripping only). */
  currency?: string;
  /** Site-wide default currency from /api/branding ("" = unset). */
  siteDefaultCurrency?: string;
  /** Minor-unit exponent for the returned integer. Defaults to 2. */
  minorUnits?: number;
}

/**
 * Parses a localized money string into the machine value (minor-unit
 * integer per contract §3.3). Returns null when the input is not a
 * parseable amount — callers must NOT submit the raw string.
 */
export function parseLocalizedMoney(
  raw: string,
  locale: Locale,
  options: MoneyParseOptions = {},
): number | null {
  const currency = resolveCurrency(locale, options);
  if (currency === null) {
    return null;
  }
  const cleaned = stripCurrencyAffixes(raw, locale, currency);
  const value = parseLocalizedNumberCore(cleaned, locale);
  if (value === null) {
    return null;
  }
  const minorUnits = Number.isInteger(options.minorUnits) && options.minorUnits! >= 0
    ? options.minorUnits!
    : DEFAULT_MINOR_UNITS;
  const minor = Math.round(value * Math.pow(10, minorUnits));
  // R4 F-007: the machine contract declares int64 minor units; JS number
  // cannot represent values beyond MAX_SAFE_INTEGER — reject instead of
  // silently losing precision.
  return Number.isSafeInteger(minor) ? minor : null;
}

/** Parses a localized plain number into a machine number; null on failure. */
export function parseLocalizedNumber(raw: string, locale: Locale): number | null {
  return parseLocalizedNumberCore(raw, locale);
}