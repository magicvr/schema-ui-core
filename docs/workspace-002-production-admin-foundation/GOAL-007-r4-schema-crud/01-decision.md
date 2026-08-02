---
title: 决策 · R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.4.0
---

# 决策 · GOAL-007

## D-001 · 用一个端到端目标实施 Root D-010

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 以单一 R4 子目标承载 records 精确契约、SQLite 持久化、CRUD API、Schema 读写交互、权限/错误状态及重启回归；六个成功标准按依赖顺序推进。
  2. 延续 Root D-010：代表实体固定为 `records`；生产默认迁入 SQLite；重启保持为 required 验收；错误响应继续使用 HTTP status + 稳定 `code` + message 的统一 envelope。
  3. 立项只确定路线与门禁，不在缺少证据时枚举新的 error `code`、DDL、并发策略或 Schema action 映射。
  4. `I-007-001`、`I-007-002`、`I-007-003`、`I-007-004` 均为 required；每项必须在表中所列首个受影响实施或验收动作前由证据关闭并记录后续决策。
- **理由**：API、持久化、Schema action 与重启证据共同构成一个可验证的业务生命周期；拆成多个并列目标会形成无法独立验收的中间态。把未知项显式登记为 required，可在保留端到端交付边界的同时防止方案被代码隐式冻结。
- **实施门禁**：立项时四项均为 `open`。D-002/D-003 已关闭 `I-007-001`/`I-007-002` 并完成 S1/S2 契约冻结；`I-007-003`/`I-007-004` 仍 open，分别阻断首个 Schema 写交互代码与 S6 验收。

### 未选方案

- **按 API / DB / Web / 测试拆成四个并列目标**：依赖紧密且成功边界不可独立成立，会增加跨目标门禁与中间态。
- **沿用进程内 records 并只补 Schema 页面**：无法满足 D-010 的 SQLite 与重启保持 required 边界。
- **立项时先猜精确 error code / DDL / action 形状**：会把尚未收集和验证的信息伪装为决定，违反 P-005。

## D-002 · 冻结 records 精确 API 与错误契约（S1）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 对外实体保持五字段 `id`/`name`/`status`/`owner`/`updatedAt`；不新增 `createdAt`。`id` 与 `updatedAt` 仅服务端管理；可编辑字段仅 `name`/`status`/`owner`（trim 后非空）；`status` 不做枚举白名单。
  2. 继承既有 list/detail/PATCH/DELETE 路径、查询参数、list envelope、权限键与错误 envelope `{"error","message"}`；稳定 code 全表见附件，含已有 `UNAUTHENTICATED`/`FORBIDDEN`/`INVALID_SORT_*`/`INVALID_PAGE*`/`RECORD_NOT_FOUND`/`INVALID_PATCH_*`。
  3. 新增 `POST /api/records`（`records.write`）：body 必填 `name`/`status`/`owner`；成功 **201** + 完整 record；`id` = `rec-` + 16 位小写 hex（`crypto/rand`）；失败 code 冻结为 `INVALID_CREATE_BODY` / `INVALID_CREATE_FIELD`（400），稀有内部失败 `INTERNAL`（500）。禁止把 create 错误复用 `INVALID_PATCH_*`，也不引入 R4 的 `RECORD_CONFLICT`/409。
  4. DELETE 保持 **204** 空体；并发语义为 last-write-wins（无乐观锁/version）；PATCH 忽略未知 JSON 键；body 上限 4 KiB。
  5. 正反矩阵 T-API-01～13 为 S3 实施与回归的最低 API 断言集。
- **理由**：对照现 handler/测试与 Root I-004 M-R4-01～06/08 可完整继承读改删与权限基线；唯一结构性缺口是 create，其 code 与 PATCH 分离可避免前端/测试歧义。不在契约层发明唯一约束或枚举，以免超出当前产品语义。
- **信息门禁**：`I-007-001` → `verified`；证据 [I-007-001-api-error-contract.md](attachments/I-007-001-api-error-contract.md)。本决策完成 **S1** 契约冻结，并放行 S3 中受 API/错误契约约束的代码变更；不构成 S3 已实现，也不关闭 `I-007-003`/`I-007-004`。

> **修订（2026-08-02 · 响应 A-001 F-001）**：`updatedAt` 时间语义精度由 RFC3339 秒级统一为**含毫秒**（固定 3 位小数）；「严格晚于」断言保留。详见 D-004 与 I-007-001 v0.2.0。

### 未选方案

- **create 复用 `INVALID_PATCH_BODY/FIELD`**：混淆操作语义，矩阵与前端映射更难维护。
- **引入 status 枚举或 name 唯一 / 409 CONFLICT**：当前实现与演示数据无此约束，扩大范围且无用户要求。
- **客户端指定 id 或乐观锁**：增加冲突面；与现 PATCH 形状不一致。
- **只收集不冻结**：无法勾选 S1，S3 仍会被 `I-007-001` 阻断。

