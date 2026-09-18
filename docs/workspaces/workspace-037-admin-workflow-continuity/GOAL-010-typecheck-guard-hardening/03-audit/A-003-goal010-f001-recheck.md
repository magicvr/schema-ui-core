---
id: A-003-goal010-f001-recheck
doc: audit-entry
status: recorded
goal_id: GOAL-010-typecheck-guard-hardening
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# A-003 · GOAL-010 A-002 F-001/F-002 finding-closure 复审（2026-09-18）

- **source**：independent
- **日期**：2026-09-18
- **auditor**：grok-build (grok-4.6 · reasoning xhigh)
- **类型** / **scope**：`GOAL-010 A-002 F-001/F-002` 的 finding-closure 复审（不是全量重审 C1～C3）
- **verdict**：`pass`
- **audit_type**：finding-closure

## 范围与区间

- 工作区：`workspace-037-admin-workflow-continuity`；`workspace.md` 的 `root_goal` = `GOAL-001-admin-workflow-continuity`，`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`（无资料引用可当事实或关闭证据）。
- 本目标为非纲领整改子目标。本条不审 Root / VP / `R5-I-004`，不改任何 `status` / `progress` / 方案正文。
- 信息项：`I-010-001` / `I-010-002` / `I-010-003` 均为 required 且已 `verified`；本 scope 无到期未关闭的 required 信息门禁。
- 关闭声明来源：`02-execution/E-004-closeout-and-projection.md`（编排器响应）。本条独立复测，不采信该声明。

## 实际执行的命令与观察

未运行 `npm run build`，未 `git commit`。N3 变异在复测后已还原。

### 1. 工作树 vs `HEAD`（改动边界）

| 命令 | 观察 |
|------|------|
| `git status --porcelain` | 空（复审开始时与 N3 还原后均空） |
| `git diff --stat HEAD` / `git diff --numstat HEAD` | 空 |
| `git diff --check` | 通过 |
| `git rev-parse HEAD` | `51a262f74fc4dbb59d34066536d16c65364a82d8` |
| `git show --stat HEAD` | 4 文件：`apps/web/src/typecheck-convention.guard.test.ts` **+16**；`E-004-closeout-and-projection.md` +48；`03-audit.md` +1/−1；`A-002-goal010-independent.md` +191 |
| `git diff --stat cc33da25 HEAD -- apps/web/src/typecheck-convention.guard.test.ts` | 仅该文件 **+16** |
| `git diff --name-only cc33da25 HEAD -- apps/web .github`（排除守卫测试） | 空 |

**判定**：当前工作树相对 `HEAD` **没有**未提交差异，因此「工作树 vs HEAD 仅测试文件 +16」这一表述**不成立**——+16 行已在 `51a262f7` 落盘。该提交的**产品/CI 边界**成立：`.github/**`、`apps/web/package.json`、产品代码无改动；测试文件相对加固提交 `cc33da25` 确为 +16。同提交还写入了本目标治理文件（`A-002`、`E-004`、索引），不构成产品越界。

守卫文件 blob：`git hash-object` = `HEAD:apps/web/src/typecheck-convention.guard.test.ts` = `bf2071069db5fb10539b3fc9e4dbc6db2fc363c2`（N3 还原后相同）。

### 2. 守卫 / typecheck / 全量测试（还原前，干净树）

在 `apps/web` 下：

| 命令 | 观察 |
|------|------|
| `npx vitest run src/typecheck-convention.guard.test.ts` | **9 passed / 9**，exit 0 |
| `npm run typecheck`（`tsc -b && tsc -p e2e/tsconfig.json`） | **TYPECHECK_EXIT=0** |
| `npm test`（`vitest run`） | **Test Files 112 passed (112) / Tests 1429 passed (1429)**，`NPM_TEST_EXIT=0`（相对 `A-002` 的 1428，+1 即本轮新增 `it`） |

### 3. 重放 N3（必须还原）

将 `dynamicTargetSelectsSources` 的

```ts
return (
  matches.length > 0 &&
  matches.every((relativePath) => projectTargetSelectsSources(relativePath))
);
```

改为 `return matches.length > 0;`，再跑 `npx vitest run src/typecheck-convention.guard.test.ts`：

- **2 failed / 7 passed (9)**，`N3_GUARD_EXIT=1`
- 失败用例 1：`rejects every vacuous command form, including -p at the root config`
  - 误判数组**只有一行**：`tsc --noEmit -p $CFG -> expected vacuous=true`
  - `tsc --noEmit --project ${CFG}` **未**出现在误判数组（该表行在无 `every()` 时仍被判空转，原因见 §4）
- 失败用例 2：`fails closed when an interpolated -p target could denote the root config`
  - 第一句 `expect(dynamicTargetSelectsSources("$CFG")).toBe(false)` 得到 `true`；后续 `${CFG}` / `./$CFG` 断言因失败中止未执行

随后把 `every(...)` 写回。`git hash-object` 回到 `bf2071069db5fb10539b3fc9e4dbc6db2fc363c2`（= `HEAD`）；`git status --porcelain` 空。还原后再跑守卫：**9 passed / 9**。

### 4. 新增用例是否覆盖 `every` 分支（独立命中集，非采信）

用与守卫相同的 `escapeRegExp` / `INTERPOLATION` / `allTsConfigs` / `selectsOwnSources` 在仓库外脚本复算（`apps/web` 10 个 `tsconfig*.json`：9 个顶层 + `e2e/tsconfig.json`）。根配置 `tsconfig.json` 为 `files: []`、无 `include`，`selectsOwnSources === false`；其余 8 个顶层 `tsconfig.*.json` 与 `e2e/tsconfig.json` 均有非空 `include`。

