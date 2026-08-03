---
title: 审计台账 · 语义化 Admin 资源替换与双实体验证
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-010-a002-schema-adapter
version: 0.4.0
---

# 审计台账 · GOAL-011

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | self | 2026-08-03 | S1 契约冻结（D-002 · I-011-001/002） | pass | 无 required；3 条 recommended |
| A-002 | independent | 2026-08-03 | S1 契约冻结（D-002 · I-011-001/002） | conditional | 2 required + 4 recommended，均已响应（D-003，见响应节） |
| A-003 | self | 2026-08-03 | S2 后端 users/roles 闭环（I-011-001 v0.2.0） | pass | 无 required；3 条 recommended |
| A-004 | independent | 2026-08-03 | S2 后端 users/roles 闭环（I-011-001 v0.2.0；progress 1/5 → 2/5） | pass | 无 required；F-001～F-003 fixed、F-004 handled（见响应节） |

## 当前审计边界

- S1 契约已冻结（D-002；`I-011-001`/`I-011-002` → verified，v0.2.0）；S2 后端闭环已实施（A-003 self · pass；A-004 independent · pass，recommended 已响应）；S3～S5 未实施。
- A-002（independent）对 S1 给出 **conditional**，两条 required（F-001 actor 通道、F-002 快照语义）经 **D-003** 走 `fixed` 闭合，F-003～F-006 采纳为 `handled`；A-001（pass）与 A-002 已趋同。
- A-003 / A-004 对 S2 同向 **pass**；本 scope 无开放 required；A-004 F-001～F-003 已 fixed、F-004 handled（见「响应 A-003 + A-004」节）；A-003 F-001（快照）/F-002（重启搜索排序）/F-003（password 长度）随 S3/S4 落实。
- `I-011-003` 仍 open required（最晚 S4）；到达最晚阶段前未 `verified` 或合规 `accepted-residual` 时阻断对应门禁。
- GOAL-010 与 Root A-002 的既有独立意见不复制到本台账。

## A-001 · S1 契约冻结自审（2026-08-03）

- **source**：self
- **auditor**：Claude Code（govern orchestrator）
- **类型 / scope**：stage · S1 契约冻结（D-002；I-011-001/002 → verified；S1 检查点达成）
- **verdict**：pass

### 范围与区间

审 GOAL-011 S1「语义资源与退场契约冻结」：两份版本化契约（`I-011-001-users-roles-contract.md` v0.1.0、`I-011-002-records-retirement.md` v0.1.0）、D-002 用户裁决、meta/decision/execution/goal-tree 同步；**不含** S2～S5 实施（未开始）。工作区绑定（workspace-002 / Root GOAL-001 / canonical）与共享资料（无）已核对，未 fail closed。

### 成果（有证据）

- **I-011-001 冻结**：users/roles 资源契约（端点、公开字段、敏感字段隔离、角色分配、self/最后管理员保护、system role 与 grant 约束、权限键/菜单/操作日志、错误码）+ 通用工厂最小扩展（`JSONFields` + `DomainError`）。依据现有代码逐项核对（`store.go` 表结构、`auth.go` 身份投影、`resources.go` 工厂、`operations.go` 事件 CHECK）。
- **I-011-002 冻结**：records 足迹盘点（11 面）、fresh/in-place 迁移矩阵（0005 operation_log 扩展 + 0006 DROP TABLE + 清理权限/菜单行）、硬退场数据处置（pre-v0006 快照兜底）、S3 验收口径。
- **用户裁决留痕（P-004）**：三项关键取舍均经用户确认（通用工厂+最小契约扩展、操作日志纳入、硬退场 DROP TABLE），记录于 D-002。
- **同步一致性**：meta v0.2.0（S1 勾选、progress 1/5、I-011-001/002 verified、I-011-003 open）、decision D-002、execution S1 条目、goal-tree v0.51.0（树 + 状态表 + 注记）四处一致。
- **未越权**：未修改任何产品代码（S1 为文档冻结）；未关 Root A-002 F-002-001；GOAL-010 保持 3/5。

### 对照成功标准

| S1 标准 | 状态 | 证据 |
|---------|------|------|
| 关闭 I-011-001/002 | ✅ verified | 两契约 v0.1.0 + D-002 |
| 冻结 users/roles 最小领域边界与安全不变量 | ✅ | I-011-001 §2/§3/§4/§6 |
| records 版本化退场与既有数据库迁移策略 | ✅ | I-011-002 §2/§3/§5 |

### Findings

- **F-001 · password 最小长度未冻结**（severity: low；建议: recommended；status: open）
  - 描述：users 资源 `password` 仅冻结「必填非空 string」，未设最小长度。既有系统（种子 admin dev 密码 "admin"）无长度先例；完整密码策略属非目标（I-011-001 §8）。
  - 证据：I-011-001 §2 create 字段、§8 非目标。
  - 影响：不阻断 S1/S2；S2 实施可自行加 ≥8 长度校验（不违背契约，仅收紧），或保持现状待后续扩展。

- **F-002 · fixture 文案残留风险**（severity: low；建议: recommended；status: open）
  - 描述：`data-table.json` 文案含 "records from the Go list API (/api/records)"；S3 改指 dataSource 时需同步改文案，否则 `api/records` 字符串残留。
  - 证据：`apps/api/internal/handler/fixtures/schema/data-table.json`；I-011-002 §3.3。
  - 影响：I-011-002 §5 验收口径「grep 无 `api/records` 残留」已覆盖；S3 执行时注意文案一并更新。

