---
id: A-002
goal: GOAL-002-r1-bounded-research
source: independent
auditor: grok build（grok-4.6 · xhigh）
date: 2026-08-14
scope: R1 调研结果（D-001 判据 / E-002 基架对照 / I-011-001 分档清单 / Root 路线图回写）
audit_type: design-plan
verdict: conditional
parent: GOAL-002-r1-bounded-research
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-002 · 独立审计意见（grok build · R1 有界调研结果）

> 代贴（保留 source: independent）：本意见由 grok build（grok-4.6 · xhigh）于独立会话生成，本仓库 /govern 编排器代为落盘，未改写原文。响应与闭合由 /govern 处理。

---

独立交叉审计先按审计技能与原则读齐证据，再对照 GOAL-002 的调研产出做只读核对。证据面先通读工作区、GOAL-002 五件套、分档清单和既有 A-001，再对照 Charter / VP / 基架代码。接下来核对基架代码与协议面，验证「已覆盖」和「真·待建」是否站得住。# A-002 · 独立审计意见（R1 有界调研结果）

- **source**：independent ｜ **auditor**：grok build（grok-4.6 · xhigh） ｜ **date**：2026-08-14 ｜ **scope**：R1 调研结果——调研方法与三档判据（D-001）、基架对照（E-002，C-01～C-11）、三档分档清单（`attachments/I-011-001-tiered-inventory.md`）、Root 路线图回写（E-003 / Root `00-meta` R1～R4） ｜ **audit_type**：design-plan ｜ **verdict**：**conditional**

工作区：`workspace-011-admin-functional-modules`（canonical：`docs/workspaces/workspace-011-admin-functional-modules/`；`vision_role: delivery`；`primary_plan: VP-011-admin-functional-modules`）。被审目标：`GOAL-002-r1-bounded-research`（`status: done` 5/5）。本会话按任务约束**只出意见、不改任何文件**；正式台账落盘须由 `/govern` 写入 `03-audit/A-002-*.md` 并更新索引。

---

### 范围与区间

- **已校验工作区绑定**：`workspace.md` 的 `id` / `root_goal` / `canonical_scope` / `plan_refs`+`primary_plan` 与 Root `GOAL-001-admin-functional-modules` 一致；`shared_materials_catalog: none`，清单未把共享资料当成事实。
- **已通读**：GOAL-002 五件套与 ledger（D-001、E-001～E-003、A-001、I-011-001）、Root `00-meta`、VP-011、Charter `@0.2.0`。
- **已只读核对基架**（本区允许的代码/愿景层，不审他区目标状态）：`apps/api/internal/modules/` 名册、users/settings/activity schema、`handler/users.go` / `auth.go` / `upload.go`、`composition.go` home 推导、`docs/vision/protocol-inventory-v2.7.0.md`。
- **不在本 scope**：R2 实施方案、实现代码、VP-011 激活本身、他区目标推进。
- **P-005**：GOAL-002 的 I-001 标 `verified`；**Root I-001 仍为 `open`**（影响门禁 = R2 立项，最晚阶段 = R1 结束）。I-002 为 non-blocking，两端均仍 `open`/「已核对」。

---

### 成果（有证据）

1. **方法与分层约定自洽（方向层）**  
   D-001 三档判据与 VP-011「三档方法论」逐字对齐：一等公民 = 业界普遍 + 几乎所有 Admin 需要 + 基架未覆盖；常用 = 高频非普遍；增补 = 低频按需。产出落 Root 路线图、VP 不逐条改写，符合分层约定。边界写明：只出清单、不实现、不改 Profile/模块矩阵/Manifest 装配语义/协议 pin/共同门禁。

2. **基架对照主体成立**  
   模块名册可核对：`users` / `roles` / `settings` / `activity` / `operationlog` / `dev/examples`。  
   - C-02～C-11 的「已覆盖、不重复立项」成立：roles+RBAC、生产面 users CRUD（含 batch delete）、settings 四类（General/Branding/Localization/Appearance）、zh-CN/en-US/auto、tokens+light/dark、activity≠通知、mvp/admin/demo 与 home 推导、`/api/upload` 控件级上传、管理员改密（`changeUserPassword`）+ `token_version` 吊销、登录锁定 423 / 限流 / 改密吊销。  
   - **C-03 与 F-01 口径可区分**：C-03 是渲染器/协议 CRUD **平台能力**（生产 users/roles + demo 代表页）；F-01 是 **生产 Profile home 产品面**。e2e/`composition` 固定：demo → `overview`，mvp/admin → `users`。此区分成立。

