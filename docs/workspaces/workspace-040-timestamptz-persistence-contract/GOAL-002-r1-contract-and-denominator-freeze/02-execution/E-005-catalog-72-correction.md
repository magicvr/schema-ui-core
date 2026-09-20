---
id: E-005-catalog-72-correction
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-005 · compiled catalog v72 口径纠正

## 已发生事实

- independent A-004 复核确认：现行 compiled catalog full length = **72**；现有 `len(applied) != 66` 是 `applied[:66]` 之后的历史 prefix assertion，不是全量目录长度。
- v0.3 inventory 已补充 catalog v1–v72 coverage statement，并核对 v67–v72 无额外 live 时间列；v0.3.1 attachment 已落盘。
- `login_failures.locked_until` / `login_failures.updated_at` 已在逐列清单 #20/#21，90 live columns 机械计数保持。

## 证据

- `attachments/r1-time-column-inventory-v0.3.md` §Compiled catalog coverage correction
- `03-audit/A-004-r1-independent-reaudit-after-a003.md` F-I-001

## 门禁

本事实修正仅为 inventory/catalog 口径；C2/C3 与 F-I-002～F-I-006/F-I-010 仍开放，不放行 R2。
