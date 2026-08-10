---
id: A-007-grok-records-retirement-surface
doc: audit-entry
goal: GOAL-005-r4-full-module-migration
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: Independent Records retirement surface review (claimed GOAL-007 handoff); apps/api + apps/web product surface, migrations, historical compatibility, test rename safety
audit_type: close-out
verdict: fail
---

# A-007 · Grok Records 退场运行面独立复审（GOAL-007 handoff）

## 声明

本意见 `source: independent`。不修改目标 `status` / `progress` / 方案正文 /
goal-tree 状态列。

用户请求对 `GOAL-007-r4-records-retirement-closure` 做独立退场复审。**该目标在
workspace-003 不存在**（无文件夹、无 goal-tree 行）。按 P-003，意见必须落在
被审目标台账；本轮将运行面核验证据与「能否关闭 GOAL-007」结论落在 handoff 源
目标 GOAL-005（E-015 声明 GOAL-007 承接）。**未**创建 GOAL-007，也**未**读取其他
工作区的同名 GOAL-007。

## 范围与区间

- 工作区：`workspace-003-modular-admin-architecture`
- 代码范围：`apps/api`、`apps/web`（含 e2e 与 uncommitted 测试命名泛化 diff）
- 成功问题：
  1. 是否仍有 Records **产品** API/handler/store/seed/manifest/专属 hook/专属 fixture
  2. 本轮测试命名泛化是否安全
  3. 是否误删 0003/0006、历史 `records.*` operation-log、通用
     `recordView`/`recordSource`/`RecordID`、负向防复活测试
  4. 是否可关闭 GOAL-007；required findings 是什么

## 成果（有证据）

### 产品面：Records 专属实现基本已退场

| 检查项 | 结果 | 证据 |
|--------|------|------|
| HTTP product handler `/api/records` | **未发现** | `handler.Register` 仅挂 users/roles + schema 等；无 records handler 文件 |
| Store 产品 CRUD / seedRecords | **未发现** | `seed.go` 仅 users/roles/settings/activity 权限与菜单；无 records seed |
| Schema 专属 page fixture | **未发现** | `handler/fixtures/schema/` 为 users/roles/overview/…，无 records.json |
| Manifest / menu_records | **未发现** | web public manifest 无 records；seed 无 menu-records |
| 专属 frontend hook / page | **未发现** | web src 无 `useRecords` / Records 产品页引用 |
| 0006 退场仍在 ledger | **保留** | `migrate.go` `0006:records-retire:v1` DROP records + 删 perm-records-* |
| 0003 历史 DDL 仍在 ledger | **保留** | `recordsPersistDDL` + `migrate0003` 注释标明 immutable legacy |
| 历史 `records.*` event CHECK | **保留** | 0004/0005/0008 CHECK 仍含 records.create/update/delete |
| 通用 `RecordID` / `record_id` | **保留** | `store.Operation.RecordID`；handler 输出 `recordId` |
| 通用 `recordView` | **保留** | `renderer/render.ts(x)` 白名单与 case |
| 通用 `recordSource` 协议约束 | **保留** | `request-construction.ts` RECORD_SOURCE_REF 校验 |
| 负向：fresh 后无 records 表 | **保留** | `migrate_test.go`、`operations_test.go` 断言 drop |
| 负向：未知 event `records.purge` | **保留** | `operations_test.go` CHECK 拒绝 |
| 历史行 `records.create` 可读写 | **保留** | `operations_test.go` legacy event + RecordID |
| dataSource 拒绝伪 records 路径 | **保留** | `resource.test.ts` `isValidDataSource("records")` / `"/api\\records"` → false |

### 本轮测试命名泛化（含工作区 uncommitted diff）

观察到的改动模式：

- 注释：`records page` → `retired demo page` / `legacy demo`
- fixture 符号：`RECORDS` → `SAMPLE_ROWS`，`recordsFetcher` → `rowsFetcher`
- 文案：`No records` → `No rows`，`records fetch failed` → `resource fetch failed`
- table id：`records-table` → `schema-table`
- 驱动面仍指向 **users** e2e / users dataSource，而非恢复 records API

**判定**：命名泛化**安全**。未删除 0003/0006、未删历史 event、未删通用
recordView/RecordID、未删负向断言；仅降低「演示实体仍是产品资源」的误导性。
`SAMPLE_ROWS` 内仍可有 `rec-1` 类假 id，属通用表格样例，不构成产品 API。

### 仍可见但不构成产品复活的引用

1. **协议 conformance 示例 URL**  
   `apps/web/src/protocol/conformance/stage3-fixtures.test.ts` 使用
   `url: "/api/records"` 作为 **action schema 形状** 合法样例，不是挂载产品
   handler。推荐澄清注释（F-IND-R4-REC-003），**不**等于 CRUD 复活。
2. **迁移/测试中的历史表名与 event 名**  
   符合 D-003「保留 0003/0006 与历史 records.*」约束。
3. **`apps/api/data/*.pre-v0006*.sqlite` 快照**  
   迁移恢复 backstop 数据文件，不是产品 API。

