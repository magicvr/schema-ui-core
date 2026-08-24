---
id: A-002
doc: audit-entry
goal: GOAL-001-account-email-identity
source: independent
status: recorded
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
scope: Root 关门审计（R1～R4 四阶段汇总 · 子目标审计链 · 五判据兑现 · 信息门禁闭环 · 边界不越界）
verdict: conditional
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-002 · Root 关门独立审计（independent · 2026-08-24）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · 思考强度 high）
- **类型** / **scope**：close-out · Root 关门——R1～R4 四阶段汇总、子目标审计链、五判据兑现、信息门禁闭环、边界不越界（代码基准 `6c6496d4`；Root 台账收口至「待 A-002」，见 E-009）
- **verdict**：**conditional**
- **开放 required**：1（F-001）

### 范围与区间

| 项 | 值 |
|----|------|
| 工作区 | `workspace-018-account-email-identity`（`workspace.md`：`root_goal` = `GOAL-001-account-email-identity`；`canonical_scope` 匹配本目录；`shared_materials_catalog: none`；`primary_plan` = `VP-018-account-email-identity`） |
| 被审目标 | `GOAL-001-account-email-identity`（Root；`parent: null`） |
| audit_type | close-out |
| 对照 | Root `00-meta` 成功标准 1–4；VP-018 方向级退出判据 1–5；`GOAL-005/attachments/r4-evidence.md` §1–§5 |
| 信息门禁 | I-001～I-006（见下「信息项」）；N-1 有界残余声明 |
| 共享资料 | 无（`none`；未把资料目录当事实或关闭证据） |
| 代码基准 | 产品代码 = `6c6496d4a72eb33774c9b18aa11832b7e5161644`（R4 证据 + 两阶段派发）。本会话 `HEAD` = `665e1e47`，相对 `6c6496d4` 仅 8 个治理文档（Root 台账收口 / A-001 self / E-008 / E-009），**产品代码零 diff** |
| 本审计未改 | 目标 `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码或测试 |

未读取或比较其他工作区上下文。项目级 `docs/architecture/independent-audit-execution.md` 要求 self 之后 grok independent；A-001 self 已落盘。本条按用户指定独立关门执行。

### 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md` id / root_goal / canonical_scope / `plan_refs`+`primary_plan` 与本目标路径一致；资料目录 `none` |
| R1 合同冻结 | GOAL-002 `done` 3/3；D-001 七条款；A-001 self **pass**（0 required） |
| R2 schema | GOAL-003 `done` 4/4；迁移 0054；A-001 independent **pass**（0 required）；F-003 已响应补登 goal-tree |
| R3 绑定流 | GOAL-004 `done` 4/4；迁移 0055；A-001 independent **conditional** → F-001 **fixed**（`bd1cdff9`，本会话复跑 `TestUsersPatchEmailPrefillFlows` PASS） |
| R4 证据 | GOAL-005 正文/E-002/goal-tree 主张 done 3/3；A-001 self **pass**；e2e `TestR4EndToEndBindVerifyThroughMockChannel` 本会话 PASS（YAML 滞后见 F-001） |
| 两阶段派发 + 补偿 | `email_identity.go` `BindEmail`/`ResendEmailCode`：Send 在占槽事务外；失败走 `compensateBind` / 挑战快照恢复。`TestSendFailureRollsBindBack` 本会话 PASS |
| I-006 HTTP 代填 | `users.go` `RawStringFields: ["password","email"]`；PATCH 非空 → pending、`""` → unbound、非字符串 400、冲突 409 `EMAIL_TAKEN` |
| 本会话复跑 | `go test ./internal/modules/authsession/ -count=1` **ok**（3.879s）；`go test ./internal/store/ -run 'TestMigrate0054\|TestMigrate0055\|Postgres' -count=1` **ok**（17.921s；无 `TestMigrate0055` 函数，0055 专项仍缺，见 N-4）；`TestR4EndToEnd*` / `TestRebindOverwriteReleasesOldSlot` / `TestSendFailureRollsBindBack` / `TestUsersPatchEmailPrefillFlows` **PASS**；web `vitest run src/components/email-identity.test.tsx` **3 passed** |

### 对照成功标准

对照 Root `00-meta` 成功标准 1–4 与 VP-018 退出判据 1–5（权威映射：`GOAL-005/attachments/r4-evidence.md`）。

