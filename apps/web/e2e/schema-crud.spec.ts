import { expect, test, type Page } from "@playwright/test";

// A-010 R-004 · real browser Schema CRUD lifecycle against Go + SQLite.
// GOAL-011 S3 repoints the driver from the retired records page to the users
// resource page. Boots via playwright webServer (same as shell.spec.ts): Go API
// + Vite, then drives the users representative page through create → edit →
// delete with confirm. T-UI-01～10 cover Renderer behavior with an in-memory
// API emulator; this file proves the browser → proxy → Go/SQLite path.

async function signInAsAdmin(page: Page): Promise<void> {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill("admin");
  await page.getByRole("button", { name: "Sign in" }).click();
  await expect(page).toHaveURL(/\/overview$/);
}

test("users drives real create / edit / delete against Go SQLite", async ({
  page,
}) => {
  const stamp = Date.now().toString(36);
  const createdName = `E2E CRUD ${stamp}`;
  const editedName = `E2E CRUD edited ${stamp}`;

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
  await createDialog.getByLabel("Username").fill(createdName);
  await createDialog.getByLabel("Name", { exact: true }).fill("E2E Bot");
  await createDialog.getByLabel("Password").fill("e2e-password");
  await createDialog.getByRole("button", { name: "Create user" }).click();

  await expect(page.getByText("Item created")).toBeVisible();
  await expect(page.getByRole("cell", { name: createdName })).toBeVisible();

  // --- Edit (row-scoped) ---
  const createdRow = page.getByRole("row").filter({ hasText: createdName });
  await createdRow.getByRole("button", { name: "Edit" }).click();
  const editDialog = page.getByRole("dialog", { name: "Edit" });
  await expect(editDialog).toBeVisible();
  await editDialog.getByLabel("Name").fill(editedName);
  await editDialog.getByRole("button", { name: "Save changes" }).click();

  await expect(page.getByText("Item updated")).toBeVisible();
  await expect(page.getByRole("cell", { name: editedName })).toBeVisible();

  // --- Delete with confirm ---
  const editedRow = page.getByRole("row").filter({ hasText: editedName });
  await editedRow.getByRole("button", { name: "Delete" }).click();
  const confirmDialog = page.getByRole("dialog", { name: "Confirm action" });
  await expect(confirmDialog).toBeVisible();
  await expect(confirmDialog.getByText("Delete this user?")).toBeVisible();
  await confirmDialog.getByRole("button", { name: "Confirm" }).click();

  await expect(page.getByText("Item deleted")).toBeVisible();
  await expect(page.getByRole("cell", { name: editedName })).toHaveCount(0);

  await page.screenshot({ path: "test-results/r4-schema-crud.png", fullPage: true });
});
