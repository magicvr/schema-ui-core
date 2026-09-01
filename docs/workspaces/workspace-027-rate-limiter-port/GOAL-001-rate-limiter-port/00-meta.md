---
id: GOAL-001-rate-limiter-port
title: 通用限流器端口
status: active
parent: null
created: 2026-09-01
updated: 2026-09-01
version: 0.1.0
progress: 3/4
plan_refs:
  - VP-027-rate-limiter-port
primary_plan: VP-027-rate-limiter-port
serves_summary: 通用限流器端口（架构分支 · H-002 同进程基座早期化 · 承接 RT-Q05）：RateLimiter 端口 + 滑动窗口内存供应商（演进 loginRateLimiter）+ 7 处使用点迁移 + Redis 接缝声明（不实现）
---

# GOAL-001 · 通用限流器端口

## 概述

承接 [VP-027-rate-limiter-port](../../vision/plans/VP-027-rate-limiter-port.md)（active v0.2.0 · [VRev-062](../../vision/reviews/VRev-062-vp027-rate-limiter-port-activation.md) self `pass` · 架构类 freshness PASS `54fb57e7`→`5744868d`）：交付通用限流器端口。**对象面**：内核级限流端口（与 Cache / Store / ObjectStore / Mail 同级）+ 滑动窗口内存供应商（演进既有 `loginRateLimiter`）+ 7 处使用点完整迁移 + Redis 接缝声明。**红线（激活即生效）**：不预制 Redis（不引入客户端依赖 / **不消耗 RT-Q05 trigger**）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；Redis 轨道共享约定继承 owner 文档 `cache-redis-seam-and-track.md`（单一所有者，不跨区绑 Goal D-001）；W12 D-002 窗口常量（15min/20/`IP|identifier`/`Retry-After`）保持；GOAL-014 账号分层锁定（DB 行锁）显式排除；停机语义继承 VP-021。

## 成功标准（对应 VP-027 七条方向级退出判据）

