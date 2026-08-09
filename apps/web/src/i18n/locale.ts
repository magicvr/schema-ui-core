/**
 * Locale pure-logic unit (S1 · C1).
 *
 * `resolveLocale` is a side-effect-free function that computes the effective
 * locale from the user's explicit choice, the system default, and the browser
 * language preferences. Keeping the decision logic in a plain function lets
 * vitest exercise every branch without a browser.
 *
 * Frozen priority (VP-007 / D-002 §I-L10N-002, user-confirmed 2026-08-09):
 *
 *   user explicit choice → system default (non-auto) → browser preference
 *   (auto) → en-US safe fallback
 */

export const SUPPORTED_LOCALES = ["zh-CN", "en-US"] as const;

export type Locale = (typeof SUPPORTED_LOCALES)[number];

/** The user-facing choice; "auto" defers to system/browser preference. */
export type LocalePreference = Locale | "auto";

export const DEFAULT_LOCALE: Locale = "en-US";

export interface LocaleResolutionInput {
  /**
   * localStorage["schema-ui:locale"] — the user's explicit choice.
   * null / undefined / invalid → no explicit choice.
   */
  stored: string | null;
  /**
   * System default from the public bootstrap (/api/branding defaultLocale).
   * "auto" or null → no system default.
   */
  systemDefault: string | null;
  /** Browser language preferences in order (navigator.languages). */
  browserLanguages: readonly string[];
}

export function isSupportedLocale(raw: string | null | undefined): raw is Locale {
  return raw === "zh-CN" || raw === "en-US";
}

/**
 * Normalizes a BCP 47-ish candidate to a supported locale.
 * Accepts exact tags, case variants, underscore separators, and language-only
 * prefixes ("zh", "zh-cn", "zh_CN", "en-US", "en", "en-us", "en-GB" → en-US).
 * Returns null for anything else (including "auto").
 */
export function normalizeLocaleCandidate(raw: string | null | undefined): Locale | null {
  if (raw === null || raw === undefined) {
    return null;
  }
  const trimmed = raw.trim();
  if (trimmed === "") {
    return null;
  }
  const lower = trimmed.toLowerCase().replace(/_/g, "-");
  if (lower === "zh-cn" || lower === "zh") {
    return "zh-CN";
  }
  if (lower === "en-us" || lower === "en") {
    return "en-US";
  }
  if (lower.startsWith("en-")) {
    return "en-US";
  }
  return null;
}

/**
 * Resolves the effective locale using the frozen priority. Pure — no I/O.
 */
export function resolveLocale(input: LocaleResolutionInput): Locale {
  const explicit = normalizeLocaleCandidate(input.stored);
  if (explicit !== null) {
    return explicit;
  }
  const system = normalizeLocaleCandidate(input.systemDefault);
  if (system !== null) {
    return system;
  }
  for (const candidate of input.browserLanguages ?? []) {
    const normalized = normalizeLocaleCandidate(candidate);
    if (normalized !== null) {
      return normalized;
    }
  }
  return DEFAULT_LOCALE;
}

/**
 * Normalizes a raw stored value into a LocalePreference.
 * Any value that is not a supported locale resolves to "auto".
 */
export function normalizePreference(raw: string | null | undefined): LocalePreference {
  const normalized = normalizeLocaleCandidate(raw);
  return normalized === null ? "auto" : normalized;
}

/** Real browser inputs for the boot path (provider default). */
export function defaultBrowserLanguages(): readonly string[] {
  if (typeof navigator === "undefined" || typeof navigator.languages !== "object") {
    return [];
  }
  return navigator.languages as readonly string[];
}
