---
id: D-004-r6-c8-control-height-and-toggle-token
doc: decision-entry
status: accepted
goal_id: GOAL-007-list-page-visual-alignment
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-004 · R6 追加 C8：控制高度、折叠开关语义 token 与空展开抑制

## 决定

用户在 C7 交付后追加三项 UI 反馈，仍属 R6 同一列表视觉分母。继续在同一 `GOAL-007-list-page-visual-alignment` 内追加 C8（不新建 GOAL-008），R6 检查点由 7 项扩为 8 项。

C8 合同：

1. **统一页面 actions 高度**：页面级 actions 行内的 schema toolbar 触发器由 `h-9`/`text-sm` 降到 `h-8`/`text-xs`，与“列配置”触发器同高（实测均为 32px）；“列配置”触发器同步由 `px-2.5 py-2` 改为 `h-8` 显式高度，使两者由“接近”变为“精确相等”。
2. **折叠开关语义 token**：展开/收起筛选按键改用新增语义 token `--control` / `--control-foreground`（经 `@theme inline` 映射为 `bg-control` / `text-control-foreground`），与“重置”这类普通描边操作在视觉上明确区分。该 token 走既有“`:root`/`.dark` 原始语义值 + `@theme inline` 别名”双层纪律，深浅色各有取值。
3. **空展开抑制**：当折叠后首行已能容纳全部筛选项（即没有任何控件被折叠隐藏）时，不渲染展开/收起按键，也不渲染它会占用的操作单元；此时若宿主未提供 actionSlot，整个操作单元都不存在。

## 理由

第 1 项是同一行内的视觉一致性问题：`h-9` 的 primary 触发器与 `h-8` 的列配置触发器并排会形成高度阶梯。统一到 `h-8` 后整行基线一致，且不改变任何按钮的语义层级（primary 仍是 primary）。

第 2 项用户明确指出“这应该是一个新的语义 token”。范例页确实把两者做成不同语义：重置是 `text-zinc-400` + 透明底 + hover `bg-zinc-800/60` 的普通描边操作；展开筛选是 `text-zinc-300` + **填充底 `bg-[#121215]`** + hover `bg-white/[0.06]` 的控件 chip。填充底是“可交互控件”这一语义，不能复用 `background`（会被读成普通操作）、也不宜复用 `muted`/`accent`（它们是状态/悬停语义，且 fork 常整体调色）。因此新增 `--control` 独立表达该语义，并纳入 fork 可覆盖面。

第 3 项是纯逻辑缺陷：只有一行筛选项时“展开”无内容可展开，按钮存在即误导。判定必须与实际折叠可见性同源，否则会在窄屏（确实有隐藏项）与宽屏（无隐藏项）之间出现按钮与实际状态不一致。

## 未选方案

- 不把 toolbar 触发器降到 `h-7`：与同行“列配置”及 Saved View 管理按键的 `h-8`/`h-9` 混合尺度相比，`h-8` 已是该行最小一致尺度，继续降低会牺牲可点击面积。
- 不为折叠开关复用 `bg-muted`/`bg-accent`：会把“控件 chip”语义混入状态语义，且 `accent` 已是 hover 底色，同元素 hover 态将无法区分。
- 不用纯 CSS（`:has()`/媒体查询计数）实现第 3 项：折叠可见性是响应式槽位计算的结果，用 JS 与可见性类共用同一张槽位表才能保证按钮与实际隐藏项永不漂移；纯 CSS 无法知道“当前断点下是否真有隐藏项”。
- 不新建 GOAL-008：同一视觉合同的第三次局部纠偏，回开同一目标保留历史证据边界。

## 门禁

C8 需实现、定向回归与自审证据完成后，才可进入 C6 修订审计；C5/C7 的 E-004/E-005 不再单独代表当前交付状态。R5-I-004 仍开放，R6/C8 不得关闭 R5、Root 或 VP-037。审计模式按常规、可逆的 UI 纠偏继续记为 `self`。

本决策局部修订 D-002 §5 与 D-001 的“不新增全局语义 token”边界：允许**新增**一个语义 token（`--control`），仍禁止重命名或重定义既有 token 值。新增 token 必须同时进入 `index.css` 双层声明、`@theme inline` 映射与 `theme.test.ts` 结构守卫。
