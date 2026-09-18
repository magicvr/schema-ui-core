---
id: E-006-cross-workspace-typecheck-errata
doc: execution-entry
status: recorded
goal_id: GOAL-008-typecheck-evidence-convention
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-006 · 跨工作区类型检查空转条目追溯更正（用户授权）

## 授权

2026-09-18，用户就 `I-008-004`（其他工作区历史条目是否追溯更正）选择**「授权追溯更正」**，解除本目标原先「只登记、不代改」的边界（AGENTS §6c 的跨区写入限制由用户明确授权放行，动作范围限定为**加勘误注记**）。Root `D-015` 记录该裁决。

## 处置原则

- **只加勘误注记，不重写历史**：原命令行、原结论、原 verdict 一律保留；勘误以行内括注或紧随其后的引用块追加，并注明追加日期与「原记录保留不改」。
- **跨区引用用 Q2 路径**：注记统一指向 `docs/workspaces/workspace-037-admin-workflow-continuity/GOAL-008-typecheck-evidence-convention/`（D-001、E-006、A-002）。
- **不改任何目标的 `status`/`progress`/verdict/finding 状态**；不改 `goal-tree.md` 之外的他人台账索引。
- **不冒充审计**：追加在 `03-audit/A-*.md` 中的注记明确写「不改本条 verdict 与结论」，其性质是事后勘误，不是第二条意见。

## 追加的勘误注记（11 处，覆盖 4 个工作区）

| 工作区 | 文件 | 原条目 | 性质 |
|--------|------|--------|------|
| workspace-009 | `GOAL-005-w4-security-audit-remediation/02-execution/E-001-w4-remediation.md` | `npx tsc --noEmit`（apps/web）：clean | 空转 |
| workspace-009 | `GOAL-005-w4-security-audit-remediation/03-audit/A-001-w4-self.md` | 复跑行 `tsc --noEmit` clean | 空转 |
| workspace-009 | `GOAL-016-w15-api-web-audit-remediation/03-audit/A-001-w15-independent-intake.md` | 验证基线 `tsc --noEmit` 通过 | 空转 |
| workspace-010 | `GOAL-004-w3-schema-host-protocol-conformance/02-execution/E-005-formal-v2-8-0-identity-repin.md` | 门禁重跑 `tsc --noEmit` 无错误 | 空转 |
| workspace-010 | `GOAL-004-w3-schema-host-protocol-conformance/02-execution/E-007-s4-remediation-and-s5-verification.md` | S5 证据表 TypeScript 行 | 空转 |
| workspace-010 | `GOAL-038-w26-email-display-and-mail-pages/02-execution/E-002-s2-implementation.md` | `npx tsc --noEmit` **0**（同行 `npm run build` 成功） | 空转；`tsc -b` 由 build 覆盖 |
| workspace-010 | `GOAL-038-w26-email-display-and-mail-pages/02-execution/E-003-s3-regression-and-go-verdict.md` | 回归表「类型检查」行（同表有 `npm run build`） | 空转；`tsc -b` 由 build 覆盖 |
| workspace-010 | `GOAL-039-w27-invite-outbox-filter-sort/02-execution/E-001-s2-s3-implementation-and-regression.md` | S3 回归表 `tsc --noEmit`（同表有 `npm run build`） | 空转；`tsc -b` 由 build 覆盖 |
| workspace-011 | `GOAL-018-mfa-manager-ui/03-audit/A-004-s5-reaudit.md` | `tsc --noEmit -p tsconfig.json` 及两处简写 | **空转（`-p` 指向根 solution-style 配置）** |
| workspace-011 | `GOAL-022-my-wallet-self-service/02-execution/E-003-s2-implemented.md` | `tsc --noEmit` 无错 | 空转 |
| workspace-011 | `GOAL-022-my-wallet-self-service/02-execution/E-004-s3-verified.md` | 验证表 `tsc --noEmit（apps/web）` | 空转 |

## 复核后**确认有效、未改动**（2 处）

| 工作区 | 文件 | 条目 | 判定 |
|--------|------|------|------|
| workspace-002 | `GOAL-008-r5-engineering-fork/02-execution.md` | `npx tsc --noEmit -p e2e/tsconfig.json` | **有效**：`e2e/tsconfig.json` 自身含 `include: ["./**/*.ts"]` |
| workspace-002 | `GOAL-008-r5-engineering-fork/03-audit.md` | 同上（A-020 复核） | **有效**，未改动 |

同工作区-037 自身 R2/R3/R4 使用的 `tsc -p tsconfig.app.json --noEmit`（R5 `E-004` 同口径）亦为**有效**检查：`tsconfig.app.json` 含 `"include": ["src"]`。

## 本次新增的第二层事实：`-p tsconfig.json` 同样空转

F-005 原判据只说「裸 `tsc --noEmit`（不带 `-b`/`-p`）空转」，并把「带 `-p`」留作「需复核」。本次以注入错误实测（`apps/web/src/renderer/list-surface.tsx` 追加 `const __probeErr: number = "definitely a string";`，随后 `git checkout` 还原、工作树干净）：

| 命令（workdir = `apps/web`） | 结果 |
|------------------------------|------|
| `npx tsc --noEmit` | 无输出，exit **0**（空转） |
| `npx tsc --noEmit -p tsconfig.json` | 无输出，exit **0**（**空转**） |
| `npx tsc --noEmit -p tsconfig.app.json` | `TS2322` + `TS6133`，非零（真实检查） |
| `npx tsc -b` | `TS2322` + `TS6133`，非零（真实检查） |

结论：判据须收紧为「**只有 `tsc -b`，或 `-p` 指向自身选择源文件的项目配置（`tsconfig.app.json` / `e2e/tsconfig.json`），才构成类型检查证据**」。该收紧暴露了防复发守卫的一处同类缺口，记为 `A-002 F-001`。

## 边界（本次未做）

- **未**逐条裁定全仓所有「`tsc` 简写」（如「tsc clean / tsc 0 / tsc 全绿」，共 300+ 行）：本次范围限定为字面 `tsc --noEmit` 与 `-p tsconfig.json` 条目。凡同条目另有 `npm run build`（`tsc -b && vite build`，`apps/web/package.json` 自建仓起即此形态）或显式 `-b`/有效 `-p` 者，实质类型检查成立，已在注记中写明。
- **未**改 `apps/web/src/typecheck-convention.guard.test.ts`：守卫缺口作为新 finding 记录，是否修复待用户裁决（新目标或并入后续轮次）。

## 结论

`I-008-004` 由 `deferred` 转入 `verified`（用户授权 + 11 处注记已落盘，2 处复核确认无需改动）。本条目不改变 GOAL-008 的 `done · 4/4`，也不关闭 `A-002 F-001`（recommended，保持 open）。
