---
id: GOAL-003-r2-timezone-semantics
doc: audit-entry
record_id: A-001
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## A-001 · R2 时区语义关门审计（2026-08-26）

- **source**：self
- **auditor**：govern 编排器（本会话）
- **类型 / scope**：close-out · GOAL-003（R2 时区语义）全量：合同 `GOAL-002 D-001` §2 / §4.2 一致性、L1～L4 解析、用户覆盖通道、站点默认接入、统一语义接线、越界守卫、证据可核对
- **verdict**：**pass**

### 范围与区间

2026-08-26 立项至本审计时点（C1～C4 done · progress 4/5）。检查点：C1 解析器、C2 用户级覆盖 UI、C3 站点默认接入、C4 统一语义接线（含输入面核对结论）、C5 关门。共享资料引用：`shared_materials_catalog: none`，无引用项。

### 成果（有证据）

1. **C1 解析器**：`apps/web/src/i18n/timezone.ts`（`resolveEffectiveTimezone` L1→L4；IANA 校验；单通道读写；`auto` 移除 key）+ 15 快测（E-002；commit `e6de919f`）。
2. **C2 用户覆盖 UI**：`apps/web/src/components/timezone-switcher.tsx` 挂载头部 `<LocaleSwitcher/>` 右侧（`App.tsx`）；选项 = auto + 常用 IANA 集；双语 catalog（`timezone.switcher.*`）；4 快测（E-003；commit `a4c9d048`）。
3. **C3 站点默认接入**：`runtime.tsx` `/api/branding` fetch 增读 `siteTimezone`（`fetchedSiteTimezone`）+ `siteTimezone` prop 测试缝；Localization tab 字段与 schema 未动（核对 `startup-config` 既有断言仍绿）。
4. **C4 统一语义**：`I18nState` 暴露 `timezone` / `timezonePreference` / `setTimezonePreference`；`formatDate` 默认注入生效时区、显式 `options.timeZone` 优先（快测覆盖）；**输入面结论**：renderer 日期控件为 date-only（YYYY-MM-DD 本地日语义，机器合同 §3.3），不涉 epoch/时区偏移；含时间（epoch）输入控件本波不存在，§2.3 未来引入时须按其实现（E-003 已文档化）。
5. **证据**：`npx vitest run` web 全量 **87 files / 1155 tests / 全绿零回归**；R2 新增快测 25 用例（C1 15 + C2 4 + C3/C4 6）。
6. **越界守卫**：本波代码改动仅 `apps/web/src/i18n/timezone.ts|runtime.tsx|messages/*|components/timezone-switcher*|app/App.tsx` 与快测；未改 API / DB DDL / 迁移台账 / Profile 默认集 / 模块矩阵 / Manifest / `docs/contracts/`（git 提交范围可核对）。

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. 生效时区解析符合 §2（L1→L4 + 降级不抛错） | ✅ | `timezone.ts` + `timezone.test.ts`（15 用例，敌意输入不抛错） |
| 2. 用户覆盖持久化单通道；auto=移除；登录/登出不清除 | ✅ | `writeStoredTimezone`（auto 移除）；runtime-timezone 持久化翻转用例；与 locale `schema-ui:*` 同通道先例一致 |
| 3. 站点默认消费（L3）；无配置 zh-CN+auto 可运行 | ✅ | `fetchedSiteTimezone` + C3 fetch 快测；L4 兜底用例 |
| 4. 展示/输入统一语义（双向无偏移） | ✅ | `formatDate` 默认生效时区 + 显式覆盖用例；输入面 = date-only 本地日（E-003 结论） |
| 5. 无越界（DDL/迁移/Profile 默认集/docs/contracts/多时区排程等） | ✅ | 提交范围 `e6de919f`/`a4c9d048`/`b1eeefa7` + 方案 D-001 越界节 |
| 6. 关门自审落盘且 open required = 0 | ✅（本条即审计；required = 0） | 必改项汇总为空 |

### Findings

- **F-001 · 含时间（epoch 交换）输入控件若在后续波次引入，须按合同 §2.3 实现统一语义**（low · recommended · open → 移交 R4 核对项）
  - 描述：本波 renderer 无含时间输入控件（日期控件为 date-only 本地日语义），§2.3 的「时间输入统一语义」尚未被任何控件消费；E-003 已留痕结论。R3/R4 若引入 datetime 输入控件（如 `datetime-local`、时间选择器），必须按 §2.3 实现并快测，否则视为合同违约。
  - 证据：`apps/web/src/renderer/form-controls.tsx`（`DateField`/`DateRangeField`）；合同 §2.3；E-003「输入面结论」。
- **F-002 · TIMEZONE_OPTIONS 常用集为可核对列表，扩展须留痕**（low · recommended · open → 随 R3 立项/关门跟踪）
  - 描述：头部菜单时区集（上海/东京/纽约/伦敦/UTC）按合同 §6「可核对、可扩展列表」落地；后续扩展（或按 locale 动态裁剪）须更新组件常量并留痕，避免两处事实漂移。
  - 证据：`apps/web/src/components/timezone-switcher.tsx`（`TIMEZONE_OPTIONS`）；合同 §6。

### 必改项汇总（required 列表）

无（0 条）。

### 结论 + 建议下一步

R2 时区语义交付完整且证据可核对；scope 内无 required 必改项，无到期 required 信息项（I-001/I-005 已 verified）；verdict **pass**。F-001/F-002 为 recommended，不阻断关门，随 R3/R4 跟踪。建议：用户确认 GOAL-003 `status: done`（R2 关门）→ 立项 GOAL-004（R3 数字/货币语义，承接 GOAL-002 F-001/F-002 与合同 §3/§4.3）。如需可在关门确认前追加本地 grok build（grok-4.6 · high）独立审（本审计为 medium-risk 前端实施，self 证据已足，非必需）。