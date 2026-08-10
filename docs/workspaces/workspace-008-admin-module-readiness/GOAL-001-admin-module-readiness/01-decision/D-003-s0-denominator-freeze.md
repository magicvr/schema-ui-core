---
id: D-003-s0-denominator-freeze
doc: decision-entry
goal: GOAL-001-admin-module-readiness
status: accepted
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
---

# D-003 · S0 准入分母与门禁冻结

## 用户确认

用户于 2026-08-10 按 P-004 确认开设 `GOAL-002-s0-denominator-freeze` 并推进 S0 分母冻结，候选基线为 **commit `852ee7e`（clean checkout）**。本决策落盘 S0 冻结的准入分母与门禁字段；冻结完成后 Root 上 S0 到期的 required 信息项 `I-READINESS-001/004/005/006/007/008/009` 据此 `verified`。

## 1. 代码与环境分母

| 字段 | 冻结值 | 证据 |
|------|--------|------|
| 候选 Git commit | `852ee7e9dc2dbabb72542c6a255fb095aa3b57c4`（short `852ee7e`；分支 `dev`） | `git rev-parse HEAD`；`git status --porcelain` 空 |
| 来源身份 | **clean checkout**（porcelain 为空，无 staged/未跟踪/工作树修改） | VP-008 §证据基线有效性 默认规则 |
| Go | `1.26.0`（`apps/api/go.mod` `go 1.26`） | `go version` / `go.mod` |
| Node | `22.17.0`（CI `node-version: "22"`、`node:22` Dockerfile） | `node --version` / CI / Dockerfile |
| npm | `10.9.2` | `npm --version` |
| Docker | `29.6.2`（compose v2） | `docker --version` |
| Go 依赖锁 | `apps/api/go.sum`（committed） | 文件 |
| npm 依赖锁 | `apps/web/package-lock.json`（committed） | 文件 |
| 运行时工具链固定 | 无 `.nvmrc`/`.node-version`/`engines`；Node 22 仅由 CI 与 Dockerfile 固定（观察项，S1 登记） | glob 确认 |
| 配置模板 | `apps/api/.env.example`（参考；Go API 不自动加载）；compose 读取仓库根 `.env` | 文件 / README |
| 数据库 | SQLite 内嵌（`modernc.org/sqlite v1.55.0` 纯 Go 驱动），无外部 DB 服务 | `go.mod` |
| DB 起始形态 | fresh 库：`openStore` 执行待迁移，`WasFresh()` → `authsessiondata.Bootstrap` 种子 `user-admin`（username `admin`，roles `admin`+`editor`）与系统角色 `admin/editor/viewer`；随后 `Reconcile` 应用权限/导航授权；`MarkSystemDataReady` | `composition.go` `openStore` / `authsession/systemdata/bootstrap.go` |
| 迁移账本 | **0001–0010**（0001 r2_baseline · 0002 rbac_expand · 0003 records_persist · 0004 operation_log · 0005 operation_log_expand · 0006 records_retire · 0007 site_settings · 0008 operation_log_settings · 0009 system_data_reconcile · 0010 site_settings_v2）；全局唯一、版本连续、checksum 64-hex | `modules/{authsession,corepersistence,operationlog,settings}/migration/*.go` / `migration/collector.go` |
| 数据迁移快照 | 每个待执行数据迁移在非空库旁生成 `schema-ui.db.pre-vNNNN-<UTC>.sqlite` 并做完整性检查 | `store/migrate.go` / QUICKSTART §4 |

## 2. 验证命令矩阵（机器可复跑）

| # | 命令 | 层 | 覆盖 | 候选 `852ee7e` 实测 |
|---|------|----|------|---------------------|
| V-001 | `cd apps/api && go build ./...` | build | API 编译 | ✅ exit 0 |
| V-002 | `cd apps/api && go test ./...` | unit/integration | API 全部包（含 cmd/server 重启、store 迁移、kernel、composition、modules） | ✅ 全包 pass |
| V-003 | `cd apps/api && go vet ./...` | static | API vet | ✅ exit 0 |
| V-004 | `cd apps/web && npm test` | unit/integration | Web vitest（含协议 conformance、renderer、i18n、D-PERM） | ✅ 40 files / 728 tests pass |
| V-005 | `cd apps/web && npm run build` | build | Web tsc + vite 构建 | ✅ vite 6.4.3，1834 modules，3.61s |
| V-006 | `cd apps/web && npm run test:e2e`（APP_PROFILE=mvp / admin） | browser e2e | shell / schema-crud / localization | ✅ mvp：3 pass + 1 skip（admin-only 用例）；admin：3 pass + 1 skip（mvp-only 用例） |
| V-007 | `bash scripts/smoke.sh`（SM-001~005 + 可选 SM-007） | smoke | readiness/登录/身份/代表页/Profile-Manifest 合同 | ✅ mvp 与 admin 均 SM-001~005+SM-007 PASS（exit 8 部分绿，符合非 disposable 语义） |
| V-008 | `bash scripts/smoke.sh --disposable`（SM-001~007 + SM-006，隔离 compose） | smoke+种子 | 种子可重复性 + 重启持久化 | ✅ exit 0（SM-001~006 完整绿；`ci-smoke-s0` 隔离 project，重启后种子断言通过） |
| V-009 | CI `r6-basic-matrix.yml` 四作业 | CI | web / api / browser-e2e / container-smoke | 本地重跑等价项见 V-001~V-008 |

