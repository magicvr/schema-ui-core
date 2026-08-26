---
title: E-002 · 生产装配回归锁实施与全量回归
status: active
created: 2026-08-26
updated: 2026-08-26
parent: null
version: 1.0.0
---

# E-002 · 生产装配回归锁实施与全量回归

按 [D-001](../../01-decision/D-001-w14-scope-and-freeze.md) 措施 B 实施（2026-08-26）。

## 变更清单

| 文件 | 变更 | 说明 |
|------|------|------|
| `apps/web/src/app/AuthGate.tsx` | 新增 | `AuthGate` / `BootScreen` / `useResourceFetcher` 自 main.tsx 原样提取；模块头注明提取动机与回归锁位置 |
| `apps/web/src/main.tsx` | 精简 | 移除上述定义与随之失效的导入（useCallback / createConfigAwareFetcher / useI18n / NavigationContext / App / LoginPage / InviteAcceptPage / ForcePasswordChange / lockedFailure / reauthFailure / AppManifest 类型），改为 `import { AuthGate, BootScreen } from "@/app/AuthGate"`；boot gate、failure boundary、组件自注册副作用全部原位保留 |
| `apps/web/src/app/auth-gate.wiring.test.tsx` | 新增 | 生产装配回归锁，jsdom 环境 |

## 回归锁设计要点

- **部分 mock 保真**：`vi.mock("@/account/auth-client")` 以 `importOriginal` 展开真实模块、仅替换 `restoreSession`——被断言的 `authFetch` 是生产传输层的**同一引用**；
- **捕获生产装配点**：`vi.mock("@/app/App")` 捕获 props，挂真实 `AuthProvider` + `I18nProvider`（无 `systemDefaultUrl`，零网络）+ 真实 `AuthGate` + 真实 manifest fixture（`test-fixtures/app-manifest.admin-dogfood.json`）；
- **断言**（用例 1）：`props.schemaFetcher === authFetch` 且 `!== globalThis.fetch`；`resourceFetcher` 为函数且非裸 fetch；navigationContext.user/features、onLogout、currentUser 接线完整；
- **分支冒烟**（用例 2）：未登录时不挂 `<App>`（props 为 null）、渲染 LoginPage 表单。

## 验证记录

| 验证 | 命令 | 结果 |
|------|------|------|
| 新锁 | `npx vitest run src/app/auth-gate.wiring.test.tsx` | 2 passed |
| 全量 web 回归 | `npm test`（vitest run） | **1130/1130 passed（84 文件）**，exit 0（基线 1128 + 本波 2） |
| 类型检查 | `npx tsc -b --pretty false` | exit 0 |
| 变异验证（锁有效性） | 临时移除 `schemaFetcher={authFetch}` 后复跑新锁 | **1 failed**，精确失败于 `expect(props.schemaFetcher).toBe(authFetch)`；恢复后 2 passed（详见 [A-001](../../03-audit/A-001-w14-self-closeout.md) §3） |

## 备注

- 以上改动均未提交 git；如需 checkpoint 由用户或关门轮次指示。
- vite dev 处于运行态，HMR 已生效；用户侧浏览器验证正常。
