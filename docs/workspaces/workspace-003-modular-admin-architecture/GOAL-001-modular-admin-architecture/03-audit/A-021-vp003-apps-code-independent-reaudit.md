---
id: A-021-vp003-apps-code-independent-reaudit
doc: audit-entry
goal: GOAL-001-modular-admin-architecture
source: independent
auditor: Grok Build / grok-4.5 / high（独立复审，未加载任何 skill）
date: 2026-08-06
scope: >
  VP-003 exit #1~#7 与 apps/api、apps/web 代码事实逐条对照；阻断性 bug 排查；
  双 Profile 运行、数据生命周期、聚合 Manifest、fail-closed、退役路径残留。
  动态执行：go build / go test / go vet / npm build / vitest /
  Playwright（mvp+admin）/ 本地冒烟（双 Profile、同库升级、毒库迁移）
audit_type: code-audit
verdict: pass
status: recorded
parent: null
created: 2026-08-06
updated: 2026-08-06
version: 0.2.0
---

# A-021 · apps/api + apps/web 代码对 VP-003 意图的独立复审

- **source**：independent
- **auditor**：Grok Build / grok-4.5 / high（用户指定独立审计；**未加载**
  audit/govern 等任何 skill，直接按 AGENTS 治理规则执行）
- **类型 / scope**：code-audit；VP-003 七条方向级退出判据 + 模块契约 + 数据
  生命周期 + 聚合 Manifest + 安全边界 + 旧路径移除，对照 `apps/api` 与
  `apps/web` 工作树代码事实；含阻断性 bug 排查
- **verdict**：**pass**（无 required / 必改 findings；无阻断 bug；无目标漂移）
- **审计对象**：工作树 HEAD `6ed88241fe1cea20b174b224467107ab99f91e81`
  （`docs(governance): record root closeout checkpoint`）
- **方法约束**：与 A-019 只读静态核验不同，本条为**动态复核**：构建、测试、
  浏览器 E2E 与冒烟在本次审计中实际执行（见 §2）。本意见**不**修改
  status / progress / goal-tree / 代码 / `00-meta` / `01-decision` /
  `02-execution` / VP 状态；按用户裁决条件，未发现漂移/未达标/阻断 bug，
  **不**回退 Root 关门状态。

## 1) 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-003-modular-admin-architecture` |
| canonical | `docs/workspaces/workspace-003-modular-admin-architecture/` |
| Root | `GOAL-001-modular-admin-architecture`（`parent: null`，`done / 6/6`） |
| 对照历史 | A-010（独立，apps 内聚债，已全部 fixed）、A-018/A-019/A-020（Root close-out） |
| 代码范围 | `apps/api`（kernel/composition/modules/handler/store/manifest/config/server/auth）、`apps/web`（app/account/protocol/renderer/main） |
| 排除 | 本条**不**放行 VP-003 `closed`（归 `/vision`）；**不**重写任何历史审计意见 |

## 2) 本次实际执行的动态验证（证据）

| # | 验证项 | 命令/方式 | 结果 |
|---|--------|-----------|------|
| V-1 | API 编译 | `go build ./...`（go1.26.0） | ✅ exit 0 |
| V-2 | API 全量测试 | `go test ./...`（kernel/composition/store/handler/modules 全覆盖） | ✅ 全 ok |
| V-3 | API 静态检查 | `go vet ./...` | ✅ exit 0 |
| V-4 | Web 构建 | `tsc -b && vite build` | ✅ exit 0 |
| V-5 | Web 单元/集成 | `vitest run` | ✅ 24 文件 / 495 用例全过 |
| V-6 | 浏览器 E2E · mvp | Playwright chromium，`WEB_PORT=9999`，`APP_PROFILE=mvp` | ✅ 2/2（schema-crud + shell） |
| V-7 | 浏览器 E2E · admin | 同一前端，`APP_PROFILE=admin` | ✅ 2/2 |
| V-8 | 冒烟 · mvp | 全新库：登录前 Manifest 200；pages=data-table,form-controls,form-with-reactions,overview,roles,search-form-table,users（**无 settings/activity**）；`/api/users`→401、`/api/settings`→404、`/api/operations`→404；登录后 `operation_log` 有 `auth.login` 行（无 Activity UI 仍记日志） | ✅ |
| V-9 | 冒烟 · admin | 全新库：pages 含 activity/settings；`/api/branding`→200（公开）；`/api/settings`、`/api/operations`→401 | ✅ |
| V-10 | 同库升级 mvp→admin | 同一 SQLite：mvp 启动后台账含迁移 1–9（含 `site_settings`/`records_retire`/`system_data_reconcile`）→ admin 重启后 settings/activity 页面出现 | ✅ 编译全局迁移与启用状态无关 |
| V-11 | fail-closed × 4 | 未知 Profile（`PROFILE_UNKNOWN`）、custom 缺模块（`PROFILE_MODULES_REQUIRED`）、未知模块 ID（`MODULE_UNKNOWN`）、台账注入未知已应用迁移 version 10（`LIFECYCLE_START_FAILED … unknown applied migration version 10`） | ✅ 全部阻断启动 |
| V-12 | Manifest 缓存 | ETag 稳定；`If-None-Match` → 304 | ✅ |
| V-13 | 退役残留扫描 | `MountProviderRoutes` / `RegisterSettings` / `RegisterActivity` / `staticSchemaDocuments` / `schemaDocumentsForPlan` / `seedRBAC` 于 `apps/api/internal`（非测试） | ✅ 无活动符号 |
| V-14 | 静态 Manifest 残留 | `apps/web/public/` 无 manifest **文件**；Dockerfile `test ! -e dist/.../app-manifest.json`；nginx 精确 `location =` 反代 API | ✅（空目录见 R-021-001） |

