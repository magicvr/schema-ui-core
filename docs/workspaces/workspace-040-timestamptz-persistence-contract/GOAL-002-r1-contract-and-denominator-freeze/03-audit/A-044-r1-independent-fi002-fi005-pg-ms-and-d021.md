---
id: A-044-r1-independent-fi002-fi005-pg-ms-and-d021
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · finding-closure + execution-facts · commit 6ab57dd0（A-042 响应 / 负值分档 + 表数 44）对照 F-I-002 · commit d67497d7（PG 首次实测 + D-021）对照毫秒族修正式与 F-I-005 拆分 · 本机临时 postgres:16/15-alpine/17-alpine 独立复现 · 不是实施审计 · freeze-candidate ≠ 已实施 · 本地临时容器 ≠ 生产就绪证据
verdict: conditional
open_required: 2
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-044 · R1 independent · F-I-002 负值分档 / PG 毫秒族复现 / D-021 F-I-005 拆分

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure + execution-facts（用户指定三块：① `6ab57dd0` 四份 C2 载体负值分档 + 表数 44 对照 A-042 的 F-I-002 剩余；② `d67497d7` 毫秒族原式缺陷与整数拆分式，本机临时容器独立复现；③ `D-021` 对 F-I-005 的 R1/R2 拆分是否合法闭合。对照基线 A-042：开放 required = 2，F-I-002 / F-I-005）。
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区作为审计上下文**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree / `apps/` / `compose.yaml`。
- `git show --stat 6ab57dd0`：10 个文件均在 workspace-040（execution / A-042 / A-043 / 四份 C2 + superseded runbook）；**`apps/` 无变更**。
- `git show --stat d67497d7`：8 个文件均在 workspace-040（D-021 / E-037 / 五份 PG 表达式载体）；**`apps/` 无变更**。
- **证据窗口**：上述两 commit 已提交材料 + Root `D-012` / Root `D-015` / child `D-017`/`D-019`/`D-021` + 本机临时容器 `postgres:16`（16.15）/ `postgres:15-alpine`（15.19）/ `postgres:17-alpine`（17.11）。容器名 `w040-a044-pg16/15/17`，端口 **15441/15442/15443**（避开 5432 与既有 `gf-pg` 的 15432 映射），验证后 `docker rm -f` 销毁。未改 `compose.yaml`。本地临时容器 **不是** 生产就绪证据。

## 核对方法

1. 通读 A-042 剩余项原文、A-043 / E-036 / E-037 / D-021、四份点名 C2 载体、C3 边界 §5、Root D-012 / D-015。
2. 对 `attachments/r1-time-column-inventory-v0.3.md` 90 行按「倒数第二段 = 表名」独立去重。
3. 全扫冻结包 `负值` / `一律` / 旧毫秒式 `date_trunc('microseconds', TIMESTAMPTZ 'epoch' + … INTERVAL '1 millisecond')`。
4. 临时起三版本 PG，对原式 / `::numeric` / 整数拆分式 / 秒族 / `information_schema` 独立求值；自寻原式误差区间。
5. 逐项核对 D-021「R1 关门所需」五条是否落盘，以及 R2 移交是否满足 P-003 三路径与 P-005 范围+复审触发。

## 成果（有证据）

