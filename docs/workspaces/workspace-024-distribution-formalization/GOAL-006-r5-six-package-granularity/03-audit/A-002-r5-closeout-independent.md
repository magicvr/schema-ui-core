---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-006-r5-six-package-granularity
version: 0.1.0
---

# A-002 · GOAL-006 关门独立审计（source: independent · 2026-08-29）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · GOAL-006（C1–C5 · E-002 + `attachments/freeze-face-v1.4.0.md` · D-001 落实度 · 残余登记 · 版本修正链口径）
- **verdict**：**conditional**
- **工作区**：`workspace-024-distribution-formalization`（Root `GOAL-001-distribution-formalization` · `canonical_scope` 本区 · `shared_materials_catalog: none`）

## 范围与区间

核对 `00-meta` 成功标准 C1–C5 是否被 E-002 / 冻结面 v1.4.0 / npmjs 实发 tarball / golden-field 五探针 **可重复**支持；并核 D-001 重写表、I-001/I-002、shell 类型面残余、版本修正链。不改 `status` / `progress` / 方案正文 / goal-tree。

下游仓 `github.com/magicvr/golden-field` 仅作为 workspace.md 写明的实验消费实证对象读取，不引入其他工作区目标状态。VP-024 判据 #5/#6 以愿景层文本为准（`docs/vision/plans/VP-024-distribution-formalization.md`）。

## 范围与区间（覆盖 / 排除）

- **覆盖**：C1 renderer external 化产物与 peer 声明；C2 exports / tsc 子路径 / 终版发布；C3 ui 纯原子与独立消费；C4 冻结面 v1.4.0；C5 五探针 + 无凭据安装；shell 类型残余；版本修正链。
- **排除**：不重跑六包构建/发布；不改 dist-lib；不以 `progress: 0/5` 或检查点勾选作为完成证据。

## 成果（有证据）

| 项 | 独立核对 | 证据路径 |
|----|----------|----------|
| C1 JS external 化 | npmjs `renderer@0.3.7` tarball `index.js` = **187 750 B**（≈187.8 kB SI；E-002 写作 187.5 kB）· **17** 处 `from "@magicvr/schema-ui-*"`（protocol / lib / ui 子路径，含扩展名）· 产物 js+d.ts **0** 处 `@/` | tarball `magicvr-schema-ui-renderer-0.3.7.tgz`；本仓 `apps/web/dist-lib/@schema-ui/renderer/index.js` 同字节 |
| C1 子路径可解析 | 17 处 import 对应的 protocol/lib/ui 文件在各自 tarball **均存在**（如 `protocol/protocol/conformance/component-format.js`、`lib/i18n/runtime.js`、`ui/components/ui/card.js`） | npm pack 解包核对 |
| C2 exports `./` 合法 | 六包终值（及 shell/theme 0.1.2）`exports["."]` 的 `types`/`import` 均为 `./` 前缀；`"./*": "./*"` 通配存在 | `npm view` + tarball `package.json` |
| C2 tsc 子路径 | protocol **0.2.11** / lib **0.1.9** / ui **0.1.7** 为 js+d.ts 树（非单 bundle） | 同 tarball 目录树 |
| C2 npmjs 实发 | dist-tags：protocol `0.2.11` · lib `0.1.9` · renderer `0.3.7` · ui `0.1.7`；修正链时间戳均在 2026-08-29（renderer `0.2.0`→`0.3.7`，protocol `0.2.0`→`0.2.11` 等） | `npm view … time dist-tags` |
| C3 `components/ui` 文件面 | 9 个原子组件文件 + `index`（源码目录 12 文件含 2 测试）：async-state / badge / breadcrumbs / button / card / input / label / skeleton / textarea；**无** renderer/protocol import | `apps/web/src/components/ui/`；`dist-lib/@schema-ui/ui/components/ui/` |
| C5 五探针（现装） | 全绿：probe-r5 `imports=17` · probe protocol **2.9** · probe-render **1573 B** · probe-six PASS · token-check `brand=2 ⊆ index=5` | `C:\Users\magicvr\Documents\Code\golden-field\web\` |
| C5 无凭据隔离安装 | 空 userconfig + 独立 cache + `--registry https://registry.npmjs.org`：`npm install` 20 packages / exit 0；lock `resolved` 全部 `https://registry.npmjs.org/@magicvr/…tgz`；五探针再跑仍全绿 | `%TEMP%\r5-audit-iso\app`（2026-08-29 本审计） |
| C5 lockfile vs GH | golden-field `pnpm-lock.yaml` 无 `npm.pkg.github` / 无 tarball URL；integrity 与 npmjs dist 一致（含 shell/theme **0.1.2**） | lockfile L38–54；`npm view dist.integrity` |
| 残余（运行时） | shell 类型面 `@/account/*`、`@/host/*` 仍在；JS 五探针绿，运行时自包含不受影响 | 见 F-005；冻结面 §3 / E-002 残余 1 已登记 R7 |
| I-001 内容 | 源码 `src/renderer` 对 `@/{i18n,protocol,components/ui,lib}` 的 41 处前缀清单，构建后落入 17 个唯一包子路径；renderer 产物无 `@/` | 本审计扫描；I-001 **台账仍 open**（见 F-006） |
| 共享资料 | `none`，无引用被当成证据 | `workspace.md` |

