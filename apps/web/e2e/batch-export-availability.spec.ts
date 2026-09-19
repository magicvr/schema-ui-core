import { expect, test } from "@playwright/test";

import { openSidebarGroup, signInAsAdmin } from "./sign-in";

// W34 (workspace-010 · GOAL-046) — the availability contract for the list pages'
// 「导出所选」 trigger, asserted in BOTH profiles.
//
// Why this spec exists: the 2026-09-19 user report was a trigger that rendered
// as if it worked while its route was not mounted at all, so every click ended
// in a bare 404 (NOT_FOUND / 「未找到」). The unit tests pin the component's
// decision, but only a browser run proves the decision survives the real
// assembly (schema node → slot host → permissions from /me → module set).
//
// The contract, and the reason it is profile-dependent:
//   - admin: `admin.jobs` + `admin.data-transfer` are in the profile, so both
//     gates the submit route requires are granted → the control is usable.
//   - mvp/demo: neither module is in the profile → the control must stay VISIBLE
//     but disabled and explained (W33 D-001 §3 fail-open: silently removing an
//     operation entry point is the worse failure mode, and an empty slot host
//     also breaks the page-actions height contract).
const appProfile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();
const adminCapable = appProfile === "admin" || appProfile === "custom";

test("the batch-export trigger never pretends to work when the profile cannot serve it", async ({
  page,
}) => {
  await signInAsAdmin(page);

  // Same navigation contract as the other list specs: sidebar, not page.goto()
  // (a raw reload drops the in-memory session and lands on the host failure).
  await openSidebarGroup(page, "manifest.nav.group.identityAccess");
  await page.getByRole("link", { name: "Users" }).click();

  const trigger = page.locator("[data-jobs-batch-export-submit]");
  await expect(trigger).toBeVisible({ timeout: 15000 });

  // W33 D-001 §3: the entry point must never disappear.
  await expect(page.locator("[data-jobs-batch-export]")).toHaveCount(1);

  // Select a row so the ONLY remaining reason to be disabled would be
  // availability (the empty-selection disable is pinned by unit tests).
  const firstRow = page.locator('tbody input[type="checkbox"]').first();
  await expect(firstRow).toBeVisible({ timeout: 15000 });
  await firstRow.click();

  if (adminCapable) {
    // Both gates granted: the control is genuinely usable, and carries no
    // "unavailable" marker.
    await expect(trigger).toBeEnabled();
    await expect(trigger).not.toHaveAttribute("data-jobs-batch-export-unavailable", "true");
    await expect(page.locator("[data-jobs-batch-export-unavailable-note]")).toHaveCount(0);
  } else {
    // Neither gate can be granted: the control explains itself instead of
    // failing at request time.
    await expect(trigger).toBeDisabled();
    await expect(trigger).toHaveAttribute("data-jobs-batch-export-unavailable", "true");
    await expect(page.locator("[data-jobs-batch-export-unavailable-note]")).toBeVisible();
  }
});
