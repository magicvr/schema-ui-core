---
id: D-005-t04-handoff-w011
doc: decision-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-005 · T-04「我的钱包」移交 workspace-011（本波不做）

## 背景

VP-010 不承载钱包业务实现。用户点名要「我的钱包」自服务页。P-004 选择题中用户选 **Other** 并书面：**相关工作在工作区 11 开子目标承载，当前子目标不做。**

## 决定

1. **本波（GOAL-013）不实现**「我的钱包」页面、自服务 API、或 `navigation.user` 钱包项。
2. **承接**：`docs/workspaces/workspace-011-admin-functional-modules/GOAL-022-my-wallet-self-service/`（Q2）。对话标签：[workspace-011] GOAL-022-my-wallet-self-service。
3. **T-01 仍有效**：用户下拉为个人中心 →（钱包项若存在）→ 设置 → 退出。钱包项由 GOAL-022 挂上后自然出现；本波缺席不算 T-01 失败。
4. **闭合 I-003**（归属已裁；自服务只读 vs 管理留给 GOAL-022 的 I 项）。

## 理由

- 账本/开户/扣款权威已在 [workspace-011 GOAL-019～021](../../../workspace-011-admin-functional-modules/GOAL-019-r3-s14-wallet-ledger/00-meta.md)。自服务页是同一领域的产品面，不应在符合性波次里开第二套钱包实现。
- 用户书面排除本区实施，避免与 VP-010 非目标冲突。

## 未选方案

- 本区做只读页+入口：与用户「当前子目标不做」冲突。
- 本区做可管理自服务：同样冲突，且扩大资金面。

## 影响

- GOAL-013 S3 不再含 T-04。W12 成功标准不因「无我的钱包页」而失败。
- 入口排序的最终可见性依赖 GOAL-022 关门或至少挂上 user-nav。

## 后续

- 下一项本波裁决：T-05 回收站删除时间。
