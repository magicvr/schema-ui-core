---
id: GOAL-004-w3-schema-host-protocol-conformance
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-12
updated: 2026-08-13
version: 0.1.1
---

# 决策记录 · GOAL-004

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 2.7.0 覆盖/偏离基线完整 | S2 | S2 冻结前 | upstream 与 API/Web 逐项对照 | collecting | — | `attachments/I-HOST-APP-001-protocol-gap-catalog.md` |
| I-002 | required | 全部 Host/App 候选有明确协议处置 | S2 | S2 冻结前 | 逐项 adopt/reserve/out | collecting | — | 上游 H0 处置已同步（附件 §1b/§1c/§6；ADR-0034 D10，proposed，95/95）；ADR-0034～0037 已 accepted（2026-08-13） |
| I-003 | required | 上游新协议已发布/固定并进入本仓 | S4 | S4 开始前 | 版本、provenance 与工件验证 | verified | — | 上游 H2 机器契约固定于 commit `453008d` 并 pin 入本仓（`provenance-v2.8.json`）；S4 生产实现完成（E-004）。正式 2.8.0 发布（tag `521cff8`，上游审计 0080 V379 权威）后已按正式身份重 pin 并重生成 claim（E-005；此前 `593f625` 为 H4 预备身份） |
| I-004 | required | cross 审视 provider 已指定 | S2/S6 | 首次 cross 审视前 | 用户指定并落 independent A 条目 | verified | — | 用户指定 `grok build`（grok 4.5，reasoning high）；self=A-001（pass）；independent=A-002（conditional，BLOCKING_COUNT=0）已落盘 |
| I-005 | required | 兼容、迁移、弃用与 fail-closed 规则 | S3/S4 | S3 固定前 | 兼容矩阵与正反 fixtures | collecting | — | 上游已交付（`migrations/2.7-to-2.8.md`、registry 弃用机制、host/app-manifest 正反 fixtures）；本仓 S3 台账录入收尾中 |
| I-006 | required | `recordView` 行上下文等争议语义归属 | S2/S4 | S2 冻结前 | upstream 裁定 | verified | — | 上游已裁定（ADR-0034 D6 IMP-004 保留独立 overlay ADR；D7 reserve 不冒充 capability） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-12 | 协议优先与实现停止线 | accepted | `01-decision/D-001-protocol-first-remediation-gate.md` |
