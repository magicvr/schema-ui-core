---
title: 执行记录 · R5 · 工程化、fork 体验与集成关门
status: done
created: 2026-08-02
updated: 2026-08-20
parent: GOAL-001-production-admin-foundation
version: 0.1.17
---

# 执行记录 · GOAL-008

## 2026-08-02 · 立项（承接 Root D-013）

- 用户通过 `/govern` 确认 R5 方案边界（Root D-013：部署基线 A、建议 15 分钟口径、复现方法、I-006 方案甲）后，立项本 R5 子目标；五件套与 `attachments/` 齐全。
- 成功标准 S1～S5（环境/配置基线、容器一键启动、fork 文档与 15 分钟体验、可复现 smoke 验收、阶段审计与 Root 关门条件评估）为五个核心检查点，`progress: 0/5`；S6（最小操作日志）为可选加分项，不进进度分母。
- 登记三项实施前 required 信息门禁：`I-008-001`（环境/配置/容器契约）、`I-008-002`（计时复现协议与 smoke 判据）、`I-008-003`（operation_log 契约，仅当 S6 实施）；当前均 `open`。
- **未做**：未修改产品代码、配置、文档、容器或脚本；未运行应用/测试；未勾选任何检查点。Root R5 检查点未勾选，Root 保持 `active / 4/5`。
- **计划（非事实）**：先在 `GOAL-008` 收集并冻结 `I-008-001`，再判断 S1/S2 实施边界。

## 2026-08-02 · 响应 A-001（F-001 fixed · R-001/R-002 handled）

- **A-001（independent · conditional）响应**：采纳 `conditional`。F-001 → **fixed**——新增 `01-decision` **D-002**：Docker Compose 为 R5 **必须交付和验收的第二启动路径**（S2 核心检查点、计入进度分母），非 S6 式可选加分项；D-001 边界修订删除「可选加分路径」表述；`00-meta` S2 同步。R-001 → **handled**：I-005 附件更新至 v0.2.1（§2～§6 时态清理为「D-013 前历史候选 · 已裁决」，frontmatter `related_decisions: D-012, D-013`）。R-002 → **handled**：`00-meta` 信息表 `I-008-001/002` 补入 R-002 最低收集清单（env/secrets、DB volume、SPA fallback 与 `/api` 反代、readiness、依赖/超时/失败行为、CI 入口；工具基线、依赖缓存前提、计时起止、失败/重试规则、证据字段、smoke 退出码）。
- **未做**：未冻结 `I-008-001`；未放行 S1/S2 实施；未勾选检查点；Root R5 未勾选，Root 保持 `active / 4/5`；本目标保持 `active / 0/5`。
- **计划（非事实）**：按 P-004 §3.1 冻结 `I-008-001` 前询问用户是否补 self 审计；随后收集并冻结 `I-008-001`，再判断 S1/S2 实施边界。

## 2026-08-02 · 响应 A-002（pass 采纳 · R-003 handled）

- **A-002（independent · finding-closure · pass）响应**：采纳 `pass`——独立复核确认 A-001 F-001 `fixed` 关闭成立（D-002 + D-001 修订 + S2 + Root/I-005 投影对齐）、R-001/R-002 handled；本 scope 无开放 required。
- **R-003 → handled（投影/历史短句消歧）**：① `GOAL-008 00-meta` 概述改为「文档双进程为默认；Docker Compose 为 R5 必须交付的第二启动路径，fork 使用者可选」；② Root `00-meta` 进度说明由「R5 待立项」改为「R5 已立项 `GOAL-008-r5-engineering-fork`，待实施」；③ I-005 附件 v0.2.2 §2 末句改为过去时，明确精确镜像/Compose 契约由 `I-008-001` 冻结。
- **未做**：未冻结 `I-008-001`；未放行 S1/S2；未勾选检查点；Root R5 未勾选，Root 保持 `active / 4/5`；本目标保持 `active / 0/5`。
- **计划（非事实）**：冻结 `I-008-001` 前按 P-004 §3.1 询问用户是否补 self 审计；随后收集并冻结 `I-008-001`。

## 2026-08-02 · self 审计（A-003）+ 冻结 I-008-001（D-003）

- **用户裁决**：补同 scope `source: self` 审计并继续推进下一步（P-004 §3.1 闭环）。
- **A-003（self · goal-definition + design-plan 复核）**：落盘 `03-audit.md`——复核 GOAL-008 立项、I-005/D-013 方案边界与 A-001/A-002 响应闭合证据，`verdict: pass`；无开放 required；确认可进入 `I-008-001` 方案冻结。
- **I-008-001 收集与冻结**：
  - 静态核对 `apps/api/internal/config/config.go`、`.env.example`、`handler/health.go`、`apps/web/vite.config.ts`、`auth-client.ts`（相对 `/api` 路径）、`.github/workflows/r6-basic-matrix.yml`；落盘 [I-008-001-engineering-contract.md](attachments/I-008-001-engineering-contract.md)（v1.0.0）。
  - 契约冻结：env 键全集与 dev/prod 行为（含 fail-closed、DB_PATH 容器挂载 `/app/data`）、health/启动验证（`/healthz` + Web 静态页 + 登录终态）、Compose 服务（`api` 多阶段镜像 + `web` node→nginx，`compose.yaml` 于仓库根，DB 命名卷）、SPA fallback 与 `/api` 反代（nginx `try_files` + `proxy_pass http://api:8080`，同源免 CORS）、依赖/超时/失败行为（`depends_on service_healthy`、优雅停机、502、DB 持久化）、CI 容器/smoke 入口、验收清单 C-001～C-007。
  - 记录 **D-003**：`I-008-001` → **`verified`**（S1/S2 方案冻结门禁解除）。
