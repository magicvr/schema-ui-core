---
id: GOAL-013-w12-product-surface-intent
doc: audit-entry
record_id: A-002
source: independent
scope: S3 实施 ～ S4 验证/关门
verdict: conditional
status: recorded
auditor: grok build（grok-4.6 · high）
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# A-002 · S3/S4 独立交叉审计（2026-08-16）

- **source**：independent
- **auditor**：grok build（grok-4.6 · high）
- **类型** / **scope**：close-out · S3 实施（T-05 / T-01 / T-03 / T-02 / T-06）→ S4 验证/关门
- **verdict**：**conditional**
- **工作区**：`workspace-010-design-implementation-conformance`（`root_goal` = `GOAL-001-design-implementation-conformance`；`canonical_scope` 已核对；`shared_materials_catalog: none`；`primary_plan` = `VP-010-design-implementation-conformance`）

## 范围与区间

- **covered**：D-002～D-008 冻结方案 vs as-built；E-002～E-004（及已落盘的 E-005）主张 vs 代码；P-005 I-001～I-006；独立复跑 `apps/api` `go test ./... -count=1` 与 `apps/web` `npm test -- --run`。
- **不 covered**：T-04 / [workspace-011] GOAL-022 实施（D-005 本波不做）；浏览器活栈点验；其他工作区上下文。
- **审计时点**：本文件写就前，`03-audit/A-002-*.md` **不存在**，但 `00-meta.md` 已 `status: done`、S4 检查点已勾选「independent A-002」、`goal-tree.md` / `workspace.md` 波次表已写「A-002 grok independent pass」、`03-audit.md` 索引已预登记 A-002 pass。本意见不承认该预写 pass。

## 成果（有证据）

| 主张 | 独立核对 |
|------|----------|
| T-05：`deletedAt` / `restoredAt` 为 UTC ISO-8601 | `recyclebin.go` `recycleItemToMap` 使用 `2006-01-02T15:04:05.000Z07:00`；`recyclebin_test.go` `time.Parse` 断言。存储层仍为 Unix 秒，仅 HTTP 投影变更。符合 D-006。 |
| T-01：全断点单一 `UserMenu`；抽屉不再重复用户链 | `App.tsx` `UserMenu` 常驻顶栏（非 `hidden lg:flex` 横铺）；菜单 = `projection.user` 声明序 + 分隔 + 退出；抽屉只渲染 `top + sidebar`。`user-menu.test.tsx` 4 例在全量套件中通过。符合 D-002 §2–§3 的可达性主线。 |
| T-03：个人中心三档 Tabs + `labelKey` | `account.json` `type: tabs`，资料 / 安全（改密+MFA）/ 会话；`TabsView` `resolveTextProp("labelKey","label")`。en/zh `schema.account.tab.*` 已登记。符合 D-004。 |
| T-02：`searchFormSubmit` 非 `q` 字段入 `filters` | `render.tsx` owned 字段先删后写；`buildResourceQuery` 序列化为顶层 query。`search-form-filters.test.tsx` 断言 `q=ali&enabled=true`（全量套件中通过）。 |
| T-02：12 页矩阵 + ExtraQuery 接线 | users `enabled`/`locked`；roles `system`；scheduled-tasks `enabled`；task-runs `status`；notifications `q`+`read`（handler 直读 query，非 ExtraQuery 工厂，语义等价）；wallet `ownerType`、recycle-bin `resource` 走既有 query。wallet-entries **无** `entryType` select（E-003 降级留痕，符合 D-003）。数据权限 / 系统监控无 search form（符合「不做」）。 |
| T-06：模块启用只认 YAML | `config.go` 不再 `envOr("APP_PROFILE")` / `APP_MODULES_ENABLED`；`resolveModulesFromYAML` preset/list 互斥；`loadPresetFile` KnownFields + 空列表 fail closed；`config_test.go` 含「旧 env 被忽略」。`compose.yaml` 移除模块 env、只读挂载 `./apps/api/configs`；密钥仍走 `AUTH_JWT_SECRET` / `ADMIN_INITIAL_PASSWORD`（D-008）。 |
| T-06：默认三档成员未变（go 不暂挂依据） | `kernel/profile.go` `profileDefaults` 仍为既有 mvp / admin / demo 集；`TestBuiltinProfilesResolveDeterministically` 锁定 mvp/demo 精确列表且本轮 Go 全量绿。Manifest 装配与模块矩阵未见本波改写。部署契约变、默认集不变 → 同意 E-005「不暂挂」判定。 |
| I-001～I-006 | 均有 D-002～D-008 书面闭合；无到期未验的 required 信息项。I-005 non-blocking 已随 D-002 闭合。 |
| 独立回归（Go） | `apps/api` `go test ./... -count=1`：**0 FAIL**（2026-08-16，本会话）。 |
| 独立回归（Web） | 全量 `npm test -- --run`：一次结果 **1025 passed / 2 failed / 1027**；失败均在 `s5-denominator-render.test.tsx`（timeout + settings 仍「正在加载」）。同文件单独复跑 **5/5 通过**。T-01/T-02 新测在全量中通过。 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S3 P0 T-05 + T-01 | 达成 | E-002 + 上表代码/测试 |
| S3 P1 T-03 + T-02 | 达成（有 recommended 标签精度缺口） | E-003 + schema/handler；见 F-003 |
| S3 P2 T-06 | 达成 | E-004 + config/kernel/compose；go 判定可核对 |
| S4 回归绿 | **部分** | Go 可重复绿；Web 全量一次未复现 1027/1027，单文件可绿 |
| S4 自审 + independent 后关门 | **未合法完成** | A-001 存在；本 A-002 此时才落盘。预写 `done` / 「A-002 pass」不是合法关门（F-001） |
| T-04 不阻塞本波 | 达成 | D-005；本波无钱包自服务页 |
| P-005 关门信息项 | 达成 | I-001～I-006 verified |

