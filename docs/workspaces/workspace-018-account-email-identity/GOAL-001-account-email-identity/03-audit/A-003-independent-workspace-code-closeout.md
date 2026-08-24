---
id: A-003
doc: audit-entry
goal: GOAL-001-account-email-identity
source: independent
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
scope: 工作区 018 整体完成（GOAL-001～005）代码层独立审计——迁移 0054/0055、服务层状态机、HTTP 面、Web 绑定卡、测试复跑、双端构建、git 基准
verdict: pass
auditor: 独立审计会话（dsh harness · 代码为判据，非治理文档）
---

# A-003 · 工作区 018 整体完成（代码层 closes-out·independent · 2026-08-25）

- **source**：independent
- **auditor**：独立审计会话（用户指令：独立审计代码，不以治理文档作为判据）
- **类型** / **scope**：close-out · 工作区 018（GOAL-001 Root 及 GOAL-002～005 子目标）**产品代码**验收：迁移 0054/0055 与 checksum、绑定/校验/重发状态机、HTTP 三端点与 7 错误码、I-006 代填、Web 最小绑定卡、golden 断言、测试复跑与双端构建、git 基准与台账一致性的代码侧核对
- **verdict**：**pass**（代码层无未闭合 high/med required；1 条 recommended + 2 条 note）
- **开放 required**：0

### 范围与区间

| 项 | 值 |
|----|------|
| 被审范围 | workspace-018 全部五个目标的产品代码面（GOAL-002 合同冻结 → GOAL-003 双方言 schema → GOAL-004 绑定流 → GOAL-005 证据） |
| 审查基准 | 工作树干净；`HEAD = 76bd21aa`（Root 关门提交）；产品代码与 A-002 记录的 `6c6496d4` 基准相关 commit 链（0cbe3242 / 0ae17f09 / bd1cdff9 / 6c6496d4）可核对 |
| 判据 | 仅代码事实：文件存在、语义与契约注释一致、测试跑绿、构建跑通；治理文档仅用于提取「应交付什么」，不以其表述为通过理由 |
| 未改 | 目标 `status` / `progress` / goal-tree / 方案正文 / 产品代码 / 测试 |

### 成果（全部经本会话代码走读 + 复跑，逐条有证据）

