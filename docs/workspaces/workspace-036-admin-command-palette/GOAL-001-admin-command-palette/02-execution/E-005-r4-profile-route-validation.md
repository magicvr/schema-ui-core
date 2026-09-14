---
doc_type: goal-execution
id: E-005-r4-profile-route-validation
status: recorded
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# E-005 · R4 Profile × permission × route 回归与关门准备（2026-09-14）

## 事实

- 已执行 `apps/web/src/app/searchable-profile-matrix.test.ts`，mvp/admin/demo/custom 四个 profile case 全部通过，分别核对目标页与 action 分母 `5/6`、`18/16`、`13/6`、`22/18`；动态/inner page、row/batch scope 与 provider errors 均符合 R1 oracle。
- 已执行完整 Web Vitest：**104 test files / 1366 tests passed**；TypeScript project build exit 0；Vite production build exit 0（只报既有 chunk size warning）。
- 已执行完整 Go API 测试：`go test ./...` exit 0，包含 Web/Go expression fail-closed 回归对应的 `apps/api/internal/account` 测试。
- 已在真实 API + Vite proxy 下执行新增 command-palette Playwright spec 的四个组合：mvp/SQLite、admin/SQLite、mvp/Postgres scratch、admin/Postgres scratch，每次 2 tests 均 passed；覆盖 page search、内部 History API 路由、Users page action modal、mobile functional row 与 Escape。
- 已核对 permission/feature/route 边界：visible menu 由 `projectNavigation` 投影；action 由 Schema mount 过滤并在 `invokeAction` 重检；未绑定 navigate、malformed row mapping、unknown/malformed expressions 均有专门拒绝证据；导航分组 active 展开有 App 集成断言。
- 已核对 VP-036 红线：未读实体行/数据库表、未扩展 pinned AppManifest、未引入搜索基础设施、未加入 recent/pinned/Saved Views/批量结果/未保存保护/Toast 重做/第二业务域。
- 证据已汇总于 `attachments/r4-profile-route-evidence.md`；R4 尚等待关门 self 审与 grok independent close-out，Root 当前保持 `active · 3/4`。

## 阻塞 / 风险

- 当前没有开放 required 信息项或 required Goal finding；R4 关门仍受最终 self/independent close-out 审计顺序约束，不把测试结果直接等同于 `done`。
- `pnpm test` 包装层仍受本机依赖构建脚本策略阻断；直接测试工具链已完整通过，需在最终审计中如实保留该环境限制。

## 下一步（计划）

- 记录 R4 self close-out，调用本地 grok build（grok-4.6 · reasoning high）做最终 independent close-out；响应所有意见后用户确认/编排器将 Root 由 active 关门为 done。
