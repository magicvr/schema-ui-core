---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-cli
version: 0.1.0
---

# D-001 · CLI 形态定案（用户裁决）

## 裁决（2026-08-29 用户）

| 项 | 定案 |
|----|------|
| 命名 | **`schema-ui`** 单命令（`create` / `add` / `upgrade` 子命令） |
| 实现 | **Go 单二进制**：`apps/api/cmd/schema-ui`（同 apps/api 模块 → 分发 = `go install github.com/magicvr/schema-ui-core/apps/api/cmd/schema-ui@vX.Y.Z`，复用已实证的模块 tag/proxy 链） |
| 命令集 | **create + add + upgrade**（对齐判据 #2 双轨对照） |

## 实施约束

- **零新增依赖**：标准库 `flag`/手写子命令解析；模板 = `go:embed`（templates/ 静态资产）。
- create 模板源 = golden-field 已实证骨架的**物化**（组合根/探针/README/npmrc/gitignore——不做二次发明）。
- add/upgrade：操作 registry（`go get` / `pnpm add` 子进程 + 版本管理与探针回归），available 模块清单 = `kernel.BuiltinModules`（CLI 同模块可直接 import）。
- 发布 = 随 apps/api 模块 tag（`apps/api/v0.2.0` 含 CLI 载荷）。

## 未选方案

- Node CLI（npm 分发）：未选（用户选 Go）。
- 独立 CLI 模块（apps/cli 独立 go.mod）：未选——同模块分发最简、共享 kernel/assembly 类型。