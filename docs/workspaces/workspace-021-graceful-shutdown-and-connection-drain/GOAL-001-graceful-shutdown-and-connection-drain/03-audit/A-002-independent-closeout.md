---
id: A-002
source: independent
auditor: grok-build（grok-4.6 · reasoning high）
date: 2026-08-27
scope: workspace-021 根目标 GOAL-001-graceful-shutdown-and-connection-drain 关门独立审（VP-021 退出判据 1–5 · 合同 v0.1.0 §1–§8 · R1–R3 全链）
verdict: conditional
project: workspace-021 · GOAL-001-graceful-shutdown-and-connection-drain
---

# A-002 · Root 关门独立审（2026-08-27 · independent · grok-build grok-4.6 high）

> 誊入说明：本条由编排器自本地 grok build（grok-4.6 · reasoning high）headless 会话原样誊入（会话运行于 2026-08-27；grok 按指令只出报告文本、未落盘——落盘与索引由编排器完成，`source: independent` 保持不变）。grok 当场独立复跑：`go vet ./...` exit 0；config/jobs/composition/store `-count=1` exit 0；`TestHTTPShutdownTimeout` 7 子测 PASS；`TestShutdownInterruptLeaseReclaim` PASS；`TestShutdownDrainHarness` A PASS / B PASS（~2.0s deadline）；**PG 变体 `TestShutdownDrainHarnessPostgres` 实测 PASS（非 skip）**；`cmd/server` 仅 `TestServerProcessRestartPersistsUsers` 编译（`shutdown_harness_test.go` 为 `//go:build !windows` 构建排除）；全量 `go test ./... -count=1` exit 0 无 FAIL。

## 核定结论

- **verdict：conditional**（0 required 可核销路径明确但当下 2 条 required 待编排器响应闭合）。
- 实现面（停机顺序、drain 预算、配置 fail-closed、Job 中断重跑、双方言 Close、checksum 未改、越界为零）**足以支撑 VP-021 合同主体**。
- 关门面：台账自相矛盾（F-001）、GOAL-004 关闭声明超前（F-002），闭合前**禁止 Root `status: done`**。

## 核对摘要（独立复跑事实）

| 区间 | `9e9a8979`（激活/开区）→ HEAD `c974f2d0` |
|------|------|
| 越界 | 代码面仅 `main.go`、`shutdown_harness_test.go`（新）、`internal/config/*`、`configs/*`、`internal/composition/shutdown_drain_test.go`（新）、`internal/jobs/shutdown_reclaim_test.go`（新）、`compose.yaml`。**无** Profile/模块矩阵/Manifest/迁移账本/运行时迁移入口/新 API 路由 |
| 退出判据 1 | 顺序/超时有实证（A/B 绿）；退出码靠 `main.go` 审查 + 进程级 A′/B′（linux 构建，本 OS 未跑）→ **部分成立** |
| 退出判据 2 | 合同 §4 + runner 语义 + `TestShutdownInterruptLeaseReclaim` PASS → **成立** |
| 退出判据 3 | SQLite A + **本会话 PG drain PASS** + store 无 diff + checksum 冻结绿 → **成立（强于自审「本机 skip」主张）** |
| 退出判据 4 | 越界核对零 → **成立** |
| 退出判据 5 | 子目标 self 0 required；本意见新增 2 med required → **本审后不成立**（响应后闭合） |

## Findings