- **未做**：未实施 S1（env 文档）或 S2（Dockerfile/compose）；未运行应用/测试/Docker；未勾选检查点；Root R5 未勾选，Root 保持 `active / 4/5`；本目标保持 `active / 0/5`。
- **计划（非事实）**：实施 S1（env 清单 + health/启动验证说明 + dev/prod 区分文档，对照 C-001/C-002）；随后实施 S2（Dockerfile × 2 + compose.yaml + nginx 反代 + CI smoke 入口，对照 C-003～C-007）；S2 完成后建议一次实施向审计。

## 2026-08-02 · 实施 S1 + S2（契约 C-001～C-007 验证通过）

- **S1 · 环境与配置基线（对照 C-001/C-002）**：
  - `apps/api/.env.example` 重写为 dev/prod 注解参考（全键与 `config.go` 一致；dev 默认 vs production fail-closed/必填）。
  - `apps/api/README.md` 新增「开发 vs 生产」表与「启动与健康验证」（`/healthz` 200 + 登录 + `/me` 判据）；`apps/web/README.md` 新增生产/compose 注记（相对 `/api`、同源反代）；根 `README.md` 新增「工程化与一键启动」节（Compose 第二启动路径、fail-closed 密钥）。
  - 验证：API `go build`/`go vet`/`go test ./...` 全绿（S1 未改 Go 源码）；`/healthz` 200 `{"status":"ok"}`（本机与容器内）。
- **S2 · 容器与一键启动（对照 C-003～C-007）**：
  - `apps/api/Dockerfile`（golang:1.26-alpine 多阶段 → alpine 运行镜像；`modernc.org/sqlite` 纯 Go → `CGO_ENABLED=0`；BuildKit cache mount + `go mod download` 重试；非 root `app` 用户；`DB_PATH=/app/data`）。
  - `apps/web/Dockerfile`（node:22 → nginx:1.27-alpine；**build context = 仓库根**，因 Vite/tsc `@schemas` 别名解析 `../../docs/schemas` 在 `apps/web` 上下文外；`COPY docs/schemas /docs/schemas`）；`apps/web/nginx.conf`（SPA `try_files` fallback + `location /api → proxy_pass http://api:8080` 同源反代）。
  - 根 `compose.yaml`（`name: schema-ui-core`；`api` 服务 + `web` 服务 `depends_on api healthy`；`db-data` 命名卷挂 `/app/data`；`/healthz` 与静态页 200 探针；`${AUTH_JWT_SECRET:?}`/`${ADMIN_INITIAL_PASSWORD:?}` fail-closed 插值）。
  - `.dockerignore`（根，服务 web 仓库根 context）+ `apps/api/.dockerignore`。
  - `.github/workflows/r6-basic-matrix.yml` 新增 `container-smoke` job（build → up → 等待 ready → 经 nginx `/api` 登录 + `/me` + SPA/route fallback + 记录跨 `api` 重启持久化 → teardown）。
  - **验证（本机 Docker 29.6 / Compose v5.3.1）**：`docker compose build` 两镜像成功（api 23.7MB）；`up -d` 后 api healthy、web 200；`/healthz` 200；登录 admin → `/me`；nginx 代理 `/api` 登录 + `/me` OK；`/list-edit-lifecycle` 刷新 fallback → index；`docker compose restart api` 后记录 `rec-9a440d7950daf745` 保持；`down`→`up` 后同一记录保持（DB volume 持久化）。**C-001～C-007 全部通过**。
- **实施细节留痕**：web 镜像构建曾因 `@schemas` 别名上下文外而失败 → 改为仓库根 context + `COPY docs/schemas`；`COPY nginx.conf` 路径按仓库根 context 修正；`go mod download` 遇 proxy.golang.org 瞬断 → BuildKit cache mount + 重试。
- **未做**：未实施 S3（fork 文档 + 15 分钟计时复现）、S4（`scripts/smoke.sh` 正式化）；`I-008-002`/`I-008-003` 仍 open；Root R5 未勾选，Root 保持 `active / 4/5`。
- **计划（非事实）**：下一拍收集并冻结 `I-008-002`（15 分钟计时复现协议 + smoke 判据），再实施 S3（QUICKSTART/fork 文档 + ≥1 次独立复现记录）与 S4（`scripts/smoke.sh` 正式化）；S2 已完成，建议一次实施向审计（self 或 `/audit`）。

## 2026-08-03 · 响应 A-004/A-005 + self 审计（F-002/F-003 fixed）

