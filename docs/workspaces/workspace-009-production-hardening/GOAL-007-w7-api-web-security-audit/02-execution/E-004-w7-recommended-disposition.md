---
id: E-004
goal: GOAL-007-w7-api-web-security-audit
title: A-003 recommended F-002/F-003 修复 + F-004/F-005 处置（非阻断）
date: 2026-08-19
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-004 · A-003 recommended 处置

## 触发

独立复核 A-005（pass）后按用户指令「同时处理 A-003 留下的 4 条非阻断 recommended」处置 A-003 recommended F-002～F-005。这 4 条均为 `low` · `recommended`，**非本波 required 闭合门禁**，但按用户指令本轮一并处理。

## 已发生事实

### A-003 F-002 · Compose 文档示例 CIDR 过宽 → **fixed（文档）**

- `compose.yaml` 注释不再建议 `HTTP_TRUSTED_PROXIES=172.16.0.0/12`（避免把信任面扩回整段 Docker 网桥），改为建议「compose 网络具体网段 或 web 容器内部 IP /32」。默认实现不变（loopback-only + 无 API 宿主端口发布）。
- 证据：`compose.yaml` L44–48。

### A-003 F-003 · PATCH 清理 `DeleteOrphan` 不校验 owner → **fixed（代码）**

- `raster_assets.go` 新增 `DeleteOrphanOwnedBy(raw, owner)`：仅在 URL 指向本 store **且** 资产 owner meta 匹配调用者时才删除（防御深度，与 `dropPreviousAvatar` 对齐）。
- `account_self.go` PATCH profile 清空/替换 avatar 改用 `DeleteOrphanOwnedBy(oldAvatar, user.ID)`。
- 原 IDOR 链（新绑定强制 owner）已断；本修复覆盖「历史 profile 已写入他人资产 URL」的防御深度场景。
- 证据：`raster_assets.go` `DeleteOrphanOwnedBy`；`account_self.go` L192–199；回归测试 `TestDeleteOrphanOwnedBy`（`account_avatar_test.go`）。
- 构建与测试：`go build ./...` + `go test ./internal/handler -run 'TestDeleteOrphanOwnedBy|TestAccountAvatar'` 通过。

### A-003 F-004 · 若干闭合路径缺回归测试 → **部分处置 + residual**

本轮新增两处对准原 finding 的高价值回归测试：

| 新增测试 | 覆盖 | 地址 |
|----------|------|------|
| `TestMFAResetAdminTargetBoundary` | 委派账号（持有 `users.mfa-reset` 但非 admin）打 admin 目标 → 403 `ADMIN_ACCOUNT_FORBIDDEN`；打非 admin → 204 + 会话撤销 | `handler/mfa_test.go` |
| `TestDeleteOrphanOwnedBy` | F-003 owner 守卫删除（他人资产不删、自有可删） | `handler/account_avatar_test.go` |

其余缺口记录为 **recorded-residual**（非阻断）：

| 缺口 | 处置理由 |
|------|----------|
| `Required()` 存储错误 fail-closed 单测 | `Service.repo` 为具体 `store.Repository`、非接口，注入失败 store 需重构；F-001 已由 A-004/A-005 独立代码复核 + `TestMFALoginTwoStep`（MFA 门必需路径）+ `BeginChallenge` 失败 → 500 覆盖。低风险可逆，不为此重构 |
| RFC1918 非信任断言 | 默认 `127.0.0.1/8` 实现 + `TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer`（loopback 信、公网不信）已覆盖主路径；补充 10/8·172.16/12 显式断言为证据增强，非阻断 |
| 非 sessions 请求不含 `X-Refresh-Token` 断言 | `auth-client.ts` `withAuth` 仅 sessions 路径设头；web 单测已覆盖 auth-client 主路径，补「非 sessions 缺头」断言为证据增强，非阻断 |

### A-003 F-005 · 锁定/禁用登录跳过 bcrypt 时序侧信道 → **recorded-residual**

- 错误码枚举（A-001 F-009）已关；锁定/禁用仍在验密前返回，未知用户仍烧 dummy bcrypt → 远程时序可能区分「存在且锁定/禁用」与「不存在」。
- 处置：**recorded-residual**，非阻断。理由：A-001 F-009 的账号状态枚举口已闭合；对锁定/禁用也烧 bcrypt 会在每次这些状态登录再增加固定成本，属产品/性能取舍；远程时序区分需在真实网络测量下才可能利用。owner = GOAL-007 / VP-008 lead；复核触发 = 后续安全审计或用户扩 scope。
- 证据：`auth.go` L156–190（锁定/禁用路径跳过 bcrypt；未知用户走 timingDummyHash）。

## 回归证据

- `go build ./...` 通过。
- `go test ./internal/handler -run 'TestMFAResetAdminTargetBoundary|TestDeleteOrphanOwnedBy|TestAccountAvatar|TestMFASelfService|TestCaptchaPreflightRateLimited' -count=1` 通过。

## 备注

- A-003 recommended 均为非阻断，本次处置不改变 GOAL-007 的 12/12 required 闭合状态（A-004/A-005 pass 维持）。
- VP-008 go 宣称恢复见 D-003（A-005 pass 已确立恢复条件）；本 E-004 不影响该结论。