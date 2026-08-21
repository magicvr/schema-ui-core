---
id: A-003
goal: GOAL-013-nav-order-config
source: independent
auditor: grok-build
date: 2026-08-14
scope: S1~S4 全量记录 + 实现代码（关门前独立审计）
audit_type: close-out
verdict: conditional
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
status: active
---

# A-003 · independent 审计（S5 关门前 · 导航顺序方案 A）

## 结论

**verdict: conditional**

- **UI 真实路径（manifest 聚合排序）**与 D-002 冻结方案一致，证据充分，可支撑产品级导航顺序目标。
- **kernel → system-data `menu_items.sort_order` 路径未真正接通**：`sortNavigation` 只重排内存切片，Reconcile 仍写入模块声明的 `Order` 且 `ON CONFLICT DO NOTHING`，与 E-003/A-002「两层都加排序 / kernel 驱动系统数据顺序」主张不符。
- 在 F-001 以 `fixed` 或 `accepted-residual` / `user-overruled` 合法闭合前，**不宜无条件关门**。

## 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-011-admin-functional-modules`（Root `GOAL-001-admin-functional-modules`；`shared_materials_catalog: none`） |
| 目标 | `GOAL-013-nav-order-config` |
| 记录 | D-001/D-002；E-001~E-005；A-001/A-002（self） |
| 代码 | kernel / manifest / config / composition + 对应测试；`configs/config.yaml` |
| 信息项 | I-001~I-004 均已标 closed（与 S1/S4 记录一致）；本审未发现新增 required 信息门禁 |
| 共享资料 | 无 |

本意见 **不** 修改 `status` / `progress` / goal-tree。

## 成果（有证据）

### 1. 默认清单 12 项与用户确认一致，快照锁定

D-002 §2 与用户确认顺位：

Dashboard → Users → Roles → Settings → Activity → Account → Notifications → File library → Data dictionary → System monitoring → Scheduled tasks → Recycle bin

与 `kernel.DefaultNavigationOrder`（`provider.go:401-414`）及 `TestDefaultNavigationOrderSnapshot`（`navigation_order_test.go:11-28`）**逐项一致**。

### 2. 排序语义（清单优先 / 未列追加 / Parent 优先）— kernel 与 manifest 单元层正确

| 层 | 行为 | 证据 |
|----|------|------|
| kernel `sortNavigation` | 清单 rank → 未列 Order/NodeID 兜底；Parent 分组优先 | `provider.go:422-450`；`TestSortNavigationDefaultOrder` / `PartialCustomOrderAppendsRest` / `CustomOrder` |
| manifest `SortNavigation` | 按槽（top/sidebar/user）独立重排；未列 stable 追加 | `manifest.go:222-255`；`TestSortNavigationOrdersSlots` |
| config | 空/非法 YAML → nil（默认）；env 覆盖 | `config.go:220-236,283-301`；`TestLoadNavigationOrder` 5 子用例 |

非法覆盖整体回退：`NormalizeNavigationOrder`（`provider.go:475-489`）unknown NodeID → `DefaultNavigationOrder` + WARN；kernel `resolveNavigationOrder` 与 composition 构建 manifest 前共用（`composition.go:336-337`）。E-004 所述「plan 原样传 manifest」绕过已修复。

### 3. 覆盖载体与 W7 优先级链一致

- YAML：`navigation.order`（`configs/config.yaml:80-84`）
- env：`NAVIGATION_ORDER` 逗号分隔；`TrimSpace` 后非空才覆盖 YAML（`config.go:225-236`）
- 宽松解析：非 Sequence / 非 `!!str` 项 → nil + 告警（`parseNavigationOrder`）
- 空列表 / 缺省 → nil → kernel 默认

### 4. manifest 为 UI 权威路径 — 实测与默认清单一致

独立探针（审计会话内临时 composition 测试，已删除，不入库）：

- `GET /.well-known/schema-ui/app-manifest.json`（admin profile）
- sidebar：`Dashboard | Users | Roles | Activity | File library | Data dictionary | System monitoring | Scheduled tasks | Recycle bin`
- user：`Settings | Account`
- 与 D-002 / E-003 描述一致（Notifications 不在 UI 槽：见 F-003）

NodeID 提取：`features.([A-Za-z0-9_]+)` → id/pageRef 兜底（`manifest.go:188-216`）。现有 admin fragment 均使用 `$context.features.menu_*`，与 kernel `NodeID` 对齐。槽独立排序，不混 top/sidebar/user。

### 5. 向后兼容与回归

