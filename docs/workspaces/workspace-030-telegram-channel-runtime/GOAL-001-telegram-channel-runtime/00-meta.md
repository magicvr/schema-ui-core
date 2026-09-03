---
id: GOAL-001-telegram-channel-runtime
title: Telegram Bot 通道运行时
status: active
parent: null
created: 2026-09-03
updated: 2026-09-03
version: 0.1.3
progress: 2/4
plan_refs:
  - VP-030-telegram-channel-runtime
primary_plan: VP-030-telegram-channel-runtime
serves_summary: 同进程 Telegram Bot 通道运行时（架构分支 · C 端 ingress · 对标 VP-017）：webhook + Update 分发端口 + SendMessage 文本 + issuer=telegram 主体映射 + Admin bot 设置
---

# GOAL-001 · Telegram Bot 通道运行时

## 概述

承接 [VP-030-telegram-channel-runtime](../../vision/plans/VP-030-telegram-channel-runtime.md)（active v0.2.0 · [VRev-070](../../vision/reviews/VRev-070-vp030-telegram-channel-runtime-activation.md) self `pass` · 架构类 freshness PASS `b5c39dfb`→`42036a3c`）：交付同进程 **C 端 Telegram 通道运行时**。**对象面**：HTTPS webhook（secret fail-closed）+ 内核级 Update 分发（命令/callback Register）+ `SendMessage` 文本端口 + `GetOrCreateSubject("telegram", id)` + Admin bot 设置。**红线（激活即生效）**：不进 `mvp`/`admin` 默认集；不做 Mini App / Stars / 对话 FSM / 付费命令；不引入独立 Bot 进程 / 长轮询生产 / 多 bot；不把 Bot 用户写入 `admin.users`；密钥 fail-closed；不消耗 RT-Q05 Redis trigger；不重开 VP-017/026/027/028/029。

## 成功标准（对应 VP-030 八条方向级退出判据）

