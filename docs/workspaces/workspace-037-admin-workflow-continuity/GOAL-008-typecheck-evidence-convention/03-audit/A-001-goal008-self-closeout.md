---
id: A-001-goal008-self-closeout
doc: audit-entry
status: recorded
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
goal_id: GOAL-008-typecheck-evidence-convention
source: self
auditor: dsh / deepseek-v4.1-flash（编排器自审）
scope: GOAL-008 C1-C4：F-005 承接、口径固化、防复发守卫与投影
verdict: pass
---

# A-001 · GOAL-008 self close-out

## 头字段

- **source**：self
- **auditor**：dsh / deepseek-v4.1-flash（编排器自审）
- **类型**：`close-out`
- **scope**：GOAL-008 的 C1（影响面/口径基线）、C2（入口固化）、C3（防复发守卫）、C4（交付与投影）
- **verdict**：`pass`

## 范围与区间

被审区间：`ff6169f7`（GOAL-008 开设与 C1/C2）→ 本轮 C3/C4 工作树。核对方式为独立复核：重新读取实现与配置、对守卫做变异测试、对 `typecheck` 入口做注入错误验证、逐一核对可执行面覆盖与文档一致性。

## 成果（有证据）

1. **F-005 已闭环（本目标承接部分）**。根因（solution-style `tsconfig.json` 使裸 `tsc --noEmit` 空转）已由 `E-001` 的注入错误证明；正确口径 `tsc -b` 与仓库既有 `build` 脚本约定一致。
2. **正确入口可发现且真实有效**。`npm run typecheck` 存在；注入 `const __probe: number = "str"` 时返回非零并报 `TS2322`/`TS6133`，干净时 exit 0。
3. **发现并修复了同类的第二层缺口（本轮新增）**：`e2e/tsconfig.json` 存在但**不在根配置的 `references` 图内**，故单独的 `tsc -b` 不检查 e2e 规格文件——这是同一缺陷下沉一层的复现。实测：注入 e2e 错误后 `tsc -b` 仍 exit 0，而 `tsc -p e2e/tsconfig.json` 报 `TS2322` exit 2。已把 `typecheck` 扩为 `tsc -b && tsc -p e2e/tsconfig.json`，注入验证通过（见 `E-004`）。
4. **防复发守卫已实施且经变异验证**。`src/typecheck-convention.guard.test.ts`（6 断言）随 `npm test` 自动运行，而 `npm test` 已在 CI 中；另有 CI 显式 `npm run typecheck` 门禁（守卫只断言约定，门禁才真正检查）。
5. **守卫非空转**：对 5 种真实失效模式做变异测试，**5/5 被捕获**（vacuous typecheck、build 失去检查、CI 混入裸命令、tsconfig 前提改变、e2e 覆盖丢失），每次变异后均还原文件。
6. **回归**：`npm run typecheck` exit 0；全量 Web Vitest **112 文件 / 1426 测试**通过；`git diff --check` 干净。
7. **范围纪律**：未写入其他工作区台账（AGENTS §6c）；跨区影响仅登记于 `E-001`，`I-008-004` 保持 deferred。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 结构与空转事实、影响面基线 | 达成 | `D-001`、`E-001`（注入对比表 + 跨区登记清单） |
| C2 `npm run typecheck` 入口 + README 约定，实测有效 | 达成 | `E-002`；注入错误返回非零 |
| C3 守卫已实施、可核对、空转口径不可无声沿用 | 达成 | `E-003`、`E-004`；5/5 变异捕获；CI 门禁 |
| C4 self 审计与 Root 投影 | 达成 | 本条 + `E-005` |
| 不跨区写入 | 达成 | 仅登记，无其他工作区文件改动 |
| 不改变 Root 六阶段分母 | 达成 | 整改子目标，不计入分母；Root 在 R6 完成后投影为 `active · 5/6`（Root `E-021`） |

## Findings

### F-001 · `typecheck` 初始版本漏检 e2e 项目（本轮自纠）

- 严重度：med
- 建议：**required**
- 描述：C2 的 `typecheck` 初始为 `tsc -b`，而 `e2e/tsconfig.json` 不在根配置 `references` 图内，因此 e2e 规格文件仍未被检查——修复了一层空转却在相邻层留下同类缺口。若未在 C4 前发现，本目标会以「已修复类型检查证据」结项，而 e2e 仍处于未检查状态。
- 证据：注入 `e2e/sign-in.ts` 错误后 `npm run typecheck`（当时的 `tsc -b`）exit 0；`tsc -p e2e/tsconfig.json` 报 `TS2322` exit 2。
- 状态：**fixed**（`typecheck` 扩为 `tsc -b && tsc -p e2e/tsconfig.json`；守卫新增第 3 条断言锁定 e2e 覆盖；README 同步说明两层陷阱；变异测试 5 号验证该断言有效）

### F-002 · 守卫只覆盖「命令写法」，不覆盖「命令是否被执行」

- 严重度：low
- 建议：recommended
- 描述：守卫断言可执行面里不存在非检查型 `tsc` 调用，也能锁定 `typecheck` 脚本内容，但它无法断言 CI **确实调用了**该脚本。本轮已加 CI 步骤，若后续有人删除该步骤，守卫不会失败（CI 步骤本身不在断言范围内——守卫只读工作流文本，未断言 `npm run typecheck` 出现在 web job 中）。
- 证据：守卫 `uses no non-checking tsc invocation in any executable surface` 只做否定断言；无「CI web job 含 typecheck 步骤」的正向断言。
- 影响：低——删除 CI 门禁属显式编辑，评审可见；且 `npm test`（守卫本身）仍在 CI 中运行。
- 状态：open（建议后续加固：为工作流增加正向断言）

## 必改项汇总（required）

| finding | 严重度 | 主张 | 状态 |
|---------|--------|------|------|
| F-001 | med | `typecheck` 须覆盖 e2e 项目 | **fixed** |

F-002 为 recommended，不阻断本目标关门。

## 信息就绪核对

| 项 | 状态 | 备注 |
|----|------|------|
| I-008-001（结构与空转事实） | verified | `E-001` 注入对比 |
| I-008-002（影响面范围） | verified | `E-001` 跨区登记清单 |
| I-008-003（守卫形态） | **verified** | `D-002` 决策 + `E-003`/`E-004` 实施与变异验证 |
| I-008-004（跨区追溯） | deferred non-blocking | 用户另行路由；本目标不跨区写入 |
| 到期 required 信息项 | 无 | 无阻断关门的信息门禁 |
| 资料引用 | 无 | 本区 `shared_materials_catalog: none` |

## 结论 + 建议下一步

C1～C4 全部达成，F-005 在本目标承接范围内的闭环成立（含本轮发现并修复的 e2e 第二层缺口）。开放 required finding = 0（F-001 已 fixed），无到期 required 信息项，verdict 为 **`pass`**。

**可关门**：GOAL-008 可投影 `done · 4/4`。

**对 R6 的影响**：按用户 2026-09-18 裁决「等 F-005 处置落定后再关门」，F-005 现已闭环，R6 `GOAL-007` 的 C6 可判完成并投影 `done · 8/8`，Root 相应投影 `active · 5/6`。R5-I-004 用户书面关门确认仍开放，Root/VP 不因本轮关门。

**建议**：F-002 可在后续加固；`I-008-004`（跨区历史条目追溯）仍待用户路由。
