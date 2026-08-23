---
date: 2026-08-23
scope: GOAL-036 S6（I-001 浏览器 e2e 回归）；发现并修复后端关联清理缺陷
---

# E-004 · I-001 浏览器 e2e 回归：暴露并修复"删用户遗留 user_roles 孤儿"缺陷

## 执行事实

**首轮 e2e（修复前）**：admin × sqlite **8 通过 / 1 跳过 / 1 失败**；mvp × sqlite **9 通过 / 1 跳过 / 0 失败**。失败用例 `schema-crud.spec.ts:12`（admin 两次复跑均稳定失败）：新建角色并分配给用户、删除用户后，Roles 页该角色的 **Delete 按钮被禁用**（`title="You don't have permission for this action."`），断言 `toBeEnabled()` 超时。

**门控链路定位**（roles.json）：Delete 行动作 `disabledWhen: {field: "deletable", equals: false}`；后端 `roles.go:80`：`deletable = !System && AssignedUsers == 0`；`roles_repository.go:351`：`assignedUsers` 为裸 `COUNT(*) FROM user_roles WHERE role_id = ?`。

**探针实证（admin profile 实跑 API 状态机）**：

```
创建角色    → deletable=true  users=0
创建用户并赋角色 → deletable=false users=1
删除用户    → deletable=false users=1   ← 缺陷：user_roles 孤儿行残留
```

**根因**：`authsession/users_repository.go` — `DeleteUser` 与 `DeleteUsersBatch` 只删除 `refresh_tokens` 与 `users`，**从不清理 `user_roles`**（创建/更新路径写入，`users_repository.go:108-114 / 220-225`）。孤儿行使任何曾被分配过的角色 `assignedUsers` 恒 ≥ 1、`deletable=false` 永久化；`user_mfa`（TOTP 密文/恢复码凭据）同理残留。既有测试未覆盖"删除拥有角色链接的用户"场景。

**修复**：`DeleteUser` 与 `DeleteUsersBatch` 在同一事务内补删 `user_roles` 与 `user_mfa`（含注释引用本条目）。新增单元回归 `TestDeleteUserCleansRoleAndMfaLinks`、`TestDeleteUsersBatchCleansRoleAndMfaLinks`（`users_repository_test.go`）。

**修复后验证**：
- `go test ./... -count=1`：**全包 ok**（含 authsession 新增 2 用例）。
- e2e 重跑：admin × sqlite **9 通过 / 1 跳过 / 0 失败（exit 0）**；mvp × sqlite **9 通过 / 1 跳过 / 0 失败（exit 0）**；`schema-crud` 在两 profile 下均完整通过（含角色删除断言）。

## 诚实边界（不编造成因）

修复前 admin 稳定失败、mvp 首轮通过，同代码不同结果。可直接核验的因果链为"admin 失败 ← 孤儿行 ← DeleteUser 不清理"（探针实证 + 修复后双侧复绿）；mvp 首轮未暴露的具体时序成因**未能从既有痕迹重建（待确认）**——修复后该场景对两个 profile 均已由单元回归钉死，不再构成差异。附：探针脚本曾两次因 `go run` 包装进程被强杀而泄漏子服务占用 25080，后续已清理；第二版探针（含改密流程）在 admin profile 完成全状态机验证，mvp 侧复用同一后端代码路径，未单独重跑探针。

## 结论

**I-001 关闭**（证据：本条目 + `attachments/I-001-evidence.md`）。本缺陷属后端数据一致性缺陷（用户删除的关联清理），由 W25 的 e2e 回归暴露并修复；防复发由单元回归 + e2e 双侧绿共同保障。