---
id: E-001-typecheck-vacuity-baseline
doc: execution-entry
status: recorded
goal_id: GOAL-008-typecheck-evidence-convention
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-001 · 确认类型检查空转事实与影响面

2026-09-18，完成 C1 的影响面与正确口径基线核对：

## 结构与机制

- `apps/web/tsconfig.json` 为 solution-style：`{"files": [], "references": [{"path": "./tsconfig.app.json"}, {"path": "./tsconfig.node.json"}]}`，自身不含 `include`/`files` 源文件。
- `apps/web/tsconfig.app.json` 才含 `"include": ["src"]` 与 `strict: true` 等实际编译选项。
- `apps/web/package.json` 的 `build` 脚本为 `tsc -b && vite build`；`apps/web/README.md` 亦记录 `npm run build # tsc -b && vite build`。**仓库既有约定本就是 `tsc -b`**。

## 空转证明（决定性）

在 `apps/web/src/components/list-filter-panel.tsx` 末尾注入 `const __probe: number = "definitely a string";`：

| 命令 | 输出 | 退出码 |
|------|------|--------|
| `npm exec -- tsc --noEmit` | 无输出 | **0**（未发现错误） |
| `npm exec -- tsc -b` | `error TS2322: Type 'string' is not assignable to type 'number'` + `error TS6133` | **2** |

注入后已还原文件（校验哈希一致）。另有一次实践复现：C8 修复 F-004 过程中曾出现 `ReferenceError: hiddenAtMobile is not defined`（5 项测试失败），而同期 `tsc --noEmit` 仍报 exit 0。

## 影响面

`apps/web/tsconfig.json` 自脚手架提交 `a3e1e5ad` 起即为此形态，此后仅 `e4751e70` 改动过（未触及 `files`/`references` 结构），故**所有使用裸 `tsc --noEmit` 的历史条目都同样空转**。

已检索到的裸 `tsc --noEmit`（未带 `-p`）条目文件：

| 工作区 | 文件 | 本目标处置 |
|--------|------|-----------|
| workspace-037 | `GOAL-007.../E-004`、`E-005`、`E-006` | 已由 R6 更正为 `tsc -b` |
| workspace-002 | `GOAL-008-r5-engineering-fork/02-execution.md`、`03-audit.md` | 只登记（跨区） |
| workspace-009 | `GOAL-005-w4.../E-001`、`A-001`；`GOAL-016-w15.../A-001` | 只登记（跨区） |
| workspace-010 | `GOAL-004-w3.../E-005`、`E-007`；`GOAL-038-w26.../E-002`、`E-003`；`GOAL-039-w27.../E-001` | 只登记（跨区） |
| workspace-011 | `GOAL-018-mfa-manager-ui/A-004`；`GOAL-022-my-wallet.../E-003`、`E-004` | 只登记（跨区） |

注：部分条目可能位于当时有效的工作目录（如带 `-p` 或当时 `tsconfig` 形态不同），本目标不逐一断言其有效性，只登记为「需复核」清单。

C1 完成；C2 见 `E-002`。跨区追溯更正保持 `I-008-004` deferred。
