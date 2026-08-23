---
id: GOAL-001-outbound-mail
doc: execution-entry
record_id: E-005
goal: GOAL-001-outbound-mail
status: recorded
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-005 · R4 执行：readyz 扩依赖与显式路径证据（2026-08-22）

## 已发生事实

- 开设子目标 `GOAL-005-r4-readyz-evidence`（五件套齐全，parent = 本 Root），承接 R4 治理上下文。
- 落地代码：`internal/mail/smtp.go` Ping probe（与 Send 共用冻结拨号形态）；composition 三元接线 + probe 进 `RegisterWithMFAProbes`；env-gated live 投递测试镜像 s3_live 先例；README 出站邮件节。
- 实跑证据：mail/composition 全绿（含双 profile lifecycle 启动不变量）；vet 干净。
- I-005/I-006 关门叙事留痕（GOAL-005 D-002）；Root 路线图 R4 → 已完成，progress 4/4——四阶段全部完成，进入关门审计。

## 证据

| 主张 | 路径 |
|------|------|
| 子目标决策/执行/审计 | `../GOAL-005-r4-readyz-evidence/` |
| 本 Root 决策 | [D-005-r4-readyz-and-closeout.md](../01-decision/D-005-r4-readyz-and-closeout.md) |

## 未做

- live env-gated 测试未实跑（无真实凭据；离线等价 harness 已覆盖，见 GOAL-005 A-001 N-001 残余留痕）。
- VP-017 关门记录属愿景层，走 `/vision` 收尾（本区不代写）。