3. **「真·待建」主缺口可复核**  
   - **F-01**：`overview.json` 仅欢迎文案；且只随 `dev.examples`。生产面无 dashboard。  
   - **F-02**：`apps/api` 无 CSV/Excel 导入导出端点（仅有 multipart 单文件上传）。  
   - **F-03**：改密是管理员视角的 users PATCH；`account.go` 仅 `$context`/`me`；无会话列表/自助吊销 API。  
   - **F-04**：代码面无站内通知模块；`operationlog` 事件为 `auth.*` / `users.*` / `roles.*` / `settings.update`，不是收件箱。  
   通知 ≠ 操作日志，判断正确。

4. **入池/不入池与 Charter**  
   多租户/白标、运行时插件市场、私增协议语义明确不入池，与 Charter 非目标一致。B-05 订阅/套餐标注须先 `/vision`，正确。订单/钱包**入池**符合 VP-011「常用业务领域典型」与 Charter「后续 VP 候选能力」（本 VP 即该后续波次）；问题在**档位**，不在「能不能进候选池」。

5. **路线图回写主表一一对应**  
   Root R2 = F-01～F-06，R3 = S-01～S-12，R4 = B-01～B-09；R2/R3 依赖 R1/R2，R4 依赖 R1（可按需拉起），与 VP-011 退出判据 3（增补不由 VP 强行关闭）一致。R4 各项有触发条件建议。

---

### Findings

#### F-001 · required · 严重度 **high** · 关联 I-001 / 清单 F-05、F-06、S-07、S-08

**主张**：一等公民档对 **F-05 订单、F-06 钱包** 未按 D-001 / VP-011 书面判据「几乎所有 Admin 都需要」应用；实际用的是未写入判据的第二轴「用户点名典型 + 电商/平台惯例」。同类领域能力 **S-07 类目、S-08 商品** 却留在常用档，同一判据漂移。

**证据**：

- 书面判据（D-001 / VP-011）：一等公民 = 业界普遍存在、**几乎所有 Admin 都需要**、且基架未覆盖。
- I-011-001 对 F-05/F-06 的理由是「电商/平台 fork 普遍；**用户点名代表**」「用户点名典型；平台后台普遍」——没有论证「几乎所有 Admin」。
- 同一来源轴上：类目同为用户候选池、商品为「业界惯例补充」，却在 R3。
- 所引通用样本（Appwrite 综述、vue-element-admin / Ant Design Pro、simple-admin-core）支撑的是仪表盘/个人中心/导入导出/通知/字典/日志，**不是**订单/钱包。订单/钱包只被电商/多商户样本覆盖。
- 影响：Root R2 与 VP-011 退出判据 1 会把「交付订单 + 钱包模块」锁进**第一波与 VP 关门条件**。这不是调研笔误，是范围承诺。

**影响门禁**：Root I-001 关闭、R2 立项（F-05/F-06 是否进入第一波）。

**建议闭合**：三选一并留痕——(a) **fixed**：F-05/F-06 降为常用（R3），与 S-07/S-08 同族；或 (b) **改写 D-001** 显式增加第二轴（例如「用户点名的典型领域可入第一波」）并重算 S-07/S-08 是否跟进；或 (c) 用户书面 `accepted-residual` / `user-overruled`（写明 VP-011 退出仍绑定这两项）。

---

#### F-002 · required · 严重度 **med** · 关联 清单 S-12

**主张**：S-12「复用 records retirement 基建（0006 `records_retire`）」**事实错误**。`0006 records_retire` 是历史演示实体 `records` 的**退场 DROP**，不是软删除/回收站产品基建。

**证据**：`apps/api/README.md` 写明 records 已从产品运行面退场，仅留历史迁移账本与 `records.*` operation-log 兼容值，「不恢复产品 CRUD」。迁移测试将 `records_retire` 与 DROP 绑定。仓库内 tombstone 语义是「退役迁移保留描述符」，不是回收站。

**影响门禁**：I-011-001 作为 R3 依据的可信度；S-12 方案不得把 0006 当可复用能力。档位（常用）可保留。

**建议闭合**：fixed——删掉/改写「复用 0006」；写明软删除/回收站需新持久化与管理 UI，或指出真实可复用点（若有）并给出路径。

---

