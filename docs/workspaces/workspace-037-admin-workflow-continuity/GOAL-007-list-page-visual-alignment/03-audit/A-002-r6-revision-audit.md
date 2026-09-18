---
id: A-002-r6-revision-audit
doc: audit-entry
status: recorded
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
goal_id: GOAL-007-list-page-visual-alignment
source: self
auditor: dsh / deepseek-v4.1-flash（编排器自审）
scope: R6 C5/C7/C8 implementation, regression, and governance projection
verdict: conditional
---

# A-002 · R6 修订审计（C5/C7/C8）

## 头字段

- **source**：self
- **auditor**：dsh / deepseek-v4.1-flash（编排器自审）
- **类型**：`close-out`（修订版阶段关门审计）
- **scope**：R6 `GOAL-007-list-page-visual-alignment` 的 C5 布局修订、C7 控制位纠偏、C8 控件语义与高度，及其回归证据与治理投影
- **verdict**：`conditional`

## 范围与区间

被审区间：`0cc18418`（C5 快照）→ `dc4b5c5b`（C8 实现）→ `d9602b2f`（记录）。核对方式为独立复核，不采信 E-004/E-005/E-006 的结论本身：重新读取实现文件、重跑构建与测试、在真实 Chromium 中重新测量几何与计算样式、逐条对照 D-002/D-003/D-004 与 VP-037/vision 投影。

## 成果（有证据）

1. **C7 搜索配对已真实恢复**。`render.tsx` 重新传入 `searchButtonSlot`，按钮类名含 `-ml-px` 与 `rounded-l-none`；真实 Chromium 实测按钮与输入水平重叠 1px、垂直偏移 0。`form-controls.tsx` 的 `paired`/`rounded-r-none` 机制在 C5 期间未被破坏，C7 复用即可。
2. **C8 高度统一属实**。`schema-table.tsx` 中 toolbar 触发器为 `h-8 … text-xs`，“列配置”触发器为显式 `h-8`；真实 Chromium 在 480/700/900/1280px 实测两者均为 32px。
3. **C8 token 边界修订合规**。`--control`/`--control-foreground` 在 `index.css` `:root`（第 32–33 行）与 `.dark`（第 82–83 行）双态声明，并在 `@theme inline`（第 151–152 行）映射为 `bg-control`/`text-control-foreground`；`theme.test.ts` 已加结构守卫（双层声明 + 别名 + 非自引用）。`git diff` 确认既有 token **未被重命名或重定义**，仅新增。生产构建产物 `dist/assets/*.css` 中确实生成了 `.bg-control` 工具类（实测 YES），即该 token 在发布路径上真实可用，不只是源码层声明。
4. **C8 空展开抑制按预期工作**。真实 Chromium：2 个筛选项时 900/1280px 无展开按键、480/700px 有；5 个筛选项时 1280px 有。判定与可见性类确实共用 `COLLAPSED_SLOTS`（`collapsedVisibility` 经 `slotCounts` 派生）。
5. **回归与构建通过**。`npm exec -- tsc -b` 通过（更正后口径；见 F-005）；全量 Web Vitest 111 文件 / 1420 测试通过；`npm run build`（含 `tsc -b`）exit 0；`git diff --check` 干净。
6. **shell 边界未被 R6 突破**。`apps/web/src/app/App.tsx` 的改动仅限页头响应式与 page-actions host 插槽（C5 轮次引入），顶部功能栏与左侧导航结构未改；R6 未触及后端 schema 与 Saved View 存储格式。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C5 页面 actions、视图标签、筛选网格、列表内分页 footer | 达成（C7/C8 后已修订） | E-004 + 本轮复核 |
| C7 图标强调 / 搜索配对 / 视图表单归属 | 达成 | E-005 + 真实 Chromium 几何测量 |
| C8 高度统一 / 折叠开关语义 token / 空展开抑制 | 达成 | E-006 + 真实 Chromium 多断点 + 构建产物核对 |
| 查询、重置、Saved View 存储、分页状态逻辑不变 | 达成 | 相关测试全绿；`render.tsx` handler 未改 |
| 顶部功能栏 / 左侧导航不变 | 达成 | `App.tsx` diff 范围核对 |
| 治理投影与实现层一致 | **未达成** | 见 F-001 |
| R5-I-004 未被本目标关闭 | 达成 | 未改 R5 台账 |

## Findings

### F-001 · VP-037 与 vision 工作区索引的 R6 投影落后两轮

