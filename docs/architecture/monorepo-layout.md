---
title: Monorepo 应用布局与包管理约定
status: active
created: 2026-07-31
updated: 2026-07-31
parent: null
version: 0.1.0
---

# Monorepo 应用布局与包管理约定

> **权威来源**：工作区 `workspace-001-mvp-admin-foundation` · Root `D-004` · GOAL-002 `D-002`。  
> 本文件是 **应用代码树** 的布局真相；治理目录见 [directory-layout.md](directory-layout.md)。

## 1. 顶层分区

```text
schema-ui-core/
├── AGENTS.md                 # AI / 治理规则入口
├── README.md                 # 最短运行与文档入口
├── apps/
│   ├── api/                  # Go 后端（GOAL-003 首次建树）
│   └── web/                  # React 前端（GOAL-004 首次建树）
├── docs/                     # 愿景 / 工作区 / architecture / contracts
└── skills/                   # 治理 Skills 分发包与镜像
```

| 路径 | 职责 | 包管理 |
|------|------|--------|
| `apps/web` | Admin 前端可运行树 | **npm** + `package-lock.json` |
| `apps/api` | Admin 后端可运行树 | **Go modules**（独立 `go.mod`） |
| `docs/` | 愿景、工作区目标、方法论、契约 | 文档；非应用 runtime |
| `skills/` | 消费适配器 / install | 非业务应用 |

**边界**：

- **禁止**把订单 / 钱包 / 通知等业务域目录当作本仓 MVP **默认**树。
- **禁止**在治理文档中把 `docs/vision/` 当作 goal-tree 或 progress 权威。
- 平行仓（`allinme.core-api` / `allinme.web-client`）仅结构参考，**禁止整仓拷贝**。

## 2. `apps/*` 创建权（GOAL-002 D-002）

| 动作 | 负责目标 |
|------|----------|
| 布局 / 包管理 / 边界**文档** | GOAL-002 |
| `apps/api` **首次实质建树**（`go.mod`、源码、Makefile） | GOAL-003 |
| `apps/web` **首次实质建树**（`package.json`、源码、Vite） | GOAL-004 |

GOAL-002 **不**交付可运行服务；空目录壳（若有）不改变路径语义，由 003/004 原地填充。

## 3. 后端约定（`apps/api`）

| 项 | 约定 |
|----|------|
| Module path | `github.com/magicvr/schema-ui-core/apps/api` |
| Go 版本 | `go 1.26`（R1 本机实测；README 声明） |
| 入口 | `cmd/server` |
| 分层 | `internal/`（config / server / handler…）、`pkg/`（可复用小库） |
| 默认端口 | `:25080`（`HTTP_ADDR`） |
| 探活 | `GET /healthz` |
| 期望命令 | `make run` 或 `go run ./cmd/server`；另推荐 `make build` / `make test` |

R1 **不**默认挂业务鉴权中间件；账号权限属 R4。

## 4. 前端约定（`apps/web`）

| 项 | 约定 |
|----|------|
| 工具链 | Vite + React 19 + TypeScript |
| 包管理 | npm + `package-lock.json`（工作目录 `apps/web`） |
| UI | Tailwind CSS + shadcn/ui（须有 `components.json` 或文档记载的 init 痕迹 + `components/ui`） |
| 主题 | 浅/深色**最小占位**（完整外壳属 R3） |
| 分层（R1 预建空壳） | `src/app/`、`src/host/`、`src/protocol/`、`src/renderer/`、`src/components/ui/` |
| 期望命令 | `npm install`；`npm run dev`；`npm run build` |

R1 **不含** App manifest 导航壳、多业务路由（属 R3）；**不含**协议 Renderer 全量（属 R5）。

## 5. 本地运行入口契约（名称级）

> **可执行性 owned-by GOAL-003 / GOAL-004**。本节约定命令名称；细节以各 app README 为准。

### API（GOAL-003）

```bash
cd apps/api
# 可选：copy configs/.env.example → configs/.env
make run
# 或：go run ./cmd/server
# 探活：curl http://localhost:25080/healthz
```

### Web（GOAL-004）

```bash
cd apps/web
npm install
npm run dev
# 构建：npm run build
```

骨架阶段：**业务能力未实现**；仅验证工程可启动。

## 6. 与治理目录的关系

- 目标状态只在 `docs/workspaces/workspace-*/`；应用代码不建立第二套 goal-tree。
- 改 `docs/architecture` 白名单文件后，按仓库 stage 规则同步 `skills/core` 镜像（若脚本可用）。

## 7. Git 忽略策略

| 层级 | 路径 | 作用 |
|------|------|------|
| 仓库根 | [`.gitignore`](../../.gitignore) | OS/IDE、env 密钥、Node/Go 构建产物、临时目录 |
| API | [`apps/api/.gitignore`](../../apps/api/.gitignore) | `bin/`、本地 `.env`、测试二进制 |
| Web | [`apps/web/.gitignore`](../../apps/web/.gitignore) | `node_modules/`、`dist/`、`.vite/` |

**应提交**：源码、`go.mod`、`package.json`、**`package-lock.json`**（npm 锁定，可复现安装）、`.env.example`、文档。  
**勿提交**：`node_modules/`、`dist/`、`bin/`、`*.exe`、真实 `.env`、本地 IDE 配置。