- **A-004（independent · execution-facts · pass）响应**：采纳 `pass`；R-001 → handled（与 A-005 F-002 同事实的分歧按用户裁决采用 `required/high` 口径，由 F-002 修复承载）；R-002 → fixed（根 README「Docker Compose 一键启动」段补 `.env` 免重复 export 注记）。
- **A-005（independent · execution-facts · fail）响应**：采纳 `fail`；**F-002 → fixed**、**F-003 → fixed**（详见下）。
- **self 审计（A-006 · execution-facts · pass）**：按 P-004 §3.1 用户裁决「需要补 self」补齐 S1/S2 实施 scope 的 `source: self` 覆盖（A-003 self 仅覆盖立项/方案边界）。
- **F-002 修复（生产运行时守卫，required/high）**：
  - `apps/api/internal/config/config.go` 新增 `ValidateProd()`：`AppEnv != "development" && AuthDevSessionEnabled` → 返回 `AUTH_DEV_SESSION_ENABLED must be false when APP_ENV=...` 启动错误；`cmd/server/main.go` 于 `config.Load()` 后立即校验，错误 → `logger.Error` + `os.Exit(1)`。
  - 新增 `apps/api/internal/config/config_test.go`：4 用例（development 允许 / production+flag fail-closed / production 无 flag 通过 / staging fail-closed）。
  - 运行时复验：`APP_ENV=production AUTH_DEV_SESSION_ENABLED=true` → `startup failed`（exit 1）；`APP_ENV=development AUTH_DEV_SESSION_ENABLED=true` → 正常启动（`dev_session: true`）。契约 §1/§5「生产禁止启用」现由硬门禁成立。
- **F-003 修复（进度投影同步，required/medium）**：`00-meta.md` frontmatter `progress: 0/5 → 2/5`，与勾选成功标准、派生进度段、`goal-tree.md`（`active / 2/5`）一致；复核无其它残留 `0/5`。
- **验证**：`go build` / `go vet` / `go test ./...`（apps/api）全绿（config 包 4 新增用例通过）；compose 与 api 镜像显式 `AUTH_DEV_SESSION_ENABLED=false`，守卫不影响第二启动路径。
- **未做**：未实施 S3/S4；`I-008-002`/`I-008-003` 仍 open required（阻断 S3/S4、S6 若实施）；未勾选新检查点（保持 `2/5`）；Root R5 未勾选，Root 保持 `active / 4/5`。
- **计划（非事实）**：建议对 F-002/F-003 关闭证据做一次 `/audit` finding-closure 复审；下一拍收集并冻结 `I-008-002`（计时复现协议 + smoke 判据），再实施 S3/S4。

## 2026-08-03 · 响应 A-007（finding-closure · pass）

- **A-007（independent · finding-closure · pass）响应**：采纳 `pass`——独立复核确认 **F-002**（生产环境开发会话硬门禁 `Config.ValidateProd()`）与 **F-003**（`00-meta` progress `2/5` 投影）的 `fixed` 关闭证据成立，两项维持闭合；本 scope 无开放 required、无新 required/recommended finding。
- **P-004 §3.1 处置**：同 scope（F-002/F-003 关闭证据）已有 **A-006（self）** 覆盖，无需再补自审；本 scope 无冲突、无 veto 触发。
- **未做**：未冻结 `I-008-002`；未实施 S3/S4；未勾选检查点（保持 `2/5`）；Root R5 未勾选，Root 保持 `active / 4/5`。
- **计划（非事实）**：下一拍收集并冻结 `I-008-002`（15 分钟计时复现协议 + smoke 判据，按 A-001 R-002 最低清单），再实施 S3（QUICKSTART/fork 文档 + ≥1 次独立复现记录）与 S4（`scripts/smoke.sh` 正式化）；建议对 `I-008-002` 冻结做一次方案冻结审计（self 或 `/audit`）。

## 2026-08-03 · 冻结 I-008-002 计时复现与 smoke 验收协议

- 新增 [I-008-002-fork-reproduction-protocol.md](attachments/I-008-002-fork-reproduction-protocol.md) v0.1.0，记录 D-013/A-001 R-002 所要求的工具/平台基线、依赖缓存前提、计时起止、失败/重试、独立复现字段、SM-001～SM-006、脚本输入/退出码和 disposable seed-reset 边界。
- 记录冻结时事实：repository revision 为 `5e27019482eb8d0695c402b784860233bbc90c39`，工作树 clean；Windows 11 x64、Go `1.26.0`、Node `22.17.0`、npm `10.9.2`、Docker `29.6.2`、Compose `v5.3.1` 和 curl `8.21.0`。WSL `bash` 启动失败，故未运行未来的 smoke 脚本；协议要求 S4 改在 Linux CI 或已验证可用的 Git Bash/WSL 执行。
- 记录 **D-004**：`I-008-002` → **`verified`**，解除 S3/S4 的方案冻结/实施前信息门禁。该状态只表示协议问题已得到可核对答案，不表示 S3/S4 已实施、`scripts/smoke.sh` 已存在、已通过 CI，或已完成 ≥1 次独立复现。
- **未做**：未改 QUICKSTART/README；未新建 `scripts/smoke.sh`；未执行 Docker/应用 smoke、浏览器验收或 15 分钟计时；S3/S4 未勾选，目标保持 `active / 2/5`，Root R5 保持 `4/5`。
- **计划（非事实）**：按 v0.1.0 实施 S3 的 fork 文档与独立复现记录，再实施 S4 脚本及本地/CI 证据；完成后做阶段审计。

## 2026-08-03 · 响应 A-008 + self A-009（协议 v0.1.1 · F-004 fixed）

