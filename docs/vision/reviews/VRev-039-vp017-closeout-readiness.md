---
doc_type: vision-review
id: VRev-039
status: active
source: self
created: 2026-08-24
updated: 2026-08-24
version: 0.2.0
parent: null
---

# VRev-039 · VP-017 关门就绪度审视（2026-08-24）

| 字段 | 值 |
|------|-----|
| source | self |
| auditor | `/vision` · 会话编排（grok-4.6）；本轮对代码 / 测试 **独立复验**，不以 Goal 台账为关门充分条件 |
| scope | `VP-017-outbound-mail` 组合层关门就绪 · 退出判据 1–6 · 代码实现 · 本会话测试 · Vision required · 有界 residual · 组合索引 |
| audit_type | vision-plan（关门就绪度 · 代码成果独立核验） |
| verdict | pass |
| 建议 class | editorial（组合层关门 + 索引同步 + residual 点名；不改 Charter 方向） |
| open required | 0 |

## 范围与结论

只读核对：`docs/architecture/principles.md` P-006、`docs/vision/alignment.md` §6/§7/§9、`charter.md` `@0.2.0`、[VP-017-outbound-mail](../plans/VP-017-outbound-mail.md)（审视时 `active` v0.2.0）、`roadmap.md`、`workspaces.md`、`revisions.md`（至 VR-041）、`reviews.md` 与 `reviews/VRev-001`～`VRev-038`、lead 工作区 `workspace-017-outbound-mail`（绑定与 Root 状态只作**指针**，不替代代码核验）。

**本轮独立核验（非转录治理记录）**：

1. **源码**：`apps/api/internal/kernel/mail.go`（`MailMessage` 单收件人纯文本、无 From 字段、`MailSender.Send`）；`internal/mail/capture.go`（容量 1 + `Last()`/`Reset()` + slog）；`internal/mail/smtp.go`（隐式 TLS 465、证书校验恒开、仅 AUTH PLAIN over TLS、对端未广告 AUTH 则 fail-closed、`Ping` 不进端口）；`internal/config` `MailSMTPConfigured` / `validateMail`（未触摸 = 合法默认；显式则四键必填）；`internal/composition` `newMailSender`（未配置 → capture + probe `nil`；显式 → SMTP + `Ping` 进 `readyz`）。`net/smtp` 仅出现在 `internal/mail/smtp.go`。
2. **测试（本会话复跑，`apps/api`，go1.26 windows/amd64）**：`go test ./internal/kernel -run Mail -count=1` ok（0.685s）；`./internal/mail -count=1` ok（0.816s，含 loopback TLS harness；`TestSMTPLiveDelivery` 因未设 `MAIL_SMTP_TEST_*` skip）；`./internal/config -run Mail|ValidateProd -count=1` ok（0.702s）；`./internal/composition -run Mail|SMTP|Capture|Ready -count=1` ok（2.911s）。`go vet ./internal/kernel ./internal/mail ./internal/config ./internal/composition` 0 finding。
3. **live**：本会话**未**对真实 SMTP 对端投递（无 `MAIL_SMTP_TEST_*`）。退出 3 允许「与生产合同等价的 harness」；本轮以 `smtp_test.go` 隐式 TLS loopback 全会话为准。不把 Goal N-001 叙述当作本轮 live。

未把 Goal `progress=4/4` 或 Root A-001/A-002 正文当作退出判据的充分证据。治理记录仅用于定位实现与对照开放 required。

**总判：pass（0 open required · 1 open recommended）。**

**关门的实质证据已齐备**（代码 + 本会话测试 + 等价 harness），可按 alignment §7 做**有界 closed**。Vision Review open required = 0；对齐链成立；激活后 Charter 仅 editorial（VR-040/VR-041），无 strategic 宽阻断。组合索引仍写 VP-017 `active`，VP 信息表 I-017-001～005 仍 `open`、I-017-006 仍 `registered`（实现层 F-003 已 delegated `/vision`）——这是待用户书面确认后的投影同步，**不是**实现缺口。本轮用户意图为「按 1→4 顺序执行：先有界关闭 VP-017」。

本意见原文**不**把组合索引改写成 `closed`。

### 核对事实

| 核对项 | 结论 | 证据（本会话独立核验，除非标明指针） |
|--------|------|------|
| 单愿景 / `vision_ref` | **pass** | 唯一 active Charter `schema-ui-core-admin-foundation@0.2.0`；VP-017 `vision_ref` 精确匹配 |
| 工作区绑定 | **pass** | `workspace-017` 唯一 lead / delivery；`plan_refs` / `primary_plan` / `vision_role: delivery` 合规；Charter `primary_workspace` 仍为 workspace-001 |
| 区证据指针（§7.1） | **pass（指针）** | goal-tree Root `done 4/4`、GOAL-002～005 `done`。**不**据此放行；放行依据是下表退出 1–5 的代码/测试 |
| 实现层开放 required | **pass（指针）** | Root A-002 F-001/F-002 → `fixed`；F-003 delegated `/vision`；开放 required = 0。本轮未改 Goal finding |
| 退出 1 · 内核发送端口；公共面无 SMTP 客户端类型 | **pass** | `kernel.MailSender` / `MailMessage`；`rg net/smtp` 仅 `internal/mail/smtp.go` |
| 退出 2 · 未配置 SMTP 时默认仍能开发与快测 | **pass** | YAML `mail.smtp.*` 缺省空；`newMailSender` → `CaptureSink` + probe `nil`；composition 本轮绿 |
| 退出 3 · 显式 SMTP 可核对至少一封投递；不完整 fail-closed | **pass** | 本会话 `./internal/mail` 全绿（loopback TLS harness）；`validateMail` / `NewSMTP` 缺键 fail-closed。真实对端 live 未跑 |
| 退出 4 · 仅显式配置后 `readyz` 扩依赖 | **pass** | `newMailSender`：未配置 probe=`nil`；显式 probe=`sender.Ping` |
| 退出 5 · 未进 SMS / 账号 email / 邀请 / 恢复 / 模板 / 业务域；未改 Charter | **pass** | `users` 基线 DDL 仍无 email 列；Charter 仍 `@0.2.0`；无第二邮件方言 |
| 退出 6 · required = 0 | **pass** | Vision Review 与实现层开放 required 均为 0 |
| Vision required（§6 门禁 8） | **pass** | `reviews.md` open required = 0；VRev-037/038 为激活审视，本条为关门就绪首份 |
| Charter strategic 后 re-align | **pass** | 激活后仅 VR-040/VR-041（editorial）；无宽阻断 |
| 组合索引当前陈述 | **pass（待同步）** | Charter 关系节 / `roadmap.md` 第 17 行与 RT-M01「registered」/ `reviews.md` 摘要 / 区 `workspace.md` 仍写 VP-017 `active`；VP 信息表滞后于 Root |

