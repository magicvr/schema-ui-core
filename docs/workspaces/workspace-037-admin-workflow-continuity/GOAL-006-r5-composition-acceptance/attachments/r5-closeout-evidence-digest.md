---
id: r5-closeout-evidence-digest
doc: audit-evidence-digest
status: provided
goal_id: GOAL-006-r5-composition-acceptance
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-006-r5-composition-acceptance
version: 1.0.0
---

# R5 关门证据一页（供用户 P-004 裁决）

- 对象：`GOAL-006-r5-composition-acceptance`（R5）· 信息项 **R5-I-004**（required）
- 生成：2026-09-18 · 应你要求「先给一页证据再定」而产出
- 本页**不**改变任何 `status`/`progress`，**不**构成关门；裁决仍归你
- 依据文件：`00-meta.md`（信息表）、`03-audit.md`（三份意见台账）、`02-execution/E-001～E-006`

## 0 · 你要裁决的唯一一条

| 项 | 内容 |
|----|------|
| ID | `R5-I-004`（required；最晚需要阶段 C4） |
| 问题 | **用户是否书面确认 Root/VP 关门？** |
| 同一条门禁 | `A-001 R5-GATE-001`（self）、`A-002 F-001`（grok independent，high · required） |
| 当前状态 | `collecting` — 无用户书面确认 |
| 合法闭合路径 | 仅你的书面确认（P-003：禁 `fixed` 冒充、禁审计员 residual/overrule、禁编排器静默推断） |
| 确认前约束 | `GOAL-006`、Root `GOAL-001`、VP-037 与 workspace 必须保持 `active` |

## 1 · 台账快照（生成时事实）

| 层 | 对象 | status · progress | 阶段审计结论 | checkpoint |
|----|------|-------------------|--------------|-----------|
| Root | `GOAL-001-admin-workflow-continuity` | `active · 5/6` | 开放 required = 0（A-005 为历史 R6 版本投影） | — |
| R1 | `GOAL-002-r1-scope-semantics-freeze` | `done · 3/3` | A-001 self `pass`、A-002 grok independent `pass`、A-003/A-004 响应 `pass` | — |
| R2 | `GOAL-003-r2-saved-views` | `done · 4/4` | A-001 self `pass`、A-002 grok independent `pass`、A-003 close-out `pass` | `39c744ef` |
| R3 | `GOAL-004-r3-unsaved-change-protection` | `done · 4/4` | A-002 independent `conditional`（F-001）→ A-003 recheck `pass` → A-004 `pass` | `d2b39189` |
| R4 | `GOAL-005-r4-unified-feedback-recovery` | `done · 4/4` | A-002 independent `conditional`（F-001 transport 分类）→ A-003 recheck `pass` → A-004 `pass` | `89666e5c` |
| R5 | `GOAL-006-r5-composition-acceptance` | `active · 3/4` | A-001 self `conditional`、A-002 grok independent `conditional`、A-003 响应 `conditional`；唯一开放 required = **R5-I-004** | — |
| R6 | `GOAL-007-list-page-visual-alignment` | `done · 8/8` | A-001 `pass`；A-002 `conditional` 的 F-001～F-005 全部合法闭合；A-003 `pass` | `81801b30`/`0cc18418`/`e4b8c5fa`/`dc4b5c5b` |
| 整改 | `GOAL-008-typecheck-evidence-convention` | `done · 4/4` | A-001 `pass`（开放 required = 0）；2026-09-18 追加 A-002 `conditional`（跨区勘误可核对；F-001 守卫缺口、F-002 简写未裁定，均 recommended open） | `09a36ebf` |
| 整改 | `GOAL-009-list-visual-e2e-guard` | `done · 4/4` | A-001 `pass`，开放 required = 0（非纲领，不计入 Root 分母） | `2a9ec682` |
| VP | `VP-037-admin-workflow-continuity` | `active` v1.5.0 | VRev-094（计划）self `pass`、VRev-095（激活）self `pass`；open required = 0；V-F124 recommended | — |

## 2 · VP-037 方向级退出判据 1～7 逐条

