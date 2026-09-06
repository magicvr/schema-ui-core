---
id: GOAL-040-w28-admin-passwd-convention
doc: execution-entry
record_id: E-001
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 1.0.0
---

# E-001 · 决策冻结 + ADMIN_PASSWD 约定登记 + TEST_ADMIN 移除实施

## 2026-09-06 · 决策冻结与实施

### 已发生事实

1. **S1 决策冻结（D-001）**：用户三项裁决——权威位置 `apps/api/configs/.env`（I-001）、本轮包含 TEST_ADMIN 移除（I-002）、`ADMIN_PASSWD` 为纯声明约定、API 不读取不重置（I-003）。I-001～I-003 全部 verified。
2. **TEST_ADMIN 机制移除**（代码面）：
   - `apps/api/internal/config/config.go`：删除 `TestAdminUsername`/`TestAdminPassword` 字段与注释块，删除 `TEST_ADMIN_USERNAME`/`TEST_ADMIN_PASSWORD` 两行 env 读取。
   - `apps/api/modules/authsession/systemdata/bootstrap.go`：删除 `EnsureTestAdmin` 函数（含 upsert/重置/角色链接分支）。
   - `apps/api/internal/composition/composition.go` `openStore`：删除 TEST_ADMIN upsert 分支；`auth`/`strings` 导入仍被其它代码使用，无死导入。
   - `apps/api/modules/authsession/systemdata/reconcile_test.go`：删除 `TestEnsureTestAdmin`（两个子测试）；移除不再使用的 `slices` 导入。
   - `apps/api/configs/.env.example`：删除 TEST_ADMIN 块，替换为 `ADMIN_PASSWD` 声明约定说明。
3. **ADMIN_PASSWD 约定登记**（文档面）：
   - `apps/api/configs/.env.example`：新增 `# ADMIN_PASSWD=` 占位与语义说明（声明当前库 admin 现密码；API 不读取不重置；只写本地 gitignored 文件）。
   - `apps/api/README.md`：配置键表新增 `ADMIN_PASSWD` 行（约定键说明）。
   - `QUICKSTART.md`：新增「连已有库的 admin 凭据」条目。
   - 仓库根 `AGENTS.md`：新增「本地开发环境与 admin 凭据约定」小节——AI 助手与自动化测试连现有库时从 `.env` / 环境变量读取 `ADMIN_PASSWD` 并以 admin 登录；秘密不得输出/提交。这是 AI 助手可发现性的权威落点。
4. **smoke 回退**：`scripts/smoke.sh` `SMOKE_PASSWORD` 缺省时回退 `ADMIN_PASSWD`（环境变量）；头部输入说明同步。`pre-release-smoke.sh` 显式传 `SMOKE_PASSWORD`（disposable 新库路径），不受回退影响。
5. **canonical 模板守卫**：`apps/api/internal/config/env_example_test.go` `TestCanonicalEnvExample` 新增 `declarationOnlyKeys` 白名单（`ADMIN_PASSWD`）——该键是 API 刻意不读取的约定键，属合法例外并注释说明。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 构建干净 | `go build ./...`（apps/api）exit 0 |
| 静态检查干净 | `go vet ./internal/config/... ./internal/composition/... ./modules/authsession/...` exit 0 |
| 受影响包测试绿 | `go test ./internal/config/... ./modules/authsession/systemdata/... ./internal/composition/...` 全 ok |
| 全量回归绿 | `go test ./...`（apps/api 全包）exit 0，0 FAIL |
| 移除面无残留（代码） | `rg -n "TEST_ADMIN|EnsureTestAdmin" --glob '*.go'` 仅剩无关测试函数名（TestAdminDisable*/TestAdminUnlock*/TestAdminPrefill*） |
| 约定落点 | `apps/api/configs/.env.example`、`apps/api/README.md`、`QUICKSTART.md`、`AGENTS.md` |

### 评估

S2/S3 完成：4 个成功检查点全部达成（约定登记 ✓ / AI 可发现性 ✓ / TEST_ADMIN 移除 + 测试绿 ✓ / smoke 回退 ✓）；关门自审 A-001 落盘（pass，0 required），meta progress 更新为 4/4，等待用户确认置 done。
