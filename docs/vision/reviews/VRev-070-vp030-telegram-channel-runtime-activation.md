---
id: VRev-070-vp030-telegram-channel-runtime-activation
doc_type: vision-review
title: VP-030 激活就绪 · Telegram Bot 通道运行时（架构 · C 端通道）
source: self
date: 2026-09-03
scope: VP-030-telegram-channel-runtime 意图 / 退出判据 / 非目标 / P-005 / 架构类 freshness（`42036a3c`）/ 限流评估（判据 6）/ VP-029 硬前置
verdict: pass
open_required: 0
status: active
created: 2026-09-03
updated: 2026-09-03
parent: null
version: 0.1.0
---

# VRev-070 · VP-030 激活就绪（Telegram Bot 通道运行时）

## 背景与触发

用户 2026-09-03 指令：「/vision 激活vp-030，然后交 /govern 开设工作区」。VP-030（架构分支 · C 端 Telegram 通道运行时 · 对标 VP-017）经 [VRev-065](VRev-065-c-end-paid-services-planned-self.md) 计划阶段 self `pass`（0 required；**不是**激活许可），本次为激活就绪审视：意图/退出判据/非目标/P-005 + **架构类 freshness** + **限流评估落盘**（判据 6）+ 硬前置 VP-029 主体接缝。

## 审视要点

### 1. 意图与退出判据可判定性

**pass**。VP-030 v0.1.0 八条退出判据（Webhook secret fail-closed / 分发 Register / SendMessage mock / 身份映射 / 设置与密钥 / 限流评估落盘 / 边界保持 / 审计闭合）均可核验。本审视即判据 6 的激活前评估（见 §6）；其余判据属 R1–R4 工作区证据，不使方向不可判定。

### 2. 非目标与红线

**pass**。明文排除 Mini App / Stars / 对话 FSM / 付费命令（归 VP-029/031 注册者）/ 独立 Bot 进程 / 长轮询生产 / 多 bot / 类目商品订单 / 改 Charter / 重开 VP-017/026/027/028。模块建议 id `channel.telegram` **不进** `mvp`/`admin` 默认集；业务导航豁免。不消耗 RT-Q05 Redis trigger。

### 3. 信息需求（P-005）

**pass**。I-030-001（无 token 启动策略）/ I-030-002（标准库 HTTP vs SDK）/ I-030-003（桶分母）required → 最晚 R1；I-030-004（模块 id）/ I-030-005（token 热切换）non-blocking。均未伪装已知。本审视新增 I-030-006（请求计数 vs 失败预算映射，required · R1）与 I-030-007（主体 Store 消费路径，non-blocking · R2），见 findings 响应。

### 4. 硬前置 VP-029 主体接缝

**pass**。VP-029 已于 2026-09-02 **`closed` v0.5.0**（[VRev-069](VRev-069-vp029-wallet-prepaid-instrument-r5-close-out.md)）。代码面：`apps/api/modules/wallet/subject.Store.GetOrCreateSubject(ctx, issuer, externalID, now)` 已交付；issuer=`telegram` 幂等与跨 issuer 隔离有测试。通道后续 handler 只见到 `subject_id`，不写 `admin.users`。V-F109 约束仍有效：查询/get-or-create **不依赖** `admin.wallet` HTTP 已启（见 V-F115）。

### 5. 架构类轻量 freshness（`b5c39dfb` → `42036a3c`）

**PASS**，不暂挂 `go`：

| 域 | 变更 | 判定 |
|----|------|------|
| 协议 pin / provenance（`apps/web/src/protocol/upstream`） | 零变更 | ✅ |
| 依赖锁（go.mod / go.sum / package.json / lockfiles） | 零变更 | ✅ |
| 迁移台账 | `modules/wallet/migration/migration.go` + `migrate_test.go` | ✅ 区间 = VP-029 已审结目（subjects / vouchers 0064） |
| Profile 装配（`kernel/profile.go`） | `mvp`/`admin` **默认模块 ID 列表零变更**；`BuiltinModules` 中 `admin.wallet` 贡献键扩展（voucher / me/redeem 路由与权限） | ✅ 属 VP-029 已审结目；**未**把新模块塞进默认集 |
| 区间代码变更 | wallet 主体接缝 + 预付凭证 + 我的钱包自助核销（VP-029 R1–R5 已关门） | ✅ 不涉及内核端口面以外新契约 / Store 方言 / Manifest 装配语义 |

