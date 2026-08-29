#!/usr/bin/env node
/**
 * pack-npm-packages.mjs — Web 包分发试点「发布可复现」（GOAL-006 S1 / 判据 #5）。
 *
 * 职责：把 apps/web/dist-lib/@schema-ui/* 产物打为 npm tarball（.tgz），
 * 输出到 dist-lib/artifacts/。tarball = registry 安装的真实载荷：
 * golden-web 以 `pnpm add <tarball>` 安装（与 `pnpm add @schema-ui/protocol@0.2.0`
 * 语义一致，仅来源为本地文件）。
 *
 * 用法：node scripts/pack-npm-packages.mjs
 */
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync, copyFileSync, existsSync, rmSync } from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));
const webRoot = path.resolve(here, "../apps/web");
const distLib = path.join(webRoot, "dist-lib");
const artifacts = path.join(distLib, "artifacts");

if (!existsSync(distLib)) {
  console.error("dist-lib 不存在：先跑 apps/web 的 vite lib 构建（protocol/renderer）");
  process.exit(1);
}

mkdirSync(artifacts, { recursive: true });
rmSync(path.join(artifacts, "*"), { recursive: true, force: true });

const packages = readdirSync(distLib, { withFileTypes: true })
  .filter((e) => e.isDirectory() && e.name.startsWith("@"))
  .flatMap((e) => {
    const pkgDir = path.join(distLib, e.name);
    return readdirSync(pkgDir, { withFileTypes: true })
      .filter((n) => n.isDirectory())
      .map((n) => ({ parent: e.name, dir: path.join(pkgDir, n.name) }));
  });

const built = [];
for (const { dir } of packages) {
  // pnpm pack 在包目录内执行，产出 <name>-<version>.tgz，再移入 artifacts。
  const cwd = dir;
  const out = execFileSync("pnpm", ["pack", "--pack-destination", artifacts], {
    cwd,
    encoding: "utf8",
    shell: process.platform === "win32",
  });
  const tgzName = out.trim().split(/\r?\n/).pop().trim();
  built.push({ pkg: path.basename(dir), tgz: tgzName, path: path.join(artifacts, tgzName) });
}

console.table(built.map((b) => ({ artifact: b.tgz, source: b.pkg })));
console.log("artifacts →", artifacts);