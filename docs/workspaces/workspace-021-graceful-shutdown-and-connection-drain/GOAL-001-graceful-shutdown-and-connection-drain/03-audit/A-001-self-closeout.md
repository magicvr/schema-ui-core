---
id: A-001
source: self
date: 2026-08-27
scope: workspace-021 根目标关门自审（VP-021 退出判据 1～5 · 全链证据 · 边界与台账一致性）
verdict: pass
project: workspace-021 · GOAL-001-graceful-shutdown-and-connection-drain
---

# A-001 · Root 关门自审（2026-08-27 · self）

## 范围

- 对象：Root `GOAL-001-graceful-shutdown-and-connection-drain` 关门；对照 VP-021 方向级退出判据 1～5 与合同 v0.1.0 §1–§8。
- 模式：`cross`（Root 关门：self + grok build independent；provider 按项目决策 independent-audit-execution）。

## 退出判据核对（VP-021）

| 判据 | 证据 | 核对 |
|------|------|------|
| 1. 停机顺序 / 超时 / 退出码合同落盘，单进程 + Compose 下可核对 | 合同 v0.1.0（GOAL-002 D-002）；进程内 harness A/B（本机绿）+ 进程级 A′/B′（linux/CI：exit 0 + `shutdown.complete` / exit 1 + `shutdown.timeout`）；compose `stop_grace_period: 15s` | ✓ |
| 2. 运行中 Job 停机语义冻结并有行为证据 | 合同 §4（用户裁决：中断标记重跑）；`TestShutdownInterruptLeaseReclaim`（Stop 无终态 → reclaim attempt+1 → succeeded） | ✓ |
| 3. 双方言 Store 排空语义一致可核对 | 合同 §5（无运行时迁移窗口 + fail-closed 启动期校验）；SQLite harness 实测 + PG 变体（CI 门控）+ store 包双方言 Open/Close/重启测试；迁移 checksum 回归锁（全量 suite 绿） | ✓（PG 进程级为 CI 责任，见残留） |
| 4. 未进 A3 余项；未改 Charter；未改 Profile 默认集 | git diff 越界核对：本波仅改 compose.yaml、`cmd/server/main.go`、`internal/config/*`、`configs/*`、4 个新增/更新测试文件 + 治理文档；无 Profile/模块矩阵/Manifest/迁移改动 | ✓ |
| 5. 开放 required = 0（或已合法闭合） | Vision VRev-046 pass 0 required；三个子目标审计 0 required；Root 本审 0 required | ✓ |

## 带头核对（合同 ↔ 实现）

| 合同条款 | 实现/证据 | 核对 |
|----------|-----------|------|
| §1 停机顺序 9 步 | OnStop 顺序（既有）+ main.go 接线 + 日志事件（`shutdown.starting` → ... → `shutdown.complete`/`shutdown.timeout`） | ✓ |
| §2 HTTP drain / 部署对齐 | `http.Server.Shutdown` + `http.shutdown_timeout`（默认 10s）+ compose `stop_grace_period: 15s` | ✓ |
| §3 退出码 | main.go：error/timeout → exit 1；clean → exit 0（进程级 harness断言） | ✓ |
| §4 Job 中断标记重跑 | runner 既有语义 + reclaim 证据测试 | ✓ |
| §5 Store × 迁移 fail-closed | R1 裁决 + 契约 + 双测锁 | ✓ |
| §6 配置键 fail-closed | `HTTP_SHUTDOWN_TIMEOUT` 非法/<=0 → LoadError（config 测试锁 7 子测） | ✓ |
| §7 日志事件 | main.go 三事件 | ✓ |
| §8 验收 harness | A/B（本机）+ A′/B′/PG（CI 门控）+ C | ✓ |

## 台账一致性

- goal-tree：Root 2/3 → 3/3 待 R3 关后同步（本审后更新）；子目标链 parent/status 一致。
- 信息台账：I-001～004 全部 verified（用户裁决 R1）。
- workspace.md 绑定：lead VP-021 active；plan_refs/primary_plan 合法。
- Vision：VRev-046 pass；VP-021 信息项 verified。

## Findings

- `F-001`：recommended。进程级 A′/B′（linux）+ PG 变体（PG_TEST_*）为门控证据：本机（Windows、无 PG）无法实跑——列为**有界 residual**：Range = 两文件的 CI 实跑核销；复审触发 = CI 首次 run 或下一架构 VP 激活前；接受者 = 用户书面（见 Root 关门记录）。→ **accepted-residual（用户确认后）**。
- `F-002`：recommended。`bye` 打印（main.go 末尾）为 legacy 输出，非合同条目；建议 R3 后保持现状不删（避免无关 churn），留待后续清理。→ 记录不处理。

## 结论

**pass**（0 required；2 recommended → 1 项 residual 待用户书面接受、1 项记录不处理）。Root 关门条件满足：五件套链条完整、证据可核对、越界为零、required 归零。允许在 grok independent 核销后标 Root `done`，并向 `/vision` 提 VP-021 关门提案。