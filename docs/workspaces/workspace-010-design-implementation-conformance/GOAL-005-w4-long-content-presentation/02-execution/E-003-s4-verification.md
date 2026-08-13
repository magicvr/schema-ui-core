---
id: E-003
goal_id: GOAL-005-w4-long-content-presentation
title: 执行 · S4 符合性验证事实
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-001-design-implementation-conformance
version: 0.1.0
---

# E-003 · S4 符合性验证（2026-08-13）

## 门禁结果

| 门禁 | 结果 |
|------|------|
| `apps/web` vitest 全量 | **48 文件 879 通过**（含本波新增 4 用例：DataTable truncate ×2、SchemaTable 透传 ×1、recordView 换行 ×1） |
| `apps/web` build（`tsc -b` + vite） | **通过**（`tsc -b` 0 错误；vite 产物正常；chunk 体积警告为既有） |
| `apps/api` `go test ./...` | **23 包全部 ok，无 FAIL**（roles 包含 schema 服务路径） |
| 上游 conformance fixtures | 全量在 vitest 内回归通过（48 文件含 `upstream-host-fixtures`、`upstream-fixtures`、`stage3-fixtures` 等） |
| conformance claim 再生成 | `prebuild` 已随 build 再生成（`conformance-claim.json` + sha256 + local-report 一致） |

## 本波新增测试用例（W4 验收标准对照）

1. `data-table.test.tsx`：`truncate` 列存在 `data-table-cell="truncated"` 容器、类含 `truncate` + `max-w-[16rem]`、`title` 为全文（`users.read,users.write,roles.read,roles.write`）；非 truncate 列无该容器（行为不变）。
2. `schema-table.test.tsx`：`truncate: true` 列规范透传到 shipped 单元格（`title` 含逗号连接全文）。
3. `render.test.tsx`：recordView 静态 record 数组值按 `", "` 连接；值列轨道类为 `sm:grid-cols-[8rem_minmax(0,1fr)]`；`dd` 保留 `break-words`。

## S4 中发现并处置的额外事实（超出 D-001 范围的必要处置）

- `tsc -b` 暴露 **3 个 W3 遗留类型错误**（`host/boot.ts` ×2、`main.tsx` ×1，均为 W3 account-locked 生产源提交引入；`git status` 确认本波改动前未被我方修改）：
  - `boot.ts` 的 `reauthFailure()`/`lockedFailure()` 合成 `BootstrapEvaluation` 时把 `"REAUTH_REQUIRED"`/`"ACCOUNT_LOCKED"`（`BootstrapResult` 值）写进了 `code`（`BootstrapResultCode`）字段。
  - `main.tsx` 使用 `SessionAdapterState` 未导入（W3 未解 issue）。
- **处置（最小修复，与 W3 台账 `tsc 0 错误` claim 一致）**：
  - `recoveryActionsFor` 形参收窄为 `Pick<BootstrapEvaluation, "fetchClassification">`（函数只读该字段）；两处合成调用改为 `{ fetchClassification: null }`。
  - `main.tsx` 补 `type SessionAdapterState` 导入。
- 行为语义零变化：`fetchClassification: null` 与 W3 字面量完全一致；`recoveryActionsFor` 只依赖 `fetchClassification`（S3 代码事实核对）。
- 该处置**不是**对 D-001 范围的扩张，而是 S4 验证门禁（D-001 §3.3）的使能修复；E-002 已预告「W3 遗留类型门禁若拦截将最小修复并留痕」。

## 验收标准核对（D-001 §3）

1. ✓ roles 列表截断容器 + title 全文（schema-table 透传测试 + DataTable 测试）。
2. ✓ recordView 数组 `", "` 连接 + 可收缩轨道（render.test 断言）。
3. ✓ 全量 vitest + build + Go 回归 + conformance fixtures（见上表）。
4. ✓ diff 面与 D-001 §2 一致（+ W3 类型门禁最小修复，见上）。
