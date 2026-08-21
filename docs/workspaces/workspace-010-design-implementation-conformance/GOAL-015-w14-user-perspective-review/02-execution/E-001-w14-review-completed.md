---
id: GOAL-015-w14-user-perspective-review
doc: execution
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-001 · S1 审视完成（2026-08-17）

## 做了什么

以真实用户/管理员视角对 `apps/api` 与 `apps/web` 已实现功能做全量审视，并把「并非很小」的改进点落盘为台账 F-01～F-14（D-001 §2）。

- **编排器亲自复核**（读代码验证证据）：通用资源工厂 `resources.go`（分页/排序/搜索/批量删除/回收站/范围/错误本地化）、`errorcatalog.go`（zh/en 双语）、`App.tsx` 壳层、`LoginPage`、`render.tsx`、`schema-table.tsx`、`data-table.tsx`、`form-controls.tsx`、`resource.ts`、`boot.ts`、`AuthContext.tsx`、用户/账户/定时任务/钱包/回收站/通知/数据权限 schema。
- **三个并行独立审视面**（读实际代码，非推测）：
  1. API 面：通用工厂 + 全部命名 handler + 仓库层 + 错误目录。
  2. Web UX 面：壳层/登录/渲染器/表单/表格/弹窗/通知/MFA/i18n 双目录。
  3. 页面 schema 面：17 份页面 schema。

## 对子代理产出的证据复核（修正项）

- **修正 1**：schema 子代理称「markAllRead 是死代码、无批量已读入口」——复核 `notification-center.tsx:120-155` 发现批量已读按钮由自定义组件直接 POST `/api/notifications/read-all` 提供，功能存在；仅 schema 中 `markAllRead` 动作声明未被引用，降级为 schema 卫生（未单列 F 项）。
- **修正 2**：schema 子代理称「钱包 adjust/freeze/unfreeze/deductFrozen 无确认」——复核 `wallet.json` 发现这些为 `type: modal` 表单（点击→填表→提交，本身含二次步骤）；仅工具栏 `Reconcile` 为直接 request 无任何确认，故仅将其计入台账（F-14 相关与 B 类之外，未单列）。
- **关键发现独立复核确认**：F-01（scheduled-tasks 表单无 handler + `scheduler.go` 仅 `system.noop`）、F-02（`updateScopes` 无引用）、F-08（`App.tsx:669-672` 调试框）、F-09（toast 前缀错误码）均亲自读源码确认。

## 整体结论（事实）

- 产品成熟度高：通用资源工厂的分页/排序/搜索/批量/回收站/范围/错误本地化完备；错误目录 zh/en 双语且 `messageKey` 稳定；壳层/auth（刷新一次 401、reauth 终端、账户锁定、返回意图）成熟；表格/表单加载、空、错误状态与移动端卡片、焦点管理（抽屉/ModalHost）基本完备；i18n en/zh 键集对称。
- 改进点集中在：功能面补全（F-01/F-02）、审计与通知本地化（F-03/F-04）、面向用户的调试痕迹（F-08～F-10）、表单与无障碍（F-11～F-14）、列表一致性硬化（F-05～F-07）。无 P0 阻塞项。
