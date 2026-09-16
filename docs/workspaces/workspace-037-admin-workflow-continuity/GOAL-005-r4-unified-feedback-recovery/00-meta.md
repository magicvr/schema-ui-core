---
id: GOAL-005-r4-unified-feedback-recovery
title: R4 统一反馈与恢复
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.4.0
progress: 4/4
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-005 · R4 统一反馈与恢复

## 概述

在 R1 D-005 反馈/恢复语义和 R2/R3 已完成的 Saved View、dirty-state 基础之上，补齐 Admin 现有 schema-driven 页面在成功、失败、显式重试、维护/不可用和无障碍呈现上的统一证据。目标是让用户知道发生了什么、哪些动作可安全恢复，以及哪些写入失败必须保留草稿等待用户再次操作。

## 范围与边界

- 统一现有 `FeedbackRegion`、列表/指标/图表数据错误和表单/行操作结果的分类、文案优先级、诊断信息边界与无障碍角色。
- 幂等读取可提供用户触发的显式 retry；retry 只能发起新的读取，不自动重复写入、不形成 retry loop。
- 写入成功给出一次成功反馈；客户端校验、fieldErrors、409/403/5xx、离线/超时等失败保留表单值和 dirty 状态，不自动重发。
- 普通资源的 maintenance/unavailable/offline/timeout 与现有 HostFailureScreen 的认证、协议、渲染终态保持边界；维护不伪装为业务成功。
- 不新增 API、错误 envelope、第二持久化栈或全局状态管理；不把 R5 组合验收、全量业务域覆盖或敏感诊断输出纳入本目标。

## 高层路线图

1. **C1 · 分类合同与信息就绪**：核对 HTTP status、错误 code/messageKey/correlation 与 Host kind 的现有事实，冻结安全用户文案和恢复策略映射。
2. **C2 · 共享反馈表面**：收敛 Toast、列表错误/retry 与表单错误的共享可访问呈现，确保 dismiss/retry 键盘可达且同一结果不重复宣告。
3. **C3 · 跨页面恢复回归**：覆盖读取 retry、成功/失败/fieldErrors、维护/不可用/离线边界、dirty 保留和不重复提交。
4. **C4 · 审计与检查点**：完成 self audit；按风险发起 independent audit，响应所有 required finding；建立 Git checkpoint 并投影 Root R4。

## 成功检查点

- [x] C1：错误分类、catalog 文案、诊断 code/correlation 和恢复策略形成可核对合同；required 信息 R4-I-001 已 verified。
- [x] C2：共享反馈表面覆盖成功/错误、dismiss/retry、role/status/alert 与键盘路径；R4-I-002～003 已 verified。
- [x] C3：schema page 的读重试、写失败保留、maintenance/unavailable/offline/timeout 分类与 Host 边界有列表/指标/图表/表单回归及既有 Host fixtures 证据；R4-I-004 已 verified。
- [x] C4：self + independent audit、required finding 响应与 Git checkpoint 完成；R4 关闭并投影 Root。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| R4-I-001 | required | HTTP status、resource envelope 与 Host failure kind 如何映射到安全用户文案和恢复策略？ | C1/C3 | C1 | 读取 `resource.ts`、`host/failure.ts`、D-005，建立纯函数/矩阵 | verified | 2026-09-17：`feedback-policy.ts` 分类矩阵、E-002、E-004；A-003 待 independent recheck | status/code、Host 边界和 AbortError/网络错误的安全 catalog 文案已形成可核对合同 |
| R4-I-002 | required | 哪些 retry 是幂等读取，如何阻止写入重复提交和 retry loop？ | C2/C3 | C2 | 盘点 DataTable、display data、form/action 调用链，补重复调用测试 | verified | 2026-09-17：E-002/E-003、A-001 | 仅读取路径暴露显式 retry；写入无自动重试且有回归证据 |
| R4-I-003 | required | toast、列表错误、表单错误的 role、focus、dismiss/retry 和宣告去重如何保持可访问？ | C2/C3 | C2 | 共享组件/Renderer 集成测试，核对既有 Host failure 规则 | verified | 2026-09-17：`FeedbackNoticeView`、E-002、A-001 | role/status/alert、键盘路径、dismiss 与双击 guard 已验证 |
| R4-I-004 | required | maintenance/unavailable/offline/timeout 在普通资源与 Host 终态之间的边界及恢复动作是什么？ | C3/C4 | C3 | 对照 D-005、`HostFailureScreen`、resource/display fetch 回归 | verified | 2026-09-17：E-003、E-004、A-002；A-003 待 independent recheck | 普通资源 policy 与 Host/page schema 终态保持边界；transport catch 已保留 timeout/offline，retry 仅作用于幂等读取 |
| R4-I-005 | non-blocking | correlation/code 是否需要更丰富的用户可展开诊断细节？ | UX/支持体验 | R5 或真实支持需求 | 仅显示安全 code/correlation；敏感原始 payload 不进入 UI | deferred | owner=`/vision`；出现支持需求时复核 | D-005 安全边界 |

## 父目标

- `GOAL-001-admin-workflow-continuity`（Root 当前 `active · 4/5`；R1/R2/R3/R4 已完成，本目标承载 R4 阶段）。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。
