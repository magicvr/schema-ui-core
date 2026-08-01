import { expect, test } from "@playwright/test";

// R6 browser E2E (I-008-005 forward path) updated for the R2 auth closed loop
// (GOAL-005): boots the Go API + Vite dev server via playwright webServer;
// verifies that an unauthenticated visit shows the login page, that signing in
// with the dev seed (admin / admin) renders the shell, and that the real auth
// chain (login -> /me -> records) works through the Web /api proxy.

test("login gates the shell and the real auth chain works through the proxy", async ({ page, request }) => {
  // Unauthenticated visit → login page, not the shell.
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  await expect(page.getByRole("button", { name: "Sign in" })).toBeVisible();

  // Sign in with the dev seed.
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill("admin");
  await page.getByRole("button", { name: "Sign in" }).click();

  // Shell renders and redirects home -> overview.
  await expect(page).toHaveURL(/\/overview$/);
  await expect(page.getByRole("heading", { name: "Overview" })).toBeVisible();
  await expect(page.getByText("Schema UI Core")).toBeVisible();

  // Manifest-driven navigation slots render (top / sidebar / user).
  await expect(page.getByRole("link", { name: "Overview" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Catalog" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Data table" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Settings" })).toBeVisible();

  // No manifest failure surface, and the sign-out control is present.
  await expect(page.getByText("MANIFEST_LOAD_FAILED")).toHaveCount(0);
  await expect(page.getByRole("button", { name: "Sign out" })).toBeVisible();

  // Real auth chain through the Web /api proxy -> Go API (independent context).
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

  // F-001 evidence: anonymous 401 fails closed through the Web /api proxy
  // (M8/M10). The request-level identity middleware must reject a missing
  // Bearer token on both the /me route and the records write routes.
  const anonMe = await request.get("/api/accounts/me");
  expect(anonMe.status()).toBe(401);
  expect((await anonMe.json()).error).toBe("UNAUTHENTICATED");

  const anonPatch = await request.patch("/api/records/rec-1", {
    data: { name: "Acme Console" },
  });
  expect(anonPatch.status()).toBe(401);

  // Control: the same write route passes once the seeded admin's access token
  // is attached, so the 401 above is a gate denial, not a missing route.
  const adminPatch = await request.patch("/api/records/rec-1", {
    headers,
    data: { name: "Acme Console" },
  });
  expect(adminPatch.status()).toBe(200);

  // Non-admin 403 is covered by API automation (records write gate):
  // apps/api/internal/handler/records_test.go::TestRecordsWriteDeniedWithoutAdminRole
  // exercises PATCH/DELETE with an editor-role token -> 403 FORBIDDEN. The
  // browser E2E exercises the same middleware path through the proxy (401
  // above); role denial is pure API logic already pinned by that test.

  await page.screenshot({ path: "test-results/r6-overview.png", fullPage: true });
});