1. **负值 fail-closed 政策已按列分档写入点名的四份 C2 载体**，与 Root D-012（仅 voucher）/ Root D-015（负 epoch 合法 instant）**政策语义同一**。A-042 的「全局负值一律 m0 fail closed」**作为政策表**已被改掉。
2. **表数 20 → 44 成立**：本审对 inventory v0.3 的 90 行独立去重得 **44** 张唯一表。F-I-025 的计数子项可 `fixed`；§3.4/§3.6 陈旧句仍在。
3. **PG 毫秒族原式缺陷独立复现成立**（15.19 / 16.15 / 17.11 同形）：`253402300799999` → `.999008`（**+8 µs**）；`::numeric` 无效。整数拆分式在用户点名样本上 **split_err_us = 0**。秒族精确。`timestamptz(6)` 的 `information_schema` 断言成立。
4. **用户点名的五份 C2/草案载体已同步修正式并标注弃用**。但 **Root D-015 与 child D-012 仍把已证伪的原式写成现行权威**，未标「已更正/已弃用」。
5. **D-021 所列 R1 约定/名/算法/表范围/append-only 边界均有落盘载体**。哈希值移交有范围与复审触发，**但不得读成 F-I-005 已完成**；整条不能按 `fixed` 闭合。
6. **无新 required 编号。无新 recommended 编号。** 开放 required **维持 2**（F-I-002 / F-I-005）。**R1 不具备关门条件。**

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| F-I-002 关闭要求（A-042：冻结包负值政策与 Root D-012/D-015/测试同一） | **政策表已分档；整条仍 open** | 见第一块；C3 §5 理由句未改；Root D-015 仍写已证伪 E2 |
| F-I-025 计数 20→44 | **计数子项 `fixed`；整条仍 open** | 44 已复核；conversion §3.4/§3.6 仍陈旧 |
| 毫秒族修正式可冻进 C2 | **C2 五份点名载体已改；权威决策未改** | 见第二块 |
| F-I-005 合法闭合 | **不闭合** | 见第三块；D-021 自身声明不闭合 |
| C2 冻结 / R1 关门 / R2 放行 | **未满足** | F-I-002 / F-I-005 仍开 |

---

## 第一块：F-I-002 剩余项（`6ab57dd0`）

### 1. 四份点名 C2 载体 — 政策表已分档，路径措辞仍有「不进表达式/USING」残留

| 载体 | 政策表 / 分档 | 「全局负值 m0 fail closed」 | 「不进表达式 / 不进 USING」残留 |
|------|----------------|------------------------------|----------------------------------|
| `r1-c2-predicate-exact-sql-v1.0-fc.md` §1+§6 | **有。** 仅 `#72/#73` 的 `bucket_negative > 0` 回滚；其余负值记录不阻断 | **已去掉作为判定规则** | §1 仍写「`< 0` **不进**任何 USING 分支」；同时又写其余列必须**正常转换** |
| `r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` §2 | **有。** voucher 两列 fail closed / 其余正常转换，并指向 `TestNegativeMustFailClosed` | **已不是判定规则** | L100 **仍写**「负值 `< 0` **一律不进表达式**」，紧接「但政策按列分档」 |
| `r1-c2-per-table-rebuild-ddl-v1.0-fc.md` §0 | **有。** 与上表同一分档 + Root 限定 | **已去掉作为判定规则** | 「**不进**任何 `USING`/rebuild 的 `CASE` 分支」仍在；后接分档 |
| `r1-c2-per-column-conversion-contract-v1.0-fc.md` §3.5 | **有。** 仅 `#72/#73` fail closed；其余正常转换 | **已去掉作为判定规则** | 「`< 0` 不进任何 `USING` 分支」仍在；`#72/#73` 行保留 voucher 专属「`< 0` 不进 USING」——**这行正确** |

**本审对「全局负值 fail closed」的判定**：作为 **m0 判定规则**，四份点名载体已不再把负值 fail closed 写成全局。`#72/#73` 行保留 voucher 专属 fail closed，与 A-043 声明一致。

**残留（路径措辞，不是政策表）**：mechanism L100 的「**一律**不进表达式」是 A-042 **逐句点名**的原文，本轮只在其后追加分档表，**没有删掉「一律」**。若 R2 只读该句、不读政策表，仍可能把普通列负值排除出转换表达式。同文件 §2.3 / §2.4 实际是 NN 列直接套 floor 表达式、CASE 不得出现 `< 0`——这才是「普通列走表达式」。两句互相拉扯。

### 2. 分档政策 vs Root D-012 / D-015 — 政策语义同一，不是 D-012 三段的逐字复制

| 权威 | 原文要点 | 冻结包是否同一 |
|------|----------|----------------|
| Root `D-012-voucher-invalid-value-policy.md` | **仅** `vouchers.expires_at` / `redeemed_at`：legacy `0`→`NULL`；**负值数据损坏 fail closed**；正值按 Unix seconds 转换；预检分桶 0/负/正；runtime 不得继续用 `>0` 把负值当 absence | **同一。** 点名载体把 fail closed 收成 `#72/#73`；USING 只处理 `=0→NULL`+正值；runtime `if exp.Valid`。三桶预检仍对 sentinel 列计数，但对非 voucher 的 `bucket_negative` **不阻断**——这是对 D-012 的**范围收窄执行**，不是把 D-012 扩到普通列 |
| Root `D-015-negative-instant-truncation.md` | **负 epoch 不是 sentinel**；合法 instant；仅 Root D-012 的 0 sentinel 映射 NULL | **负值政策同一。** 四份载体写「其余全部时间列正常转换」 |