- 严重度：med
- 建议：**required**
- 描述：`docs/vision/plans/VP-037-admin-workflow-continuity.md` 的状态行、实现边界行、R6 路线图行与 Workspaces 表（第 26、30、62、82 行）仍写 R6 `active · 5/6`，且规划短史只记录到 C5（第 105 行）；`docs/vision/workspaces.md`（第 51、80 行）同样停在 `active · 5/6`。实现层 `GOAL-007` 实为 `active · 7/8`（C7/C8 已完成）。C5 轮次曾同步这两个文件，C7/C8 两轮未同步，属本轮引入的漂移。
- 证据：`git diff 81801b30 HEAD -- docs/vision/` 无输出（C5 之后 vision 层零改动）；而 `GOAL-007/00-meta.md` 的 `progress: 7/8` 与 `goal-tree.md` 均已更新。
- 影响：决策层（VP/vision）与实现层对同一阶段状态给出不同数字，违反 P-006 对齐递归与“投影同步”要求；VP 的规划短史缺失 C7/C8 两轮事实。
- 状态：open

### F-002 · C5 曾直接改写用户已冻结的契约测试

- 严重度：med
- 建议：recommended
- 描述：C5 提交 `0cc18418` 修改了 `search-form-filters.test.tsx`，把用户于 2026-08-16 冻结的 A-003 / W13 T-03 配对断言（`-ml-px`、`rounded-l-none`、输入与按钮同格）**反转**为断言按钮不在输入同格（`expect(searchButton?.className).not.toContain("-ml-px")`、`filterActions.contains(searchButton)`）。即该回归不是“未被测试覆盖”，而是**被测试主动固化**，且当时没有对应决策条目记录这一契约变更。
- 证据：`git show 0cc18418:apps/web/src/renderer/search-form-filters.test.tsx`；C7（`e4b8c5fa`）已恢复原断言。
- 影响：冻结契约测试被静默改写会让回归从“被捕获”变为“被背书”，是比单纯漏测更高的风险模式。行为已由 C7 修正，风险在于流程。
- 状态：**fixed**（C7 恢复配对断言，E-005；本轮复核确认 `-ml-px`/`rounded-l-none`/`rounded-r-none` 断言均已回到同格配对语义）

### F-003 · 列表视觉面缺持久化浏览器级回归

- 严重度：med
- 建议：recommended
- 描述：`apps/web/e2e/` 中 `data-filter-`、`data-saved-view`、`data-list-page-actions`、`data-table-footer` 的匹配数为 **0**，即 R6 全部结论依赖 jsdom 与一次性临时预览页。jsdom 无法求值响应式 `matchMedia` 分档与真实计算高度，而这正是 C5→C7→C8 连续三轮出现布局回归、且均由用户先于测试发现的直接原因。临时预览页在核对后即删除，不构成可复跑的持久证据。
- 证据：`Select-String apps/web/e2e/*.spec.ts` 对上述选择器 0 匹配；E-004/E-005/E-006 的浏览器核对均自述“核对后已删除”。
- 影响：同类视觉/响应式回归在下次改动时仍只能靠人工发现；R6 的“视觉基线”缺少可复跑的机器守卫。
- 状态：open

### F-004 · 隐藏项提示常量重复了槽位表数值

- 严重度：low
- 建议：recommended
- 描述：C8 引入 `COLLAPSED_SLOTS` 作为折叠容量与可见性的单一真相源，但 `hiddenAtSmall`/`hiddenAtMedium`/`hiddenAtLarge`（`list-filter-panel.tsx` 第 158–160 行）仍以三元表达式手写同一组数值（1/2、2/3、3/4）。二者当前一致，但正是 C8 想消除的那类“第二份拷贝”，后续若调整断点容量极易漂移。
- 证据：`COLLAPSED_SLOTS` = base{1,1} sm{1,2} md{2,3} lg{3,4}；`hiddenAt*` 展开后数值相同。
- 影响：低——仅影响“已隐藏的生效筛选项”提示计数，不影响按键渲染与查询语义。
- 状态：**fixed**（本轮改为由 `slotCounts` 派生 `hiddenFromMobile/Small/Medium/Large`，已无手写副本；修复过程中一度漏改一处引用导致 5 项测试失败，已修正并复跑全绿）

### F-005 · 各 E 条目声称的 `tsc --noEmit` 校验实为空转（范围超出本工作区）

- 严重度：**high**
- 建议：**required**
- 描述：`apps/web/tsconfig.json` 为 `{"files": [], "references": [...]}` 的 solution-style 配置，本身不含任何源文件。因此 `npm exec -- tsc --noEmit`（不带 `-b`，也不带 `-p`）**不检查任何文件**，恒返回 exit 0。而 workspace-037 的 R6 各 E 条目（E-004/E-005/E-006）与多处审计均以「`npm exec -- tsc --noEmit --pretty false`：通过」作为类型安全证据。该证据是空转的，不构成类型校验。
- 证据（决定性）：在 `list-filter-panel.tsx` 末尾注入 `const __probe: number = "definitely a string";` 后——
  - `tsc --noEmit` → 无输出、exit 0（**未发现任何错误**）；
  - `tsc -b`（`npm run build` 实际使用）→ 报 `error TS2322: Type 'string' is not assignable to type 'number'`、`TS6133`，exit 2。
  本次修复 F-004 时确实出现过一次 `ReferenceError: hiddenAtMobile is not defined`（5 项测试失败），而同期 `tsc --noEmit` 仍报 exit 0，从实践上复现了该空转。
