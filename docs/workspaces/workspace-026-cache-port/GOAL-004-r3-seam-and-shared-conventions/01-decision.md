---
id: GOAL-004-r3-seam-and-shared-conventions
doc: decision
status: active
parent: GOAL-001-cache-port
created: 2026-09-01
updated: 2026-09-01
version: 0.1.0
---

# 决策记录 · GOAL-004 R3 接缝与共享约定

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-026-004 | non-blocking | mail `cachedAdapter` 迁移与否（版本戳 vs 通用 TTL） | 判据 #2（评估面） | R3 | lead 评估 + 用户确认 | **verified** | — | 2026-09-01 用户确认：**不迁移，评估留痕**（评估见 `03-audit/../attachments/`；D-001） |

（I-026-001/002/003 已 verified；R3 无新 required 信息项。）

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-01 | R3 裁决：I-026-004 不迁移 + F-002 fx 容器挂载（用户裁决 ×2） | accepted | `01-decision/D-001-r3-adjudication.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。