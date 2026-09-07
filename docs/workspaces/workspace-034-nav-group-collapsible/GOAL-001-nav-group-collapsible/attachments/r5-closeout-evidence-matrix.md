---
doc_type: goal-attachment
id: r5-closeout-evidence-matrix
parent: GOAL-001-nav-group-collapsible
status: recorded
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# R5 · VP-034 / Root 关门证据矩阵

> 本矩阵把 VP-034 方向级退出判据映射到 workspace-034 Root 的实现、测试、Profile/route 证据。R5 self A-010、Grok independent A-011 与编排响应 A-012 已完成；`verified` 仍只表示证据存在，Root `done` 由 Goal meta/goal-tree 保存。

| # | 退出判据 | 证据 | 状态 |
|---:|---|---|---|
| 1 | Shell sidebar 支持 group 渲染、折叠/展开、Enter/Space 键盘切换 | `apps/web/src/app/App.tsx`；`apps/web/src/app/nav-groups.test.tsx`；Web full regression | verified |
| 2 | NavigationContribution 可选 group；既有权限、route、slot、NodeID 兼容 | `apps/api/kernel/contribution.go`；17 个模块 Provider；`apps/api/kernel/provider_test.go`；API full regression | verified |
| 3 | 跨模块共组聚合稳定，无重复/丢失，无中央业务注册分支 | `apps/api/internal/manifest/manifest.go`；`manifest_test.go`；composition + serve 双 assembly；R4 custom matrix | verified |
| 4 | 当前已注册 sidebar 全量迁移：五组、Dashboard 单例、Examples 保留 | `attachments/r1-navigation-profile-slot-matrix.md`；`attachments/r4-navigation-route-profile-matrix.md`；`nav_group_r4_test.go` | verified |
| 5 | 直接 URL 进入已分组页面及登记内页/动态路径时父组自动展开 | `apps/web/src/app/navigation.ts`；`nav-groups.test.tsx`；`nav-groups-r4.test.ts`；`navigation.test.ts` | verified |
| 6 | 未声明 group 平铺兼容；top/user slot、default/admin/demo/custom、权限与直接 URL 回归 | `apps/api/internal/composition/nav_group_r4_test.go`；`apps/web/src/app/nav-groups-r4.test.ts`；API/Web full regression | verified |
| 7 | module-contribution-playbook 更新 group key、顺序、共组、未分组例外与迁移规范 | `docs/architecture/module-contribution-playbook.md` v1.2.0 | verified |

## 信息门禁

| ID | 状态 | 证据 |
|---|---|---|
| I-034-001 | verified | R1 静态分母 + R4 runtime matrix |
| I-034-002 | verified | 用户 D-002 + R2/R4 Manifest 行为 |
| I-034-003 | verified | 用户 D-005 + R3 sessionStorage tests |
| I-034-004 | verified | D-005 + R3/R4 route projection matrix |
| I-034-005 | verified | R4 custom/demo runtime matrix |

## 审计与追溯

- R1/R2：A-001 self、A-002 Grok independent、A-004 required closure、A-005 self、A-006 Grok independent、A-007 response。
- R3：A-008 self；R3 checkpoint `6e581ca9`。
- R4：A-009 self；R4 checkpoint `b25bd777`。
- R2 checkpoint：`41e89f47`。
- R5 close-out：A-010 self `pass`、A-011 Grok independent `pass`、A-012 response；最终验证 `go vet`/`go test`/Vitest 99/1339/tsc/vite 均通过；Root 状态 checkpoint `b1d569a1`，最终台账同步 checkpoint `b4efabda`。

## 残余 / 非目标

- 不引入服务端分组持久化、分组独立权限、多级嵌套或拖拽排序。
- 现行协议 NavGroup 不增加 key/id；sessionStorage 是会话内 UI 状态，不是服务端事实。
- 本 Root 矩阵无已知 required residual；A-002 F-007 是 VP-034 计划层的历史 recommended 文案卫生项，保留给 `/vision`，不阻断 Root `done`。
- 本矩阵已完成 R5 审计响应；后续若导航契约变更，仍须按 P-003 重新审计。
