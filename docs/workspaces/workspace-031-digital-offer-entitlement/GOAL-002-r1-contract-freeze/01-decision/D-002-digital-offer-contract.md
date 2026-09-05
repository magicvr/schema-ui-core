---
doc_type: goal-decision
id: D-002-digital-offer-contract
parent: GOAL-002-r1-contract-freeze
date: 2026-09-05
status: draft
version: 0.1.0
---

# D-002 · 数字 Offer 业务域合同 v0.1.0

> R1 合同正文。R2/R3/R4 的实施、验收与关门以本文为分母；偏离本文需先以 02 决策修订合同再实施。
> 框架性裁决（权益形态 / 购买状态机 / 人工发放边界 / 模块 id / 命令清单）见 D-001，本文不得与之冲突。
> 落盘状态：`draft`（C2 待 C3 审计闭合后置 `accepted` 并升 1.0.0）。

## 0. 对齐与判据映射

| VP-031 判据 | 合同条款 |
|-------------|----------|
| 1 · Offer CRUD + Admin 协议页面 + 权限键 + 审计 + C 端上架列表 | §2、§5、§6、§7 |
| 2 · 余额不足拒绝、freeze→deduct_frozen、失败 fail-closed、凭证与权益同事务、并发测试 | §4、§9 |
| 3 · 权益有效/过期/耗尽可测，服务提供前统一核验 | §3、§5 |
| 4 · 只挂 VP-029 `subject_id`，不创建 `admin.users` | §4.2、§5 |
| 5 · 可选 Telegram 命令；未启用不依赖 Bot API | §6 |
| 6 · 激活门禁留痕（freshness / RT-Q03/Q05） | Root meta（已由 VRev-080 留痕） |
| 7 · 边界保持 | §10 |
| V-F119 · 业务限流桶冻结，请求计数禁 key-wide Clear | §8 |

## 1. 模块与装配

- 模块 id：**`biz.digital-offer`**（D-001）。
- 代码位置：`apps/api/modules/digitaloffer/`（Go 包名 `digitaloffer`），布局对齐 wallet 先例：`provider.go`（ModuleID、New、路由/权限/页面贡献）、`store/`、`migration/`、`schema/`、`manifest/`。
- 装配：`composition.go` 以 `plan.HasModule("biz.digital-offer")` 门控，**不进** `mvp`/`admin` 默认 Profile 集。
- 模块输入依赖（构造注入）：`*walletstore.Repository`（钱包资金原语）、`*subject.Store`（VP-029 主体）、`operationlog.Recorder`（审计）、`kernel.RateLimiterProvider`（限流）、可选 `kernel.TelegramDispatcher`（通道缝）。
- 对 `admin.wallet` / `channel.telegram` 的依赖均为**消费依赖**：本模块不修改其行为，不要求其 HTTP 面启用（直接消费 store/service 层，先例：wallet subject 消费）。

## 2. Offer 模型（表 `digital_offers`）

| 列 | 类型 | 约束 / 说明 |
|----|------|-------------|
| id | TEXT PK | 模块内 `newID(now)`（对齐 wallet 时间有序 id 先例） |
| name | TEXT NOT NULL | 展示名 |
| description | TEXT NOT NULL DEFAULT '' | 展示描述 |
| price_amount | INTEGER NOT NULL CHECK (> 0) | 标价，**最小货币单位** |
| currency | TEXT NOT NULL | 币种代码（与钱包账户币种对齐才可购） |
| entitlement_form | TEXT NOT NULL CHECK IN ('duration','count') | 权益形态，**创建后不可改**（防存量权益语义漂移） |
| duration_seconds | INTEGER NULL | duration 型必填且 ≥ 1；count 型必须 NULL |
| count_per_purchase | INTEGER NULL | count 型必填且 ≥ 1；duration 型必须 NULL |
| status | TEXT NOT NULL DEFAULT 'draft' CHECK IN ('draft','on_sale','off_sale') | 上架状态 |
| version | INTEGER NOT NULL DEFAULT 0 | 乐观锁（对齐 wallet_accounts） |
| created_at / updated_at | TIMESTAMP NOT NULL | UTC |

