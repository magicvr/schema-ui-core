---
id: E-014
doc: execution
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-014 · S3 T-07 回归

- `search-form-filters.test.tsx` 重写为新交互契约（单测逐条断言）：
  1. select 选择 → **立即**触发请求（含 enabled=true、无 q），chips 即刻出现 Enabled；
  2. 文本框输入 → **不触发**请求、chips 不出现关键词；
  3. 点击搜索按键 → 请求含 q+enabled，chips 出现关键词与 Enabled；
  4. chip 移除（×）→ 立即重新请求并去除该条件；
  5. reset → 清空全部并重新请求；配对断言（-ml-px/直角）保留。
- **Go**：全量 0 FAIL（本任务无 Go 改动，回归确认）。
- **Web**：vitest 全量 **1037/1037**（65 文件）；`tsc -b` 0。
- **e2e**：admin **8/8** + mvp **8/8**（各 1 跳过跨 profile 用例）——搜索表单改动未破坏浏览器全流程（含 users/roles 列表、schema-crud、w4 角色页）。