- **A-008（independent · design-plan · conditional）响应**：采纳 `conditional`。**F-004 → fixed**、**R-008-001～004 absorbed**（详见下）。
- **self 审计（A-009 · design-plan 复核 · pass）**：按 P-004 §3.1 用户裁决「补 self」补齐 A-008 同 scope 的 `source: self` 覆盖（A-003 self 仅覆盖立项/方案边界），复核协议 v0.1.1 补丁与 F-004 关闭证据成立。
- **协议 v0.1.1 修订（D-005）**：
  - **F-004 fixed**：§5.1「S4 验收必须包含至少一次 disposable/隔离运行且 `SM-006=PASS`；非 disposable 默认路径的 exit 0 不得单独作为 S4「种子可重复」关闭证据」；§5.2 SM-006 标为「disposable 模式 · S4 必检」；§5.3 退出码 `0` =「SM-001～SM-005 通过，且 SM-006（disposable）通过——S4 完整绿」；保持「拒绝对普通开发库 destructive reset」安全边界（不安全 destructive → exit 2）。
  - **R-008-001**：§3.2 默认 URL 按路径区分（compose → API `:8080`/Web `:8081`；local-dual-process → API `:8080`/Web `${WEB_PORT:-5173}`）。
  - **R-008-002**：§3.2 操作化说明——D-013「后台首页可交互（列表加载）」= R4 代表页 `list-edit-lifecycle`（非 `overview`），QUICKSTART 必须覆盖该路由。
  - **R-008-003**：§5.2 SM-005 通过条件写为「响应体含 `id="root"`（或等价稳定标记）」。
  - **R-008-004**：§5.3 退出码 `2` 并入「不安全 destructive 模式」；`6` 仅指 SM-006 断言失败。
- **`I-008-002` 维持 verified（v0.1.1 权威）**；协议 frontmatter `related_decisions: D-004, D-005` + §7 修订记录。
- **未做**：未实施 S3/S4；未勾选检查点（保持 `2/5`）；Root R5 未勾选，Root 保持 `active / 4/5`。
- **计划（非事实）**：实施 S3（QUICKSTART/fork 文档 + ≥1 次独立复现记录，终点 4 = `list-edit-lifecycle` 列表加载）与 S4（`scripts/smoke.sh` + 本地/CI 全绿，含 ≥1 次 disposable `SM-006=PASS` 证据）；S4 完成后建议一次实施向审计。

## 2026-08-03 · 实施 S3 + S4（D-006）

- **S3 · fork 文档与 15 分钟体验**：
  - 根 [QUICKSTART.md](../../../../QUICKSTART.md) 交付：前置（Go/Node/Docker、依赖缓存不计时）、获取与准备（clone/checkout/`.env`）、双启动路径（A `docker compose` / B 本地双进程）、验收四终点（healthz → login → `/me` → 浏览器 `list-edit-lifecycle` 列表加载 `Acme Console`）、命令行与完整 smoke 用法、接业务指引。终点 4 操作化对齐 I-008-002 §3.2（`list-edit-lifecycle`，非 `overview`）。
  - **独立复现记录 [R5-S3-REPRO-001](attachments/R5-S3-REPRO-001.md)**（协议 §4 全字段）：`compose` 路径，`same-operator-clean-session`，先 `docker compose down -v` 清卷再全新 `up -d`；四终点全 PASS，单次计时 **34.5s ≤ 900s**（`03:04:45.596Z` 起点 → `03:05:20.141Z` 终点 4）。浏览器终点由 Playwright Chromium 驱动（[r5-repro-endpoint4.mjs](attachments/r5-repro-endpoint4.mjs)），标题 `List + edit lifecycle` + cell `Acme Console` 可见，截图 [r5-repro-endpoint4.png](attachments/r5-repro-endpoint4.png)。
- **S4 · 可复现验收（`scripts/smoke.sh`）**：
  - 新建 [scripts/smoke.sh](../../../../scripts/smoke.sh)（bash）：SM-001 参数/工具 → SM-002 readiness（30s）→ SM-003 代理登录 → SM-004 `/me` → SM-005 代表页 `id="root"` → SM-006 种子重复性（仅 `--disposable`，从空 DB 断言 `SMOKE_EXPECTED_SEED_TOTAL`（默认 8）条且含 `SMOKE_RECORD_ID`（默认 `rec-1`）/`Acme Console`，经 `SMOKE_RESTART_CMD` 重启后再次断言不重复播种）；退出码 `0/2/3/4/5/6/70`；不输出 token/password/secret；无 `--disposable` 不得执行种子 reset。
  - **本机验证（Git Bash + Docker，disposable）**：`bash scripts/smoke.sh --disposable` → `SM-001~006 全 PASS`，`EXIT=0`；log 见 [r5-smoke-disposable-local.txt](attachments/r5-smoke-disposable-local.txt)。非 disposable 默认路径亦验证（SM-001～005 PASS，SM-006 SKIP）。
  - **CI 接入**：`.github/workflows/r6-basic-matrix.yml` `container-smoke` 新增「Run S4 smoke (SM-001~006, disposable)」step：`SMOKE_USERNAME=admin SMOKE_PASSWORD=ci-admin SMOKE_RECORD_ID=rec-1 SMOKE_EXPECTED_SEED_TOTAL=8 bash scripts/smoke.sh --disposable`（runner 每次为隔离 project/volume → 满足「CI 默认 disposable」）；原 C-006 持久化 step 保留。
  - 说明：CI step 依赖 runner 上 job 级 `AUTH_JWT_SECRET`/`ADMIN_INITIAL_PASSWORD` 满足 compose fail-closed 插值；本地需在调用 shell 导出密钥。
