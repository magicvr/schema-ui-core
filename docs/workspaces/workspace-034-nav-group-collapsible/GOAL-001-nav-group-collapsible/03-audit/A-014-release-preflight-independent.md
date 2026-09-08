---
id: A-014-release-preflight-independent
goal_id: GOAL-001-nav-group-collapsible
doc: audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-08
scope: release candidate preflight for dev HEAD, apps/api/v0.6.0 version pin, navigation-group API/Web changes, PR/Hosted CI/merge/tag/Release gates
verdict: pass
parent: null
version: 0.1.0
---

# A-014 · 发布候选预检 independent 审计

## 范围与基线

- 工作区：`workspace-034-nav-group-collapsible`；目标：`GOAL-001-nav-group-collapsible`。
- 候选分支：`dev`，相对 `origin/dev` 的新增提交为 `06e4c2a4`、`b0af19b0`、`ba01663f`；当前尚未创建开放 PR、`apps/api/v0.6.0` tag 或 GitHub Release。
- 审计重点：导航分组/Manifest 聚合交付是否有可核对的发布前验证；CLI/API 版本钉与 QUICKSTART 是否一致；E2E fallback-login 竞态修复是否合理；Hosted CI、合并、tag 和资产发布是否被错误地提前宣称完成。
- 本条由项目约定的本地 Grok Build 独立审计路径执行，意见由编排器代贴到正式 Goal 审计台账；未修改实现、目标状态、进度或 goal-tree。

## 成果（有证据）

| 检查项 | 结果 | 证据 |
|---|---|---|
| 导航分组实现边界 | pass | API `NavigationContribution`/Manifest 聚合与 Web `projectNavigation`、分组交互/回归测试均在当前候选范围内；未发现把分组能力错误写入上游协议包版本的证据 |
| 版本钉一致性 | pass | CLI 模板 `apiVersion` 与 `QUICKSTART.md` 均指向 `apps/api/v0.6.0`；六个 Web 包版本保持既有冻结面，未见本轮包源码变更需要升版的证据 |
| API 本地验证 | pass | 独立审计会话复核 `apps/api` 的 `go vet ./...`、`go test ./... -count=1`、`go build ./...` 均通过 |
| Web 本地验证 | pass | 当前候选已有 Vitest 99 files / 1343 tests、`npm run build` 通过证据；构建生成的 claim 文件未被当作版本钉改动提交 |
| 浏览器 E2E 稳定性 | pass | fallback login 竞态修复后，SQLite `npm run test:e2e` 完成 14 tests：10 passed、4 skipped、exit 0；首次失败属于修复前的异步登录提交态竞态，修复将“等待初次 POST 结束”置于替换密码之前 |
| 生产冒烟边界 | 证据不足但不伪装 | 本地未执行 `pre-release-smoke.sh`：Git Bash 可运行，但根 `.env` 未提供 `ADMIN_INITIAL_PASSWORD`/`AUTH_JWT_SECRET`；该环境事实未被猜测，Hosted container-smoke 仍必须作为发布门禁 |
| CI/PR/merge/tag/Release | 未发生 | 审计确认这些是后续外部事实门禁，不将其写成已完成 |

## 对照成功标准

1. **进入 PR/Hosted CI 门禁**：满足。当前候选有版本钉、导航交付回归、API/Web 本地证据和 E2E 竞态修复，可提交 PR。
2. **Hosted CI 全绿**：尚未发生；必须等待 `r6-basic-matrix` 的 web、api、api-postgres、browser-e2e 四腿矩阵以及 mvp/admin container smoke 全部 `success`。
3. **合并到 `main`**：尚未发生；仅允许在 PR 所有检查完成且 head 未漂移后合并。
4. **tag 与新资产**：尚未发生；应在合并后的 `main` merge commit 上创建 annotated `apps/api/v0.6.0`，构建并校验 release assets 后再创建 GitHub Release。

## Findings

| finding | level | 状态 | 证据 / 响应 |
|---|---|---|---|
| — | — | — | **无开放 required/recommended finding**。Hosted CI、PR 合并、tag、资产 SHA-256 与 GitHub Release 是尚未发生的发布门禁事实，不作为 finding 关闭，也不作为已完成事实。 |

## 必改项汇总

- 无。发布编排必须继续保持 CI 全绿 → merge → tag → 资产校验 → Release 的顺序；任何一步失败均不得宣称发布完成。

## 与 A-013 self 的异同

- 与 A-013 self 同向：版本选择为 `apps/api/v0.6.0`，API/Web 基础验证通过，Hosted CI/merge/tag/Release 尚未发生。
- 独立审计额外核对了六包冻结面与导航分组能力的边界，并确认 E2E fallback-login 修复针对的是异步提交竞态，而不是把失败登录伪装为成功。
- 本条保留本地 production smoke 的证据边界：缺少必要 secret 时不运行、不猜测；这不替代 Hosted container-smoke。

## 结论与给编排器/用户的下一步

本 scope 判定 **pass**，无开放 required/recommended finding，可以进入 PR 与 Hosted CI 门禁。后续只允许按“PR 检查全绿 → 合并 main → annotated tag → 构建/校验资产 → GitHub Release”顺序推进；`pass` 不代表这些外部事实已经完成。建议由 `/govern` 汇总本条与 A-013 后继续推进，并在所有外部事实完成后补写执行台账。

## 声明

本意见不修改 `status`/`progress`；响应由 `/govern` 处理。
