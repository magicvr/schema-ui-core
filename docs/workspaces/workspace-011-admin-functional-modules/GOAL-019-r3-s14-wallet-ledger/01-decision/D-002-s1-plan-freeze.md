---
id: D-002
goal: GOAL-019-r3-s14-wallet-ledger
title: 方案冻结：钱包/账务（账本：余额、流水、对账）设计（S1）
date: 2026-08-16
status: accepted
parent: GOAL-019-r3-s14-wallet-ledger
created: 2026-08-16
updated: 2026-08-16
version: 1.1.0
---

# D-002 · 方案冻结（S-14 钱包/账务 · admin.wallet）

> 依据：I-011-001 §4 S-14（余额、流水、对账；余额变动审计 + 迁移基建）、§7 协议对照口径；GOAL-019 00-meta 边界与 I-001~I-004；A-002 019-F-001（关联单据 vs S-13）与 019-F-002（权限键拆分）。
> 证据：modules/datapermission/（provider.go / manifest/fragment.json / schema/data-permission.json / migration/migration.go ——模块五段式先例）、modules/recyclebin/migration（领域表 DDL 先例）、modules/operationlog/repository.go（EventXxx 常量）、internal/store/store.go WithTx（平台事务边界）、kernel/profile.go（ProfileAdmin）、kernel/provider.go（DefaultNavigationOrder）、internal/composition/composition.go（组合根接线）、handler/error_contract_test.go（错误码契约）、docs/vision/protocol-inventory-v2.7.0.md §2.5。

## 1. 账务领域模型（I-001 闭合 · 019-F-001 裁定）

- **账户实体** `wallet_accounts`：id（服务端生成 ULID）｜owner_type ∈ {user, business, system} ｜ owner_id ｜ currency（v1 默认 `CNY`，单币种，枚举留扩展位）｜ 三余额口径 **总额 balance_total = 可用 balance_available + 冻结 balance_frozen**（同表 CHECK 恒等式 + 对账链校验双保障，杜绝浮点）｜ status ∈ {active, disabled} ｜ version 乐观锁。UNIQUE(owner_type, owner_id, currency)。
- **精度与并发**：所有金额以**整数最小单位（分）** INTEGER 存储，禁止浮点；展示层格式化（web i18n/format 先例）。并发控制 = **version 乐观锁**（UPDATE … WHERE version = ?）+ `store.WithTx` 事务内完成「校验 → 更新余额 → 插流水」；冲突返回 409 `LEDGER_VERSION_CONFLICT`。幂等 = 可选 `idempotency_key`（**复合唯一 (account_id, idempotency_key)**；同账户同 key 同载荷 → 返回既有流水，异载荷 → LEDGER_IDEMPOTENCY_CONFLICT；金额敏感写操作防重放。权威裁定见下方幂等小节与 §4 DDL）。
- **流水实体** `wallet_ledger_entries`（**不可变账本**）：id（ULID）｜ account_id ｜ entry_type ∈ {adjust, freeze, unfreeze} ｜ amount_delta（带符号整数，**语义按 entry_type 见下方 apply 表（A-004 F-001 勘误）**）｜ **变动后三余额快照** balance_after_total / _available / _frozen（链式对账证据）｜ ref_type/ref_id **可选空引用**（019-F-001 裁定，见下）｜ idempotency_key ｜ memo（必填）｜ actor_id/actor_name ｜ created_at。**无 UPDATE/DELETE 路径**。

### 金额变动原语（A-004 F-001 勘误 · 可执行 apply 表）

| entry_type | amount_delta 含义 | 符号 | total | available | frozen | 拒绝条件 |
|------------|-------------------|------|-------|-----------|--------|----------|
| adjust | 总额变动 | ≠ 0（可正可负） | += d | += d | 不变 | 任一余额 < 0 或恒等式破坏；账户 disabled |
| freeze | available→frozen 转移额度 | > 0 | 不变 | -= |d| | += |d| | available < |d|；账户 disabled |
| unfreeze | frozen→available 转移额度 | > 0 | 不变 | += |d| | -= |d| | frozen < |d|；账户 disabled（F-004 裁定：冻结资金随停用锁定） |