**不是逐字粘贴**：D-012 全文三句，冻结包是表格+判定伪 SQL。本审要的是政策同一，不是字面拷贝。

**D-015 的毫秒族 SQL 字面不在本块「负值政策」比对范围内**——它属于第二块：D-015 仍写已证伪的原式，并声称「全程不经二进制浮点」。这与负值分档不是同一子项，但会阻断 F-I-002 的表达式唯一性（见第二块）。

另：exact SQL §1 标题仍写「Sentinel 0→NULL 列（5 列，**Root D-012 政策**）」。D-012 **不是** 五列 sentinel 的权威（sentinel 0→NULL 是 D-008 族）。政策块正文已把 D-012 收到 voucher。属卫生，不单独升格 required。

exact SQL §1 可执行断言路径写成 `apps/api/internal/wcontracttest`（缺 `040`）。实际包是 `w040contracttest`。卫生。

### 3. 表数 20 → 44 — 独立复算成立

方法（与 A-043 声明同一）：对 inventory v0.3 90 行，取反引号标识符的 **倒数第二段 = 表名** 去重。

| 口径 | 本审 | 附件 |
|------|------|------|
| 行数 | **90** | 对 |
| 唯一表 | **44** | rebuild §0 现写 **44 张不同的表**（机械去重；分母仍 90 列）。**对** |
| ms 列 | 10（`#33/#34/#35/#36/#37/#41–#45`） | 与 conversion 90 行一致 |

44 张表（字母序）：`captcha_challenges` `captcha_config` `data_scope_policies` `dict_entries` `dict_types` `digital_entitlements` `digital_offers` `digital_purchases` `email_verification_challenges` `jobs` `login_failures` `mail_config` `mail_outbox` `menu_items` `mfa_proofs` `notifications` `operation_log` `operation_log_archive` `password_recovery_challenges` `permissions` `recycle_items` `refresh_tokens` `roles` `scheduled_tasks` `schema_migrations` `service_credentials` `site_settings` `subjects` `system_data_reconcile` `task_runs` `telegram_config` `telegram_inbound_messages` `telegram_outbound_messages` `telegram_sessions` `user_data_scopes` `user_invites` `user_mfa` `user_password_history` `users` `voucher_batches` `vouchers` `wallet_accounts` `wallet_ledger_entries` `wallet_reconciliation_runs`。

推导可复核。F-I-025 **计数子项 `fixed`**。

F-I-025 仍 open：conversion contract §3 标题仍写「必须收口的 3 项（本候选仍未满足）」且 **§3.4「逐表 exact SQLite rebuild DDL 仍未写出」**、**§3.6「双方言 checksum 约定二选一：未选定」** 仍在。前者与已落盘的 rebuild 附件矛盾；后者与 child `D-017` / ledger §2 已裁决选项 A 矛盾。A-042 已把这两句并入 F-I-025 剩余，本轮未划掉。

### 4. A-042 点名但本轮未改的句子

| A-042 点名 | 本轮 |
|------------|------|
| mechanism L100「一律不进表达式」 | **句子仍在**，后接分档表 |
| rebuild §0 负值行 | 已加分档；「不进 CASE」仍在 |
| exact SQL §1 全局句 | 已加分档；「不进任何 USING」仍在 |
| conversion §3.5 | 已加分档；「不进任何 USING」仍在 |
| **C3 边界 §5 断言 4 理由句** | **未改。** `r1-c3-backup-recovery-boundary-v1.0-fc.md` L148 仍写：「转换是 fail-closed，**B 内不应再有负值**」。A-043 只改了已 `superseded` 的 runbook `negative-invalid` 标注，**没有**改这条现行 C3 权威句。该句对普通列为假：负 epoch 转换后以 pre-1970 TEXT/`timestamptz` 出现在 B 中（A-042 已用测试证明） |

### 5. F-I-002 现在能否闭合？

**不能。** A-042 的「全局 fail closed 作为政策」已被四份 C2 政策表收口，**政策子项可记为再收窄**。整条仍 open。**全部剩余项**：

