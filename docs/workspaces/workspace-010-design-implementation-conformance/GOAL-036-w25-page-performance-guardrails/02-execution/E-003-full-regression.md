---
date: 2026-08-23
scope: GOAL-036 S5 收尾（全量回归补跑）
---

# E-003 · 全量回归（S5 收尾）

## 回归事实

- **Go**：`go test ./internal/store/ -count=1` 通过（含新增 `store_wal_test.go` 4 用例）；`go build ./...` exit 0。
- **Web**：vitest 全量 **77 文件 / 1096 测试全绿**（相对 S4 的 76/1093 净增：`custom-components.schema.test.ts` 注册校验 1 用例 + `render.test.tsx` statCard+chart 合并、`refreshList` 定向刷新 2 用例）。
- **TypeScript**：`tsc -b` exit 0。

## 说明

- 首轮定向测试暴露并修复两处测试自身问题（非产品问题）：Go `Open` 需显式 `Dialect`、SQLite `PRAGMA synchronous` 读回数值 1（NORMAL）；注册校验测试需镜像 `main.tsx` 导入组件模块使注册表填充。均已留痕于 E-002。
- 监控页定向刷新行为变更（tick 只刷 `/status`）已由 `refreshList` 回归测试覆盖契约；e2e 断言面核对归入 S6（I-001）。