---
id: GOAL-001-full-protocol-contract-v2-7-0
doc: decision
status: active
parent: null
created: 2026-08-08
updated: 2026-08-10
version: 0.1.1
---

# 决策记录 · GOAL-001-full-protocol-contract-v2-7-0

## 信息需求与阶段门禁

> 权威信息项表见 `00-meta.md`；本索引同步关键门禁项。长决策写入 `01-decision/D-NNN-*.md`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-PROTO-FULL-001 | required | 整份契约覆盖 disposition 冻结 | S1；S2–S5 分母 | S1 前 | 新文件 + Root 决策 | **closed** | — | `attachments/I-PROTO-FULL-001-coverage-v2-7-0.md`（v1.0.1）+ D-002 冻结 + D-003 勘误（2026-08-10） |
| I-001 | required | v0.1.3 vs 全量差集 | S0；S1 | S0 结束前 | 盘点 | **closed** | — | `02-execution/E-002-s0-gap-analysis.md` + `attachments/I-S0-001-*` |
| I-002 | required | 范围收缩 residual（若有） | S1 | S1 决策时 | 用户书面 | **N/A** | 无域级收缩 | S0 差集全部可纳入（E-002 §4）；2 个 local adapter execution exclusion 不改变协议承诺面（D-003） |
| I-003 | non-blocking | vendor 上游策略 | 验证策略 | S2 前可决 | 策略决策 | **closed** | 继承历史（I-PROTO-004 = vendor） | 6 schema + 17 fixture 套件全部 vendor + SHA pin（`provenance.json`；uploads `aaeb9683…`、permissions-inheritance `ac124fa1…` 于 S1 补入） |
| I-004 | non-blocking | S2 批次边界 | S2 立项 | S2 | Root 决策 | **closed** | — | D-002 §4 批次 B1–B6 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-08 | 开区与纲领路线图采纳 | accepted | `01-decision/D-001-workspace-open-and-roadmap.md` |
| D-002 | 2026-08-08 | 整份契约覆盖表 I-PROTO-FULL-001 冻结（S1） | accepted | `01-decision/D-002-full-coverage-freeze.md` |
| D-003 | 2026-08-10 | I-PROTO-FULL-001 执行分母勘误 | accepted | `01-decision/D-003-i-proto-full-errata.md` |
