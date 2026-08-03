---
title: 审计台账 · 语义化 Admin 资源替换与双实体验证
status: active
created: 2026-08-03
updated: 2026-08-04
parent: GOAL-010-a002-schema-adapter
version: 0.11.0
---

# 审计台账 · GOAL-011

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | self | 2026-08-03 | S1 契约冻结（D-002 · I-011-001/002） | pass | 无 required；3 条 recommended |
| A-002 | independent | 2026-08-03 | S1 契约冻结（D-002 · I-011-001/002） | conditional | 2 required + 4 recommended，均已响应（D-003，见响应节） |
| A-003 | self | 2026-08-03 | S2 后端 users/roles 闭环（I-011-001 v0.2.0） | pass | 无 required；3 条 recommended（随 S3/S4 落实） |
| A-004 | independent | 2026-08-03 | S2 后端 users/roles 闭环（I-011-001 v0.2.0；progress 1/5 → 2/5） | pass | 无 required；F-001～F-003 fixed、F-004 handled（见响应节） |
| A-005 | self | 2026-08-03 | S3 records 产品运行面退场（I-011-002 v0.2.0） | pass | 无 required；3 条 recommended（随 S4/S5 落实） |
| A-006 | independent | 2026-08-03 | S3 records 产品运行面退场（I-011-002 v0.2.0；progress 2/5 → 3/5） | pass | 无 required；F-001/F-002 fixed、F-003/F-004 handled（见响应节） |
| A-007 | independent | 2026-08-03 | I-011-003 冻结就绪性（双资源验收矩阵） | conditional | 3 条 medium required，均已 fixed（D-004，见响应节） |
| A-008 | independent | 2026-08-03 | A-007 F-001～F-003 finding-closure 复审（I-011-003 冻结） | pass | 三项 required 均可重核为 fixed；不等于 S4 完成 |
| A-009 | self | 2026-08-03 | S4 双语义实体 Schema 接入验证（I-011-003 v0.2.0；S4 验收收据） | pass | 无 required；2 条 recommended |
| A-010 | self | 2026-08-03 | S5 关门审计（GOAL-011 全目标 close-out） | pass | 无 required；F-001（grok 服务取消）为可复核待办 |
| A-011 | independent | 2026-08-04 | S5 关门独立交叉审计（grok build 文本意见代贴） | conditional | 2 required（F-001/F-002）已 fixed；3 recommended handled（见响应节） |
| A-012 | independent | 2026-08-04 | S5 关门独立代码审计（VP-002 users/roles 生产可用性 + records 残留） | fail | 5 required open；1 recommended open；D-005 已响应但未闭合 |

## 当前审计边界

- S1 契约已冻结（D-002；`I-011-001`/`I-011-002` → verified，v0.2.0）；S2 后端闭环已实施（A-003 self · pass；A-004 independent · pass，recommended 已响应）；S3 records 退场已实施（A-005 self · pass；**A-006 independent · pass**）；I-011-003 v0.2.0 已由 D-004 冻结并 verified；**S4 双资源 Schema 接入验证已完成（A-009 self · pass，验收收据落盘）**；**S5 回归、关门审计与父级交接已执行（A-010 self · pass；A-011 independent · conditional 经 F-001/F-002 修复后趋同）**。
- A-002（independent）对 S1 给出 **conditional**，两条 required（F-001 actor 通道、F-002 快照语义）经 **D-003** 走 `fixed` 闭合，F-003～F-006 采纳为 `handled`；A-001（pass）与 A-002 已趋同。
- A-003 / A-004 对 S2 同向 **pass**；A-005 / A-006 对 S3 同向 **pass**；本 scope 无开放 required；recommended 见 A-003/A-004 响应节、A-005 F-001～F-003 与 A-006 F-001～F-004（随 S4/S5 或文档清理落实）。
- A-007 的 F-001～F-003 经用户裁决全部 `fixed`（见响应节），`I-011-003` 信息门禁已解除；该响应不是同范围自审，也不构成 S4 阶段通过或 progress 推进。
- A-009 对 S4 给出 **pass**（验收收据 + Renderer 零 diff + 全维度证据）；S4 无开放 required，可进入 S5。三个 required 信息项（I-011-001/002/003）均已 verified。
- A-012 在既有文档/契约闭合之外进行代码与流水线复核，发现 5 条新的 required：用户写权限可越权分配 admin、改密不撤销 refresh 且密码输入不安全、users/roles Schema 不能完成实际角色语义、活动 CI 仍依赖退役 records、删除用户的最后管理员检查存在事务竞态。D-005 已完成第一轮 `/govern` 响应：五项仍 open、未选择 residual/overruled，GOAL-011 已 fail closed 恢复为 `active / 5/5`；A-010/A-011 与 A-012 的 close-out 冲突仍待用户裁决。
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

## A-005 · S3 records 产品运行面退场自审（2026-08-03）

- **source**：self
- **auditor**：Claude Code（govern orchestrator）
- **类型 / scope**：stage · S3 records 产品运行面退场（I-011-002 v0.2.0；S3 检查点达成）
- **verdict**：pass

### 范围与区间

审 GOAL-011 S3「移除 records API 注册、Schema fixture、菜单/权限、种子、操作日志耦合、前端 records 专名与当前测试依赖；保留不可改写历史治理事实与迁移链证据」。对照 I-011-002 v0.2.0 验收口径；**不含** S4 双资源 Schema 接入验证（I-011-003 open）与 S5 关门。工作区绑定与共享资料（无）已核对。

### 成果（有证据）

- **迁移 0006**：`DROP TABLE records` + 清理 records 权限/菜单行（先 join 后父行，FK 安全）+ checksum 冻结；**per-pending 快照**落地（`migrate()` 每待应用数据变更迁移前快照，0005+0006 同批时 `pre-v0006` 存在——`TestMigrateExistingV3ToV4` 断言 pre-v0005/pre-v0006）。
- **后端退场**：`/api/records` 路由移除；`handler/records.go`、`store/records.go`、`seed_records.go` 删除；`seed.go` 去 records 权限/菜单/grants；`StaticDevSession` 去 records 键；共享辅助（`jsonQuote`/`newOperationID`/`ErrRecordExists`）迁至通用位置。
- **fixture/manifest**：`users.json`/`roles.json` CRUD 页新增（I-011-002 §3.3）；records/catalog 页与 `menu_list_edit_lifecycle` 移除；`data-table`/`search-form-table` 改指 `/api/users`/`/api/roles`；users/roles 页接入 manifest（`menu_users`/`menu_roles`）。
- **前端专名退场**：`RecordItem`/`RecordList` 别名删除、`RecordsQuery`→`ResourceQuery`；`use-records.ts` 删除；`recordsFetcher`→`resourceFetcher`；`schema-table` footer/empty/caption 与 `render` 成功 toast 去 records 化。
- **测试/文档退场改指**：records 测试删除或改指 users/roles；进程级重启测试改指 users；QUICKSTART/smoke.sh/playwright 更新。
- **验收口径**：fresh install 无 `records` 表（`TestMigrateFreshDB`）；产品代码 grep 无 `api/records`/`records.read`/`records.write`/`menu_list_edit_lifecycle`/`list-edit-lifecycle` 残留（仅 0006 清理语句与历史注释）；users/roles 操作日志生效。
- **回归**：`go test ./...` 全绿 + `go vet` 干净；web `vitest` 481/481 + `tsc -b` + `vite build` 干净；e2e `playwright` 2/2（users CRUD 真实 Go/SQLite 往返 + shell/auth 链）。
- **未越权**：历史 GOAL-004/007/010 文档、I-007-001 契约、0001～0004 迁移 checksum 均未改写；records 退场走新迁移 + 代码演进（I-011-002 硬约束满足）。

### 对照成功标准

| S3 标准 | 状态 | 证据 |
|---------|------|------|
| 移除 records API 注册 / 种子 / 权限菜单 / 操作日志耦合 | ✅ | health.go/seed.go/operations.go 无 records 写路径；0006 清理 |
| 移除 records Schema fixture / 前端专名 / 测试依赖 | ✅ | fixture 删 records 页、前端去 records 化、测试改指 |
| 保留历史治理事实与迁移链证据 | ✅ | 0001～0004 checksum 未变；历史文档保留 |
| 既有库升级（0004→0006）可追溯 + 数据处置快照 | ✅ | TestMigrateExistingV3ToV4 pre-v0006 断言 + recordsRetire 迁移 |
| 验收口径（fresh 无 records 表 / grep 无残留 / 操作日志 users/roles） | ✅ | TestMigrateFreshDB + grep + 操作日志测试 |

### Findings

- **F-001 · S4 验收矩阵（I-011-003）待 S4 冻结**（severity: low；建议: recommended；status: open；关联 I-011-003）
  - 描述：S3 已移除 records 并接入 users/roles 页；但「双资源 Schema-only 接入 + fresh/upgrade/restart/401-403 完整边界」的正式验收矩阵（I-011-003）与 Renderer diff 证据属 S4。S3 完成定义已由 §5 口径覆盖。
  - 影响：不阻断 S3；S4 首步冻结 I-011-003 并产出双资源证据。