## D-003 · 冻结 records SQLite DDL、0003 迁移、seed 与 repository（S2）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 在既有迁移链追加 **`0003` / `records_persist` / transformID `0003:records-persist:v1`**；只创建 `records` 表与 `name`/`updated_at`/`owner` 索引，**不**改写 0001/0002 的 SQL 或 checksum 输入。
  2. DDL：`id TEXT PK`，`name`/`status`/`owner` TEXT NOT NULL（trim 非空 CHECK），`updated_at INTEGER` Unix 毫秒；无 FK、无 name 唯一、无 soft-delete。API RFC3339 含毫秒 ↔ DB Unix 毫秒在 repository 映射。
  3. 迁移 up 只建空表；业务种子走 **`seedRecords`**：在 `seedAdmin=true` 路径于 `seedRBAC` 之后执行；**仅当表行数为 0** 时插入与现 `staticRecords()` 对齐的 8 行（`rec-1`…`rec-8`）；非空则整段跳过，避免撤销用户删除或覆盖变更。
  4. 生产默认唯一数据源为 SQLite repository；废除进程内切片作为生产路径。写并发 = SQLite 单写者 + last-write-wins。Open 签名保持不变。
  5. 非空文件库在应用 pending（含 0003）前必须有一致性快照 + `integrity_check`；checksum 漂移与迁移事务失败 fail closed。T-DB-01～09 为 S3/S6 最低持久化断言。
- **理由**：复用 R3 runner/ledger/seed 模式可将 records 持久化纳入同一启动路径与恢复口径；空表才 seed 才能同时满足「新库有演示数据」与「删除/更新重启保持」。
- **信息门禁**：`I-007-002` → `verified`；证据 [I-007-002-sqlite-migration-plan.md](attachments/I-007-002-sqlite-migration-plan.md)。本决策完成 **S2** 结构冻结，并放行 S3 持久化代码变更；不构成 repository 已实现或 S6 重启证据，也不关闭 `I-007-003`/`I-007-004`。

> **修订（2026-08-02 · 响应 A-001 F-001）**：`updated_at` 由 Unix 秒改为 Unix **毫秒**（INTEGER，UTC）；DDL 列语义、seed 数据与 repository 映射见 I-007-002 v0.2.0 与 D-004。

### 未选方案

- **迁移内插入 8 条种子**：把演示数据绑死在 checksum 上，后续改种子即漂移；与 R3「迁移 vs seed」分层不一致。
- **按 id 永续 ensure 种子行**：用户 DELETE 后重启会插回，破坏删除持久化与 D-010 重启保持。
- **保留 handler 进程切片作为生产回落**：无法证明 SQLite 默认路径，且双源状态不一致。
- **改写 0001/0002 塞入 records**：破坏已部署库 checksum 与 R3 证据链。

## D-004 · 统一 `updatedAt` 精度与断言（响应 A-001 F-001）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. `updated_at` 存储精度由 Unix **秒** 提升为 Unix **毫秒**（INTEGER，UTC）；API `updatedAt` 返回 RFC3339 **含毫秒**（固定 3 位小数，格式 `2006-01-02T15:04:05.000Z07:00`，如 `2026-08-02T03:04:05.123Z`）。
  2. 保留「每次成功 update 后 `updatedAt` **严格晚于**更新前值」的断言（I-007-001 §3.1 与 T-API-05 语义不变）；create 与每次成功 update 仍写入 `time.Now().UTC()`。
  3. 严格递增的实现确定性：写入时若该行新时间戳 ≤ 其前一 `updated_at`，则钳制为 `prev + 1ms`（**单调钳制**，仅在同一毫秒内的快速连续更新时触发）；**禁止**整秒/整毫秒跳变，也不退回 Unix 秒。该钳制使读回时间戳仍为本次变更的单调时间，测试无需 sleep 即可稳定断言严格递增。
  4. seed 数据 `updated_at` 改为与现 `staticRecords()` 对齐的 Unix 毫秒（`2026-07-31T00:00:00Z` 起每条 +11h）。
  5. 修订范围：D-002（API 时间语义）、D-003（DDL 列精度）与两份 I-007 附件同步更新；`idx_records_updated_at` 索引与 `updatedAt` 排序语义不变。
- **理由**：响应 A-001 F-001（independent）。`updatedAt` 严格递增是既有 API 契约（`records_test.go` `.After()`）的既定语义，秒级存储使其在同一秒内不可满足；统一到毫秒级使「每次 update 写入 `time.Now().UTC()`」与「严格晚于」同时成立，并保留客户端可检测的同秒变更。生产级基架下秒级 `updatedAt` 过粗。
- **信息门禁**：`I-007-001`/`I-007-002` 附件更新为 v0.2.0 并继续 `verified`；A-001 F-001 按 `fixed` 闭合（见 03-audit 响应节与 A-002 self 复核）；**S3 实施放行**。

### 未选方案

