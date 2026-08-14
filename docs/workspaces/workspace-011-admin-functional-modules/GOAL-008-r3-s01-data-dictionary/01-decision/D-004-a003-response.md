---
id: D-004
goal: GOAL-008-r3-s01-data-dictionary
title: A-003 响应：2 required + 1 recommended 全 fixed
date: 2026-08-14
status: accepted
parent: GOAL-008-r3-s01-data-dictionary
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-004 · A-003 响应（关门）

## 结论

A-003（grok-build independent，security/data）verdict **conditional**、2 required + 1 recommended —— **全部 fixed**（无 residual / overruled）：

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| F-001 非白名单表单控件 number + defaultValue 缺 advanced | required | **fixed** | 两页 schema：type number → inputNumber（白名单控件）；requiredCapabilities + form.controls.advanced（defaultValue 门禁）；AJV 校验 + form-controls 能力检查通过 |
| F-002 行 navigate 缺 actions.row.navigate 且 host 不执行 navigate | required | **fixed** | render.tsx invokeAction 实现 type:navigate 执行（window.location.assign，仅单斜杠同源路径，非法 URL fail-closed）+ 声明 actions.row.navigate（ADR-0021 host 路径）；单测覆盖（download-behavior.test.tsx：合法导航 assign + 绝对 URL 拒绝） |
| F-003 类型级联删除无条目级审计 | recommended | **fixed** | store.DeleteType 同事务收集条目 id 并返回；handler 在类型删除事件 detail 记录 {"entries":[...]}（json.Marshal） |

## 验证

- 修复后回归：go test ./... 全绿、vitest 898/898（新增 navigate 测试）、e2e mvp/admin 8/8。
- D-002 §1 修订：navigate 为 host 已执行的 ADR-0021 路径（非「无 renderer 扩展」——A-003 修正）；inputNumber/advanced 为呈现自由内的既有白名单控件。