- 无覆盖时 = 默认清单（产品冻结意图，非字母序契约）— 与 D-002 / E-005 一致
- 相关单测包本审运行：`kernel` / `config` / `manifest` / `composition` **全绿**
- web fixture 旧序：A-002 已注；属静态测试数据，非 API 输出

### 6. go（VP-008）判定

E-005 表：装配顺序 / Profile 矩阵 / 协议形状 / 权限集合 / 门禁语义均不变；仅导航槽内条目顺序可配置。对照实现：`Plan.NavigationOrder` 附加字段 + `ForModulesWithFragments` 可选变参 + config 段 — **同意 go 不 held**。未见遗漏的装配/门禁语义变化。

### 7. 边界（代码路径核对）

| 边界 | 行为 | 判定 |
|------|------|------|
| 空 order `[]` | nil → 默认 | 有测 |
| 单条目 / 部分清单 | 未列追加末尾 | kernel+manifest 有测 |
| 重复 NodeID | rank map 后者覆盖；未专项测试 | 见 F-004 |
| NodeID 大小写 | 大小写敏感；不匹配 → 整表回退默认 | 见 F-004 |
| `NAVIGATION_ORDER=""` | 视为未设置，YAML 生效 | 代码正确 |
| YAML + env 同时 | env 胜 | 有测 |
| 非法 env `menu_bogus` | Normalize 回退 | kernel 有测；E-004 进程实测 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 默认清单 + 覆盖载体冻结 | 达成 | D-002；A-001 |
| S2 默认清单 + 排序 + 配置 | **部分** | 常量/config/manifest 达成；kernel→DB 未接通（F-001） |
| S3 单测 + 三场景实测 | 达成（UI/解析路径） | E-004；本审复跑单测 |
| S4 go 不 held | 达成 | E-005；本审复核 |
| 产品：无覆盖 = 默认清单 UI 序 | 达成 | manifest 探针 |
| 产品：配置可覆盖且非法回退 | 达成（manifest 路径） | Normalize + config 测 + E-004 |

## Findings

### F-001 · kernel 排序未写入 `menu_items.sort_order`（双层主张不实）

- **严重度**：med
- **建议**：**required**
- **状态**：open
- **描述**：

  E-003 / A-002 主张排序在 **kernel（系统数据）+ manifest（UI）** 两层生效。实现上：

  1. `sortNavigation` 仅 `sort.Slice` 重排贡献切片，**不改写** `NavigationContribution.Order`（`provider.go:422-450`）。
  2. `Reconcile` 在写入前按 `ModuleID`/`Key` **再次排序**，丢弃切片序（`reconcile.go:38-43`）。
  3. `ensureNavigation` 将 **`n.Order`（模块声明的旧 Order）** 写入 `sort_order`，且 `ON CONFLICT(id) DO NOTHING` 永不按清单更新（`reconcile.go:204-208`）。
  4. 运行时读路径 `FeaturesForUser` 按 `feature_key` 投影 feature map，**不按 `sort_order` 驱动 UI**（`accounts.go:140-148`）。

  **审计探针证据**（admin 新库，12 行 `menu_items`）：多数 `sort_order` 仍为模块 Order（如 `menu_settings=1` 而非清单位 3；`menu_account=1` 而非 5；`menu_recycle_bin=8` 而非 11），**不等于** `DefaultNavigationOrder` 下标。同会话 manifest sidebar/user 序则与默认清单一致。

  **影响**：产品 UI 顺序正确（manifest）；若任何消费者依赖 `menu_items.sort_order` 表达产品默认序，则仍为冲突的旧 Order。双层完成的关门表述目前不成立。

  **闭合路径（择一，须书面留痕）**：

  | 路径 | 动作 |
  |------|------|
  | `fixed` | 在 Reconcile 前把清单 rank 写入 `Order`（或更新 `sort_order` 策略），并加集成测试锁定 admin `menu_items` 序；注意与「保留运营者手改 sort_order」既有测试意图协调 |
  | `accepted-residual` | 用户书面接受：产品权威序 = **仅 manifest**；kernel `sortNavigation` 不保证 system-data 序；`menu_items.sort_order` 仍为模块 Order / 运营者字段；复审触发 = 若出现读 `sort_order` 的导航消费者 |
  | `user-overruled` | 用户驳回本 finding（接受现状不等价于 residual 范围说明） |

### F-002 · composition 缺少「发布 manifest 序」集成锁定

- **严重度**：low
- **建议**：non-blocking（recommended）
- **状态**：open
- **描述**：`SortNavigation` 与 kernel 排序均有单元测；E-004 为进程手测。仓库内无稳定 composition 测试断言 admin 发布文档 sidebar/user 序 = 默认清单（或 env 覆盖后的序）。回归依赖人工或间接测。建议补 1 条 `testMux` + `/.well-known/schema-ui/app-manifest.json` 断言。

