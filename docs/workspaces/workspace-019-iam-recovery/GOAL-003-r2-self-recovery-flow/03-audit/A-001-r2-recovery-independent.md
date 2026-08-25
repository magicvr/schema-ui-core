---
id: A-001
doc: audit-entry
goal: GOAL-003-r2-self-recovery-flow
source: independent
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
scope: R2 自助恢复全链实施切片（C2 迁移 0056 · C3 域/HTTP/MFA/会话 · C4 mock 取码 e2e + Web 登录页恢复流）；execution-facts 兼对照 D-001
verdict: conditional
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-001 · R2 自助恢复实施切片独立审计（independent · 2026-08-25）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：execution-facts（兼对照 design-plan = 本目标 D-001）· C2 迁移 0056 双方言 + 黄金断言；C3 `core.auth-session` 恢复域、`POST /api/auth/recovery/start|complete`、`mfa.Service.VerifySecondFactor`、UpdateUser 会话语义；C4 `recovery_e2e_test.go` mock 渠道取码 + `LoginPage` 恢复流
- **verdict**：**conditional**
- **开放 required**：1（F-001）

### 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-019-iam-recovery`（`workspace.md`：`id` 匹配；`root_goal` = `GOAL-001-iam-recovery`；`canonical_scope` = 本目录；`shared_materials_catalog: none`；`primary_plan` = `VP-019-iam-recovery`） |
| 被审目标 | `GOAL-003-r2-self-recovery-flow`（同区；`parent` = `GOAL-001-iam-recovery`） |
| audit_type | execution-facts（兼对照 D-001） |
| 对照 | GOAL-003 `00-meta` 成功标准 1～3 与 C2–C4；本目标 D-001；Root D-002；GOAL-002 D-001 §1/§4/§5 |
| 信息门禁 | I-001 / I-002 / I-009 **verified**（Root D-002，2026-08-25）。I-006 在 Root `00-meta` 仍为 **registered**（产品事实投影，合同 §1 已写入）；无到期且影响本切片的 required 未闭项。I-003/I-007 属策略产品面（R3）；本阶段设密走现行 8–72 基线，与 GOAL-003 边界一致 |
| 共享资料 | 无（`none`；未把资料目录当事实或关闭证据） |
| 代码基准 | `HEAD` `cf296e3c7a11037b8c89e9b1af1dc78eb44e5874` |
| 本审计未改 | 目标 `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码或测试 |

未读取或比较其他工作区上下文。本目标此前无 self / independent 条目（`03-audit.md` 仅占位行）。项目级 `docs/architecture/independent-audit-execution.md` 写「先自审再 independent」；AGENTS P-003 写 `independent` 不固定要求 self。本条按用户指定独立执行，不把缺 self 标为 required。

### 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md` id / root_goal / canonical_scope / `plan_refs`+`primary_plan` 与本目标路径一致；资料目录 `none` |
| C2 · 0056 成对 DDL（SQLite INTEGER / PG BIGINT）+ `ApplyPostgres` | `apps/api/internal/modules/authsession/migration/migration.go` `passwordRecoveryDDL` / `passwordRecoveryPGDDL`；Descriptors v56 `password_recovery_challenges`；`migratePasswordRecovery` + `migratePasswordRecoveryPG`。未改既有迁移 checksum 载体（`r2BaselineDDL` 未塞入新表） |
| C2 · 冻结 checksum 与 DDL 实算一致 | 台账 want = `e19db1a293a013e801abbf60c47be55174dec7f8722fd2a5cd05eb971b4520c3`（`store/migrate_test.go`）。本会话按 `kernel.MigrationChecksum` 规则（normalizeSQL + `"0056:password-recovery-challenges:v1"`）独立 sha256 复算 **一致** |
| C2 · 黄金断言随 catalog 头 = 56 | `identity.go` `completeFingerprintCatalogHead = 56`；`completeLostLedgerTables` / `postV1CatalogTables` 含 `password_recovery_challenges`；`identity_test.go` `lockedHeadExtraTables[56]` + complete-no-ledger 夹具；`migrate_test.go` `len==56`；`operations_test.go` / `restart_test.go` 尾断言 |
| C3 · 合同数值与证明形态 | `recovery.go` 常量 `recoveryCodeTTL=10m` / `recoveryResendCooldown=60s` / `recoveryMaxFailedAttempts=5`；码经既有 `generateEmailCode`（`%06d`）+ `hashCode`（sha256）+ `subtle.ConstantTimeCompare`。只消费 `kernel.MailSender` |
| C3 · start 防枚举 + 两阶段补偿 | 未知/无 verified 邮箱/disabled → `ErrRecoveryNotAvailable` → HTTP **同一** `202 {"status":"dispatched"}` 且不发信。发信失败删挑战。账号定位：用户名精确优先，其次 `lower(email)` 且 `email_status='verified'`。投递目标恒为账号已绑定地址 |
| C3 · complete 形状与 MFA 门序 | `POST /api/auth/recovery/complete` → 204、不签发会话。顺序：body → 定位（未命中统一 `RECOVERY_CODE_INVALID`）→ 评估码 → **匹配后**第二因子 → 8–72 基线 → `CompleteRecovery`。`VerifySecondFactor` 为 `requireActiveSecondFactor` 薄封装；模块关闭 true-nil 接口（`composition.go`）。错码与错误第二因子走 `ConsumeRecoveryAttempt` |
| C3 · 会话语义（I-008 / §4） | `CompleteRecovery` 先 DELETE 挑战（rows≠1 → `ErrRecoveryNotPending`）再 `UpdateUser(PasswordHash, MustChangePassword=false)` → `token_version+1` + refresh 全撤销。`TestCompleteRecoveryConsumesChallengeAndRotatesCredentials` 读回 version 前移、live refresh=0、must_change 清除 |
| C3 · 错误目录 + 运维放行 | catalog 新增 4 码；`error_contract_test.go` 冻结集已含。`operational.go` 将 start/complete 列入维护模式 allowlist，`TestOperationalGateAllowsEveryRecoveryPath` 覆盖 |
| C4 · mock 渠道取码 | `recovery_e2e_test.go`：`StartRecovery` 经真实 `OutboxSink` 出站取 6 位码 → evaluate match → `CompleteRecovery` 换密 → 再 complete 拒绝 |
| C4 · Web 两步流 + i18n | `LoginPage.tsx` 忘记密码 → start 202 进步骤 2 → complete；`RECOVERY_SECOND_FACTOR_REQUIRED` 揭示 TOTP/恢复码字段；过期回步骤 1 保留账号。`auth-client.ts` `recoveryStart`/`recoveryComplete`。zh-CN / en-US `login.recovery.*` 键齐全。本会话 `LoginPage.test.tsx` **19/19 PASS**（含三组恢复用例） |
| 本会话复跑 | `apps/api`：`go test ./internal/store/ ./internal/modules/authsession/ ./internal/handler/ ./internal/modules/mfa/ -count=1` **exit 0**（store 34.8s / authsession 6.4s / handler 39.7s / mfa 8.6s）；`go build ./...` **exit 0**。`TestFullCatalogPostgresBootstrapIntegration` 本会话 **未跑**（需 `PG_TEST_*`；sqlite 全新/升级路径由 store 包覆盖） |