- CHECK 约束同时锁定形态互斥：`(entitlement_form='duration' AND duration_seconds>=1 AND count_per_purchase IS NULL) OR (entitlement_form='count' AND count_per_purchase>=1 AND duration_seconds IS NULL)`。
- 状态机：`draft → on_sale → off_sale`；`off_sale → on_sale` 允许重新上架；**无删除**（保留历史凭证/权益引用）。仅 `on_sale` 可购。
- 可变字段：`name` / `description` / `price_amount` / `status`（均走 version 乐观锁 + 审计）；`currency`、`entitlement_form`、`duration_seconds`、`count_per_purchase` 创建后**不可变**（改价不动形态；价格快照在购买时定格）。
- 时长语义选择 `duration_seconds`（而非天数）：粒度最细、可表达试用/会员/年卡，向后兼容无需迁移。

## 3. 权益模型（表 `digital_entitlements`）

| 列 | 类型 | 约束 / 说明 |
|----|------|-------------|
| id | TEXT PK | 模块 newID |
| subject_id | TEXT NOT NULL | VP-029 `subjects.id`；**唯一主体外键，禁止关联 `admin.users`** |
| offer_id | TEXT NOT NULL | 来源 offer |
| purchase_id | TEXT NOT NULL | 来源购买凭证（本波全部权益来自购买，无人工发放） |
| form | TEXT NOT NULL CHECK IN ('duration','count') | 购买时 offer 形态快照 |
| expires_at | TIMESTAMP NULL | duration 型 = `purchase.created_at + offer.duration_seconds`；count 型 NULL |
| remaining_count | INTEGER NULL | count 型 = `offer.count_per_purchase`；duration 型 NULL |
| status | TEXT NOT NULL DEFAULT 'active' CHECK IN ('active','voided') | 存储状态仅两值 |
| created_at / updated_at | TIMESTAMP NOT NULL | updated_at 记录作废时间 |

- **有效判定（惰性派生，无后台 job）**：duration 型有效 = `status='active' AND expires_at > now`；count 型有效 = `status='active' AND remaining_count > 0`。`expired` / `exhausted` 是派生谓词，不落库变更。
- 判据 3 三态映射：有效（valid）/ 过期（expired）/ 耗尽（exhausted）均可经 §5 核验 API 测得；另含 `voided`（作废）与 `no_entitlement`（从未持有）。
- 叠加语义：同 subject 多份有效权益**并存**（每份一行）；duration 型任一行在窗内即有效（取最晚 `expires_at` 展示）；count 型消耗规则见 §5.2。
- 作废：`voided` 终态，本波**不可逆**（纠错重发可另行购买或后续波次引入人工发放）。

## 4. 购买模型与事务边界（表 `digital_purchases`）

### 4.1 凭证列

| 列 | 类型 | 说明 |
|----|------|------|
| id | TEXT PK | 模块 newID |
| subject_id | TEXT NOT NULL | 购买主体（VP-029） |
| offer_id | TEXT NOT NULL | 购买的 offer |
| offer_name | TEXT NOT NULL | 名称快照 |
| amount | INTEGER NOT NULL | 成交金额快照（最小货币单位） |
| currency | TEXT NOT NULL | 币种快照 |
| freeze_entry_id | TEXT NOT NULL | 钱包 `freeze` 流水引用 |
| deduct_entry_id | TEXT NOT NULL | 钱包 `deduct_frozen` 流水引用 |
| request_id | TEXT NOT NULL | 客户端幂等键；`UNIQUE(subject_id, request_id)` |
| status | TEXT NOT NULL DEFAULT 'fulfilled' CHECK (status='fulfilled') | **本波唯一终态**（D-001：无 pending） |
| created_at | TIMESTAMP NOT NULL | UTC |