- **回归**：`go test ./... -count=1`（apps/api）全绿；web `vitest run` **458/458**、`tsc`/build 未受影响（本轮未改产品代码）；workflow YAML 解析通过。
- **未做**：未实施 S5（阶段审计与 Root R5 勾选评估）；`I-008-003` 仍 open（S6 若实施）；Root R5 未勾选，Root 保持 `active / 4/5`。
- **已做**：S3/S4 检查点勾选（`2/5 → 4/5`），同步 `00-meta`、`01-decision` D-006、`03-audit`。
- **计划（非事实）**：实施 S5（对 R5 工程化交付做阶段审计 self + 视需要 independent，评估并记录 Root R5 勾选与 Root / VP-002 关门证据口径）；建议对 S3/S4 实施证据做一次实施向审计（self 或 `/audit`）。

## 2026-08-03 · 响应 A-011（F-005～F-009 fixed 路径整改）

- **P-004 §3.2 裁决**：A-010（self · pass）与 A-011（independent · fail）同 scope（execution-facts · S3/S4）verdict 冲突；用户书面裁决「按 **fixed** 路径整改 F-005～F-009」并确认「提交并推送 dev」。落盘 `01-decision` **D-007**。
- **协议 v0.1.2（D-007）**：[I-008-002 协议](attachments/I-008-002-fork-reproduction-protocol.md) 补丁——§5.3 新增部分绿退出码 `8`（非 disposable 不得 exit 0）；§5.1 冻结 disposable 隔离守卫（强制 `SMOKE_ISOLATION_ID` + `SMOKE_DISPOSABLE_CONFIRM=yes`、机器校验 project/卷绑定、禁止外部注入重启命令）；§5.2 SM-006 重启由脚本以隔离 project 执行 + readiness 重判。
- **F-007 修复（required/medium）**：`scripts/smoke.sh` 非 disposable 分支不再 `exit 0`——SM-001～005 通过但 SM-006 未运行时输出 `SMOKE RESULT: PARTIAL` 并 `exit 8`（部分绿，不是 S4 完整绿）；CI 只以 `--disposable` 完整绿（exit 0）作为 S4 证据。
- **F-008 修复（required/high）**：`scripts/smoke.sh` 重写隔离守卫——disposable 必须 `SMOKE_DISPOSABLE_CONFIRM=yes` 且 `SMOKE_ISOLATION_ID` 非空；`check_isolation()` 机器校验 `docker compose -p <id> ps -q api` 非空且 `docker inspect` 挂载含 `<id>_db-data`（或等价绑定），不满足 → `SM-001=FAIL` + exit 2；**删除 `SMOKE_RESTART_CMD` 任意 `eval`**，SM-006 重启改为脚本直接执行 `docker compose -p <id> restart api`；重启后 readiness 重新判定（失败 exit 3，吸收 R-011）。CI `container-smoke` 改用显式隔离 project `-p ci-container-smoke`（job env `SMOKE_ISOLATION_ID=ci-container-smoke`，smoke step 注入 `SMOKE_DISPOSABLE_CONFIRM=yes`），teardown 用同 project `down -v`。
- **本地守卫验证（四路径，log `r5-smoke-disposable-local-v0.1.2.txt`）**：① 非 disposable → SM-006 SKIP + `PARTIAL` + **exit 8**；② disposable 缺 `SMOKE_ISOLATION_ID` → `SM-001=FAIL` + **exit 2**；③ disposable 指定不存在 project → `SM-001=FAIL` + **exit 2**；④ disposable 隔离 project（`ci-smoke-local-verify` 全新卷）→ SM-001～006 全 PASS + `SMOKE RESULT: PASS` + **exit 0**（含 `isolation: project=ci-smoke-local-verify` 留痕）。
- **F-005/F-006 修复（required/high + required/medium）**：干净 worktree（`git worktree add` 于 commit `a086872`，detached clean）执行 **clean-ref 计时复现** [R5-S3-REPRO-002](attachments/R5-S3-REPRO-002.md)——计时起点 = `.env` 写入 + `docker compose up -d --build`（build/配置/迁移/种子/登录/页面加载全部计入，仅镜像层获取按 §3.1 排除）；四终点全 PASS，**13.5s ≤ 900s**（T0 `01:15:05.557Z` → T4 `01:15:19.092Z`）；记录按协议 §4 全字段落盘（clean ref、逐项 checks、smoke 输出引用、result 日志/截图路径、失败/重试记录）；完整命令与 BuildKit 输出见 [r5-repro-002-run.txt](attachments/r5-repro-002-run.txt)，截图 [r5-repro-002-endpoint4.png](attachments/r5-repro-002-endpoint4.png)。`R5-S3-REPRO-001` 加取代注记（历史保留）。
- **CI run 证据**：推送 commit `a086872` 触发 `r6-basic-matrix` run（初始 run 成功，见 `03-audit` A-011 响应节最终证据）。
- **未做**：未实施 S5；`I-008-003` 仍 open（S6 若实施）；Root R5 未勾选，Root 保持 `active / 4/5`；本目标保持 `active / 4/5`（S3/S4 检查点不因整改回退，验收有效性以 F-005～F-009 关闭证据为条件）。
- **计划（非事实）**：落盘 A-011 响应节（含 F-009 最终 CI run 证据）；建议对 F-005～F-009 关闭证据做一次 `/audit` finding-closure 复审，再进入 S5。