- **F-002 · 进程级重启持久化已改指 users 但 I-011-003 矩阵未列**（severity: low；建议: recommended；status: open）
  - 描述：`cmd/server` 重启测试已改指 users（create/patch/delete + 重启），但 formal 双资源（users+roles）重启矩阵列于 S4。
  - 影响：不阻断；S4 补 roles 重启路径。

- **F-003 · smoke.sh SM-006 未在本会话执行（需 --disposable compose）**（severity: low；建议: recommended；status: open）
  - 描述：smoke.sh 已改指 users/roles（`/users`、`/api/users`、`user-admin` 种子断言），但 SM-006 种子可重复性需 `--disposable` Compose 隔离环境，本会话未运行；`bash -n` 语法校验通过。
  - 影响：不阻断 S3；S5 全量回归可跑 smoke（含 disposable 环境）。

### 必改项汇总（required 列表）

无。

### 结论 + 建议下一步

S3 records 退场完整、验收口径达成、回归全绿、未改写历史事实；无未闭合 required。**pass**。下一步：S4 双语义实体 Schema 接入验证（冻结 I-011-003 验收矩阵 + fresh/upgrade/restart/401-403 双资源证据 + Renderer diff 边界）。按用户指令，本自审后将调用 **grok build 独立交叉审计**（scope: S3 records 退场），等待其意见后合并响应。

## 响应 A-005 + A-006（self · 编排响应 · 2026-08-03）

S3 同向 pass（A-005 self / A-006 independent），无 verdict 冲突、无 required；recommended 按实施成本即时落实。

### 关闭证据表

| Finding | 严重度 | 状态 | 证据路径 |
|---------|--------|------|----------|
| A-006 F-001 · apps/api/README + apps/web/README 仍写 records 为现行 API | low · recommended | **fixed** | `apps/api/README.md`（端点表/鉴权边界/测试覆盖改指 users/roles + records 退场注记）、`apps/web/README.md`（页面/鉴权改指 users/roles） |
| A-006 F-002 · web 单测仍用 `/api/records` 路径字符串 | low · recommended | **fixed** | `records.test.ts`/`schema-table.test.tsx`/`render.test.tsx` 示例 URL 改指 `/api/users`（`rec-3`→`usr-3`）；产品代码 grep 现已无 `api/records`（测试文件亦无残留） |
| A-006 F-003 · 承接 A-005（I-011-003 / roles 重启 / smoke disposable） | low · recommended | **handled** | A-005 F-001（I-011-003）随 S4 冻结；F-002（roles 重启）随 S4 补；F-003（smoke SM-006）S5 全量回归跑 disposable compose |
| A-006 F-004 · grok 会话 Playwright EACCES :5173 未重跑 | low · recommended | **handled** | 本会话已以 `WEB_PORT=9999` 重跑 `playwright test` **2/2 通过**（users CRUD 真实往返 + shell/auth 链） |
| A-005 F-001/F-002/F-003 | low · recommended | open | 随 S4（I-011-003 矩阵 + roles 重启）与 S5（smoke disposable）落实 |

### 结论

S3 无开放 required；A-005/A-006 趋同为 **pass**，可放行 S4。下一步：S4 双语义实体 Schema 接入验证——冻结 **I-011-003** 验收矩阵（fresh fork / 既有库升级 0004→0006 / 重启持久化 / 401-403 双资源 + Renderer diff 边界），产出 users/roles 双资源产品证据并关闭 I-011-003。

## A-006 · S3 records 产品运行面退场独立交叉审计（2026-08-03）

- **source**：independent
- **auditor**：grok build
- **类型 / scope**：execution-facts · S3 records 产品运行面退场（GOAL-011 S3；I-011-002 v0.2.0；S3 检查点达成；progress 2/5 → 3/5）
- **verdict**：pass

### 范围与区间

审 GOAL-011 S3「按冻结策略移除当前产品默认运行面中的 records API 注册、Schema fixture、菜单/权限、种子、操作日志耦合、前端 records 专名与当前测试依赖；保留不可改写的历史治理事实与迁移链证据」是否与 I-011-002 v0.2.0 **真实落地**（非仅声明）。关注点：

1. migration `0006 records_retire`（DROP TABLE + 权限/菜单清理）与 **per-pending 快照**（0005+0006 同批时 `pre-v0006` 必存在）是否真实落地，且不改写 0001～0004 checksum/历史文档；
2. records API/种子/权限/菜单/fixture/前端专名/测试是否从产品运行面完整退场；
3. users/roles fixture 与 manifest 接入是否正确（I-011-002 §3.3）；
4. S3 验收口径（fresh 无 records 表、grep 无 `api/records` 等残留、操作日志 users/roles 生效）；
5. 回归证据充分性（go test / vet / web vitest / tsc / vite build / e2e playwright）。

**不含** S4 双资源 Schema 接入正式验收矩阵 / `I-011-003`（open，最晚 S4，已到期待冻结）与 S5 关门。

工作区：`workspace-002-production-admin-foundation` / Root `GOAL-001-production-admin-foundation` / `canonical_scope` 已校验；`shared_materials_catalog: none`，未将共享资料当作关闭证据。未读取或比较其他工作区。

**只读核验**：GOAL-011 五件套与附件（I-011-001/002 v0.2.0）；`migrate.go`（0006 + per-pending 快照循环）；`seed.go` / `store.go`；`handler/{health,resources,users,roles}.go`；`account/session.go`；fixtures `users.json`/`roles.json`/`data-table.json`/`search-form-table.json`；`app-manifest.json`；`renderer/{records.ts,schema-table.tsx,render.tsx}`、`App.tsx`；相关 store/handler/cmd 测试；`QUICKSTART.md`、`scripts/smoke.sh`。本机重跑：`go test ./...`（apps/api）全绿 + `go vet ./...` 干净；web `vitest run` **481/481** + `tsc -b` + `vite build` 干净。e2e Playwright 本机启动失败（`listen EACCES 127.0.0.1:5173`，环境端口权限，非产品断言失败）——见 F-004。

### 成果（有证据）

| 主张 | 证据 | 核验结论 |
|------|------|----------|
| **① 0006 `records_retire`** | `compiledMigrations` 追加 version=6、name=`records_retire`、transformID=`0006:records-retire:v1`；`recordsRetireDDL`：`DROP TABLE IF EXISTS records` + 先删 join（`role_permissions`/`role_menu_items`）再删父行（`permissions`/`menu_items`，FK 安全、幂等） | **成立** |
| **① per-pending 快照**（A-002 F-002 / I-011-002 §2.3 v0.2.0） | `migrate()` 对每个 pending 且 version≥2 在 `applyMigration` **前**调用 `snapshotBeforePending`；注释明确 0005+0006 同批 → `pre-v0005`+`pre-v0006` | **成立**（相对 first-pending-only 已改正） |
| **① pre-v0006 回归** | `TestMigrateExistingV3ToV4`：v3 基线 Open 升级后断言 `pre-v0004`/`pre-v0005`/`pre-v0006` 各 1 个 + applied 含 0006 `records_retire` + records 表不存在 | **成立** |
| **① 0001～0004 账本未改写** | 0001～0004 的 `transformID`/`stmts` 仍为历史职责（0003 仍 CREATE records、0004 仍含 records.* CHECK）；退场仅靠 **0006 新迁移**；`migrationChecksum` 机制未改 | **源码级成立**（未改历史迁移字面） |
| **② API 注册退场** | `health.go` 仅 `registerResource(users/roles)`；`handler/records.go` / `store/records.go` / `store/seed_records.go` **文件不存在** | **成立** |
| **② 种子 / StaticDev** | `seed.go` 仅 users/roles 四权限 + 两菜单 + grants（admin 4rw+menus、editor/viewer ro）；无 records 权限/菜单；`StaticDevSession` 同步 `users.*`/`roles.*` + `menu_users`/`menu_roles` | **成立** |
| **② fixture / manifest** | 删除 `list-edit-lifecycle.json`/`catalog.json`；新增 `users.json`/`roles.json`（dataSource `/api/users`/`/api/roles`，permissions 键 users/roles.write）；`data-table`→users、`search-form-table`→roles 且文案已去 `api/records`；manifest 含 users/roles 页 + sidebar `visibleWhen` `menu_users`/`menu_roles`，无 list-edit-lifecycle/catalog | **成立**（§3.3） |
| **② 前端专名** | `RecordItem`/`RecordList` 别名已删 → `ResourceItem`/`ResourceList`；`RecordsQuery`→`ResourceQuery`；`use-records.ts` 不存在；`App.tsx` prop `resourceFetcher`；toast `Item *`；schema-table empty/caption 去 records 化 | **成立**（`fetchRecords`/`records.ts` 文件名按契约保留为通用 transport） |
| **② 测试改指** | `cmd/server` 进程级重启改 users；`seed_test` 权限集仅 users/roles；records handler/store 测试文件已删；web schema-crud/e2e 驱动 users 页 | **运行面成立** |
| **③ 操作日志 users/roles** | `EventUser*`/`EventRole*` 常量 + `usersOnWrite`/`rolesOnWrite`；`TestMigrateFreshDB` 写 `users.create`；handler users/roles 操作日志测试 | **成立**（0005 生效；records.* 保留历史合法值、无运行时写路径） |
| **④ fresh 无 records 表** | `TestMigrateFreshDB` 断言 applied=[1..6] 且 `records` 表不存在 | **成立** |
| **④ grep 运行面残留** | `apps/api` `*.go`：仅 `migrate.go` 0006 清理 SQL 含 `menu-list-edit-lifecycle` + 历史注释；无运行时 `api/records` 路由。fixtures/manifest/seed 无 `api/records`/`records.read`/`records.write`/`list-edit-lifecycle` | **运行面成立**；测试字符串/README 见 F-001/F-002 |
| **⑤ 本机回归** | `go test ./...` 全绿 + `go vet` 干净；`vitest` 481/481 + `tsc -b` + `vite build` 干净 | **成立**（与 02-execution 同向） |
| **⑤ e2e** | 规格已改指 users CRUD 真实往返；本机 Playwright webServer `EACCES :5173` 未完成重跑 | **代码改指成立**；独立重跑收据见 F-004 |
| **文档 S3 对象** | `QUICKSTART.md` 已写 records 0006 退场 + 终点 `/users`；`smoke.sh` SM-005→`/users`、SM-006→`/api/users` | **成立** |
| **未越权** | 本意见不修改 status/progress/契约/产品代码；历史 GOAL 文档与 0001～0004 未改写 | **成立** |

