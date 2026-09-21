---
id: A-005-independent-fi101-closure
doc: audit-entry
status: active
parent: GOAL-008-r3-exit-matrix-and-root-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---
# A-005 · independent 定向复审：`F-I-101` 闭合

- **source**: independent
- **provider / auditor**: grok build · model grok-4.6 · reasoning high
- **日期**: 2026-09-21
- **scope**: `GOAL-008-r3-exit-matrix-and-root-closeout` · finding-closure · `F-I-101`（required，用户报告 `.\dev.cmd start` 启动失败 → 声称 `fixed`）；附带核验同类相对路径隐患、退出矩阵 `L-1` 限定、`I-041-011` / GOAL-004 `D-002` 用户裁决留痕、Root 关门就绪度
- **audit_type**: finding-closure
- **verdict**: **pass**
- **开放 required（本审新开 + 既有未闭合 findings）**: **0**
- **是否同意 `F-I-101` 已 `fixed`**: **同意**
- **是否同意现在再次把 Root 交给用户确认关门**: **同意**（判据 6 / `I-041-010` 仍待用户书面确认；本意见不代替确认、不同意把 Root 标 `done`）
- **新的 P-004 待裁点**: **无必须新裁项**。`.env` 指向共享维护库 `postgres` 已由 `A-004`/`E-003` 提出，属用户环境决策，本审不升格为 finding。

本意见不修改 `status` / `progress` / 方案正文 / goal-tree。按本轮用户书面要求，**不改仓库内任何文件**；请编排器原样落盘为本条目并更新 `03-audit.md` 索引。响应归 `/govern`。

基线：`git log --oneline -3` = `0f9b16bc`（本轮修复）在 `1ca0d90a`（A-002 落盘后提交）之上；`git diff 1ca0d90a..HEAD --stat` 共 11 文件、+327/−16。结束时 `git status` = **working tree clean**。

---

## 1. `F-I-101` 逐项闭合判定

| 项 | 判定 | 独立证据（非采信叙述） |
|----|------|------------------------|
| **根因是否成立** | **成立** | 旧 `recoveryArtifactsDir` 在 `db.path` 非空时返回 `filepath.Join(filepath.Dir(cfg.DBPath), "recovery")`，对 `./data/schema-ui.db` 即为相对 `data\recovery`。`recoveryWiring` 把它赋给 `backup.PgProvider.WorkDir`。`PgProvider.Create`/`Restore` 原样拼接 `"-v", p.WorkDir+":"+/vp040-backup`。本审在仓库外复刻旧函数后执行 `docker run --rm -v data\recovery:/vp040-backup postgres:15-alpine …` → **exit 125**，报错原文与 `E-003` 一致：`"data\\recovery" includes invalid characters for a local volume name … use absolute path`。`configs/config.yaml` 确为 `path: ./data/schema-ui.db`；`loadEnvFile` 默认读 CWD 相对的 `configs/.env`（进程环境优先）；`dev.cmd` 以 `-WorkingDirectory %API_DIR%` 启动 API。`TestCompositionPostgresStartup` 使用 `t.TempDir()`（绝对路径），覆盖缺口属实。 |
| **修复是否完整** | **对所报缺陷：完整** | 两条返回分支都经 `absoluteArtifactDir` → `filepath.Abs`。本审探针：相对 `data/…`、字面 `./data/schema-ui.db`、字面 `.\data\schema-ui.db`、`filepath.Join(".", "data", …)`、仅文件名、已是绝对路径、Windows `C:\data\…`、UNC `\\server\share\…`、仅 DSN（走 `UserCacheDir`）、空配置、`os.TempDir()` 回退，**新函数在非空结果上均为 `filepath.IsAbs == true`**；空配置仍返回 `""`（锚点禁用，docker 不会被调用）。`recoveryWiring` 把该绝对目录同时交给 `NewService` 与 `PgProvider.WorkDir`；`Create` 的 `-v` 源即该字段。生产构造点只有 `composition.go` 这一处。 |
| **测试是否决定性** | **对所报缺陷：是** | 本审 `go test -count=1 -run "TestRecoveryArtifactsDirIsAbsolute\|TestRecoveryWiringHandsPgProviderAnAbsoluteWorkDir" ./internal/composition/` → **PASS（0.201s）**。把 `absoluteArtifactDir` 退回相对后的**等价旧函数**：`./data/schema-ui.db` → `data\recovery` 且 `IsAbs==false`，现有 `filepath.IsAbs` 断言会失败。测试未直接调 docker，但 `Create` 对 `WorkDir` 无再加工；`IsAbs(WorkDir)` 与 `-v` 源是否绝对**等价**。字面 `./data` 与 `Join(".", "data", …)` 经 `filepath` Clean 后同为 `data\recovery`，覆盖等价。 |
| **端到端是否可复现** | **本审已复跑，可复现** | 一次性库 `vp040_a005`（**未**触碰共享 `postgres` 库）：`healthz 200`、`readyz 200`（`timestamp` 为 fixed-6：`2026-09-21T00:43:59.097189Z` / `…59.792537Z`）；`openStore` 约 8.98s 成功（原失败点）；产物 `data/recovery/rp-1789951432020-5d6e89d3f49edb4a.artifact` **101,764 B**（class-B 走同一 `PgProvider.Create` / `docker run -v`）。随后 `DROP DATABASE vp040_a005 WITH (FORCE)`、删除该产物、`:25081` 无监听。 |
| **是否触碰冻结物** | **否** | `1ca0d90a..HEAD` 代码仅 `apps/api/internal/composition/composition.go` + 新测试。未改 v1–v87 canonical SQL/checksum、codec、wire 合同、`kernel` C3 Port 形状、`PgProvider` 字段/CLI 形状（仍 `pg_dump -F c` / `pg_restore --exit-on-error --no-owner`）。不构成新的 required。 |

