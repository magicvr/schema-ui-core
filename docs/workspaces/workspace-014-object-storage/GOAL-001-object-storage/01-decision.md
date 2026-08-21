---
id: GOAL-001-object-storage
doc: decision
status: active
parent: null
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | S3 API 子集与驱动（MinIO / R2 / AWS 公约数） | R2 方案 | R2 实施前 | R2 决策 | open | — | 待确认 |
| I-002 | required | 桶模型与三类落盘 key 隔离 | R1 方案 | R1 冻结 | R1 决策 | open | — | 待确认 |
| I-003 | required | 配置键名与凭证注入 | R2 方案 | R2 实施前 | R2 决策 | open | — | 待确认 |
| I-004 | non-blocking | 存量本地 → 对象存储搬运器 | R5 叙事 | R5 | 点名 residual | recorded | 不进退出分母 | 用户已裁决：无产品搬运器 |
| I-005 | required | List/GC 是否进端口 | R1 方案 | R1 冻结 | R1 决策 | open | — | 待确认 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-21 | 开区 scaffold 与 A2 纲领路线图 | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
