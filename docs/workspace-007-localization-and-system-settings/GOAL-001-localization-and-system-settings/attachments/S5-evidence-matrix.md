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
| overview | N/A¹ | ✓ | N/A¹ | ✓ | `schema-keys.structural.test.ts`；渲染：`ui-bilingual` / e2e 登录后「总览」 |
| admin-list-batch | N/A¹ | ✓ | N/A¹ | ✓ | structural + `representative-pages.integration.test.tsx` |
| data-display | N/A¹ | ✓ | N/A¹ | ✓ | `schema-keys.structural.test.ts` |
| data-table | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| search-form-table | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| form-controls | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| form-with-reactions | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| form-with-upload | N/A¹ | ✓ | N/A¹ | ✓ | 同上 |
| users | N/A¹ | ✓ | N/A¹ | ✓ | `ui-bilingual.test.tsx` + `schema-crud.test.tsx` |
| roles | N/A¹ | ✓ | N/A¹ | ✓ | `schema-keys.structural.test.ts`（44 keys） |
| settings | N/A² | N/A² | N/A¹ | ✓ | structural + `startup-config.test.tsx` + e2e M3 |
| activity | N/A² | N/A² | N/A¹ | ✓ | `schema-keys.structural.test.ts`（16 keys） |

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
| M1 | `src/i18n/locale.test.ts`（优先级）；`ui-bilingual.test.tsx`（登录双语 + 错误码）；e2e M1（切换 zh、`lang`、登录成功 shell） |
| M2 | `ui-bilingual.test.tsx`（users 写表单）；`schema-crud.test.tsx`；e2e 登录后 overview |
| M3 | `startup-config.test.tsx`（PATCH + 刷新头 + branding 投影）；Go `handler/settings_test.go`；e2e 站点标题保存 → header/title |
| M4 | `ui-bilingual.test.tsx` M4；`catalog.test.ts`（事件去重/回退链） |
| 权限正反 | `startup-config.test.tsx`（write 缺失禁用保存；有权限真实 PATCH） |
| 配置刷新 | `startup-config.test.tsx` + Go settings PATCH header `X-Schema-UI-Config-Changed`；e2e 可见投影 |
| 错误回退 | S4：`handler/localize_test.go`、`error_contract_test.go`；`src/renderer/error-localization.test.tsx`；e2e Accept-Language zh 登录失败 envelope |

## 2. 真实入口与浏览器证据（C2 交叉引用）

| 项 | 路径 / 命令 |
|----|-------------|
| API 双启动 `/api/branding` 体一致 | scratch `s5-launch/run1.json` + `run2.json` + 比对日志 |
| Web production build | scratch `s5-launch/web-build.log` |
| Playwright e2e（admin） | `apps/web/e2e/localization.spec.ts`；截图 `apps/web/test-results/s5-settings-zh.png`；scratch `s5-launch/e2e-localization.log` |
| Unit suites（矩阵引用） | scratch `s5-tests/` 摘要（按需刷新） |

## 3. 分母完整性声明

- 行分母 8 格全覆盖；列分母与 F-V029 / VP-007 exit 6 一致。  
- 非 N/A 单元格均指向 shipped 测试名或捕获日志/截图，无「代表性」收缩。  
- 未列入 F-V029 的用户可见文案面：本波次无额外 residual 排除项。
