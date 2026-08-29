---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-007-r6-migration-tooling
version: 0.1.0
---

# E-002 · S2/S3 实现与实测（2026-08-29）

## 交付物

- `apps/api/cmd/schema-ui/migrate.go` + main.go 挂载（usage + case）：`schema-ui migrate-fork [--dir <path>] [--dry-run]`。
  - 检查面：go.mod require 版本 · web/.npmrc scope 映射 · cmd/server/main.go 形态（薄封装 / 旧组合根 / 手搓 kernel）
  - 类型判定（指南 §1 程序化 · 校准）：A 纯装配/薄封装 · B 轻度定制（旧组合根 · 无 kernel 覆盖）· C 深度定制（手搓 kernel 且未走组合路径 → 建议保持 fork）
  - 非破坏原则：实跑仅 go.mod require bump（`go get @latest` registry 语义）+ .npmrc 钉 npmjs（备份 `.npmrc.migrate.bak`）；用户代码只引导不覆盖
- `go build ./cmd/schema-ui` exit 0 · `go vet` exit 0。

## 实测（C1–C3 · 9510023 = golden-field v0.3.0 旧态）

| 阶段 | 结果 |
|------|------|
| dry-run | 类型判定 **B**（旧组合根 · assembly.OpenStore 标准件不算 kernel 覆盖）· 4 步清单（bump /.npmrc/引导/验证）· 未写文件 |
| 实跑 migrate-fork | `go get apps/api@latest`：**v0.3.0 → v0.4.0**（registry 语义）· `.npmrc` GH 映射 → `@magicvr:registry=https://registry.npmjs.org`（备份 `.npmrc.migrate.bak`）· main.go 引导（未覆盖）· migrate exit 0 |
| 迁移后验证 | `go mod tidy` exit 0 · `go build ./cmd/server` **exit 0**（旧组合根对 v0.4.0 server 面兼容构建）· 备份存在 ✓ |
| 清理 | worktree 已 remove + prune |

## 核销

- VP-024 判据 #7（fork→包迁移工具化）✅（指南成文 → 工具化 CLI · A/B 型实测）
- go 后清单「fork→包迁移指南（工具化）＝go 后」→ `schema-ui migrate-fork` ✅

## 边界注记

- C 型（深度定制 fork）：工具建议保持 fork（Charter fork 并存），包化承载面（assembly 扩展 + 六包 external 组合）登记为 R7 后候选。
- 迁移引导的完整迁移（薄封装替换 + 业务模块迁入）仍由 `schema-ui create` 演练路径执行（R2 冒烟已证）；migrate-fork 负责**判定 + 低侵入改写 + 引导**。