- **保留秒级 + 断言放宽为非递减**：弱化既有「严格晚于」契约，客户端无法区分同一秒内变更；与 VP-002 生产级语义不符。
- **人为加整秒/整毫秒跳变**：脱离 `time.Now().UTC()`，A-001 F-001 明确指出不可取。
- **微秒精度**：与毫秒在本语义等价，但毫秒更贴近常见前端展示粒度且格式确定；微秒可后续按需升级，不在本决策冻结。

## D-005 · 冻结 Schema CRUD 读写交互契约（I-007-003）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. **代表页面**：`list-edit-lifecycle` 作为唯一代表性 CRUD 生命周期页（承接 D-010 的 records 代表实体与 R3 种子菜单 `list-edit-lifecycle`）。其 fixture 演进为「列表 table（含行操作 `edit`/`delete` 与工具栏 `create`）+ 新建/编辑 form（`submitAction`）+ 详情 recordView + 搜索绑定」；行为全部由 fixture 数据驱动，**不修改 Renderer 主路径代码**（S4 成功标准）。
  2. **Node/action 绑定**：`table.props.actions`（rowAction `edit`/`delete`，`permissionIntent` edit/delete）→ PATCH / DELETE；`table.props.toolbar`（toolbarTrigger `create`，`permissionIntent` edit）→ 打开新建 form → POST；`form` `submitAction`（formSubmit；create 模式 POST `/api/records`，edit 模式 PATCH `/api/records/{id}`）；搜索 form → table `q` 绑定（把 `search-form-table` 现行「form-to-query 绑定 out of scope」纳入 S4）。渲染层**一次性**补齐「table actions/toolbar 渲染 + form submit 绑定 + 成功/错误/确认反馈」；此后新页面仅改 fixture，不改 Renderer 主路径。
  3. **字段映射**：`id`/`updatedAt` 仅服务端、只读展示；`name`/`status`/`owner` 可编辑 string。控件：`name`→input、`owner`→input、`status`→select（options active/pending/archived **仅 UI 提示**，API 不做枚举白名单，与 I-007-001 一致）。create body `{ name, status, owner }` 全必填；PATCH body 仅 present 键；wire kind 全 string。
  4. **交互状态**：加载（DataTable `loading` + 提交中禁用提交）；空态（"No records match."）；成功（create/edit/delete 后刷新列表 + 页级/行级成功提示）；错误（统一 envelope `{error,message}` → `role=alert` 页级提示，稳定 code 映射见附件 §4）；删除确认（复用冻结 `executeAction` `confirm=true` → 未确认即 `CONFIRM_CANCELLED`）。
  5. **权限矩阵**：后端为权威（`records.read` 门禁 GET，`records.write` 门禁 POST/PATCH/DELETE；匿名 401 / 缺权限 403）；**前端隐藏不是安全边界**（S5：后端 403 不被前端隐藏替代）。admin 全量 CRUD；editor/viewer 只读（写 affordance 隐藏/禁用，直接调用仍 403）；匿名 → LoginPage。表达式门禁用冻结语法：`$context.user.permissions contains "records.write"`（写）、`"records.read"`（读）、`$context.features.menu_list_edit_lifecycle`（菜单）。使用 permission 字段的页面 `meta.requiredCapabilities` 须含 `permissions.inheritance`；复用冻结 executeAction 序列（visible → permission → disabled → confirm）与 target kinds（rowAction/toolbarTrigger/formSubmit/actionButton）。
  6. **测试矩阵**：T-UI-01～10（附件 §6）为 S4/S5 验收最低断言；API 负向已由 T-API-08/09 承担，UI 只负责正确呈现。
- **理由**：前端已具备冻结的 Renderer node whitelist、$context 表达式/反应引擎与 `permissions.inheritance` 执行引擎（rowAction/toolbarTrigger/formSubmit/actionButton 目标已建模），缺的是「页面文档把 records 写路径绑定到这些原语」以及「一次性渲染补齐」。把代表页面固定为 `list-edit-lifecycle` 可让 fixture 演进与既有种子菜单/导航投影对齐，避免新增页面破坏已冻结的导航与菜单证据。
- **信息门禁**：`I-007-003` → `verified`；证据 [I-007-003-schema-crud-interaction.md](attachments/I-007-003-schema-crud-interaction.md)。本决策完成 S4/S5 交互契约冻结并**放行首个 Schema 写交互代码变更**；**不构成 S4/S5 已实现**，不关闭 `I-007-004`，不勾选 Root R4。

### 未选方案

- **拆分为 list / create / edit 三个页面**：增加导航与跨页状态传递复杂度，与单一代表菜单（`list-edit-lifecycle`）不对齐；单页生命周期更贴近 D-010「一个代表实体完整闭环」。
- **在 Renderer 主路径硬编码 records 特定逻辑**（如组件内写死 POST/PATCH/DELETE 调用）：违反 S4「新增/调整代表页面不修改 Renderer 主路径代码」；改为 fixture 驱动的通用 action 绑定。
- **用 status 枚举白名单约束 API**：违反 I-007-001 冻结（status 非枚举）；UI select 选项仅作提示，非服务端约束。
