---
id: GOAL-004-w3-security-audit-remediation
doc: decision
status: active
parent: GOAL-001-production-hardening
created: 2026-08-11
updated: 2026-08-11
version: 0.1.0
---

# 决策记录 · GOAL-004

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 审计 finding 清单 | 方案 | 方案前 | 会话四路审计 | verified | — | D-001 |
| I-002 | required | batch 原子实现面 | 实施 | 实施前 | entity/store 边界 | verified | — | D-001 |
| I-003 | non-blocking | 反代真实 IP 规则 | 验收 | 验收前 | private peer + X-Real-IP | collecting | — | D-001 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-11 | W3 修复范围与技术取舍 | accepted | `01-decision/D-001-w3-remediation-scope.md` |
