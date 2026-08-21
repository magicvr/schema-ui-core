---
title: 决策 · 语义化 Admin 资源替换与双实体验证
status: active
created: 2026-08-03
updated: 2026-08-04
parent: GOAL-010-a002-schema-adapter
version: 0.4.0
---

# 决策 · GOAL-011

## D-001 · 采用 users + roles，替换 records 默认语义并独立立项

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**：同意创建新目标承接 GOAL-010 S4；将 `records` 替换为对绝大多数系统具有实际语义的设计，并把 S4 的“新增”改为第二种同样具有实际语义的资源。用户确认采用编排器推荐的 `users + roles` 结构。
- **决定**：
  1. `users` 替换 `records` 作为默认代表实体；`roles` 作为新增的第二个语义资源。
  2. 创建 `GOAL-011-s4-semantic-admin-resources`，`parent: GOAL-010-a002-schema-adapter`，以五个顺序检查点承载契约、后端双资源、records 退场、Schema 集成与审计交接。
  3. GOAL-010 S4 改为本目标的父级验收门：本目标完成后，才判断 records 是否已从当前产品默认运行面退场、users/roles 是否完成双资源闭环，以及前端 Renderer 主路径是否保持不变。
  4. “Schema-only 接入”只约束前端页面接入：后端资源注册、持久化、权限与领域规则仍须显式实现和验证。
  5. users/roles 只交付最小可用 Admin 管理边界；密码哈希、refresh token 等敏感字段不得进入通用响应，完整 IAM、SSO、SCIM 与复杂策略编排保持非目标。
  6. GOAL-010 D-002/I-010-001 的 records 零 API 变更继续作为 S1～S3 的历史实施与回归事实；本目标的终态替换必须通过新契约和迁移版本演进，不静默改写既有决策、迁移或审计记录。
- **理由**：
  - `records` 已证明 transport/CRUD 功能，但其 `name/status/owner` 模型没有稳定业务语义；fork 项目需要主动移除 API、表、种子、权限、菜单、操作日志、fixture 与测试，构成产品基线污染。
  - users 与 roles 已属于本仓真实认证/RBAC 持久化域，绝大多数 Admin 系统均有实际价值，并直接服务 VP-002 的真实认证、最小权限与 fork 即用边界。
  - 替换涉及多个独立门禁域和至少两个可独立验收交付块，已超过 GOAL-010 单个 S4 检查点可以诚实承载的范围，按 P-001 建立子目标更可追踪。
- **未选方案**：
  - **仅在 GOAL-010 S4 内直接实施**：会把领域契约、数据迁移、双资源交付和退场清理压进一个检查点，门禁与证据混杂。
  - **继续使用 records，只新增 catalog**：仍要求 fork 项目清理两套无普遍语义的示例域，不能解决用户指出的代码污染。
  - **roles + menu_items**：实现风险较低，但缺少账户管理闭环，通用价值与 VP-002 真实认证链的结合弱于 users + roles。
- **信息门禁**：`I-011-001`、`I-011-002`、`I-011-003` 均为 required，初始 `open`；当前只放行 S1 信息收集，不放行 S2～S4 产品实施或验收。
- **影响**：本目标 `active / 0/5`；GOAL-010 保持 `active / 3/5`、S4 未勾选；Root A-002 F-002-001 继续 `open`，Root/VP-002 关门继续阻断。
- **后续**：先收集并提交 `I-011-001`/`I-011-002` 的方案裁决，冻结 S1 后再进入任何 users/roles 或 records 退场代码变更。

## D-002 · 冻结 S1：users/roles 领域契约与 records 退场策略

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**（P-004.4）：用户于契约冻结裁决点逐项确认三项关键取舍（均采纳编排器推荐）：
  1. **通用工厂 + 最小契约扩展**：users/roles 均走通用资源工厂五路由；扩展 `Resource.JSONFields`（任意 JSON 值字段透传，承载 users.roles）与 `DomainError{Status,Code,Message}`（实体返回领域错误，工厂逐字映射）；对 I-010-001 §5「不引入 409」做限定范围偏离（仅账号域 409，envelope 形状不变）。
  2. **操作日志纳入**：migration `0005` 重建 `operation_log` event CHECK，新增 `users.*`/`roles.*` 事件（保留 `records.*`/`auth.*` 历史合法值）；users/roles 写路径挂 `OnWrite`。
  3. **records 硬退场 DROP TABLE**：migration `0006` `DROP TABLE records` + 清理 records 权限/菜单行；既有库升级自动 `pre-v0006` 快照兜底，records 数据随表删除（可由快照恢复）。
