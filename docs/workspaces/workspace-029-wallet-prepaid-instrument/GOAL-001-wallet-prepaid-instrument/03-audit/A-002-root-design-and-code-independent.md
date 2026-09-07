---
id: GOAL-001-wallet-prepaid-instrument
doc: audit-entry
record_id: A-002
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-002 · Root 方案设计与代码实现独立交叉审计（2026-09-02）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：close-out / design-plan / execution-facts
- **scope**：`[workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument]` 根目标完成情况——对照 VP-029 七条退出判据，**独立审查方案设计与代码实现**；不以五件套 status、goal-tree、前序 A 条目闭合声明或执行台账勾选作为关门证据
- **verdict**：**conditional**
- **工作区**：`workspace-029-wallet-prepaid-instrument` · root `GOAL-001-wallet-prepaid-instrument` · `shared_materials_catalog: none`（无共享资料引用）

### 范围与区间

- **审什么**：主体接缝、凭证生命周期、核销原子/幂等、账本不变式、Admin 协议面（生成/导出/作废）、红线边界。
- **不审什么**：不改 `status`/`progress`/goal-tree；不把治理文档「已 done」当作实现完成。
- **本会话实证**：
  - 通读实现：`modules/wallet/subject`、`voucher`、`store`、`migration`（0064）、`provider.go`、`schema/wallet-vouchers.json`、`handler/wallet.go` 凭证路由、`composition.go` 装配、`kernel/profile.go` 默认集、`apps/web/src/renderer/render.tsx` 请求/表单提交路径。
  - 复跑测试（exit 0）：`go test ./modules/wallet/... ./internal/handler`；`go test ./internal/store -run TestCompiledMigrationCatalogOwnership|TestMigrateFreshDB`；`go test ./internal/handler -run TestVoucher`。
  - 本会话无 `PG_TEST_*`，PostgreSQL 运行时核销未实测。

### 成果（有证据）

1. **主体接缝（判据 #1，模块 API）**  
   `subjects` 表 `UNIQUE(issuer, external_id)`（0064）。`GetOrCreateSubject` 空输入拒绝、冲突后重读、15 并发文件库测试恰好 1 条新建。`CreateAccount(owner_type=subject)` 走 `SubjectExists`，未登记 → `ErrNotFound`（`provider_test.go`）。未见写入 `admin.users`。持久化由 compiled catalog 贡献，不随 `HasModule("admin.wallet")` 过滤（`compiled/persistence.go`）。无公开主体 HTTP，符合「模块 API」合同。

2. **凭证生命周期（判据 #2，存储与 API 切片）**  
   24 字符 Crockford 类字母表（120 bit）+ SHA-256 hex；库列仅 `code_hash`/`code_prefix`；`UNIQUE(code_hash)`。生成 HTTP 201 返回一次性 `code`；列表/详情手工组包不含 `code`/`codeHash`。作废 CAS：`UPDATE ... AND status='unused'`；已核销 → `ErrVoucherAlreadyRedeemed`；过期 Redeem → `ErrVoucherExpired`。审计详情断言不含明文。无公开 HTTP Redeem。

3. **核销原子且幂等（判据 #3，SQLite 路径）**  
   `Redeem` 单事务：按 hash 读取 → CAS `unused→redeemed`（`RowsAffected==0` → `ErrVoucherConflict`）→ 同事务 `GetOrCreateSubjectAccountInTx` → `MutateInTx(entry_type=adjust, ref_type=voucher, idempotency_key=voucher.id)`。重复核销 → `ErrVoucherAlreadyRedeemed`，余额不双记。文件 SQLite（WAL + pool）20 并发：1 成功 / 19 冲突，流水 1 条。本会话复跑 `voucher` 包测试通过。

4. **账本不变式（判据 #4，类型面）**  
   未新增 `entry_type`；`Apply` 仍走 `adjust`。`TestRedeemSuccess` 断言三余额。既有 `modules/wallet/store` 测试本会话通过。未引入 Telegram / 支付 SDK（`apps/api/go.mod` 无匹配）。

