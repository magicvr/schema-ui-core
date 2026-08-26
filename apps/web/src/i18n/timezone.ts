/**
 * Timezone pure-logic unit (workspace-020 · R2 · C1).
 *
 * `resolveEffectiveTimezone` computes the effective timezone from the user's
 * explicit override, the session probe, and the site default, per contract
 * GOAL-002 D-001 §2 (user-confirmed I-001, 2026-08-26):
 *
 *   L1 user override (localStorage "schema-ui:timezone") → L2 session probe
 *   (Intl) → L3 site default (siteTimezone) → L4 "auto" fallback
 *
 * Keeping the decision logic in plain functions lets vitest exercise every
 * branch; the probe is injectable so tests do not depend on the host zone.
 */

export const TIMEZONE_STORAGE_KEY = "schema-ui:timezone";

export const AUTO_TIMEZONE = "auto";

/** The user-facing choice; a timezone IANA name or "auto". */
export type TimezonePreference = string | "auto";

export interface TimezoneResolutionInput {
  /**
   * localStorage["schema-ui:timezone"] — the user's explicit override.
   * null / undefined / invalid / "auto" → no override (skip to L2).
   */
  stored: string | null;
  /**
   * Site default from the public bootstrap (/api/branding siteTimezone).
   * "auto" / "" / null → unset (skip to L4 when L2 is empty).
   */
  siteDefault: string | null;
  /**
   * Session probe of the host zone (injectable; real default is
   * `detectBrowserTimezone`). Invalid or empty results are skipped.
   */
  detect: () => string;
}

/**
 * Validates an IANA timezone name via Intl (RangeError on invalid names).
 * Side-effect-free apart from the noexcept Intl probe; returns false for
 * empty / non-string / unknown zones instead of throwing.
 */
export function isValidIanaTimeZone(raw: string | null | undefined): raw is string {
  if (typeof raw !== "string") {
    return false;
  }
  const trimmed = raw.trim();
  if (trimmed === "") {
    return false;
  }
  try {
    // The format call itself never renders; constructing with the option is
    // the validation probe (RangeError for unknown "Foo/Bar" style names).
    new Intl.DateTimeFormat("en-US", { timeZone: trimmed });
    return true;
  } catch {
    return false;
  }
}

/**
 * Normalizes a raw stored/site value into a TimezonePreference.
 * Any value that is not a valid IANA name resolves to "auto".
 */
export function normalizeTimezonePreference(raw: string | null | undefined): TimezonePreference {
  if (typeof raw !== "string") {
    return AUTO_TIMEZONE;
  }
  const trimmed = raw.trim();
  if (trimmed === "" || trimmed.toLowerCase() === AUTO_TIMEZONE) {
    return AUTO_TIMEZONE;
  }
  return isValidIanaTimeZone(trimmed) ? trimmed : AUTO_TIMEZONE;
}

/** Reads the stored user override; best-effort (privacy mode → null). */
export function readStoredTimezone(): string | null {
  try {
    if (typeof localStorage === "undefined") {
      return null;
    }
    return localStorage.getItem(TIMEZONE_STORAGE_KEY);
  } catch {
    return null;
  }
}

/** Persists the user override; "auto" removes the key (single channel). */
export function writeStoredTimezone(preference: TimezonePreference): void {
  try {
    if (typeof localStorage === "undefined") {
      return;
    }
    if (preference === AUTO_TIMEZONE) {
      localStorage.removeItem(TIMEZONE_STORAGE_KEY);
    } else {
      localStorage.setItem(TIMEZONE_STORAGE_KEY, preference);
    }
  } catch {
    // Best-effort persist; the in-memory preference still applies this page.
  }
}

/** Real session probe: the host zone from Intl.resolvedOptions(). */
export function detectBrowserTimezone(): string {
  try {
    const zone = new Intl.DateTimeFormat().resolvedOptions().timeZone;
    return typeof zone === "string" ? zone : "";
  } catch {
    return "";
  }
}

/**
 * Resolves the effective timezone per contract §2 (L1 → L2 → L3 → L4).
 * Returns an IANA name, or "auto" when nothing is configured/detectable —
 * consumers then fall back to the locale's default zone.
 */
export function resolveEffectiveTimezone(input: TimezoneResolutionInput): string {
  const stored = normalizeTimezonePreference(input.stored);
  if (stored !== AUTO_TIMEZONE) {
    return stored; // L1 · user override wins
  }
  let detected = "";
  try {
    detected = String(input.detect() ?? "");
  } catch {
    detected = "";
  }
  if (isValidIanaTimeZone(detected)) {
    return detected; // L2 · session probe
  }
  const site = normalizeTimezonePreference(input.siteDefault);
  if (site !== AUTO_TIMEZONE) {
    return site; // L3 · site default
  }
  return AUTO_TIMEZONE; // L4 · embedded default
}