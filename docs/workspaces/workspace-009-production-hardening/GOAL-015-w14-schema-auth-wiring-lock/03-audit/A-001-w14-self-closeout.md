---
title: A-001 · W14 关门前自审（self）
source: self
date: 2026-08-26
scope: GOAL-015 全部产物（hotfix 追认 / AuthGate 提取 / 生产装配回归锁 / 台账与树同步）
verdict: pass
status: active
created: 2026-08-26
updated: 2026-08-26
parent: null
version: 1.0.0
---

# A-001 · W14 关门前自审（self）

- **source**：`self` · **日期**：2026-08-26 · **verdict**：`pass`
- **模式依据**：见 [00-meta 审计模式记录](../../00-meta.md)——恢复既定鉴权语义 + 测试补充，security 高影响门禁语义不成立，不强制 independent。

## 1. 成果核对（逐项指回证据）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| hotfix 正确性：页面文档传输恢复认证语义，与 F-010 意图一致 | ✅ | D-001 §2；E-001 §5–7；匿名 401 实测 ×5 |
| 措施 B 为 verbatim 提取、无行为漂移 | ✅ | E-002 变更清单；全量 vitest 1130/1130 + tsc 0 兜底 |
| 回归锁真实覆盖生产装配点（而非注入式假阳性） | ✅ | E-002 设计要点：部分 mock 保真 + props 捕获真实 AuthGate 装配 |
| 未选方案有据、风险残余已登记 | ✅ | D-001 §3–4 |
| 文档同步：goal-tree 树+表、workspace.md 波次行 | ✅ | 两文件本波次更新记录 |

## 2. Findings

无 required。

- **R-001（recommended → 已按用户指令并入本波，闭合 = fixed）**：回归锁的 fetch 边界为 stub，未覆盖真实网络路径。**处置**：新增 `apps/web/e2e/schema-auth-transport.spec.ts`（真实 Chromium + 真实 API/vite，断言登录态全部 `/api/schema/*` 请求携带 Bearer 且 200、失败面不出现），实施与验证见 [E-003](../02-execution/E-003-w14-r001-e2e-bearer-smoke.md)。闭合路径：`fixed`。
- 其余观察：main.tsx 现仅承担 boot/boundary 职责，职责面收窄为正向副产物。

## 3. 变异验证（锁有效性硬证据）

1. 临时将 `AuthGate.tsx` 中 `schemaFetcher={authFetch}` 替换为占位注释；
2. 复跑 `auth-gate.wiring.test.tsx` → **1 failed / 1 passed**，失败点精确落在身份断言 `expect(props.schemaFetcher).toBe(authFetch)`（用例 2 分支冒烟不受影响，符合预期）;
3. 恢复接线 → 复跑 → **2/2 passed**。
   结论：若装配再次丢失认证传输，该锁必然先于生产变红。

## 4. 结论

S1–S3 事实齐备、证据可回指；开放 required = 0。**verdict: pass**。

**关门响应（2026-08-26 补记）**：R-001 已按用户指令并入本波并闭合 = `fixed`（[E-003](../02-execution/E-003-w14-r001-e2e-bearer-smoke.md)）；关门验证另暴露并修复 shell.spec 匿名 schema 探测陈旧契约（[E-004](../02-execution/E-004-w14-shell-spec-contract-fix.md)，定性为 F-010 后 e2e 未完整运行所致遗留破损，非本波回归）。全量证据与用户书面关门授权见 [D-002](../01-decision/D-002-w14-closeout.md) → 目标 `status: done`（4/4）。verdict 维持 `pass`。
