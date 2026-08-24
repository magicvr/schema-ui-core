---
doc_type: vision-review
id: VRev-040
status: active
source: self
created: 2026-08-24
updated: 2026-08-24
version: 0.2.0
parent: null
---

# VRev-040 · VP-018 意图完备 / 可行性 / 激活就绪（2026-08-24）

| 字段 | 值 |
|------|-----|
| source | self |
| auditor | `/vision` · 会话编排（grok-4.6） |
| scope | `VP-018-account-email-identity` 意图完备、Charter 对齐、退出分母、与 roadmap Admin 功能消费链一致性、激活与开区就绪、VP-008 `go` 消费前新鲜度（Admin 类） |
| audit_type | vision-plan（意图 / 激活就绪） |
| verdict | pass |
| 建议 class | editorial |
| open required | 0 |

## 范围与结论

只读核对：`docs/architecture/principles.md` P-006、`docs/vision/alignment.md`、Charter `@0.2.0`、[VP-018-account-email-identity](../plans/VP-018-account-email-identity.md)（审视起点为同日 planned 草案）、[VP-017-outbound-mail](../plans/VP-017-outbound-mail.md)（本轮步骤 1 有界 `closed`）、roadmap Admin 功能消费链、VR-042、`apps/api/internal/modules/authsession/migration/migration.go`（`users` 基线无 email 列）、`kernel.MailSender`（VP-017 已交付）。本报告落盘时用户已书面要求执行步骤 2–4：落盘 VP → freshness → 激活并开区。

**总判：pass（0 open required）。** 单愿景与 `vision_ref` 精确匹配；新 VP 承接 Admin 功能「账号邮箱身份」的结构选型合法；退出分母与用户已确认消费链同构（出站邮件 → 邮箱身份 → IAM）；方向足以激活并开新 delivery 工作区。两条 recommended 约束 Root 纲领/信息项与 freshness 留痕，不阻断激活。

## 核对事实

| 核对项 | 结论 | 证据 |
|--------|------|------|
| 单愿景 | **pass** | 唯一 active Charter `schema-ui-core-admin-foundation@0.2.0` |
| VP→Charter 机读 | **pass** | `vision_ref` 精确匹配 |
| 语义对齐 | **pass** | 可 fork 的 Admin 基架账号能力；不把业务领域当成功条件；不改 Charter 非目标 |
| 最小完备 | **pass** | 意图、冻结表、非目标、退出 1–5、邻接 VP、I-018-001～006、工作区表、短史均在 |
| 结构选型 | **pass** | 同愿景新纲领波次 → 新 VP；不重开 VP-017；不改 Charter；新 delivery 区 |
| 与消费链 | **pass** | roadmap：A6 出站邮件 → 账号邮箱身份 → IAM。VP-017 本轮已 `closed`，运输面可消费 |
| 退出分母有界 | **pass** | 明确排除忘记密码状态机、邀请、密码策略、SMS、模板、Profile 改默认集 |
| 可行性 | **pass（工作量中、边界清）** | `users` 无 email 列，接入点已知（`core.auth-session` 已在 mvp/admin/demo）。运输端口已存在。工作是补列+状态机+消费 `MailSender`，不是换认证叙事 |
| 开放 VRev required | **pass** | 本报告前 open required = 0（VRev-039 V-F072 随关门闭合） |
| 过早交付主张 | **无** | 激活 ≠ 邮箱列已迁移 |

## VP-008 `go` 消费前新鲜度（Admin 类 · V-F074）

VP-008 正文强制 freshness 的对象是**后续业务 VP**。VP-018 是 Admin 功能身份面，不是 Tier D（Catalog/订单/CMS）。按激活门闩做 **Admin 类**复核：核对 `go` 是否仍可消费、共同门禁有无未恢复暂挂、本 VP 意图是否改 Profile/模块矩阵/Manifest/协议 pin。

| 项 | 结论 |
|----|------|
| 原 `go` 候选 | `ed99e88`（2026-08-10，clean）；解锁 scope = 标准业务模块框架能力，不是账号邮箱身份本身 |
| 现行 HEAD | `092bf37`（`docs(goal-036): 记录 A-004 建议2 浏览器端 sqlite e2e 双 profile 复跑事实 E-008`） |
| 工作树 | **clean** |
| 比对区间 | `ed99e88` → `092bf37`：含 VP-011 模块交付、VP-012～017 架构/横切、VP-009 W7–W11、VP-010 W1–W25。源代码已变；消费模式 = 各波次暂挂后**恢复** `go`，不是仍钉在 `ed99e88` 字节 |
| VP-009 | W1–W4、W6–W11 **done**；W11 D-004 **恢复** VP-008 `go` 宣称；无现行暂挂 |
| VP-010 | W1–W25（含 GOAL-037）**done**；无现行 `go` 暂挂宣称 |
| Vision open required | 0 |
| F-007 residual | 上传授权深度仍 **deferred**（owner=VP-008 lead）。本 VP 不得借邮箱面扩张上传授权 scope |
| 本 VP 是否改 Profile / 模块矩阵 / Manifest / 协议 pin | **意图否**。归属既有 `core.auth-session`。实施若新增模块或改默认集，按消费有效性暂挂 |
| 迁移台账 | **意图会变**（`users` 加列）。这是本 VP 分母内工作，不是激活前已存在的失效；R2 须双方言 checksum 台账 |
| 复核结果 | **PASS（Admin 激活）**。不消费 Tier D 业务解锁；不暂挂 `go` |

