---
id: GOAL-003-r2-impl-and-test
title: R2 实现与测试（shutdown_timeout 配置键 / main 接线 / compose 对齐 / 测试锁）
status: done
parent: GOAL-001-graceful-shutdown-and-connection-drain
created: 2026-08-27
updated: 2026-08-27
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-021-graceful-shutdown-and-connection-drain
primary_plan: VP-021-graceful-shutdown-and-connection-drain
serves_summary: 承载 VP-021 R2 阶段（合同 v0.1.0 §1–§7 的实施）：http.shutdown_timeout 配置键（YAML + env + fail-closed）、main.go 停机预算接线、compose stop_grace_period 对齐、配置测试锁。不改 Job runner / Store / 迁移。
---

# GOAL-003 · R2 实现与测试

## 概述

按合同 v0.1.0（GOAL-002 D-002）实施**本合同波次允许的代码面**：停机预算从 `main.go` 硬编码 10s 改为配置驱动（`http.shutdown_timeout` / `HTTP_SHUTDOWN_TIMEOUT`，默认 10s，非法值 fail-closed，§6）；Compose 显式 `stop_grace_period ≥ shutdown_timeout`（§2 部署对齐，A-001 F-001 承接）；配置测试锁。**不改** Job runner / Store / 迁移台账 / Profile 默认集。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **方案冻结**：实现计划（config.go 字段/yaml/env/校验、main.go 接线、compose、测试矩阵）落盘 D-001 | **已关门**（D-001 planned → accepted，2026-08-27） |
| C2 | **实施**：code diff 落地（config.go / config.default.yaml / configs/config.yaml / .env.example / main.go / compose.yaml / config_test.go） | **已关门**（2026-08-27；`go test ./internal/config/...` 绿） |
| C3 | **测试与关门**：全量 `go test ./...` 绿；自审 A-001；R2 关门 | **已关门**（2026-08-27：全量回归 exit 0；A-001 self `pass` 0 required） |

`progress` = 已关门检查点数 / 3。当前 **3/3**（R2 已关门）。

## 成功标准

1. `http.shutdown_timeout`（YAML）+ `HTTP_SHUTDOWN_TIMEOUT`（env）可解析；缺省 10s；非法值（解析失败 / <=0）→ 启动 fail-closed。
2. `cmd/server/main.go` 停机预算改用配置值（不再是硬编码 10s）。
3. `compose.yaml` 显式 `stop_grace_period: 15s`（≥ 默认 10s 预算）。
4. 配置测试锁覆盖：默认值 / YAML 覆盖 / env 覆盖 / 非法值 fail-closed。
5. 未越界：不改 Job runner、Store、迁移、Profile/Manifest 语义。

## 信息就绪与未知项

无开放 required（I-001/002/003 已于 R1 verified；I-004 属 R3 验收面）。本目标只实施配置面，不引入新信息项。