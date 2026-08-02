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
| [docs/workspace-001-mvp-admin-foundation/goal-tree.md](docs/workspace-001-mvp-admin-foundation/goal-tree.md) | 当前工作区目标树 |
| [AGENTS.md](AGENTS.md) | AI 协作强制规则 |

## 仓库布局（摘要）

```text
apps/web   # React 前端（npm）
apps/api   # Go 后端（Go modules）
docs/      # 愿景 / 工作区 / architecture
skills/    # 治理 Skills 包
```

完整约定与边界见 [monorepo-layout.md](docs/architecture/monorepo-layout.md)。

## 本地运行（骨架阶段）

> 业务未实现；仅工程可启动。命令契约详情见 monorepo 约定与各 app README。

### API · `apps/api`（GOAL-003）

```bash
cd apps/api
make run
# 或：go run ./cmd/server
# 探活：curl http://localhost:8080/healthz
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
> `docs/workspace-002-production-admin-foundation/GOAL-008-r5-engineering-fork/attachments/I-008-001-engineering-contract.md`。

### Docker Compose 一键启动（第二启动路径）

```bash
# 需先提供生产必填密钥（fail-closed）；可写入仓库根 .env（gitignored）或 export
AUTH_JWT_SECRET=<强随机串>
ADMIN_INITIAL_PASSWORD=<初始 admin 密码>

docker compose up --build
#  API: http://localhost:8080  (GET /healthz 探活)
#  Web: http://localhost:8081  (nginx 服务 SPA + /api 反代；同源免 CORS)
#  登录种子 admin → 后台首页
```

- 本地开发仍为默认双进程路径（见上文 API / Web 段）；fork 使用者可选本地双进程或 Compose。
- `docker compose down` / 重启后 SQLite 数据由命名卷 `db-data` 保持。
- 将密钥写入仓库根 `.env`（gitignored）可避免新 shell 里 `docker compose config` / `down` 因 fail-closed 插值重复 export。
- 完整生产运维 / CI-CD 部署流水线、TLS、多实例为**非目标**。

## 状态说明

- R2 MVP 协议覆盖子集已按 Root `I-PROTO-001` v0.1.3 冻结；这不是「支持全部协议功能」或 R3-R5 已实现的声明。
- R1 目标：可运行前后端骨架 + 布局约定；账号权限属 R4；Admin 外壳属 R3。
