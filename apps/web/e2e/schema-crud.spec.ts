import { expect, test, type Page } from "@playwright/test";

// A-010 R-004 · real browser Schema CRUD lifecycle against Go + SQLite.
// GOAL-011 S3 repoints the driver from the retired demo page to the users
// resource page. Boots via playwright webServer (same as shell.spec.ts): Go API
// + Vite, then drives the users representative page through create → edit →
// delete with confirm. T-UI-01～10 cover Renderer behavior with an in-memory
// API emulator; this file proves the browser → proxy → Go/SQLite path.

async function signInAsAdmin(page: Page): Promise<void> {
  // Home redirect follows the manifest home: demo -> overview, else dashboard (F-01).
  const profile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill("admin");
  await page.getByRole("button", { name: "Sign in" }).click();

  // W16-F01: a fresh seed forces an initial password replacement.
  const forced = page.getByRole("heading", { name: "Change your password" });
  if (await forced.isVisible().catch(() => false)) {
    await page.getByLabel("Current password").fill("admin");
    await page.getByLabel("New password").fill("admin-e2e-pass");
    await page.getByLabel("Confirm new password").fill("admin-e2e-pass");
    await page.getByRole("button", { name: "Change password" }).click();
  } else {
    // If an earlier spec already replaced the seed password, retry with it.
    await page.getByLabel("Password").fill("admin-e2e-pass");
    await page.getByRole("button", { name: "Sign in" }).click();
  }
  await expect(page).toHaveURL(profile === "demo" ? /\/overview$/ : /\/dashboard$/);
}

