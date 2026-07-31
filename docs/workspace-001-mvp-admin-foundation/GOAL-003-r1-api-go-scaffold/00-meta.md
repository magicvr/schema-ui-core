---
id: GOAL-003-r1-api-go-scaffold
title: R1 · Go API 工程骨架
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# GOAL-003 · R1 · Go API 工程骨架

## 概述

在 `apps/api` 建立可本地运行的 Go 服务骨架：目录分层、模块路径、健康检查与构建/运行约定。  
**不**实现业务域 API，**不**宣称协议兼容完成。

结构参考平行仓 `allinme.core-api`（`dev`）的 `cmd/` + `internal/` + `pkg/` 分层；择优移植通用层，禁止整仓拷贝订单/钱包/通知等业务。

## 成功标准

- [ ] `apps/api/go.mod` 存在且模块路径属于本仓约定（非 `allinme.core-api`）
- [ ] `cmd/server`（或等价）可 `go run` / 文档记载的命令启动
- [ ] 至少提供健康探活（如 `/healthz`）可手工或脚本验证
- [ ] 分层目录与 Makefile（或等价）最小命令齐备：`run` / `test` / `build` 之一组
- [ ] README 或约定文档写明默认端口与 env 示例；无业务域路由作为 MVP 默认

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-003-001 | non-blocking | Go 版本下限（平行仓为 1.26；本机/CI 是否对齐） | 脚手架可移植性 | 骨架可运行前 | 探测本机 `go version` 并写入约定 | open | — | 待实施时记录 |
| I-003-002 | non-blocking | 模块 path 最终字符串（如 `github.com/magicvr/schema-ui-core/apps/api`） | go.mod | 首次 `go mod` 前 | 与用户/仓库远程 URL 对齐后写入 | open | — | 待确认远程 canonical |

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 纲领阶段：**R1**。依赖 GOAL-002 路径约定（已由 Root D-004 预锁定为 `apps/api`）。
- 协议目标版本为 `schema-ui-docs@v2.7.0`；平行仓 README 曾出现 page schema `2.4` —— **不得**原样搬协议版本声明。
