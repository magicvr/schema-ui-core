---
id: GOAL-004-evidence-closure-and-closeout
doc: audit-entry
record_id: A-002
status: recorded
parent: GOAL-004-evidence-closure-and-closeout
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-002 · R4 证据闭环与根目标关门独立交叉审计（2026-09-02）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：close-out / execution-facts
- **scope**：`[workspace-029-wallet-prepaid-instrument/GOAL-004-evidence-closure-and-closeout]` 与 Root `[workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument]` —— VP-029 七条方向级退出判据、红线越界核账、前序 required findings 闭合证据、R4 自审主张
- **verdict**：**conditional**
- **完整意见**：本条由编排器自本地 grok build（grok-4.6 · reasoning high · `/audit`）独立审计会话原样誊入，`source: independent` 保持不变

### 成果（有证据）

1. **工作区绑定合格**：Charter `@0.4.0`；VP-029 `active` v0.2.0；lead 即本区；R1–R3 子目标 GOAL-002/003 `done` 4/4；I-029 系列信息项全部 closed；资料目录 `none`。
2. **判据 #1 主体接缝（模块 API）**：`(issuer, external_id) -> subject_id` 幂等 get-or-create；空输入拒绝；未登记开户拒绝；持久化由 compiled catalog 不过滤 Profile，表结构不依赖钱包 HTTP 启用。
3. **判据 #2 / #3 凭证与核销**：24 字符 Crockford 类字母表（120 bit 熵）+ SHA-256 hex；库列仅 `code_hash`/`code_prefix`；`Redeem` 单事务 CAS 更新 + 同事务调账入金（`adjust` + `ref_type=voucher` + `idempotency_key=voucher.id`）；无公开核销 HTTP；文件 SQLite 20 并发防双花实测 1 成功 19 拦截。
4. **判据 #4 账本不变式**：复用 `adjust` 账本原语，Apply 表未扩新类型，三余额恒等式完全保持，不重开 VP-011。
5. **判据 #5 HTTP / 权限 / 审计切片**：四条路由交付；生成/作废细粒度鉴权 `wallet.voucher.issue`；审计详情无明文泄露；泛化 INTERNAL 文案无底层驱动诊断泄露；权限隔离测试通过。
6. **判据 #6 红线（越界核账）**：Charter 零 diff；默认模块集未改动；`go.mod` 无 Telegram/支付依赖；未重开 VP-011。

### 对照成功标准（VP-029 七条）

| # | 判据 | 判定 | 说明 |
|---|------|------|------|
| 1 | 主体接缝可用 | **pass** | 幂等/未登记不开户/不建 admin.users 有代码与测试 |
| 2 | 凭证生命周期 | **pass** | 生成/作废/过期拒绝有测试；明文列不存在 |
| 3 | 核销原子且幂等 | **pass** | CAS + 幂等键设计；并发测试已纠正 |
| 4 | 账本不变式 | **pass** | 复用 adjust；三余额恒等保持 |
| 5 | Admin 可操作 | **conditional** | HTTP+权限+审计成立；协议页需具备可操作生成/导出与导航入口 |
| 6 | 边界保持 | **pass** | Charter / 默认模块集 / 无支付或 Telegram 依赖 / 未重开 VP-011 |
| 7 | 审计闭合 | **conditional** | 待本意见 F-001 required 闭合后归零 |

### Findings

#### F-001 · required · med · 关联 VP-029 判据 #5
**协议驱动页面需具备可操作的批次生成与导航入口**：在 `wallet-vouchers.json` 中配置 toolbar、生成弹窗表单 `openGenerate`（绑定 `POST /api/wallet/vouchers/batches`），并在 manifest 与 provider 注册侧栏导航 `menu_wallet_vouchers`。

#### F-002 · recommended · low
测试规范性强化：`TestNoPlaintextInDatabase` 采用真实 raw SQL 扫描，`TestVoucherSchemaRegistration` 验证完整 Schema 结构与 Action 绑定。

#### F-003 · recommended · med
自审台账事实对照规范化。

#### F-004 · recommended · low
PG 方言下的实证边界声明保持客观。

#### F-005 · recommended · low
Root 台账与 VP 信息表闭合状态一致性同步。
