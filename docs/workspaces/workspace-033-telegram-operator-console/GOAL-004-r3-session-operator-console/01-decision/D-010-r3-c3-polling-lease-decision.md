---
doc_type: goal-decision
id: D-010-r3-c3-polling-lease-decision
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: user
status: done
version: 0.1.0
---

# D-010 · R3 C3 polling lease 与专用权限裁决

## 用户已裁决

未绑定 polling 在无控制台心跳时可以处于 `idle`/`receiver=none`，但 C3 operator
surface 继续要求 `running`、有效 receiver 和已确认 `bot_id`。现有 polling lease
的 acquire/heartbeat/release 授权扩展为接受 `telegram.operator.read`，与原有
`settings.read` 并列；因此 operator 读取权限持有者可取得并维持既有 lease，使
接收器进入 `running`，不必拥有 settings 权限。`settings` API 的权限与行为不变。

## 影响与边界

- 保留 VP-033 的未绑定按心跳启停和 `Dispatcher.HasBusinessHandlers()` 占用位；
  不由每个 operator API 请求隐式自启 polling。
- webhook 已处于 `running` 时不需要 lease；polling 下 C4 页面必须先用
  `telegram.operator.read` 获取 lease，再访问 C3 operator API。发送仍额外要求
  `telegram.operator.write`。
- 本裁决响应 A-018 F-002；不把 `idle` 直接视为人工台可用，也不改变 D-002 的
  专用权限选择或 C4 的 UI/发言权缓存范围。

## 依据

用户通过裁决工具选择“专用权限接管 lease (Recommended)”。该选择解决默认
polling 的 `settings.read` 回绑，同时保留现有运行时 active receiver 门禁。
