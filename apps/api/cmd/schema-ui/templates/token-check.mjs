// Token 覆盖纪律断言（复刻主仓 brand-example.test.ts 语义）：
// brand.css 只允许覆盖 index.css 已声明的 CSS 变量；新增变量 = 违规（防 token 体系漂移）。
import { readFileSync } from "node:fs";
import assert from "node:assert";
import { fileURLToPath } from "node:url";
import path from "node:path";

const here = path.dirname(fileURLToPath(import.meta.url));
const indexCss = readFileSync(path.join(here, "index.css"), "utf8");
const brandCss = readFileSync(path.join(here, "brand.css"), "utf8");

const tokens = (css) => [...css.matchAll(/--[a-z0-9-]+(?=\s*:)/gi)].map((m) => m[0]);
const indexTokens = new Set(tokens(indexCss));
const brandTokens = tokens(brandCss);

const violations = brandTokens.filter((t) => !indexTokens.has(t));
assert.equal(violations.length, 0, `brand.css 新增未声明 Token: ${violations.join(", ")}`);
assert.ok(brandTokens.length >= 1, "brand.css 应有覆盖示例");

console.log(`golden-web token override PASS · brand=${brandTokens.length} tokens ⊆ index=${indexTokens.size}`);