---
doc_type: goal-audit
id: A-001-r4-self-audit
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-03
source: self
scope: R4 证据矩阵、退出判据与关门自审（全工作区范围）
audit_type: stage-closeout
verdict: pass
open_required: 0
---

# A-001 · R4 证据矩阵、退出判据与关门自审（self）

## 1. 审计基本信息

- **目标**：[GOAL-005-r4-evidence-closeout](../00-meta.md)
- **审视范围**：
  - VP-030 退出判据 1～8 证据链完备性（对照 `attachments/r4-evidence-matrix.md`）。
  - R1、R2、R3 子目标关门事实与审计记录（GOAL-002, GOAL-003, GOAL-004）。
  - 架构红线合规：未污染默认 Profile、无第三方 Telegram SDK 泄漏、无 Redis 依赖、无 Mini App/Stars 越界、不污染 `admin.users`、未改动 Charter。
  - 全量自动化测试：`go test ./...` 全绿。
- **审计模式**：`cross`（自审 A-001 + 本地 grok build 独立交叉审计 A-002）。
- **结论**：**PASS**（开放必改 0，建议 0）。

## 2. 判据达成度核验

1. **判据 #1（Webhook 合同）**：`POST /api/channel/telegram/webhook`，Secret 常时比对，未配 token 503，错误 secret 401 fail-closed，解析失败 400，成功 200。**PASS**。
2. **判据 #2（分发端口）**：命令去 `/` 与 `@BotName` 精确分发，未知命令发 `DefaultTelegramUnknownCommandText` 常量回落；Callback 精确分发。**PASS**。
3. **判据 #3（出站端口）**：`HTTPSender` 基于 stdlib `net/http`，10s 超时预算，支持文本与 InlineKeyboard 按钮，未配置 token 自动降级至 Mock 记录，公共类型无 SDK 泄漏。**PASS**。
4. **判据 #4（身份映射）**：`GetOrCreateSubject("telegram", user_id)` 幂等创建并填入 `upd.SubjectID`，存储失败返回 500 促使重试，不写入 `admin.users`，不依赖 `admin.wallet` HTTP 路由。**PASS**。
5. **判据 #5（设置与密钥）**：`RuntimeManager` 线程安全热切换，`SettingsHandler` 只读掩码脱敏，密钥 fail-closed 不进导出明文。**PASS**。
6. **判据 #6（限流评估落盘）**：VRev-070 评估通过，IP 60/m、Chat 30/m、User 20/m 三桶限流落地且在 Secret 失败时正确记账。**PASS**。
7. **判据 #7（边界保持）**：未进 `mvp`/`admin`/`demo` 默认集，无 Mini App/Stars/FSM，无 Redis/SDK 依赖。**PASS**。
8. **判据 #8（审计闭合）**：GOAL-002/003/004 均已完成审计闭合，无遗留必改项。**PASS**。

## 3. 下一步

自审通过。调用本地 grok build（grok 4.6 · 思考强度 high）执行 `/audit` 独立交叉审计。
