---
id: GOAL-003-r2-timezone-semantics
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## D-001 · R2 实施方案（lead 方案冻结）

### 触发

GOAL-002（R1 合同冻结）已关门（A-001 self pass · 用户 2026-08-26 确认）。用户确认「直接推进 R2」。本方案把合同 §2 / §4.2 落为可实施设计；前置信息门禁 I-001 / I-005 均已 verified。

### 决定（方案）

1. **生效时区解析器**（新模块 `apps/web/src/i18n/timezone.ts`，与 `locale.ts` 同构）：
   - `TIMEZONE_STORAGE_KEY = "schema-ui:timezone"`；`readStoredTimezone()` / `writeStoredTimezone(pref)`（`"auto"` = 移除 key；隐私模式抛错兜底，同 `runtime.tsx` 先例）。
   - `resolveEffectiveTimezone({ stored, siteDefault, detect }) → string`：L1 用户覆盖（IANA 名，非 `"auto"` 时直接生效）→ L2 会话探测（`Intl.DateTimeFormat().resolvedOptions().timeZone`，空串跳过）→ L3 站点默认（`siteTimezone`，非空且非 `"auto"`）→ L4 `"auto"`。
   - `detectBrowserTimezone()`：纯探测函数（测试可注入）；无效 IANA 名不在此层判错（降级在消费层）。
   - 快测矩阵：L1 覆盖 / L2 探测 / L3 站点默认 / L4 兜底 / 空值跳过 / auto=移除。
2. **用户级覆盖 UI**（C2）：头部 locale 切换旁新增时区选择（`apps/web/src/components/timezone-switcher.tsx`，仿 `locale-switcher`）；选项 = `auto` + IANA 常用集（Asia/Shanghai、Asia/Tokyo、America/New_York、Europe/London、UTC 等，可编辑列表可后续扩展）；持久化走 `schema-ui:timezone`。
3. **站点默认接入**（C3）：`I18nProvider` 新增 `siteTimezone` 输入（`/api/branding.siteTimezone` 已有字段，`startup-config` 已投影）；解析时传入 L3。Localization tab 不改字段、不改 schema。
4. **统一语义接线**（C4）：`I18nProvider` 的 `formatDate` 默认注入生效时区（`formatDateImpl(value, locale, { timeZone })`）；时间输入解析组件（R2 范围内若涉及）使用同一生效时区；展示/输入双向无偏移由快测断言。
5. **越界**：全程不改 API / DB / DDL / Profile 默认集 / `docs/contracts/`。

### 为什么

- 解析器纯函数 + 注入式探测：可测、无副作用、与既有 `locale.ts`/`runtime.tsx` 模式一致（P-002 可验证事实）。
- 用户覆盖通道与 locale 完全同构：登录/登出不清除、`auto` 移除 key，避免第二套偏好语义。
- 站点默认走既有 `/api/branding` 投影，零 API 变更。
- 快测优先：每个检查点带可核对测试，关门审计可指回证据。

### 未选方案

- 服务端会话时区（请求头探测 + 响应头回写）：引入 API/中间件改动，超出合同「前端解析」落点；无用户裁决依据，弃。
- 独立时区设置页 / tab：合同已定「头部通道」，避免设置面膨胀。
- 用户级货币/数字格式覆盖：合同 §4.2 明确不含，弃。