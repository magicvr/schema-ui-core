---
id: A-005
goal: GOAL-007-w7-api-web-security-audit
title: W7 独立代码复核（F-001/F-002/F-006 代码层面修正验证）
source: independent
auditor: claude-sonnet-4 · thinking · 本会话独立代码复核（用户指令：独立复核 A-001 F-001/F-002 是否已在代码层面正确修正，并判断是否足以恢复 VP-008 go 宣称）
date: 2026-08-19
verdict: pass
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# A-005 · W7 独立代码复核（independent）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | claude-sonnet-4 · thinking · 本会话独立代码复核（用户指令：独立复核代码层面修正，不加载 skills；按 P-003 代贴） |
| **类型** | close-out / code-level verification |
| **scope** | A-001 F-001（MFA fail-closed）、F-002（MFA admin reset boundary）、F-006（captcha generate limiter）在现行代码中是否 genuine fixed；并判断 VP-008 go 宣称恢复条件是否满足 |
| **verdict** | **pass** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 已校验；`shared_materials_catalog: none`） |

## 范围与区间

- **覆盖**：A-001 F-001、F-002、F-006 三条 required 在 `apps/api` + `apps/web` 现行实现中的代码级闭合证据。逐条回读源码与关键测试，不以 A-002/A-004 的结论作为本条的闭合依据。
- **方法**：源码通读 + 测试阅读；未做动态 exploit / 渗透；未跑全量回归。
- **不覆盖**：不改 `status` / `progress` / 方案正文 / goal-tree。不把 A-001 recommended F-014～F-016 升格为本波 required。不把本意见当作已关门。
- **排除**：GOAL-002/004 已书面接受残余（refresh localStorage、匿名 schema/manifest、Compose 无 TLS、bcrypt cost、data-permission v1 未接线、development JWT、单会话吊销不 bump `token_version`）。

## P-005 / 工作区核对

| 核对项 | 结论 |
|--------|------|
| I-001（finding 清单） | verified；本条不重开 |
| I-002（high 是否暂挂 VP-008 `go`） | verified（D-002：F-001/F-002 闭合前暂挂）。本条独立确认 **A-001 F-001/F-002 已 genuine fixed**；VP-008 go 宣称恢复条件已满足，见下 |
| 到期 required 信息项 | 无到期未关闭项阻断本 close-out scope |
| 共享资料 | 无（`shared_materials_catalog: none`） |
| 工作区绑定 | Root / canonical / `plan_refs`+`primary_plan` 与 `workspace.md` 一致 |

## 代码级复核

### F-001 · MFA 登录门在存储出错时 fail-open → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| `Required()` 仅 `ErrNotFound` 视为不需要 MFA，其余存储错误 `return true` | `apps/api/internal/modules/mfa/service.go` L78–92 | ✅ 存储异常不再 fail-open。`s == nil` / `s.repo == nil` 防御性返回 false；`ErrNotFound` → false（无 MFA）；其他 err → true（需要 MFA，fail-closed） |
| Login 在 MFA 要求时进入第二因子分支，`BeginChallenge` 失败 → 500 `LOGIN_FAILED`，不发 token | `apps/api/internal/auth/auth.go` L196–201；`apps/api/internal/handler/auth.go` L143–161 | ✅ `Login` 在 `a.mfa.Required(u.ID)` 为 true 时返回 `MFARequiredError`；handler 调用 `BeginChallenge`，失败 → 500 `LOGIN_FAILED`，不发任何 token。存储异常时攻击者无法绕过 MFA 拿到 access/refresh |

**原缺陷**：`Required()` 把任何 `GetState` 错误（含 SQLite busy / 读失败）当成「不需要 MFA」，密码通过后 `Login` 直接 `issue()`。  
**现行为**：存储异常 → `Required()` 返回 true → 进入 MFA 分支 → `BeginChallenge` 也失败 → 500，无 token 发出。**fail-closed 成立。**

### F-002 · `users.mfa-reset` 无管理员目标边界 → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 非 admin 对 admin 目标 → 403 `ADMIN_ACCOUNT_FORBIDDEN` | `apps/api/internal/handler/mfa.go` L228–244 | ✅ `!slices.Contains(user.Roles, "admin") && slices.Contains(target.Roles, "admin")` → 403。与 `users.go` `authorizeAdminTargetBoundary` 镜像 |
| `AdminReset` 返回 `removedActive`，仅 true 时 `BumpTokenVersionAndRevokeAll` | `apps/api/internal/handler/mfa.go` L245–258 | ✅ `removedActive, err := service.AdminReset(targetID)` → `if removedActive { revoker.BumpTokenVersionAndRevokeAll(...) }`。未开通 MFA 的用户不会被当通用踢会话 |
| `AdminReset` 实现：`ErrNotFound` → false；pending → false；active → true | `apps/api/internal/modules/mfa/service.go` L262–274 | ✅ `st.Status == "active"` 才返回 true。pending 状态删除但不踢会话 |

**原缺陷**：MFA 重置无管理员目标边界检查；`DELETE FROM user_mfa` 后无条件 `BumpTokenVersionAndRevokeAll`。  
**现行为**：管理员边界已建立（非 admin 不得重置 admin）；`removedActive` 精细控制会话撤销（仅 active enrollment 移除才踢）。**两条子缺陷均已闭合。**

### F-006 · 登录验证码生成接口无计量 → **fixed**（E-003 后）

