#!/usr/bin/env node
/**
 * publish-npm-packages.mjs — 包分发试点「registry 发布」（R1 S2 · VP-023 判据 #1 npm 侧）。
 *
 * 前置（凭据）：GH_TOKEN（scope write:packages）或 npm 认证配置；目标 =
 * GitHub Packages（https://npm.pkg.github.com，包名 scope 须 = owner `@magicvr`）。
 *
 * 用法：
 *   set GH_TOKEN=ghp_xxx        （或已配置 npm auth）
 *   node scripts/publish-npm-packages.mjs [--dry-run]
 *
 * 流程：pack-npm-packages.mjs（tgz）→ 逐包改名 @magicvr/schema-ui-* →
 * npm publish --registry https://npm.pkg.github.com。
 */
import { execFileSync } from "node:child_process";
import { readdirSync, readFileSync, writeFileSync, mkdtempSync, mkdirSync, cpSync, existsSync, rmSync } from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(here, "..");
const artifacts = path.join(root, "apps/web/dist-lib/artifacts");
const registry = process.env.NPM_REGISTRY || "https://npm.pkg.github.com";
const scope = process.env.PUBLISH_SCOPE || "@magicvr";
const dryRun = process.argv.includes("--dry-run");

if (!existsSync(artifacts)) {
  console.error("artifacts 不存在：先跑 node scripts/pack-npm-packages.mjs");
  process.exit(1);
}
if (!process.env.GH_TOKEN && !dryRun) {
  console.error("缺少发布凭据：设置 GH_TOKEN（write:packages）或先 --dry-run 验证流程");
  process.exit(1);
}

const tgzs = readdirSync(artifacts).filter((f) => f.endsWith(".tgz"));
if (tgzs.length === 0) {
  console.error("artifacts 无 tgz");
  process.exit(1);
}

const stage = mkdtempSync(path.join(os.tmpdir(), "gf-publish-"));
try {
  for (const tgz of tgzs) {
    const base = path.basename(tgz, ".tgz"); // schema-ui-protocol-0.2.0
    const [name, version] = base.match(/^(.*)-(\d+\.\d+\.\d+)$/).slice(1);
    const pkgName = `${scope}/${name}`;
    const dir = path.join(stage, base);
    cpSync(path.join(artifacts, tgz), path.join(stage, `${base}.tgz`));
    // 解包改名（npm publish 以 package.json name 为准的发布侧改名为最稳路径：
    // 本地重新解包改名后 pack 再 publish）。
    mkdirSync(dir, { recursive: true });
    execFileSync("tar", ["-xzf", path.join(stage, `${base}.tgz`), "-C", dir], { shell: process.platform === "win32" });
    const pkgPath = path.join(dir, "package/package.json");
    const pkg = JSON.parse(readFileSync(pkgPath, "utf8"));
    pkg.name = pkgName;
    writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));
    // GH Packages 认证：临时 .npmrc（token 仅存在于 stage 临时目录，成功后随 stage 删除）
    if (process.env.GH_TOKEN) {
      writeFileSync(
        path.join(dir, "package/.npmrc"),
        `//npm.pkg.github.com/:_authToken=${process.env.GH_TOKEN}\n${scope}:registry=${registry}\n`,
      );
    }
    const args = ["publish", "--registry", registry];
    if (dryRun) args.push("--dry-run");
    // 版本存在性检查：已发布则跳过（registry 禁止覆盖同一版本）；404 = 未发布 → 继续。
    let published = "";
    try {
      published = execFileSync(
        "npm",
        ["view", `${pkgName}@${version}`, "version", "--registry", registry],
        { cwd: path.join(dir, "package"), encoding: "utf8", shell: process.platform === "win32" },
      ).trim();
    } catch {
      published = "";
    }
    if (published.length > 0) {
      console.log(`skip ${pkgName}@${version}（已发布）`);
      continue;
    }
    execFileSync("npm", args, { cwd: path.join(dir, "package"), stdio: "inherit", shell: process.platform === "win32" });
    console.log(`published ${pkgName}@${version}${dryRun ? " (dry-run)" : ""}`);
  }
} finally {
  rmSync(stage, { recursive: true, force: true });
}

function mkdirSync2(p) {
  const { mkdirSync } = require("fs");
  mkdirSync(p, { recursive: true });
}