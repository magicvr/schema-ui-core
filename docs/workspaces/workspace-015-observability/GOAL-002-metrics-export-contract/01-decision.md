---
id: GOAL-002-metrics-export-contract
doc: decision
status: active
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

本目标不另立信息项；继承 Root GOAL-001-observability 的 I-001 / I-003 / I-004（见其 `00-meta.md` 信息表）。三者由本目标 D-001 关闭，证据链：D-001 决定 → config 实现与测试（E 条目）→ A-001 自审。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| （继承）I-001 | required | 指标面：路径/端口、绑定/鉴权、基数、内核 vs 模块最小集合、标签不得含秘密 | R1 方案冻结 / 实施 | R1 合同冻结 | D-001 | verified（继承） | — | `01-decision/D-001-metrics-export-contract.md` |
| （继承）I-003 | required | Store / 对象存储 / Job 是否进本波分母 | R1 方案冻结 | R1 合同冻结 | D-001 §7 | verified（继承） | — | `01-decision/D-001-metrics-export-contract.md` §7 |
| （继承）I-004 | required | `/metrics` 或 OTLP 是否进入 `readyz` | R2/R3 方案冻结 | R2/R3 接入前 | D-001 §8（提前闭合） | verified（继承） | — | `01-decision/D-001-metrics-export-contract.md` §8 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-21 | 指标导出合同与配置面冻结（R1） | accepted | `01-decision/D-001-metrics-export-contract.md` |

> 编号在本目标内单调不复用。合同修订（新系列/新键/分母扩展）必须新增 D 记录，不得原地改写已接受决定。
