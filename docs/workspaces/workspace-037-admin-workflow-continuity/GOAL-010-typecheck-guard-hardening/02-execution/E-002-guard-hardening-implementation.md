---
id: E-002-guard-hardening-implementation
doc: execution-entry
status: recorded
goal_id: GOAL-010-typecheck-guard-hardening
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-002 · 守卫加固实现与合成变异用例（C2）

## 事实

2026-09-18，完成 C2：只改 `apps/web/src/typecheck-convention.guard.test.ts`（纯测试文件，+239/−14），未改任何产品代码、`package.json` 脚本或 CI 配置。

## 实现要点

1. **判定改为按配置内容**：新增 `projectTargetSelectsSources(target)`——`-p` 目标解析为路径（目录按 `tsc` 语义落到 `./tsconfig.json`）后读取 JSON，用既有 `selectsOwnSources()` 判定；读不到或解析失败返回 `false`（fail closed）。
2. **参数解析**：新增 `tokenizeArgs` + `projectTargetsFrom` + `tscCallTokens`，支持 `-p x`、`-p=x`、`--project x`、`--project=x`，并在 `&&`/`||`/`;`/`|` 处截断，避免把后续命令的 `-p` 算到当前调用头上。
   - **引号处理按位置判定**：remainder 从 `tsc` 之后开始（行中段），因此「前一个字符是分隔符」的引号才是开引号，紧贴前字符的引号视为收尾引号并丢弃。此规则由 `["tsc", "-p", "x"]`（JSON 数组形态）与 `"typecheck": "tsc -b && tsc -p e2e/tsconfig.json",`（JSON 值形态）两类真实行确定。
   - 反引号区域内的 `${…}` 插值（含 `"`、`)`、`]`）保持为一个 token。
3. **非字面目标**：`scripts/build-lib-packages.mjs` 以 `` `tsconfig.${name.split("/")[1]}.json` `` 逐包构建（`@schema-ui/protocol|lib|ui`）。新增 `dynamicTargetSelectsSources`：把插值换成通配后与包内全部 `tsconfig*.json` 匹配，要求「至少命中一个，且命中项全部选择源文件」；否则 fail closed。当前通配命中 8 个 `tsconfig.*.json`，全部有效 → 该调用判为真实检查（不是误报）。
4. **合成变异用例表**：新增 `COMMAND_FORMS`（16 行）与断言，逐条固定「必须判空转」与「必须判有效」两侧，含 `-p tsconfig.json`、`--project tsconfig.json`、`--project=tsconfig.json`、`-p .`、`-p tsconfig.missing.json`、`-p`（缺目标）、`npm exec --` 前缀，以及「散文提及 `tsc` 不算调用」。
5. **另加定位断言**：`projectTargetSelectsSources` 的 7 条单元断言（根配置/`.`/`./tsconfig.json`/不存在/空串 → false；`tsconfig.app.json`/`e2e/tsconfig.json` → true）。
6. **既有语义未弱化**：`typecheck` 与 `build` 脚本改由同一个行级检查器断言（`nonCheckingTscOnLine(...) === []`），比原先的 `CHECKING_FLAG.test()` 更强；根配置仍须 solution-style、`typecheck` 仍须显式覆盖 e2e 项目、README 约定断言不变。

## 实现过程中的自查（已修正，不是遗留缺陷）

- 初版 tokenizer 把「rc 中的收尾引号」当开引号，导致 `"typecheck": "...e2e/tsconfig.json",` 的目标被吞成 `e2e/tsconfig.json,` → 误报；改为按位置判定引号。
- 初版通配判定复用了一个「按已转义字符串设计」的正则，对原始字符串恒不命中 → `-p` 动态目标被误判为空转；改为先转义再替换插值，并以「替换结果是否变化」判定是否含插值。
- 两者均由守卫自身用例暴露（`uses no non-checking tsc invocation in any executable surface` 报出真实行号），修正后 8/8 通过。

## 测试结果

`npx vitest run src/typecheck-convention.guard.test.ts` → **8 passed**（含新增 2 项：`rejects every vacuous command form...`、`resolves -p targets by config content...`）。C2 完成；`I-010-002` 为 `verified`。C3 见 `E-003`。
