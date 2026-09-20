---
id: GOAL-002-r1-contract-and-denominator-freeze
doc: audit
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.15
---

# 审计 · GOAL-002

## 信息就绪核对

| 核对项 | 状态 | 备注 |
|--------|------|------|
| I-040-001～003 | collecting | A-006 接受 F-I-001 closed（90 列 + catalog 72 + v1–v72）；A-014 接受 F-I-014 closed；A-016 接受 leftover 列名表已列出及 F-I-015 碰撞 closed；A-018 接受 Root D-012/D-013 方向唯一及 F-I-004 开放标记准确；A-020 接受 A-019 已把 guardrails/column-contract/matrix 收成同一 `date_trunc`+整数 interval（F-I-002 表达式子项 `fixed`）；A-022 接受 A-021 对 F-I-016 的关闭（当时 `02-execution.md` 为 E-001～E-023 严格递增）；A-025 确认 A-024 的 E-020～E-023 四份 proposed 草案可收窄但不可闭合 F-I-002～006；A-027 确认 A-026 的 Root D-014/D-015 与 child D-012 可收窄 F-I-005（未发布 allocation baseline）与 F-I-002（负瞬间 Go Truncate 解释），仍不可闭合五条 required；A-029 确认 A-028 未关 F-I-002～006：已接受 allocation 列号不相交，但 owner spec v74「ledger/reconcile」仍在；E-017 双文件已改为唯一 E-024（F-I-017 closed）；SQLite runbook 已区分 rollback snapshot 与 RecoveryPoint，PG restore 仍绑预转换 `<artifact>`；D-015「integer interval」与秒列 `to_timestamp(double)` 仍未收口；F-I-018 保持 open（冻结包仍无限定 D-012）；新增 F-I-019（E-024 插在 E-016 与 E-017 之间）；草案/baseline 不是实施证据；A-030 接受 E-025/E-026 后 **F-I-003 closed**（90 列 mapping）、**F-I-018/F-I-019 closed**、F-I-002 的 D-015 字面子项 `fixed`、F-I-005 的 v74 同文子项 `fixed`；**仍 open required = 4**（F-I-002/004/005/006）；freeze-candidate 不是实施证据 |
| I-040-004 | open | R3 回归矩阵；R1 接口仍未登记（A-002/A-004/A-006/A-010/A-012/A-014 F-I-009） |
| 资料引用 | 无 | 工作区 `shared_materials_catalog: none` |

