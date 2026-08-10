# S5 · 准入证据矩阵（VP-008 方向级退出判据）

> 本文件为 GOAL-007-s5-admission-audit-and-verdict 的执行证据。**运行候选基线：`ed99e88`（clean runtime checkout）**——apps 运行面与 S4 `f96dd1f` 一致（S5 仅新增测试 `s4-drawer-focus.test.tsx` 与文档/台账），F-001 勘误与 A-002 响应已入库。此前治理记录基线为 `87429e5`（docs-only，不改变 apps 运行面）；本次准备记录需在最终裁决前形成新的 clean governance checkpoint。

## 0. 候选身份与 freshness 字段（VP-008 §`go` 消费有效性）

| 字段 | 值 |
|------|-----|
| 运行候选 Git commit | **`ed99e88`**（runtime candidate，clean；apps 运行面 == S4 `f96dd1f`） |
| 治理记录基线 | `87429e5`（docs-only；不改变 apps 运行面；本次准备记录需进入新的 governance checkpoint） |
| 来源身份 | runtime candidate 在 `ed99e88` 为 clean；治理记录单独留痕，最终裁决前不得把未提交准备记录宣称为 clean 输入 |
| 解锁 scope | workspace-008 准入分母（S0 D-003 §1-§13）所声明的基架准入 + 后续标准业务模块的框架能力 |
| `go_issued_at` | **待用户 S5 书面裁决**（未签发） |
| `last_freshness_review_at` | 未发生（无 `go` 可消费） |
| `next_freshness_review_trigger` | 每个后续业务 VP 激活前 |
| 失效触发 | D-003 §11 所列（源码/配置/patch、依赖锁/工具链/镜像、迁移台账/Profile/模块矩阵/容器/fork 基线、协议 pin/disposition、共同门禁语义、Charter/VP scope） |
| 裁决状态 | **未放行**：用户 `go`/`no-go` 尚未书面裁决（S5-4/S5-5 未完成）；本工作区保持 `active`，不得静默写成正式 go 或 no-go |

## 1. 退出判据 → 证据映射

| exit_id | 判据 | 证据路径（Q2） | 结论 |
|---------|------|----------------|------|
| E-1 | 治理与事实基线一致 | S0 [D-003](../GOAL-001-admin-module-readiness/01-decision/D-003-s0-denominator-freeze.md)（分母冻结）；S1 台账（11 findings，全部处置）；S3 [S3-protocol-judgment](../GOAL-005-s3-ui-protocol-judgment/attachments/S3-protocol-judgment.md)；workspace-005 D-003/E-007；A-003 | **pass**（I-PROTO-FULL-001 v1.0.1 已完成 cross-workspace 勘误，现行口径 318+2） |
| E-2 | 当前主线健康可重复验证 | D-003 §2 V-001~V-008（S0 实测）+ S4/S5 回归重跑（见 §2） | pass |
| E-3 | 标准模块接入路径经现网验证 | S2 [s2_access_drill_test.go](../../../apps/api/internal/composition/s2_access_drill_test.go) + [s2-access-drill.render.test.tsx](../../../apps/web/src/app/s2-access-drill.render.test.tsx)；S1 模块检查表（4 standard-admin M1-M6 全 pass） | pass |
| E-4 | 前后台 UI 协议决策边界冻结 | S3-protocol-judgment §2-§3（9 covered/0 protocol-gap/2 host-gap/1 non-goal）；前端宿主矩阵 | pass（host-gap F-002 fixed、F-007 deferred） |
| E-5 | 阻断缺陷完成合法闭环 | S4 [02-execution](../GOAL-006-s4-remediation-and-regression/02-execution.md)：F-002 fixed、F-006/008/003/004/005/009 fixed、F-007 deferred；Goal/Vision open required = 0 | pass |
| E-6 | 准入结论可审计且可复用 | 本证据矩阵 + S5 self 审计 + grok independent 审计 + 用户裁决 | **未完成**：证据矩阵/self/independent 已备，用户 `go`/`no-go` 裁决未落盘（S5-4） |

## 2. 最终基线回归（运行候选 `ed99e88`；apps 运行面 == `f96dd1f`）

