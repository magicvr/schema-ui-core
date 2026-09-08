---
id: A-013-release-preflight-self
goal_id: GOAL-001-nav-group-collapsible
doc: audit-entry
source: self
auditor: current-session
date: 2026-09-08
scope: release candidate preflight for dev HEAD, apps/api/v0.6.0 version pin, PR/CI/merge/tag/release gates
verdict: pass
parent: null
version: 0.1.0
---

# A-013 · 发布候选预检 self 审计

## 范围与基线

- 当前工作区：`workspace-034-nav-group-collapsible`；目标：`GOAL-001-nav-group-collapsible`（已完成 Root，当前审计只追加发布证据，不改状态或进度）。
- 发布候选分支：`dev`；版本钉提交：`06e4c2a4`（`chore(release): pin apps/api/v0.6.0 package surface`）；随后补充 E2E 竞态修复提交 `b0af19b0`、`ba01663f`。
- 版本选择：在既有 `apps/api/v0.5.0` tag 之后，本轮新增导航分组/Manifest 聚合能力属于向后兼容的可选能力扩展，采用下一 minor `apps/api/v0.6.0`；QUICKSTART 与 CLI 模板版本钉已同步。

## 成果（有证据）

| 检查项 | 结果 | 证据 |
|---|---|---|
| API 全量测试 | pass | `apps/api`：`go test ./... -count=1` exit 0 |
| API 静态检查 | pass | `apps/api`：`go vet ./...` exit 0 |
| API 构建 | pass | `apps/api`：`go build ./...` exit 0 |
| Web 全量单测 | pass | `apps/web`：Vitest 99 files / 1343 tests 全部通过 |
| Web 生产构建 | pass | `apps/web`：`npm run build` exit 0；Vite 产物构建成功 |
| Web 浏览器 E2E | pass | 修复 fallback login 请求竞态后，`npm run test:e2e`：14 tests，10 passed，4 skipped，exit 0；SQLite 浏览器矩阵完成 |
| 版本文档一致性 | pass | `apps/api/cmd/schema-ui/main.go` 与 `QUICKSTART.md` 均指向 `apps/api/v0.6.0` |
| 工作树边界 | pass | 版本钉与 E2E 稳定性修复均按 owned paths 提交；构建生成 claim 文件未作为本次版本钉改动提交 |
| 本地生产冒烟 | 未执行 | Git Bash 可用但根 `.env` 未提供 `ADMIN_INITIAL_PASSWORD`/`AUTH_JWT_SECRET`；未猜测或输出秘密，Hosted container-smoke 仍为发布门禁 |

## 对照发布成功标准

1. **候选可提交 PR**：满足。代码变更已在当前 `dev` 历史，版本钉已形成独立提交。
2. **Hosted CI 全绿**：本条在 PR 创建后由 GitHub Actions `r6-basic-matrix` 验证；本审计不把本地结果冒充 Hosted CI 证据。
3. **合并到 `main`**：须在 PR 全部检查成功后执行；本审计不宣称已合并。
4. **tag 与 Release 资产**：须在合并后的 `main` commit 上创建 annotated `apps/api/v0.6.0` tag，并在资产构建/校验后创建 GitHub Release；本审计不宣称已发布。

## Findings

| finding | level | 状态 | 证据 / 响应 |
|---|---|---|---|
| — | — | — | 本 scope 未发现 required 或 recommended finding。Hosted CI、PR merge、tag 与 Release 仍是后续事实门禁，不以“待执行”伪装为完成。 |

## 结论

对“版本钉 + 本地回归 + 可进入 PR/Hosted CI 门禁”的 scope 判定为 **pass**。在 Hosted CI 全绿、PR 合并成功、tag 指向合并后的 `main` commit、资产 SHA-256 校验和 GitHub Release 可访问之前，不得宣称本次发布完成。

## 声明

本意见只记录当前已发生的预检事实，不修改 `GOAL-001-nav-group-collapsible` 的 `status`、`progress` 或 `goal-tree`；后续独立意见与门禁响应由编排流程处理。
