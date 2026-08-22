---
id: GOAL-001-outbound-mail
doc: audit-entry
record_id: A-002
source: independent
goal: GOAL-001-outbound-mail
status: recorded
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## A-002 · Root 关门独立审计（R1～R4 全阶段 + 方向级退出判据）

> 代贴说明：本条由编排器自本地 grok build（grok-4.6 · reasoning high）headless 会话原样誊入（原始输出：[attachments/audit-independent-raw.md](../attachments/audit-independent-raw.md)）；意见只出，不改 status。

- **source**: independent
- **auditor**: grok-build (grok-4.6 · reasoning high)
- **日期**: 2026-08-22
- **scope**: `[workspace-017-outbound-mail] GOAL-001-outbound-mail` Root 关门——VP-017 六条方向级退出判据、R1～R4 全阶段、四个子目标五件套与 A-001 self、信息门禁 I-001～I-006 / I-017-001～006、代码合同与接线、`internal/store` 环境失败是否阻塞
- **verdict**: pass

### 核对面与结论

| 核对点 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 合格：`id` / `root_goal` / `canonical_scope` / `plan_refs`+`primary_plan`=VP-017 / `vision_role=delivery` / `shared_materials_catalog: none`；Charter primary 仍为 workspace-001 | `workspace.md`；`docs/vision/charter.md` |
| 治理结构 | 合格：Root=`GOAL-001` `parent: null`；子目标 002～005 单调不复用；`id`=文件夹名；五件套 + 三个 ledger 目录齐全；路线图 4/4 与 `progress` 一致；goal-tree 树/表同步 | 各 `00-meta.md`；`goal-tree.md` |
| I-00N / P-005 | **无到期未关闭 required**。Root I-001～I-006 均 **verified**（D-002～D-005）。映射：I-001↔I-017-003、I-002↔I-017-004、I-003↔I-017-001、I-004↔I-017-002、I-005↔I-017-005、I-006↔I-017-006 | Root `00-meta.md` / `01-decision.md`；GOAL-002～005 D-001/D-002 |
| VP-017 判据 1 · 端口落地、公共面无 SMTP 类型 | **成立**。`kernel.MailSender`/`MailMessage` 为唯一发送合同（单收件人 `To`、纯文本、From 不进消息体）；`net/smtp` / `tls.Dialer` **仅** `internal/mail/smtp.go`；handler/modules 对 mail 类型零引用 | `apps/api/internal/kernel/mail.go`；本轮 rg；GOAL-004 E-001 |
| VP-017 判据 2 · 未配置仍能启动；capture 可取最后一封 | **成立**。未配置 → `CaptureSink` 容量 1 + slog 双写 + `Last()`；`newMailSender` 缺省 `probe=nil`；composition 双 profile lifecycle 测试在本轮复跑包内为绿 | `internal/mail/capture.go`；`composition.go` `newMailSender`；`composition_mail_test.go` |
| VP-017 判据 3 · 显式路径可核对投递；不完整 fail-closed | **成立（等价 harness，非 live）**。loopback TLS 断言 AUTH 身份 / MAIL FROM / RCPT TO / DATA；明文端点 Send/Ping 必败；部分块 `validateMail` + 构造器双保险拒收。live 测试 `MAIL_SMTP_TEST_*` gated，本轮未实跑——VP 原文允许「live **或** 与生产合同等价的 harness」 | `smtp.go` / `smtp_test.go` / `smtp_live_test.go`；`config.go` `validateMail`；`config_mail_test.go` |
| VP-017 判据 4 · 仅显式配置后 readyz 扩依赖 | **成立**。显式 → `sender.Ping` 经 `RegisterWithMFAProbes` 变参进入；nil probe 被忽略；Ping 与 Send 共用 `tlsConfig()`（隐式 TLS、MinVersion 1.2、校验恒开、无 `InsecureSkipVerify`） | `composition.go:375-380,691-705`；`handler/health.go`；`smtp.go` `Ping`/`tlsConfig` |
| VP-017 判据 5 · 非目标未越界 | **成立**。`git diff --name-only be23164..HEAD -- apps/api` 仅 mail/config/composition/README/configs 文件；无 users email 列、邀请、自助恢复、模板、SMS、第二拨号、Charter 改写 | 本轮 `git diff --stat be23164..HEAD -- apps/api` |
| VP-017 判据 6 · 开放 required = 0 | **成立**。四子目标 A-001 均 self `pass`；GOAL-002/003 的 minor F-001 闭合成立（D-001 已写 bare addr-spec / 传输守卫，与代码一致）。本意见无 required | 各 `03-audit/A-001-*.md`；对照 D-001 与代码 |
| 拨号路径与配置面（R2） | **成立**。唯一拨号形态 = `tls.Dialer` → `smtp.NewClient`，默认端口 465；YAML `mail.smtp.*` + `MAIL_SMTP_*`；全空合法、任一非空则四键必填；secret 不回显；已提交 `configs/config.yaml` 无字面密码 | `smtp.go`；`config.go`；`configs/config.yaml`；`configs/env.example` |
| I-005 / I-006 关门叙事 | **成立**。合同仅 `TextBody`；README 明示重启生效、HTML/MIME 与热加载不进分母 | GOAL-005 D-002；`kernel/mail.go`；`apps/api/README.md` §出站邮件 |
| 本波运行时接线 | **与书面决定一致，不构成缺口**。R4 起 `MailSender` 不再进 fx（GOAL-005 D-001：本波无消费方，构造 = 启动 fail-closed + Ping）；`_ = mailSender` 只把 probe 交给 readyz。下一消费 VP 须把 sender 注入模块 | `composition.go:370-380`；GOAL-005 D-001 |
| `internal/store` 两测失败 | **不阻塞本次关门**。本轮确认 `git diff be23164..HEAD -- apps/api/internal/store` 为空；两测为 live-Postgres 门控集成测试（`postgres_test.go`），属共享 probe 库环境/遗留库状态，不在本区 diff 与退出分母内 | `apps/api/internal/store/postgres_test.go:461-547`；本轮 git diff |
| N-001 live 未实跑 | **不阻塞**。VP 判据 3 允许等价 harness；本轮独立复跑 harness 全绿。self 将其标 `accepted-residual` 用词偏重（P-003 residual 须用户书面接受），实质是分母外 note，无需升级 P-004 | GOAL-005 A-001 N-001；`smtp_test.go` |
| 本轮独立复跑 | kernel / mail / config `-count=1` 全绿；composition 包测试绿；vet 干净（go module `apps/api`） | 本轮 `go test` / `go vet` |
| 子目标 self 已闭合 minor | 闭合成立，不重开 | GOAL-002/003 A-001 F-001 |
| 外围指针 | 见 Findings F-001 / F-003（不阻断） | `workspace.md` 纲领表；VP-017 信息表 |