### 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| 成功标准 1：全新库与升级路径干净应用 0056；`go build ./...` 与相关包测试绿 | **达成**（sqlite 路径本会话核对；PG 全 catalog 本会话未复跑） | C2 黄金断言 + store/authsession/handler/mfa 测试绿 + `go build ./...` exit 0 |
| 成功标准 2：start mock 取码 → complete 设密 → 旧会话失效 → 新密码可登录；MFA 账号被第二因子门拦住直至有效因子 | **部分达成** | mock 取码+设密：`recovery_e2e_test.go`。旧会话失效：域测试 token_version/refresh（e2e 未覆盖）。「新密码可登录」无 HTTP 登录回证（e2e 写入的是字面 hash `recovered-hash`，不是可登录 bcrypt）。MFA 门：handler 用 **fakeGate**；无真实 `mfa.Service` 同链。见 F-002 |
| 成功标准 3：independent 开放 required = 0；Root progress → 2/4 | **本条落盘后意见侧未满足** | 开放 required = 1（F-001）。Root 1/4 与检查点不由本审计改写 |
| D-001 §1 迁移同构 | **达成** | 与 0055 同形；双方言时间列；checksum 入台账 |
| D-001 §2 API 形状 / 防枚举 / 错误码 | **部分达成** | start 202 同形、complete 204、4 新码 + 复用码均在。**complete 失败未 `record` 限流**（F-001）。冷却走已有 `EMAIL_RESEND_COOLDOWN`（§6 未列该复用码；Web 已映射，作注释而非 finding） |
| D-001 §3 MFA 门 | **达成**（实现；测试为 fake） | 匹配后、设密前；缺字段 → `RECOVERY_SECOND_FACTOR_REQUIRED`（故意不烧尝试，与 Web 两步一致）；错因子烧尝试。薄封装 + true-nil |
| D-001 §4 会话 / 不签发 | **达成** | UpdateUser 语义；HTTP 204 无 token |
| D-001 §5 审计/通知 | **部分达成** | event=`account.password-change`，detail action=`password-recovery`；通知 `account.password-changed`。detail **携带 `userId` 而非 username**（F-003） |
| GOAL-002 D-001 §1 6 位码 / TTL / 冷却 / 错 5 次 / 无邮箱不自助 / MFA 完成前第二因子 | **达成**（实现） | 常量与测试：cooldown 30s 拒绝 / 61s 换码；第 5 次作废；pending 邮箱不可 resolve；无路径 start 不发信 |
| I-001 / I-002 / I-009 | **无到期未闭 required** | Root D-002 verified；实现投影可核对 |

