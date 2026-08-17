import { readFileSync } from "node:fs";

import Ajv, { type ValidateFunction } from "ajv";
import { expect, test } from "@playwright/test";

// R2 hygiene regression guard: the published runtime manifest must satisfy the
// pinned protocol schema (docs/schemas/app-manifest.schema.json — every block
// is additionalProperties: false). The browser host refuses non-protocol
// fields (UNKNOWN_MANIFEST_FIELD) and every shell surface fails; a dashboard
// fragment once leaked a non-protocol "order" nav field and took the whole
// suite down with it. This assertion turns that leak class into a targeted
// failure at the exact surface the host consumes.
const appManifestSchema = JSON.parse(
  readFileSync(new URL("../../../docs/schemas/app-manifest.schema.json", import.meta.url), "utf8"),
) as object;
const nodeSchema = JSON.parse(
  readFileSync(new URL("../../../docs/schemas/node.schema.json", import.meta.url), "utf8"),
) as object;
const pageSchema = JSON.parse(
  readFileSync(new URL("../../../docs/schemas/page.schema.json", import.meta.url), "utf8"),
) as object;
const actionSchema = JSON.parse(
  readFileSync(new URL("../../../docs/schemas/action.schema.json", import.meta.url), "utf8"),
) as object;
const reactionSchema = JSON.parse(
  readFileSync(new URL("../../../docs/schemas/reaction.schema.json", import.meta.url), "utf8"),
) as object;
const ajv = new Ajv({ allErrors: true, strict: false, validateSchema: false });
// The app-manifest schema has its own absolute $id (schema-ui.dev), so its
// relative $ref "node.schema.json#/definitions/VisibleWhen" resolves against
// that base URI. The referenced schemas use a different $id namespace
// (internal/schema-ui), so besides the filename registrations they must also
// be registered under the app-manifest base (mirrors how
// src/protocol/conformance/runtime-schema-validate.ts keeps refs resolvable).
const manifestSchemaBase = String(appManifestSchema.$id).replace(/[^/]*$/, "");
for (const [name, schema] of [
  ["node", nodeSchema],
  ["page", pageSchema],
  ["action", actionSchema],
  ["reaction", reactionSchema],
] as const) {
  ajv.addSchema(schema, `${name}.schema.json`);
  ajv.addSchema(schema, `${manifestSchemaBase}${name}.schema.json`);
}
const validateManifestSchema: ValidateFunction = ajv.compile(appManifestSchema);

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
  // The runtime manifest must be structurally valid per the pinned protocol
  // schema; a fragment leaking a non-protocol field would be refused by the
  // host (UNKNOWN_MANIFEST_FIELD) and break every browser surface.
  const schemaValid = validateManifestSchema(manifest);
  expect(
    schemaValid,
    validateManifestSchema.errors
      ? validateManifestSchema.errors.map((error) => `${error.instancePath || "/"}: ${error.message}`).join("; ")
      : "runtime manifest failed app-manifest.schema.json validation",
  ).toBe(true);
  const manifestPageIds = manifest.pages.map((page: { pageId: string }) => page.pageId);
  // W1 (GOAL-002 / workspace-010): production defaults ship no dev.examples.
  expect(manifestPageIds).toEqual(expect.arrayContaining(["users", "roles"]));
  expect(manifestPageIds.includes("overview")).toBe(isDemoProfile);
  expect(manifestPageIds.includes("data-table")).toBe(isDemoProfile);
  // F-01 (GOAL-003): mvp/admin production home is now the dashboard.
  expect(manifest.app.homePageRef).toBe(isDemoProfile ? "overview" : "dashboard");
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

  // W16-F01: a fresh seed forces an initial password replacement. If an
  // earlier spec already replaced it, retry with the shared e2e password.
  const forced = page.getByRole("heading", { name: "Change your password" });
  if (await forced.isVisible().catch(() => false)) {
    await page.getByLabel("Current password").fill("admin");
    await page.getByLabel("New password").fill("admin-e2e-pass");
    await page.getByLabel("Confirm new password").fill("admin-e2e-pass");
    await page.getByRole("button", { name: "Change password" }).click();
  } else {
    await page.getByLabel("Password").fill("admin-e2e-pass");
    await page.getByRole("button", { name: "Sign in" }).click();
  }

  // Shell renders and redirects home -> manifest home (demo: overview; else dashboard).
  await expect(page).toHaveURL(isDemoProfile ? /\/overview$/ : /\/dashboard$/);
  await expect(page.getByRole("heading", { name: isDemoProfile ? "Overview" : "Dashboard" })).toBeVisible();
  // The shell brand renders the siteTitle from the public startup config. Read
  // it from /api/branding rather than hardcoding the default: an earlier spec
  // (localization M3) may have PATCHed a custom siteTitle into the shared
  // playwright SQLite, so the header value is data-dependent, not constant.
  const branding = await request.get("/api/branding");
  expect(branding.status()).toBe(200);
  const brandBody = await branding.json();
  expect(typeof brandBody.siteTitle).toBe("string");
  // W13 T-02: the brand text now appears twice in the DOM — the mobile-only
  // brand bar (hidden at the desktop viewport) and the desktop single-row
  // brand link — so target the LAST occurrence (the visible desktop one).
  await expect(page.getByText(brandBody.siteTitle).last()).toBeVisible();

  // Manifest-driven navigation slots render (top / sidebar / user). Production
  // profiles expose no dev.examples navigation (S5 hygiene); demo does.
  await expect(page.getByRole("link", { name: "Users" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Overview" })).toHaveCount(isDemoProfile ? 1 : 0);
  await expect(page.getByRole("link", { name: "Data table" })).toHaveCount(isDemoProfile ? 1 : 0);
  if (isAdminProfile) {
    // T-01 (GOAL-013 D-002): Settings lives in the topbar user dropdown
    // (user-chain), not the sidebar — asserted below via the menu.
    await expect(page.getByRole("link", { name: "Activity" })).toBeVisible();
    // S-02/S-01/S-03/S-04 (GOAL-007/008/009/010): admin-only surfaces.
    await expect(page.getByRole("link", { name: "File library" })).toBeVisible();
    await expect(page.getByRole("link", { name: "Data dictionary" })).toBeVisible();
    await expect(page.getByRole("link", { name: "System monitoring" })).toBeVisible();
    await expect(page.getByRole("link", { name: "Scheduled tasks" })).toBeVisible();
    // S-11 (GOAL-011): login captcha settings page (admin-only surface).
    // S-12 (GOAL-012): recycle bin page (admin-only surface).
    await expect(page.getByRole("link", { name: "Recycle bin" })).toBeVisible();
  } else {
    await expect(page.getByRole("link", { name: "Settings" })).toHaveCount(0);
    await expect(page.getByRole("link", { name: "Activity" })).toHaveCount(0);
    await expect(page.getByRole("link", { name: "File library" })).toHaveCount(0);
    await expect(page.getByRole("link", { name: "Data dictionary" })).toHaveCount(0);
    await expect(page.getByRole("link", { name: "System monitoring" })).toHaveCount(0);
    await expect(page.getByRole("link", { name: "Scheduled tasks" })).toHaveCount(0);
    await expect(page.getByRole("link", { name: "Recycle bin" })).toHaveCount(0);
  }

  // No manifest failure surface. The user chain (个人中心 / 我的钱包 / 设置
  // + 退出登录) lives in the topbar user dropdown (W12 T-01) — open it and
  // assert the menu items.
  await expect(page.getByText("MANIFEST_LOAD_FAILED")).toHaveCount(0);
  await page.getByRole("button", { name: "User menu" }).click();
  await expect(page.getByRole("menuitem", { name: "Settings" })).toHaveCount(isAdminProfile ? 1 : 0);
  await expect(page.getByRole("menuitem", { name: "Sign out" })).toBeVisible();
  await page.keyboard.press("Escape");

  // W13 T-04: the theme toggle renders LEFT of the language switcher.
  const themeBox = await page.getByRole("button", { name: "Toggle color theme" }).boundingBox();
  const langBox = await page.getByRole("button", { name: "Language" }).boundingBox();
  expect(themeBox).not.toBeNull();
  expect(langBox).not.toBeNull();
  expect(themeBox!.x).toBeLessThan(langBox!.x);

  // W13 T-02: on a mobile viewport the logo + site title own a dedicated
  // brand bar on top, and the functional area keeps its own row below.
  await page.setViewportSize({ width: 390, height: 844 });
  await expect(page.locator('[data-shell-region="mobile-brandbar"]')).toBeVisible();
  await expect(
    page
      .locator('[data-shell-region="mobile-brandbar"]')
      .getByText(brandBody.siteTitle),
  ).toBeVisible();
  await expect(page.getByRole("button", { name: "Open navigation menu" })).toBeVisible();

  // Real auth chain through the Web /api proxy -> Go API (independent context).
  const login = await request.post("/api/auth/login", {
    data: { username: "admin", password: "admin-e2e-pass" },
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

  // W13 T-05: avatar self-service through the real API — upload a PNG, commit
  // it to the profile, reload, and the user-menu trigger shows the avatar.
  const avatarUpload = await request.post("/api/account/avatar", {
    headers,
    multipart: {
      file: {
        name: "avatar.png",
        mimeType: "image/png",
        buffer: Buffer.from(
          "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==",
          "base64",
        ),
      },
    },
  });
  expect(avatarUpload.status()).toBe(200);
  const avatarBody = await avatarUpload.json();
  expect(avatarBody.url).toMatch(/^\/api\/account\/avatars\//);
  const avatarPatch = await request.patch("/api/account/profile", {
    headers,
    data: { name: "Admin", avatarUrl: avatarBody.url },
  });
  expect(avatarPatch.status()).toBe(200);
  await page.reload();
  await expect(page.locator(`img[src="${avatarBody.url}"]`).first()).toBeVisible();

  // W13 T-05 follow-up (user 2026-08-16): the header must refresh IMMEDIATELY
  // after a UI-driven avatar save — no reload. Drive the account page upload
  // control, save the profile, and assert the user-menu trigger swaps to the
  // new avatar via the account.profile session refresh.
  await page.getByRole("button", { name: "User menu" }).click();
  await page.getByRole("menuitem", { name: "Account" }).click();
  await expect(page).toHaveURL(/\/account$/);
  await page.locator("#field-avatarUrl").setInputFiles({
    name: "avatar-ui.png",
    mimeType: "image/png",
    buffer: Buffer.from(
      "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==",
      "base64",
    ),
  });
  // The upload settles and the form preview shows the new avatar.
  const profileForm = page.locator("form").filter({ has: page.getByLabel("Avatar") });
  const preview = profileForm.locator('img[src^="/api/account/avatars/"]');
  await expect(preview.first()).toBeVisible();
  const uiAvatarUrl = (await preview.first().getAttribute("src"))!;
  await page.getByRole("button", { name: "Save profile" }).click();
  // The header avatar updates WITHOUT a reload (account.profile session refresh).
  await expect(page.locator(`img[src="${uiAvatarUrl}"]`).first()).toBeVisible();

  // W13 T-06: the bell dropdown's "View all" reaches the notifications page,
  // which renders the interactive center (custom node) + settings form. The
  // fresh SQLite DB has no notifications → the empty state is shown.
  await page.getByRole("button", { name: "Notifications" }).click();
  await page.getByRole("menuitem", { name: "View all" }).click();
  await expect(page).toHaveURL(/\/notifications$/);
  await expect(page.locator("[data-notification-center]")).toBeVisible();
  await expect(page.getByText("No items match")).toBeVisible();

  await page.screenshot({ path: "test-results/r6-shell-users.png", fullPage: true });
});