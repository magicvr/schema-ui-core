---
id: A-002
doc: audit-entry
goal: GOAL-001-key-rotation-and-backup
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
auditor: grok-build (grok-4.6 · reasoning high)
---

> 独立审计员按本轮书面要求**未落盘**；由编排器代贴至 Root `03-audit/A-002` 并更新索引。`source: independent` 保持不变，内容逐字未改（仅 Markdown 包裹与文件名）。项目级决策依据：`docs/architecture/independent-audit-execution.md`。

# A-002 · Root 关门独立审计（R1～R5 全阶段 + 判据 1–6）（2026-08-22）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型 / scope**：close-out · Root GOAL-001 整体（workspace-016 纲领 R1～R5；VP-016 方向级退出判据 1–6；信息门禁 I-001～I-005；五阶段意见台账；越界核对；goal-tree / meta 一致性）
- **verdict**：**conditional**
- **工作区**：`workspace-016-key-rotation-and-backup`（Root `GOAL-001-key-rotation-and-backup`；`shared_materials_catalog: none`；本意见未把共享资料当证据）

本意见不修改 `status`/`progress`；响应由 `/govern` 处理。

### 范围与区间

| 项 | 值 |
|----|----|
| 被审目标 | GOAL-001-key-rotation-and-backup（Root 关门） |
| 工作区 / 合同 | workspace-016-key-rotation-and-backup / VP-016-key-rotation-and-backup（`active`，架构 A5） |
| 子目标 | GOAL-002～GOAL-005（committed `done`）+ GOAL-006-r5-dual-path-evidence（工作树未提交，`active`） |
| 代码基线 | 开区提交 `5195104` → HEAD `1dc6975`；工作树另有 Root A-001 与 GOAL-006 未提交 |
| 四份 checkpoint | `c96e963`（R1 配置面）/ `8116565`（R2 双密钥）/ `1b8e9b0`（R3 恢复测试）/ `1dc6975`（R4 证据整合） |
| 代码切片 | `apps/api` 10 文件 + 根 `compose.yaml`（见判据 5）。`apps/web`、`internal/modules`、`docs/vision/charter.md` 相对基线 **零 diff** |
| 不在本 scope | 把 Root 标 `done`、改 VP status、KMS/kid/refresh 模型/热加载/第二套 dump 的实现审查（只核有无越界） |
| P-005 | I-001～I-004 required 已 verified（有决策+代码+本审复跑）；I-005 non-blocking 仍 collecting，默认「previous 可验」已实施，不阻断。见 F-002 镜像滞后 |
| 共享资料 | catalog = none；无固定引用被当作证据 |

### 成果（有证据）

| 成果 | 证据 |
|------|------|
| R1 配置面：current 沿用、previous 新增、生产同强度 + 同值守卫、缺省单密钥 | Root D-002；`config.go` 字段 / YAML `jwt_secret_previous` / env `AUTH_JWT_SECRET_PREVIOUS` / `ValidateProd`；本审 `TestJWTSecretPreviousConfig` **8/8 PASS**、`TestValidateProd` **9/9 PASS** |
| R2 双密钥：签发只用 current；校验 current→previous；过期不延长；无 `kid` | `issue()` 仍 `SignAccessToken(a.secret, …)`；`verifyAccess`；`SignAccessToken` 未设 `kid`；本审 `TestDualKeyRotationOverlapWindow` **4/4 PASS** |
| R2 composition 接线（含 A-002 F-001 钉死） | `newAuthenticator` → `NewWithRepositoryAndPrevious`；`NewApp` 签名未改；本审 `TestNewAuthenticatorWiresPreviousSecret` **PASS** |
| R3 双方言恢复循环 | `post_rotation_recovery_test.go`：VACUUM INTO / `pg_dump -F c`→`pg_restore` + ledger 指纹；本审 `TestSQLitePostRotationRecovery` **PASS（0.54s）**、`TestPostgresPostRotationRecovery` **PASS（5.41s，`R16_PG_DUMP_CONTAINER=r2-pg-probe`）** |
| R4 缺省单密钥 | 本审 config 8/8 + ValidateProd 9/9 + composition 空 previous 拒绝 old-key + `go test ./cmd/server/ -count=1` **ok（11.5s）** + `docker compose config` 输出 `AUTH_JWT_SECRET_PREVIOUS: ""` |
| R5 双路径新鲜实跑（产品面） | 本审独立复跑 GOAL-006 E-001 四项载体 **全 PASS**（轮换路径 4+1、恢复路径 SQLite+PG） |
| 四份 git checkpoint 名实相符 | `c96e963` 21 files 配置面；`8116565` auth+composition；`1b8e9b0` 仅新增恢复测试 + 台账；`1dc6975` 仅 GOAL-005 文档 |
| GOAL-003 recommended 闭合可核对 | F-001 钉死测试本审 PASS；F-002 索引已列 E-002/E-003；F-003 E-002 v1.1 已收窄。三路径 = **fixed** |
| 越界未发生（产品 diff） | 见下判据 5 |
| vet | 本审：`apps/api` `go vet ./...` **0 finding** |
| JWT 相关包回归 | 本审：config / auth / composition / handler / `cmd/server` 指定套件 **ok**（未把整包 `go test ./...` 当作本审已复现事实） |

