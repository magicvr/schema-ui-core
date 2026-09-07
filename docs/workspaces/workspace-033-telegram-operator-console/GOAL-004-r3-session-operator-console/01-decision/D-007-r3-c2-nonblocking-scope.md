---
doc_type: goal-decision
id: D-007-r3-c2-nonblocking-scope
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: user
status: done
version: 0.1.0
---

# D-007 · R3 C2 非阻断项处理范围与 polling 失败策略

## 用户裁决

用户要求将当前 C2 的非阻断审计项与实现一并处理，并接受此前裁决工具中“本轮补齐”范围：

- A-010 的推荐项纳入 C2 验证：gated PostgreSQL 首次写入、重复投递、并发唯一竞争，以及共同路径成功语义必须有明确测试证据；PG 未配置时只记录为 gated skip，不把 skip 当成通过。
- A-008 F-003～F-006 一并处理：空文本/未建模媒体明确跳过；私聊会话在缺少 chat title 时使用发送者姓名；v68 迁移与重复 dispatch 现钉更新；polling 持久化失败行为补测试。
- polling 收到更新但入站持久化失败时，选择**进入 error 状态**：保持当前 offset 不推进，退出当前 polling 循环，不在数据库故障期间热循环；后续 reconcile 或重启负责恢复。限流拒绝仍按 D-005 的既有兼容语义单独推进 offset。

## 边界与依据

- 本裁决不改变 D-004 的双表、规范化字段、不保存 raw JSON 或 C2 不建 outbound 状态机等选择。
- 本裁决不提前裁决 C3 的人工发送 API、显式重试身份、权限运行时细节或 C4 UI 行为；这些仍按既有信息门禁进入后续合同。
- 依据：A-008 F-003～F-006、A-010 F-001 的原始 independent 意见；用户在本轮要求“非阻断项一起处理一下”，并在裁决工具中选择“进入错误态”。

## 落地约束

- 代码与测试必须把该策略写成可核对事实；不得用文档把未运行的 PG 证据写成已通过。
- A-008/A-010 原始意见保持不变；本记录只是用户裁决和实施范围，不直接改审计 verdict 或目标状态。