## Findings

#### V-F072 · 组合层关门须同步索引，并显式映射 exit 1–6 ↔ 本轮独立证据、点名有界 residual

- level: `recommended`
- status: `open`
- severity: low
- impact: alignment §7.2 允许有界 closed，但 residual 必须点名到具体 workspace / goal id。若只让 Root `done` 而组合索引仍称 `active`、且 VP 信息表仍 `open`，后续读者会把 A6 读成未交付。
- finding: |
  1. 用户确认组合层关门时一次写清 exit 1–6 ↔ **本轮独立**证据（源码路径、本会话 `go test`/`go vet`、loopback TLS harness；治理 A 条目仅作指针）。
  2. residual 至少点名：`workspace-017` / `GOAL-001-outbound-mail` / **env-gated live 未实跑**（`TestSMTPLiveDelivery` 无 `MAIL_SMTP_TEST_*`；离线 harness 已满足「与生产合同等价」）。不把「本波无 handler 消费端口」写成缺口——退出分母消费者 = 测试/harness。
  3. 同步 `roadmap.md`（VP 行 + **RT-M01 → delivered**）/ `workspaces.md` / Charter 关系节 / `reviews.md` 摘要 / 区 `workspace.md`：VP-017 → `closed`。将 VP 信息表 I-017-001～006 改为 verified（对齐 Root I-001～I-006）。Root `done` 不能冒充 VP 层用户确认。
- evidence:
  - 本会话测试：kernel / mail / config / composition 指定套件全绿；vet 0
  - `git diff --name-only 250cb9c..53b64a5` 产品面含 kernel/mail/config/composition + YAML；`apps/web` 无邮件产品页
  - Root A-002 F-003：VP 信息表滞后，收 VP 时回写
  - alignment.md §7.2
- closure: |
  `/vision` 在用户书面确认组合层关门时按上列三项一并完成。本 finding 不阻断「就绪」结论，只约束关门落盘形状。
- 建议 class: `editorial`

### 不构成 fail / 不新开 required 的诚实边界

1. 本 `pass` **不是**「组合索引已 closed」：用户书面确认与索引原子同步仍待发生（本轮用户已给出「按 1234 顺序全部执行」，步骤 1 = 关 VP-017）。
2. `/vision` 本入口只写 `source: self`。用户要求独立符合代码实现：本报告用源码+测试兑现；不是 `/vision-audit` 的 independent VRev。无独立 Vision Review 不是 alignment 强制项（强制时机仅为 Charter 初建与 strategic）。
3. 真实 SMTP 对端未投递是取证边界，不是「显式路径未实现」。loopback TLS harness 覆盖拨号形态（隐式 TLS、AUTH、MAIL/RCPT/DATA）。
4. `composition` 当前 `_ = mailSender`：本波冻结的消费者是测试/harness，不是账号邮箱。后续 Admin VP 才接线。
5. 不把 progress=`4/4` 当作关门权威。
6. 架构 A3、SMS、账号 email、邀请、自助恢复本就不在退出分母。

### 声明

本意见不修改 Charter / VP / Goal status 或 progress；required/recommended finding 的响应由 `/vision` 追加在本报告中；实现层执行仍交 `/govern`。原 verdict 与 finding 原文不得改写。本入口不关闭 Goal finding。

### 门禁含义

- Vision Review **open required = 0**。
- **允许**：用户确认后，`/vision` 按 V-F072 执行 VP-017 有界组合层关门与索引同步。
- **禁止**：在无用户书面确认时把组合索引改成 VP-017 `closed`；把 Root `done` 或 Goal 审计原文冒充本轮代码核验；把账号 email / 自助恢复 / SMS 写成已交付。

### 响应（对 self 意见 · VRev-039 findings 闭合 · 2026-08-24）

| date | actor | summary |
|------|-------|---------|
| 2026-08-24 | `/vision` · 用户书面「按照你的建议的 1234 顺序，帮我全部执行」（步骤 1 = 有界关闭 VP-017） | **不回溯改写**原 verdict `pass` 与 finding 正文。**V-F072 → `fixed`**：VP-017 组合层确认 **有界 `closed`**（架构 A6）。关门记录含 exit 1–6 ↔ 本轮独立证据；residual 点名 `workspace-017` / `GOAL-001` / env-gated live 未实跑。`roadmap.md` / `workspaces.md` / Charter 关系节 / `reviews.md` / 区 `workspace.md` 原子同步（VR-042）。I-017-001～006 → verified。本 scope **0 open required、0 open recommended**。 |
