---
id: GOAL-046-w34-batch-export-availability
doc: decision
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# 决策记录 · GOAL-046（W34）

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-046-001 | required | 「未找到」的确切来源 | 方案 | C1 | operator config 原样复现 + `curl -i` 读状态与 `messageKey`；admin preset 对照 | **verified** | — | `D-001` §1 |
| I-046-002 | required | 不可用时隐藏 vs 可见禁用 | 实施 | C1 | 复核 W33 `D-001` §3 冻结取向 + page-actions 高度契约 | **verified** | — | `D-001` §3 |
| I-046-003 | non-blocking | 既有覆盖是否触及「节点已发布但路由未挂载」 | 验收 | C2 | 全仓检索批量导出断言 | **verified** | — | `D-001` §5 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-19 | W34 冻结：operator config 补齐 + 触发面可用性门禁（fail-open 可见禁用） | accepted | `01-decision/D-001-w34-availability-freeze.md` |
