---
id: E-002
goal: GOAL-002-audit-findings-remediation
title: A-002 审计 findings 响应（F-001 fixed + F-002～F-005）
date: 2026-08-10
status: recorded
---

# E-002 · A-002 审计 findings 响应

响应 grok build 独立审计 [A-002](../03-audit/A-002-goal-002-independent.md)（verdict: conditional）的 findings。commit `01b7202`。

## F-001（required）· fixed

**问题**：`http.DetectContentType` 把 SVG 嗅探为 `text/plain`/`text/xml`，填充/夹带 HTML 可绕过拒绝；`dangerousInlineTypes` 拒绝面形同虚设。

**修复**（`apps/api/internal/handler/upload.go`）：
- 新增 `activeContentMarkers` 字节标记启发式：`<svg`/`<SVG`/`<script`/`<SCRIPT`/`<?xml`，任一体内出现即拒绝（`containsActiveContent`）。
- 与 MIME 检测叠加为硬门：`dangerousInlineTypes[base] || containsActiveContent(body)` → 415。
- 大/小写双覆盖；512 字节头后夹带（`strings.Repeat("A", 600)`）与 GIF89a 夹带均被标记命中。

**回归测试**（`upload_test.go`）：`TestUploadRejectsHtmlAndForcesAttachment` 增补 4 个夹带用例（svg-plain / svg-xml / html-padded / gif-html），全部断言 415。

## F-002（recommended）· 已处理

C1 SVG/嗅探旁路专项测试已随 F-001 补齐（上述 4 用例）。

## F-003（recommended）· 已处理

跨标签页并发刷新误清会话：`apps/web/src/account/auth-client.ts` `doRefresh` 失败且 `localStorage` 中 refresh token 已被其他标签轮换时，用新 token 重试一次；再失败才 `clearTokens()`。

## F-004（recommended）· 已处理

`NeedsBootstrap` 专项测试：`systemdata/reconcile_test.go` `TestNeedsBootstrapTracksUserPresence`（空库 true → Bootstrap 后 false）。

## F-005（recommended）· 已处理

C5–C8 前端专项回归：`apps/web/src/renderer/render.test.tsx` 新增 4 测试（网络失败不卡死 / 空搜索覆盖旧过滤 / 无权限声明默认放行 / recordSource route.query 贯通）。

## F-006（recommended）· 运维边界已知

D2 限流为进程内 best-effort（多实例不共享、不信任 X-Forwarded-For），注释与文档已说明，非阻塞。

## 回归验证（2026-08-10）

- `go test ./...`：21 包全绿（含新增夹带用例与 NeedsBootstrap）。
- `npx tsc -b`：无错；`npx vitest run`：739 全过（734 基线 + 5 新增）。
- 候选 commit：`01b7202`（HEAD）。

## 证据

- commit `01b7202`（5 文件，+223/-1）。
- 待 grok 复审确认 F-001 闭合。
