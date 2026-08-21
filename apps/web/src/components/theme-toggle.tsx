/**
 * Theme toggle (S1 · C3) — W8 visual unification (GOAL-009):
 *
 * The trigger now matches the neighbouring header icon buttons (language
 * switcher and notification bell): size-9 rounded-md ghost with
 * text-muted-foreground + hover:bg-accent — instead of the form-styled
 * outline Button. The tooltip / aria-label follow the active locale via
 * the i18n catalog (shell.theme.toggle) — no hardcoded English.
 */

import { Moon, Sun } from "lucide-react";

import { useTranslate } from "@/i18n/runtime";
import { setTheme } from "@/theme/theme";

function isDark() {
  return document.documentElement.classList.contains("dark");
}

export function ThemeToggle() {
  const t = useTranslate();
  const label = t("shell.theme.toggle");
  return (
    <button
      type="button"
      aria-label={label}
      title={label}
      className="inline-flex size-9 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
      onClick={() => {
        // setTheme persists the choice and applies it via applyThemeToElement,
        // which also syncs CSS color-scheme so native controls match (D8).
        setTheme(isDark() ? "light" : "dark");
      }}
    >
      <Sun aria-hidden="true" className="size-4 dark:hidden" />
      <Moon aria-hidden="true" className="hidden size-4 dark:block" />
    </button>
  );
}
