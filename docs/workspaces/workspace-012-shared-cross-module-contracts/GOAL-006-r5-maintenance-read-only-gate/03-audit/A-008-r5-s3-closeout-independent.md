---
id: A-008-r5-s3-closeout-independent
goal: GOAL-006-r5-maintenance-read-only-gate
doc: audit-entry
record_id: A-008
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R5 S3 close-out；四模式配置与 fail-closed、统一写门禁、maintenance vs degraded/read-only Host/API 消费、system-monitoring projection、正常态兼容、Profile/module/Manifest/protocol/readiness 不变式、A-002 F-001/F-002 实施闭合、全量 API/Web/build 证据与 git checkpoints
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-006-r5-maintenance-read-only-gate
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
reviews: A-007
---

# A-008 · R5 S3 independent close-out（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：close-out / S3
- **scope**：GOAL-006 最终 S3 cross independent 关门审计。核对 D-003 修订契约、四条成功标准、I-001～I-005、A-001～A-007 开放项，以及当前实现/测试/提交。
- **verdict**：**pass**
- **required findings**：0

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）。
- **covered**：GOAL-006 `00-meta` / D-001 / D-002（历史） / D-003（现行 accepted） / E-001～E-005 / A-001～A-007；Root `GOAL-001` 00-meta 与 R5 指针；VP-012 / alignment；实现：`config.go` + default/prod YAML、`operational.go`、`bootstrap.go`、`systemmonitoring.go`、`route_envelope.go`、`health.go`、`localize.go`、`errorcatalog`、`composition.go`、`systemmonitoring/provider.go` + page schema、相关 tests、`apps/web/src/host/bootstrap.ts` / `boot.ts` / `failure.ts`、`renderer/resource.ts`、`docs/schemas/capability-registry.json`。
- **excluded**：不改 `status` / `progress` / D-003 / `00-meta` / goal-tree / 业务代码；不读取或比较其他工作区上下文；不审 R6 API Token；不把 A-007 F-001 的后续 UI 横幅升为 required。
- **本轮复验**：独立复跑 `go test -timeout 15m -count=1 ./...`（apps/api，237s，全部 ok）；定向 `config` / `composition` / `docscheck` / `systemmonitoring` / handler R5+health；`npm test -- --run src/host src/renderer/error-localization.test.tsx src/renderer/resource.test.ts`（5 files / 60 tests）；`npm run build`（tsc + vite 成功）。验证 build 会把 claim `buildId` 推到当前 HEAD，已立刻还原那 3 个 protocol 文件，不把验证副作用留在工作区。

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-006 `parent`、`primary_plan` 一致；`goal-tree.md` 含本目标且 `status: active` |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 对齐链 | 未发现与 Root R5 / VP-012 / Charter 的明显冲突 | VP-012 `vision_ref` = `schema-ui-core-admin-foundation@0.2.0`；R5 是 maintenance 门控交付，无 Tier D、无 Profile/协议 pin 改写 |
| Vision Review required | 本 scope 未见开放 required | 本意见不审 Vision Review 本身 |
| 既有 Goal 审计 | A-001～A-007 无开放 required；A-002 F-001/F-002 已由 A-004 `fixed`；A-007 self = pass | `03-audit.md` |
| P-004 冲突 | 无 | A-007 与本条均为 pass；recommended 可叠加，无必改互否 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| 四模式单一来源；YAML 缺省 `normal`；`RUNTIME_MODE` 覆盖；显式空/未知 fail-closed | `config.go` L81–95、L190、L333–344、L651–652；`config.default.yaml` L93–95；`configs/config.yaml` L105–107；`config_test.go` L123–144 | 通过。`LookupEnv` 不用 `envOr`，空字符串不会静默回落 normal |
| 非法 mode 在 Load 写入 `LoadError`，`ValidateProd` 与 bootstrap 注册 fail-closed | `config.go` L341–343、L642–645；`bootstrap.go` L84–86；`bootstrap_test.go` L105–107 | 通过 |
| 门禁位于 request-id/CORS 之后、JSON 404/405 与 mux 之前 | `server.go` L24、L35–53；`composition.go` L505–507；`operational.go` L10–19、L39–47 | 通过 |
| 已注册 POST/PUT/PATCH/DELETE 在三种受控态拒绝；GET/HEAD 与探针放行 | `operational.go` L17、L50–57；`health.go` L47–48、L59–77、L78–107 | 通过 |
| allowlist 精确匹配五条认证/改密路径 | `operational.go` L61–71；公开写仅为 login/refresh/logout（`auth.go` L67–69）与 MFA verify（`mfa.go` L67）；`POST /api/account/password` 为强制改密 | 通过；匹配是 `r.URL.Path` 精确 switch，无前缀 |
| 未知路径/方法不匹配仍为 404/405，不伪装成运行态拒绝 | `route_envelope.go` L10–28；`operational_test.go` L64–80；`r5_operational_gate_test.go` L101–111 | 通过 |
| 拒绝统一 HTTP 503 + `SERVICE_*` catalog + `correlation_id`；无 Retry-After、不写 operation log | `operational.go` L22–35；`errorcatalog.go` L33–35；`localize.go` L17–22；`requestid.BodyName` = `correlation_id`；composition `assertOperationalError` | 通过 |
| Provider 真实写路径在 composition 黑盒中被门禁 | `r5_operational_gate_test.go` L68–73：`POST /api/data-dictionary/types` → 503 + code + request-id | 通过 |
| 正常态业务写仍走既有 401 | `r5_operational_gate_test.go` L92–99 | 通过 |
| bootstrap 投影：normal/maintenance 原样；degraded/read-only → 既有 `degraded`；省略 `disabledCapabilities` | `bootstrap.go` L76–86；`bootstrap_test.go` L75–104；struct `omitempty` 且未赋值 | 通过 |
| system-monitoring `availabilityMode` 为后端原始 mode；`status`/`ready` 仍只表达存储+模块图 | `systemmonitoring.go` L27–30、L107–135；`provider.go` L37–42、L70；composition 传入 `string(cfg.RuntimeMode)`；page schema 增加 Availability stat card | 通过 |
| Host：仅 bootstrap `maintenance` 进 `MAINTENANCE` 终态；degraded/read-only 为 `READY_DEGRADED` 进应用 | `bootstrap.ts` L273–276、L317–318；`boot.ts` L201–206、L256 | 通过 |
| 应用内写失败走 `ResourceApiError` 包络，不进 `classifyHostFetch` 的 5xx→unavailable | `resource.ts` L45–75、L168–182、L333–364；`failure.ts` L124–126 仅用于 Host fetch 分类；全仓生产路径无 `classifyHostFetch(` 调用 | 通过 |
| A-002 F-001 实施闭合：未把 `form.controls.readonly` 放入 `disabledCapabilities` | `bootstrap.go` L82–83；`capability-registry.json` L134–140 语义未改；R5 未改 `form-controls.ts` / Host registry | 通过 |
| A-002 F-002 实施闭合：degraded/read-only 写拒绝是应用内 503 catalog error | `operational.go` L22–35；D-003 L21；Host 5xx 映射未接到 Resource 写路径 | 通过 |
| Profile / 模块矩阵 / Manifest 聚合 / 协议 pin / health 探针未改 | `git diff 77d0b4c..3537779` 不含 `kernel/profile.go`、`manifest/`、`health.go`、`capability-registry.json`、`host-bootstrap.schema.json`、Host/resource 源码；claim 仅刷新 `buildId`，`artifactVersion` 仍为 `2.9.0` | 通过 |
| Git checkpoints 与 E-003～E-005 一致 | `c4856f2` 实现；`ffacbfb` composition 黑盒；`a687c05` claim；`3537779` E-005/A-007 | 通过 |

