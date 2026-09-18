---
id: E-002-typecheck-entrypoint
doc: execution-entry
status: recorded
goal_id: GOAL-008-typecheck-evidence-convention
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-002 · 固化 `typecheck` 脚本入口与 README 约定

2026-09-18，完成 C2：

- `apps/web/package.json` 新增脚本 `"typecheck": "tsc -b"`，使正确命令成为最省力的默认可发现入口（与既有 `build` 脚本的 `tsc -b` 一致）。
- `apps/web/README.md` 在 Run 区块补 `npm run typecheck` 说明，并新增「类型检查口径（GOAL-008 · 必读）」小节，写明：`tsconfig.json` 为 solution-style、真正编译选项在 `tsconfig.app.json`/`tsconfig.node.json`、**不带 `-b` 的 `tsc --noEmit` 不检查任何文件且恒返回 exit 0、不构成类型检查证据**，以及正确写法与「`tsc -b` 遇真实错误返回非零」的门禁用法。

## 验证事实

| 场景 | 命令 | 结果 |
|------|------|------|
| 干净工作树 | `npm run typecheck` | 输出 `tsc -b`，exit 0 |
| 注入 `const __probe: number = "str";` | `npm run typecheck` | `error TS2322`（Type 'string' is not assignable to type 'number'）+ `TS6133`，非零退出 |

注入后已还原文件（校验哈希一致）。这证明新入口**确实执行类型检查**，而裸 `tsc --noEmit` 不会。

C2 完成。C3（防复发守卫形态，`I-008-003`）待决策；跨区追溯更正保持 `I-008-004` deferred。