5. **Admin HTTP / 权限 / 审计（判据 #5 的 API 切片）**  
   四条路由挂在 `admin.wallet`：`POST /batches`、`GET /`、`GET /{id}`、`POST /{id}/void`。生成/作废 `wallet.voucher.issue`；列表/详情 `wallet.read`。仅 `wallet.adjust` 的角色 403。错误码已进 `errorcatalog` 与 `error_contract_test.go` 冻结集。协议页有 toolbar `openGenerate`、行内作废、侧栏 `menu_wallet_vouchers`。

6. **边界（判据 #6）**  
   `profileDefaults` 未向 `mvp`/`admin` 塞入新模块 id；`admin.wallet` 原已在 admin 默认集，本波是该模块增量。`OwnerUser` 自动开户 HTTP 仍走 `UserByID`，不把 C 端主体做成可登录用户。无支付网关依赖。

### 对照成功标准（VP-029 七条；以代码与运行时为准）

| # | 判据 | 判定 | 说明 |
|---|------|------|------|
| 1 | 主体接缝可用 | **pass（模块 API）** | 幂等/未登记不开户/不建 admin.users 有测试；查询不依赖钱包 HTTP 启用 |
| 2 | 凭证生命周期 | **conditional** | 生成/作废/过期/哈希存储有测试；「导出」在协议页不可用（F-001） |
| 3 | 核销原子且幂等 | **conditional** | SQLite CAS+防双花有测试；PG 同事务 UNIQUE 重读沿用已文档化的失败模式（F-003） |
| 4 | 账本不变式 | **conditional** | 复用 adjust、三余额测试通过；凭证 `currency` 不参与入账账户选择（F-002） |
| 5 | Admin 可操作 | **fail 切片** | HTTP/权限/审计/作废成立；协议页生成成功后丢弃一次性明文，导出未交付（F-001） |
| 6 | 边界保持 | **pass** | 无新默认模块、无 Telegram/支付依赖、未改 entry_type apply 表 |
| 7 | 审计闭合 | **fail（本条）** | 本意见开放 required > 0；不得以 A-001/子目标闭合声明抵销 |

### Findings

#### F-001 · 协议页生成成功后丢弃一次性明文，Admin 无法导出卡密
- 严重度：**high**
- 建议：**required**
- 关联：VP-029 判据 #2 / #5（生成/**导出**/作废须有协议驱动页面）
- 状态：**open**
- 描述：预付凭证明文只在 `POST /api/wallet/vouchers/batches` 的 201 响应出现一次，之后不可再取。协议页把该请求接成表单 `submitAction: generateBatch`，`onSuccess.behavior = reload`，按钮文案却是 “Generate & Export”。渲染器成功路径**不把响应体交给页面**：
  1. `runRequest`（`apps/web/src/renderer/render.tsx`）解析 JSON 只为采集 `fieldErrors`，成功返回 `{ ok: true }`，明文数组被丢弃。
  2. `submitForm` 在 `result.ok` 后固定：通用成功 toast → `reloadList()` → `setActiveModal(null)`。模态关闭后表内只有 `codePrefix`。
  3. 仓库已有的导出能力是白名单 `custom` handler（`export.users` 等 blob 下载）。凭证页未走该路径，也没有二次下载端点（二次下载在哈希存储模型下本来就不可行）。
- 证据：
  - `apps/api/modules/wallet/schema/wallet-vouchers.json`：`generateBatch.onSuccess = reload`；`submitLabel = Generate & Export`；无 download/custom handler。
  - `apps/web/src/renderer/render.tsx`：`runRequest` 成功分支丢弃 body；`submitForm` 成功分支关模态并 reload。
  - `apps/api/internal/handler/wallet.go`：201 `items[].code` 为唯一明文出口。
