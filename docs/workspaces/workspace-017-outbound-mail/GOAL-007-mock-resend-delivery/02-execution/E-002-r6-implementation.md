---
id: E-002
doc: execution-entry
goal: GOAL-007-mock-resend-delivery
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-002 · R6 实施执行（2026-08-24）

## 已发生事实

按 D-001 范围实施 GOAL-006 D-002 合同的代码落地，全部测试绿：

1. **配置层（C1）**：`internal/config` 新增 `mail.channel`（env `MAIL_CHANNEL`，值归一化小写）与 `mail.resend.api-key` / `mail.resend.from`（env `MAIL_RESEND_API_KEY`/`MAIL_RESEND_FROM`）；`ResolveMailChannel()` 实现冻结解析算法（显式选择→块可用性校验；空值推导：单生产块胜出、双全配歧义 fail-closed、均未配 mock）；`validateMail` 扩为双块配对 + 解析门禁（每个环境生效）。`config.default.yaml` 与 `configs/.env.example` 同步键位文档。测试：`config_mail_channel_test.go`（解析矩阵、fail-closed、secret 不泄漏）。
2. **持久层（C2）**：编译迁移 **0051 `mail_outbox`**（owner `core.persistence`；SQLite INTEGER / PG BIGINT 双方言 DDL + created_at 索引）；store 侧冻结账目同步（catalog ownership want 表、identity head 50→51、lockedHeadExtraTables[51]、fresh/restart/operations 期望 51 条）。
3. **mock 发布器**：`internal/mail/outbox.go` `OutboxSink` 实现 `kernel.MailSender`——写表与有界淘汰同事务（默认保留最近 `DefaultOutboxCap=500` 条，`created_at DESC, id DESC`）；`List`（新→旧分页、无正文）/`Count`/`Get`（含正文；未知 id → `ErrOutboxRecordNotFound`）。id 生成沿用 (UnixNano, 单调序, rand 扰动) 进程内唯一先例。测试覆盖写入/排序/淘汰/重启持久化/端口校验。
4. **Resend 适配器**：`internal/mail/resend.go` `Resend` 实现 `kernel.MailSender`——POST `{base}/emails`，Bearer 鉴权，2xx 即成功；非 2xx 报状态码不泄漏密钥；构造器 fail-closed（api-key/from 必填、from bare 地址）；BaseURL/Client 为 harness 测试缝。readyz 探针留空待 R8。
5. **面层（C3）**：composition `newMailSender(cfg, st, logger)` 改为按 `mail.channel` 三路解析（mock 默认 → OutboxSink；CaptureSink 不再是装配默认，仅存内部测试替身）；新增 `GET /api/mail/outbox`(+`/{id}`)（handler `RegisterMailOutbox`，Bearer + `settings.read` 门禁，统一 `{items,total,page,pageSize}` 包络，limit 上限 200）；路由独立于 `/api/settings/*`（落实 I-012）。错误码复用 pinned 集（INVALID_PAGE_SIZE / INVALID_PAGE / NOT_FOUND / INTERNAL），未扩错误码合同。README 出站邮件节与端点表更新。
6. **回归**：`go build ./...` 通过；config/store/mail/composition/handler 全部包测试绿；全仓 `go test ./...` 复跑见 E-003。

## 证据

| 主张 | 路径 |
|------|------|
| 配置解析 | `apps/api/internal/config/config.go`（`ResolveMailChannel` / `validateMail`） |
| 迁移 0051 | `apps/api/internal/modules/corepersistence/migration/migration.go` |
| mock 发布器 | `apps/api/internal/mail/outbox.go` |
| Resend 适配器 | `apps/api/internal/mail/resend.go` |
| 接线与管理 API | `apps/api/internal/composition/composition.go`、`apps/api/internal/handler/mail_outbox.go` |
| 冻结账目同步 | `apps/api/internal/store/migrate_test.go`、`identity.go`、`identity_test.go` 等 |

## 未做

- 未动设置页 / 热切换 / 试发（R7）；未做 Resend live 投递与生产探针（R8）；VP-018 保持冻结。