**required（阻断 C2 冻结）**

1. **C3 边界 §5 断言 4 理由句仍把转换写成全局 fail-closed、并断言 B 内不应有负值**（A-042 点名，本轮未改现行 C3 权威文件）。
2. **表达式权威不唯一（本轮 PG 复现后新点名，仍归本条）**：点名 C2 载体已把毫秒族改为整数拆分式，但 **Root `D-015` L18 与 child `D-012-v73-allocation-negative-truncation.md` L13 仍把已证伪原式写成现行毫秒族**，并声称「全程不经二进制浮点」。C2 不得在「权威决策写原式、冻结包写修正式」并存时冻结 USING。见第二块。
3. freeze-candidate ≠ 实施（边界；单独不阻断，与 F-I-003/F-I-006 同一尺子）。

**recommended / 卫生（不单独挡本条设计层闭合，但应一并改）**

4. mechanism L100「一律不进表达式」与「普通列走表达式」并存。
5. exact SQL §1「不进任何 USING 分支」全局句；`wcontracttest` 路径笔误；§1 标题把 5 列 sentinel 归到 Root D-012。
6. conversion §3.4 / §3.6 陈旧句（并入 F-I-025）。

关闭要求：改 C3 §5 理由句（普通列负 epoch **会**出现在 B；负值非法只属于 voucher `m0`）；把 Root D-015 / child D-012 的毫秒族改成与冻结包同一整数拆分式或显式标「已更正/已弃用」并指向 pg-ddl §0.1；删掉「一律不进表达式」这类与政策表相反的路径句。不要把临时 PG 容器当生产证据，也不要把测试包当 C2 已实施。

---

## 第二块：PG 毫秒族修正（独立复现）

### 实验装置

| 项 | 本审 |
|----|------|
| 镜像 | 本地已有 `postgres:16` / `postgres:15-alpine` / `postgres:17-alpine` |
| 实测版本 | **16.15**（Debian）、**15.19**、**17.11** |
| 容器 | `w040-a044-pg16/15/17`，`--rm`，端口 **15441/15442/15443**（非 5432；不占用 `gf-pg` 历史 15432） |
| 时区 | `SET TIME ZONE 'UTC'` |
| 事后 | `docker rm -f` 三容器；`docker ps` 无 `w040-a044-*`；`compose.yaml` 无变更；未留常驻服务 |
| 边界 | **本地临时验证，不是生产就绪证据** |

三版本在下列结论上**同形**（数值一致）。下文以 16.15 为代表。

### 1. 原式误差 — 复现成立，且不止一个点

原式：`date_trunc('microseconds', TIMESTAMPTZ 'epoch' + x * INTERVAL '1 millisecond')`。

用户点名样本 `253402300799999`（公元 9999-12-31 23:59:59.999）：

| 版本 | 原式 | 期望 | 误差 |
|------|------|------|------|
| 15.19 / 16.15 / 17.11 | `9999-12-31T23:59:59.999008` | `.999000` | **+8 µs** |

`x::numeric * INTERVAL '1 millisecond'` **同为 `.999008`**。`::numeric` 无效，与 E-037 一致。

**本审自寻的其它受影响点（原式 − 期望，单位 µs）**：

| x (ms) | 约略年代 | orig_err_us |
|--------|----------|-------------|
| ≤ `5e13` 的抽样（含 remainder 0/1/123/500/998/999） | ≲ 3554 | **0** |
| `8e13` remainder 1 / 123 / 999 | ~4504 | **−8 / +8 / +8** |
| `1e14` remainder 1 / 123 / 999 | ~5138 | **−8 / +8 / +8** |
| `2e14` remainder 998 / 999 | ~8307 | **+16 / +8** |
| `253402300799990` … `…999` | 公元 9999 末 | **+16 / +8 / 0 / −8 / −16** 随余数振荡 |

同网格里 **整数拆分式 split_err_us 全 0**。

结论：E-037「大数值 +8 µs、`::numeric` 无效」**成立**，但「只在公元 9999 这一个输入」**不完整**。原式在合同上界（公元 9999）邻域与至少从 ~8×10¹³ ms（约公元 4504）起的部分 remainder 上出现 ±8/±16 µs。近端生产量级（~1.7×10¹² ms / 2025）抽样为 0。负向抽样到 −2×10¹⁴（公元 436 年量级）在本网格为 0；`−253402300799999` 超出 PG timestamp 范围，原式报 `timestamp out of range`。15/17 与 16 **同一缺陷**。

