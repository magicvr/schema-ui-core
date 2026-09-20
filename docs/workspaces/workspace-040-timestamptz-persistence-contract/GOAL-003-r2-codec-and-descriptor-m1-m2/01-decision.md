---
id: GOAL-003-r2-codec-and-descriptor-m1-m2
doc: decision
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# 决策记录 · GOAL-003

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 状态 | 证据 / 决策 |
|----|------|-----------------|----------|------|-------------|
| I-041-001 | required | Go codec 公共 API 形态（签名 / 错误分类 / 导出面） | A | **collecting** | 本目标内定稿并落盘；技术细节由独立审计复审 |
| I-040-001 | required（继承） | 逐列编解码与精度合同 | A/B/C | **verified** | R1 关门时已 verified |
| I-040-003 | required（继承） | 双方言原地转换 / 失败恢复 / 备份依赖 | B/C | **verified** | R1 关门时已 verified |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| — | — | 尚无本目标内决策 | — | — |

> 新决策从本目录 `D-NNN-<slug>.md` 写入；编号在本目标内单调不复用，允许空洞。
> **继承边界**：本目标的合同权威均来自 Root 决策（`D-004`/`D-014`/`D-015`/`D-016`）与 `GOAL-002` 的 `D-017`–`D-021`；本目标**不得**静默改写它们。引用时一律写 **Root** / **GOAL-002** 限定。