消费候选 HEAD `42036a3c`（workspace-029 结项 commit）。本 VP 属架构分支、**不是**业务域 VP，H-002「业务域 VP 激活前确认同进程主要形态」发现机制不适用；H-002 仍为 Charter 冻结假设（用户 2026-09-02 书面确认同进程，本波无反证）。

### 6. 限流评估落盘（判据 6 · RT-Q05 C 端 ingress 义务）

**结论：进程内 RateLimiter 覆盖本波桶空间；本波不需要 Redis。评估不可缺，现已落盘。**

| 项 | 评估 |
|----|------|
| 端口 | `kernel.RateLimiter` + `RateLimiterProvider`（VP-027 `closed` v0.3.0）。键为不透明字符串；调用方可组 `tg:webhook:{ip}` / `tg:chat:{chat_id}` / `tg:user:{telegram_user_id}`。 |
| 进程内是否够用 | H-002 同进程单实例。`DefaultRateLimiterCapacity = 1<<16`；`Record` 达容量 FIFO 驱逐，喷雾 distinct key 不能撑爆 map。Webhook 流量与现有登录/核销桶同进程、可独立 `NewRateLimiter` 实例，不共用失败计数。 |
| Redis | **不需要**。RT-Q05 Redis 实现仍 trigger-gated；本 VP **不消耗** trigger（触发仍是多实例 **或** C 端**业务域**模块接入且进程内不够用）。VP-030 是通道运行时，不是业务域。 |
| 语义缺口 | 端口是 **失败预算**（Allow 无副作用；Record 记失败；Clear 成功清零）。C 端 webhook 洪水是 **请求计数**。R1 必须冻结映射（I-030-006）：建议独立 limiter 实例对每次入站 Update `Record`、不 `Clear`，把滑动窗口当计数器；secret 失败可另桶。此缺口不使「能否覆盖桶」不可判定，但不可静默当成登录限流照搬。 |
| 桶哪些本波必做 | 仍属 I-030-003（R1）。评估只保证端口**能**覆盖三桶，不替 R1 选定分母。 |

### 7. 组合对齐

**pass**。VP-030 `vision_ref` = `schema-ui-core-admin-foundation@0.4.0` 精确匹配现行 Charter；roadmap 组合表第 30 行 / RT-M03 / RT-Q05 C 端评估注记为本波同步对象；lead `workspace-030-telegram-channel-runtime` 沿用 VP-013～029 slug 惯例。三分支并行：VP-029 已 closed，允许本架构 VP 为唯一 active 交付；VP-031 保持 planned。

## Verdict

**pass**（0 required）。

VP-030 意图/判据/非目标/信息需求已就绪，硬前置主体接缝已交付，架构类 freshness PASS（不暂挂 `go`），限流评估已落盘（进程内够用、不需要 Redis）。可激活并交 `/govern` 开区。

## Findings

### 必改（required）

无。

### 建议（recommended）

| id | finding | 建议 | 状态 |
|----|---------|------|------|
| V-F114 | VP-027 端口是失败预算，webhook 需要请求计数；若不在 R1 冻结映射，实施易照搬登录 `Allow`/`Record`/`Clear` 而放空洪水。 | 增补 I-030-006 required（最晚 R1）：入站 Update 的 Record/Clear 映射；可与 I-030-003 桶分母同裁决。 | **fixed**（激活事务内写入 VP-030 / Root 信息表） |
| V-F115 | 主体 Store 现位于 `modules/wallet/subject`。若通道 import 绑死「钱包 HTTP 已启」，无钱包 Profile 的 Bot 会断（V-F109 再现）。 | 增补 I-030-007 non-blocking（最晚 R2）：消费路径（直接 import vs 中性端口）须 **不要求** `admin.wallet` HTTP 已启；Persistence 随编译候选，不随 Profile 过滤。 | **fixed**（激活事务内写入 VP-030 相邻表 / Root 信息表） |

## Finding 响应（2026-09-03 · `/vision` 激活事务内）

| id | 路径 | 说明 |
|----|------|------|
| V-F114 | **fixed** | VP-030 P-005 增补 **I-030-006 required（最晚 R1）**；投影 lead Root 信息表 |
| V-F115 | **fixed** | VP-030 相邻 VP 表 + **I-030-007 non-blocking（最晚 R2）**；投影 lead Root 信息表 |

原 verdict（pass）与 finding 原文未改写；本响应为 append-only 补充。

## 声明

本意见不直接修改 Charter / VP / Goal status。required finding 的响应由 `/vision` 追加在本报告中；原 verdict 与 finding 原文不得改写。
