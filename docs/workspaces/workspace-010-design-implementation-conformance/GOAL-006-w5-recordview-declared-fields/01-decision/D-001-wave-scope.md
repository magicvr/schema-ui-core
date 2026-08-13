---
id: D-001
doc: decision-entry
goal: GOAL-006-w5-recordview-declared-fields
status: accepted
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-001 · 波次范围：recordView declared-fields 契约 + fail-open + dev/文档卫生

## 决策

本波（VP-010 W5）范围：

1. **declared-fields 契约**：renderer 的 recordView 按 schema 声明的字段元数据渲染（字段标题 / 顺序 / 包含集）；users / roles / activity schema 增加声明；i18n 双语标题同步。
2. **fail-open 兜底**：声明缺失 / 异常时回退默认渲染（不崩、不黑屏），并修正健壮性缺口（提交 `a831754` 语义）。
3. **dev/文档卫生**（同批未归档提交）：dev 脚本等待 API ready 后启动 Web、stop 按 PID 精确停止（`5c309ff`）；QUICKSTART 修正 dev.cmd 前缀与排版（`c420e5d`）。

**边界**：不改 Profile 默认集 / 模块矩阵 / Manifest 装配 / 协议 pin / 共同门禁语义；declared-fields 为本地渲染器契约（I-001 verified），不新增协议 capability。

## 审计模式

**self**（低风险、可逆、已提交；回归以 HEAD 冻结矩阵为证据）。

## 未选方案

- 把 declared-fields 上升为协议 capability：超出本波范围；协议侧无对应定义（I-001），本地契约 + fail-open 足够。
- 单独为 dev 脚本卫生开目标：与波次同批提交，归入本波 S3，避免过度拆目标。
