---
id: A-001-goal009-self-closeout
doc: audit-entry
status: recorded
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
goal_id: GOAL-009-list-visual-e2e-guard
source: self
auditor: dsh / deepseek-v4.1-flash（编排器自审）
scope: GOAL-009 C1-C4：列表视觉浏览器级守卫的范围、实现、双 profile 回归与变异验证
verdict: pass
---

# A-001 · GOAL-009 self close-out

## 头字段

- **source**：self
- **auditor**：dsh / deepseek-v4.1-flash
- **类型**：`close-out`
- **scope**：GOAL-009 的 C1（可达性与基线）、C2（守卫实现）、C3（双 profile 回归与变异验证）、C4（交付与投影）
- **verdict**：`pass`

## 范围与区间

被审对象：`apps/web/e2e/list-visual-surface.spec.ts` 及其对 GOAL-007 `A-002 F-003` 的闭合。核对方式为独立复核：重跑两 profile、对 6 类真实回归做变异测试、检查类型检查与工作树状态、核对断言是否覆盖 F-003 所指的缺口。

## 成果（有证据）

1. **守卫已实现并持久化**：`e2e/list-visual-surface.spec.ts` 覆盖 C7（搜索配对、图标、视图表单归属）、C8（高度一致、`--control` token、空展开抑制）、C5（布局顺序、列表内 footer）。它随现有 `browser-e2e` CI job 自动运行，无需改 CI 契约。
2. **双 profile 通过**：`admin` 2 passed（33.2s）、`mvp` 2 passed（35.4s）。CI 矩阵为 `profile × dialect`，两 profile 均覆盖。
3. **守卫非空转**：6 种变异全部被捕获（见 `E-003` 表）。其中变异 1 复现 C5 的搜索配对回归、变异 6 复现 C7 的视图表单位置问题——**这两个问题当年都由用户发现，测试未捕获**。
4. **断言风格避免了噪音**：关系型断言 + "等于解析出的 `--control`"，使 token 改值与间距微调不触发失败，只有破坏合同才失败。
5. **未修改被守卫的实现**：仅新增规格文件；`git status` 最终只有 `?? e2e/list-visual-surface.spec.ts`。变异测试期间的改动全部还原并校验。
6. **类型检查覆盖新规格**：`npm run typecheck`（`tsc -b && tsc -p e2e/tsconfig.json`）exit 0——这直接受益于 GOAL-008 修复的 e2e 类型检查缺口，新 e2e 代码受类型守卫保护。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 可达性与基线 | 达成 | `E-001`（挂具冒烟 1 passed；两 profile 探测；1440/700 几何基线） |
| C2 守卫实现 | 达成 | `E-002`（覆盖映射表） |
| C3 双 profile + 变异验证 | 达成 | `E-003`（6/6 捕获；2+2 passed） |
| F-003 覆盖缺口闭合 | 达成 | F-003 主张 "0 个 e2e 选择器覆盖列表视觉面"；现该面由 2 个浏览器测试覆盖，且对 6 类回归敏感 |
| 不修改实现 / 不新增 CI job | 达成 | 工作树仅新增规格；未改 `playwright.config.ts` |
| 不改变 Root 六阶段分母 | 达成 | 整改子目标，不计入分母 |

## Findings

### F-001 · 守卫未覆盖 `.dark` 下的开关实际底色（只断言 token 有覆盖值）

- 严重度：low
- 建议：recommended
- 描述：C8 item 2 的暗色一侧目前断言"`.dark` 下 `--control` 存在且不同于浅色值"，但没有断言暗色下开关**计算背景**确实等于该值。若将来开关的类名被改成不消费 `--control`（例如硬编码浅色背景），浅色断言会捕获，但暗色路径无独立断言。
- 证据：规格中暗色检查读取 CSS 变量后即移除 `dark` 类，未在暗色下重读开关的 `background-color`。
- 影响：低——浅色路径的"等于 `--control`"断言已钉住 token 驱动关系；暗色是同一 CSS 变量作用域内的覆盖。
- 状态：**fixed**（2026-09-18，经 [workspace-010 GOAL-043](../../../workspace-010-design-implementation-conformance/GOAL-043-w31-cross-workspace-residual-closeout/00-meta.md) 修复：暗色下新增四条不变量——根 token = 暗色覆盖值、开关继承同一值、`--color-control` 别名解析到该值、开关**计算背景**等于该值；变异验证＝把类名改成 `bg-control dark:bg-[oklch(0.955_0_0)]`（暗色硬编码浅色）时浅色断言全过、新增暗色断言失败并指明原因。证据：`GOAL-043 E-002` §1）

> 修复过程中的一个坑（记录备查）：开关带 `transition-colors`，加 `.dark` 后立即 `getComputedStyle` 读到的是过渡插值帧（浅色），会被误判为缺陷；正确读法是读取前内联 `transition: none`。这不是产品缺陷。

### F-002 · 目标页仅 `/roles` 一个，覆盖面受页面差异限制

- 严重度：low
- 建议：recommended
- 描述：守卫只跑 roles 页。其他列表页（users、account sessions 等）若出现同类回归不会被捕获。选择单页是 `D-001` 的显式取舍（roles 是唯一同时具备 search form 与 toolbar 的页面），但覆盖面确实窄于"通用列表页"这一合同表述。
- 证据：`D-001` §2；`E-001` schema 检索结果（table `filters` 仅 account 一处；search+toolbar 交集仅 roles）。
- 影响：低——被守卫的是**共享实现**（`SchemaTable` / `ListFilterPanel` / `render.tsx` 的 search slot），单页即可回归共享层；页面差异主要影响数据而非这些合同面。
- 状态：**fixed**（2026-09-18，经 workspace-010 `GOAL-043` 修复：分页契约用例改为 **roles + users 双页面参数化**，断言默认显示 20、首个列表请求不带 `pageSize`、选 10 后确实发出 `pageSize=10`、跳转按钮文案；两 profile 各 4 passed。证据：`GOAL-043 E-002` §2）

## 必改项汇总（required）

无。本目标无 required finding。

## 信息就绪核对

| 项 | 状态 | 备注 |
|----|------|------|
| I-009-001（挂具可跑性） | verified | `E-001` 冒烟 1 passed |
| I-009-002（目标页与 profile 可达性） | verified | `E-001` 两 profile 探测 |
| I-009-003（几何基线） | verified | `E-001` 1440/700 实测表 |
| I-009-004（是否引入像素快照） | deferred non-blocking | 用户未要求；`D-001` 未选方案 |
| 到期 required 信息项 | 无 | 无阻断关门的信息门禁 |
| 资料引用 | 无 | 本区 `shared_materials_catalog: none` |

## 结论 + 建议下一步

C1～C4 全部达成，无 required finding，verdict 为 **`pass`**。GOAL-007 `A-002 F-003` 所指的覆盖缺口已闭合：列表视觉合同现有持久化浏览器级守卫，且在 mvp/admin 两 profile 下通过，对 6 类真实回归（含当年由用户发现的两类）敏感。

**可关门**：GOAL-009 投影为 `done · 4/4`。它为非纲领整改子目标，Root 六阶段分母与 `progress: 5/6` 不变；R5-I-004 用户书面关门确认仍开放，本目标不关闭 Root 或 VP-037。

**建议**：F-001/F-002 为 recommended，可择机加固（暗色下重读开关背景；如需更宽覆盖再增一个页面）。R5 `GOAL-006` 现为 Root 六阶段中唯一未完成项，其关门取决于用户书面确认。
