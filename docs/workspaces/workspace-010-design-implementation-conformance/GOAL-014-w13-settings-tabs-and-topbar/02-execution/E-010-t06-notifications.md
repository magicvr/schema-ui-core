---
id: E-010
doc: execution
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-010 · S2 T-06 通知中心交互修正实施

## ① 铃铛下拉条目可点击（问题 1）

- `apps/web/src/app/notification-bell.tsx`：条目改为 `role=menuitem` 按钮——点击 → best-effort `POST /api/notifications/{id}/read` → 关闭下拉 → `onOpenItem(id)`；未读圆点视觉（已读透明）。
- `apps/web/src/app/App.tsx`：`onOpenItem` 跳转 `/notifications?open=<id>`（深链）。

## ② 列表页点击条目 → 展开详情 + 标记已读（问题 2）

- 新增 GOAL-018 自定义组件 `apps/web/src/components/notification-center.tsx`（`main.tsx` 副作用注册）：
  - 读取共享查询状态（props.targetTable = notifications-table，搜索表单继续绑定 q/read 筛选）；
  - 行点击 → 行内展开详情面板（全文/事件/时间/已读状态）→ 未读条目自动 POST /{id}/read（已读条目只展开）；
  - 深链 `?open=<id>` 加载后自动展开并标记该条（铃铛跳转落点）；
  - 保留工具栏【全部标为已读】（POST /read-all）+ 分页 + 空态/错误态；
  - 静默刷新（行点击后不闪 loading，首载才显示）。
- `apps/api/internal/modules/notifications/schema/notifications.json`：table 节点替换为 custom 节点；**移除行内 markRead action（问题 3）**。

## ③ 未读数即时刷新（问题 4）

- Go `handler/notifications.go`：`POST /api/notifications/{id}/read` 与 `POST /api/notifications/read-all` 响应增加 `X-Schema-UI-Config-Changed: notifications.read`（复用 config-change 通道；组件/铃铛的 fetcher 均为配置感知）。
- `config-events.ts`：新增 `NOTIFICATIONS_READ_NAMESPACE`；铃铛订阅 → 徽标即时重查，下拉打开时同步重查条目。

## 测试

- Go：`notifications_test.go` 断言 read/read-all 响应头。
- Web：`notification-bell.test.tsx`（新增 2 例：下拉条目点击 → onOpenItem + read POST；notifications.read 事件 → 徽标重查）；`notification-center.test.tsx`（新，3 例：行点击展开+标记、已读行不重复标记、深链展开+标记）。
- e2e：`shell.spec.ts` 通知页冒烟（铃铛 → View all → /notifications 渲染 notification-center 空态 + 设置表单）。