| 核对项 | 路径 | 结论 |
|--------|------|------|
| `allow()` 检查 → 拒绝时 429；通过后调用 `record()` 真正计数 | `apps/api/internal/handler/captcha.go` L57–65 | ✅ `allow()` 拒绝 → 429 `RATE_LIMITED`；通过后 `record()` 创建滑动窗口条目。`allow()` 只读不建条目（`rate_limit.go` L41–66），`record()` 才建条目（`rate_limit.go` L93–107）。契约正确 |
| 第 11 次匿名生成 → 429 已锁回归 | `apps/api/internal/handler/captcha_test.go` L64–94 | ✅ `TestCaptchaPreflightRateLimited`：独立 limiter、10 次 200 后第 11 次 429 + `error=RATE_LIMITED` |

**原缺陷**（A-003）：E-002 声称 F-006 fixed，但 `captchaGenerateLimiter` 只调 `allow()` 不调 `record()`，限流是空操作。  
**现行为**（E-003 后）：`allow()` → `record()` 正确配对，滑动窗口真正计数。题目仍为明文 1–50 加减（A-004 已记录为 OR 路径残余，不重开 required）。

## 对照 A-003 recommended（4 条非阻断 · 本波不升格 required）

本条独立确认 A-003 的 4 条 recommended 在现行代码中**仍为 recommended**，不阻断代码闭合：

| A-003 | 严重度 | 建议 | 现行状态 | 处理建议 |
|-------|--------|------|----------|----------|
| F-002 · Compose 文档示例 CIDR 过宽 | low | recommended | `compose.yaml` L44–48 注释仍写 `HTTP_TRUSTED_PROXIES=172.16.0.0/12`；默认实现正确（loopback-only + 无端口发布） | 修正注释为具体网段或 nginx 容器 IP `/32`；见下 E-004 |
| F-003 · PATCH 清理 `DeleteOrphan` 不校验 owner | low | recommended | `account_self.go` L196 `DeleteOrphan(oldAvatar)` 不查 owner；新绑定已强制 owner，原 IDOR 链已断 | 防御深度：对齐 `dropPreviousAvatar`；见下 E-004 |
| F-004 · 若干闭合路径缺少对准原 finding 的回归测试 | low | recommended | captcha 429 已补（E-003）；仍无 `Required()` 存储错误、委派 mfa-reset 打 admin、RFC1918 不信任、非 sessions 不含 `X-Refresh-Token` 断言 | recorded-residual；低风险可逆，不阻断本波关门 |
| F-005 · 锁定/禁用登录仍跳过密码哈希（时序枚举残余） | low | recommended | `auth.go` L165–172 锁定/禁用仍在验密前返回；未知用户烧 dummy bcrypt；错误码枚举已关 | recorded-residual；远程时序区分风险低，不阻断本波关门 |

## VP-008 go 宣称恢复判断

### 条件回顾（D-002）

> 在 F-001/F-002 两条 high required 闭合前，**不对外宣称 VP-008 `go` 消费有效性**；闭合后恢复宣称前应复核。

### 本条判断

| 条件 | 状态 | 证据 |
|------|------|------|
| F-001（MFA fail-closed）genuine fixed | ✅ | 本条第 F-001 复核；A-004 独立确认 |
| F-002（MFA admin reset boundary）genuine fixed | ✅ | 本条第 F-002 复核；A-004 独立确认 |
| 独立复核已完成 | ✅ | 本条（A-005）+ A-004（grok-4.6 independent）构成 cross 复核 |
| GOAL-007 全部 12/12 required 闭合 | ✅ | A-004 pass；本条未发现新 open required |

**结论**：VP-008 `go` 消费有效性恢复条件已满足。F-001/F-002 两条 high required 均已 genuine fixed 并经独立复核；A-001 全部 12/12 required 可核对闭合。**VP-008 `go` 宣称应从暂挂恢复为有效。**

恢复后：
- 后续业务 VP 在激活前仍需完成消费前 freshness review（VP-008 §go 消费有效性）
- VP-008 自身 `status: closed` 不变；go 消费有效性规则不变
- F-006（captcha generate limiter）已不再构成披露缺口（A-004 确认）

## 与既有意见的异同

| 来源 | 异同 |
|------|------|
| A-001 independent fail | 十二条 required 当时均 open。本条确认 F-001/F-002/F-006 现已 genuine fixed |
| A-002 self pass | 当时声称开放 required = 0，A-003 驳回 F-006。E-003 后 12/12 可复核闭合。本条**同意** A-002 对 F-001/F-002 的 fixed 判定（以现行代码为准）；对 F-006 以 E-003 + 现行代码为准 |
| A-003 independent conditional | 11/12 fixed、F-006 关闭声明不实。本条**同意**当时判定；**确认** E-003 已闭合 |
| A-004 independent pass | 12/12 required 均可核对闭合。本条**独立确认** F-001/F-002/F-006 的代码级修正，与 A-004 结论一致 |

无 P-004 冲突。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass** — A-001 F-001（MFA fail-closed）和 F-002（MFA admin reset boundary）在代码层面已 genuine fixed；F-006（captcha generate limiter）经 E-003 后已 genuine fixed。VP-008 go 宣称恢复条件满足。

建议 `/govern`：

1. 响应本条：将 A-001 F-001/F-002/F-006 维持 `fixed`（证据：本条 + A-004）。
2. **恢复 VP-008 `go` 宣称**：I-002 从「暂挂」更新为「已恢复」；D-002 追加恢复记录。
3. A-003 recommended F-002（compose CIDR 注释）→ 顺手修正，E-004 记录。
4. A-003 recommended F-003（DeleteOrphan owner 校验）→ 顺手修正，E-004 记录。
5. A-003 recommended F-004/F-005 → recorded-residual，不阻断本波关门。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。