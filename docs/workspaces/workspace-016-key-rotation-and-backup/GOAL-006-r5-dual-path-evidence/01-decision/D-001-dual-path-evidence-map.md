---
id: D-001
doc: decision-entry
goal: GOAL-006-r5-dual-path-evidence
status: accepted
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# D-001 · R5 判据 4 双路径证据映射与关门顺序

## 判据 4 → 双路径证据映射

| 路径 | 显式双密钥配置 | 证据载体（可重复） | 判定 |
|------|----------------|---------------------|------|
| **轮换路径** | current=K2, previous=K1 | `TestDualKeyRotationOverlapWindow`（auth 包 4 子用例：重叠窗可验/窗关闭拒绝/签发仅 current/过期不延长）+ `TestNewAuthenticatorWiresPreviousSecret`（生产装配路径） | 4+1 全 PASS |
| **轮换后恢复路径** | 恢复库 + current=K2, previous=K1 | `TestSQLitePostRotationRecovery`（SQLite VACUUM INTO）+ `TestPostgresPostRotationRecovery`（PG pg_dump→pg_restore；pgtest+工具门控） | 2 全 PASS |

两路径均以**显式双密钥**运行（K1≠K2，previous 非空），满足「生产向验收以显式双密钥配置为准」。

## 判据 5 越界核对单（关门时逐条）

1. 未进入 A3 / KMS / PITR / Admin 功能 / 业务域 —— 以全波次 diff 范围核对（config/auth/composition 测试与接线 + 台账，无业务模块改动）。
2. 未改 Charter —— `docs/vision/charter.md` 版本仍 `0.2.0`。
3. 未假装交付热加载 —— 密钥仅构造器注入，无运行期换钥 API。
4. 未做第二套 dump —— 备份仅消费 VP-013 既有合同与官方客户端。

## 关门顺序

1. E-001 新鲜实跑登记（四项测试 + 结果）。
2. Root 关门自审（Root `03-audit/A-001`，self）：成功标准 1–5 逐条对照 + 意见/信息台账核对。
3. Independent 审计（grok build · grok-4.6 · `/audit` 流程）：意见落盘 Root `03-audit/A-002`。
4. 编排器合并响应全部意见；required 三路径闭合后 Root `done` 5/5；goal-tree/workspace 终态同步；git checkpoint。

## 未选方案

- 把 R5 做成新的"演示脚本"：既有测试即证据载体，重复建设稀释可维护性。
- 跳过 independent 直接关门：违反 Root D-001 §5（生产路径实施按 independent）与 P-003。