## 对照成功标准

| 标准 | 本轮 | 证据 |
|------|------|------|
| 1. 四种运行态具有单一后端来源、启动时校验与确定性公开投影；非法配置 fail closed | **达成** | Load/env/YAML 测试 + bootstrap 非法 mode 注册失败 + 四模式投影表 |
| 2. maintenance/read-only 的写拒绝覆盖核心及模块贡献路由；GET/HEAD、探针与错误优先级符合冻结契约 | **达成** | 统一 `WithOperationalGate` 包住整个 mux；Provider 真实写路径黑盒；handler 覆盖 GET 放行与 404/405；healthz 不碰库、readyz 不读 runtime mode |
| 3. bootstrap 与 system-monitoring 消费同一 `cfg.RuntimeMode`；正常态兼容；degraded 不产生未知 capability | **达成** | 同一 config 字段分别进入 `RegisterBootstrapWithAvailability` 与 monitoring provider；省略 `disabledCapabilities`；Host enum 未扩 `read-only` |
| 4. 定向与全量测试通过；Profile 默认集、模块矩阵、Manifest bytes/装配语义和既有协议 pin 不变 | **达成** | 本轮全量 API PASS；Web 5/60 PASS；`npm run build` PASS；R5 未改 profile/manifest/health/protocol schema |

## Findings

本轮 **无新 required finding**。

### F-001 · YAML 空/未知子测试被空 `RUNTIME_MODE` 污染

| 字段 | 值 |
|------|-----|
| level | **recommended** |
| 严重度 | low |
| status | open |
| 影响门禁 | 不阻断 S3 关门；测试隔离质量 |
| evidence | `apps/api/internal/config/config_test.go` L145–165 |

`empty yaml` / `unknown yaml` 用例在未声明 `env` 时零值为 `""`，随后无条件 `t.Setenv("RUNTIME_MODE", tc.env)`。`LookupEnv` 视其为已设置，LoadError 由空 env 触发，不能独立证明 YAML 空串或 `paused` 本身 fail-closed。

实现侧仍正确：`yf.Runtime.Mode != nil` 写入后再走同一 `ValidRuntimeMode`（`config.go` L333–344）。合法 YAML、env 覆盖、空 env、未知 env 已独立覆盖。后续把 YAML 用例改成「不设置 `RUNTIME_MODE`」即可补齐隔离。

### F-002 · composition 黑盒未钉死其余 allowlist / 中心注册核心写 / HEAD

