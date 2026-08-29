---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-go-library-consumption
version: 0.1.0
---

# A-001 · S2 重构回归审计（source: self · 2026-08-29）

## scope

S2 外移重构（方案 A）等价性：编译回归、全量测试、装配首验证据；S1 实验结论复核。

## verdict

**pass**（无 required；2 条 recommended）

## 核对点

| # | 项 | 证据 | 结论 |
|---|----|------|------|
| 1 | 外移范围 = 冻结面 §1/§2（kernel + modules 全量） | git mv 记录 + C 层 20 包未动 | ✅ |
| 2 | 残留引用 = 0 | 全仓 grep | ✅ |
| 3 | 编译等价 | `go build ./...` exit 0 | ✅ |
| 4 | 测试等价 | 全量仅 1 个时序敏感失败（见 F-003）；单跑/composition 全包复跑 PASS | ✅ |
| 5 | 装配闭环首验 | golden-consumer `go run` exit 0（kernel 2.0.0 · profile admin · dashboard 2.0.0 · 23 模块） | ✅ |

## findings

- **F-003（recommended）**：`TestShutdownDrainHarnessPostgres` 在全量并发下偶发 `EOF`（单跑 3.20s PASS / 全包复跑 ok）——**环境时序敏感，非重构缺陷**。建议：该 PG drain harness 移出默认全量并发（`-short` 分流或独立包），随 R5 或 VP-009 波次处理；不阻塞本目标。
- **F-004（recommended）**：golden-consumer 首验仅覆盖零依赖模块（dashboard）；`users`（含 `auth.Authenticator` C 层依赖）装配与双方言迁移台账 apply 尚未实测——S3 验证矩阵必办项（F-001 联动）。

## 响应

- F-003/F-004：采纳；F-003 登记随 R5/VP-009，F-004 纳入 GOAL-003 S3 验证矩阵（E-003 计划）。