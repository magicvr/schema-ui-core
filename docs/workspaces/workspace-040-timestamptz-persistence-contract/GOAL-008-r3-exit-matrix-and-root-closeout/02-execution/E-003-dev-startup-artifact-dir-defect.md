---
id: E-003-dev-startup-artifact-dir-defect
doc: execution-entry
status: active
parent: GOAL-008-r3-exit-matrix-and-root-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-003 · 用户报告的 dev 启动失败：`recoveryArtifactsDir` 相对路径导致 docker 挂载失败（已修复）

## 来源

用户 2026-09-21 在 Root 关门确认包中报告：**`.\dev.cmd start` 目标启动失败，需要先修复这个问题**（用户同时裁定 `I-041-011` 只保留在退出矩阵/R3-C 附件、并授权为 GOAL-004 的 `user-overruled` 补一条 `D-` 条目）。

## 复现（编排器独立复现，非采信叙述）

以启动器同样的环境运行 API（`APP_ENV=development`、`CONFIG_FILE=apps/api/configs/config.yaml`、`HTTP_ADDR=:25080`）：

```text
go run ./cmd/server
→ WARN AUTH_JWT_SECRET not set; generated ephemeral random key
→ [Fx] ERROR Failed to start: ... LIFECYCLE_START_FAILED [core.auth-session]:
   open store: postgres apply catalog: backup: pg create (pg_dump): ToolFailure:
   docker run --rm -v data\recovery:/vp040-backup postgres:15-alpine pg_dump ... :
   exit status 125: docker: Error response from daemon: create data\recovery:
   "data\\recovery" includes invalid characters for a local volume name, only
   "[a-zA-Z0-9][a-zA-Z0-9_.-]" are allowed. If you intended to pass a host directory,
   use absolute path
```

### 根因

1. `configs/.env`（gitignored，但**启动器与 API 都会读**：`internal/config.loadEnvFile` 自动加载 `configs/.env`，进程环境优先）设置了 `DB_DIALECT=postgres`、`DB_NAME=postgres`、`DB_HOST=192.168.31.213` …，因此即便 `configs/config.yaml` 写 `dialect: sqlite`，**生效方言仍是 postgres**。
2. `configs/config.yaml` 的 `db.path: ./data/schema-ui.db` 是**相对**路径。
3. `composition.recoveryArtifactsDir` 返回 `filepath.Join(filepath.Dir(cfg.DBPath), "recovery")` —— 对相对 `db.path` 即 `data\recovery`（**相对**）。
4. `composition.recoveryWiring` 把该值直接交给 `backup.PgProvider.WorkDir`，provider 用它构造 `docker run -v <WorkDir>:/vp040-backup`；**docker 拒绝相对源路径**，`pg_dump` 退出 125。
5. 该 dump 是 C3 的 **class-A rollback 产物**，创建失败即 `openStore` 失败 → 整个 API 启动失败。

### 为什么既有测试没抓到

`internal/composition` 的 PG 启动测试使用 `t.TempDir()`（**绝对**路径），因此 `WorkDir` 一直是绝对路径；只有「配置里 `db.path` 为相对路径」的真实 dev 入口才会触发。这是一条**测试夹具与真实配置之间的覆盖缺口**（不是断言写错）。

## 修复（本轮）

| 文件 | 改动 |
|------|------|
| `internal/composition/composition.go` | `recoveryArtifactsDir` 的所有返回分支经新增 `absoluteArtifactDir` 解析为**绝对路径**（`filepath.Abs`，失败则保留原值）；同时更正函数文档注释：`db.path` **优先**（含 postgres；对 postgres 而言它是文件存储根），仅当 `db.path` 为空时才落到用户缓存目录按 DSN 派生 —— 与 `GOAL-005/A-004` 记载的实际行为一致 |
| `internal/composition/recovery_wiring_test.go`（新增） | `TestRecoveryArtifactsDirIsAbsolute`（5 例：相对 `./data`、带 `./` 前缀、绝对、仅 DSN、全空）断言**必为绝对路径**；`TestRecoveryWiringHandsPgProviderAnAbsoluteWorkDir` 断言 `PgProvider.WorkDir` 与 wired dir 绝对、镜像默认值正确、空配置禁用锚点、sqlite 无 rollback creator |

## 端到端验证（编排器实测）

在**一次性**测试库上复跑真实启动路径（不触碰共享 `postgres` 库）：

```text
psql: CREATE DATABASE vp040_devsmoke
DB_NAME=vp040_devsmoke APP_ENV=development CONFIG_FILE=...config.yaml HTTP_ADDR=:25080 \
  go run ./cmd/server

→ healthz 200 within ~5s : {"status":"ok","timestamp":"2026-09-21T00:29:26.557243Z",...}
→ readyz  200 within ~10s: {"status":"ok","timestamp":"2026-09-21T00:29:31.874018Z",...}
→ data/recovery/rp-1789950559694-587a18ec47e0d8e4.artifact (101,787 bytes)  ← class-A/B 产物落盘成功
```

清理：停止进程（`:25080` 无监听）→ `DROP DATABASE vp040_devsmoke WITH (FORCE)` → 删除 `data/recovery/*`（该目录在 `apps/api/.gitignore` 内，本就不入库）。工作树仅剩本次修复的两个文件。

**附带确认**：`/healthz` 与 `/readyz` 的 `timestamp` 在**真实运行**中已是 R3-A 的规范 fixed-6 形状（`…T00:29:26.557243Z`），这是 R3-A 收敛在活体进程上的旁证。

## 未决（需用户裁决，已随确认包一并提出）

- 该 `.env` 使 **dev 启动器指向共享常驻实例的 `postgres` 库**（维护库），而不是专用 dev 库（`config.yaml` 的 `name: schema_ui` 被 `.env` 覆盖）。本轮验证刻意改用一次性库避免污染；**用户日常 `dev.cmd start` 会在 `.env` 指向的库上跑完整迁移链**。是否给 dev 单独指定库名属用户环境决策。

## 进度评估

修复 + 回归测试 + 端到端验证完成；`GOAL-008` 仍 `1/3`（检查点 B 已完成，检查点 C 待用户确认）。该缺陷已被记为 `A-004` 中的新 finding 并以 `fixed` 闭合。
