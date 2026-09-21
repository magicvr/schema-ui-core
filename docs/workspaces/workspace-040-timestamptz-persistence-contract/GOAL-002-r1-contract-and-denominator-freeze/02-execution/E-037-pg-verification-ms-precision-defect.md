---
id: E-037-pg-verification-ms-precision-defect
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-037 · PG 侧首次经验验证：发现并修正毫秒族精度缺陷（+ D-021）

## 事实

1. **用户 2026-09-20 P-004 裁决两项**，落盘 child `01-decision/D-021-fi005-gate-and-pg-verification.md`：
   - **F-I-005 拆分**：R1 关「checksum 计算约定 / descriptor 名 / `transform_id` / 算法与输入结构 / 唯一表范围 / append-only 边界」（均已落盘）；**真实哈希值**作为**显式 R2 验收项**移交（附范围与复审触发条件），不得静默当已完成。
   - **PG 侧允许临时起容器做经验验证**（用本地已有镜像、非 5432 端口、验证后销毁）。
2. **PG 侧首次经验验证执行**：起 `postgres:16` 临时容器（映射 **15432**，命名 `w040-pg-verify`），以 `SET TIME ZONE 'UTC'` 建 legacy 形状表并**实际执行冻结的 PG `ALTER … TYPE … USING`** 语句，再回读比对。**验证后容器已 `docker rm -f` 销毁**；未改 `compose.yaml`；无残留端口。
3. **发现真实缺陷（本轮主要产出）**：毫秒族原式

   ```sql
   date_trunc('microseconds', TIMESTAMPTZ 'epoch' + <col> * INTERVAL '1 millisecond')
   ```

   在**大数值**上产生**非零误差**：

   | 输入（ms） | 原式 | 期望 |
   |-----------|------|------|
   | `253402300799999`（公元 9999 年） | `9999-12-31T23:59:59.999008` | `…59.999000` |

   误差 **+8 微秒**。根因：`BIGINT * INTERVAL` 路径内部经**浮点**转换，在 ~2.5×10¹⁴ 量级丢失微秒精度。**`::numeric` 强制转换无效**（实测同为 `.999008`）。
4. **修正为整数拆分式**（无浮点参与），实测零误差：

   ```sql
   TIMESTAMPTZ 'epoch'
     + ((<col> - CASE WHEN <col> >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second'
     + ((<col> % 1000 + 1000) % 1000) * INTERVAL '1 millisecond'
   ```

   实测（`postgres:16`，UTC）：`1758320000123` → `…20.123000`；`1758320000999` → `…20.999000`（无进位）；`-1` → `1969-12-31T23:59:59.999000`；`-999` → `…59.001000`；`-1000` → `…59.000000`；`-1001` → `…58.999000`；`-86400000` → `1969-12-31T00:00:00.000000`；`-1758320000123` → `1914-04-14T01:46:39.877000`；`253402300799999` → **`9999-12-31T23:59:59.999000`** ✅。
   PG 整数除法同样**向零截断**（实测 `-1/1000 = 0`、`-1%1000 = -1`），故修正与 SQLite 侧同构。
5. **秒族无需更改**：`to_timestamp(value::double precision)` 在 `253402300799` 与 `-1` 上实测精确。
6. **精度断言通过**：`information_schema` 实测 `data_type = timestamp with time zone`、`datetime_precision = 6`。
7. **改动的载体**：`r1-c2-per-table-pg-ddl-v1.0-fc.md` §0.1（更正块 + 更正的 E2 + §1/§3 的 mail_outbox / mail_config / operation_log 展开式）；`r1-c2-per-column-conversion-contract-v1.0-fc.md` §1 E2 与 `#34` 行；`r1-c2-c3-guardrails-v0.1.md` §2 毫秒行。

## 证据

- 验证方式：`docker run -d --name w040-pg-verify -p 15432:5432 postgres:16`；`psql` 经 `docker exec -i` 执行；容器已销毁（`docker rm -f w040-pg-verify` 返回容器名）。
- 容器列表复核：剩余 `gf-pg Exited (255) 3 weeks ago` —— **非本轮产物**（3 周前既存）。
- `compose.yaml` 未修改；`git status` 无 `compose.yaml` 变更。
- 本轮 `apps/**` **未修改**。

## 状态评估

- **开放 required = 2**（F-I-002、F-I-005）；**本条未闭合任何 required**。
- 该缺陷**若不经验证不会暴露**：`pgTimeColRe` 派生只做 `INTEGER → BIGINT` 文本替换，任何断言都看不到 +8 µs 偏差。这直接强化了 `D-019` §6「PG DDL 必须显式书写」的正当性。
- `D-021` 的 F-I-005 拆分**尚未经 independent 复审**——R1 关门前须由审计确认「约定/名/算法」确实完整、且移交给 R2 的哈希项有明确范围与复审触发条件。
- 下一步：`/audit` 复审（a）F-I-002 负值政策分档与表数 44，（b）PG 毫秒式修正，（c）`D-021` 的 F-I-005 拆分口径。
