import { type Page, expect } from "@playwright/test";

// Shared W16-F01-aware sign-in for the browser E2E suite.
//
// A fresh seed forces the initial admin password to be replaced (must_change_password=1).
// On a clean SQLite this helper performs the real forced-change flow; on a shared
// server where an earlier spec already replaced the password it falls back to the
// shared e2e password. This keeps the whole suite consistent with I-008-002 v0.1.3 —
// the first-login password change is a real user step, never a test-side bypass.
export const E2E_INITIAL_PASSWORD = "admin";
export const E2E_PASSWORD = "admin-e2e-pass";

const homeUrl = (profile: string) => {
  const homeRe = profile === "demo" ? /\/overview$/ : /\/dashboard$/;
  return homeRe;
};

export async function signInAsAdmin(page: Page): Promise<void> {
  const profile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();
  const homeRe = homeUrl(profile);

  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill(E2E_INITIAL_PASSWORD);
  await page.getByRole("button", { name: "Sign in" }).click();

  // Already replaced (a prior spec changed the password): land straight on home.
  try {
    await page.waitForURL(homeRe, { timeout: 8000 });
    return;
  } catch {
    // Not signed in yet — next, the forced-change screen or a failed initial login.
  }

  const forced = page.getByRole("heading", { name: "Change your password" });
  try {
    await forced.waitFor({ state: "visible", timeout: 8000 });
    await page.getByLabel("Current password").fill(E2E_INITIAL_PASSWORD);
    await page.getByLabel("New password", { exact: true }).fill(E2E_PASSWORD);
    await page.getByLabel("Confirm new password", { exact: true }).fill(E2E_PASSWORD);
    await page.getByRole("button", { name: "Change password" }).click();
    await page.waitForURL(homeRe, { timeout: 15000 });
    return;
  } catch {
    // Fresh-initial login failed because the password was already replaced:
    // fall back to the shared e2e password.
    await page.getByLabel("Password").fill(E2E_PASSWORD);
    await page.getByRole("button", { name: "Sign in" }).click();
    await page.waitForURL(homeRe, { timeout: 15000 });
  }
}