整十、整千的「圆整」大数（如 `2.5e14`、`253402300799000`）可以为 0，**不能**用圆整网格代替 `.999` 余数。

### 2. 整数拆分式 — 点名样本零误差

公式与 pg-ddl §0.1 同一。`UPDATE` 写入 `timestamptz(6)` 列后回读：

| 输入 (ms) | 回读（UTC） | 独立期望（Python floor `divmod`） | 误差 |
|-----------|-------------|-------------------------------------|------|
| 0 | `1970-01-01T00:00:00.000000` | 同 | 0 |
| 1758320000123 | `2025-09-19T22:13:20.123000` | 同 | 0 |
| 1758320000999 | `2025-09-19T22:13:20.999000` | 同（无进位） | 0 |
| −1 | `1969-12-31T23:59:59.999000` | 同 | 0 |
| −999 | `…59.001000` | 同 | 0 |
| −1000 | `…59.000000` | 同 | 0 |
| −1001 | `…58.999000` | 同 | 0 |
| −86400000 | `1969-12-31T00:00:00.000000` | 同 | 0 |
| −1758320000123 | `1914-04-14T01:46:39.877000` | 同 | 0 |
| 253402300799999 | `9999-12-31T23:59:59.999000` | 同 | 0 |

`EXTRACT(EPOCH)*1e6 − x*1000` 在上表及 9×10¹² / 1×10¹³ / 1×10¹⁴ / 2×10¹⁴ 上均为 0。

PG 整数除法向零：`−1/1000 = 0`、`−1%1000 = −1`，与 E-037 及 SQLite 侧同构。负值必须保留 `−999` 修正与 `(x%1000+1000)%1000` 归一化。

**不把「无浮点参与」升格为已证明的实现层不变量**（`bigint * interval` 在 PG 内仍可能走 `float8`）；本审只接受：**合同边界与本网格上无可见误差**。

### 3. 秒族 — 精确（含负值与公元 9999）

`date_trunc('microseconds', to_timestamp(x::double precision))`：

| 输入 (sec) | 输出 | 期望 | 误差 |
|------------|------|------|------|
| 0 | `1970-01-01T00:00:00.000000` | 同 | 0 |
| −1 | `1969-12-31T23:59:59.000000` | 同 | 0 |
| −86400 | `1969-12-31T00:00:00.000000` | 同 | 0 |
| 1758320000 | `2025-09-19T22:13:20.000000` | 同 | 0 |
| −1758320000 | `1914-04-14T01:46:40.000000` | 同 | 0 |
| 253402300799 | `9999-12-31T23:59:59.000000` | 同 | 0 |

`sec_err_us = 0`。秒族无需改动，与 E-037 一致。`253402300799 * INTERVAL '1 second'` 与 `to_timestamp` 在该点亦差 0。

### 4. `timestamptz(6)` information_schema — 成立

`CREATE TABLE a044_prec (c timestamptz(6))` 与拆分式 `UPDATE` 后的 `ts` 列：

| 列 | data_type | datetime_precision |
|----|-----------|-------------------|
| `c` / `ts` | `timestamp with time zone` | **6** |

三版本同一。

### 5. 修正式载体同步 — 点名五份已改；权威决策未改

| 载体 | 本审 |
|------|------|
| `r1-c2-per-table-pg-ddl-v1.0-fc.md` | **已改。** §0.1 更正块 + E2 整数拆分式；mail_outbox / mail_config / operation_log 已展开。jobs / archive 仍「同形 / 骨架」（A-036 已接受骨架形式），骨架指向 §0.1 |
| `r1-c2-per-column-conversion-contract-v1.0-fc.md` | **已改。** §1 E2 + `#34` ELSE；更正注记 |
| `r1-c2-c3-guardrails-v0.1.md` | **已改。** 毫秒行整数拆分 +「已弃用」 |
| `r1-c2-column-contract-draft-v0.1.md` | **已改。** 整数拆分 +「已弃用」 |
| `r1-c2-column-contract-matrix-v0.2.md` | **已改。** 整数拆分 +「已弃用」 |

**仍持有旧式、且未标「已更正/已弃用」的现行权威**（不是历史审计正文）：

