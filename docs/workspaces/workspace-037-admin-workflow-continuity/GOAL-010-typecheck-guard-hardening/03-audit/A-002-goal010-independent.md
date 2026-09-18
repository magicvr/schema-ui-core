---
id: A-002-goal010-independent
doc: audit-entry
status: recorded
goal_id: GOAL-010-typecheck-guard-hardening
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# A-002 · GOAL-010 独立交叉审计（加固正确性 / 非空转 / 回归）

- **source**：independent
- **日期**：2026-09-18
- **auditor**：grok-build (grok-4.6 · reasoning xhigh)
- **类型** / **scope**：`GOAL-010-typecheck-guard-hardening` 的 C1～C3（空转事实、`-p` 目标判定、非空转变异、误报面、改动边界与回归）；承接同工作区 `GOAL-008` `A-002` `F-001`
- **verdict**：`pass`
- **audit_type**：close-out / finding-closure（对 `GOAL-008 A-002 F-001` 的整改核验）

## 范围与区间

- 工作区：`workspace-037-admin-workflow-continuity`；`workspace.md` 的 `root_goal` = `GOAL-001-admin-workflow-continuity`，`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`（无资料引用可当事实）。
- 本目标为非纲领整改子目标；不审 Root / VP / `R5-I-004`，不改任何 `status` / `progress` / 方案正文。
- 信息项：`I-010-001` / `I-010-002` / `I-010-003` 均为 required 且声称 verified；本次独立复测 C1/C2 证据，C4 的 provider 即本条。无到期未关闭的 required 信息门禁。
- 其余工作区文档中的 `tsc --noEmit` 历史注记不在范围。

## 实际执行的命令与观察

工作树在注入/变异前后均用 `git hash-object` 与 `HEAD` 对照；最终 `git status --porcelain` 为空，`git diff --check` 通过。未运行 `npm run build`，未 `git commit`。

### 0. 改动边界

| 命令 | 观察 |
|------|------|
| `git show --stat --format=fuller cc33da25` | 仅 `apps/web/src/typecheck-convention.guard.test.ts`，**+239 / −14** |
| `git show --name-only cc33da25` | 同一文件 |
| `git diff --stat cc33da25 HEAD -- apps/web/src/typecheck-convention.guard.test.ts` | 空（HEAD 即该提交的守卫文件） |
| `git show --stat 5ba514c6` | 仅治理文档（GOAL-010 五件套、Root D-016/E-024、goal-tree） |
| 终态 `git diff --stat -- apps/web .github` | 空；`.github/**`、`apps/web/package.json`、产品代码无未声明改动 |

### 1. 空转事实（断言 1）

`apps/web/tsconfig.json` 现为 solution-style：`files: []`，`references` 指向 `tsconfig.app.json` / `tsconfig.node.json`，无 `include`。其余 8 个顶层 `tsconfig.*.json` 与 `e2e/tsconfig.json` 均有非空 `include`；没有任何 `apps/web` tsconfig 使用 `extends`。

**仓库内注入**（与 `E-001` 同一探针）：向 `apps/web/src/renderer/list-surface.tsx` 追加 `const __probeErr: number = "definitely a string";`，在 `apps/web` 下跑四条命令，随后删除探针。文件 `git hash-object` 恢复为 `94d174faab1cce1d5c23ac87dcef1293560761dd`（= `HEAD`）。

| 命令 | 编译器输出（仓库内注入） |
|------|--------------------------|
| `npx tsc --noEmit` | 无 `error TS*` |
| `npx tsc --noEmit -p tsconfig.json` | 无 `error TS*` |
| `npx tsc --noEmit -p tsconfig.app.json` | `TS2322` + `TS6133`（`list-surface.tsx:94`） |
| `npx tsc -b` | `TS2322` + `TS6133`（同上） |

该次 `cmd /c` 未开 delayed expansion，**仓库内四条命令的数字退出码未写入**（见「残余」）。退出码在 `%TEMP%\goal010-tsc-scratch` 用**同构** solution-style 布局 + 同一探针、`apps/web/node_modules/typescript/bin/tsc` 复测：

