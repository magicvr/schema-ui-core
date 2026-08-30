---
id: GOAL-001-config-export-diff-dryrun-import
doc: decision
status: active
parent: GOAL-001-config-export-diff-dryrun-import
created: 2026-08-30
updated: 2026-08-30
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

> 本文件是稳定索引。长决策和独立决策记录放在 `01-decision/D-NNN-<slug>.md`，每条记录必须保持可独立阅读。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

信息台账权威见 `00-meta.md`「信息就绪与未知项」；此处投影门禁要点：

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-025-001 | required | 配置包内容边界与密钥处理 | 方案冻结 + 判据 #1 | R1 | 用户裁决 | open（待裁决） | — | 待确认 |
| I-025-002 | required | 落地形态（CLI vs 管理面 vs 两者） | 方案冻结 | R1 | 用户裁决 | open（待裁决） | — | 待确认 |
| I-025-003 | non-blocking | diff 语义与输出格式 | 判据 #2 | R2 | lead 建议 + 用户确认 | open | — | 待确认 |
| I-025-004 | required | 导入失败快照/回滚语义 | 判据 #4 | R3 | 用户裁决 | open（待裁决） | — | 待确认 |
| I-025-005 | required（投影） | Profile 默认集/Manifest 红线 | 退出分母 | R1 | 冻结不进（台账投影） | **registered** | — | VP-025 §边界 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-30 | 工作区 Root 建立（绑定 / freshness 三字段 / 审计模式 / 信息门禁 / 现状锚点） | accepted | `01-decision/D-001-workspace-root-establishment.md` |

> 编号在本目标内单调不复用。