| # | 标准 | 状态 | 证据 |
|---|------|------|------|
| 1 | `users` 可持有可空邮箱与可核对校验状态；无邮箱账号仍能登录 | **达成** | 0054 `email TEXT` 可空 + `email_status` CHECK；e2e 中 bob 全程 unbound 且 `UserByID` 正常；`apps/api/internal/auth` 无 email 依赖（登录面不读该列） |
| 2 | 绑定/校验流落地；校验信经 `kernel.MailSender`；无生产渠道时从 VP-017 当时默认渠道取信 | **达成** | e2e 走真实 `mail.NewOutboxSink`（非测试桩）写入 `mail_outbox`，从记录 body 取 6 位码完成 Verify；常量 TTL 10 分钟 / 冷却 60 秒 |
| 3 | 非空邮箱唯一性 fail-closed；换绑走同一校验合同 | **达成** | e2e 他号大小写折叠 → `ErrEmailTaken`；`TestRebindOverwriteReleasesOldSlot` PASS；同址 pending 重绑套用 60s 冷却（F-002 后） |
| 4 | 未引入忘记密码状态机 / 邀请 / 密码策略产品 / SMS / 第二运输 / 模板中心；未改 Charter；未改 Profile 默认集 | **达成** | 本波 diff 无恢复/邀请/SMS/模板中心产品面；Charter 仍 `schema-ui-core-admin-foundation@0.2.0`、`primary_workspace` 仍 workspace-001；`profileDefaults` 模块列表未增删，仅既有 `admin.account` 路由贡献追加 bind/verify/resend 三端点（不是默认集扩模块）；`settings.json` `tab-mail` / `mail-admin-tab` 仍为 017 出站渠道管理，非身份绑定复制品 |
| 5 | 开放 required finding = 0（VP-018 判据 5） | **本条落盘后未满足** | 子目标历史 required 已合法闭合；**本条新开 F-001 required/med**（GOAL-005 / Root 台账 YAML 与关门主张冲突）。关闭前不得把判据 5 写成已兑现 |

`progress` 百分比不作为本条闭合或放行依据（P-001）。

### 信息项（P-005）

| ID | 级别 | 台账状态 | 本审计结论 |
|----|------|----------|------------|
| I-001 | required | verified（GOAL-002 D-001） | **门禁关闭成立**。会话留痕 `i001_slot` / `i001_norm`（E-002）；实现 = 绑定即占槽 + `lower(email)` 唯一 + 原样存储 |
| I-002 | required | verified（GOAL-002 D-001） | **门禁关闭成立**。会话留痕 `i002_form`；实现 = 6 位验证码，时效归 I-005 |
| I-003 | required | **registered**（VP 冻结投影；Root `00-meta` 权威表） | **不构成到期未关**。VP-018 已冻结「可空」；D-003 维持 registered 投影。Root `03-audit.md` 写「registered→verified 同步」**过称**（见 F-002） |
| I-004 | required | **registered**（VP 冻结投影） | 同上；换绑进分母为 VP 冻结，GOAL-002 D-001 §5 为单列机械推导。不需要另一次会话 verified |
| I-005 | required | verified（GOAL-004 D-001） | **门禁关闭成立**。会话留痕 `i005_ttl` / `i005_cooldown`；常量 `emailCodeTTL=10m` / `emailResendCooldown=60s` 与测试一致 |
| I-006 | non-blocking | verified（GOAL-004 D-001） | **正向 HTTP 面已接通**（相对 GOAL-004 A-001 F-001）。会话留痕 `i006_admin` + `scope_ui`；PATCH 代填 → pending，无直达 verified |

「三次用户书面裁决」可核对落款为：**两轮 `ask_user_question`**（R1 三项 `i002_form`/`i001_slot`/`i001_norm`；R3 四项 `i005_ttl`/`i005_cooldown`/`i006_admin`/`scope_ui`）外加 **D-003 解冻确认**。台账「三次 / R1 两项」与会话次数不完全同构，但不构成门禁缺口——条款已写入 D-001 且与实现一致。无独立会话誊本文件；本仓库以决策/执行条目中的问题 id 为书面留痕，与既有实践一致。

**N-1**（SQLite `lower()` 仅 ASCII）：GOAL-005 `attachments/r4-evidence.md` 已作有界残余声明，复核触发 = 「出现非 ASCII 邮箱真实需求或 PG→SQLite 双向迁移产品化」；归属后续 IAM 波次。原 finding 为 recommended/note 而非 required，**不**触发 P-003 `accepted-residual` 用户书面接受强制路径。声明本身范围与触发可核对。