- **F-003 · 工厂 DomainError 检查优先级**（severity: low；建议: recommended；status: open）
  - 描述：通用工厂新增 DomainError 映射后，create/update 错误路径须在 `store.ErrNotFound`/通用 500 兜底**之前**先识别 DomainError，否则 users/roles 领域拒绝退化成 INTERNAL。
  - 证据：`resources.go` create/update/delete/detail 现有错误映射；I-011-001 §7.2。
  - 影响：S2 实现顺序要求；建议作为 S2 自审复核点。

### 必改项汇总（required 列表）

无。

### 结论 + 建议下一步

S1 契约冻结完整、可实施、与既有事实一致；无未闭合 required。**pass**。下一步：S2 后端 users/roles 资源闭环（通用工厂扩展 + store 领域方法 + 双资源 CRUD + 401/403 负向路径）。F-001～F-003 为 recommended，随 S2/S3 实施落实。按用户指令，本自审后将调用 **grok build 独立交叉审计**（scope: S1 契约冻结），等待其意见后合并响应。

## A-002 · S1 契约冻结独立交叉审计（2026-08-03）

- **source**：independent
- **auditor**：grok build
- **类型 / scope**：design-plan · S1 契约冻结（D-002；I-011-001 / I-011-002 → verified；S1 检查点达成，progress 0/5 → 1/5）
- **verdict**：conditional

### 范围与区间

审 GOAL-011 S1「语义资源与退场契约冻结」是否足以无条件支撑后续 S2/S3 实施，关注点：

1. users/roles 领域契约是否完整、安全、可实施；
2. records 退场与迁移策略是否可追溯、不破坏既有迁移账本/历史事实；
3. 对 I-010-001 §5「不引入 409」的限定偏离是否被清晰记录；
4. 信息门禁闭合是否合规（P-005）；
5. 用户裁决（P-004）是否留痕。

**不含** S2～S5 实施事实（产品代码未改，与 `02-execution` 一致）。工作区：`workspace-002-production-admin-foundation` / Root `GOAL-001-production-admin-foundation` / `canonical_scope` 已校验；`shared_materials_catalog: none`，未将共享资料当作关闭证据。

**只读核验**：GOAL-011 五件套与附件；父契约 I-010-001 v0.2.1；`resources.go` / `records.go` / `migrate.go` / `store.go` / `operations.go` / `auth.go`（及 `account/session.go` StaticDevSession 对照）。

### 成果（有证据）

| 主张 | 证据 | 核验结论 |
|------|------|----------|
| 两份版本化契约已落盘并被 D-002 采纳 | `attachments/I-011-001-users-roles-contract.md` v0.1.0；`attachments/I-011-002-records-retirement.md` v0.1.0；`01-decision` D-002 | 成立 |
| users 公开字段 / 敏感隔离 / username UNIQUE / 角色双写依据真实代码 | users 表含 `password_hash` UNIQUE username（`migrate.go` r2BaselineDDL）；`CreateUser`+`linkUserRole` 双写；`userWithRoles` 集合一致性；`accountFromUser` 仅投影 ID/Name/Roles/Permissions（无 hash） | 领域边界与代码对齐；敏感隔离方向正确 |
| roles system / ON DELETE RESTRICT / roleKeyRe / id=`role-{key}` | `roles.system` CHECK；`user_roles.role_id ON DELETE RESTRICT`；`roleKeyRe`；`ensureRole` id 派生；seed 将 admin/editor/viewer 升 `system=1` | 保护与格式主张可落地 |
| 操作日志扩展路径正确 | `operations.go` 现事件仅 records.* + auth.*；0004 CHECK 字面匹配；契约 0005 重建 CHECK 并**保留** records.* / auth.* | 可追溯；不改写 0001～0004 |
| records 退场不改历史账本 | I-011-002 硬约束 + 0006 新迁移 DROP；足迹 11 面与仓库现状（handler/store/fixture/manifest/e2e 等）方向一致 | 0001～0004 不可变策略成立 |
| 409 限定偏离有书面记录 | I-011-001 §6 表格式列出 409 码；D-002 第 1 项明确「仅账号域 409，envelope 形状不变」 | 关注点③满足 |
| P-004 三项裁决留痕 | D-002「用户裁决」三条（工厂扩展 / 操作日志纳入 / 硬退场 DROP）均 accepted | 关注点⑤满足 |
| P-005 S1 门禁 | `I-011-001`/`I-011-002` verified + 证据路径；`I-011-003` open 且最晚 S4、未伪关 | 关注点④主体合规 |
| meta / decision / execution / goal-tree 同步 | progress `1/5`、S1 勾选、信息表 verified、goal-tree 注记 | 文档投影一致（本意见不将其当作放行依据） |
| 未越权改代码 / 未关 Root F-002-001 | `02-execution` + 代码仍为 records 工厂实例 | 成立 |

### 对照成功标准

| S1 标准 | 本意见 | 说明 |
|---------|--------|------|
| 关闭 I-011-001/002 | 有条件成立 | 契约 + D-002 齐备；但契约内存在阻断 S2/S3 无歧义实施的缺口（见 F-001/F-002） |
| 冻结 users/roles 最小领域边界与安全不变量 | 大体成立 | users §2 较完整；roles 响应形状与工厂 actor 通道不足 |
| records 版本化退场与既有库迁移策略 | 大体成立 | 0001～0004 保护与硬退场清晰；**pre-v0006 验收与现有快照机制不一致** |

### Findings

