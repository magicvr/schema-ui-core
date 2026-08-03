---
id: GOAL-008-r5-engineering-fork
title: R5 · 工程化、fork 体验与集成关门
status: active
created: 2026-08-02
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.3.0
progress: 5/5
---

# GOAL-008 · R5 · 工程化、fork 体验与集成关门

## 概述

承接 Root **D-012 / D-013**（Root `I-005` = verified、`I-006` = closed）：按部署基线 **A**（文档双进程为默认；**Docker Compose 为 R5 必须交付的第二启动路径**，fork 使用者可选本地双进程或 Compose）、建议 15 分钟 fork 计时口径与复现方法，补齐环境/配置/文档、容器一键启动、fork 上手体验复现与可重复 smoke 验收，并为 Root R5 检查点勾选与关门提供阶段证据。最小操作日志为**可选加分项**（D-013 方案甲），不阻断核心验收。

## 成功标准

- [x] **S1 · 环境与配置基线**：补齐环境变量清单（dev/prod、fail-closed 键、TTL、种子密码策略）、健康检查与启动验证说明（`/healthz` + Web 侧启动验证）、dev/prod 区分文档及工程化 README 段。（2026-08-02 已实施：`apps/api/.env.example` dev/prod 注解；`apps/api/README.md`「开发 vs 生产」+「启动与健康验证」；`apps/web/README.md` 生产/compose 注记；根 `README.md`「工程化与一键启动」；对照契约 C-001/C-002 验证通过）
- [x] **S2 · 容器与一键启动**：交付 Docker Compose（**R5 必须交付和验收的第二启动路径**，非可选加分项；fork 使用者可选择本地双进程或 Compose）——api 多阶段镜像 + web 静态构建（nginx 服务 `dist/` 并把 `/api` 反代到 api service）+ DB volume + 探针；`docker compose up` 后 healthz + 登录 + 代表页可用。（2026-08-02 已实施：`apps/api/Dockerfile` + `apps/web/Dockerfile`（BuildKit cache mount）+ 根 `compose.yaml` + `apps/web/nginx.conf`（SPA fallback + `/api` 反代）+ `.dockerignore` + CI `container-smoke` job；对照契约 C-003～C-007 验证通过——本机 Docker 29/Compose v5 `docker compose up` 后 `/healthz` 200、web 200、登录 admin、`/me`、nginx 代理 + `/list-edit-lifecycle` fallback、api restart 与 down/up 后 DB volume 保持）
- [x] **S3 · fork 文档与 15 分钟体验**：交付 QUICKSTART/README fork 上手文档；≥1 次独立复现记录（日期/仓库 ref/版本/耗时），终点 = 登录成功 + 后台首页可交互（列表加载），耗时 ≤15 分钟（不含依赖下载）。（2026-08-03 已实施：根 `QUICKSTART.md` fork 上手段（双路径 + 终点 4 = `list-edit-lifecycle` + smoke 用法）；**clean-ref 独立复现记录 [R5-S3-REPRO-003](attachments/R5-S3-REPRO-003.md)**（compose 路径，四终点全 PASS，**64.833s ≤ 900s**，run log + 截图附件，**预 T0 禁用 BuildKit 结果缓存、`go build`/`npm run build` 编译层均实际执行非 `CACHED`**——响应 A-012 F-005；取代 `R5-S3-REPRO-002`（BuildKit 编译层 `CACHED`，保留历史）与 `R5-S3-REPRO-001`（build 排除、字段不足，保留历史；响应 A-011 F-005/F-006））
- [x] **S4 · 可复现验收**：交付 `scripts/smoke.sh`（`/healthz` → login → `/api/accounts/me` → 代表页 → 种子可重复），本地与 CI smoke 全绿。（2026-08-03 已实施：`scripts/smoke.sh` SM-001～SM-006 + 退出码 0/2/3/4/5/6/8/70 + `--disposable` 隔离守卫（`SMOKE_ISOLATION_ID` + `SMOKE_DISPOSABLE_CONFIRM=yes` 机器校验，不满足 exit 2；非 disposable 部分绿 exit 8；重启后 readiness 重判）；本地四路径验证（部分绿 8 / fail-closed 2 / 完整绿 0）**SM-006=PASS**（log `r5-smoke-disposable-local-v0.1.2.txt`）；CI `container-smoke` 显式隔离 project `-p ci-container-smoke` 接入 `--disposable`，并有对应成功 run 证据（响应 A-011 F-007/F-008/F-009））
- [x] **S5 · 阶段审计与 Root 关门条件评估**：对 R5 工程化交付做阶段审计（self + 视需要 independent），评估并记录 Root R5 勾选条件与 Root / VP-002 关门证据口径。（2026-08-03 已实施：**A-014（self · finding-closure · pass）** 按 P-004 §3.1 用户裁决补 F-005～F-009 关闭证据 self 覆盖；**A-015（self · stage-audit · pass）** 阶段审计——S1～S5 证据链成立、无开放 required、无意见冲突；**Root R5 勾选条件口径已记录**（S1～S5 全勾选 + 无开放 required + 关门向审计；Root R5 勾选与 Root/VP-002 关门须用户确认）；`I-008-001/002` verified、`I-008-003` 仅 S6 适用）

> 可选加分 **S6 · 最小操作日志**（D-013 方案甲，**非成功标准、不进进度分母**）：若用户决定在本目标内实施，则交付 SQLite `operation_log`（迁移 + repository），覆盖 records 写（create/update/delete）与 auth 关键事件（登录/登出/刷新），并补测试；若不实施，则在本目标 `01-decision` 记录「不纳入本波次」及其理由，作为后续加分项，**不阻断**本目标或 Root R5 关门。

## 派生进度