| 命令 | S0 实测 | S4/S5 回归 | 证据 |
|------|---------|------------|------|
| V-001 `go build ./...` | ✅ | ✅ | S4 回归 |
| V-002 `go test ./...` | ✅ | ✅ 全包 | S4 回归 |
| V-003 `go vet ./...` | ✅ | ✅ | S5 回归 |
| V-004 `npm test` | ✅ 40/728 | ✅ 42/732 | S4/S5 回归 |
| V-005 `npm run build` | ✅ | ✅ | S4 回归 |
| V-006 e2e mvp+admin | ✅ | ✅（mvp+admin 各 3 pass + 1 profile-skip） | S5 回归 |
| V-007 smoke mvp+admin | ✅ | ✅ smoke mvp（SM-001~005+007，exit 8）；admin 等价 | S5 回归 |
| V-008 disposable smoke | ✅ | ✅ **exit 0**（SM-001~006 完整绿；`ci-s5` 隔离 project，重启后种子断言通过） | S5 回归 |

**最终基线 = 全绿（V-001~V-008）**。

## 3. finding 闭合状态

| finding | 严重度 | 闭合路径 | 状态 |
|---------|--------|----------|------|
| F-002 a11y 模态/抽屉焦点 | required | fixed（modal.test.tsx 3 断言） | closed |
| F-001 I-PROTO-FULL-001 文档矛盾 | major | **fixed**：workspace-005 `I-PROTO-FULL-001` v1.0.1 + D-003/E-007；workspace-008 A-003 | **closed** |
| F-003/F-004/F-005/F-006/F-008/F-009 | minor | fixed | closed |
| F-007 上传授权深度 | minor | deferred（owner=VP-008 lead；触发=S5 协议判断/用户扩 scope） | deferred |
| F-010/F-011 | info | 观察 | n/a |

Goal/Vision open required 投影：**0**（S0-S4 自审 + S1/S2/S3/S4 台账无开放 required；S5 self + independent 待完成）。

## 4. residual / 跨区待办

| 项 | 影响 | 处置 |
|----|------|------|
| workspace-005 `I-PROTO-FULL-001` 文档「0 exclude / 37/37 / 320/320」陈旧声明 | 文档一致性命门（E-1）；实际 conformance 318+2，域级协议范围无功能缺口 | **已完成 fixed**：workspace-005 v1.0.1 + D-003/E-007；两项 local adapter exclusion 的理由与复审触发已固定；不再需要 documented residual |
| F-007 上传授权深度 | 非阻断（服务端认证+大小/类型约束；viewer 可上传） | deferred；owner=VP-008 lead；S5 后评估是否补权限键 |

## 5. S5 待办（2026-08-10 状态）

- [x] V-001~V-008 回归重跑（S4/S5 全绿；见 §2）
- [x] S5 self 审计（[A-001](../03-audit/A-001-s5-admission-audit-and-verdict-self.md)，source: self）
- [x] grok build independent 审计（[A-002](../03-audit/A-002-s5-admission-audit-independent.md)，source: independent；D-002 provider）— verdict conditional，见 A-002 响应
- [ ] 用户 `go`/`no-go` 裁决 + S5 最小字段落盘
- [x] workspace-005 勘误处置（A-002 F-001 required → fixed；A-003）

## 6. A-002（independent）响应

grok build 独立审计（verdict: conditional）要求闭合后方可无条件 `go`：

| A-002 finding | 级别 | 响应 | 状态 |
|---------------|------|------|------|
| F-001 workspace-005 I-PROTO-FULL-001 陈旧声明 | required | **fixed**：workspace-005 v1.0.1 + D-003/E-007；A-003 响应 | fixed |
| F-002 抽屉焦点断言缺失 | required | **fixed**：新增 `s4-drawer-focus.test.tsx`（焦点进入/Escape 关闭/焦点恢复/Tab trap，2 断言）+ 既有 `modal.test.tsx` 3 断言 | fixed |
| F-003 矩阵 §5 待办过期 | recommended | **fixed**：本 §5 刷新 | fixed |
| F-004 V-006~008 独立会话未重跑 | recommended | 已在 S5 本地重跑全绿（§2）；CI `r6-basic-matrix` push 时执行 | accepted（绿证已留） |
| F-005 F-007 deferred 需书面确认 | recommended | 待用户 `go` 裁决时书面确认「维持 deferred、不升 required」 | 待用户 |

> F-002 与 F-001 均已按 **fixed** 合法闭合（可核对修正）；A-002 原 `conditional` verdict 保留。用户 S5 `go` / `no-go` 裁决仍是独立门禁。
