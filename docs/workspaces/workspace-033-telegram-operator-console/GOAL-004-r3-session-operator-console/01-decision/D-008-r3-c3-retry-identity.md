---
doc_type: goal-decision
id: D-008-r3-c3-retry-identity
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: user
status: done
version: 0.1.0
---

# D-008 · R3 C3 显式重试身份裁决

## 用户已裁决

C3 人工发送的每次显式重试都生成新的 `request_id` 与新的发送记录；重试记录通过 `retry_of` 关联其原始请求。首发记录没有 `retry_of`，同一 `request_id` 仍不得重复外发；失败只允许显式重试，不做后台自动重试。

## 影响与边界

- 发送记录必须保留每次尝试的 `pending` → `sent`/`failed` 状态与可追溯的原始请求关系。
- C3 合同需要定义 `retry_of` 的引用约束、状态转换、并发冲突和列表/成绩单聚合，但不得退回为原地状态迁移或无关联新 request。
- 本裁决只冻结重试身份模型；权限键、发送顺序、消息字段、分页参数和发言权缓存仍按 D-002 与 C3 实施合同处理。

## 依据

用户通过裁决工具选择“新 request + retry_of（推荐）”；该选择响应 A-003 F-002 关于显式重试身份未钉死的 recommended finding，并将 I-033-022 的 C3 发送合同从待定项推进为可写合同。