- **F-001 · 通用工厂扩展未覆盖 SELF_OPERATION 所需的操作者身份通道**（severity: medium；建议: **required**；status: open；关联 I-011-001）
  - 描述：I-011-001 §2.4 要求「不能删除自己 / 不能移除自己的 admin 角色」→ `SELF_OPERATION`(409)。§7 仅冻结 `JSONFields` + `DomainError` 两项工厂扩展。当前 `ResourceEntity` 签名为 `Create/Update/Delete` **无** `context` / actor；工厂 `create/update/delete` 虽从 `requirePermission` 取得 `account.User`，但只传给事后 `OnWrite`，**不**传入 entity。结果：在「users 仍走通用五路由」前提下，S2 无法按冻结清单实现 self 保护，只能静默再扩接口、塞全局状态或旁路工厂——削弱 D-002「均走通用工厂」主张。
  - 证据：`apps/api/internal/handler/resources.go`（`ResourceEntity` 接口 L49–55；create/update/delete 仅 `OnWrite` 持有 user）；I-011-001 §2.4、§7；D-002 第 1 项。
  - 建议闭合：在 I-011-001 增补工厂扩展（例如 entity 方法接收 `context.Context` 且允许 `auth.IdentityFrom`，或 `ResourceEntity` 显式 `Actor account.User`），并写明 DomainError 检查优先级（承接 A-001 F-003）；**或**书面裁定 self 保护落在工厂层钩子并冻结钩子形状。修订后升版本 + D 响应。

- **F-002 · pre-v0006 快照验收与 migrate 现有「仅 first-pending 快照」机制冲突**（severity: medium；建议: **required**；status: open；关联 I-011-002）
  - 描述：I-011-002 §2.3 Path B 写「`pre-v0005` → 0005 → `pre-v0006` → 0006」；§5 验收要求「从 0004 基线到 0006 后 … `pre-v0006` 快照存在」。但 `migrate()` 仅在 `pending[0]` 前调用一次 `snapshotBeforePending`。当 0005 与 0006 **同批 pending**（S3 合入后、库仍停在 0004 的常见升级）时，只会生成 **pre-v0005**，不会生成 pre-v0006，按字面验收将失败。pre-v0005 虽仍含 records 表（数据可恢复），但契约与验收口径名不副实。
  - 证据：`apps/api/internal/store/migrate.go` L249–263（仅 first pending ≥2 快照）；I-011-002 §2.3、§2.4、§5。
  - 建议闭合（三选一，须写入契约并改验收句）：
    1. S2/S3 实现改为**每个**待应用版本前快照（或至少在 0006 前强制快照）；或
    2. 验收改为「存在可恢复 records 的升级前快照（pre-v0005 或 pre-v0006，以实际 first-pending 为准）」；或
    3. 强制分开发布/分步升级路径并在验收夹具中固定「先只 pending 0006」。

- **F-003 · roles 公开响应形状未像 users 字面冻结**（severity: low；建议: recommended；status: open；关联 I-011-001）
  - 描述：users §2.1 给出固定 JSON 行形状（含毫秒 `createdAt`/`updatedAt`）；roles §3 仅表格式字段/保护，**未**冻结 list/detail 响应是否暴露 `system`、时间戳形状、以及 `system` 的 JSON 类型（bool vs 0/1）。S2 负向/正向断言与 S4 Schema 列字段缺少金标准；UI 也无法在无 `system` 字段时禁用 system 角色的编辑控件（仅能依赖 409）。
  - 证据：I-011-001 §2.1 vs §3；`roles` 表含 `system`/`created_at`/`updated_at`（`migrate.go` rbacExpandDDL）。
  - 影响：不否定 S1 主体，但建议在进 S2 前补一小节响应形状（推荐暴露 `system: boolean` + 毫秒时间戳）。

- **F-004 · 既有 `CreateUser`/`linkUserRole` 会隐式 ensureRole，与 API「不隐式建角色」冲突**（severity: low；建议: recommended；status: open；关联 I-011-001）
  - 描述：契约 §2.3 正确要求未注册 key → `INVALID_ROLE_REF`、users 写不隐式创建角色。但现有 `linkUserRole` → `ensureRole` 会 `INSERT … system=0`。S2 若直接复用 `CreateUser` 做管理 API 将违反契约。契约未点名「禁止复用 linkUserRole 于 API 路径 / 须另写仅链接已存在角色的事务」。
  - 证据：`store.go` CreateUser L181–184；`migrate.go` linkUserRole L582–596、ensureRole L568–577；I-011-001 §2.3。
  - 影响：recommended 实施警示；闭合可选——契约加一句「API 路径不得调用 ensureRole」或 S2 实现时另写方法并在自审引用本 finding。

- **F-005 · 父契约 I-010-001 §5 仍写「不引入 409」，跨目标双真相残留**（severity: low；建议: recommended；status: open）
  - 描述：I-011-001 §6 / D-002 已清晰记录账号域 409 限定偏离，**本目标内**充分。但权威父契约 I-010-001 v0.2.1 §5 仍写「不引入 409/业务唯一冲突」，无指向 GOAL-011 的修订注记。实施者若只读父契约会与子契约冲突。
  - 证据：I-010-001 §5 L64；I-011-001 §6；D-002 第 1 项。
  - 影响：recommended——经 GOAL-010 决策追加修订注记（不改 S1～S3 records 零 API 历史事实），声明账号域 409 由 I-011-001 限定扩展。

- **F-006 · 承接 A-001 recommended（password 长度 / fixture 文案 / DomainError 优先级）**（severity: low；建议: recommended；status: open）
  - 描述：独立复核同意 A-001 F-001～F-003 的 low recommended 判断；其中 DomainError 优先级与本意见 F-001 工厂扩展补丁应一并冻结。password 最小长度属非目标边界，可不升 required。fixture 文案属 S3 grep 验收子集。
  - 证据：A-001 Findings；`resources.go` create 错误路径 L329–331 现将非 `ErrRecordExists` 一律 INTERNAL。

