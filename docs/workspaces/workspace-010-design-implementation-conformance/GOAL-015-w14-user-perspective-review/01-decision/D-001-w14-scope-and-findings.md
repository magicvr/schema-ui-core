---
id: GOAL-015-w14-user-perspective-review
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · W14 开波：真实用户视角审视结论与改进项台账

## 1. 审视范围与方法

- **覆盖**：`apps/web`（壳层 `App.tsx` / 登录 / `renderer` 全部 / 表单控件 / 数据表格 / 通知铃铛 / MFA / i18n 目录）与 `apps/api`（通用资源工厂 `resources.go` / 各 handler / 模块页面 schema 17 份 / `errorcatalog` / 仓库层）。
- **方法**：以真实用户/管理员视角走查 as-built（读代码 + schema 声明）。三个并行独立审视面（API / Web UX / 页面 schema）+ 编排器对关键发现逐一证据复核（详见 E-001）。全部发现均有文件:行证据。

## 2. 改进项台账（F-01～F-14，按优先级分组）

### A 类 · 功能面补全（P1，非小改动）

| ID | 问题 | 证据 | 用户影响 |
|----|------|------|----------|
| F-01 | 定时任务页无法指定 handler：v1 仅注册 `system.noop`，create/edit 表单无 `handler` 字段，UI 建的任务永远空转（无实际执行目标） | `scheduledtasks/schema/scheduled-tasks.json`：`createTask`/`updateTask` 的 `bodyMapping` 含 `handler`，但 `openCreate`/`openEdit` 的 `fields` 无 handler；`scheduledtasks/scheduler.go:45-47` 仅 `system.noop` | 页面核心功能不可用（建了不执行） |
| F-02 | 数据权限页只能「注册策略」，没有任何入口给用户设置数据范围（`updateScopes` 声明但无 `actionRef`/按钮/行操作引用） | `datapermission/schema/data-permission.json:32-43` 声明 `updateScopes`（PATCH `/api/data-permission/scopes`）但 `body` 无引用；`policies` 表无 `actions` | 页面主目的（按用户授范围）无法达成 |
| F-03 | 审计/活动日志无结构化过滤（事件/操作者/时间范围）与导出 | `activity/schema/activity.json`：`ops-search` 仅 `q`；`handler/operations.go:14-26` Resource 仅 `QSearch`；`operationlog/repository.go:84-91` 仅关键字过滤 | 审计「谁何时做了什么」只能单关键字翻页 |
| F-04 | 系统通知文案英文硬编码，zh-CN 下不本地化 | `handler/notifications.go:273-285` 英文 `title`/`body`（"Account locked" 等）；文件头注明 deferred | 中文操作员收到英文安全通知 |

### B 类 · 数据一致性与健壮性（P2）

| ID | 问题 | 证据 | 用户影响 |
|----|------|------|----------|
| F-05 | 手写列表端点校验不一致：recycle-bin/wallet 吞掉非法 `page`/`pageSize`；wallet `pageSize` 无上限；data-permission policies 不分页且伪造 `pageSize`；per-task runs 硬上限 50 且忽略分页参数 | `handler/recyclebin.go:53-55`、`handler/wallet.go:57-59,276-278`、`handler/datapermission.go:53-79`、`handler/scheduledtasks.go:271-286` | 参数错误静默 / 大数据响应 / 分页能力名不副实 |
| F-06 | 错误码复用产生误导消息（`INVALID_SCOPE`/`INVALID_WALLET_BODY` 被多种校验复用，目录只给一条泛化文案）；`OPERATION_NOT_FOUND` 不在错误目录（无 messageKey） | `handler/datapermission.go:98-110,132-135`、`handler/wallet.go:80-82`、`errorcatalog.go`、`handler/operations.go:24`、`handler/systemmonitoring.go:95` | 表单错误难定位 / 404 无法本地化 |
| F-07 | 搜索/排序不一致：通知 `q` 大小写敏感；wallet 搜索仅匹配 `owner_id`；recycle-bin 不暴露 sort/order（仓库层已实现）；钱包 ledger 无 entry-type 过滤 | `handler/notifications.go:118` + `notifications_repository.go:136-138`、`wallet/store/repository.go:148-151`、`handler/recyclebin.go:53-55`、`wallet/schema/wallet-entries.json` | 搜索命中不可预期 / 排序缺口 |

