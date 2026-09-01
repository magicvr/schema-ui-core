---
id: GOAL-001-cache-port
title: 通用缓存端口
status: active
parent: null
created: 2026-08-31
updated: 2026-09-01
version: 0.3.0
progress: 2/4
plan_refs:
  - VP-026-cache-port
primary_plan: VP-026-cache-port
serves_summary: 通用缓存端口（架构分支 · H-002 同进程基座早期化 · 承接 RT-Q03）：Cache 端口 + 绝对/滑动过期 + 可插拔策略 + 内存供应商（默认）+ Redis 接缝声明（不实现）
---

# GOAL-001 · 通用缓存端口

## 概述

承接 [VP-026-cache-port](../../vision/plans/VP-026-cache-port.md)（active v0.2.0 · [VRev-060](../../vision/reviews/VRev-060-vp026-cache-port-activation.md) self `pass` · 架构类 freshness PASS `055da2fd`→`54fb57e7`）：交付通用缓存端口。**对象面**：内核级缓存端口（与 Store / ObjectStore / Mail 同级）+ 双过期策略 + 可插拔策略接口 + 内存供应商 + Redis 接缝声明。**红线（激活即生效）**：不预制 Redis（不引入客户端依赖 / **不消耗 RT-Q03 trigger**）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；Redis 轨道共享约定（VP-026/027）单一所有者登记；停机语义继承 VP-021。

## 成功标准（对应 VP-026 八条方向级退出判据）

- [x] 判据 #1（端口契约冻结）：Cache 端口（Get/Set/Delete + TTL + 命名空间 + 并发安全）供应商无关、快测可断言——R1（2026-09-01：合同 D-002 冻结 + `kernel/cache.go` + 快测 33 例绿）
- [x] 判据 #2（双策略 + 可插拔）：绝对过期 + 滑动过期 + 策略接口（含自定义策略测试样例）——R2（2026-09-01：Absolute/Sliding 专测 + nextMidnightPolicy 样例 · GOAL-003 done 3/3）
- [x] 判据 #3（内存供应商可用）：有界容量 + TTL 清理 + 驱逐 + 并发边界测试——R2（2026-09-01：**进程总预算**（用户裁决）+ 全局 FIFO + 惰性清理 + `-race` · GOAL-003 done 3/3）
- [ ] 判据 #4（Redis 接缝声明落盘）：供应商边界（端口不变）+ 连接管理约定 + key 前缀/命名空间约定；`go.mod` 无 Redis 客户端——R3
- [ ] 判据 #5（共享约定登记）：Redis 轨道约定（VP-026/027）单一所有者文档落地（本区为 owner）——R3
- [x] 判据 #6（停机语义）：后台清理协程（若选）声明 SIGTERM 排空；否则惰性清理——R1/R2（2026-09-01：I-026-002 裁决惰性清理 · 无新生命周期 · 合同 §5）
- [ ] 判据 #7（边界保持）：未改 Charter；未改 Profile 默认集 / Manifest；未预制 Redis；未重开历史 VP——全程
- [ ] 判据 #8（审计闭合）：开放 required finding = 0（或已合法闭合）——R4

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 合同冻结（判据 #1/#6 边界）：Cache 端口 API 形态（I-026-001）· TTL/清理语义（I-026-002）· 命名空间形态（I-026-003）· 策略接口形态 | **已关门**（2026-09-01 · GOAL-002 done 3/3：三信息项用户裁决 · 合同 D-002 v0.1.1 · 端口落地 · A-001 self + A-002 grok independent 双审 pass · 开放 required=0） |
| R2 | 内存供应商 + 双策略（判据 #2/#3）：有界 + TTL 清理 + 驱逐 + 并发安全 | **已关门**（2026-09-01 · GOAL-003 done 3/3：FIFO 用户裁决 · **进程总预算**（A-002 F-001 用户裁决）· internal/cache 21 测试（-race）· config 键 · A-001 self pass + A-002 grok independent conditional→fixed · 开放 required=0） |
| R3 | 接缝与共享约定（判据 #4/#5 + I-026-004）：Redis 接缝声明 + Redis 轨道约定 owner 文档 + mail 迁移评估（依赖 R2 ✅） | 计划 |
| R4 | 证据与关门（判据 #8；依赖 R1–R3） | 计划 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-026-001 | required | Cache 端口 API 形态：Go 泛型 vs `[]byte` vs 结构化值；零值/未命中语义。 | 方案冻结 + 退出判据 #1 | R1 | 用户裁决（R1 合同冻结前置） | **verified** | — | 2026-09-01 用户裁决：**`[]byte` 负载 + 非泛型端口 + 类型化封装**（GOAL-002 D-001 accepted；合同 §1） |
| I-026-002 | required | TTL 清理语义：惰性（读时清理） vs 后台协程清理；边界与容量来源。 | 退出判据 #3/#6 | R1 | 用户裁决（R1 合同冻结前置；停机语义随选） | **verified** | — | 2026-09-01 用户裁决：**惰性清理 + 配置化容量驱逐**（GOAL-002 D-001 accepted；合同 §5/§6） |
| I-026-003 | non-blocking | 命名空间 / key 前缀约定：模块 ID 前缀 vs 独立命名空间参数。 | 退出判据 #1/#4 | R1 | lead 建议 + 用户确认 | **verified** | — | 2026-09-01 用户确认：**显式命名空间 scoped 视图**（GOAL-002 D-001 accepted；合同 §2） |
| I-026-004 | non-blocking | 既有 mail runtime `cachedAdapter` 是否迁移到端口（评估，不强制；版本戳失效语义可能不匹配通用 TTL）。 | 退出判据 #2 | R3 | lead 评估 + 用户确认 | 待确认 | — | — |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001 已首条落盘，后续按编号递增。

## 备注

- 审计模式（D-001 已定）：阶段关门 default self；实证门禁（R4 证据 / 关门）可按需 independent（grok build 先例，项目级默认执行路径）。
- freshness 三字段与激活锚点（VRev-060）：消费候选 `55da2fd` 不适用——本 VP 首次激活，消费候选 = HEAD `54fb57e7`；next trigger = 首个 C 端业务域 VP 激活或多实例部署评估（H-002）。
- I-026-001 / I-026-002 为 R1 合同冻结前置用户裁决点（P-004）。