### `filepath.Abs` 与 CWD / `dev.cmd`

`filepath.Abs` 相对**进程 CWD** 解析。`db.path` 本身也是 CWD 相对（SQLite `Open`、`.env` 的 `configs/.env`、文件存储根 `filepath.Dir(db.path)` 皆然）。本审核对：`Abs(Dir(db.path)/recovery)` = `Dir(Abs(db.path))/recovery`。

`dev.cmd` 即使用户在**仓库根**执行 `.\dev.cmd start`，API 仍是 `Start-Process -WorkingDirectory %API_DIR%`（`apps\api`）。因此启动器路径下 Abs 与 `db.path` 解析**一致**，结论仍成立。若有人在仓库根直接 `go run ./cmd/server` 且不改 CWD，相对路径会落到仓库根而非 `apps/api`，但 Abs 与 SQLite 仍彼此一致；那不是本缺陷，也不由本次修复引入。

### 残余（不推翻 `fixed`）

`absoluteArtifactDir` 在 `filepath.Abs` 失败时返回原值（可为相对）。本审用含 NUL 的非法路径使 Abs 返回 `invalid argument`，回退相对。真实 `db.path` 还要过 `validateDBPathShape`；Abs 失败的现实场景主要是 `Getwd` 失败（进程已不健康）。`PgProvider.Create`/`Restore` **自身不** `Abs`、也不拒绝非绝对 `WorkDir`。生产唯一构造点已 Abs；这是防御纵深缺口，见 `F-I-102`（recommended），**不**把 `F-I-101` 打回未闭合。

---

## 2. 同类隐患扫描（路径交给外部进程或容器）

| # | 位置 | 用法 | 相对路径是否安全 | 依据 |
|--:|------|------|------------------|------|
| 1 | `internal/backup/provider_pg.go` `Create`/`Restore` | 生产 `docker run -v <WorkDir>:/vp040-backup` | **修复后生产路径安全** | WorkDir 来自已 Abs 的 `recoveryArtifactsDir`。Provider 自身不校验 `IsAbs`（`F-I-102`）。 |
| 2 | `internal/backup/pg_cross_version_matrix_test.go` `matrixDump`/`matrixRestore` | 测试 `docker run -v workDir:/vp040` | **安全** | `workDir := t.TempDir()`，绝对。 |
| 3 | `internal/store/c3_pg_anchors_test.go`、`provider_pg_test.go` | 测试构造 `PgProvider{WorkDir: t.TempDir()}` | **安全** | 夹具绝对；这正是 `F-I-101` 原先漏掉真实相对 `db.path` 的原因，矩阵 `L-1` 已写明。 |
| 4 | `internal/composition/post_rotation_recovery_test.go` | `docker exec` + `pg_dump` stdout，无 host bind | **安全** | 不把宿主目录挂进容器。 |
| 5 | `compose.yaml` | `./apps/api/configs:/app/configs:ro`；命名卷 `db-data` | **安全** | Compose 相对路径相对 compose 文件目录解析，不是 `docker run -v` 的 volume-name 语法。 |
| 6 | `store/migrate.go` `VACUUM INTO`、`store/recovery.go` `VACUUM INTO` | SQLite 进程内 SQL，目标 = `s.path` + 后缀 | **安全** | 不经 docker；相对路径由 SQLite/CWD 解析。 |
| 7 | `backup/provider_sqlite.go` `VACUUM INTO` / `os.Open` | 进程内文件 | **安全** | 同上；生产 `artifactPath` 现为绝对。 |
| 8 | `provider_pg.go` `--file /vp040-backup/<basename>`；矩阵 `-f /vp040/<archive>` | 容器内绝对路径 | **安全** | 宿主侧危险点只在 `-v` 源。 |
| 9 | `CREATE DATABASE` | `admin.ExecContext` / `quoteIdent`，不是 shell `createdb` | **安全** | 无路径参数。 |
| 10 | `config.loadEnvFile` `os.ReadFile("configs/.env")` | CWD 相对 | **对启动器安全** | `WorkingDirectory=apps\api`。从仓库根裸 `go run` 可能读不到该文件——既有行为，不是 docker 挂载缺陷。 |
| 11 | 对象存储/上传根 `filepath.Dir(cfg.DBPath)` | `os` 文件 API | **安全** | 不交给 docker。 |
| 12 | `cmd/schema-ui/migrate.go` `exec.Command` | `go`/`npm`，工作目录显式 `Dir` | **安全** | 无 `docker -v`。 |
| 13 | `scripts/*.sh` `docker compose … down -v` | 删命名卷，不是 bind 源 | **安全** | `-v` 在此不是宿主导出路径。 |

