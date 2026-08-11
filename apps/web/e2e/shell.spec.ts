import { expect, test } from "@playwright/test";

// R6 browser E2E (I-008-005 forward path) updated for the R2 auth closed loop
// (GOAL-005): boots the Go API + Vite dev server via playwright webServer;
// verifies that an unauthenticated visit shows the login page, that signing in
// with the dev seed (admin / admin) renders the shell, and that the real auth
// chain (login -> /me -> users) works through the Web /api proxy.

test("login gates the shell and the real auth chain works through the proxy", async ({ page, request }) => {
  const profile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();
  const isAdminProfile = profile === "admin";
  // W2 (GOAL-003 / workspace-010): demo = mvp capability + dev.examples, so it
  // exposes the example pages and homes to overview; mvp/admin stay examples-free.
  const isDemoProfile = profile === "demo";

  // The same browser build must consume the runtime Manifest selected by the
  // API profile. The proxy must not silently fall back to a static Web file.
  const manifestResponse = await request.get("/.well-known/schema-ui/app-manifest.json");
  expect(manifestResponse.status()).toBe(200);
  expect(manifestResponse.headers()["x-schema-ui-manifest-source"]).toBe("api");
  const manifest = await manifestResponse.json();
  const manifestPageIds = manifest.pages.map((page: { pageId: string }) => page.pageId);
  // W1 (GOAL-002 / workspace-010): production defaults ship no dev.examples.
  expect(manifestPageIds).toEqual(expect.arrayContaining(["users", "roles"]));
  expect(manifestPageIds.includes("overview")).toBe(isDemoProfile);
  expect(manifestPageIds.includes("data-table")).toBe(isDemoProfile);
  expect(manifest.app.homePageRef).toBe(isDemoProfile ? "overview" : "users");
  expect(manifestPageIds.includes("settings")).toBe(isAdminProfile);
  expect(manifestPageIds.includes("activity")).toBe(isAdminProfile);

  const settingsSchema = await request.get("/api/schema/settings");
  const activitySchema = await request.get("/api/schema/activity");
  expect(settingsSchema.status()).toBe(isAdminProfile ? 200 : 404);
  expect(activitySchema.status()).toBe(isAdminProfile ? 200 : 404);
  const settingsRoute = await request.get("/api/settings");
  const activityRoute = await request.get("/api/operations");
  expect(settingsRoute.status()).toBe(isAdminProfile ? 401 : 404);
  expect(activityRoute.status()).toBe(isAdminProfile ? 401 : 404);

  // Unauthenticated visit → login page, not the shell.
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  await expect(page.getByRole("button", { name: "Sign in" })).toBeVisible();

  // Sign in with the dev seed.
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill("admin");
  await page.getByRole("button", { name: "Sign in" }).click();

  // Shell renders and redirects home -> manifest home (demo: overview; else users).
  await expect(page).toHaveURL(isDemoProfile ? /\/overview$/ : /\/users$/);
  await expect(page.getByRole("heading", { name: isDemoProfile ? "Overview" : "Users" })).toBeVisible();
  // The shell brand renders the siteTitle from the public startup config. Read
  // it from /api/branding rather than hardcoding the default: an earlier spec
  // (localization M3) may have PATCHed a custom siteTitle into the shared
  // playwright SQLite, so the header value is data-dependent, not constant.
  const branding = await request.get("/api/branding");
  expect(branding.status()).toBe(200);
  const brandBody = await branding.json();
  expect(typeof brandBody.siteTitle).toBe("string");
  await expect(page.getByText(brandBody.siteTitle).first()).toBeVisible();

  // Manifest-driven navigation slots render (top / sidebar / user). Production
  // profiles expose no dev.examples navigation (S5 hygiene); demo does.
  await expect(page.getByRole("link", { name: "Users" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Overview" })).toHaveCount(isDemoProfile ? 1 : 0);
  await expect(page.getByRole("link", { name: "Data table" })).toHaveCount(isDemoProfile ? 1 : 0);
  if (isAdminProfile) {
    await expect(page.getByRole("link", { name: "Settings" })).toBeVisible();
    await expect(page.getByRole("link", { name: "Activity" })).toBeVisible();
  } else {
    await expect(page.getByRole("link", { name: "Settings" })).toHaveCount(0);
    await expect(page.getByRole("link", { name: "Activity" })).toHaveCount(0);
  }

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

  // Users read route remains reachable through the same proxy (GOAL-011).
  const users = await request.get("/api/users", { headers });
  expect(users.status()).toBe(200);
  const list = await users.json();
  expect(list.items.length).toBeGreaterThan(0);

  // F-001 evidence: anonymous 401 fails closed through the Web /api proxy
  // (M8/M10). The request-level identity middleware must reject a missing
  // Bearer token on both the /me route and the users write routes.
  const anonMe = await request.get("/api/accounts/me");
  expect(anonMe.status()).toBe(401);
  expect((await anonMe.json()).error).toBe("UNAUTHENTICATED");

  const anonPatch = await request.patch("/api/users/user-admin", {
    data: { name: "Admin" },
  });
  expect(anonPatch.status()).toBe(401);

  // Control: the same write route passes once the seeded admin's access token
  // is attached, so the 401 above is a gate denial, not a missing route.
  const adminPatch = await request.patch("/api/users/user-admin", {
    headers,
    data: { name: "Admin" },
  });
  expect(adminPatch.status()).toBe(200);

  // Non-admin 403 is covered by API automation (users write gate):
  // apps/api/internal/handler/users_test.go::TestUsersAuthGates exercises POST
  // with a viewer token -> 403 FORBIDDEN. The browser E2E exercises the same
  // middleware path through the proxy (401 above); role denial is pure API logic
  // already pinned by that test.

  await page.screenshot({ path: "test-results/r6-shell-users.png", fullPage: true });
});
