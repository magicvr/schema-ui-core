---
id: A-002
goal: GOAL-013-nav-order-config
source: self
date: 2026-08-14
scope: S2/S3 实施与验证
verdict: pass
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-002 · self 审计（S2/S3）

## 结论

**verdict: pass**（E-003/E-004）。

## 核对

- 默认清单 12 项冻结为 kernel 常量 + 快照测试（I-003 维护规则落地：新模块入清单 + 测试更新）。
- 排序语义（清单优先 → 未列追加末尾 → Parent 分组优先）与 D-002 §3 一致；三层都有测试。
- manifest 聚合层排序补齐（实施发现 UI 载体）；kernel 与 manifest 共用 NormalizeNavigationOrder，非法覆盖一致回退。**A-003 F-001 修正（用户 accepted-residual）**：产品权威序 = 仅 manifest；kernel 排序不保证 system-data 落库顺序（sort_order 保持模块声明 Order + 运营者契约）。
- 覆盖载体（YAML navigation.order + NAVIGATION_ORDER env）与 W7 优先级链一致；宽松解析回退 + 告警。
- 实测三场景（默认 / env 覆盖 / 非法回退）全部符合 D-002。

## Findings

- 无 required。
- 备注（非必改）：web 静态 fixture（app-manifest.admin.json）仍保留旧顺序，是测试数据非 API 输出快照；若未来需要 fixture 与 API 顺序对齐可另行同步（S5 审计可复核）。
