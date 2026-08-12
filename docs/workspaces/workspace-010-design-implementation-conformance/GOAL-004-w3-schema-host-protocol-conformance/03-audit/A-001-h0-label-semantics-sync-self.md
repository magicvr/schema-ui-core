---
id: A-001
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: 自审 · H0 标签语义同步（ADR-0034 D10 → 候选目录）
source: self
scope: attachments/I-HOST-APP-001-protocol-gap-catalog.md §1b/§1c/§6 同步
verdict: pass
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# A-001 · 自审：H0 标签语义同步（source: self）

## 审计范围与方法

- 范围：候选目录附件 `I-HOST-APP-001-protocol-gap-catalog.md` 新增 §1b（H0/S2 双语义映射）、
  §1c（95 项逐项上游处置对照）、§6（H0 同步状态说明）；上游 `schema-ui-docs` 提案 H0 勾选。
- 方法：机械比对脚本（ADR-0034 D10 表 ↔ 目录 §1c）、相对链接解析校验、语义核对。

## 核对结果

| 核对项 | 结果 |
|--------|------|
| §1c 处置取值与 ADR-0034 D10 逐字一致 | **pass**：95/95，0 差异（脚本比对） |
| ID 集合一致（§3=91 + §2.2 IMP=4 = §1c=95） | **pass**：无缺失、无多余 |
| §1b 跨仓相对链接可解析 | **pass**：解析至 `schema-ui-docs/docs/decisions/0034-host-app-interoperability-boundary.md` |
| H0/S2 三标签映射与 D10 前言表一致 | **pass** |
| 不以 deferred 冒充已保留 capability | **pass**：§1b 明确 `reserve-extension` 一律记“上游 deferred”；§6 明示 S2 门禁未达成、复选框保持未勾选 |
| 上游提案 H0 勾选与事实一致 | **pass**：勾选引用消费者 commit，描述与落盘内容相符 |

## Findings（审计中发现并已修复）

| # | 发现 | 处置 |
|---|------|------|
| F-1 | §1b 初稿把 95 个候选全部归属 §3，实际为 §3 的 91 个 + §2.2 的 IMP-001～004 | **fixed**：措辞修正为“目录全部 95 个候选（§3 的 91 个能力候选 + §2.2 的 IMP-001～004）” |
| F-2 | §1c 注记声称“处置取值与残余说明逐字对照裁定列”，残余说明实为压缩投影 | **fixed**：改为“处置取值逐字一致（95/95）；残余/对齐说明为裁定列的压缩投影，全文以上游 D10 为权威” |

## 未决行动项（非 required finding，随 S2 推进）

- GOAL-004 台账（00-meta / 01-decision）I-002 证据更新为“H0 处置已同步（ADR-0034 D10，proposed）”，
  状态保持 `collecting`，S2 冻结前待 ADR accepted 确认。
- I-004 independent provider 按用户指令指定为 `grok build`（grok 4.5，reasoning high），
  待独立审计条目落盘后闭合。

## 结论

自审未发现遗留阻断项；两处措辞问题已在同一变更集内修复。同步满足上游提案 H0 门禁第 5 条
（“不以 deferred 冒充已保留 capability”）。等待独立审计（A-002，source: independent）。