## 对照成功标准

| 标准 | 状态 | 独立证据 | 缺口 |
|------|------|----------|------|
| C1 renderer external 化（体积 + 17 import + 0 `@/` + **peerDependencies**） | **部分达成** | JS/体积/import/`@/` 清零可重复 | **已发布 `renderer@0.3.7` 无 `peerDependencies` 字段**（C1 已勾选「peer 声明」）。见 F-001。入口 `types: ./index.d.ts` 在 tarball **不存在**。见 F-002 |
| C2 exports 子路径 + tsc 产物 + 终版 npmjs | **基本达成** | `./` 合法；三包 tsc 子路径文件齐；protocol/lib/renderer/ui 的 latest = 冻结面终值 | C2/冻结面终值 **shell/theme = 0.1.2**，但 npm `latest` = **0.1.3**（且 0.1.3 丢掉 `files` 与 shell `peerDependencies`）。见 F-004。D-001「files 收窄」已被 C2/冻结面改写为「全量」，决策未回写 |
| C3 ui 纯原子 + 独立消费 | **部分达成** | `components/ui` 无 renderer/protocol 反向依赖；probe-six 在**六包齐装**下可消费 DataTable | breadcrumbs 仍 import i18n（D-001 明文禁止）；`DataTable` 仍从 ui 入口导出（ui tarball 含 `components/data-table.js`）；**仅装 ui** 在 `badge.js` 即 `ERR_MODULE_NOT_FOUND @magicvr/schema-ui-lib`。见 F-003 |
| C4 冻结面 v1.4.0 定稿 | **部分达成** | 附件存在，记录导出面 / 版本终值 / shell 残余 / 0.2.0→0.3.0 契约变化 | peer 矩阵与 registry **不一致**；renderer 入口 types 路径与冻结面、与 tarball 三者互不一致；shell/theme 终值与 `latest` 不一致。见 F-001/F-002/F-004 |
| C5 golden-field 五探针 + 无凭据 + 独立审计 | **探针/安装达成；关门未放行** | 现装 + 本审计隔离 npmjs 安装五探针全绿；`.npmrc` 钉 `@magicvr:registry=https://registry.npmjs.org` | 本条即独立意见。C1/C3/C4 缺口未闭合前不得把 C5「关门」视为完成 |

## Findings

### F-001 · 已发布包未声明 peer 矩阵，C1/C4「peerDependencies 声明」不实

- 严重度：high
- 建议：**required**
- 状态：open
- 描述：C1（已勾）、E-002、A-001、冻结面 v1.4.0 §2 均主张 renderer peer = `react` / `react-dom` / `@magicvr/schema-ui-{protocol,lib,ui}`，ui/shell peer = react 系，lib peer = react。独立核对：
  1. `npm view @magicvr/schema-ui-renderer@0.3.7 peerDependencies` **空**；tarball `package.json` 无该字段，亦无 `dependencies` 指向三包子面。
  2. `ui@0.1.7`、`lib@0.1.9` 同样无 peer。`shell@0.1.2` **有** react peer（golden-field 所钉）；`shell@0.1.3`（npm latest）又丢掉。
  3. 根因可复现：`scripts/rewrite-lib-aliases.mjs` 的 `peers`/`versions` 以短名 `"renderer"` 为键，写入时却用 `peers[old.name]`（`@magicvr/schema-ui-renderer`）查找 → **永远 miss**，peer 不会落进 `package.json`。
  4. 本审计在隔离目录只装 `renderer@0.3.7` + react：npm 只装入 `@magicvr/schema-ui-renderer`（不自动拉取 protocol/lib/ui）。这正是 peer 字段要防止的「静默缺面」。