| 命令 | `STATUS` | `error TS*` |
|------|----------|-------------|
| `tsc --noEmit` | **0** | 0 |
| `tsc --noEmit -p tsconfig.json` | **0** | 0 |
| `tsc --noEmit -p tsconfig.app.json` | **2** | `TS2322` |
| `tsc -b` | **1** | `TS2322` |

与 `E-001`「裸/`-p` 根配置空转；`-p tsconfig.app.json` 与 `-b` 非零」一致。干净树上的 `apps/web` 四条命令在无探针时均为 `STATUS=0`（不构成空转证明，只说明当前树无类型错误）。

### 2. 判定正确性（断言 2）

| 命令 | 结果 |
|------|------|
| `cd apps/web ; npm run typecheck` | **TYPECHECK_EXIT=0**（`tsc -b && tsc -p e2e/tsconfig.json`） |
| `npx vitest run src/typecheck-convention.guard.test.ts` | **8 passed / 8**，`GUARD_EXIT=0` |

**自构反例**（不依赖 `COMMAND_FORMS` 表）：在 `apps/web/scripts/generate-claim.mjs` 第 2 行插入 `npx tsc --noEmit -p tsconfig.json`，再跑守卫：

- 7 项通过；**可执行面扫描失败**
- 精确命中：`apps/web/scripts/generate-claim.mjs:2: npx tsc --noEmit -p tsconfig.json`
- `GUARD_EXIT=1`

随后还原该文件，hash 回到 `6b443a5b1fd7e2024201c6b54e924a028ce9739b`（= `HEAD`）。旧判定（只看 `-p` 令牌）会放行这一行；现判定按根配置 `selectsOwnSources === false` 判违规。**`tsc --noEmit -p tsconfig.json` 会被守卫抓住，且走的是扫描路径，不是只在用例表里生效。**

### 3. 非空转 · 新变异（断言 3；均不在 `E-003` M1–M5）

每次变异后跑 `npx vitest run src/typecheck-convention.guard.test.ts`，再还原守卫文件。终态 hash `8413e791958a00ea1eb5abd541fe413580662c59` = `HEAD`。

| ID | 变异（`E-003` 未列） | 结果 | 失败点 |
|----|----------------------|------|--------|
| **N1** | `projectTargetSelectsSources` 两个 `catch` 从 `return false` 改为 `return true`（fail-closed 取反） | **捕获**（2 failed / 6 passed） | `tsc --noEmit -p tsconfig.missing.json` 不再被判空转；`projectTargetSelectsSources("tsconfig.nope.json")` 期望 `false` 得到 `true` |
| **N2** | `-p` 目标改为**只比文件名**：basename 为 `tsconfig.json` 或 `.` 则否，否则视为有效，**不读配置内容** | **捕获**（4 failed / 4 passed） | `tsconfig.missing.json` 漏判；`e2e/tsconfig.json` 因 basename 也是 `tsconfig.json` 被误杀，连 `package.json` 的 `typecheck` 脚本（`tsc -p e2e/tsconfig.json`）也被可执行面扫描打出 |
| **N3** | `dynamicTargetSelectsSources` 去掉 `matches.every(projectTargetSelectsSources)`，改为「命中至少一个即通过」 | **未捕获**（**8 passed / 8**） | 当前通配 `^tsconfig\.[^/]*\.json$` 命中的 8 个顶层 `tsconfig.*.json` 全部选择源文件，削弱「全匹配」后行为不变 |

N1/N2 证明 fail-closed 与「按配置内容判定」两条承重断言不是空转。N3 见 **F-001**。

### 4. 误报面与漏检形态（断言 4）

干净树上可执行面扫描为空（守卫 8/8 含该项）。抽查：

