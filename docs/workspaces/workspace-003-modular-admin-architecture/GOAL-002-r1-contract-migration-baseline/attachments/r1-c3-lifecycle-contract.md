---
id: R1-C3-EVIDENCE
title: R1 C3 模块公共契约与生命周期边界
status: recorded
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-002-r1-contract-migration-baseline
version: 0.1.0
source: self
---

# R1 C3 · 模块公共契约与生命周期边界

## 现状事实

| 项目 | 当前事实 | 证据 |
|------|----------|------|
| Go baseline | API module 声明 `go 1.26`；当前直接依赖没有 `go.uber.org/fx` | `apps/api/go.mod:1-9` |
| Fx implementation | 未发现 Fx import、`fx.App`、`fx.Module`、`fx.Provide`、`fx.Invoke` 或 Fx lifecycle hook | `apps/api` 源码检索；当前仅有集中式 `handler.Register` |
| Startup | `config.Load` → production validation → JWT/seed hash → `store.Open` → `handler.Register` → `server.New` | `apps/api/cmd/server/main.go:23-69` |
| Listen failure | `ListenAndServe` 非 `http.ErrServerClosed` 时记录 `server failed` 并退出 | `apps/api/cmd/server/main.go:71-83` |
| Stop | signal context 等待 SIGINT/SIGTERM，10 秒 timeout 调用 `srv.Shutdown`，Store 通过 defer close | `apps/api/cmd/server/main.go:85-96,56` |
| Health/readiness | `/healthz` 只返回 liveness；`/readyz` 只检查 `st.Ping`，并非模块图/迁移/reconcile 聚合 | `apps/api/internal/handler/health.go:39-74`; `docs/architecture/module-architecture.md:109-114` |

## R1 冻结的公共契约语义

R1 冻结语义字段和边界，不冻结具体 Go type 名称或实现代码。R2 负责把这些语义落成框架无关 API、Fx 组合根和验证实现。

1. **组合根边界**：选用 Uber Fx 作为 Go 组合根的依赖装配与生命周期候选；Fx 只存在于组合根实现，`fx.In`、`fx.Out`、`fx.Option` 等不得进入模块公共 API。当前没有 Fx 依赖，不把候选写成现状事实。
2. **模块描述**：每个模块必须声明稳定且不可复用的 `id`、模块版本、内核 API 兼容范围和显式 `DependsOn`；无依赖也必须写空列表。模块只实现框架无关契约，不修改中央业务注册表。
3. **核心六项**：标准 Admin 功能模块必须贡献 `HTTP`、`Schema`、`Authorization`、`Navigation`、`Manifest`、`Persistence` 六项。`Persistence` 包括迁移与 system-data reconciliation 边界。
4. **按需能力**：`Configuration`、`Lifecycle`、`Observability` 可按需声明；`Lifecycle` 的语义包括启动、就绪、停止钩子。按需能力不能覆盖或降级核心六项。
5. **能力协商**：Profile 展开为显式 `modules.enabled` 后，模块描述和聚合结果声明模块 API 版本、协议版本及前端 `requiredCapabilities`；未知、重复、缺失、循环、未启用依赖、冲突贡献或 capability 不兼容均 fail closed，不静默自动启用或降级未知页面。
6. **启动顺序**：先完成模块注册和依赖图/贡献冲突校验，再按依赖拓扑顺序启动；停止按反向拓扑执行。任何启动失败必须清理已启动资源并保留可定位的 module id。
7. **健康与就绪**：`/healthz` 只表示进程存活；`/readyz` 必须反映模块图、迁移、system-data reconciliation 和必需依赖是否就绪。当前 SQLite ping 不能被误写成目标 readiness。
8. **错误分类**：R1 约定至少区分配置/启动、模块图与依赖、协议/capability、迁移/reconcile、监听运行、就绪、停止/清理类别，并要求错误包含可定位模块/阶段信息；稳定错误类型和 code 在 R2 实现时冻结。

## R2 明确承接

R2 必须选择并验证具体 Fx 版本与 Go 1.26 兼容性，建立框架无关模块 API、Profile 展开、依赖图校验、capability negotiation、模块级 lifecycle hook、聚合 readiness、失败清理和结构化错误分类。R1 不添加依赖、不修改启动代码，不宣称架构完成。

## 检查点结论

C3 的现状与 R1 语义边界已形成，Root I-003 仍为 `open`；只有 R1 freeze/stage-gate 审计及 `/govern` 响应完成后，才可提议 Root 信息项验证或创建 R2 子目标。
