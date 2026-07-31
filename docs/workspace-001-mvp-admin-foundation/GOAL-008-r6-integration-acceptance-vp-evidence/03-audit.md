---
id: GOAL-008-r6-integration-acceptance-vp-evidence
doc: audit
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.5.0
---

# 审计 · GOAL-008

> 本文件是本目标的唯一正式意见台账（P-003）。正式意见须从 `A-001` 起编号，并包含 `source`、`scope` 与 `verdict`。本轮已追加阶段 1 计划 self 审视；尚无 independent 审计，且 self 审视不替代阶段 2 执行证据。

## 信息就绪核对（R6 规划基线）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| 影响本 scope 的 I-00N | 已登记 | `I-008-001`～`I-008-005` 均为 required |
| 到期 required 是否已 verified / residual | 本轮全部闭合 | 五项均「有证据/已决定」；CI 首跑 green 闭合 I-008-002/005；无需 residual |
| 资料引用是否固定且用户确认 | 无 | `shared_materials_catalog: none`；本目标不使用共享资料目录 |
| Vision required finding | 0 open | VRev-003 `pass`；`F-V003` recommended 不阻断本 R6 规划 |
| 相关 Goal required finding | 1 open → 本轮响应后 0 | A-001 `F-008-001` 经 A-002 计划审视通过后由响应节关闭 |

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required |
|------|------|--------|-------|---------|---------------|
| A-001 | 2026-08-01 | self | 阶段 1 验收合同与证据计划 | fail | 1 |
| A-002 | 2026-08-01 | self | 阶段 1 冻结候选计划审视（矩阵/环境/oracle/schema/CI） | pass | 0 |
| A-003 | 2026-08-01 | self | 阶段 2 → 3 门禁（集成验收执行证据） | pass | 0 |
| A-004 | 2026-08-01 | self | 阶段 3 → 4 门禁（VP 证据汇编） | pass | 0 |
| A-005 | 2026-08-01 | self | R6 关门审计（close-out） | pass | 0 |

## A-001 · 阶段 1 计划与信息门禁自审（2026-08-01）

- **source**：self
- **auditor**：Codex `/govern`
- **类型**：stage / design-plan
- **scope**：GOAL-008 阶段 1「验收合同与证据计划冻结」；I-008-001～I-008-005；draft evidence schema 与本地能力基线
- **verdict**：fail

### 范围与区间

本审视只核对阶段 1 计划是否已经具备进入阶段 2 的信息条件，不把本轮本地测试/构建结果当作 R6 验收、VP 关门或跨平台证据。

### 成果（有证据）

- 现有命令、cwd、运行态入口与环境身份已记录于 [02-execution.md](02-execution.md) 和 [R6-acceptance-plan.md](attachments/R6-acceptance-plan.md)。
- R4 单端/fixture/HTTP 事实与 R5 逐域登记输入已被识别，但仍按“输入”而非已冻结 R6 oracle/coverage evidence 处理。
- `attachments/evidence-index.schema.json` 与 `attachments/evidence-index.dry-run.json` 可作为 `I-008-004` 的候选形状；dry-run 明确为 `planning`、`blocked`、`not-captured`。

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| VP 三条判据有明确主张、证据、排除与 residual 形状 | 部分 | 计划 v0.1.1 已映射三条判据；最低矩阵、跨层 oracle 和正式 evidence contract 未冻结 |
| React + Go 本地命令可复跑 | 部分 | 本轮 revision/runtime 与 15/395 Web、Go test/build 结果；未形成 clean-install / 双服务 R6 artifact |
| R2/R5 覆盖可追溯 | 部分 | R5 `I-007-001` 和 R2 `I-PROTO-001 v0.1.3` 已识别；R6 coverage map/conformance result 尚不存在 |
| 账号权限前后端集成 oracle | 未完成 | R4 事实可复用；API→Web/Renderer/动作链正向/拒绝场景尚未冻结 |
| 机器可读 evidence index、hash 与重跑规则 | 部分 | draft schema/dry-run 已新增；结果 artifact、正式 schema 验证和文件摘要仍缺 |

### Findings

#### F-008-001 · 阶段 1 required 信息尚未闭合