> **2026-09-20 本轮（A-030）后的口径**：A-030 为最新 independent 意见，`open required = 4`（F-I-002 / F-I-004 / F-I-005 / F-I-006）。E-025/E-026 的三份 `freeze-candidate` 被接受为设计证据扩容：**F-I-003 closed**（90 列 mapping）；**F-I-018 / F-I-019 closed**；F-I-002 的 D-015 字面子项 `fixed`（裁决 B）；F-I-005 的 owner-spec v74 同文子项 `fixed`。**F-I-002/004/005/006 仍未闭合；C2/C3 未冻结；R2 未放行。** freeze-candidate 不是实施证据。

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required | 文件 |
|------|------|--------|-------|---------|---------------|------|
| A-001 | 2026-09-20 | self | C1/C2/C3 readiness | conditional | 3（历史 self；independent 不接受其 F-R1-001 fixed） | `03-audit/A-001-r1-self-readiness.md` |
| A-002 | 2026-09-20 | independent | C1/C2/C3 readiness + F-R1-001 关闭复审 | conditional | 6 | `03-audit/A-002-r1-independent-readiness.md` |
| A-003 | 2026-09-20 | self | response to A-002 / inventory re-audit | conditional | 5 | `03-audit/A-003-r1-self-response-to-independent.md` |
| A-004 | 2026-09-20 | independent | F-I-001 关闭复审 + F-I-002..006 再评估 + D-003 公共 wire | conditional | 7 | `03-audit/A-004-r1-independent-reaudit-after-a003.md` |
| A-005 | 2026-09-20 | self | response to A-004 / catalog 72 + wire inventory correction | conditional | 5 | `03-audit/A-005-r1-self-response-to-a004.md` |
| A-006 | 2026-09-20 | independent | F-I-001 关闭复审 after A-005 + F-I-002..006 / F-I-010 / D-004 | conditional | 6 | `03-audit/A-006-r1-independent-reaudit-after-a005.md` |
| A-007 | 2026-09-20 | self | response to A-006 / F-I-001 closure + E-006 index fix | conditional | 6 | `03-audit/A-007-r1-self-response-to-a006.md` |
| A-008 | 2026-09-20 | self | C2/C3 guardrails readiness / Backup SPI surface | conditional | 7 | `03-audit/A-008-r1-self-guardrails-readiness.md` |
| A-009 | 2026-09-20 | self | response to A-008 / Backup Port surface user decision | conditional | 6 | `03-audit/A-009-r1-self-response-backup-port.md` |
| A-010 | 2026-09-20 | independent | C2/C3 guardrails freeze gate after A-009 | conditional | 6 | `03-audit/A-010-r1-independent-c2-c3-guardrails.md` |
| A-011 | 2026-09-20 | self | response to A-010 / C2/C3 user decisions and wire coverage | conditional | 6 | `03-audit/A-011-r1-self-response-to-a010.md` |
| A-012 | 2026-09-20 | independent | C2/C3 freeze-gate follow-up after A-011 / D-008 / D-009 | conditional | 6 | `03-audit/A-012-r1-independent-after-a011-d008-d009.md` |
| A-013 | 2026-09-20 | self | response to A-012 / freeze-package alignment | conditional | 6 | `03-audit/A-013-r1-self-response-to-a012.md` |
| A-014 | 2026-09-20 | independent | C2/C3 design-evidence follow-up after A-013 | conditional | 5 | `03-audit/A-014-r1-independent-after-a013-c2-c3-evidence.md` |
| A-015 | 2026-09-20 | self | response to A-014 / concrete C2/C3 evidence | conditional | 5 | `03-audit/A-015-r1-self-response-to-a014.md` |
| A-016 | 2026-09-20 | independent | C2/C3 design-evidence follow-up after A-015 | conditional | 5 | `03-audit/A-016-r1-independent-after-a015-c2-c3-evidence.md` |
| A-017 | 2026-09-20 | self | response to A-016 / precision-voucher-monotonic-owner evidence | conditional | 5 | `03-audit/A-017-r1-self-response-to-a016.md` |
| A-018 | 2026-09-20 | independent | C2/C3 design-evidence follow-up after A-017 | conditional | 5 | `03-audit/A-018-r1-independent-after-a017-c2-c3-evidence.md` |
| A-019 | 2026-09-20 | self | response to A-018 / matrix expression + E-index correction | conditional | 5 | `03-audit/A-019-r1-self-response-to-a018.md` |
| A-020 | 2026-09-20 | independent | C2/C3 design-evidence follow-up after A-019 | conditional | 5 | `03-audit/A-020-r1-independent-after-a019-c2-c3-evidence.md` |
| A-021 | 2026-09-20 | self | response to A-020 / E-index correction | conditional | 5 | `03-audit/A-021-r1-self-response-to-a020.md` |
| A-022 | 2026-09-20 | independent | C2/C3 design-evidence follow-up after A-021 / F-I-016 closure | conditional | 5 | `03-audit/A-022-r1-independent-after-a021-c2-c3-evidence.md` |
| A-023 | 2026-09-20 | self | response to A-022 / E-index closure | conditional | 5 | `03-audit/A-023-r1-self-response-to-a022.md` |
| A-024 | 2026-09-20 | self | C2/C3 design evidence expansion | conditional | 5 | `03-audit/A-024-r1-self-response-design-evidence-expansion.md` |
| A-025 | 2026-09-20 | independent | C2/C3 design-evidence follow-up after A-024 / E-020..E-023 drafts | conditional | 5 | `03-audit/A-025-r1-independent-after-a024-c2-c3-evidence.md` |
| A-026 | 2026-09-20 | self | response to A-025 / accepted allocation baseline | conditional | 5 | `03-audit/A-026-r1-self-response-to-a025.md` |
| A-027 | 2026-09-20 | independent | C2/C3 design-evidence follow-up after A-026 / D-014 / D-015 | conditional | 5 | `03-audit/A-027-r1-independent-after-a026-c2-c3-evidence.md` |
| A-028 | 2026-09-20 | self | response to A-027 / owner and backup boundary hygiene | conditional | 5 | `03-audit/A-028-r1-self-response-to-a027.md` |
| A-029 | 2026-09-20 | independent | C2/C3 design-evidence follow-up after A-028 / owner / E-ID / D-012 / backup boundary | conditional | 5 | `03-audit/A-029-r1-independent-after-a028-c2-c3-evidence.md` |
| A-030 | 2026-09-20 | independent | C2 freeze-candidate follow-up after E-025/E-026 / A-029 F-I-002..006 close-or-narrow | conditional | 4 | `03-audit/A-030-r1-independent-after-e025-e026-c2-freeze-candidates.md` |
| A-031 | 2026-09-20 | self | response to A-030 / lock-predicate direction, two-source closure, disposition counts, jobs indexes, E3 alignment | conditional | 4 | `03-audit/A-031-r1-self-response-to-a030.md` |

## 结论状态

