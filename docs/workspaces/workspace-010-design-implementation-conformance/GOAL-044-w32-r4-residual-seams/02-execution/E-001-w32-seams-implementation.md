---
id: E-001-w32-seams-implementation
doc: execution-entry
parent: GOAL-044-w32-r4-residual-seams
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-001 · W32 C1～C3 实施（三项通用 seam）

## 事实（2026-09-19）

### 1. C1 · 方案冻结

`01-decision/D-001-w32-seams-freeze.md` v1.0.0 落盘：① 列级 `valueLabels`（本地扩展，fail-open 回落原始值，与 pinned `tagMap` 划界）；② `refreshTable(tableId)`（保留选择；明确不删 in-flight 键的取舍理由）；③ `activeStatuses` + 行注册表的空闲判定；另含用户授权与跨区可写范围（§5）与仍不可写清单。

信息项 `I-044-001`/`002`/`003` 均由 `D-001` 关闭为 `verified`（决策即结论）。

### 2. C2 · 实施产物

| 产物 | 位置 |
|------|------|
| `SchemaCrudValue` 新增 `refreshTable` / `tableRefreshToken` / `publishTableRows` / `tableRows`（ref 注册表，发布不触发重渲染） | `apps/web/src/renderer/render.tsx` |
| 列 `valueLabels` 解析（`labeledCellValue`；badge 列与普通列均生效，键缺失回落原始值） | `apps/web/src/renderer/schema-table.tsx` |
| 取数 effect 依赖 `tableRefreshToken`；发布当前行到注册表 | `apps/web/src/renderer/schema-table.tsx` |
| `jobs-auto-refresh` 改用 `refreshTable(targetTable)`；按 `activeStatuses` 跳过空闲拍（行不可得时保守刷新） | `apps/web/src/components/jobs-auto-refresh.tsx` |
| `status` 列六态 `valueLabels`；组件节点声明 `targetTable`/`statusField`/`activeStatuses` | `apps/api/modules/jobs/schema/jobs.json` |

### 3. C3 · 测试与回归

| 验证 | 结果 |
|------|------|
| `go build ./...` / `go vet ./...`（apps/api） | exit 0 / 无输出 |
| `go test ./...`（apps/api） | **全绿**（含 `modules/jobs`、`docscheck`、`internal/handler`） |
| `npm run typecheck` | exit 0 |
| `npm test`（vitest） | **118 files / 1463 tests 全绿**（W32 前基线 117/1458 → +1 file / +5 tests） |

**新增测试**

| 测试 | 断言 |
|------|------|
| `table-refresh-seam.test.tsx`（3） | ① `refreshTable` 重取该表格且**选择集保持不变**；② `reloadList()` 仍清空全部选择（ADR-0022 D2 对照，未被削弱）；③ 行注册表可读、未知 tableId 为 no-op（不触发全页刷新） |
| `jobs-result-center.test.tsx` +2 | ① 六态在 en-US（Running）与 zh-CN（执行中）下本地化，且**未映射值原样保留**（fail-open）；② 全部终态时 5s 档 **30s 内零新增请求**，正对照（存在 running 行）确实轮询 |

**变异验证（均实测变红后还原，工作树无残留）**

| 变异 | 结果 |
|------|------|
| 从 `jobs.json` 移除 `valueLabels` | 本地化例 **变红**（状态单元格回落原始码） |
| 让 `refreshTable` 同时 `setSelections({})` | seam 第 1 例 **变红**（选择被清空） |
| 禁用 `activeStatuses` 空闲判定 | 空闲例 **变红**（30s 内仍在请求） |

### 4. 边界

- **未**触碰 pinned 工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）；未新增 pinned 能力 id。
- **未**改 `reloadList()` 的 ADR-0022 D2 语义（以对照测试钉住）；未改 `refreshList`（display-only）。
- **未**改 Job 六态合同、同步 `batch-delete`、权限键或后端路由。
- `apps/api/modules/jobs/schema/jobs.json` 的改动仅为列 `valueLabels` 与组件节点 props。

### 5. Git checkpoints

| hash | 内容 |
|------|------|
| `c2ea042b` | W32 C2/C3 实施与测试（渲染器 seam + 组件 + jobs.json + 2 个测试文件） |

> `D-001` 与台账同步的 checkpoint 见 C4 条目。

### 6. 未做（移交 C4）

- **未**执行审计（C4：self；是否追加 independent 见 C4 判定）。
- **未**回填 `[workspace-038] GOAL-005` `A-003`（C4 完成）。