## 2026-08-03 · 响应 A-012（F-005 fixed：无项目编译缓存重做 S3 计时复现；R-012 handled）

- **A-012（independent · finding-closure · fail）响应**：采纳 `fail` 中 **F-005（required/high）** 与 **R-012（recommended/low）**；F-006～F-009 关闭证据维持（A-012 确认）。用户书面指示按 **fixed** 路径重做，完成后请求仅针对 F-005 的关闭复审。落盘 `01-decision` **D-008**。
- **禁用/隔离项目编译缓存（预 T0，排除项内）**：`docker rmi schema-ui-core-api:local schema-ui-core-web:local` + `docker builder prune -a -f`（清理 21.52GB BuildKit 结果缓存；基础镜像保持本地）。`docker compose up` 无 `--no-cache` 直传，故采用协议可陈述、可复核的 rmi + prune 做法。
- **重做 S3 计时复现 [R5-S3-REPRO-003](attachments/R5-S3-REPRO-003.md)**：clean worktree（`git worktree add --detach` 于 `1961e5a`，预 T0 `git status --short` 为空）；计时起点 = `.env` 写入 + `docker compose up -d --build`（协议 §3.2，`.env` 写入在计时内）；**四终点全 PASS，64.833s ≤ 900s**（T0 `02:09:49.734Z` / monotonic `403981233142700` → T4 `02:10:54.565Z` / `404046066135300`）；BuildKit 归档输出（[r5-repro-003-run.txt](attachments/r5-repro-003-run.txt)）中 **`go build`（#29，DONE 12.8s）与 `npm run build`（#38，DONE 6.1s，真实 vite 构建输出）均实际执行，正式 retry #3 的项目编译层均非 `CACHED`（该次仅一条非编译 `WORKDIR` 缓存）**——直接回应 A-012「编译层不得 CACHED」（响应 A-013 R-013 表述收窄）；截图 [r5-repro-003-endpoint4.png](attachments/r5-repro-003-endpoint4.png)（sha256 `89171fb1…809f8`，写于 worktree 外 gitignored 目录后归档）。
- **R-012 → handled**：预 T0 与运行后 `git status --short` 分别留痕（均空）；单调计时原始读数逐终点落盘（node `process.hrtime.bigint()` ns）；截图产物 hash 与路径记录。
- **失败/重试留痕（协议 §3.3，run log 内全记录）**：attempt #1——驱动脚本 PowerShell NativeCommandError（stderr 管道）中断，未测终点，复位后重试；attempt #2——T1/T4 窗口内 PASS 但 T2/T3 因 PowerShell 5.1 向 `curl.exe` 传内联 `-d` JSON 引号被吞、登录体未送达而失败（`login_http_ok=0`、`/me` UNAUTHENTICATED），对同一已构建栈以 `--data @file` 修正调用复验 `token_len=176` + `/me` 完整响应，按 §3.3 记为失败并复位重试；attempt #3——正式尝试单次通过。
- **未做**：未改协议 v0.1.2 正文、未改 `scripts/smoke.sh`、未改产品代码、未重开 F-006～F-009；未实施 S5；`I-008-003` 仍 open（S6 若实施）；Root R5 未勾选，Root 保持 `active / 4/5`；本目标保持 `active / 4/5`。
- **计划（非事实）**：同步 goal-tree；请求 `/audit` 仅针对 **F-005** 的 finding-closure 关闭复审；复审 pass 前不得推进 S5、勾选 Root R5 或关门。

## 2026-08-03 · 响应 A-013（F-005～F-009 维持闭合 · R-013 handled）

- **A-013（independent · finding-closure · pass）响应**：用户书面指示「复审通过，响应审计意见 A-013」。**verdict 采纳 `pass`**——A-013 独立二次复审确认 F-005 的 REPRO-003 关闭证据（clean checkout、预 T0 `git status --short` 空、单调原始读数、正式 retry #3 实际 `go build` 12.8s / `npm run build` 6.1s 非 `CACHED` 输出、64.833s ≤ 900s）补足 A-012 唯一 required/high 缺口；F-006～F-009 维持成立；**F-005～F-009 全部 `fixed` 闭合，本 scope 无开放 required**。A-011 关闭主张经 A-012（重做 F-005 证据）+ A-013（pass）完整闭环。
- **R-013（recommended/low · 缓存计数表述收窄）→ handled**：REPRO-003、A-012 响应节、D-008、02-execution、goal-tree 的「全程仅 1 条 `CACHED`」表述统一修正为「正式 retry #3 的项目编译层均非 `CACHED`，且该次仅一条非编译 `WORKDIR` 缓存」，与完整归档（含 retry #2 同类缓存）一致。
- **未做**：未推进 S5、未勾选 Root R5、未关门；`I-008-002` 维持 `verified`（v0.1.2）；`I-008-003` 仍 open（S6 若实施）；本目标保持 `active / 4/5`；Root 保持 `active / 4/5`。
- **计划（非事实）**：进入 S5（阶段审计与 Root R5 勾选评估）前，按 P-004 §3.1 询问用户是否补 self finding-closure 审计（F-005～F-009 关闭证据 scope 现无 `source: self` 覆盖）。