### 对照成功标准（VP-016 方向级判据 1–6）

| # | 判据 | 判定 | 证据 |
|---|------|------|------|
| 1 | JWT 轮换合同：可配置 current+previous；新签发只用 current；重叠窗 previous 可验。立即失效仅作 I-016-005 有界残余 | **达成** | R1 配置面 + R2 `verifyAccess` + 4/4 重叠窗测试。I-005 未改为立即失效；默认 previous 可验已冻结并实施 |
| 2 | 未配置 previous 时本地/Compose 默认仍能开发与快测；轮换不是 mvp/dev 启动硬依赖 | **达成** | 缺省 previous 为空；生产缺 previous 不 fail-closed；compose 可选透传默认 `""`；`cmd/server` 启动测试 ok。`resolveJWTSecret` 仍只解析 current |
| 3 | 轮换后恢复：SQLite `VACUUM INTO` **与** PG `pg_dump`/`pg_restore` 上都有轮换后启动 + 鉴权证据 | **达成** | 本审双方言循环均 PASS（非 skip）。T4 = `NewApp(K2, prev=K1).Start`；T5 = A1/A2/A3（重叠窗旧 access / 新签发仅 K2 / opaque refresh 连续）。T5 走恢复库上独立构造的 Authenticator，不是 T4 仍在跑的 HTTP 服务——与剧本「恢复库 + 轮换后密钥集」相符；生产接线另由 composition 钉死 |
| 4 | 显式双密钥下：一轮换路径 **与** 一轮换后恢复路径都有可核对证据 | **产品面达成；台账面未齐** | 四项载体本审全 PASS。但交付目标 GOAL-006 未入 tree、五件套不完整、自身检查点仍 pending（F-001）。不得把 Root A-001「五阶段全部关门」写成已与 canonical 状态一致 |
| 5 | 未进入 A3/KMS/PITR/Admin 功能/业务域；未改 Charter；未假装热加载或第二套 dump | **达成** | `git diff --name-only 5195104`：`apps/web`、`internal/modules`、`docs/vision/charter.md` 为空。代码 = config/auth/composition 测试与接线 + YAML 样例 + README + `compose.yaml`。Charter 仍 `schema-ui-core-admin-foundation@0.2.0`，`primary_workspace` 仍 workspace-001。无运行期换钥 API；备份只消费 `VACUUM INTO` 与官方 `pg_dump`/`pg_restore`。GOAL-006 E-001「共 10 文件」漏计根 `compose.yaml`（第 11 个，仍在分母内） |
| 6 | 开放 required finding = 0（或已合法闭合） | **本意见之前为 0；本意见打开 F-001 required** | 既有台账：GOAL-002 A-001、GOAL-003 A-001+A-002（F-001/2/3 recommended 已 fixed）、GOAL-004 A-001、GOAL-005 A-001、Root A-001 —— 开放 required 均为 0，闭合路径合法。本条 independent 新增 F-001（med required）未闭合 ⇒ 判据 6 当前不满足 |