- 凭证**只追加、不可更新/删除**；失败购买不落凭证（钱包流水由钱包自身记录）。
- `UNIQUE(subject_id, request_id)`：幂等重放返回既有凭证（见 §4.3）。

### 4.2 单事务购买（判据 2 的「同事务」实现）

```text
store.Run(ctx, func(tx kernel.Tx) error {
  1. 读 offer（tx 内）：status 必须 = 'on_sale'，否则 ErrOfferNotOnSale
  2. 校验 subject 存在（subjects 表，tx 内）；不存在 → ErrSubjectNotFound
  3. 幂等：按 (subject_id, request_id) 查既有凭证；命中 → 原样返回（成功幂等重放）
  4. 币种校验：subject 钱包账户 currency == offer.currency，否则 ErrCurrencyMismatch
  5. walletRepo.MutateInTx(tx, freeze)        // available → frozen（余额不足 → ErrInsufficient 透出）
  6. walletRepo.MutateInTx(tx, deduct_frozen) // frozen → deducted
  7. INSERT digital_purchases (status='fulfilled')
  8. INSERT digital_entitlements（duration → expires_at；count → remaining_count）
})
```

- 任一步失败 → **整体回滚**：无购买凭证、无权益、无冻结残留。冻结从未对外可见，等价于「失败 unfreeze」的更强 fail-closed 形态（VP-031 判据 2 的「同事务」支路）。
- **本波不存在「freeze 已提交但购买失败」状态**，故购买路径无 unfreeze 补偿分支；`unfreeze` 保留为 wallet Admin 资金纠错入口（已有）。若未来引入分段提交（如异步履约），必须修订本合同：补 unfreeze 补偿 + 孤儿冻结对账 + pending 状态机。
- 并发正确性依赖钱包事务内余额校验与 `ErrVersionConflict`；offer 状态以事务内读取为准，`version` 乐观锁兜底管理面并发改价/下架。

### 4.3 入口与鉴权

| 入口 | 主体来源 | 说明 |
|------|----------|------|
| 模块 Service API `Purchase(ctx, subjectID, offerID, requestID, ...)` | 显式 subject_id | 测试面与未来通道封装；VP-031「无 Telegram 时 HTTP/模块 API 仍可测」 |
| Telegram 命令 `buy`（§6） | `TelegramUpdate.SubjectID`（VP-030 身份映射已填充，webhook.go） | 真实用户路径；未启用通道时零依赖 |
| C 端 HTTP 购买端点 | — | **本波不设**：尚无 subject 会话鉴权设施，发明 token 方案属范围扩张；future seam，不阻塞判据 |

- `GET /api/biz/offers`（§5.3）为公开只读上架列表（判据 1「C 端可列出上架项」）。

## 5. 权益 API（模块 Service 面，判据 3）

### 5.1 核验（服务提供前统一走此 API）

`Check(subjectID, offerID, now) → {valid bool, reason}`，reason ∈ `none | no_entitlement | expired | exhausted | voided`。

- duration 型：存在任一 `active` 且 `expires_at > now` 的行 → valid。
- count 型：存在任一 `active` 且 `remaining_count > 0` 的行 → valid。
- 全部行均 `voided` → `voided`；有行但全部过期 → `expired`；count 型全部 `remaining_count=0` → `exhausted`；无行 → `no_entitlement`。

### 5.2 消耗（count 型）

`Consume(ctx, subjectID, n, now) → error`：

- 单事务内按 `created_at ASC` 逐行原子扣减：`UPDATE digital_entitlements SET remaining_count = remaining_count - ? WHERE id = ? AND remaining_count >= ?`，以 RowsAffected 确认；总可用 < n → `ErrEntitlementInsufficient`，整体回滚。
- duration 型不参与消耗（窗内不限次）。
- 并发消耗不得超扣（ RowsAffected 守卫 + 事务串行化）。

