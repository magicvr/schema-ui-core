// FOUC bootstrap (S1 · C3 + W8 F-002): synchronous external script applied
// before any paint. Reads localStorage.theme and prefers-color-scheme to add
// the `dark` class and set color-scheme on <html> in the very first frame.
// This file is served as /theme-init.js from the web root and is allowed by
// the production CSP (script-src 'self') without an inline-script hash/nonce.
(function () {
  try {
    var stored = localStorage.getItem("theme");
    var prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    var isDark = stored === "dark" || (!stored && prefersDark);
    if (isDark) {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
    document.documentElement.style.colorScheme = isDark ? "dark" : "light";
  } catch (e) {
    // localStorage may be unavailable in some privacy contexts.
  }
})();