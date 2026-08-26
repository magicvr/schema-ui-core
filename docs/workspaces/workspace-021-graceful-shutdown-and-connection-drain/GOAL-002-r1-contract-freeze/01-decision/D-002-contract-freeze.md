---
id: D-002
title: 优雅停机 / 连接排空合同 v0.1.0（冻结 · RT-D02 · 单进程基线）
date: 2026-08-27
status: accepted
---

# D-002 · 优雅停机 / 连接排空合同 v0.1.0（2026-08-27 冻结）

> 责任文件（frozen）。实现（R2）与验收（R3）以本合同为分母。本波不实施任何超出本合同的改动；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；不改迁移台账。

## 0. 适用与验收基线

- **单进程**（API 进程内 Job runner）+ **Compose**；协议信号 **SIGINT / SIGTERM**。
- 双方言 Store（**SQLite / PostgreSQL**）同一合同、同一 Close 路径。
- 范围外（仍 trigger-gated / default-non-goal）：A3 多实例、就绪探针扩依赖、PG 锁 vs Redis vs 队列评估、API 与 worker 进程分离（RT-D03）、Job 租约 / leader election（RT-Q04）、外部队列（RT-Q02）、K8s / TLS 终止。

## 1. 停机顺序（强制）

收到 SIGINT / SIGTERM 后按序执行；任一步骤在停机预算（§6）内发生：

1. **开始**：结构化日志 `shutdown.starting`（含信号类型）。
2. **停止接收新请求**：`http.Server.Shutdown` 关闭 listener——新连接直接被拒；`/readyz` 不再 ready。
3. **存量请求排空（grace）**：预算内等待 in-flight 请求完成；预算耗尽 → 强制退出（§3 退出码 1）。
4. **后台面收尾**：operation-log retention sweep 停止；metrics listener 停止（VP-015；`/metrics` 停机后将不可达，属预期）。
5. **运行中 Job 处理**：对所有 active job 发出 ctx 取消 + 工作进程预算内收尾（语义见 §4）。
6. **模块 runtime hooks** 按注册逆序 Stop。
7. **Store 排空**：`db.Close()`（无并发迁移窗口，见 §5）。
8. **Tracing flush**：OTLP exporter 预算内刷完 pending spans（VP-015 已交付）。
9. **结束**：全部无错误 → 结构化日志 `shutdown.complete` → **exit 0**；任一错误或预算耗尽 → `shutdown.timeout|error`（含错误明细）→ **exit 1**。

## 2. HTTP drain 合同

- **拒绝语义**：`http.Server.Shutdown` 关闭监听即拒绝新连接；已有连接（in-flight 请求）在预算内完成，或预算耗尽被截断。
- **预算**：`http.shutdown_timeout`，默认 `10s`（§6）。超时 → Shutdown 返回 context deadline → 强制退出。
- **请求级超时不变**：既有 `read_timeout` / `write_timeout` / `idle_timeout` 语义维持（请求级防护，不属于停机预算）。
- **部署对齐（RT-D01）**：本地双进程与 Compose 路径一致可核对；fork 部署须保证容器编排 stop 宽限 ≥ `http.shutdown_timeout`（本仓 Compose 建议显式 `stop_grace_period: 15s`，R2 落地；否则编排层可能在排空完成前 SIGKILL）。

## 3. 退出码语义

| 码 | 含义 |
|----|------|
| `0` | 排空完成、无错误（§1 全序成功） |
| `1` | 停机错误（任一 join 错误）或预算超时强制退出 |
| 其它 | 启动 / 配置阶段失败（沿用现状；与停机合同无关） |

SIGKILL 无合同（部署层恢复策略负责；Compose `restart: on-failure` 维持）。

## 4. 运行中 Job 停机语义（I-001 · 冻结：**中断标记重跑**）

