---
id: GOAL-001-design-system-and-ui-experience
doc: decision-entry
record_id: D-004
status: accepted
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

## D-004 · 视觉方向冻结（Stitch 定稿）

### 触发

- VP-005 / Root 已 active；S1 Token 约定（D-002/D-003）已 accepted。
- 用户经 [stitch.withgoogle.com](https://stitch.withgoogle.com) 多轮生成/纠偏（含响应式、移动卡片列表、暗色对比度），导出至本地  
  `raw/stitch-vp005-visual-refs/exports/stitch_schema_ui_core_admin_console/`。
- 第二轮评审：13 屏 `screen.png` 均健康；结论 **可用作 VP-005 视觉方向输入**（非实现完成证据）。
- 用户确认：将视觉方向 **落盘冻结** 到本 Root，再推进 S1 实施。

### 已采纳决定

1. **视觉方向参考权威（过程）**  
   - 定稿目录（本地，**gitignored**）：`raw/stitch-vp005-visual-refs/exports/stitch_schema_ui_core_admin_console/`  
   - 评审摘要：`raw/stitch-vp005-visual-refs/exports/notes.md`  
   - 提示词/约束材料：`raw/stitch-vp005-visual-refs/`（00–07）  
   - 仓库内摘要指针：`attachments/visual-direction-stitch-summary.md`  
   - **不**把 Stitch `code.html` 当作生产源；生产仍为 React + Schema + `apps/web/src/index.css` Token。

2. **气质**  
   - Linear + Vercel Dashboard 式：克制、高密度、工作导向。  
   - shadcn/ui new-york 语汇；中性灰；近黑主按钮；Geist/Inter 系无衬线。  
   - 禁止营销 hero、底部 Tab 消费 App 壳、Material 默认紫/亮蓝主色。

3. **壳层（S3 主约束；S1 主题须兼容）**  
   - 桌面（≥1024）：sticky 顶栏 + 常驻左栏（~256px）+ 主内容。  
   - 移动（&lt;768）：无常驻侧栏；**汉堡 + 全高导航抽屉**；单列。  
   - 深/浅色：class 策略（与 D-002 一致）；暗色主文案须高对比可读（定稿 dark Overview 为下限参考）。

4. **列表呈现（S2 主约束）**  
   - **桌面**：高密度多列 **数据表**（`table` 语义）。  
   - **移动**：同一数据语义的 **卡片列表**（主标题 + 1～2 行次要 + `⋯` 操作）；**禁止**把多列表格拉成窄屏主路径。  
   - 横滑表仅次要备选，不作为 Stitch/产品默认英雄布局。

5. **详情与编辑（S2/S3）**  
   - 详情（`recordView`）：桌面 **右侧详情栏或 Drawer**；移动 **全高 Sheet**。  
   - **居中 Modal** 仅用于新建、短字段编辑、Confirm——**不作**完整详情主表面。  
   - 实现字段以协议/范例 schema 为准；Stitch 字段仅为布局示意。

6. **实施参考优先级（非阶段强制序）**  
   1. Overview 浅色 / 暗色 / 移动 → Shell + Token 气质  
   2. Data table 桌面 + 移动卡片 → 列表双端  
   3. Users 桌面 + 移动 → CRUD + 详情栏  
   4. Sign in → 登录  
   5. Search + table / Data display / Form controls / Form with reactions → 其余 S2 表面  

7. **明确不因本决策而改变**  
   - **不**勾选 S1–S5；`progress` 仍 `0/5` 直至实施与门禁证据。  
   - **不**将本决策读成 Charter #3 或 VP-005 exit 已满足。  
   - **不**扩张 `I-PROTO-FULL-001`；**不**引入 `Detail`/`Filter` 杜撰 Node。  
   - Token 命名与权威落点仍以 **D-002 / D-003** 为准；Stitch 仅约束「目标观感」。  
   - F-002（Shadow 实施闭合）仍待代码证据。

### 为什么

- 无冻结方向则 S1–S3 易来回改审美；有冻结则实施可对齐可复核截图。  
- 本地 `raw/` 不进 git：以本决策 + 附件摘要作仓库内权威指针，避免第二状态源。  
- 双端列表与详情模式已在多轮 Stitch 中验证「可画」，适合写入 S2/S3 验收注记。

### 未选方案

| 方案 | 未选原因 |
|------|----------|
| 把 Stitch HTML 直接接入 apps/web | 破坏 Schema 驱动与单主线；非 VP-005 交付形态 |
| 将截图整夹 commit 进仓库 | 体积大且 raw 已 ignore；摘要 + 本地路径足够 |
| 因定稿勾选 S1 或宣称 exit 1–3 | 无运行时代码证据；违反过程诚实 |

### 闭合 I-005

| 日期 | 动作 |
|------|------|
| 2026-08-09 | 用户确认落盘 → D-004 **accepted**；I-005 **closed**（证据：本文件 + E-004） |

### 依赖

- D-001 纲领；D-002/D-003 Token  
- E-004 导出与评审事实  
- VP-005 v0.4.1；I-PROTO-FULL-001  
