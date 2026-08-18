---
id: D-001-r5-operational-state
goal: GOAL-006-r5-maintenance-read-only-gate
status: proposed
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# D-001 · R5 运行态与统一写边界初始方案

## 已确认方向

1. 运行态采用启动配置的单一来源；首版不新增管理写 API或热切换机制。
2. 统一门禁位于最终 HTTP handler 边界，覆盖核心路由和 Provider 贡献路由；不在各业务 handler 重复实现。
3. 拒绝响应复用既有本地化 JSON error envelope、request-id/correlation 和稳定 catalog code。
4. public bootstrap 与 system-monitoring 投影同一运行态；`/healthz` 保持纯存活语义，`/readyz` 继续表达数据库与模块图 readiness。
5. 不新增 Profile/module ID，不改变 Manifest bytes 或聚合算法；degraded capability 必须先通过既有 Host registry 校验。

## 待冻结契约

- `normal / maintenance / degraded / read-only` 的写允许矩阵。
- 门禁在认证/权限之前或之后执行，以及 login/refresh/logout 等 public mutation 的处理。
- maintenance/read-only 的 HTTP status、error code、retry metadata 与 method 范围。
- `read-only` 在现有 Host availability 不含同名 mode 时的兼容投影。
- `degraded.disabledCapabilities` 的来源、合法集合与后端 route enforcement 边界。

## 未选方案

- 逐 handler 加门禁：无法保证新增贡献路由默认受控，且跨模块重复。
- 修改 Profile 或 Manifest 来表达维护状态：会把运行态与静态装配契约混合，越过 VP-012 边界。
- 首版提供数据库持久化/管理 UI：扩大安全与审计面，不是建立共享门禁所必需。

## 门禁

I-002～I-004 须由精确矩阵关闭并经过 cross 设计审计；D-001 在此之前保持 `proposed`，不得进入 S1/S2 实施。
