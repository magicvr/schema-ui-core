---
id: A-001
goal: GOAL-005-w4-security-audit-remediation
title: W4 八项修复 self 审计
source: self
date: 2026-08-11
verdict: conditional
status: recorded
---

# A-001 · W4 修复 self 审计

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | self |
| **auditor** | Claude Code · 主路径编排 + 第一手核对 |
| **类型** | execution-facts / close-out 前自审 |
| **scope** | GOAL-005 W4 八项修复实现正确性与回归风险（限流驱逐、上传权限门+配额、改密吊销 access token、前端异常捕获、URL 校验统一、web 启动加固、登录文案、autoComplete） |
| **verdict** | **conditional**（F-001 已 fixed；N-001/N-002 recommended） |

## 范围与区间

- 代码：`apps/api`（handler/rate_limit、upload、auth、authsession migration+repo、composition、account、testsupport、errorcatalog）+ `apps/web`（render、request-construction、theme、LoginPage、auth-client、form-controls、i18n）。
- 过程：`02-execution/E-001-w4-remediation.md`、`01-decision/D-001-w4-scope.md`、`00-meta.md`。
- 信息项：I-001/I-002 均已 verified（D-001）。本 scope 无到期 open required 信息门禁。

## 成果（有证据 · 本人核对）

| # | 项 | self verdict | 证据摘要 |
|---|----|--------------|----------|
| 1 | 限流容量驱逐 | **pass** | `allow()` 不再预建 key；`TestLoginRateLimiterAllowDoesNotRegisterKey` 验证 allow 后 map=0、登录路径驱逐后容量保持 2 |
| 2 | 上传权限门 | **pass** | `files.write` 中央注入 + `uploadPermissionGate`；`TestUploadRequiresFilesWritePermission` viewer 403 / 0 落盘 |
| 3 | 上传配额 | **pass** | `quotaReached` + 环境变量；`TestUploadPerUserQuota` 第 3 文件 429；frozen 集/双语/前端目录同步 |
| 4 | 改密吊销 access token | **pass** | 迁移 v11 + claims `tv` + middleware 比对；`TestUsersPasswordChangeRevokesAccessToken` + 既有断言更新 |
| 5 | 前端异常捕获 | **pass** | 三处 constructRequest 包 try/catch + `.catch` 兜底；renderer 205 测试绿、tsc clean |
| 6 | URL 校验统一 | **pass** | 三 builder 统一 `isRelativeProtocolUrl`；protocol 343 测试绿 |
| 7 | web 启动加固 | **pass** | initTheme/applySystemDefaultTheme/setTheme 全 try/catch；tokens.ts 同款 |
| 8 | 登录文案 + autoComplete | **pass** | AuthError.status → t() params 插值；`LoginPage.test` 无 `{status}` 断言；BaseInput `new-password` |

**复跑（2026-08-11）**：`go test ./...` 23 包 ok、`go vet` clean、`tsc --noEmit` clean、`vitest run` 747 测试过、i18n 62 过、protocol 343 过、renderer 205 过。

## Findings

### F-001 · required · low · 「allow 返回 false 仍保留过期剪枝」语义澄清 — **fixed（审计中已核对实现）**

`allow()` 对已存在 key 的窗口内列表写回 `kept`（含拒绝分支）。拒绝时写回保留窗口内失败（正确：record 会在下次失败追加）；放行时写回剪掉过期项（正确：不预建新 key）。实现与测试 `TestLoginRateLimiterUnit` 既有断言一致。无缺陷。

### N-001 · recommended · 配额为扫描式 O(files)

`quotaReached` 每次上传全目录扫描（含 Stat）。当前 `files.write` 权限门 + 1000 文件默认上限下可接受；未来上传量级增长应改为持久化计数或定期清理。不阻断。

### N-002 · recommended · migration 计数断言三处 10→11

`store` 包 `migrate_test.go`/`operations_test.go`/`restart_test.go` 的迁移计数硬编码更新为 11（含 `access_token_revocation` checksum 冻结 `c3ea720a…`）。未来新增迁移需同步三处 —— 属已知维护面，不阻断。

## 必改项汇总

- 开放 required：**0**。F-001 无缺陷。

## 与既有意见的异同

本波为新增波次，无历史意见可比对。审计输入的四路 agent 报告（api 认证面 / api 内核面 / web 协议面 / web 组件面）均在本目标 D-001 与 E-001 中留痕；主路径对关键 finding 逐一复核属实。

## 结论 + 建议

- **verdict：conditional**（无开放 required；N-001/N-002 为 recommended，不阻断关门）。
- 建议编排器：8 项成功标准全部达成；开放 required = 0；可推进关门。独立审（VP-009 provider）可另行开 A-002 复核。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
