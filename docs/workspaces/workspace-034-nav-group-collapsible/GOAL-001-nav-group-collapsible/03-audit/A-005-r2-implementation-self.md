---
id: A-005-r2-implementation-self
doc: audit-entry
parent_goal: GOAL-001-nav-group-collapsible
source: self
auditor: /govern（schema-ui-core 编排器）
type: execution-facts
scope: R2 分组注册与 Manifest 聚合契约实施
date: 2026-09-07
verdict: pass
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# A-005 · R2 分组注册与 Manifest 聚合契约实施自审

## 范围与证据

本自审覆盖 E-004 的实施事实、D-004 冻结契约与 R2 相关回归：

- kernel `NavigationGroup`、Group key/metadata 校验、finalize 冲突 fail-closed 与稳定错误码；
- 17 个当前 sidebar Provider 的五组声明；
- `internal/manifest` group projection/explicit GroupOrder/sidebar-only/mixed output；
- composition 与 `server/serve.go` 两条 Manifest assembly；
- API/Web/Build 验证。

证据路径：E-004、D-004、`apps/api/internal/manifest/manifest_test.go`、`apps/api/kernel/provider_test.go`、`apps/api/internal/composition/composition_test.go`、`apps/api/internal/composition/composition_digitaloffer*_test.go`。

## 事实核对

1. Group 是可选 presentation metadata；未声明 group 的贡献仍可通过 Manifest 平铺。
2. 同 key 元数据比较不含 NodeID，仅比较 D-004 冻结的 key/order/label/labelKey/icon；冲突在 kernel finalize fail closed。
3. Manifest 不输出内部 key/id；Dashboard、结构化组、未分组叶子、既有 authored group 的顺序符合 D-004；top/user 不参与。
4. `menu_items`/system-data checksum 未扩展，Group 不改变权限/角色授权身份。
5. composition 与 serve assembly 均接收 `GroupingsFromContributions`，未出现单路径遗漏。
6. 本轮未改 Shell 折叠逻辑，因此不把 R3 成功标准提前写成已完成。

## 验证结果

- API：`go test ./... -count=1` pass。
- Web：直接 Vitest `97/97` files、`1332/1332` tests pass；直接 `tsc -b` pass；直接 `vite build` pass。
- `git diff --check` pass。

## Findings

本 self scope 未发现新的 required/recommended finding；R2 代码没有越过 D-004 的协议兼容边界。R3 的 F-005 recommended（内页/动态路径自动展开范围）仍按原 A-002 响应保持 open，不属于 R2 放行阻断。

## 审计模式与下一步

该 scope 仍是跨 API kernel、Manifest、serve assembly 与 Web consumer 的 compatibility cross scope。按 D-002/D-004 与项目级独立审计决策，本条 self 之后调用本地 grok build（grok-4.6 · high）执行 independent implementation audit；在 independent 结果合并前不将 R2 检查点标为完成、不创建 checkpoint commit。

## 声明

本意见为 `source: self`，不修改目标 status/progress，不替代 independent 意见。
