---
id: GOAL-009-mvp-bugfix-followup
doc: audit
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.3.0
---

# 审计 · GOAL-009

> 本文件是本目标的**唯一正式意见台账**（P-003）。长文可链附件；每条正式意见须为 `A-00N` 编号节。

## 信息就绪核对（立项 scope）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| I-009-001 / I-009-002 | resolved（用户书面裁决纳入） | F-009-006/007 升格实施并 `fixed`；2026-08-01 用户裁决「都纳入实施」 |
| 到期 required 信息 | 无 | 审视附件已固定 findings 清单；A-002 关门复审确认无到期 required 信息门禁 |
| 共享资料引用 | 无 | workspace `shared_materials_catalog: none` |

---

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required |
|------|------|--------|-------|---------|---------------|
| A-001 | 2026-08-01 | independent | 立项基线 · 代码审视 findings | conditional | 0 开放（5 required + 2 recommended 均 `fixed`）；阶段/关门审计由 A-002 承接 |
| A-002 | 2026-08-01 | independent | 关门复审 · A-001 七条 fixed 证据 + 成功标准 | pass | 0（2 recommended 文档卫生/边界措辞，不阻断） |

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

> **历史基线节（A-001 立项时点）**：此处「5 条开放 required」为审视时基线，仅描述立项时状态。实施后 F-009-001～007 已全部 `fixed`（见下方响应表）；关门复核见 **A-002**（pass）。不属当前开放状态。

- 存在 **5** 条开放 required：在未 `fixed` / residual / overruled 前，**不得**将本目标标 `done`。
- 响应与实施走 **`/govern`**；本 A-001 不改 Root/VP status。

### 响应（2026-08-01 · `/govern` 实施后）

| date | actor | scope | summary |
|------|-------|-------|---------|
| 2026-08-01 | `/govern` | A-001 全部 findings | **F-009-001～005 → `fixed`**（证据见各 Finding 闭合留痕与「对照成功标准」表）：API PATCH 刷新 `updatedAt`；list-edit 接入真实 context + 可失败权限表达式 + 拒绝路径组件测；account 失败可观察横幅；sessionProvider nil fail-closed；两 README 更新。**F-009-006/007 → `fixed`**（用户裁决纳入，I-009-001/002 resolved）：records 写路由挂 fail-closed admin 会话鉴权；PATCH body ≤ 4 KiB、`pageSize ≤ 100`。回归：`go test ./...` 绿、web vitest **398/398**、`tsc -b` 与 `npm run build` 通过。 |

**门禁状态**：A-001 全部 7 条 findings 已按 P-003 `fixed` 合法闭合（5 required + 2 recommended，用户书面裁决纳入 006/007）。剩余开放项：**阶段/关门审计未做**。**本目标仍未 `done`**：关门需先经阶段/关门审计（self 或 `/audit`）复核 fixed 证据。

---

## A-002 · 关门复审 · A-001 fixed 证据与成功标准（2026-08-01）

- **source**：`independent`
- **auditor**：GitHub Copilot（独立交叉审计 / `/audit`）
- **类型**：`close-out`（兼 `execution-facts` / `finding-closure`）
- **scope**：`GOAL-009-mvp-bugfix-followup` 关门前复核 — A-001 的 F-009-001～007 闭合证据、成功标准 5/5、I-009-001/002、工作区绑定与共享资料边界；**不**重开 VP-001 / Root `done`，**不**扩 schemaUrl/真实 IAM 等排除项
- **verdict**：`pass`
- **工作区上下文**：`workspace-001-mvp-admin-foundation` / Root `GOAL-001-mvp-admin-foundation` / canonical `docs/workspace-001-mvp-admin-foundation/` / `shared_materials_catalog: none` / `primary_plan: VP-001-mvp-admin-foundation`

### 范围与区间

| 项 | 内容 |
|----|------|
| **覆盖** | 00-meta 成功标准与路线图；01-decision D-001～D-004 与 I-009；02-execution 实施时间线；A-001 findings 闭合留痕；`apps/api` records/account 实现与测试；`apps/web` list-edit / account error / README；本轮独立复跑相关测试 |
| **不覆盖** | Root/VP 重开；生产 IAM / 请求级令牌；浏览器手测矩阵；附件基线长文重写（附件仍保留 A-001 时点 open 快照属预期） |
| **方法** | 只读通读五件套 + 源码定点核对 + `go test ./...` + 聚焦 vitest（list-edit + App.integration） |

### 工作区 / P-005 / 共享资料