无到期未关的 required 信息项。无共享资料引用。

### Findings

#### F-001 · GOAL-005 与 Root `00-meta` YAML 与关门主张冲突

- 严重度：med
- 建议：**required**
- 状态：open
- 描述：E-002 / E-008 / goal-tree 主张 GOAL-005 `done · 3/3`、Root 纲领 **4/4**。但 GOAL-005 `00-meta.md` YAML 仍为 `status: active`、`progress: 0/3`（正文检查点表已标 3/3 已关门）。Root `00-meta.md` YAML 自开区起一直 `progress: 0/4`，从未随 R1～R4 正文检查点重算；P-001 要求 `00-meta` 与 `goal-tree` 保存同一派生结果。这不是「independent 落盘前的时序滞后」（R2/R3 F-003 类），而是 **/govern 已宣称子目标关门与 Root 4/4 之后仍未改 YAML**。`status` 是生命周期字段（GOAL-005 YAML=`active` vs goal-tree=`done`），不只是派生百分比。在 YAML 与五件套/goal-tree 对齐之前，**不得**把 Root 标 `done`（判据 5 亦未满足）。本条**不**否定产品实现与测试事实。
- 证据：`GOAL-005-r4-evidence/00-meta.md` YAML L4/L9 vs 正文 L32「已关门」与 `goal-tree.md` 状态表；`GOAL-001-account-email-identity/00-meta.md` YAML L9 `progress: 0/4` vs 正文 L35「当前 **4/4**」与 `goal-tree.md`；`git log -p` 显示 R1～R4 历次只改正文分数、从未改 YAML `progress`。关联 P-001 / AGENTS §7。

#### F-002 · I-003/I-004 与 GOAL-004 审计索引口径滞后

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：(a) Root `03-audit.md` 信息表写 I-001～I-006「全部 verified」且「I-003/I-004 registered→verified 同步」。权威表 `00-meta` 仍为 **registered**（VP 投影），与 D-003「维持 registered 投影」一致。过称不构成到期未关，但关门响应应改回 registered，避免把投影项伪装成会话 verified。(b) GOAL-004 `03-audit.md` 索引行仍写 A-001 开放 required = 1；结论节已记录 F-001 `fixed`（`bd1cdff9`），本会话 HTTP 测试 PASS。索引未刷新不等于 F-001 重新打开。
- 证据：Root `00-meta.md` I-003/I-004 行；Root `03-audit.md` 信息就绪表；GOAL-004 `03-audit.md` 索引 vs 结论状态；`handler/users_email_test.go`。

#### N-1 · SQLite `lower()` ASCII 有界残余（维持）

- 严重度：low（note）
- 建议：recommended（无需本门禁动作）
- 状态：open（有界残余；非本波 required）
- 描述：唯一性仍走 SQL `lower(email)`。证据包已声明范围、理由与复核触发。应用层 `EqualFold` 只用于同号已 verified 幂等，不参与他号占槽。不升 required。
- 证据：`GOAL-005/attachments/r4-evidence.md` N-1 节；`email_identity.go` 他号 COUNT；`migration.go` 0054 注释。

#### N-2 · 两阶段派发补偿失败分支（维持 GOAL-005 A-001 F-1）

- 严重度：low（note）
- 建议：recommended（无需本门禁动作）
- 状态：open（已知边界）
- 描述：Send 失败后补偿成功路径有测试（状态回滚到先前身份）。补偿**自身**失败时返回 `fmt.Errorf("%w: … (compensation failed: %v)", ErrEmailSendFailed, …)`，`errors.Is(..., ErrEmailSendFailed)` 仍成立，HTTP 可归一为 `EMAIL_SEND_FAILED`。无补偿失败专项测试；概率低。GOAL-005 A-001 已留痕。不升 required。
- 证据：`email_identity.go` `BindEmail` L211–215、`ResendEmailCode` L440–454；`email_identity_test.go` `TestSendFailureRollsBindBack`（只覆盖补偿成功）。

#### N-3 · VP-018 投影未随实现层解冻/关门更新