- [x] 判据 #1（Webhook 合同）：secret 校验 fail-closed；无/错 secret 不可被当成合法 Update；有测试——已由 GOAL-003 完成
- [x] 判据 #2（分发端口）：命令与 callback 的 Register/Unregister + 分发有测试；未知命令有确定回落（不 panic）——已由 GOAL-003 完成
- [ ] 判据 #3（出站端口）：`SendMessage` 文本可测（mock 供应商）；生产供应商不把 Bot 客户端类型漏进模块公共契约——R3（R1 冻结 HTTP vs SDK）
- [x] 判据 #4（身份映射）：同一 `telegram_user_id` 多次 get-or-create 得到同一 `subject_id`；不写 `admin.users`——已由 GOAL-003 完成
- [ ] 判据 #5（设置与密钥）：Admin 可配置 token/secret；密钥 fail-closed；不进配置包明文——R3
- [x] 判据 #6（限流评估落盘）：激活前书面评估 VP-027 进程内 limiter 对 webhook/`chat_id`/`telegram_user_id` 是否足够——**已核销**（2026-09-03 · VRev-070 §6：进程内够用、不需要 Redis）。桶分母与 Record/Clear 映射已随 GOAL-002/003 落地。
- [ ] 判据 #7（边界保持）：未改 Charter；未进默认集；未做 Mini App / Stars / 对话引擎 / 付费命令；未重开历史 VP——R4
- [ ] 判据 #8（审计闭合）：开放 required finding = 0（或已合法闭合）——R4

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。`progress` = 已完成纲领阶段 / 4。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 合同冻结（判据 1/2/3/6 + I-030-001/002/003/006）：无 token 启动策略 · Bot HTTP vs SDK · 桶分母 · 请求计数 vs 失败预算映射 · 分发 API / mock | **已关门**（GOAL-002 done 3/3 · D-002 冻结 + E-002 端口代码与快测全绿 + A-001 pass） |
| R2 | webhook + 分发 + 身份（判据 1/2/4 + I-030-007）：secret 路由 · Register 分发 · 主体映射（不要求钱包 HTTP 已启）；**入站三桶限流随 webhook 落地** | **已关门**（GOAL-003 done 3/3 · grok 独立审 A-002 pass/闭合 + Webhook/分发/身份/限流落地） |
| R3 | 出站 + 设置 + 限流核账（判据 3/5 + I-030-005）：SendMessage mock/生产隔离 · Admin bot tab · 入站限流核账（使用点已随 R2） | **进行中**（待立项 GOAL-004） |
| R4 | 证据与关门（判据 7/8；依赖 R1–R3）：证据矩阵 / 越界核账 / 审计闭合 | 待 R1–R3 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-030-001 | required | 无 token 时：拒绝进程启动 vs 启动但 webhook 503。 | 方案冻结 + 退出判据 1/3 | R1 | 用户裁决（R1 前置） | **verified** | — | 2026-09-03：进程可启动；webhook 503；出站 mock（GOAL-002 D-001） |
| I-030-002 | required | Bot API 调用：标准库 HTTP vs 引入第三方 SDK。默认倾向标准库，避免 SDK 泄漏进公共面。 | 退出判据 3 | R1 | 用户裁决（R1 前置） | **verified** | — | 2026-09-03：stdlib `net/http`，不引入 Telegram SDK（GOAL-002 D-001） |
| I-030-003 | required | 限流桶分母：webhook IP / chat_id / telegram_user_id 哪些本波必做。 | 退出判据 6 实施 | R1 | 用户裁决（R1 前置） | **verified** | — | 2026-09-03：三桶全做（GOAL-002 D-001） |
| I-030-004 | non-blocking | 模块 id 最终字符串（建议 `channel.telegram`）。 | 装配 | R1 | lead 建议 + 用户确认 | **verified** | — | 2026-09-03：`channel.telegram`（GOAL-002 D-001） |
| I-030-005 | non-blocking | 设置是否热切换 token（mail 有热切换先例）还是重启生效。 | 退出判据 5 | R3 | 用户裁决或沿用 mail 先例 | open | — | 待确认 |
| I-030-006 | required | 入站 Update 如何映射 VP-027 失败预算：独立 limiter 对每次请求 `Record`（不 `Clear`）当计数器 vs 只 Record secret/parse 失败。须与 I-030-003 同裁决。（V-F114） | 退出判据 6 实施 | R1 | 用户裁决（R1 前置） | **verified** | — | 2026-09-03：每次入站 Record，永不 Clear（GOAL-002 D-001） |
| I-030-007 | non-blocking | 主体 Store 消费路径：直接 import `modules/wallet/subject` vs 抽中性端口。无论哪条，**不得要求** `admin.wallet` HTTP 已启。（V-F115） | 退出判据 4 | R2 | 用户裁决（GOAL-003 D-001） | **verified** | — | 2026-09-03 用户裁决：直接复用 subject.Store，纯 TxRunner 依赖，不依赖 admin.wallet HTTP |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001 已首条落盘，后续按编号递增。

## 备注

- **开区（2026-09-03 · 用户指令）**：VP-030 `planned → active` v0.2.0（VRev-070 self `pass` 0 required · 架构类 freshness PASS `b5c39dfb`→`42036a3c` · 限流评估 = 进程内够用、不需要 Redis · 不暂挂 `go`）；lead `workspace-030-telegram-channel-runtime`。
- 审计模式（D-001 已定）：阶段关门 default self；R2 webhook secret / R4 证据与关门按需 independent（grok build 先例，项目级默认执行路径）。
- freshness 三字段与激活锚点见 D-001：消费候选 = HEAD `42036a3c`；next trigger = 首个 C 端业务域 VP 激活（H-002）或多实例部署评估。
- **R1 信息裁决（2026-09-03）**：I-030-001/002/003/004/006 用户书面全部采纳建议项；合同正文 [GOAL-002 D-002](../GOAL-002-r1-contract-freeze/01-decision/D-002-telegram-channel-contract.md) v0.1.0。入站限流使用点随 R2 webhook，不拆到公开之后。
