import { expect, test, type Page } from "@playwright/test";

// A-010 R-004 · real browser Schema CRUD lifecycle against Go + SQLite.
// Boots via playwright webServer (same as shell.spec.ts): Go API + Vite, then
// drives the list-edit-lifecycle representative page through create → edit →
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

test("list-edit-lifecycle drives real create / edit / delete against Go SQLite", async ({
  page,
}) => {
  const stamp = Date.now().toString(36);
  const createdName = `E2E CRUD ${stamp}`;
  const editedName = `E2E CRUD edited ${stamp}`;

  await signInAsAdmin(page);

  // Menu projection (R3): admin seed grants menu_list_edit_lifecycle; login
  // must resolve /me features so the link is present after sign-in.
  await expect(page.getByRole("link", { name: "List + edit" })).toBeVisible();
  await page.getByRole("link", { name: "List + edit" }).click();
  await expect(page).toHaveURL(/\/list-edit-lifecycle$/);
  await expect(page.getByRole("heading", { name: "List + edit lifecycle" })).toBeVisible();
  await expect(page.getByRole("button", { name: "New record" })).toBeEnabled();

  // --- Create ---
  await page.getByRole("button", { name: "New record" }).click();
  const createDialog = page.getByRole("dialog", { name: "New record" });
  await expect(createDialog).toBeVisible();
  await createDialog.getByLabel("Name").fill(createdName);
  await createDialog.getByLabel("Status").selectOption("active");
  await createDialog.getByLabel("Owner").fill("e2e-bot");
  await createDialog.getByRole("button", { name: "Create record" }).click();

  await expect(page.getByText("Record created")).toBeVisible();
  await expect(page.getByRole("cell", { name: createdName })).toBeVisible();

  // --- Edit (row-scoped) ---
  const createdRow = page.getByRole("row").filter({ hasText: createdName });
  await createdRow.getByRole("button", { name: "Edit" }).click();
  const editDialog = page.getByRole("dialog", { name: "Edit" });
  await expect(editDialog).toBeVisible();
  await expect(editDialog.getByLabel("Name")).toHaveValue(createdName);
  await editDialog.getByLabel("Name").fill(editedName);
  await editDialog.getByRole("button", { name: "Save changes" }).click();

  await expect(page.getByText("Record updated")).toBeVisible();
  await expect(page.getByRole("cell", { name: editedName })).toBeVisible();
  await expect(page.getByRole("cell", { name: createdName })).toHaveCount(0);

  // --- Delete with confirm ---
  const editedRow = page.getByRole("row").filter({ hasText: editedName });
  await editedRow.getByRole("button", { name: "Delete" }).click();
  const confirmDialog = page.getByRole("dialog", { name: "Confirm action" });
  await expect(confirmDialog).toBeVisible();
  await expect(confirmDialog.getByText("Delete this record?")).toBeVisible();
  await confirmDialog.getByRole("button", { name: "Confirm" }).click();

  await expect(page.getByText("Record deleted")).toBeVisible();
  await expect(page.getByRole("cell", { name: editedName })).toHaveCount(0);

  await page.screenshot({ path: "test-results/r4-schema-crud.png", fullPage: true });
});