- 严重度：low（note）
- 建议：recommended（交 `/vision`，非本 Root 产品缺口）
- 状态：open
- 描述：`VP-018-account-email-identity.md` 仍写 lead Root **`blocked`**；信息表 I-018-001～006 仍 collecting/registered。实现层已解冻并跑完 R1～R4。不阻断本条产品核对；VP 关门是独立 `/vision` 流程。
- 证据：`docs/vision/plans/VP-018-account-email-identity.md` 状态表 / 信息需求表。

#### N-4 · 0055 无语义专项；R4 e2e 未在 PG 方言单独实跑

- 严重度：low（note）
- 建议：recommended（可选，承 GOAL-003 F-001 / GOAL-004 N-6 / GOAL-005 F-2）
- 状态：open
- 描述：无 `TestMigrate0055`。本会话 store 过滤含 `Postgres` 的套件 **PASS**（17.9s，含 live-PG 集成的时间量级）。R4 e2e 仍仅 SQLite + 同一 `OutboxSink` 抽象。不构成合同名不副实。
- 证据：本会话 `go test ./internal/store/ -run 'TestMigrate0054|TestMigrate0055|Postgres'`；仓库无 `TestMigrate0055`；`email_identity_e2e_test.go`。

### 必改项汇总

1. **F-001（required / med）**：对齐 GOAL-005 `00-meta` YAML（`status: done`、`progress: 3/3` 与正文检查点一致）以及 Root `00-meta` YAML `progress` 与正文纲领表同一派生结果（当前应为 `4/4`，在 Root 仍 `active` 等待本条响应期间）。同步后 goal-tree 与 YAML 不得再互否。在此之前 **不得**将 Root 标 `done`，也 **不得**把 VP-018 判据 5 写成已兑现。

开放 **required** = **1**。

### 与既有意见的异同

对照 Root A-001 self **pass**：同意五判据的**产品**证据链（e2e / 0054 / 0055 / I-006 HTTP / 边界）。**不同意**「台账已收口、判据 5 已满足」——self 未核对 YAML `status`/`progress` 与 goal-tree 互否，也把 I-003/I-004 过称为 verified。

对照 GOAL-004 A-001 independent **conditional**：F-001（HTTP 代填）本会话确认为 **fixed**（`bd1cdff9` + `TestUsersPatchEmailPrefillFlows`）。F-002/F-004 已在 D-001 v1.1.0 口径对齐。N-1/N-2 维持残余，不假装闭合。

对照 GOAL-003 A-001 independent **pass**：F-003 类台账滞后当时发生在 independent **之前**；本条 F-001 发生在 **/govern 已宣告关门之后**，故升为 Root 关门 required，而不是再标 recommended。

对照 GOAL-002 A-001 self：合同条款与实现映射可核对；F-2 ASCII `lower()` 已在 R4 证据包声明复核触发。

### 结论 + 建议给编排器/用户的下一步

**conditional** —— R1～R4 产品交付、子目标审计链、五判据 1–4、I-001/I-002/I-005/I-006 会话裁决与实现、N-1 残余声明、边界红线，本会话均可核对；**不能**无条件把 Root 标 `done`，因为 GOAL-005 YAML 仍 `active`、Root YAML `progress` 仍 `0/4`，与关门主张和 goal-tree 冲突（F-001 required/med）。

建议 `/govern`：

1. **先响应 F-001（fixed）**：只改 GOAL-005 与 Root `00-meta` YAML（及若需要的索引口径 F-002），使 `status`/`progress` 与正文检查点、goal-tree 同一派生结果。不要用正文分数或 goal-tree 单独声称覆盖 YAML。
2. F-002 / N-1～N-4 按 recommended/note 处理。N-1 维持有界残余（已有复核触发）；N-3 交 `/vision` 更新 VP-018 投影（Root `done` 之后）。
3. 不要由本审计代改 `status` / 检查点 / `progress` / goal-tree。F-001 闭合且开放 required = 0 后，Root 才具备 `done` 条件；VP-018 `closed` 为独立愿景流程。
4. 本条与 A-001 self 在产品事实上同向、在台账完备性上冲突。冲突范围限于 YAML/索引口径，**不**否定实现。按 P-004：无「一要一否」的产品必改项；用户若要在 YAML 未改时强行 Root `done`，须对 F-001 走 `user-overruled` 并留痕（不建议）。

### 声明

本意见 `source: independent`，**不修改** `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码。响应、finding 闭合与关门状态变更由 **`/govern`** 与用户书面裁决处理。
