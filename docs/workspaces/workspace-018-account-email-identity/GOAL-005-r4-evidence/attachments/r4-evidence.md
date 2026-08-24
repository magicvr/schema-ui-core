# R4 证据包 · 账号邮箱身份面（2026-08-24）

对照 VP-018 方向级退出判据与 Root 成功标准逐条固化。核心证据 = 端到端测试
`apps/api/internal/modules/authsession/email_identity_e2e_test.go`
（`TestR4EndToEndBindVerifyThroughMockChannel`）。

## 判据 1 · 可空邮箱 + 无邮箱账号可登录

- `users.email TEXT` 可空（迁移 0054）；存量行升级后 (NULL, NULL) = unbound
  （`store/migrate_0054_test.go::TestMigrate0054UpgradePathLandsUnboundRows`）。
- 多账号 NULL 共存：唯一表达式索引对 NULL 互异（同文件 Fresh Catalog 语义用例；
  R4 e2e 再次断言 bob 全程 unbound 且 `UserByID` 正常）。
- email 不是启动硬依赖：无任何 Profile/装配改动依赖该列（GOAL-003 A-001 已核对边界）。

## 判据 2 · 绑定/校验流落地，校验信经 kernel.MailSender，且可从 017 默认渠道取信

链路（全部走真实适配器，非测试桩）：

```
BindEmail ──(phase1: 占槽+挑战落库)──▶ phase2: mail.OutboxSink.Send
                                          └─ INSERT mail_outbox（管理员可检视）
VerifyEmail ◀── 从出站记录 Get(body) 取 6 位码 ──┘
```

- 出站记录经 `mail.NewOutboxSink`（VP-017 GOAL-007 冻结契约的 mock 渠道适配器）
  写入 `mail_outbox`（迁移 0051）；`List/Get` 为管理端读取面原语。
- 两阶段派发（R4 发现并修正）：渠道适配器自持存储事务，嵌套 Run 禁止
  （R1 v1.4 §2）——派发移出占槽事务，失败以补偿恢复先前身份/挑战
  （`TestSendFailureRollsBindBack` + resend 快照恢复路径）。
- 生产渠道（SMTP/Resend）走同一端口，验收形态不变。

## 判据 3 · 唯一性 fail-closed + 换绑同合同

- 他号大小写折叠冲突 → `EMAIL_TAKEN`（e2e 步骤 3 + 服务测试矩阵）。
- 换绑 = 覆写重置 pending，旧址槽随覆写释放
  （`TestRebindOverwriteReleasesOldSlot`）。
- 同址已 verified 再绑定幂等不重发；同址 pending 重绑受 60 秒冷却
  （A-001 F-002 后补齐）。

## 判据 4 · 未越界

无忘记密码状态机 / 邀请 / 密码策略产品 / SMS / 第二运输 / 模板中心；
未改 Charter、未改 Profile 默认集。验证码仅作邮箱归属校验，不作恢复证明。

## N-1 边界声明（SQLite lower() ASCII）

唯一性判定统一走 SQL `lower(email)` 表达式索引：SQLite 侧 lower() 仅折叠
ASCII 字母。非 ASCII 大小写（如土耳其语 İ/ı）在该方言下归一不完全——
接受为**有界残余**：(a) 电子邮件本地部分实践中 ASCII 为主，IDN 域名为
punycode；(b) 应用层入口 trim + 折叠归一已覆盖常规输入；(c) PostgreSQL
生产方言为 locale 感知 lower()，不受此限。复核触发：出现非 ASCII 邮箱
真实需求或 PG→SQLite 双向迁移产品化时重开（归属：后续 IAM 波次台账）。

## 判据 5 · 开放 required finding

GOAL-002/003/004 各审计台账开放 required = 0（A-001×3：self pass /
independent pass / conditional→响应后归零）。本目标 C3 self 审计见其台账。