用户已完成 A-010/A-012 点名的关键方案裁决。A-014 independent 接受 A-013 对 **F-I-014** 的 `fixed`；A-016 independent 接受 A-015 对 **F-I-015** 碰撞的 `fixed`；A-018 independent 确认 Root D-012/D-013 方向已唯一但发现 matrix 表达式不一致；A-020 independent 接受 A-019 已把三份 C2 载体收成同一 `date_trunc`+整数 interval（F-I-002 表达式子项 `fixed`）；A-022 independent 接受 A-021 对 **F-I-016** 的 `fixed`（当时 `02-execution.md` E-001～E-019 严格递增且路径/`id` 一致）；A-023 已响应并维持 F-I-002～F-I-006 open。A-025 independent 确认当时索引为 E-001～E-023 单调，且 A-024 的 E-020～E-023 proposed 草案可收窄但不可闭合五条 required。A-027 independent 确认 A-026 的 Root D-014（v73–v87 未发布 baseline）与 D-015（负瞬间 Go Truncate / 整数来源不经 typmod）为可核对的方向收窄，**仍不可闭合** F-I-002～006：逐列 codec、90 列 mapping、Port 调用点、唯一表范围+descriptor 名+checksum、exact SQL 仍缺；D-015「integer interval」与秒列 `to_timestamp(double)` 未收口；被接受 allocation 仍含 v74 ledger 歧义；A-026 后出现 E-017 双文件（F-I-017 recommended）与无限定 D-012 同号不同义（F-I-018 recommended）。Root D-012/D-013 无政策矛盾。A-028 已修正已接受 allocation 的 v74 范围、唯一 E-024 与 SQLite rollback/RecoveryPoint 散文；A-029 independent 接受 F-I-017 closed（唯一 E-024，不接受 A-028 把 owner overlap 写进本条），维持 F-I-018 open（冻结包仍无限定 D-012），新增 F-I-019（索引把 E-024 插在 E-016 与 E-017 之间），并确认 owner spec v74「ledger/reconcile」与 PG restore `<artifact>` 仍未消掉。F-I-010 planning 仍 closed。R1 仍处于证据收集阶段。A-006 接受 F-I-001 `fixed`；A-012/A-014/A-016/A-018/A-020/A-022/A-025/A-027/A-029 接受精度截断、config D0 0→NULL、`pg_dump -F c`/`pg_restore`、Port 仅 `CreateRecoveryPoint`、`schema_migrations` owner = `core.persistence` 为方向已选，leftover 列名表已列出，PG 式已同一，D-014 baseline 已接受，已接受 allocation 列号不相交，但确认 C2/C3 仍不可冻结：草案与未发布 baseline 不是实施证据。90 列、catalog 72、v1–v72 / v67–v72 扫描、`login_failures` 与 retired `records` 已处理。A-030 independent 接受 E-025/E-026 后 **F-I-003 closed**（90 列 old→new→read/write mapping + voucher 三桶预检）、**F-I-018 closed**（点名五份载体无无限定 D-012/D-013）、**F-I-019 closed**（`02-execution.md` E-001→E-026 严格递增）；F-I-002 的 D-015 字面子项 `fixed`（用户裁决 B，秒族保留 `to_timestamp(double)`）；F-I-005 的 owner-spec v74 同文子项 `fixed`。**仍不可闭合** F-I-002（exact SQLite rebuild DDL、`#34`/`#72/#73` 与 E3 模板/D-015 不完全同一）、F-I-004（本轮未触及 C3；转换合同引用了不存在的 boundary 文件）、F-I-005（无已记录 checksum、无可执行测试改写、双方言 checksum 二选一未选定）、F-I-006（readwrite spec 仍 proposed；`#5` 锁谓词 new SQL 方向写反；§5 none vs ORDER BY 散文不一致）。新增 recommended **F-I-020**（悬空 C3 引用 / readwrite 未 superseded / 无限定 D-011）。freeze-candidate 不是实施证据。**开放 required = F-I-002 / F-I-004 / F-I-005 / F-I-006（4 条）**。存在未合法闭合的 required findings 时，不得冻结 C2/C3、不得关闭本子目标、不得将 Root R1 标 completed、不得放行 R2。

**A-031（self 响应 A-030）备注**：A-030 判定 **F-I-003 / F-I-018 / F-I-019 closed**，D-015 字面与 owner-spec v74 两个子项 `fixed`，**开放 required 5 → 4**（F-I-002、F-I-004、F-I-005、F-I-006），并新增 recommended **F-I-020**。A-030 同时查出编排器在 exact SQL 表中引入的一处**实质缺陷**——`#5 users.locked_until` 的 new SQL 方向写反（会把锁语义翻转，把过期锁当已锁定）；A-031 已按 `users_repository.go:496-503` 对位基线改正，并一并修正：readwrite spec 标 `superseded`（两源收口）、§5 `none` vs ORDER BY 归属（`#22/#31/#33/#60/#67/#74/#79` 上移，`#62` 下移，新计数 5+3+2+15+21+44=90）、jobs 四索引 exact `CREATE INDEX` 补入、`#6` callsite 改为 `accounts.go:194-201` 并列入写 0、`#34` E3 与毫秒族（含 `date_trunc`）同一、`#72/#73` 负值单路径（只走 `m0` 预检 fail closed，USING 只处理 `=0→NULL` 与正值）、悬空 C3 引用显式标注未落盘、无限定 D-011 改为 **Root** D-011、删除陈旧「D-015 待落盘」句。**C2/C3 仍未冻结，R2 仍未放行**（4 条 required 未闭合）。F-I-005 的双方言 checksum 约定按 P-004 待用户裁决（A-031 §4 列出 A/B 选项并倾向沿用 v1–v72 单 checksum）。
