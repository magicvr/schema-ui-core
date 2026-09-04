---
id: VRev-075-vp033-telegram-operator-console-activation
doc_type: vision-review
title: VP-033 激活就绪 · Telegram Bot 人工控制台
source: self
date: 2026-09-04
scope: VP-033-telegram-operator-console 意图 / 退出判据 / I-033-007/008 冻结 / Admin freshness（42036a3c → dd1edade）/ 工作区对齐
verdict: pass
open_required: 0
status: active
created: 2026-09-04
updated: 2026-09-04
parent: null
version: 0.1.0
---

# VRev-075 · VP-033 激活就绪（Telegram Bot 人工控制台）

## 背景与触发

用户先指令“`/vision` 走流程激活 VP-033，然后交 `/govern` 开设工作区”，随后于 2026-09-04 书面“接受建议，继续”。本审视承接 VRev-072 的计划阶段 `pass`，核对激活前 required 信息裁决、Admin 类 freshness、持续程序开放意见投影与工作区对齐。

## 1. 意图、边界与结构选型

**pass**。VP-033 消费 VP-030 已交付的 Telegram runtime，在 Admin 功能分支提供连接状态、互斥 webhook/polling、业务占用位与未绑定人工文本控制台；不重开 VP-030，不把本意图塞入 workspace-030，不进入业务域、默认 Profile、SSE/WebSocket、多 bot 或多实例 polling。新 delivery 工作区与 Root slug 已由用户接受：`workspace-033-telegram-operator-console` / `GOAL-001-telegram-operator-console`。

## 2. I-033-007/008 激活冻结

| id | 用户书面裁决 | 判定 |
|----|--------------|------|
| I-033-007 | 首波不要求关闭 Telegram Privacy Mode；只收录 Telegram 实际投递给 bot 的消息。私聊全文；群内以命令、回复 bot 等 bot 可见更新为准。 | **verified**；判据 5 的“群选项卡”不再暗含群消息全量可见性 |
| I-033-008 | webhook 公网 URL 使用显式配置项，不做运行时/代理头猜测；`setWebhook` URL = 公网 base URL + `/api/channel/telegram/webhook`。本地以可注入 Fake Bot API 验收；真实公网 tunnel/live 为可选验证。 | **verified**；判据 1/2 的 URL 来源与本地验收可判定 |

I-033-009/010 保持 `non-blocking open`，分别最晚在 R1/R3 冻结；它们不阻断激活或 Root 建立。

## 3. Admin 类 freshness（`42036a3c` → `dd1edade`）

候选身份：clean HEAD `dd1edade7e6f741bc2b5abf58bd95f6d3a0d8bfa`；本轮写入前 `git status --short` 为空。

| 域 | 区间事实 | 判定 |
|----|----------|------|
| 协议 pin / provenance | `apps/web/src/protocol/upstream` 区间零变更；现行 provenance 为 `schema-ui-docs` v2.9.0 / `81aa1d8` | ✅ 未改变消费语义 |
| 依赖锁 / compose | `apps/api/go.mod`、`go.sum`、`apps/web/package.json`、`package-lock.json`、`compose.yaml` blob 均未变 | ✅ |
| 迁移台账 | 新增 VP-030 已交付的 Telegram `0066 telegram_config`；Store identity 与 Telegram migration 定向测试通过 | ✅ 变化可追溯至已审结实现层；不是未绑定输入 |
| Profile / 装配 | `channel.telegram` 进入编译候选及显式 config，但未进入 `mvp`/`admin` 默认列表；`dev.cmd` 改为读取显式 config | ✅ 不改默认 Profile 集；VP-033 延续“不进默认集”红线 |
| 关键证据可执行 | `go test -count=1 ./kernel ./internal/store ./internal/composition ./modules/channel/telegram/...` PASS；`go test -count=1 ./internal/docscheck` PASS | ✅ |
| Goal / Vision 意见投影 | Vision open required = 0；workspace-030 Root / R5 required = 0；限定扫描未发现 workspace-009/010 对 VP-033 的当前开放 required | ✅ 不暂挂 VP-008 `go` 消费 |

适用边界：workspace-030 `R-009`（默认 master key 文件与 DB 同目录）保留其原有 `accepted-residual` 与复审触发；本审视只引用该已存在边界，不把它自动扩张为 VP-033 的新风险接受。若 R1/R2 改动密钥存储或生产部署隔离，必须回流 P-004 重新裁决。

## 4. V-F116 与 VP-030 状态

VRev-072 的 V-F116 为 **recommended**：建议激活 VP-033 前先按现行分母关闭 VP-030，以减少两个 Telegram 意图同时 active。用户已确认本次不夹带 VP-030 关门；VP-030 保持 `active`，其 Root 已 `done`。该建议不构成 required 门禁，也不改变 VP-030 的退出判据或历史事实。

## 5. 对齐与开区许可

**pass**。VP-033 `vision_ref = schema-ui-core-admin-foundation@0.4.0` 精确匹配唯一 active Charter；结构为单一 lead delivery 工作区，`primary_plan = VP-033-telegram-operator-console`。工作区建立只表示实现上下文与 Root 路线图就位，不构成任何退出判据已完成。

## Verdict

**pass**（open required = 0）。I-033-007/008 已由用户书面冻结；Admin freshness PASS；可将 VP-033 `planned → active`（v0.2.0）并交 `/govern` 建立 `workspace-033-telegram-operator-console` 与 Root。

## Findings

### 必改（required）

无。

### 建议（recommended）

| id | finding | 状态 |
|----|---------|------|
| V-F118 | VP-030 Root 已 done、VP 仍 active；后续宜另轮 `/vision` 按现行分母关门，避免组合叙事长期并存。 | open recommended；不阻断 VP-033 激活/开区 |

## 声明

本意见为 `/vision` self Review；不冒充 independent。用户已另行确认激活与开区写入。工作区/Goal 生命周期由 `/govern` 与目标台账承接。
