---
id: GOAL-003-r2-audit-event-model
doc: audit-entry
record_id: A-001
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R2 S0 · 审计事件模型范围、I-001/I-002 信息门禁、D-001 决策与 E-001 扫描事实
audit_type: design-plan+execution-facts
verdict: conditional
status: recorded
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# A-001 · R2 S0 设计/计划 + 扫描事实独立审计（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：stage（S0）/ design-plan + execution-facts
- **scope**：R2 S0：审计事件模型范围、I-001/I-002 信息门禁、D-001 决策与 E-001 扫描事实
- **verdict**：conditional

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）
- **covered**：GOAL-003 定义与非目标、S0 路线图、D-001 范围/门禁、E-001 已记录扫描事实、I-001/I-002 登记与阶段门禁、对照代码中的 operationlog 写入/读取面
- **excluded**：S1 实施、D-002 方案正文（尚不存在）、S2 接入验证、S3 关门、未重新执行测试套件、其他工作区上下文、共享资料内容（目录为 `none`）
- **本轮未复验**：未运行 `go test` / Web 测试；「测试已通过」仅核到测试代码与断言存在，运行结果标为证据不足

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-003 `parent`、`primary_plan` 一致 |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none`；GOAL-003 未引用 `material_id`/`sha256` |
| 对齐链 | 未发现与 Root R2 / VP-012 方向的明显冲突 | Root R2 = 结构化 diff / 脱敏 / correlation；GOAL-003 是该波次有界切片。VP-012 另列 session/effective actor、保留/归档，见 F-006 |
| Vision Review required | 本 scope 未见开放 required | `docs/vision/reviews.md` 索引声明 open required = 0；本意见不审 Vision Review 本身 |
| 既有 Goal 审计 | 无 | `03-audit.md` 索引与 `03-audit/` 在本条之前为空 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| 目标范围与非目标已写清，且挂 VP-012 | `00-meta.md` 范围/非目标/`plan_refs` | 通过 |
| D-001 冻结的是**范围与门禁**，不是 detail schema | `01-decision/D-001-r2-audit-event-model.md` 决策 1–4 | 通过；勿与 S0「冻结 D-001」误读为 schema 已冻（F-007） |
| I-001/I-002 已按 P-005 最小列登记，且均仍为 `collecting` | `00-meta.md` 信息表 | 通过；未伪装为 `verified` |
| `Operation` 含 event/actor/record/detail/`CorrelationID`/`CreatedAt` | `apps/api/internal/modules/operationlog/repository.go` `Operation` + `RecordOperation` | 通过；correlation 写入独立表 `operation_log_correlation` |
| `event` 受 SQLite CHECK 约束，扩事件需新 migration | `apps/api/internal/modules/operationlog/migration/migration.go` 各 `operationLog*DDL` | 通过 |
| auth login/logout/refresh detail 仅 `{username}` | `handler/auth.go` `authEvent`；`handler/operations_test.go` `TestOperationLogAuthEvents` | 写入路径通过；测试存在且断言「exactly {username} (no token/password/secret)」；**本轮未跑测试** |
| settings mutation detail 为 `{siteTitle, action}` | `handler/settings.go` `recordSettingsOperation` | 通过 |
| users create/update detail 主要为 `{username}`；delete 无 detail | `handler/users.go` `usersOnWrite` | 通过 |
| 抽样 MFA / 账户改密路径当前不把 secret/recovery/password 写入 detail | `handler/mfa.go` detail=`{userId}`；`handler/account_self.go` password-change detail 为空 | 通过（抽样，不是全量清单） |
| 读取 API 原样暴露 `detail` | `handler/operations.go` `operationToMap`；`handler/operations_export.go` CSV 含 `detail` | 通过；E-001 只写了 JSON 读取面，未写 CSV / 检索面（F-002） |
| R1 auth correlation 持久化测试存在 | `handler/operations_test.go` `TestR1CorrelationIDPersistsOnAuthOperation` | 测试代码可核对；**运行结果本轮未复验** |
| I-002 未被执行记录静默关闭 | E-001「不能由执行记录静默代替用户裁决」；I-002 仍 `collecting` | 通过 |

## 对照成功标准（S0 适用部分）

GOAL-003 四条成功标准均属 S1–S3 交付物。S0 只评估「是否已具备冻结 schema / 进入实施的信息」。

| 标准 | S0 状态 | 证据 |
|------|---------|------|
| 1. 新写入可被统一 schema 解析且带版本 | 未开始 | 无 envelope / schema 实现；D-002 不存在 |
| 2. 敏感字段无法经新事件明文写入或读取 | 未开始（仅有局部手工约定） | 无统一脱敏器；现有路径靠各 handler 手写 JSON |
| 3. auth/settings/users 至少各一条真实 mutation 消费模型并保留 correlation | 部分基线，未达标准 | auth/settings 已写 CorrelationID；**users 写路径不写**（F-003） |
| 4. 兼容读取、迁移/回滚与全量验证证据 | 未开始 | 读取面现状已部分扫描（F-002） |

## 信息门禁核对（P-005）

| ID | 级别 | 最晚阶段 | 状态 | 是否到期 | 本轮结论 |
|----|------|----------|------|----------|----------|
| I-001 | required | S0 结束前 | collecting | S0 仍「进行中」，**尚未逾期**；但**不足以关闭** | 影响门禁 = S1 方案冻结。E-001 自承 settings 其余字段与非核心事件仍待扫。00-meta 证据栏仍写「待 E-001」，与已落盘 E-001 不一致（F-005） |
| I-002 | required | S1 实施前 | collecting | **未到期**（不阻断 S0 结束） | 影响门禁 = S1 实施 / S3 关门。D-001 写「independent **或** cross」，模式未唯一；项目级默认 provider = grok-build，但仍缺本目标书面裁决（F-004） |

无 `deferred` 项。无用户书面 `accepted-residual`。不得把本条 A-001 的存在自动解释为 I-002 已关闭：本条只覆盖 S0 design-plan + execution-facts，不是 S3 关门独立审。

## Findings

### F-001 · I-001 现有扫描不足以关闭，S0 不得结束 / schema 不得冻结

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | S0 结束；I-001；S1 方案冻结 |
| evidence | `00-meta.md` I-001；`D-001` 决策 3 与不确定项；`E-001`「扫描结论」第 2 条；`settings/repository/repository.go` `SiteSettings`；`operationlog/repository.go` 事件常量 |

I-001 问的是：「现有事件 detail、敏感字段与 API 兼容边界**是否足以冻结 schema**」。E-001 的诚实结论是「否」：仍须补 settings 全字段与非核心事件敏感字段清单。`SiteSettings` 除 `siteTitle` 外还有 `LogoURL`/`LogoURLLight`/`LogoURLDark`/`FaviconURL`/`DefaultLocale`/`SiteTimezone`/`DefaultTheme`/`CopyrightText`/`ICPNumber` 等，均未进入 I-001 证据。事件常量覆盖 auth/users/roles/settings/account/data-transfer/files/dictionary/tasks/captcha/recycle/data-permission/mfa/wallet 等，E-001 只抽样了 auth/settings/users。

这不是「未知被伪装成已验证」——登记仍为 `collecting`，值得肯定。但按 D-001 决策 3 与 P-005 规划门禁，**在清单与兼容边界未落盘前不得把 I-001 标 verified，不得结束 S0，不得冻结 S1 schema**。

### F-002 · 读取/导出/检索面的兼容与泄露边界未写入扫描结论

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | I-001；S1 方案冻结 |
| evidence | `handler/operations.go` `operationToMap`；`handler/operations_export.go` headers；`operationlog/repository.go` `operationsWhere` Q 检索 |

E-001 只写「operation API 当前原样暴露 detail」。独立核验还表明：

1. `operationToMap` **不输出** `correlationId`（仓库已 JOIN 读出，但 API/CSV 丢弃）。
2. `GET /api/operations/export` 把原始 `detail` 写入 CSV。
3. 列表 `q` 对 `detail` 做 `instr(lower(COALESCE(detail,'')))`，即原始 detail 可被检索命中。

这三项直接回答 I-001 的「API 兼容边界」，也决定 schema 版本化后旧客户端、导出与检索如何 fail-closed。缺它们则 I-001 不能关。

### F-003 · D-001 消费切片中 users（及多数写路径）当前不持久化 correlation

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | I-001；S1 方案冻结（成功标准 3） |
| evidence | `handler/users.go` `usersOnWrite`（无 `CorrelationID`）；`handler/users_state.go`；`handler/mfa.go` `recordMFAEvent`；对比 `handler/auth.go`、`handler/settings.go` |

D-001 选定 auth / settings / users 三类真实 mutation。E-001 把「`Operation` 已有 CorrelationID + R1 auth round-trip 测试」写成已完成事实，但未记录：**users 写路径不写 CorrelationID**；MFA / account / roles / wallet 等同样不写。成功标准 3 要求三类路径都「保留 R1 correlation 关联」。S1 方案必须把 users（及明确列入切片的路径）接线列为交付，而不是假定 R1 已覆盖。

### F-004 · I-002 审计模式未唯一确定，S1 实施前仍须书面裁决

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | I-002；S1 实施；S3 关门 |
| evidence | `00-meta.md` I-002；`D-001` 不确定项；`docs/architecture/principles.md` P-003/P-004.1；`docs/architecture/independent-audit-execution.md` |

R2 含敏感字段脱敏与审计读取面，按 P-003 风险表属于 security/data，**最低模式可唯一判定为 `independent`**（`cross` 仅当用户要求多工具、或将其视为不可逆/跨边界时）。D-001 写成「independent **或** cross」，使门禁成本（是否强制 self）无法唯一确定，触发 P-004.1。项目级默认 provider 已指定 grok-build（grok-4.6 · high），但 I-002 仍写「待用户确认」，且无本目标书面裁决。

本条 A-001 **不能**关闭 I-002：它只证明 S0 已发生一次 independent 审，不是 S3 关门审，也不是用户对模式/provider 的书面确认。

### F-005 · I-001 证据栏未回写 E-001（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| 影响门禁 | 台账可追溯性（不单独阻断） |
| evidence | `00-meta.md` I-001「证据 / 结论」=「待 E-001」；`02-execution/E-001-r2-audit-surface-scan.md` 已 `recorded` |

E-001 已存在，信息表仍写「待 E-001」。应回写「扫描部分完成，结论见 E-001；仍 collecting」。不得借回写把状态改成 `verified`。

### F-006 · VP-012 审计能力中的 session/effective actor 与保留/归档未在 R2 显式延期（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| 影响门禁 | 无（有界切片，未见 P-006 明显冲突） |
| evidence | `docs/vision/plans/VP-012-shared-cross-module-contracts.md` 方向表；Root / GOAL-003 / D-001 均未点名延期 |

VP-012 方向级「审计事件模型增强」含 correlation/**session/effective actor** 与 **保留/归档触发**。Root R2 与 D-001 只收结构化 diff / 脱敏 / correlation，这是合法有界切片，但未写「其余能力不在 R2 / 待后续波次」。建议在 D-002 或 D-001 补记，避免后续把 VP 全文读成 R2 必交付。

### F-007 · 「D-001 accepted」不等于 schema 已冻结（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| 影响门禁 | 防误放行 S1 |
| evidence | `00-meta.md` S0「冻结 D-001」；`D-001` `status: accepted` 同时要求先关 I-001 |

路线图用语容易让编排器把 D-001 accepted 当成 S0 完成。独立核验：D-001 只冻结消费切片、兼容原则、fail-closed 默认与信息门禁顺序。

## 必改项汇总

1. **F-001**：补齐 I-001 可核对清单（至少：settings 全字段敏感性、非核心事件 detail 抽样/分类、敏感键默认集合），再决定能否 `verified`；在此之前 **禁止结束 S0 / 冻结 schema**。
2. **F-002**：把读取 API、CSV 导出、`q` 检索对 raw `detail` 的行为，以及 **correlationId 不在读出面**，写入 I-001/D-002。
3. **F-003**：在方案中写明 users（及任何列入 R2 切片的写路径）如何写入/校验 correlation；不得只引用 auth 测试。
4. **F-004**：S1 实施前由用户书面确认审计模式（建议 `independent`）与 provider（建议 grok-build）；`cross` 仅在用户明确要求时采用。闭合 I-002 后才能实施。

## 与既有意见的异同

无既有 self / independent 条目。本条为 GOAL-003 第一条意见。不与 GOAL-002 A-001（R1 self close-out，同区短 id）冲突。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** S0 扫描主干事实成立，且未把 I-001/I-002 伪装成已验证；但 I-001 仍不足以支持 schema 冻结，读取面与 users correlation 缺口未进清单，I-002 模式也未唯一裁决。

建议 `/govern`：

1. 响应本条 A-001；**不要**把 S0 标完成，**不要**进入 S1 实施。
2. 先补 I-001 清单（含 F-002/F-003），回写信息表（F-005），必要时落 D-002。
3. 按 P-004.1 请用户确认 I-002：建议模式 `independent`、provider grok-build；确认后才能实施。
4. 本意见不改 `status`/`progress`/goal-tree。

## 声明

本意见 `source: independent`，不修改目标 `status` / `progress` / 检查点 / 方案正文 / `goal-tree.md`。响应、finding 闭合与阶段推进由 `/govern` 处理。保证等级为框架默认 **L0**（入口分离），不得表述为第三方鉴证。
