---
doc_type: goal-decision
id: D-001-r2-implementation-parameters
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
status: accepted
version: 0.1.0
---

# D-001 · R2 实施参数与接缝裁决

## 用户裁决

2026-09-04，用户对 R2 C1 的三项 required 信息作出书面裁决。以下选择只适用于 `[workspace-033-telegram-operator-console]` 的 R2，不改变 D-002 + D-003 已冻结的产品边界。

| 信息项 | 已接受方案 | 未选方案 | 理由与影响 |
|--------|------------|----------|------------|
| `I-033-014` | `mode` 与 `webhook_public_base_url` 使用 Telegram 专属持久化字段；首次无 DB 行时由 YAML/env seed，行存在后 DB authoritative；Admin PATCH 可更新 mode/URL；token/secret 仍加密保存 | 每次启动由 YAML/env 覆盖；或 DB 只读、PATCH 不更新 mode/URL | 与现有 Telegram token/secret 的单行 authoritative 模型一致；mode/URL 的来源可重启回读、可审计且不复用全局 `runtime.mode` |
| `I-033-015` | 未绑定 polling 使用活跃控制台会话的引用计数；每个会话 heartbeat 持有独立 TTL，TTL 至少覆盖两个 10 秒心跳（实现基线 20 秒）；引用计数归零或全部 lease 失效后 Stop 并 drain | 单共享 lease；或只允许显式业务绑定 | 多会话关闭其中一个不会误停；保留 D-002 的未绑定人工台懒启动边界，并明确失效收敛 |
| `I-033-016` | `getUpdates` 请求 timeout 为 30 秒；使用独立 polling HTTP client，client timeout 为 40 秒；继续遵守 client timeout 严格大于请求 timeout | 25 秒/35 秒；或新增部署必填 timeout 配置 | 空轮询次数与取消延迟可控；明确禁止复用现有 10 秒 `sendMessage` client，减少长轮询误报 |

上述是用户对参数方案的正式接受，不属于 `accepted-residual` 或 `user-overruled`。mode 缺省仍按 D-002 解释为 `polling`；webhook 仍须显式 URL 与非空 secret。

## 实施合同

- DB 行存在时，mode/URL 的空值也属于 authoritative 值；只有首次无行时才使用 YAML/env seed。PATCH 采用部分更新，未提供字段保留当前值。
- `mode` 只接受 `webhook` / `polling`；`webhook_public_base_url` 仍须满足 D-002 的绝对 origin 约束，polling 可不配置 URL。
- 未绑定 polling 的 lease 在内存中按会话身份计数；单个会话到期只移除自己的引用，不影响其他有效会话。无引用时 manager 进入 D-003 的 `idle` + `receiver=none`，并完成 receiver drain。
- 已绑定 polling 不受 heartbeat lease 约束而常驻；webhook 不启动 polling。所有切换和 Stop 仍由唯一 connection manager 串行化。
- `getUpdates` 的正常空结果/正常等待结束继续同一 loop；context 取消退出；Bot API/传输/协议错误 fail closed。polling client 与现有 `sendMessage` client 分离。

## recommended 接缝回应

- `HasBusinessHandlers()` 暂不扩展 `kernel.TelegramDispatcher` 公共端口；R2 以具体 `*Dispatcher` 或内部 adapter 提供只读探测，避免重开 VP-030 kernel contract。
- Telegram HTTP route 继续遵守现有 module/Profile gating；“不保留 webhook receiver”指 Telegram 侧 `deleteWebhook` 与进程内 polling loop，不运行时卸载 module route。
- R2 实施入口、计划与测试必须同时引用 D-002 + D-003 + 本 D-001；若实现需要改变上述用户选择，必须回到 `/govern` 请求新裁决。

## R2 C1 边界

本条关闭 I-033-014～I-033-016 的信息门禁并允许进入 C2/C3 计划与实现；本条不宣称 mode/URL 迁移、Bot API、connection manager、UI 或测试已完成。A-002/A-004 recommended 仍需在 R2 对应阶段以代码/测试证据回应。
