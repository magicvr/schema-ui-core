---
id: GOAL-003-r1-api-go-scaffold
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.0
---

# 决策记录 · GOAL-003

## 信息需求与阶段门禁

| ID | 级别 | 最晚阶段 | 状态 | 阻断 |
|----|------|----------|------|------|
| I-003-001 | non-blocking | 骨架可运行前 | **verified** | 不阻断；README / go.mod 声明 `go 1.26` |
| I-003-002 | **required** | 首次 go.mod 前 | **verified**（D-002） | 已确认 module path；可写 go.mod |

布局与复用策略服从父目标 [D-004](../GOAL-001-mvp-admin-foundation/01-decision.md)；目录交界服从 [GOAL-002 D-002](../GOAL-002-r1-repo-layout-conventions/01-decision.md)。

## D-001 · 骨架范围与平行仓复用边界

**日期**：2026-07-31  
**状态**：accepted

**决定**：

1. 目标路径：`apps/api/`（Go module 根）。
2. 目录取向：`cmd/server`、`internal/`（config/server/handler 等最小集）、`pkg/`（可复用 envelope/version 等，按需）。
3. 复用：参考 `../allinme.core-api`（本地平行，`dev`）的分层、Makefile、`.env.example`、health 模式；**移植时改名/改 module**，去掉业务域。
4. **明确不搬**：`internal/domain` 中 order/wallet/notification、对应 handler/repository、demo 业务 seed、page schema 业务页（属后续阶段或非目标）。
5. 鉴权/JWT/SQLite 可作为**后续 R4 候选模式**记在备注，R1 不强制完整 auth 闭环；**禁止** R1 默认挂业务鉴权中间件为必选。

**为什么**：

- Charter 非目标排除特定业务终端模块；平行仓已含演示业务，整树拷贝会污染 MVP 边界。
- R1 只要可运行基架，账号权限在 R4。

**未选方案**：

- **整仓拷贝再删业务**：快但易残留协议 2.4 声明与业务路由。
- **R1 直接上完整 auth+RBAC**：超出 R1「不实现业务能力」边界，且 `I-PROTO-002` 未就绪。

## D-002 · 闭合 A-001 F-001：module path 门禁 + Makefile 硬度

**日期**：2026-07-31  
**状态**：accepted  
**响应**：independent A-001 · F-001（required）；顺带采纳 F-002 recommended 的 `run` 必达  
**用户意图**：`/govern` 明确要求闭合 module path 门禁后推进 R1

**决定**：

1. **I-003-002** 级别由 `non-blocking` **升为 `required`**；最晚阶段 = **首次写入 `go.mod` 前**（与门禁一致）。
2. **module path 定稿**：`github.com/magicvr/schema-ui-core/apps/api`  
   - 与远程 `origin` = `git@github.com:magicvr/schema-ui-core.git` 对齐  
   - **禁止**使用 `github.com/magicvr/allinme.core-api` 或其它平行仓 path
3. **I-003-002 → verified**（本决策 + 后续 `go.mod` 证据）。
4. **I-003-001**：本机探测 `go1.26.0`；R1 在 `go.mod` 写 `go 1.26` 并在 README 声明实测版本 → **verified**。
5. **Makefile**：R1 必达至少 **`run`**（或等价文档命令）+ health 验证；`build` / `test` 推荐一并提供（响应 A-001 F-002 recommended）。
6. **目录交界**（A-001 F-003）：若 `apps/api` 已有空目录壳，**原地填充**，不改路径语义（服从 GOAL-002 D-002）。

**为什么**：

- 无 module path 无法合法 init；标 non-blocking 与「首次 go.mod 前」自相矛盾（A-001）。
- 仅 `test`/`build` 无 `run` 削弱「可本地运行」意图。

**未选方案**：

- **保持 I-003-002 non-blocking 实施时自拟 path**：易 rename 债，用户已要求闭合门禁。
- **module path 用短名 `schema-ui-core/api`**：与 GitHub 远程惯例不一致，后续 replace 成本高。

**影响**：

- 放行 `apps/api` 骨架实施；F-001 → `fixed`。
- 不放行业务 API / auth / 协议兼容主张。