## Findings

### F-001 · 关门台账先于独立意见落盘

- 严重度：**med**
- 建议：**required**
- 状态：open
- 描述：P-003 要求独立意见完成前不得将目标标为 `done`，且不得用预写 pass 代替本意见。本审计写就时：
  - `00-meta.md`：`status: done`；S4 检查点已勾选「independent A-002（grok）」；frontmatter `progress: 2/4` 与正文「当前 4/4」互相矛盾。
  - `goal-tree.md` / `workspace.md`：W12 已写成「关门 4/4」且「A-002 grok independent pass」。
  - `03-audit.md`：在 `03-audit/A-002-*.md` 尚不存在时已登记 A-002 **pass** / 开放 required「无」。
  - `02-execution/E-005` 将「A-002 核销」写成已发生事实。
- 证据路径：上列文件 vs 本目录当时仅有 `A-001-s3s4-self.md`。
- 要求 `/govern`：按本意见真实 verdict 改写 `03-audit.md` 索引行；在响应本意见并合法闭合 required 之前，**不得维持** `status: done` 或把预写 pass 当关门证据。本独立审不改 status / progress / goal-tree。

### F-002 · Web 全量 1027/1027 未能在独立一次复跑中复现

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：E-005 / A-001 主张 `npx vitest run` **1027/1027**。本会话全量一次为 **1025/1027**，两失败均在既有 `s5-denominator-render.test.tsx`（union titleKey 5s timeout；admin `/settings` 仍处「正在加载页面 Schema」，未匹配 `/常规|站点标题/`）。同文件单独复跑 5/5。形态像并行负载 flake，**不是** T-01/T-02/T-03 新测失败，也未证明 T-06 改了默认集。但「全量已绿」不能当作不可重复关门证据。
- 证据：本会话全量日志；随后 `npm test -- --run src/i18n/s5-denominator-render.test.tsx` 5/5。

### F-003 · users 关键词标签写了 ID，后端不搜 `id`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-003 矩阵与 i18n `schema.users.search.q`（「用户名 / 显示名 / ID」）声明搜 ID。`usersWhere` 仅 `username` / `name` 的 `instr`，不含 `id`。与「label 必须写清搜什么」不完全一致。未改 Extra 接线（enabled/locked 已落地）。
- 证据：`apps/web/src/i18n/messages/zh-CN.json`；`authsession/users_repository.go` `usersWhere`。

### F-004 · T-01 窄屏隐藏显示名（可达性主线仍成立）

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-002 §1 写「头像圆标 + 显示名 + 小箭头；全断点同一控件」。实现是同一 `UserMenu`（符合「同一控件」），但显示名 `hidden … sm:inline`，`<sm` 只见头像+箭头。触发器有 `aria-label`（`shell.userMenu`）、`min-h-9`，抽屉已去用户链，移动端可从顶栏打开个人中心/设置/退出。不构成 D-002 可达性失败。
- 证据：`App.tsx` `UserMenu` 触发器 className。

### F-005 · wallet-entries `entryType` 按 D-003 降级（与 A-001 R-001 同）

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-003 允许 list 无谓词则仅关键词。`wallet-entries.json` 仅 `q`；`ListFilter` 无 entryType。E-003 已留痕。后续账本 list 加谓词后再挂即可，无需本波重开目标。

## 必改项汇总

| ID | 级别 | 要做什么 |
|----|------|----------|
| F-001 | **required** | `/govern` 纠正预写关门：索引 verdict 以本文件为准；响应本意见前不得维持 `done` / 「A-002 pass」。 |

无其他 required。F-002～F-005 为 recommended，不单独阻断 S3 主张。

## 与既有意见的异同

- **A-001 self（pass）**：同意其对 T-05/T-01/T-03/T-02/T-06 代码主张与 I-001～I-006 闭合、以及 T-06「默认集不变 → 不暂挂」。同意 R-001（本文件 F-005）。
- **分歧**：A-001 把 S4 写成「等待 A-002 后关门」。编排器随后在 A-002 文件不存在时勾选 S4 并标 `done`、预写本条 pass。本意见 **不** 追认该闭环；verdict 为 conditional，不是 pass。
- 无与 A-001 在产品实现结论上的冲突需要 P-004 裁断。F-001 是过程门禁，由编排器响应，不必用户在产品方案上另裁。

## 结论 + 建议给编排器/用户的下一步

S3 五条实施（T-05/T-01/T-03/T-02/T-06）与 D-002～D-008 **可核对**；T-06 部署契约变化与默认三档不变一致，**同意不暂挂** go 消费。独立 Go 全量绿。Web 新测绿；全量一次有 2 个既有 i18n 套件 flake，单文件可复现绿。

**不能无条件关门**：存在未闭合 required F-001（预写独立 pass + 提前 `done`）。

建议 `/govern` 下一句：响应 A-002；按 F-001 修正台账（索引 verdict=conditional；在闭合 F-001 前撤回/不得维持 `done`）；可选再跑一次 Web 全量或将 F-002 接受为 flake residual；F-003～F-005 可本波 overruled / 后续波次。

## 声明

本意见不修改 status / progress / goal-tree / 方案正文 / 业务代码。响应由 `/govern` 处理。