- 严重度：med
- 建议：required
- 描述：`I-008-001`～`I-008-005` 均在阶段 1 最晚需要；当前分别为 `collecting`、`collecting`、`collecting`、`collecting`、`collecting`，没有用户书面 `accepted-residual`，且 D-002 仍为 `proposed`、计划附件仍为 draft。因此阶段 1 的“计划冻结”退出条件未满足，不能进入阶段 2。
- 证据：[00-meta.md](00-meta.md) 信息表；[01-decision.md](01-decision.md) D-002/D-003；[R6-acceptance-plan.md](attachments/R6-acceptance-plan.md) 阶段门禁；[02-execution.md](02-execution.md) 本轮基线。
- 状态：open

### 必改项汇总（required 列表）

- `F-008-001`：完成/审视五项 required 信息，或获得明确范围、期限、缓解与复审触发的用户书面 `accepted-residual`；在此之前保持阶段 2 关闭。

### 结论 + 建议下一步

本轮确实推进了阶段 1 的信息收集，但没有达到冻结条件。建议继续完成 `I-008-001`～`I-008-005` 的证据/决策闭环：先把跨层账号权限 oracle、最低环境矩阵和正式 evidence schema 形成候选，再由用户裁决任何 residual/有界实验，随后重新做同 scope 计划审视。当前没有触发 P-004.1（尚无 independent 审计）；若请求接受 residual 或以有限实验放行，必须触发 P-004.4 并留痕。

## 下一审视点

- A-001 已确认本地能力基线与 draft evidence schema 可作为规划输入，但 `I-008-001`～`I-008-005` 仍未闭合；不得进入阶段 2。
- 继续收集五项 required，并在 D-002 从 proposed 进入冻结候选后重新执行同 scope 计划审视。
- 若先出现 independent 审计而无同 scope self audit，进入后续门禁前按 P-004.1 询问用户是否需要自审。
- 任何 required finding 未按 `fixed` / `accepted-residual` / `user-overruled` 合法闭合前，不得进入 R6 阶段 2 或关门。

## A-002 · 阶段 1 冻结候选计划审视（2026-08-01）

- **source**：self
- **auditor**：Claude Code `/govern`
- **类型**：stage / design-plan（同 scope 计划审视）
- **scope**：GOAL-008 阶段 1「验收合同与证据计划冻结」；`I-008-001`～`I-008-005`；验收矩阵 C-001～C-008；环境矩阵（D-004）；账号权限 oracle；evidence schema dry-run；CI 首跑
- **verdict**：pass

### 范围与区间

本审视核对阶段 1 是否具备冻结条件（五项 required 有证据结论或合规 residual；矩阵/环境/oracle/schema 可审视；无开放 required finding）。本审视**不**放行阶段 2 执行，也不把本地/CI 通过写成 R6 验收完成。

### 成果（有证据）

- **I-008-001**：验收矩阵 C-001～C-008 已落盘 [R6-acceptance-plan.md](attachments/R6-acceptance-plan.md) v0.2.0 §2b，覆盖 VP 三条退出判据，每条有主张/入口/预期/证据/排除。
- **I-008-002**：本地双服务/health/proxy/账号上下文/records 实测（`evidence/planning/results/runtime-probes.log`）+ GitHub Actions 首跑 green（`npm ci` 干净安装 + Linux 等价）。
- **I-008-003**：账号权限跨层 oracle 已登记 [account-permission-oracle.md](attachments/account-permission-oracle.md) v0.1.0（正向 P-1～P-4、拒绝 D-1～D-6、排除与边界）。
- **I-008-004**：evidence schema 经 [validate-evidence-dry-run.mjs](attachments/validate-evidence-dry-run.mjs) 校验 **可解析、5 artifact SHA-256 可重算**，dry-run 持久化为 `evidence/planning/evidence-index.dry-run.json`。
- **I-008-005**：用户裁决（D-004）搭建最小 CI + 浏览器矩阵；CI 首跑 run `30666932343` **success**（api 22s / web 27s / browser-e2e 53s）。

### 对照阶段 1 退出条件

