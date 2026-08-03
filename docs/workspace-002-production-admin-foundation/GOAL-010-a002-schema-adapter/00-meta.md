---
id: GOAL-010-a002-schema-adapter
title: A-002 · Schema 驱动通用数据适配层
status: active
created: 2026-08-03
updated: 2026-08-04
parent: GOAL-001-production-admin-foundation
version: 0.6.0
progress: 3/5
---

# GOAL-010 · A-002 · Schema 驱动通用数据适配层

## 概述

承接 [Root A-002 审计响应](../GOAL-001-production-admin-foundation/03-audit.md) F-002-001（required / high）：当前 Renderer/CRUD 硬编码单一 records 实体（`schema-table.tsx` 无条件调用 `fetchRecords`；`records.ts` 固定 `id/name/status/owner/updatedAt` 解析形状），不满足 VP-002 产品级成功标准 1/4/6「通过修改 Schema 新增业务页面，无需修改前端 Renderer 主路径」。本目标将表格/表单 transport、字段模型与 response mapping 提升为 Schema 驱动的通用适配层，并提供明确的后端资源契约；records 成为该适配层的注册实例。

## 路线图 / 成功标准

五个检查点默认等权、原则上串行（P-001；Root D-014 用户裁决走通用适配层改造，不降级 VP-002 主张）。**串行偏差留痕（2026-08-03）**：S3 为纯前端、不依赖 S2，因 A-001 F-001/F-002 的关闭证据在 S3（dataSource/rowKey 正反测试）而先于 S2 实施；**S2 已于同日补做，串行顺序已恢复（S1/S2/S3 全勾选）**。

- [x] **S1 · 资源契约与方案冻结**：定义 Schema 驱动的通用资源契约——`dataSource` 资源标识、字段模型、response mapping、后端通用 CRUD 端点/注册形态与错误 envelope 扩展边界；冻结方案（决策 + 附件契约，解除 `I-010-001` 门禁）。
- [x] **S2 · 后端通用资源 CRUD**：按资源契约提供通用 CRUD 入口（records 作为已注册资源），保持 `records.read` / `records.write` 权限键与现有错误 envelope；`go test ./...` 全绿。
- [x] **S3 · 前端通用适配层**：`schema-table` / 表单 transport 与 response mapping 通用化，去除 `RecordItem` / `RecordList` 固定解析（records 降为泛化实例）；web `vitest run` 全绿 + `tsc -b` / 生产构建干净。
- [ ] **S4 · 语义化双实体验证**：子目标 `GOAL-011-s4-semantic-admin-resources` 完成；`users` 替换 records 默认代表实体、`roles` 作为第二个语义资源，二者在后端资源注册完成后仅通过修改前端 Schema 接入列表/CRUD 页面（不修改 Renderer 主路径）；records 从当前产品默认运行面按版本化兼容策略退场。
- [ ] **S5 · 回归、审计与关闭**：全量回归（api + web + build）+ 阶段/关门审计；Root A-002 F-002-001 关闭证据经 `/audit` finding-closure（或 self + 独立复核）确认后按 `fixed` 闭合。

阶段子目标：`GOAL-011-s4-semantic-admin-resources`（active / 5/5；A-012 required findings 仍 open，当前 fail-closed）承载 S4 的语义资源、records 退场与双实体验证；其原 S1～S5 检查点事实只用于评估本目标 S4，不自动勾选本目标 S4 或关闭 Root finding。

## 派生进度

`progress` 由 S1～S5 五个检查点等权派生（`0/5` 起）。检查点不替代审计 finding 或关门结论。

## 信息需求与阶段门禁

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-010-001` | 通用资源契约的精确形状（字段模型、response mapping、端点注册形态、错误 envelope 扩展边界） | required | S1 方案冻结与 S2 实施 | S2 首个实施变更前 | 对照现有 records 契约与 A-002 建议关闭路径，形成版本化适配层契约附件，提交用户裁决 | **verified** | 已关闭（D-002）；v0.2.0 为 A-001 F-001/F-002 响应修订（D-003）；v0.2.1 为 D-004 S4 交接附注；v0.2.2 为 D-005（响应 GOAL-011 A-002 F-005）账号域 409 限定扩展注记，均维持 verified | [I-010-001-schema-resource-contract.md](attachments/I-010-001-schema-resource-contract.md) v0.2.2：S1～S3 技术契约不变；`catalog` 降为 genericity 历史示例，S4 产品终态交由 GOAL-011 `users + roles` 与 records 退场契约承接；§5 账号域 409 由 I-011-001 限定扩展 |
| `I-010-002` | 向后兼容与迁移策略（现有 records fixture/API/权限键在通用化后的迁移或双轨形态） | required | S2/S3 实施 | S3 首个前端变更前 | 评估 records 注册为实例的迁移边界与回归影响，记录决策 | **verified** | 已关闭（D-002，提前于最晚阶段） | 契约 §6：后端零对外变更收敛；前端一次性泛化（无双轨）；fixture/emulator/测试形状保持 |

> 未关闭的 required 信息项不得实施受影响范围；允许先进行收集与方案冻结（Root P-005）。

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)（A-002 响应；Root D-014；VP-002 成功标准 1/4/6） |
| 阶段子目标 | [GOAL-011-s4-semantic-admin-resources](../GOAL-011-s4-semantic-admin-resources/00-meta.md)（users + roles、records 退场、双资源 S4 证据） |
| In | 通用 transport / 字段模型 / response mapping、后端资源契约、S1～S3 records 泛化历史、GOAL-011 语义双实体验收、回归与审计 |
| Out | 完整业务模块；扩大 `I-PROTO-001 v0.1.3` 覆盖；F-002-002/003（归 GOAL-009）；Root / VP-002 关门（Root 层面独立裁决） |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)
