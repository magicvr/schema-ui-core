---
id: GOAL-001-object-storage
doc: decision
status: active
parent: null
created: 2026-08-21
updated: 2026-08-21
version: 0.2.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | S3 API 子集与驱动（MinIO / R2 / AWS 公约数） | R2 方案 | R2 实施前 | R2 决策 | verified | — | D-004：aws-sdk-go-v2；Put/Get/Delete/Head/List+HeadBucket 子集 |
| I-002 | required | 桶模型与三类落盘 key 隔离 | R1 方案 | R1 冻结 | R1 决策 | verified | — | D-002：单桶 + 命名空间前缀 |
| I-003 | required | 配置键名与凭证注入 | R2 方案 | R2 实施前 | R2 决策 | verified | — | D-005：键名沿 D-001；static credentials 显式构造，禁默认链 |
| I-004 | non-blocking | 存量本地 → 对象存储搬运器 | R5 叙事 | R5 | 点名 residual | recorded | 不进退出分母 | 用户已裁决：无产品搬运器 |
| I-005 | required | List/GC 是否进端口 | R1 方案 | R1 冻结 | R1 决策 | verified | — | D-003：List+Stat 进端口 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-21 | 开区 scaffold 与 A2 纲领路线图 | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
| D-002 | 2026-08-21 | 桶模型：单桶 + 命名空间前缀（I-002 闭合） | accepted | [D-002-bucket-model.md](01-decision/D-002-bucket-model.md) |
| D-003 | 2026-08-21 | 枚举能力进端口：List + Stat（I-005 闭合） | accepted | [D-003-list-gc-in-port.md](01-decision/D-003-list-gc-in-port.md) |
| D-004 | 2026-08-21 | S3 驱动与 API 子集裁定（I-001 闭合） | accepted | [D-004-s3-driver-subset.md](01-decision/D-004-s3-driver-subset.md) |
| D-005 | 2026-08-21 | 配置键名与凭证注入确认（I-003 闭合） | accepted | [D-005-credential-injection.md](01-decision/D-005-credential-injection.md) |
