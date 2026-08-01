import { expect, test } from "@playwright/test";

// R6 browser E2E (I-008-005 forward path) updated for R2 real auth (GOAL-005):
// boots the Go API + Vite dev server via playwright webServer; verifies the
// shell renders from the real manifest, then walks the real auth chain
// (login -> /me -> records) through the Web /api proxy in a real browser.
//
// The Go API seeds the dev admin (admin / admin) when APP_ENV defaults to
// development, so no test-only env is needed.

test("shell renders and the real auth chain works through the proxy", async ({ page, request }) => {
  // Home redirect resolves from manifest.homePageRef.
  await page.goto("/");
  await expect(page).toHaveURL(/\/overview$/);
  await expect(page.getByRole("heading", { name: "Overview" })).toBeVisible();
  await expect(page.getByText("Schema UI Core")).toBeVisible();

  // Manifest-driven navigation slots render (top / sidebar / user).
  await expect(page.getByRole("link", { name: "Overview" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Catalog" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Data table" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Settings" })).toBeVisible();

  // No manifest failure surface.
  await expect(page.getByText("MANIFEST_LOAD_FAILED")).toHaveCount(0);

  // The boot /me without a token is now 401 (no dev-session fallback), so the
  // shell surfaces the non-blocking account banner instead of a session.
  await expect(page.getByText("Account session failed to load")).toBeVisible();

  // Real auth chain through the Web /api proxy -> Go API.
  const login = await request.post("/api/auth/login", {
    data: { username: "admin", password: "admin" },
  });
  expect(login.status()).toBe(200);
  const tokens = await login.json();
  expect(tokens.accessToken).toBeTruthy();
  expect(tokens.refreshToken).toBeTruthy();
  expect(tokens.user.id).toBe("user-admin");

  const headers = { Authorization: `Bearer ${tokens.accessToken}` };

  // Request-level identity: /me resolves to the seeded admin, not a static dev
  // session.
  const me = await request.get("/api/accounts/me", { headers });
  expect(me.status()).toBe(200);
  const session = await me.json();
  expect(session.user.id).toBe("user-admin");
  expect(session.user.roles).toEqual(expect.arrayContaining(["admin"]));

  // Records read route remains reachable through the same proxy.
  const records = await request.get("/api/records", { headers });
  expect(records.status()).toBe(200);
  const list = await records.json();
  expect(list.items.length).toBeGreaterThan(0);

  await page.screenshot({ path: "test-results/r6-overview.png", fullPage: true });
});
