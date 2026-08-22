---
id: GOAL-005-r4-readyz-evidence
doc: execution-entry
record_id: E-002
goal: GOAL-005-r4-readyz-evidence
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-002 · 关门响应：A-002 F-002 修复（AUTH 未广告即拒发）（2026-08-22）

## 已发生事实

- 独立关门审计（Root A-002，grok build）verdict `pass`、无 required；三条 recommended 中 F-002 由本轮立即修复：
  - `internal/mail/smtp.go`：`Send` 在对端 EHLO 未广告 `AUTH` 时 **fail-closed 拒发**（配置合同强制凭证，静默跳过认证 = 显式端点降级为匿名投递）。
  - `smtp_test.go`：新增 `TestSMTPSendFailsClosedWithoutAuthAdvertisement`（无 AUTH 广告的 fake 对端 → Send 必败且零 envelope/payload 到达）；harness 抽出 `startFakeSMTPVariant(advertiseAuth)`。
- 实跑：`go test ./internal/mail/ -count=1` 全绿；vet 干净。

## 证据

| 主张 | 路径 |
|------|------|
| 独立意见原文 | `../../GOAL-001-outbound-mail/03-audit/A-002-independent-closeout.md`（F-002 行） |
| 修复实现 | `apps/api/internal/mail/smtp.go`（Send 的 AUTH 门） |