- `F-001`：**required（med）** · 台账与 goal-tree 不一致——GOAL-004 `00-meta` 已 `done 3/3` 而 goal-tree 仍 `active 0/3`、Root meta 仍「R3 待立项（GOAL-004）」、GOAL-004 `01-decision.md` / `02-execution.md` frontmatter 仍 `active`、Root `03-audit.md` 索引 A-002 指向**尚不存在**的文件。建议一次性对齐（AGENTS §7）。
- `F-002`：**required（med）** · GOAL-004 C3 / E-003 关闭声明超前——C3 写「Root done；VP-021 关门提案」且标已关门、E-003 写「Root 纲领进度 3/3」；事实：独立审本轮才发生且未落盘。建议改写并对齐。
- `F-003`：**recommended（med）** · 进程级退出码/日志未在本 OS 核销（`shutdown_harness_test.go` 为 `!windows` 构建排除，非 `t.Skip`；`docker compose stop` 未实跑）。建议：范围 = linux CI 首次绿（可选一次 compose stop）；复审触发 = CI 失败或下一架构 VP 激活前；闭合 = CI 证据 → fixed 或用户书面 accepted-residual。未闭合前不把「退出码可核对」写成已在本环境实证。
- `F-004`：**recommended（low）** · harness A 弱于合同 §8 字面——未断言 Shutdown 后新请求被拒 / `/readyz` 不可达 / Store 已关 / 进程 exit 0。建议补一条 Shutdown 后 `Dial`/`readyz` 失败断言。
- `F-005`：**recommended（low）** · `shutdown.starting` 未含信号类型（合同 §1 步骤 1）。建议 `signal.Notify` 信道 + `signal` 字段，或修订合同删括号。
- `F-006`：**recommended（low）** · `net/http.Shutdown` 状态机与合同措辞精度——`StateIdle`（及 >5s 无首包头的 `StateNew`）立即关；`StateActive`/近期 `StateNew` 预算内等待；超时后 `os.Exit(1)` 截断。建议合同改为该精确表述，实现不必改。

## 风险点专项（a–e，摘要）

- a) harness 探测读规避 accept 竞态：B（budget hole）~2.0s deadline 行为为强证；「已入 handler body 读取」场景假阴性已被锁住；未覆盖纯 StateNew 窄窗口（F-006 措辞层面）。
- b) Shutdown 状态机 vs 「存量请求排空」：不构成行为性违约，属合同精度问题（F-006）。
- c) Windows skip 可接受为有界残余（D-001 预登记）；PG 门控残余**收窄**——PG drain 本会话实测 PASS。
- d) `shutdown.timeout` vs `shutdown.error` 分支与合同一致；`shutdown.error` 无进程级断言为 recommended 覆盖缺口。
- e) 双方言证据链充分（同一 Close 路径 + SQLite/PG 双实测 + checksum 冻结）。

## 必改项汇总

- F-001 / F-002（required）：对齐台账并收回超前声明；未 fixed 前不得将 Root 标 `done`。
- F-003～F-006（recommended）：F-003 建议用户书面残余或 CI 核销；F-004～F-006 可进后续清理，不阻断本波。

## 声明

- `source: independent`；本轮**未修改、未创建任何文件**（按用户指示由编排器落盘）。
- 不修改 `status` / `progress` / goal-tree / 方案正文；响应与放行由 `/govern` 处理。

---

## 编排器响应（2026-08-27 · /govern · P-003）

| finding | 闭合路径 | 证据 |
|---------|----------|------|
| **F-001**（required） | **fixed** | goal-tree 树/表/文首、Root `00-meta`（R3 已关门 / progress 3/3 / status done）、workspace.md、GOAL-004 `01-decision.md`/`02-execution.md` frontmatter 一次性对齐；Root `03-audit.md` 索引 A-002 落盘本文件。commit（见 Root E-006） |
| **F-002**（required） | **fixed** | GOAL-004 C3 改写为「证据自审已过；Root 关门双审另记 Root 03-audit」；E-003 删除「Root done / 3/3」超前句（修订注留痕） |
| F-003（recommended） | **登记残差（不阻断）** | 范围 = `cmd/server/shutdown_harness_test.go` 在 linux CI 首次绿（可选一次 compose stop）；复审触发 = CI 失败或下一架构 VP 激活前；闭合 = CI 证据 → fixed 或用户书面 accepted-residual。已写入 Root 结项记录 |
| F-004（recommended） | **fixed** | harness A 增加 Shutdown 后 `net.Dial` 拒绝断言（`shutdown_drain_test.go`） |
| F-005（recommended） | **fixed** | `main.go` 改 `signal.Notify` 信道，`shutdown.starting` 携带 `signal` 字段 |
| F-006（recommended） | **fixed（合同级）** | GOAL-002 D-002 追加 §9 勘误（v0.1.1 · editorial）：Shutdown 状态机精确表述（拒新连接 / idle 立即关 / StateActive 预算内等 / 超时 exit 1 截断）；原条款按勘误理解 |

**闭合结果**：开放 required = **0**。Root 关门放行成立（连同 A-001 self、子目标审计与全量回归）。