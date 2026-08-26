---
id: D-001
title: 信息裁决：I-001 / I-002 / I-003（用户 2026-08-27 采纳建议）
date: 2026-08-27
status: accepted
---

# D-001 · 信息裁决（2026-08-27 · P-004 / P-005）

> 2026-08-27 用户裁决：三条 required 全部**采纳建议**（界面裁决记录）：I-001 = 中断标记重跑；I-002 = 默认 10s + `http.shutdown_timeout` / `HTTP_SHUTDOWN_TIMEOUT`（非法值 fail-closed）；I-003 = fail-closed 启动期校验（无运行时迁移窗口）。I-004（non-blocking）= lead 口径：结构化日志断言，指标不进分母。对应 I-001/002/003 → `verified`；I-004 → `verified`（lead 提案，符合 VP-021 首波冻结）。合同正文见 D-002。

## 意图

R1 方案冻结（合同正文）前，关闭 C1 的三条 required 信息项。以下全部基于**现行实现事实**（只读扫描，未改任何代码），每项给出建议；采纳后本条目转 `accepted`，对应 I-00N 转 `verified`。

## I-001 · 运行中 Job 停机语义（required）

**现行事实**（`internal/jobs/runner.go`）：

1. `Runner.Stop(ctx)`：置 `stopping`、close `stop`、**逐个 cancel 所有 active job 的 context**，然后 `workers.Wait()` 或等待 `ctx`。
2. `finish()`：若 `isStopping()` 为真 → **直接返回**（不做 durable 终态转移）；被中断的 Job 保持 `running`，其租约自然过期。
3. 重启后（`ScanOnce`）：`ListRunnable` 以 `status='running' AND cancel_requested=0 AND lease_expires_at <= ? AND attempt < max_attempts` 重新认领（**attempt+1**）；`ExhaustExpiredJobs` 对 `attempt >= max_attempts` 的置 `failed / JOB_ATTEMPTS_EXHAUSTED`。
4. 六态（queued / running / succeeded / failed / cancelled / expired）与 `DefaultMaxAttempts = 3`（`model.go`）。

**结论**：现行可持久语义 = **中断标记重跑**（interrupt + lease-reclaim，attempt 计入重试预算，耗尽即 failed）。

**建议（采纳即冻结）**：合同定为 **中断标记重跑**——

- 停机 = 向所有运行中 Job 发出取消（ctx cancel）+ 工作进程在停机预算内收尾；
- 未在停机窗口内完成/落终态的 Job：保持 `running`、租约过期，重启后在重试预算（默认 3 次）内**重跑**（attempt+1）；预算耗尽 → `failed / JOB_ATTEMPTS_EXHAUSTED`；
- 窗口内明确完成或按取消协议落终态（succeeded / failed / cancelled）的 Job 尊重其终态；
- **等完成（drain 等待 Job 全部完成再关 Store）不冻结为默认**：会阻塞停机预算，且与现行 runner 语义冲突；留待 A3 多实例评估（trigger-gated）。
- 不新增 Job 类型分流；`Progress`/`Cancelled()` 协议维持（handler 合作取消）。

## I-002 · grace / 超时默认值与配置键（required）

**现行事实**：

1. `cmd/server/main.go`：`shutdownCtx` 硬编码 **10s**；`app.Stop` 出错 → `os.Exit(1)`，否则 exit 0（`fmt.Println("bye")`）。
2. `internal/server/server.go`：`http.Server` 已有 `ReadTimeout` / `WriteTimeout` / `IdleTimeout`（配置驱动，`http:` 段）。
3. `internal/config`：段式 YAML + `HTTP_*` 环境变量覆盖（如 `HTTP_ADDR` / `HTTP_TRUSTED_PROXIES`）。
4. `compose.yaml`：无 `stop_grace_period`（docker compose stop 默认 ~10s）。

**建议（采纳即冻结）**：

- **停机总预算默认 10s**（与现状一致），含义 = 信号收到 → `http.Server.Shutdown` grace 开始计数 → 预算内完成排空 → 预算耗尽未完成即强制退出；
- **配置键**：YAML `http.shutdown_timeout`（duración Go，如 `10s`）+ 环境变量 `HTTP_SHUTDOWN_TIMEOUT`；缺省 = 10s；非法/<=0 值 fail-closed（拒绝启动）或回落默认（二选一，建议 **fail-closed**，与 `ValidateProd` 风格一致）；
- **退出码语义**：drain 成功（无错误）→ exit **0**；drain 错误/预算超时 → exit **1**（与现状一致并合同化）；信号处理覆盖 SIGINT + SIGTERM；
- 分面计时（HTTP vs Job vs Store 各自子预算）**不冻结**为默认配置——单进程基线一个总预算即可核对；细分留给 A3。

## I-003 · Store 排空 × 迁移窗口（required）

**现行事实**：

1. 迁移**只在启动期**执行：`store.Open` → startup plan（fresh / noop / apply-pending / restore-ledger），一个迁移一个事务（`migrate.go` `applyMigration`），ledger 带 checksum；`postgres.go` 同合同（`actionApplyPending`）。
2. 服务期**无**运行时迁移入口（`/api/schema/*` 是 schema 文档面，不是迁移执行面）；`MaxOpenConns=1` 单连接。
3. `Store.Close()` = `db.Close()`（SQLite / PG 同）。
4. 停机顺序现状：`srv.Shutdown` → metrics/jobs/runtime stop → `st.Close()`。

**结论**：**迁移窗口与停机在运行时不可能重叠**——迁移只发生在「未监听 / 未就绪」的启动期；进程若在迁移中途被杀，账本按事务 + checksum 完整性在下次启动校验（fail-closed：drift 即拒绝启动）。

**建议（采纳即冻结）**：

- 合同写明：**无运行时迁移窗口**；停机期间唯一 Store 语义 = 存量查询/事务（经 `srv.Shutdown` 与 Job 收尾）完成后 `db.Close()`；
- 迁移完整性 = **fail-closed**（启动期 ledger 校验 + 单迁移事务幂等重试由启动计划保证）；不存在「排队」场景；
- SQLite / PG 双方言对上述语义**一致**（同一端口合同、同一 close 路径）；checksum 台账不变（本波不改迁移）。

## I-004 · 日志 / 指标断言（non-blocking）

**建议**：验收面 = 结构化日志断言（停机开始/完成/超时三行，走 VP-015 已交付的结构化日志与 correlation 通道）；**指标断言不进本波分母**（VP-021 首波冻结已声明）。

## 裁决结果（2026-08-27 · 用户界面裁决，全部采纳建议）

- **I-001 → `verified`**：中断标记重跑（上方建议原样冻结）。
- **I-002 → `verified`**：默认 10s + `http.shutdown_timeout` / `HTTP_SHUTDOWN_TIMEOUT`，非法值 fail-closed；退出码 0（clean）/ 1（error|timeout）。
- **I-003 → `verified`**：无运行时迁移窗口；迁移完整性 fail-closed = 启动期 ledger 校验；停机 Store 语义 = 排干后 `Close`；双方言一致。
- **I-004 → `verified`（lead 口径）**：结构化日志三事件断言（shutdown.starting / complete / timeout|error）；指标断言不进分母。

合同正文 = `01-decision/D-002-contract-freeze.md`（C2）。