- **决定**：采纳两份版本化契约冻结 S1——
  1. [I-011-001-users-roles-contract.md](attachments/I-011-001-users-roles-contract.md) **v0.1.0**：users/roles 资源契约（公开字段、敏感字段隔离、角色分配、self/最后管理员保护、system role 保护、grant 约束、权限键/菜单/操作日志、错误码）、最小 IAM 边界、通用工厂最小扩展。
  2. [I-011-002-records-retirement.md](attachments/I-011-002-records-retirement.md) **v0.1.0**：records 足迹盘点、fresh install 与 in-place upgrade 迁移矩阵（0005/0006）、硬退场数据处置、代码/种子/fixture/前端退场动作、S3 验收口径。
- **理由**：三裁决点均为「先例契约改写/数据处置/范围取舍」，须用户书面确认而非编排器静默推断；采纳推荐路径保持 S2「通用工厂之上」、审计链一致与 fork 基线干净。
- **未选方案**：users 自定义 handler（不改工厂契约但 S2 主张打折、双套门禁逻辑）；操作日志不纳入（省 0005 但账号变更无审计留痕）；软退场保留死表（fork 基线不干净、数据处置不明）。
- **信息门禁**：`I-011-001` → **verified**（契约 v0.1.0 + 本决策）；`I-011-002` → **verified**（契约 v0.1.0 + 本决策）。S1 方案冻结门禁解除，S2 实施与 S3 退场可放行。`I-011-003` 保持 `open`（最晚 S4 前）。
- **影响**：**S1 检查点达成，GOAL-011 `0/5 → 1/5`**；GOAL-010 保持 `active / 3/5`；Root A-002 F-002-001 仍 `open`，Root/VP-002 关门继续阻断。未修改任何产品代码（S1 为文档冻结）。
- **后续**：S2 后端 users/roles 资源闭环（通用工厂扩展 + store 领域方法 + 双资源 CRUD + 401/403）→ S3 records 退场 → S4 双资源 Schema 接入 → S5 回归审计关门。

## D-003 · 响应 A-002：契约修订至 v0.2.0（fixed 闭合 F-001/F-002 + 采纳 F-003~F-006）

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**（P-004 §3.2）：A-001（self · pass）与 A-002（independent · conditional）同 scope verdict 冲突；用户裁决闭合路径「**全部 fixed**」——修订两契约 + 采纳 recommended，不补 residual/overruled。
- **决定**：
  1. **I-011-001 v0.1.0 → v0.2.0**：§7 补 `ResourceEntity` Create/Update/Delete 的 `account.User` actor 通道（F-001，SELF_OPERATION/LAST_ADMIN 可诚实实现）+ DomainError 检查先于 ErrNotFound/INTERNAL（F-006）；§2.3 禁 API 路径复用 `linkUserRole`→`ensureRole` 隐式建角色（F-004）；§3.0 冻结 roles 公开响应形状 `system:boolean` + 毫秒时间戳（F-003）。
  2. **I-011-002 v0.1.0 → v0.2.0**：§2.1/§2.3 快照语义改为「每个待应用数据变更迁移前快照」（至少 0006 前强制），0005+0006 同批时 `pre-v0006` 必然存在（F-002）；§5 验收句对齐。
  3. **GOAL-010 D-005 / I-010-001 v0.2.2**：父契约 §5 追加账号域 409 限定扩展注记（F-005，消除跨目标双真相）。
  4. A-001 F-001（password 长度）与 F-002（fixture 文案）维持 recommended，随 S2/S3 落实（F-006 承接）。