### 必改项汇总

| ID | 严重度 | 摘要 | 建议闭合路径 |
|----|--------|------|--------------|
| **F-001** | medium · required | 工厂扩展未规定 actor/context → SELF_OPERATION 不可在通用工厂内诚实实现 | 修订 I-011-001 §7（+ 可选 §2.4 实现注记）升版本；D 响应 |
| **F-002** | medium · required | pre-v0006 验收与 first-pending-only 快照机制不一致 | 修订 I-011-002 §2.3/§5 或规定实现改快照行为；D 响应 |

无 high required。F-003～F-006 为 recommended，不单独阻断，但 F-001 闭合时宜顺带处理 DomainError 优先级（F-006/A-001 F-003）。

### 与既有意见的异同

| 项 | A-001（self · pass） | A-002（independent · conditional） |
|----|----------------------|-------------------------------------|
| 契约落盘 / D-002 / P-004 / P-005 主体 | 通过 | **同意** |
| 409 偏离记录 | 隐含于成果 | **显式复核通过**（关注点③） |
| 0001～0004 不可变 + 硬退场 | 通过 | **同意** |
| 工厂可实施性 | 视为完整（仅 F-003 DomainError 优先级 recommended） | **不同意完整**：缺 actor 通道 → **F-001 required** |
| 迁移快照验收 | 未挑战 pre-v0006 | **F-002 required**（与 `migrate.go` 行为冲突） |
| password / fixture 文案 | F-001/F-002 recommended | 并入 F-006，维持 recommended |
| roles 响应形状 / linkUserRole / 父契约 409 注记 | 未提 | 新增 F-003/F-004/F-005 recommended |

同 scope 下 self=pass 与 independent=conditional **verdict 冲突** → 按 P-004 §3.2 须用户裁决（展示冲突 + 建议 + 留痕）；本意见不自动改状态。

### 结论 + 建议给编排器/用户的下一步

**verdict: conditional**——S1 主体工作（契约落盘、用户三项裁决、I-011-001/002 verified 投影、409 偏离书面化、迁移账本保护）**成立且可核对**；但存在 **2 条 medium required**：工厂扩展对 self 保护不完备（F-001）、退场快照验收与现实现机制不一致（F-002）。按 P-003，在 required 合法闭合前**不得**将 S1 视为无条件「可无歧义开工 S2/S3」的放行态（进度数字本身不构成放行）。

**建议下一步（/govern）**：

1. 汇总 A-001 + A-002；因 verdict 冲突执行 P-004：建议走 **fixed** 修订两契约（不建议 residual——缺口小、修文档成本低）。
2. 闭合 F-001/F-002（及可选 F-003 响应形状）→ 契约升小版本 → 本台账响应节留痕。
3. 再评估是否放行 S2（users/roles 后端闭环）；S3 退场实施须以 F-002 闭合后的验收句为准。
4. A-001 F-001～F-003 / 本意见 F-003～F-006 随 S2/S3 落实或 handled。

### 声明

本意见 **source: independent**，仅追加审计台账；**不修改** `00-meta` 的 status / 检查点 / 派生 progress，**不修改** goal-tree 状态列，**不修改**契约正文。响应、finding 闭合与阶段推进归 **`/govern`**。

## 响应 A-002（self · 编排响应 · 2026-08-03 · GOAL-011 D-003）

响应 A-001 与 A-002（同 scope，verdict 冲突按 P-004 §3.2 用户裁决「全部 fixed」）。

### 关闭证据表

| Finding | 严重度 | 状态 | 证据路径 |
|---------|--------|------|----------|
| A-002 F-001 · actor 通道 | medium · required | **fixed** | I-011-001 **v0.2.0** §7.2（`ResourceEntity` Create/Update/Delete 增传 `account.User`）+ §2.4 实现依托；GOAL-011 D-003 |
| A-002 F-002 · pre-v0006 快照语义 | medium · required | **fixed** | I-011-002 **v0.2.0** §2.3（每待应用数据变更迁移前快照，0005+0006 同批时 `pre-v0006` 必存在）+ §5 验收对齐；GOAL-011 D-003 |
| A-002 F-003 · roles 响应形状 | low · recommended | **handled** | I-011-001 **v0.2.0** §3.0（`system:boolean` + 毫秒时间戳）；GOAL-011 D-003 |
| A-002 F-004 · 禁 ensureRole 隐式建角色 | low · recommended | **handled** | I-011-001 **v0.2.0** §2.3 实现约束（API 路径另写仅链接已存在角色的事务）；GOAL-011 D-003 |
| A-002 F-005 · 父契约 409 双真相 | low · recommended | **handled** | GOAL-010 **D-005** + I-010-001 **v0.2.2** §5 注记（账号域 409 限定扩展，不改 S1～S3 历史） |
| A-002 F-006 · 承接 A-001 recommended | low · recommended | **handled** | A-001 F-003（DomainError 优先级）并入 I-011-001 v0.2.0 §7.3；A-001 F-001（password 长度）/F-002（fixture 文案）维持 recommended，随 S2/S3 落实 |

### 趋同说明

A-001（pass）与 A-002（conditional）的差异（工厂可实施性、快照验收）经 v0.2.0 修订闭合：A-002 两条 required 已 fixed，F-003~F-006 已 handled；S1 无开放 required，A-001/A-002 趋同为「可无歧义开工 S2/S3」。