未发现第二处「相对宿主路径 → `docker run -v`」的生产调用。`F-I-101` 不是家族性漏网的冰山一角，是 C3 PG provider 与相对 `db.path` 的交点。

---

## 3. Findings

### F-I-101 · 相对 `db.path` → `docker run -v` 启动失败
- **级别**: required（用户要求先修；PG 真实入口无法启动）
- **闭合**: **`fixed`**（本审同意）
- **证据**: 见 §1 与 §5。代码 + 回归测试 + 本审 docker 反例（相对 exit 125 / 绝对成功）+ 一次性库 e2e（healthz/readyz 200 + 产物落盘）。

### F-I-102 · `PgProvider` 不强制绝对 `WorkDir`；`Abs` 失败回退相对
- **级别**: recommended
- **问题**: (1) `Create`/`Restore` 把 `WorkDir` 直接拼进 `-v`，调用方若传入相对路径仍会 exit 125；(2) `absoluteArtifactDir` 在 Abs 失败时原样返回。注释写「ALWAYS absolute」对 Abs 失败分支不成立。
- **证据**: `provider_pg.go:136,186`；`composition.go:314-318`；本审 NUL 路径 Abs → `invalid argument`、回退相对。
- **关闭要求**: 任选或组合：`Create`/`Restore` 在非绝对 `WorkDir` 上 fail closed；Abs 失败时不要把相对路径交给 docker（返回错误或禁用锚点）。不要求重开 `F-I-101`。
- **为何不是 required**: 生产唯一构造点已 Abs；Abs 失败不是所报缺陷的触发条件。

### F-I-103 · `00-meta` 信息表仍把 `I-041-011` 标为 `open`
- **级别**: recommended
- **问题**: `01-decision.md` 与 `A-004` 已记 `I-041-011` = **verified（用户 2026-09-21 裁决：只留退出矩阵/R3-C 附件）**；检查点 C 叙述也写 verified；**同文件** `00-meta.md` 信息表仍为 `open`。
- **证据**: `GOAL-008/00-meta.md` 信息表 vs 检查点 C 行；`01-decision.md` L18；`A-004` §1。
- **关闭要求**: 把 `00-meta` 信息表该行改为 `verified`，与决策索引一致。不阻断 Root 确认（该项 non-blocking，且裁决已另有载体）。
- **为何不是 required**: 不是未裁；是投影未同步。

**必改项汇总**: 无。开放 required = **0**。

---

## 4. 治理核验（本轮要求的第 6–8 点）

