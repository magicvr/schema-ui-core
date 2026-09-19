---
id: A-001-w33-self
doc: audit-entry
parent: GOAL-045-w33-list-actions-slot-and-roles-trigger
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-001 · W33 自审（GOAL-045 C1～C3）

## A-001 · W33 C1～C3 自审（2026-09-19）

- **source**：self
- **auditor**：`/govern` 编排器（DeepSeek Harness 会话）
- **类型** / **scope**：`stage` · GOAL-045 C1～C3（方案冻结 / 插槽与 roles 触发面实施 / 回归）
- **verdict**：**pass**（0 required；4 recommended，其中 F-004 已在本区间 fixed）
- **审计模式判定**：**`self`**。理由：改动为渲染器**本地扩展 + 页面 schema 补齐**，无后端/权限/迁移/协议面变更；跨页回归风险在实施中**已被浏览器 e2e 实际捕获并修复**（§成果 5），并留下回归锁。跨工作区可见效果（VP-038 已交付页）由全量 vitest + 双 profile e2e 复核。

### 成果（有证据）

| # | 主张 | 复核 |
|---|------|------|
| 1 | 插槽按用户选型 A 落地：`data-list-page-actions` 行拆左右两段，custom 节点以 `props.slot` + `props.targetTable` 投放左段 | `render.tsx`（`declaresListPageActionsSlot` + `ListActionsSlotNode` + 4 个 seam）、`schema-table.tsx`（左段宿主 + `justify-between`）；`list-actions-slot.test.tsx` 断言节点确实落在左段 |
| 2 | 兼容边界成立：未声明 slot 的 custom 节点**逐字原路径**（纯谓词分流，不新增 context 订阅） | `render.tsx` 的 custom 分支；测试「leaves a node without the declaration in document flow」 |
| 3 | 回落为 fail-open（目标表缺失/拼写错 → 原地渲染，控件不消失） | `D-001` §3；测试「falls back …」与「unknown slot value」两例 |
| 4 | roles 页补齐：`selection.mode=multiple` + `roles-batch-export`（resource=roles，同插槽声明）；**后端零改动** | `roles.json`；`jobs-batch-export-roles.test.tsx`（左段落位、空选择禁用、提交体 `{resource:"roles",ids:[…]}`、不 reloadList） |
| 5 | **实施中被浏览器 e2e 捕获并修复一处跨页状态串扰**：roles 页新增 custom 节点后，两页 table 落在同一子节点索引 → React 复用同一 `SchemaTable` 实例 → 上一张表的 `visibleColumns` 渗入另一张表（用户名列消失） | 归因：干净树该 spec **通过**、带 W33 **失败**；instrumented 打印 8 列 → 8 列 → 3 列。修复：渲染器按 table 节点 id 作 key（不同表即不同实例）。新增回归锁「keeps two different tables' column state apart across a page switch」 |
| 6 | 变异验证 | `roles.json` 的 `slot` 拼写改成 `list-page-action` → shipped 声明例与 roles 落位例**双双变红**；还原后 7/7 绿 |
| 7 | 既有 flake 的诚实归因与处置 | `s5-denominator-render` 2 例在**干净树**上同样失败（基线实验），`--testTimeout=20000` 全套件转绿 → 判定既有负载敏感；为 3 个真实 App 渲染用例补显式超时（仅测试可靠性，无产品变更） |
| 8 | 判据/回归证据 | `npm test` 121 files / **1476 tests 全绿**（默认超时）；`typecheck` / `build` exit 0；e2e **mvp 16 passed / 5 skipped / 0 failed**、**admin 17 passed / 4 skipped / 0 failed**；Go 相关包与 docscheck 全绿 |
| 9 | **用户报告的运行期控制台循环已定位并修复**（`Maximum update depth exceeded`）：经浏览器探针 + 插桩计数 + **W33 前基线对照**（`c6346f13` 同样复现）确认**先于本波次存在**；根因 = 保存视图 page-actions 宿主的 claim 效应把 `availabilityVersion` 列为依赖、其 cleanup 又 `release()` 递增该值 → 自持循环；修复 = 拆分为「claim 效应（只依赖宿主元素 + tableId + 重试计数）」+「重试效应（宿主空闲时 +1）」 | `E-001` §5 全过程；`schema-table.tsx` 注记 |
| 10 | **新增永久守卫**：`e2e/console-health.spec.ts` 走真实外壳（users → roles → users，admin 另加 jobs），任何 `console.error` / `pageerror` 即失败（仅静态资源白名单；网络错误带 URL 便于定位） | 变异验证：把 `availabilityVersion` 加回 claim 依赖 → **守卫立即报出成片 max-update-depth**；还原后两 profile 全绿（mvp 17 passed / admin 18 passed） |
| 11 | 判据/回归证据 | `npm test` 121 files / **1476 tests 全绿**（默认超时）；`typecheck` / `build` exit 0；e2e **mvp 17 passed / 5 skipped / 0 failed**、**admin 18 passed / 4 skipped / 0 failed**（均含 console-health）；Go 相关包与 docscheck 全绿 |
| 12 | 边界 | 未改后端（导出分母/job 运行时/权限键）；未改 `props.toolbar` 语义、权限与禁用规则（右段样式不变）；未触 pinned 工件；**未改 VP-038 `status`**（保持 `closed`）与 workspace-038 台账正文 |