## 3. 模块集合与适用矩阵

| 模块 id | Version | KernelAPIRange | 分级 | 核心六项 | Profile 默认 |
|---------|---------|----------------|------|----------|--------------|
| `core.server-registration` | 2.0.0 | >=2.0 <3.0 | infra | HTTP | mvp/admin |
| `core.auth-session` | 2.0.0 | >=2.0 <3.0 | infra | Authorization, Persistence | mvp/admin |
| `core.schema-render` | 2.0.0 | >=2.0 <3.0 | infra | Schema, Validation | mvp/admin |
| `core.manifest-route` | 2.0.0 | >=2.0 <3.0 | infra | Manifest | mvp/admin |
| `core.navigation-capability` | 2.0.0 | >=2.0 <3.0 | infra | Navigation, Expressions | mvp/admin |
| `core.operationlog` | 2.0.0 | >=2.0 <3.0 | infra | OperationLog | mvp/admin |
| `admin.users` | 2.0.0 | >=2.0 <3.0 | standard-admin | HTTP+Schema+Authz+Nav+Manifest+Persistence | mvp/admin |
| `admin.roles` | 2.0.0 | >=2.0 <3.0 | standard-admin | 同六项 | mvp/admin |
| `admin.settings` | 2.0.0 | >=2.0 <3.0 | standard-admin | 同六项 | admin |
| `admin.activity` | 2.0.0 | >=2.0 <3.0 | standard-admin | 同六项 | admin |

- **Profile 默认启用集**：`mvp` = 6 core + users + roles（8）；`admin` = mvp 8 + settings + activity（10）；`custom` 无显式 `APP_MODULES_ENABLED` 时 fail-closed（`PROFILE_MODULES_REQUIRED`）。解析优先级 `compiled-profile-default → modules.enabled → environment`。
- **标准 Admin 六项**（`kernel.StandardAdminCapabilities()`）：`http, schema, authorization, navigation, manifest, persistence`；standard-admin 模块 `Requires` 全套。
- **`other` 分级**：当前编译名册无未归类模块；`other` 标签仅作临时发现标签，S2 前必须收敛。
- **迁移 ownership**：`core.auth-session`（0001/0002/0009）、`core.persistence`（0003/0006）、`core.operationlog`（0004/0005/0008）、`admin.settings`（0007/0010）；全局收集不按 Profile 过滤（`compiled.PersistenceProviders()`）。

## 4. 运行形态分母

- **真实 Profile**：`mvp`（8 模块）与 `admin`（10 模块）为编译候选集；同一 Web 构建随 Profile 切换页面集（`app-manifest.mvp.json` 10 页 / `app-manifest.admin.json` 12 页）。
- **`custom` 失败语义**：未知 Profile → `PROFILE_UNKNOWN`；`custom` 无模块 → `PROFILE_MODULES_REQUIRED`；依赖缺失/配置冲突 → `MODULE_*`/`PROFILE_*` fail-closed。这些进入内核测试分母（`kernel/profile_test.go` 等）。
- **Manifest 发布**：API 聚合经 `/.well-known/schema-ui/app-manifest.json`（header `X-Schema-UI-Manifest-Source: api`）；禁止生产静态 Manifest 兜底（Dockerfile 断言 `! -e dist/.well-known/...`）。

## 5. 协议面分母