| 载体 | 旧式角色 |
|------|----------|
| **Root `D-015-negative-instant-truncation.md` L18** | 现行毫秒族定义，并写「全程不经二进制浮点」。**已被本审证伪** |
| **child `D-012-v73-allocation-negative-truncation.md` L13** | 现行毫秒族，指向 Root D-015 |
| `D-021` 决定 2 验证目标句 | 仍引用原式作为「要验证的两条表达式」——作为当时验证任务陈述可理解，但未回写「原式已弃用」 |

历史 A-018/A-020/A-022/A-036 正文保留旧式是台账，不要求改写。

**本块结论**：E-037 的缺陷发现与修正式在临时容器上成立；点名五份载体已同步。C2 **仍不能**把 E2 冻成唯一式，因为 Root D-015（F-I-002 一直引用的截断权威）还在发布已弃用原式。并入 F-I-002 剩余第 2 项。

---

## 第三块：F-I-005 拆分口径（D-021）

### 1. 「R1 关门所需」逐项落盘核对

| D-021 项 | 声称载体 | 本审 |
|----------|----------|------|
| checksum 计算约定（单 checksum / SQLite DDL 切片 / PG 不进哈希 / `transform_id` 无方言后缀） | child `D-017-v73-checksum-convention.md` | **已落盘。** 用户书面选 A；与 `apps/api/kernel/persistence.go:14-17` 现行 `MigrationChecksum` / `normalizeSQL` **同一**。D-017 自身写明「不等于 checksum 已记录」 |
| 15 个 descriptor 的 `Name` 与 `transform_id` | ledger `r1-c2-descriptor-ledger-v1.0-fc.md` §1 | **已落盘。** v73–v87 共 15 行，形如 `0073:vp040-temporal-core-persistence:v1` |
| `MigrationChecksum` 算法与输入结构（有序语句 + `m0–m5`） | 同台账 §2 | **已落盘。** 逐字引用 `persistence.go`；`m0→m5` 序位表在 |
| 唯一表范围（v73/v74/v85 不相交；v74 = `system_data_reconcile` + auth） | 同台账 §1.1 | **已落盘。** `schema_migrations` 仅 v73；列分配 3+31+…+6=90 |
| append-only 边界（v1–v72 不可变；`rebuildOperationLog` 断言不改 stmts） | `D-019` §5 | **已落盘，口径略窄于标题。** §5 禁止改 `0004–0071` DDL 切片、改 Go 控制流不进哈希。v1–v72 不可变另见 D-017 影响段与 ledger §3。合在一起覆盖 D-021 所写边界 |

**卫生**：`01-decision.md` 索引仍止于 **D-012**，**未登记 D-017～D-021**（文件在 `01-decision/` 目录内）。不使上述载体「未落盘」，但索引不是完整台账。`02-execution.md` 已含 E-036/E-037。

### 2. 移交 R2 的验收项 — 哈希项有范围+触发；另两项写在清单里但未进入「范围」句

D-021「移交 R2 的显式验收项」列了三条：

1. 15 个 descriptor 的 canonical SQL 落码 + **真实 `MigrationChecksum` 哈希**写入 ledger；
2. `migrate_test.go` / `postgres_test.go` 的 v73+ 追加行与**金额列断言拆分**；
3. leftover 21 名补入 PG 断言集合。

P-005 字段：

| 要求 | 哈希项 | 测试改写 / leftover 21 |
|------|--------|------------------------|
| 明确范围 | **有。** 「本拆分只适用于『哈希值』这一子项」 | 出现在清单，**未进入「范围」句** |
| 可操作复审触发 | **有。** 「R2 落码后首次记录哈希时，须由 independent 复审该哈希与约定一致；若约定被修订，本拆分的 R1 部分随之回退」 | **无对应触发** |
| 不是「以后再收集」 | 是（具体产物 + 何时复审） | 清单具体，但未写成 residual 字段 |

哈希项满足 P-005「范围 + 复审触发」。测试改写与 leftover 21 若仍算 F-I-005 剩余（A-030/A-042 原文含「可执行测试改写」），则 **residual 覆盖不完整**。

### 3. 会不会被读成「F-I-005 已完成」？

**会被误读，若编排器或后续审计只看「R1 关门所需均已有」。** D-021 自己写了两道闸：

