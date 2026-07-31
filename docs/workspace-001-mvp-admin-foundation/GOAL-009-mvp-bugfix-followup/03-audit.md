---
id: GOAL-009-mvp-bugfix-followup
doc: audit
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.2.0
---

# 审计 · GOAL-009

> 本文件是本目标的**唯一正式意见台账**（P-003）。长文可链附件；每条正式意见须为 `A-00N` 编号节。

## 信息就绪核对（立项 scope）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| I-009-001 / I-009-002 | resolved（用户书面裁决纳入） | F-009-006/007 升格实施并 `fixed`；2026-08-01 用户裁决「都纳入实施」 |
| 到期 required 信息 | 无 | 审视附件已固定 findings 清单 |
| 共享资料引用 | 无 | workspace `shared_materials_catalog: none` |

---

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required |
|------|------|--------|-------|---------|---------------|
| A-001 | 2026-08-01 | independent | 立项基线 · 代码审视 findings | conditional | 0 开放（5 required + 2 recommended 均 `fixed`）；阶段/关门审计未做 |

---

## A-001 · 代码审视 findings 基线（2026-08-01）

- **source**：`independent`
- **auditor**：代码审视（会话内对照 VP-001 + 源码通读 + 测试绿）
- **类型**：`ad-hoc`（立项前代码审视 → 转入本目标台账）
- **scope**：`apps/api` + `apps/web` 实现 vs VP-001 意图；真实 bug / 集成失真 / 文档；本目标成功标准基线
- **verdict**：`conditional`（VP MVP 意图大体符合；存在须修 findings，未实施修正前不得对本目标宣称 done）
- **长文**：[attachments/audit-code-review-bugs-2026-08-01.md](attachments/audit-code-review-bugs-2026-08-01.md)

### 范围与区间

- **覆盖**：API records/account/permission；Web boot/account/list-edit/README；与 VP 退出判据对齐摘要。
- **不覆盖**：完整 schema 宿主、生产 IAM、已排除域（upload/batch 等）的实现完整性验收。
- **与既有意见**：非 GOAL-006/007/008 历史 A-00N 的复审；新目标新序列。Root/子目标历史 pass **不**被本意见推翻。

### 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工程可构建、单测绿（审视时） | `go test ./...`；Web 395 tests |
| VP-001 方向大体符合 + 有界 residual | 附件 §1–§2；VP 关门记录 |
| Findings 已枚举并可引用 | 附件 §3 F-009-001～007 |

### 对照成功标准（2026-08-01 实施后）

| 成功标准 | 状态 | 证据 |
|----------|------|------|
| F-009-001 updatedAt | **done** | `records.go` update()；`TestRecordsUpdateRefreshesUpdatedAt` |
| F-009-002 list-edit 权限/context | **done** | 真实 context prop + 表 cascade/表达式；`list-edit-lifecycle.test.tsx` 拒绝路径 |
| F-009-003 account error 可观察 | **done** | main.tsx 保留 error + App alert 横幅；`App.integration.test.tsx` |
| F-009-004 sessionProvider 一致 | **done** | account.go nil fail-closed；`TestAccountsMeNilProviderFailsClosed` |
| F-009-005 README | **done** | 两 README 重写与端点/边界一致 |

> 5 条 required 已按 P-003 `fixed` 合法闭合（见下方 Findings 状态与闭合留痕）。F-009-006/007 经用户裁决纳入实施（I-009-001/002 resolved）并 `fixed`。**未**执行阶段/关门审计，本目标保持 `active`。

### Findings

#### F-009-001 · PATCH 不更新 updatedAt

| 字段 | 值 |
|------|-----|
| **级别** | required |
| **严重度** | med |
| **影响门禁** | 实施 / 验收 / 关门 |
| **状态** | fixed |
| **描述** | 见附件 §3 F-009-001 |
| **证据** | `apps/api/internal/handler/records.go` update 路径 |
| **关闭要求** | 成功 PATCH 刷新时间 + 测试 |
| **闭合留痕** | 2026-08-01 `/govern`：update() 成功 PATCH 后写 `time.Now().UTC()`；`TestRecordsUpdateRefreshesUpdatedAt` 断言 updatedAt 后移；`go test` 绿 |

#### F-009-002 · list-edit 权限演示失真

| 字段 | 值 |
|------|-----|
| **级别** | required |
| **严重度** | med |
| **影响门禁** | 实施 / 验收 / 关门 |
| **状态** | fixed |
| **描述** | 硬编码 admin context；无真实 permissions 表达式；门禁恒过 |
| **证据** | `apps/web/src/app/examples/list-edit-lifecycle-page.tsx` |
| **关闭要求** | 真实 context + 可失败表达式 + 拒绝路径测试 |
| **闭合留痕** | 2026-08-01 `/govern`：页面接入真实 `navigationContext` prop；`PAGE_DOCUMENT` 表节点挂 `permissionCascade` + `permissions` 表达式；`list-edit-lifecycle.test.tsx` 覆盖 admin 启用 / viewer 拒绝 disabled；web 398 测试绿 |