### Findings

| F-ID | 级别 | 内容 | 建议 | 是否可 fixed 闭合 |
|------|------|------|------|-------------------|
| F-001 | recommended | `workspace.md`「纲领阶段」表仍写 R1～R4 **未开始**，与 Root `00-meta` / goal-tree **4/4 已完成** 不一致。canonical 路线图在 Root，不构成关门名不副实。 | `/govern` 响应时把该指针表改为已完成（或改为「见 Root 00-meta」）。 | 是（改指针表即可） |
| F-002 | recommended | `SMTP.Send` 在对端 **未广告 AUTH** 时跳过 `PlainAuth`，仍投递。配置面强制 username/password，但发送路径可能在无 AUTH 的端点上静默不认证。这不是 D-001「凭证可选」的复活（配置仍 fail-closed），也不在 VP-017 退出分母（open-relay 归 VP-009）。 | 后续 hardening 或下一消费 VP 接入前：若 EHLO 未广告 `AUTH`（或未含 `PLAIN`）则 fail-closed，并补一条「无 AUTH 必败」测试。 | 是（改 `smtp.go` + 测试）；**不作为本次关门条件** |
| F-003 | recommended | VP-017 信息表 I-017-001～005 仍为 `open`、I-017-006 仍为 `registered`，而 Root 对应项已 verified。属愿景层台账滞后；Root 关门不代替 `/vision` 收 VP。 | VP 收尾时把 I-017-001～006 改为 verified，并填写关门记录。 | 是（`/vision` 收 VP-017）；**不作为 Root `done` 条件** |

无 required / 必改项。

### 结论

**同意将 `[workspace-017-outbound-mail] GOAL-001-outbound-mail` 标为 `done`。** 六条方向级退出判据均有可回指代码与测试证据；R1～R4 子目标 done 且 self 闭合成立；影响本 scope 的 required 信息项已 verified；`internal/store` live-Postgres 失败与 N-001 live 未实跑均不进入本波分母。

无条件门禁。保留意见仅限上表三条 recommended：工作区指针与 VP 信息表滞后、以及 AUTH 未广告时的发送姿态——均可在 `/govern` 响应与后续 `/vision` VP-017 收尾中处理，**不阻断** Root 关门。
