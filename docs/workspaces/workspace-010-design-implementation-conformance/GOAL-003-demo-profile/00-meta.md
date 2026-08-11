---
id: GOAL-003-demo-profile
title: W2 · `demo` Profile：mvp + 范例页面
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-11
updated: 2026-08-11
version: 0.1.0
progress: 6/6
---

# GOAL-003 · W2 · `demo` Profile：mvp + 范例页面

## 概述

本子目标是 VP-010 / workspace-010 的**第二波**：把 W1 预留的「开发/dogfood 显式启用演示」正式化为一个**编译 Profile `demo`**。`demo` 的启用集 = **mvp 模块集 + `dev.examples`**：`APP_PROFILE=demo` 即可启动，展示 mvp 生产能力面（users/roles 等）**同时**带出 8 个协议范例页 + Examples 导航，home 指向 `overview`（`deriveHomePageRef` 对 `dev.examples` 启用已返回 `overview`，W1 D-003 §2）。

不改变 mvp/admin 生产向默认（仍不含 `dev.examples`）；不改变 Charter 目的/边界/非目标；`custom` 语义不变。

### 与 W1 的关系

- W1 已把范例面拆为可选模块 `dev.examples`，启用路径 = `APP_MODULES_ENABLED` 显式列表（D-002/D-003）。
- 本波把该路径收敛为一个**命名的演示 Profile**，降低开发者拼模块列表的心智负担，并使演示形态可测试、可文档化。

## 成功标准

- [x] **S1 · 编译 Profile**：`kernel` 新增 `ProfileDemo = "demo"` 常量与 `profileDefaults[demo]`（= mvp 集 + `dev.examples`）；`ResolveProfile("demo", nil)` 成功；`kernel_test`/`config_test` 分母更新（E-001）
- [x] **S2 · 产品面**：`APP_PROFILE=demo` 启动展示 mvp 能力面（users/roles）+ 8 范例页 + Examples 导航；`homePageRef` = `overview`（`TestDemoProfileManifest` + demo e2e；E-001）
- [x] **S3 · 卫生保持**：mvp/admin 默认仍**不含** `dev.examples`（不回归 W1）；`custom` 语义不变；W1 S5 卫生断言保持绿（`TestDemoProfileIsNonProduction`；E-001）
- [x] **S4 · 回归与烟测**：API 测试（demo 解析 + manifest/home 断言 + mvp/admin 回归绿）+ **demo 纳入 playwright e2e**（白名单放开 `demo`；浏览器烟测展示范例页 + users CRUD）（E-001）
- [x] **S5 · 文档**：README/QUICKSTART 标注 `demo` 为非生产向演示 Profile（用途 + 启用方式）（apps/api/README + 根 README + apps/web/README；E-001）
- [x] **S6 · go 接口**：新增 Profile = 模块矩阵变更 → 判定留痕：**mvp/admin 生产默认未变、`demo` 非生产向 → `go` 保持有效、不触发暂挂**；业务 VP 以 demo 为候选时触发 freshness（E-001 §go）

## 高层路线图（P-001）

1. **方案确认**：Profile id、e2e 白名单、文档落点。 **（2026-08-11 用户确认：`demo` / API+e2e demo / workspace-010 W2）**
2. **实施**：`kernel` ProfileDemo + `profileDefaults[demo]`；config/解析无回归；playwright 白名单放开 `demo`；web README/QUICKSTART 标注。 **（2026-08-11 完成 · E-001）**
3. **回归**：API 测试（demo 解析 + manifest/home）+ mvp/admin 回归 + demo e2e 烟测。 **（2026-08-11 完成 · E-001：go/web 全绿 + 三 Profile e2e）**
4. **go 影响留痕**：矩阵变更说明 + 触发/恢复判定。 **（2026-08-11 完成 · E-001 §go：不触发暂挂）**
5. **波次审计**：self（必要）+ 触及 Profile 矩阵 → `cross`/`independent`（grok-build@grok-4.5，P-004 已定模式）。 **（待办）**

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 新增 Profile 的 id 命名 | 方案 | 方案前 | 用户确认 | verified | — | `demo`（用户 2026-08-11 确认） |
| I-002 | required | 烟测范围（API 级 vs 纳入 e2e） | 方案/回归 | 方案前 | 用户确认 | verified | — | API + e2e demo（用户 2026-08-11 确认） |
| I-003 | non-blocking | 是否需 web 端 dogfood 入口/文档页（除 README 标注外） | 验收整洁度 | 验收 | 可 residual | open | deferred：README 标注已覆盖核心；复核=S5 | 可接受无独立入口 |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