### 仍开放项

- `I-011-003`（双资源集成验收契约）保持 `open`（最晚 S4 前，未到期）。
- A-001 F-001（password 最小长度）、A-001 F-002（fixture 文案）为 recommended，随 S2/S3 实施落实（非门禁阻断）。

### 后续

S2 后端 users/roles 资源闭环（按 I-011-001 v0.2.0 落地工厂扩展 + store 领域方法 + 双资源 CRUD + 401/403 负向路径）；S2 完成后按关键节点自审 + grok build 独立审计。

## A-003 · S2 后端 users/roles 闭环自审（2026-08-03）

- **source**：self
- **auditor**：Claude Code（govern orchestrator）
- **类型 / scope**：stage · S2 后端 users/roles 资源闭环（I-011-001 v0.2.0；S2 检查点达成）
- **verdict**：pass

### 范围与区间

审 GOAL-011 S2「在通用资源工厂之上实现 users 与 roles 的持久化 list/search/detail/create/update/delete、字段校验、敏感字段隔离、关系/系统角色保护、稳定错误 envelope 与 401/403 负向路径」。对照 I-011-001 v0.2.0；**不含** S3 退场（records 仍注册）与 S4 前端接入（未开始）。工作区绑定与共享资料（无）已核对。

### 成果（有证据）

- **工厂扩展**：`ResourceEntity` actor 通道（A-002 F-001）、`Resource.JSONFields`、`DomainError` + `writeEntityError`（DomainError → ErrNotFound → INTERNAL 顺序，A-002 F-006）；records 实体兼容（零对外变化，`go test` 既有 records 测试全绿）。
- **users 资源**：CRUD + 敏感字段隔离（`password_hash` 永不出响应，负向断言）、username 唯一（409 USERNAME_TAKEN）、self 保护（409 SELF_OPERATION）、last-admin 保护（store 层测试，非 admin actor）、角色分配校验（400 INVALID_ROLE_REF，不隐式建角色 A-002 F-004）、password 仅写（bcrypt，创建后可登录）、users.* 操作日志。
- **roles 资源**：CRUD + key 格式（400 INVALID_ROLE_KEY）、重复（409 ROLE_KEY_TAKEN）、system 保护（409 ROLE_SYSTEM）、in-use 保护（409 ROLE_IN_USE）、响应形状 `system:boolean` + 毫秒时间戳（A-002 F-003）、roles.* 操作日志。
- **迁移 0005**：operation_log event CHECK 扩展（users/roles 事件，保留 records/auth 历史）；重建迁移保留既有行（`TestMigrateExistingV3ToV4` 升级路径验证）。
- **种子/注册**：users/roles 权限、菜单、grants（admin rw、editor/viewer ro）+ `/api/users`、`/api/roles` 注册 + `StaticDevSession` 同步；records 种子保持（S3 退场）。
- **回归**：`go test ./...` 全绿（151 测试函数，含新增 20+ 用例）+ `go vet` 干净；web `vitest` 481/481 + `tsc -b` 干净（后端变更未破坏前端）。
- 修复实现中发现的两个缺陷：单连接嵌套查询死锁（`ListUsers` 先收行再 reconcile roles）、`updated_at` INTEGER 扫描 time.Time 类型错误（`UpdateUser`）——均已在 S2 内修复并有测试覆盖。

### 对照成功标准

| S2 标准 | 状态 | 证据 |
|---------|------|------|
| 通用工厂之上五路由 CRUD | ✅ | users/roles 经 `registerResource`；genericity 测试保持 |
| 字段校验 | ✅ | INVALID_CREATE_FIELD/PATCH_FIELD、INVALID_ROLE_REF、INVALID_ROLE_KEY |
| 敏感字段隔离 | ✅ | password_hash 负向断言（list/detail/create） |
| 关系/系统角色保护 | ✅ | SELF_OPERATION/LAST_ADMIN/ROLE_SYSTEM/ROLE_IN_USE |
| 稳定错误 envelope | ✅ | DomainError → `{error,message}`，409 账号域扩展 |
| 401/403 负向路径 | ✅ | anonymous 401、viewer write 403 |

### Findings

- **F-001 · 每待应用版本快照（A-002 F-002）顺延 S3**（severity: low；建议: recommended；status: open；关联 I-011-002 v0.2.0）
  - 描述：per-pending-version 快照（0005+0006 同批时 `pre-v0006` 必存在）需在 0006 存在时才可验收；S2 仅新增 0005，单快照行为自洽。S3 落地 0006 时必须改 `snapshotBeforePending` 为每待应用数据变更迁移前快照。
  - 证据：I-011-002 §2.3 v0.2.0；`migrate.go` `snapshotBeforePending`。
  - 影响：不阻断 S2/S3；作为 S3 实施要求记录。

- **F-002 · users/roles 重启持久化与排序/搜索正例测试顺延 S4**（severity: low；建议: recommended；status: open；关联 I-011-003）
  - 描述：users/roles 与 records 共用同一 SQLite store，重启持久化机制由 `records_restart_test.go` 覆盖；formal 双资源重启/排序/搜索正例属 I-011-003 S4 验收矩阵，S2 未单列。
  - 证据：`records_restart_test.go`；I-011-003（open）。
  - 影响：不阻断 S2；S4 验收矩阵须包含 users/roles 重启 + list sort/search 正例。

- **F-003 · password 最小长度仍未冻结（承接 A-001 F-001）**（severity: low；建议: recommended；status: open）
  - 描述：users 资源 password 仍为「必填非空」，未设最小长度；S2 按契约执行。非目标边界内的未来收紧项。
  - 影响：不阻断；如需可后续在契约小版本加 min-length 并补反例。