#### F-003 · required · 严重度 **med** · 关联 清单 C-01、C-11；缺口未入池

**主张**：C-01「账号管理（…状态/锁定）」**过标**。生产 users 页/表无 `status`/`enabled`，无管理员启停动作。存在的是 C-11 的**登录失败自动锁定**（`locked_until` / 423），不是账号生命周期产品能力。该真缺口未进入 F/S/B 任一档。

**证据**：

- `users.json` 动作仅为 create / update name / roles / **管理员改密** / delete；字段无状态。
- users 查询列为 `id, username, name, roles, password_hash, token_version, failed_login_count, locked_until, …`，无 enabled/disabled。
- 自动锁定在 `auth.go` / `account_lock` 迁移，属于 C-11。

**影响门禁**：I-001 完整性（候选池是否漏掉业界普遍的账号启停）；避免 R2 把 C-01 当成「启停已有」。

**建议闭合**：fixed——将「管理员启用/停用（及手动解锁）」补入清单（建议一等公民或并入 F-03），并修正 C-01 表述为「CRUD + 角色分配；不含产品态启停」；或书面说明不入池理由。

---

#### F-004 · recommended · 严重度 **med** · 关联 I-002、清单 F-01/F-02、A-001 F-002

**主张**：I-011-001 §7 用「S3-protocol-judgment 已冻结 9 covered / 0 protocol-gap」推出 F-01/F-02「预计呈现自由 + fail-open」，是**分母用错**。该 9/0 是 VP-008 准入波次对**当时共性能力**的分类，不是对「生产 dashboard / CSV 导出」的新判定。

补充证据（A-001 未写）：

- 愿景层协议清单 `docs/vision/protocol-inventory-v2.7.0.md` §2.5 已列信息性场景 **`grid-dashboard`**（非语义权威，但不是「协议面毫无 dashboard 痕迹」）。
- 上游样例点名 `user-profile-*`、`order-list-*`，与 F-03/F-05 相关，清单未引用。
- `node.schema.json` 将 `export` 列为 permissions **扩展动作键示例**（协议不限制键名），不是导出能力契约。

**影响门禁**：F-01/F-02 的 R2 **方案冻结**（不自动阻断立项，但方案不得把 9/0 当已放行）。

与 A-001 F-002 **同向加强**，不升格为与 self 冲突的必改项。

---

#### F-005 · recommended · 严重度 **med** · 关联 清单 F-04；A-001 F-001

**主张**：同意 A-001 F-001。F-04 把「业界通用站内通知」与用户点名的业务领域「通知」合成一个模块，边界未冻（系统通知 / 业务通知 / 与 S-05 公告 / B-09 模板的切分）。不否定「通知中心」作为通用能力的档位倾向。

**影响门禁**：F-04 方案冻结。

---

#### F-006 · recommended · 严重度 **low** · 关联 候选池完整性

**主张**：在已引用的中文企业后台/Go admin 样本轴上，**组织/部门/岗位**、**登录日志独立视图**未入池也未写入「明确不入池」。登录事件已在 operation log（`auth.login/logout/refresh`），缺的是独立产品面，不是日志源。组织架构则完全缺席。

不要求为每个低频项建子目标；R2 前应显式补入或排除。

---

#### F-007 · recommended · 严重度 **low** · 关联 回写一致性、Root I-001

**主张**：E-003 声称 I-001 → verified 且路线图已回写，但：

- Root `00-meta` **I-001 仍 `open`**（R2 立项门禁仍开）；
- `workspace.md` 纲领表 R2/R3/R4 仍写「待调研产出」，与 Root 已细化的 F/S/B 不一致；
- Root R1 单元格写 `attachments/I-011-001`，易被读成 Root 自身附件；权威路径应是  
  `docs/workspaces/workspace-011-admin-functional-modules/GOAL-002-r1-bounded-research/attachments/I-011-001-tiered-inventory.md`。

路线图**内容**回写已发生；指针与信息项未对齐。Root I-001 在本意见 required 项闭合前**不应**标 verified。

---

#### F-008 · recommended · 严重度 **low** · 关联 F-01、F-05

**主张**：

