/**
 * I18n React runtime (S1 · C4/C5).
 *
 * - Resolves the effective locale via `resolveLocale` (frozen priority).
 * - Persists the user's explicit choice in localStorage["schema-ui:locale"]
 *   (single channel, same pattern as the theme mechanism; login/logout never
 *   clears it — D-002 §I-L10N-002).
 * - Applies `document.documentElement.lang` on locale change.
 * - Exposes `t` / `formatDate` / `formatNumber` to components.
 */

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";

import { createTranslator, type MessageParams } from "./catalog";
import { formatDate as formatDateImpl, formatNumber as formatNumberImpl } from "./format";
import {
  defaultBrowserLanguages,
  DEFAULT_LOCALE,
  normalizePreference,
  resolveLocale,
  type Locale,
  type LocalePreference,
} from "./locale";
import {
  AUTO_TIMEZONE,
  detectBrowserTimezone,
  normalizeTimezonePreference,
  readStoredTimezone,
  resolveEffectiveTimezone,
  writeStoredTimezone,
  type TimezonePreference,
} from "./timezone";
export const LOCALE_STORAGE_KEY = "schema-ui:locale";

export function readStoredLocale(): string | null {
  try {
    if (typeof localStorage === "undefined") {
      return null;
    }
    return localStorage.getItem(LOCALE_STORAGE_KEY);
  } catch {
    // Privacy mode / disabled storage throws SecurityError — same posture as
    // tokens.ts / theme.ts. Never let locale boot white-screen the tree.
    return null;
  }
}

export function writeStoredLocale(preference: LocalePreference): void {
  try {
    if (typeof localStorage === "undefined") {
      return;
    }
    if (preference === "auto") {
      localStorage.removeItem(LOCALE_STORAGE_KEY);
    } else {
      localStorage.setItem(LOCALE_STORAGE_KEY, preference);
    }
  } catch {
    // Best-effort persist; the in-memory preference still applies this page.
  }
}

/** Applies the effective locale to <html lang>. No-op outside a browser. */
export function applyLocaleToDocument(locale: Locale): void {
  if (typeof document !== "undefined") {
    document.documentElement.lang = locale;
  }
}

// ── active-locale registry ────────────────────────────────────────────────────
// Module-level current locale for non-React consumers (API fetchers attach it
// as Accept-Language so the server negotiates the same language — VP-007 S4).

let activeLocale: Locale = "en-US";

/** Returns the provider's currently effective locale (defaults en-US). */
export function getActiveLocale(): Locale {
  return activeLocale;
}

/** Internal: keeps the registry in sync with the provider's effective locale. */
export function setActiveLocale(locale: Locale): void {
  activeLocale = locale;
}

export interface I18nState {
  /** Effective (resolved) locale — always a supported locale. */
  locale: Locale;
  /** User preference; "auto" defers to system/browser defaults. */
  preference: LocalePreference;
  /** Sets the user preference and persists it (localStorage single channel). */
  setPreference: (preference: LocalePreference) => void;
  /** Effective timezone (IANA name or "auto" per contract §2 / L1–L4). */
  timezone: string;
  /** User timezone preference; "auto" defers to session/site defaults. */
  timezonePreference: TimezonePreference;
  /** Sets the user timezone preference and persists it (single channel). */
  setTimezonePreference: (preference: TimezonePreference) => void;
  /** Translate a catalog key with the effective locale. */
  t: (key: string, params?: MessageParams) => string;
  /** Locale-aware date formatting (defaults to the effective timezone). */
  formatDate: (value: Date | string | number, options?: { timeZone?: string }) => string;
  /** Locale-aware number formatting. */
  formatNumber: (value: number, options?: Intl.NumberFormatOptions) => string;
}

const I18nContext = createContext<I18nState | null>(null);

export interface I18nProviderProps {
  children: ReactNode;
  /** System default locale from the public bootstrap; null/"auto" = none. */
  systemDefault?: string | null;
  /** Test seam: explicit stored preference (defaults to localStorage). */
  stored?: string | null;
  /** Test seam: browser language list (defaults to navigator.languages). */
  browserLanguages?: readonly string[];
  /**
   * Site default timezone from /api/branding (contract §2 · L3). The provider
   * also captures it from `systemDefaultUrl` when present; this prop wins.
   */
  siteTimezone?: string | null;
  /** Test seam: explicit stored timezone preference (defaults to localStorage). */
  storedTimezone?: string | null;
  /** Test seam: session timezone probe (defaults to detectBrowserTimezone). */
  detectTimezone?: () => string;
  /**
   * When set, the provider fetches this public startup endpoint once and
   * re-resolves the system default locale from `defaultLocale` (VP-007 S3:
   * the shell/login apply the site-wide default when the user has no
   * explicit choice) and the site default timezone from `siteTimezone`
   * (workspace-020 · contract §2 · L3). Applies after the initial resolve.
   */
  systemDefaultUrl?: string;
}

