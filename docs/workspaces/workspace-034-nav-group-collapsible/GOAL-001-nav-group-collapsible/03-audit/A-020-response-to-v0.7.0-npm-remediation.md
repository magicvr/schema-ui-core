---
id: A-020-response-to-v0.7.0-npm-remediation
goal_id: GOAL-001-nav-group-collapsible
doc: audit-entry
source: self
auditor: current-session
type: response
scope: response to A-019 independent audit of v0.7.0 npm remediation before PR
date: 2026-09-21
verdict: pass
created: 2026-09-21
updated: 2026-09-21
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# A-020 · 响应 v0.7.0 npm 修复 independent 审计（2026-09-21）

- **source**：self（编排器响应）
- **类型 / scope**：response · 汇总 A-018 self 与 [A-019 independent](A-019-v0.7.0-npm-remediation-independent.md)，核对合并后发现的真实 npm 消费故障修复、PR 前门禁与发布后顺序。
- **verdict**：**pass**（开 PR 前 open required = 0；不代表 Hosted CI、修复包发布或 Release 已完成。）

## 审计意见响应

| finding | 响应 | 状态 |
|---|---|---|
| A-018 F-S-002：finalize / ESM 包面和真实 tarball consumer 缺少门禁 | 接受修复。alias rewrite 现在自动运行 finalizer，compat normalizer 有单测，CI smoke 对六个本地 tarball 在 OS 临时消费者中 npm install 后做 13 项导入。A-019 独立复跑确认 8/8 单测、13/13 导入与 protocol schema tarball 资产。 | `fixed`（A-018 / A-019） |
| A-019 F-001：Hosted PR CI 未运行 | 接受。PR 建立后必须等待新增 web job（build + finalize + isolated tarball consumer smoke）与其余 required checks 全绿，之后才合并；本地审计不替代 Hosted 结果。 | adopted；merge 前门禁 |
| A-019 F-002：修复 patch 必须先于 v0.7.0 tag 可见 | 接受并更新 A-016 F-001。合并后等 main CI 全绿，发布 protocol `0.2.17` / lib `0.1.16` / renderer `0.3.15` / ui `0.1.13`；逐包 registry 可见并通过干净消费者验证后，再建 annotated `apps/api/v0.7.0` tag 和 GitHub Release。 | adopted；发布顺序门禁 |

## 执行状态

- PR #16 已合并；此修复候选尚未创建 PR。
- A-019 `pass` 且 open required = 0，可进入创建修复 PR 阶段。
- `A-019` 指出 A-018/E-016 先前使用「tarball 解包」表述；已在这两条 self 记录中改为准确描述：临时消费者对六个本地 npm tarball 执行 `npm install`，并安装 peers / dependencies。此措辞更正不改变独立审计意见。
- 修复版 npm patch、tag、GitHub Release 尚未发生；不得用 dry-run 或本地 tarball smoke 标成已发布。

## 结论

A-019 的 `pass` 与 F-S-002 `fixed` 有独立验证支持。A-019 F-001/F-002 均采纳为后续门禁与执行顺序。本响应不修改 Goal/workspace/VP 状态，不更新 goal-tree。