| 主张（文档） | 代码证据 | 复跑 |
|------|------|------|
| 迁移 0054 account_email_identity（email TEXT NULL / email_status CHECK('pending','verified') / lower(email) 唯一索引；双方言可移植） | `authsession/migration/migration.go` L159-163 DDL + L411-419 注册（ApplyPostgres nil 理由注释）；checksum `f9a0bc65…2b0b` 入 `store/migrate_test.go` L699-701 冻结目录 | `TestMigrate0054UpgradePathLandsUnboundRows` / `TestMigrate0054FreshCatalogSemantics` PASS（升级路径 NULL 落地、大小写唯一拒绝、多 NULL 共存、CHECK 拒绝） |
| 迁移 0055 email_verification_challenges（成对方言 INTEGER/BIGINT） | `migration.go` L170-188 成对 DDL + L420-427 注册；checksum `1556bda2…b754` 冻结目录 L702-704；对 PG 走独立 ApplyPostgres | `TestCompiledMigrationCatalogOwnership` PASS（55 条目录逐项 identity+checksum+Apply） |
| golden 断言随动（head 55、lockedHeadExtraTables 54/55） | `store/identity.go` L93 `completeFingerprintCatalogHead = 55`；`identity_test.go` L98-107 含 54/55 条目 | `TestCompleteFingerprintTracksCatalogHead` PASS |
| 服务层：Bind 占槽 fail-closed / Verify 常量时间 + 失败计数 / Resend 冷却；TTL 10m、冷却 60s、5 次作废；换绑=覆写释放旧址；同址 verified 幂等；两阶段派发+补偿 | `email_identity.go` 全文：常量 L32-35；BindEmail L141-218（两阶段 + `compensateBind`）；VerifyEmail L240-269（`subtle.ConstantTimeCompare` L314、独立事务失败记账、5 次作废）；ResendEmailCode L385-457（快照补偿）；换绑/幂等/冷却 L161-177 | `go test ./internal/modules/authsession/ -count=1` **ok**（含 `TestSendFailureRollsBindBack`、`TestRebindOverwriteReleasesOldSlot`、`TestBindSamePendingAddressHonorsCooldown` 等 10 项服务用例） |
| R4 e2e：经真实 mock 渠道适配器（OutboxSink→mail_outbox 迁移 0051）取码闭环，非测试桩 | `email_identity_e2e_test.go`：`mail.NewOutboxSink(st, 0)` 真实适配器 L62；`codeFromOutbox` 从 `sink.List/Get` 记录 body 提取 6 位码 L23-45；全链路 bind→他号 EMAIL_TAKEN→通道取码 verify→bob 全程 unbound | 随 authsession 包 **ok**（5.574s） |
| HTTP：POST /api/account/email/bind\|verify\|resend（identity-only，无权限键） | `handler/email_identity.go` 三端点 + 错误映射（EMAIL_INVALID 400 / EMAIL_TAKEN 409 / EMAIL_NOT_PENDING 409 / EMAIL_CODE_INVALID 400 / EMAIL_CODE_EXPIRED 400 / EMAIL_RESEND_COOLDOWN 429 / EMAIL_SEND_FAILED 502）；`account/provider.go` L77 挂载 + Descriptor L53；`kernel/profile.go` L177 admin.account v2.0.0 路由清单含三端点 | — |
| 7 新错误码入 Catalog 并冻结 | `errorcatalog/errorcatalog.go` L192-199 七个 EMAIL_* 码（zh/en 双语）；`handler/error_contract_test.go` L112-115 冻结清单七码齐 | — |
| I-006 管理员代填：PATCH 非空→pending、""→unbound、非字符串 400、冲突 409 EMAIL_TAKEN、无 verified 旁路 | `handler/users.go` L70 RawStringFields `["password","email"]` + L220-226 类型守卫；`users_repository.go` L194-225 代填语义（lower(email) 冲突检查、清空删挑战）；`handler/account_self.go` L121-127 profile 读回 email/emailStatus | `TestUsersPatchEmailPrefillFlows` **PASS**（200 pending → 200 清空 → 400 非字符串 → 409 EMAIL_TAKEN） |
| Web 最小绑定卡：自注册 custom 组件、account tab 挂载、i18n zh/en | `web/src/components/email-identity.tsx`（unbound/pending/verified 三态 + bind/verify/resend 表单，`registerCustomComponent("email-identity")`）；`apps/api/internal/modules/account/schema/account.json` L116-120 profile tab custom 节点；`main.tsx` L16 导入；i18n en/zh 各 12+ 键（含 errCooldown「Too many requests」） | `vitest run src/components/email-identity.test.tsx` **3 passed** |
| 双端构建 | — | `go build ./...` **OK**；`npm run build`（tsc -b && vite build）**成功** |

### 对照成功标准（代码侧）

1. users 可持可空邮箱与可核对状态；无邮箱账号仍可登录 —— e2e bob 全程 (NULL,NULL) 且 `UserByID` 正常（e2e L92-99）；登录面无该列依赖（A-002 已核）✓
2. 绑定/校验流落地且校验信经 `kernel.MailSender`；无生产渠道时能从 017 默认渠道取信 —— kernel/mail.go 端口契约（L67-75）+ 真实 OutboxSink e2e ✓
3. 唯一性 fail-closed + 换绑同合同 —— 服务层 COUNT + 唯一索引双保险；换绑测试 PASS ✓
4. 未越界 —— 本波 commit 面（0cbe3242/0ae17f09/bd1cdff9/6c6496d4）无恢复/邀请/SMS/模板产品面 ✓

`progress` 百分比不作为本条闭合依据（P-001）。

### Findings

#### F-001 · 未绑定账号调用 verify 返回 500 INTERNAL 而非受控 EMAIL_NOT_PENDING