| 主张 | 冻结投影 |
|------|----------|
| 协议 pin | `schema-ui-docs` `v2.7.0` / commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`（`I-PROTO-FULL-001`） |
| 能力域 | 12/12 |
| registry type | 24/24（`component-registry.json` 24 types） |
| conformance 套件 | 16/16（`apps/web/src/protocol/upstream/*.cases.json`） |
| case 总量 | 320/320（318 执行 + 2 排除） |
| 本地 disposition | `upstream-fixtures.test.ts`：2 项 app-manifest case exclude（`CAPABILITY_REQUIRED` vs `MISSING_REQUIRED_CAPABILITY` 差异）；`stage3-fixtures.test.ts` 零排除 |
| Manifest `protocolVersion` | `"2.7"`（非 `2.7.0`；一致） |

## 6. 主流程与用例选取规则

- **smoke 下限**：`scripts/smoke.sh` SM-001（参数/工具/安全前提）→ SM-002（readiness 30s）→ SM-003（代理登录）→ SM-004（当前身份）→ SM-005（代表页 SPA）。`--disposable` + 隔离 compose project 才运行 SM-006（种子可重复性 + 重启持久化）；SM-007（Profile/Manifest 合同）在提供 `SMOKE_EXPECTED_PROFILE` 时运行。退出码 0/2/3/4/5/6/8/70。
- **代表页清单**：每个实际可达 Runtime Manifest `pageId` 与 `schemaUrl` 进入清单。admin 12 页（overview / data-table / search-form-table / form-controls / form-with-reactions / form-with-upload / data-display / admin-list-batch / users / roles / settings / activity）；mvp 10 页（不含 settings/activity）。
- **CRUD 覆盖规则**：每个声明 CRUD 的资源（users、roles、settings）至少覆盖 list、detail、create/update、delete（或不支持时记录理由）；每个可达写操作至少覆盖成功、未授权/权限失败或校验失败路径。写入权限：`users.read/write`、`roles.read/write/assign`、`settings.read/write`、`operations.read`。
- **用例固定字段**：每条用例固定 `profile`、`pageId`/资源 id、`schemaUrl`、权限键、预期错误码、证据路径；未列入项只能以明确 residual 留痕。

## 7. 消费路径与升级边界

- **compose/容器启动**：仓库根 `compose.yaml`（api + web；SQLite 命名卷 `db-data`；fail-closed `AUTH_JWT_SECRET`/`ADMIN_INITIAL_PASSWORD`；`/readyz` healthcheck）。
- **fork bootstrap**：QUICKSTART（≤15 分钟四终点验收）；本地双进程（API :25080 / Web dev :25173）或 compose（:25080/:25081）。
- **升级边界**：支持升级来源/目标 = 迁移账本 0001→0010（版本化 reconcile）；升级前 `pre-vNNNN` 快照/恢复路径已定义；**本 VP 不要求降级**（用户未扩 scope）。
- **fork/compose 基线**：指向文档化基线分支/commit；超出升级窗口/恢复边界/fork 基线的兼容诉求记录为 `N/A`/residual 或回 `/vision`。

## 8. 跨模块 UI 可访问性下限（S0 冻结）

| 共享宿主 | 最低可判定项 | 证据形式（S0 冻结） |
|----------|--------------|---------------------|
| Renderer/Shell 布局、导航与移动导航 | 键盘可达全部交互控件；焦点可见、顺序稳定；移动导航开/关后焦点去向可预期 | 静态 a11y 断言（`shell.test.ts`：hamburger/close aria-label 经 i18n catalog）+ 必要人工核对 |
| schema-driven 表单/列表/详情/动作 | 控件有可计算名称、角色、状态；校验/禁用/加载状态可观察；错误不只依赖颜色 | 渲染器语义/状态断言（`render.test.tsx`、`schema-crud.test.tsx`、`schema-table.test.tsx`）+ 代表页人工核对 |
| 模态与动态反馈 | 模态焦点进入/约束/恢复成立；成功/错误/异步状态可感知 | `renderer/modal.tsx`（`role=dialog`、`aria-modal`、`aria-label`、close aria-label）+ 焦点/状态断言 + 人工核对 |
| 语言切换与共享文案 | 切换后名称/错误/状态/焦点仍可感知；无仅依赖语言或颜色的隐含状态 | 双 Profile/locale 断言（`ui-bilingual.test.tsx`、`error-localization.test.tsx`、`localization.spec.ts`）+ 人工核对 |

- **N/A/延期触发**：每个未使用宿主能力记录理由、影响范围、重新纳入触发；暂缓项记录 owner、复核日期/触发。
- **失败严重度映射**：跨模块键盘不可达/焦点丢失/共享状态不可感知且影响标准模块 → `blocker/required` 候选；可隔离宿主缺口 → `major` gate 判定；局部体验 → `minor/info`。

## 9. 严重度量尺冻结（I-READINESS-006）

- **量尺版本**：采纳 VP-008 v0.10.0 §阻断与严重度量尺（本决策绑定该版本，不得在 S1 重写）。
- **等级**：`blocker`（required；阻断启动/构建/依赖闭包/Manifest/Schema/fail-closed、破坏认证授权/数据隔离/迁移完整性/协议边界、证据不可复现、全局 protocol-gap）→ 必须修复/用户书面 residual/`no-go`；`major`（依 gate 判定，影响跨模块必需能力或退出判据，涉生产边界 → required，否则 non-blocking+延期）；`minor`（non-blocking；局部质量）；`info`（non-blocking；观察项）。
- **适用 scope**：本 VP 冻结基架准入分母 + 未来标准业务模块共同门禁；领域特有项默认不进 required。
- **台账映射**：信息未知用 `I-READINESS-*`；实现缺陷/整改事实落 Goal ledger；VP 级阻断/意见响应留 Vision Review；单一 canonical finding id，其余引用。
- **S1 只应用规则**：S1 只能按本量尺分类，不得重写定义；S2 后未经用户书面裁决不得扩大/收缩/改写 `required` 定义。

## 10. 证据基线有效性（I-READINESS-007）

- **基线字段**：候选 commit（`852ee7e`）；构建 artifact digest（如有，登记）；依赖锁 digest（`apps/api/go.sum`、`apps/web/package-lock.json`）；Go/Node/npm/基础镜像版本（§1）；Profile 默认集（§3）；迁移台账版本（0001–0010）；`schema-ui-docs` pin（`ca9e5fe…`）；adapter/exclude disposition（§5）。
- **来源身份**：候选默认 clean checkout；若验证有意使用未提交输入，必须记录输入用途/scope、owned-path manifest 及 digest、patch/diff digest、未跟踪/生成文件清单及 digest、纳入/排除理由。
- **变更分类**：仅影响已冻结条目的代码/配置/工具链/锁/迁移/Profile/容器 → 标记受影响分母项重跑并更新 digest；改变分母范围/模块适用规则/协议 pin/disposition/风险语义 → 回流 S0 更新分母并经用户裁决；S5 后候选与裁决不一致 → 重跑受影响项，否则保持 `no-go`。

## 11. `go` 消费有效性（I-READINESS-009）

- **`go` 适用候选**：S5 证据矩阵指向的候选 commit/artifact、绑定的 patch/manifest/input digest、Profile/模块集合、协议 pin/disposition、升级与 fork/compose 基线、明确解锁 scope。
- **消费前 freshness review**：每个后续业务 VP 激活前必须完成；最低字段：候选 commit/artifact + digest 身份；冻结命令与关键证据可执行性；外部输入/证书/镜像/包源与验证环境可用性；最新 Goal/Vision open finding 与 accepted residual 投影。
- **失效触发**：源代码/配置/有意 patch 改变；依赖锁/工具链/基础镜像/artifact 改变；迁移台账/Profile 默认集/模块适用矩阵/容器或 fork 基线改变；`schema-ui-docs` pin/disposition/协议语义改变；认证授权/数据隔离/fail-closed/可访问性等共同门禁语义改变；Charter/VP scope 或退出判据改变。触发后受影响项完成重验证并留痕前，后续业务实现门闩关闭。
- **S5 记录字段**：`decision`、日期、证据矩阵链接、Goal/Vision finding 闭合状态、accepted residual、受影响/解锁 scope、适用候选、来源身份、`go_issued_at`、`last_freshness_review_at`、`next_freshness_review_trigger`、失效触发、roadmap 业务门闩生效语句。

## 12. 审计模式与 scope（I-READINESS-005 的 S0 段）

- **模式**：`cross`（self + independent）。
- **independent provider**：**grok build（模型 `grok-4.5`、思考强度 high）执行 `audit`**（[D-002](D-002-independent-audit-provider-grok-build.md)）。
- **覆盖 scope**：compatibility、data、migration、production/release 与跨边界治理语义；跨模块 UI 可访问性下限。
- **S5 审计**：由 grok 独立会话产出可核对 Goal 审计意见（`03-audit/A-NNN-*.md`，`source: independent`）；self 复核基线一致性。provider 不可用/超时/无可核对输出时独立门禁保持未满足。

## 13. 共享能力分母（I-READINESS-004）

从订单/钱包/类目/通知候选抽取的**框架级共性能力**（不预设领域模型）：list/detail 读表面、写操作（create/update/delete）、状态流转、权限（读/写/分配）、操作审计、迁移与 system-data reconcile、表单/校验/反馈、导航与 Manifest 发布、双语与设置。领域特有项默认不进本 VP required。

---

## 未选方案

- **不冻结分母直接进入 S1**：违反 VP-008 最小可枚举证据面与 S0 门禁；S1 只按冻结量尺登记。
- **候选基线含未提交输入**：当前树 clean，无需 patch/manifest/digest 绑定路径。
- **降低可访问性下限为纯人工**：VP-008 要求「可复跑断言 + 必要人工核对」；静态/运行时断言已存在于 vitest/e2e，故按 §8 冻结。

## 遗留与触发

- V-001~V-008 已在候选 `852ee7e`（clean）实测通过；V-009 由 CI 在 push/PR 时执行（与 V-001~V-008 本地重跑等价）。S1 扫描以本矩阵为基线。
- 本决策冻结后，若后续发现改变分母范围/协议/风险语义的证据，回流本决策按 P-004 重新裁决。
