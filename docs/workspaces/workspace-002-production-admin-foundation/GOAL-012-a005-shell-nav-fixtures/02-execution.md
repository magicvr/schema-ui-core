---
title: 执行 · Shell 导航与 Schema fixture 洁净度
status: done
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-001-production-admin-foundation
version: 0.2.0
---

# 执行 · GOAL-012

## 时间线

### 2026-08-04 · 立项

- 独立代码审计 A-005（未使用 skill）确认：`app-manifest.json` 声明 `activity`、`settings`，但 `apps/api/internal/handler/fixtures/schema/` 仅有 `overview`、`data-table`、`search-form-table`、`form-controls`、`form-with-reactions`、`users`、`roles`。
- 用户指令：阻断问题 → 新设子目标修正 + Root 回退关门。
- 本目标创建为 `active / 0/4`；尚未修改产品代码。

### 2026-08-04 · S1～S5 实施

- **D-002**：S1 采用移除占位（非补假页）。
- **S1 · 消除死链**：
  - 更新 `apps/web/public/.well-known/schema-ui/app-manifest.json`：删除 `activity`/`settings` pages；删除 sidebar Workspace 组；`user` 置 `[]`。
  - 同步 `apps/web/src/protocol/app-manifest.test.ts` 期望 pageId 列表。
  - 同步 `STATIC_MANIFEST_SHA256` = `a30eb1c1fb368d20663e2d657d9421df39f4b0f1ff39e0f2f000559c38a54a83`（`upstream-fixtures.test.ts`）。
- **S2 · 一致性门禁**：
  - `apps/api/internal/handler/schema_test.go` 新增子测 `checked-in web manifest pageIds all have embed fixtures`：读取 web checked-in manifest，断言每个 `pageId` 有 embed fixture 且 `GET /api/schema/{id}` 200；并强制 `schemaUrl == /api/schema/{pageId}`。
- **S5 · QUICKSTART（R-001）**：§4 改为真实落点 `fixtures/schema/*.json` + rebuild API + manifest 挂载；明确 `docs/schemas/` 是协议 schema 非页面文档。
- **S3 · 回归**：
  - `cd apps/api && go vet ./...` 干净；`go test ./... -count=1` 全绿（含新 schema 子测）。
  - `cd apps/web && npx vitest run` **491/491**；`tsc -b` 干净。
- **S4 · 关门**：见 03-audit A-001（self · close-out · pass）；Root A-005 F-001 → `fixed`。