## 2026-08-03 · 补 self 审计（A-014）+ 实施 S5（A-015 阶段审计 · 4/5 → 5/5）

- **P-004 §3.1 用户裁决**：「补一次 self 审计，没有阻断项的话，推进 S5」。A-013 响应后本 scope 无开放 required、无冲突 → **A-014（self · finding-closure · F-005～F-009 · pass）** 落盘：self 复核 REPRO-003 正式 retry #3（pre-T0 `git status --short` 空、`go build` #29 DONE 12.8s / `npm run build` #38 DONE 6.1s 均非 `CACHED`、64.833s ≤ 900s、四终点 PASS、截图 sha256）、smoke.sh（非 disposable `PARTIAL` + exit 8；disposable 强制确认 + 隔离身份 + `check_isolation()` 机器校验 project/卷 + 去 eval 固定重启）、workflow 显式 `ci-container-smoke` project、CI run 30776646293 归档（SM-006=PASS + `SMOKE RESULT: PASS` + `down -v`）、本地四路径守卫日志——全部成立，与 A-013 同向无冲突，**无开放 required**。
- **实施 S5（阶段审计与 Root 关门条件评估）= A-015（self · stage-audit · pass）**：S1（C-001/C-002）、S2（C-003～C-007）、S3（QUICKSTART + REPRO-003 64.833s）、S4（smoke.sh + 本地四路径 + CI run 30776646293）逐项核对成立；意见台账 A-001～A-014 全部 responded/自审、F-001～F-009 全部 `fixed` 闭合、R-001～R-013 全部 handled/非阻断；信息门禁 `I-008-001/002` verified、`I-008-003` 仅 S6 适用（S6 可选不进分母）；**Root R5 勾选条件口径已记录**（S1～S5 全勾选 + 无开放 required + 至少一次关门向审计；Root R5 勾选、Root close-out 与 VP-002 关门须用户确认）。
- **检查点**：S5 勾选（`4/5 → 5/5`），同步 `00-meta`、`03-audit`、goal-tree；Root R5 检查点**未**勾选（Root 保持 `active / 4/5`）；本目标保持 `active / 5/5`（**未 `done`**——关门待 close-out 审计 + 用户裁决 + Root R5 勾选）。
- **未做**：未勾选 Root R5、未关门、未实施 S6（可选加分）；`I-008-003` 仍 open（S6 若实施）；未改协议/脚本/产品代码。
- **计划（非事实）**：由用户确认：① 是否勾选 Root R5 检查点；② 是否对 GOAL-008 做 close-out 关门审计（self 或 `/audit`）以进入 Root/VP-002 关门流程；③ S6（可选加分）是否纳入。

## 2026-08-03 · 实施 S6（可选加分 · 最小操作日志 · D-009 + I-008-003 v1.0.0）

- **用户裁决**：「推进 S6」。**`I-008-003` → verified**（D-009 + [I-008-003-operation-log-contract.md](attachments/I-008-003-operation-log-contract.md) v1.0.0——事件 6 个、0004 `operation_log` 迁移、repository 边界、接线与 best-effort 失败语义；HTTP 查询/清理/轮转为非目标）。
- **迁移 0004 `operation_log`**（`apps/api/internal/store/migrate.go`）：`operation_log(id, event CHECK 枚举, actor_id, actor_name, record_id, detail, created_at Unix ms)` + `idx_operation_log_created_at DESC`；既有库升级走 `pre-v0004` 快照（迁移运行器自动）。
- **repository**（`store/operations.go`）：`RecordOperation`（追加）+ `ListOperations(limit)`（created_at DESC, id DESC；limit ≤ 0 空）；事件常量 `records.create/update/delete`、`auth.login/logout/refresh`。
- **handler 接线**：
  - `records.go`：create/update/delete 成功写响应前记录（actor 取自 `requirePermission` 的 `account.User`；detail = `{"name":...}` 摘要）；`logOperation` best-effort（失败 → slog.Error，不阻断业务）。
  - `auth.go`：login/refresh 成功记录（detail = `{"username":...}`）；logout 记录（`auth.Authenticator.Logout` 改为返回被撤销 token 的 userID，I-008-003 §5）；`authsHandler` 注入 store。
- **测试**（全绿）：store `operations_test.go`（追加/排序/limit/CHECK 拒绝未知事件/0004 升级 `pre-v0004` 快照与台账）；handler `operations_test.go`（records 三写事件 + auth 三事件接线与 actor/record_id/detail、失败写不记日志）；迁移版本断言同步 `{1,2,3}→{1,2,3,4}`（migrate/restart/records_test）；`go test ./...`（apps/api）全绿、`go vet` 干净、`gofmt -l` 无输出、web vitest 458/458。
- **未做**：S6 不进 `progress` 分母（`5/5` 不变）；未新增 HTTP 端点/UI/清理策略；未勾选 Root R5、未关门；`I-008-003` 关闭后无其余信息门禁。
- **计划（非事实）**：S6 实施完成标注（00-meta）；Root R5 勾选与 GOAL-008 close-out 关门仍待用户确认。

## 2026-08-03 · 响应 A-016（F-010 fixed · R-014 handled）