- **fixed 关闭证据**：两契约 v0.2.0 + 本决策 + GOAL-010 D-005/I-010-001 v0.2.2；A-001/A-002 差异经此趋同（见 03-audit 响应节）。
- **理由**：F-001/F-002 均为真实可核对缺口（工厂无 actor 通道则 self 保护不可诚实实现；快照机制与验收字面不一致），修文档成本低、无 residual 必要；recommended 采纳后 S2/S4 金标准更明确。
- **未选方案**：residual（缺口小、应修）；overruled（拒绝合理必改项无依据）；仅闭 required 延后 recommended（用户裁决全部采纳，提高 S2 可实施性）。
- **信息门禁**：`I-011-001`/`I-011-002` 维持 `verified`（v0.2.0 为响应修订，不改变冻结结论）；`I-011-003` 保持 open。
- **影响**：GOAL-011 保持 `active / 1/5`，S1 契约以 v0.2.0 为准；A-002 conditional 经 F-001/F-002 闭合后与 A-001 pass 趋同；S2 实施门禁保持解除；Root A-002 F-002-001 仍 open。
- **后续**：进入 S2 后端 users/roles 资源闭环（按 v0.2.0 契约落地工厂扩展 + store 领域方法 + 双资源 CRUD + 401/403）。

## D-004 · 响应 A-007：fixed 闭合三项必改并冻结 I-011-003 v0.2.0

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**（P-004.1 / P-004 §3.2）：A-007 为同范围唯一正式意见，且存在三条 required finding。用户明确裁决“**不用补自审计，直接 fix**”：跳过同范围自审，不采用 residual 或 overruled，F-001～F-003 全部走 `fixed`。
- **决定**：
  1. 将 [I-011-003-acceptance-matrix.md](attachments/I-011-003-acceptance-matrix.md) 从 A-007 审计时的候选 v0.1.0 修订并冻结为 **v0.2.0**；只有本决策落盘后，附件才使用“冻结”与 `verified` 的既成事实表述（F-001）。
  2. 冻结 Renderer/App 基线为提交 `adfe15a17da770699d5e109f22402c41ece5eeea`，并在契约 §3 固定受限生产文件清单和可执行 `git diff --exit-code` 命令；测试文件不属于受限生产路径，但不能替代零 diff 证明。
  3. Web 补齐 roles 真实 fixture 的 create/update/delete action 断言、users/roles action id 无 Renderer 硬编码断言，以及真实 manifest + roles fixture 的页面级渲染断言（F-002）。
  4. API 补齐 users/roles 五路由匿名 401 与 viewer 读 200/写 403、双资源 operation-log actor/record/detail、roles 进程重启 list/detail/毫秒时间戳往返，以及升级库重开后 users/roles 事件持久化断言（F-003）。
  5. `I-011-003` → **verified**，解除 S4 集成验收的信息门禁；S4 仍须按冻结矩阵形成独立验收收据并通过阶段审视，不能由信息项关闭自动推导完成。
- **fixed 关闭证据**：I-011-003 v0.2.0；`schema-crud.test.tsx` T-UI-10；`representative-pages.integration.test.tsx` roles 页面断言；`users_test.go` / `roles_test.go`；`server_restart_test.go`；`operations_test.go`；本机 `go test ./...`、Web 485/485、Web build 与 Renderer 基线 diff 均通过（详见 `02-execution.md` 和 A-007 响应节）。
- **理由**：三项 finding 均指向可核对的契约真实性或证据缺口，修复成本可控，且直接影响 `I-011-003` 的核心主张；没有接受残余或驳回意见的合理依据。
- **未选方案**：补同范围自审（用户明确跳过）；收窄为单资源/代表性路径（不足以支持 GOAL-011 双实体成功标准）；accepted-residual / user-overruled（没有必要保留已可直接修复的缺口）。
- **影响**：A-007 F-001～F-003 均 `fixed`，`I-011-003` 信息门禁解除；GOAL-011 仍为 `active / 3/5`，S4/S5 均未勾选；GOAL-010、Root A-002 F-002-001 与 VP-002 状态均不变。本次未触发 status/progress/parent 变更，`goal-tree.md` 无需状态投影更新。
- **后续**：按 I-011-003 v0.2.0 执行 S4 完整验收并形成收据；阶段审视通过且无开放 required 后，才可勾选 S4 并同步 goal-tree。

## D-005 · 响应 A-012：恢复关门门禁并保留五项必改为 open

