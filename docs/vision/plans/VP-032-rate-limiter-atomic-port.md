---
doc_type: vision-plan
id: VP-032-rate-limiter-atomic-port
title: 限流器端口原子化（AllowRecord）
status: planned
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace:
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
parent: null
---

# VP-032 · 限流器端口原子化（AllowRecord）

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`planned`**（2026-09-03 · v0.1.0 · 0 区 · 用户裁决登记） |
| lead_workspace | 未绑定（激活时按惯例 `workspace-032-rate-limiter-atomic-port`） |
| Vision required | 计划阶段审视（self = VRev-071）；激活前须：架构类 freshness、`/vision` 正式冻结退出分母后交 `/govern` 开区 |
| 组合位置 | **架构分支** · VP-027 后续端口语义强化（不重开已 closed 的 VP-027 关门事实，只承接其 residual R-007） |

## 意图

消除 `kernel.RateLimiter` 的 **Allow/Record 两次调用之间的 TOCTOU 窗口**（GOAL-001 A-008 R-007 residual）：当前 webhook/auth/recovery/captcha/mfa/wallet 等 10+ 使用点均为「先 `Allow` 再 `Record`」，两调用非原子，并发下预算可被穿透（进程内单实例，µs 级窗口）。

本 VP 在端口层新增**原子方法**（如 `AllowRecord(key, now) bool`，check+record 一次加锁完成），并迁移全部使用点；内存供应商实现原子语义，Redis 接缝（RT-Q05）保持 trigger-gated。

## 首波冻结（退出分母）

| 项 | 本 VP 交付 | 不进本 VP |
|----|-----------|-----------|
| 端口 | `kernel.RateLimiter` 新增原子 `AllowRecord`（或等价签名），`Allow`/`Record` 保留兼容 | 删除现有 Allow/Record（向后兼容保留） |
| 供应商 | `Memory` 实现原子 check+record（单锁内完成）；`RetryAfterSeconds`/`Clear` 语义不变 | Redis 实现（RT-Q05 仍 trigger-gated） |
| 使用点迁移 | 全仓 10+ 处 Allow→Record 调用点迁移到 `AllowRecord`（webhook/auth/recovery/captcha/mfa/wallet/invites 等） | 其它内核端口变更；VP-027 关闭事实改写 |
| 测试 | 并发穿透回归测试（两个 goroutine 同时 Allow+Record 不得超预算）+ 各使用点行为等价测试 | 性能基准 |
| Profile | 不改 Profile 默认集 | 改变装配红线 |

## 非目标

- 不重开 / 改写已 closed 的 VP-027 关门事实
- 不实现 Redis / 分布式限流（RT-Q05 触发条件不变）
- 不改其它内核端口
- 不把 R-007 之外的 recommended（R-004/R-009）卷入

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-027** | 承接其关门后 residual R-007；端口语义强化，不重开关门 |
| **VP-030** | 消除的 TOCTOU 直接影响 telegram webhook 三桶限流 |
| **VP-009** | 共享基架安全程序正交；本 VP 是端口级修复波 |
| **RT-Q05** | Redis 实现仍 trigger-gated；本 VP 只做内存供应商原子化 |

## 方向级退出判据

1. **原子性**：`AllowRecord` 在并发下 check+record 原子，无穿透窗口（有并发回归测试）。
2. **行为等价**：全部使用点迁移后单请求语义与 Allow→Record 等价（测试覆盖）。
3. **兼容**：`Allow`/`Record` 保留（非破坏性），文档标注新方法为推荐路径。
4. **边界保持**：未重开 VP-027；未实现 Redis；未改 Profile 默认集。
5. **审计闭合**：开放 required finding = 0（或已合法闭合）。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-032-001 | `AllowRecord` 精确签名与返回值语义（bool 是否足够，是否需返回剩余额度）。 | required | 判据 1 | 方案冻结 | open（激活前 `/vision` 裁决） |
| I-032-002 | 是否所有使用点都应迁移（如 Clear-on-success 调用点是否需要原子变体）。 | required | 判据 2 | 方案冻结 | open（激活前 `/vision` 裁决） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | — | — | 未绑定（planned · 0 区） |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-03 | 用户书面裁决（GOAL-001 A-008 R-007 处置）：新建 VP 下一波做端口原子化，承接 `kernel.RateLimiter` Allow/Record TOCTOU residual。登记 `planned` v0.1.0（0 区），退出分母草案待 `/vision` 正式冻结。 |
