---
id: GOAL-041-w29-api-web-protocol-conformance
doc: decision-entry
record_id: D-003
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# D-003 · S3 决策 — C-009 custom 边界固定（用户 P-004 书面裁决）

## 触发

S2 分类将 C-009 定为 `custom-extension-candidate`（15 个 custom 键无统一 namespace，D-002 §2）；按 D-001 §4 与 00-meta S3，custom 候选必须**逐项向用户展示「推动上游增补 / 本仓 custom / explicitly-out」互斥方案与风险，经用户书面裁决后**才能固定边界。S3 开启即执行该裁决。

## 用户裁决（P-004 · 2026-09-06）

| 裁决项 | 选项 | 结果 |
|--------|------|------|
| C-009 处置路径 | 本仓合法 custom（推荐）/ 推动上游增补 / explicitly-out | **本仓合法 custom** |
| namespace 策略 | 保留现有 15 键 + 登记所有权 + 新键规范（推荐）/ 立即批量重命名 | **保留 + 登记 + 新键规范** |
| C-005 子项（digitaloffer 未使用能力声明） | 删除（推荐）/ 保留 | **删除** |

## 决定

1. **处置路径 = 本仓合法 custom**：依据上游 `08-renderer-spec.md §1.1` Host Extension 模型（显式注册、非核心协议、宿主私有）；15 键均为本应用业务域控件，不属于上游通用组件语义。推动上游增补会污染通用协议注册表（拒绝）；explicitly-out 与 schema 内 `type: custom` 声明矛盾（不选）。
2. **custom 边界固定**：以 `attachments/custom-extension-boundary.md` 为权威规范，固定 namespace（§3.2）、capability（§3.3）、schema/validator（§3.4）、failure（§3.5）、compatibility（§3.6）、fixtures（§3.7）、退出/迁移触发（§3.8）。现有 15 键保留；新键按 `<module-scope>-<feature>` 规范 + 防碰撞规则；不新增协议 capability。
3. **C-005 子项 = 删除**：S4 实施时移除 `digitaloffer-entitlements.json` 的 `data.route-binding` 与 `digitaloffer-offers.json` 的 `form.controls.readonly` 声明（页面无对应使用），并补「声明 ⊆ 使用或文档化意图」守卫断言。
4. **未选方案**：
   - 推动上游增补（业务域控件入通用协议，周期长、定位冲突）；
   - 批量重命名 15 键（破坏性变更，兼容成本高，留作未来独立迁移）；
   - 保留未使用能力声明（用户已选删除，避免声明-使用漂移）。
5. **联动**：C-010 未知 custom 呈现对齐（S4）与 C-011 动作层白名单均在 §3.5 边界内；F-001（页面级能力门禁 + claim/HOST_SUPPORT）不属本边界，S4 另行实施；上游增补报告不创建（I-004 不适用）。

## 影响

- **I-005**（custom 门禁）：`open` → **verified**（用户书面裁决 + 边界规范落盘，D-003 + `attachments/custom-extension-boundary.md`）。
- S4 custom 相关实施项：C-010（failure 对齐）、C-005 子项（digitaloffer 声明删除 + 守卫）、新键命名规范执行。
- S6 关门前按 meta 对 custom 边界复审（independent 腿）。
