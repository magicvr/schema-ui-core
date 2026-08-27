---
title: D-001 · W14 范围与方案冻结
status: active
created: 2026-08-26
updated: 2026-08-26
parent: null
version: 1.0.0
---

# D-001 · W14 范围与方案冻结

日期：2026-08-26 · scope：本目标全部产物 · 决策性质：实现方案取舍（用户指令「补上避免【测试装配 ≠ 生产装配】的测试。连同本次修复一起，在工作区9新增一个子目标承载治理上下文」）

## §1 决定

**双措施闭环**，缺一不可：

| 措施 | 内容 | 状态 |
|------|------|------|
| A（hotfix 追认） | `main.tsx` AuthGate 装配补 `schemaFetcher={authFetch}`——页面文档请求恢复与其它 API 一致的 Bearer / refresh-on-401 传输 | 已实施（用户确认；见 E-001） |
| B（防复发） | `AuthGate` 原样提取为 `apps/web/src/app/AuthGate.tsx`（main.tsx 仅保留 boot/boundary 并回导 `BootScreen`）；新增 `src/app/auth-gate.wiring.test.tsx` 生产装配回归锁：挂真实 `AuthProvider` + 部分 mock（`restoreSession`、`@/app/App` prop 捕获），断言生产 gate 实际传入 `<App>` 的 `props.schemaFetcher === authFetch`（真实模块引用），并覆盖 resourceFetcher / navigationContext / onLogout / currentUser 接线与未登录分支 | 已实施（见 E-002） |

## §2 为什么是这两条

- F-010 之后 `/api/schema/{pageId}` 是认证端点；页面文档传输层必须带 Bearer，这是**语义修复**，不是行为变更。
- 该缺陷能上线，是因为既有测试全部显式注入 fetcher（34 处 `schemaFetcher=`），生产入口的装配点没有任何测试可达。只有把装配点本身变成被测对象（措施 B），同类缝隙才从结构上不可再隐身。

## §3 未选方案（有据否决）

| 未选方案 | 否决理由 |
|----------|----------|
| 源码文本正则断言 main.tsx 含 `schemaFetcher={authFetch}` | 不验证行为、与格式强耦合、重构即碎；对「传错 transport」（如换成裸 fetch）不敏感 |
| 测试直接 import `main.tsx` 再断言 | 入口 import 即执行 `createRoot().render()` 与 bootstrap discovery 副作用，需 mock react-dom/client + 整条 host boot 链，复杂度高且仍难精确触达 AuthGate→App 的装配点 |
| 只修 bug 不加锁（维持现状） | 「30+ 测试注入 fetcher」的结构性假阳性未消除，下一次入口装配改动（或新增认证端点）会以同样方式静默回归 |

## §4 风险与残余

- 措施 B 为 verbatim 搬移 + 导入调整：以 `tsc -b` exit 0 与全量 vitest 1130/1130 兜底（E-002）。
- 回归锁的网络层真实性有限（fetch 边界为 stub）：真实网络路径由后端契约测试与 e2e 承担，登记 recommended R-001（A-001），不阻断关门。
  - **2026-08-26 更新**：R-001 已按用户指令并入本波实施（e2e Bearer 冒烟，见 [E-003](../02-execution/E-003-w14-r001-e2e-bearer-smoke.md)），A-001 R-001 闭合 = `fixed`；本节残余清零。
