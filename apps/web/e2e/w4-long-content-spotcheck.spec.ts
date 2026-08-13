import { expect, test } from "@playwright/test";

// W4 · GOAL-005 (workspace-010) browser spot-check — answers A-003 F-3:
// the roles list must no longer let long permissions/menuItems values crowd
// out sibling columns, and the recordView drawer must wrap long values
// without horizontal scrolling. Runs under APP_PROFILE=admin (roles page
// with full seed grants). jsdom tests assert classes/attributes only; this
// spec asserts real table layout in a desktop viewport.

test.use({ viewport: { width: 1440, height: 900 } });

test("roles list truncates long columns and the detail drawer wraps", async ({ page }) => {
  // Sign in with the dev seed (same flow as shell.spec.ts).
  await page.goto("/");
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill("admin");
  await page.getByRole("button", { name: "Sign in" }).click();
  await expect(page.getByRole("link", { name: "Roles" })).toBeVisible();

  await page.getByRole("link", { name: "Roles" }).click();
  await expect(page.getByRole("heading", { name: "Roles" })).toBeVisible();

  // Long-value columns (permissions/menuItems) ship the truncate affordance.
  const truncated = page.locator('[data-table-cell="truncated"]');
  await expect(truncated.first()).toBeVisible();
  const count = await truncated.count();
  expect(count).toBeGreaterThanOrEqual(2);

  // Every truncated cell is capped at 16rem (256px) and carries the full
  // text as its native title affordance.
  for (let i = 0; i < count; i++) {
    const cell = truncated.nth(i);
    const width = await cell.evaluate((el) => el.getBoundingClientRect().width);
    expect(width, `truncated cell #${i} width`).toBeLessThanOrEqual(257);
    const title = await cell.getAttribute("title");
    const text = (await cell.textContent()) ?? "";
    if (text.length > 0) {
      expect(title, `truncated cell #${i} title`).toBe(text);
    }
  }

  // The long columns are capped in real auto table layout and the sibling
  // columns stay readable — the pre-fix failure mode was permissions/menus
  // blowing out to their full-string width and squeezing ID/Key/Name.
  const desktop = page.locator('[data-table-presentation="desktop-table"]');
  await expect(desktop).toBeVisible();
  const colWidths = await desktop.evaluate((el) => {
    const out: Record<string, number> = {};
    for (const th of Array.from(el.querySelectorAll("th"))) {
      out[th.textContent ?? ""] = th.getBoundingClientRect().width;
    }
    return out;
  });
  expect(colWidths["Permissions"] ?? 0).toBeLessThanOrEqual(257);
  expect(colWidths["Menus"] ?? 0).toBeLessThanOrEqual(257);
  for (const name of ["ID", "Key", "Name"]) {
    expect(colWidths[name] ?? 0).toBeGreaterThanOrEqual(40);
  }

  // Open the detail drawer and verify long values wrap instead of
  // overflowing horizontally.
  const row = desktop.locator("tbody tr").first();
  await row.click();
  const panel = page.locator('[data-record-view="panel"][data-record-view-mode="drawer"]');
  await expect(panel).toBeVisible();
  const dl = panel.locator("dl");
  const noHorizontalOverflow = await dl.evaluate((el) => el.scrollWidth <= el.clientWidth + 1);
  expect(noHorizontalOverflow).toBe(true);

  // Array values render as a wrapped, space-joined list (not the raw
  // `a,b,c` single run).
  const permissionText = await panel.locator("dd").filter({ hasText: "users.read" }).textContent();
  expect(permissionText ?? "").toMatch(/users\.read, /);

  await page.screenshot({ path: "test-results/w4-roles-truncate-detail.png", fullPage: true });
});