- **日期**：2026-08-04
- **状态**：accepted
- **用户指令边界**：用户明确调用 `/govern` 响应 A-012；该指令授权本轮完成意见汇总、fail-closed 状态恢复和待裁决信息登记，但没有选择 `fixed` / `accepted-residual` / `user-overruled`，也没有授权把尚未实施的修复写成既成事实。
- **冲突**：A-010（self · pass）与经修复后趋同的 A-011（independent · conditional）支持原 close-out；A-012（independent · fail）在同一关门范围新增五条 required。按 P-004，该冲突在用户选择闭合路径前保持未决，不取乐观侧。
- **决定**：
  1. 采纳 A-012 为正式相关意见，并将 F-001～F-005 全部保留为 `required / open`；F-006 保留 `recommended / open`。本轮不声明任何 finding 已闭合。
  2. 因开放 required 与 `status: done` 不相容，将 GOAL-011 从 `done` 恢复为 `active`；原 S1～S5 五个检查点仍为已发生事实，故派生 `progress` 保持 `5/5`。状态恢复不是进度回退，也不是 finding closure。
  3. 新增 `I-011-004`（required）：裁决 F-003 的产品边界——保留 seed-only 并移除/隐藏无有效授权语义的管理面，或补齐角色分配、grant 来源与 roles 管理流程。该信息项阻断 F-003 整改方案冻结与重新关门。
  4. F-001、F-002、F-004、F-005 的推荐主路径为 `fixed`；静态复核已确认其实现/流水线落点，但实际修复、测试和独立复审尚未发生，不在本决策中冒充完成。
  5. A-010/A-011 原文和原 S1～S5 执行收据保留为历史；A-012 响应只追加新状态，不重写旧 verdict 或关闭记录。
- **理由**：安全授权、凭据生命周期、活动 CI 和最后管理员原子性均影响 VP-002 的生产边界，不能以既有绿测或 `5/5` 覆盖；F-003 同时触及既有“完整 IAM 非目标”边界，必须先由用户明确选定产品口径。
- **未选方案**：
  - **维持 `done / 5/5` 等以后再处理**：违反开放 required 关门门禁。
  - **直接将五项写为 fixed**：当前没有修复代码、回归收据或 finding-closure 复审。
  - **静默把 F-003 视为完整 IAM 范围内/范围外**：两种路径都会代替用户做 P-004 裁决。
  - **立即新建 GOAL-012**：finding 仍权威落在 GOAL-011；在用户确认整改结构前，新目标不能替代本目标的响应与关门阻断。
- **影响**：GOAL-011 → `active / 5/5`；GOAL-010 保持 `active / 3/5`；Root 与 VP-002 均保持 active，既有 Root finding 不在本轮关闭。`goal-tree.md` 同步状态与本响应注记。
- **后续**：用户确认 A-012 的闭合路径，并对 `I-011-004` 二选一；若选择推荐路径，则先冻结整改契约/测试矩阵，再实施 F-001～F-005，随后请求限定 scope 的 `/audit` finding-closure 复审。F-006 单独保持 recommended，不得因可选卫生项阻断必改整改。

## D-006 · 裁决 A-012 全部 fixed，并选择完整角色授权/grant 管理路径

- **日期**：2026-08-04
- **状态**：accepted
- **用户书面裁决**（P-004.2 / P-004.3 / P-005）：A-012 F-001、F-002、F-004、F-005 按编排器建议走 `fixed`；F-003 选择“**补齐角色授权/grant 管理路径**”，同样走 `fixed`。未选择 `accepted-residual` 或 `user-overruled`。
- **决定**：
  1. F-001～F-005 的合法闭合路径统一确定为 `fixed`；路径选择已完成，但在实现、回归与限定范围独立复审完成前，五项状态仍为 `open`，不得提前写成已闭合。
  2. 冻结 [I-011-004-a012-remediation-contract.md](attachments/I-011-004-a012-remediation-contract.md) **v0.1.0**：以 `roles.assign`、操作者权限子集与 admin 保护形成角色委派边界；以密码专用控件/原字节传输、8～72 bytes、改密同事务撤销 refresh token 形成凭据边界；为自定义角色补齐 permission/menu grants 的可验证 CRUD 与有效权限投影；清理活动 CI/生产 transport 的 records 残留；把最后管理员检查与删除放入同一事务并验证受影响行数。
  3. `I-011-004` → **verified**，只解除 F-003 产品边界与整改方案冻结门禁。I-011-001 v0.2.0 中“grant 管理界面非目标/自定义角色无 grants”的历史冻结条款，在 A-012 整改范围内由 I-011-004 v0.1.0 限定取代；仍不扩展为完整 IAM、SSO/SCIM、复杂策略编排或多租户。
  4. 前端可为这次修复增加通用、结构化能力（密码字段不 trim、基于行布尔字段禁用动作），并把 `records.ts` 通用 transport 改为资源中性命名；这些是 A-012 修复，不再受 I-011-003 v0.2.0 的历史 S4 零 Renderer diff 断言约束。历史 S4 收据保持不改写。
  5. F-006 保持 `recommended / open / non-blocking`，不与五项 required 的实施和复审捆绑。