### 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 权威表（00-meta） | 本审核对 |
|----|------|----------|-------------------|----------|
| I-001 | required | R1 | **verified**（D-002） | 键名/熵/fail-closed/同值守卫已落地；secret 错误只点名键名。成立 |
| I-002 | required | R1 | **verified**（D-002） | `NewServiceCredentialToken` = CSPRNG + SHA-256，不派生自 JWT secret。服务凭证书面出局成立 |
| I-003 | required | R2 接入前 | **verified**（GOAL-003 D-001） | 重叠窗 = previous 配置存续期；无 `kid`；refresh opaque 不受签名密钥影响。成立 |
| I-004 | required | R3 接入前 | **verified**（GOAL-004 D-001） | T0–T5 + A1/A2/A3 + 双方言既有工具。核心问题已答。D-001 正文「17↔17 无 GUC」与实跑/测试代码不一致（F-003），**不**足以把 I-004 打回 collecting |
| I-005 | non-blocking | R2 | **collecting** | 默认措辞随 VRev-035 冻进退出 1；本波按 previous 可验交付；无用户书面「立即失效」残余。处理合规，不阻断关门 |

Root `01-decision.md` 自称镜像却把 I-004 仍写成 collecting（F-002）。VP-016 信息表 I-016-001～004 仍 collecting，未回流 Root 结论（F-005）。

### Findings

#### F-001 · R5 交付目标未入 tree、五件套不完整，与 Root A-001「五阶段关门」名实不符（required · med · open）

- **严重度**：med  
- **建议**：required  
- **状态**：open  
- **关联**：判据 4 台账面、判据 6、AGENTS §3/§7、GOAL-006 检查点 1–4  
- **描述**：产品证据（四项测试）本审已独立复现，**不是**虚构。但 R5 的交付目标 `GOAL-006-r5-dual-path-evidence` 在 canonical 状态上尚未可关门：  
  1. `goal-tree.md` 仍写「R5 子目标待开」，状态表无 GOAL-006 行（新建目标未同步树 = AGENTS §7 硬约束未满足）。  
  2. 缺五件套索引文件 `01-decision.md`（`01-decision/D-001-*.md` 在，索引文件不在）。  
  3. GOAL-006 `00-meta` 检查点 1–4 全部 **pending**，`progress` 未写；`02-execution.md` 索引为「暂无」，但 `02-execution/E-001-fresh-evidence-runs.md` 已存在；`03-audit.md` 尚无本目标条目。  
  4. Root A-001（self）写「五阶段全部关门 / 判据 4 达成 / GOAL-006 证据登记完成」；Root `00-meta`、`workspace.md` 仍写 R5 **未开始**、progress **4/5**。A-001 相对 canonical 记录超前。  
  5. GOAL-006 整目录与 Root A-001 均未提交（HEAD 停在 R4 `1dc6975`）。  
  D-001 关门顺序把 tree 终态同步放在 independent 之后，这可以解释 Root 仍 4/5，**不能**解释「目标已创建却不在树上」和「缺 `01-decision.md`」。Root **不得**在 F-001 闭合前标 `done`。  
- **证据**：`goal-tree.md` 树与状态表；`GOAL-006-r5-dual-path-evidence/00-meta.md`、`02-execution.md`、`03-audit.md`；`01-decision.md` 不存在；`git status`：`?? GOAL-006-r5-dual-path-evidence/`、`?? …/A-001-root-closeout-self.md`、`M …/03-audit.md`；Root A-001 vs Root `00-meta` 路线图 R5 行。  
- **建议动作**：补 `01-decision.md` 索引并登记 D-001；把 E-001 写入执行索引；按事实把检查点 1–2 标 done（3–4 等本意见响应后）；goal-tree 增 GOAL-006 行；`workspace.md` R5 行改为进行中/完成（与响应后状态一致）。无需重做测试。

