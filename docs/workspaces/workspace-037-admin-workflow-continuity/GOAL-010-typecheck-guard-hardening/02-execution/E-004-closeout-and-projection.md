---
id: E-004-closeout-and-projection
doc: execution-entry
status: recorded
goal_id: GOAL-010-typecheck-guard-hardening
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-004 · C4 双审响应、关门与投影

## 事实

### 1. self 审计

`A-001`（self）verdict `pass`，开放 required = 0（详见该条）。

### 2. 独立审计（用户指定路径）

按用户 2026-09-18 指令与项目级决策 `docs/architecture/independent-audit-execution.md`，调用本地 grok build 执行独立交叉审计：`grok --prompt-file <audit 指令> -m grok-4.6 --reasoning-effort xhigh`（本目标的 `/audit` 流程）。意见落盘为 `A-002`（`source: independent`，`auditor: grok-build (grok-4.6 · reasoning xhigh)`），verdict **`pass`**，开放 required = **0**。

审计员独立执行的可核对动作包括：仓库内注入类型错误复测四命令；同构 scratch 复测退出码（0 / 0 / 2 / 1）；`npm run typecheck` exit 0；守卫 8/8；自构反例（在 `apps/web/scripts/generate-claim.mjs` 插入 `npx tsc --noEmit -p tsconfig.json`）被扫描路径精确打出并还原；16 种额外命令形态（`-p=./tsconfig.json`、`--project ./`、`sh -c`、`-p $CFG` 等）13 违规 / 2 合法；三个 `E-003` 未列出的新变异 N1/N2 被捕获；`npm test` 112 文件 / 1428 测试通过。

### 3. 意见响应（P-003：全部意见）

| 意见 | 级别 | 状态 | 响应 |
|------|------|------|------|
| `A-002 F-001` 动态目标「命中项必须全部选择源文件」没有回归用例（N3 变异未被捕获） | medium · recommended | **fixed** | 补 `COMMAND_FORMS` 两条整值通配用例（`-p $CFG`、`--project ${CFG}` → 必须判空转）与新增用例 `fails closed when an interpolated -p target could denote the root config`（`$CFG`/`${CFG}`/`./$CFG` → false；`tsconfig.${...}.json` → true）。复算 `A-002` 的 N3 变异：删去 `matches.every(...)` 后**现被捕获**（2 failed / 7 passed），撤销后 9/9 通过。 |
| `A-002 F-002` 执行索引把尚未落盘的 `E-004` 标成 recorded | low · recommended | **fixed** | 本条即 `E-004`，已落盘并覆盖 C4 全部事实；索引与文件一致。 |

`A-001` 与 `A-002` 结论一致（均 `pass`、开放 required = 0），无 P-004 冲突需用户裁决。`A-002` 的残余与不确定性（仓库内注入的数字退出码、YAML 跨行 argv、JSONC/`extends`、扫描边界）不构成本目标缺口，已在该条与 `A-001` 留痕。

### 4. 复审与关门

- `A-003`（independent · same provider）对 F-001/F-002 的修正做 finding-closure 复审。
- 复审通过后：`GOAL-008 A-002 F-001` 按 P-003 的 `fixed` 路径闭合（证据指向本目标 `E-002`/`E-003`/`E-004`、`A-002`/`A-003`）；`GOAL-010` 投影为 `done · 4/4`。

## 投影（复审通过后生效）

| 层 | 变化 |
|----|------|
| `GOAL-008` | `A-002 F-001` `open → fixed`（不改其 `done · 4/4`） |
| `GOAL-010` | `active · 3/4 → done · 4/4`（非纲领整改子目标） |
| Root `GOAL-001` | 不变：`active · 5/6`（整改子目标不计入六阶段分母） |
| VP-037 / workspace | 保持 `active`；`R5-I-004` 用户书面关门确认**未**被本轮关闭或推断 |
| goal-tree / workspace.md | 同步 `GOAL-010` 行与概述 |