`consumer_vp` = VP-018；`last_freshness_review_at` = 2026-08-24；`next_freshness_review_trigger` = 若实施改变共同门禁 / Profile / 模块矩阵 / Manifest / 协议 pin 则重做。

## Findings

### V-F073 · recommended · Root 须写出纲领阶段并登记 I-00N

- level: `recommended`
- status: `open`
- severity: medium
- impact: 若不开区就改 `users` DDL，或把唯一性、验证码 vs 链接、过期冷却混成未登记未知，身份面会在 R1/R3 才爆。
- finding: |
  激活后 lead Root **必须**：
  1. 写出串行纲领（身份合同冻结 → 双方言 schema + 唯一性 → 绑定/校验消费 `MailSender` → 证据：capture + 唯一性 + 无 IAM 恢复）。
  2. 把 I-018-001/002 登记为 required（最晚 R1）；I-018-003/004 投影 VP 已冻结决策；I-018-005 required（最晚 R3）；I-018-006 non-blocking。
  3. 不得把验证码 vs 链接在 VP 层假装已选。
- evidence: VP-018 退出 1–5；I-018-001～006；`users` 基线无 email。
- closure: |
  Root `00-meta` 含 P-001 阶段表 + 上述 I-00N。不要求本 Review 落盘时已经有唯一性/投递形态答案。
- 建议 class: `editorial`

### V-F074 · recommended · 激活记录须留下 Admin 类 freshness 结论，避免被读成已消费 Tier D `go`

- level: `recommended`
- status: `open`
- severity: low
- impact: 若激活记录只写「开区」而不点名：本 VP 非业务域、不改 Profile 意图、F-007 不升格、W11 已恢复 `go`、HEAD `092bf37`，后续读者会把邮箱身份误读成 Catalog/订单解锁。
- finding: 激活时在 VP 短史或 lead Root D-001 写入上表复核结论与候选/HEAD 指针。
- evidence: VP-008 §`go` 消费有效性（业务 VP）；VP-018 激活门闩；GOAL-007 D-001 原候选 `ed99e88`；workspace-009 GOAL-011 D-004 最近一次恢复。
- close requirement: D-001 或 VP 激活短史含 freshness 表；不要求重开 VP-008。
- 建议 class: `editorial`

### 不构成 fail / 不新开 required 的诚实边界

1. I-018-001/002/005 仍 collecting。唯一性占用、验证码 vs 链接、过期冷却不是激活阻断；最晚 R1/R3。
2. 本 `pass` 允许激活与开区，**不是** R1 合同已冻结，也不是可以开始无设计地改 `users` DDL。
3. 无独立 Vision Review 不是 alignment 强制项。若用户要求交叉审视，另走 `/vision-audit`。
4. 用户本轮确认开区但未另写 slug。不把惯例 slug 写成用户已口头点名；开区记录须留痕推导。
5. 架构 A3、IAM 恢复、邀请、密码策略本就不在退出分母。

### 声明

本意见不直接修改 Charter / VP / Goal status。required finding 的响应由 `/vision` 追加在本报告中；原 verdict 与 finding 原文不得改写。

### 门禁含义

- Vision Review **open required = 0**。
- **允许**：用户确认后激活 VP-018、开新 delivery 工作区、按 V-F073/V-F074 写 Root 纲领 / freshness / I-00N。
- **禁止**：把本 `pass` 写成邮箱列已迁移或自助恢复已交付；重开 workspace-017；把 VP-018 当 Tier D 业务 VP 消费订单/CMS 解锁 scope。

### 响应（对 self 意见 · VRev-040 findings 闭合 · 2026-08-24）

| date | actor | summary |
|------|-------|---------|
| 2026-08-24 | `/vision` · 用户书面「按照你的建议的 1234 顺序，帮我全部执行」（步骤 2–4 = 落盘 VP-018 → freshness → 激活开区） | **不回溯改写**原 verdict `pass` 与 finding 正文。**V-F073 → `fixed`**：Root `GOAL-001-account-email-identity` P-001 R1–R4 + I-001～I-006。**V-F074 → `fixed`**：D-001 Admin 类 freshness 表（`ed99e88` → `092bf37`；W11 已恢复 `go`；非 Tier D）。VP `active`；lead `workspace-018-account-email-identity`。本 scope **0 open required、0 open recommended**。 |
