---
id: E-011
doc: execution
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-011 · S3 T-06 回归

- **Go**：`go test ./...` 全量 **0 FAIL**（含通知头断言）。
- **Web**：vitest **1037/1037**（65 文件；新增铃铛 2 例 + notification-center 3 例）；`tsc -b` 0；D-VAL 与 schema-keys 结构测试全绿（custom 节点经 GOAL-018 本地扩展合法）。
- **e2e（Playwright + 真实 Go API + SQLite 新鲜库）**：admin **8/8** + mvp **8/8**（各 1 跳过跨 profile 用例）；shell.spec 新增通知页冒烟（铃铛 → View all → 通知中心空态渲染）通过。
- 说明：未做通知「有数据」的全链路 e2e——系统无管理端造数端点，唯一事件源是账户事件（锁/禁用/改密），造数需改管理员密码（会吊销浏览器会话且污染后续用例），改由组件级单测覆盖交互（D-003 未选方案 ③）。
