import { expect, test } from "@playwright/test";

// S5 browser evidence (VP-007): the same shipped build, real Go API, and real
// manifest serve the localization + settings surfaces end-to-end:
// - M1: locale switch on the anonymous login page, HTML lang follows, login
//   feedback + post-login shell stay in the switched locale.
// - M3: a settings save projects to the shell header through the config
//   refresh event.
// - S4: the API negotiates error messages per Accept-Language (zh-CN).
// - F-003: when APP_PROFILE=mvp, settings edit surface is absent while locale
//   switch + branding bootstrap still work.

const appProfile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();

test("S5 localization: zh switch, lang, error negotiation, settings projection", async ({ page, request }) => {
  test.skip(appProfile !== "admin", "admin-only: settings edit surface requires admin profile");
  const pageErrors: string[] = [];
  page.on("pageerror", (error) => pageErrors.push(String(error)));

  // M1 · anonymous login page (en default) → switch to zh.
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  // 2026-08-14 refactor: the switcher is a header dropdown (no <select>).
  await page.getByRole("button", { name: "Language" }).click();
  await page.getByRole("menuitemradio", { name: "简体中文" }).click();
  await expect(page.getByRole("heading", { name: "登录" })).toBeVisible();
  expect(await page.evaluate(() => document.documentElement.lang)).toBe("zh-CN");

  // S4 · the API negotiates the login failure message in zh when asked.
  const failedLogin = await request.post("/api/auth/login", {
    headers: { "Accept-Language": "zh-CN" },
    data: { username: "admin", password: "wrong-password" },
  });
  expect(failedLogin.status()).toBe(401);
  expect(failedLogin.headers()["content-language"]).toBe("zh-CN");
  const failedBody = await failedLogin.json();
  expect(failedBody.error).toBe("UNAUTHORIZED");
  expect(failedBody.message).toBe("用户名或密码错误");
  expect(failedBody.messageKey).toBe("error.unauthorized");

  // M1 · login in zh: the same seed works and the shell stays zh.
  await page.getByLabel("用户名").fill("admin");
  await page.getByLabel("密码").fill("admin");
  await page.getByRole("button", { name: "登录" }).click();
  await expect(page).toHaveURL(/\/dashboard$/);
  await expect(page.getByRole("heading", { name: "仪表盘" })).toBeVisible();
  expect(await page.evaluate(() => document.documentElement.lang)).toBe("zh-CN");

  // M3 · admin settings: edit the inline General form, save, and the shell
  // header + document title project the new value (config refresh).
  // W12 T-01: Settings is a user-menu item, not a sidebar link.
  await page.getByRole("button", { name: "用户菜单" }).click();
  await page.getByRole("menuitem", { name: "设置" }).click();
  await expect(page.getByRole("heading", { name: "设置" })).toBeVisible();
  const generalForm = page.locator("form").filter({ has: page.getByLabel("站点标题") });
  await generalForm.getByLabel("站点标题").fill("Acme 管理台");
  await generalForm.getByRole("button", { name: "保存设置" }).click();
  // W13 T-02: the brand text appears twice in the DOM (mobile-only brand bar
  // hidden at the desktop viewport + the desktop single-row brand link) —
  // target the LAST occurrence.
  await expect(page.getByText("Acme 管理台").last()).toBeVisible();
  expect(await page.title()).toBe("Acme 管理台");

  // W13 T-01: the settings surface now switches by functional unit through
  // tabs (same shape as the account page); Restore defaults stays a button
  // outside the tabs, reachable from any tab.
  await expect(page.getByRole("tab", { name: "常规" })).toBeVisible();
  await page.getByRole("tab", { name: "品牌" }).click();
  await expect(page.getByRole("heading", { name: "品牌" })).toBeVisible();
  await page.getByRole("tab", { name: "本地化" }).click();
  await expect(page.getByRole("heading", { name: "本地化" })).toBeVisible();
  await page.getByRole("tab", { name: "外观" }).click();
  await expect(page.getByRole("heading", { name: "外观" })).toBeVisible();
  await expect(page.getByRole("button", { name: "恢复默认" })).toBeVisible();

  await page.screenshot({ path: "test-results/s5-settings-zh.png", fullPage: true });

  expect(pageErrors).toEqual([]);
});

test("S5 mvp profile: locale switch + no settings surface + branding public", async ({ page, request }) => {
  test.skip(appProfile !== "mvp", "mvp-only boundary evidence");

  const pageErrors: string[] = [];
  page.on("pageerror", (error) => pageErrors.push(String(error)));

  // Public branding available without admin.settings module (exit 4).
  const branding = await request.get("/api/branding");
  expect(branding.status()).toBe(200);
  const body = await branding.json();
  expect(body.siteTitle).toBeTruthy();
  expect(body.supportedLocales).toEqual(expect.arrayContaining(["zh-CN", "en-US"]));

  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  // 2026-08-14 refactor: the switcher is a header dropdown (no <select>).
  await page.getByRole("button", { name: "Language" }).click();
  await page.getByRole("menuitemradio", { name: "简体中文" }).click();
  await expect(page.getByRole("heading", { name: "登录" })).toBeVisible();
  expect(await page.evaluate(() => document.documentElement.lang)).toBe("zh-CN");

  await page.getByLabel("用户名").fill("admin");
  await page.getByLabel("密码").fill("admin");
  await page.getByRole("button", { name: "登录" }).click();
  await expect(page).toHaveURL(/\/dashboard$/);
  await expect(page.getByRole("heading", { name: "仪表盘" })).toBeVisible();

  // Settings edit surface must not appear under mvp.
  await expect(page.getByRole("link", { name: "设置" })).toHaveCount(0);
  await page.goto("/settings");
  await expect(page.getByText(/找不到|not found|Page not found|页面/i).first()).toBeVisible();

  await page.screenshot({ path: "test-results/s5-mvp-dashboard-zh.png", fullPage: true });
  expect(pageErrors).toEqual([]);
});