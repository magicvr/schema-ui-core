---
id: GOAL-040-w28-admin-passwd-convention
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

## D-001 · ADMIN_PASSWD 声明约定 + TEST_ADMIN 机制退役

### 触发

管理员手动测试后，登录强制改密（`must_change_password`）使 `admin` 现密码不可知。`ADMIN_INITIAL_PASSWORD` 只对零用户库的 fresh bootstrap 生效；自动化测试与 AI 助手连**现有库**时不知道 admin 现密码而撞墙。既有 TEST_ADMIN 测试账户机制（[workspace-030 E-011](../../../workspace-030-telegram-channel-runtime/GOAL-001-telegram-channel-runtime/02-execution/E-011-operator-config-and-test-admin.md) 引入）不可发现、AI 助手与自动化测试通常不知道怎么用，且 boot-time 每次启动 upsert 一个可用 admin 形如后门。

用户 2026-09-06 确认三项裁决（I-001/I-002/I-003）：权威位置 = `apps/api/configs/.env`；本轮包含 TEST_ADMIN 移除；`ADMIN_PASSWD` 为纯声明约定。

### 决定

1. **`ADMIN_PASSWD` 约定（新增）**：维护者在本地 gitignored 的 `apps/api/configs/.env` 设置 `ADMIN_PASSWD`，声明**当前库 `admin` 的现密码**。语义：
   - 供**自动化测试**与 **AI 助手**在连接现有库时读取并登录 `admin`，避免撞墙；
   - **API 不读取、不据此重置密码**（不是 seed/reset 通道，不改变 `must_change_password` 语义）；
   - 只写本地 gitignored 文件（`apps/api/configs/.env`），canonical 模板 `apps/api/configs/.env.example` 登记占位说明；
   - 若消费方遇到 `must_change_password=1`（admin 尚未改密），按真实首登改密流程处理（smoke 已有 `SMOKE_PASSWORD_NEW` 通道）。
2. **AI 助手可发现性**：在仓库根 `AGENTS.md`（AI 强制规则，助手每次会话加载）新增「本地开发环境」小节/快速链接，说明连现有库时从 `.env` / 环境变量的 `ADMIN_PASSWD` 读取 admin 密码；`QUICKSTART.md` 与 `apps/api/README.md` 同步登记。
3. **TEST_ADMIN 机制退役（删除）**：移除 `TEST_ADMIN_USERNAME/TEST_ADMIN_PASSWORD` 全部代码与文档面：
   - `apps/api/internal/config/config.go`（字段 + env 读取）
   - `apps/api/modules/authsession/systemdata/bootstrap.go`（`EnsureTestAdmin`）
   - `apps/api/internal/composition/composition.go`（openStore upsert 分支）
   - `apps/api/modules/authsession/systemdata/reconcile_test.go`（`TestEnsureTestAdmin`）
   - `apps/api/configs/.env.example`（TEST_ADMIN 块 → 由 ADMIN_PASSWD 说明替代）
4. **自动化测试回退**：`scripts/smoke.sh` 允许 `SMOKE_PASSWORD` 缺省时回退到 `ADMIN_PASSWD`（环境变量已含该值时），连现有库的 smoke 不再强制手工传密码。

### 为什么

- **`ADMIN_PASSWD` 是声明而非重置**：TEST_ADMIN 之所以失败，是因为它自造了一个可用凭据（后门形状）且不可发现；声明「当前 admin 密码」让测试/AI 走真实登录路径，语义清晰、零后门、不触碰 `must_change_password` 门禁（W16-F01 设计意图）。
- **权威位置选 `apps/api/configs/.env`**：与 E-011 既有测试凭据约定一致；`configs/.env` 是 config.Load 读取的权威 env 文件；canonical 模板 `.env.example` 登记后对人与 AI 都可发现。
- **AI 可发现性落 AGENTS.md**：AGENTS.md 是助手每次会话必读的强制规则，是"AI 助手知道能从 .env/环境变量拿到密码"的唯一可靠落点。
- **TEST_ADMIN 移除**：机制不好用（AI/自动测试不知道怎么用它）、不可发现、且 boot-time upsert 每次重置密码与 `must_change_password` 门禁张力大；退役后由 `ADMIN_PASSWD` 约定取代。

### 未选方案

- **`ADMIN_PASSWD` 由 API 读取并用于重置/同步 admin 密码**：否定。等于再造一个 TEST_ADMIN 式后门，且与 `must_change_password`（W16-F01）门禁冲突；只做声明、消费方自行登录。
- **保留 TEST_ADMIN 并继续推广**：否定。机制本身不可发现、AI/自动测试不会用；继续保留只会增加攻击面与心智负担。
- **ADMIN_PASSWD 权威位置放仓库根 `.env`**：否定（未选）。仓库根 `.env` 是 compose 插值来源，与 API `configs/.env` 语义不同；用户裁决以 `apps/api/configs/.env` 为准，便于与既有测试凭据约定合并管理。

### 影响

- 代码：`config.go`、`systemdata/bootstrap.go`、`composition.go`、`reconcile_test.go`、`.env.example`、`smoke.sh`。
- 文档：`AGENTS.md`、`QUICKSTART.md`、`apps/api/README.md`。
- 历史：workspace-030 E-011 五件套不改写；本目标 00-meta 备注留痕退役事实。

### 后续

S1 冻结 → S2 实施（代码 + 文档）→ S3 回归验证（`go test`）→ S4 关门自审（A-001）。
