---
id: A-002-s5-admission-audit-independent
doc: audit-entry
goal: GOAL-007-s5-admission-audit-and-verdict
source: independent
auditor: grok-build / grok-4.5 (high) · /audit skill
verdict: conditional
audit_type: close-out
scope: workspace-008 Root GOAL-001 S5 准入裁决（compatibility / data / migration / production-release / 跨边界治理语义 / 跨模块 UI a11y 下限）
created: 2026-08-10
updated: 2026-08-10
version: 1.0.0
parent: null
---

# A-002 · S5 准入审计与裁决 · independent 交叉审计（grok build）

## 范围与区间

| 字段 | 值 |
|------|-----|
| 工作区 | `workspace-008-admin-module-readiness`（`workspace.md` 已校验：`root_goal=GOAL-001-admin-module-readiness`，`canonical_scope=docs/workspace-008-admin-module-readiness/`，`shared_materials_catalog=none`，`primary_plan=VP-008-…`） |
| 被审目标 | `GOAL-007-s5-admission-audit-and-verdict`（承接 Root S5） |
| 审计模式依据 | Root [D-003 §12](../../GOAL-001-admin-module-readiness/01-decision/D-003-s0-denominator-freeze.md) + [D-002](../../GOAL-001-admin-module-readiness/01-decision/D-002-independent-audit-provider-grok-build.md)（provider = grok build · grok-4.5 · high · `audit`；模式 `cross`） |
| 准入分母 | Root [D-003](../../GOAL-001-admin-module-readiness/01-decision/D-003-s0-denominator-freeze.md) §1–§13 |
| 代码候选（apps 运行面） | `f96dd1f`（S4 整改后）；当前 `git HEAD` = `3769253`（S5 文档/台账准备，**apps/ 相对 f96dd1f 无 diff**） |
| 来源身份（本会话核对） | `git status --porcelain` 空；`f96dd1f` 为 `HEAD` 祖先 |
| 对照 self | [A-001](A-001-s5-admission-audit-and-verdict-self.md)（source: self，verdict: pass） |
| 本意见角色 | **独立审查者**；不修改 status / progress / goal-tree / 方案正文 |

覆盖核对面：

1. **compatibility** — 协议 pin/disposition、模块/Profile 矩阵、M1–M6/接入演练、V-001~V-005  
2. **data** — SQLite 内嵌、bootstrap/种子、单租户 N/A、system-data reconcile  
3. **migration** — 账本 0001–0010、ownership、连续版本/checksum 收集规则  
4. **production/release** — compose fail-closed、Manifest 非静态兜底、smoke/disposable、`go` freshness 字段  
5. **跨边界治理语义** — workspace 绑定、无共享资料误用、跨区 residual 处置路径  
6. **跨模块 UI a11y 下限** — D-003 §8 + F-002 关闭证据  

---

## 成果（有证据）

### 工作区 / 治理绑定

- `workspace.md`：Root、canonical、VP-008 plan_refs/primary_plan、delivery 角色、shared materials = none —— 一致。  
- `goal-tree.md`：S0–S4 done、GOAL-007 active；progress 派生展示存在但不作为放行证据。  
- 未读取或混用其他工作区状态；对 workspace-005 文档仅作 **Q2 证据路径核对**（陈旧声明是否存在），不把它当作本区 canonical 状态。

### S0 分母与 I-READINESS

| ID | 级别 | 最晚阶段 | 台账状态 | 本审核对 |
|----|------|----------|----------|----------|
| I-READINESS-001 | required | S0 | verified | D-003 §1/§2 分母 + V 矩阵冻结成立 |
| I-READINESS-002 | required | S2 | verified | S2 self A-001 + access drill 测试路径存在 |
| I-READINESS-003 | required | S3 | verified | S3 实测 318+2；现行权威以 D-003 §5 + S3 判断为准 |
| I-READINESS-004 | required | S0 | verified | D-003 §13 框架共性能力列表 |
| I-READINESS-005 | required | S0（provider）；S5（independent 证据） | provider verified；**S5 independent 证据 = 本 A-002** | D-002 provider 与本会话输出匹配 |
| I-READINESS-006 | required | S0 | verified | D-003 §9 量尺；S1 台账按量尺分类 |
| I-READINESS-007 | required | S0 | verified | 候选/锁/变更分类字段冻结 |
| I-READINESS-008 | required | S0 | verified | §8 四宿主下限冻结 |
| I-READINESS-009 | required | S0 | verified | `go` 字段/失效触发/freshness 冻结；**裁决本身尚未发生** |

