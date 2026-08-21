---
id: D-002
goal: GOAL-006-r2-f04-notification-center
title: S1 · 方案冻结 — admin.notifications 模块（必办-2 边界冻结 + 持久化/已读模型 + 首批系统事件 + 铃铛 UI）
date: 2026-08-14
status: accepted
parent: GOAL-006-r2-f04-notification-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-002 · S1 方案冻结（F-04 通知中心 · 站内通知）

> 依据：I-011-001 `3 F-04、`8 必办-2；GOAL-006 00-meta 边界与 S1 门禁。

## 1. 必办-2 · 模块边界冻结（A-001 F-001 / A-002 F-005 处置）

| 面 | 定义 | 归属 |
|----|------|------|
| **系统通知** | 平台安全/账号事件（锁定、停用、解锁、改密）——本模块 R2 交付 | `admin.notifications`（本模块） |
| **业务通知** | 领域实体事件（订单状态、钱包变动等）——投递面复用本模块 API，事件源归领域模块 | R3 领域波次（S-13/S-14）经 API 投递 |
| **公告（S-05）** | 广播型全量通知（含未读计数语义差异：公告一读全读） | R3 S-05 独立立项 |
| **消息模板（B-09）** | 模板化渲染（title/body 占位符、邮件通道） | R4 B-09 |
| 操作日志（C-07） | 通知 ≠ 操作日志；本模块不写 operationlog（事件源侧已有各自审计） | 不变 |

**R2 明确不做**：业务事件消费、公告、模板、邮件/推送通道、通知设置里的通道选择（仅站内开关）。

## 2. 持久化与已读模型（I-002 关闭）

- **迁移 0016**（`admin.notifications` owner）：`notifications` 表：
  `id TEXT PK / user_id TEXT NOT NULL / event TEXT NOT NULL / title TEXT NOT NULL / body TEXT NOT NULL / read_at INTEGER NULL / created_at INTEGER NOT NULL` + 索引 `(user_id, created_at DESC)`。
- **迁移 0017**：`users ADD COLUMN notifications_enabled INTEGER NOT NULL DEFAULT 1`（站内通知总开关）。
- **已读模型**：`read_at` 单值（已读/未读，无分级）；已读不可逆（无「标未读」——R2 不做）。
- **保留期/上限（清理策略）**：每用户上限 **500 条**（on-create 同事务裁剪最旧已读）；时间保留期（90 天）清理任务归 R3（S-03 定时任务）。R2 无后台任务。
- 通知设置：`PATCH /api/notifications/settings {enabled: bool}`（写 `users.notifications_enabled`）；开关关闭时停止**新通知产生**（事件钩子检查）。

## 3. 事件源（首批系统事件，I-003 关闭）

| 事件 | 触发点 | 内容 |
|------|--------|------|
| `account.locked` | auth.Login 锁定开窗（C-11） | 账号已锁定 15 分钟 |
| `account.disabled` | SetUserEnabled(disable) | 账号已被管理员停用 |
| `account.unlocked` | UnlockUser | 账号已由管理员解锁 |
| `account.password-changed` | 自助改密 + 管理员改密（UpdateUser password） | 密码已修改（安全提示） |

产生方式：事件钩子在 auth/repository 调用点旁触发 `notification.Recorder`（best-effort：失败只打日志，不阻断业务——与 operationlog 同纪律）。

## 4. 模块与端点

| 项 | 冻结值 |
|----|--------|
| 模块 ID | `admin.notifications` |
| 依赖 | core.auth-session / core.navigation-capability / core.schema-render / core.operationlog |
| 端点 | `GET /api/notifications`（分页 + `unreadOnly`）、`POST /api/notifications/{id}/read`（owner-only、幂等）、`POST /api/notifications/read-all`、`GET /api/notifications/unread-count`、`PATCH /api/notifications/settings` |
| 权限键 | 无（自服务面，任何已认证用户） |
| 页面 | `notifications`（navigation.user 区铃铛入口） |
| 持久化 | 迁移 0016 + 0017 |

**铃铛 UI（shell 级）**：App.tsx 顶栏右侧新增铃铛按钮——拉取 `/api/notifications/unread-count`（注入 fetcher），未读数徽标；点击弹出下拉（最近 5 条 + 「查看全部」链接到 `/notifications` 页）。R2 交付轻量下拉（不嵌套路由）。

## 5. 前端设计

- `notifications` 页（schema 驱动）：table（title/body/event/createdAt/read 状态列）+ 行操作 `read`（`POST /api/notifications/{id}/read` + confirm 不需要）+ 工具栏 `read-all` + `unreadOnly` 筛选（search-form 模式绑定？——R2 用工具栏 read-all + 行 read，筛选列以 query 参数直接呈现）。
- i18n：`manifest.title/nav.notifications`、`schema.notifications.*`、`shell.notifications.*`（铃铛 aria/徽标）en/zh。
- fail-open：加载失败 → 页面错误态；铃铛不可用时静默隐藏（徽标请求失败不阻断 shell）。

## 6. Profile 声明

- `admin.notifications` 加入 **mvp + admin** 默认集（账号安全通知是安全基线——与 F-03 同理由）。**Profile 内容扩展**声明同前（不改装配语义）。
- `adminFunctionalOrder` **不追加**（无首页贡献语义变更——notification 页不进 order；account 之后）。
- smoke.sh SM-007 页面集：mvp += notifications；admin += notifications；demo += notifications。

## 7. 必办核对（I-011-001 `8）

| 必办 | 适用 | 处置 |
|------|------|------|
| **必办-2（通知边界）** | **适用** | **✅ `1**（系统/业务/公告/模板四切分 + R2 范围声明） |
| 必办-1/3/4/5 | 其它目标 | 不适用 |

## 8. 未选方案（留痕）

- 不做「标未读」/已读分级；不做公告广播（S-05）；不做模板/邮件（B-09）；不做 90 天清理任务（R3 S-03）；不做推送通道。
- 铃铛为轻量下拉（最近 5 条），不做完整通知中心抽屉（页面承担详情）。

## 9. 实现范围（S2 清单）

1. 迁移 0016/0017 + compiled 注册。
2. authsession：`ListNotifications` / `MarkNotificationRead`（owner 幂等）/ `MarkAllRead` / `UnreadCount` / `SetNotificationsEnabled` + 裁剪逻辑。
3. 事件钩子：auth.Login（锁开窗）、users_state disable/unlock、account_self 改密、users.go 管理员改密 → `notification.Recorder`。
4. 端点 + 模块 provider + fragment（user 区）+ 页面 schema。
5. App.tsx 铃铛（fetcher 注入、徽标、下拉、a11y）+ 测试。
6. 装配 + smoke.sh 页面集 + i18n。
7. 测试：Go（迁移/已读/owner/裁剪/事件钩子/设置开关）+ Web（铃铛、页面、i18n）。
