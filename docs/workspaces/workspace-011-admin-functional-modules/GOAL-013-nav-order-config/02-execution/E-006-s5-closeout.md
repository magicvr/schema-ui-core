---
id: E-006
goal: GOAL-013-nav-order-config
date: 2026-08-14
status: recorded
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-006 · S5 关门（A-003 响应）

## 事实

- 2026-08-14：grok 独立审计（A-003）verdict conditional，1 条 required（F-001）。用户裁决 + 全部 findings 响应后关门。

## Findings 响应

| Finding | 级别 | 处置 | 证据 |
|--------|------|------|------|
| F-001 kernel 排序未写入 menu_items.sort_order | required | **accepted-residual（用户书面裁决）**：产品权威序 = 仅 manifest；sort_order 保持模块声明 Order + 运营者手改契约；复审触发 = 出现读 sort_order 的导航消费者 | A-003 响应节 + E-003/A-002 表述修正 |
| F-002 缺 manifest 序集成测试 | non-blocking | **fixed**：TestPublishedManifestNavigationOrder（默认序 + env 覆盖序） | composition_test.go |
| F-003 Notifications 清单位 vs 空 fragment | non-blocking | **accepted**（既有产品形态，非本目标引入） | A-003 响应节 |
| F-004 重复/大小写边界无测 | non-blocking | **accepted-residual**（行为符合非法回退语义；运维文档注明精确匹配） | A-003 响应节 |
| F-005 web fixture 旧序 | non-blocking | **accepted**（静态测试数据非 API 输出） | A-003 响应节 |

- 回归：composition 全量 + 新集成测试绿。

## 关门条件

- A-003 required（F-001）已按用户裁决合法闭合（accepted-residual，范围 + 复审触发书面留痕）；其余 non-blocking 已处置。S5 关门成立。