- `apps/web/package.json`：`typecheck` = `tsc -b && tsc -p e2e/tsconfig.json`；`build` = `tsc -b && vite build`。
- `.github/workflows/r6-basic-matrix.yml`：web job 为 `npm test` + `npm run typecheck`；注释中的 `tsc --noEmit` 被 `isCommentLine` 跳过。
- `scripts/build-lib-packages.mjs:127`：`["tsc", "-p", \`tsconfig.${name.split("/")[1]}.json\`]` 未被打出 —— 插值通配命中 8 个选择源文件的配置，判为真实检查。
- 仓库 `scripts/` 其他 `tsc` 出现处为注释 / `tsc-built` 日志（`TSC_COMMAND` 不匹配 `tsc-built`）；`docker compose -p` 不含 `tsc` 命令词。
- `apps/web/scripts/` 三文件无 `tsc` 调用。

在 `generate-claim.mjs` **额外插入**下列形态（随后全部还原）：

| 行 | 形态 | 守卫 |
|----|------|------|
| 2 | `npx tsc --noEmit -p tsconfig.json` | 打出 |
| 3 | `npx tsc --noEmit -p=./tsconfig.json` | 打出 |
| 4 | `npx tsc --noEmit --project ./` | 打出 |
| 5 | `npx tsc --noEmit -p ./tsconfig.json` | 打出 |
| 6 | `sh -c "tsc --noEmit -p tsconfig.json"` | 打出 |
| 7–8 | `-p $CFG` / `-p "$CFG"` | 打出（fail closed） |
| 9 | `npx tsc -p tsconfig.node.json` | **未打出**（合法） |
| 10 | `npx tsc --noEmit -p tsconfig.app.json` | **未打出**（合法） |
| 11–14 | `--project ./`、`--project=./tsconfig.json`、`-p .`、`--project tsconfig.json` | 打出 |
| 15 | `bash -lc "tsc --noEmit"` | 打出 |
| 16 | `node -e "...execSync('tsc --noEmit -p tsconfig.json')"` | 打出 |

共 13 条违规、2 条合法。**未发现「应当捕获却漏过」的已测形态**；合法 `-p tsconfig.app.json` / `tsconfig.node.json` 与插值构建脚本均未误伤。

fail-closed 取向（不可解析 / 缺目标 / 变量目标 → 违规）与 `D-001` 一致，且当前可执行面没有因此误报。可接受。

### 5. 回归（断言 5）

| 命令 | 结果 |
|------|------|
| `cd apps/web ; npm run typecheck` | exit **0**（见上） |
| `cd apps/web ; npm test` | **Test Files 112 passed (112) / Tests 1428 passed (1428)**，`NPM_TEST_EXIT=0`，含 `src/typecheck-convention.guard.test.ts (8 tests)` |
| `git diff --check` | 通过 |

未跑 `npm run build`（按审计指令）。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| `-p tsconfig.json` 空转是事实 | 达成 | 仓库内注入无 `TS*`；scratch `STATUS=0`；根配置仍 `files: []` |
| 守卫按配置内容判定，该形态违规 | 达成 | 自构反例打出 `generate-claim.mjs:2`；`COMMAND_FORMS` 与 7 条 `projectTargetSelectsSources` 断言 |
| fail-closed 非空转 | 部分 | N1 捕获；动态「全匹配」N3 **未**被捕获（F-001） |
| 不误伤现有可执行面 | 达成 | 干净树扫描空；插值脚本与 `typecheck`/`build` 仍绿 |
| 只改守卫测试文件 | 达成 | `cc33da25` `--stat` |
| 回归绿 | 达成 | typecheck 0；Vitest 112/1428 |
| `I-010-001`～`003` 就本 scope 可维持 verified | 达成 | 本条复测 + 本独立审计即 provider |

## Findings

### F-001 · 动态目标「命中项必须全部选择源文件」没有回归用例

