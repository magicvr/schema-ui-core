---
id: GOAL-003-r1-api-go-scaffold
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# 决策记录 · GOAL-003

## 信息需求与阶段门禁

| ID | 级别 | 最晚阶段 | 状态 | 阻断 |
|----|------|----------|------|------|
| I-003-001 | non-blocking | 骨架可运行前 | open | 不阻断开工；影响版本声明 |
| I-003-002 | non-blocking | 首次 go.mod 前 | open | 实施时确认 module path |

布局与复用策略服从父目标 [D-004](../GOAL-001-mvp-admin-foundation/01-decision.md)。

## D-001 · 骨架范围与平行仓复用边界

**日期**：2026-07-31  
**状态**：accepted

**决定**：

1. 目标路径：`apps/api/`（Go module 根）。
2. 目录取向：`cmd/server`、`internal/`（config/server/handler 等最小集）、`pkg/`（可复用 envelope/version 等，按需）。
3. 复用：参考 `../allinme.core-api`（本地平行，`dev`）的分层、Makefile、`.env.example`、health 模式；**移植时改名/改 module**，去掉业务域。
4. **明确不搬**：`internal/domain` 中 order/wallet/notification、对应 handler/repository、demo 业务 seed、page schema 业务页（属后续阶段或非目标）。
5. 鉴权/JWT/SQLite 可作为**后续 R4 候选模式**记在备注，R1 不强制完整 auth 闭环。

**为什么**：

- Charter 非目标排除特定业务终端模块；平行仓已含演示业务，整树拷贝会污染 MVP 边界。
- R1 只要可运行基架，账号权限在 R4。

**未选方案**：

- **整仓拷贝再删业务**：快但易残留协议 2.4 声明与业务路由。
- **R1 直接上完整 auth+RBAC**：超出 R1「不实现业务能力」边界，且 `I-PROTO-002` 未就绪。
