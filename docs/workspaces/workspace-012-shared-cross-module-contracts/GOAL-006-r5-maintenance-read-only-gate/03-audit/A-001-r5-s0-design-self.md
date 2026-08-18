---
id: A-001-r5-s0-design-self
goal: GOAL-006-r5-maintenance-read-only-gate
source: self
verdict: pass
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# A-001 · R5 S0 设计自审

## 审计头

| 项 | 值 |
|----|----|
| source | self |
| scope | R5 S0：现状扫描、D-001/D-002、I-001～I-005；统一写边界与 Host 投影 |
| verdict | pass |
| required findings | 0 |

## 核对结论

1. composition 将核心和 Provider 路由汇入同一 mux，`WithJSONRouteErrors` 可保留 404/405；最终 server handler 链能覆盖所有已注册贡献路由。
2. request-id/CORS 先于门禁，allowlist 内的认证生命周期和强制改密保留既有认证错误；业务写在受控态 fail closed，不暴露权限差异。
3. `/healthz`、`/readyz` 的职责与现有架构边界清楚，system-monitoring 的 status envelope 可追加字段而不改变既有字段语义。
4. Host registry 已登记 `form.controls.readonly`；把 read-only/degraded 投影为既有 `degraded` + 已知 capability 避免扩展上游协议 mode。
5. 配置只在启动读取且不改变 Profile/Manifest，满足 VP-012 的静态装配不变边界。

## Findings

| ID | 等级 | finding | disposition |
|----|------|---------|-------------|
| F-001 | recommended | S1 必须用配置加载测试证明 YAML/env 优先级与非法 mode fail closed。 | implementation gate：S1 |
| F-002 | recommended | S2 必须对核心和 Provider 路由各取一个真实写路径，并验证未知路径仍为 404/405。 | implementation gate：S2 |

## 结论

Self 设计审计通过，开放 required = 0；I-002～I-004 与 D-002 仍等待 cross independent 复核后才可放行 S1。
