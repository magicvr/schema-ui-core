---
id: D-001
doc: decision
title: S5 · 关门流程与证据矩阵口径
status: accepted
parent: GOAL-006-s5-evidence-and-closeout
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# D-001 · S5 关门流程与证据矩阵口径（2026-08-09）

## 决策

S5 阶段按 `C1 → C2 → C3 → C4` 顺序推进，范围纪律 = 只验证与收口，不新增业务功能、不重开已关闭阶段：

1. **C1 证据矩阵**：落 Root `attachments/`，**复用 F-V029 冻结表同一分母**（行 = `zh-CN`/`en-US` × `mvp`/`admin` × 匿名/已认证；列 = 固定 UI / 冻结 pageId-schema 并集 / M1～M4 / 权限正反例 / 缺失翻译 / 配置刷新 / 错误回退）。非 N/A 单元格必须有可核对证据路径（测试名 / 产物路径 / 捕获日志）；N/A 仅限 Profile 不可达单元格并注明模块边界（`admin.settings`/`admin.activity` 不在 mvp 模块集）。**禁止**以「代表性」收缩分母。
2. **C2 真实入口验证**：API 二进制经 `go build ./cmd/server` 真实构建后启动，`GET /api/branding` 断言响应体内容（siteTitle/supportedLocales/defaultLocale 等字段值）而非仅 200，同一构建启动 ≥2 次结果一致；`npm run build`（apps/web）成功；playwright 可用 → serve + 加载页面断言零页面错误、`lang` 随语种切换、一次设置保存产生可见变化并截图；不可用 → 捕获失败输出并走静态/结构性回退 + 单元测试为验收线。输出统一捕获到 `{SCRATCH}`。
3. **C3 关门独立审计**：由 grok CLI（`-m grok-4.5 --effort high`，提示词 = `skills/prompts/05-independent-audit.md`，即 `/audit`）对 S5 关门范围执行；意见落盘 GOAL-006 `03-audit/A-NNN-*` 并更新 `03-audit.md` 索引（`source: independent`）；required findings 按 P-003 三路径（fixed / accepted-residual / user-overruled）合法闭合后才可放行关门。
4. **C4 关门**：用户书面关门确认（P-004 留痕，含日期与范围）→ Root `status: done`、`progress: 6/6`；GOAL-006 `done 4/4`；VP-007 关门记录填写（outcome/summary/evidence_links/residuals）并按 alignment §7 置 `closed`；`goal-tree.md` 树与表最终同步；`workspaces.md` / `roadmap.md` 按既有惯例原子同步。

## 为什么

- Root D-002 §4 已冻结：S0 契约冻结与 **S5 关门 = `independent`**（grok CLI `-m grok-4.5 --effort high` 执行 `/audit`），常规阶段 `self` 兜底——本决策按此执行，不降级。
- VP-007 exit 6 要求证据矩阵复用「最小可枚举证据面」同一分母，Profile 不可达单元格标 N/A 并写明模块边界，不算 pass——矩阵口径直接沿用 F-V029，不另行定义覆盖面。
- alignment §7（规划关门轻量）：退出判据方向满足 + 证据链接指向工作区目标 done 路径 + 用户确认；本流程第 4 步即该轻量流程。
- P-004：关门属必须用户书面确认的裁决点；C4 的「用户书面关门确认（含日期与范围）」即 P-004 留痕。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 矩阵另行定义覆盖面（「代表性」行/列） | VP-007 exit 6 与 F-V029 冻结纪律明确禁止静默收缩分母 |
| S5 关门审计用 `self` 代替 independent | Root D-002 §4 已冻结 S5 = independent；P-003 禁止静默降级或编排器冒充 independent |
| 浏览器不可用时伪造启动证据 | 诚实降级：捕获失败输出 + 静态/结构性回退 + 单元测试为验收线 |