- **`F-I-101` 的 `fixed` 有可核对修正**: 代码（`absoluteArtifactDir`）+ 测试（两函数）+ 端到端（编排器 `E-003` 与本审 `vp040_a005` 各一次）。三路径用的是 `fixed`，不是口头。
- **冻结物**: 未触碰。见 §1。
- **Root 关门就绪度**: `F-I-101` 记为 required 并在修复后**不再阻塞**——同意。退出矩阵判据 3 的 `L-1` 已追加「测试夹具绝对路径不能代表配置里相对 `db.path`」；本审同意该限定现在写明了覆盖缺口。跨目标开放 required 不因本轮而增加。`I-041-010` 仍是**用户确认门禁**（信息项，不是 finding）。可以再次把 Root 交给用户确认关门；**不可以**自行标 `done`。
- **`I-041-011`**: 与确认包选项及 `A-004` 记录一致（只留矩阵/R3-C，不进发布说明）。未改结论。`00-meta` 表未跟上 = `F-I-103`。
- **GOAL-004 `D-002`**: 与 `GOAL-004/A-003` 2026-09-20 同期记录一致（接受 `*time.Time` + JSON `null` 为 R2 必要后果）；2026-09-21 只补书面载体。本审未见台账以外的逐字原话稿，但 `D-002` **没有引入新结论、没有改代码**。与 `A-002` 的 `F-I-004` 关闭要求匹配。

---

## 5. 反例 / 复跑记录

### 5.1 仓库内回归
```
cd apps/api
go test -count=1 -timeout 60s -run "TestRecoveryArtifactsDirIsAbsolute|TestRecoveryWiringHandsPgProviderAnAbsoluteWorkDir" ./internal/composition/
ok  github.com/magicvr/schema-ui-core/apps/api/internal/composition  0.201s
```

### 5.2 仓库外旧/新函数 + 真 docker（未改仓库）
探针目录 `$TEMP\vp040-a005-probe`（用后已删）。旧函数对 `./data/schema-ui.db` → `data\recovery`（相对）；新函数 → 绝对路径。`Abs(recovery)` 与 `Dir(Abs(db.path))/recovery` 相等。

```
DOCKER relative-defect src="data\recovery"  err=exit status 125
  docker: Error response from daemon: create data\recovery:
  "data\\recovery" includes invalid characters for a local volume name ...
  If you intended to pass a host directory, use absolute path

DOCKER absolute-fix  src="<abs>\data\recovery"     err=<nil>
DOCKER abs-tempdir   src="<abs>\probe-recovery"    err=<nil>
CHECKS failed=0 total=18
```

Windows 盘符绝对路径可作为 `docker run -v` 源（Docker Desktop 29.7.2 接受 `C:\…:/vp040-backup`）。

### 5.3 端到端（一次性库 `vp040_a005`）
- 凭证从 `apps/api/configs/.env` 读取，未打印秘密；`DB_NAME` 进程环境覆盖为 `vp040_a005`。
- `CREATE DATABASE vp040_a005` → `SELECT datname` 核验存在 → 启动 `go run ./cmd/server`（CWD=`apps/api`，`APP_ENV=development`，`HTTP_ADDR=:25081`）。
- `GET /healthz` **200**；`GET /readyz` **200**；class-B 产物 101,764 B。
- 停止进程；`:25081` 无监听；`DROP DATABASE vp040_a005 WITH (FORCE)`；删除新产物；recovery 目录剩余 0。

### 5.4 清理
```
git status → nothing to commit, working tree clean
HEAD     → 0f9b16bc
parent   → 1ca0d90a
```
仓库外探针目录已删。未改仓库文件。

---

## 6. 明确结论

| 问题 | 回答 |
|------|------|
| 开放 required 数量 | **0**（本审新开 0；`F-I-101` 同意 `fixed`；`F-I-102`/`F-I-103` 为 recommended） |
| 是否同意 `F-I-101` 已 `fixed` | **同意** |
| 是否同意现在再次把 Root 交给用户确认关门 | **同意** |
| 是否同意把 Root 标 `done` | **不同意**（`I-041-010` 仍 open） |
| 新的 P-004 待裁点 | **无必须新裁项**。可选环境问题（dev 是否专用库名）已在 `A-004` §4 提出，不阻塞确认包。 |

与上一轮 `A-002`（`conditional` / 开放 required = 0 / 同意交给用户确认）相比：所报启动缺陷是确认之后新出现的 required，现已独立复现并闭合；Root 关门就绪度在 findings 意义上恢复为「可再提请确认」。判据 6 仍未满足，这是设计如此。

---

## 7. 建议给编排器 / 用户的下一步

1. 把本意见落盘为 `GOAL-008/03-audit/A-005-independent-fi101-closure.md`，索引表追加 A-005（source=independent，verdict=pass）。
2. 用 **`/govern`** 响应：将 `F-I-101` 维持 `fixed`；`F-I-102`/`F-I-103` 可在检查点 C 卫生中处理或登记，**不**因此否决 `I-041-010`。
3. 再次把 Root 确认包交给用户（`I-041-010`）。不要静默标 `done`。

### 声明
本意见不修改 status/progress；响应由 /govern 处理。
```
