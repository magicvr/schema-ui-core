---
id: D-004-t03-account-tabs
doc: decision-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-004 · T-03 个人中心选项卡（S2 分项冻结）

## 背景

`account.json` 纵向平铺资料 / 改密 / MFA / 会话。W11 U-11 已点名 Tabs（P2 未做）。渲染器已有 `type: "tabs"`（`TabsView`）。用户本轮采纳三档分组。

## 决定

1. **三个选项卡**：资料（改名表单）｜ 安全（改密 + MFA 管理器）｜ 会话（已登录会话表）。
2. **实现路径**：改 `account.json` 用已有 `tabs` 节点包裹现有子树；不新协议、不加路由。
3. **i18n**：`TabsView` 目前只读 `props.label`。实施时补 `labelKey`（与其它节点一致），避免中英混用。
4. **闭合 I-002**。

## 理由

- 安全相关（密码、第二因素）同属「保护这个账号」，放一档；会话是设备清单，单独一档。
- 四档会把改密和 MFA 拆开，用户改安全设置要切两次。
- 两档几乎不解决长页问题。

## 未选方案

- **资料 / 密码 / MFA / 会话**：过碎。
- **资料与安全 / 会话**：资料+改密+MFA 仍长滚。

## 影响

- 仅 account schema + 可能的 Tabs `labelKey` 小补；`self` 审计。
- 不改认证语义。

## 后续

- T-03 进 S3 P1。下一项裁决：T-04「我的钱包」归属与范围。
