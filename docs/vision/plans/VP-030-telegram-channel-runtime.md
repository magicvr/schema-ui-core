---
doc_type: vision-plan
id: VP-030-telegram-channel-runtime
title: Telegram Bot 通道运行时
status: closed
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-030-telegram-channel-runtime
created: 2026-09-02
updated: 2026-09-05
version: 0.3.0
parent: null
---

# VP-030 · Telegram Bot 通道运行时

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`closed`**（2026-09-05 · v0.3.0 · 用户指令授权关门 · lead `workspace-030-telegram-channel-runtime`） |
| lead_workspace | `workspace-030-telegram-channel-runtime`（2026-09-03 开区） |
| Vision required | 计划阶段 self = [VRev-065](../reviews/VRev-065-c-end-paid-services-planned-self.md)；激活就绪 = [VRev-070](../reviews/VRev-070-vp030-telegram-channel-runtime-activation.md) self `pass`；关门就绪 = [VRev-076](../reviews/VRev-076-vp030-telegram-channel-runtime-close-out.md) self `pass`；状态同步复核 = [VRev-079](../reviews/VRev-079-vp030-vp033-closeout-status-sync-response-self.md) self `pass`（响应 [VRev-078](../reviews/VRev-078-vp030-vp033-closeout-status-sync-independent-gpt-sol.md) independent fail，2 个聚合 required 已 fixed；当前 0 required） |
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

## 限流评估（激活前 · 判据 6 · 2026-09-03）

权威全文：[VRev-070](../reviews/VRev-070-vp030-telegram-channel-runtime-activation.md) §6。摘要：

| 项 | 结论 |
|----|------|
| 端口覆盖 | `kernel.RateLimiter` 不透明键可组 webhook IP / `chat_id` / `telegram_user_id` 三桶 |
| 进程内 | 同进程单实例 + `1<<16` 容量 FIFO 驱逐 → **本波够用** |
| Redis | **不需要**；不消耗 RT-Q05 trigger |
| R1 仍须冻 | I-030-003 哪些桶必做；I-030-006 请求计数 vs 失败预算的 Record/Clear 映射 |

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
| **VP-027** | 必消费限流端口；激活前评估已落盘（进程内够用、不需要 Redis）；桶分母与 Record/Clear 映射仍 R1 |
| **VP-026** | 无硬依赖。Update 去重若用缓存，评估后可选，不预制 Redis |
| **VP-028** | 不把 Bot Update 当领域事件总线；需要 fan-out 时再评估 typed event（仍 gated） |
| **VP-029** | **硬前置已交付**（`closed` v0.5.0）：`GetOrCreateSubject("telegram", id)`。消费路径 **不得要求** `admin.wallet` HTTP 已启（V-F109 / V-F115）；主体 Persistence 随编译候选 |
| **VP-031** | 业务命令的注册者，不是本 VP 的退出分母 |
| **VP-033** | Admin 运营台（连接状态 / 入站模式 / 占用位 / 人工 IM）。**不重开本 VP**、不改本 VP 八条判据；单实例 `getUpdates` 有界解禁归 033 |

## 方向级退出判据

1. **Webhook 合同**：secret 校验 fail-closed；无/错 secret 不可被当成合法 Update；有测试。
2. **分发端口**：命令与 callback 的 Register/Unregister + 分发有测试；未知命令有确定回落（不 panic）。
3. **出站端口**：`SendMessage` 文本可测（mock 供应商）；生产供应商不把 Bot 客户端类型漏进模块公共契约。
4. **身份映射**：同一 `telegram_user_id` 多次 get-or-create 得到同一 `subject_id`；不写 `admin.users`。
5. **设置与密钥**：Admin 可配置 token/secret；密钥 fail-closed；不进配置包明文。
6. **限流评估落盘**：激活前书面评估 VP-027 进程内 limiter 对 webhook/`chat_id`/`telegram_user_id` 是否足够；结论允许「不需要 Redis」，评估不可缺。**（本条已由 VRev-070 §6 核销）**
7. **边界保持**：未改 Charter；未进 `mvp`/`admin` 默认集；未做 Mini App / Stars / 对话引擎 / 付费命令；未重开历史 VP。
8. **审计闭合**：开放 required finding = 0（或已合法闭合）。

