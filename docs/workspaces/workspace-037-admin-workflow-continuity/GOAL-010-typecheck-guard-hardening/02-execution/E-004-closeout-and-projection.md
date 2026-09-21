---
id: E-004-closeout-and-projection
doc: execution-entry
status: recorded
goal_id: GOAL-010-typecheck-guard-hardening
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.1.0
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

- `A-003`（independent · 同一 provider）对 F-001/F-002 做 finding-closure 复审：verdict **`pass`**，`F-001` → **fixed**，`F-002` → **fixed**，**无新增 finding**，开放 required = 0。该条独立重放 N3（2 failed / 7 passed）、复算通配命中集（`$CFG` 命中 9 个顶层配置含根配置 → `every` 为 false）、核对 `E-001`～`E-004` 索引/目录/frontmatter 一致，并确认产品/CI 边界未破。
- `A-003` 的两点精度已记入其残余（非 finding）：`` `tsc --noEmit --project ${CFG}` `` 表行经扫描路径会被 `{`/`}` 拆词，故 N3 下该行仍判空转（覆盖 `every` 的是 `-p $CFG` 表行与函数级 `${CFG}` 断言）；`./$CFG` 命中集为空，走的是「至少命中一个」分支。两者当前行为均为 fail closed。
- 三方意见（`A-001` self、`A-002` independent、`A-003` independent）结论一致，无 P-004 冲突；无未合法闭合的 required/必改 finding，满足用户设定的关门条件。

### 5. 关门（已执行）

- `GOAL-008 A-002 F-001` 按 P-003 的 `fixed` 路径闭合（证据：本目标 `E-002`/`E-003`/`E-004`、`A-002`、`A-003`）。
- `GOAL-010` 投影为 `done · 4/4`（C4 勾选），`00-meta`/`goal-tree`/`workspace.md` 同步。

## 投影（已生效 · 2026-09-18）

| 层 | 变化 |
|----|------|
| `GOAL-008` | `A-002 F-001` `open → fixed`（不改其 `done · 4/4`） |
| `GOAL-010` | `active · 3/4 → done · 4/4`（非纲领整改子目标） |
| Root `GOAL-001` | 不变：`active · 5/6`（整改子目标不计入六阶段分母） |
| VP-037 / workspace | 保持 `active`；`R5-I-004` 用户书面关门确认**未**被本轮关闭或推断 |
| goal-tree / workspace.md | 同步 `GOAL-010` 行与概述 |