| 退出条件 | 状态 | 证据 |
|----------|------|------|
| 五项 required 有证据结论或合规 residual | **达成** | 均「有证据/已决定」，无需 residual |
| 验收矩阵 / 环境矩阵 / oracle / 证据格式可审视 | **达成** | R6-acceptance-plan v0.2.0 §2b/§4c；account-permission-oracle v0.1.0；schema dry-run |
| 计划审视无开放 required finding | **达成** | 本轮 A-002 pass；F-008-001 由响应节关闭 |
| D-002 由 proposed 冻结为 accepted | **本轮完成** | 见 [01-decision.md](01-decision.md) D-002 更新 |

### Findings

#### F-008-002 · CI 非阻断注解（recommended，不阻断）

- 严重度：low
- 建议：recommended
- 描述：CI run `30666932343` 两条非阻断注解——`actions/checkout`/`setup-node`/`setup-go` 触发 Node 20 弃用强制跑 Node 24（GitHub 侧行为）；`setup-go` 缓存因 `apps/api/go.sum` 缺失 skip（API 无外部依赖）。均不影响 job 成功。
- 状态：open（recommended，不阻断阶段 1 冻结）

### 必改项汇总（required 列表）

无 open required。

### 结论 + 建议下一步

阶段 1 冻结候选经审视通过（pass）：五项 required 均有实际证据且 CI 首跑 green；验收矩阵、环境矩阵、账号权限 oracle 与 evidence schema 可支持阶段 2 执行。D-002 已按用户授权冻结为 accepted。**阶段 1 冻结完成，可进入阶段 2**；阶段 2 执行时须按 evidence index 持久化结果、显式记录失败/排除，并按 `I-008-004` 的 schema 产出正式 acceptance index。

---

## 响应节 · 响应 A-001 · 关闭 F-008-001（2026-08-01）

**响应**：A-001（self · fail · F-008-001 required「阶段 1 required 信息尚未闭合」）。

**关闭证据表**：

| finding / I-00N | 状态 | 证据路径 |
|-----------------|------|----------|
| `F-008-001` | **fixed** | 五项 required 均已闭合；A-002 计划审视 pass（见上）；CI 首跑 green |
| `I-008-001` | verified | [R6-acceptance-plan.md](attachments/R6-acceptance-plan.md) v0.2.0 §2b 验收矩阵 |
| `I-008-002` | verified | `evidence/planning/results/runtime-probes.log` + run `30666932343`（npm ci + Linux） |
| `I-008-003` | verified | [account-permission-oracle.md](attachments/account-permission-oracle.md) v0.1.0 |
| `I-008-004` | verified | `validate-evidence-dry-run.mjs`（SCHEMA_VALIDATION_OK）+ `evidence-index.dry-run.json` |
| `I-008-005` | verified | D-004（用户裁决）+ run `30666932343` browser-e2e pass |

**闭合路径**：`fixed`（可核对修正；用户 `/govern` 授权阶段 1 冻结并响应 A-001）。

**仍开放项**：`F-008-002`（recommended，CI 注解清理，不阻断）。

**结论**：A-001 的必改项已满足——五项 required 信息完成/审视且无合规 residual 需要；阶段 2 门禁解除，等待阶段 2 执行。

## 下一审视点

- 阶段 2 执行开始：按 evidence index 持久化 Web/API/协议回归/账号权限集成结果，失败与排除显式列出。
- 阶段 2 完成后做阶段 2→3 门禁审视；全部 required 闭合后，再由用户授权 Root R6 / `progress` / status，VP 关门另走 `/vision`。

## A-003 · 阶段 2 → 3 门禁审视（2026-08-01）

- **source**：self
- **auditor**：Claude Code `/govern`
- **类型**：stage / execution-facts
- **scope**：GOAL-008 阶段 2「集成验收执行」→ 阶段 3「VP 证据汇编与缺口整改」门禁；验收矩阵 C-001～C-008；`evidence-index.json`（mode: acceptance）
- **verdict**：pass

### 范围与区间

本审视核对阶段 2 退出条件是否满足（冻结矩阵 required item 全运行并落盘；失败/未执行/排除显式；evidence index 可解析且文件摘要可重算；新发现关键未知已回流信息表），以决定是否放行阶段 3。本审视不把 evidence index 写成 VP 关门证据，也不改 Goal/VP status。

