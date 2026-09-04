---
doc_type: goal-audit
id: A-007-r3-c2-contract-self
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: self
audit_type: contract-review
scope: R3 C2 用户裁决、入站双表/规范化字段、共同 webhook/polling 顺序与幂等合同
verdict: pass
open_required: 0
version: 0.1.0
---

# A-007 · R3 C2 入站合同自审（2026-09-04）

## 审视范围

核对 D-005 是否忠实承接 D-004 用户裁决、D-003 持久化确认顺序、C1 已确认的内部 `UpdatePayload` 边界，以及当前 webhook/polling 接缝。此意见只审合同，不把尚未存在的代码、迁移或测试写成完成事实。

## 结果

| 核对项 | 结论 | 证据 |
|---|---|---|
| 双表最小面 | 通过 | D-004 的会话表 + 入站消息表在 D-005 §持久化对象中落地；明确不建 outbound 表 |
| 规范化字段 | 通过 | D-005 仅列 chat/user/message/update、文本、callback 等字段，明确不保存 raw JSON、媒体或历史回灌 |
| 重复更新不重复落盘/分发 | 通过 | D-005 使用 `(bot_id, update_id)` 复合主键，唯一冲突不更新会话、不调用 Dispatcher |
| 事务与确认顺序 | 通过 | D-005 规定收据/会话同事务提交后才分发；webhook 失败不 2xx，polling 失败不推进 offset |
| command/callback 语义 | 通过 | 所有支持分发更新先有收据，普通 text 才进入成绩单；command/callback 仍可按既有 Dispatcher 分发且只分发一次 |
| kernel 与模块边界 | 通过 | D-005 保持内部 `UpdatePayload`、既有 Store/TxRunner 与 Telegram kernel port，不引入第二状态源 |
| 现有失败语义 | 通过 | handler 错误保持告警/不自动重试；限流 polling 的非持久化跳过例外被显式记录，未伪装成成功落盘 |

## 自审结论

`verdict: pass`，`open_required: 0`。D-005 足以作为 C2 代码实施合同；但它不是 C2 实现证据。A-003 的后续 recommended 项仍保持原意见状态，C2 代码、migration、并发和 webhook/polling 运行时证据必须在独立审计后补齐。

## 门禁

本意见不修改目标 `status` 或 `progress`，不关闭 C2，不接受 residual，不 overrule 任何既有意见。按用户要求，下一步调用本地 Grok Build（grok-4.6，reasoning high）进行 independent contract audit。

