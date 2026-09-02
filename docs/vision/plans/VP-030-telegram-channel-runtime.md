---
doc_type: vision-plan
id: VP-030-telegram-channel-runtime
title: Telegram Bot 通道运行时
status: planned
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace:
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
parent: null
---

# VP-030 · Telegram Bot 通道运行时

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`planned`**（2026-09-02 · v0.1.0 · 0 区） |
| lead_workspace | 未绑定（激活时按惯例 `workspace-030-telegram-channel-runtime`） |
| Vision required | 计划阶段 self = [VRev-065](../reviews/VRev-065-c-end-paid-services-planned-self.md)；激活前另做**架构类** freshness + 激活审视 |
| 组合位置 | **架构分支 · C 端通道**（对标 VP-017 出站邮件：内核端口 + 一方模块 + Admin 设置）。**不是**业务域，**不是**付费产品 |

## 意图

在同进程基座上提供**可扩展的 Telegram Bot 通道运行时**，让业务模块（钱包核销、数字 Offer 等）以注册方式挂接命令/回调，而不是每个 fork 重写 webhook。

对标 VP-017：

1. **入站端口**：HTTP webhook（secret token fail-closed）→ 校验 → 解析 Update → 分发给已注册 handler。
2. **出站端口**：`SendMessage`（首波文本；媒体 gated）。
3. **身份**：把 `telegram_user_id` 经 VP-029 主体接缝登记为 `issuer=telegram` 的 `subject_id`。不创建 `admin.users`。
4. **Admin 设置**：bot token / webhook secret（密钥 fail-closed，不进配置包明文）。
5. **限流**：消费 VP-027 端口，按 webhook IP / `chat_id` / `telegram_user_id` 分桶（R1 冻结桶分母）。

本模块是 **C 端 ingress**，即使尚未激活 VP-031，只要 webhook 对公网开放就有 C 端流量。因此 **激活前必须评估**进程内 RateLimiter 是否覆盖上述桶；可结论「不需要 Redis」，但评估本身不可跳过（H-002 / RT-Q05 的 C 端接入精神适用于通道，不把本 VP 误标成业务域）。

## 首波冻结（退出分母）

| 项 | 本 VP 交付 | 不进本 VP |
|----|-----------|-----------|
| 入站 | HTTPS webhook 路由；`secret_token` 校验 fail-closed；无 token 拒绝启动或拒绝请求（R1 冻结二选一，生产不得明文默认） | long polling 生产路径；多 bot 进程；独立 Bot 微服务 |
| 分发 | 内核级 Update 分发：模块可 `Register` 命令与 callback query；静态编译候选，不热插拔 | 对话状态机 / 场景引擎 / 中间件市场 |
| 出站 | `SendMessage` 文本端口；供应商 = Telegram Bot API HTTP；无 token 时 mock/记录（R1 冻结，对标 mail mock） | 媒体/文件/贴纸；主动广播运营；Bot API 全覆盖 |
| 身份 | 首次可见的 telegram user → VP-029 `GetOrCreate("telegram", id)`；后续 handler 只见到 `subject_id` | Admin JWT 登录；Telegram Login Widget |
| 设置 | Settings 增加 Bot 渠道 tab：token / secret / 试发（可选） | 把 token 写入可导出配置包明文 |
| 模块形态 | `channel.telegram`（建议 id，R1 冻结）：HTTP + 设置 Schema/Nav/Manifest + Persistence（若有本地状态）。**豁免**业务导航。不进 `mvp`/`admin` 默认集 | 一方标准 Admin 六项里的业务列表页 |
| 限流 | 使用点接入 VP-027；桶分母 R1 冻结并写进退出判据 | Redis 实现（仍 gated） |

## 非目标

- Mini App / WebApp、Telegram Login Widget
- Telegram Stars / Payments / 发票
- 对话 FSM、多轮表单框架、自然语言路由
- 具体付费命令（`/buy` `/redeem` `/balance` 由 VP-029/031 注册，本 VP 只提供注册表）
- 独立 Bot 进程、长轮询生产、多 bot 账号（单 token）
- 类目/商品/订单业务
- 改 Charter；重开 VP-017/026/027/028

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-017** | 形态参照（端口 + 设置 + mock 默认）。邮件端口不复用 |
| **VP-003 / VP-004** | 通道模块按「横切 + 设置面」显式豁免业务导航；内核不得 import `channel.telegram` 实现细节 |
| **VP-008 `go`** | 架构类 freshness（pin / 锁 / 迁移 / Profile 默认集 / provenance） |
| **VP-021** | webhook handler 遵守停机 drain；不得在 SIGTERM 后继续调 Bot API 死循环 |
| **VP-027** | 必消费限流端口；激活前评估进程内是否够用并登记路线图位置 |
| **VP-026** | 无硬依赖。Update 去重若用缓存，评估后可选，不预制 Redis |
| **VP-028** | 不把 Bot Update 当领域事件总线；需要 fan-out 时再评估 typed event（仍 gated） |
| **VP-029** | **硬前置**：主体接缝必须已可用（本 VP 不自建用户表） |
| **VP-031** | 业务命令的注册者，不是本 VP 的退出分母 |

## 方向级退出判据

1. **Webhook 合同**：secret 校验 fail-closed；无/错 secret 不可被当成合法 Update；有测试。
2. **分发端口**：命令与 callback 的 Register/Unregister + 分发有测试；未知命令有确定回落（不 panic）。
3. **出站端口**：`SendMessage` 文本可测（mock 供应商）；生产供应商不把 Bot 客户端类型漏进模块公共契约。
4. **身份映射**：同一 `telegram_user_id` 多次 get-or-create 得到同一 `subject_id`；不写 `admin.users`。
5. **设置与密钥**：Admin 可配置 token/secret；密钥 fail-closed；不进配置包明文。
6. **限流评估落盘**：激活前书面评估 VP-027 进程内 limiter 对 webhook/`chat_id`/`telegram_user_id` 是否足够；结论允许「不需要 Redis」，评估不可缺。
7. **边界保持**：未改 Charter；未进 `mvp`/`admin` 默认集；未做 Mini App / Stars / 对话引擎 / 付费命令；未重开历史 VP。
8. **审计闭合**：开放 required finding = 0（或已合法闭合）。

建议 Root 纲领：R1 合同（路由/secret/分发 API/桶分母/mock）→ R2 webhook + 分发 + 身份 → R3 出站 + 设置 + 限流接入 → R4 证据与关门。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-030-001 | 无 token 时：拒绝进程启动 vs 启动但 webhook 503。 | required | 判据 1/3 | R1 | open |
| I-030-002 | Bot API 调用：标准库 HTTP vs 引入第三方 SDK。默认倾向标准库，避免 SDK 泄漏进公共面。 | required | 判据 3 | R1 | open |
| I-030-003 | 限流桶分母：webhook IP / chat_id / telegram_user_id 哪些本波必做。 | required | 判据 6 | R1 | open |
| I-030-004 | 模块 id 最终字符串（建议 `channel.telegram`）。 | non-blocking | 装配 | R1 | open |
| I-030-005 | 设置是否热切换 token（mail 有热切换先例）还是重启生效。 | non-blocking | 判据 5 | R3 | open |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | lead | — | `planned` 0 区；硬前置 VP-029 主体接缝 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-02 | 初创 `planned`：用户确认同进程 + 需要一方 Telegram 通道运行时（通道而非业务域）。Offer 顺延为 VP-031。 |