### 必改项汇总（required 列表）

无。

### 结论 + 建议下一步

S2 后端 users/roles 闭环完整、契约对齐、回归全绿；无未闭合 required。**pass**。下一步：S3 records 产品运行面退场（0006 DROP TABLE + 每待应用版本前快照 F-001 落地 + 权限/菜单/fixture/前端/测试退场）。按用户指令，本自审后将调用 **grok build 独立交叉审计**（scope: S2 后端闭环），等待其意见后合并响应。

## A-004 · S2 后端 users/roles 闭环独立交叉审计（2026-08-03）

- **source**：independent
- **auditor**：grok build
- **类型 / scope**：execution-facts · S2 后端 users/roles 资源闭环（GOAL-011 S2；I-011-001 v0.2.0；S2 检查点达成；progress 1/5 → 2/5）
- **verdict**：pass

### 范围与区间

审 GOAL-011 S2「在通用资源工厂之上实现 users 与 roles 的持久化五路由 CRUD、字段校验、敏感字段隔离、关系/系统角色保护、稳定错误 envelope 与 401/403 负向路径」是否与 I-011-001 v0.2.0 **真实落地**（非仅声明）。关注点：

1. 通用工厂扩展（actor 通道 / `JSONFields` / `DomainError` 映射顺序）是否按契约落地且 records **零对外变更**；
2. users/roles 领域不变量（敏感字段隔离、self/last-admin、system/in-use、不隐式建角色）是否真实实现；
3. 错误码 / envelope 与 I-011-001 §6 一致；
4. 401/403 负向路径与权限投影（种子 grants、`StaticDevSession`）一致；
5. migration `0005` 是否正确保留既有 `operation_log` 行与 checksum 账本；
6. 测试证据是否充分（`go test` / `vet` / web）。

**不含** S3 records 退场（records 仍注册，符合计划）与 S4 前端双资源接入 / `I-011-003`（open，最晚 S4，未到期）。

工作区：`workspace-002-production-admin-foundation` / Root `GOAL-001-production-admin-foundation` / `canonical_scope` 已校验；`shared_materials_catalog: none`，未将共享资料当作关闭证据。未读取或比较其他工作区。

**只读核验**：GOAL-011 五件套与附件（I-011-001/002 v0.2.0）；父契约 I-010-001 v0.2.2 §5 注记；`resources.go` / `users.go` / `roles.go` / `records.go` / `health.go`；`store/{users,roles,migrate,seed,operations}.go`；`account/session.go`；handler/store 测试套件；本机重跑 `go test ./...`（apps/api 全绿）+ `go vet ./...` 干净 + web `vitest run` **481/481** + `tsc -b` 干净。

### 成果（有证据）

| 主张 | 证据 | 核验结论 |
|------|------|----------|
| **① 工厂 actor 通道**（A-002 F-001） | `resources.go` `ResourceEntity.Create/Update/Delete(..., user account.User)`；create/update/delete 将 `requirePermission` 取得的 user 传入 entity；records 签名补 `_ account.User` 并忽略 | **成立**；SELF_OPERATION 可诚实实现 |
| **① `JSONFields`** | `Resource.JSONFields`；`decodeResourceCreate/Patch` 原始 JSON 透传；`usersResource` `JSONFields: ["roles"]` | **成立** |
| **① `DomainError` 映射顺序** | `writeEntityError`：先 `errors.As(*DomainError)` → 再 `store.ErrNotFound` → 最后 INTERNAL；create 路径 `ErrRecordExists` 重试后走 `writeEntityError` | **成立**（A-001 F-003 / A-002 F-006） |
| **① records 零对外变更** | records 实体仅补 actor 参数；既有 `records_test` 全绿（本机重跑 handler/store records 套件 PASS） | **成立** |
| **② 敏感字段隔离** | `userToMap` 永不含 `password`/`password_hash`；`TestUsersListAndDetail` / create 负向断言；password bcrypt 仅写 + `TestUsersPasswordWriteOnly` 可登录 | **成立** |
| **② self / last-admin** | `DeleteUser`/`UpdateUser` 用 `actorID`；self-delete/self-demote → `ErrSelfOperation`；最后 admin → `ErrLastAdmin`；HTTP `TestUsersSelfProtection`；store `TestUsersLastAdminProtection` | **成立**（LAST_ADMIN HTTP 层见 F-002） |
| **② system / in-use** | `UpdateRole`/`DeleteRole`：`system` → `ErrRoleSystem`；`user_roles` 计数 → `ErrRoleInUse`；handler/store 保护测试全绿 | **成立** |
| **② 不隐式建角色**（A-002 F-004） | `CreateUserManagement`/`UpdateUser` 先 `SELECT COUNT(*) FROM roles`，再 `INSERT user_roles`，**不**调用 `ensureRole`；`TestUsersCreateManagementRoleValidation` 断言 ghost 角色 count=0 | **成立** |
| **③ 错误码 / envelope** | users/roles `map*StoreError` → §6 码表（`USERNAME_TAKEN`/`SELF_OPERATION`/`LAST_ADMIN`/`INVALID_ROLE_REF`/`ROLE_*`/`INVALID_ROLE_KEY`/`USER|ROLE_NOT_FOUND`）；`writeError` 形状 `{error,message}` | **成立**；与 I-010-001 v0.2.2 账号域 409 注记一致 |
| **④ 401/403 + 种子 / StaticDev** | `seed.go`：admin 四权限 + 两菜单；editor/viewer 只读 `users.read`/`roles.read`；`StaticDevSession` 同步 permissions + `menu_users`/`menu_roles`；`TestUsersAuthGates`/`TestRolesAuthGates` 匿名 401 + viewer 写 403 | **成立**（dev session 断言覆盖见 F-003） |
| **⑤ migration 0005 实现** | `compiledMigrations` 追加 `0005:operation-log-expand:v1`；`migrate0005` = RENAME → CREATE expanded CHECK（保留 records.* / auth.*）→ `INSERT…SELECT` 全列 → DROP old → 重建索引；`up` 在单事务 + ledger checksum 插入 | **实现正确**；账本 checksum 机制未改写 0001～0004 |
| **⑤ checksum 账本** | `migrationChecksum(stmts+transformID)`；`TestMigrateFailClosedChecksumDrift` 等 fail-closed 仍 PASS；升级路径 `TestMigrateExistingV3ToV4` 断言 applied 含 v4+v5 | **账本侧成立**；行级保留证据见 F-001 |
| **⑥ 回归** | 本机：`go test ./...`（apps/api）全绿 + `go vet ./...` 干净；web `vitest run` **481/481** + `tsc -b` 干净；users/roles handler 与 store 新增套件均 PASS | **成立**（与 02-execution 主张同向；测试函数计数口径见 F-004） |
| 注册面 | `health.go` 注册 records + users + roles | **成立**；records 保留至 S3 |
| 操作日志 | `EventUser*`/`EventRole*`；`usersOnWrite`/`rolesOnWrite` detail 仅 username/key；handler 操作日志测试 | **成立** |

