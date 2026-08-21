---
id: GOAL-004-s3-settings-productization
doc: decision
status: done
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# 决策记录 · GOAL-004（S3）

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 四类字段语义（合并/清空/恢复默认/校验/时区失败）是否齐备 | 实施（C1–C6） | 实施前 | 读 Root D-002 + 现有 settings 栈 | **closed** | — | Root D-002 冻结（2026-08-09 用户裁决） |
| I-002 | required | Settings 页 schema 形态与预览/恢复默认承载 | C4 | 实施前 | 盘点 Renderer custom action 现状 | **closed** | — | 见 D-001 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-09 | S3 · 四类字段语义、Settings 页形态与预览/恢复默认承载 | proposed | `01-decision/D-001-s3-field-semantics-and-page-shape.md` |
