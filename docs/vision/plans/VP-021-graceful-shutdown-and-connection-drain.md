---
doc_type: vision-plan
id: VP-021-graceful-shutdown-and-connection-drain
title: 优雅停机 / 连接排空合同
status: closed
vision_ref: schema-ui-core-admin-foundation@0.3.0
lead_workspace: workspace-021-graceful-shutdown-and-connection-drain
created: 2026-08-26
updated: 2026-08-27
version: 0.3.0
parent: null
---

# VP-021 · 优雅停机 / 连接排空合同（RT-D02）

## 状态与门闩（2026-08-27 · **closed**）

| 项 | 值 |
|----|-----|
| status | **`closed`**（2026-08-26 用户确认立项 → 2026-08-27 激活 → **2026-08-27 关门** v0.3.0 · VRev-047 self `pass` · 用户指令授权） |
| **lead_workspace** | `workspace-021-graceful-shutdown-and-connection-drain`（2026-08-27 `/govern` 开区并结项；Root `GOAL-001-graceful-shutdown-and-connection-drain` `done 3/3`） |
| **Vision required** | 无（VRev-046 intent-activation self `pass`；VRev-047 closeout self `pass` · 0 required） |
| **推进门闩** | 激活前置满足 = **架构类 freshness PASS**（`ed99e88` → `fddaf638`，不暂挂 `go`）；R1 方案冻结前关闭 I-021-001/I-021-002，R2 前关闭 I-021-003（全部 `verified`） |
| **组合位置** | 架构分支 · RT-D02（`registered` → 本 VP 承接 → `active` → **`delivered`**，2026-08-27） |
| **完整 ≠ A3** | 只做**单进程基线**的优雅停机 / 连接排空合同。A3 多实例、就绪探针扩依赖、PG 锁 vs Redis vs 队列评估仍 `trigger-gated` |

## 意图

把现行「进程生命周期有、但无明确 drain 合同」的后端收成可核对的**优雅停机 / 连接排空合同**（架构 RT-D02），以单进程 + Compose 为验收基线：

1. **停机顺序合同**：停止接收新请求 → 存量请求排空（grace period）→ Store / 连接排空 → 超时与退出码语义可核对。
2. **HTTP drain**：`http.Server` 关闭语义合同化；与既有本地双进程 / Compose 路径（RT-D01 `delivered`）对齐。
3. **运行中 Job / 后台任务语义**：停机时运行中的 Job（VP-012 六态）行为有明确定义（等完成 / 中断标记重跑，R1 冻结其一）。
4. **Store 连接排空**：在 VP-013 双方言 Store 之上，连接关闭顺序与迁移窗口重叠时的停机语义合同化；SQLite / PG 一致，checksum 台账不变。
5. **内嵌默认**：无多实例、无外部代理仍能在 dev / 快测 / Compose 核对（SIGTERM / SIGINT → 排空 → 退出码）。

本 VP 属 **架构分支**。它是注册项 RT-D02 的落地，与 VP-009 / VP-010 正交；不承接 A3 的其余内容。

## 配置面与模块归属

- 走既有进程生命周期代码路径（`cmd` / 进程入口），**不是**新模块、不改 Profile 默认集。
- **缺省**：默认顺序与超时即合同；显式配置仅作覆盖（键名由 lead Root R1 冻结）。
- **生产 / 本 VP 验收**：单进程 / Compose 下可核对信号 → 排空 → 退出码；不要求多实例。
- **生效方式**：随进程启动生效。热加载不进退出分母。

## 首波冻结（退出分母 = RT-D02 合同）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| 停机顺序 | 新请求停止 → 存量排空 → Store 排空 → 超时 / 退出码合同 | A3 多实例 / 水平扩展；API 与 worker 进程分离（RT-D03 仍 gated） |
| HTTP drain | `http.Server` Shutdown 合同化（grace / 超时 / 退出码可核对） | TLS 终止（RT-D05 仍 gated）；热加载；零停机部署产品 |
| Job 语义 | 运行中 Job（VP-012 六态）停机行为冻结（等完成 / 中断标记其一，I-021-001） | Job 租约 / leader election（RT-Q04 仍 gated）；外部队列 |
| Store 排空 | 双方言连接关闭顺序 + 迁移窗口重叠时停机语义，SQLite / PG 一致 | 连接池改造（RT-P04 仍 gated）；PITR |
| 可观测 | 停机事件走既有结构化日志与 correlation（VP-015 已交付） | 新监控产品 / Admin 监控页 |

## 非目标

- **A3 余项**：多实例、水平扩展、就绪探针扩依赖、PG `SKIP LOCKED` vs Redis vs 队列评估（仍 `trigger-gated`；本 VP 不等待它们，也不把它们拉进来）
- **API 与 worker 进程分离**（RT-D03）、**分布式锁 / Job 租约**（RT-Q04）、**外部消息队列**（RT-Q02）
- **K8s / Helm / Operator**（RT-D06 `default-non-goal`）、**TLS 终止**（RT-D05）
- 热配置 / 零停机部署产品；PITR；改 Profile 默认集 / 模块矩阵 / Manifest 装配语义
- 重开 VP-012 / VP-013；替代 VP-009 / VP-010；改变 Charter 边界；业务域

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003** | 遵守薄内核。停机合同是进程生命周期能力，不是模块级平行注册 |
| **VP-012** | 消费其 Job 六态信封定义停机时 Job 行为；不重开横切契约 |
| **VP-013** | 在双方言 Store（A1 已 closed）之上定义排空顺序；不重做迁移台账 |
| **VP-015** | 停机事件复用结构化日志 / correlation；不新造观测面 |
| **VP-009 / VP-010** | 安全与符合性 finding 仍归持续程序，本 VP 不替代 |
| **A3（未立项）** | 本 VP 与 A3 解耦：合同以单进程为基线；A3 触发后复用本合同扩展 |
| **业务域** | 不引入领域长任务语义；领域 Job 的停机行为走本 VP 合同 |

