---
id: E-005
goal: GOAL-006-r2-f04-notification-center
date: 2026-08-14
status: recorded
parent: GOAL-006-r2-f04-notification-center
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-005 · schema capability 缺陷修复（通知页报错）

## 事实

- 用户报告通知页面报错：`form.recordSource requires capability "form.record.load" in meta.requiredCapabilities`。
- 根因：notifications.json 的 `notification-settings` 表单使用 recordSource + select，但 meta.requiredCapabilities 未声明 `form.record.load` / `form.controls.extended`（host 渲染器 fail-closed 校验，renderer/form-controls.ts + recordSource 校验点）。
- 全量审计所有模块 schema（recordSource/select/textarea/switch/upload vs meta capabilities）并同批修复：
  - notifications.json：+ form.record.load、form.controls.extended（本次用户报错）；
  - account.json：+ form.record.load（个人中心 recordSource 表单）；
  - users.json：+ actions.upload（上传字段）；
  - file-library.json：+ actions.upload（上传字段）；
  - settings.json 已齐备（对照基线）；dev/examples form-with-upload 已齐备（审计脚本误报，未改）。
- 验证：vitest 903/903 全绿（s5-denominator 真实渲染上述全部页面，capability 门禁通过）；go 相关模块测试全绿。
- 台账：本次为 done 目标后的维护修复，不改变各目标 status/progress；commit 见工作区历史。