### Findings

#### F-001 · 插槽机制只定义了 `list-page-actions` 一个值，其它 slot 值静默等同未声明

- 严重度：low
- 建议：**recommended**
- 描述：未知 `slot` 值按「未声明」处理（原地渲染）是有意的 fail-open，但**没有任何提示**：若后续波次新增 slot 名称却忘了同步渲染器，表现是「控件位置没变」而非报错，排查成本高。
- 证据：`D-001` §1；`declaresListPageActionsSlot` 只比较一个常量。
- 状态：**open（recommended）** —— 处置建议：由「本地扩展登记」（roadmap 未决项）一并记录支持的 slot 值清单。

#### F-002 · 左段目前靠 flex 顺序天然靠左，未固定「左段始终在右段之前」的结构约束

- 严重度：low
- 建议：**recommended**
- 描述：左段渲染在右段之前（DOM 顺序）由当前 JSX 顺序保证，测试断言了 `row.firstElementChild === left`；若将来有人把右段提前（例如为了移动端换行顺序），左段会跑到右边，而 `justify-between` 会让它看起来像右对齐组的一部分。
- 证据：`schema-table.tsx` 行结构；`list-actions-slot.test.tsx` 的 `firstElementChild` 断言（已提供一定保护）。
- 状态：**open（recommended）** —— 已有断言覆盖主要退化形态，可接受为低危。

#### F-003 · 移动端/窄屏未做视觉验证

- 严重度：low
- 建议：**recommended**
- 描述：左段 + 右段在窄屏会换行；`flex-wrap` 下左段先换行、右段可能单独成行——视觉上是否可接受未验证（本次只在桌面视口跑了 e2e）。
- 证据：e2e 均为桌面视口（`playwright.config.ts` 默认 1280×720 类视口）；`list-visual-surface.spec.ts` 的移动端用例只覆盖既有控件。
- 状态：**open（recommended）** —— 处置建议：由后续符合性波次或用户实机确认。

#### F-004 · 测试矩阵对「控制台错误」这一类回归完全盲（已用守卫闭合）

- 严重度：**med**
- 建议：**recommended**
- 描述：`Maximum update depth exceeded` 在列表页持续刷屏，而 vitest（渲染片段、无 App 外壳）与全部 e2e（断言 UI 行为）**长期全绿**——"页面能用" 掩盖了 "状态自激循环"。这不是 W33 引入的缺陷（基线可复现），但它暴露的是**验证面的结构性缺口**。
- 处置：**已 fixed** —— 新增 `e2e/console-health.spec.ts` 永久守卫（见成果 10），并以变异验证证明其有效性（重新引入循环 → 守卫失败）。
- 证据：`E-001` §5；`apps/web/e2e/console-health.spec.ts`。
- 状态：**fixed**

### 必改项汇总（required）

**无。**

### 结论 + 建议下一步

两项用户诉求均按选型 A 落地：**roles 页具备「导出所选」**（后端零改动），**按钮进入列表控件行左侧**（与右对齐的既有工具栏按钮同一行）。实施过程中发现并修复了一处**跨页表格状态串扰**的既有隐患，并留下回归锁；同时把一处既有 flake 归因清楚并加固。

**verdict = pass**。**建议下一步**：C4 响应三条 recommended（登记 / 记录 / 交后续波次），同步 `goal-tree.md`/`workspace.md`，在 `roadmap.md`「未决项统一登记」把「roles 触发面补齐」标为 `fixed`、并把 `list-page-actions` 插槽登记进「本地扩展」。

**本条不修改** `status`、检查点或派生 `progress`。