### 对照成功标准

| S3 标准（00-meta / I-011-002 §5） | 本意见 | 说明 |
|----------------------------------|--------|------|
| 移除 records API 注册 / 种子 / 权限菜单 / 操作日志耦合 | ✅ | health/seed/StaticDev 干净；无 records 写 OnWrite |
| 移除 records Schema fixture / 前端专名 / 测试依赖 | ✅ | fixture 删 records 页；专名泛化；测试改指 |
| 保留历史治理事实与迁移链证据 | ✅ | 0001～0004 transformID/stmts 保留；0006 新迁移 |
| 既有库升级 + per-pending 快照（含 pre-v0006） | ✅ | 实现 + `TestMigrateExistingV3ToV4` |
| fresh 无 records 表 | ✅ | `TestMigrateFreshDB` |
| grep 产品运行面无关键残留 | ✅（有条件） | 运行面/fixture 干净；单测 URL 字符串与 README 见 recommended |
| 操作日志 users/roles | ✅ | 常量 + OnWrite + 测试 |
| 回归 go/web/build | ✅ | 本机重跑通过 |
| e2e playwright | 有条件 | 规格正确；本独立会话未能重跑通过（环境） |

### Findings

- **F-001 · 产品面 README 仍将 records 描述为现行 API**（severity: low；建议: **recommended**；status: open）
  - 描述：`apps/api/README.md` 端点表与权限说明仍列出 `/api/records` + `records.read`/`records.write` 与 `seedRecords`；`apps/web/README.md` 仍写 `data-table` 走 `/api/records`、`list-edit-lifecycle` 与 records 权限键。契约 §3.7 仅强制更新 `QUICKSTART.md`（已更新），故**不升 required**；但对 fork 用户构成「产品运行面已退场、文档仍宣称在线」的双真相。
  - 证据：`apps/api/README.md` L3/L93–105；`apps/web/README.md` L53–68；对照 `health.go` 已无 records 注册。
  - 影响：不阻断 S3 代码验收；建议 S4/S5 或即时将 README 改指 users/roles，并标注 records 为历史/已退场。
  - 建议闭合：同步两份 README 端点/示例页表 → fixed。

- **F-002 · 前端单测仍大量使用 `/api/records` 作为路径字符串；A-005「grep 无残留」表述过宽**（severity: low；建议: **recommended**；status: open）
  - 描述：运行面/fixture/manifest 已无 `api/records`。但 `apps/web/src` 内 `records.test.ts`、`schema-table.test.tsx`、`render.test.tsx`、`stage3-fixtures.test.ts` 等仍以 `/api/records` 作通用 transport/dataSource 用例。`records.ts` 注释仍含 `(records.write)`。这些**不构成运行时耦合**（无真实路由、无权限键），故不升 required；A-005「仅 0006 清理语句与历史注释」未覆盖上述单测字符串，证据表述宜收窄。
  - 证据：`apps/web/src/renderer/records.test.ts` 等多处 `"/api/records"`；对照 I-011-002 §5 grep 句。
  - 影响：recommended 精度项；S4 可顺手改为 `/api/users` 或中性 `/api/example`。
  - 建议闭合：单测改指 + 收窄执行/自审 grep 口径 → fixed。

- **F-003 · 承接 A-005 recommended（I-011-003 / roles 重启 / smoke disposable）**（severity: low；建议: **recommended**；status: open；关联 I-011-003）
  - 描述：独立复核同意 A-005 **F-001**（S4 验收矩阵 I-011-003 待冻结）、**F-002**（进程级重启已改指 users，formal 双资源含 roles 属 S4）、**F-003**（`smoke.sh` 已改指 users，SM-006 disposable 本会话未跑）。不阻断 S3。
  - 证据：A-005 Findings；`cmd/server/server_restart_test.go` 标题/路径 users；`scripts/smoke.sh` SM-005/006。
  - 影响：随 S4/S5 落实。

- **F-004 · 本独立会话未能重跑 Playwright e2e（环境 EACCES）**（severity: low；建议: **recommended**；status: open）
  - 描述：本机 `npx playwright test` 因 webServer `listen EACCES: permission denied 127.0.0.1:5173` 未能启动，**不是**断言失败。规格 `e2e/schema-crud.spec.ts` 已驱动 `/users` 真实 create/edit/delete。A-005 主张 2/2 作为执行方历史收据保留；本意见**不**将其当作本会话可重复关闭证据。
  - 证据：本会话 Playwright 启动错误；`apps/web/e2e/schema-crud.spec.ts` L19–70。
  - 影响：不否定规格改指；S5 全量回归或修复端口后重跑可闭合本 finding。
  - 建议闭合：可复现的 e2e 2/2 收据 → fixed。

### 必改项汇总

| ID | 严重度 | 摘要 | 建议闭合路径 |
|----|--------|------|--------------|
| （无） | — | 本 scope **无 required / high** | — |

无 high/medium required。F-001～F-004 均为 recommended，不单独阻断 S3 放行或 S4 信息收集开工（`I-011-003` 仍 open 且最晚 S4，属独立信息门禁）。

### 与既有意见的异同

| 项 | A-005（self · pass） | A-006（independent · pass） |
|----|----------------------|------------------------------|
| 0006 + per-pending 快照 + pre-v0006 测试 | 通过 | **同意**（代码 + 测试复核） |
| 0001～0004 未改写 | 通过 | **同意**（transformID/stmts 源码核对） |
| API/种子/fixture/manifest 退场与 users/roles 接入 | 通过 | **同意** |
| 前端专名 / resourceFetcher | 通过 | **同意**（transport 文件名按契约保留） |
| fresh 无 records + 操作日志 users/roles | 通过 | **同意** |
| go/web 回归 | 通过 | **同意**（本机重跑） |
| grep 无残留 | 「仅 0006 + 历史注释」 | **收窄**：运行面成立；单测 URL 与 README 另列 F-001/F-002 |
| e2e 2/2 | 主张通过 | 规格同意；**本会话未重跑** → F-004 |
| A-005 F-001～F-003 recommended | open | **同意维持**（并入 F-003） |
| verdict | pass | **同向 pass**（无 verdict 冲突） |

同 scope 下 self 与 independent **verdict 一致（pass）**，无 P-004 §3.2 冲突；差异为 residual 文档/测试字符串精度与 e2e 独立重跑收据。

### 结论 + 建议给编排器/用户的下一步

**verdict: pass**——S3 records 产品运行面退场在 I-011-002 v0.2.0 下**真实落地**：0006 DROP + 权限/菜单清理、per-pending 快照（含 0005+0006 同批 `pre-v0006`）、0001～0004 账本保护、API/种子/fixture/manifest/前端专名/主测试路径退场、users/roles 页接入、操作日志与本机 API/Web 回归均可核对。无未闭合 required finding；无到期阻断 S3 的 required 信息项（`I-011-003` 最晚 S4，阻断的是 S4 验收而非 S3 完成定义）。

**建议下一步（/govern）**：

1. 汇总 A-005 + A-006（同向 pass）；可选将 F-001/F-002 即时修 README + 单测路径字符串，或标 open recommended 随 S4/S5。
2. **放行 S4 准备**：先冻结 `I-011-003`（双资源 Schema-only + fresh/upgrade/restart/401-403 矩阵 + Renderer diff），再实施/验收；勿以 S3 progress 数字代替 I-011-003 门禁。
3. S5 前补 e2e 可复现收据（F-004）与可选 smoke disposable（A-005 F-003）。
4. 不修改本意见中的 status/progress；阶段推进与 finding 响应归编排器。

### 声明

本意见 **source: independent**，仅追加审计台账；**不修改** `00-meta` 的 status / 检查点 / 派生 progress，**不修改** goal-tree 状态列，**不修改**契约正文或产品代码。响应、finding 闭合与阶段推进归 **`/govern`**。
## A-007 · I-011-003 冻结就绪性独立交叉审计（2026-08-03）

- **source**：independent
- **auditor**：Codex（GPT-5）
- **类型 / scope**：design-plan · `I-011-003` 双语义资源 Schema 接入验收矩阵是否足以冻结，并放行 S4 集成验收与 S5 关门所需的信息门禁
- **verdict**：conditional