## 方向级退出判据

1. 停机顺序 / 超时 / 退出码合同落盘，且单进程 + Compose 下可核对（信号测试或等价 harness）。
2. 运行中 Job 的停机语义已冻结并有明确行为证据。
3. 双方言（SQLite / PG）Store 排空语义一致可核对。
4. 未进 A3 余项；未改 Charter；未改 Profile 默认集作为本波成功条件。
5. 开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root `GOAL-001-graceful-shutdown-and-connection-drain`（P-001）书写：R1 合同冻结（顺序 / 超时 / Job 语义）→ R2 实现与测试 → R3 证据。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

允许带未知立项。下列不影响「本 VP 意图已冻结」，但必须在对应阶段前关闭或经用户接受残余。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-021-001 | 停机时运行中 Job 语义：等完成 vs 中断标记重跑（或二者按 Job 类型分流）。 | required | 方案冻结 | R1 合同冻结 | **verified**（2026-08-27 用户裁决：**中断标记重跑**；合同 v0.1.0 §4 = workspace-021 GOAL-002 D-002） |
| I-021-002 | grace period / 超时默认值与可配置键（含超时后的强制退出语义）。 | required | 方案冻结 | R1 合同冻结 | **verified**（2026-08-27 用户裁决：默认 `10s` + `http.shutdown_timeout` / `HTTP_SHUTDOWN_TIMEOUT`，非法值 fail-closed；合同 §6） |
| I-021-003 | Store 排空与迁移窗口重叠时的停机语义（fail-closed？排队？）。 | required | 方案冻结 / 实施 | R2 | **verified**（2026-08-27 用户裁决：**fail-closed 启动期校验**，无运行时迁移窗口；合同 §5） |
| I-021-004 | 停机是否需日志 / 指标断言（消费 VP-015 已交付面）。 | non-blocking | 验收 | R3 | **verified**（lead 口径：结构化日志三事件断言；指标不进分母；合同 §7） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-021-graceful-shutdown-and-connection-drain | GOAL-001-graceful-shutdown-and-connection-drain | lead（delivery） | 2026-08-27 | `active` 1 区；开区经 `/govern`（Root 五件套 + P-001 纲领 R1～R3 + I-021-001～004 投影台账） |

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-27 | **closed**（v0.3.0 · `active → closed`） | 退出判据 1～5 全部成立（VRev-047 self `pass` · 0 required）：合同 v0.1.1 冻结（含 §9 Shutdown 状态机勘误）+ 进程内 harness A/B 实测 + 进程级 A′/B′（linux/CI）+ Job 中断标记重跑实测（reclaim attempt+1）+ 双方言（SQLite 实测 + **PG drain 实测 PASS**）+ 迁移 checksum 回归锁；lead workspace-021 Root `done 3/3`；关闭双审 A-001 self `pass` + A-002 grok independent `conditional` → required fixed → 0 开放；越界为零（未进 A3 / 未改 Charter / Profile）。RT-D02 → **delivered** | `docs/workspaces/workspace-021-graceful-shutdown-and-connection-drain/`（Root E-006 结项 + GOAL-002/003/004 五件套）；closeout 审计 `GOAL-001/03-audit/A-001`（self）+ `A-002`（grok independent）；VRev-047 本报告 | 有界 residual：进程级 SIGTERM harness（`cmd/server/shutdown_harness_test.go` · `!windows`）与 `docker compose stop` 实跑以 **linux CI** 核销（V-F083 登记；复审触发 = CI 失败或下一架构 VP 激活前；闭合 = CI 证据或用户书面 accepted-residual） |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-26 | 初创 `planned`：用户确认立项（架构分支 RT-D02 · 优雅停机 / 连接排空合同）；退出分母 = 单进程基线合同；A3 余项、Job 租约、进程分离仍 gated 不进。roadmap 索引原子同步 |
| 2026-08-27 | **v0.2.0 `planned → active`**：VRev-046（self）`pass`（0 required；V-F081/V-F082 → **fixed**）；架构类 freshness **PASS**（`ed99e88` → `fddaf638`，不暂挂 `go`）；lead `workspace-021-graceful-shutdown-and-connection-drain` 同日开区（Root `GOAL-001-graceful-shutdown-and-connection-drain`，P-001 纲领 R1～R3 + I-021-001～004 台账）。VR-048 |
| 2026-08-27 | **v0.3.0 `active → closed`**：VRev-047（self）`pass`（0 required；V-F083 → **fixed**）；退出判据 1～5 全部成立（lead Root `done 3/3`；关闭双审闭合；PG drain 实测 PASS）；RT-D02 → **delivered**；关门记录与有界 residual 落盘。VR-049 |