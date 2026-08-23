---
id: GOAL-001-observability
doc: decision
status: active
parent: null
created: 2026-08-21
updated: 2026-08-22
version: 0.2.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

> 状态以 `00-meta.md` 信息表为准（本表为镜像，须保持同号同状态）。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 指标 scrape 路径/端口/绑定鉴权/基数/最小集合/标签卫生 | R1 方案 | R1 冻结 | R1 决策 | verified（2026-08-21） | — | GOAL-002 D-001 §1–§6 |
| I-002 | required | OTLP 协议、采样、no-op | R3 方案 | R3 接入前 | R3 决策 | verified（2026-08-22） | — | GOAL-004 D-001 |
| I-003 | required | Store / 对象存储 / Job 是否进本波 | R1 方案 | R1 冻结 | R1 决策 | verified（2026-08-21，出局） | — | GOAL-002 D-001 §7 |
| I-004 | required | metrics/OTLP 是否扩 `readyz` | R2/R3 方案 | R2/R3 接入前 | R2/R3 决策 | verified（2026-08-21，不扩） | — | GOAL-002 D-001 §8 |
| I-005 | required | request-id 写入 span 的属性名 / baggage | R4 方案 | R4 关联前 | R4 决策 | verified（2026-08-22） | — | GOAL-005 D-001 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-21 | 开区 scaffold 与 A4 纲领路线图 | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