`00-meta` 已将 C2–C4 标完成、progress 4/5。本意见 **不批准或改写** 这些检查点；C3「限流」完成声明被 F-001 打折。

### Findings

#### F-001 · complete 失败路径未写入 loginRateLimiter（D-001 §2 字面未落地）

- 严重度：med
- 建议：required
- 状态：open
- 描述：D-001 §2 冻结「两端点共用 loginRateLimiter 模型（IP\|identifier 桶；**complete 失败也 record**）」。`recovery.go` start/complete 均 `allow()`，但 `record()` **仅**出现在 start 的 `ErrRecoveryNotAvailable` 分支（`handler/recovery.go` 约 L99–L103）。complete 的错码 / 过期 / MFA 失败 / `INVALID_PASSWORD` 均不 `record`。限流器是恢复 handler **自有实例**（不与登录共用），因此对已知账号 complete 的 429 路径实质上不可达。I-002 的 5 次作废仍是主控制；本条是冻结合同的防御纵深缺口，不是「没有限流函数」。
- 证据：`apps/api/internal/handler/recovery.go`（complete 无 `rateLimiter.record`）；对照 `auth.go` 登录失败必 `record`。handler 测试无限流用例。E-002 / `00-meta` C3 将「loginRateLimiter」记为完成，与字面不符。
- 建议修复：complete 在消耗尝试的失败路径（及未知账号的统一 invalid 响应）调用 `record`；补测试断言 20 次失败后 429 `RATE_LIMITED`。未知账号是否 record 须保持与 start 防枚举同形（不要用 429 当存在性预言机）。

#### F-002 · 成功标准 2 的「新密码可登录」与真实 MFA 服务缺少同链证据

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：C4 检查点（mock 出站取码 + Web 两步流 + i18n）有直接证据。成功标准 2 额外要求旧会话失效、**新密码可登录**、MFA 账号被门拦住直至有效因子。后两项目前分别是：域层写入字面 hash（不能拿去登录）；handler 用 fakeGate 返回 `handler.ErrMFAInvalid`（碰巧与 `mfa` 包 sentinel 别名同一变量，实现可接线，但没有真实 TOTP/恢复码消费的恢复链测试）。不构成 C4 造假，但关门前成功标准 2 不能称为全链可核对。
- 证据：`recovery_e2e_test.go` 使用 `"recovered-hash"`；`recovery_test.go` 会话语义；`recovery_test.go`（handler）`fakeGate`；`mfa` 包无 `VerifySecondFactor` 专项测试（方法为薄封装，覆盖落在 `requireActiveSecondFactor` 既有路径）。

#### F-003 · 审计 detail 未按 D-001 §5 携带 username

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-001 §5：「detail action=`password-recovery`，**携带 username**」。实现 `recordRecoveryEvent` 用 `UserByID` 的 **Name** 填 `ActorName`，detail after 为 `userId`。事件族正确、不需 CHECK 迁移；可检索性弱一档。
- 证据：`handler/recovery.go` `recordRecoveryEvent`；`operationlog.NewDetail("password-recovery", nil, map[string]any{"userId": userID})`。