- 为何不能用「API 返回数组即导出」放行：判据写的是**协议驱动页面**可操作导出。管理员按页面操作会永久失去卡密。这是方案与实现的结构缺口，不是文档口径问题。
- 闭合方向（示例，非本条实施）：生成请求必须在**同一用户手势**把 201 的 `code` 落成 CSV/TXT 下载（custom handler 或新的 outcome），禁止以 reload/关模态作为成功副作用；补协议页或 e2e 断言「生成后浏览器获得含明文的文件，且文件未进审计原文」。

#### F-002 · 核销忽略凭证币种，非 CNY 面额会记入默认 CNY 账户
- 严重度：**med**
- 建议：**required**
- 关联：VP-029 判据 #4；生成 API 接受任意 `currency`
- 状态：**open**
- 描述：`GenerateBatch` 把调用方 `currency` 写入 `vouchers.currency`（空则 `CNY`）。`Redeem` → `GetOrCreateSubjectAccountInTx` **写死** `DefaultCurrency`（`CNY`），`MutateInTx` 把凭证 `amount` 记入该 CNY 账户，**不比较** `voucher.currency` 与 `account.currency`。协议生成表单允许编辑 Currency。结果：面额为 USD 的凭证核销后增加 CNY 余额，三余额恒等仍成立但资金语义错误。
- 证据：
  - `modules/wallet/voucher/service.go` `GenerateBatch` 持久化 `v.Currency`；`Redeem` 入账不读该字段做账户选择。
  - `modules/wallet/store/repository.go` `GetOrCreateSubjectAccountInTx`：`OwnerSubject, subjectID, DefaultCurrency`。
  - `handler/wallet.go` 生成接口未校验币种白名单。
  - 无「异币种 fail-closed」测试。
- 闭合方向：异币种 fail-closed（推荐，与现行单币种账本一致），或按凭证币种开户并禁止串币种入账；生成 API/表单与账户选择同一规则；补测试。

#### F-003 · Redeem 同事务内账户 UNIQUE 冲突重读，复用已被否决的 PostgreSQL 模式
- 严重度：**med**
- 建议：**required**
- 关联：VP-029 判据 #3（双方言）；`GetOrCreateUserAccount` 已记录的 W9 F-001
- 状态：**open**
- 描述：`GetOrCreateUserAccount` 明确写着：PostgreSQL 上失败的 INSERT 会 abort **当前事务**，失败者不得在同一 tx 里重读（W9 F-001）。`GetOrCreateSubjectAccountInTx` 却在 **Redeem 的同一 `runner.Run` 内** INSERT，UNIQUE 冲突后立刻 `GetSubjectAccountByOwnerInTx`。SQLite 约束失败通常不 abort 事务，故本会话并发测试通过不能外推到 PG。并发「新主体首次入金」或账户行冲突时，PG 路径会得到含 aborted transaction 的不透明错误；CAS 随事务回滚（fail-closed，不双花），但核销对调用方不可重试语义化，且与用户账户路径的既有修复不一致。本会话无 `PG_TEST_*`，未做 PG Redeem 实证。
- 证据：
  - `store/repository.go` `GetOrCreateUserAccount` 注释 vs `GetOrCreateSubjectAccountInTx` 实现。
  - `voucher/service.go` `Redeem`：CAS 与 `GetOrCreateSubjectAccountInTx` 同 tx。
- 闭合方向：同事务内 `INSERT ... ON CONFLICT DO NOTHING` 再 SELECT，或 SAVEPOINT；错误映射为可重试；补 PG 方言 Redeem（含并发首次开户）。不得仅用 SQLite 测试关闭本条。