| 目标 | 通配 | 命中集 | `any`（`length>0`） | `every`（全部选择源） |
|------|------|--------|---------------------|------------------------|
| `$CFG` | `^[^/]*$` | 9 个顶层配置，**含** `tsconfig.json` | true | **false**（根配置拖垮） |
| `${CFG}` | 同上 | 同上，**含** `tsconfig.json` | true | **false** |
| `./$CFG` | `^\./[^/]*$` | **空**（`relative()` 路径无 `./` 前缀） | **false** | false |
| `tsconfig.${name.split("/")[1]}.json` | `^tsconfig\.[^/]*\.json$` | 8 个 `tsconfig.*.json`，**不含**根 `tsconfig.json` | true | true |

旧逻辑（无 `every`，只要求命中至少一个）会漏判的原因：`$CFG` / `${CFG}` 的命中集既含 solution-style 根配置，也含真实项目配置。`any === true` 会把整值通配当成「真实检查」，于是 `tsc --noEmit -p $CFG` 不再被判空转——这正是 `every` 存在的理由。新用例里 **`$CFG` 的表行 + 对 `dynamicTargetSelectsSources("$CFG")` 的直接断言**会在删掉 `every` 后失败，因此**不是碰巧通过**。

两点精度（不构成新 finding，见残余）：

1. `COMMAND_FORMS` 的 `` `tsc --noEmit --project ${CFG}` `` 在扫描路径上**没有**钉住 `every`。`ARG_SEPARATOR` 含 `{` / `}`，该命令被拆成 `--project` + `$` + `CFG`，取出的目标是 `$`；`dynamicTargetSelectsSources("$")` 因「无插值标识符」直接 `false`，无 `every` 时该表行仍判空转。函数级 `dynamicTargetSelectsSources("${CFG}")` 才会走到混合命中集。
2. `./$CFG` 命中集为空，无 `every` 时仍为 `false`（走的是 `matches.length > 0`，不是 `every`）。

### 5. `02-execution.md` 索引 vs `02-execution/` 目录

| E-ID | 索引状态 | 文件存在 | 文件 `status` |
|------|----------|----------|---------------|
| E-001 | recorded | `E-001-vacuity-boundary-and-rule-baseline.md` | recorded |
| E-002 | recorded | `E-002-guard-hardening-implementation.md` | recorded |
| E-003 | recorded | `E-003-mutation-verification-and-regression.md` | recorded |
| E-004 | recorded | `E-004-closeout-and-projection.md` | recorded |

四个文件均存在，索引行与目录、frontmatter 状态相符。`A-002 F-002` 所描述的「索引有 E-004、目录无文件」已不成立。

## 成果（有证据）

- `A-002 F-001` 的闭合证据可重复：N3 现捕获 2/9；还原后 9/9；命中集证明 `$CFG` 覆盖混合根配置。
- `A-002 F-002` 的闭合证据可重复：`E-004` 已落盘且与索引一致。
- 回归仍绿：typecheck exit 0；Vitest 112 文件 / 1429 用例。
- 产品/CI 边界未破。

## 对照成功标准（本 scope = 关闭复审）

| 标准 | 状态 | 证据 |
|------|------|------|
| F-001 关闭证据充分可重复 | 达成 | 本条 N3；命中集；9/9 |
| F-002 关闭证据充分可重复 | 达成 | 目录列举 vs `02-execution.md` 索引；`E-004` frontmatter `recorded` |
| 关闭声明未把产品/CI 改动混进来 | 达成 | `cc33da25..HEAD` 对 `apps/web`/`.github` 仅守卫测试 +16 |
| 无到期 required 信息项阻断本 scope | 达成 | `I-010-001`～`003` verified |

## Findings

无新增 finding。

## 必改项汇总

无 open required finding。

## 残余与不确定性

- **工作树表述**：编排器「仅测试文件 +16」描述的是 `51a262f7` 里守卫文件的 delta，不是当前 `git diff HEAD`（工作树干净；同提交含治理文件）。
- **`${CFG}` 表行不经扫描路径覆盖 `every`**：见 §4；函数级断言仍覆盖该形态。当前行为是 fail closed，方向安全。
- **`./$CFG` 不经 `every`**：空命中集。若将来 `relative()` 改为带 `./` 前缀，该断言才会变成混合命中集测试。
- 本条未重做 `A-002` 的仓库内类型错误注入、scratch 退出码、N1/N2。那些不在本 finding-closure scope。
- `00-meta.md` 的 C4 仍未勾选；本意见不改检查点。`E-004` 中「复审通过后」的投影（`GOAL-010` → `done · 4/4`、`GOAL-008 A-002 F-001` → `fixed`）须由 `/govern` 在用户确认后执行。

## 与既有意见的异同

- `A-002`（independent，pass）：本条同意其 F-001/F-002 的原判定，并确认编排器响应后两条均可按 P-003 `fixed` 闭合。
- 差异：`E-004` 把两条 `COMMAND_FORMS` 都写成 N3 捕获面；本条实测 N3 的表驱动失败**只有** `-p $CFG`。不构成 P-004 冲突（同向、建议级精度，不改变可闭合结论）。
- `A-001`（self，pass）：无冲突。

## 结论 + 建议给编排器 / 用户的下一步

verdict **`pass`**。开放 required = 0。

- `A-002 F-001` → **fixed**
- `A-002 F-002` → **fixed**
- 无新增 finding

建议 `/govern`：按 P-003 将上述两条标 `fixed` 并留痕；用户确认无未闭合 required 后勾 C4、投影关门。不要用本条改 `status`/`progress`。`${CFG}` 表行与 `./$CFG` 的精度可作后续测试整理，不阻断本目标关门。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