- 「**不得**把本拆分读作『F-I-005 已完成』——它是**有范围的移交**，不是闭合声明。」
- 「本决定**不闭合任何 required**；F-I-005 的 R1 侧仍需 independent 复审确认约定/名/算法确实完整。」

本审确认 R1 侧五条**有载体且内容对位**，同时 **拒绝**把整条 F-I-005 标 closed。

### 4. P-003 三路径落在哪一条？

| 路径 | 是否适用 |
|------|----------|
| **`fixed`** | **否。** 原 finding 要的「已记录 canonical SQL / 真实哈希 / 可执行测试改写」并未发生。`apps/` 未改。 |
| **`user-overruled`** | **否。** 用户没有驳回或降级 F-I-005，只是把哈希（及清单上的测试/leftover）挪到 R2。 |
| **`accepted-residual`** | **哈希子项形态接近，整条还不完整。** 用户书面 P-004 落在 D-021；哈希有范围与复审触发。缺：D-021 **明确说不闭合任何 required**；测试改写/leftover 21 未进入 residual 范围句与复审触发；P-003 residual 还建议写清为何可接受、适用期限、缓解/监控、责任人——D-021 只部分覆盖。 |

**三条都不完全适用到「整条 F-I-005 合法闭合」。**

建议给编排器（本审不代裁）：

1. **本轮不要把 F-I-005 标 closed。** 把 R1 约定/名/算法/表范围/append-only 记为**子项 `fixed`**（本审接受这些子项已落盘）。
2. 若用户意图是 R1 关门、哈希等 R2 再核：用 `/govern` 把 D-021 补成完整 **`accepted-residual`**——范围覆盖哈希 **以及**（若仍属本条）测试改写与 leftover 21；每项复审触发；责任人；约定被改则 R1 子项回退。补完后下一次 independent 才可将整条按 residual 闭合。
3. 不要把 D-021 读成 `fixed`。不要在 residual 未写全时放行 R1 关门。

---

## Findings

### F-I-002 · 逐列 USING / rebuild / codec

- **严重度**：high · **建议**：required
- **状态**：open（**再收窄，仍不关闭**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修**：A-042「冻结包全局化负值 fail closed」作为 **m0 判定规则** → 四份点名 C2 载体已按列分档，与 Root D-012 / D-015 **政策语义同一**；表数 20→44 已写入 rebuild §0 且本审复算 44。
- **仍不闭合（全部剩余项）**：见第一块 §5（C3 §5 理由句；Root D-015/child D-012 仍发布已证伪 E2；路径措辞「一律不进表达式」；freeze-candidate ≠ 实施）。
- **关闭要求**：C3 §5 改成「普通列负 epoch 会出现在 B；仅 voucher `m0` fail closed」；Root D-015 / child D-012 与冻结包 E2 收成同一整数拆分式或标弃用；删掉与政策表相反的「一律不进表达式」。

### F-I-003 / F-I-004 / F-I-006 … F-I-024 / F-I-026 / F-I-027

- 维持既有 closed（本轮未复审这些条的关闭证据，不重开）。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high · **建议**：required
- **状态**：open（**R1 约定子项可记 `fixed`；整条不关闭**）
- **本轮**：D-021 拆分落地。R1 五条有载体（D-017、ledger §1/§1.1/§2、D-019 §5 + ledger §3）。真实哈希 / 测试改写 / leftover 21 在 R2 清单。
- **不按 `fixed` 闭合。不按完整 `accepted-residual` 闭合。** 见第三块。
- **关闭要求**：要么 R2 落码后记录哈希并经 independent 复审（`fixed`），要么把 D-021 补成覆盖全部延期子项的书面 residual（范围+触发+责任人）后再复审 residual 闭合。

### F-I-025 · 冻结包卫生（计数 + 陈旧句）

- **严重度**：low · **建议**：recommended
- **状态**：open（**计数子项 `fixed`；整条仍 open**）
- **本轮**：20→**44** 独立复算通过。**仍 open**：conversion §3.4「DDL 仍未写出」、§3.6「checksum 二选一未选定」。
- **关闭要求**：划掉这两句（指向 rebuild 附件与 D-017）。

### F-I-008 / F-I-009

