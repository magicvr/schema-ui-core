---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-web-package-consumption
version: 0.1.0
---

# E-001 · S1 Web 结构扫描事实（2026-08-29）

## apps/web 结构（89 个 .ts/.tsx）

| 面 | 目录 | 说明 |
|----|------|------|
| 协议 | `src/protocol/` | schema 文档加载/校验（load-page.ts）、manifest 协商（app-manifest.ts）、conformance/、upstream/（provenance + fixtures） |
| 渲染器 | `src/renderer/` | schema 驱动渲染（render.tsx、schema-table/row-action/form-controls/reactions/permissions/资源绑定等 37 文件） |
| Shell | `src/app/` | App.tsx / LoginPage / AuthGate / navigation / branding / ManifestFailure（产品面） |
| Host | `src/host/` | v2.8+ Host/App 互操作层（boot/bootstrap/claim/failure/nginx-proxy） |
| UI 原子 | `src/components/` | shadcn/ui 风格组件（依赖 radix/cva/clsx/tailwind-merge/lucide） |
| 主题 | `src/theme/` | Token 语义化 + `brand.example.css` 覆盖机制（fork 定制面，现成） |
| 工具 | `src/lib/` | fetch-timeout 等通用工具（protocol 依赖其一） |

## 关键事实

1. **package.json**：`schema-ui-core-web` v0.1.0 private；React 19.2 · Tailwind 4.1（@tailwindcss/vite）· ajv 8（devDeps）· Vite 6.3 · vitest/playwright。无现成 workspace/monorepo。
2. **vite.config**：alias `@` → `./src`、`@schemas` → `../../docs/schemas`；dev proxy 直连 API。
3. **protocol 切片消费性**：load-page.ts 依赖 `@/lib/fetch-timeout` 与 `@/protocol/conformance/runtime-schema-validate`（ajv）——**拆包需同步拆 lib 或相对化 import**；纯 TS 无 React 依赖 → 最轻首包候选。
4. renderer/ui/shell 相互交织（renderer → components/theme/i18n）→ 完整拆包 = 源码迁移 + import 改造 + Tailwind 产物处理（CSS 面）。

## ⚠️ 协议 pin 漂移发现（I-007）

- 代码事实：`src/protocol/app-manifest.ts` → `APP_MANIFEST_PROTOCOL_VERSION = "2.9"`、`APP_MANIFEST_SOURCE = ".../81aa1d8"`（v2.9.0 formal release commit）；支持窗 `["2.7","2.8","2.9"]`。
- 文档事实：Charter/roadmap `v2.8.0`（pin `521cff8`；`I-PROTO-FULL-001` v1.0.1 历史分母）。
- 影响：按 VP-008 §go 消费触发规则「协议 pin 改变」→ 影响消费基线 freshness 与 R4 演练基线选择；**需要 `/vision` 裁决**（pin bump 是否升级为 2.9.0 + provenance 更新 + 容器门禁链），本目标内仅登记不裁决。
- 初证范围：差异 = v2.9 additive（支持窗仍含 2.7/2.8，前瞻性兼容保持）；具体 diff 面（ADR-0039/0040 等）在 `upstream/` fixtures 中可核对——留 I-007 收集。