test("users and roles drive real authorization management against Go SQLite", async ({
  page,
  request,
}) => {
  const stamp = Date.now().toString(36);
	const createdUsername = `e2e-${stamp}`;
  const createdName = `E2E CRUD ${stamp}`;
  const editedName = `E2E CRUD edited ${stamp}`;
	const roleKey = `e2e_${stamp}`;
	const replacementPassword = "  e2e-password-new  ";

  await signInAsAdmin(page);

  // Menu projection (GOAL-011 S4): admin seed grants menu_users; login must
  // resolve /me features so the link is present after sign-in.
  await expect(page.getByRole("link", { name: "Users" })).toBeVisible();
  await page.getByRole("link", { name: "Users" }).click();
  await expect(page).toHaveURL(/\/users$/);
  await expect(page.getByRole("heading", { name: "Users" })).toBeVisible();
  await expect(page.getByRole("button", { name: "New user" })).toBeEnabled();

  // --- Create ---
  await page.getByRole("button", { name: "New user" }).click();
  const createDialog = page.getByRole("dialog", { name: "New user" });
  await expect(createDialog).toBeVisible();
	await createDialog.getByLabel("Username").fill(createdUsername);
  await createDialog.getByLabel("Name", { exact: true }).fill("E2E Bot");
  await createDialog.getByLabel("Password").fill("e2e-password");
  await createDialog.getByRole("button", { name: "Create user" }).click();

  await expect(page.getByText("Item created")).toBeVisible();
	await expect(page.getByRole("cell", { name: createdUsername })).toBeVisible();

  // --- Edit (row-scoped) ---
	const createdRow = page.getByRole("row").filter({ hasText: createdUsername });
  await createdRow.getByRole("button", { name: "Edit" }).click();
  const editDialog = page.getByRole("dialog", { name: "Edit" });
  await expect(editDialog).toBeVisible();
  await editDialog.getByLabel("Name").fill(editedName);
  await editDialog.getByRole("button", { name: "Save changes" }).click();

  await expect(page.getByText("Item updated")).toBeVisible();
  await expect(page.getByRole("cell", { name: editedName })).toBeVisible();

	// --- Custom role + grants ---
	await page.getByRole("link", { name: "Roles" }).click();
	await expect(page).toHaveURL(/\/roles$/);
	const adminRoleRow = page.getByRole("row").filter({ hasText: "role-admin" });
	await expect(adminRoleRow.getByRole("button", { name: "Edit" })).toBeDisabled();
	await expect(adminRoleRow.getByRole("button", { name: "Delete" })).toBeDisabled();
	await page.getByRole("button", { name: "New role" }).click();
	const createRoleDialog = page.getByRole("dialog", { name: "New role" });
	await createRoleDialog.getByLabel("Key").fill(roleKey);
	await createRoleDialog.getByLabel("Name", { exact: true }).fill("E2E Support");
	await createRoleDialog.getByLabel("users.read").check(); // W11 U-02: catalog labels = permission keys
	await createRoleDialog.getByLabel("Users", { exact: true }).check();
	await createRoleDialog.getByRole("button", { name: "Create role" }).click();
	await expect(page.getByRole("cell", { name: roleKey, exact: true })).toBeVisible();

	const roleRow = page.getByRole("row").filter({ hasText: roleKey });
	await roleRow.getByRole("button", { name: "Edit" }).click();
	const editRoleDialog = page.getByRole("dialog", { name: "Edit" });
	await expect(editRoleDialog.getByLabel("users.read")).toBeChecked();
	await editRoleDialog.getByLabel("roles.read").check();
	await editRoleDialog.getByRole("button", { name: "Save changes" }).click();

	// --- Assign the custom role and rotate password through distinct forms ---
	await page.getByRole("link", { name: "Users" }).click();
	const managedRow = page.getByRole("row").filter({ hasText: createdUsername });
	await managedRow.getByRole("button", { name: "Roles" }).click();
	const rolesDialog = page.getByRole("dialog", { name: "Roles" });
	// W11 U-01: the roles assignment is a checkboxGroup of role names (optionsSource /api/roles).
	await rolesDialog.getByLabel("E2E Support").check();
	await rolesDialog.getByRole("button", { name: "Save roles" }).click();
	await expect(rolesDialog).toBeHidden();

	// W11 U-05: row actions beyond the primary two collapse into a "⋯ More" menu.
	await managedRow.getByRole("button", { name: "More actions" }).click();
	await page.getByRole("menuitem", { name: "Password" }).click();
	const passwordDialog = page.getByRole("dialog", { name: "Password" });
	await expect(passwordDialog.getByLabel("New password")).toHaveAttribute("type", "password");
	await passwordDialog.getByLabel("New password").fill(replacementPassword);
	await passwordDialog.getByRole("button", { name: "Change password" }).click();
	await expect(passwordDialog).toBeHidden();

	const oldLogin = await request.post("/api/auth/login", {
	  data: { username: createdUsername, password: "e2e-password" },
	});
	expect(oldLogin.status()).toBe(401);
	const trimmedLogin = await request.post("/api/auth/login", {
	  data: { username: createdUsername, password: replacementPassword.trim() },
	});
	expect(trimmedLogin.status()).toBe(401);
	const login = await request.post("/api/auth/login", {
	  data: { username: createdUsername, password: replacementPassword },
	});
	expect(login.status()).toBe(200);
	const tokens = await login.json();
	expect(tokens.user.permissions).toEqual(expect.arrayContaining(["users.read", "roles.read"]));
	const session = await request.get("/api/accounts/me", {
	  headers: { Authorization: `Bearer ${tokens.accessToken}` },
	});
	expect(session.status()).toBe(200);
	expect((await session.json()).features.menu_users).toBe(true);

  // --- Delete with confirm ---
  const editedRow = page.getByRole("row").filter({ hasText: editedName });
  await editedRow.getByRole("button", { name: "More actions" }).click();
  await page.getByRole("menuitem", { name: "Delete" }).click();
  const confirmDialog = page.getByRole("dialog", { name: "Confirm action" });
  await expect(confirmDialog).toBeVisible();
  await expect(confirmDialog.getByText("Delete this user?")).toBeVisible();
  await confirmDialog.getByRole("button", { name: "Confirm" }).click();

  await expect(page.getByText("Item deleted")).toBeVisible();
  await expect(page.getByRole("cell", { name: editedName })).toHaveCount(0);

	await page.getByRole("link", { name: "Roles" }).click();
	const freeRoleRow = page.getByRole("row").filter({ hasText: roleKey });
	// Roles rows carry exactly two actions (Edit + Delete) — inline, no overflow.
	await expect(freeRoleRow.getByRole("button", { name: "Delete" })).toBeEnabled();
	await freeRoleRow.getByRole("button", { name: "Delete" }).click();
	await page.getByRole("dialog", { name: "Confirm action" }).getByRole("button", { name: "Confirm" }).click();
	await expect(page.getByRole("cell", { name: roleKey, exact: true })).toHaveCount(0);

  await page.screenshot({ path: "test-results/r4-schema-crud.png", fullPage: true });
});