- 维持 recommended open（本轮未触及 R3 回归矩阵）。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮 |
|----|------|------------|------|
| F-I-002 | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec | **再收窄**：负值政策表已分档；表数 44。**仍缺** C3 §5 理由句、Root D-015/child D-012 与修正式同一 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；不得把「无哈希」读成已完成 | **R1 约定子项落地；整条仍 open**。哈希 = 有范围的 R2 移交，不是 `fixed` |
| F-I-025 | 卫生 | 不单独挡冻结 | 计数 `fixed`；§3.4/§3.6 仍 open |

F-I-001、F-I-003、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015、F-I-016、F-I-017、F-I-018、F-I-019、F-I-006、F-I-020、F-I-021、F-I-022、F-I-023、F-I-024、F-I-026、F-I-004、F-I-027 为 closed。F-I-008、F-I-009、**F-I-025** 为 recommended open。

**开放 required = 2**（F-I-002、F-I-005）。在这些合法闭合前：不得冻结 C2、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。**R1 不具备关门条件。**

## 与既有意见的异同

| 项 | A-042 independent | A-043 / E-036 / D-021 / E-037 自称 | A-044 independent（本条） |
|----|-------------------|--------------------------------------|---------------------------|
| verdict | conditional；open required=2 | 不自证闭合；待本审 | **conditional**；open required=**2** |
| F-I-002 | 测试 `fixed`；整条 open（全局 fail closed） | 政策已分档；待复审 | 政策表同一 Root D-012/D-015；**整条仍 open**（C3 §5 + D-015 旧 E2） |
| F-I-005 | open | D-021 拆分，不自称闭合 | **仍 open**；R1 子项落地；非 `fixed` / 非完整 residual |
| F-I-025 | open（20≠44） | 已改 44 | 计数 `fixed`；§3.4/§3.6 仍 open |
| PG 原式 | 未测 | +8 µs @ 9999；split 零误差；秒族精确 | **复现成立**（15/16/17）；受影响区间宽于单点；split 点名样本 0 误差 |
| 新 finding | 无 | 无 | **无新 required / 无新 recommended 编号** |
| R1 关门 | 禁止 | 待本审 | **不具备关门条件** |

无「一要一否」需用户在 finding 之间裁。F-I-005 走哪条 P-003 路径见第三块建议，**本审不代用户把整条标 residual closed**（D-021 自己说不闭合；测试/leftover 覆盖不全）。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | F-I-002 仍开（C3 理由句 + E2 权威不唯一），阻断 C2 |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列分母未扩大；44 张表口径已写入 rebuild §0 且本审复算 |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 设计层 closed 维持；C3 §5 理由句仍误；实施仍在 R2 |
| I-040-004 | required | R3 | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001 仍开放，阻断 C2。无用户书面 residual 覆盖 F-I-002。D-021 对 F-I-005 哈希项是有界移交，不是 I-040-* 的 verified。

## 结论 + 建议给编排器/用户的下一步

**conditional。** `6ab57dd0` 把 A-042 的负值 **fail-closed 政策**收成 voucher 专属，表数 44 可复核；**F-I-002 整条仍 open**（C3 §5 理由句未改；Root D-015 仍发布已证伪毫秒式）。`d67497d7` 的原式 +8 µs 与整数拆分式零误差在 15/16/17 上独立成立，点名五份载体已同步，**权威决策未同步**。D-021 **不能**让 F-I-005 按 `fixed` 闭合；哈希移交接近 residual 但整条覆盖不全。开放 required **2**。**R1 不具备关门条件。**

建议 `/govern`：

1. 响应本 A-044；**不要**把 F-I-002 或 F-I-005 标 closed；可接受 F-I-025 计数子项 `fixed`、F-I-005 的 R1 约定子项 `fixed`。
2. 最小文档收口：C3 §5 理由句；Root D-015 / child D-012 毫秒族与 pg-ddl §0.1 同一或标弃用；mechanism「一律不进表达式」；conversion §3.4/§3.6；顺手补 `01-decision.md` 索引 D-017～D-021。
3. F-I-005：按第三块补全 residual 字段后再请 independent 复审，或保持 open 直到 R2 哈希落地。
4. **不要**冻结 C2，**不要**启动 R2 生产 schema/codec，**不要**改 formatter/历史 DDL。不要把本审临时容器结果当生产就绪证据。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree / `apps/`。响应、finding 闭合与是否推进由 `/govern` 处理。
