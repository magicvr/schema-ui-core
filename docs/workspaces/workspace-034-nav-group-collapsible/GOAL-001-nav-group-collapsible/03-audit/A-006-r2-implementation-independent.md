---
id: A-006-r2-implementation-independent
doc: audit-entry
parent_goal: GOAL-001-nav-group-collapsible
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
type: execution-facts
scope: R2 分组注册与 Manifest 聚合契约实施（A-005 self / E-004）
date: 2026-09-07
verdict: pass
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# A-006 · R2 分组注册与 Manifest 聚合契约实施独立审计

## 范围与区间

本意见由项目指定的本地 grok build（grok-4.6 · reasoning high）只读复核：

- 工作区 `workspace-034-nav-group-collapsible`、Root、D-004、A-005、E-004；
- `apps/api/kernel` `NavigationGroup` / finalize 冲突校验与 17 个 Provider；
- `internal/manifest` `NormalizeSidebarGroups` 与 tests；
- `internal/composition` 与 `apps/api/server/serve.go` 两条 assembly；
- composition / digitaloffer / telegram 测试；Web translation keys、strict protocol 边界与当前 diff。

本意见未修改任何文件、status、progress、goal-tree 或方案正文；未把 R3 折叠/自动展开写成已完成。

## Verdict

**`pass`**；开放 required = **0**。

D-004 的 Dashboard/五组/未分组/Examples 混排、同 key 精确全等、sidebar-only fail closed、协议/checksum 兼容、optional 模块聚合、两条 assembly 一致，均可在代码与测试中独立核对。A-005 self 的 pass 结论成立。

## 成果（有证据）

1. `NavigationGroup` 仅含 `Key/Order/Label/LabelKey/Icon`，不含 NodeID；同 key 冲突在 kernel finalize 以 `CodeModuleNavigationGroupConflict` fail closed；注释明确不改变 menu_items、授权或 system-data checksum。
2. 17 个当前 sidebar Provider 的五组声明真实存在且元组一致：identity-access×3、content-data×2、operations×4、communications×3、commerce×5。Dashboard、top/user、通知铃面与 Examples authored group 未标 Group。
3. `NormalizeSidebarGroups` 真实实现：Dashboard 第一；结构化组按 GroupOrder；未分组叶子随后保持相对顺序；已有协议组最后保留；不输出内部 key/id。测试断言 `Dashboard, Identity, Content, Commerce, Future, Examples`。
4. 分组 NodeID 不在 sidebar 叶子中会 fail closed；真实 user-slot 节点未标 Group。
5. kernel 比较的结构不含 NodeID；Manifest `sameNavigationGroupingMetadata` 只比较 GroupKey/GroupOrder/GroupLabel/GroupLabelKey/GroupIcon。
6. `NormalizeSidebarGroups` 只改 sidebar；navigation checksum/ensureNavigation 字段集未扩展；现行 strict NavGroup schema 不新增 key/id；协议仍 2.7。
7. optional digitaloffer/telegram 的 Provider 元组与五组一致；真实 custom mux assembly 成功，禁用模块不在 set.Navigation，不产生空组。
8. composition 与 serve 都调用 `ForModulesWithFragmentsAndGroups(..., GroupingsFromContributions(set.Navigation))`，没有生产 assembly 分叉。
9. R3 未实施：App.tsx/navigation.ts/protocol parser 不在本轮 diff；Shell 仍静态渲染组，无折叠、键盘或直接 URL 自动展开。
10. en-US/zh-CN 均有五个 `manifest.nav.group.*` labelKey。

## 对照成功标准

| D-004 / R2 标准 | 结论 |
|---|---|
| Dashboard 第一、五组显式顺序 | 成立 |
| 未分组叶子不丢失、Examples 保留 | 成立 |
| 同 key 精确全等、冲突 fail closed | 成立 |
| sidebar-only、top/user 不误改 | 成立 |
| Group 不进 Parent/menu_items/checksum | 成立 |
| strict protocol 不新增 key/id | 成立 |
| optional telegram/digitaloffer 入组、未启用无空组 | 成立 |
| 两条 assembly 一致 | 成立 |
| R3 不越权宣称 | 成立 |

## 独立验证事实

本独立会话复现：

- `git diff --check`：exit 0（仅 Windows CRLF 提示）；
- `go test ./kernel ./internal/manifest -count=1`：pass；
- R2 相关 composition / optional module mux tests：pass。

本独立会话未复跑全量 `go test ./...`、直接 Vitest 97/1332、`tsc -b`、`vite build`；这些验证已由 E-004/A-005 记录，本条不把它们冒充为独立复跑。

## Findings

### F-001 · optional 入组与默认组内顺序的回归钉不足

- 级别：recommended，状态：open。
- digitaloffer/telegram 测试递归收集 pageRef，但没有直接断言节点位于 commerce/communications group；若误删 Provider Group，部分测试仍可能通过。现行实现和 assembly 证据成立，不是 required 缺陷。

### F-002 · kernel 冲突测试未逐字段覆盖完整元组

- 级别：recommended，状态：open。
- 现有 kernel 测试覆盖 Order 冲突，Manifest 测试覆盖 Label 冲突；D-004 的 key/order/label/labelKey/icon 全字段等价语义由实现保证，但缺少逐字段单测与一致元数据成功路径单测。

### F-003 · 治理索引/notes 需要与 E-004/A-006 合并状态同步

- 级别：recommended，状态：open。
- 独立扫描时发现部分历史/当前投影仍保留“R2 待开始”或“F-002～F-004 开放”的文案。若编排器已在本轮后续同步，这一 finding 只需以当前文件证据核对并记录 fixed；不要修改历史 A-002/A-003/A-004 原文。

## 信息门禁

| ID | 级别 | 最晚阶段 | 本轮判定 |
|----|---|---|---|
| I-034-001 | required | R1 | verified（静态分母），不阻断 R2 |
| I-034-002 | required | R1 | 决策与 Manifest 行为均有证据 |
| I-034-003 | non-blocking | R2 | open，折叠状态选择留到 R3 |
| I-034-004 | required | R3 | 未到期，自动展开仍未实施 |
| I-034-005 | required | R4 | 未到期，完整 Profile harness 仍属 R4 |

共享资料：`none`，未作为事实或关闭证据。

## 与 A-005 self 的异同

- 同意：R2 实施符合 D-004；17 项五组声明真实；同 key 比较排除 NodeID；sidebar-only、checksum/协议/top-user、两条 assembly、optional 模块与 R3 边界均成立。
- 补充：三个 recommended（optional 入组回归钉、逐字段 kernel 测试、治理投影同步）；未提出 required。
- 无 verdict 冲突；A-005 `pass` 与本条 `pass` 同向。

## 结论与建议

R2 代码契约可以按 D-004 视为已实施。建议 `/govern`：

1. 落盘本意见；required 无需整改。
2. 响应 F-001/F-002：可作为 R2 收尾测试增强，或记录为 R4 harness 前的有界 recommended residual；不得把它们升级成 R3 阻断。
3. 响应 F-003：更新当前 goal-tree、execution、decision 投影，但不改历史审计原文。
4. A-002 F-006/F-008/F-009/F-010 已由 D-004/E-004 约束覆盖，可由编排器另条响应标记 fixed。
5. 合并意见后创建 Git checkpoint，再进入 R3；I-034-003 仍需 R3 前做折叠状态存储选择。

## 声明

本意见为 `source: independent`，auditor: grok-build (grok-4.6 · reasoning high)。响应由 `/govern` 处理。
