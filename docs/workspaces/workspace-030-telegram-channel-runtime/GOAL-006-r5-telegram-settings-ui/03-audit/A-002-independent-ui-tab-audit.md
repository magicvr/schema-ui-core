---
doc_type: goal-audit
id: A-002-independent-ui-tab-audit
parent: GOAL-006-r5-telegram-settings-ui
date: 2026-09-03
source: independent
auditor: grok-4.6 (reasoning high)
scope: VP-030 判据 #5 Admin UI tab 交付（GOAL-006 C1/C2；不以 self A-001 为证据）
audit_type: execution-facts
verdict: pass
open_required: 0
status: recorded
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
---

# A-002 · 判据 #5 UI tab 独立复审（independent）

## 审计基本信息

| 字段 | 值 |
|------|-----|
| 被审目标 | [GOAL-006-r5-telegram-settings-ui](../00-meta.md) |
| 工作区 | `workspace-030-telegram-channel-runtime`（`root_goal` 匹配；`shared_materials_catalog: none`） |
| source | `independent` |
| auditor | grok-4.6（reasoning high） |
| 类型 | `execution-facts` |
| scope | VP-030 判据 #5「Admin 可配置 token/secret」的 **Admin UI tab 交付**（C1 Schema/Nav/Manifest + C2 前端组件/i18n） |
| 对照 | VP-030 退出判据 5；GOAL-006 D-001；[A-001 self](A-001-r5-closeout-audit.md) 仅作声称清单 |
| 方法 | 只读现码。本会话：`go test` telegram 三包 + kernel + composition ok；`vitest` `telegram-admin-tab` 2/2 + i18n structural 4/4 ok。不改 status |
| verdict | **pass** |
| 开放 required | **0** |

---

## 范围与区间

用户：`/audit（对判据 #5 UI tab 交付做独立复审）`。

判据 #5 原文（VP-030）：**Admin 可配置 token/secret；密钥 fail-closed；不进配置包明文。** 首波冻结表写「Settings 增加 Bot 渠道 tab」。GOAL-006 D-001（accepted）把形态冻成 **mail 同款独立页**（非 settings.json 内嵌 tab）。本审按 D-001 + 判据 #5 核交付，不把「必须是 settings 页里的一个 tab」当必改。

不审 F-001/F-002 装配/加密（已由 GOAL-001 A-008 闭合）。不改 Root/本目标 status。

---

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| Schema 页存在且挂 custom 组件 | `modules/channel/telegram/schema/telegram-settings.json`：`custom` / `component: telegram-admin-tab` |
| Provider 贡献 Pages + Nav + Manifest | `provider.go` L50–55、L108–139；Descriptor 与 `kernel/profile.go` BuiltinModules 对齐（含 schema-render / navigation-capability） |
| 组合根把 Pages 交给 schema handler | `composition.go` `handler.RegisterSchemas(mux, a, set.Pages)`；telegram 启用时 `telegrammodule.New` 进 providers |
| 独立菜单 `menu_telegram` | Nav `NodeID=menu_telegram`，`PageID=telegram-settings`，`Visibility=PolicyAdmin`；fragment `visibleWhen` `$context.features.menu_telegram == true` |
| 前端组件 write-only | `telegram-admin-tab.tsx`：GET 只读 `token_set`/`secret_set`；password 输入空；PATCH 只发非空字段；`main.tsx` 侧效注册 |
| i18n 双目录 | `en-US.json` / `zh-CN.json` 的 `schema.telegram.*`、`manifest.title.telegramSettings`、`manifest.nav.telegram`；structural 测试已扫 telegram fragment |
| 设置 API 密钥门 | GET `settings.read`；PATCH `settings.write`；响应走 `RuntimeStatus`（无明文/末四位） |
| 未进默认集 | mvp/admin/demo 不含 `channel.telegram`（composition / provider 测试） |
| 测试 | 本会话 Go + vitest 如上，均绿 |

「试发」在 VP 冻结表标可选，本波未做，不构成缺口。

---

## 对照成功标准

