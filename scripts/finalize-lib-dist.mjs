#!/usr/bin/env node
/**
 * finalize-lib-dist.mjs — dist 产物终版整理（R5 · 幂等可复跑）。
 * 对 apps/web/dist-lib/@schema-ui/<pkg> 下 .js/.d.ts/.json 执行：
 *   1) `@schema-ui/` → `@magicvr/schema-ui-`（发布实态全名）
 *   2) 相对导入扩展名规范化：`.tsx.js`/`.ts.js` → `.js`；无扩展相对导入 → `.js`
 * 用法：node scripts/finalize-lib-dist.mjs
 */
import { readdirSync, readFileSync, writeFileSync, copyFileSync, mkdirSync, existsSync } from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));
const distRoot = path.resolve(here, "../apps/web/dist-lib/@schema-ui");
const schemasSrc = path.resolve(here, "../docs/schemas");

// 3) JSON 资产面：protocol 的 @schemas/*.json import → 包内 schemas/ 资产 + import attributes
const schemaImports = ["action", "node", "page", "reaction"];
for (const name of schemaImports) {
  const src = path.join(schemasSrc, `${name}.schema.json`);
  if (!existsSync(src)) continue;
  const outDir = path.join(distRoot, "protocol/schemas");
  mkdirSync(outDir, { recursive: true });
  copyFileSync(src, path.join(outDir, `${name}.schema.json`));
}
// F2 (GOAL-042 · host-support 移入 protocol 包)：tsc 不复制 JSON 资产，构建时
// 把版本/能力协商单源 host-support.json 拷到 host-support.js 同目录
// （protocol/protocol/ · load-page 运行时依赖 ./host-support.json）。
const hostSupportJsonSrc = path.join(here, "../apps/web/src/protocol/host-support.json");
if (existsSync(hostSupportJsonSrc)) {
  const outDir = path.join(distRoot, "protocol", "protocol");
  mkdirSync(outDir, { recursive: true });
  copyFileSync(hostSupportJsonSrc, path.join(outDir, "host-support.json"));
}
const protocolDir = path.join(distRoot, "protocol");
const walkProto = (d) => {
  for (const ent of readdirSync(d, { withFileTypes: true })) {
    const full = path.join(d, ent.name);
    if (ent.isDirectory()) walkProto(full);
    else if (/\.js$/.test(ent.name)) {
      const t = readFileSync(full, "utf8");
      // @schemas/*.json → 包自子路径 @schema-ui/protocol/schemas/*（generic 段再转
      // @magicvr/schema-ui-protocol/schemas/*）。0.2.11 实证形态：相对 ./schemas/
      // 在 protocol/conformance/ 等子目录下解析失败（消费端 ERR_MODULE_NOT_FOUND）。
      let out = t.replace(
        /@schemas\/([a-z-]+)\.schema\.json/g,
        (m, n) => `@schema-ui/protocol/schemas/${n}.schema.json`,
      );
      // 包内 JSON import（包自子路径 schemas/*、相对 host-support.json）→ import
      // attributes；负向断言保证幂等（可复跑不重复追加）。
      out = out.replace(
        /from\s+(["'][^"']+\.json["'])(?!\s+with\b)/g,
        'from $1 with { type: "json" }',
      );
      if (out !== t) writeFileSync(full, out);
    }
  }
};
if (existsSync(protocolDir)) walkProto(protocolDir);

let files = 0;
for (const pkg of readdirSync(distRoot)) {
  const dir = path.join(distRoot, pkg);
  const walk = (d) => {
    for (const ent of readdirSync(d, { withFileTypes: true })) {
      const full = path.join(d, ent.name);
      if (ent.isDirectory()) walk(full);
      else if (/\.(js|d\.ts|json)$/.test(ent.name)) {
        const t = readFileSync(full, "utf8");
        let out = t;
        out = out.split("@schema-ui/").join("@magicvr/schema-ui-");       // 1)
        out = out.replace(/\.tsx?\.js"/g, '.js"');                        // 2a) .tsx.js → .js
        out = out.replace(/from\s+["'](\.[^"']*?)["']/g, (m, p) => {      // 2b) 无扩展相对 → .js
          if (/\.(js|json|css|svg|png|mjs|cjs|woff2?)$/.test(p)) return m;
          return `from "${p}.js"`;
        });
        out = out.replace(/from\s+["'](@magicvr\/schema-ui-[^"']*?)["']/g, (m, p) => { // 2c) 无扩展包子路径 → .js
          if (/\.(js|json|css|svg|png|mjs|cjs|woff2?)$/.test(p)) return m;
          return `from "${p}.js"`;
        });
        // 2d) 包内相对 JSON import（i18n/messages/*.json 等，lib/ui tsc 产物保留
        // 裸 JSON import）→ import attributes；负向断言幂等。Node ESM 必须。
        out = out.replace(
          /from\s+(["'][^"']+\.json["'])(?!\s+with\b)/g,
          'from $1 with { type: "json" }',
        );
        if (out !== t) {
          writeFileSync(full, out);
          files++;
        }
      }
    }
  };
  walk(dir);
}
console.log(`finalize-lib-dist: ${files} files normalized`);