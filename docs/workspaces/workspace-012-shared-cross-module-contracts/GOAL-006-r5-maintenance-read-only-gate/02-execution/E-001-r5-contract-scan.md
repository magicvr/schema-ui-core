---
id: E-001-r5-contract-scan
goal: GOAL-006-r5-maintenance-read-only-gate
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# E-001 · R5 运行态与写边界现状扫描

## 已核对事实

- API `Config` 与 YAML 当前没有 maintenance、degraded 或 read-only 字段，也没有配置热加载路径。
- composition 将核心路由与所有 Provider route contributions 汇入同一 mux，最终经 `handler.WithJSONRouteErrors` 传给 `server.New`；该 handler 链是统一覆盖边界。
- `/healthz` 固定表达进程存活；`/readyz` 表达数据库与模块图 readiness，不应被运行态静默改写。
- public bootstrap 当前固定 `availability.mode = normal`；Web Host 已验证 `normal / maintenance / upgrade-required / degraded`，但没有 `read-only` mode。
- `admin.system-monitoring` 已提供只读 status endpoint 与页面，可作为受认证的真实状态消费者。
- `RouteContribution.Public` 当前不参与 composition 的运行时 middleware 装配，不能作为豁免或优先级依据。
- 既有错误包络支持本地化、稳定 error code 与 request-id；统一门禁可以复用该契约。

## 结论

I-001 已验证。I-002～I-004 仍需以精确模式/路由/前端投影矩阵关闭，当前不放行实现。