#### F-002 · Root 决策镜像与 workspace 纲领表滞后（recommended · low · open）

- **严重度**：low  
- **建议**：recommended  
- **状态**：open  
- **关联**：I-004 镜像、workspace 纲领表  
- **描述**：`00-meta` 为 I-00N 权威，I-004 已 verified；`01-decision.md` 镜像表仍 collecting，违反该文件自己「须保持同号同状态」的声明。`workspace.md` 纲领表 R5 仍「未开始」，与磁盘上已存在的 GOAL-006 矛盾。不否定 I-004 关闭本身。  
- **证据**：Root `01-decision.md` I-004 行 vs `00-meta` I-004 行；`workspace.md`「纲领阶段」R5 行。  
- **建议动作**：镜像与纲领表改到与 `00-meta`/tree 一致。GOAL-003 A-002 已指出过同类 I-003 镜像滞后，I-003 已改、I-004 未跟。

#### F-003 · I-004 剧本正文与实跑/测试代码的 PG 版本组合不一致（recommended · low · open）

- **严重度**：low  
- **建议**：recommended  
- **状态**：open  
- **关联**：I-004 文档精度（不重开信息项）  
- **描述**：GOAL-004 D-001 写「容器内 17 客户端 ↔ 17 服务端，同主版本无 GUC 告警」。GOAL-004 E-001 / A-001 与 `post_rotation_recovery_test.go` 则按跨版本客户端 + **1 条** `unrecognized configuration parameter` ignored error 放行，并以 `schema_migrations` count+checksum 指纹补强。本审 PG 循环 **PASS**，合同消费的是官方 `pg_dump -F c` / `pg_restore`，不是第二套 dump。决策正文若保持「无 GUC」会误导后续操作者。  
- **证据**：GOAL-004 `01-decision/D-001-recovery-playbook.md`「两方言证据命令」表；`post_rotation_recovery_test.go` `dump`/`restore` 注释与 ignored-error 分支；GOAL-004 E-001。  
- **建议动作**：勘误 D-001 为「允许 VP-013 已记录的跨版本 GUC 告警类，以 ledger 指纹为准」，不要改测试。

#### F-004 · 跨区裸 id 与空指针引用（recommended · low · open）

- **严重度**：low  
- **建议**：recommended  
- **状态**：open  
- **描述**：本区文档多次裸写 `GOAL-006 D-002` 指 **另一工作区** 的备份合同（Root `00-meta` I-004 证据列、GOAL-004 E-001/A-001、恢复测试注释）。在本区该裸 id 会撞到 `GOAL-006-r5-dual-path-evidence`。GOAL-005 D-001 第 6 行写「GOAL-005 D-004 既有合同」，本区 GOAL-005 只有 D-001。GOAL-004 E-001 写「回归结果见 E-002」，本目标无 E-002。均不损害产品结论。  
- **证据**：上列文件；workspace-protocol §2.6（文档默认 Q2 路径）。  
- **建议动作**：改为 Q2 路径限定引用；删或改空指针。

#### F-005 · VP-016 信息表未回流 Root 已关闭项（recommended · low · open）

- **严重度**：low  
- **建议**：recommended  
- **状态**：open  
- **描述**：VP-016 `I-016-001`～`I-016-004` 仍为 collecting；Root 对应 I-001～I-004 已 verified。Root 门禁用 Root 表。VP 关门或下一次 Vision 回流时应同步，以免组合层误读「信息项仍开」。  
- **证据**：`docs/vision/plans/VP-016-key-rotation-and-backup.md` 信息需求表 vs Root `00-meta` 信息表。  
- **建议动作**：VP 信息行改为 verified 并链到 Root/GOAL-00N 决策；属 VP 卫生，不作为本 Root 产品重开条件。

无 high required。无到期未处理的 required 信息项。无对同一必改项的一要一否。

### 必改项汇总（required）

