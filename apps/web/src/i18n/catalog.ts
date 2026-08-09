/**
 * Translation catalog (S1 · C2/C3).
 *
 * Catalogs are pure data files under `messages/`; `en-US` is the canonical
 * baseline. Resolution order for a key in locale L:
 *
 *   catalog[L] → catalog[en-US] → observable missing-key event → key itself
 *
 * A key is "missing" only when neither the current catalog nor the en-US
 * catalog has it; the en-US fallback is silent (designed behavior). Missing
 * keys never render empty, never throw, and never block the flow; they are
 * observable via the `schema-ui:missing-translation` window event (deduped
 * per locale+key, so the first occurrence always reports).
 */

import enUS from "./messages/en-US.json";
import zhCN from "./messages/zh-CN.json";
import { DEFAULT_LOCALE, type Locale } from "./locale";

export type MessageParams = Record<string, string | number>;

export interface MissingTranslationDetail {
  locale: Locale;
  key: string;
  /** Optional rendering context (e.g. "nav.sidebar", "page.users.form"). */
  path?: string;
}

export const MISSING_TRANSLATION_EVENT = "schema-ui:missing-translation";

const catalogs: Record<Locale, Record<string, string>> = {
  "en-US": enUS as unknown as Record<string, string>,
  "zh-CN": zhCN as unknown as Record<string, string>,
};

const reportedMissing = new Set<string>();

/** True when the key exists in the given locale catalog. */
export function hasTranslation(key: string, locale: Locale): boolean {
  return Object.prototype.hasOwnProperty.call(catalogs[locale], key);
}

/** Raw catalog text for a key, or null when the locale catalog lacks it. */
export function lookupTranslation(key: string, locale: Locale): string | null {
  if (hasTranslation(key, locale)) {
    return catalogs[locale][key];
  }
  return null;
}

/** Replaces `{name}` placeholders with params; unknown placeholders stay. */
export function interpolate(template: string, params?: MessageParams): string {
  if (params === undefined) {
    return template;
  }
  return template.replace(/\{([a-zA-Z0-9_]+)\}/g, (match, name: string) =>
    Object.prototype.hasOwnProperty.call(params, name) ? String(params[name]) : match,
  );
}

/** Publishes a deduped missing-key report to the window event bus. */
export function reportMissingTranslation(detail: MissingTranslationDetail): void {
  const dedupeKey = `${detail.locale}:${detail.key}`;
  if (reportedMissing.has(dedupeKey)) {
    return;
  }
  reportedMissing.add(dedupeKey);
  if (typeof window !== "undefined") {
    window.dispatchEvent(
      new CustomEvent<MissingTranslationDetail>(MISSING_TRANSLATION_EVENT, { detail }),
    );
  }
}

/** Resets the missing-key dedupe set (test seam). */
export function resetMissingTranslationReports(): void {
  reportedMissing.clear();
}

/**
 * Resolves a message key for a locale with the frozen fallback chain.
 * Never throws, never returns an empty string for a missing key.
 *
 * Fallback order: catalog[locale] → catalog[en-US] → `literalFallback`
 * (protocol literal text, when supplied) → key itself.
 */
export function translate(
  key: string,
  params?: MessageParams,
  locale: Locale = DEFAULT_LOCALE,
  path?: string,
  literalFallback?: string,
): string {
  const direct = lookupTranslation(key, locale);
  if (direct !== null) {
    return interpolate(direct, params);
  }
  const fallback = lookupTranslation(key, DEFAULT_LOCALE);
  if (fallback !== null) {
    return interpolate(fallback, params);
  }
  reportMissingTranslation({ locale, key, path });
  return literalFallback !== undefined && literalFallback !== "" ? literalFallback : key;
}

/** Binds a locale (+ optional context path) to a translate function. */
export function createTranslator(
  locale: Locale,
  options?: { path?: string },
): (key: string, params?: MessageParams, literalFallback?: string) => string {
  return (key: string, params?: MessageParams, literalFallback?: string) =>
    translate(key, params, locale, options?.path, literalFallback);
}

/**
 * Resolves a schema/manifest text prop pair — the `*Key` field wins over the
 * literal protocol text, and the literal text is the last fallback before the
 * key itself (frozen chain: 当前语种 → en-US → 字面文本 → key).
 */
export function resolveTextProp(
  props: Record<string, unknown> | undefined,
  keyProp: string,
  literalProp: string,
  t: (key: string, params?: MessageParams, literalFallback?: string) => string,
  fallback = "",
): string {
  if (props === undefined) {
    return fallback;
  }
  const key = props[keyProp];
  if (typeof key === "string" && key !== "") {
    const literal = typeof props[literalProp] === "string" ? props[literalProp] : undefined;
    return t(key, undefined, literal);
  }
  const literal = props[literalProp];
  return typeof literal === "string" ? literal : fallback;
}