## 3) VP-003 exit #1～#7 逐条对照

| 退出判据 | 代码事实（独立结论） | verdict |
|----------|----------------------|---------|
| #1 单主线与 Profile | `kernel/profile.go`：`mvp`/`admin`/`custom` 版本化配置；`ResolveProfile` 优先级 compiled-profile-default → modules.enabled → environment；`Resolve` 校验依赖闭包/循环/未知 ID fail closed；同一 `cmd/server` 二进制双 Profile 运行（V-8/V-9） | pass |
| #2 薄内核、组合根与模块契约 | `kernel` 包零导入 `internal/modules`；`composition` 为唯一 Fx 组合根；模块 provider 无 `go.uber.org/fx` 运行时导入；`StandardAdminCapabilities()` 强制核心六项；users/roles/settings/activity 经 `RegisterContributions` 自注册；Descriptor 与 Plan 全字段 `descriptorsMatch` | pass |
| #3 数据生命周期 | `store/migrate.go` 平台 runner + 模块拥有 Apply；全局台账/checksum/未知已应用迁移 fail closed（V-11）；升级前快照 + integrity_check；`compiled/persistence.go` 编译全局收集（V-10 mvp 库已含 settings 迁移）；`records_retire` tombstone 保留；bootstrap 与 reconcile 分离，`system_data_grants` 追踪授权 | pass |
| #4 后端聚合运行时契约 | `manifest` + `handler/manifest`：确定性聚合、冲突 fail closed、登录前无秘密、`/.well-known` + ETag/304（V-12）；Vite/nginx 精确代理；前端 `validateAppManifest` protocol 2.7 精确匹配、capability 缺失 fail closed；生产无静态 manifest 文件兜底（V-14）；同一前端 build 双 Profile（V-6/V-7） | pass |
| #5 安全、横切与生命周期 | auth：HS256 短时 access + 旋转 refresh、bcrypt、401/403；`ValidateProd` 生产密钥强度；`operationlog` 两 Profile 恒启用（V-8 mvp 有 login 日志且 `/api/operations` 404）；activity 只读模块；Settings 经 `X-Schema-UI-Config-Changed` 通用事件；Start/Ready/Stop 带 `module_id`；`/readyz` 模块图 + system-data 就绪 | pass |
| #6 现有能力迁移与旧路径退出 | users/roles/settings/activity 为 kernel.Provider；核心六项齐全；records 无产品面（迁移 0006）；中央业务 Register 仅剩 core auth/health；静态生产 Manifest、Shell 特例、退役符号无活动残留（V-13/V-14；CI retirement 扫描同口径） | pass |
| #7 可 fork、可运维、可回归 | `QUICKSTART.md`、`compose.yaml`、双 Dockerfile、`.github/workflows/r6-basic-matrix.yml`（单元/构建/vet + mvp+admin 浏览器矩阵 + 容器 smoke + 退役扫描）、`scripts/smoke.sh` | pass |

## 4) Findings

### 无 required / 必改 findings

逐条核验未发现与 VP-003 意图不符的目标漂移、退出判据未达标或阻断性 bug
（含安全边界、数据完整性、fail-closed 语义）。历史 A-010 内聚债
（F-001/F-002/F-003b/F-005 等）在代码中已按 R6 承诺落地：store 不再上帝对象
（平台 runner 与领域仓库分离）、`CollectPersistence` 已由 composition 生产接线、
Schema 文档字节来自 ContributionSet、seed 由 bootstrap/reconcile 贡献驱动。

### recommended（不阻断；不构成关门回退理由）

- **R-021-001（recommended）**：`apps/web/public/.well-known/schema-ui/` 与
  `dist/.well-known/schema-ui/` 为 R6 移除静态 fixture 后的**本地空目录残留**
  （git 未跟踪；Vite 会拷贝空目录）。生产路径不受影响（无 manifest 文件、
  nginx 精确 location 反代、Dockerfile 断言），建议顺手删除以消除误解。
- **R-021-002（recommended · 措辞澄清）**：`module-architecture.md` §7 与
  VP-003 exit #5 提到「日志**与指标**均携带 `module_id` 语义」，但当前实现
  无任何指标基础设施（grep `metric|prometheus|expvar` 零命中）；日志（带
  `module_id`）、健康诊断（healthz/readyz 模块图门控）均已落实。按 §2.2
  Observability 为按需能力、且无指标贡献契约，不构成退出判据缺口；建议在
  Root 决策或 VP 修订中显式写明「指标 = 按需，当前无指标贡献」，避免后续
  审计对 exit #5 措辞产生歧义。

## 5) 结论

- **verdict：pass**（required 0；recommended 2：R-021-001 / R-021-002）
- **未发现**目标漂移、未达标或阻断性 bug → 按用户裁决条件，**不**回退
  Root `done / 6/6` 关门状态，**不**改动 status/progress/goal-tree。
- 本意见不自动放行 VP-003 `closed`（该门禁仍归 `/vision`，以七条退出判据
  的 Q2 证据台账为准）。
- 响应归 `/govern` 编排器（若有）；R-021-001/R-021-002 为 recommended，
  不阻断任何阶段。