### 成果（有证据）

- **C-001～C-008 全部执行并落盘**（revision `a941bedb1fc2cd4859a408df50653e867da35ff2`，worktree clean）：web test/build（15 files / 395 tests pass）、api test/build、双服务+health/proxy/账号上下文/records 探测、浏览器 E2E（1 passed / 0 unexpected）、D-PERM fixtures、stage3 conformance。原始输出见 `attachments/evidence/acceptance/results/`。
- **正式 evidence index**：`attachments/evidence/acceptance/evidence-index.json`（`mode: acceptance`，**7 artifact SHA-256 verified，overallOutcome=pass**），经 [build-acceptance-index.mjs](attachments/build-acceptance-index.mjs) 用 ajv 2020 校验通过、可重跑。
- **排除显式记录**：reactions multi-round 16/16（Root D-008）、request-construction batch 11（Root D-010 Q1=否）、D-UPLOAD 整域（I-PROTO-001 v0.1.3）、本地非干净安装（`npm ci` 由 CI run `30667596846` 覆盖）、浏览器级拒绝未断言（真实 manifest 无权限门控项，拒绝以 renderer/组件层 oracle 断言）——均列入 evidence-index exclusions，未用总体 pass 掩盖。

### 对照阶段 2 退出条件

| 退出条件 | 状态 | 证据 |
|----------|------|------|
| 冻结矩阵全部 required execution item 已运行并落盘 | **达成** | C-001～C-008 均有结果 artifact（evidence-index 7 项） |
| 失败、未执行、排除与平台缺口显式 | **达成** | evidence-index exclusions 5 项 + 各结果 outcome |
| evidence index 可解析、文件摘要可重算 | **达成** | ajv 校验 OK；7 SHA-256 可重算（build-acceptance-index.mjs） |
| 新发现关键未知已回流信息表，未静默扩域 | **达成** | 未发现新关键未知；边界未扩大 |

### Findings

无 open required / recommended finding 影响本门禁。`F-008-002`（CI 注解，recommended）不阻断。

### 必改项汇总（required 列表）

无 open required。

### 结论 + 建议下一步

阶段 2 退出条件全部满足（pass）：集成验收已执行、机器可读证据已按 schema 落盘、失败/排除显式。**可进入阶段 3「VP 证据汇编与缺口整改」**：把 VP 三条退出判据逐条指向工作区 Q2 证据，required 缺口按 P-003 闭合，边界主张一致。

## A-004 · 阶段 3 → 4 门禁审视（2026-08-01）

- **source**：self
- **auditor**：Claude Code `/govern`
- **类型**：stage / execution-facts
- **scope**：GOAL-008 阶段 3「VP 证据汇编与缺口整改」→ 阶段 4 门禁；VP-001 三条退出判据的 Q2 证据映射
- **verdict**：pass

### 范围与区间

本审视核对阶段 3 退出条件：VP 三条判据逐条指向工作区 Q2 证据；无未闭合 required 缺口；边界主张一致。本审视不把汇编写成 VP 已可关门，也不改 Goal/VP status。

### 成果（有证据）

- **汇编产物**：[vp-evidence-assembly.md](attachments/vp-evidence-assembly.md) v0.1.0 把 VP-001 三条退出判据逐条映射到工作区 Q2 证据：
  - 判据 1 → evidence-index C-001～C-004（web/api test+build、双服务+proxy、浏览器 E2E）+ CI run `30666932343` + protocol pin；
  - 判据 2 → I-PROTO-001 v0.1.3、I-007-001 registry、stage3 conformance、docs/schemas；
  - 判据 3 → runtime-probes `/api/accounts/me`、E2E session 断言、account-permission-oracle、dperm 17 例。
- **缺口**：无 required 缺口；浏览器级拒绝未断言（C-006）与 reactions/batch/D-UPLOAD 均为登记排除，边界与 I-PROTO-001 v0.1.3、R5 一致。

### 对照阶段 3 退出条件