### F-003 · `menu_notifications` 在默认清单中，但 manifest fragment 无导航槽条目

- **严重度**：low
- **建议**：non-blocking
- **状态**：open
- **描述**：`admin.notifications` 向 kernel 注册 `menu_notifications`（`notifications/provider.go:72-81`），fragment 的 `navigation.user` 为空数组（`notifications/manifest/fragment.json:22-23`）。故 UI 不出现 Notifications 项，清单第 7 位仅影响 system-data 贡献存在性。属既有产品形态，非 GOAL-013 引入；维护清单时知悉即可。

### F-004 · 重复 NodeID / 大小写边界无测试

- **严重度**：low
- **建议**：non-blocking
- **状态**：open
- **描述**：重复 id 时 rank 后者覆盖；大小写不匹配整表回退。行为可接受且与「非法回退」一致，但无用例锁定。运维文档可注明 NodeID 必须精确匹配。

### F-005 · web 静态 fixture 仍为旧序

- **严重度**：low
- **建议**：non-blocking
- **状态**：open（与 A-002 备注同旨）
- **描述**：静态 fixture 非 API 输出；vitest 绿不证明与默认清单对齐。需要时另任务同步。

## 必改项汇总

| ID | 级别 | 一句话 |
|----|------|--------|
| **F-001** | **required** | 双层排序中 kernel→`menu_items.sort_order` 未接通；须 fixed 或 residual/overruled 后才能无条件关门 |
| F-002 | non-blocking | 建议补 manifest 序集成测试 |
| F-003 | non-blocking | Notifications 清单位 vs 空 fragment |
| F-004 | non-blocking | 重复/大小写边界测 |
| F-005 | non-blocking | web fixture 旧序 |

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-001 (self, S1, pass) | 同意方案冻结质量；本审不回退 S1 |
| A-002 (self, S2/S3, pass) | **分歧**：A-002 称「两层共用 Normalize、系统数据顺序」已就绪；本审证实 Normalize 双处统一属实，但 **system-data 落库序未跟随清单**。故升级为 conditional，新增 F-001 |

## 结论与建议下一步

1. **不可**在 F-001 开放时将 GOAL-013 标 `done`（P-003 开放必改门禁）。
2. 建议 `/govern` 响应 A-003：
   - 优先与用户确认 F-001 走 **fixed** 还是 **accepted-residual（manifest-only 权威）**；
   - residual 时写清范围与复审触发，并修正 E-003/A-002 表述以免后续误读；
   - fixed 时补 Reconcile/`Order` 策略 + 集成测后再请复审。
3. F-002~F-005 不阻断关门（在 F-001 闭合后）。
4. I-001~I-004 无需重开；go 不 held 维持。

## 声明

本意见 `source: independent`，不修改目标 `status` / `progress` / goal-tree。响应、finding 闭合与是否关门由 **`/govern`** 处理。


---

## 编排器响应（2026-08-14）

### F-001 · accepted-residual（用户书面接受）

- 用户裁决：**accepted-residual**（ask_user_question，2026-08-14）。
- 范围：产品权威导航顺序 = **仅 manifest**（UI 渲染路径）；kernel `sortNavigation` 不保证 system-data 顺序；`menu_items.sort_order` 保持模块声明 Order，且维持「运营者手改不覆盖」既有契约（reconcile_test 断言）。
- 复审触发：出现任何读取 `menu_items.sort_order` 的导航消费者时，重开本 finding 并评估 fixed（rank 落库 + SystemDataVersion bump）。
- 文档修正：E-003/A-002 中「kernel 驱动 system-data 顺序」表述已改为「kernel 排序保证贡献切片确定性顺序（system-data 落库仍为模块声明 Order）」。

### F-002 · fixed（补 composition 集成断言）

- 新增 composition 测试：admin profile 下 `/.well-known/schema-ui/app-manifest.json` sidebar/user 槽顺序 = 默认清单（含 env 覆盖后顺序断言）。

### F-003 · accepted（既有产品形态，非本目标引入）

- `menu_notifications` 无 fragment 槽：保持现状；清单维护知悉。

### F-004 · accepted-residual（边界行为符合非法回退语义）

- 重复 NodeID（rank 后者覆盖）与大小写不匹配（整表回退默认 + WARN）行为符合 D-002；无专项测试锁定，运维文档注明 NodeID 精确匹配。

### F-005 · accepted（web 静态 fixture 非 API 输出）

- 保持现状；需要时另任务同步。