export function I18nProvider({
  children,
  systemDefault = null,
  stored,
  browserLanguages,
  siteTimezone = null,
  storedTimezone,
  detectTimezone,
  systemDefaultUrl,
}: I18nProviderProps) {
  const [preference, setPreferenceState] = useState<LocalePreference>(() => {
    const raw = stored !== undefined ? stored : readStoredLocale();
    return normalizePreference(raw);
  });
  const [fetchedSystemDefault, setFetchedSystemDefault] = useState<string | null>(null);
  const [fetchedSiteTimezone, setFetchedSiteTimezone] = useState<string | null>(null);
  const [timezonePreference, setTimezonePreferenceState] = useState<TimezonePreference>(() => {
    const raw = storedTimezone !== undefined ? storedTimezone : readStoredTimezone();
    return normalizeTimezonePreference(raw);
  });

  const browserList = browserLanguages !== undefined ? browserLanguages : defaultBrowserLanguages();

  useEffect(() => {
    if (systemDefaultUrl === undefined) {
      return;
    }
    let cancelled = false;
    fetch(systemDefaultUrl)
      .then((response) => (response.ok ? response.json() : null))
      .then((body: unknown) => {
        if (cancelled) {
          return;
        }
        const record = body as Record<string, unknown> | null;
        const locale = typeof record?.defaultLocale === "string" ? record.defaultLocale : null;
        setFetchedSystemDefault(locale === "auto" ? null : locale);
        const zone = typeof record?.siteTimezone === "string" ? record.siteTimezone : null;
        setFetchedSiteTimezone(zone === "auto" ? null : zone);
      })
      .catch(() => {
        if (!cancelled) {
          setFetchedSystemDefault(null);
          setFetchedSiteTimezone(null);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [systemDefaultUrl]);

  const effectiveSystemDefault = systemDefault ?? fetchedSystemDefault;
  const effectiveSiteTimezone = siteTimezone ?? fetchedSiteTimezone;

  const locale = useMemo<Locale>(
    () =>
      resolveLocale({
        stored: preference === "auto" ? null : preference,
        systemDefault: effectiveSystemDefault,
        browserLanguages: browserList,
      }),
    [preference, effectiveSystemDefault, browserList],
  );

  const timezone = useMemo(
    () =>
      resolveEffectiveTimezone({
        stored: timezonePreference === AUTO_TIMEZONE ? null : timezonePreference,
        siteDefault: effectiveSiteTimezone,
        detect: detectTimezone ?? detectBrowserTimezone,
      }),
    [timezonePreference, effectiveSiteTimezone, detectTimezone],
  );

  useEffect(() => {
    applyLocaleToDocument(locale);
    setActiveLocale(locale);
  }, [locale]);

  const setPreference = useCallback((next: LocalePreference) => {
    writeStoredLocale(next);
    setPreferenceState(next);
  }, []);

  const setTimezonePreference = useCallback((next: TimezonePreference) => {
    writeStoredTimezone(next);
    setTimezonePreferenceState(next);
  }, []);

  const t = useMemo(() => createTranslator(locale), [locale]);

  const formatDate = useCallback(
    (value: Date | string | number, options: { timeZone?: string } = {}) => {
      const requested =
        options.timeZone !== undefined && options.timeZone !== "" ? options.timeZone : null;
      const zone = requested ?? (timezone === AUTO_TIMEZONE ? null : timezone);
      return formatDateImpl(value, locale, zone === null ? {} : { timeZone: zone });
    },
    [locale, timezone],
  );

  const formatNumber = useCallback(
    (value: number, options?: Intl.NumberFormatOptions) =>
      formatNumberImpl(value, locale, options),
    [locale],
  );

  const value = useMemo<I18nState>(
    () => ({
      locale,
      preference,
      setPreference,
      timezone,
      timezonePreference,
      setTimezonePreference,
      t,
      formatDate,
      formatNumber,
    }),
    [
      locale,
      preference,
      setPreference,
      timezone,
      timezonePreference,
      setTimezonePreference,
      t,
      formatDate,
      formatNumber,
    ],
  );

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n(): I18nState {
  const value = useContext(I18nContext);
  if (value === null) {
    throw new Error("useI18n must be used within an I18nProvider");
  }
  return value;
}

/**
 * Tolerant translator hook for deep renderer internals.
 *
 * Returns the provider's translator, or a safe default (en-US resolution +
 * missing-key observable fallback) when no provider is mounted. Production
 * always mounts I18nProvider; bare component tests and pre-provider surfaces
 * degrade to the documented safe fallback instead of throwing.
 */
export function useTranslate(): (key: string, params?: MessageParams) => string {
  const value = useContext(I18nContext);
  return useMemo(() => value?.t ?? createTranslator(DEFAULT_LOCALE), [value]);
}
