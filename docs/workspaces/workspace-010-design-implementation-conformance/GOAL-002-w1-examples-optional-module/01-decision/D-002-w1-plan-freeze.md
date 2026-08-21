---
id: D-002-w1-plan-freeze
doc: decision-entry
goal: GOAL-002-w1-examples-optional-module
status: accepted
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# D-002 · W1 方案冻结：范例面 dev.examples 可选模块化

## 决定

1. **模块 id = `dev.examples`**（I-001 用户确认）：加入 compiled 候选集；`mvp`/`admin` 默认 Profile **不含**。
2. **homePageRef 策略 = 首个启用的 admin 功能页**（I-002 用户确认）：Manifest 装配期按启用模块集**确定性推导**；`dev.examples` 启用时可经 fragment 覆写回 `overview`。推导优先级与「无 admin 功能页可用」的 fallback 在拆分步钉死。
3. **mvp/admin 默认关闭演示模块**（I-003 用户确认）：生产向 Profile 不暴露演示产品面；启用路径 = `APP_MODULES_ENABLED` 显式列表（未来可选 `dogfood` profile）。
4. **目标形态**：
   - 8 个范例 pageId（`overview`/`data-table`/`admin-list-batch`/`data-display`/`search-form-table`/`form-controls`/`form-with-reactions`/`form-with-upload`）与 Examples 导航，从 `core.schema-render` 贡献与 manifest baseline 拆出，归 `dev.examples`。
   - `core.schema-render` 保留 CapabilitySchema / CapabilityValidation；`admin.*` 与 `core.manifest-route` **永不 DependsOn `dev.examples`**。
   - 组合根按 `plan.HasModule` 装配：`core.schema-render` 与 `dev.examples` 均作条件 provider。
   - manifest baseline 不再含 Examples 组与 8 范例 pageId。
5. **测试分母调整**：更新假定范例在默认集 / `homePageRef=overview` 的既有测试；新增断言——`dev.examples` 禁用时 Manifest 无 Examples 组与 8 范例 pageId、schema 404；启用时恢复。
6. **对 VP-008 `go` 的影响**：本波改变 Profile 默认集、模块矩阵与 Manifest 装配语义 → 按 VP-008 §`go` 消费有效性**触发暂挂/重验证路径**；W1 回归证据落盘后由 `/govern` 留痕恢复或用户 P-004 裁决。

## 为什么

- 用户 2026-08-11 书面确认 I-001～I-003 三项方案参数。  
- 符合 VP-010 方向级范围「Profile 与产品包装」「Manifest/Schema/导航聚合」「模块贡献符合性」；纠正 G1–G4 已核实的 as-built 偏差。  
- 动态「首个启用的 admin 页」比固定 `users` 对缺 `admin.users` 的 custom profile 更通用（装配期即可判定）。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| `core.examples` / 无前缀 `examples` | 保留 core 前缀或裸名弱化「dev-only 演示」语义，与薄内核及命名空间区分不符 |
| homePageRef 固定 `users` | 用户选择「首个启用的 admin 页」；固定值在缺 users 的 custom profile 下失效 |
| 默认保留演示、仅提供可关能力 | 不满足 S4 产品面卫生目标；生产 Profile 仍暴露演示面 |
| 物理裁剪 / 删除全部范例页 | 损害协议 dogfood 与回归分母；可选化而非消灭 |

## 影响与后续

- 信息项：I-001 / I-002 / I-003 → `verified`（用户确认）。
- 下一步（roadmap 阶段 2/3）：**拆分与迁移**（kernel BuiltinModules、composition 装配、manifest baseline 与 fragment、schema 归属、web 代表路径）→ **回归**（composition / manifest / profile 测试 + 双 Profile 烟测）。
- 审计模式：触及模块矩阵 / Manifest 装配语义 → 倾向 `cross`（self + independent）；provider 与触发时机于实施前按 P-004 定。
