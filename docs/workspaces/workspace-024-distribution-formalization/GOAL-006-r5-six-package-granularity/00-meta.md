---
id: GOAL-006-r5-six-package-granularity
title: R5 · 六包形态细化（renderer 依赖图 external 化 + ui 纯原子断言 + 冻结面 v1.4.0）
status: active
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
progress: 0/5
---

# GOAL-006 · R5 · 六包形态细化

## 概述

承接 Root R5 与 VP-024 判据 #5/#6：把 renderer 从「自包含 bundle」升级为**依赖图 external 化**（产物对 `@schema-ui/{protocol,lib,ui}` 的 import 显式化 + peer 声明，消费端解析）；`ui` 包**纯原子拆分断言**（原子组件 vs 业务组件边界 + 独立消费）；六包 **exports 子路径 + files 收窄 + 版本推进发布**；**冻结面升格 v1.4.0**（六包导出面 + peer 矩阵定稿）。核销 go 后清单「renderer external 化 + 纯原子拆分」两项。

## 成功标准（可验证检查点）

- [ ] C1：renderer 包 external 化：重建产物 `index.js` 对 `@schema-ui/{protocol,lib,ui}` 的 import 断言存在（静态探针）+ package.json 声明对应 peerDependencies；`@/i18n → @schema-ui/lib/i18n`、`@/protocol → @schema-ui/protocol`、`@/lib → @schema-ui/lib/lib`、`@/components/ui → @schema-ui/ui/components/ui` 重写表生效
- [ ] C2：六包 exports 通配子路径（`"./*": "./*"`）+ files 收窄为实际目录 + 版本推进（renderer 0.3.0 · protocol 0.2.1 · lib 0.1.1 · ui/theme/shell 0.1.1）→ npmjs 实发
- [ ] C3：ui 纯原子断言（`components/ui` 无 renderer/protocol/i18n 反向依赖）+ ui 包独立消费实证（probe-six 既有 + 收窄后重装）
- [ ] C4：冻结面 v1.4.0 定稿（六包导出面 + peer 矩阵）落 `attachments/`
- [ ] C5：golden-field 升级（六包新版 + renderer 包子路径探针断言）全绿；独立审计（grok）→ 关门（Root 5/7）

## 方案与路线（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | 设计定档：重写映射表 + exports/files 布局 + 版本推进（D-001） | 未开 |
| S2 | 构建链改造：build-lib-packages.mjs 增 renderer external 构建（重写插件 + external 函数）· tsc declaration 重建 | 依赖 S1 |
| S3 | 六包 pack + 新版本 npmjs 发布 + golden-field 升级 + 探针（含 renderer import 断言） | 依赖 S2 |
| S4 | 冻结面 v1.4.0 定稿 + 独立审计（grok）→ 关门 | 依赖 S2/S3 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | renderer 对内部面包化面的完整 import 清单（重写表覆盖性） | C1 | S2 | 静态扫描 `src/renderer/**` | open | S2 前闭合 | 待确认（勘察已得：i18n 16 · protocol 13 · ui 9 · lib 3） |
| I-002 | non-blocking | 旧包（renderer 0.2.0 自包含）消费者的兼容口径 | 发布 | S3 | changelog 注记（external 化 = 消费契约变化 → minor bump） | open | — | 待确认 |

## 父目标

- `GOAL-001-distribution-formalization`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- 审计模式：S4 关门 = independent（grok build 先例）。
- 版本语义：external 化改变 renderer 消费契约 → minor（0.2.x→0.3.0）；纯元数据包 additive → patch（0.1.1）。