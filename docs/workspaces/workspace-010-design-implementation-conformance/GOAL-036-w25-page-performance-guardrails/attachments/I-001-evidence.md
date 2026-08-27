# I-001 证据 · 浏览器 e2e 回归（admin/mvp × sqlite）

## 复现（修复前，2026-08-23）

- admin × sqlite：8 通过 / 1 跳过 / **1 失败**——`schema-crud.spec.ts:12`，失败点 `:132`（角色 Delete 按钮 `toBeEnabled()` 超时，按钮带 `title="You don't have permission for this action."` → 判定"无权限"实际来源为行级 `deletable=false`）。独立复跑该 spec 依旧失败（确定性）。
- mvp × sqlite：9 通过 / 1 跳过 / 0 失败（首轮通过成因待确认，见 E-004 诚实边界）。

## 根因链（探针实证，admin profile API 状态机）

| 步骤 | assignedUsers | deletable |
|------|---------------|-----------|
| 创建角色 `e2e_probe_*` | 0 | true |
| 创建用户并赋该角色 | 1 | false |
| 删除该用户 | **1（残留）** | **false（卡死）** |

代码链：`handler/roles.go:80` `deletable = !System && AssignedUsers==0` ← `roles_repository.go:351` 裸 COUNT ← `users_repository.go` `DeleteUser`/`DeleteUsersBatch` **不清理 `user_roles`（及 `user_mfa`）**。

## 修复

- `apps/api/internal/modules/authsession/users_repository.go`：单删 + 批量删均在同一事务补 `DELETE FROM user_roles` 与 `DELETE FROM user_mfa`（后者为 TOTP/恢复码凭据）。
- 回归测试：`TestDeleteUserCleansRoleAndMfaLinks` / `TestDeleteUsersBatchCleansRoleAndMfaLinks`。

## 修复后验证

- `go test ./... -count=1` 全包 ok。
- e2e：admin 9/9（另 1 profile 专属跳过）、mvp 9/9（另 1 profile 专属跳过），均 exit 0；`schema-crud` 双 profile 完整通过。