### 范围与区间

审当前工作区 `workspace-002-production-admin-foundation` 的 GOAL-011：`workspace.md`、`goal-tree.md`、目标五件套、[I-011-003 候选矩阵](attachments/I-011-003-acceptance-matrix.md) v0.1.0，以及其列出的 API/Web 测试和当前 Git 证据。`shared_materials_catalog: none`，未把共享资料作为证据，未读取或比较其他工作区。

本审计判断的是**信息冻结就绪性**，不是 S4 已实施、S4 检查点已勾选、GOAL-010 S4 已交接，亦不是 Root A-002 F-002-001 的关闭复核。

### 成果（有证据）

- S1/S2/S3 已有同向审计结论：A-003/A-004 对 S2 为 pass，A-005/A-006 对 S3 为 pass；本目标没有承接到 S4 的开放 required finding。
- 候选矩阵正确覆盖了应回答的主要域：fresh DB、0004 态升级、进程级重启、401/403、操作日志，以及 Renderer 主路径边界。
- 本审计重跑候选矩阵所列的 API 目标测试：`go test ./internal/store ./internal/handler ./cmd/server -run 'TestMigrateFreshDB|TestMigrateExistingV3ToV4|TestMigrate0005PreservesOperationLogRows|TestUsersListAndDetail|TestUsersCreateUpdateDeleteLifecycle|TestUsersAuthGates|TestUsersOperationLogEvents|TestRolesListAndDetail|TestRolesWriteLifecycleAndProtection|TestRolesAuthGates|TestRolesOperationLogEvents|TestServerProcessRestartPersistsUsers' -count=1`，三个包均通过。
- 本审计重跑 Web 相关证据：`npm test -- src/renderer/schema-crud.test.tsx src/renderer/representative-pages.test.tsx src/app/representative-pages.integration.test.tsx src/app/navigation.test.ts src/protocol/app-manifest.test.ts`，5 个文件、47 项均通过。
- 复核时当前提交为 `adfe15a17da770699d5e109f22402c41ece5eeea`；仅候选矩阵为未跟踪文件，受限路径没有已跟踪产品代码 diff。此事实可作为后续冻结时选择基线的输入，但候选矩阵尚未固定它。

### Findings

- **F-001 · 候选契约将未发生的冻结写成既成事实**（severity: medium；建议: **required**；status: open；关联 `I-011-003`）
  - 描述：[I-011-003 候选矩阵](attachments/I-011-003-acceptance-matrix.md) v0.1.0 已填写 `related_decision: D-004`、标题“冻结”、正文“由 D-004 置为 verified”及修订记录“关闭 I-011-003”。但当前 `00-meta.md` 仍将该项列为 **open**，`01-decision.md` 只存在 D-001～D-003，尚无 D-004。候选与 canonical 状态相反，违反 P-005 对未知/待裁决不得伪装为既成事实的要求。
  - 证据：候选矩阵 frontmatter/§7；[00-meta.md](00-meta.md) 信息表 `I-011-003`；[01-decision.md](01-decision.md) D-001～D-003。
  - 影响：不能把候选附件本身视为冻结或门禁解除依据。
  - 建议闭合：在 D-004 获 `/govern` 正式采纳前，将附件明确标为 candidate/draft，并把 `D-004`、`verified`、`关闭`改为条件性措辞；采纳后再由编排器一次性写 D-004、信息表状态/证据、执行记录与 goal-tree 投影。

- **F-002 · 双资源 Schema-only 证明没有可重复基线，且 roles 缺少页面级证据**（severity: medium；建议: **required**；status: open；关联 `I-011-003`）
  - 描述：候选 §3 要求 S4 的 Renderer/App diff 为空，并以 T-UI-10 和 representative-pages 作为证明；但未冻结可比较的 Git revision 或命令。现有 T-UI-10 只检查 `createUser`/`updateUser`/`deleteUser` 等 users action id；representative 页面测试把 roles 纳入结构加载循环，却只直接渲染 users CRUD 页面；浏览器 e2e 同样只走 users。因而不能证明 roles 也以 Schema-only 方式完成页面 CRUD，亦不能在 S4 后可重复地证明 Renderer 未改。
  - 证据：候选矩阵 §3；`apps/web/src/renderer/schema-crud.test.tsx` T-UI-10；`apps/web/src/renderer/representative-pages.test.tsx`；`apps/web/src/app/representative-pages.integration.test.tsx`；`apps/web/e2e/schema-crud.spec.ts`。
  - 影响：直接触及 I-011-003 的“双资源”和“Renderer 主路径无修改”核心主张，不能只用 users 证据替代。
  - 建议闭合：冻结时记录明确 baseline revision（本审计复核的 `adfe15a17da770699d5e109f22402c41ece5eeea` 可供用户采纳）及受限路径的可执行 diff 命令；在 S4 矩阵中增加 roles 的真实 manifest/fixture 页面渲染与 CRUD action 断言，并把 roles action id 的无硬编码检查纳入 T-UI-10 或新的具名测试。尚未存在的测试应写为 S4 必交付证据，不得标为“已实施”。

- **F-003 · 后端“完整双资源边界”行的文字超过现有断言**（severity: medium；建议: **required**；status: open；关联 `I-011-003`）
  - 描述：候选矩阵把现有测试表述为“双资源进程重启 list/detail + 毫秒往返”、“双资源五路由 401/403”及“双资源操作日志 actor + 非敏感 detail”。实际进程级测试对 roles 只做 create 后重启 detail，不做 roles list 或时间戳往返；两份 AuthGates 只覆盖匿名 list 与 viewer POST；roles 操作日志测试只断言事件序列，不断言 actor/detail；0004→0006 升级夹具只在迁移后的表上显式写入 `users.create`，没有 roles 事件的升级后断言。
  - 证据：候选矩阵 §2；`apps/api/cmd/server/server_restart_test.go`；`apps/api/internal/handler/{users_test.go,roles_test.go}`；`apps/api/internal/store/operations_test.go`。
  - 影响：已通过的测试支持 S2/S3 和部分 S4 输入，但不足以按候选的“完整边界”文字关闭此 required 信息项。
  - 建议闭合：二选一并在矩阵中固定：收窄验收口径为已证实的代表性读/写路径；或增加具名 S4 断言，覆盖 roles restart 的 list/detail 与毫秒时间戳、两资源五路由的 401/403 策略（可说明工厂共享门禁的推导）、两资源日志 actor/detail，以及升级库上的 `roles.*` 写入。后者更符合 GOAL-011 S4 的双实体成功标准。

### 必改项汇总

| ID | 严重度 | 摘要 | 影响门禁 |
|----|--------|------|----------|
| **F-001** | medium · required | 候选把 D-004/verified 写成既成事实 | `I-011-003` 冻结 |
| **F-002** | medium · required | roles 页面级 Schema-only 与 Renderer diff 基线不充分 | S4 集成验收 |
| **F-003** | medium · required | 后端验收矩阵超过其当前测试断言 | S4 集成验收 / S5 关门 |

### 与既有意见的异同

A-005/A-006 的 pass 仅确认 S3 records 退场；两者均将 `I-011-003` 留作 S4 的独立 required 信息门禁。A-007 不否定 S2/S3 事实或其 pass 结论，而是复核候选矩阵能否将该门禁从 open 合法转为 verified。

### 结论 + 建议给编排器/用户的下一步

**verdict: conditional**。`I-011-003` **尚不能冻结**：候选覆盖方向正确且列出的当前测试通过，但三项 required 缺口使其既不能诚实声称 D-004 已发生，也不能把单资源/部分路径证据升级为“双资源完整边界”。

先通过 `/govern` 修订候选以关闭 F-001～F-003；修订后的信息契约可作为 D-004 的候选输入。D-004、`I-011-003 → verified`、S4 放行及所有目标树投影仍由 `/govern` 和用户决定，本独立意见不作这些变更。

### 声明

本意见 **source: independent**，仅追加审计台账；**不修改** `00-meta` 的 status / 检查点 / 派生 progress，**不修改** goal-tree 状态列，**不修改**候选契约或产品代码。响应、finding 闭合与阶段推进归 **`/govern`**。

## 响应 A-007（self · 编排响应 · 2026-08-03 · GOAL-011 D-004）

- **响应性质**：这是 `/govern` 对既有独立意见的 finding closure 记录，不是新增的同范围 self audit，不产生新的 A 编号或 verdict。
- **用户裁决（P-004.1 / P-004 §3.2）**：用户明确选择“**不用补自审计，直接 fix**”。因此跳过同范围自审，F-001～F-003 全部按 `fixed` 闭合；没有 `accepted-residual` 或 `user-overruled`。