| 标准 | 判定 |
|------|------|
| C1 后端 Schema/Nav/Manifest + settings.read/write 沿用 | **达成（API 面）**。Nav 未引用 `settings.read` 字符串，见 R-001 |
| C2 前端 tab + 注册 + i18n + write-only | **达成** |
| 判据 #5 Admin 可配置 token/secret | **达成**：启用模块后有页面、菜单、GET/PATCH。密钥不回显 |
| D-001 独立菜单、无新权限键、不进默认集 | **达成** |
| 密钥不进配置包明文 | 非本目标代码；GOAL-001 A-008 已核 export `sensitiveFields`。本审不重开 |

---

## Findings

无 required。

### R-001 · 导航未绑 `settings.read`（PolicyAdmin 授权菜单）

- **严重度**：low · **建议**：recommended · **状态**：open
- **证据**：`provider.go` L119–131 `Permission` 留空，只设 `Visibility: PolicyAdmin`。对比 `modules/settings/provider.go` 的 `menu_mail`：`Permission: "settings.read"`。
- **原因**：`kernel/provider.go` L344–347 要求 nav.Permission 必须是**本贡献集已声明**的权限。`settings.read` 由 `admin.settings` 声明，跨模块引用会 `undeclared permission`。这是无新权限键红线的代价，不是漏写一行。
- **影响**：admin 角色会拿到 `menu_telegram`；去掉 `settings.read` 仍可能看见菜单，打开页后 GET 403。API 仍 fail-closed。
- **建议**：在 D-001/注释把该约束写成残余；或让 `channel.telegram` DependsOn `admin.settings` 并接受「无 settings 模块则无 UI」（仍不能跨模块填 Permission，除非改 kernel）。

### R-002 · `menu_telegram` 不在 `DefaultNavigationOrder`

- **严重度**：low · **建议**：recommended · **状态**：open
- **证据**：`kernel/provider.go` L403–431 与 `navigation_order_test.go` 快照含 `menu_mail` / `menu_mail_outbox`，**无** `menu_telegram`。未知节点追加在末尾（L59–60 测试说明）。
- **影响**：启用后菜单能出现，但落在侧栏底部，不在设置/邮件簇。
- **建议**：把 `menu_telegram` 插在 `menu_mail_outbox` 之后，并更新快照测试。

### R-003 · 无经组合根的 schema 200 断言

- **严重度**：low · **建议**：recommended · **状态**：open
- **证据**：`provider_test.go` 断言 Register 的 Page/Nav/Fragment；`composition_telegram_test.go` 覆盖 webhook/settings 401 与 Fx 同实例。没有 `NewApp`/`newMux` 上 `GET /api/schema/telegram-settings` 200、禁用模块 404。装配路径是 `RegisterSchemas(set.Pages)`，**不是** A-006 那种双工厂，假绿风险低于 F-001。
- **建议**：加一条启用/禁用模块的 schema 探测（对标 `s2_access_drill_test`）。

### R-004 · UI 不能清空已保存的 token/secret

- **严重度**：low · **建议**：recommended · **状态**：open
- **证据**：空输入被省略；后端对 JSON `null` 省略保持原值，显式 `""` 可清空。控制台没有「清除」动作。
- **建议**：若运营需要关 bot，补「留空=保持 / 显式清除」或二次确认清空。非判据 #5 必改。

---

## 必改项汇总

无。开放 required = **0**。

---

## 与既有意见的异同

| 既有 | 本审 |
|------|------|
| GOAL-006 A-001 self pass、findings 无 | **同意 pass / required=0**。不接受「零 recommended」：R-001～R-004 是真残余 |
| GOAL-001 A-006/A-008 把 Schema/Nav 降为 recommended | 用户书面补做后，**UI 面已落地**，判据 #5 不再是 API-only |
| VP 冻结「Settings tab」vs D-001 独立页 | 按 **D-001** 核。mail 同样是独立页 + custom component，不是 settings.json 子 tab |

---

## 结论 + 建议给编排器/用户的下一步

**verdict: pass。** 判据 #5 的 Admin 配置面（页 + 菜单 + write-only 表单 + i18n + API 权限）在现码上成立。建议 `/govern` 接受本 pass；R-001～R-004 可 residual 或顺手修，不阻断 GOAL-006 关门。

---

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
