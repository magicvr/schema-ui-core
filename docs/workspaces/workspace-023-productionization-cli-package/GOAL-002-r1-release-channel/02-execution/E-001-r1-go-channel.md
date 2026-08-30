---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-release-channel
version: 0.1.0
---

# E-001 · S1 Go 通道实证（2026-08-29）

## 实证序列（真实操作）

| 步骤 | 动作 | 结果 |
|------|------|------|
| 1 | `git tag v0.1.0` + push origin | ✅ 推送成功，但**消费失败**：`unknown revision apps/api/v0.1.0` |
| 2 | **发现**：module path = `github.com/magicvr/schema-ui-core/apps/api`（子目录）→ Go 要求 tag 形如 **`apps/api/v0.1.0`**（子目录模块 tag 约定；G1 单模块粗粒度也不例外） | 教训：tag 命名与 module path 前缀绑定 |
| 3 | `git tag apps/api/v0.1.0` + push；删除误导 tag `v0.1.0`（origin 与本地） | ✅ |
| 4 | golden-field `go mod tidy`（公共 proxy） | 初败：**sumdb 收录时延**（`sum.golang.org/lookup … 404`）——新 tag 索引延迟 |
| 5 | `GOSUMDB=off go mod tidy`（临时绕行）→ `go run ./cmd/server` | ✅ `go: downloading github.com/magicvr/schema-ui-core/apps/api v0.1.0` · `go.sum` 记 `h1:DAkEPGy…` · 运行全绿（kernel=2.0.0 · fresh=true · contrib 与基线一致） |

## 知识项（产线流程必须）

1. **tag 命名**：子目录 module path → tag 前缀 = 模块路径相对根的子目录（`apps/api/`）；发布脚本须按 `$(module-dir)/vX.Y.Z` 打 tag。
2. **sumdb 收录时延**：新 tag 公开后 sumdb 校验 404 属正常（分钟~小时级）；CI/脚本首次发布后应等待 sumdb 收录再消费（或临时 `GOSUMDB=off` 并在收录后重跑默认校验补齐 go.sum）。
3. **发布顺序**：tag 必须指向包含完整发布载荷的 commit（本次 `apps/api/v0.1.0` = HEAD `c63492b3`，含 pack 脚本与 GOAL-023 文档——R4 审计曾指出 v0.0.2 旧 tag 不含载荷的问题，本次已规避）。

## 现状

- I-023-001 **verified**；golden-field `go.mod` 无 replace（registry 语义）✓。
- **npm 侧（S2）阻塞于凭据**：GitHub Packages 目标定案（scope 须 = owner `@magicvr`；registry `https://npm.pkg.github.com`），等待用户提供发布凭据（GH token with `write:packages`）或配置方式。