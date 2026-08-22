---
id: GOAL-003-dual-key-jwt
doc: audit
status: done
parent: GOAL-001-key-rotation-and-backup
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
---

# 审计 · GOAL-003

> 本文件是稳定索引。正式意见写在 `03-audit/A-NNN-*.md`。

## 信息就绪核对（按 scope）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| 影响本 scope 的 I-00N | I-003 已 verified（D-001） | 无到期开放 required；I-005 non-blocking 保持 collecting（默认「previous 可验」已按 VP 冻结措辞实施） |

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required | 文件 |
|------|------|--------|-------|---------|---------------|------|
| A-001 | 2026-08-22 | self | stage close-out · 目标整体（R2 双密钥） | pass | 0 | [A-001-r2-dual-key-stage-self.md](03-audit/A-001-r2-dual-key-stage-self.md) |
| A-002 | 2026-08-22 | independent | close-out · 目标整体（D-001 语义/verifyAccess/接线/证据/越界） | pass | 0 | [A-002-r2-dual-key-closeout-independent.md](03-audit/A-002-r2-dual-key-closeout-independent.md) |

## 编排器响应记录（2026-08-22 · /govern）

- **意见合并**：A-001（self pass）+ A-002（independent pass）无 verdict 冲突、无对同一必改项的一要一否；P-004 冲突裁决点不触发。开放 required finding = 0。
- **recommended findings 响应**（详见 E-003）：
  - F-001 composition 接线缺钉死 → **fixed**：新增 `TestNewAuthenticatorWiresPreviousSecret`（生产构造路径双密钥通过 / 空 previous 单密钥拒绝），`-count=1` PASS。
  - F-002 E-002 索引漏登 → **fixed**：索引已补 E-002/E-003。
  - F-003 整包 exit 0 不可作为已复现事实 → **fixed（收窄）**：E-002 v1.1 改引「vet 0 + JWT 相关包 ok（双方独立复现）」；store PG 两条集成失败确认为共享 probe DB 残留（`postgresWasFresh` 对任意用户表敏感），非本切片回归；R3 注记使用专用 DB。
- **关门判定**：required = 0（三路径闭合齐备）；I-003 门禁 verified；方向级 R2 对应标准逐条可核对。GOAL-003 `done`（4/4）。

## 结论状态

GOAL-003 关门：self + independent 双 pass，recommended 全部 fixed，0 开放 required。
