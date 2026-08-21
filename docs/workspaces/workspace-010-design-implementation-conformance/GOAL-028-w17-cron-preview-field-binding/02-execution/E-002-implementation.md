---
id: E-002
goal: GOAL-028-w17-cron-preview-field-binding
title: S2 实施字段绑定与 describeCron 本地化
status: completed
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: GOAL-001-design-implementation-conformance
---

# E-002 · S2 实施（2026-08-18）

## 已发生事实

1. `FormControlField` / `gateRenderFormFields` 透传本地属性 `afterComponent`。
2. `FormControls` 在字段下方渲染已注册 custom component，传入 `node.props.bindValue`。
3. `cron-preview` 在存在 `bindValue` 时进入绑定模式（无独立输入，400ms 防抖）。
4. `scheduled-tasks.json`：create/edit 的 `cron` 字段挂 `afterComponent: "cron-preview"`；移除页面块 `cron-preview-block`。
5. `describeCron` 按 `Accept-Language`（`errorcatalog.Negotiate`）输出中文/英文人话；覆盖每分钟、每 N 分钟、每小时第 N 分钟、每天、每周、每月，其余回退「5 段 Cron 计划」。

## 阻塞

无。

## 下一步（计划）

S3 记录定向测试结果。
