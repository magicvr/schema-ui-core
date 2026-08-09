/**
 * Theme pure-logic unit (S1 · C3 + VP-007 S3 system default).
 *
 * `applyTheme` is a side-effect-free function that computes the desired theme
 * and then applies it to the provided document root.  Keeping the decision
 * logic in a plain function lets vitest exercise it without a browser.
 *
 * `initTheme` is the top-level boot call used by the synchronous inline
 * script in index.html and by main.tsx (no longer needed after inline
 * migration, kept for backwards-compat import path).
 *
 * VP-007 S3 priority (D-002, user-confirmed): user explicit choice →
 * system default theme (from the public startup configuration, non-auto) →
 * OS preference. The user's explicit choice always wins.
 */

export type Theme = "light" | "dark";

export interface ThemeInput {
  /** Value of `localStorage.getItem("theme")` — null when absent. */
  stored: string | null;
  /** `window.matchMedia("(prefers-color-scheme: dark)").matches` */
  prefersDark: boolean;
  /** System default theme from /api/branding (auto/light/dark); null = auto. */
  systemDefault?: string | null;
}

export interface ThemeOutput {
  theme: Theme;
  /** CSS `color-scheme` value for the root element. */
  colorScheme: "light" | "dark";
}

/**
 * Pure function: decides the effective theme from storage + system default +
 * OS preference. No DOM side-effects — safe to call in vitest.
 */
export function resolveTheme(input: ThemeInput): ThemeOutput {
  if (input.stored === "dark" || input.stored === "light") {
    return { theme: input.stored, colorScheme: input.stored };
  }
  if (input.systemDefault === "dark" || input.systemDefault === "light") {
    return { theme: input.systemDefault, colorScheme: input.systemDefault };
  }
  const isDark = !input.stored && input.prefersDark;
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
 * Applies the system default theme (from the public startup configuration)
 * when the user has no explicit choice. VP-007 S3: the explicit user theme
 * always wins; otherwise the system default (non-auto) beats the OS
 * preference. Called after branding loads on the login page and the shell.
 */
export function applySystemDefaultTheme(systemDefault: string | null | undefined): void {
  const stored = localStorage.getItem("theme");
  if (stored === "dark" || stored === "light") {
    return;
  }
  if (systemDefault !== "light" && systemDefault !== "dark") {
    return;
  }
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const output = resolveTheme({ stored, prefersDark, systemDefault });
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