| # | 判据（VP-037 §方向级退出判据） | 状态 | 可回指证据 |
|---|--------------------------------|------|-----------|
| 1 | Saved Views 分母、用户隔离、字段序列化、恢复/失效与权限边界形成可核验矩阵，且不把共享/协作偷偷纳入 | 达成 | `GOAL-002` D-003（localStorage 方案 A，用户确认）/D-004/D-005；`attachments/r1-denominator-matrix.json`（24/58 分母）、`r1-form-matrix.json`、`r1-state-feedback-matrix.md`；A-002 independent 独立复算 |
| 2 | Saved Views 保存/选择/恢复/更新/删除及空态/错误态闭环；无效或越权视图 fail closed 且不污染页面 | 达成 | `apps/web/src/renderer/saved-views.ts`、`schema-table.tsx`、`render.tsx`；`saved-views.test.ts`、`saved-views.ui.test.tsx`；`attachments/r2-saved-view-acceptance-matrix.md`；`GOAL-003` A-002 independent `pass` |
| 3 | dirty-state 状态机覆盖内部导航、刷新/关闭、提交、重置、取消；误离开保护可由回归验证 | 达成 | `apps/web/src/renderer/dirty-state.ts`、`app/App.tsx`；`dirty-state.test.ts`、`r3-dirty-state.ui.test.tsx`、`App.integration.test.tsx`；`attachments/r3-dirty-state-acceptance-matrix.md`；`GOAL-004` A-003 independent recheck `pass` |
| 4 | 统一反馈语义与可访问呈现；不重复提交、不吞服务端错误、不泄露敏感信息 | 达成 | `apps/web/src/components/ui/feedback.tsx`、`renderer/feedback-policy.ts`、`components/data-table.tsx`；`feedback.test.tsx`、`feedback-policy.test.ts`、`representative-pages.integration.test.tsx`；`GOAL-005` A-002 F-001（transport 分类）`fixed` + A-003 independent recheck `pass` |
| 5 | 首波不解除实体全文搜索 / `RT-X01`/`RT-X02` / 批量结果中心 / 组织与数据权限 / 架构分支 gated 项，边界证据可追溯 | 达成 | `GOAL-006 02-execution/E-003`（边界与对齐核对）；Root `00-meta.md`「明确非目标」；`apps/api` 自 R5 基线 `89666e5c` 起**零文件改动**（`git diff --name-only 89666e5c..HEAD -- apps/api` 为空） |
| 6 | VP 工作区阶段链、Goal 审计、required 信息与 Vision Review 全部闭合，关门前组合投影同步，开放 required = 0 | **待你的确认**（判据的最后一环） | 阶段链/审计/Vision Review 已闭合（见 §1）；R1～R4、R6、GOAL-008/009 台账开放 required = 0；C4 投影清单见 §5 |
| 7 | 通用列表页视觉与交互基线由共享实现覆盖；顶部功能栏/左侧导航不变、查询/重置合同不变、单页有效列表仍显示分页 | 达成 | `GOAL-007` `done · 8/8`（C5/C7/C8 三轮纠偏）；`apps/web/src/components/data-table.tsx`、`renderer/list-surface.tsx`、`renderer/schema-table.tsx`、`index.css`（`--control` token）；`GOAL-009` 浏览器级守卫 `apps/web/e2e/list-visual-surface.spec.ts` |

## 3 · 本次复跑（2026-09-18，裁决前证据刷新）

| 命令 / 检查 | 结果 |
|-------------|------|
| `apps/web`: `npm run typecheck`（`tsc -b` + `tsc -p e2e/tsconfig.json`） | exit 0（非空转口径，见 GOAL-008） |
| `apps/web`: `npm test`（Vitest） | **112 文件 / 1426 测试全部通过** |
| `apps/web`: `npm run test:e2e -- list-visual-surface`（`APP_PROFILE=admin`，真实 Go API + Chromium） | **2 passed**（32.9s） |
| `apps/web`: 同上（`APP_PROFILE=mvp`） | **2 passed**（33.3s） |
| `git diff --check` | 通过（仅 LF/CRLF 转换提示，无 whitespace error） |
| `apps/api` 漂移 | 自 `89666e5c` 起 0 文件改动 |
| 结构扫描 | Root 与 R1～R6 七目标五件套 + 三个 ledger 目录 + `attachments/` 齐全；`GOAL-008`/`GOAL-009` 同样齐全 |