`progress` 由上方 S1～S5 五个核心检查点等权派生（`0/5` 起）。S6 为可选加分，不进入进度分母；是否实施由用户书面决定（见 `01-decision`）。检查点不替代审计 finding 或关门结论。**2026-08-03**：`I-008-001` 已 verified（D-003 + 契约附件 v1.0.0）；`I-008-002` 已 verified（D-004 + [复现与 smoke 协议 v0.1.2](attachments/I-008-002-fork-reproduction-protocol.md)，D-005/A-008 与 D-007/A-011 响应后权威版本 v0.1.2），解除 S3/S4 的方案冻结/实施前信息门禁；S1/S2 已实施（`0/5 → 2/5`）；**S3/S4 已实施（2026-08-03，`2/5 → 4/5`）**——QUICKSTART fork 文档 + **clean-ref 独立复现记录 `R5-S3-REPRO-003`（预 T0 禁用 BuildKit 结果缓存，`go build`/`npm run build` 编译层实际执行非 CACHED，64.833s ≤ 900s；响应 A-012 F-005）** + `scripts/smoke.sh`（SM-001～006 + 隔离守卫 + 部分绿 exit 8）+ 本地 disposable `SM-006=PASS` 证据（`r5-smoke-disposable-local-v0.1.2.txt`）+ CI `container-smoke` 显式隔离 project 与成功 run 证据。**A-011 F-005～F-009 关闭闭环完成（2026-08-03）**：A-012（independent · fail）重做 F-005 证据 → REPRO-003；**A-013（independent · finding-closure · pass）二次复审确认 F-005～F-009 `fixed` 全部成立、无开放 required**（R-013 handled：CACHED 计数表述收窄至正式 retry #3）；**A-014（self · finding-closure · pass）** 按 P-004 §3.1 补 self 覆盖；**S5 已实施（2026-08-03，`4/5 → 5/5`）**——A-015（self · stage-audit · pass）阶段审计与 Root R5 勾选条件评估。`I-008-003` 仍 open（S6 若实施时阻断）；**本目标五项核心成功标准全部达成（5/5），仍未 `done`**——关门需 close-out 审计 + 用户裁决 + Root R5 勾选（Root 层）。

## 信息需求

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-008-001` | 环境/配置与容器部署的精确契约是什么？（env 键全集与 dev/prod 行为、health/启动验证口径、compose 服务/镜像/DB volume/探针/`/api` 反代形状） | required | S1/S2 实施与验收 | S1 首个文档/配置变更前 | 对照 Root D-013 方案边界 A（Compose = 必须交付的第二启动路径）与现有 `config.go`/`healthz`/web 构建，按 A-001 R-002 最低清单形成版本化工程化契约与验收清单，记录决策 | **verified** | 已关闭（D-003） | [I-008-001-engineering-contract.md](attachments/I-008-001-engineering-contract.md) v1.0.0：env 键全集与 dev/prod、health/启动验证、Compose/镜像/DB volume/探针、SPA fallback 与 `/api` 反代、依赖/超时/失败行为、CI 入口 + 验收清单 C-001～C-007；**不**勾选 S1/S2 |
| `I-008-002` | 15 分钟计时复现协议与 smoke 验收判据是什么？（终点定义、独立复现记录字段、smoke 检查项与通过条件） | required | S3/S4 实施与验收 | 首个 fork 文档/smoke 变更前 | 将 D-013 建议口径**按 A-001 R-002 最低清单**（工具/平台基线、依赖缓存前提、计时起止、失败/重试规则、证据字段、`scripts/smoke.sh` 机器可判定退出码）固化为可执行协议与脚本判据，记录决策 | **verified** | 已关闭（D-004）；v0.1.1（D-005 · A-008 响应）；**v0.1.2（D-007 · A-011 响应：退出码 8 + disposable 隔离守卫）** | [I-008-002-fork-reproduction-protocol.md](attachments/I-008-002-fork-reproduction-protocol.md) **v0.1.2**：计时、独立复现字段、SM-001～SM-006（S4 强制 ≥1 次 disposable 通过）、输入/退出码（含部分绿 8）与安全/隔离守卫；**不**勾选 S3/S4。S3/S4 实施后按协议对照：QUICKSTART + [R5-S3-REPRO-003](attachments/R5-S3-REPRO-003.md)（clean-ref，**预 T0 禁用 BuildKit 结果缓存，编译层非 CACHED**，64.833s）+ `scripts/smoke.sh` + 本地四路径验证（`r5-smoke-disposable-local-v0.1.2.txt`）+ CI `container-smoke` 成功 run |
| `I-008-003` | 最小操作日志（若实施）的事件与存储契约是什么？（事件类型、`operation_log` 表结构/迁移、查询/清理边界） | required（仅当 S6 实施） | S6 实施与验收 | 首个 operation_log 变更前 | 若用户决定在 R5 实施 S6，先冻结事件/存储契约；若不实施则记录为后续加分项 | open | 以用户对 S6 的决定为准 | 待用户决定 S6 后处理 |

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)（D-012 / D-013；Root `I-005` = verified、`I-006` = closed） |
| 前置证据 | [I-005-engineering-fork-collection.md](../GOAL-001-production-admin-foundation/attachments/I-005-engineering-fork-collection.md)（当前工程面事实、候选比较与裁决口径） |
| In | env/配置清单、dev/prod 区分、health/启动验证、容器一键启动（Compose）、fork 文档、15 分钟计时复现、`scripts/smoke.sh`、阶段审计与 Root 关门条件评估；可选操作日志（S6） |
| Out | 完整生产运维 / CI-CD 部署流水线；复杂 IAM / 业务模块；扩大 `I-PROTO-001 v0.1.3`；重开已关闭的 R1～R4 子目标；以「能启动」代替生产级验收 |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)
