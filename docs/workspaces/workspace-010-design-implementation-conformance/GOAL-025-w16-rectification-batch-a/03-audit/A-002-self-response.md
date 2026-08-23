---
id: A-002
goal: GOAL-025-w16-rectification-batch-a
title: 自审响应 · 闭合 A-001 required findings
source: self
date: 2026-08-17
verdict: pass
scope: 响应 A-001（F-001/F-002/F-004）
---

# A-002 · 自审响应 · 闭合 A-001 required findings

## 1. 响应对象

- A-001（independent · conditional）

## 2. 关闭证据表

| finding | 状态 | 证据 |
|---------|------|------|
| A-001 F-001 | **fixed** | `account-session-toolbar.tsx` 增加 `crud?.reloadList()`；成功后刷新会话列表。 |
| A-001 F-002 | **fixed** | `account_self.go` 增加 `newPassword == currentPassword → 400 INVALID_PASSWORD`；`w16_batch_a_test.go` 增加同密码拒绝断言。 |
| A-001 F-004 | **fixed** | `shell.spec.ts` / `schema-crud.spec.ts` 登录流程兼容强制改密；API 登录使用 `admin-e2e-pass`。 |
| A-001 F-003 | accepted-residual（recommended） | 老库种子 admin 不自动置 1；作为升级残余，由管理员重置密码或后续迁移处理，不阻断关门。 |
- **2026-08-23 兑现复核：fixed**——迁移 0049 幂等回填落地并验证（`TestMigrate0049*` PASS；store 包全绿）；见 [workspace-010 GOAL-033 E-005](../../../workspace-010-design-implementation-conformance/GOAL-033-w22-residual-closeout/02-execution/E-005-s2-completion-facts.md)。
| A-001 F-005 | open（recommended） | 并行 401 refresh 竞态；非 required，后续增强。 |
| A-001 F-006 | open（recommended） | 强制改密成功提示残留；非 required，后续增强。 |
| A-001 F-007 | open（recommended） | rotate 后下载按钮；非 required，后续增强。 |

## 3. 回归证据

- Go 全量 `go test ./...`：通过。
- Web 全量 vitest：1055/1055 通过；`tsc -b`、`npm run build` 通过。
- 新增/更新定向测试：Go `w16_batch_a_test.go`、Web `force-password-change`/`account-session-toolbar`/`LoginPage`/`mfa-manager` 通过。

## 4. 结论

A-001 的两条 required findings 已按 `fixed` 闭合；recommended 项留痕。S3 测试与回归门禁满足。
