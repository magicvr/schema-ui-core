---
id: GOAL-034-w23-admin-login-home-redirect
title: W23 · admin 登录后 home 推导回归修复（N-001 承接）
status: done
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-001-design-implementation-conformance
version: 0.3.0
progress: 4/4
---

# GOAL-034-w23-admin-login-home-redirect · W23 · admin 登录后 home 推导回归修复

## 概述

承接 GOAL-033 A-001 **N-001**（显式引用，非自动继承）：admin profile 登录成功（含首登改密完成）后停留在 `/` 而未跳转 `/dashboard`，`e2e/localization.spec.ts:62` 稳定失败。GOAL-033 基线实验（stash 全部 W22 改动后 HEAD 同败）当时判定为「先于 W22 的 W14–W21 路由漂移」。用户指令：立项本波承接排查修复。

**根因结论（D-001，2026-08-23）**：重定向链**从未断开**。真实断点在登录第一步——本机 gitignored `apps/api/configs/.env`（2026-08-21 建立，`DB_DIALECT=postgres`）经 `config.Load` 的 env-file 层把 e2e 挂具本应隔离的临时 SQLite 静默改写为开发者 Postgres；全新种子 `admin/admin` 直接 401，界面必然停留在 `/`。W22 基线实验 `git stash` 无法移除 gitignored 文件，其「先于 W22」结论不成立。修复 = 挂具显式钉死 `DB_DIALECT=sqlite` + 登录/菜单两处等待竞态硬化（连带 F-1/F-2）。

## 成功标准（达成，A-001 pass）

1. admin（mvp 不回退）登录/首登改密后落地 `/dashboard`，`localization.spec.ts` M1 断言恢复绿 —— **达成**（凭证级 200 全链 + e2e admin 连续 5 轮 9/9）；
2. 根因以代码级证据落盘（D-001），含「为何此前各波全量绿未拦截」的说明 —— **达成**（密闭单测无甄别力；`.env` 建于 W21 窗口，此后本地 e2e 全部在污染下运行）；
3. 补一条防回退测试（home 推导分支或 e2e 断言）—— **达成**（既有 fresh-seed 强制改密 + `/dashboard` 断言组为强制回归护栏，隔离失效第一触点 = `force-password-change.spec.ts` 全新种子 401；D-001 D5 记录复审触发）；
4. 全量回归绿：vitest + go + tsc + build + e2e admin（+ mvp）—— **达成**（E-004 矩阵）。

## 路线图（P-001 · 分母 = 4）

```text
S1 冻结   → C1 根因定位 + 修复方案（D-001）
S2 实施   → C2 修复 + 防回退测试
S3 回归   → C3 全量回归绿（含 e2e admin localization）
S4 关门   → C4 关门审计 + 台账/goal-tree 终态同步
```

## 信息需求登记（P-005）

| 编号 | 问题 | 级别 | 影响门禁 | 最晚需要阶段 | 收集动作 | 状态 | 备注 |
|------|------|------|----------|--------------|----------|------|------|
| I-001 | `/` → `/dashboard` 重定向链在哪一层断开（deriveHomePage / StampHomePageRef / 路由守卫 / manifest home 字段）？W14–W21 何波引入？ | required | S2 实施 | C1 前 | 读 `deriveHomePage*`、`App.tsx` 登录后导航路径、admin-dogfood/upstream manifest 的 home 字段；必要时二分波次 | **closed（D-001）** | 重定向链**未断开**：根因 = 本机 gitignored `configs/.env`（2026-08-21 建立）`DB_DIALECT=postgres` 劫持 e2e 挂具的临时 SQLite，全新种子 `admin/admin` 直接 401，界面停在 `/`；非 W14–W21 代码回归 |

## 边界与审计声明

- 仅修重定向回归；不顺手重构路由。
- 关门审计默认 self；若修复触及协议/manifest 契约语义则升级 independent。