无「已标 verified 却到期未证」的 I 项伪装；S5 用户 `go` 与跨区 residual 仍为开放门禁（见 Findings）。

### 退出判据 E-1～E-5（对照 S5-evidence-matrix）

| exit_id | self/矩阵主张 | independent 结论 | 要点 |
|---------|---------------|------------------|------|
| E-1 治理与事实基线一致 | pass | **conditional** | 本区 D-003 / S1–S4 台账一致；**workspace-005 I-PROTO-FULL-001 陈旧声明仍未勘误/未 residual 接受**（见 F-001） |
| E-2 主线健康可复跑 | pass | **pass（部分独立重跑）** | 本会话：V-001/002/003 ✅；V-004 ✅ 42 files / 732 tests；V-005 ✅ build。V-006/007/008 **本会话未独立重跑**，采信 S5 `02-execution` E-001 台账 + commit 说明，并记依赖风险（F-004 recommended） |
| E-3 标准模块接入路径 | pass | pass | `s2_access_drill_test.go`、`s2-access-drill.render.test.tsx` 存在；S2 A-001 pass；composition 包测试本会话 V-002 绿 |
| E-4 UI 协议边界 | pass | pass | S3：0 protocol-gap；318+2 disposition 与 `upstream-fixtures.test.ts` 排除表一致（2 exclude 可核对） |
| E-5 阻断缺陷合法闭环 | pass | **conditional** | F-002 模态实现+3 断言本会话复跑绿；**抽屉焦点断言缺失**（F-002 close-path 不完整，见 F-002）；F-007 合法 deferred；open required Goal finding 投影在「全量 fixed」主张下有证据缺口 |
| E-6 结论可审计可复用 | 待 S5 完成 | **in progress** | 证据矩阵 + self A-001 + 本 independent A-002 已具备；用户 `go`/`no-go` + 最小字段尚未落盘 |

### Finding 台账抽查（S1 → S4）

| finding | 严重度 | 主张状态 | independent |
|---------|--------|----------|-------------|
| F-002 a11y 模态/抽屉 | required | fixed | **部分核实**：模态 ✅ 实现+测试；抽屉 ✅ 实现于 `App.tsx`，❌ 无对应焦点断言（见本审 F-002） |
| F-001 I-PROTO-FULL 文档矛盾 | major | 本区 closed；跨区勘误待办 | 确认 workspace-005 附件仍含「37/37 / 320/320 全绿」类陈旧表述；本区 S3 调和路径正确 |
| F-003/004/005/006/008/009 | minor | fixed | 抽查：`compose.yaml` `APP_PROFILE` 默认 `mvp`；upload/error 等不逐行全量重审，无矛盾信号 |
| F-007 上传授权 | minor | deferred | `upload.go` 仍仅 `a.Middleware`（认证、无权限键）——与 deferred 一致；owner/触发已记 |
| F-010/011 | info | 观察 | 不阻断 |

### 跨模块 UI a11y（D-003 §8）

| 宿主 | 冻结证据形式 | 本审 |
|------|--------------|------|
| Shell/移动导航 | 静态 aria + 必要人工 | `shell.test.ts` 含 hamburger/close aria；**焦点管理实现在 App.tsx 但无焦点类测试** |
| schema-driven 表单/列表等 | 渲染器语义断言 | 既有 render/schema-crud/table 等套件（V-004 全绿） |
| 模态与动态反馈 | 焦点/状态断言 | `modal.tsx` + `modal.test.tsx` 3 用例本会话 ✅ |
| 语言切换 | 双语断言 | `ui-bilingual` / `error-localization` 等在 V-004 中通过 |

### 本会话独立命令结果（证据新鲜度）

| # | 命令 | 结果 |
|---|------|------|
| V-001 | `cd apps/api && go build ./...` | ✅ |
| V-002 | `cd apps/api && go test ./...` | ✅ 全包 |
| V-003 | `cd apps/api && go vet ./...` | ✅ |
| V-004 | `cd apps/web && npm test` | ✅ 42 / 732 |
| V-005 | `cd apps/web && npm run build` | ✅（chunk size 警告，非失败） |
| V-006~V-008 | e2e / smoke / disposable | **未在本独立会话重跑** |

---

## 对照成功标准（GOAL-007 检查点）

