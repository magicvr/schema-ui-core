---
id: GOAL-001-design-system-and-ui-experience
doc: decision-entry
record_id: D-003
status: accepted
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

## D-003 · 合并响应 A-001/A-002/A-003：S1 Token 映射与实施清单取舍

### 触发

- 三份独立审计已落盘：A-001（pass）、A-002（conditional / F-002）、A-003（conditional / 复审 F-002，并集 recommended）。
- 原则层三审一致：**不另建第二套 Token 真相源**合理；无 P-004 verdict/必改项冲突；开放 findings 取并集。
- 用户 `/govern` 指令：**合并响应工作区 6 的三个审计**（2026-08-09）。

### 冲突与裁决

| 项 | 结论 |
|----|------|
| 原则「不另建第二套 Token 真相源」 | **维持** D-002；三审同向，无需推翻 |
| F-002 required | **无冲突**（A-002 首报，A-003 维持）；本决策锁定无自引用映射方案；**完整 fixed 仍须 S1 实施证据** |
| recommended 并集 | 同向可叠加，全部纳入 S1 实施清单（见下） |

### 已采纳决定

#### 1. 原则维持（响应 A-001 / A-002 / A-003）

- Design Token **唯一运行时权威**仍为 `apps/web/src/index.css`（`:root` / `.dark` 原始语义 + `@theme inline` 消费别名）。
- **禁止**并行 `tokens.json` / Style Dictionary / `ds.*` 第二真相源（文档发现入口 `docs/architecture/design-tokens.md` 仍非第二真相源）。
- 「单一真相源」**不等于**禁止 alias 层；须复用现行 Color 双层纪律（A-003 重点）。

#### 2. F-002 · Shadow 无自引用映射（决策半程 · 实施后才能合法闭合）

| 层 | 命名 | 职责 |
|----|------|------|
| 原始语义（可被 `.dark` 覆盖） | `--elevation-sm`、`--elevation-md`、`--elevation-lg` | 存实际阴影值（深色下降低不透明度） |
| Tailwind theme 别名 | `@theme inline` 中 `--shadow-sm: var(--elevation-sm)` 等 | 生成 `shadow-sm\|md\|lg` utility |

- **禁止** `@theme { --shadow-sm: var(--shadow-sm) }` 同名自引用。
- **修订** D-002 §5 原文「新增 `--shadow-sm|md|lg`」：运行时**原始语义名**改为 `--elevation-*`；消费面仍为 Tailwind `shadow-*`（与 shadcn/utility 习惯一致）。
- **闭合条件（F-002 完整 fixed）**：S1 实施证据证明 (a) 上述变量存在于 `index.css`；(b) `shadow-sm|md|lg` utility 可生成；(c) 深/浅色下取值符合预期；(d) confirm/modal 等至少一处从硬编码 `shadow-xl` 迁到语义 shadow（或书面边界说明为何保留）。

#### 3. Color 双层映射纪律 = S1 实施模板（F-005）

| 规则 | 内容 |
|------|------|
| 权威值 | 仅写在 `:root` / `.dark` |
| `@theme` | 仅 alias（例：`--color-primary: var(--primary)`） |
| 禁止 | 同名自引用；文档/JSON 成为第二运行时权威 |
| 适用范围 | Color 已满足；Typography / Shadow / 未来 token **一律套用** |

#### 4. Typography 字阶路径收口（F-006）

- **选定路径**：字号阶梯 **依赖 Tailwind 默认 `text-*` scale**；S1 **不**在 `@theme` 重映射 `--text-xs`…`--text-lg`。
- **仍新增**：`--font-sans`、`--font-mono`（及必要的 `@theme` 字体族别名，遵循双层纪律）。
- **禁止**：Renderer / 业务组件硬编码 px 字号。
- 验收证据形态：消费示例用 `text-sm` / `text-base` / `text-lg` + 字体 token；非第二套字号表。

