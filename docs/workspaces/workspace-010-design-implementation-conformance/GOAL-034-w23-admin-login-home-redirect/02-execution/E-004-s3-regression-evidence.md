---
title: E-004 · S3 回归事实与连带整改（含双 flake 调查）
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-034-w23-admin-login-home-redirect
version: 0.1.0
---

# E-004 · S3 回归事实（2026-08-23）

## 全量回归矩阵（HEAD = 本波全部改动）

| 项 | 结果 | 证据 |
|----|------|------|
| `go test ./...`（apps/api） | 全包 ok（store 60.4s，其余 cached） | `e2e-w23-go.log`（job exit 0） |
| `npx vitest run`（apps/web） | **76 文件 / 1088 测试全绿**（含新写的菜单滚动契约单测） | `e2e-w23-vitest2.log` |
| `npx tsc -b --noEmit` + `npm run build` | TSC_EXIT=0 / BUILD_EXIT=0（vite built 4.91s） | `e2e-w23-build2.log`（job exit 0） |
| e2e admin 全量（10 用例 / 1 mvp-skip） | 自修复后 **连续 5 轮全绿**（run3 / run6 / run7 / run8 / run9-clean，每轮 9 passed） | `e2e-w23-full2/full6.log`、pwsh-52 双轮、最终 clean 轮 |
| e2e mvp 全量 | **9/9 passed**（1.7m；同一挂具，隔离语义下含 mvp localization 断言） | 最终矩阵轮（clean） |

## 连带发现与整改（S3 过程中命中 C3 门禁的既有问题）

**F-1（schema-crud 行操作菜单 flake，产品面）→ fixed**
- 现象：完整套件中 `schema-crud.spec.ts:90` 偶发 60s 超时（`menuitem "Password"` 不出现）；单跑 3/3 绿、刮扫探针证实菜单本体正常渲染。
- 根因：`RowActionsMenu`（W11 引入）在**任意** scroll（capture）时无条件关闭；紧随其后的模态框（`modal.tsx`）卸载焦点回恢复（F-002）引发的聚焦滚动与该关闭竞速——菜单刚打开即被撕掉，永不重开。
- 修复：`schema-table.tsx` 关闭语义收敛为「仅当触发钮真实位移（rect 变化 >1px）才关闭」；不足 1px 的无关滚动（内层容器/焦点布局滚动）不再误关。单测同步改写为「不动 → 保持；移动 → 关闭；外部 pointerdown → 关闭」三断言（`schema-table.test.tsx`）。
- 验证：vitest 1088 全绿；e2e admin 全量 5 连绿。

**F-2（sign-in fallback 等待竞态，测试面）→ fixed**
- 现象：`host-failure.spec.ts:71` 偶发：fallback 分支点击「Sign in」时按钮持续 disabled（首次登录 POST 未落定的窗口），点击动作空等到 60s 测试超时；后续轮未在插桩下复现。
- 修复：fallback 分支在点击前先 `toBeEnabled({timeout:15000})` 显式等待首 POST 落定（`sign-in.ts` 与 `localization.spec.ts` 两处同构），把不确定的空等转换为有界确定等待；保留 W23TRACE 插桩代码已移除。
- 验证：修复后全量 2 连绿 + 最终 clean 轮绿。

## 结论

- C3（S3 全量回归绿，含 e2e admin localization）达成：admin 套件自修复后无失败记录（5 连绿），mvp 套件见最终回填；产品代码改动仅 `schema-table.tsx`（菜单关闭契约，有单测锁定）。
- 无协议/manifest 契约语义变化 → 按 00-meta 边界声明，关门审计模式 `self`。