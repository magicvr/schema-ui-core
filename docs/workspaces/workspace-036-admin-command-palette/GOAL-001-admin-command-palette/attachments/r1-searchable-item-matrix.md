---
doc_type: goal-attachment
id: r1-searchable-item-matrix
status: recorded
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# R1 · SearchableItem 分母与 Profile 覆盖矩阵

> 这是 R1 的可机器核对基线。它描述已确认的候选分母，不把页面 Schema 的所有顶层 action 定义误记为可执行的全局命令。

## 1. 首波分母规则

1. **页面/导航目标**：来自当前 `projectNavigation(manifest, currentPath, navigationContext, t)` 的已可见 top/sidebar/user 叶子。带 `pageRef` 的叶子按 `pageRef` 去重；无 `pageRef` 的站内 `url` 叶子按规范化 URL 去重。只接受可解析的非参数 href；未挂导航页、inner/detail 参数页和无法解析的动态目标不进入全局目标分母。
2. **Shell 入口**：认证上下文存在且 manifest 注册 `notifications` 页面时，补入现有 NotificationBell 的 `/notifications` 入口；它不是 Manifest navigation leaf，不改变 Manifest 或权限模型。
3. **声明式动作**：只从上述可见、可解析页面的已验证 Schema body 中收录直接触发器：`table.props.toolbar[]` 与 `actionButton`。行级 `table.props.actions[]`、`requiresSelection`/batch 触发器、需要 `$row` 或未绑定路径参数的动作不进入全局命令分母。动作保留原 page action 的 `modal` / `navigate` / `custom` / `request` 类型，但选择时必须交给现有页面级执行器，禁止直接调用原始 URL 或 actionRef。
4. **安全**：页面 Schema 获取与 D-VAL 失败时不暴露该页面的动作；显式 `visibleWhen`、`permissionIntent`、局部 `permissions` 和现有 permission cascade 在收录与执行时均重新核对；后端鉴权仍是最终边界。
5. **实体搜索边界**：不读取或索引 API 返回的实体行、数据库表、Saved Views、最近/固定项；`RT-X01` / `RT-X02`、Redis、MQ、跨进程索引和多实例均保持 gated。

## 2. Profile 覆盖矩阵

| profile / 当前候选 | 启用模块数 | 页面注册数 | 持久化导航贡献 | 发布 Manifest 导航叶子 | Shell 目标数（叶子 + 通知） | Schema 顶层 action 定义数 | 可收录页面级触发器 | 浏览器证据策略 |
|---|---:|---:|---:|---:|---:|---:|---:|---|
| `mvp` | 11 | 6 | 5 | 4 | 5 | 32 | 6 | 现有 Playwright mvp / SQLite + Vitest fixture 矩阵 |
| `admin` | 23 | 22 | 18 | 17 | 18 | 93 | 16 | 现有 Playwright admin / SQLite + Postgres + Vitest fixture 矩阵 |
| `demo` | 12 | 14 | 5 | 12 | 13 | 34 | 6 | Vitest/fixture 矩阵；不扩展现有 Playwright profile |
| 当前 `config.yaml` custom（admin + Telegram + digital-offer） | 25 | 27 | 22 | 21 | 22 | 101 | 18 | Vitest/fixture 矩阵；浏览器门禁仍由 mvp/admin 双方言承担 |

- 页面注册数包含已注册但不满足全局目标规则的 inner/detail 页面；因此它不等于全局页面目标数。
- 持久化导航贡献包含没有公开 Manifest leaf 的 `menu_notifications`；`notifications` 通过现有 Shell NotificationBell 计入 Shell 目标数。
- 顶层 action 定义数来自所有模块 Schema 的 `actions` 对象；可收录触发器只计算直接 toolbar/actionButton 入口，故两列不相等。
- 当前 `config.yaml` 的 custom 候选明确包含 `channel.telegram` 与 `biz.digital-offer`，不包含 `dev.examples`；Profile/optional module 的事实来源是 `apps/api/kernel/profile.go` 与 `apps/api/configs/config.yaml`。

## 2.1 稳定 item ID 清单（R2/R4 oracle）

以下是按当前 profile 候选、以管理员可见菜单 projection 计算的稳定 ID。不同用户角色/自定义菜单 grant 仍由 `projectNavigation` 在运行时取交集；清单不把隐藏项当作可见事实。

| profile | page/navigation IDs | action IDs |
|---|---|---|
| `mvp` | `page:dashboard`, `page:users`, `page:roles`, `page:account`, `page:notifications` | `action:users:create`, `action:users:invites`, `action:users:export`, `action:users:import`, `action:roles:export`, `action:roles:create` |
| `admin` | `page:dashboard`, `page:users`, `page:roles`, `page:file-library`, `page:data-dictionary`, `page:system-monitoring`, `page:scheduled-tasks`, `page:recycle-bin`, `page:data-permission`, `page:activity`, `page:mail`, `page:mail-outbox`, `page:wallet`, `page:wallet-vouchers`, `page:my-wallet`, `page:settings`, `page:account`, `page:notifications` | `action:users:create`, `action:users:invites`, `action:users:export`, `action:users:import`, `action:roles:export`, `action:roles:create`, `action:settings:reset`, `action:file-library:upload`, `action:data-dictionary:create`, `action:scheduled-tasks:create`, `action:recycle-bin:purgeAll`, `action:data-permission:register`, `action:my-wallet:redeem`, `action:wallet-vouchers:generate`, `action:wallet:create`, `action:wallet:reconcile` |
| `demo` | mvp IDs + `page:overview`, `page:data-table`, `page:admin-list-batch`, `page:data-display`, `page:search-form-table`, `page:form-controls`, `page:form-with-reactions`, `page:form-with-upload` | mvp action IDs; the `admin-list-batch` trigger is excluded because it requires a selection |
| current custom | admin IDs + `page:telegram-settings`, `page:digitaloffer-offers`, `page:digitaloffer-entitlements`, `page:digitaloffer-purchases` | admin action IDs + `action:telegram-settings:telegram-operator-entry-button`, `action:digitaloffer-offers:create` |

> The custom `page:`/`action:` IDs are intentionally namespaced by page and trigger key. A provider collision on any ID is an error and removes the conflicting candidates.

## 3. 排序、匹配与去重

- 契约版本：`SearchableItem` / provider v1。
- 查询：不持久化；大小写与重音不敏感；空查询按 provider id、item rank、稳定 item id 展示；非空查询依次为完整标签精确命中、标签前缀、关键词/标签包含，最后按稳定 provider/item 顺序裁剪。
- 上限：12 条。
- 冲突：provider/item 重复稳定 id 不静默覆盖；冲突 id 的全部候选从结果中排除并报告非敏感 provider 错误。
- 证据：用户 2026-09-14 书面确认本附件所述推荐范围、契约和 UX 口径，决策见同目标 `D-002`。

## 4. 已知实现缺口（进入 R2/R3）

- 当前 Manifest 协议没有 module/profile/provider/action 字段；不得扩展 pinned AppManifest envelope。内置 provider 从已验证 Manifest + authenticated page Schema 组装，外部 provider 通过前端注入 seam 接入。
- 当前 `SchemaCrudProvider.invokeAction` 需要补齐 programmatic modal/navigate/custom/request 的统一权限复核，并修正 actionButton 无显式 `props.key` 时的 node id 目标传递；否则 palette 可能绕过页面 UI 入口的既有 gating。该缺口是实现任务，不是对现有权限语义的静默放宽。
