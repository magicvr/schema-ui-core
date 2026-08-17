---
id: GOAL-019-w14-rectification-batch-b
doc: decision
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · GOAL-019 S1 方案冻结（F-05～F-07）

## 决策

### F-05 · 列表端点校验与分页

- recycle-bin / wallet accounts / wallet entries / reconcile runs / per-task runs 全部校验 `page`、`pageSize`（1..100），非法返回 `INVALID_PAGE` / `INVALID_PAGE_SIZE`。
- data-permission policies 改为真实分页（内存切片分页），不再伪造 pageSize。
- recycle-bin 暴露 `sort`/`order`（deletedAt/resource/actorName）。

### F-06 · 错误码与目录

- 新增 `OPERATION_NOT_FOUND` 目录条目（messageKey + en/zh）。
- 将复用误导的错误码细分：
  - `INVALID_SCOPE_BODY`：policy/scopes 请求体与必填字段错误。
  - `INVALID_WALLET_OWNER`：ownerId 缺失。
  - `INVALID_WALLET_ACCOUNT`：accountId 缺失。
  - `INVALID_WALLET_STATUS`：status 非法。
- 更新 error_contract_test 冻结集合与扫描（含 NotFoundCode）。

### F-07 · 搜索/排序/过滤一致性

- 通知 q 搜索大小写不敏感（repository lower(q)）。
- wallet 账户 q 搜索扩展到 owner_id/owner_type/currency。
- wallet ledger 支持 `entryType` 过滤与 `q` 搜索（memo/ref_type/ref_id），schema 增加 entryType 筛选。
- recycle-bin 支持 sort/order。
- 不新增业务域模块；不改 Profile 默认集。

## 信息项更新

| ID | 状态 | 说明 |
|----|------|------|
| I-001 | **closed** | data-permission policies 采用真实分页（当前数据量小，内存分页足够，不引入 SQL pushdown） |
| I-002 | **closed** | wallet ledger entry-type 取值集合 = adjust/freeze/unfreeze/deduct_frozen |