### 5.3 C 端列表（HTTP）

- `GET /api/biz/offers`：无鉴权，仅返回 `on_sale` 项；字段：`id, name, description, price_amount, currency, entitlement_form, duration_seconds | count_per_purchase`。不含 `version` 等内部字段。限流见 §8。

## 6. 可选 Telegram Register（判据 5，I-031-004）

- 命令（经 `kernel.TelegramDispatcher.RegisterCommand` 注册；`channel.telegram` 未启用时注入 `DisabledDispatcher` no-op，模块零感知，测试不依赖 Bot API）：

| 命令 | 行为 |
|------|------|
| `price` | 列出全部 `on_sale` offer：名称、价格（最小货币单位→展示串）、权益形态（时长/次数）与 id |
| `buy` | `buy <offer_id>`：同步一拍购买（§4.2），成功回复凭证摘要与权益详情；失败回复原因（余额不足/已下架/形态不存在/币种不符），**无确认回调流**（本波最小） |
| `entitlements` | 列出本人有效权益：形态、剩余次数或到期时间；作废/过期不列 |

- 命令处理内 subject 为空（身份映射缺失）→ 明确报错提示，不得退化为匿名购买。
- 回调按钮（`RegisterCallback`）本波不冻结、不实现。

## 7. Admin 面（判据 1，I-031-003）

- 权限键（`kernel.PermissionContribution`，Resource = `digitaloffer`，PolicyID = admin，对齐 wallet 先例）：

| 权限键 | 动作 |
|--------|------|
| `digitaloffer.read` | Admin 读取 offers / purchases / entitlements 列表与详情 |
| `digitaloffer.offer.manage` | 创建/编辑/上架/下架 offer（全部写操作，version 乐观锁） |
| `digitaloffer.entitlement.void` | 作废权益（终态，写审计） |

- 路由（对齐 wallet 路由风格，`a.Middleware` + 权限键门控 + `operationlog` 审计）：

| 方法 | 路径 | 权限 |
|------|------|------|
| GET | `/api/digitaloffer/offers` | digitaloffer.read |
| POST | `/api/digitaloffer/offers` | digitaloffer.offer.manage |
| PATCH | `/api/digitaloffer/offers/{id}` | digitaloffer.offer.manage |
| GET | `/api/digitaloffer/purchases` | digitaloffer.read |
| GET | `/api/digitaloffer/entitlements` | digitaloffer.read |
| POST | `/api/digitaloffer/entitlements/{id}/void` | digitaloffer.entitlement.void |

- **无人工发放端点**（D-001）；作废走 `POST .../void`（幂等：已作废再作废返回成功）。
- Admin UI：schema 驱动页面（`schema/digitaloffer-offers.json`、`schema/digitaloffer-entitlements.json` + 导航/manifest fragment），本波**无自定义 React 组件**（判据 1「协议页面」由协议渲染满足；先例：wallet 页面）。
- 购买列表为只读审计视图（Admin 不代购、不改凭证）。

## 8. 限流桶（V-F119 · 请求计数语义）

经 `kernel.RateLimiterProvider` 工厂创建（常量声明于 handler，对齐 `wallet_self` 先例）：

| 桶 key 形状 | 覆盖 | 窗口 / 阈值（首波冻结值） |
|-------------|------|---------------------------|
| `bizoffer|purchase|<subject_id>` | Telegram `buy` 与 Service API 购买调用 | 1 分钟 / 10 次 |
| `bizoffer|price|<subject_id>` | Telegram `price` / `entitlements` 查询 | 1 分钟 / 30 次 |
| `bizoffer|list|<ip>` | C 端 `GET /api/biz/offers` | 1 分钟 / 60 次 |

- 语义：**请求计数桶，非失败预算桶**——每次请求 `AllowRecord`；返回 false → HTTP 429 + `Retry-After`（对齐 W12 语义）或 Telegram 文本提示。
- **永不调用 key-wide `Clear`**（V-F119 硬约束；请求计数无「成功后清零」语义）。`Reserve/Cancel` 不用于本模块。
- 阈值为合同冻结值；调整需修订本合同（02 决策），不得实现期静默改。