#### F-004 · D-001「任一失败路径消耗 attempt」与两处实现例外未回写方案

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-001 §2：「任一失败路径消耗挑战 attempt_count（≤5 次作废，含第二因子失败）」。实现中：（1）缺第二因子字段 → `RECOVERY_SECOND_FACTOR_REQUIRED` **不消耗**（handler 测试显式断言；Web 靠此揭示字段）；（2）邮箱码已匹配后 `INVALID_PASSWORD` **不消耗**。第（1）与防 TOTP 爆破意图不冲突（错因子仍消耗），且是已实现 UX 的前提。第（2）在 MFA 恢复码已在 `VerifySecondFactor` 内一次性消费之后发生时，会出现「恢复码烧掉、挑战仍活、密码未改」。应把例外写回 D-001，而不是为迁就字面去烧「缺字段」尝试。
- 证据：`handler/recovery.go` complete 门序；`TestRecoveryCompleteSecondFactorGate` 注释 “NO attempt burned for an empty field”；`TestRecoveryCompletePasswordBaseline` 不调用 consume。

### 必改项汇总

| ID | 严重度 | 摘要 |
|----|--------|------|
| F-001 | med | complete 失败必须按 D-001 §2 `record` 到 IP\|identifier 桶，并补可核对测试 |

F-002～F-004 为 recommended，不阻断本切片放行，但 F-002 应在 C5 关门前决定是补 HTTP 回证还是书面收窄成功标准 2。

### 与既有意见的异同（若有 self/independent 历史）

本目标此前无 A 条目。不与 GOAL-002 A-001（R1 合同 self `pass`）冲突。E-002 将 C2–C4 记为完成：C2 与本审计一致；C3 限流完成声明被 F-001 否定其字面；C4 检查点本身可核对，成功标准 2 全链过宽于 C4 证据。

注释（不升 finding）：goal-tree 文本树把 GOAL-003 画在 GOAL-002 下，状态表 `parent` 仍为 `GOAL-001-iam-recovery`（与 `00-meta` 一致）。E-002 写 i18n「各 19 键」，zh/en 实际为 20 个 `login.recovery.*` 键。`TestStartRecoverySilentOnNoPath` 用标识 `"u-bob"`（用户 id）而非 username `"bob"`，无邮箱已有账号的 start 静默路径靠代码阅读而非该用例命中。

### 结论 + 建议给编排器/用户的下一步

实施切片主体与冻结合同同构（0056、6 位码、TTL/冷却/5 次作废、无邮箱不发信、MFA 门序、UpdateUser 撤会话、mock 取码、登录页两步流）。**不能无条件放行 C5 关门**：F-001 为开放 required。

建议 `/govern`：响应 A-001；对 F-001 落地 `record` + 测试后标 `fixed`；F-002 选择补一条「complete 后用新密码登录」(+ 可选真实 MFA) 的包测试，或书面收窄成功标准 2 的「可登录」为「走与自愿改密同一 UpdateUser/HashPassword 路径」；F-003/F-004 按需改代码或回写 D-001。不要在 required 未闭合时把 GOAL-003 标 `done` 或把 Root 推到 2/4。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。

---

## 编排器响应（2026-08-25 · /govern · git ddd20500）

> 本节为编排器响应留痕，不改动上方独立意见正文。

| Finding | 处置 | 证据 |
|---------|------|------|
| F-001（required） | **fixed** | complete 的消耗型失败路径（错码/无挑战/未知账号/过期/第二因子错误）全部 `rateLimiter.record`；新增 `TestRecoveryCompleteRateLimitedAfterTwentyFailures`（20 次后 429 RATE_LIMITED）与 `TestRecoveryCompleteNonGuessFailuresDoNotRecord`（非猜测失败不记桶）；handler/authsession/mfa 包全绿 |
| F-002（recommended） | **fixed**（补证路径） | ①e2e 改真实 bcrypt：`bcrypt.GenerateFromPassword` 设密 + `CompareHashAndPassword` 双向断言（新密码可验证、错误密码不可）；②新增 mfa 包 `recovery_gate_test.go`：真 `Service.VerifySecondFactor` 走 TOTP 正确/错误 + 恢复码一次性消费——成功标准 2 的 MFA 分支同链可核对。选择补证而非收窄标准 |
| F-003（recommended） | **fixed** | `recordRecoveryEvent` detail 改为携带 `username`（D-001 §5 原口径），ActorName 仍为 Name |
| F-004（recommended） | **fixed（回写方案）** | D-001 §2 已回写两条例外（缺第二因子字段不消耗；匹配后 INVALID_PASSWORD 不消耗；恢复码消费不回滚的可接受损耗注明） |

注释三项（goal-tree 树形 vs 表格 parent、i18n 键数笔误、SilentOnNoPath 用例标识）核对属实：树形为展示层级、表格 parent 为权威且一致，不改；E 文本不再修正计数（以本节为准）；用例标识已顺带在后续无行为差异，保留。

开放 required = 0（F-001 fixed）。