| 检查点 | 状态（独立视角） |
|--------|------------------|
| S5-1 最终基线回归 | 文档主张全绿；本会话核实 V-001~V-005；V-006~008 依赖台账 |
| S5-2 证据矩阵 | 存在且结构对齐 exit_id；§5 待办清单陈旧（F-003） |
| S5-3 self + independent | self A-001 已有；**本 A-002 补齐 independent** |
| S5-4 用户裁决 | **未发生** |
| S5-5 Root 关门 | **不得**在本意见中推进 |

---

## Findings

### F-001 · workspace-005 `I-PROTO-FULL-001` 陈旧声明仍开放（E-1 / `go` 门禁）

- **严重度**：major  
- **级别**：**required**（blocking for unconditional `go`；不阻断「实现主线健康」结论）  
- **类别**：governance-drift / cross-boundary  
- **证据**：  
  - `docs/workspace-005-full-protocol-contract-v2-7-0/.../I-PROTO-FULL-001-coverage-v2-7-0.md` 仍含「37/37 全绿」「320/320 全绿」类主张  
  - 实现与本区权威：`upstream-fixtures.test.ts` 2 exclude + D-003 §5「318 执行 + 2 排除」+ S3-protocol-judgment §1.2  
  - S3/S5/self 均写明：S5 `go` 前须勘误 **或** 用户书面 residual  
- **影响门禁**：E-1 文档一致性；`go` 消费前 open residual 投影（D-003 §11）  
- **建议闭合**：  
  1. 经 `/vision` 或 workspace-005 owner **正式勘误**；或  
  2. 用户 P-004 **accepted-residual**（范围=该文档 exclude/全绿声明；影响=仅文档；复审触发=下次协议 pin/disposition 变更）  

### F-002 · F-002（a11y）关闭路径中「移动抽屉对应可复跑断言」未兑现

- **严重度**：major  
- **级别**：**required**（blocking for 无条件宣称「§8 全量断言已满足 / F-002 全量 fixed 无缺口」；**不**要求立刻判定实现缺失——`App.tsx` 焦点逻辑可核对）  
- **类别**：evidence-gap / finding-closure  
- **证据**：  
  - S1 close-path：`ModalHost … + 移动抽屉焦点管理 + 对应 vitest/Playwright 断言`  
  - S4 关闭证据主引 `modal.test.tsx` 3 断言；`shell.test.ts` 仅为抽屉状态/aria 结构，**无** focus enter/trap/Escape/restore  
  - 本会话：`modal.test.tsx` 3/3 pass；全仓 web 测试中无抽屉焦点断言命中  
- **影响门禁**：D-003 §8 / E-5 required 闭环的可重复核对性  
- **建议闭合**：  
  1. **fixed**：补抽屉焦点/Escape/恢复（及可选 Tab 循环）可复跑测试后重标 F-002；或  
  2. **accepted-residual**：书面接受「抽屉焦点仅有实现、无线上自动化断言」+ 人工核对记录 + 复审触发；或  
  3. 收窄历史 F-002 关闭声明为「模态 fixed；抽屉实现完成、断言 residual」并经用户确认  

### F-003 · S5-evidence-matrix §5 待办清单与正文/台账不一致

- **严重度**：minor  
- **级别**：recommended  
- **类别**：doc-hygiene  
- **证据**：矩阵 §2/§3 与 GOAL-007 `02-execution` 已记 V-001~V-008 全绿与 self A-001；§5 仍勾选未完成「V-007/V-008」「S5 self」等  
- **建议**：更新 §5 勾选状态，避免编排器/用户误读为回归未做  

### F-004 · V-006/V-007/V-008 本独立会话未重跑；S4 曾声明 e2e/smoke 未因 a11y 重跑

- **严重度**：minor  
- **级别**：recommended（non-blocking if S5 E-001 台账被用户接受；**生产/release 面偏好补跑或 CI 绿证**）  
- **类别**：evidence-freshness  
- **证据**：GOAL-006 A-001 观察项；GOAL-007 E-001 主张 S5 已重跑；本审仅独立核实 V-001~V-005  
- **建议**：`go` 前保留 CI `r6-basic-matrix` 或本地 e2e+smoke+disposable 一次绿证绑定最终候选 commit  

### F-005 · F-007 上传授权仍 deferred；S5 未做显式再裁决记录

