---
doc_type: goal-audit
id: A-003-self-response-a002
parent: GOAL-003-r2-offer-purchase-wallet
date: 2026-09-05
status: closed
version: 1.0.0
---

# A-003 · 响应 A-002（self · response）

## A-003 · 响应 A-002 实施独立审计（2026-09-05）

- **source**：self（编排器响应记录，非独立审）
- **模式**：response · 响应 A-002（independent · codex gpt-5.6-sol · verdict **fail** · 2 high required + 1 med recommended + 1 low recommended）
- **verdict**：**pass**（作为响应记录：4 项 finding 全部按 `fixed` 处理；放行以 A-004 independent closure 复审为准）

### 前提承认

A-002 对 A-001 的两处纠正**均成立并接受**：(1) purchase Service API 限流是 D-002 §8/§11 的 R2 明文义务，A-001 将其延期到 R3 属误判；(2) A-001 所称「双数据库验证」超出实际覆盖（PG 当时仅同 request 场景）。self 审计口径已按「能否证明」而非「是否有条款」校准。

### 关闭证据表

| Finding | 级别 | 闭合 | 证据 |
|---------|------|------|------|
| A-002 F-001 · purchase Service API 请求计数桶未接入 | high required | fixed | `service.NewService` 增 `kernel.RateLimiterProvider` 参数并创建 §8 冻结桶（`bizoffer\|purchase\|<subject_id>`、1 min/10、`AllowRecord`、永不 Clear）；`Purchase` 入口逐请求计数，拒绝返回 `ErrRateLimited`（映射冻结通用码 RATE_LIMITED/429）。composition 传 `rateLimiters`（生产面强制）。测试 `TestPurchaseRateLimited`：预算内 10 次成功、第 11 次 429 且域零残留、subject 隔离独立桶、滑窗恢复后再次放行、全程无 Clear。**未修订合同延期义务，而是按合同实现** |
| A-002 F-002 · §4.4 双数据库并发与 retry-exhaustion 零残留验收未完成 | high required | fixed | 新增 `SetPurchaseFaultHookForTest` 注入 seam（freeze/deduct/purchase/entitlement 四步）；验收矩阵 `runPurchaseMatrix` 统一覆盖：retry exhaustion（恰 3 次 attempt、无流水/无凭证/无权益/余额不动）、终态错误不重试（注入 wallet ErrInsufficient → 恰 1 次调用）、瞬时失败后成功（无重复扣款/凭证/权益）、admin 审计 fail-closed、多字节/超长 Q。矩阵在 **SQLite 与真 PostgreSQL** 双库执行（`TestPurchaseMatrixSQLite` + `TestPurchasePostgresAcceptance` 全子测试），PG 另含 4 路并发同 request 收敛 |
| A-002 F-003 · searchQ 按字节截断破坏 UTF-8 | med recommended | fixed | `searchQ` 改为按 **rune** 截断（100 runes）；测试覆盖 60 个 CJK 字符（180 字节）、emoji、101 ASCII、多字节匹配 |
| A-002 F-004 · 审计索引与版本投影未同步 | low recommended | fixed | `03-audit.md` 补登记 A-001 行；workspace.md / goal-tree.md / GOAL-003 00-meta 的 D-002 版本投影统一为 v1.2.0（见下） |

### 实施期新发现（按合同修订纪律处理）

PG 并发验收暴露真缺陷：并发同请求各调用预生成不同 purchaseID、共享幂等键，PG READ COMMITTED 下败者以「同键不同 RefID」触达 ledger 触发 `ErrIdempotencyConflict` 终态，破坏 §4.4 收敛。按「先修合同再实施」处理：**D-002 v1.2.0**（GOAL-003 D-002 附录）——购买事务内凭证 INSERT 前置为请求守卫，败者在守卫处失败并回读重放，永不触碰钱包。SQLite 写者串行未暴露该窗口的事实同时记录。

### 仍开放项

- 无（A-002 全部 finding 已处理；GOAL-003 保持 active，待 A-004 independent closure 复审通过后关门）。

### 声明

本响应记录不修改 A-001/A-002 原文；`fixed` 证据以现行代码、测试与 D-002 v1.2.0 文本及 git 历史可核对。
