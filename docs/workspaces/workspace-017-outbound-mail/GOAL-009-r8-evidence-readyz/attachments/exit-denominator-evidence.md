# VP-017 现行退出分母 · 证据包（R8 · GOAL-009）

> 对照 [VP-017-outbound-mail](../../../../../vision/plans/VP-017-outbound-mail.md) v0.4.0「方向级退出判据（现行 · 再关门用）」1～7 与 Root `GOAL-001-outbound-mail` 成功标准逐条登记。日期：2026-08-24。本文件为关门审计的证据索引，不构成放行本身。

## 判据 1 — 内核发送端口仍是唯一合同；公共面无供应商客户端类型

| 证据 | 指针 |
|------|------|
| R1 冻结合同未改写 | `apps/api/internal/kernel/mail.go`（`MailSender` / `MailMessage`；From 由适配器盖章） |
| 公共面 sweep 结论保持 | R3 历史（GOAL-004）+ 现有 handler/module 无 mock/resend/smtp 类型引用（R6/R7 新增面均只消费 kernel 端口或 mail 内部接口） |
| 测试 | `kernel/mail_test.go` 端口面测试持续通过 |

## 判据 2 — 具名渠道落地；默认 mock 发布到管理员可检视站内出站记录；未配置生产渠道进程可启动

| 证据 | 指针 |
|------|------|
| 渠道 id 封闭集 + 解析算法 | `internal/config/config.go` `ResolveMailChannel`（GOAL-006 D-002 §2）；`config_mail_channel_test.go` |
| mock 发布器 | `internal/mail/outbox.go`（写入+淘汰同事务，默认保留 500） |
| 管理员检视面 | `GET /api/mail/outbox(+/{id})`（`handler/mail_outbox.go`；settings.read 门禁） |
| 未配置可启动 | composition 默认路径测试 `TestMailRuntimeDefaultsToSwitcherOverMockWithoutProbe`；cmd/server 全量启动测试 |

## 判据 3 — 显式 Resend 后可核对至少一封投递（live 或等价 harness）；配置不完整 fail-closed

| 证据 | 指针 |
|------|------|
| 等价 harness（离线） | `internal/mail/resend_test.go`：httptest 断言 POST /emails 请求形状（Bearer、from/to/subject/text）、2xx 即成功、非 2xx 报状态码 |
| **live 已实跑（2026-08-24，用户裁决补跑）** | `internal/mail/resend_live_test.go` 实跑 PASS：Ping 2xx + POST /emails 2xx，真实投递至 magicvr@hotmail.com。首次尝试 from=eshowy.top 被 403 拒（域名未验证——运营项：生产启用前须在 resend.com/domains 验证）；改用 resend.dev 沙箱发件人成功 |
| fail-closed | 触碰 resend 块即要求 api-key/from 齐备（双层：validateMail + 构造器）；`TestMailResendMisconfigDoesNotLeakSecret` 等 |

## 判据 4 — 管理员可在设置面选择渠道、填写配置、热切换、同一 MailSender 试发

| 证据 | 指针 |
|------|------|
| web 邮件 tab | settings.json tab-mail（渠道 select、参数表单、试发表单、mock 记录 table）；i18n 双语各 21 键 |
| 热切换语义 | `internal/mail/runtime.go` Switcher（DB 权威、updated_at 缓存失效、候选先验证后落库=切失败保留旧 sender）；Root D-007 |
| 试发同端口 | `POST /api/mail/test-send` 走 Switcher.Send；handler 测试断言 mock 下 test-send 落 outbox 表 + 审计事件 |
| 密钥 write-only | PublicView 仅 *Set 布尔；AES-GCM 加密落库（secrets.go 测试） |

## 判据 5 — 仅显式生产渠道后 readyz 扩依赖；无越界引入；R1～R4 实施史未回退

| 证据 | 指针 |
|------|------|
| 探针接线 | `composition.newMailRuntime`：boot=resend→Resend.Ping；boot=smtp→ESMTP Ping（R4 原样）；mock/空→nil（`TestMailRuntime*` 三态断言） |
| Resend 探针 | `resend.go Ping`（GET /domains，5s 超时，仅报状态码）+ `TestResendPing` |
| 史不回退 | SMTP/CaptureSink 适配器原样保留；既有 smtp/composition 测试全绿 |
| 越界检查 | 无 SMS/用户通知/账号 email/邀请/自助恢复代码与 schema 变更；VP-018 保持冻结 |

## 判据 6 — 开放 required finding = 0（或已合法闭合）

| 证据 | 指针 |
|------|------|
| 子目标台账 | GOAL-002～009 各 03-audit：self pass、required 全闭合（GOAL-007 F-001 accepted-residual 经用户书面裁决；GOAL-008 F-001 residual 由 R8 承接——见下） |
| 本目标 | GOAL-009 关门审计后归零 |

## 判据 7 — VP-018 仅在本 VP 再次 closed 之后解冻

| 证据 | 指针 |
|------|------|
| 冻结状态 | workspace-018 / VP-018 未推进；本工作区五件套无 018 改动 |
| 解冻动作 | 归愿景层收口（VP-017 再 closed 之后由 `/vision` 处理），不在实现层预支 |

## Root 成功标准对照（方向级 1～5）

| # | 结论 |
|---|------|
| 1 | 满足（判据 1 同源） |
| 2 | 满足（判据 2） |
| 3 | 满足（判据 3：harness 已绿；live 为 opt-in 缝，补跑与否随关门问询留痕） |
| 4 | 满足（判据 4） |
| 5 | 满足（判据 5） |
