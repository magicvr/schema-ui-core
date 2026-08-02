// S3 复现终点 4 · 浏览器登录 → /list-edit-lifecycle → 列表加载 Acme Console
// 用法: node e2e-repro-endpoint4.mjs <webBaseUrl> <username> <password>
import { chromium } from "playwright";

const [webBaseUrl, username, password] = process.argv.slice(2);
if (!webBaseUrl || !username || !password) {
  console.error("usage: node e2e-repro-endpoint4.mjs <webBaseUrl> <user> <pass>");
  process.exit(2);
}

const browser = await chromium.launch();
const page = await browser.newPage();
try {
  await page.goto(webBaseUrl + "/", { waitUntil: "networkidle" });
  await page.getByRole("heading", { name: "Sign in" }).waitFor({ timeout: 15000 });
  await page.getByLabel("Username").fill(username);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: "Sign in" }).click();
  await page.waitForURL(/\/overview$/, { timeout: 15000 });
  await page.getByRole("link", { name: "List + edit" }).click();
  await page.waitForURL(/\/list-edit-lifecycle$/, { timeout: 15000 });
  await page.getByRole("heading", { name: "List + edit lifecycle" }).waitFor({ timeout: 15000 });
  await page.getByRole("cell", { name: "Acme Console" }).waitFor({ timeout: 15000 });
  console.log("ENDPOINT4=PASS title=list-edit-lifecycle cell=Acme Console");
  await page.screenshot({ path: "test-results/r5-repro-endpoint4.png", fullPage: true });
} finally {
  await browser.close();
}
