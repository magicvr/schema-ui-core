---
id: GOAL-003-r1-api-go-scaffold
title: R1 · Go API 工程骨架
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.0
---

# GOAL-003 · R1 · Go API 工程骨架

## 概述

在 `apps/api` 建立可本地运行的 Go 服务骨架：目录分层、模块路径、健康检查与构建/运行约定。  
**不**实现业务域 API，**不**宣称协议兼容完成，**不**默认挂业务鉴权中间件。

结构参考平行仓 `allinme.core-api`（`dev`）的 `cmd/` + `internal/` + `pkg/` 分层；择优移植通用层，禁止整仓拷贝订单/钱包/通知等业务。

## 成功标准

- [x] `apps/api/go.mod` 存在且模块路径 = `github.com/magicvr/schema-ui-core/apps/api`（非 `allinme.core-api`）
- [x] `cmd/server` 可 `go run` / `make run` 启动
- [x] 至少提供健康探活 `/healthz` 可手工或脚本验证
- [x] Makefile 至少含 **`run`**（必达）；`build` / `test` 推荐齐备
- [x] README 与 `.env.example` 写明默认端口（建议 `:8080`）与 env；无业务域路由作为 MVP 默认
- [x] 未默认挂业务鉴权中间件为必选（R4 再议）

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-003-001 | non-blocking | Go 版本下限 | 脚手架可移植性 / README 声明 | 骨架可运行前 | 探测本机 `go version` 并写入约定 | **verified** | — | 本机 `go1.26.0`；`go.mod` 写 `go 1.26`；README 声明 |
| I-003-002 | **required** | 模块 path 最终字符串 | 首次 `go.mod` 写入前 | 首次 `go mod` 前 | 与远程 `magicvr/schema-ui-core` 对齐后写入决策 | **verified** | — | D-002：`github.com/magicvr/schema-ui-core/apps/api` |

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 纲领阶段：**R1**。路径 `apps/api` 由 Root D-004 锁定；目录所有权服从 GOAL-002 D-002。
- 协议目标版本为 `schema-ui-docs@v2.7.0`；平行仓 README 曾出现 page schema `2.4` —— **不得**原样搬协议版本声明。