- 证据：tarball `package.json`；`scripts/rewrite-lib-aliases.mjs` L42–46、L104；`%TEMP%\r5-audit-renderer-only` 安装图；冻结面 §2；`00-meta` C1 勾选句。
- 关闭条件：修正脚本查找键（短名或全名一致）后重建；在 **新版本** 写入与冻结面一致的 peer 矩阵并 npmjs 实发；冻结面版本号与 `latest` 对齐。在此之前不得把 C1/C4 的 peer 子句视为已核销，不得 GOAL-006 `done`。

### F-002 · renderer 0.3.7 入口 types 指向不存在的 `./index.d.ts`

- 严重度：high
- 建议：**required**
- 状态：open
- 描述：发布物 `package.json`：`"types": "./index.d.ts"` 且 `exports["."].types` 同值。tarball 清单仅有 `package/renderer/index.d.ts`，**根路径 `package/index.d.ts` 不存在**。冻结面写入口 types = `./renderer/index.d.ts`，与 package.json、与实文件布局三者不一致。probe-r5 按文件系统读 `renderer/index.d.ts`，**不会**发现包入口 types 断裂。另：4 份 renderer d.ts 把协议面写成 `@magicvr/schema-ui-protocol/app-manifest.js` / `…/conformance/…`（缺 `protocol/` 段），而 JS 与 tsc 产物实际在 `protocol/protocol/…`。
- 证据：`npm pack @magicvr/schema-ui-renderer@0.3.7` 文件清单；`renderer/renderer/{row-action,render,permissions,form-controls}.d.ts`；冻结面 §1 renderer 行。
- 关闭条件：二选一并发布新版本：(a) 增加根 `index.d.ts` re-export `./renderer/index.d.ts`；或 (b) 把 `types` / `exports["."].types` 改为 `./renderer/index.d.ts`。同步修正 d.ts 协议子路径并回写冻结面。probe-r5 应改为读 `package.json` 入口 types 而非硬编码相对路径。

### F-003 · C3「纯原子 / 独立消费」声明过宽

- 严重度：med
- 建议：**required**
- 状态：open
- 描述：D-001 / C3 要求 `components/ui` **禁止** import `@/i18n`（及 renderer/protocol），且「业务组件出 ui 包」+ ui 独立消费。独立核对：
  1. `breadcrumbs.tsx` / 产物 `breadcrumbs.js` 仍 `from "@magicvr/schema-ui-lib/i18n/…"`.
  2. ui 包仍含并再导出 `components/data-table.js`；`src/components/ui/index.ts` 注释仍写「原子 + DataTable 核心」。E-002「业务组件（data-table）留在 renderer」与 ui tarball **并存**（renderer 把 data-table **打进** `index.js`，ui 仍独立导出）。
  3. 本审计 **仅装 ui@0.1.7 + react**：`import("@magicvr/schema-ui-ui")` → `ERR_MODULE_NOT_FOUND Cannot find package '@magicvr/schema-ui-lib'`（`badge.js` 的 `lib/utils`）。ui 未把 lib 列为 dependency/peer。probe-six 是六包齐装，不能当作「ui 包独立消费」。
  4. 「12 原子组件」实际是源码目录 12 **文件**（9 组件 + index + 2 测试），不是 12 个组件。
- 证据：`apps/web/src/components/ui/breadcrumbs.tsx` L21–22；ui tarball `components/data-table.js` + `components/ui/index.js` L13；`%TEMP%\r5-audit-ui-only`；D-001 条 4；VP-024 判据 #6。
- 关闭条件：三选或组合并留痕：(1) breadcrumbs 改为相对/自含或把 lib 列为 ui peer，并修正 D-001「禁止 i18n」；(2) DataTable 从 ui 公共入口移除（或书面修订 #6 / D-001：data-table 保留为 ui 设计系统面）；(3) 用「只装 ui（+ 已声明 peer）」重跑独立消费探针。在修订或修正落地前不得把 C3 / 判据 #6 标为已核销。

### F-004 · 冻结面终值 shell/theme 0.1.2 与 npm `latest` 0.1.3 分叉