- **快照链重放规则**（对账/测试可执行）：对每条流水，按上表 `apply(prev_after, entry) == entry.balance_after_*`；首笔 prev = (0,0,0)；**末笔快照 == 账户当前三余额**；每笔快照自身满足恒等式 total = available + frozen。**链序 = (created_at ASC, id ASC)**（F-003 勘误：秒级时间戳同秒并列，id 为并列键）。
- **幂等（A-004 F-002 勘误）**：UNIQUE 范围 = **(account_id, idempotency_key)** 复合；同账户 + 同 key + **同载荷**（entry_type/amountDelta/memo 一致）→ 返回既有流水（幂等成功）；同账户 + 同 key + 异载荷 → `LEDGER_IDEMPOTENCY_CONFLICT`；查找必须带 account_id，**禁止按裸 key 取他户流水**。
- **019-F-001 裁定（关联单据）**：v1 `ref_type/ref_id` 为**可选空引用**（NULL 合法），不做存在性校验；对账仅**账本内部勾稽**（快照链 + 恒等式）。S-13 订单立项后若需单据级联/存在性校验，作为扩展登记（触发 = S-13 立项），本波不实现。
- **账户生命周期**：管理面创建（初始零余额）/ 启停（disabled 时拒绝调账/冻结/**解冻**——冻结资金随停用锁定，流水只读）；**不实现转账（transfer）与外部支付/结算**（未选方案）。

## 2. 余额变动审计与迁移基建（I-002 闭合）

- **双层审计**：(a) **业务账本**——每笔变动 = 不可变 ledger entry（带三余额快照 + actor + memo），账本本身即审计；负向调账即冲正（不另设 reverse 动作）；(b) **operationlog 事件**——沿用 operationlog.Recorder 先例（EventMFAEnroll 同款），新增 `wallet.account-create` / `wallet.account-update` / `wallet.adjust` / `wallet.freeze` / `wallet.unfreeze` / `wallet.reconcile`。
- **残余（文档化）**：operationlog 事件与账本变更**非同一事务**（handler 层独立记录，与 recyclebin/MFA 先例一致）；账本 entry 同事务写入 + 快照链校验兜底，审计强度不依赖事件事务性。
- **迁移基建**：**0031**（admin.wallet 三表 + CHECK/UNIQUE/索引 + compiled persistence 注册，沿用 datapermission/migration + compiled 注册先例）；**0032**（core.operationlog CHECK 事件超集重建，沿用 0028/0030 先例）；迁移冻结校验测试（GOAL-004 A-003 冻结校验先例）。

## 3. 端点、权限键与 Profile 归属（I-004 闭合 · 019-F-002 响应）

| 端点 | 门禁 | 说明 |
|------|------|------|
| GET /api/wallet/accounts | wallet.read | 账户列表（分页/搜索） |
| POST /api/wallet/accounts | wallet.write | 创建账户（owner_type/owner_id/currency 唯一）；审计 wallet.account-create |
| PATCH /api/wallet/accounts/{id} | wallet.write | 启停（status active/disabled）；审计 wallet.account-update |
| GET /api/wallet/accounts/{id}/entries | wallet.read | 流水列表（分页，created_at desc） |
| POST /api/wallet/accounts/{id}/adjust | **wallet.adjust** | 调账（带符号金额 + memo 必填 + 可选 idempotencyKey）；审计 wallet.adjust |
| POST /api/wallet/accounts/{id}/freeze | **wallet.adjust** | 冻结（available→frozen）；审计 wallet.freeze |
| POST /api/wallet/accounts/{id}/unfreeze | **wallet.adjust** | 解冻（frozen→available）；审计 wallet.unfreeze |
| POST /api/wallet/reconcile | wallet.read | 对账校验（恒等式 + 快照链），写 reconciliation_runs 留痕；审计 wallet.reconcile |
| GET /api/wallet/reconcile/runs | wallet.read | 对账运行记录 |

- **权限键（019-F-002 响应：写路径权限键从 I-004 拆出并冻结为 required 设计）**：`wallet.read`（查看）/ `wallet.write`（账户生命周期）/ `wallet.adjust`（**金额变动专用键**：调账/冻结/解冻）。金额敏感写操作不与生命周期写操作共用权限键。
- **模块命名确认**：ModuleID `admin.wallet`；页面 `wallet`；菜单 `menu_wallet`（fragment visibleWhen features.menu_wallet）；路由 /wallet。
- **Profile 归属**：admin.wallet 进入 **admin 默认集**（profile.go ProfileAdmin 追加 + 注释；S 系列内容扩展先例 file-library/data-dictionary）；mvp/demo 不含；DefaultNavigationOrder 尾部追加 menu_wallet（GOAL-013 产品权威序 = 仅 manifest，快照测试同步更新）。
- Descriptor：DependsOn core.auth-session / core.navigation-capability / core.schema-render / core.operationlog（datapermission 同款）；Contributions 声明 Routes/Pages/Navigation/Permissions/Fragments。

## 4. 迁移 DDL（0031 · admin.wallet）

```sql
CREATE TABLE wallet_accounts (
  id                TEXT PRIMARY KEY,
  owner_type        TEXT NOT NULL CHECK (owner_type IN ('user','business','system')),
  owner_id          TEXT NOT NULL,
  currency          TEXT NOT NULL DEFAULT 'CNY',
  balance_total     INTEGER NOT NULL DEFAULT 0 CHECK (balance_total >= 0),
  balance_available INTEGER NOT NULL DEFAULT 0 CHECK (balance_available >= 0),
  balance_frozen    INTEGER NOT NULL DEFAULT 0 CHECK (balance_frozen >= 0),
  status            TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','disabled')),
  version           INTEGER NOT NULL DEFAULT 0,
  created_at        INTEGER NOT NULL,
  updated_at        INTEGER NOT NULL,
  UNIQUE (owner_type, owner_id, currency),
  CHECK (balance_total = balance_available + balance_frozen)
);
CREATE TABLE wallet_ledger_entries (
  id                     TEXT PRIMARY KEY,
  account_id             TEXT NOT NULL,
  entry_type             TEXT NOT NULL CHECK (entry_type IN ('adjust','freeze','unfreeze')),
  amount_delta           INTEGER NOT NULL CHECK (amount_delta != 0),
  balance_after_total    INTEGER NOT NULL CHECK (balance_after_total >= 0),
  balance_after_available INTEGER NOT NULL CHECK (balance_after_available >= 0),
  balance_after_frozen   INTEGER NOT NULL CHECK (balance_after_frozen >= 0),
  CHECK (balance_after_total = balance_after_available + balance_after_frozen),
  ref_type               TEXT,
  ref_id                 TEXT,
  idempotency_key        TEXT,
  UNIQUE (account_id, idempotency_key),
  memo                   TEXT NOT NULL,
  actor_id               TEXT NOT NULL,
  actor_name             TEXT NOT NULL,
  created_at             INTEGER NOT NULL
);
CREATE INDEX idx_wallet_ledger_account ON wallet_ledger_entries(account_id, created_at DESC);
CREATE TABLE wallet_reconciliation_runs (
  id             TEXT PRIMARY KEY,
  account_id     TEXT,  -- NULL = 全局
  result         TEXT NOT NULL CHECK (result IN ('consistent','inconsistent')),
  mismatch_count INTEGER NOT NULL DEFAULT 0,
  details        TEXT NOT NULL DEFAULT '{}',
  actor_id       TEXT NOT NULL,
  created_at     INTEGER NOT NULL
);
```

- 0032（core.operationlog）：事件集合超集重建，加入 `wallet.account-create` / `wallet.account-update` / `wallet.adjust` / `wallet.freeze` / `wallet.unfreeze` / `wallet.reconcile`（0028/0030 同款 CHECK 重建）。
- 无 FK（沿用既有模块无外键先例，应用层校验 account_id 存在）。

## 5. 协议对照（I-003 闭合 · 独立口径）

- **独立对照**（I-011-001 §7 口径，不沿用准入波次 9/0 外推）：protocol-inventory-v2.7.0 §2.5 信息性场景——Manifest 点名样例 `order-list-batch` / `order-detail-lifecycle` 为上游 `_samples/` 清单，**无钱包/账本专属协议面**；D-DATA（列表/详情 API 形状）、D-ACT（actions）、D-PERM（权限继承/intent）能力面已由基架覆盖并有范例（search-form-table / row-backend-actions / permission-inheritance）。
- 页面 requiredCapabilities（app.manifest / app.navigation / permissions.inheritance / actions.row.request / actions.page.trigger / table.sort / form.controls.extended / form.controls.advanced）与 data-permission 页同类，全部已有覆盖先例。
- **处置**：S-14 = **本地领域模块**（admin.wallet），不新增协议 capability、不改协议 pin（v2.8.0）、不改 Manifest 装配语义；**「呈现自由 + fail-open」处置留痕**——页面呈现（表格/表单/统计卡）属呈现自由，协议无强制形态；调账/冻结/对账为本地后端动作，不声明协议覆盖。wallet 数据**不接入 data-transfer 导出面**（v1 不登记导出资源；后续接入按 §7 口径评估并留痕）。

## 6. 测试与验证

- store：WithTx 原子性（余额更新 + 流水插入同事务）、version 乐观锁冲突、idempotency_key 幂等、负向调账 available 不足拒绝、disabled 账户拒绝写。
- 对账：恒等式 + 快照链重放校验（apply 表驱动，链序 (created_at ASC, id ASC)；连续条目 balance_after 衔接；末笔 == 账户三余额）、不一致清单与 mismatch_count、reconciliation_runs 留痕。
- handler：门禁 401/403（read/write/adjust 分键）、审计事件、错误码（WALLET_NOT_FOUND / INVALID_LEDGER_ENTRY / INSUFFICIENT_BALANCE / LEDGER_VERSION_CONFLICT / INVALID_WALLET_BODY / WALLET_DISABLED / LEDGER_IDEMPOTENCY_CONFLICT 候选，S2 按 error_contract 收敛）。
- 组合根：admin 权限键 **27→30**（+wallet.read/write/adjust；实测基线 27，A-004 F-005 勘误）、导航 **13→14**（composition_test 快照，实施按 live snapshot 断言）；迁移 **30→32**。
- web：i18n（manifest.title.wallet / manifest.nav.wallet / schema.wallet.*）+ fixture（schema-keys / s5-denominator / e2e 双 profile）回归 + wallet 页冒烟。

## 7. 未选方案（留痕）

- **转账（transfer）与外部支付/结算**：超出 I-011-001 §4 S-14 最小面（余额/流水/对账）与 Charter 边界；后续按业务需要立项。
- **冲正（reverse）独立动作**：负向调账已覆盖，避免第二写路径。
- **单据级联/存在性校验（F-001 备选）**：S-13 未立项，依赖不成立；可选空引用 + 触发登记替代。
- **对账引入外部流水/银行流水**：v1 仅账本内部勾稽；外部对账源后续按需立项。
- **多币种/汇率**：currency 枚举留扩展位，v1 单币种。
- **角色级/数据权限集成**：wallet 账户不登记 data-permission 生产资源（GOAL-016 v1 无登记资源）；如需行级范围，后续目标按 ScopeAware 契约登记。

## 8. S2 实现清单

1. migration 0031/0032 + 迁移 provider 注册（compiled persistence 先例）。
2. modules/wallet：store（accounts/ledger/reconcile 仓库，WithTx + 乐观锁 + 幂等）+ provider（Descriptor/路由/权限/fragment）+ schema（wallet 页 JSON）。
3. handler：WalletRoutes（九端点；错误码契约；六类审计事件）。
4. 组合根接线（composition.go）+ profile.go（ProfileAdmin += admin.wallet）+ DefaultNavigationOrder（+= menu_wallet）+ 快照测试更新。
5. web：i18n + fixture + wallet 页回归。
6. 测试：store 原子性/并发/幂等、对账链校验、handler 门禁与错误码、组合根快照（权限 **27→30**、导航 13→14）、迁移 30→32 冻结校验、e2e 双 profile。
7. A-002 recommended 落地：F-001 可选空引用已裁定（§1）；F-002 权限键拆分已冻结（§3）。
8. 组合根快照按 **§6 基线（权限 27→30、导航 13→14、迁移 30→32）** 断言，禁止采用旧 26→29（A-005 F-005 残留勘误）。