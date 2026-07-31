---
title: I-008 · R6 阶段 3 · VP 证据汇编（三条退出判据 → Q2 工作区证据）
status: active
doc_type: vp-evidence-assembly
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-008-r6-integration-acceptance-vp-evidence
version: 0.1.0
related_info: I-008-001
---

# VP 证据汇编 · VP-001 三条退出判据 → 工作区 Q2 证据

> **性质**：阶段 3「VP 证据汇编与缺口整改」的登记产物。把 [VP-001 方向级退出判据](../../../vision/plans/VP-001-mvp-admin-foundation.md#方向级退出判据) 三条逐条指向当前工作区可复核证据；本表**支持提出 R6/VP 决定，不自动改变 Goal 或 VP status**。证据来自阶段 2 验收（revision `a941bedb1fc2cd4859a408df50653e867da35ff2`）与 R1–R5 历史记录。

## 逐条映射

### VP-001 判据 1 · React + Go 可运行、可 fork、固定协议边界

| 主张 | 工作区 Q2 证据 | 说明 |
|------|---------------|------|
| React Web 与 Go API 可构建、可测试 | `GOAL-008-*/attachments/evidence/acceptance/results/{web-test,web-build,api-test,api-build}.log`（15 files / 395 tests；Vite build；`go test ./...`；`go build ./...` 全 pass） | evidence-index C-001/C-002 |
| 双服务可启动、health、Web→API proxy 成立 | `runtime-probes.log`（`/healthz` 200、`/` 200、`/api/accounts/me` 200、`/api/records` 200、`/api/records?sort=INVALID` 400、manifest 200） | evidence-index C-003 |
| 浏览器关键路径成立 | `browser-e2e.json`（1 passed / 0 unexpected）+ `r6-overview.png` 截图；`apps/web/e2e/shell.spec.ts` | evidence-index C-004 |
| 干净安装 + Linux/CI 等价 | `.github/workflows/r6-basic-matrix.yml`；CI run `30666932343`（api/web/browser-e2e 三 job success） | D-004；I-008-002/005 |
| 固定协议边界 | `docs/vision/protocol-inventory-v2.7.0.md`（`schema-ui-docs@v2.7.0` @ `ca9e5fe…`）；Charter H-001 ①；`apps/web/src/protocol/app-manifest.ts` 的 `APP_MANIFEST_SOURCE` | 协议 pin 可复核 |

### VP-001 判据 2 · 受控协议清单定义 MVP 覆盖范围，每项有实现/范例/验证路径

| 主张 | 工作区 Q2 证据 | 说明 |
|------|---------------|------|
| 受控覆盖清单冻结 | `GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md` v0.1.3（Root D-009 冻结） | include/include-partial/exclude 边界 |
| 每纳入域可追到实现/范例/验证入口 | `GOAL-007-r5-examples-contract-verification/attachments/I-007-001-registry.md`（逐域登记 + §2b include-suite 执行矩阵） | R5 登记 v0.8.1 |
| 集成级回归不扩大边界 | `GOAL-008-*/attachments/evidence/acceptance/results/web-test.log`（stage3 conformance 222 项，含 request-construction non-batch 64）；`R6-acceptance-plan.md` §2b C-007 | evidence-index C-007 |
| 结构 conformance（Ajv schemas） | `docs/schemas/*`（6 schema vendor + SHA pin）；stage3 structural 断言 | I-PROTO-004=vendor |

### VP-001 判据 3 · 核心账号权限前后端集成，不依赖未声明业务模块

| 主张 | 工作区 Q2 证据 | 说明 |
|------|---------------|------|
| 账号上下文从 API 到 Web 可观察 | `runtime-probes.log`（`/api/accounts/me` 返回 dev-001 session：`roles:[admin,editor]`、`features.beta:true`）；`browser-e2e.json`（E2E session 断言）；`apps/web/src/account/context.ts` + `main.tsx` | evidence-index C-005 |
| 权限求值 / 拒绝路径 oracle 已登记 | `GOAL-008-*/attachments/account-permission-oracle.md`（正向 P-1～P-4、拒绝 D-1～D-6） | I-008-003 |
| D-PERM fixtures 可复核 | `GOAL-006-r4-account-permission/attachments/dperm/cases.json`（17 例，SHA-256 pin 于 `permissions-inheritance.test.ts`）；web-test.log 内 `permissions-inheritance.test.ts` pass | R4 复用 |
| 不依赖未声明业务模块 | 依赖边界：Web 仅 React + Tailwind/shadcn 风格组件；API 无外部业务模块；evidence-index exclusions 明确 D-UPLOAD 等排除 | I-PROTO-001 v0.1.3 |

## 缺口与残余检查

- **浏览器级拒绝路径未断言**（C-006）：真实 manifest 无权限门控导航项，拒绝以 renderer/组件层 oracle（D-1～D-6）断言。属**登记排除**（evidence-index exclusion），非 required 缺口；复审触发：若真实 manifest 新增权限门控项。
- **reactions multi-round / request-construction batch / D-UPLOAD**：既有冻结排除（D-008 / D-010 / v0.1.3），不构成缺口。
- **无其他 required 缺口**：三条退出判据均有可复核 Q2 证据；无需用户书面 residual/overruled。

## 结论

三条 VP-001 退出判据均已有工作区 Q2 证据指向；无未闭合 required 缺口；边界主张（include-partial/exclude）与 I-PROTO-001 v0.1.3、R5 登记一致。本汇编支持向 `/vision` 提出 VP-001 关门提案输入；是否关门由 `/vision` + 用户确认决定。
