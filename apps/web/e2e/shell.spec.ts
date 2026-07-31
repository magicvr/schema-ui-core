import { expect, test } from "@playwright/test";

// R6 minimal browser E2E (I-008-005 / I-008-003 forward path).
// Boots the Go API and Vite dev server via playwright webServer; verifies the
// shell renders from the real manifest and that the account context flows
// through the Web /api proxy to the Go API (dev session) in a real browser.

test("shell renders from the real manifest and proxy serves the dev session", async ({ page }) => {
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

  // Cross-layer forward path: account context via Web /api proxy -> Go API.
  const res = await page.request.get("/api/accounts/me");
  expect(res.status()).toBe(200);
  const session = await res.json();
  expect(session.user.id).toBe("dev-001");
  expect(session.user.roles).toEqual(expect.arrayContaining(["admin", "editor"]));

  // Records example API reachable through the same proxy (D-DATA forward path).
  const records = await page.request.get("/api/records");
  expect(records.status()).toBe(200);
  const list = await records.json();
  expect(list.items.length).toBeGreaterThan(0);

  await page.screenshot({ path: "test-results/r6-overview.png", fullPage: true });
});
