#!/usr/bin/env node
/**
 * rewrite-lib-aliases.mjs — 六包 d.ts alias 重写 + package.json 定稿（R5 · VP-024 判据 #5/#6）。
 *
 * 背景：B 路径产物 d.ts 保留 `@/` 别名 import，消费端 TS 无法解析（renderer 0.2.0 类型面实际不可消费）。
 * R5 把包子路径化映射应用到各包 .d.ts；package.json 统一 exports（"." + "./*" 子路径）、
 * files 收窄（index.js + 实际目录）、版本推进、renderer peer 矩阵定稿。
 *
 * 用法：node scripts/rewrite-lib-aliases.mjs（在 build-lib-packages.mjs 之后执行）
 */
import { readdirSync, readFileSync, writeFileSync, existsSync } from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(here, "..");
const distRoot = path.join(root, "apps/web/dist-lib/@schema-ui");

// 映射表：src 面 → 包子路径（D-001-r5-granularity §1；发布实态全名 @magicvr/schema-ui-*）
// 注：tsc 产物按 rootDir=src 镜像（src/protocol → outcprotocol/），包子路径须带对应段（protocol/）。
const faces = {
  "@/i18n/": "@magicvr/schema-ui-lib/i18n/",
  "@/protocol/": "@magicvr/schema-ui-protocol/protocol/",
  "@/lib/": "@magicvr/schema-ui-lib/lib/",
  "@/components/data-table": "@magicvr/schema-ui-renderer/components/data-table",
  "@/components/ui/": "@magicvr/schema-ui-ui/components/ui/",
  "@/theme/": "@magicvr/schema-ui-theme/",
  "@/app/": "@magicvr/schema-ui-shell/app/",
  "@/renderer/": "@magicvr/schema-ui-renderer/renderer/",
};

// 包终名 = 发布实态全名（@magicvr/schema-ui-<pkg>）；versions/peers 以包短名为键
const versions = {
  "renderer": "0.3.4",
  "protocol": "0.2.3",
  "lib": "0.1.3",
  "ui": "0.1.3",
  "theme": "0.1.2",
  "shell": "0.1.2",
};

const peers = {
  "renderer": { react: "^19.0.0", "react-dom": "^19.0.0", "@magicvr/schema-ui-protocol": "^0.2.3", "@magicvr/schema-ui-lib": "^0.1.3", "@magicvr/schema-ui-ui": "^0.1.3" },
  "ui": { react: "^19.0.0", "react-dom": "^19.0.0" },
  "shell": { react: "^19.0.0", "react-dom": "^19.0.0" },
};

// 无法映射面（无对应包）计数（shell 的 host/account 面 → 残余登记）
const unmapPattern = /from\s+["']@\/(?!i18n\/|protocol\/|lib\/|components\/ui\/|theme\/|app\/|renderer\/)/g;

let unmapTotal = 0;
for (const pkg of readdirSync(distRoot)) {
  const dir = path.join(distRoot, pkg);
  const pkgPath = path.join(dir, "package.json");
  const old = existsSync(pkgPath) ? JSON.parse(readFileSync(pkgPath, "utf8")) : {};
  const topDirs = readdirSync(dir, { withFileTypes: true })
    .filter((e) => e.isDirectory())
    .map((e) => e.name);
  const files = pkg === "renderer" ? ["index.js", "renderer", "components"] : ["index.js", ...topDirs];
  const defaultTypes = topDirs.length > 0 ? `./${topDirs[0]}/index.d.ts` : "./index.d.ts";
  const types = old.types || defaultTypes; // 缺 package.json 包（renderer clean 后）按默认模板生成
  const rewritten = [];
  const walk = (d) => {
    for (const ent of readdirSync(d, { withFileTypes: true })) {
      const full = path.join(d, ent.name);
      if (ent.isDirectory()) walk(full);
      else if (ent.name.endsWith(".d.ts") || ent.name.endsWith(".js")) {
        const text = readFileSync(full, "utf8");
        let out = text;
        for (const [k, v] of Object.entries(faces)) out = out.split(k).join(v);
        // ESM 扩展名转写：tsc 产物为无扩展/`.ts(x)` 相对导入（Node ESM 不可解析）→ `.js`
        out = out.replace(/from\s+["'](\.[^"']*?)["']/g, (m, p) => {
          if (/\.(js|json|css|svg|png|mjs|cjs|woff2?)$/.test(p)) return m;
          return `from "${p}.js"`;
        });
        if (out !== text) {
          writeFileSync(full, out);
          rewritten.push(path.relative(dir, full));
        }
        const um = (out.match(unmapPattern) || []).length;
        if (um > 0) {
          unmapTotal += um;
          console.log(`  [unmapped ${um}] ${path.relative(distRoot, full)}`);
        }
      }
    }
  };
  walk(dir);
  if (rewritten.length > 0) console.log(`${pkg}: rewrote ${rewritten.length} d.ts`);

  // package.json 定稿
  const exported = { ".": { types, import: "./index.js" }, "./*": "./*" };
  const next = {
    name: old.name || `@schema-ui/${pkg}`,
    version: versions[old.name || `@schema-ui/${pkg}`] || old.version || "0.1.0",
    type: "module",
    description: old.description || `schema-ui-core 包面（${pkg}）`,
    main: old.main || "index.js",
    types,
    exports: exported,
    files,
    license: old.license || "UNLICENSED",
  };
  if (peers[old.name || `@schema-ui/${pkg}`]) next.peerDependencies = { ...peers[old.name || `@schema-ui/${pkg}`] };
  writeFileSync(pkgPath, JSON.stringify(next, null, 2));
  console.log(`${next.name} -> v${next.version} · exports "./*" · files [${files.join(", ")}]`);
}

console.log(`alias 重写完成；无法映射引用总数（残余登记）= ${unmapTotal}`);