### C 类 · 真实用户可见的调试痕迹与本地化（P1/P2）

| ID | 问题 | 证据 | 用户影响 |
|----|------|------|----------|
| F-08 | 每页标题右侧无条件显示 `pageId` + `route` 技术信息框 | `app/App.tsx:669-672` | 面向用户的内部标识泄露，观感像调试脚手架 |
| F-09 | 反馈 toast 前缀技术错误码 + 多处英文硬编码消息（如 `select at least one row first`、`action "x" is not defined`） | `renderer/render.tsx:1169` 前缀及 316/391/401/426/555/748/855/897 等；`schema-table.tsx:548-568` | 中文用户看到机器码/英文而非友好文案 |
| F-10 | 页面 schema 加载失败时把错误码当主标题 + 显示原始 URL/issue.path | `app/App.tsx:451-475` | 故障时无人类可读解释与恢复引导 |

### D 类 · 表单与无障碍（P1/P2）

| ID | 问题 | 证据 | 用户影响 |
|----|------|------|----------|
| F-11 | 必填字段无任何标记（无 `*` / `required` / `aria-required`），仅提交时才知道 | `renderer/form-controls.tsx` label 渲染；`form-controls.ts:485-501`；`render.tsx:1518-1528` | 填完整个表单提交后才被告知必填 |
| F-12 | 确认对话框无焦点圈 / 无 ESC / 无初始焦点（与 ModalHost 不一致） | `renderer/confirm.tsx:24-49` | 键盘/读屏用户可 Tab 到对话框之后，可能忽略阻塞确认 |
| F-13 | 桌面表格行不可键盘选中（无 `tabIndex`/`role`/`onKeyDown`；移动卡片已支持 Enter/Space） | `components/data-table.tsx:329-347` | 键盘用户无法打开记录详情抽屉 |
| F-14 | 若干小缺口：移动卡片列表静默丢弃第 4+ 列；单选 select 空值显示第一项（误导必填）；列排序无法清除；Tabs 无方向键导航/`aria-controls`；禁用按钮无原因说明；字段错误未 `aria-describedby` 关联；<sm 移动端无语言切换；通知空收件箱仍复用搜索空态文案「No items match.」（**已本地化**，属语义措辞欠佳，非英文硬编码） | `data-table.tsx:121-124,221-228`、`form-controls.tsx:292-310`、`render.tsx:1882-1910`、`App.tsx:902`、`notification-bell.tsx:179`、`schema-table.tsx:714-731` | 移动端/键盘/读屏体验与文案缺口 |

## 3. 建议实施顺序（供用户裁决；**下列阶段号属于未来整改波次，不是本波 S1～S4 检查点**）

- **整改波次冻结**：用户对 F-01～F-14 逐项 in-scope / defer 裁决 + 优先级；F-01 需同时确定 handler 目录暴露方式（新增端点 / 静态选项 / fork 扩展点），F-04 确定「存 messageKey」还是「存成品文案」方案。
- **整改波次实施（分批）**：
  1. A 类（功能面补全，价值最高，F-01/F-02/F-03/F-04）；
  2. C 类（调试痕迹清理，低成本高观感，F-08/F-09/F-10）；
  3. D 类（表单/无障碍，F-11～F-14）；
  4. B 类（一致性硬化，F-05/F-06/F-07）。
- **整改波次收尾**：回归（Go 全量 + vitest/tsc + e2e）+ 审计 + goal-tree 同步。

## 4. 审计风险分级（整改波次 S3 实施时按项升级）

| F 项 | 风险 | 建议模式 |
|------|------|----------|
| F-01/F-02 | 功能面，可能改端点/schema | self（改端点契约时按兼容性升级 independent） |
| F-04 | 通知文案存本地化键（数据语义变化，旧文案迁移） | independent |
| F-05/F-06 | 错误码/分页契约（兼容性） | self + 契约测试；改错误语义时 independent |
| 其余 | 呈现/无障碍，可逆 | self |