### 对照成功标准

| S2 标准（00-meta） | 本意见 | 说明 |
|--------------------|--------|------|
| 通用工厂之上 list/search/detail/create/update/delete | ✅ | users/roles 经 `registerResource`；sortFields/qSearch 按契约 |
| 字段校验 | ✅ | INVALID_CREATE_FIELD / INVALID_ROLE_REF / INVALID_ROLE_KEY 等 |
| 敏感字段隔离 | ✅ | password_hash 响应负向断言 |
| 关系/系统角色保护 | ✅ | SELF_OPERATION / LAST_ADMIN / ROLE_SYSTEM / ROLE_IN_USE |
| 稳定错误 envelope | ✅ | DomainError 逐字映射；409 仅账号域 |
| 401/403 负向路径 | ✅ | 匿名 401 UNAUTHENTICATED；viewer 写 403 |

### Findings

- **F-001 · 0005 既有 operation_log 行保留缺少专用升级回归**（severity: low；建议: **recommended**；status: open；关联 I-011-001 §5 / I-011-002）
  - 描述：`migrate0005` 的 `INSERT…SELECT` **代码路径正确**，可保留既有行；但测试套件**没有**「0004 态写入 `records.*`/`auth.*` 行 → Open 升级到 0005 → 断言行内容与 id 完整保留、且新事件可写入」的专用夹具。A-003 将 `TestMigrateExistingV3ToV4` 表述为「重建迁移保留既有行」**证据指针偏宽**：该测只断言 applied 版本链到 0005 与表存在，**未**插入/核对 operation_log 行。
  - 证据：`migrate.go` `migrate0005` L246–265；`operations_test.go` `TestMigrateExistingV3ToV4` L87–127（无 op 行往返）；`TestOperationLogAppendAndList` 仅在 fresh 终态写行。
  - 影响：不否定实现正确性，不阻断 S2；建议 S3 前或随 0005 加固补一条升级行保留测试，并收窄自审证据表述。
  - 建议闭合：补 `TestMigrate0005PreservesOperationLogRows`（或等价）→ 本 finding → fixed。

- **F-002 · LAST_ADMIN 无 HTTP 层 409 断言**（severity: low；建议: **recommended**；status: open；关联 I-011-001 §2.4/§6）
  - 描述：store `TestUsersLastAdminProtection` 覆盖 demote/delete 最后 admin → `ErrLastAdmin`；handler 有 `mapUserStoreError` → `LAST_ADMIN` 409。但 **HTTP 负向用例仅覆盖 `SELF_OPERATION`**（`TestUsersSelfProtection`），未构造「非 self 的最后 admin 删除/降级」→ 409 `LAST_ADMIN`。映射路径与 SELF_OPERATION 同构，漏映射风险低。
  - 证据：`store/users_test.go` `TestUsersLastAdminProtection`；`handler/users_test.go` 无 `LAST_ADMIN` 字符串断言。
  - 影响：recommended 覆盖加固；不阻断 S2。

- **F-003 · StaticDevSession users/roles 投影缺回归断言**（severity: low；建议: **recommended**；status: open；关联 I-011-001 §4）
  - 描述：`account.StaticDevSession` 已同步 `users.read/write`、`roles.read/write` 与 `menu_users`/`menu_roles`（实现与契约一致）。但 `TestAccountsMeDevSessionFallback` 仅断言 `menu_list_edit_lifecycle`，**未**断言 users/roles 权限键或新菜单 feature——静默回退时测试不会失败。
  - 证据：`account/session.go` L35–45；`handler/account_test.go` `TestAccountsMeDevSessionFallback` L105–109。
  - 影响：recommended；种子 grants 与真实登录 401/403 路径已另有覆盖。