- 严重度：med
- 建议：**required**
- 状态：open
- 描述：C2 / E-002 / 冻结面 §1 写终版 shell/theme **0.1.2**；golden-field 亦钉 0.1.2（integrity 与 npmjs 0.1.2 一致）。但 `npm view dist-tags.latest` = **0.1.3**。0.1.2 有 `files` 白名单且 shell 有 react peer；0.1.3 两字段皆无。本仓 `dist-lib` 本地 package.json 已是 0.1.3。未钉版本的 `npm i @magicvr/schema-ui-shell` 会装到比冻结面更弱的契约。版本修正链本身（0.3.0–0.3.7 / 0.2.1–0.2.11 / 0.1.1–0.1.9）时间戳可核对，问题是 **「终值」未钉在 latest**。
- 证据：`npm view @magicvr/schema-ui-shell dist-tags time`；两版本 tarball `package.json`；golden-field `package.json` L12–13；冻结面 §1；本地 `dist-lib/@schema-ui/{shell,theme}/package.json`。
- 关闭条件：冻结面、C2、golden-field、npm `latest` 四者同一组版本。可选：`npm dist-tag` 把 latest 收回 0.1.2 并 deprecate 0.1.3；或发布 0.1.4（恢复 files/peers）并升冻结面终值。

### F-005 · shell 类型面残余登记成立，但「7 处含 data-table」计数口径不准

- 严重度：low
- 建议：**recommended**
- 状态：open
- 描述：E-002 / 冻结面 §3 / A-001 R-001 将残余登记为 R7 复核，方向正确：JS 运行时自包含、五探针绿、消费端 tsc 类型面未验证。独立计数（0.1.2 / 0.1.3 / 本地 dist-lib 相同）：**7 处 `@/` 出现**，分布在 4 个文件（`account/AuthContext.d.ts` 3、`app/HostFailureScreen.d.ts` 1、`app/LoginPage.d.ts` 1、`host/boot.d.ts` 2），目标仅为 `@/account/*` 与 `@/host/*`。`components/data-table.d.ts` **没有** `@/` import（是一份自含类型副本）。把 data-table 算进「7 处 `@/`」会误导 R7 去修一个不存在的 alias。
- 证据：本审计对三份 shell 产物的 `@/[a-z]` 扫描。
- 关闭条件：冻结面 §3 / E-002 改为「4 文件 / 7 处 `@/account|@/host`；data-table.d.ts 无 `@/`」。R7 范围维持：消费端 tsc 类型面。

### F-006 · 信息项与台账未随 S2/S3 闭合；A-001 曾未入索引

- 严重度：low
- 建议：**recommended**
- 状态：open
- 描述：
  - I-001 required（最晚 S2）在 `00-meta` / 本索引写入前仍 `open`「待确认」，而 S2 已声称关门。本审计认为 **内容已满足**（41 处前缀 → 17 个唯一外部 import、renderer 无 `@/`），缺的是台账闭合。
  - I-002 non-blocking：包内无 CHANGELOG；冻结面 §3 已有「0.2.0 自包含 → 0.3.0 external 化」迁移句，可作为注记，但 `00-meta` 仍 open。
  - 本条写入前 `03-audit.md` 索引为空（A-001 文件已在）；`progress: 0/5` 与 C1–C5 已勾并存（独立审不改 progress）。
  - D-001 映射表仍写 `@schema-ui/*`，实发 `@magicvr/schema-ui-*`（构建脚本注释同样）。
- 证据：`00-meta` 信息表；本条写入前的 `03-audit.md`；D-001 条 1。
- 关闭条件：`/govern` 按证据把 I-001 标 verified、I-002 指向冻结面或补 changelog；D-001 补发布全名；按真正闭合的检查点重算 progress（不得在 F-001～F-004 开放时把 C1/C3/C4 勾成完成）。

## 必改项汇总

| id | 严重度 | 阻断 |
|----|--------|------|
| F-001 | high required | C1 peer 子句、C4 peer 矩阵、VP-024 判据 #5「peer 矩阵定稿」 |
| F-002 | high required | C1/C4 类型入口 / 冻结面导出面 |
| F-003 | med required | C3、VP-024 判据 #6 |
| F-004 | med required | C2/C4 终值口径 vs npm `latest` |

F-005 / F-006 为 recommended，不单独阻断，但 F-006 的检查点勾选应在 required 闭合后才算数。

## 与既有意见的异同

