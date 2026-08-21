---
id: A-002-grok-r4-c4-schema-migration-audit
doc: audit-entry
goal: GOAL-010-r4-c4-schema-other-migration
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: C4 settings/activity provider 化、Manifest 全 fragment、Schema owner 贡献驱动、门禁（secrecy/Ready/校验器）
audit_type: execution-facts
verdict: conditional
---

# A-002 · Grok R4-C4 Schema 与其他能力迁移独立交叉审计

## 声明

本意见 `source: independent`，只读。未修改任何文件、status/progress/checkpoint/
goal-tree/台账。响应与关门由 `/govern` 处理。

## 结论摘要

C4.1 settings provider 化、C4.2 activity 只读 provider 化、Manifest 全 fragment 化
（`adminModules` 全仓删除）、secrecy 与部分校验器成立，回归全绿（API `go test ./...`
+ `go vet` + Web 495）。**verdict: conditional**：C4.4 字面含 ledger drift/unknown 使
C4.4 勾选/关门受 F-IND-C4-002（required）阻断；C4.3 owner 仍为中心固定 contributor
表（F-IND-C4-001）；另有 recommended 若干。

## 核验成果（摘要）

1. settings provider：4 路由（branding 公开 + list/detail/patch 鉴权）、Schema、
   权限、导航、fragment；中心 RegisterSettings 清除（生产）。
2. activity provider：只读 2 路由、POST→405、writer 横切保持。
3. Manifest 基线核心专用（无 admin 页/导航）；`ForModulesWithFragments` = core +
   provider fragments；全仓无 `adminModules`。
4. C4.3 owner map 模块驱动（ModuleID/PageIDs + plan 门禁）——但 handler 内固定四模块
   枚举，非 runtime ContributionSet.Pages 驱动。
5. C4.4：secrecy `rejectFragmentSecrets` 成立；PolicyID/Visibility/JSON 校验器有实现
   （PolicyID/Visibility 过弱）；Ready 失败反向清理在 composition（缺双 Profile 矩阵）；
   ledger drift 诚实登记 C5 residual；Records historical-only 保持。

## Findings

### F-IND-C4-001 · Schema owner 仍中心固定 contributor 表（recommended；若宣称「完全贡献驱动」则 required）
- evidence: `handler/schema.go` 硬编码 users/roles/settings/activity contributor 切片
- closure: 从 ContributionSet.Pages 派生 owner，或用户 residual 接受模块常量 + 中心枚举

### F-IND-C4-002 · C4.4 含 ledger residual，不得无条件勾选/关门（**required**）
- evidence: `00-meta` C4.4 字面含 ledger drift/unknown；E-002 已诚实登记 C5 residual
- closure: fixed / accepted-residual / 缩窄 C4.4 成功标准并决策留痕

### F-IND-C4-003 · Ready 失败反向清理缺双 Profile 矩阵（recommended）
- closure: 补 mvp+admin Ready 失败清理测，或 residual 接受 composition 代码审边界

### F-IND-C4-004 · PolicyID/Visibility 校验器过弱（recommended）
- closure: allowlist/表达式语法并测，或决策记录「R4 最小 trim 规则」accepted residual

### F-IND-C4-005 · 中心 RegisterSettings/RegisterActivity 双路径仍在测试/legacy（recommended）
- closure: 测试改 provider finalize；C5/R6 终态删除

### F-IND-C4-006 · branding 路由未置 Public: true（recommended）
- closure: `SettingsRoutes` branding 置 `Public: true`

## 独立结论

| 问题 | 结论 |
|------|------|
| C4.1/C4.2 实施 | 可成立 |
| C4.3 完全贡献驱动 | 有条件成立（F-IND-C4-001） |
| C4.4 勾选/关门 | **不可无条件**（F-IND-C4-002 required） |
| GOAL-010 能否关门 | **当前否**（required 未闭合） |

**明确声明：本独立审计员未修改任何 status / progress / goal-tree / 文件内容。**
