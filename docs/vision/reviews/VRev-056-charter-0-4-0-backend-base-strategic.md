---
id: VRev-056-charter-0-4-0-backend-base-strategic
doc_type: vision-review
title: Charter 0.4.0 Strategic · 同进程基座定位
source: self
date: 2026-08-31
scope: Charter strategic（目的陈述 / 成功边界 #6 / 非目标新条 / H-002）
verdict: pass
open_required: 0
status: active
created: 2026-08-31
updated: 2026-08-31
parent: null
version: 0.1.0
---

# VRev-056 · Charter 0.4.0 Strategic · 同进程基座定位

## 背景与触发

用户于 2026-08-31 会话中提出：下游 fork 项目大概率对外提供 C 端服务，而当前 Charter 仅描述"Admin 框架"，隐含 Admin 独立服务假设。若 C 端 API 须另建后端，基架对下游的价值将被显著削减。用户明确裁决为"B——同进程基座"形态，触发本次 strategic 修订（VR-052）。

## 审视范围

- Charter 目的陈述修订（补充"后端基座"定语与同进程预期）
- 成功边界新增第 6 条（基础设施端口开放 + 缓存/限流不预制原则）
- 非目标新增条款（不预制 C 端业务逻辑 + 双形态合法声明）
- 战略假设表新增 H-002（同进程基座冻结假设）
- re-align 负担评估（当前无 active VP）
- roadmap RT-Q03/RT-Q05 触发条件扩展（随配套 roadmap 更新）

## Verdict

**pass**

## Findings

### 必改（required）

无。

### 建议（recommended）

| id | finding | 建议 |
|----|---------|------|
| V-F092 | 成功边界 #6 表述"缓存、限流扩展、队列等有状态横切能力在业务域模块触发时由架构分支按需交付，不预制"——这句话在 Charter 层面方向正确，但未明确 roadmap 对应登记义务（RT-Q03/RT-Q05 须在首个业务域 VP 激活前完成触发条件评估和路线图声明，而非留作事后补登）。建议在 roadmap 更新中同步注明"业务域 VP 激活即视为触发条件成立，架构分支须在该 VP 开区前完成评估并登记路线图位置"。 | 配套 roadmap 更新时同步落实；可在本 Review 响应事务内 fixed。 |
| V-F093 | H-002 的"主要预期形态"措辞无量化标准，若后续 fork 均选独立部署，该假设将成为方向性误导。建议在 H-002 状态栏补充"如有具名 fork 选择独立部署形态，H-002 触发复核"以明确假设生命周期。 | 在 Charter H-002 状态栏追加复核触发条件；可在 roadmap 更新后配套修订。 |

## 说明

1. **Strategic 分类正确**：修订触及目的陈述的"Admin 框架 → 后端基座与 Admin 管理框架"措辞，以及成功边界的新增第 6 条，属于愿景成功条件定义的扩展，符合 strategic 门槛（至少 minor 版本升）。

2. **单愿景不变量**：仅一个 active Charter（`schema-ui-core-admin-foundation@0.4.0`），未破坏。

3. **Re-align 负担为零**：本次修订时无 active VP，所有已 closed VP 历史封存不变，re-align 无额外工作。

4. **内核代码天然支持 B 形态**：`kernel/Provider`、`Registrar`、`ContributionKeys` 接口均无路由前缀限制，业务域模块可直接以标准 module 契约注册 C 端路由，无需改动内核代码。Charter 修订与代码现状一致。

5. **缓存不预制的保留**：RT-Q03/RT-Q05 保持 trigger-gated 性质正确——同进程基座不等于必须预制 Redis；触发条件从"多实例"扩展为"多实例 **或** C 端业务域模块接入"，仍须业务触发而非预制。

6. **非目标双形态声明**：新增非目标条款明确"单独部署 Admin 服务仍是合法 fork 形态"，避免 Charter 被误读为强制同进程，保留 fork 选择权。

## 建议下一步

1. **响应 V-F092**：在 roadmap 更新（RT-Q03/RT-Q05 注记）中同步加入"业务域 VP 激活即视为触发条件成立"表述——可在本轮 roadmap 更新事务内 fixed。
2. **响应 V-F093**：在 Charter H-002 状态栏追加复核触发条件——可在 roadmap 更新后配套修订。
3. 完成 roadmap RT-Q03/RT-Q05 更新后，交 `/vision-audit` 调用 grok build 独立审视本次 strategic 修订。

## Finding 响应（2026-08-31 · self · 本轮事务内）

| finding | 响应 | 状态 |
|---------|------|------|
| V-F092 | roadmap RT-Q03/RT-Q05 备注栏已更新："业务域 VP 激活即视为触发条件成立，架构分支须在该 VP 开区前完成评估并在路线图中登记位置（可选「不需要」结论，但评估本身不可跳过）" | **fixed** |
| V-F093 | Charter H-002 状态栏已追加"复核触发：如有具名 fork 以单独部署形态为主要形态，H-002 须经 `/vision` 复核并决定是否修订" | **fixed** |
