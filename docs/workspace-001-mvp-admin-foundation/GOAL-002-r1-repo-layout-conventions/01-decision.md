---
id: GOAL-002-r1-repo-layout-conventions
doc: decision
status: done
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.1
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

本目标无独立开放 required 信息项。布局权威决策见父目标 [GOAL-001 `D-004`](../GOAL-001-mvp-admin-foundation/01-decision.md)。

## D-001 · 承接 Root D-004，本目标只落约定不写业务

**日期**：2026-07-31  
**状态**：**superseded**（部分条款由 D-002 修正）

**决定**（历史）：

1. 在本目标内把 monorepo 路径、包管理、边界写入可维护文档（优先扩展/新建 architecture 约定，并在根 README 给最短入口）。
2. ~~可创建空的 `apps/web`、`apps/api` 占位~~ → 见 **D-002**（创建权收紧）。
3. 可运行骨架分别由 GOAL-003 / GOAL-004 交付。

**为什么**：

- 分开「约定真相」与「可运行骨架」便于并行与验收。
- 避免在约定未落盘时两边 scaffold 路径漂移。

**未选方案**：

- **约定与双端 scaffold 挤进同一目标**：文档少，但失败面大、难并行关门。

## D-002 · 闭合 A-001：`apps/*` 创建权 + 运行入口成功标准（F-001 / F-002）

**日期**：2026-07-31  
**状态**：accepted  
**响应**：independent A-001 · F-001、F-002  
**用户意图**：`/govern` 明确要求闭合创建权与运行入口成功标准后推进 R1

**决定**：

### F-001 · `apps/*` 首次创建权 → 方案 **(A)**

1. **GOAL-002 只写约定文档**（`docs/architecture/`、根 README 等），**禁止**在本目标内首次创建可运行的 `apps/web` / `apps/api` 工程树（含 `go.mod`、`package.json`、脚手架源码）。
2. **`apps/api` 首次实质建树**由 **GOAL-003** 负责；**`apps/web` 首次实质建树**由 **GOAL-004** 负责。
3. 仓库中若已存在空目录壳（无 manifest / 无源码），**不视为** 002 交付物，也**不**改变路径语义；003/004 可原地填充，**不得**改 D-004 已锁路径（`apps/api`、`apps/web`）。
4. 003/004 决策须服从本交界：若目录已存在则原地填充，不改路径。

### F-002 · 「本地运行入口」成功标准两档

1. **必达（002）**：布局 / 包管理 / 边界 / 非业务默认树 文档；期望命令**名称契约**（不要求 002 能独立执行出服务）。
2. **运行命令可执行性**：**owned-by GOAL-003 / GOAL-004**。002 只写期望契约并链到姊妹 README；不以「命令已可跑」作为 002 独有关门门禁。
3. 期望契约（骨架阶段，业务未实现）：
   - API：`cd apps/api && make run` 或 `go run ./cmd/server`；探活 `/healthz`；默认端口建议 `:8080`
   - Web：`cd apps/web && npm install && npm run dev`；构建 `npm run build`

**为什么**：

- 消除并行双重 init（空 README vs `go mod`/`npm create`）与约定目标独吞可运行验收的风险（A-001）。
- 与 Root D-004「002 约定、003/004 骨架」切分一致。

**未选方案**：

- **(B) 002 仅 `.gitkeep`/一行 README 占位**：可接受，但本轮选 A 更清晰；空壳若已存在由 003/004 填充即可，无需 002 再写占位交付。
- **保持「命令级可找到且可执行」为 002 硬标准**：会与「002 不交付服务」矛盾，易假性阻塞或编造命令。

**影响**：

- 成功标准见 `00-meta` v0.2.0；F-001/F-002 → `fixed`。
- 放行 002 文档实施与 003/004 建树；不改变 monorepo 路径。

**后续**：

1. 落盘 monorepo 约定 + 根 README。
2. 003/004 按 D-002 交界并行 scaffold。
