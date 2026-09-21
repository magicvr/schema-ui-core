---
id: E-003-mutation-verification-and-regression
doc: execution-entry
status: recorded
goal_id: GOAL-010-typecheck-guard-hardening
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-003 · 变异验证与全量回归（C3）

## 变异验证（5/5 全部被捕获）

对本轮加固逐项施加「回到旧行为」的变异，运行守卫并记录失败点，随后还原：

| # | 变异 | 结果 | 失败用例 |
|---|------|------|----------|
| M1 | `projectTargetIsRealCheck` 恒返回 `true`（回到「只看有没有 `-p` 令牌」） | 捕获 | `rejects every vacuous command form, including -p at the root config` |
| M2 | `selectsOwnSources` 恒返回 `true`（任何配置都当成选择源文件） | 捕获（3 项） | root solution-style 前置条件 / vacuous command forms / resolves -p targets |
| M3 | 在 `.github/workflows/r6-basic-matrix.yml` 插入 `- run: npx tsc --noEmit -p tsconfig.json` | 捕获 | 可执行面扫描，精确定位 `r6-basic-matrix.yml:31` |
| M4 | `package.json` 的 `typecheck` 改为 `tsc --noEmit -p tsconfig.json` | 捕获（3 项） | typecheck 脚本 / e2e 项目覆盖 / 可执行面扫描 |
| M5 | 引号处理回退为「遇引号即开」（丢弃按位置的收尾引号判定） | 捕获 | 可执行面扫描（`package.json` 的 `typecheck` 行被误报，证明 tokenizer 判定是承重的） |

M3/M4 同时证明：扫描路径确实走加固后的判定，而不是只在新用例表里生效——旧守卫会放行这两处（`-p tsconfig.json` 含 `-p` 令牌）。

**还原核验**：每次变异后均以 `git checkout --`（对 `package.json`／workflow）或反向编辑（对守卫文件）还原；最终 `git diff --stat -- apps/web .github` 仅含 `apps/web/src/typecheck-convention.guard.test.ts`，`.github/**` 与 `apps/web/package.json` 无差异。

## 全量回归

| 检查 | 结果 |
|------|------|
| `npx vitest run src/typecheck-convention.guard.test.ts` | **8 passed**（新增 2 项 + 原 6 项语义不弱化） |
| `npm run typecheck`（`tsc -b` + `tsc -p e2e/tsconfig.json`） | exit **0** |
| `npm test`（全量 Vitest） | **112 文件 / 1428 测试全部通过**（较 GOAL-009 时点 112/1426 增加 2 项，为本轮新增用例） |
| `git diff --check` | 通过 |
| 代码面改动 | 仅守卫测试文件（+239/−14）；无产品代码、无脚本、无 CI 变更 |

C3 完成。C4（self 审计 → grok 独立审计 → 响应与投影）见 `A-001`、`A-002`、`E-004`。
