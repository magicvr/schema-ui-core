---
title: D-002 · W14 关门记录
status: done
created: 2026-08-26
updated: 2026-08-26
parent: null
version: 1.0.0
---

# D-002 · W14 关门记录

日期：2026-08-26 · 决定：`GOAL-015-w14-schema-auth-wiring-lock` → **`status: done`（4/4）**

## §1 用户书面关门授权

用户指令原文：「**把 r-001 并入本波。处理完再关门。**」——即以「R-001 并入并完成」为关门前置条件的有条件书面授权；条件已满足，本记录落盘关门。

## §2 R-001 并入（闭合 = fixed）

- 实施：`apps/web/e2e/schema-auth-transport.spec.ts`（真实 Chromium + 真实 API/vite，登录态全部 `/api/schema/*` 请求断言 Bearer 头 + 200，失败面双语标题 count=0），见 [E-003](../02-execution/E-003-w14-r001-e2e-bearer-smoke.md)；
- A-001 §2 R-001 闭合路径：`fixed`。

## §3 关门验证证据（全绿）

| 验证面 | 结果 |
|--------|------|
| 生产装配回归锁（vitest） | `auth-gate.wiring.test.tsx` 2/2（含变异验证红→绿，A-001 §3） |
| web 全量单测 | vitest **1130/1130**（84 文件）、`tsc -b` exit 0 |
| 浏览器 e2e 全量 | Playwright chromium **10 passed / 1 skipped（admin 档用例于 mvp 档按设跳过）· exit 0**，含新 Bearer 冒烟与全部既有 spec |
| 关门验证附带修复 | `shell.spec` 匿名 schema 探测陈旧契约对齐 F-010（401），见 [E-004](../02-execution/E-004-w14-shell-spec-contract-fix.md)——定性为 F-010 落地后 e2e 未再完整运行所致的遗留破损，非本波回归 |

## §4 开放项清点

- required findings：0；
- 信息门禁：无未决 required 信息项；
- 残余移交：无新增。git checkpoint 未提交（工作树含本波全部改动），由用户决定提交时机。

Root 保持 active；VP-009 无 go 宣称变动需要（本波为装配缺陷修复 + 测试加固，未触及 Profile/模块矩阵/协议 pin 语义）。