- 停机触发：所有 **active job** 收到 ctx 取消（合作取消协议维持：`Progress` / `Cancelled()`）。
- 预算内落终态者（succeeded / failed / cancelled）尊重其终态。
- **未收尾者 = 中断标记重跑**：保持 `running`，租约（默认 30s）过期；进程重启后由 runner 认领重跑——`attempt+1`，且 `attempt < max_attempts`（默认 3）；预算耗尽 → `failed / JOB_ATTEMPTS_EXHAUSTED`。
- `queued` 未运行 Job 不受停机影响，重启后正常调度。
- **不冻结**：等完成（Job drain 等待）、按类型分流——留 A3 评估（trigger-gated）。
- 本语义与现行 `jobs.Runner` 行为零漂移（`runner.go`：Stop cancel + lease-reclaim 路径已存在）。

## 5. Store 排空与迁移（I-003 · 冻结：**fail-closed 启动期校验**）

- **无运行时迁移窗口**：迁移只在启动期 `store.Open` 的 startup plan（fresh / noop / apply-pending / restore-ledger）中执行；服务期不存在任何迁移执行入口（`/api/schema/*` 为 schema 文档面，非迁移执行面）。
- **停机期间 Store 语义**：存量查询 / 事务经 §1 步骤 2/3/5 排干后，最后 `db.Close()`。
- **迁移完整性 fail-closed**：一个迁移一个事务 + 全局 ledger checksum；进程若在迁移中途被杀，下次启动由 ledger 校验兜底（版本缺失 / 名称不符 / checksum drift → 拒绝启动）。
- SQLite / PG 双方言一致（同一端口合同、同一 close 路径）；checksum 台账**不变**（本波不改迁移）。

## 6. 配置键与默认（I-002 · 冻结）

| 项 | 值 |
|----|-----|
| YAML | `http.shutdown_timeout`（Go duration，如 `10s`） |
| env | `HTTP_SHUTDOWN_TIMEOUT` |
| 默认 | `10s`（与现行 `main.go` 硬编码一致） |
| 非法值 | `<=0` 或解析失败 → **fail-closed**：拒绝启动（与 `ValidateProd` 风格一致） |
| 显式配置 | 仅作覆盖；缺省即本合同默认 |

## 7. 可观测（I-004 · 口径）

- 结构化日志（走 VP-015 通道）三事件：`shutdown.starting` / `shutdown.complete` / `shutdown.timeout|error`。
- **指标断言不进本波验收分母**（VP-021 首波冻结已声明；metrics listener 停机后不可达属预期）。

## 8. 验收方式（R3 预告）

- **harness A（clean drain）**：起单进程 → 注入存量慢请求 + 运行中 Job → SIGTERM → 断言：新请求被拒、存量在预算内完成、Job 按 §4 落终态或留待重跑、Store 已关闭、**exit 0**；SQLite 与 PG 各跑一遍。
- **harness B（timeout）**：`http.shutdown_timeout=1s` + 慢于预算的请求 → SIGTERM → **exit 1** + `shutdown.timeout` 日志。
- **harness C（重启回收）**：停机时运行中 Job 未收尾 → 重启后按租约 reclaim（attempt+1）直至成功或耗尽。
- Compose 路径：`docker compose stop`（或等宽 signal）可重复；`stop_grace_period ≥ shutdown_timeout` 已落地。
- 迁移台账 checksum 前后不变（回归锁）。

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| Job 语义 | 等完成 / 按类型分流 | 与现行 runner 冲突、阻塞停机预算；A3 gated |
| 超时非法值 | 回落默认 | 静默掩盖配置错误；合同要求可核对 |
| 迁移窗口 | 运行时排队 / 迁移收敛 | 迁移仅在启动期，无并发场景；属 A3 级新机制 |
| 可观测 | 指标断言进分母 | VP-021 首波冻结声明指标不进分母 |
| 配置 | 分面子预算（HTTP/Job/Store 各自） | 单进程基线一个总预算即可核对；细分留 A3 |

---

**引用链**：证据 → `GOAL-002/01-decision/D-001`（信息裁决）；实施责任 → R2（GOAL-003）；验收 → R3（GOAL-004）。