- 严重度：low
- 建议：**recommended**（不阻断关门）
- 描述：`evaluateVerification`（email_identity.go L287-294）对**完全未绑定**账号（`email IS NULL`）的查询 `SELECT email_status FROM users WHERE id = ? AND email IS NOT NULL` 返回无行 → `Scan` 错误（kernel.ErrNoRows）经 `!errors.Is(err, errNotSentinel)` 分支 L338-341 包装为硬错误 → `VerifyEmail` 返回非领域错误 → handler `writeDomainError` 走 default → **500 INTERNAL**。对照：已绑定 pending 但无挑战行（如 I-006 代填未 resend）走 errNotSentinel → 受控 `EMAIL_NOT_PENDING` 409。同一「当前无可校验验证码」情形两种 HTTP 语义。UI 在 unbound 态不渲染 verify 表单，实际用户面不可达；但 API 契约上「无待校验邮箱」应以 EMAIL_NOT_PENDING 报之。修复建议：L290-294 将无行错误判为 `verificationNotPending`（返回 errNotSentinel 或直接受控），并补一个 unbound-verify 的 HTTP/服务用例。
- 证据：`authsession/email_identity.go` L287-294 + L338-341；`handler/email_identity.go` L122-140（default→INTERNAL）；对照 `email_identity_test.go` `TestAdminPrefillPendingCannotVerifyDirectly`（L222 断言 ErrEmailNotPending 达成的路径是挑战行缺失而非无邮箱）。

#### N-1 · verify 成功路径 rowsAffected 守卫与失败计数双事务（best-effort，已文档化）

- 严重度：low（note）
- 描述：`evaluateVerification` L325-328 用 `RowsAffected()!=1` 兜并发消费；`registerFailedAttempt` 走两个独立事务、读回 best-effort——并发下阈值作废可能延迟一个轮次，最终由下次 evaluate 收敛。均与注释一致，无未兑现承诺。
- 证据：email_identity.go L318-331、L344-375。

#### N-2 · 两阶段派发进程崩溃窗口（与既有 N-2 补偿失败分支互补）

- 严重度：low（note）
- 描述：phase-1 提交挑战后、phase-2 投递前进程崩溃 → 落库 pending 且挑战在册但无邮件出站。恢复路径完整：用户 resend（60s 冷却后）即补投；管理员可清空回 unbound；无死锁态。属两阶段语义固有窗口，既有 N-2（补偿自身失败）同族，本期不升 required。
- 证据：email_identity.go L197-217（挑战先落库）、L385-457（resend 补投路径）。

### 必改项汇总

无 required。F-001（recommended/low）供 /govern 择期响应；N-1/N-2 为 note 维持。

### 与既有意见的异同

- 与 A-002（independent conditional）**同向**：产品事实一致；A-002 的 F-001（YAML 漂移）已 fixed——本会话复核 GOAL-001～005 `00-meta` YAML 现全部 `status: done`、`progress` 与正文/goal-tree 同源（001:4/4、002:3/3、003:4/4、004:4/4、005:3/3），A-002 的核对结论（五判据 1-4、门禁闭环、边界）在代码侧即 A-003 成果表。
- 与 A-001（self pass）同向。
- 新增点：F-001 为本次代码走读新发现的低危 API 语义瑕疵（HTTP 契约缺口），非台账问题。

### 结论 + 建议给编排器/用户的下一步

**pass** —— 工作区 018 声称交付的产品代码全部存在、与契约注释一致、测试与构建复跑全绿。治理台账（YAML/goal-tree/进度）与代码事实现一致，Root `done 4/4` 的代码前提成立。残余仅 3 条低危 note/recommended，不构成放行或关门障碍。

建议 /govern：

1. F-001 按 recommended 处理：可顺手修正（一行受控路径 + 一个用例），或登记为后续 IAM 波次承接项；无需 `accepted-residual` 强制路径。
2. N-1/N-2 维持 known-boundary 台账即可。
3. 本审计不改 status/progress；Root/VP-018 的关门状态已由既有流程落定，无需重开。

### 声明

本意见 `source: independent`，**不修改** `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码。响应、finding 闭合与关门由 **`/govern`** 处理。