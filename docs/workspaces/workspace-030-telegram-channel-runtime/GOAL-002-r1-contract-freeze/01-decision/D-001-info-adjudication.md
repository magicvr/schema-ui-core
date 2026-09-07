---
doc_type: goal-decision
id: D-001-info-adjudication
parent: GOAL-002-r1-contract-freeze
date: 2026-09-03
status: accepted
version: 0.1.0
---

# D-001 · 信息裁决：I-030-001 / 002 / 003 / 004 / 006 + 分发/mock 包（2026-09-03 用户裁决）

## 上下文

R1 合同冻结前置信息项（P-005 / P-004）。编排器按 VP-017 邮件先例、VP-027 失败预算端口语义、VRev-070 §6（V-F114）与 VP-030 建议包提出带建议的选项。2026-09-03 经用户书面裁决**全部采纳建议项**。

## 裁决记录

| ID | 级别 | 选项（用户所见） | 裁决 |
|----|------|------------------|------|
| I-030-001 | required | ① **进程可启动；webhook 拒绝入站（503）；出站走 mock** ② 模块已启用且无 token → 拒绝整个进程启动 ③ 进程启动且 webhook 仍接 Update（仅出站 mock） | **采纳①**：对标 VP-017 未配生产渠道仍能启动；入站 fail-closed；生产不得明文默认 token。 |
| I-030-002 | required | ① **stdlib `net/http`，不引入 Telegram SDK** ② 引入第三方 Bot SDK（类型仍不得漏进公共面） | **采纳①**：首波只有 webhook 接收 + SendMessage 文本；内核端口只暴露薄类型。 |
| I-030-003 | required | ① **三桶全做：IP + chat_id + telegram_user_id** ② 仅 IP ③ IP+user ④ IP+chat | **采纳①**：独立 limiter 实例；IP 挡解析前洪水，chat 挡群洪水，user 挡单用户刷命令。 |
| I-030-004 | non-blocking | ① **`channel.telegram`** ② 其它字符串 | **采纳①**：不是 `admin.*` 业务模块，不是 `core.*` 内核能力；不进 `mvp`/`admin` 默认集。 |
| I-030-006 | required | ① **请求计数：每次入站 Record，永不 Clear** ② 混合（请求计数 + secret 失败另桶） ③ 仅失败预算（只 Record secret/parse 失败，成功 Clear） | **采纳①**：VRev-070 / V-F114；不照搬登录 Allow/Record/Clear，避免合法 Update 放空洪水。 |
| 分发/mock 包 | — | ① **采纳建议包** ② 要改 | **采纳①**：见影响；合同正文 D-002。 |

## 建议包（用户一并冻结）

- 内核端口：`Register`/`Unregister` 命令 + callback；未知命令确定回落、不 panic；`SendMessage`（文本）；薄类型不含 SDK。
- Webhook：`POST /api/channel/telegram/webhook` + 头 `X-Telegram-Bot-Api-Secret-Token` fail-closed。
- 无 token：mock 记管理员可检视的出站记录（对标 mail outbox）。
- 实现在模块内；内核不 import `channel.telegram` 实现细节。

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 无 token 拒绝整个进程启动 | I-030-001 ② | 未配 Bot 的开发/测试 Profile 无法起 Admin；与 mail 默认 mock 先例冲突 |
| 无 token 仍接 Update | I-030-001 ③ | 入站合同空心，违背判据 #1 |
| 第三方 Bot SDK | I-030-002 ② | 多余依赖；类型泄漏风险；首波 HTTP 面极窄 |
| 本波只做 IP 桶（或少一桶） | I-030-003 ②③④ | 端口已覆盖三桶；少做会让 VP-031 立刻缺分母 |
| 仅失败预算 / 混合另桶 | I-030-006 ②③ | ① 已覆盖 secret 失败（IP 桶在解析前 Record）；另桶增加复杂度且非本波必需 |
| `admin.telegram` / `core.telegram` | I-030-004 其它 | 通道不是 Admin 业务列表，也不是内核能力 |

## 影响

- 合同正文 D-002 按本裁决冻结。
- Root / VP-030 信息台账 I-030-001/002/003/004/006 → `verified`（证据 = 本条目）。
- I-030-005 / I-030-007 保持 open（最晚 R3 / R2）。
