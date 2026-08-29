---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-compose-cicd
version: 0.1.0
---

# A-002 · GOAL-004 关门独立审计（source: independent · 2026-08-29）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · GOAL-004（C1–C4 · E-002 证据 · D-001 落实度 · 残余登记）
- **verdict**：**pass**
- **工作区**：`workspace-024-distribution-formalization`（Root `GOAL-001-distribution-formalization` · `canonical_scope` 本区 · `shared_materials_catalog: none`）

## 范围与区间

核对 `00-meta` 成功标准 C1–C4 是否被 E-002 / 主仓 compose 制品 / golden-field `consumer-regression.yml` **可重复**支持；并核 D-001 落实度、I-024-002 有界核销、hosted 触发 R7 残余是否已登记且未主张 hosted acceptance。不改 `status` / `progress` / 方案正文 / goal-tree。

下游仓 `github.com/magicvr/golden-field` 仅作为 workspace.md 写明的实验消费实证对象读取，不引入其他工作区目标状态。

## 成果（有证据）

| 项 | 独立核对 | 证据路径 |
|----|----------|----------|
| C1 compose 全服务 | 审计员复跑：`docker compose up -d --build` **exit 0**（镜像 cache hit = 实施侧 17:42 构建 `schema-ui-core-api:local` `f533c8e563e0` / `schema-ui-core-web:local` `c3eab22c45be`）；api **healthy** + web **healthy**；`compose exec api wget` `/readyz` → `{"status":"ok",...}`；宿主 `http://127.0.0.1:25081/` **HTTP 200**（body 849 B）；`compose stop` → api **ExitCode 0** · web **ExitCode 0**；日志 `shutdown.starting`（`signal=terminated`）→ `shutdown.complete` → `bye` | 本审计复跑（2026-08-29T09:52:51Z–09:53:03Z UTC）；`compose.yaml` healthcheck=`/readyz` · `stop_grace_period: 15s` |
| C2 workflow 重构 | golden-field `c4d14eacb6e702bfadad6a196d018e0d036eb607`（2026-08-29 17:41:31 +0800）：删 `NODE_AUTH_TOKEN` 与 GH Packages `npm config set`；`setup-node` `cache: pnpm` + `cache-dependency-path: web/pnpm-lock.yaml`；`go get @latest` + `go build`；serve 后台 + healthz/readyz 轮询 + 四探针 + `kill -TERM` + `grep shutdown.complete` | `C:\Users\magicvr\Documents\Code\golden-field\.github\workflows\consumer-regression.yml`；`git show c4d14ea` |
| C2 本地等价 | 四探针全绿：protocol **2.9** / render **1573 B** / six-package PASS / token `brand=2 ⊆ index=5`；`go build ./cmd/server` + `serve -addr 127.0.0.1:25118` → `/healthz` 200 `{"status":"ok"}` · `/readyz` 200 `{"status":"ok"}`。Windows 未投递 SIGTERM（与 D-001 一致） | golden-field `web/probe.mjs` 等；本审计复跑 |
| C3 harness A | = C1 compose stop 路径：exit 0 + `shutdown.starting`/`shutdown.complete` | 同上 C1 |
| C3 harness B | 审计员复跑：`docker run` api 镜像 + `HTTP_SHUTDOWN_TIMEOUT=1s` + 宿主 TcpClient 不完整 `POST /api/auth/login` 保持在途 → `docker stop` → 日志 `shutdown.starting`（terminated）→ `shutdown.timeout`（`context deadline exceeded`，约 1s）→ **ExitCode 1** | 本审计复跑（容器 `b4eb330f07da` · 2026-08-29T09:53:48Z–09:53:51Z UTC） |
| C4 / I-024-002 | 本地等价 + linux 容器实跑成立；**未**触发 GitHub-hosted runner；E-002 残余 1 / A-001 R-001 / D-001 条 5 / C4 均登记 R7，未主张 hosted acceptance | D-001；E-002；A-001 |
| 凭据卫生 | 仓库根 `.env` 仅有 `github_token`/`npm_token` 键（无 AUTH 值落入 compose 文件）；无 env 时 `docker compose ps` fail-closed（`ADMIN_INITIAL_PASSWORD is required`）。与 D-001「命令环境注入、不写 .env」一致 | `.env` 键名扫描（不记录值）；compose 插值错误 |
| 共享资料 | `none`，无引用被当成证据 | `workspace.md` |

## 对照成功标准

| 标准 | 状态 | 独立证据 | 缺口 |
|------|------|----------|------|
| C1 compose 全服务 healthy + readyz ok + web 200 + stop exit 0 + drain 日志 | **达成** | 本审计复跑与 E-002 表一致（`signal=terminated` 与 E-002「terminated」吻合） | E-002 无原始 transcript / container id（F-001） |
| C2 workflow 免凭据重构 + 本地等价四探针全绿 | **达成（有界）** | 文件 diff + 四探针 + serve healthz/readyz；Windows 无 SIGTERM 由 C3 容器 + R7 覆盖（D-001） | `setup-node cache: pnpm` 位于 `corepack enable` **之前**，hosted 触发会因找不到 pnpm 失败（F-002）；本环境未能核对 `c4d14ea` 已推 origin |
| C3 drain harness A/B（linux 容器） | **达成** | A = C1 stop；B = 1s 预算 + 在途慢请求 → timeout + exit 1 | 对象是 compose `cmd/server` 镜像，不是 `server.Serve` 下游壳（F-003） |
| C4 I-024-002 核销；hosted 触发登记 R7、不主张 hosted acceptance | **达成（有界）** | 等价实跑成立；R7 登记在 D-001 / E-002 / C4 / A-001 R-001。Root `00-meta` I-024-002 仍为 `open`（待 `/govern` 关门时同步，不构成本条过声明） | 不阻本 scope |

