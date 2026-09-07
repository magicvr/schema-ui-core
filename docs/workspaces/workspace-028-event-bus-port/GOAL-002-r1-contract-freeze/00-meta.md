---
id: GOAL-002-r1-contract-freeze
title: R1 契约冻结（EventBus 端口 / 注册表 / 异步投递 / 错误语义）
status: done
parent: GOAL-001-event-bus-port
created: 2026-09-01
updated: 2026-09-01
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-028-event-bus-port
primary_plan: VP-028-event-bus-port
serves_summary: 承载 VP-028 R1 阶段（判据 #1/#6）：冻结 EventBus 端口契约——类型化注册表 + 可序列化约束（I-028-001）、异步投递与缓冲满阻塞（I-028-002）、handler 吞掉+日志与 panic 隔离（I-028-003）；端口本体落 kernel/eventbus.go + 合同级快测。
---

# GOAL-002 · R1 契约冻结

## 概述

执行 Root 纲领 **R1**：在仓库既有内核端口先例（`kernel.Cache` / `RateLimiter` / `Store` / `ObjectStore` / `MailSender`）之上，冻结 VP-028 **EventBus 运输端口契约**。用户 2026-09-01 P-004 裁决（未采纳编排器建议项 ①②）：注册表 topic→type + 可序列化约束；异步 channel + 缓冲满阻塞；handler 吞掉+日志 + panic 隔离。**合同正文 = GOAL-002 D-002 产物**；端口本体（`apps/api/kernel/eventbus.go`）+ 合同级快测在本目标落地。进程内实现（判据 #2）归 R2；outbox/MQ 接缝与 topic 命名登记（判据 #3/#4/#5）归 R3。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **信息裁决**：I-028-001 / I-028-002 / I-028-003 用户裁决（P-004 required） | **已关门**（2026-09-01 用户裁决：注册表+可序列化 / 异步+缓冲满阻塞 / 吞掉+日志+panic 隔离——D-001） |
| C2 | **合同正文 + 端口落地**：D-002 合同冻结；`kernel/eventbus.go` 实现；合同级快测绿 | **已关门**（D-002 v0.1.0 冻结 + kernel/eventbus.go + eventbus_test.go；`go vet` 0 / `go test ./kernel/...` 绿 / `go build ./kernel/` 通过；2026-09-01） |
| C3 | **审视与关门**：合同自审（self A-001）+ independent（本地 grok build grok-4.6 · high A-002）合并响应；R1 关门、Root 信息台账回写 | **已关门**（A-001 self `pass` 0 required + A-002 grok independent `pass` 0 required；A-003 合并响应 F-001～F-004 全处置；R1 关门 3/3 · Root progress 1/4；2026-09-01） |

`progress` = 已关门检查点数 / 3。当前 **3/3**。

## 成功标准（方向级）

1. 端口契约冻结：EventBus（Register / Publish / Subscribe / Unsubscribe / Stop）按用户裁决冻结、快测可断言（判据 #1）。
2. 类型化按用户裁决落盘：topic→type 注册表；Publish 类型必须匹配已注册类型；负载必须 JSON 可序列化（fail-closed）。
3. 投递按用户裁决落盘：异步（Publish 不等待 handler 完成）；每订阅独立有界缓冲；缓冲满时发布者阻塞；Stop 须 SIGTERM 取消订阅/排空（判据 #6）。
4. 错误语义按用户裁决落盘：handler 无 error 返回（吞掉为结构）；供应商必须 recover panic 并记日志；多订阅独立失败。
5. 未越界：不改 Profile 默认集 / 模块矩阵 / Manifest；不引入 broker 客户端；不实现 outbox；不改 Charter；不解除 Admin typed domain event gated。

## 信息就绪与未知项

与 Root / VP-028 同号镜像。I-028-004 因 I-028-001 选注册表升为 **required**，最晚阶段仍为 R3，本目标不关闭。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-028-001 | required | 事件类型化：注册表 topic→type + 可序列化约束 | 方案冻结 + 判据 #1 | C1 | 用户裁决 | **verified** | — | 2026-09-01 用户裁决：注册表 + JSON 可序列化（D-001） |
| I-028-002 | required | 投递语义：异步 + 缓冲满阻塞 | 判据 #2/#6 | C1 | 用户裁决 | **verified** | — | 2026-09-01 用户裁决：异步 channel + 缓冲满阻塞 + Stop 排空（D-001） |
| I-028-003 | required | handler 错误语义 | 判据 #2 | C1 | 用户裁决 | **verified** | — | 2026-09-01 用户裁决：吞掉+日志 + panic 隔离；handler 无 error 通道（D-001） |
| I-028-004 | required | 事件类型注册权属（业务域 VP vs Admin 功能 VP）；不解除 Admin typed domain event gated | 判据 #4 | R3 | lead 建议 + 用户确认 | 待确认 | 因 I-028-001 选注册表由 non-blocking 升 required | 本目标不关闭 |

## 父目标

- `GOAL-001-event-bus-port`（Root · 纲领 R1）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式：本目标落 kernel 公共面（兼容性门禁）→ **C3 走 cross**：A-001 self + A-002 本地 grok build（grok-4.6 · high）independent。
- R1 合同为本 VP 首波冻结分母；D-002 冻结后实施（R2）与验收（R4）以本合同为准。