## 对照「能否关闭 GOAL-007」

| 问题 | 判定 |
|------|------|
| GOAL-007 是否存在且可关门 | **否** — 目标未建立 |
| 若仅问「当前代码是否已无 Records 产品面」 | **基本是** — 见上表 |
| 是否已有 GOAL-007 五件套 + self + independent 关门台账 | **否** |
| 是否可把 E-015「GOAL-007 承接」视为已完成 | **否** |

## Findings

### F-IND-R4-REC-001 · 声称的 GOAL-007 在 workspace-003 不存在

- **level**: `required`
- **severity**: high
- **status**: `open`
- **impact**: 不能关闭 GOAL-007；不能把 R4-I003「运行面核验」记为已由子目标完成；
  handoff 悬空
- **evidence**:
  - `docs/workspaces/workspace-003-modular-admin-architecture/` 无
    `GOAL-007-r4-records-retirement-closure/`；
  - `goal-tree.md` 仅至 GOAL-006；
  - D-003 / E-015 / GOAL-005 meta 均引用 GOAL-007。
- **closure**: `/govern` 按五件套建立 GOAL-007 并挂 parent GOAL-005（或明确改写
  handoff，把核验收进 GOAL-005/006 并删除悬空引用）；完成后再做 close-out 审计。

### F-IND-R4-REC-002 · 退场核验目标缺失导致「可关闭」主张失败

- **level**: `required`
- **severity**: high
- **status**: `open`
- **impact**: 任何「GOAL-007 done / Records 退场已关门」声明
- **evidence**: 无 00-meta / 03-audit 可关闭；本意见只能证明 **代码切片**，不能
  替代目标关门。
- **closure**: 先有合法目标与成功标准，再以本 A-007 代码结论作 evidence 附件响应。

### F-IND-R4-REC-003 · stage3 协议 fixture 仍用 `/api/records` 作合法 URL 样例

- **level**: `recommended`
- **severity**: low
- **status**: `open`
- **evidence**: `stage3-fixtures.test.ts` L175/L191
- **closure**: 改为 `/api/users` 等现行资源，或注释标明「协议形状样例，非产品
  资源仍存在」。

### F-IND-R4-REC-004 · 缺少「请求 /api/records → 404/未注册」的显式防复活 HTTP 测试

- **level**: `recommended`
- **severity**: med
- **status**: `open`
- **evidence**: 负向覆盖主要在 migration（表不存在）与 dataSource 校验；未找到
  mux 级 `GET /api/records` 必须失败的专用测试。
- **closure**: 在 API 集成测试增加 retired route 负向用例，加固防复活。

### F-IND-R4-REC-005 · README 等文档仍用跨工作区 GOAL-007 叙事指 records 契约

- **level**: `recommended`
- **severity**: low
- **status**: `open`
- **evidence**: `apps/api/README.md` 提及 records 退场与「历史契约见 GOAL-007
  I-007-001」（语境像其他工作区历史编号，易与本区未建 GOAL-007 混淆）。
- **closure**: 文档改为 Q2 路径或「historical 0006」表述，避免裸 GOAL-007 误指。

## 必改项汇总

1. **F-IND-R4-REC-001 / 002**：不得关闭不存在的 GOAL-007；先建目标或改 handoff。
2. 代码产品面 **无** 必须立刻回滚的 Records CRUD 复活证据；推荐补 REC-003/004。

## 与既有意见 / 决策的关系

- D-003 historical-only 与本轮代码切片 **一致**。
- E-015「运行面已删产品 handler… GOAL-007 承接关门审计」：前半（产品面删除）
  **大体属实**；后半（GOAL-007 承接）**未落地**。
- 本意见 **不** 放行 GOAL-005 C1/C2，也 **不** 关闭 R4-I003 的运行面门禁；
  仅提供独立代码核验意见。

## 结论

**verdict: fail**（相对用户问题「是否可关闭 GOAL-007」）

- **不能关闭 GOAL-007**：目标在 workspace-003 不存在，台账与 goal-tree 均无。
- **代码层 Records 产品面退场**（API/handler/store/seed/manifest/专属 hook/专属
  fixture）在本轮抽查下 **成立**；0003/0006、历史 `records.*`、通用
  recordView/recordSource/RecordID、既有负向迁移测试 **未被误删**；本轮测试命名
  泛化 **安全**。
- required findings：**F-IND-R4-REC-001**、**F-IND-R4-REC-002**（目标缺失导致无法关门）。

## 建议给编排器 / 用户的下一步

1. `/govern` 决定：新建 `GOAL-007-r4-records-retirement-closure` 并挂 GOAL-005，
   或取消该 id 引用、把退场核验正式收进 GOAL-005/006。
2. 将本 A-007 代码表作为新建目标的 evidence 输入；补 recommended HTTP 404 与
   fixture 澄清。
3. 与 GOAL-006 A-004 一并响应：C1 最终冻结仍受 A-004 required 阻断。

本意见不修改 status/progress；响应由 `/govern` 处理。