| 项 | A-001 self | 本条 independent |
|----|------------|------------------|
| verdict | conditional（self 侧 pass；待 A-002） | **conditional**（运行时/探针可重复；契约面未定稿） |
| C1 JS/17 import/0 `@/` | ✅ | **同意**（187 750 B、17、0 `@/`） |
| C1 peer 声明 | ✅（未拆开核对 registry） | **不同意** → F-001 required |
| C2 终版发布 / exports / tsc | ✅ | **基本同意**；shell/theme latest 分叉 → F-004 |
| C3 纯原子 + 独立消费 | ✅ | **部分不同意** → F-003（i18n / DataTable / 仅装 ui 失败） |
| C4 冻结面 | ✅ | **部分不同意**（附件有，但与 registry/types 不一致） |
| C5 五探针 | ✅ | **同意**（现装 + 隔离 npmjs 重跑全绿） |
| 残余 shell 类型面 | R-001 recommended → R7 | **同意方向**；计数口径 → F-005 recommended |
| 版本修正链 | R-002 recommended | **同意链存在**；终值 vs latest → F-004 required |
| required | 0 | **4**（F-001～F-004） |

与 A-001 无「一要一否」的 P-004 意见冲突（双方均为 conditional、均不主张立即 `done`），但本条新增 4 条 required：编排器不得在未闭合前把 GOAL-006 标 `done` 或核销判据 #5/#6。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional**。R5 的**运行时外部化**是真的：renderer 从自包含体积降到 187 750 B、17 处 `@magicvr/schema-ui-*` import、renderer 产物无 `@/`、protocol/lib/ui tsc 子路径齐、exports `./` 合法、golden-field 五探针在无凭据 npmjs 安装下全绿。shell `@/account|@/host` 类型残余已诚实登记 R7，方向可接受。

不能无条件关门的原因是 **消费契约面（C1 peer + C4 冻结面 + C3 独立消费）尚未与 registry 一致**：peer 矩阵只写在附件、没写进已发布 `package.json`；renderer 入口 types 指向缺失文件；ui 不能单独安装；冻结面终值与 `latest` 分叉。VP-024 判据 #5「六包 peer 矩阵定稿」与判据 #6「业务组件出 ui 包 / ui 可独立消费」在本审计下 **内容未核销**。

建议 `/govern`：响应 A-001 + 本条 → 先修 F-001～F-004（脚本键、peer 实发、types 入口、ui 边界或书面修订 D-001、latest 对齐）→ 再决定是否复审 A-002 或追加 A-003。不要用已勾的 C1–C5 或 `progress` 百分比代替闭合。I-001 可用本条证据标 verified；I-002 可用冻结面 §3 闭合或补包内 changelog。

## 声明

本意见不修改 status / progress / 方案正文 / goal-tree 状态列；响应由 `/govern` 处理。

---

## 响应（2026-08-29 · /govern · source: self）

- **F-001（peer 未实发）→ fixed**：根因 = rewrite 脚本 `peers[old.name]` 键错位（短名表 vs 全名键）；修复 = 元数据注入 + 新版本实发——renderer **0.3.8**（peer = react/react-dom/protocol/lib/ui）· lib **0.1.10**（react）· ui **0.1.8**（+lib peer）· shell **0.1.4** · theme **0.1.4**（latest 齐平）。
- **F-002（renderer types 指向缺失文件 + d.ts 协议段）→ fixed**：types/exports → `./renderer/index.d.ts`（0.3.8 内含）；renderer 4 份 d.ts 协议引用补 `protocol/` 段（form-controls/permissions/render/row-action）。
- **F-003（ui 边界）→ fixed（用户 P-004 裁决 · 2026-08-29）**：判据 #6 口径修订——**data-table 属 ui 设计系统面**（留在 ui 包；业务组件 form-controls 等留在 renderer）；breadcrumbs 的 i18n 经 `@magicvr/schema-ui-lib` 包子路径（ui peer 已声明 lib）；ui 独立消费 = 仅装 ui + peer → **UI-ONLY 实测 PASS（exports=18）**。D-001 §7 修订注记。
- **F-004（shell/theme latest 分叉）→ fixed**：终值升 0.1.4（恢复 peers 契约）实发 → 冻结面/golden-field/latest 四者齐平。
- **F-005 → fixed**（E-002/冻结面残余口径改「4 文件 / 7 处 `@/account|@/host`；data-table.d.ts 无 `@/`」）。
- **F-006 → fixed**（I-001 verified（41 前缀 → 17 唯一包子路径 · 产物零 `@/`）· I-002 指向冻结面 §3 迁移句 · 03-audit 索引齐全 · progress 5/5（检查点全闭）· D-001 补发布全名注记）。

golden-field 五探针 + UI-ONLY 全绿（无凭据 npmjs）。全部 required 闭合 → **GOAL-006 done 5/5 · Root 5/7（判据 #5/#6 核销）**。