#### 5. 消费闭环清单（F-003 · recommended → S1 计划必跟）

S1 消费矩阵至少列入既有硬编码点（来自 I-001 / 审计抽样）：

| 消费点 | 现状（审计时） | S1 期望 |
|--------|----------------|---------|
| success 反馈 | `emerald-*` 等 | 明确是否进公共语义（建议 `--success` / `--success-foreground`）或书面「非 Token」边界 |
| chart | `hsl(...)` stroke | 消费 `--chart-1`…`--chart-5`（D-002 已定） |
| overlay | `bg-black/40` | 语义 overlay token 或边界理由 |
| shadow | `shadow-xl` | 语义 `shadow-*` / elevation（F-002） |
| destructive | 缺口 | D-002 已定增量 |

不进公共 Token 的类别必须写边界理由；禁止「定义有、消费无」后勾选 S1。

#### 6. 单一权威升级触发条件（F-004）

以下**任一**成为 required 时，重新评估「**生成源 → 唯一产物 `index.css`**」单向架构（仍禁止并行双真相）：

1. 多前端包 / 跨包共享 Token 消费者  
2. 非 Web 消费者  
3. 设计工具导出同步  
4. 跨仓 Token 同步  
5. 机器可验证的 Token 发布契约  

当前阶段 **不**引入 Style Dictionary / JSON 作为必须项。

#### 7. index.css 结构可读性（F-001）

S1 实施时按区块组织（Color / Typography / Elevation·Shadow / Radius / 主题引导注释等），弥补无独立 JSON 管线时的可发现性；可选配合 `docs/architecture/design-tokens.md` 发现入口。

### 未选方案

| 方案 | 未选原因 |
|------|----------|
| 推翻 D-002、引入第二套 Token 真相源 | 三审否决；与 VP-005 交付形态不符 |
| 仅 `:root` 使用 `--shadow-*` 同名且不写 `@theme` 映射 | utility/主题意图仍不清；不满足 F-002 可核对要求 |
| F-002 以 accepted-residual 放行 S1 完成 | 用户未要求 residual；映射可低成本固定，应 fixed |
| Typography 在 S1 强制重映射全部 `--text-*` | 增加与默认 scale 漂移面；默认 scale + 禁硬编码 px 已够 |

### Finding 响应状态（本决策时刻）

| ID | 来源 | 级别 | 本决策动作 | 合法闭合？ |
|----|------|------|------------|------------|
| F-001 | A-001 | recommended | 纳入 S1 结构约定 §7 | 计划采纳；非阻断 |
| **F-002** | A-002/A-003 | **required** | §2 锁定 elevation→shadow 映射 | **否**（待实施证据） |
| F-003 | A-002/A-003 | recommended | §5 消费矩阵 | 计划采纳 |
| F-004 | A-002/A-003 | recommended | §6 升级触发 | 决策已写死触发条件 |
| F-005 | A-003 | recommended | §3 双层模板 | 决策已写死 |
| F-006 | A-003 | recommended | §4 字阶路径 | 决策已收口 |

### 门禁影响

| 门禁 | 本决策后 |
|------|----------|
| S1 方案冻结 | **可继续**（I-001/I-002 已 closed；映射歧义已决策澄清） |
| S1 **完成** / 勾选检查点 / 宣称 exit 1 | **仍阻断**直至 F-002 实施证据合法闭合 |
| Root `status` / `progress` | **不变**（`active` / `0/5`） |

### 后续动作（计划 · 非事实）

1. S1 实施：按 §2–§5 改 `apps/web/src/index.css` 与消费点。  
2. 实施后用 `/govern` 或阶段自审关闭 F-002（`fixed` + 证据路径）。  
3. 可选：S1 末补 `docs/architecture/design-tokens.md` 发现入口。

### 依赖证据

- A-001 / A-002 / A-003  
- D-002（本决策修订其 Shadow 命名语义）  
- I-001 盘点；`apps/web/src/index.css` Color 双层先例  
