---
doc_type: goal-decision
id: D-001-workspace-root-establishment
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
status: active
version: 0.1.0
---

# D-001 · 工作区 / Root 建立与开区决策

## 上下文

用户 2026-09-03 指令「/vision 激活vp-030，然后交 /govern 开设工作区」；slug 已确认（按惯例：`workspace-030-telegram-channel-runtime` / `GOAL-001-telegram-channel-runtime`，沿用 VP-013～029）。激活门禁已满足：VRev-070 self `pass`（0 required；V-F114/115 → 开区事务内 fixed）+ 架构类轻量 freshness PASS（`b5c39dfb` → `42036a3c`，不暂挂 `go`）+ 限流评估落盘（进程内够用、不需要 Redis）。

## 决策

| 项 | 决定 |
|----|------|
| 工作区 | `workspace-030-telegram-channel-runtime`（canonical `docs/workspaces/workspace-030-telegram-channel-runtime/`） |
| Root | `GOAL-001-telegram-channel-runtime`（`parent: null`；primary_plan = `VP-030-telegram-channel-runtime`） |
| 愿景角色 | `delivery`（不改变 Charter primary workspace） |
| 纲领阶段 | R1 合同冻结 → R2 webhook+分发+身份 → R3 出站+设置+限流接入 → R4 证据与关门（串行；阶段内可并行子目标） |
| 审计模式 | 阶段关门 default **self**；R2 webhook secret fail-closed 与 R4 证据/关门实证门禁按需 **independent**（grok build 先例 · 项目级默认执行路径） |
| 红线 | 不进 `mvp`/`admin` 默认集；不做 Mini App / Stars / 对话 FSM / 付费命令；不引入独立 Bot 进程 / 长轮询生产 / 多 bot；不把 Bot 用户写入 `admin.users`；密钥 fail-closed、不进配置包明文；不消耗 RT-Q05 Redis trigger；不重开 VP-017/026/027/028/029；内核不得 import `channel.telegram` 实现细节 |

## freshness 三字段（VRev-070 · 先例执行惯例）

| 字段 | 值 |
|------|-----|
| consumer_vp | `VP-030-telegram-channel-runtime`（vision_ref `schema-ui-core-admin-foundation@0.4.0`） |
| last_freshness_review_at | 2026-09-03（`b5c39dfb` → `42036a3c` · 架构类轻量 PASS · 协议 pin / 依赖锁 / Profile 默认集 / provenance 零变更） |
| next_freshness_review_trigger | 首个 C 端业务域 VP 激活（H-002 发现机制）或 多实例部署评估 |

## 限流评估投影（判据 6 · 已核销）

权威：[VRev-070](../../../../vision/reviews/VRev-070-vp030-telegram-channel-runtime-activation.md) §6。

- 进程内 `kernel.RateLimiter` 不透明键可覆盖 webhook IP / `chat_id` / `telegram_user_id`。
- 本波 **不需要 Redis**；不消耗 RT-Q05 trigger。
- R1 仍须冻结 I-030-003（哪些桶必做）与 I-030-006（请求计数 vs 失败预算的 Record/Clear 映射）。

## 未选方案

- 不与 VP-031 合并（通道运行时不得被业务域拖死；关门独立）。
- 不把本区挂到 workspace-017（VP-017 已 closed，禁止重开吸收本意图；邮件端口不复用）。
- 不把本区挂到 workspace-029（主体接缝只消费、不重开钱包 VP）。
- 不在本波引入 Redis / 独立 Bot 进程 / Mini App。
