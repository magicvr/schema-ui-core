---
id: A-009-r4-matrix-self
doc: audit-entry
parent_goal: GOAL-001-nav-group-collapsible
source: self
auditor: /govern（schema-ui-core 编排器）
type: execution-facts
scope: R4 当前 sidebar 全量迁移、Profile/slot/route matrix、optional/custom/demo 回归与 Playbook
date: 2026-09-07
verdict: pass
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# A-009 · R4 当前 sidebar 全量迁移与矩阵自审

## 范围与证据

本自审覆盖 E-008 与附件 `r4-navigation-route-profile-matrix.md`：

- 17 个当前 sidebar NodeID 的五组成员与 route；
- Dashboard 顶层单例、top/user slot、通知铃面、Examples authored group；
- `mvp`、`admin`、`demo`、custom Telegram、custom Digital Offer + Telegram 的真实 Manifest assembly；
- 所有已登记内页/动态路径的 Web 父链接 active/组自动展开；
- module-contribution-playbook v1.2.0 的 group 规范；
- API/Web 全量回归。

## 事实核对

1. API R4 matrix 真实经 composition mux 读取发布 Manifest，五组成员与 Profile 组合符合 D-005/D-004；禁用 optional 模块不产生空组。
2. Web R4 matrix 覆盖所有当前 grouped NodeID 的直接 route，并覆盖 users-invites、wallet-entries、dictionary-entries、task-runs、telegram-operator 深链父级激活。
3. Dashboard、user slot 与 Examples 组在 API 与 Web 断言中保持独立；Settings/Account/My wallet 没有被搬进 sidebar。
4. Playbook 已包含 group key 命名、同 key 精确一致、跨模块共组、未分组兼容、slot 边界、迁移和验证要求。
5. 未发现新的 required/recommended finding；I-034-004/I-034-005 的证据已足以关闭其 R4 信息门禁。

## 验证结果

- API `go test ./... -count=1`：通过。
- API R4 matrix `go test ./internal/composition -run TestR4NavigationGroupProfileMatrix -count=1`：通过。
- Web 全量 Vitest：99 files / 1339 tests 通过。
- Web R4 route matrix 2/2、R3 interaction 3/3、navigation 7/7：通过。
- Web `tsc -b`：通过。
- `git diff --check`：通过。

## Findings

本 self scope 无开放 required/recommended finding。R5 仍需形成最终证据矩阵、执行最终回归、关门审计与 required 意见汇总；本条不将 Root 标记 done。

## 结论

**verdict: `pass`**。R4 当前已注册 sidebar 全量迁移与 Profile/slot/route 回归达到检查点完成条件，可同步 Root 4/5；R5 关门准备仍未完成。

## 声明

本意见为 `source: self`，R4 是兼容性回归与可逆 UI 交付，不另触发 independent provider；关门阶段按 R5 另行确定 cross 审计。
