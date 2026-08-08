/**
 * Theme pure-logic unit (S1 · C3).
 *
 * `applyTheme` is a side-effect-free function that computes the desired theme
 * and then applies it to the provided document root.  Keeping the decision
 * logic in a plain function lets vitest exercise it without a browser.
 *
 * `initTheme` is the top-level boot call used by the synchronous inline
 * script in index.html and by main.tsx (no longer needed after inline
 * migration, kept for backwards-compat import path).
 */

export type Theme = "light" | "dark";

export interface ThemeInput {
  /** Value of `localStorage.getItem("theme")` — null when absent. */
  stored: string | null;
  /** `window.matchMedia("(prefers-color-scheme: dark)").matches` */
  prefersDark: boolean;
}

export interface ThemeOutput {
  theme: Theme;
  /** CSS `color-scheme` value for the root element. */
  colorScheme: "light" | "dark";
}

/**
 * Pure function: decides the effective theme from storage + OS preference.
 * No DOM side-effects — safe to call in vitest.
 */
export function resolveTheme(input: ThemeInput): ThemeOutput {
  const isDark =
    input.stored === "dark" || (!input.stored && input.prefersDark);
  const theme: Theme = isDark ? "dark" : "light";
  return { theme, colorScheme: theme };
}

/**
 * Applies the resolved theme to a DOM element (typically
 * `document.documentElement`).  Adds or removes the `dark` class and sets
 * the `color-scheme` style so native browser controls match.
 */
export function applyThemeToElement(
  el: HTMLElement,
  output: ThemeOutput,
): void {
  if (output.theme === "dark") {
    el.classList.add("dark");
  } else {
    el.classList.remove("dark");
  }
  el.style.colorScheme = output.colorScheme;
}

/**
 * End-to-end theme boot: reads real browser APIs and mutates the DOM.
 * Called from the inline <script> in index.html and from main.tsx.
 */
export function initTheme(): void {
  const stored = localStorage.getItem("theme");
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const output = resolveTheme({ stored, prefersDark });
  applyThemeToElement(document.documentElement, output);
}

/**
 * Sets the user's preferred theme, persists it to localStorage, and applies
 * it immediately.  Pass `null` to revert to OS preference.
 */
export function setTheme(theme: Theme | null): void {
  if (theme === null) {
    localStorage.removeItem("theme");
  } else {
    localStorage.setItem("theme", theme);
  }
  initTheme();
}
