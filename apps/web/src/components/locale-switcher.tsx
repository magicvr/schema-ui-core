/**
 * Language switcher (S1 · C4).
 *
 * Reachable from the Shell and the anonymous login page — no settings
 * permission required (VP-007 requirement). Uses the native <select> for
 * accessibility; the visible labels come from the translation catalog.
 */

import { useI18n } from "@/i18n/runtime";
import { SUPPORTED_LOCALES, type LocalePreference } from "@/i18n/locale";

export interface LocaleSwitcherProps {
  className?: string;
}

export function LocaleSwitcher({ className = "" }: LocaleSwitcherProps) {
  const { preference, setPreference, t } = useI18n();
  return (
    <label className={className}>
      <span className="sr-only">{t("locale.switcher.label")}</span>
      <select
        aria-label={t("locale.switcher.label")}
        title={t("locale.switcher.label")}
        value={preference}
        onChange={(event) => setPreference(event.target.value as LocalePreference)}
        className="h-9 rounded-md border border-border bg-background px-2 text-xs text-muted-foreground outline-none transition-colors hover:bg-accent/60 focus-visible:ring-2 focus-visible:ring-ring"
      >
        <option value="auto">{t("locale.switcher.auto")}</option>
        {SUPPORTED_LOCALES.map((locale) => (
          <option key={locale} value={locale}>
            {t(`locale.name.${locale}`)}
          </option>
        ))}
      </select>
    </label>
  );
}
