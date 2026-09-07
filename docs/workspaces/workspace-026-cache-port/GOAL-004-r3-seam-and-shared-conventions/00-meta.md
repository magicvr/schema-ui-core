---
id: GOAL-004-r3-seam-and-shared-conventions
title: R3 接缝与共享约定（Redis 接缝声明 + 轨道 owner 文档 + mail 迁移评估）
status: done
parent: GOAL-001-cache-port
created: 2026-09-01
updated: 2026-09-01
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-026-cache-port
primary_plan: VP-026-cache-port
serves_summary: 承载 VP-026 R3 阶段（判据 #4/#5 + I-026-004）：Redis 供应商接缝声明（端口不变/供应商边界/连接管理/key 前缀）+ Redis 轨道共享约定 owner 文档（VP-026/027 单一所有者）+ mail cachedAdapter 迁移评估（用户确认：不迁移）+ F-002 义务兑现（fx 容器持有 kernel.Cache 单一实例 + newMux 注入点）。
---

# GOAL-004 · R3 接缝与共享约定

## 概述

执行 Root 纲领 **R3**：按 D-002 合同与 VP-026 冻结要求落盘——① **Redis 接缝声明**（判据 #4：端口不变 / 供应商边界 / 连接管理约定 / key 前缀 `<ns>:<key>` / 不引入客户端依赖）与 **Redis 轨道共享约定**（判据 #5：VP-026/027 单一所有者，不跨区绑 Goal D-001）合入架构短文（`docs/architecture/cache-redis-seam-and-track.md`）；② **I-026-004 mail `cachedAdapter` 迁移评估**（2026-09-01 **用户确认：不迁移**，评估留痕于 attachments）；③ **F-002 义务**（2026-09-01 **用户裁决：fx 容器持有 + newMux 注入点**）——组合根以 `fx.Provide(newCache)` 将单一 `kernel.Cache` 实例挂入 Fx 容器（进程级长生命周期持有 + 构造 eager + fail-closed），`newMux` 依赖注入即为首个消费者显式接入点。**不预制 Redis 实现（RT-Q03 保持 trigger-gated）**。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **信息裁决**：I-026-004（mail 迁移与否）用户确认；F-002 挂载方式用户裁决 | **已关门**（2026-09-01 用户裁决：**不迁移，评估留痕** / **fx 容器持有 + newMux 注入点**——D-001） |
| C2 | **落盘 + 实现**：架构短文（接缝 §2 + 轨道约定 §3）；mail 评估附件；组合根 fx 改造（fx.Provide(newCache) + newMux 注入 + 4 调用点更新）；`go vet`/`go test` 全绿；`go.mod` 无 redis 复核（判据 #4/#5） | **已关门**（2026-09-01：短文 v1.0.0 + 评估 + fx 改造落地；`go vet` 0 / 全模块回归 exit 0 / go.mod+go.sum redis 0 命中 / mail git 空 diff） |
| C3 | **审视与关门**：A-001 self + A-002 grok build（grok-4.6 · high）independent 合并响应；R3 关门、Root/VP 台账回写 | **已关门**（A-001 self `pass` + A-002 grok build independent `pass`（0 required）；A-003 合并响应 3+2 findings 全处置（fixed ×1 · fixed-recording ×2 + 合并）；开放 required = 0；台账回写 F-003 完成；2026-09-01） |

`progress` = 已关门检查点数 / 3。当前 **3/3**（R3 已关门）。

## 成功标准（方向级）

1. **判据 #4（接缝声明落盘）**：供应商边界（端口不变）、连接管理约定、key 前缀/命名空间约定写入架构短文；`go.mod` 无 Redis 客户端（可核对）。
2. **判据 #5（共享约定登记）**：Redis 轨道约定（VP-026/027）单一所有者文档落地（本区为 owner；VP-027 激活时继承；VP-028 不属 Redis 轨道）。
3. **I-026-004 闭合**：mail cachedAdapter 评估留痕 + 用户确认（不迁移）；判据 #2 评估面闭合。
4. **F-002 闭合**：fx 容器持有单一实例（进程级长生命周期）；newMux 注入点 = 首个消费者接入点；无 blank-holder 谎言注释。
5. **边界保持**：未预制 Redis 实现（不消耗 RT-Q03 trigger）；未改端口合同 / Profile 默认集 / Manifest / Charter；mail 行为零漂移。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-026-004 | non-blocking | 既有 mail runtime `cachedAdapter` 是否迁移到端口（版本戳失效语义 vs 通用 TTL） | 判据 #2（评估面） | R3 | lead 评估 + 用户确认 | **verified** | — | 2026-09-01 用户确认：**不迁移，评估留痕**（attachments/mail-cached-adapter-evaluation-2026-09-01.md；D-001） |

（R3 无新 required 信息项；I-026-001/002/003 已 verified。）

## 父目标

- `GOAL-001-cache-port`（Root · 纲领 R3）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式（Root D-001）：阶段关门 default self；R3 涉及组合根改造（生产接线面）→ **C3 走 cross**：A-001 self + A-002 本地 grok build（grok-4.6 · high）independent。
- 架构短文为跨 VP 共享资产（VP-026/027 轨道 owner 文档），落 `docs/architecture/`（不落工作区目标内，不跨区绑 Goal D-001）。