## Findings

### F-001 · 执行索引与派生 progress 未随 E-002 对齐

- 严重度：low
- 建议：**recommended**
- 状态：open
- 描述：E-002 已落盘且 `00-meta` C1–C4 已勾选，但 `02-execution.md` 索引仍只有 E-001；`progress` 仍为 `0/4`。不否定 C1–C4 事实，属台账卫生。独立审不改 status/progress。
- 证据：`02-execution.md` 索引表；`00-meta` `progress: 0/4` vs C1–C4 `[x]`。
- 关闭条件：`/govern` 响应时挂上 E-002 索引行，并按 P-001 由检查点重算派生 progress（不得用手填百分比兜底）。

### F-002 · `consumer-regression.yml` 的 pnpm cache 步骤顺序会使 hosted job 在 setup-node 失败

- 严重度：med
- 建议：**recommended**
- 状态：open
- 描述：D-001 / C2 要求 `setup-node` 挂 pnpm cache。现行顺序是 `actions/setup-node@v4`（`cache: pnpm`）→ 随后 `corepack enable`。`setup-node` 在解析 cache 路径时需要 PATH 上已有 `pnpm`，否则典型失败为 `Unable to locate executable file: pnpm`（actions/setup-node#530）。`web/package.json` 亦无 `packageManager` 字段。C4 明确不主张 hosted acceptance，故不升格 required；但 R7 实触发前必须先改顺序（`pnpm/action-setup` 或 `corepack enable` 在 `setup-node` 之前），否则登记的 hosted 复核会在第一项 cache 步失败，无法验证 SIGTERM 断言。
- 证据：golden-field `consumer-regression.yml` L24–31 vs L30–31；`web/package.json` 无 `packageManager`。
- 关闭条件：调整 workflow 使 pnpm 在 `cache: pnpm` 之前可用；R7 触发时该步不再因缺少 pnpm 失败。

### F-003 · R1「信号级 drain harness」核销口径应写清对象（cmd/server 容器，而非 serve 壳进程）

- 严重度：low
- 建议：**recommended**
- 状态：open
- 描述：D-001 把 C3 定为 compose api 镜像（`apps/api/Dockerfile` ENTRYPOINT = `./cmd/server`）的 linux 容器 A/B，本审计按该口径复现通过。GOAL-002 E-003 残余 1 原文是 serve 面 SIGTERM → exit 0/1 在 linux runner 补齐；`server.Serve` 的 SIGTERM 仅写入 golden-field workflow（未 hosted 实跑），Windows 本地等价按 D-001 不伪造 SIGTERM。E-002「R1 残余 1 → 核销」过宽，易被读成 serve 壳进程级 harness 已在 linux runner 执行。C3 本身不因此失败。
- 证据：`apps/api/Dockerfile` L22 / L37；`golden-field/cmd/server/main.go` 调 `server.Serve`；workflow L61–63 未在本环境执行；GOAL-002 `02-execution/E-003-s2-s3-evidence.md` 残余 1。
- 关闭条件：编排响应把核销句改为「compose `cmd/server` 容器 A/B 已证；serve 面 SIGTERM = workflow 文件交付，实跑随 R7 hosted」。

## 必改项汇总

无 required。

## 与既有意见的异同

| 项 | A-001 self | 本条 independent |
|----|------------|------------------|
| C1–C4 主体 | 达成（摘要表） | **同意**；并独立复跑 C1/C3B/四探针/serve 冒烟 |
| hosted R7 | R-001 recommended 登记 | **同意**，不新开 finding；不主张 hosted acceptance |
| I-024-002 | verified（有界） | **同意**有界核销；Root 信息表仍 `open` 交 `/govern` |
| 本条新增 | — | F-001 台账索引；F-002 workflow cache 顺序；F-003 R1 核销口径 |
| verdict | conditional（待独立审） | **pass**（0 required） |

## 结论 + 建议给编排器/用户的下一步

C1–C4 在 D-001 有界口径下成立：compose 全服务与 drain A/B 可独立复现；workflow 免凭据重构与本地四探针/serve 冒烟成立；hosted runner 实触发已登记 R7、未被写成 acceptance。无未关闭 high/med **required** finding，无到期且影响本 scope 的未核销 required 信息项（I-024-002 按用户环境等价 + 容器实跑核销，hosted 为书面残余）。

建议用户执行：`/govern` 响应 GOAL-004 A-002（及 A-001 R-001）。可直接关门 GOAL-004，并把 F-001～F-003 与 hosted 触发一并记入 Root/R7 复核清单；**不要**把本意见写成 GitHub-hosted acceptance。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。

---

## 响应（2026-08-29 · /govern · source: self）

- **F-001 → fixed**：02-execution 索引挂 E-002；meta progress 4/4（P-001 检查点重算）· status done。
- **F-002 → fixed**：`consumer-regression.yml` 改 `pnpm/action-setup@v4` 先于 `setup-node`（cache: pnpm 需要 PATH 上的 pnpm）；`web/package.json` 声明 `packageManager: pnpm@11.11.0`（golden-field commit `3f2a5c2`）。
- **F-003 → fixed**：E-002 核销口径精化为「compose `cmd/server` 容器 A/B 已证；serve 面 SIGTERM = workflow 文件交付，实跑随 R7 hosted 触发核销」。
- R-001（hosted 触发）保持登记（R7）。全部 required 闭合 → GOAL-004 done（用户框架授权 · 无必改项）。