| 退出条件 | 状态 | 证据 |
|----------|------|------|
| VP 三条判据逐条指向工作区 Q2 证据 | **达成** | vp-evidence-assembly.md 三条映射表 |
| 所有 required 缺口 fixed 或用户书面 residual/overruled | **达成** | 无 required 缺口；排除项已登记 |
| 边界主张一致 | **达成** | 与 I-PROTO-001 v0.1.3、R5 登记一致 |

### Findings

无 open required / recommended finding 影响本门禁。

### 必改项汇总（required 列表）

无 open required。

### 结论 + 建议下一步

阶段 3 退出条件满足（pass）：VP 三条退出判据均有 Q2 工作区证据，无 required 缺口，边界一致。**可进入阶段 4「R6 关门审计与 VP 提案输入」**：跑 close-out 自审，开放 required=0 后由用户授权 Root R6/`progress`/status 变化；VP 关门提案另走 `/vision`。

## A-005 · R6 关门审计（close-out · 2026-08-01）

- **source**：self
- **auditor**：Claude Code `/govern`
- **类型**：close-out
- **scope**：GOAL-008「R6 · 集成验收与 VP 证据」全目标关门；VP-001 三条退出判据的 Q2 证据；`I-008-001`～`I-008-005`
- **verdict**：pass

### 范围与区间

本审计核对 R6 关门条件：相关意见无未合法闭合的 required；相关信息项无未处理关门 required；成功标准逐条可核对；并有阶段/关门向审计。本审计**不**自行改 Root `progress`/status 或 VP status（须用户授权 / `/vision`）。

### 成果（有证据）

- **意见台账**：A-001（F-008-001 已 fixed）、A-002/A-003/A-004 均 pass；**开放 required = 0**。无 independent 意见、无冲突。
- **信息台账**：`I-008-001`～`I-008-005` 均 verified（矩阵/运行环境/账号权限 oracle/evidence schema/环境矩阵）。
- **验收证据**：`evidence-index.json`（mode: acceptance，revision `a941bed`，worktree clean）7 artifact SHA-256 verified、overallOutcome=pass；5 项排除显式、residuals=0。
- **VP 证据汇编**：[vp-evidence-assembly.md](attachments/vp-evidence-assembly.md) 三条退出判据逐条映射 Q2 证据。
- **CI**：`.github/workflows/r6-basic-matrix.yml` 首跑 green（run `30666932343`）；后续 dev 推送均 green。

### 对照成功标准

| 成功标准 | 状态 | 证据 |
|----------|------|------|
| 受控验收矩阵映射 VP 三条判据（主张/入口/预期/证据/排除/residual） | **达成** | R6-acceptance-plan v0.2.0 §2b C-001～C-008 |
| React+Go 干净环境可复现启动 + 浏览器/API 关键路径，证据绑定 revision/runtime/worktree | **达成** | evidence-index（revision `a941bed`/clean/env）+ C-003/C-004 |
| R2 v0.1.3 每个纳入域可从 R5 登记追到实现/范例/验证，R6 回归不扩大边界 | **达成** | I-007-001 registry v0.8.1 + stage3 conformance + C-007 |
| 账号权限链路含可核对正向/拒绝路径，不依赖未声明业务模块 | **达成** | account-permission-oracle（P-1～P-4 / D-1～D-6）+ runtime-probes + C-005 |
| 机器可读索引完整记录命令/退出码/时间/环境/结果/排除/文件摘要，稳定寻址 | **达成** | evidence-index.json（schema 校验）+ build-acceptance-index.mjs |
| R6 关门审计无开放 required；Root/VP 变化仍分别等用户与 `/vision` | **达成** | 本 A-005；用户授权项在响应节说明 |

### Findings

无 open required / recommended 影响本关门。`F-008-002`（CI 注解，recommended）不阻断。

### 必改项汇总（required 列表）

无 open required。

### 结论 + 建议下一步

R6 目标满足关门条件（pass）：四条验收证据链 + 三条 VP 判据映射均可核对，意见/信息台账均闭合，开放 required=0。**GOAL-008 可关门（status → done）**。后续受控动作：
1. **用户授权** Root R6 检查点完成、Root `progress` 5/6 → 6/6、Root status 变化（用户 `/govern` 决定）；
2. VP-001 关门提案另走 `/vision`（读取 R6 工作区证据、用户确认）。