- [x] 判据 #1（端口契约冻结）：RateLimiter 端口（Allow/Record/Clear/RetryAfterSeconds + key 寻址 + 供应商无关）冻结并可用；快测可断言——R1（2026-09-01：D-002 v0.1.1 冻结 + `kernel/ratelimit.go` 端口落地 + 快测 15 子例绿；A-001 self + A-002 grok independent 双审 pass · 0 required）
- [x] 判据 #2（内存供应商可用）：滑动窗口 + 容量边界 + 驱逐语义实现并有测试（并发、窗口边界、驱逐、RetryAfter 计算）——R2（2026-09-01：`internal/ratelimit` 供应商落地；allow 不注册直查 / FIFO 驱逐 / RetryAfter 不剪枝 / provider 默认容量 65537-key / `-race` 并发 · GOAL-003 done 3/3 · A-001 + A-002 双审 pass）
- [x] 判据 #3（使用点迁移不回归 · 完整分母 V-F099）：7 处构造点（登录 / 验证码生成 / 密码修改 / 自助恢复 / MFA verify 独立桶 / MFA step-up / 邀请接受）全部接入端口；回归证据形态 = 各迁入点既有 handler 测试套件全量通过（`go test ./...` exit 0）+ `newLoginRateLimiter` 0 残留 + W12 D-002 窗口常量保持；GOAL-014 分层锁定显式排除——R2（2026-09-01）
- [x] 判据 #4（Redis 接缝声明落盘）：供应商边界（端口不变）+ 原子窗口语义（INCR + EXPIRE）+ 连接管理约定写入；不引入 Redis 客户端依赖——R3（2026-09-01：owner 短文 v1.1.0 §2.6.1～2.6.5 · GOAL-004 done 3/3 · A-001 + A-002 双审 pass · `go.mod` redis 0）
- [x] 判据 #5（共享约定登记）：Redis 轨道约定（VP-026/027：key 前缀 / 命名空间 / 连接管理 / 测试 harness）继承 owner 文档登记（单一所有者；VP-028 不属 Redis 轨道）——R3（2026-09-01：短文 §3.3 首条 `rl` 登记 · 026 登记义务闭环 · 修订史 v1.1.0）
- [ ] 判据 #6（边界保持）：未改 Charter；未改 Profile 默认集 / 模块矩阵 / Manifest 装配；未预制 Redis 实现；未重开历史 VP——全程
- [ ] 判据 #7（审计闭合）：开放 required finding = 0（或已合法闭合）——R4

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 合同冻结（判据 #1 + I-027-001/003/004）：端口 API 形态（Allow/Record 拆分 vs 内聚 Allow）· RetryAfter 语义 · 窗口语义默认（滑动保持）· key 维度扩展 · 供应商无关面 | **已关门**（2026-09-01 · GOAL-002 `done` 3/3：三信息项用户裁决 · D-002 v0.1.1 合同 + kernel.RateLimiter 端口落地 · A-001 self + A-002 grok independent 双审 pass · 开放 required=0） |
| R2 | 内存供应商 + 使用点迁移（判据 #2/#3 + I-027-002）：演进 `loginRateLimiter`（迁移 vs 并存双轨）· 7 处构造点接入 · 回归（D-001 P1 防暴破 + W12 D-002 窗口常量保持） | **已关门**（2026-09-01 · GOAL-003 `done` 3/3：I-027-002 用户裁决方案 A + `internal/ratelimit` 供应商 + 7 处注入 + `rate_limit.go` 删除 + 全量回归绿 · A-001 self + A-002 grok independent 双审 pass · 开放 required=0） |
| R3 | 接缝与共享约定（判据 #4/#5）：Redis 接缝声明（端口不变 / INCR+EXPIRE / 连接管理）+ Redis 轨道约定继承登记（owner = cache-redis-seam-and-track.md） | **已关门**（2026-09-01 · GOAL-004 `done` 3/3：短文 v1.1.0 §2.6 + §3.3 `rl` 首条登记（026 义务闭环）· 零 Go 变更 · A-001 self + A-002 grok independent 双审 pass · 开放 required=0） |
| R4 | 证据与关门（判据 #7；依赖 R1–R3 ✅）：证据矩阵 / 越界核账 / 审计闭合 | **待启动** |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-027-001 | required | 端口 API 形态：Allow/Record 拆分 vs 内聚 Allow（内部记数）；RetryAfter 语义（现有 `retryAfterSeconds` 演进）。 | 方案冻结 + 退出判据 1 | R1 合同冻结 | 用户裁决（R1 合同冻结前置） | **verified** | — | 2026-09-01 用户裁决：**语义拆分保持**（Allow 不注册 + Record 失败计数 + RetryAfterSeconds + Clear；now 注入；capacity≤0 默认 1<<16）（GOAL-002 D-001 accepted；合同 §1） |
| I-027-002 | required | 既有 `loginRateLimiter` 迁移策略：演进为内存供应商（推荐） vs 保留并存（双轨）；多实例 limiter 实例的 key 维度是否扩展。 | 退出判据 3 | R2 | 用户裁决（R2 前置） | **verified** | — | 2026-09-01 用户裁决：**方案 A 演进为内存供应商 + 全量注入**（GOAL-003 D-001 accepted；key 维度不扩展 / 多实例随 W12 D-002 单实例边界保持；证据 = GOAL-003 E-002 + 全量回归绿） |
| I-027-003 | non-blocking | 窗口语义默认：滑动窗口（现状保持） vs 固定窗口 vs 混合；策略接口是否与缓存 VP-026 共用形态。 | 退出判据 2 | R1 | lead 建议 + 用户确认 | **verified** | — | 2026-09-01 用户确认：**滑动窗口保持 + 策略接口独立**（GOAL-002 D-001 accepted；合同 §3） |
| I-027-004 | non-blocking | 限流 key 维度扩展：是否新增"路由+用户"复合 key（C 端 API 防刷典型）；或留给业务域 VP 自行定义维度。 | 退出判据 1 | R1 | lead 建议 + 用户确认 | **verified** | — | 2026-09-01 用户确认：**本波不新增复合 key**（GOAL-002 D-001 accepted；合同 §2） |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001 已首条落盘，后续按编号递增。

## 备注

- **开区（2026-09-01 · 用户指令）**：VP-027 `planned → active` v0.2.0（VRev-062 self `pass` 0 required · 架构类 freshness PASS `54fb57e7`→`5744868d` 五域零变更 · 不暂挂 `go`）；lead `workspace-027-rate-limiter-port`。
- 审计模式（D-001 已定）：阶段关门 default self；实证门禁（R4 证据 / 关门）可按需 independent（grok build 先例，项目级默认执行路径）。
- freshness 三字段与激活锚点见 D-001：消费候选 = HEAD `5744868d`；next trigger = 首个 C 端业务域 VP 激活或多实例部署评估（H-002）。
- 三端口第二个：Redis 轨道共享约定（key 前缀 / 命名空间 / 连接管理 / 测试 harness）由 workspace-026 交付的 `docs/architecture/cache-redis-seam-and-track.md` owner 文档承载，本区激活为命名空间登记义务触发点之一（细则以短文为准，R3 按短文登记）。