---
id: E-001
title: 立项：R1 契约基线扫描（main.go / OnStop / runner / store / compose 事实）
date: 2026-08-27
status: done
---

# E-001 · GOAL-002 立项与基线扫描（2026-08-27）

## 事实

1. **立项**：Root 纲领 R1 → `GOAL-002-r1-contract-freeze`（五件套 + ledger 目录 + 检查点 C1～C3 + 信息台账）。goal-tree 追加子树；Root / workspace.md 索引同步。
2. **基线扫描（只读，未改代码）**：
   - `cmd/server/main.go`：`signal.NotifyContext(SIGINT, SIGTERM)` → `app.Stop(10s ctx)`；错误 exit 1，成功 exit 0。
   - `composition.registerLifecycle` OnStop 顺序：`srv.Shutdown` → listener.Close → retention sweep stop → `metrics.Stop` → `jobs.Stop` → `runtime.Stop` → `st.Close` → `tracing.Shutdown`；错误 join。
   - `jobs.Runner.Stop`：cancel 所有 active job ctx + workers.Wait（受 ctx 限制）；`finish` 在 stopping 时放弃 durable 转移；重启经 lease-reclaim（attempt+1）或 attempts 耗尽失败。
   - `store.Open`：迁移仅启动期（fresh/noop/apply-pending/restore-ledger；一迁移一事务 + checksum）；服务期无迁移入口；`Close` = `db.Close`（SQLite/PG 同）。
   - `server.New`：`http.Server` 超时全部配置驱动（`http:` 段）。`compose.yaml`：`restart: on-failure`，无 `stop_grace_period`。
3. **产出**：D-001 信息裁决提案（I-001/I-002/I-003 证据 + 建议 + I-004 口径），交用户裁决（P-004）。

## 验证 / 后续

- C1 等待用户裁决。裁决后：D-001 → accepted，I-00N → verified，合同正文 D-002 落盘（C2）。