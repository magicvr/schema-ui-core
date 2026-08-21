---
id: E-003-w10-s4-closeout-implementation
goal: GOAL-010-w10-api-web-security-audit
status: done
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-010-w10-api-web-security-audit
version: 0.1.0
---

# E-003 · W10 S4 收尾：A-003 recommended ×3 修复 + 滞后索引同步（2026-08-21）

## 事实

1. 用户 `/govern` 指令："响应 GOAL-010 A-003：将 F-001/F-002/F-007 标 fixed、F-003～F-006 按 D-003 作废闭合，同步滞后索引。修正审计报告提出的 recommended 意见，然后关门并恢复go宣称。"
2. A-003（grok-build · grok-4.6 · independent · **pass**）三条 recommended 全部修正：
   - **A-003·F-001（索引滞后）**：`01-decision.md` I-002 行、`03-audit.md` 信息核对表 I-002 行、`workspace.md` W10 行全部同步至 S3 后事实（本条 + 关门写入一并完成）。
   - **A-003·F-002（listener 清理）**：`lib/fetch-timeout.ts` finally 中 `removeEventListener("abort", relayAbort)`；新增测试「removes the caller-signal relay listener after the request settles」（spy 断言 add/remove 同一 handler 各一次）。
   - **A-003·F-003（opener 置空纵深）**：`render.tsx` library.preview 在持有窗口引用后立即 `previewWindow.opener = null`（注释引用 A-003）；`download-behavior.test.tsx` 预览用例扩展 `opener` 哨兵断言置空。
3. 回归：web `npx vitest run` **76 files / 1084 tests 全绿**（+1 listener 用例）；`npx tsc -b` exit 0。Go 本轮零改动，A-003 已复跑 `go vet ./...` exit 0 + `go test ./...` 全绿（含 `-count=1`），不重复。

## 产物

- `apps/web/src/lib/fetch-timeout.ts`（+test）、`apps/web/src/renderer/render.tsx`、`apps/web/src/renderer/download-behavior.test.tsx`
- 本目标全部索引/状态文件关门同步（见各文件 updated 版本）

## 后续

- A-004 闭合记录 + D-004 go 宣称恢复（同日落盘）；密码轮换仍为用户侧残余（A-002/A-003 一致定性）。