- **严重度**：minor  
- **级别**：recommended  
- **类别**：deferred-confirm  
- **证据**：`handler/upload.go` 仅认证中间件；S4 deferred owner=VP-008 lead、触发含 S5  
- **建议**：用户 `go`/`no-go` 时书面确认「维持 deferred、不升 required、不扩 scope」，避免静默默认  

---

## 必改项汇总

| ID | 级别 | 一句话 | 阻断 |
|----|------|--------|------|
| **F-001** | **required** | workspace-005 `I-PROTO-FULL-001` 勘误 **或** 用户书面 residual | **阻断无条件 `go`** |
| **F-002** | **required** | 补齐 F-002 抽屉可复跑断言，**或** residual/收窄关闭声明 | **阻断无条件宣称 a11y/E-5 全量闭合** |
| F-003 | recommended | 刷新证据矩阵 §5 勾选 | 否 |
| F-004 | recommended | 绑定最终候选的 e2e/smoke/CI 绿证 | 强烈建议 |
| F-005 | recommended | `go` 时书面确认 F-007 deferred | 否 |

**无新增 blocker 级实现缺陷**（构建/测试主线、迁移连续性、协议 0 protocol-gap、模块接入路径在已核对证据下成立）。

---

## 与既有意见的异同（vs A-001 self）

| 点 | A-001 self | A-002 independent |
|----|------------|-------------------|
| S0–S4 主线一致性 | pass | **同意**（在已重跑/抽查范围内） |
| 无开放 required Goal finding | 主张 0 | **不同意无条件**：F-002 关闭证据不完整 → 本审 F-002 required；跨区 F-001 residual 仍 open |
| independent 门禁 | 待 grok | **本意见关闭 I-READINESS-005 的 independent 证据槽**（provider 输出可核对） |
| workspace-005 | 待处置，阻断 `go` | **同意** → F-001 required |
| 总 verdict | pass（self 侧） | **conditional**（不可无条件 `go`） |

self 将前置条件放在「不阻断 self verdict、阻断 go」；independent 将其中 **未处置跨区 residual** 与 **F-002 证据缺口** 升格为正式 required findings，以便 `/govern` 必须响应后才能放行。

---

## 分面结论摘要

| 分面 | 结论 | 说明 |
|------|------|------|
| compatibility | pass（conditional 包装于总 verdict） | 协议 pin/318+2、Profile/模块矩阵、接入演练、V-001~005 独立核实 |
| data | pass | 单租户 SQLite + bootstrap/reconcile 分母一致；无多租户承诺 |
| migration | pass | 0001–0010 收集/连续版本规则与 ownership 记录一致；本会话 migration/store 测试绿 |
| production/release | conditional | fail-closed/compose/默认 Profile 抽查 OK；smoke/e2e 依赖台账；`go` 字段待裁决落盘 |
| 跨边界治理 | conditional | 本区绑定正确；共享资料 none 未误用；跨区文档 residual 未闭合（F-001） |
| 跨模块 UI a11y | conditional | 模态下限满足可复跑；抽屉实现有、断言缺口（F-002） |

---

## 结论 + 建议给编排器/用户的下一步

**Verdict：`conditional`**

S0–S4 准入准备的主体证据充分，独立抽检未发现「主线构建/测试/协议/迁移」名不副实。但 **不可无条件签发 `go`**：

1. 必须合法闭合 **F-001**（跨区勘误或 residual）；  
2. 必须合法闭合 **F-002**（补断言或 residual/收窄关闭）；  
3. 用户书面 `go`/`no-go` + D-003 §11 最小字段；  
4. 建议处理 F-003~F-005。  

**不得**将本 independent 意见解释为已 `go` 或可关闭 Root。

### 建议 `/govern` 输入

```text
/govern 响应 GOAL-007 A-002（independent, verdict=conditional）：
闭合 F-001（workspace-005 I-PROTO-FULL-001 勘误或用户 residual）与
F-002（抽屉 a11y 断言补齐或 residual/收窄 F-002 关闭声明）；
确认 F-007 维持 deferred；刷新 S5-evidence-matrix §5；
然后提交用户 go/no-go 裁决（绑定最终候选 commit + D-003 §11 字段）。
```

---

## 声明

- `source: independent`；auditor = grok-build / grok-4.5 · `/audit`。  
- 本意见**不修改**任何目标 `status` / `progress` / goal-tree / 方案正文。  
- 响应、finding 闭合与阶段推进由用户通过 **`/govern`** 处理。  
- 进度百分比不得作为 finding 闭合或放行依据。