- 范围（超出本目标）：该 solution-style 配置**自仓库脚手架提交 `a3e1e5ad` 起即存在**且此后未被修改，故所有使用裸 `tsc --noEmit`（未带 `-p`）的历史条目都同样空转。仓库既有约定实为 `tsc -b`（`apps/web/package.json` 的 `build` 脚本与 `apps/web/README.md` 均如此），workspace-037 的 R6 轮次是偏离该约定的少数。已识别到的其他工作区文件（workspace-002/009/010/011）不在本目标 canonical 范围内，本审计只登记、不代改。
- 影响：高。类型错误可在“已通过 tsc”的表象下进入代码库，直到运行时或测试才暴露。对本 R6 而言：C5/C7/C8 的**实现正确性不受影响**（本轮已用 `tsc -b` 复核全绿），受影响的只是「类型检查通过」这一**证据有效性**。
- 状态：**in-scope 部分 fixed / cross-workspace 部分 open（待用户裁决）**
  - **fixed（本目标范围内）**：E-004/E-005/E-006 的证据表述已更正为 `tsc -b` 并注明原口径空转；本轮及后续 R6 验证一律使用 `tsc -b`。该更正属对**自身失实证据的纠错**，非取舍决策，故无需用户裁决即可闭合。
  - **open（跨工作区）**：workspace-002/009/010/011 的同类条目不在本目标 canonical 范围，本审计只登记不代改。是否追溯修正、立新目标处置或接受为残余，超出 R6 范围，**须用户按 P-004 裁决**。

## 必改项汇总（required）

| finding | 严重度 | 主张 | 状态 |
|---------|--------|------|------|
| F-001 | med | VP-037 与 `docs/vision/workspaces.md` 的 R6 投影须同步到 C7/C8 后的真实状态 | **fixed**（本轮已同步四处 VP 行、规划短史与两处 vision 索引；`git diff` 可核） |
| F-005 | high | 以裸 `tsc --noEmit` 作为类型校验证据属空转，须改用 `tsc -b`；本目标条目须更正，跨工作区影响另行处置 | **in-scope fixed / cross-workspace open**（E-004/E-005/E-006 已更正为 `tsc -b`；其他工作区条目待用户裁决） |

F-002、F-004 已 `fixed`；F-003 为 recommended，不阻断本阶段，但建议在 R6 关门后或下一次列表面改动前处理。

## 信息就绪核对

| 项 | 状态 | 备注 |
|----|------|------|
| I-007-001～004 | verified | 沿用既有基线，C5/C7/C8 未新增未知项 |
| I-007-005 | deferred non-blocking | 未实装多选，按用户指令忽略 |
| 到期 required 信息项 | 无 | 无阻断本 scope 的信息门禁 |
| 资料引用 | 无 | `shared_materials_catalog: none`；范例为仓库内路径，非共享资料 |

## 结论 + 建议下一步

C5/C7/C8 的实现与回归经独立复核**属实**，C8 的 token 边界修订合规且已进入生产构建产物。本轮发现两项 required：F-001（决策层投影落后，**已闭合**）与 F-005（`tsc --noEmit` 证据空转，high）。

F-005 的**本目标部分已闭合**（E-004/E-005/E-006 证据表述已更正为 `tsc -b`，属对自身失实证据的纠错）；**跨工作区部分仍开放**，且其范围与处置方式超出 R6，须用户裁决。

故 verdict 为 `conditional`：**R6 自身的实现与证据已达标**，但 F-005 的跨工作区部分作为已知缺陷仍需用户决定处置路径，在此之前不将 R6 标为 `done` 以保持台账诚实。

**需用户按 P-004 裁决（F-005 跨工作区部分）**，选项与建议：

| 选项 | 说明 | 影响 |
|------|------|------|
| A · 立独立目标处置（**推荐**） | 在 workspace-037 或对应工作区开设信息/整改目标，系统性核对并更正历史 `tsc --noEmit` 证据，明确后续约定 | 一次性厘清跨阶段证据有效性；成本中等，但消除同类风险 |
| B · 只改约定不追溯 | 仅声明后续一律用 `tsc -b`，历史条目保留但加注 | 成本最低；历史证据仍失实，审计追溯时需重复解释 |
| C · accepted-residual | 书面接受该历史证据缺陷，明确范围与复审触发 | 需用户书面接受；不追溯修正 |

无论选哪项，**后续 R6 及 workspace-037 的验证一律使用 `tsc -b`** 已自本轮生效。

F-003（列表视觉面缺持久化浏览器级回归）为 recommended，建议作为 R6 后续或独立小目标处理。R5-I-004 用户书面关门确认仍开放，本审计不关闭 Root/VP。
