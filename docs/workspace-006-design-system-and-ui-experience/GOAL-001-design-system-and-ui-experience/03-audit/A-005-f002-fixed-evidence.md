---
id: GOAL-001-design-system-and-ui-experience
doc: audit-entry
record_id: A-005
status: active
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

## A-005 · F-002 实施证据核查与 fixed 闭合（2026-08-09）

- **source**：self（编排响应 — 对应 A-002 required finding F-002）
- **scope**：F-002 Shadow 映射实施证据 / GOAL-001 S1 完成门禁
- **verdict**：pass（F-002 fixed）

### 核查对象

A-002 required finding F-002：「Shadow 原始语义变量与 Tailwind `--shadow-*` theme namespace 的无自引用映射，须以 S1 实施证据验证」。

### 实施证据

| 证据项 | 路径 / 描述 |
|--------|-------------|
| `--elevation-sm|md|lg` 定义 | `apps/web/src/index.css` §3 `:root`，原始 oklch shadow 值，无别名循环 |
| `@theme inline` alias | `--shadow-sm: var(--elevation-sm);` `--shadow-md: var(--elevation-md);` `--shadow-lg: var(--elevation-lg);`（§4） |
| 无自引用断言 | `theme.test.ts`：`toContain("--shadow-sm: var(--elevation-sm)")` + `not.toContain("--shadow-sm: var(--shadow-sm)")` 三对均通过 |
| confirm/modal 迁移 | `renderer/confirm.tsx` → `shadow-md`；`renderer/modal.tsx` → `shadow-lg`（C5 消费闭环） |
| vitest 全绿 | 595 tests / 0 failures（含 Token 结构断言） |
| build 通过 | `vite build` exit 0；dist 产物正常 |

### F-002 闭合路径

**fixed**：实施采用 `--elevation-*` 保存原始语义值 + `@theme inline` 单向 alias 至 `--shadow-*`，完全符合 A-002 recommended 方案。vitest 结构断言证明无自引用，`npm run build` 证明 utility 可生成。消费点（confirm/modal）已从 `shadow-xl` 迁移到语义 utility。

### 结论

**F-002 → fixed**（2026-08-09）。A-002 开放 required finding 清零。S1 完成门禁已满足，可勾选 GOAL-001 S1 检查点。