| 核对项 | 结果 |
|--------|------|
| workspace Root / canonical | 匹配；本目标 `parent: GOAL-001-mvp-admin-foundation`，在 canonical 内 |
| plan_refs / primary_plan | `VP-001-mvp-admin-foundation`（与 workspace 一致） |
| 共享资料引用 | 无；catalog=`none`，未发现伪引用 |
| I-009-001 / I-009-002 | **resolved**（用户书面「都纳入实施」+ D-004 + 代码/测试）；无到期 required 信息项阻断关门 |
| progress 字段 | `5/5` 仅派生自 5 条 required 成功标准；**不得**单独推导 `done`（本意见亦不改 status） |

### 成果（有证据 · 本轮复核）

| 主张 | 证据路径 / 命令 |
|------|-----------------|
| F-009-001 PATCH 刷新 `updatedAt` | [records.go](../../../../apps/api/internal/handler/records.go) `rec.UpdatedAt = time.Now().UTC()`；`TestRecordsUpdateRefreshesUpdatedAt` |
| F-009-002 真实 context + 可失败权限 + 拒绝路径 | [list-edit-lifecycle-page.tsx](../../../../apps/web/src/app/examples/list-edit-lifecycle-page.tsx) `permissionCascade` + admin 表达式；App/registry 透传 `navigationContext`；[list-edit-lifecycle.test.tsx](../../../../apps/web/src/app/examples/list-edit-lifecycle.test.tsx) admin 启用 / viewer disabled |
| F-009-003 account 失败可观察 | [main.tsx](../../../../apps/web/src/main.tsx) 保留 error + `accountError`；[App.tsx](../../../../apps/web/src/app/App.tsx) `role="alert"` 横幅；`App.integration.test.tsx`；失败时 `loadAccountContext` → 空 context fail-closed |
| F-009-004 nil provider fail-closed | [account.go](../../../../apps/api/internal/handler/account.go) 注释与 nil→401；`TestAccountsMeNilProviderFailsClosed` |
| F-009-005 README 与现状一致 | [apps/api/README.md](../../../../apps/api/README.md)、[apps/web/README.md](../../../../apps/web/README.md) 端点 / 写路由鉴权 / 上限 / 测试命令 |
| F-009-006 写路由挂 Allow | `writeGate()` + PATCH/DELETE 调用；`TestRecordsWriteRequiresSession` / `TestRecordsWriteDeniedWithoutAdminRole` |
| F-009-007 body/pageSize 上限 | `MaxBytesReader` 4 KiB + `pageSize ≤ 100`；`TestRecordsUpdateBodyTooLarge` / `TestRecordsListPageSizeCap` |
| 回归（本轮独立执行） | `apps/api`: `go test ./...` **ok**（handler/account 包）；`apps/web`: vitest `list-edit-lifecycle.test.tsx` + `App.integration.test.tsx` → **7/7 pass**（2026-08-01） |

### 对照成功标准

| 成功标准 | 复核结论 | 备注 |
|----------|----------|------|
| F-009-001 | **满足 · fixed 可重复核对** | 严格 `After` 断言 |
| F-009-002 | **满足 · fixed 可重复核对** | 组件测覆盖拒绝路径；默认 `context={}` 亦 fail-closed（按钮 disabled） |
| F-009-003 | **满足 · fixed 可重复核对** | UI + console + 空 context 回落 |
| F-009-004 | **满足 · fixed 可重复核对** | 行为与注释一致 |
| F-009-005 | **满足 · fixed 可重复核对** | 与代码边界一致；见 recommended F-A002-001 措辞精度 |
| F-009-006/007（用户升格） | **满足 · fixed 可重复核对** | 达 A-001 关闭要求（dev session + 角色门槛 / 上限 + 测试） |

### Findings

#### F-A002-001 · 写路由「鉴权」措辞易被读成请求级凭证（recommended）

| 字段 | 值 |
|------|-----|
| **级别** | recommended |
| **严重度** | low |
| **影响门禁** | 不单独阻断关门 |
| **状态** | fixed |
| **描述** | `writeGate` 评估的是**服务端注入**的 `sessionProvider`（生产接线为 `StaticDevSession`，恒含 admin），**不是**按 HTTP 请求凭证/令牌鉴权。匿名客户端在默认进程配置下仍可 PATCH/DELETE 成功。A-001 关闭要求（挂 Allow + 可注入测拒绝）已满足；但 README「需 admin 会话」若脱离「静态开发会话」上下文，易被高估为网络侧鉴权。 |
| **证据** | `records.go` `recordsHandler` → `sessionProvider`；`account.StaticDevSession` roles 含 admin；写测通过 `recordsMuxWith` 注入 nil/editor 才得 401/403 |
| **建议** | `/govern` 响应时在两 README 鉴权边界各加一句：**写路由 gate 绑定进程内会话提供者，非请求头身份；默认静态 admin 会话下 HTTP 客户端无凭证仍可写**。可选：不改代码。 |
| **关联** | A-001 F-009-006 闭合边界澄清；非推翻 fixed |
| **闭合留痕** | 2026-08-01 `/govern`：`apps/api/README.md` 与 `apps/web/README.md` 鉴权边界各加「gate 绑定进程内会话提供者、非请求头身份、默认 admin 会话下无凭证仍可写」说明 |

