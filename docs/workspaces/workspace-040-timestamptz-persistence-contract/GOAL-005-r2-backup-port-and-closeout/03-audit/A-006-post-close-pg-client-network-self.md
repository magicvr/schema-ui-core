---
id: A-006-post-close-pg-client-network-self
goal_id: GOAL-005-r2-backup-port-and-closeout
doc: audit-entry
source: self
auditor: current-session
date: 2026-09-21
scope: PR #16 post-close PostgreSQL CI failure, Docker helper network configuration, and hosted revalidation
verdict: pass
parent: null
version: 0.1.0
---

# A-006 · PR #16 post-close PostgreSQL helper network self 审计

## 范围与基线

- 工作区：`workspace-040-timestamptz-persistence-contract`；目标：`GOAL-005-r2-backup-port-and-closeout`（已 `done · 3/3`）。本意见只补记 R2 backup provider 在最终集成 CI 中暴露并修复的环境接线缺口，不重开或改写目标状态。
- 候选提交：代码修复 `20026ffd`；最终 PR 候选 `971ebbce`。
- 首次 Hosted CI：run `35567407333`。`api + postgres` 中 4 个 PG recovery 用例因 helper Docker 容器连接 runner `localhost:5432` 被拒绝而失败，其他 8 个 job 通过。
- 复验：run `35568956651`（修复提交 `20026ffd`）和最终候选 run `35569624987`（提交 `971ebbce`）；两轮各自 9 个 job 均 `success`。

## 成果（有证据）

| 检查项 | 结果 | 证据 |
|---|---|---|
| helper 容器网络能力 | pass | `apps/api/internal/backup/provider_pg.go` 仅在配置非空时将 `--network <name>` 放在 `docker run` 参数中；创建与恢复共用该参数构造 |
| 配置传递 | pass | `config.yaml`、`DB_CLIENT_DOCKER_NETWORK`、composition wiring 与恢复 wiring 测试已同步 |
| CI 环境接线 | pass | `.github/workflows/r6-basic-matrix.yml` 的 `api-postgres` 与 browser E2E PostgreSQL service job 设置 `host`；SQLite 情况无须调用 PG helper |
| 本地针对性回归 | pass | `go test ./internal/backup ./internal/config ./internal/composition` exit 0 |
| PostgreSQL 真实集成 | pass | run `35568956651` 的 `api + postgres` job 成功；最终 run `35569624987` 的 PG API 与两组 PG browser E2E 均成功 |
| 最终集成门禁 | pass | run `35569624987` 对最终 PR HEAD `971ebbce` 的 9/9 job 全部 `success` |

## Findings

| finding | level | 状态 | 证据 / 响应 |
|---|---|---|---|
| F-S-001 · PG helper 容器默认 bridge 无法访问 host-published PostgreSQL | required / high | **fixed** | `20026ffd` 加入可配置网络并为 CI service 选择 host 网络；真实 PG hosted job `35568956651` 通过，最终 run `35569624987` 9/9 通过 |

## 信息门禁与边界

- `I-041-006` 仍是用户裁决后的 `verified`；`I-041-004` 仍为 R3 前的 `non-blocking / deferred`，与本轮 helper 网络缺口无关。
- `DB_CLIENT_DOCKER_NETWORK` 默认为空，保持未显式配置网络时的既有命令行为；Compose 中可配置为可访问 DB service 的 network 名。
- 未以跳过 PG 测试替代修复；最终 CI 的 PostgreSQL 测试实际运行并通过。

## 结论

本 scope 判定 **pass**，`F-S-001` 已按 `fixed` 路径闭合，open required = 0。独立复核应核对 `--network` 位置、配置默认行为及最终 Hosted PG job 证据。

## 声明

本意见为 `source: self`，只记录代码修复和验证事实，不修改目标 `status` / `progress` 或 `goal-tree`；独立意见与任何后续响应仍由项目编排流程处理。
