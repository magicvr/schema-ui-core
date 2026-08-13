# schema-ui-core

面向中型项目、可 fork 的 **Admin 基架**（React + Go），协议兼容边界对齐 [`schema-ui-docs`](https://github.com/magicvr/schema-ui-docs) `v2.7.0`。

本仓库同时承载 **Goal Governance** 核心文档与 Skills 分发（`docs/`、`skills/`）。

## 文档入口

| 文档 | 说明 |
|------|------|
| [docs/README.md](docs/README.md) | 文档体系 |
| [docs/architecture/monorepo-layout.md](docs/architecture/monorepo-layout.md) | **Monorepo 布局与包管理约定（R1）** |
| [docs/architecture/directory-layout.md](docs/architecture/directory-layout.md) | 治理目录布局 |
| [docs/vision/charter.md](docs/vision/charter.md) | 现行愿景 Charter |
| [docs/workspaces/workspace-008-admin-module-readiness/goal-tree.md](docs/workspaces/workspace-008-admin-module-readiness/goal-tree.md) | 当前准入与基架收敛工作区（VP-008）目标树 |
| [docs/workspaces/workspace-003-modular-admin-architecture/goal-tree.md](docs/workspaces/workspace-003-modular-admin-architecture/goal-tree.md) | 模块化架构工作区目标树（历史） |
| [AGENTS.md](AGENTS.md) | AI 协作强制规则 |

## 仓库布局（摘要）

```text
apps/web   # React 前端（npm）
apps/api   # Go 后端（Go modules）
docs/      # 愿景 / 工作区 / architecture
skills/    # 治理 Skills 包
```

完整约定与边界见 [monorepo-layout.md](docs/architecture/monorepo-layout.md)。

## 本地运行

命令契约详情见 monorepo 约定与各 app README。API 不会自动读取 `.env` 文件；本地进程
必须通过 shell 环境传入配置，Compose 才会读取仓库根 `.env`。

### API · `apps/api`（GOAL-003）

```bash
cd apps/api
export APP_ENV=development   # 必须显式设置；未设置时启动 fail-closed（生产无公开弱默认）
export APP_PROFILE=mvp       # 或 admin / demo；custom 必须同时设置 APP_MODULES_ENABLED
make run
# 或：go run ./cmd/server
# 探活：curl http://localhost:25080/healthz
# 就绪：curl http://localhost:25080/readyz
```

详见 [apps/api/README.md](apps/api/README.md)。

### Web · `apps/web`（GOAL-004）

```bash
cd apps/web
npm install
npm run dev
# 构建：npm run build
```

详见 [apps/web/README.md](apps/web/README.md)。

## 工程化与一键启动（R5 · GOAL-008）

> 生产级工程化交付（环境/配置、容器一键启动、健康检查、dev/prod 区分）随 R5 推进。契约见
> `docs/workspaces/workspace-002-production-admin-foundation/GOAL-008-r5-engineering-fork/attachments/I-008-001-engineering-contract.md`。

### Docker Compose 一键启动（第二启动路径）

```bash
# 需先提供生产必填密钥（fail-closed）；可写入仓库根 .env（gitignored）或 export
AUTH_JWT_SECRET=<强随机串>
ADMIN_INITIAL_PASSWORD=<初始 admin 密码>
APP_PROFILE=mvp                 # 或 admin / demo
APP_MODULES_ENABLED=            # 可选，逗号分隔的显式模块覆盖

docker compose up --build
#  API: http://localhost:25080  (GET /healthz 探活)
#  Web: http://localhost:25081  (nginx 服务 SPA + /api 反代；同源免 CORS)
#  登录种子 admin → 后台首页
```

- 本地开发仍为默认双进程路径（见上文 API / Web 段）；fork 使用者可选本地双进程或 Compose。
- `docker compose down` / 重启后 SQLite 数据由命名卷 `db-data` 保持。
- 将密钥写入仓库根 `.env`（gitignored）可避免新 shell 里 `docker compose config` / `down` 因 fail-closed 插值重复 export。
- `APP_PROFILE` 默认为 `mvp`（users/roles/account/notifications + Dashboard 首页）；选择 `admin` 会在同一 Web build 上增加 Settings/Activity/Data-Transfer；选择 `demo`（**非生产向**）会额外启用 `dev.examples`，在同一 build 上展示 8 个协议范例页 + Examples 导航（home 指向 `overview`）。
- `APP_MODULES_ENABLED` 非空时覆盖 Profile 默认模块集合；`custom` Profile 没有显式模块时 fail-closed。
- 完整生产运维 / CI-CD 部署流水线、TLS、多实例为**非目标**。

## 模块化 Admin 架构与 Profile（R4 起）

后端以**薄内核 + 模块 Provider + 启动时 Profile** 组装（workspace-003 模块化架构）：

- **模块 Provider**：一方标准 Admin 模块（`admin.users` / `admin.roles` /
  `admin.settings` / `admin.activity` / `admin.dashboard` / `admin.account` /
  `admin.notifications` / `admin.data-transfer`）以 `kernel.Provider` 结构化贡献
  HTTP、Schema、授权、Navigation、Manifest 与 compiled-global Persistence；
  composition 消费 finalize，冲突 fail-closed。
- **Profile**：`mvp`（core + users/roles + dashboard/account/notifications，
  home = `dashboard`）与 `admin`（+ settings/activity/data-transfer）为编译候选集；
  `demo`（W2，非生产向）= mvp + `dev.examples`（范例页演示面，home = `overview`）；
  `APP_MODULES_ENABLED` 显式覆盖。**同一 Web 构建**随 Profile 切换页面集，无需改前端。
- **数据**：迁移账本 `0001`-`0017` 全局唯一；fresh 与 versioned reconcile 分离；
  operationlog best-effort；`/api/records` 已退场（`0006` historical-only）。
- **探测**：`/healthz`（liveness）与 `/readyz`（store ping + 模块图 Start/Ready
  readiness，R5）。

fork 起点：选 Profile + `APP_MODULES_ENABLED` + 模块贡献接入业务，不修改 Renderer/Shell
主路径。

## 状态说明

- R2 MVP 协议覆盖子集已按 Root `I-PROTO-001` v0.1.3 冻结；这不是「支持全部协议功能」或 R3-R5 已实现的声明。
- R1 目标：可运行前后端骨架 + 布局约定；Admin 外壳属 R3（历史路线图编号）。
- R2 一等公民波次（[workspace-011](docs/workspaces/workspace-011-admin-functional-modules/goal-tree.md)）：`admin.dashboard` / `admin.account` / `admin.notifications` / `admin.data-transfer` 四个一方标准 Admin 模块已交付；订单/钱包等业务域降档至 R3（S-01～S-14）与 R4（B-01～B-11）。
