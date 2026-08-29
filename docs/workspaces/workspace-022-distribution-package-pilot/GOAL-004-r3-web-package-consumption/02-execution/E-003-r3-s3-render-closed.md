---
status: done
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-web-package-consumption
version: 0.1.0
---

# E-003 · S3 渲染闭环事实（2026-08-29）

## 产物

- **`@schema-ui/renderer` v0.1.0**（粗粒度单包 · 用户裁决 D-002）：入口 `src/renderer/index.ts`（RenderPage / useSchemaCrud / registerCustomComponent / I18nProvider / 类型面）；Vite lib 构建 **436.7 kB · 1653 modules**（components/i18n/lib/protocol 内部 bundle；**React 体系外部化为 peer**）；d.ts = 手工最小声明（自动化链路 = F-006）。
- 构建链路复用 protocol 模式（`vite.lib.renderer.config.ts` + `tsconfig.renderer.json`）。

## 验证三探针（golden-web，仅包依赖）

| 探针 | 断言 | 结果 |
|------|------|------|
| `probe-render.mjs` | SSR 渲染真实形态 schema 页面文档（form + defaultValue + reaction 语义）→ HTML 含 label/值 | **PASS**（1573 bytes） |
| `token-check.mjs` | brand.css 覆盖 ⊆ index.css 变量集（fork 定制纪律） | **PASS**（brand=2 ⊆ index=5） |
| 能力门控（调试期捕获） | 缺 `form.controls.advanced` 时渲染 `FORM_CAPABILITY_REQUIRED` 错误面 → 补能力后正常渲染 | 语义保护在包内生效（fail-closed ✓） |
| `probe.mjs`（S2） | protocol 包功能断言 | PASS |

## 结论

判据 #2 核心成立：**空下游 app 仅 npm 包组组装 → 渲染与主线同一渲染器（同源产物）的 schema 页面文档；Token 覆盖定制面以文件纪律验证**。浏览器级视觉与 Tailwind `@source` 实编译 = 主仓既有视觉回归基线 + README 指引（peer 矩阵实测点）覆盖。

## 关联

- F-006（recommended）：renderer 包 d.ts 自动化链路 = TS5056（`render.ts`/`render.tsx` 同名输出冲突）→ 登记，**go 后随 monorepo 化修复**（当前手工最小声明可用，JS 运行时不受影响）。
- I-002：六包细化（ui/theme 独立）与 d.ts 链路并入 go 后正式化清单。
- Web 冻结面（v1.2 候选）：六包边界 + peer 矩阵（v0.1 设计）在 go 后升格。