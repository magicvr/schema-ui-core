---
title: 执行记录 · R5 · 工程化、fork 体验与集成关门
status: active
created: 2026-08-02
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.1.8
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
