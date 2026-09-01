---
id: GOAL-018-w17-refresh-token-httponly
doc: decision
status: done
parent: GOAL-001-production-hardening
created: 2026-09-01
updated: 2026-09-01
closed: 2026-09-01
version: 1.0.0
---

# 决策记录 · GOAL-018

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账可放在这里；长决策和独立决策记录放在 `01-decision/D-NNN-<slug>.md`，每条记录必须保持可独立阅读。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | Cookie 属性配置（SameSite/Secure/Path/Domain） | 方案冻结 | S1 | 参考业界最佳实践 + 现有 CORS 配置 | verified | — | D-001: HttpOnly, Secure (dev自适应), SameSite=Lax, Path=/api/auth |
| I-002 | required | 非浏览器环境（移动客户端/CLI）兼容性策略 | 方案冻结 | S1 | 明确 header 回退逻辑 + 文档 | verified | — | D-001: 三层回退 Cookie → Header → Body |
| I-003 | required | token 轮换时的 cookie 更新策略（是否每次 refresh 都更新） | 方案冻结 | S1 | 现有轮换逻辑分析 + 安全最佳实践 | verified | — | D-001: 每次 refresh 更新 cookie，响应仍含 JSON 字段 |
| I-004 | non-blocking | 开发环境 cookie 配置（Secure 属性在 HTTP 下的行为） | 实施 | S2 | 代码实现时确认 | verified | — | D-001: DevMode 检测自动禁用 Secure |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-01 | S1 方案冻结 · httpOnly Cookie 双模式架构 | frozen | [D-001-s1-design-freeze.md](01-decision/D-001-s1-design-freeze.md) |

## 信息就绪状态（更新）

所有 required 信息项（I-001/I-002/I-003）已在 D-001 中 verified，方案冻结完成。等待用户授权进入 S2 实施。
