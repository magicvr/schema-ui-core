---
title: 执行记录 · R5 · 工程化、fork 体验与集成关门
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.1.4
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
