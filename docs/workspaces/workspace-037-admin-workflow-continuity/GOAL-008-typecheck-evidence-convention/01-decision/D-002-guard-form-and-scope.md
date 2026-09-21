---
id: D-002-guard-form-and-scope
doc: decision-entry
status: accepted
goal_id: GOAL-008-typecheck-evidence-convention
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-002 · 防复发守卫形态：结构守卫测试 + CI 显式门禁

## 决定

`I-008-003` 的守卫形态确定为**组合方案**，而非三选一：

1. **结构守卫测试**（主守卫）：新增 `apps/web/src/typecheck-convention.guard.test.ts`，沿用本仓既有 `.guard.test.ts` 约定（同 `fixture-root.guard.test.ts` / `nginx-proxy.test.ts`：读仓库文件做结构断言）。它断言四件事：
   - 根 `tsconfig.json` 保持 solution-style（这是「裸 `tsc` 空转」的前提，前提变了必须回来重审约定）；
   - `package.json` 存在 `typecheck` 脚本且使用 `-b`/`-p` 这类真正检查的形式；
   - `build` 脚本仍含检查型 `tsc`；
   - **所有可执行面**（`package.json` 脚本、`.github/workflows/*.yml`、仓库 `scripts/`、`apps/web/scripts/`）中不存在非检查型 `tsc` 调用；
   - README 仍记录该约定。
2. **CI 显式门禁**（补强）：`.github/workflows/r6-basic-matrix.yml` 的 web job 在 `npm test` 之后增加 `npm run typecheck`。守卫测试只断言**约定**，这一步才真正**类型检查** `src/**`。

## 理由

**为什么选结构守卫测试而非另两种**：

- 本仓已有成熟的同型先例（`fixture-root.guard.test.ts` 正是「让某条约定成为 CI 期事实」的写法），复用它使守卫的形态、位置与失败信息对后续维护者都是熟悉的，不需要引入新机制或新工具。
- 它随 `npm test` 自动运行，而 `npm test` **已在 CI 中**——无需新增 CI 配置即可生效，降低「守卫写了但没人跑」的风险。
- 相比「仅约定文档权威化」，结构断言是**机器可核对**的，符合 P-002「审计须能指回证据」；文档会漂移，断言不会。

**为什么还要加 CI 显式门禁**：这是本轮最关键的判断。守卫测试断言的是「命令写对了没有」，它**本身不做类型检查**——`src/**` 里真实的类型错误仍然只能靠 `tsc -b` 发现。若只加守卫而不加门禁，就会出现「守卫全绿但代码有类型错误」的空档，等于用一种新的空转替换旧的空转。因此两者必须同时存在。

**为什么守卫不扫描治理正文**：历史 E/A 条目在**解释这个缺陷时**必然要引用那条错误命令（如 `E-001` 的证据表、`A-002 F-005` 的描述）。若守卫扫描 `docs/`，它会强迫改写历史记录或引入豁免清单，而 P-002 要求「实施事实」保持可追溯、不得美化。正文纪律由 README 约定 + 「正确命令已是最省力选项」共同承担。

## 未选方案

- **只加 CI 门禁不加守卫**：能挡住类型错误，但挡不住「有人在脚本/工作流里新写一条 `tsc --noEmit` 并把它当作证据」——那正是本缺陷的历史成因。
- **只加守卫不加门禁**：见上，会留下「约定正确但代码未检查」的空档。
- **守卫扫描 `docs/` 治理正文**：会与历史记录的诚实性冲突，且需要豁免清单，脆弱且违背 P-002。
- **改成 ESLint/自定义 npm script 做守卫**：本仓无 ESLint 配置（`apps/web` 无 `eslint.config.*`），引入新工具链超出本目标必要范围。
- **修改 `tsconfig.json` 使其非 solution-style**：会牵动 `tsconfig.app.json`/`tsconfig.node.json` 与构建链，且 solution-style 本身正当；本目标只纠正**命令选择**，不改配置结构。

## 门禁

守卫已实施并通过变异验证（4 种变异均被捕获，见 `E-003`）。C4 自审须核对：守卫的非空转性（变异测试证据）、CI 门禁的存在、以及 README 约定的保留。