#### F-A002-002 · 目标文档台账卫生滞后（recommended）

| 字段 | 值 |
|------|-----|
| **级别** | recommended |
| **严重度** | low |
| **影响门禁** | 不单独阻断关门（代码事实以闭合留痕与源码为准） |
| **状态** | fixed |
| **描述** | (1) [01-decision.md](01-decision.md)「信息需求」节仍写 I-009-001/002「仍 open」，与同文件 D-004 及 00-meta resolved 表矛盾。(2) A-001「编排提示」仍写「存在 5 条开放 required」，与后续响应表及 Findings `fixed` 矛盾（历史节未加删除线/响应覆盖说明）。(3) 附件 [audit-code-review-bugs-2026-08-01.md](attachments/audit-code-review-bugs-2026-08-01.md) findings 仍标 open——作为 A-001 **时点基线**可接受，但若读者只看附件会误判未修。 |
| **证据** | 01-decision 文首信息需求段；本文件 A-001「编排提示」；附件 §3 各 finding `状态: open` |
| **建议** | `/govern` 关门前轻量修订：决策文首与 A-001 编排提示对齐 resolved/fixed；附件顶加「实施后状态以 03-audit 闭合留痕为准」一句即可，无需重写长文 |
| **闭合留痕** | 2026-08-01 `/govern`：01-decision 文首 I-009 对齐 resolved；A-001「编排提示」加历史基线说明；附件顶加「实施后状态以 03-audit 为准」注 |

### 必改项汇总

- **开放 required：0**
- **开放 recommended：2**（F-A002-001、F-A002-002）— 文档边界/台账卫生，**不阻断**本目标关门

### 与既有意见的异同

| 对照 | 说明 |
|------|------|
| vs A-001 | A-001 为立项基线 `conditional`（未修前不得 done）。本 A-002 在实施后独立复核 7 条 `fixed` 证据，**同意** required 闭合充分；**不**推翻 A-001 历史 verdict，仅给出关门 scope 新 verdict=`pass` |
| vs 执行记录 | 02-execution 所述改动与源码一致；本轮复跑测试支持其回归主张（web 全量 398 未在本轮重跑，聚焦测 7/7 + go 全绿已覆盖本目标关键路径） |
| 冲突 | **无**与 self 意见冲突（本目标尚无 self 关门意见） |

### 结论 + 建议给编排器/用户的下一步

**结论**：在声明的 MVP bugfix 范围内，GOAL-009 成功标准与 A-001 七条 findings 的 `fixed` 闭合**证据充分、可重复核对**；信息门禁清空；无未闭合 required。独立交叉审计 **verdict = pass**。

**建议 `/govern`**：

1. **P-004**：已有 independent 关门意见（A-002 pass）、尚无 self 关门意见 → **询问**是否还要 self 关门审，或接受仅独立意见后授权 `done`。
2. 可选：顺手闭合 F-A002-001/002（README 一句 + 台账卫生）或 `accepted-residual` / 关门后跟踪。
3. 用户书面授权后：GOAL-009 → `status: done`；同步 goal-tree；**不**改 Root `done` / 纲领 6/6 / VP closed。
4. 浏览器手测仍为 optional，非本意见 required。

### 声明

本意见 `source: independent`，**不**修改目标 `status` / `progress` / 方案正文 / goal-tree 状态列。响应、finding 闭合与关门放行由 **`/govern`** 处理。

### 响应（2026-08-01 · `/govern` 关门）

| date | actor | scope | summary |
|------|-------|-------|---------|
| 2026-08-01 | `/govern` | A-002 关闭项 | 用户 P-004 裁决：**不补 self 关门审**，接受 A-002 `pass`；**F-A002-001 → `fixed`**（两 README 鉴权边界各加「gate 绑定进程内会话提供者，非请求头身份；默认静态 admin 会话下无凭证仍可写」）；**F-A002-002 → `fixed`**（01-decision 文首 I-009 对齐 resolved；A-001「编排提示」加历史基线说明；附件顶加「实施后状态以 03-audit 为准」注）。 |

**关门授权（P-004 留痕）**：用户书面授权 GOAL-009 `status: done`（2026-08-01）。依据：A-001 七条 findings 全部 `fixed`、A-002 独立关门复审 `pass`（开放 required 0）、I-009-001/002 resolved、回归 `go test` 全绿 + web 398/398 + build 通过。**不**改 Root `GOAL-001` `done` / 纲领 `6/6` / VP-001 `closed`；浏览器手测保持 optional。
