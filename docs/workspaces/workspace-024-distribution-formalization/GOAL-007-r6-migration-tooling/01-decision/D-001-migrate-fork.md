---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-007-r6-migration-tooling
version: 0.1.0
---

# D-001 · migrate-fork 工具设计定档（2026-08-29）

## 决策

1. **形态**：`schema-ui migrate-fork [--dir <path>] [--dry-run]` 子命令（Go CLI，apps/api/cmd/schema-ui 内新增）。
2. **非破坏原则**：除两处低侵入改写外零修改——① `go.mod` 的 `require apps/api` 版本行（bump 至最新已发布 tag）② `web/.npmrc` 的 GH Packages 映射替换为 npmjs 钉死（**先备份 `.npmrc.migrate.bak`**）。用户代码（cmd/server/main.go 等）一律**检查 + 引导**，不覆盖。
3. **类型判定（指南 §1 程序化）**：
   - A · 纯装配型：`apps/api` 源码存在 + `cmd/server/main.go` 无 kernel 覆盖特征（无 `kernel.` / 渲染器 override 导入）→ 提示直接迁移
   - C · 深度定制：main 含 kernel 契约面 import（`kernel.`）或渲染器主路径符号 → 建议**保持 fork**
   - 其余（有覆盖但非 kernel）→ B
4. **步骤清单（dry-run 输出 · 实跑执行前两项）**：
   1. go.mod require 版本检查（< 最新已发布 → 报告/bump）
   2. web/.npmrc 检查（含 `@magicvr` 映射 → 报告/改写为 `@magicvr:registry=https://registry.npmjs.org`）
   3. main.go 检查（薄封装 vs 旧组合根 → 引导：薄封装 OK；旧组合根提示 `schema-ui create` 重建 + 业务迁入）
   4. 验证引导（go build + healthz 探针命令输出）
5. **最新已发布版本来源**：`go list -m -json github.com/magicvr/schema-ui-core/apps/api@latest`（registry 语义 · 免 .env token；github proxy 版——主仓 tag 已 push；migrate 目标仓 go get 会用 proxy ✓）。注意：`apps/api` 消费 = `github.com/magicvr/schema-ui-core/apps/api`（实际 module 路径）。
6. **验收**：golden-field `9510023`（v0.3.0 旧态）worktree 实测——dry-run 报 A 型 + 步骤清单；实跑后 go.mod bump ✓ + .npmrc 钉死（备份）✓ + main 引导 ✓；`go build ./cmd/server` 绿。

## 7. 校准注记（2026-08-29 · A-002 F-001 响应）

- **类型判定终版**：A = `server.Serve(` 薄封装或尚未依赖包面；B = 旧组合根（`assembly.OpenStore` 标准组合路径存在 · 无 kernel 覆盖——组合根模板标准件 `kernel.JoinIdentifiers` 等**不算**覆盖）；C = 手搓 kernel 面且未走组合路径。
- **9510023 实测定型 = B**（dry-run 原 §6 写作「报 A 型」为执行偏差，以本注记为准）；A/C 由独立审计夹具复现。
- **F-002 补丁**（同日）：`require ( … )` 块形态解析 + 实跑无条件执行 `go get @latest`（含原未依赖场景）。

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 全自动迁移（含 main.go 重写） | 引导替代覆盖 | 用户代码所有权 + fork 定制保留原则 |
| 独立脚本（非 CLI 子命令） | CLI 子命令 | 与 create/add/upgrade/serve 同族，发现性一致 |
| 支持 C 型强制迁移 | 建议保持 fork | 指南 §1 C 类结论 + Charter fork 并存 |