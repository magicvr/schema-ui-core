---
id: D-001-r2-audit-event-model
doc: decision-entry
status: accepted
parent: GOAL-003-r2-audit-event-model
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# D-001 · R2 审计事件模型范围与信息门禁

## 决策

1. R2 先以 auth、settings、users 三类真实 mutation 为消费切片，建立版本化 detail envelope、统一脱敏规则与 correlation 关联；不一次性改造所有事件目录。
2. 历史 detail 保持读取兼容，新增事件明确 schema 版本；任何无法确认安全性的字段默认脱敏或拒绝写入。
3. S0 必须先完成现有 detail/敏感字段/API 兼容扫描并关闭 I-001；S1 实施前必须完成 I-002 的审计模式/provider 裁决。
4. R2 依赖已关闭的 R1 `GOAL-002-r1-correlation-error-contract`，不重复实现 request-id 生成。

## 依据与取舍

- operationlog 已有独立 `operation_log_correlation` 关系表；R2 复用该关系，不把 correlation 再复制进业务 detail。
- 既有 auth 测试要求 detail 不包含 token/password/secret；该约束升级为统一脱敏器的回归基线。
- 全量事件目录改造会扩大风险与审计面，先用三类 mutation 验证模型，再按后续阶段扩展。

## 不确定项

- I-001：settings 其余配置字段及非核心模块 detail 的敏感性仍待扫描。
- I-002：R2 属于 security/data 影响范围，按 P-003 需要 independent 或 cross 审计；provider 需按 P-004 取得用户裁决。
