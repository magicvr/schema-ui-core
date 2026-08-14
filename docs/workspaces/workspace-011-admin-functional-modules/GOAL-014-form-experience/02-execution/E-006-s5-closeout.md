---
id: E-006
goal: GOAL-014-form-experience
date: 2026-08-14
status: recorded
parent: GOAL-014-form-experience
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-006 · S5 关门（A-003 响应）

## 事实

- 2026-08-14：grok 独立审计（A-003）verdict fail，2 条 required（F-001/F-002）+ 5 条 non-blocking。全部响应后关门。

## Findings 响应

| Finding | 级别 | 处置 | 证据 |
|--------|------|------|------|
| F-001 schema 约束解析时被丢弃（required/pattern/length 运行时无效） | required | **fixed**：gateRenderFormFields 透传 + 回归测试 | render.ts + render.test.ts |
| F-002 fieldErrors 未进 ActionResult（回显断裂） | required | **fixed**：ActionResult/runRequest 转发 + handleSubmit 真实访问 + 回归测试 | render.tsx + error-localization.test.tsx |
| F-003 内联文案未 catalog 化 | non-blocking | **fixed**：form.validation.* 6 键 + messageKey + t() 翻译 | en-US/zh-CN.json + form-controls.ts + render.tsx |
| F-004 account_self 同码路径无 fieldErrors | non-blocking | **fixed**：writeLocalizedFieldError | account_self.go |
| F-005 Tailwind 动态类拼接 | non-blocking | **fixed**：静态 GRID_COL_CLASSES 查找 | form-controls.tsx |
| F-006 dateRange 空值语义 | non-blocking | **accepted-residual**（既有语义；后续波次可声明 required） | A-003 响应节 |
| F-007 缺自动化锁 | non-blocking | **fixed**：F-001/F-002 回归测试 + 全量复跑 | render.test.ts + error-localization.test.tsx |

- 回归：web 913/913 绿；go test ./... 全绿。

## 关门条件

- A-003 2 条 required 全部 fixed；5 条 non-blocking 已处置（1 条 accepted-residual 用户范围明确）。S5 关门成立。
