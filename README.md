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

## 状态说明

- R2 MVP 协议覆盖子集已按 Root `I-PROTO-001` v0.1.3 冻结；这不是「支持全部协议功能」或 R3-R5 已实现的声明。
- R1 目标：可运行前后端骨架 + 布局约定；账号权限属 R4；Admin 外壳属 R3。
