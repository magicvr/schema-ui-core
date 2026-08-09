# S5 · 双 Profile 证据矩阵（F-V029 分母，2026-08-09）

> **状态**：`filled`（GOAL-006 C1）。分母权威 = [F-V029-coverage-table-s0-freeze.md](./F-V029-coverage-table-s0-freeze.md)。  
> **行**：`zh-CN` / `en-US` × `mvp` / `admin` × 匿名 / 已认证。  
> **列**：固定 UI / 冻结 pageId–schema 并集 / M1～M4 / 权限正反例 / 缺失翻译 / 配置刷新 / 错误回退。  
> **N/A**：仅 Profile 不可达单元格，并写明模块边界；不算 pass。  
> **证据路径约定**：相对 `apps/web/` 或 `apps/api/` 的测试文件名 + 可选用例语义；浏览器端到端见 `apps/web/e2e/localization.spec.ts`（admin Profile）。

## 1. 行 × 列总表

图例：`✓` = 有可核对证据路径；`N/A` = Profile 模块边界不可达。

### 1.1 固定 UI 列

| 行（locale × profile × auth） | U1 登录 | U2 Shell header | U3 导航 | U4 切换器 | U5 反馈 | U6 通用组件 | U7 Settings 四类 |
|-------------------------------|---------|-----------------|---------|-----------|---------|-------------|------------------|
| zh-CN × mvp × 匿名 | ✓ | N/A¹ | N/A¹ | ✓ | ✓ | N/A¹ | N/A² |
| zh-CN × mvp × 已认证 | N/A³ | ✓ | ✓ | ✓ | ✓ | ✓ | N/A² |
| en-US × mvp × 匿名 | ✓ | N/A¹ | N/A¹ | ✓ | ✓ | N/A¹ | N/A² |
| en-US × mvp × 已认证 | N/A³ | ✓ | ✓ | ✓ | ✓ | ✓ | N/A² |
| zh-CN × admin × 匿名 | ✓ | N/A¹ | N/A¹ | ✓ | ✓ | N/A¹ | N/A³ |
| zh-CN × admin × 已认证 | N/A³ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| en-US × admin × 匿名 | ✓ | N/A¹ | N/A¹ | ✓ | ✓ | N/A¹ | N/A³ |
| en-US × admin × 已认证 | N/A³ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

¹ 匿名态无 Shell/导航/表单面（未登录路由 → 登录页）。  
² `admin.settings` 不在 mvp 模块集（`kernel/profile.go`）。  
³ 已认证态不展示登录页；匿名态不展示 Settings 编辑面。

**证据路径（固定 UI）**

| 面 | 证据 |
|----|------|
| U1 | `src/i18n/ui-bilingual.test.tsx`（zh/en 登录面、错误码映射、`lang`）；`src/app/LoginPage.test.tsx`；e2e `localization.spec.ts`（zh 登录面 + 登录） |
| U2 | `ui-bilingual.test.tsx`（Shell 文本）；`src/app/shell.test.ts`；`src/app/App.integration.test.tsx` |
| U3 | `ui-bilingual.test.tsx`（titleKey/labelKey 双语 + 字面回退）；`src/i18n/schema-keys.structural.test.ts`（manifest key） |
| U4 | `src/components/locale-switcher.test.tsx`；`src/i18n/runtime.test.tsx`；e2e 切换器 → `lang=zh-CN` |
| U5 | `ui-bilingual.test.tsx`（路由回退、空态、分页）；`App.integration.test.tsx`；`src/renderer/visual-fidelity.test.tsx` |
| U6 | `ui-bilingual.test.tsx`（列头/工具栏/提交/确认）；`src/renderer/schema-crud.test.tsx`；`src/components/data-table.test.tsx` |
| U7 | `schema-keys.structural.test.ts`（settings.json keys + 双语 catalog）；`src/app/startup-config.test.tsx`（四类工具条 zh、General PATCH、权限禁用、reset、bootstrap 投影）；e2e settings 保存可见变化 |

### 1.2 Manifest / Schema 并集列（12 pageId）

| pageId | mvp 匿名 | mvp 已认证 | admin 匿名 | admin 已认证 | 证据 |
|--------|----------|------------|------------|--------------|------|
| overview | N/A¹ | ✓ | N/A¹ | ✓ | structural + **S5** `apps/web/src/i18n/s5-denominator-render.test.tsx` titleKey；e2e overview |
| admin-list-batch | N/A¹ | ✓ | N/A¹ | ✓ | structural + `apps/web/src/app/representative-pages.integration.test.tsx` + S5 titleKey |
| data-display | N/A¹ | ✓ | N/A¹ | ✓ | structural + S5 s5-denominator-render titleKey zh/en |
| data-table | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| search-form-table | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| form-controls | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| form-with-reactions | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| form-with-upload | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| users | N/A¹ | ✓ | N/A¹ | ✓ | `apps/web/src/i18n/ui-bilingual.test.tsx` + `apps/web/src/renderer/schema-crud.test.tsx` |
| roles | N/A¹ | ✓ | N/A¹ | ✓ | structural + **S5** s5-denominator-render 列头/工具栏 zh+en |
| settings | N/A² | N/A² | N/A¹ | ✓ | structural + `apps/web/src/app/startup-config.test.tsx` + e2e admin M3；mvp 反证 e2e |
| activity | N/A² | N/A² | N/A¹ | ✓ | structural + **S5** s5-denominator-render 正文/列 zh |

