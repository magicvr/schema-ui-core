---
id: E-001
doc: execution-entry
goal: GOAL-006-r5-dual-path-evidence
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-001 · R5 双路径新鲜实跑登记（2026-08-22）

## 轮换路径（显式 current=K2, previous=K1）

| 测试 | 结果 |
|------|------|
| `go test ./internal/auth/ -run TestDualKeyRotationOverlapWindow -count=1 -v` | **PASS**（3.36s；重叠窗可验 / 窗关闭拒绝 / 签发仅 current / 过期不延长） |
| `go test ./internal/composition/ -run TestNewAuthenticatorWiresPreviousSecret -count=1 -v` | **PASS**（0.47s；生产装配路径双密钥通过 / 空 previous 单密钥拒绝） |

## 轮换后恢复路径（恢复库 + 显式 K2+prev=K1）

| 测试 | 结果 |
|------|------|
| `TestSQLitePostRotationRecovery -count=1 -v` | **PASS**（0.56s；VACUUM INTO 备份 → 轮换启动 → A1/A2/A3） |
| `TestPostgresPostRotationRecovery -count=1 -v`（`R16_PG_DUMP_CONTAINER=r2-pg-probe`） | **PASS**（4.84s；pg_dump→pg_restore + ledger 指纹核对 → 恢复库启动 → A1/A2/A3） |

## 判据 5 越界核对（关门单，2026-08-22 实核）

1. 全波次代码 diff（基线 `5195104` → HEAD）：仅 config 面 / auth 核心 / composition 接线与测试 / YAML 样例，共 10 文件 —— **未进入** A3 / KMS / PITR / Admin 功能 / 业务域。
2. `docs/vision/charter.md` 零改动（`git diff --name-only` 为空）—— Charter 仍 `schema-ui-core-admin-foundation@0.2.0`。
3. 密钥仅构造器注入、无运行期换钥 API —— 无热加载宣称。
4. 备份仅消费 VP-013 合同与官方 pg_dump/pg_restore 客户端 —— 无第二套 dump。

## 判定

VP 方向级判据 4 达成：显式双密钥下轮换路径与轮换后恢复路径各有 ≥1 条可重复、当日新鲜实跑的可核对证据。判据 5 五项越界核对全部通过。
