#!/usr/bin/env node
/**
 * publish-npmjs-packages.mjs — npmjs.com 公开发布（VP-024 R2 · 判据 #2）。
 *
 * 前置：仓库根 `.env` 含 `npm_token`（.gitignore 已排除，不入库）；
 *       `apps/web/dist-lib/artifacts/` 六份 tgz 就绪（pack-npm-packages.mjs）。
 *
 * 用法：
 *   node scripts/publish-npmjs-packages.mjs [--dry-run]
 *
 * 流程：读取 .env npm_token → 临时 .npmrc（仅 stage 临时目录，随 stage 删除）→
 *       逐包解包改名为 <PUBLISH_SCOPE>/schema-ui-*（默认 @magicvr · npmjs 公开）
 *       → npm publish --registry https://registry.npmjs.org --access public
 *       （已发布版本自动跳过：npm view versions 预检 + 403「cannot publish over」兜底）。
 */
import { execFileSync } from "node:child_process";
import { readdirSync, readFileSync, writeFileSync, mkdtempSync, mkdirSync, cpSync, existsSync, rmSync } from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(here, "..");
const artifacts = path.join(root, "apps/web/dist-lib/artifacts");
const registry = "https://registry.npmjs.org";
// 实发 scope = @magicvr（D-001 §6 · 用户裁决：npmjs 公开 @magicvr 先行；
// @schema-ui 需同名 org，为 org 就绪后的正式化候选，届时覆写本默认值）。
const scope = process.env.PUBLISH_SCOPE || "@magicvr";
const dryRun = process.argv.includes("--dry-run");

// 从仓库根 .env 读取 npm_token（值不打印、不落盘）。
function readEnvToken() {
  const p = path.join(root, ".env");
  if (!existsSync(p)) return "";
  for (const line of readFileSync(p, "utf8").split(/\r?\n/)) {
    const t = line.trim();
    if (!t || t.startsWith("#")) continue;
    const eq = t.indexOf("=");
    if (eq < 0) continue;
    if (t.slice(0, eq).trim() === "npm_token") return t.slice(eq + 1).trim();
  }
  return "";
}

if (!existsSync(artifacts)) {
  console.error("artifacts 不存在：先跑 node scripts/pack-npm-packages.mjs");
  process.exit(1);
}
const token = readEnvToken();
if (!token && !dryRun) {
  console.error("缺少 npm_token：请在本仓根 .env 配置 npm_token（不入库）；或先 --dry-run");
  process.exit(1);
}

const tgzs = readdirSync(artifacts).filter((f) => f.endsWith(".tgz"));
if (tgzs.length === 0) {
  console.error("artifacts 无 tgz");
  process.exit(1);
}

const stage = mkdtempSync(path.join(os.tmpdir(), "gf-npmjs-"));
try {
  if (token) {
    writeFileSync(path.join(stage, ".npmrc"), `//registry.npmjs.org/:_authToken=${token}\n`);
  }
  for (const tgz of tgzs) {
    const base = path.basename(tgz, ".tgz"); // magicvr-schema-ui-protocol-0.2.2 | schema-ui-protocol-0.2.1
    const dir = path.join(stage, base);
    cpSync(path.join(artifacts, tgz), path.join(stage, `${base}.tgz`));
    mkdirSync(dir, { recursive: true });
    execFileSync("tar", ["-xzf", path.join(stage, `${base}.tgz`), "-C", dir], { shell: process.platform === "win32" });
    const pkgPath = path.join(dir, "package/package.json");
    const pkg = JSON.parse(readFileSync(pkgPath, "utf8"));
    // 包名 = scope + 解包 json 的 name 尾段（兼容 pack 出的两种命名：schema-ui-* 与 magicvr-schema-ui-*）。
    const namePart = (pkg.name || base).includes("/") ? (pkg.name || base).split("/").pop() : (pkg.name || base);
    const pkgName = `${scope}/${namePart}`;
    const version = pkg.version || (() => { const m = base.match(/-(\d+\.\d+\.\d+)$/); return m ? m[1] : ""; })();
    pkg.name = pkgName;
    // scoped 包默认 private：显式公开（免费账号私有包发布会被 E402 拒绝），
    // 并写 publishConfig 使后续发布默认公开。
    pkg.publishConfig = { access: "public" };
    writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));

    // 版本存在性检查：已发布则跳过（registry 禁止覆盖同一版本）。
    // npmjs CDN 传播可能有秒级延迟——404 不可靠；以 publish 403 判定兜底。
    let versions = "";
    try {
      versions = execFileSync(
        "npm", ["view", pkgName, "versions", "--json", "--registry", registry],
        { cwd: path.join(dir, "package"), encoding: "utf8", shell: process.platform === "win32" },
      ).trim();
    } catch {
      versions = "";
    }
    if (versions.includes(`"${version}"`)) {
      console.log(`skip ${pkgName}@${version}（已发布）`);
      continue;
    }
    console.log(`publishing ${pkgName}@${version}${dryRun ? " (dry-run)" : ""} ...`);
    const args = ["publish", "--registry", registry, "--access", "public"];
    if (dryRun) args.push("--dry-run");
    if (token) args.push("--userconfig", path.join(stage, ".npmrc"));
    try {
      execFileSync("npm", args, {
        cwd: path.join(dir, "package"), stdio: ["ignore", "inherit", "pipe"], encoding: "utf8",
        shell: process.platform === "win32",
      });
    } catch (err) {
      const stderr = String(err?.stderr ?? "");
      if (stderr.includes("cannot publish over the previously published versions")) {
        console.log(`skip ${pkgName}@${version}（已发布 · 403 兜底判定）`);
        continue;
      }
      throw err;
    }
    console.log(`published ${pkgName}@${version}${dryRun ? " (dry-run)" : ""}`);
  }
} finally {
  rmSync(stage, { recursive: true, force: true });
}