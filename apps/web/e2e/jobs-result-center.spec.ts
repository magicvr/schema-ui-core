import { expect, test } from "@playwright/test";

import { openSidebarGroup, signInAsAdmin } from "./sign-in";

// R5 (workspace-038 · GOAL-006) — the browser end-to-end path the VP-038 exit
// matrix was missing: submit a real async batch export, watch it progress, take
// the CSV result, and then read the SAME job from the result center where the
// row's own download action serves it again.
//
// Why the profile guard: `admin.jobs` is only in the admin default set, and the
// harness defaults to APP_PROFILE=mvp. Without this guard the spec would fail on
// a missing page in every mvp run — and silently passing would be worse, because
// it would prove nothing about the module. Run it with:
//
//   APP_PROFILE=admin npx playwright test jobs-result-center.spec.ts
const appProfile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();
test.skip(appProfile !== "admin", "admin-only: the result center requires admin.jobs in the profile");

test("batch export progresses to a downloadable result and the result center serves it", async ({ page }) => {
  await signInAsAdmin(page);

  // 1. Select real rows on the users page (the async trigger lives there).
  //    Navigation goes through the sidebar: a raw page.goto() reloads the SPA
  //    without the in-memory session and lands on the "Session expired" host
  //    failure (the suite's other specs navigate the same way).
  await openSidebarGroup(page, "manifest.nav.group.identityAccess");
  await page.getByRole("link", { name: "Users" }).click();
  const checkboxes = page.locator('tbody input[type="checkbox"]');
  await expect(checkboxes.first()).toBeVisible({ timeout: 15000 });
  const rows = await checkboxes.count();
  expect(rows, "the seeded admin profile must expose at least one user row").toBeGreaterThan(0);
  const selected = Math.min(rows, 2);
  for (let index = 0; index < selected; index += 1) {
    await checkboxes.nth(index).click();
  }

  // 2. Submit. The control is disabled until the selection is published.
  const submit = page.locator("[data-jobs-batch-export-submit]");
  await expect(submit).toBeEnabled();
  await submit.click();

  // 3. Real progress, then a terminal success. The component polls the R2 read
  //    route (1s for the first ticks), so this is the actual job lifecycle.
  await expect(page.locator("[data-jobs-batch-export-progress]")).toBeVisible({ timeout: 15000 });
  const downloadButton = page.locator("[data-jobs-batch-export-download]");
  await expect(downloadButton).toBeVisible({ timeout: 45000 });

  const csvDownload = await Promise.all([page.waitForEvent("download"), downloadButton.click()]).then(
    ([download]) => download,
  );
  // The filename is the SERVER's own `fileName` field, not a client guess.
  expect(csvDownload.suggestedFilename()).toBe("users-selection.csv");

  // 4. The result center lists the same job with a LOCALIZED terminal state.
  await openSidebarGroup(page, "manifest.nav.group.operations");
  await page.getByRole("link", { name: "Jobs" }).click();
  const exportRow = page.locator("tbody tr", { hasText: "jobs.batch-export" }).first();
  await expect(exportRow).toBeVisible({ timeout: 15000 });
  await expect(exportRow).toContainText("Succeeded");

  // 5. …and its own row action downloads the same artefact. Download is the
  //    third action, so it lives behind the overflow menu.
  await exportRow.locator("[data-row-actions-menu] button").click();
  const downloadItem = page.getByRole("menuitem", { name: "Download result" });
  await expect(downloadItem).toBeEnabled();
  const rowDownload = await Promise.all([page.waitForEvent("download"), downloadItem.click()]).then(
    ([download]) => download,
  );
  expect(rowDownload.suggestedFilename()).toBe("users-selection.csv");
});
