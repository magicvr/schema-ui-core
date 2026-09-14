import { expect, test } from "@playwright/test";

import { signInAsAdmin } from "./sign-in";

test.describe("VP-036 command palette", () => {
  test("searches visible pages and opens a declared page action", async ({ page }) => {
    await signInAsAdmin(page);

    const trigger = page.locator("[data-command-palette-trigger]");
    await expect(trigger).toBeVisible();

    await page.keyboard.press("Control+K");
    const palette = page.locator('[role="dialog"]');
    const input = palette.getByRole("combobox");
    await expect(input).toBeFocused();

    await input.fill("Users");
    const usersPage = palette.locator('[role="option"]').filter({ hasText: "Users" }).first();
    await expect(usersPage).toBeVisible();
    await usersPage.click();
    await expect(page).toHaveURL(/\/users$/);
    await expect(page.getByRole("heading", { name: "Users" })).toBeVisible();

    await trigger.click();
    await expect(input).toBeFocused();
    await input.fill("New user");
    const createUser = palette.locator('[role="option"]').filter({ hasText: "New user" }).first();
    await expect(createUser).toBeVisible();
    await createUser.click();

    // The action is handed to the owner page's existing SchemaCrudProvider;
    // the modal title and form are therefore the same as the page toolbar path.
    await expect(page.getByRole("dialog", { name: "New user" })).toBeVisible();
    await expect(page.getByRole("dialog", { name: "New user" }).getByLabel("Username")).toBeVisible();
  });

  test("keeps the trigger available on the mobile functional row", async ({ page }) => {
    await signInAsAdmin(page);
    await page.setViewportSize({ width: 390, height: 844 });
    const trigger = page.locator("[data-command-palette-trigger]");
    await expect(trigger).toBeVisible();
    await trigger.click();
    await expect(page.getByRole("combobox")).toBeFocused();
    await page.keyboard.press("Escape");
    await expect(page.locator('[role="dialog"]')).toHaveCount(0);
  });
});
