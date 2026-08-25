---
id: E-002
doc: execution-entry
goal: GOAL-002-iam-contract-freeze
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# E-002 · 第二轮五项裁决 + D-001 合同落盘（C2 满）

2026-08-25 完成：

- 用户在结构化裁决会话（`i003_policy_defaults` / `i004_invite_form` / `i005_invite_lifecycle` / `i007_existing_accounts` / `i008_session_semantics`）就剩余五项逐一书面选定，均取推荐项：
  - I-003 = 起步宽松（默认 min 8 字节、复杂度/历史关；可配置最小长度/类别数/历史深度）
  - I-004 = 双形态并存（邮件经 MailSender 或链接手动出示）+ 邀请即建号
  - I-005 = 默认 7 天有效 · 一次性 · 可撤销 · 重发撤旧发新（冷却 60 秒）
  - I-007 = 渐进生效：仅设密时刻强制，不扫描存量、不强登出
  - I-008 = 投影现行设密语义（token_version 前移撤销其余会话）
- 合同条款落盘：[D-001-iam-contract-freeze.md](../01-decision/D-001-iam-contract-freeze.md) §1～§5；Root 镜像表 I-003～I-008 同步为 verified。
- 未改动任何产品代码。

后续：C3 self 关门审计（A-001）→ GOAL-002 关门 → Root R1 记完成（1/4）。