| finding | 闭合 | 可核对证据 |
|---------|------|------------|
| **A-007 F-001** | **fixed** | I-011-003 v0.1.0 在修订记录中明确为未冻结候选；D-004 实际落盘后，v0.2.0 才使用“冻结”与 `verified` 表述。`00-meta.md`、D-004、契约 frontmatter/正文状态一致。 |
| **A-007 F-002** | **fixed** | I-011-003 §3 固定 baseline `adfe15a17da770699d5e109f22402c41ece5eeea`、受限生产文件和可执行 diff 命令，当前 exit 0；`schema-crud.test.tsx` T-UI-10 增加 roles 真实 fixture 的 create/update/delete 与双资源 action-id 反证；`representative-pages.integration.test.tsx` 增加真实 manifest + roles fixture 页面断言。 |
| **A-007 F-003** | **fixed** | users/roles AuthGates 各覆盖五路由匿名 401 与 viewer 读 200/写 403；两资源 operation-log 测试核对 actor/record/detail；server restart 核对 roles list/detail 与毫秒时间戳往返；0005 升级测试核对重开后的 users/roles 新事件及 legacy 行。 |

### 验证收据

- `go test ./internal/handler ./internal/store`：通过。
- `go test ./...`：全包通过，含 `apps/api/cmd/server` 进程级重启测试。
- `npm test -- --run src/renderer/schema-crud.test.tsx src/app/representative-pages.integration.test.tsx`：2 files / 26 tests passed。
- `npm test`：23 files / 485 tests passed。
- `npm run build`：`tsc -b` + Vite production build 通过。
- I-011-003 §3 Renderer/App baseline command：exit 0、无 diff。

### 响应结论

A-007 三条 required finding 均已按合法 `fixed` 路径闭合，当前 scope 无开放 required；I-011-003 v0.2.0 + D-004 足以将信息项置为 `verified` 并解除 S4 信息门禁。GOAL-011 仍保持 `active / 3/5`，S4/S5 未勾选；本响应不把信息契约冻结升级为 S4 实施验收或父级/Root 接受。

## A-008 · A-007 F-001～F-003 finding-closure 独立复审（2026-08-03）

- **source**：independent
- **auditor**：Codex（GPT-5）
- **类型 / scope**：finding-closure · 仅复核 A-007 的 F-001～F-003 是否已按 `fixed` 合法闭合，以及这些修复是否足以支撑既有 D-004 对 `I-011-003` 的信息冻结结论
- **verdict**：pass

### 范围与区间

审当前工作区 `workspace-002-production-admin-foundation` 的 `workspace.md`、`goal-tree.md`、GOAL-011 五件套、I-011-003 v0.2.0、A-007 响应及其本次提交的 API/Web 测试。`shared_materials_catalog: none`，未使用共享资料，未读取或比较其他工作区。

本意见只复核 finding closure 和信息门禁证据，不审 S4 实施/验收收据、S4 检查点、S5 关门、GOAL-010 交接或 Root A-002 F-002-001。未重跑浏览器 E2E；它不替代且也不属于本次 A-007 三项缺口的闭合判据。

### 成果（有证据）

- **A-007 F-001**：D-004 已作为 `accepted` 决策实际落盘；I-011-003 的修订记录将 v0.1.0 明确为候选、v0.2.0 明确为 D-004 后冻结。`00-meta.md`、决策、执行记录、A-007 响应与附件均一致写为 `verified`；S4/S5 仍为未完成、`progress: 3/5`。
- **A-007 F-002**：I-011-003 §3 固定 baseline `adfe15a17da770699d5e109f22402c41ece5eeea`、受限生产文件和可执行命令；本审计在当前 `73bc93abba52db0440bc8c70eaf89969174a00cc` 重跑该命令，exit 0 且无 diff。T-UI-10 现直接覆盖 roles fixture 的 create/update/delete 与 users/roles action-id 无生产源码硬编码；真实 manifest + fixture 集成测试覆盖 roles 页面。
- **A-007 F-003**：users/roles AuthGates 均覆盖 list/detail/create/update/delete 的匿名 401 与 viewer 读 200/写 403；两资源 operation-log 均核对 actor、record、非敏感 detail；进程重启测试核对 roles list/detail 与时间戳往返；升级测试核对重开后的 users/roles 新事件及 legacy 行。
- **独立复现**：`go test ./... -count=1`、`go vet ./...`、`npm test`（23 files / 485 tests）、`npm run build`、I-011-003 §3 Renderer/App diff 命令以及 `git diff --check` 均通过。

### Findings

无新增 required finding。`goal-tree.md` 中 S1～S3 的“open/待冻结”表述是对应历史节点的事实；状态表仍正确显示 GOAL-011 `active / 3/5`。D-004 未改变 Goal 的 status、progress 或 parent，故不以该历史注记推导当前信息项仍为 open；S4 真正完成时仍须按 I-011-003 §5 同步目标树。

### 必改项汇总

| Finding | 复审结果 | 结论 |
|---------|----------|------|
| A-007 F-001 | fixed | D-004、v0.2.0 与 canonical 信息状态可相互核对。 |
| A-007 F-002 | fixed | 双资源页面级证据与 Renderer/App 可重复零 diff 基线均成立。 |
| A-007 F-003 | fixed | 后端双资源断言已补齐并经全量 API 测试复现。 |

### 与既有意见的异同

本复审确认 `/govern` 的 A-007 响应所列三条 `fixed` 闭合证据，不产生与 A-007 的 conditional 结论相冲突的新 required finding。D-004 已记录用户选择不补同范围自审、直接修复，故本次独立复审不触发新的 P-004 裁决点。

### 结论 + 建议给编排器/用户的下一步

**verdict: pass**。A-007 F-001～F-003 均有可重复核对的 `fixed` 证据；本意见确认它们足以支撑既有 D-004 将 I-011-003 v0.2.0 置为 `verified` 并解除 **S4 信息门禁**。

该 pass **不是** S4 阶段验收或 GOAL-011 关门结论。下一步由 `/govern` 按 I-011-003 §5 执行完整 S4 验收、形成带 revision/命令/结果的收据、汇总阶段审视；只有届时才可决定是否推进 `3/5 → 4/5` 并同步 `goal-tree.md`。

### 声明

本意见 **source: independent**，仅追加审计台账；**不修改** `00-meta` 的 status / 检查点 / 派生 progress，**不修改** goal-tree 状态列、D-004、信息契约或产品代码。finding 响应、阶段推进与任何状态变更归 **`/govern`**。

## A-009 · S4 双语义实体 Schema 接入验证自审（2026-08-03）

- **source**：self
- **auditor**：Claude Code（govern orchestrator）
- **类型 / scope**：stage · S4 双语义实体 Schema 接入验证（I-011-003 v0.2.0；S4 验收收据已落盘）
- **verdict**：pass

### 范围与区间

审 GOAL-011 S4「users/roles 代表性列表/CRUD 页面由 Schema 接入、Renderer 主路径无修改；fresh fork / 升级 / 重启 / 401-403 完整边界可复核」。对照 I-011-003 v0.2.0 §5 完成定义；**不含** S5 全量回归与关门审计。工作区绑定与共享资料（无）已核对。

### 成果（有证据）

- **验收收据**：02-execution 新增「S4 阶段验收收据」节（revision `73bc93a`、Renderer 基线 `adfe15a`、完整命令与结果）。
- **§2 全维度通过**：fresh fork / 升级 / 重启（users+roles）/ 401-403 双资源 / 操作日志 / Schema-only 页面均有可指回证据（I-011-003 §2 证据索引，本会话重跑通过）。
- **§3 Renderer diff 边界**：`git diff --exit-code adfe15a -- <受限生产文件>` **exit 0 无 diff**（本会话重跑）；T-UI-10 断言 users/roles action id 无 Renderer 硬编码；真实 manifest + fixture 页面测试覆盖 users 与 roles。
- **回归**：`go test ./...` 全绿 + `go vet` 干净；web `vitest` 23 files / **485/485** + `tsc -b` + `vite build` 干净；e2e `playwright` **2/2** 通过。
- **信息门禁**：`I-011-001/002/003` 均 `verified`（前序已关闭）；S4 无到期 required 信息项。

### 对照成功标准

| S4 标准 | 状态 | 证据 |
|---------|------|------|
| 关闭 I-011-003 | ✅ verified | I-011-003 v0.2.0 + D-004（前序） |
| users/roles 页由 Schema 接入 | ✅ | users.json/roles.json + manifest + T-UI-10 |
| Renderer 主路径无修改 | ✅ | 基线 diff exit 0 |
| fresh fork / 升级 / 重启 / 401-403 可复核 | ✅ | migrate/restart/users+roles 测试套件 |

### Findings

- **F-001 · S4 阶段审视未覆盖浏览器级角色管理 UX 断言**（severity: low；建议: recommended；status: open）
  - 描述：e2e 覆盖 users CRUD 真实往返；roles 的浏览器级 CRUD 未单独驱动（roles 由 handler/store 测试 + 进程重启测试覆盖）。S5 全量回归可补 roles 浏览器路径。
  - 影响：不阻断 S4；S5 可选加固。

- **F-002 · S5 交接前 GOAL-010 S4 证据链未落盘**（severity: low；建议: recommended；status: open）
  - 描述：GOAL-010 S4 父级验收门需要本目标证据链交接（I-011-003 §5）；属 S5。
  - 影响：不阻断 S4；S5 完成交接。

### 必改项汇总（required 列表）

无。

### 结论 + 建议下一步

S4 双资源 Schema 接入验证完整、Renderer 零 diff、全维度证据可复核、回归全绿；无未闭合 required。**pass**。下一步：S5 回归、审计与父级交接。按用户指令，本自审后将调用 **grok build 独立交叉审计**（scope: S4 双资源 Schema 接入验证），等待其意见后合并响应。

