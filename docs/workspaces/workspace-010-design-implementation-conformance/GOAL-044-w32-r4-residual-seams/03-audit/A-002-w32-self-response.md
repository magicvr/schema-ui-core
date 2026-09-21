---
id: A-002-w32-self-response
doc: audit-entry
parent: GOAL-044-w32-r4-residual-seams
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-002 · W32 自审响应（A-001 三条 recommended 闭合）

- **source**：orchestrator（`/govern`）
- **scope**：GOAL-044 C4 —— 响应 `A-001`（self，pass，0 required + 3 recommended）并驱动 C4 闭合。

| # | finding | 处置 | 证据 |
|---|---------|------|------|
| 1 | `A-001` F-001：`valueLabels` 只测了「值未映射」，未测「键不在目录中」 | **fixed** | 新增用例：把 `status` 列的映射指向不存在的 `schema.jobs.status.gone`，断言单元格回落为用户可见的原始值 `running`，且**不**出现键名本身（`jobs-result-center.test.tsx`，13 例） |
| 2 | `A-001` F-002：③ 的保守分支（行不可得时仍刷新）未被直接断言 | **fixed** | 新增 `components/jobs-auto-refresh.test.tsx`（4 例，组件边界 + stub context）：有活跃行→刷新；全终态→跳过；**行未发布（`undefined`）→保守刷新**；档位 Off→零调用 |
| 3 | `A-001` F-003：本地扩展缺少单一清单登记 | **fixed（登记制）** | 在 `docs/vision/roadmap.md`「未决项统一登记 · 一、有界残余」增设「本地扩展登记」项（现状 / 触发条件 / 责任人 / 证据），并更新该节「最近更新」；与 W31 建立的登记机制一致 |

**开放 required = 0**；三条 recommended 全部 `fixed`，无冲突、无残余需用户裁决。

**回归（响应后）**：`go test ./...` 全绿；`npm test` **119 files / 1468 tests 全绿**；`npm run typecheck` exit 0。

## 回填 `[workspace-038] GOAL-005`

| 原状态 | 现状态 | 依据 |
|--------|--------|------|
| `A-001` F-002（状态/错误文本未逐值本地化） | **`fixed`** | 列级 `valueLabels` 落地 + jobs 六态接入 + 双语言断言 + 未映射/缺键两条回落路径断言（`D-001` §1；`E-001` §3） |
| `A-001` F-003（自动刷新依赖 `reloadList()` 清空选择） | **`fixed`** | `refreshTable(tableId)` 保留选择并重取该表格；`reloadList()` 的 ADR-0022 D2 语义以对照测试钉住未被削弱；组件已改用新 seam（`D-001` §2） |
| `A-001` F-004（无进行中作业时仍按档位轮询） | **`fixed`** | `activeStatuses` + 行注册表实现空闲跳过；页面级 30s 零请求断言 + 组件级保守分支断言（`D-001` §3） |

> 按 P-003，回填只加闭合注记（`fixed` + 证据），**不改** `[workspace-038] GOAL-005` 的 status/progress 与其 A-001/A-002 正文。
