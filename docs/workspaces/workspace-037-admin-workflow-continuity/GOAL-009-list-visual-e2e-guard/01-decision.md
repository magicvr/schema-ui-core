---
id: GOAL-009-list-visual-e2e-guard-decisions
doc: decision
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 0.2.0
---

# 决策台账 · GOAL-009-list-visual-e2e-guard

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-009-001 | required | e2e 挂具本机可跑性 | C1/C2 | C1 | 既有规格冒烟 + 环境检查 | verified | 2026-09-18 | `E-001` |
| I-009-002 | required | 目标页与两 profile 可达性 | C1/C2 | C1 | schema 检索 + 两 profile 实跑探测 | verified | 2026-09-18；roles 页，mvp/admin 均可达 | `E-001` |
| I-009-003 | required | 合同真实几何基线 | C2 | C2 | 真实浏览器采集 1440/700 盒模型与计算样式 | verified | 2026-09-18 | `E-001` |
| I-009-004 | non-blocking | 是否引入像素快照测试 | 范围外 | C2 | 用户未要求；与关系型断言取向冲突 | deferred | 触发：用户明确要求像素级回归 | `D-001` 未选方案 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-18 | 守卫范围、目标页与断言风格 | accepted | [D-001-guard-scope-and-assertion-style.md](01-decision/D-001-guard-scope-and-assertion-style.md) |

## 当前投影

- 守卫目标页为 `/roles`；规格为 `apps/web/e2e/list-visual-surface.spec.ts`。
- 断言为关系型（相等/贴合/确实隐藏），token 断言写成"等于解析出的 `--control`"，避免 token 调优产生噪音。
- 必须在 `mvp` 与 `admin` 两 profile 下通过；CI 现有 `browser-e2e` job 自动纳入，不新增 job。
- 不修改列表实现（除非守卫暴露真实缺陷）；不引入像素快照测试。
