#!/usr/bin/env node
/**
 * Production CSP + real-browser theme bootstrap check (W8 F-002 follow-up).
 *
 * Requires the production compose stack on http://127.0.0.1:25081
 * (docker compose up --build -d). Launches headless Chromium and verifies:
 *   - the document response carries the production Content-Security-Policy
 *     header with script-src 'self';
 *   - /theme-init.js is served as an external same-origin script;
 *   - the browser executes the theme bootstrap (dark class + color-scheme);
 *   - no CSP "Refused to execute" / "Content Security Policy" console errors.
 */
import { chromium } from "@playwright/test";

const baseURL = process.env.PROD_WEB_URL || "http://127.0.0.1:25081";
const results = {
  baseURL,
  ok: true,
  checks: {},
  consoleErrors: [],
};

function requireCheck(name, ok, detail) {
  results.checks[name] = { ok: Boolean(ok), detail };
  if (!ok) results.ok = false;
}

const browser = await chromium.launch({ headless: true });
const context = await browser.newContext({ viewport: { width: 1280, height: 800 } });
const page = await context.newPage();

page.on("console", (msg) => {
  const text = msg.text();
  if (/Content Security Policy|Refused to execute inline script/i.test(text)) {
    results.consoleErrors.push({ type: msg.type(), text });
  }
});
page.on("pageerror", (err) => {
  if (/Content Security Policy|Refused to execute inline script/i.test(String(err))) {
    results.consoleErrors.push({ type: "pageerror", text: String(err) });
  }
});

try {
  const response = await page.goto(baseURL + "/", { waitUntil: "networkidle" });
  const headers = response.headers();
  const csp = headers["content-security-policy"] || "";
  requireCheck(
    "csp-header",
    csp.includes("script-src 'self'") && csp.includes("default-src 'self'"),
    `Content-Security-Policy length=${csp.length}; present=${Boolean(csp)}`,
  );

  // Verify the production HTML references the external bootstrap script.
  const htmlRef = await page.evaluate(() => {
    const script = document.querySelector('script[src="/theme-init.js"]');
    const inlineBootstrap = Array.from(document.scripts).some(
      (s) => !s.src && s.textContent.includes('localStorage.getItem("theme")'),
    );
    return { hasExternalThemeInit: Boolean(script), hasInlineBootstrap: inlineBootstrap };
  });
  requireCheck("external-theme-init-referenced", htmlRef.hasExternalThemeInit, JSON.stringify(htmlRef));
  requireCheck("no-inline-theme-bootstrap", !htmlRef.hasInlineBootstrap, JSON.stringify(htmlRef));

  // Theme script must run on first paint (no stored preference -> light/dark set).
  const firstTheme = await page.evaluate(() => ({
    colorScheme: document.documentElement.style.colorScheme,
    darkClass: document.documentElement.classList.contains("dark"),
  }));
  requireCheck(
    "theme-bootstrap-first-paint",
    firstTheme.colorScheme === "light" || firstTheme.colorScheme === "dark",
    `firstPaint colorScheme=${firstTheme.colorScheme}, darkClass=${firstTheme.darkClass}`,
  );

  // Stored explicit dark theme must survive a reload (the bootstrap script reads
  // localStorage before the React bundle mounts).
  await page.evaluate(() => localStorage.setItem("theme", "dark"));
  await page.reload({ waitUntil: "networkidle" });
  const storedTheme = await page.evaluate(() => ({
    colorScheme: document.documentElement.style.colorScheme,
    darkClass: document.documentElement.classList.contains("dark"),
  }));
  requireCheck(
    "theme-bootstrap-stored-dark",
    storedTheme.colorScheme === "dark" && storedTheme.darkClass === true,
    `after reload colorScheme=${storedTheme.colorScheme}, darkClass=${storedTheme.darkClass}`,
  );

  // The external script itself must be served from the same origin.
  const themeInitResponse = await context.request.get(`${baseURL}/theme-init.js`);
  requireCheck(
    "theme-init-served",
    themeInitResponse.ok() && themeInitResponse.headers()["content-type"]?.includes("javascript"),
    `status=${themeInitResponse.status()} content-type=${themeInitResponse.headers()["content-type"] || "missing"}`,
  );

  requireCheck(
    "no-csp-console-violations",
    results.consoleErrors.length === 0,
    JSON.stringify(results.consoleErrors),
  );
} catch (err) {
  results.ok = false;
  results.fatal = String(err);
} finally {
  await browser.close();
}

console.log(JSON.stringify(results, null, 2));
process.exit(results.ok ? 0 : 1);