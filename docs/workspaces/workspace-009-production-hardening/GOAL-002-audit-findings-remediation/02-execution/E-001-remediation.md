---
id: E-001
goal: GOAL-002-audit-findings-remediation
title: 16 项审查发现修复与回归
date: 2026-08-10
status: recorded
---

# E-001 · 16 项审查发现修复与回归

## 范围

2026-08-10 代码审查发现的全部 16 项（C1–C8 中高危 + D1–D8 低危），输入 `raw/audit-20260810-api-web-bug-review.md`（gitignored）。

## 修复清单

| id | 修复 | commit |
|----|------|--------|
| C1 | 上传存储型 XSS：服务端 `http.DetectContentType` 检测 + 拒绝 HTML/XHTML/SVG；下载强制 `Content-Disposition: attachment` + `Content-Security-Policy: sandbox`；`UPLOAD_ALLOWED_TYPES` 按检测类型校验 | `5cbf3dc` |
| C2 | 刷新竞态双面：`RevokeRefreshToken` 单语句原子化（`WHERE revoked_at IS NULL` + RowsAffected）；`Refresh` 对 `ErrAlreadyRevoked` fail-closed（`ErrTokenRevoked`）不再签发第二对；`Logout` 容忍并发已撤销；前端 `auth-client` 共享 in-flight refresh Promise；并发测试 1/1（3 连跑稳定） | `5cbf3dc` |
| C3 | `APP_ENV` fail-open：默认值 `development` → `""`（未设置即拒绝启动，`ValidateProd` 显式报错）；显式 `development` 才允许 dev 密钥/密码；QUICKSTART/README 同步 | `5cbf3dc` |
| C4 | Bootstrap 一次性锁死：`WasFresh()` 门 → `NeedsBootstrap`（users 表为空则重试）；composition + testsupport 同步 | `5cbf3dc` |
| C5 | 异步错误：`runRequest`/`runBatchRequest` 捕获 fetch 网络失败 → `REQUEST_FAILED` ActionResult；`handleSubmit` try/finally 防按钮卡死 | `e498789` |
| C6 | 搜索清空：`q` 为空显式写入（覆盖旧过滤） | `e498789` |
| C7 | 权限缺省契约：`runRequest` 仅对已声明 permission entry 的目标 gate（与 batch 路径一致） | `e498789` |
| C8 | 路由 query 贯通：App 持 `routeQuery`（深链/popstate/onNavigate）；`SchemaPageSurface` 注入 `context.route={params,query}`；recordSource 使用 route 绑定；onNavigate 剥离 query 防 matchRoute 干扰 | `e498789` |
| D1 | `PATCH roles:null` 视为未提供（不清空角色）；显式 `[]` 仍清空；测试 | `9c1d0a7` |
| D2 | 登录限流（按 IP 滑动窗口 20/15min → 429 `RATE_LIMITED`）+ 缺用户 dummy bcrypt 时序均衡；`RATE_LIMITED` 入错误码契约 + i18n；测试 | `9c1d0a7` |
| D3 | **用户裁决（2026-08-10）**：保持 schema 匿名可读为设计决策（登录壳预认证渲染 + e2e 依赖），记录 accepted-residual，无代码改动 | — |
| D4 | `splitUrl` 畸形百分号编码安全解码（不抛 URIError）；conformance 测试 | `9c1d0a7` |
| D5 | 升级快照文件名毫秒粒度（同秒重试不冲突）；per-pending 快照契约（I-011-002 A-002 F-002）保留 | `9c1d0a7` |
| D6 | `defaultLocale`/`defaultTheme` 显式空串归一化 `auto`（存储/GET/Validate 三层一致）；测试 | `9c1d0a7` |
| D7 | `inputNumber` 清空 → `undefined`（提交时字段省略、后端默认生效）而非 0；上传 input 上传后重置 | `9c1d0a7` |
| D8 | `ThemeToggle` 改用 `setTheme`（同步 CSS `color-scheme`） | `9c1d0a7` |

已证伪/不入账（审查阶段裁定）：Commit 失败连接池中毒（modernc v1.55.0 tx.go:35-54 驱动层已防护）；devSession 默认开启（config.go 默认 false）。

## 回归验证（2026-08-10）

- `go test ./...`：21 包全绿（含新增：上传 XSS 测试、并发刷新 1/1 测试、限流测试、roles:null 测试、D6 归一化测试）。
- `npx tsc -b`：无类型错误。
- `npx vitest run`：735 测试全过（734 基线 + 1 新增 D4 测试）。
- 候选 commit：`9c1d0a7`（HEAD）。

## 证据

- 修复 commits：`5cbf3dc`（C1–C4）、`e498789`（C5–C8）、`9c1d0a7`（D1–D8）。
- 审查输入：`raw/audit-20260810-api-web-bug-review.md`（gitignored）。
