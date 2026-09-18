---
id: E-001-vacuity-boundary-and-rule-baseline
doc: execution-entry
status: recorded
goal_id: GOAL-010-typecheck-guard-hardening
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-001 · 空转边界复算与判定规则冻结（C1）

## 事实

2026-09-18，完成 C1：复算空转形态边界并冻结 `-p` 目标判定规则（`D-001`）。

1. **空转边界复算**（与 `GOAL-008 E-006` 同一方法，本轮独立重跑）：在 `apps/web/src/renderer/list-surface.tsx` 末尾注入 `const __probeErr: number = "definitely a string";`，在 `apps/web` 下执行四条命令，随后 `git checkout` 还原：

| 命令 | 结果 | 判定 |
|------|------|------|
| `npx tsc --noEmit` | 无输出，exit 0 | 空转 |
| `npx tsc --noEmit -p tsconfig.json` | 无输出，exit 0 | **空转**（`-p` 指向 solution-style 根配置） |
| `npx tsc --noEmit -p tsconfig.app.json` | `TS2322` + `TS6133`，非零 | 真实检查 |
| `npx tsc -b` | `TS2322` + `TS6133`，非零 | 真实检查 |

还原后 `git status` 对 `apps/web` 无输出，工作树干净。

2. **配置形态盘点**（`apps/web`）：`tsconfig.json` 为 solution-style（`files: []` + 2 条 `references`）；其余 8 个 `tsconfig.*.json` 全部通过 `include` 选择源文件——`app`(`src`)、`node`(`vite.config.ts`/`vitest.config.ts`)、`lib`、`protocol`、`renderer`、`shell`、`theme`、`ui`；`e2e/tsconfig.json` 选择 `./**/*.ts`。故 `-p <具体项目配置>` 是真实检查，`-p tsconfig.json` 与 `-p .` 不是。

3. **判定规则冻结**（`D-001`）：`tsc` 调用构成检查型当且仅当带 `-b`/`--build`，或 `-p`/`--project` 目标解析到自身选择源文件的配置；目标缺失/不可读/不可解析一律 fail closed。

4. **既有可执行面复算**：判定规则收紧后，`package.json`（`typecheck`、`build`）、`.github/workflows/*.yml`、仓库与 `apps/web` 的 `scripts/` 中**没有**需要改动的调用；唯一非字面目标是 `scripts/build-lib-packages.mjs` 用 `tsconfig.${name.split("/")[1]}.json` 逐包构建，该形态需要在 C2 中被正确解析（否则会成为误报）。

C1 完成；`I-010-001`、`I-010-002` 为 `verified`。C2 见 `E-002`。
