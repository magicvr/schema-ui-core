import { expect, test } from "@playwright/test";

import { E2E_INITIAL_PASSWORD, E2E_PASSWORD } from "./sign-in";

// W16-F01 (v0.1.3 reverse-carried into browser E2E): on a fresh seed the
// seeded admin MUST replace the initial password before business APIs are
// allowed. This spec runs first (lexicographic) so the forced-change surface is
// guaranteed on a clean SQLite; later specs reuse E2E_PASSWORD through their
// sign-in helper fallback branch.
//
// R5 (workspace-038 · GOAL-006) — why the "00-" prefix:
// The ordering assumption above silently broke once VP-036 added
// `command-palette.spec.ts`. Every spec shares ONE scratch database, and that
// file sorts before this one, so its sign-in helper consumed the fresh-seed
// state (the helper performs the forced change itself); this spec then failed
// with "invalid username or password" for admin/admin. Verified both ways: it
// passes in isolation (`npx playwright test 00-force-password-change.spec.ts`)
// and failed inside the full run before the rename. The prefix makes the
// documented order real instead of aspirational.
test("fresh seed forces initial password change before business access", async ({ page, request }) => {
  const profile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();
  const home = profile === "demo" ? "Overview" : "Dashboard";

  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password", { exact: true }).fill(E2E_INITIAL_PASSWORD); // W22: exact avoids matching the visibility-toggle aria-label
  await page.getByRole("button", { name: "Sign in" }).click();

  // The forced-change surface must appear on the fresh seed.
  await expect(page.getByRole("heading", { name: "Change your password" })).toBeVisible();
  await page.getByLabel("Current password").fill(E2E_INITIAL_PASSWORD);
  await page.getByLabel("New password", { exact: true }).fill(E2E_PASSWORD);
  await page.getByLabel("Confirm new password", { exact: true }).fill(E2E_PASSWORD);
  await page.getByRole("button", { name: "Change password" }).click();

  // After the real change the shell opens.
  await expect(page).toHaveURL(profile === "demo" ? /\/overview$/ : /\/dashboard$/);
  await expect(page.getByRole("heading", { name: home })).toBeVisible();

  // Business API now works with the NEW password only; the initial one is gone.
  const oldLogin = await request.post("/api/auth/login", {
    data: { username: "admin", password: E2E_INITIAL_PASSWORD },
  });
  expect(oldLogin.status()).toBe(401);

  const newLogin = await request.post("/api/auth/login", {
    data: { username: "admin", password: E2E_PASSWORD },
  });
  expect(newLogin.status()).toBe(200);
  const tokens = await newLogin.json();
  expect(tokens.user.mustChangePassword).toBe(false);

  const users = await request.get("/api/users", {
    headers: { Authorization: `Bearer ${tokens.accessToken}` },
  });
  expect(users.status()).toBe(200);

  await page.screenshot({ path: "test-results/r6-w16-force-password-change.png", fullPage: true });
});