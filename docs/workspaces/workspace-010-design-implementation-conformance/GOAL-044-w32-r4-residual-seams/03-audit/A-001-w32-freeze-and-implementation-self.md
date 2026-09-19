---
id: A-001-w32-freeze-and-implementation-self
doc: audit-entry
parent: GOAL-044-w32-r4-residual-seams
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-001 · W32 方案与实施自审（GOAL-044 C1～C3）

## A-001 · W32 C1～C3 自审（2026-09-19）

- **source**：self
- **auditor**：`/govern` 编排器（DeepSeek Harness 会话）
- **类型** / **scope**：`stage` · GOAL-044 C1～C3（方案冻结 / 三项实施 / 测试与回填）
- **verdict**：**pass**（0 required；3 recommended）
- **审计模式判定**：**`self`**。理由：改动是渲染器**本地扩展与内部 seam**，不涉安全/数据/迁移/发布面；`reloadList()` 的 ADR-0022 D2 语义以**对照测试**钉住未被削弱；无后端契约、权限键、pinned 工件或 Job 合同变更。跨工作区影响（038 的结果中心组件与页面）将由 038 的 **R5**（浏览器/自动化回归 + 独立意见）再验一次；若 R5 发现回退，本目标的残余判定需重开。

### 范围与区间

- 被审目标：`docs/workspaces/workspace-010-design-implementation-conformance/GOAL-044-w32-r4-residual-seams/`
- 审计区间：`d8d2de71`（038 R4 关门与移交）→ `c2ea042b`（W32 实施）。
- 承接输入：`[workspace-038…] GOAL-005` 的 `A-001` F-002/F-003/F-004 与 `A-003` §2（用户裁决 + 有界接受 + 复核触发）。

### 成果（有证据）

| # | 成果 | 证据 |
|---|------|------|
| 1 | ① 列值本地化落地为**通用**列能力（非 jobs 特例），badge 列与普通列都生效 | `schema-table.tsx` `labeledCellValue`；`status` 列同时使用 `badgeStyleField` 与 `valueLabels` |
| 2 | ① 的失败模式是**展示层 fail-open**：未映射值/缺键 → 原始值（不会清空单元格） | 同上；测试断言未映射的 `quiesced` 原样呈现 |
| 3 | ② `refreshTable` **保留选择集**，且用当前查询重取该表格 | `render.tsx`；seam 测试第 1 例（重取 + 选择仍为 2） |
| 4 | ② **未削弱** ADR-0022 D2：`reloadList()` 仍清空全部选择 | seam 测试第 2 例（对照） |
| 5 | ② 的 in-flight 取舍已显式记录并说明理由（只读轮询可合流；变更后仍须 reloadList） | `D-001` §2；`render.tsx` 注释 |
| 6 | ③ 空闲判定基于**行可见性**（表格发布行）+ 显式 `activeStatuses`；行不可得时保守刷新 | `render.tsx` 注册表；`jobs-auto-refresh.tsx`；空闲例 + 正对照 |
| 7 | ③ 在 jobs 页声明的活跃态 = `{queued, running}`，与 Job 六态中「非终态」一致 | `jobs.json` 节点 props；`D-001` §3 |
| 8 | 组件不再调用 `reloadList()`（改为 `refreshTable`），从根本上消除「清空选择」的耦合 | `jobs-auto-refresh.tsx`；seam 测试把两条 seam 的差异钉死 |
| 9 | 三项均经**变异验证**（见下） | `E-001` §3 |
| 10 | 未触 pinned 工件 / 权限键 / 后端契约 / Job 六态 | 本区间 `git diff --stat` 仅 `apps/web/src/renderer/**`、`apps/web/src/components/jobs-auto-refresh.tsx`、`apps/api/modules/jobs/schema/jobs.json` + 台账 |

**变异验证（本区间实测）**：① 移除 `valueLabels` → 本地化例红；② `refreshTable` 加 `setSelections({})` → seam 第 1 例红；③ 禁用空闲判定 → 空闲例红。三次变异均已还原，工作树无实现残留（`git status` 仅台账）。

### 回归证据

| 命令 | 结果 |
|------|------|
| `cd apps/api && go build ./... && go vet ./... && go test ./... -count=1` | exit 0 / 无输出 / **全绿 0 FAIL** |
| `cd apps/web && npm run typecheck` | exit 0 |
| `cd apps/web && npm test` | **118 files / 1463 tests 全绿**（较基线与 038 R4 关门前 +1 file / +5 tests） |

### Findings

#### F-001 · ① 的键缺失路径只测了「值未映射」，未测「键不在目录中」

- 严重度：low
- 建议：**recommended**
- 描述：`labeledCellValue` 的回落有两个触发条件（值未映射、键名在目录中不存在），测试只覆盖了前者。键名缺失时的行为（`translate(key, undefined, raw)` → 记 missing-translation 并返回原始值）未被断言；若将来有人把 `literalFallback` 去掉，单元格会变成键名本身（如 `schema.jobs.status.queued`）。
- 证据：`schema-table.tsx` `labeledCellValue`；`jobs-result-center.test.tsx` 的本地化例（仅未映射场景）。
- 状态：**open（recommended）** —— 处置建议：补一例「映射到不存在的键 → 回落原始值」；或接受为低危残余（当前页面无此形态）。

#### F-002 · ③ 的保守分支（行不可得时仍刷新）未直接断言

- 严重度：low
- 建议：**recommended**
- 描述：`tableRows()` 返回 `undefined` 时组件按「有活跃行」处理（宁多刷不漏刷）。该分支只由代码注释与 `D-001` §3 表达，没有测试直接钉住——若将来被改成「不可得即跳过」，作业页在首次取数前会停止刷新且不会被现有用例发现。
- 证据：`jobs-auto-refresh.tsx` 的 `rows !== undefined` 分支；现有空闲例覆盖的是 `rows` 已发布且全部终态的场景。
- 状态：**open（recommended）**

#### F-003 · `valueLabels` 未在 `docs/README` 或组件注册表文档中登记为本地扩展

- 严重度：low
- 建议：**recommended**
- 描述：仓库既有本地扩展（`badgeStyleField`/`truncate` 等）同样未在 `docs/README` 集中登记，因此本项与现状一致；但这类扩展正在增加，长期缺少单一清单会让「哪些是本地扩展、哪些是 pinned」难以分辨。`D-001` §1 已在本目标内记录边界。
- 证据：`D-001` §1；全仓无本地扩展清单文件。
- 状态：**open（recommended）** —— 处置建议：由后续符合性波次考虑建立「本地扩展登记」节（不在本目标扩范围）。

### 必改项汇总（required）

**无。**

### 结论 + 建议下一步

三项残余均以**通用能力**形态落地（列值本地化、定向刷新 seam、行可见性驱动的轮询），038 侧的三条 `accepted-residual` 可按 P-003 **回填为 `fixed`**（附证据与变异记录）。`reloadList()` 的既有语义未被削弱，且由对照测试钉住。

**verdict = pass**。**建议下一步**：C4 回填 `[workspace-038…] GOAL-005` `A-003` 三条为 `fixed`、同步 `goal-tree.md`/`workspace.md`、`GOAL-044` 关门（`done · 4/4`）；038 的 R5 在浏览器/自动化回归中复核本目标的用户可见效果。

**本条不修改** `status`、检查点或派生 `progress`。