- F-01「复用 overview 页能力」不成立：`overview.json` 只有 section+text，无指标/卡片/grid。可复用的是模块贡献与页面装配模式，不是 dashboard 能力。落地 `admin.dashboard` 并进入 mvp/admin 默认启用集，是 Profile **内容扩展**，须用既有贡献机制，并更新 `adminFunctionalOrder`（现为 users→roles→settings→activity）；这不是改装配**语义**，但必须在 R2 方案写清，以免踩中「不改变 Profile 默认集」的字面读法。
- F-05 若含商品/类目引用，与 R3 的 S-07/S-08 依赖未写。R2 单独做订单须声明最小实体（无目录外键）或允许桩。

---

### 必改项汇总

| ID | 未闭合前阻断 | 建议路径 |
|----|----------------|----------|
| **F-001** | Root I-001 关闭；**R2 立项不得把 F-05/F-06 当已冻结的一等公民** | 降档 **或** 改写 D-001 第二轴并重算类目/商品 **或** 用户书面 residual/overruled |
| **F-002** | 不得带着「复用 0006」进入 S-12 方案 | 改正清单事实 |
| **F-003** | 不得把 C-01 当作「账号启停已覆盖」 | 修正 C-01 + 补入池或书面排除 |

无其他 required。F-004/F-005 为方案期 recommended，与 A-001 同向。

---

### 与既有意见的异同（A-001 self = pass）

| 点 | A-001 self | 本意见 independent |
|----|------------|-------------------|
| 总 verdict | **pass** | **conditional**（存在 high/med required） |
| 方法、来源落盘、F-01～F-04 主缺口、不入池 | 认可 | **同意** |
| 通知模块边界 | recommended F-001 | **同意**，本意见 F-005 |
| dashboard/导出协议 | recommended F-002：立项时按 S3-protocol-judgment 留痕 | **同向加强**（F-004）：9/0 不能外推；须另做 F-01/F-02 协议对照，并引用 `grid-dashboard` / export 扩展键 |
| F-05/F-06 档位 | 仅在偏差节称「商品、数据权限等有判断空间」 | **升级为 required F-001**：订单/钱包档位才是主漂移 |
| S-12 / C-01 启停 | 未提 | **新 required** F-002、F-003 |
| 冲突？ | — | **无 P-004 级「一要一否」**。self 无必改项；independent 新增必改。编排器必须响应本意见 required，不得因 A-001 pass 放行 R2。 |

A-001 将 I-001 视为已验证，偏乐观：优先级分类（I-001 的核心问题）在 F-05/F-06 上未按书面判据闭合。

---

### 结论 + 建议给编排器/用户的下一步

**结论**：R1 作为有界调研**做了该做的结构工作**——计划冻结、来源锚点、11 项已覆盖对照、清单三档、Root R2/R3/R4 主表回写，且未触碰 Charter 非目标、未在本目标改协议/Profile 语义。不能给 pass，是因为**分档主产物在两处关键判定上不自洽或事实错误**：把领域模块抬进一等公民（F-001）、把退场迁移当成回收站（F-002）、把自动锁账号写成已覆盖的状态管理并漏掉启停（F-003）。

**建议编排器（`/govern`）**：

1. **代为落盘**本意见为 `docs/workspaces/workspace-011-admin-functional-modules/GOAL-002-r1-bounded-research/03-audit/A-002-r1-research-independent.md`，并更新 `03-audit.md` 索引（A-002 / independent / conditional）。本会话未写盘。
2. **不要改** GOAL-002 `status/progress`（已是 done）；**不要**在 F-001～F-003 合法闭合前关闭 Root I-001 或开 R2 子目标。
3. **P-004 只问一件事**（建议选项已排序）：F-05/F-06 如何处置？  
   - **建议**：降为 R3 常用，与类目/商品同族；R2 先做 F-01～F-04（及补入的账号启停，若采纳 F-003）。  
   - 若用户坚持第一波做订单/钱包：先改 D-001 写明第二轴，并决定 S-07/S-08 是否跟进，再立项。
4. 清单小修正（F-002/F-003/F-007）可与上一问同批做。
5. R2 各子目标方案仍须消化 A-001 F-001/F-002 与本意见 F-004/F-005/F-008（协议对照、通知边界、home/依赖）。

**建议用户下一句**：  
`/govern 响应 workspace-011 GOAL-002 A-002：F-05/F-06 降为常用；修正 S-12 与 C-01；再关闭 Root I-001。`

---

### 声明

本意见不修改 `status` / `progress` / 方案正文 / goal-tree。响应与 finding 闭合由 `/govern` 处理。保证等级为框架默认 **L0**（入口分离），不是第三方鉴证。