- **A-016（independent · close-out · conditional）响应**：用户书面指示「响应 A-016 / F-010，走修复」（P-004 §3.3，未选 residual/overruled）。落盘 `01-decision` **D-010**。
- **F-010 修复（required/medium · auth logout 未记录冻结 username detail）**：`handler/auth.go` 新增 `authEvent(event, userID)`——logout 与 refresh 成功路径经 `store.UserByID` 解析**真实登录用户名**，按 [I-008-003 §3](attachments/I-008-003-operation-log-contract.md) 写入 `detail={"username":"admin"}` 与 actorName；login 保持 `creds.Username`；refresh 此前误用 `user.Name` 一并修正；actor 不可解析时仍记 actor_id + slog.Error（best-effort §5 不变）。
- **测试补齐**：`operations_test.go` 对 **login/refresh/logout 三类**统一断言 `"username":"admin"`（A-016 指出此前仅 login 断言 detail）；`go test ./...`（apps/api）全绿、`go vet` 干净、`gofmt -l` 无输出、web vitest 458/458 保持。
- **R-014 → handled**：S6 变更当前为未提交工作树（HEAD `851f9b6…`，15 files +284/−57），本地执行收据（apps/api 全绿 + web 458/458 + build 通过）记录为当前树候选事实，**不冒充** CI/容器验收；容器级证据待用户确认提交后补版本化 CI 或 disposable smoke。
- **未做**：未置 GOAL-008 `done`、未勾选 Root R5、未关门；`I-008-003` 维持 `verified`（契约 v1.0.0 未改）；本目标保持 `active / 5/5`；Root 保持 `active / 4/5`。
- **计划（非事实）**：按 P-004 §3.1 询问用户是否补覆盖 S6 的 self close-out 审计；随后由用户确认推进 GOAL-008 关门与 Root R5 勾选。

## 2026-08-03 · 合并响应 A-017/A-018 + GOAL-008 关门（done） + Root R5 勾选

- **用户裁决**：「补充自设计（补自审），没问题的话，合并响应交叉审计，开始关门」。
- **A-017（independent · finding-closure · pass）响应**：F-010 `fixed` 维持闭合——clean revision `eb6ff19` 上独立复核三条 auth 成功路径均满足 `I-008-003` §3 冻结 username detail；R-014 维持 handled（revision 收据固定 `eb6ff19`，不冒充 CI/容器验收）；R-015（recommended/low）非阻断。
- **A-018（self · close-out · S1～S6 · pass）**：按 P-004 §3.1 补齐 S6 scope 的 `source: self` 覆盖——S1～S5 核心（5/5）+ S6 加分全部成立；F-001～F-010 全部 `fixed`、R-001～R-015 全部 handled/非阻断；`I-008-001/002/003` 全部 `verified`；关门向审计齐备（A-013/A-014/A-015/A-016/A-017/A-018）；**关门条件满足**。
- **R-015 → handled**：`operations_test.go` `TestOperationLogAuthEvents` 提升为 JSON 解码精确断言（detail 仅 `username` 单字段、无 token/password/secret）并覆盖**轮换后有效 refresh token 的首次 logout**；`go test ./...`（apps/api）全绿 + `go vet` 干净 + `gofmt -l` 无输出 + web vitest 458/458。
- **关门（用户确认）**：GOAL-008 已置 `done`（`00-meta` status、`03-audit`、goal-tree 同步）；**Root R5 检查点已勾选**（Root `00-meta` R5 行 + `progress 4/5 → 5/5`）。
- **未做**：Root close-out 关门审计与 VP-002 关门为独立用户裁决（Root / 愿景层），未自动推进。
- **计划（非事实）**：由用户裁决 Root close-out 关门审计与 VP-002 关门（`/vision` 或 `/audit`）。

## 2026-08-20 · 协议 v0.1.3：W16-F01 首登改密 + 对齐现行 smoke/CI/端口（D-011）

- 按用户指令「好，继续」将 [I-008-002-fork-reproduction-protocol.md](attachments/I-008-002-fork-reproduction-protocol.md) 修订至 **v0.1.3**（D-011）。
- 修订要点：§3.2 默认 URL 更新为 25080/25081/25173 并写明 compose API 默认不发布宿主端口、S4 smoke 经 `scripts/pre-release-smoke.sh` loopback override；§5.1 新增 **W16-F01 首登强制改密必须走真实 `/api/account/password`（禁止清标志/改库/跳过）**；§5.2 SM-002 改 `/healthz`+`/readyz`、SM-004 内置改密、SM-005 代表路由 `/users`、SM-006 种子检查对齐 `SMOKE_SEED_ID=user-admin`/`SMOKE_EXPECTED_SEED_TOTAL=1`；§5.3 补 `SMOKE_PASSWORD_NEW`/`SMOKE_CSP`/`SMOKE_SEED_ID` 输入与退出码 `7`（SM-008）。
- 事实依据：本轮已实跑 `scripts/pre-release-smoke.sh` 全绿（SM-004 含 W16-F01 真实改密、SM-006、SM-008、C-006 persistence PASS），CI `container-smoke` 已统一走 wrapper；协议修订属文档对齐，不改代码/脚本/GOAL-008 status。
- 未做：未改产品代码/脚本；未重开 GOAL-008 关门；`I-008-002` 维持 `verified`（v0.1.3）。