- **理由**：五项均有明确代码、Schema 或活动流水线落点；其中 F-001/F-002/F-005 是生产安全不变量，F-004 是交付证据真实性，F-003 则决定 users/roles 是否形成可实际使用的授权闭环。直接修复比接受残余风险更符合 VP-002 的生产 Admin 边界。
- **未选方案**：seed-only 并移除 roles 管理面（用户明确未选）；对任一 finding 采用 residual/overruled（用户明确选择 fixed）；把 grant 路径扩大为完整 IAM（超出本目标，仍为非目标）。
- **信息门禁**：`I-011-004` verified；F-001～F-005 实施/回归/复审门禁仍 open。GOAL-011 保持 `active / 5/5`，GOAL-010、Root 与 VP-002 状态不变。
- **后续**：按 I-011-004 v0.1.0 实施 API/store/Schema/Renderer/CI 修复，形成 API/Web/build/E2E/并发与静态洁净度收据；随后请求只覆盖 A-012 F-001～F-005 的 `/audit` finding-closure 复审。

## D-007 · 响应 A-013：确认 A-012 五项 fixed 并恢复关门

- **日期**：2026-08-04
- **状态**：accepted
- **输入**：候选提交 `fb5cd067156a39f0d879760961db2bac0d4266d0` 已按 D-006 / I-011-004 v0.1.0 实施 A-012 F-001～F-005；A-013（independent · finding-closure · pass）限定复审五项，逐项确认 `fixed`，无新增 required。
- **决定**：
  1. 采纳 A-013 `pass`。A-012 F-001～F-005 由 `open` 按合法路径闭合为 **`fixed`**；A-012 原始 `fail` 保留为发现当时的历史意见，不改写。
  2. F-006（legacy roles JSON 双写）继续为 `recommended / open / non-blocking`；不把可选后续卫生项提升为关门阻断，也不静默写成已处理。
  3. GOAL-011 当前无开放 required finding、无到期 required 信息项；原 S1～S5 均已完成，故从 `active / 5/5` 恢复为 **`done / 5/5`**。
  4. 同步 goal-tree 的树与状态表，并向 GOAL-010 追加新的子目标交接事实。GOAL-010 仍为 `active / 3/5`，S4/S5 不因子目标关门自动勾选；Root A-002 F-002-001、Root 与 VP-002 状态均不在本决策关闭。
- **fixed 关闭证据**：实现提交 `fb5cd06`；`02-execution.md` 的 API/Web/build/E2E/并发/scoped residue/Linux Compose/disposable smoke/重启持久化矩阵；A-013 对权限委派、密码/refresh、grant 管理、records 活动面及最后管理员事务的逐项独立核对。
- **证据边界**：GitHub-hosted Actions 尚未触发，不能写成远端 CI 已通过；I-011-004 要求不可用环境证据单列，并未把 hosted run 设为本次 finding fixed 的额外门禁。
- **理由**：用户已明确选择五项 fixed，实施与本地验收可核对，限定独立复审无新增 required；P-003 的 finding closure 与重新关门条件均已满足，没有 residual/overruled 或继续 fail-closed 的依据。
- **影响**：GOAL-011 `active → done`，progress 保持 `5/5`；F-001～F-005 fixed，F-006 recommended/open/non-blocking；父级及愿景层状态/进度不变。
- **后续**：由 GOAL-010 的 `/govern` 流程独立评估其 S4 验收门与 S5 close-out，不以本次子目标关门替代父级审计或 Root finding closure。