## A-010 · S5 关门审计（close-out · self · 2026-08-03）

- **source**：self
- **auditor**：Claude Code（govern orchestrator）
- **类型 / scope**：close-out · GOAL-011 全目标关门（S1～S5；users/roles 语义资源替换 + records 退场 + 双资源 Schema 接入验证 + 回归与交接）
- **verdict**：pass

### 范围与区间

关门审计覆盖 GOAL-011 全目标：S1 契约冻结、S2 后端闭环、S3 records 退场、S4 双资源 Schema 接入验证、S5 回归/交接。工作区绑定（workspace-002 / Root GOAL-001 / canonical）与共享资料（无）已核对；跨区仅 Q2 引用 GOAL-010。

### 成果（有证据）

- **五个检查点**：S1（I-011-001/002 冻结 verified）、S2（users/roles 后端闭环）、S3（records 退场 0006）、S4（I-011-003 v0.2.0 verified + 验收收据 + Renderer 零 diff）已勾选（progress `4/5`）；S5 为本关门审计覆盖（回归 + 交接），完成后 `5/5`。
- **意见台账**：A-001～A-009 全部响应；A-002（S1 conditional）两条 required 经 D-003 `fixed`、F-003～F-006 handled；A-007（I-011-003）三条 required 经 D-004 `fixed`、A-008 复审 pass；A-003/A-004（S2）、A-005/A-006（S3）、A-009（S4）同向 pass；本 scope **无开放 required**。
- **信息门禁**：`I-011-001`/`I-011-002`/`I-011-003` 全部 `verified`；无到期 required 信息项。
- **回归（S5）**：`go vet ./...` + `go test ./...`（apps/api）全绿；web `vitest` 23 files / **485/485** + `tsc -b` + `vite build` 干净；e2e `playwright` **2/2**；`git diff --check` 干净；I-011-003 §3 Renderer 基线 diff exit 0。
- **父级交接**：GOAL-010 02-execution 已落「GOAL-011 S4 证据交接」节（users/roles 替换、Renderer 零修改、records 退场、双资源验收收据）。
- **grok build 调用记录**：S1/S2/S3 独立审计（A-002/A-004/A-006）已产出；S4 与本次关门调用 grok 服务均被取消（stopReason cancelled，未产出书面意见），已在 02-execution 记录，不以未落盘意见作为放行依据。

### 对照成功标准

| S5 标准 | 状态 | 证据 |
|---------|------|------|
| API/Web/build/E2E 全量回归通过 | ✅ | go test/vet、vitest 485、vite build、playwright 2/2 |
| 阶段/关门审计通过 | ✅ | A-001～A-009 + 本关门审计；无开放 required |
| GOAL-010 S4 证据交接 | ✅ | GOAL-010 02-execution「S4 证据交接」节 |
| 目标自身无开放 required / 到期 required 信息项 | ✅ | 台账 + 信息表核对 |

### Findings

- **F-001 · grok build 对 S4/关门未产出独立书面意见**（severity: low；建议: recommended；status: open）
  - 描述：按用户指令调用 grok build 独立审计，S4 与关门各尝试多次均被 grok 服务取消（未写入 03-audit）。S1～S3 独立意见已落盘；S4/关门以 self 审计 + 既有 independent（A-007/A-008）覆盖。用户可复核后在需要时重跑 grok。
  - 影响：不阻断关门（self 关门审计满足检查清单；未落盘的独立意见不作放行依据）；作为用户可复核的待办记录。

### 必改项汇总（required 列表）

无。

### 结论 + 建议下一步

GOAL-011 五个检查点全部达成、无开放 required finding、无到期 required 信息项、回归全绿、GOAL-010 交接就绪。**pass**——符合关门条件，建议置 `done`（progress `4/5 → 5/5`）并同步 goal-tree；GOAL-010 可据此评估其 S4 勾选与 S5 关门（Root A-002 F-002-001 关闭证据链仍归 GOAL-010 S5）。

## A-011 · GOAL-011 S5 关门独立交叉审计（2026-08-04 · grok build 文本意见代贴）

- **source**：independent
- **auditor**：grok build（xAI · audit skill / 05-independent-audit.md；用户要求只读核验、意见以回复文本返回，本台账代贴并保留 source: independent）
- **类型 / scope**：close-out · GOAL-011 全目标关门（S1～S5；A-001～A-010 台账；progress 4/5）
- **verdict**：conditional

### 范围与区间

只读核验 `00-meta`、`02-execution`、`03-audit`、GOAL-010 交接节；抽查 `fixtures/schema/{users,roles}.json`、`migrate.go`（0006 + per-pending 快照）、`server_restart_test.go`（users+roles 重启）。未重跑 go/web/e2e（用户要求只读）。

### 成果（有证据）

| 主张 | 核验 |
|------|------|
| 台账 A-001～A-010 齐全；required 历史项已闭合 | A-002 → D-003 fixed；A-007 → D-004 fixed + A-008 pass；无开放 required |
| 三信息项 verified | I-011-001/002/003 均 verified |
| progress 4/5、S1～S4 勾选 | 与检查点及 goal-tree 一致 |
| S4 验收收据 | 02-execution 有 revision `73bc93a`、基线 `adfe15a`、命令表与 §2 对照 |
| GOAL-010 父级交接 | GOAL-010 02-execution「S4 证据交接」已落盘 |
| 关键代码证据存在 | users/roles fixture、0006 records_retire + per-pending 快照、server_restart_test |

### Findings

- **F-001 · S5 执行事实未写入 `02-execution`**（severity: medium；建议: **required**）
  - 描述：A-010 主张 S5 全量回归，但 02-execution 止于 S4 收据，末条仍为「计划（非事实）：S5…」。按 P-002，回归事实应以 02-execution 时间线可指回。
  - 证据：02-execution L99–123；对照 A-010「回归（S5）」节。
  - 建议闭合：补 S5 事实节（命令、目录、结果、是否复用 S4 收据或独立重跑、GOAL-010 交接指针、grok 取消记录）。
- **F-002 · 关门文档真相不一致**（severity: medium；建议: **required**）
  - 描述：① A-010 写「五个检查点全勾选」与 progress 4/5 / meta S5 `[ ]` 矛盾；② 00-meta L42 仍写「S1～S3 已完成、S4/S5 未完成、progress 3/5」；③ 03-audit「当前审计边界」仍写「S5 未实施」。
  - 建议闭合：修正 A-010 措辞 / meta 陈旧段 / 边界摘要，使「4/5 + S5 待投影 vs 关门就绪」单一真相后再 done/5/5。
- **F-003 · S4 阶段无 independent execution-facts 闭环**（severity: low；建议: recommended）——S4 验收仅 A-009 self；A-007/A-008 只管 I-011-003。本意见可作 S5 交叉输入，或事后补跑留收据。
- **F-004 · A-005 F-003 smoke disposable 仍无 S5 闭合证据**（severity: low；建议: recommended）——smoke.sh --disposable 未跑；非 S5 字面必选，但台账待办未清。
- **F-005 · roles 浏览器 e2e 未补**（severity: low；建议: recommended）——承接 A-009 F-001；后端/进程重启已覆盖 roles。

### 必改项汇总

| ID | 严重度 | 摘要 |
|----|--------|------|
| **F-001** | medium · required | 02-execution 补齐可核对的 S5 回归/交接事实 |
| **F-002** | medium · required | 消除 A-010「全勾选」/ meta `3/5` 残留 / 审计边界「S5 未实施」的三重不一致 |

无 high required。产品侧 S1～S4 主链与历史 required 闭合成立；阻断的是关门投影前的证据与文档一致性。

### 与既有意见

同意 A-010：无开放 required、三信息项 verified、GOAL-010 交接已落、S1～S4 有据。**不同意** A-010 无条件 pass 即可 done：F-001/F-002 使 close-out 证据链与文档真相未达「可无条件置 done」标准。

### 结论 + 建议

**verdict: conditional**——GOAL-011 产品交付与历史门禁主体可核对通过，尚不可无条件关门置 done/5/5。先经 /govern 闭合 F-001、F-002（补 S5 执行事实 + 对齐 meta/台账/A-010 表述），再评估放行；F-003～F-005 可 handled 或不阻断。

### 声明

本意见 source: independent；按用户要求以文本返回，由编排器代贴入本台账；未修改 00-meta/goal-tree/契约/产品代码。

## 响应 A-011（self · 编排响应 · 2026-08-04）

A-011（independent · conditional）F-001/F-002（required）经 `fixed` 闭合：

### 关闭证据表