建议 Root 纲领：R1 合同（路由/secret/分发 API/桶分母/mock）→ R2 webhook + 分发 + 身份（入站三桶随 webhook）→ R3 出站 + 设置 + 限流核账 → R4 证据与关门。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-030-001 | 无 token 时：拒绝进程启动 vs 启动但 webhook 503。 | required | 判据 1/3 | R1 | **verified**（2026-09-03：进程可启动；webhook 503；出站 mock · GOAL-002 D-001） |
| I-030-002 | Bot API 调用：标准库 HTTP vs 引入第三方 SDK。默认倾向标准库，避免 SDK 泄漏进公共面。 | required | 判据 3 | R1 | **verified**（2026-09-03：stdlib `net/http` · GOAL-002 D-001） |
| I-030-003 | 限流桶分母：webhook IP / chat_id / telegram_user_id 哪些本波必做。 | required | 判据 6 | R1 | **verified**（2026-09-03：三桶全做 · GOAL-002 D-001） |
| I-030-004 | 模块 id 最终字符串（建议 `channel.telegram`）。 | non-blocking | 装配 | R1 | **verified**（2026-09-03：`channel.telegram` · GOAL-002 D-001） |
| I-030-005 | 设置是否热切换 token（mail 有热切换先例）还是重启生效。 | non-blocking | 判据 5 | R3 | **verified**（2026-09-03：热切换 · GOAL-004 D-001） |
| I-030-006 | 入站 Update 如何映射 VP-027 失败预算：独立 limiter 对每次请求 `Record`（不 `Clear`）当计数器 vs 只 Record secret/parse 失败。须与 I-030-003 同裁决。 | required | 判据 6 实施 | R1 | **verified**（2026-09-03：每次入站 Record，永不 Clear · GOAL-002 D-001） |
| I-030-007 | 主体 Store 消费路径：直接 import `modules/wallet/subject` vs 抽中性端口。无论哪条，**不得要求** `admin.wallet` HTTP 已启。 | non-blocking | 判据 4 | R2 | **verified**（2026-09-03：直接复用 subject.Store，纯 TxRunner 依赖 · GOAL-003 D-001） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-030-telegram-channel-runtime | GOAL-001-telegram-channel-runtime | lead | 2026-09-03 | 唯一 delivery；Root done；VRev-070 激活 self `pass`；VRev-076 关门 self `pass`；open required = 0；R-009 仅保留 workspace-030 A-009 已接受的 bounded residual |

## 关门记录

- 2026-09-05 · **`active → closed` v0.3.0**（用户当前指令授权；[VRev-076](../reviews/VRev-076-vp030-telegram-channel-runtime-close-out.md) self `pass`；open required = 0）。
- 八条方向级退出判据 #1～#8 全部 verified；workspace-030 [Root GOAL-001](../../workspaces/workspace-030-telegram-channel-runtime/GOAL-001-telegram-channel-runtime/00-meta.md) 为 done，R1～R4 与 R5 均已结项，证据路径已在 VRev-076 逐条登记。
- Root 审计链 [A-008 independent](../../workspaces/workspace-030-telegram-channel-runtime/GOAL-001-telegram-channel-runtime/03-audit/A-008-independent-closure-reaudit.md) pass + [A-009 response](../../workspaces/workspace-030-telegram-channel-runtime/GOAL-001-telegram-channel-runtime/03-audit/A-009-a008-response.md) pass；I-030-001～I-030-007 全部 verified；无开放 required。
- R-009 默认 master key 文件与 DB 同目录仍按 A-009 的用户书面 accepted-residual 保留，复审触发为 KMS/HSM 或密钥管理波次；本次不扩大 residual 范围。
- VRev-078 independent 状态同步 fail 发现的 2 个聚合 required 已由 VRev-079 self response fixed；workspace-030 workspace.md 与 Root 00-meta.md 的现行 VP 投影均已同步为 closed v0.3.0，历史 active 记录不改写。

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-02 | 初创 `planned`：用户确认同进程 + 需要一方 Telegram 通道运行时（通道而非业务域）。Offer 顺延为 VP-031。 |
| 2026-09-03 | 用户指令激活：`planned → active` v0.2.0。VRev-070 self `pass`（0 required）。架构类 freshness PASS（`b5c39dfb`→`42036a3c`，不暂挂 `go`）。限流评估落盘：进程内够用、不需要 Redis。I-030-006/007 增补（V-F114/115 → fixed）。lead `workspace-030-telegram-channel-runtime` 交 `/govern` 开区。 |
| 2026-09-03 | R1 信息裁决（`/govern`）：I-030-001/002/003/004/006 **verified**（用户书面全部采纳建议项）。合同正文 = workspace-030 GOAL-002 D-002 v0.1.0。入站限流使用点随 R2 webhook。I-030-005/007 仍 open。 |
| 2026-09-03 | 边界指针：人工控制台 / 入站模式开关登记为 [VP-033](VP-033-telegram-operator-console.md) `planned`（结构选型 A）；**不**把该意图并入本 VP 分母，**不**重开本 VP。 |
| 2026-09-05 | 用户指令授权 VP-030 关门：VRev-076 self `pass`（八条判据全 verified、open required = 0）；workspace-030 Root done；R-009 继续沿用 A-009 已接受的 bounded residual；`active → closed` v0.3.0。此前 VRev-072 V-F116 与 VRev-075 V-F118 的排序建议由本次关门事实解决。 |
