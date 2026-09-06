---
doc_type: goal-execution
id: E-012-offer-price-yuan-input
parent: GOAL-001-digital-offer-entitlement
date: 2026-09-06
status: done
version: 1.0.0
---

# E-012 · 关门后维护：Offer 标价改为元（两位小数）输入

## 事实（时间线）

- 2026-09-06 · **用户意见**：新建 Offer 的标价按【元】（最多两位小数）输入更符合直觉，而非分。
- 2026-09-06 · **方案取舍**：
  - 参照 E-008 预付凭证先例（表单输入元、服务端 `parseYuanToCents` 转分），但数字 Offer 有**编辑回填**路径：编辑弹窗从表格行预填 `priceAmount`（分），若只改服务端入参则回填会显示「990」并被再次 ×100，造成数据翻倍。
  - 故采用**客户端显示单位**方案：`inputNumber` 新增 Host-local 可选 `unit: "yuan"`（非协议控制属性，与 `afterComponent` 同类的本地扩展）。显示/输入为元（step 0.01、提交四舍五入到整分），**wire 值始终为整分**——创建、编辑、列表（`format: currency`）、存储、购买/目录流程与冻结契约（D-002 priceAmount 最小货币单位）全部不变，无需改 API。
- 2026-09-06 · **实施**：
  - `apps/web/src/renderer/form-controls.types.ts`：`FormControlField.unit`（仅 inputNumber，"yuan"）。
  - `apps/web/src/renderer/form-controls.tsx`：`NumberField` 支持 `unit`——显示 `/100`、提交 `Math.round(元×100)`（防浮点漂移）、step 缺省 0.01、min/max 按同尺度换算；仅 `unit==="yuan"` 生效，其它 inputNumber 零行为变化。
  - `apps/api/modules/digitaloffer/schema/digitaloffer-offers.json`：创建与编辑弹窗的 `priceAmount` 增加 `"unit":"yuan"`、`"step":0.01`。
  - i18n：zh `标价（元，最多两位小数）`；en `Price (yuan, up to 2 decimals)`。
- 2026-09-06 · **验证**：
  - Ajv（冻结 schema）对 offers/entitlements 文档校验 VALID（`unit` 为 props 业务键，不违反约束）。
  - 新增渲染测试：wire 990 显示为 9.9、step 0.01；输入 12.5 提交 1250（分）。
  - web vitest 全包 92 文件 / 1219 测试通过；`tsc -b` 通过；API digitaloffer/handler/composition 相关测试通过（无 API 改动）。

## 产物路径

- `apps/web/src/renderer/form-controls.types.ts`、`form-controls.tsx`、`visual-fidelity.test.tsx`
- `apps/api/modules/digitaloffer/schema/digitaloffer-offers.json`
- `apps/web/src/i18n/messages/zh-CN.json`、`en-US.json`

## 进度评估

- 关门状态不变（Root done 4/4 · VP-031 closed v0.3.4）；低风险可逆维护，审计模式 none。
- 说明：`type="number"` 输入框不强制显示尾随 "9.90"（浏览器归一化），精度以 step 0.01 + 四舍五入保证「最多两位小数」，与预付凭证金额输入行为一致。若需强制定长两位显示可后续跟进。