## 9. 错误码（errorcatalog 注册，对齐 `WALLET_*` 格式）

| 码 | i18n key | 语义 |
|----|----------|------|
| `BIZOFFER_NOT_FOUND` | error.bizOfferNotFound | offer 不存在 |
| `BIZOFFER_NOT_ON_SALE` | error.bizOfferNotOnSale | 非 on_sale 不可购 |
| `BIZOFFER_CURRENCY_MISMATCH` | error.bizOfferCurrencyMismatch | offer 币种与主体钱包账户不符 |
| `BIZOFFER_INSUFFICIENT_FUNDS` | error.bizOfferInsufficientFunds | 余额不足（透出钱包 ErrInsufficient） |
| `BIZOFFER_SUBJECT_NOT_FOUND` | error.bizOfferSubjectNotFound | subject 不存在 |
| `BIZOFFER_ENTITLEMENT_INVALID` | error.bizOfferEntitlementInvalid | 权益无效（expired/exhausted/voided/no_entitlement） |
| `BIZOFFER_REQUEST_CONFLICT` | error.bizOfferRequestConflict | 同 request_id 但负载不同（对齐 wallet 幂等冲突语义） |
| `BIZOFFER_FORM_CONFLICT` | error.bizOfferFormConflict | offer 形态不可变字段被修改 |

- 域 sentinel（`store` 包）由 handler 映射到上述冻结码，先例：`writeWalletError`。

## 10. 迁移、Profile 与红线

- 迁移：`migration/migration.go` 建 `digital_offers` / `digital_purchases` / `digital_entitlements` + 索引（purchases：`UNIQUE(subject_id, request_id)`、`idx(subject_id)`、`idx(offer_id)`；entitlements：`idx(subject_id, status)`、`idx(purchase_id)`；offers：`idx(status)`）。占位符一律 `?`（dialect 中性，pg + sqlite 双跑）。
- Profile：不进 `mvp`/`admin` 默认集；装配仅经 `plan.HasModule("biz.digital-offer")`。
- 红线（判据 7）：不做类目树/SKU/税/库存/物流订单/购物车；不做营销/订阅计费/支付网关/退款编排；不把 Telegram Bot 运行时打进本模块；不解禁通用 Entitlement/Approval Gate 接缝；不改 Charter；购买与权益只挂 `subject_id`，不创建 `admin.users`。
- 事件：模块内同步调用；不要求 typed domain event 接缝解禁。

## 11. R2/R3 实施切片（预告，不改变判据）

- R2：§1、§2、§4、§7（offer 部分与购买链路）、§8（purchase 桶）、§9、§10 迁移 + 判据 2 并发测试（同 subject 并发购买余额恰够一次 → 恰一成功；幂等双发 → 一凭证一扣款一份权益；并发消耗不超扣）。
- R3：§3 核验/消耗、§5、§6 Telegram 命令、§7 权益作废 UI、§8 查询桶。
- R4：证据矩阵（判据逐条 → 测试/代码/审计路径）、边界核账、关门审计。

## 12. 未选方案（合同级）

- **分段提交（freeze 独立提交 + unfreeze 补偿）**：引入孤儿冻结窗口与对账义务；单事务回滚是更强 fail-closed（§4.2）。
- **expired/exhausted 落库状态 + 后台 job**：增加生命周期任务与迁移复杂度；惰性派生可测且无 job（§3）。
- **C 端 HTTP 购买端点**：需先发明 subject 会话鉴权，超本波范围（§4.3）。
- **权益全局库存/限量（stock）**：属库存域，红线排除；count 型限额是 per-purchase 授予次数，非全局库存。
- **价格历史版本表**：审计日志 + 购买快照已覆盖首波需求。
