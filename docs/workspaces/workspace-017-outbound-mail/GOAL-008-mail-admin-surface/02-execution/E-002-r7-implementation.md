---
id: E-002
doc: execution-entry
goal: GOAL-008-mail-admin-surface
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-002 · R7 实施执行（2026-08-24）

## 已发生事实

1. **持久层**：编译迁移 **0052 `mail_config`**（单行运行时渠道状态：channel、mock_retention、resend/smtp 参数；密钥字段 `_enc` 后缀 AES-GCM 加密存放）与 **0053 `operation_log_mail_events`**（operation_log 事件 CHECK 枚举扩入 `mail.channel-update` / `mail.test-send`）。0053 的重建同时搬运 correlation 与 session 两张 FK 侧表（rename 会把侧表外键改写到 operation_log_old，仅搬 correlation 不够——实施中发现并修正，SQLite 与 PG 双路径均验证）。store 冻结账目同步至目录头 53。
2. **加密**：`internal/mail/secrets.go` —— `MAIL_CONFIG_MASTER_KEY` env（SHA-256 派生）或首次自动生成 32 字节 key 文件（0600）；AES-256-GCM 加解密；空串直通。
3. **热切换运行时**：`internal/mail/runtime.go` `Switcher` 实现 `kernel.MailSender`——mail_config 行按 D-007 种子一次（文件层解析结果为初值），DB 即权威；每次 Send 按 updated_at 缓存解析当前适配器；`Update` 先构造候选适配器校验（SMTP 另加 5s Ping），失败即在落库前返回错误（切失败保留旧 sender）；成功后原子写行 + 刷新缓存。`PublicView` 只含密钥是否已设置的布尔值，无任何密钥值。
4. **管理 API**：`GET/PUT /api/mail/config`（settings.read / settings.write 门禁）+ `POST /api/mail/test-send`（走当前 Switcher，禁止旁路）；PUT/test-send 记 operation log（新事件常量）；错误码复用 pinned 集并新增 INVALID_MAIL_CONFIG / MAIL_SWITCH_REJECTED / MAIL_SEND_FAILED（已入错误目录双语表与 error_contract 冻结清单）。
5. **composition 接线**：`newMailRuntime(cfg, st, logger)` 构造 Switcher 并注入 handler；readyz 探针语义保持 R4 现状（仅 boot=SMTP 时挂 ESMTP Ping；Resend 探针留 R8）。
6. **web 设置页**：`settings.json` 新增「Mail」tab（schema-driven）：渠道 select、mock 保留条数、Resend api-key/from、SMTP host/port/username/password/from 表单（recordSource GET /api/mail/config 点路径回填；密钥字段 GET 不回填=空提交即保持不变）、试发表单（POST /api/mail/test-send）、mock 出站记录 table（dataSource /api/mail/outbox，to/subject/created_at 列）；en-US/zh-CN 各补 21 个 i18n 键。

## 测试证据

- api：config/store/mail/handler/composition 定向测试全绿；全仓 `go test ./...` 结果见 E-003。
- web：i18n 全套（含 schema-keys 结构校验、bilingual 渲染）+ representative-pages 绿；tsc/vite build 与全量 vitest 见 E-003。

## 未做

- Resend live 投递证据与 readyz 生产探针（R8）；多实例一致性（非目标）。