- **F-004 · 承接 A-003 recommended + 测试计数口径**（severity: low；建议: **recommended**；status: open）
  - 描述：独立复核同意 A-003 **F-001**（每待应用版本快照顺延 S3）、**F-002**（users/roles 重启/排序/搜索正例顺延 S4 / I-011-003）、**F-003**（password 最小长度仍为 recommended）。另：02-execution 写「151 测试函数」；本机按 `^func Test` 计数约 **126**（含子包合计口径可能不同）。**全绿事实成立**，计数仅为口径差异，不构成回归失败。
  - 证据：A-003 Findings；`02-execution` S2 证据句；本机 `go test ./...` PASS。
  - 影响：非阻断；S3/S4 落实 A-003 项；执行记录可选用稳定口径（包级 PASS / 或明确计数方法）。

### 必改项汇总

| ID | 严重度 | 摘要 | 建议闭合路径 |
|----|--------|------|--------------|
| （无） | — | 本 scope **无 required / high** | — |

无 high/medium required。F-001～F-004 均为 recommended，不单独阻断 S2 放行或 S3 开工。

### 与既有意见的异同

| 项 | A-003（self · pass） | A-004（independent · pass） |
|----|----------------------|------------------------------|
| 工厂扩展 actor / JSONFields / DomainError 顺序 | 通过 | **同意**（代码级复核） |
| 领域不变量真实实现 | 通过 | **同意**；不隐式建角色有负向 count 断言 |
| records 零对外变更 | 通过 | **同意**（既有 records 套件全绿） |
| 401/403 + 种子 grants | 通过 | **同意** |
| 0005 行保留证据 | 将 `TestMigrateExistingV3ToV4` 作行保留证据 | **收窄**：实现正确，但该测不足以证明行保留 → **F-001 recommended** |
| LAST_ADMIN | store 层覆盖即可 | **同意主体**；补 HTTP 断言为 F-002 recommended |
| StaticDevSession | 实现同步 | 实现同意；回归断言缺口 → F-003 |
| A-003 F-001～F-003 recommended | open | **同意维持**（并入 F-004） |
| verdict | pass | **同向 pass**（无 verdict 冲突） |

同 scope 下 self 与 independent **verdict 一致（pass）**，无 P-004 §3.2 冲突；差异仅为 recommended 粒度与证据指针精度。

### 结论 + 建议给编排器/用户的下一步

**verdict: pass**——S2 后端 users/roles 闭环在 I-011-001 v0.2.0 下**真实落地**：工厂扩展（含 A-002 F-001 actor 通道与 DomainError 优先级）、领域不变量、错误码 envelope、401/403 与种子投影、0005 迁移与 checksum 账本追加、以及本机 API/Web 回归均**可核对成立**。无未闭合 required finding；无到期阻断 S2 的 required 信息项（`I-011-003` 最晚 S4）。

**建议下一步（/govern）**：

1. 汇总 A-003 + A-004（同向 pass）；可选将 F-001～F-004 标为 open recommended 随 S3/S4 落实，或即时加固 0005 行保留测试后 closed。
2. **放行 S3** records 产品运行面退场：migration `0006` DROP TABLE + 权限/菜单清理 + **每待应用数据变更迁移前快照**（A-002 F-002 / A-003 F-001 / I-011-002 v0.2.0）+ API/种子/fixture/前端/测试退场。
3. S3 实施时优先消化 A-004 F-001（0005 行保留测试可随迁移测试包一并补）与 A-003 F-001（per-pending 快照）。
4. 不修改本意见中的 status/progress；阶段推进与 finding 响应归编排器。

### 声明

本意见 **source: independent**，仅追加审计台账；**不修改** `00-meta` 的 status / 检查点 / 派生 progress，**不修改** goal-tree 状态列，**不修改**契约正文或产品代码。响应、finding 闭合与阶段推进归 **`/govern`**。

## 响应 A-003 + A-004（self · 编排响应 · 2026-08-03）

S2 同向 pass（A-003 self / A-004 independent），无 verdict 冲突、无 required；按 recommended 加固即时落实如下。

### 关闭证据表

| Finding | 严重度 | 状态 | 证据路径 |
|---------|--------|------|----------|
| A-004 F-001 · 0005 行保留缺专用回归 | low · recommended | **fixed** | 新增 `store/operations_test.go` `TestMigrate0005PreservesOperationLogRows`（0004 态写 records./auth. 行 → Open 升级 0005 → 断言行内容/record_id 完整保留 + 新 `users.*` 事件可写入） |
| A-004 F-002 · LAST_ADMIN 缺 HTTP 层断言 | low · recommended | **fixed** | 新增 `handler/users_test.go` `TestUsersLastAdminHTTP`（dev 会话非 admin actor 删除唯一 admin → 409 `LAST_ADMIN`） |
| A-004 F-003 · StaticDevSession 投影缺回归断言 | low · recommended | **fixed** | `handler/account_test.go` `TestAccountsMeDevSessionFallback` 增补 users/roles 四权限键 + `menu_users`/`menu_roles` 断言 |
| A-004 F-004 · 承接 A-003 recommended + 计数口径 | low · recommended | **handled** | A-003 F-001（per-pending 快照）随 S3 落地；A-003 F-002（重启/排序/搜索正例）随 S4/I-011-003；A-003 F-003（password 长度）维持 open recommended；02-execution 计数口径已修正 |
| A-003 F-001/F-002/F-003 | low · recommended | open | 随 S3（快照）、S4（重启/搜索/排序）、非目标（password 长度）落实，非门禁阻断 |

### 结论

S2 无开放 required；A-003/A-004 趋同为 **pass**，可放行 S3。下一步：S3 records 产品运行面退场（0006 DROP TABLE + 每待应用迁移前快照 + 权限/菜单/fixture/前端/测试退场），优先落实 A-004 F-001 的 0005 行保留回归与 A-003 F-001 的 per-pending 快照。
