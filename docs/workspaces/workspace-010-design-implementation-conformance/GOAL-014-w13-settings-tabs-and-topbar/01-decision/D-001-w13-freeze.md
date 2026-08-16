---
id: D-001
doc: decision
status: accepted
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-001 · 开波：四项交互整改纳入本波；设计冻结

## 背景

用户 2026-08-16 点名在 workspace-010 新增子目标，处理四项产品面交互整改（对照个人中心页既有 Tabs 形态与 W12 交付的搜索/顶栏 as-built）。

## 决策

### T-01 · 设置页改为按功能单元 Tabs 切换（对齐个人中心页）

- `apps/api/internal/modules/settings/schema/settings.json` 的 body 由「单个 section 内纵向堆叠全部表单」改为：`section[描述文本 → tabs(五个功能单元) → actionButton 恢复默认]`。
- 五个功能单元沿用现有表单：常规（站点标题）/ 品牌（Logo/Favicon 上传）/ 本地化（默认语种+时区）/ 外观（默认主题）/ 安全（登录验证码）。
- tab 标签新增 `schema.settings.tab.{general,branding,localization,appearance,security}` 文案键（双语文案表同步；en-US 为基线）。
- 各 tab 内表单保留其 `titleKey`（`toolbar.*`）标题——与个人中心页「tab + 面板内标题」形态一致。
- **恢复默认（reset）按钮与页面描述文本保留在 Tabs 之外**：破坏性动作任何 tab 下都可达，且不把「动作」伪装成「功能单元」。
- 权限门禁保持 body section 的 `permissionCascade.edit`（settings.write）不变；tabs 子 section 不重复声明（避免级联双倍）。
- 未选方案：① 不把 reset 做成第六个 tab（破坏性动作与功能单元语义不同）；② 不改协议/renderer 以支持「跨 tab 常驻动作条」（成本大于收益）。

### T-02 · 移动端品牌条（logo + 网站标题单独一条）

- 断点沿用壳层既有约定：**< lg（1024px）为移动端**（与汉堡/抽屉导航、sidenav 隐藏同断点）。
- <lg：header 顶部新增独立品牌行（`data-shell-region="mobile-brandbar"`，`lg:hidden`），仅含 logo + 站点标题 + 副标题；下方主行只保留功能区控件（汉堡、亮暗、语种、通知铃、用户菜单），`ml-auto` 全宽可用。
- ≥lg：保持现有单行布局（品牌链接 `hidden lg:flex` 回到主行）。
- 品牌链接 JSX 提取为局部组件避免双份复制；两条品牌链均保留原跳转/aria 语义。

### T-03 · 搜索框组【文本框+搜索键】贴合恒同行

- search 模式下关键词 input 与搜索按钮在**同一网格单元**内：`flex flex-nowrap` 零间隙（去掉 `gap-2`），按钮 `-ml-px` 叠边、input 右侧与按钮左侧内角改直角（`rounded-r-none` / `rounded-l-none`）——视觉上贴合为一个组件。
- `flex-nowrap` + 按钮 `shrink-0` + 输入侧 `min-w-0 flex-1` 保证任何页面宽度下两者**恒在同一行**，不可能被折行拆开（语义单组件）。
- 未选方案：把输入放进带边框的外壳容器（改动面大、需重做 focus 环与 label 对齐，收益有限）。

### T-04 · 顶栏功能区亮暗/语种按键对调

- 功能区顺序由 `[汉堡, 语种, 亮暗, 铃, 用户]` 改为 `[汉堡, 亮暗, 语种, 铃, 用户]`：**亮/暗切换在左、语种切换在右**。
- 仅调换 JSX 顺序；组件自身 aria-label/title/行为不变。

## 影响

- **go 判定**：四项均为呈现/交互层整改；设置页 schema 仅 body 结构重组（字段/动作/URL 不变）；不改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义 → **VP-008 go 无影响、不暂挂**。
- **测试影响**：`startup-config.test.tsx`（四类表单同屏断言 → 改为 Tabs 切换断言）、`e2e/localization.spec.ts`（品牌/本地化/外观标题 → 点击 tab 后断言）、`shell.test.ts`（静态结构可加 T-02/T-04 守卫）。
