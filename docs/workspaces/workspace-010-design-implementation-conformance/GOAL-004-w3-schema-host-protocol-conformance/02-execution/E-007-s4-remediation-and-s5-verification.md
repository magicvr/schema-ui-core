---
id: E-007
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: S4 残余整改完成 + S5 符合性验证证据
status: recorded
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# E-007 · S4 残余整改完成 + S5 符合性验证证据

## 已完成事实

### 1. S4 工作清单完成（D-002 §2）

| # | 项目 | 结果 |
|---|------|------|
| S4-1 | return intent 登录流接入 | `apps/web/src/host/return-intent.ts`（捕获/单次消费/nonce/过期/`/login` 自循环拒绝）+ `boot.ts` reauth 动作捕获 + `AuthProvider` 登录后 `history.replaceState` 恢复（ADR-0036 D6）。测试 `return-intent.test.ts`（6 用例） |
| S4-2 | session adapter reauth-required 映射 | `auth-client.restoreSession` 判别式返回（none/reauth/session）；`AuthContext` 增加 `reauth-required` 终态；`bootstrapAuthFor` 映射 reauth-required（ADR-0035 D4/D7）。测试 `boot.test.ts`（2 用例）+ `auth-client.test.ts` reauth 用例 |
| S4-3 | hostOwnedPaths 显式集合 | `App.tsx` 声明 `["/login"]`，认证态命中 host-owned path 导航回 home（ADR-0036 D3a），不再产生 `HOST_ROUTE_NOT_FOUND` |
| S4-4 | `$deps` residual 纠错 | 引擎已 `e18edce` 实现（stage3 reactions 套件零排除）；`generate-claim.mjs` residuals 文本更新；claim 重生成 |
| S4-5 | 304/ETag | 维持 200-only 装载（ADR-0035 D6 conditional GET 为「可用于」可选优化）→ 无动作 |
| S4-6 | account-locked 生产源 | 映射层保留、fixtures 已 pin；生产源缺位为拟议 residual（D-002 §3，S6 用户 P-004 决策） |
| S4-7 | IMP-002 导航 label 单源固化 | Go 断言测试 `TestNavigationSingleProjectionWithLabelKey`：served manifest 导航单一投影（labelKey `manifest.nav.users` + label fallback），provider `NavigationContribution.Label` 不泄漏进 manifest 文档 |

### 2. Claim 重生成（S4-4）

`node apps/web/scripts/generate-claim.mjs` → `buildId=git:5e4c3848dee57763461a9a52a511d14c57bcbb94`，
canonical digest `sha256:60f75039b3256a970d10f45624cb523f6be0e746ce7218f58f38606a32bb4273`；
`conformance-claim.json` residuals 现为空（`$deps` residual 随纠错关闭），`contentSha256` 仍为正式
`4fae4605…`、`fixtureSha256` `7aacf133…`。

### 3. S5 符合性验证证据（全绿）

| 门禁 | 结果 |
|------|------|
| Web vitest | **48 文件 871 通过** |
| TypeScript | `tsc --noEmit` 0 错误 |
| Go 内部套件 | `go test ./internal/...` 全 ok（manifest/handler/composition/modules/users/roles/settings/schemarender/store） |
| 上游 fixtures | host 三 suite **96 零排除**（23+43+30）、app-manifest 41、app-navigation 16、version-negotiation（stage3 266）零排除 |
| Claim 门禁 | `claim-artifact.test.ts` 5/5 |
| 浏览器级（Playwright） | **7 通过 + 1 既有 skip**（host-failure 4、localization 1+1skip、schema-crud 1、shell 1） |
| 旧协议兼容证据 | version-negotiation（2.7 接受/拒绝 + 2.8 向量）+ app-manifest 2.7/2.8 严格协商 fixtures + migration `2.7-to-2.8.md` 双轨矩阵（I-005/E-006） |
| 代表性页面 / auth / bootstrap / shell / error 流程 | `representative-pages.integration.test.tsx`、`startup-config.test.tsx`、`shell.test.ts`、`host-failure.spec.ts`、`schema-crud.spec.ts` 覆盖 |

## 阻塞 / 风险

无阻断。唯一拟议 residual：account-locked 生产源缺位（S4-6），S6 关门时点用户 P-004 书面决策。

## 关联信息项

无新增开放 required。I-001～I-006 均 `verified`；I-007 维持 `collecting`（non-blocking）。

## progress

S4（实现整改）与 S5（符合性验证）检查点完成 → **5/6**（S1+S2+S3+S4+S5）。S6 关门待执行。
