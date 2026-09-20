---
id: E-004-checkpoint-d-checksums-and-residual-rereview
doc: execution-entry
status: active
parent: GOAL-003-r2-codec-and-descriptor-m1-m2
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-004 · 检查点 D：canonical SQL / 真实 checksum 落盘 + `D-021` residual 复审

## 事实

- **canonical 语句清单落盘**：`attachments/r2-v73-v87-generated-statements-v0.1.md` —— 15 个 descriptor 的 `transform_id`、**真实 `MigrationChecksum`**、表范围、FK 子女盘点（父表/子表/无子表 46 张）、m0（guard + sentinel 预检）/ m1–m3（重建）/ m4（校验）逐条语句文本、PG 显式 DDL 清单。
- **哈希记录与锁定**：`internal/store/migrate_test.go` 冻结目录表逐条断言 15 个 checksum；`TestCompiledMigrationCatalogOwnership` 通过。
- **可复现性**：重跑生成器（`VP040_GENERATE=1`）对同一份 v1–v72 历史产出 **16/16 输出字节级不变**（15 个 Go 文件 + 语句清单）。
- **`D-021` residual 复审（触发已满足）**：independent（grok-build grok-4.6 · high）落盘 `03-audit/A-002-independent-checkpoints-b-c-d.md`，逐项判定 residual ①②③ 均有可核对证据、可 `fixed` 闭合；编排器接受并在 **GOAL-002** 台账留下闭合记录（`GOAL-002/03-audit/A-048-r1-fi005-residual-rereview-closure.md`）。
- **独立审计暴露并已修正的实现缺陷**（A-002 → A-003）：
  1. `F-I-001`（required）：v73 缺台账点名的「retired `records` 不存在」断言 → 新增 m0 guard（SQLite / PG 双方言）+ 负例测试 `TestV73RefusesWhenRetiredRecordsTableIsPresent`；v73 checksum 重算。
  2. `F-I-002`（recommended）：m4 的两条 PRAGMA 在 checksum 输入里每表一次、执行只一次 → `VerifyStatements` 改为每 descriptor 一次，与 `RunVerify` 同构；11 个多表 descriptor 的 checksum 随之重算（4 个单表 descriptor 不变）。
  3. `F-I-003`（recommended）：缺同进程 codec↔SQL 对拍 → 新增 `TestCodecMatchesRealMigration`（秒族 7 值 + 毫秒族 14 值）。
- **未关闭**：`F-I-004`（生成器留树）与 `F-I-005`（v86 `ModuleID` 展示列更正）为记录项；`F-I-006`（元数据同步）已修正。

## 证据

- `go build ./...` exit 0；`go test -count=1 ./...` **63/63 包 ok**（含真实 PG 15.4 路径）。
- `go test ./internal/w040contracttest/` 全绿（真实迁移边界矩阵 + codec 对拍 + v73 guard 负例）。
- `go test ./internal/store/ -run TestCompiledMigrationCatalogOwnership|TestCompleteFingerprintTracksCatalogHead` ok。

## 关门判据

A～D 全部 `completed`；`03-audit` 的 self（A-001）与 independent（A-002）意见均已落盘并由 A-003 响应（`pass`，开放 required = 0）；`D-021` residual 复审结论已落盘。故本子目标按「非关键子目标关门可经交叉审计后静默执行」的既有用户裁决**静默关门**（`done · 4/4`）。

**边界**：本关门只覆盖 codec + descriptor（M1/M2）；**不**放行 R2、**不**等于 M3/M4 完成。