¹ 匿名不可达受保护 page。  
² 模块边界：`admin.settings` / `admin.activity` 不在 mvp profile。

语种维度：catalog 键集对称（`catalog.test.ts` + `schema-keys.structural.test.ts` 双语断言）；运行时渲染在 zh-CN/en-US 下由 `ui-bilingual.test.tsx` 与 e2e 覆盖。

### 1.3 M1～M4 / 权限 / 缺失翻译 / 配置刷新 / 错误回退

| 列 | zh-CN × mvp × 匿名 | zh-CN × mvp × 已认证 | en-US × mvp × 匿名 | en-US × mvp × 已认证 | zh-CN × admin × 匿名 | zh-CN × admin × 已认证 | en-US × admin × 匿名 | en-US × admin × 已认证 |
|----|--------------------|---------------------|--------------------|---------------------|----------------------|------------------------|----------------------|------------------------|
| M1 匿名启动→locale→登录反馈 | ✓ | N/A³ | ✓ | N/A³ | ✓ | N/A³ | ✓ | N/A³ |
| M2 登录→导航→读写表单 | N/A¹ | ✓ | N/A¹ | ✓ | N/A¹ | ✓ | N/A¹ | ✓ |
| M3 settings 投影 | N/A² | N/A² | N/A² | N/A² | N/A¹ | ✓ | N/A¹ | ✓ |
| M4 缺失 key fallback | ✓⁴ | ✓ | ✓⁴ | ✓ | ✓⁴ | ✓ | ✓⁴ | ✓ |
| 权限正例（有 settings.write） | N/A² | N/A² | N/A² | N/A² | N/A¹ | ✓ | N/A¹ | ✓ |
| 权限反例（缺 write 禁用） | N/A² | N/A² | N/A² | N/A² | N/A¹ | ✓ | N/A¹ | ✓ |
| 缺失翻译可观察 | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 配置刷新（X-Schema-UI-Config-Changed） | N/A² | N/A² | N/A² | N/A² | N/A¹ | ✓ | N/A¹ | ✓ |
| 错误回退（码+前端保底+协商） | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

¹ 匿名。 ² mvp 无 settings。 ³ 已认证非 M1 入口。 ⁴ 匿名态通过 catalog 单元 + 登录面缺失回退路径；完整 M4 列表面在已认证。

**证据路径（流程列）**

| 列 | 证据 |
|----|------|
| M1 | `apps/web/src/i18n/locale.test.ts`；`ui-bilingual.test.tsx`；e2e admin+mvp M1（切换 zh、`lang`、登录 shell） |
| M2 | **单元写表单** `ui-bilingual` + `schema-crud`；**浏览器** e2e 登录后 overview（写表单以单元为主，F-004 已标注） |
| M3 | `startup-config.test.tsx`；Go `handler/settings_test.go`；e2e admin 站点标题保存投影 |
| M4 | `ui-bilingual.test.tsx` M4；`catalog.test.ts` |
| 权限正反 | `startup-config.test.tsx` |
| 配置刷新 | `startup-config.test.tsx` + Go PATCH header；e2e 可见投影 |
| mvp 真实入口 | e2e mvp + dual-run mvp branding + `RegisterPublicBranding`；settings 不可达 |
| 错误回退 | S4：`handler/localize_test.go`、`error_contract_test.go`；`src/renderer/error-localization.test.tsx`；e2e Accept-Language zh 登录失败 envelope |

## 2. 真实入口与浏览器证据（C2 交叉引用）

| 项 | 路径 / 命令 |
|----|-------------|
| API 双启动 `/api/branding` 体一致 | 仓库内 `GOAL-006/.../attachments/s5-launch/run{1,2}-{admin,mvp}.json` + compare |
| Web production build | `GOAL-006/.../attachments/s5-launch/web-build.log` |
| Playwright e2e | `apps/web/e2e/localization.spec.ts`（admin + mvp）；截图 s5-settings-zh / s5-mvp-overview-zh |
| 分母渲染 | `apps/web/src/i18n/s5-denominator-render.test.tsx`（5/5） |

## 3. 分母完整性声明

- 行分母 8 格全覆盖；列分母与 F-V029 / VP-007 exit 6 一致。  
- 非 N/A 单元格均指向 shipped 测试名或捕获日志/截图，无「代表性」收缩。  
- 未列入 F-V029 的用户可见文案面：本波次无额外 residual 排除项。