| 字段 | 值 |
|------|-----|
| level | **recommended** |
| 严重度 | low |
| status | open |
| 影响门禁 | 不阻断 S3 关门；覆盖完备性 |
| evidence | `r5_operational_gate_test.go` L68–87；`operational_test.go` L28、L55–59；`operational.go` L50–71；`upload.go` L240 |

S2 矩阵已证明：Provider `POST /api/data-dictionary/types` 被拒、login 保留自身 400、health GET 200、未知 POST 404、`POST /healthz` 405。handler 单测另覆盖合成 POST 门禁与 login allowlist。

未单独黑盒：`refresh` / `logout` / `mfa/verify` / `account/password`、中心注册的 `POST /api/upload`、HEAD、以及 PUT/PATCH/DELETE。门禁是单一 mux wrapper + 精确 path switch，代码审查足以认为这些路径与已测 POST 同机制。补测可降低回归成本，不是契约缺口。

### 历史 finding 处置

| 原 finding | 原级别 | 本轮 | 说明 |
|------------|--------|------|------|
| A-002 F-001 | required / high | **保持 fixed** | 实施未把任何协议 capability 当作运行态开关 |
| A-002 F-002 | required / med | **保持 fixed** | 503 catalog + Resource 包络；maintenance 才是 Host 终态 |
| A-002 F-003～F-006 | recommended | 已吸收为实施门 | 503 避开 423；空 env 不用 `envOr`；精确方法/allowlist |
| A-005/A-006 S3 消费门 | recommended | 本轮闭合为证据 | Host/resource 路径核验 + 全量回归 |
| A-007 F-001 | recommended | 维持非阻塞 | `SERVICE_*` 若进入 HostFailureScreen 再补文案；当前 Resource + status 已足够 |

## 信息项核对（P-005）

| ID | 级别 | 最晚阶段 | 登记状态 | 本审计结论 |
|----|------|----------|----------|------------|
| I-001 | required | S0 结束前 | verified | 维持；扫描事实仍成立 |
| I-002 | required | S0 结束前 | verified | 维持；四模式 503/code、无 Retry-After、消费边界已实施 |
| I-003 | required | S0 结束前 | verified | 维持；门禁先于 handler 认证；allowlist 覆盖公开写与强制改密 |
| I-004 | required | S0 结束前 | verified | 维持；只用既有 `degraded` mode，省略能力裁剪 |
| I-005 | required | S3 关门 | verified | **本条即为所要求的 independent 关门审**；模式 cross，provider = 项目级 grok-build |

无 `deferred`。无用户书面 `accepted-residual`。无到期未关闭 required 信息项。

## 本轮复验

| 命令 | 结果 |
|------|------|
| `apps/api`: `go test -timeout 15m -count=1 ./...` | PASS（handler 232.599s；composition 30.415s；全包 ok） |
| `apps/api`: `go test ./internal/config ./internal/composition ./internal/docscheck ./internal/modules/systemmonitoring` | PASS |
| `apps/api`: `go test ./internal/handler -run "RuntimeMode\|OperationalGate\|RegisterBootstrap\|SystemMonitoring\|ErrorCode\|ErrorCatalog\|Health"` | PASS |
| `apps/web`: `npm test -- --run src/host src/renderer/error-localization.test.tsx src/renderer/resource.test.ts` | PASS（5 files / 60 tests） |
| `apps/web`: `npm run build` | PASS（tsc + vite）；claim 仅会改 `buildId`，已还原，未提交 |

## 必改项汇总

**无。** required = 0。

两条 recommended（F-001 测试隔离、F-002 矩阵补测）不阻断关门。

## 与既有意见的异同

| 点 | A-007 self | 本意见 |
|----|------------|--------|
| S3 消费边界 / F-001/F-002 实施闭合 | pass | 同意；独立复跑全量 API 与 Web 定向/build |
| 装配/协议/readiness 不变式 | pass | 同意；用 R5 提交范围 diff 复核 |
| 测试缺口 | 未开新 finding | recommended F-001/F-002（隔离与覆盖，非契约失败） |
| A-007 F-001 UI 文案 | deferred non-blocking | 维持；不升格 |
| verdict | pass | **pass** |

不是 P-004.2 冲突。

## 结论 + 建议给编排器/用户的下一步

R5 S3 关闭证据充分：四模式配置 fail-closed、统一写门禁覆盖 mux 上的核心与 Provider 路由、maintenance 与 degraded/read-only 的 Host/API 消费边界符合 D-003，system-monitoring 投影原始 mode，正常态与 Profile/Manifest/protocol/readiness 不变式保持，A-002 F-001/F-002 在实现上可重复核对为 `fixed`。全量 API 与 Web 定向/build 本轮独立通过。

建议 `/govern` 下一句：响应 A-008（pass，required=0），将 GOAL-006 标为 `done` 并同步 `goal-tree`；F-001/F-002 recommended 可记入后续清理或接受为残余测试债。不要改写 D-003。

## 声明

本意见不修改 `status` / `progress` / D-003 / `00-meta` / goal-tree / 业务代码。响应与关门由 `/govern` 处理。