#### F-004 · 凭证页 i18n 键未进入前后端文案目录
- 严重度：**med**
- 建议：**recommended**
- 状态：**open**
- 描述：`wallet-vouchers.json` / `manifest/fragment.json` 声明 `schema.walletVouchers.*`、`manifest.title.walletVouchers`、`manifest.nav.walletVouchers`，`apps/web/src/i18n/messages/{zh-CN,en-US}.json` 无对应条目。页面有英文 `label` 兜底；中文环境标题/导航/表单会退回英文或键名。HTTP 错误码在 Go `errorcatalog` 有中英，但 web catalog 未镜像 `error.invalidVoucher*`（客户端未知 key 会退回服务端 message，严重度低于页面键缺失）。
- 证据：对 `apps/web/src/i18n/messages` 检索上述 key 无命中。

#### F-005 · composition 的 `OwnerExistsFunc` 仍只查 `admin.users`
- 严重度：**low**
- 建议：**recommended**
- 状态：**open**
- 描述：主体开户门禁在 `Service.CreateAccount` / `Redeem.SubjectExists`，功能上未回退 `UserByID`。但 `composition.go` 注入的 `OwnerExistsFunc` 仍是 `authRepository.UserByID`，与 R1「OwnerExists 承认已登记主体」的字面合同不一致。当前 HTTP `by-owner` 只开 **user** 账户，故不是资金漏洞；后续若有人把该回调复用到 subject 路径会静默走错表。
- 证据：`apps/api/internal/composition/composition.go` 约 567–572 行；对比 `provider.go` `CreateAccount` 的 `SubjectExists` 分支。

#### F-006 · 生成页缺少过期字段；`batch_id` 无唯一约束
- 严重度：**low**
- 建议：**recommended**
- 状态：**open**
- 描述：模型与 API 支持 `expiresAt`，协议生成表单无该字段，管理员无法在页面设置过期。`batch_id` 非 UNIQUE，重复生成同一批次会混在同一列表过滤里，而明文只存在于各次 201 响应。与 F-001 叠加后运营风险更大。

### 必改项汇总

1. **F-001（required / high）**：协议页必须在生成当次把一次性明文变成可保存的导出（下载），不得 reload/关模态丢弃 201 body。
2. **F-002（required / med）**：核销入账必须尊重凭证币种（异币种 fail-closed 或按币种开户），禁止把非 CNY 面额记入 CNY 账户。
3. **F-003（required / med）**：去掉 Redeem 同事务 UNIQUE 冲突后的 PG-unsafe 重读；补 PG 实证。

### 与既有意见的异同

| 来源 | 异同 |
|------|------|
| 本目标 A-001 | 条文过短，把子目标闭合声明当作 Root 完成。本条按代码/渲染器重审，**不继承**其 pass。 |
| GOAL-003 A-002 F-005（recommended，「导出=API 数组」） | 该口径不能替代 VP-029 判据 #5 的协议页导出。本条 F-001 升级为 required。 |
| GOAL-002 A-002 F-002（PG） | 当时指出 0064/Redeem 缺 PG 证据。当前 `GetOrCreateSubjectAccountInTx` 仍使用 W9 已否决的同事务重读。本条 F-003 针对**现存代码**，不是追究台账是否曾标注 fixed。 |
| GOAL-004 A-002 F-001 | 当时要的是 toolbar/导航，代码里已有。本条 F-001 是下一步：**有入口仍导出失败**。 |

方案层（独立表、高熵 SHA-256、CAS+adjust 同事务、不暴露 HTTP Redeem、细粒度 `wallet.voucher.issue`）总体合理，且 SQLite 资金路径测试通过。不能关门的原因是 Admin 导出结构缺失、币种入账错误、PG 账户开户路径与已知修复不一致。

### 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** 判据 #5 的导出切片未在协议页成立；判据 #4/#3 各有一条资金/方言 required。开放 required > 0，判据 #7 本条不满足。

建议 `/govern`：响应本目标 A-002；先修 F-001 与 F-002，再修 F-003；修复后用代码+可重复测试（含协议页导出或等价 e2e，以及异币种 fail-closed）申请 finding 闭合。不要用更新治理勾选代替上述证据。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
