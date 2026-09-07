---
doc_type: goal-decision
id: D-002-r1-contract-freeze
parent: GOAL-001-wallet-prepaid-instrument
date: 2026-09-02
status: accepted
version: 0.1.0
---

# D-002 · R1 核心合同冻结与技术选型裁决

## 触发

承接 `GOAL-001` 纲领路线图 R1 阶段及 VP-029 信息门禁，在开展 R2/R3 实施前，对主体落点、账本入金类型、凭证权限键、核销 API 范围及哈希并发合同（I-029-001/002/003/005/006）进行方案选型与技术裁决。各选型经技术调研并由用户正式书面确认裁决。

## 决定

| 信息项 | 裁决决定 | 理由与技术规格 |
|--------|----------|----------------|
| **I-029-001 主体落点** | **独立通道无关主体表 `subjects`** | 独立维护 `(issuer, external_id) -> subject_id` 映射。提供 `GetOrCreateSubject` / `SubjectExists` 基础能力；查询与登记不依赖 `admin.wallet` 启用。钱包 `wallet_accounts` 的 `owner_type` 扩展支持 `subject`（迁移 0064），`OwnerExistsFunc` 适配主体存在性校验（禁止回退至 `admin.users`）。满足 VP-029 判据 1。 |
| **I-029-002 入金类型** | **复用 `adjust` + `ref_type='voucher'`** | 零侵入复用现有 `adjust` 账本原语，通过 `ref_type='voucher'` 和 `ref_id=voucher.id` 追踪凭证来源。严格遵守“不重开 VP-011”，不破坏三余额恒等、流水快照链与对账 Job。满足 VP-029 判据 3/4。 |
| **I-029-003 权限键** | **新增细粒度权限键 `wallet.voucher.issue`** | 遵循最小特权原则，将管理员批量生成/导出/作废预付凭证的权限与超管直接调账扣款（`wallet.adjust`）清晰隔离。满足 VP-029 判据 5。 |
| **I-029-005 核销 API** | **仅交付 Go 模块内部 API `Redeem(ctx, subjectID, code)`** | 首个消费者为同进程 Telegram Bot 模块（VP-030），在 Bot 会话鉴权后直接在进程内安全调用；不在公共 HTTP 暴露未登录核销端点，避免消耗 RT-Q05 精神的复杂防暴力破解限流。 |
| **I-029-006 哈希与并发** | **高熵码 + SHA-256 + 单事务 CAS 原子核销** | 采用 16~24 字符高熵随机字符（>80bit 熵），数据库存储 `code_hash`（SHA-256）与 `code_prefix`（前缀供运营检索/识别）；`UNIQUE(code_hash)`；核销时同事务 CAS `UPDATE vouchers SET status='redeemed', redeemed_by=?, redeemed_at=? WHERE id=? AND status='unused'`，影响行数=0 则 fail-closed，并发抢占失败；同事务内以 `voucher.id` 为幂等键调用钱包入账，确保原子且绝对不双记。满足 VP-029 判据 2/3。 |

## 未选方案与取舍

- **未选主体收敛于 `admin.wallet`**：若绑定在钱包模块内部，则违反判据 1“查询与 get-or-create 不依赖 `admin.wallet` 已启”，且无法优雅服务后续不依赖钱包的主体场景。
- **未选新增 `entry_type='voucher_redeem'`**：侵入已有流水 CHECK 约束与 Apply 表，增加了重开 VP-011 的风险；`adjust + ref_type` 机制在工程与审计上已经完备。
- **未选直接复用 `wallet.adjust` 权限**：发卡具有向市场发行资产凭证的性质，若与调账混用则权限粒度过粗，不符合金融安全规范。
- **未选暴露匿名公共 HTTP 核销端点**：匿名核销需要复杂的 IP 限流、验证码和防爆破机制，与当前同进程 Admin/Bot 架构不匹配。
- **未选 bcrypt / 6位低熵数字码**：bcrypt 计算缓慢且无法建立 UNIQUE 索引查重；低熵码易被撞库，高熵 SHA-256 是最佳工程实践。
