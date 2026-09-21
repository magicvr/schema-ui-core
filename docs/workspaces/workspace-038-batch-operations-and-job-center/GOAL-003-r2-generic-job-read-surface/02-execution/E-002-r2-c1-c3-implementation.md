---
id: E-002-r2-c1-c3-implementation
doc: execution-entry
parent: GOAL-003-r2-generic-job-read-surface
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-002 · R2 C1～C3 实施（模块 / 查询与索引 / 读面 API）

## 事实（2026-09-19）

### 1. 用户 P-004 裁决（本轮决策来源）

| 决策点 | 裁决 |
|--------|------|
| `I-038-007` 索引（O-3） | 新增迁移 v72 `jobs_management_indexes`；**并细化为 `(created_at DESC, id DESC)`**（二次裁决，见 D-001 §1.2a） |
| 权限策略 | `jobs.read` 用 `PolicyAdmin`（写权限 `jobs.write` 归 R4，本目标不声明） |
| `I-038-008` 结果 URL | 泛化 + 共享 helper + 各模块自申 base path；wallet 输出字符串逐字不变；**登记为字节等价重构交审计复核** |
| 与 R1 O-3 的口径调和 | R2 D-001 补登记取代关系（R1 原文保留不改写） |

### 2. 落地产物

**C2 · 查询与索引**

| 产物 | 位置 |
|------|------|
| `ResultURL(basePath, id)` 共享派生；`ListFilter` + `ListJobs`（kind/status/actor/时间过滤 + 排序白名单 + COUNT 同 WHERE + `pagination.Offset`）；`GetJob`（管理作用域详情） | `internal/jobs/list.go` |
| 回归测试（过滤/分页/排序白名单/管理读面/ResultURL 逐字不变） | `internal/jobs/list_test.go` |
| 新迁移贡献 v72（`CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(created_at DESC, id DESC)`） | `modules/jobs/migration/migration.go` |
| 描述符测试（len 1→2、Version/Key/Name、非 tombstone、`ApplyPostgres == nil`） | `modules/jobs/migration/migration_test.go` |

**C1 · 模块与权限接线**

| 产物 | 位置 |
|------|------|
| `admin.jobs` Provider（Descriptor / Register / CompiledPersistence 空） | `modules/jobs/provider.go` |
| 页面 schema（`table.sort`；**不**声明 `actions.batch.request`，符合 O-2 冻结口径） | `modules/jobs/schema/jobs.json` + `schema.go` |
| Manifest fragment（jobs 页 + `menu_jobs`） | `modules/jobs/manifest/` |
| 进 admin 默认集（`profileDefaults`）+ `BuiltinModules` 描述符（Profile 内容扩展） | `kernel/profile.go` |
| 组合根装配 + **R-1 修复**（`jobRuntime.enabled` 改为 `admin.jobs` 或 `admin.wallet` 任一存在即启用） | `internal/composition/composition.go` |

**C3 · 读面 API 与作用域**

| 产物 | 位置 |
|------|------|
| `GET /api/jobs`（列表，共享 `resourceList` 信封）、`/{id}`（详情）、`/{id}/result`（409/410/终态/附件三段语义） | `internal/handler/jobs.go` |
| 管理作用域读面测试（跨 actor 可见、过滤、非法输入、401/403、详情/结果语义） | `internal/handler/jobs_test.go` |
| wallet 结果地址改经共享 helper（输出逐字不变） | `internal/handler/wallet.go` |

### 3. 冻结断言同步（修复 store 包 RED）

实施中出现一次**真实回归**：新增 v72 后 `internal/store` 变红（3 个测试 + 身份指纹）。根因是**只加了迁移描述符、未同步冻结断言**。已修复：

| 断言 | 变更 |
|------|------|
| `internal/store/identity.go` | `completeFingerprintCatalogHead` 71 → 72 |
| `internal/store/identity_test.go` | `lockedHeadExtraTables[72] = {}` |
| `internal/store/migrate_test.go` / `operations_test.go` / `restart_test.go` | applied 尾部 71 → 72（含 `jobs_management_indexes`） |
| `internal/store/migrate_test.go` 冻结列表 | 追加新三元组（`async_jobs`(42) 行**未动**） |

**教训**（供 R3/R4 复用）：新增迁移贡献必须一次同步 **6 处** 断言（冻结列表、身份指纹头、`lockedHeadExtraTables`、3 处 applied 尾部）。

### 4. 前端副作用（C1 连带）

| 产物 | 变更 |
|------|------|
| `denominator-render.test.tsx` | 页面分母 35 → 36 |
| `schema-keys.structural.test.ts` | `jobs/schema/jobs.json` 入 SCHEMA_FILES；`jobs/manifest/fragment.json` 入 manifest 键清单 |
| i18n `zh-CN` / `en-US` | 新增 `manifest.title.jobs`、`manifest.nav.jobs` 与 `schema.jobs.*` 共 22 键（两目录键集一致） |

### 5. 计数/快照断言同步

| 断言 | 变更 |
|------|------|
| `composition_test.go` | admin profile `wantPermissions` 34 → 35、`wantNavigation` 18 → 19 |
| `nav_group_r4_test.go` | `operations` 组三处补 `jobs` |
| `s5_manifest_snapshot_test.go` | admin 22→23、admin+digitaloffer 25→26、admin+telegram 24→25 |
| `testsupport/store.go` | 测试 env 系统数据补 `jobs.read` 权限贡献 |

### 6. 验证证据

| 验证 | 结果 |
|------|------|
| `go build ./...` | exit 0 |
| `go test ./...`（apps/api） | **全绿** |
| `vitest`（apps/web） | **114 files / 1440 tests 全绿** |
| `npm run typecheck`（tsc -b + e2e tsconfig） | exit 0 |
| capability-declaration.guard + all-module-schemas-dval | 75/75 通过（新页面合规） |

**已知 flake 留痕**：全量 Go 首跑曾命中 `TestShutdownDrainHarnessPostgres` 一次失败；隔离运行与整包复跑均绿 —— 既有 VP-021 PG drain harness 并行 flake（workspace-009 `A-003` 已留痕），与本次改动无关。

### 7. Git checkpoints

| hash | 内容 |
|------|------|
| `d8532c68` | R2 查询/索引/wallet 结果 URL 泛化 + 冻结断言修复 |
| `0624b808` | 索引细化为 `(created_at DESC, id DESC)` |
| `c456cbe2` | C1 admin.jobs 模块建立与接线 |
| `aa21c411` | 前端分母与 i18n 接线 |
| `bd471b3b` | C3 管理作用域读面测试 |

### 8. 边界

- **未**改动任何 pinned 协议工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）。
- **未**放宽 `GetForActor` actor 隔离语义（wallet 冻结测试全绿）。
- **未**重开 VP-012：对 `wallet.go` 的改动为**字节等价重构**（`resultUrl` 输出字符串不变），已登记交 C4 审计复核。
- **未**声明 `actions.batch.request`（O-2 口径）。
- `jobs.write` 未声明（写操作/结果中心归 R4）。

### 9. 未做（移交 C4）

- **未**执行审计：C4 待跑（self + grok build independent）。
- **未**投影 Root R2 检查点。
