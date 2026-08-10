---
id: GOAL-002-audit-findings-remediation
title: 审查发现修正（第一个子目标）
status: active
parent: GOAL-001-production-hardening
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
progress: 0/16
---

# GOAL-002 · 审查发现修正（第一个子目标）

## 概述

本子目标承接 2026-08-10 代码审查发现的**全部 16 项**共享基架安全/健壮性缺陷（C1–C8 中高危 + D1–D8 低危）的修复与回归。输入清单见 `raw/audit-20260810-api-web-bug-review.md`（gitignored 临时记录，审查当日产物）。

## 范围（16 项）

| id | 严重度 | 位置 | 问题 |
|----|--------|------|------|
| C1 | 高危 | `apps/api/internal/handler/upload.go:126,152-160` | 上传端点存储型 XSS → 账号接管 |
| C2 | 中危 | `auth.go:99-107` + `accounts.go:238-256` + `auth-client.ts:83-110` | 刷新竞态：并发兑换双会话 / 有效会话误登出 |
| C3 | 中危 | `config.go:43` + `main.go:80-83,93-95` | `APP_ENV` 默认 development fail-open |
| C4 | 中危 | `composition.go:112-118` | 管理员种子 Bootstrap 一次性锁死 |
| C5 | 中危 | `render.tsx:1114-1140,572,606,669` | 异步错误无 catch：提交按钮卡死、动作静默失败 |
| C6 | 中危 | `render.tsx:733-744` | 清空搜索框无法清除过滤 |
| C7 | 中危 | `permissions.ts:531-536` vs `render.tsx:458-465` | 未声明权限的 action 被拦截 |
| C8 | 中危 | `App.tsx:391-398` + `render.tsx:894-899` | 路由 query 未送达渲染层 |
| D1 | 低危 | `users.go:163-172` + `users_repository.go:137-162` | `PATCH roles: null` 清空角色 |
| D2 | 低危 | `auth.go:75-87` | 登录无速率限制 + 用户名枚举时序侧信道 |
| D3 | 低危 | `composition.go:205` + `schema.go` | schema 文档匿名可读（需显式确认） |
| D4 | 低危 | `request-construction.ts:46-66` | `splitUrl` 畸形百分号编码 URIError |
| D5 | 低危 | `migrate.go:236,54-63,225-244` | 升级快照秒级粒度 + VACUUM 拷贝 |
| D6 | 低危 | `repository.go:58` vs `configuration.go:84-88` | `defaultLocale: ""` 两层校验不一致 |
| D7 | 低危 | `form-controls.tsx:344-347,451-481` | inputNumber 清空变 0；上传无法重传 |
| D8 | 低危 | `theme-toggle.tsx:17-21` | theme-toggle 未同步 color-scheme |

## 成功标准

- [ ] C1–C8 全部修复并回归（上传 MIME 服务端校验 + 下载 `Content-Disposition`/CSP；刷新轮换原子化 + in-flight 去重；`APP_ENV` fail-closed；Bootstrap 空表重试；异步 catch 全覆盖；搜索清空语义；权限缺省契约统一；路由 query 贯通）
- [ ] D1–D8 修复或按用户 P-004 裁决（accepted-residual / 明确非目标）；D3 匿名可读属性经显式确认
- [ ] 新增/更新回归测试覆盖每项修复；`go test ./...` + `vitest run` 全绿；基线不回归
- [ ] 共享基架重验证证据落盘，VP-008 `go` 消费有效性按规则恢复（或用户书面裁决继续暂挂）

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 16 项缺陷是否全部覆盖、修复顺序与受影响测试范围 | 实施 | 实施前 | 逐项核对清单与修复映射 | open | — | 待确认 |
| I-002 | non-blocking | D3（schema 匿名可读）是否设计决策 | 验收 | 验收前 | 用户 P-004 确认 | open | — | 待确认 |

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录；索引文件与目录条目共同构成正式记录。
