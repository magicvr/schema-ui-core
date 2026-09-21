---
id: E-005-post-close-pg-client-network-fix
parent_goal: GOAL-005-r2-backup-port-and-closeout
doc: execution-entry
status: recorded
date: 2026-09-21
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-005 · PR #16 PostgreSQL CI 暴露 helper network 缺口并修复

## 已发生事实

1. PR #16 首次 `r6-basic-matrix` run `35567407333` 中，`api + postgres` 失败：`pg_dump` / `pg_restore` 在临时 Docker 客户端容器内访问 `localhost:5432` 被拒绝，导致 4 个 PG recovery 用例失败；同一 run 的其他 8 个 job 通过。
2. 提交 `20026ffd` 为 `PgProvider` 增加可选 `ClientDockerNetwork`，并由 `db.client_docker_network` / `DB_CLIENT_DOCKER_NETWORK` 配置进入 composition。Docker 工具容器仅在显式设置时追加 `--network`；GitHub `api-postgres` 与 browser E2E PostgreSQL 服务 job 显式设置 `host`，默认应用行为保持不变。
3. 本地 `go test ./internal/backup ./internal/config ./internal/composition` 通过。Hosted run `35568956651` 在 `20026ffd` 上的 9 个 job 全部通过，包含真实 `api + postgres` 测试。
4. 发布钉与包名映射提交 `971ebbce` 后，最终候选 CI run `35569624987` 的 9 个 job 全部通过；其中 `api + postgres` 和两组 PostgreSQL browser E2E 均为 `success`。
5. 本条是完成目标后的证据补录，不修改本目标的 `status`、`progress` 或 Root 路线图状态。

## 验证边界

- CI 覆盖 GitHub Linux runner 上发布端口 PostgreSQL service 由 helper 容器经 host 网络访问的情况。
- `DB_CLIENT_DOCKER_NETWORK` 为空时不向 Docker 命令追加网络参数；Compose 部署可单独配置数据库所在的用户网络名。
- 首次失败不是通过跳过测试处置；最终 PR 候选的全部 9 个 job 均在成功状态。