对照 R5 C3（2026-09-17）记录的 `110/1408`：数量增长来自 GOAL-008（类型检查约定守卫）与 GOAL-009（列表视觉 e2e 用例）新增的测试文件，无既有用例被删除或跳过。

## 4 · 关门后**仍开放 / 仍 gated**（确认 ≠ 这些已验证）

| 项 | 性质 | 触发 / 归属 |
|----|------|-------------|
| `V-F124` | recommended（愿景层） | 保持 open |
| `R5-I-005` / A-002 `F-003` Host/resource 直接对照 | non-blocking deferred | 真实支持需求出现 → `/vision` |
| `I-037-005` 跨用户共享视图 / 最近使用 / 收藏 / 协作权限 | deferred | owner `/vision`；真实协作需求 |
| `I-007-005` 未实装的多选 | deferred | 按你此前指令忽略 |
| `GOAL-008 A-001 F-002` 守卫未正向断言 CI 步骤存在 | recommended | 后续加固 |
| `GOAL-008 A-002 F-001` 守卫按 `-p` 令牌判定检查型调用，`tsc --noEmit -p tsconfig.json`（实测空转）仍会被判合规 | recommended open（2026-09-18 新发现） | 当前无可执行面使用该形态；是否修复待你裁决（新目标或并入后续轮次） |
| `GOAL-008 A-002 F-002` 全仓 `tsc` 简写（300+ 行）未逐条裁定 | recommended open | 由 README 约定 + 守卫 + CI 门禁保证新增记录 |
| `GOAL-008 I-008-004` 其他工作区历史条目追溯更正 | 2026-09-18 你已授权并**已执行**：11 处勘误注记、2 处复核确认有效（`E-006`/`E-023`）；`I-008-004` → `verified` | 完成 |
| `GOAL-009 A-001 F-001` 暗色开关计算背景未断言 · `F-002` 仅 roles 页覆盖 | recommended | 后续加固 |
| gated 非目标：实体全文检索 / `RT-X01`/`RT-X02`、批量结果中心、组织·部门·岗位与 `org` 数据权限、新业务域、Redis/MQ/多实例/第二持久化栈 | 未解锁 | 不因关门解禁 |

## 5 · 确认后我会执行的关门动作与 C4 投影清单

1. 你的原话落盘为 P-004 用户书面裁决（`GOAL-006` 决策/执行台账），`R5-I-004` → `verified`；`A-001 R5-GATE-001` / `A-002 F-001` → `fixed`（附确认证据指针）。
2. `GOAL-006` C4 勾选 → `done · 4/4`；Root `GOAL-001` → `done · 6/6`；`workspace.md` 与 `goal-tree.md` 同步。
3. VP-037 `Closeout placeholder` 填实（判据 1～7 逐条证据 + 用户确认 + 投影同步），VP-037 → `closed`，版本推进。
4. Vision 层投影同步（A-002 F-002 清单）：`docs/vision/roadmap.md`（VP-037 行与三处叙述仍写 plan `v1.1.0`、R6 `done · 4/4`）、`docs/vision/workspaces.md`、Charter「现行组合投影」快照 —— Charter 属愿景层 editorial，按 A-002 F-002 由 `/vision` 同步。
5. 执行 VP 关门 Vision Review（对照 VP-032/033/036 惯例为 self `pass`）。
6. 修掉本轮核对中发现的两处**目标内当前态**滞后（不涉及关门）：`GOAL-006 00-meta.md` 父目标行的 Root `active · 4/5`、`GOAL-008 00-meta.md` 备注的 `progress: 4/6`。其余 `4/5`/`4/6` 出现处均为时间线历史条目，按事实保留不改。**（已于 2026-09-18 本轮完成）**
7. 说明：§4 中的跨区勘误（`I-008-004`）与新发现的守卫缺口（`A-002 F-001`）已在本轮先行执行/登记，不属于确认后的待办；确认后只需在 C4 投影中一并复核。

## 6 · 确认模板（可复制）

> 确认 R5-I-004：接受 VP-037 首波交付与上述 residual 边界，同意 Root `GOAL-001` 与 VP-037 关门。

若你想缩小范围，也可只确认部分，例如「接受交付但不关门，继续保留 active」，我会按你的措辞落盘而不执行关门投影。