| ID | 摘要 | 闭合路径 |
|----|------|----------|
| F-001 | 登记 GOAL-006 到 goal-tree；补 `01-decision.md`；执行索引登记 E-001；检查点 1–2 与事实对齐。完成前 Root 不得 `done` | `fixed`（台账对齐即可；不要求重跑测试） |

F-002～F-005 为 recommended，不阻断在 F-001 闭合后的 Root 关门。

### 与既有意见的异同

| 点 | Root A-001 (self) | 本意见 (independent) |
|----|-------------------|----------------------|
| verdict | pass（待 independent） | **conditional** |
| 判据 1–3、5 产品证据 | 达成 | **同意**，且本审独立复跑（含 PG 循环，非 skip） |
| 判据 4 | 「GOAL-006 E-001 四项新鲜实跑全 PASS」⇒ 达成 | 测试 **同意**；台账 **不同意已关门**（F-001） |
| 判据 6 | 开放 required = 0 | 本意见之前同意为 0；本意见打开 F-001 后不为 0 |
| I-001～I-004 / I-005 | verified / collecting 合规 | 同意；补镜像滞后（F-002）与 VP 未回流（F-005） |
| 越界 | 10 文件、charter 零改动 | 同意无越界；实际 10 个 `apps/api` + `compose.yaml` |
| GOAL-003 F-001/2/3 | 已 fixed | 同意；钉死测试本审 PASS |
| 全仓 `go test ./...` | 未再当作 Root 关门证据 | 同意不引用；本审只复现 JWT 相关包 + vet + 四项载体 |

无 verdict 冲突需要 P-004 裁同一必改项（self 未提 F-001；独立审新开）。编排器应响应 F-001，而不是把 A-001 pass 当成可以压过本条 required。

### 核对审计重点（摘要）

1. **六条退出判据**：1–3、5 名实相符且本审复现；4 产品面相符、台账面不符；6 因 F-001 当前不满足。  
2. **五阶段意见台账**：既有 A 条目开放 required = 0，GOAL-003 recommended 闭合合法。本条新增 1 条 required。GOAL-006 自身尚无 A 条目（Root 级审计按 D-001 落在 Root，可接受）。  
3. **I-001～I-004 verified 有据**；I-005 collecting + 默认实施合规。  
4. **越界**：无 A3/KMS/PITR/Admin 页/业务域/Charter 改动/热加载/第二套 dump。  
5. **goal-tree 与 meta**：R1–R4 子目标与 tree 一致（done）；R5/GOAL-006 不一致（F-001）。Root 仍 `active` 4/5 与「independent 通过前不 done」相符。  
6. **夸大/虚构**：测试结果 **未虚构**（本审全 PASS）。夸大点是 A-001 把未入树、检查点仍 pending 的 GOAL-006 写成「阶段关门完成」。

### 结论 + 建议给编排器/用户的下一步

JWT 轮换合同、缺省单密钥、双方言轮换后恢复、显式双密钥双路径的**产品证据**可核对，越界未发生，I-001～I-004 门禁真实关闭，I-005 处理合规。独立审计 **不同意无条件 Root 关门**，因为 R5 交付目标尚未进入 canonical 树且五件套/检查点/索引未与已发生事实对齐（F-001）。

**independent verdict = conditional。**

建议 `/govern` 下一句：

> 响应 Root A-001（self pass）+ A-002（independent conditional）：F-001 required — 补 GOAL-006 `01-decision.md`、登记 E-001、把检查点 1–2 改为 done、goal-tree 增 GOAL-006 并同步 workspace R5 行。F-002～F-005 recommended 一并勘误（I-004 镜像、D-001 PG 版本措辞、Q2 引用、VP 信息表回流）。台账对齐后将 GOAL-006 与 Root 标 done 5/5；不必重做实现或重跑四项载体。无 P-004 冲突。

### 声明

本意见不修改 status/progress；响应由 `/govern` 处理。独立审计员未改代码、未改 goal-tree、未写入 `03-audit`（按本轮「不要写入任何文件」执行）。
