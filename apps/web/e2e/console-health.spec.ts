import { expect, test, type ConsoleMessage, type Page } from "@playwright/test";

import { openSidebarGroup, signInAsAdmin } from "./sign-in";

// Console health guard (workspace-010 GOAL-045 W33).
//
// Why this exists: React's "Maximum update depth exceeded" was being logged
// hundreds of times on every list page while EVERY suite stayed green — vitest
// renders fragments without the App shell, and the browser specs assert UI
// behaviour, not console output. A production app can loop on state updates and
// still "work", so nothing failed. This spec walks the real shell and fails on
// any console error or page error, which is what turned an invisible regression
// into a fixable finding.
//
// Scope note: it visits representative list pages that share the table/page-
// actions machinery (dashboard has no table; users/roles/jobs do). Jobs is
// admin-only by profile, so the route is visited only when it is mounted.

const appProfile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();

/** Errors we tolerate, with the reason. Keep this list narrow and justified. */
const ALLOWED: Array<{ pattern: RegExp; why: string }> = [
  // Static asset noise: a missing favicon/brand icon is a deployment concern,
  // not a runtime defect, and the harness serves no real brand assets.
  { pattern: /favicon|\.ico\b|\.png\b|\.svg\b/, why: "static asset" },
];

/**
 * Records console/page errors. Network-level "Failed to load resource" entries
 * carry the URL in `location()`, which is what makes them actionable: a 404 on
 * a real endpoint must still fail the guard, while an asset 404 must not.
 */
function collectProblems(page: Page): string[] {
  const problems: string[] = [];
  const record = (entry: string) => {
    if (ALLOWED.some((allowed) => allowed.pattern.test(entry))) {
      return;
    }
    problems.push(entry);
  };
  page.on("console", (message: ConsoleMessage) => {
    if (message.type() !== "error") {
      return;
    }
    const location = message.location();
    const where = typeof location?.url === "string" && location.url !== "" ? ` @ ${location.url}` : "";
    record(`[console.error] ${message.text()}${where}`);
  });
  page.on("pageerror", (error) => {
    record(`[pageerror] ${error.message}`);
  });
  return problems;
}

test("no console errors while walking the list pages", async ({ page }) => {
  // Sign in FIRST and only then start collecting: the pre-auth bootstrap
  // legitimately probes session endpoints and would report its 401s here.
  await signInAsAdmin(page);
  const problems = collectProblems(page);

  await openSidebarGroup(page, "manifest.nav.group.identityAccess");
  await page.getByRole("link", { name: "Users" }).click();
  await expect(page.getByRole("heading", { name: "Users" })).toBeVisible();
  await page.waitForTimeout(1000);

  await page.getByRole("link", { name: "Roles" }).click();
  await expect(page.getByRole("heading", { name: "Roles" })).toBeVisible();
  await page.waitForTimeout(1000);

  // Back to users: the round trip is what exposed the table state bleed.
  await page.getByRole("link", { name: "Users" }).click();
  await expect(page.getByRole("heading", { name: "Users" })).toBeVisible();
  await page.waitForTimeout(1000);

  if (appProfile === "admin") {
    await openSidebarGroup(page, "manifest.nav.group.operations");
    await page.getByRole("link", { name: "Jobs" }).click();
    await expect(page.getByRole("heading", { name: "Jobs" })).toBeVisible();
    await page.waitForTimeout(1000);
  }

  expect(problems, `console errors on list pages:\n${problems.join("\n")}`).toEqual([]);
});