#### F-009-003 · Account 失败静默

| 字段 | 值 |
|------|-----|
| **级别** | required |
| **严重度** | low–med |
| **影响门禁** | 实施 / 验收 / 关门 |
| **状态** | fixed |
| **描述** | `main.tsx` 丢弃 `loadAccountContext` error |
| **证据** | `apps/web/src/main.tsx` |
| **关闭要求** | 失败可观察 |
| **闭合留痕** | 2026-08-01 `/govern`：main.tsx 保留 error + `console.error` + `accountError` prop；App 渲染非阻断 `role="alert"` 横幅；`App.integration.test.tsx` 断言横幅（含健康路径无横幅） |

#### F-009-004 · sessionProvider nil/注释

| 字段 | 值 |
|------|-----|
| **级别** | required |
| **严重度** | low |
| **影响门禁** | 实施 / 验收 / 关门 |
| **状态** | fixed |
| **描述** | 注释称 nil→static，实际 nil panic |
| **证据** | `apps/api/internal/handler/account.go` |
| **关闭要求** | 行为与文档一致 |
| **闭合留痕** | 2026-08-01 `/govern`：注释改为「nil 即 fail-closed」；`me()` nil provider 返回 `401 UNAUTHENTICATED`；`TestAccountsMeNilProviderFailsClosed`；`go test` 绿 |

#### F-009-005 · README 滞后

| 字段 | 值 |
|------|-----|
| **级别** | required |
| **严重度** | low |
| **影响门禁** | 验收 / 关门 |
| **状态** | fixed |
| **描述** | API/Web README 与现状不符 |
| **证据** | `apps/api/README.md` 等 |
| **关闭要求** | 与当前端点与 MVP 边界一致 |
| **闭合留痕** | 2026-08-01 `/govern`：两 README 重写为当前端点、session 模型、鉴权边界、测试命令 |

#### F-009-006 · records 未挂 Allow（recommended）

| 字段 | 值 |
|------|-----|
| **级别** | recommended |
| **严重度** | med |
| **影响门禁** | 默认不单独阻断；升格见 I-009-001 |
| **状态** | fixed |
| **描述** | 见附件 §3 F-009-006 |
| **关闭要求** | 写路由鉴权或用户 residual |
| **闭合留痕** | 2026-08-01 `/govern`：用户裁决纳入（I-009-001）；records.go `writeGate()` 对 PATCH/DELETE 挂 fail-closed（无会话 401 / 非 admin 403，`account.Allow` admin 角色门槛）；`TestRecordsWriteRequiresSession` + `TestRecordsWriteDeniedWithoutAdminRole` |

#### F-009-007 · body/pageSize 上限（recommended）

| 字段 | 值 |
|------|-----|
| **级别** | recommended |
| **严重度** | low |
| **影响门禁** | 默认不单独阻断；升格见 I-009-002 |
| **状态** | fixed |
| **描述** | 见附件 §3 F-009-007 |
| **关闭要求** | 上限 + 测试或 residual |
| **闭合留痕** | 2026-08-01 `/govern`：用户裁决纳入（I-009-002）；records.go `MaxBytesReader` 4 KiB + `pageSize ≤ 100`；`TestRecordsUpdateBodyTooLarge` + `TestRecordsListPageSizeCap` |

### 编排提示

- 存在 **5** 条开放 required：在未 `fixed` / residual / overruled 前，**不得**将本目标标 `done`。
- 响应与实施走 **`/govern`**；本 A-001 不改 Root/VP status。

### 响应（2026-08-01 · `/govern` 实施后）

| date | actor | scope | summary |
|------|-------|-------|---------|
| 2026-08-01 | `/govern` | A-001 全部 findings | **F-009-001～005 → `fixed`**（证据见各 Finding 闭合留痕与「对照成功标准」表）：API PATCH 刷新 `updatedAt`；list-edit 接入真实 context + 可失败权限表达式 + 拒绝路径组件测；account 失败可观察横幅；sessionProvider nil fail-closed；两 README 更新。**F-009-006/007 → `fixed`**（用户裁决纳入，I-009-001/002 resolved）：records 写路由挂 fail-closed admin 会话鉴权；PATCH body ≤ 4 KiB、`pageSize ≤ 100`。回归：`go test ./...` 绿、web vitest **398/398**、`tsc -b` 与 `npm run build` 通过。 |

**门禁状态**：A-001 全部 7 条 findings 已按 P-003 `fixed` 合法闭合（5 required + 2 recommended，用户书面裁决纳入 006/007）。剩余开放项：**阶段/关门审计未做**。**本目标仍未 `done`**：关门需先经阶段/关门审计（self 或 `/audit`）复核 fixed 证据。
