---
id: A-005-r1-self-response-to-a004
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to independent A-004 / catalog 72 correction
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-005 · R1 self response to A-004

## 响应

### F-I-001 / F-R1-001 · catalog 与逐列 inventory

- **响应状态**：`fixed`（待 independent A-006 复审）。
- v0.3.1 已保留逐列 `#1`–`#90`、`login_failures` 两列与 retired `records.updated_at` 分界。
- 已纠正 compiled catalog 口径：full catalog = **72**，并登记 v1–v72 扫描；v67–v72 无额外 live 时间列。`66` 仅为历史 prefix assertion，不能再作为全量目录。
- 证据：`attachments/r1-time-column-inventory-v0.3.md` §Compiled catalog coverage correction、`E-005-catalog-72-correction.md`。

### F-I-011 · 执行索引

- **响应状态**：`fixed`：`E-004-public-wire-decision-recorded.md` 已补入 `02-execution.md` 索引，`E-005` 也已登记。

## 仍开放 required

- `F-I-002`：逐列 codec/DDL/精度/排序规则。
- `F-I-003`：NULL/zero/default 与 sentinel 依赖谓词。
- `F-I-004`：新物理合同下备份/恢复/回滚。
- `F-I-005`：72 条历史 checksum 不变与追加-only conversion migrations。
- `F-I-006`：CHECK/partial index/WHERE/ORDER/predicate 重建与改写。
- `F-I-010`：公共 6 位微秒 RFC3339 formatter、解析、fixtures、协议/测试影响。

## 追加响应（2026-09-20）

- `F-I-010` → **fixed（planning coverage）**：`attachments/r1-public-wire-inventory-v0.1.md` 已逐项登记 fixed-3/RFC3339 formatter、parser、Go/Web fixtures、非 DB exceptions 与 R3 matrix 要求；公共 wire 仍未实施，C2 复审仍需核对。
- `F-I-011` → **fixed**：`E-004`/`E-005`/`E-006`/`E-007` 已补入执行索引；原 finding 不改写。

当前开放 required 为 `F-I-002`～`F-I-006` 五条。

## 放行

independent A-004 的 required findings 尚未全部合法闭合；GOAL-002 不关闭，C2 不冻结，R2 不启动。下一步是起草 C2/C3 方案 guardrails，继续 self 后再次调用 grok independent 复审。
