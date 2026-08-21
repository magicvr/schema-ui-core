---
id: D-002-audit-mode
doc: decision-entry
status: accepted
parent: GOAL-008-audit-log-retention-settings
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# D-002 · 关门审计模式

GOAL-008 触及审计日志生命周期、SQLite 迁移（0046/0047）与过期删除/归档，属 data/migration/security 高影响门禁。

## 决定

- 模式：**independent**（项目路径仍先 self，再 independent）。
- provider：项目级 grok-build（grok 4.6 · reasoning high；`docs/architecture/independent-audit-execution.md`）。
- 不选 `cross`：无元规则/协议冲突、无证据互否；设置页增量不是协议 pin / Manifest 装配变更。
- 不选 `self` 单轨：过期删除会永久丢掉热表行，不可仅自审关门。

I-002 → **verified**。
