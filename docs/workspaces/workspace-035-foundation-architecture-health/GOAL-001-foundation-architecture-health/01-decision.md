---
id: GOAL-001-foundation-architecture-health
doc: decision
status: active
parent: null
created: 2026-09-09
updated: 2026-09-09
version: 0.3.0
---

# 决策记录 · GOAL-001-foundation-architecture-health

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|------------------|------|-------------|-------------|
| I-035-001 | required | 对照分母覆盖哪些包、端口、Profile 与文档 | R2 | R1 | 包含/排除表 | **verified** | — | GOAL-002 D-001 + 附件 §1 |
| I-035-002 | required | 业界参照集四类 | R3 | R1 | 用户书面 | **verified** | — | VP-035 正文；VRev-086 |
| I-035-003 | required | 对照是否迫使改 Charter 非目标 | R3→R4 | R3 | 逐行检查 | collecting | — | 待 R3（GOAL-004 C4 前） |
| I-035-004 | required | closed VP residual 入册范围 | 分类完成 | R1 | 扫描 VP-013～034 | **verified** | — | GOAL-002 D-001 + 附件 §2 |
| I-035-005 | required | 现在修 vs 另立 | 任何代码整改前 | R1 | 用户 R1 冻结；默认另立 | **verified** | — | GOAL-002 D-001 + 附件 §3 |
| I-035-006 | required | R3/R4 independent 门禁的 provider 与模式（会话指令 = 本地 codex `gpt-5.6-sol`·high；项目级默认 = grok build 4.6） | R3 C4 交叉审计 / R4 关门 | R3 C4 之前 | 询问用户并留痕；不静默择一、不降级 | **待用户裁决** | 下一复核 = R3 C4 之前 | [GOAL-004 D-001](../GOAL-004-r3-industry-comparison/01-decision/D-001-r3-execution-boundary.md) 待确认项 C |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-09 | 工作区 / Root 建立与开区决策 | accepted | `01-decision/D-001-workspace-root-establishment.md` |
| D-002 | 2026-09-09 | R1 分母冻结（由 GOAL-002 承载） | accepted | `../GOAL-002-r1-denominator-freeze/01-decision/D-001-r1-denominator-freeze.md` |
| D-003 | 2026-09-09 | R2 as-built 矩阵与关门（由 GOAL-003 承载；A-002 F-001 `fixed`） | accepted | `../GOAL-003-r2-as-built-matrix/00-meta.md` 关门记录 + `../GOAL-003-r2-as-built-matrix/03-audit/A-003-r2-a002-response.md` |
| D-004 | 2026-09-09 | R3 执行边界与分类规则（由 GOAL-004 承载；**提议**，含 3 项待用户确认） | proposed | `../GOAL-004-r3-industry-comparison/01-decision/D-001-r3-execution-boundary.md` |