| Finding | 严重度 | 状态 | 证据路径 |
|---------|--------|------|----------|
| A-011 F-001 · S5 执行事实未入 02-execution | medium · required | **fixed** | `02-execution.md` 新增「2026-08-03 · S5 回归、审计与父级交接」节：全量回归命令表（go vet/test、vitest 485、vite build、playwright 2/2、Renderer 基线 diff exit 0、git diff --check）+ 关门审计 + grok 重试 A-011 + 父级交接指针 |
| A-011 F-002 · 关门文档真相不一致 | medium · required | **fixed** | ① A-010 措辞改为「S1～S4 已勾选（4/5）；S5 为本关门审计覆盖」；② `00-meta.md` L42 陈旧段更新为「S4 已达成（4/5），S5 完成后 done（5/5）」；③ `03-audit.md`「当前审计边界」更新为「S5 回归、关门审计与父级交接已执行」 |
| A-011 F-003 · S4 无 independent execution-facts 闭环 | low · recommended | **handled** | 本意见（A-011 代贴）即 S4/S5 交叉输入；S4 收据已含完整命令/结果 |
| A-011 F-004 · smoke disposable | low · recommended | **handled** | 需 `--disposable` Compose 隔离环境；S5 收据不含；作为用户可复核的待办（非 S5 字面必选，API/Web/build/E2E 已全量） |
| A-011 F-005 · roles 浏览器 e2e | low · recommended | **handled** | roles 由 handler/store/进程重启测试覆盖；浏览器级补测为可选加固 |

### 趋同说明

A-010（self · pass）与 A-011（independent · conditional）经 F-001/F-002 修复后趋同：S5 无开放 required，可置 `done`。GOAL-011 置 `done / 5/5`，GOAL-010 据此可评估其 S4 勾选与 S5 关门；Root A-002 F-002-001 仍 open（GOAL-010 S5 关闭证据链）。

## A-012 · GOAL-011 S5 关门独立代码审计（2026-08-04）

- **source**：independent
- **auditor**：Codex（GPT-5）
- **类型 / scope**：close-out · 代码与流水线专项复核（VP-002 持久化 users/roles、真实鉴权/会话生命周期、Schema 管理面、records 退场洁净度）
- **verdict**：fail

### 范围与方法

本意见复核当前 HEAD 1d86081 所代表的 GOAL-011 关门状态，先核对 workspace-002 的 workspace binding、VP-002/Charter 对齐、GOAL-011 五件套、I-011-001/002/003 和 A-001～A-011，再抽查实现与活动 CI。当前目标文档与 goal-tree 已写为 done / 5/5，工作区绑定有效，shared materials 为 none；这些状态本意见不修改。

本次关注点是用户明确要求的“作为 VP-002 预期的生产可用 users/roles 设计和实现”以及 records 历史代码是否仍污染运行面。VP-002 允许初期不提供完整可视化 IAM，但要求持久化用户/角色/权限或菜单关系、可用超级管理员/基础角色、真实请求鉴权和可支撑业务使用的接口或种子机制；本意见没有把“完整 IAM”偷换成 GOAL-011 的必做范围，而是检查当前已暴露的管理接口和 Schema 是否自洽、安全、可用。

### 已核对的正向事实

- API 使用真实 JWT/opaque refresh、请求级身份和持久化 role-permission/menu 投影；匿名和 viewer 的 401/403 基线仍通过。
- users/roles 资源的敏感字段隔离、角色存在性校验、system role 与 in-use 保护、seed 幂等和历史迁移测试均存在。
- go test ./... -count=1、go vet ./...、npm test（23 files / 485 tests）、npm run build 均通过。
- 浏览器 E2E 在本机默认 5173 受 listen EACCES 阻断；按仓库配置改用 WEB_PORT=9999 后 users CRUD 与 auth shell 两项均为 2/2 通过。现有 E2E 没有 roles 浏览器管理路径。

这些正向事实不能抵销下列会影响生产授权、凭据生命周期、关门 CI 和 records 退场洁净度的开放 required findings。

### Findings

- **F-001 · users.write 可越权授予 admin（severity: high；建议: required；status: open）**
  - **描述**：users 资源的所有写操作只经过 users.write 一把门；create/patch 的 roles 字段允许写入任意已存在 role key。UpdateUser 只检查 self demote 和最后管理员，不检查操作者是否拥有角色管理权限、目标角色的有效权限是否超过操作者，也不区分“编辑用户资料”和“分配角色”。
  - **证据**：apps/api/internal/handler/users.go:28-40,88-105,112-132；apps/api/internal/store/users.go:220-247；通用工厂写门 apps/api/internal/handler/resources.go:385-430。权限投影实际来自 normalized role_permissions（apps/api/internal/store/store.go:253-278）。
  - **独立复现**：在隔离 production 配置中创建 user-manager 角色，只给它 users.write（不含 roles.write），以该身份创建 roles: ["admin"] 的用户；请求返回 201，新用户以单字符密码即可登录并获得 admin 权限。该场景使用 VP-002 明确允许的 seed/持久化 grant 机制，不是伪造 Bearer token。
  - **影响**：任何被委派为用户管理员的角色都能建立超级管理员，形成稳定的权限提升路径；roles.write 看似独立但实际上不能保护角色分配。
  - **建议闭合**：拆出独立的 role-assignment 权限并在服务端校验；至少要求 actor 拥有 roles.write/受保护的 admin 能力，并禁止授予 actor 无权委派的有效权限集合。为“仅 users.write、无 roles.write”的身份补 create/patch 负向测试，覆盖不能创建/提升 admin。

- **F-002 · 密码与会话生命周期不满足生产安全基线（severity: high；建议: required；status: open）**
  - **描述**：users Schema 将密码字段声明为普通 input（apps/api/internal/handler/fixtures/schema/users.json:46-52），Renderer 的 input 分支渲染为 type="text"（apps/web/src/renderer/form-controls.tsx:283-303），输入密码在 UI 中明文可见。后端通用字符串解码会对密码 TrimSpace（apps/api/internal/handler/resources.go:303-322,346-370），并只要求非空，因此单字符密码被接受。PATCH 改写 bcrypt hash 后没有撤销该用户既有 refresh tokens（apps/api/internal/handler/users.go:118-124、apps/api/internal/store/users.go:253-269；token 清理只出现在删除路径 291-327）。
  - **独立复现**：单字符 m 创建的用户登录成功；管理员改密后，用改密前取得的 refresh token 调用 /api/auth/refresh 仍返回新的 access/refresh pair。
  - **影响**：管理页面暴露凭据；泄露的长期 refresh token 在改密后继续可用，无法把“改密”作为实际会话失效边界。VP-002 的真实认证标准明确包含 token 生命周期、失效与错误处理；本 finding 不要求一次性引入完整密码策略，但要求最小长度/边界、保留密码字节语义和可验证的会话撤销规则。
  - **建议闭合**：使用 password control；密码字段不得走通用 trim；服务端设置明确的最小长度与长度上限并补负向测试；密码变更在同一事务撤销该用户所有 refresh tokens，并为 access-token 的剩余 TTL 或版本号语义写明边界和回归测试。

- **F-003 · users/roles Schema 暴露的是不完整且误导性的角色管理（severity: medium；建议: required；status: open）**
  - **描述**：Users 页面 create body 只映射 username/name/password，edit 只映射 name，因而不能从 Schema 页面分配已有角色或完成密码变更（apps/api/internal/handler/fixtures/schema/users.json:15-73）。Roles 页面只能创建/改名/删除 key/name（apps/api/internal/handler/fixtures/schema/roles.json:15-71）；自定义角色按实现永远没有 grants（apps/api/internal/store/roles.go:1-5），所以新建角色对任何授权都没有实际语义。系统角色的 edit/delete 行按钮也未按 system 行字段禁用：表格对每行无条件渲染 action（apps/web/src/renderer/schema-table.tsx:244-266），点击只会得到 409，而契约已明确指出前端可据 system 禁用（I-011-001 §3.0）。
  - **边界判断**：I-011-001 §3.4 将完整 grant CRUD 列为非目标，VP-002 也允许用 seed-only 最小 IAM；问题不是强行要求完整 IAM，而是当前产品同时展示“可管理 roles”的页面，却不能把用户加入可用角色、不能表达自定义角色权限、还对不可变系统角色展示必失败操作。若坚持 seed-only，应收窄产品声明并隐藏/移除 inert custom-role 管理面；若保留该管理面，则必须补齐角色分配与有效 grant/菜单来源。
  - **影响**：业务管理员无法通过实际页面完成最基本的用户授权配置；新增角色只是标签，且角色列表的系统行操作不具备可用 UX。现有 2/2 浏览器 E2E 只覆盖 users CRUD，不覆盖 roles 管理。
  - **建议闭合**：选择并记录上述两条产品边界之一；推荐实现 users 角色分配、system/in-use 行级禁用和至少一个可验证的自定义角色授权/种子路径，并补真实浏览器 roles 流程及“无 grant 角色不能写入”的 API/UI 断言。

