import { expect, test } from "@playwright/test";

// S5 browser evidence (VP-007): the same shipped build, real Go API, and real
// manifest serve the localization + settings surfaces end-to-end:
// - M1: locale switch on the anonymous login page, HTML lang follows, login
//   feedback + post-login shell stay in the switched locale.
// - M3: a settings save projects to the shell header through the config
//   refresh event.
// - S4: the API negotiates error messages per Accept-Language (zh-CN).

test("S5 localization: zh switch, lang, error negotiation, settings projection", async ({ page, request }) => {
  const pageErrors: string[] = [];
  page.on("pageerror", (error) => pageErrors.push(String(error)));

  // M1 · anonymous login page (en default) → switch to zh.
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  await page.getByLabel("Language").selectOption("zh-CN");
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
  await expect(page).toHaveURL(/\/overview$/);
  await expect(page.getByRole("heading", { name: "总览" })).toBeVisible();
  expect(await page.evaluate(() => document.documentElement.lang)).toBe("zh-CN");

  // M3 · admin settings: open General, change the site title, save, and the
  // shell header + document title project the new value (config refresh).
  await page.getByRole("link", { name: "设置" }).click();
  await expect(page.getByRole("heading", { name: "设置" })).toBeVisible();
  await page.getByRole("button", { name: "常规" }).click();
  await page.getByLabel("站点标题").fill("Acme 管理台");
  await page.getByRole("button", { name: "保存设置" }).click();
  await expect(page.getByText("Acme 管理台").first()).toBeVisible();
  expect(await page.title()).toBe("Acme 管理台");

  // The four-category surface renders (toolbar labels in zh).
  await expect(page.getByRole("button", { name: "品牌" })).toBeVisible();
  await expect(page.getByRole("button", { name: "本地化" })).toBeVisible();
  await expect(page.getByRole("button", { name: "外观" })).toBeVisible();
  await expect(page.getByRole("button", { name: "恢复默认" })).toBeVisible();

  await page.screenshot({ path: "test-results/s5-settings-zh.png", fullPage: true });

  expect(pageErrors).toEqual([]);
});