- 严重度：medium
- 建议：recommended
- 状态：open
- 描述：`dynamicTargetSelectsSources` 的公开规则是「通配至少命中一个，且**每一个**命中配置都 `selectsOwnSources`，否则 fail closed」（`E-002` §3、`A-001` 残余、源码注释）。把 `matches.every(...)` 删掉后，守卫 **8/8 仍通过**。当前命中集（8 个顶层 `tsconfig.*.json`）恰好全是真实项目配置，所以削弱看不出来。若将来新增一个 solution-style 的 `tsconfig.<name>.json`，没有 `every()` 时插值调用仍会被放行——这正是该分支存在的理由。
- 证据：本条 N3；`typecheck-convention.guard.test.ts` 的 `COMMAND_FORMS` 与单元断言均不覆盖插值 / 混合命中集。
- 闭合要求：补一条合成断言（例如构造「通配命中含 `tsconfig.json` 或一个 `files: []` 配置」时必须返回 false）；或用户书面 `accepted-residual`（范围=仅当新增 `tsconfig.*.json` 时复核；触发=该文件布局变化）。

### F-002 · 执行索引把尚未落盘的 E-004 标成 recorded

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`02-execution.md` 索引写有 `E-004-closeout-and-projection.md`、状态 `recorded`，但 `02-execution/` 目录只有 `E-001`～`E-003`。`00-meta.md` 的 C4 仍未勾选。这不影响守卫正确性，但会让 C4 投影看起来已经发生。
- 证据：目录列举 vs `02-execution.md` 索引行。
- 闭合要求：删行或改为 planned；真正写 E-004 后再标 recorded。本独立意见不代写该条目。

## 必改项汇总

无 open required finding。

## 残余与不确定性

- **仓库内注入的数字退出码**：四条 `tsc` 的 `TS*` 有无已在仓库内核对；`STATUS=` 数字来自 TEMP 同构 scratch，不是那一次 `apps/web` 进程的 `%ERRORLEVEL%`。与 `E-001`「空转 vs 非零」结论不冲突。
- **未测形态**：YAML 跨行拆开的 argv、`tsc` 经自定义 wrapper 二进制改名、tsconfig JSONC / `extends` 继承 `include`（当前 `apps/web` 无 `extends`、无注释配置）。JSON.parse 失败会 fail closed。
- **扫描边界**（设计如此，非本目标缺口）：不扫治理散文；不扫仓库根 `package.json`（其中也无 `tsc` 脚本）；`GOAL-008 A-002 F-002`（全仓 tsc 简写）仍 open recommended。
- **N2 的附带观察**：若用「文件名是否为 `tsconfig.json`」代替读配置，会误杀合法的 `tsc -p e2e/tsconfig.json`。这支持 `D-001` 不选文件名黑名单。
- 不关闭、不投影 `GOAL-008 A-002 F-001`；按 P-003 的 `fixed` 路径留痕属编排器在 C4 响应时的工作。

## 与既有意见的异同

- `A-001`（self，pass，开放 required = 0）：本条同意 C1～C3 的核心主张（空转事实、判定收紧、扫描路径使用新规则、回归绿、未越界改动）。
- 差异：`A-001` 把插值「全匹配」写成已知残余，但未验证该残余对应的代码是否有测试。本条用 N3 证明**没有**。不构成与 `A-001` 的 P-004 冲突（建议级、同向：保留 fail-closed，补测试）。
- `A-001` 称 `sh -c` / 变量目标会 fail closed、当前仓无此形态：本条在可执行面插入后确认**会被打出**，且干净树无此形态。

## 结论 + 建议给编排器 / 用户的下一步

`GOAL-008 A-002 F-001` 所描述的缺口（`-p` 令牌即视为检查型，导致 `tsc --noEmit -p tsconfig.json` 空转被放行）在守卫实现与扫描路径上**已被修复**，且有独立注入、自构反例、N1/N2 变异与 112/1428 回归支撑。verdict **`pass`**，开放 required = 0。

建议 `/govern`：响应本条 F-001 / F-002（补动态全匹配用例，并修正 E-004 索引）；用户确认无未闭合 required 后勾 C4、投影关门，并按 P-003 将 `GOAL-008 A-002 F-001` 标 `fixed`。不要用本条改 `status`/`progress`。F-001 为 recommended，不阻断关门，除非用户要求先补测再关。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
