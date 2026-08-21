---
id: A-006-grok-r4-c3-cutover-audit
doc: audit-entry
goal: GOAL-009-r4-c3-users-roles-migration
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: C3.3 composition 消费 provider finalize、users/roles 中心特例清除、冻结包 §7 步骤 3
audit_type: execution-facts
verdict: pass
---

# A-006 · Grok R4-C3.3 中心特例清除独立交叉审计

## 声明

本意见 `source: independent`，只读。未修改任何文件、status/progress/checkpoint/
goal-tree/台账。响应与关门由 `/govern` 处理。

## 结论摘要

相对 GOAL-009 C3.3 成功标准与冻结包 §7 步骤 3 的 users/roles 切换意图，**C3.3 成立**：
生产 mux 经 `RegisterContributions` 挂载 users/roles 路由；中心 `handler.Register`
不再挂载；Schema 内容迁入模块包；Manifest 基线移除 users/roles + `ForModulesWithFragments`
合并 provider fragment；`adminModules` 仅 settings/activity；无永久双路径；API
`go test ./...` + `go vet` + Web `vitest run`（495 测试）通过。未将 C3.4 宣称为完成。
**开放 required finding：0。可进入 C3.4。**

## 核验矩阵（摘要）

1. composition 消费 RegisterContributions 挂载 users/roles 路由：pass
2. 中心 Register 无 users/roles 分支：pass
3. 生产无重复注册/双路径：pass（MountProviderRoutes 仅测试助手）
4. Schema 内容模块所有 + 中心 fixture 已删：pass
5. schemaDocumentsForPlan owner map 仍门禁（残余·已登记）：pass（文档化）
6. Manifest 基线无 users/roles：pass
7. ForModulesWithFragments + provider fragment：pass
8. adminModules 仅 settings/activity：pass
9. Provider fragment 非 stub（闭合 A-003 C32-002）：pass
10. 冻结 §7 步骤 3 字面对照：pass（owner map 删除有 recommended 残余）
11. 行为矩阵 C3.3 保真（浅层）：pass（深度属 C3.4）
12. API/Web 测试：pass
13. 未误宣称 C3.4：pass

## Findings（recommended，不阻断 C3.3）

### F-IND-C33-001 · Schema owner map 仍硬编码 users/roles（文档化残余）
- level: recommended · severity: low · status: open
- evidence: `handler/schema.go:73-86` owner map；`00-meta`/`E-005` 已登记
- closure: C4/后续改为 provider/schema 贡献驱动，或用户书面 accepted-residual

### F-IND-C33-002 · MountProviderRoutes 与生产 finalize 并存（测试旁路）
- level: recommended · severity: low · status: open
- evidence: `health.go:40-55`、`testhelpers_test.go:64-66`
- closure: C3.4 后测试助手改 provider finalize 挂载，或标注 tombstone + lint 门禁

### F-IND-C33-003 · composition 贡献注册失败时 ModuleID 固定为 admin.users
- level: recommended · severity: low · status: open
- evidence: `composition.go:105-106`
- closure: 按失败 provider 填 ModuleID 或中性聚合错误码

### F-IND-C33-004 · C3.3 生产路径行为证据偏浅（C3.4 门禁）
- level: recommended · severity: med（对 C3.4）· status: open（延至 C3.4）
- evidence: `modules/{users,roles}/provider_test.go` 仅 401/200/404
- closure: C3.4 在 provider finalize 生产路径补行为矩阵 + 双 Profile + 失败注入复审

## 独立结论

| 问题 | 结论 |
|------|------|
| C3.3 检查点是否成立 | **是**（证据充分可复核） |
| 是否可进入 C3.4 | **是**（无开放 required 阻断） |
| C3.4 是否已完成 | **否**（检查点未勾选） |
| Schema owner map 残余是否阻断 C3.3 | **否**（目标已登记；若要按冻结字面删除须 residual 或 C4 前 fixed） |

**明确声明：本独立审计员未修改任何 status / progress / goal-tree / 文件内容。**
