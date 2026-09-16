---
id: E-005-r4-recommended-regression-coverage
doc: execution-entry
status: recorded
parent: GOAL-005-r4-unified-feedback-recovery
created: 2026-09-17
updated: 2026-09-17
version: 1.0.0
goal_id: GOAL-005-r4-unified-feedback-recovery
---

# E-005 · 补齐 recommended 回归证据

## 事实

2026-09-17，针对 A-002 的 recommended F-002/F-004，在不改变生产反馈合同的前提下补充回归：

- `feedback-policy.test.ts` 覆盖 `SERVICE_MAINTENANCE` → maintenance catalog 分类。
- `schema-table.test.tsx` 覆盖资源 maintenance 文案与显式 retry 恢复，并确认 401/403 列表错误不显示 retry。
- `render.test.tsx` 覆盖 chart 读取失败后的显式 retry，调用次数由 1 增至 2；recordSource 的 timeout/offline 页面回归已由 E-004 覆盖。

HostFailureScreen 与普通 resource feedback 的直接对照断言仍未新增；既有 Host fixtures 和代码边界保持不变，该项继续作为不阻断的 recommended 备注。

## 验证

- 全量前端 Vitest：110 个测试文件、1408 项通过。
- `tsc -p tsconfig.app.json --noEmit` 通过。
- `git diff --check` 通过（仅保留既有 CRLF 转换提示）。