- **F-004 · 活动 CI 与运行时代码仍依赖已退役 records（severity: high；建议: required；status: open）**
  - **描述**：r6-basic-matrix.yml 仍在 dev/main push 和 PR 上运行；container smoke 传入已废弃的 SMOKE_RECORD_ID，并要求 seed total=8（.github/workflows/r6-basic-matrix.yml:91-96），而 scripts/smoke.sh:20-35 的文档写默认 1、实现却仍默认 8。fresh seed 只插入一个 admin 用户（apps/api/internal/store/store.go:124-138，RBAC seed 不新增用户）。更直接地，workflow 的 persistence step POST/GET /api/records（.github/workflows/r6-basic-matrix.yml:98-115），而当前注册表只挂 /api/users 和 /api/roles（apps/api/internal/handler/health.go:21-34）；隔离 API 复现该 POST 为 404。
  - **运行时代码残留**：apps/web/src/renderer/records.ts 仍是被生产组件导入的模块，导出 fetchRecords、buildRecordsQuery、parseRecordList、readRecordApiError、RecordApiError 等历史命名（records.ts:5-12,28-50,99-176），schema-table.tsx 与 render.tsx 继续依赖它。API 工厂和 auth/store 注释仍把 records 当作首个/门禁语义；apps/api/internal/store/operations.go:14-32 还保留只供历史测试使用的 EventRecord 常量。它们不是可运行的 records 资源，却构成用户要求清除的无实际语义历史污染。
  - **合法保留边界**：migrate.go 的 0003/0005/0006 迁移、operation_log 的历史 CHECK 值、迁移/升级测试和明确标注的历史文档可保留，用于既有数据库升级与审计链；不能据此保留活动 CI、生产模块名或默认运行依赖。
  - **影响**：活动流水线在当前代码上必然访问不存在的 endpoint/错误种子基线，不能作为 GOAL-011 或 VP-002 的生产交付证据；records 语义残留还会继续诱发错误接入和后续维护误判。
  - **建议闭合**：把 persistence smoke 改为 users/roles 的真实重启断言，修正 seed total/废弃环境变量；将通用 Web transport 重命名为 resource 语义并更新所有 imports/tests；移除无运行用途的 EventRecord 常量和过时注释。增加 scoped grep/CI check，明确只允许迁移历史、历史测试与已标注 fixture 命中 records。

- **F-005 · 删除用户的最后管理员不具备完整事务原子性（severity: medium；建议: required；status: open）**
  - **描述**：DeleteUser 先在事务外调用 GetUser（apps/api/internal/store/users.go:291-302），随后才 Begin，并在事务内依据事务外读取的 cur.Roles 决定是否执行最后管理员计数（users.go:304-317）。单连接只串行 SQL，不会把这两个高层操作合并为一个不可插入的临界区。
  - **可复核交错**：目标用户初始为非 admin；删除路径读完后，另一写请求把目标提升为 admin，再把原 admin 降级；删除事务使用旧的非-admin 快照而跳过计数，最终可留下零个 admin 用户。这与契约 I-011-001 §2.4/§7 所要求的“在 store 事务内判定”不一致。
  - **影响**：并发管理请求下失去最后管理员保护，锁死后续管理面；现有串行测试不能发现该竞态。
  - **建议闭合**：把目标读取、自身检查、admin 计数和 DELETE 放入同一事务/锁定语义，检查受影响行数，并加入可控并发回归测试。

- **F-006 · legacy roles JSON 双写继续扩大数据一致性面（severity: low；建议: recommended；status: open）**
  - **描述**：users.roles JSON 与 normalized user_roles 同时写入并在每次读取时比较（apps/api/internal/store/store.go:203-227、apps/api/internal/store/users.go:258-284）。这是历史迁移兼容策略，不是本次 close-out 的单独阻断项，但它让角色语义拥有两个持久化写入面，且任何外部修复都可能触发 fail-closed mismatch。
  - **建议**：在后续迁移中明确唯一权威来源并淘汰 legacy 列，或把双写/校验边界和修复工具作为正式运行手册与回归契约。

### 必改项汇总

| Finding | 严重度 | 状态 | 关门影响 |
|---------|--------|------|----------|
| **F-001** | high · required | open | 用户写权限可创建超级管理员，阻断生产授权模型 |
| **F-002** | high · required | open | 密码可见、弱校验，改密后 refresh 仍有效，阻断凭据生命周期 |
| **F-003** | medium · required | open | Schema 管理面不能完成实际角色语义，阻断“可支撑业务使用” |
| **F-004** | high · required | open | 活动 CI/运行时代码仍依赖 records，阻断退场洁净度与交付证据 |
| **F-005** | medium · required | open | 最后管理员保护存在竞态，阻断账户安全不变量 |

### 与既有意见的关系

A-010/A-011 主要复核 S5 文档事实链、历史 required closure 和既有回归收据；本意见新增了代码级授权组合、凭据生命周期、Schema 语义和活动 CI 的专项范围，且实际复现了既有测试未覆盖的路径。因此本意见不是把 A-011 已 fixed 的文档 finding 重新打开，而是对同一 close-out 的新增 independent evidence。按 P-003，新的 required finding 在合法 fixed、accepted-residual 或 user-overruled 之前仍阻断无条件关门。

### 结论 + 建议给编排器/用户的下一步

**verdict: fail**。当前 GOAL-011 虽在目标树中为 done / 5/5，但不满足用户要求的生产级 users/roles 设计复核，也不能证明 records 已从活动交付面洁净退场。建议 /govern 逐项响应 A-012 F-001～F-005，优先走 fixed：先修复授权边界、凭据/会话撤销、Schema 角色语义、CI/records 残留和最后管理员事务；完成后再进行一次限定范围独立复审。对于 F-003，如用户选择继续保持 seed-only 最小 IAM，必须在治理记录中明确收窄产品声明并移除误导性的 inert 管理面，不能把当前页面直接当作生产可用角色管理。

本意见不修改 00-meta、01-decision、02-execution、goal-tree、代码或任何 Goal status/progress；状态推进、finding 响应和是否接受残余风险归 /govern 按 P-004/P-005 裁决。

## 响应 A-012（self · 编排响应 · 2026-08-04 · GOAL-011 D-005）

- **响应性质**：这是 `/govern` 对 A-012 的第一轮意见汇总、门禁恢复与待裁决登记，不是新增 self audit，不产生新的 A 编号或 verdict。
- **工作区边界**：`workspace-002-production-admin-foundation` / Root `GOAL-001-production-admin-foundation` / canonical scope 已校验；primary plan 为 active VP-002；`shared_materials_catalog: none`，本响应未使用共享资料。
- **冲突**：A-010（self · pass）与 A-011（independent · conditional，经 F-001/F-002 fixed 后趋同）支持原 close-out；A-012（independent · fail）对同一 close-out 新开五条 required。按 P-004.2，未裁决前不取乐观侧、不维持无条件关门。

### 响应状态表

| Finding | 级别 | 本轮状态 | 响应与证据边界 |
|---------|------|----------|----------------|
| A-012 F-001 · users.write 可越权授予 admin | high · required | **open** | 静态抽查确认 `usersResource` 只要求 `users.write`，roles 值只做存在性校验；尚无权限委派修复或负向回归 |
| A-012 F-002 · 密码与会话生命周期 | high · required | **open** | 静态抽查确认普通 input、通用 trim 与改密不撤销 refresh；尚无密码控件/长度边界/同事务撤销或 access-token 边界收据 |
| A-012 F-003 · users/roles 管理语义不完整 | medium · required | **open · 待 I-011-004 裁决** | 两条合法产品路径尚未选择：seed-only + 移除/隐藏 inert 管理面；或补齐角色分配、grant 来源与 roles 流程 |
| A-012 F-004 · 活动 CI/生产专名仍依赖 records | high · required | **open** | 静态抽查确认 workflow `/api/records`、废弃 env/seed total 与生产 transport 专名；尚无 CI/grep/重启收据 |
| A-012 F-005 · 最后管理员删除竞态 | medium · required | **open** | 静态抽查确认 `GetUser` 在事务外；尚无事务内读取/删除修复或可控并发回归 |
| A-012 F-006 · legacy roles JSON 双写 | low · recommended | **open · non-blocking** | 后续迁移/运维契约建议；不得与 F-001～F-005 混成同一关门阻断 |

### 冲突裁决与状态投影

- 用户本轮明确要求“响应 A-012”，但尚未书面选择 `fixed`、`accepted-residual` 或 `user-overruled`。因此本响应只记录真实的 open 状态，不把沉默解释为 residual/overruled。
- 编排器建议：F-001、F-002、F-004、F-005 走 `fixed`；F-003 先完成 `I-011-004` 产品边界裁决，再按所选边界走 `fixed`。不建议对 high security/CI findings 采用 residual。
- 因五条 required 均未合法闭合，GOAL-011 从 `done` 恢复为 `active`。原 S1～S5 检查点事实保留，派生 progress 仍为 `5/5`；该数字不构成关门放行。
- GOAL-010 保持 `active / 3/5`，Root 与 VP-002 保持 active；本响应不关闭父级/Root finding，也不创建跨区或愿景层状态。

### 仍开放项与下一步

1. 用户裁决 A-010/A-011 与 A-012 的 close-out 冲突：推荐采纳 A-012 的五条 required 并全部走 `fixed`。
2. 用户裁决 `I-011-004`：推荐补齐用户角色分配、可验证 grant/seed 来源、system/in-use 行级禁用与 roles 浏览器路径；备选为明确保留 seed-only 并移除/隐藏无授权语义的 custom-role 管理面。
3. 裁决落盘后，先冻结整改契约/验收矩阵，再实施并记录 API/Web/CI/E2E 事实；最后请求仅覆盖 A-012 F-001～F-005 的 `/audit` finding-closure 复审。

### 响应结论

A-012 已被 `/govern` 正式认账并进入 fail-closed 整改门禁，但 **F-001～F-005 仍全部 open**，A-010/A-011 与 A-012 尚未趋同。当前不具备重新关门条件。
