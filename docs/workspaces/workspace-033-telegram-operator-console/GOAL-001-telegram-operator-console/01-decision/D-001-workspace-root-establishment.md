---
doc_type: goal-decision
id: D-001-workspace-root-establishment
parent: GOAL-001-telegram-operator-console
date: 2026-09-04
status: active
version: 0.1.0
---

# D-001 · 工作区 / Root 建立与开区决策

## 上下文

用户指令先走 `$vision` 激活 VP-033，再交 `$govern` 开设工作区，并书面“接受建议，继续”。VRev-075 self `pass`（open required = 0）确认 I-033-007/008 与 Admin freshness；slug 为 `workspace-033-telegram-operator-console` / `GOAL-001-telegram-operator-console`。

## 决策

| 项 | 决定 |
|----|------|
| 工作区 | `workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/` |
| Root | `GOAL-001-telegram-operator-console`；`parent: null`；`primary_plan = VP-033-telegram-operator-console` |
| 愿景角色 | `delivery`；不改变 Charter primary workspace |
| 纲领阶段 | R1 合同冻结 → R2 连接/热切换/占用位/设置页 → R3 会话落盘/人工 IM → R4 证据与关门 |
| 审计模式 | R1 常规合同冻结至少 self；涉及密钥、数据、生产 webhook 的实施/验收与 R4 关门按 `independent`，provider 遵循项目默认 grok build；若发生证据矛盾或元规则变化升级 cross |
| 红线 | 不重开 VP-030；不进业务域或默认 Profile；不实现多 bot/多实例 polling；不解除 SSE/WebSocket；不把群消息全量可见写成成功条件 |

## 用户冻结

| ID | 结论 |
|----|------|
| I-033-007 | 不要求关闭 Telegram Privacy Mode；只收录 bot 实际可见消息 |
| I-033-008 | 公网 base URL 显式配置；本地以可注入 Fake Bot API 验收 `setWebhook`，不做运行时/代理头猜测 |

## freshness 三字段

| 字段 | 值 |
|------|-----|
| consumer_vp | `VP-033-telegram-operator-console`（vision_ref `schema-ui-core-admin-foundation@0.4.0`） |
| last_freshness_review_at | 2026-09-04（`42036a3c` → `dd1edade` · Admin 类轻量 PASS；协议 pin/依赖锁/compose 无变；迁移 0066 与 runtime 装配可追溯至 VP-030；默认 Profile 未加入 `channel.telegram`；定向 Go 测试与 docscheck PASS） |
| next_freshness_review_trigger | Profile 默认集、Telegram migration/runtime、协议 provenance、依赖锁或生产部署/密钥边界变化；或首个 C 端业务域 VP 激活触发 H-002 复核 |

## residual 边界与未选方案

- workspace-030 `R-009 accepted-residual` 只在其原有范围内有效；本区若改变密钥存储或生产隔离，须回流 P-004，不能自动扩张。
- VP-030 本轮保持 `active`；V-F118 为 recommended，不阻断 VP-033 开区。
- 不把本意图放回 workspace-030，也不在本轮创建 R1 子目标或实施代码。
