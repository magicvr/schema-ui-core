---
id: E-001
doc: execution
title: S1 · 多语种核心实施（locale/catalog/format/runtime/切换器）+ 测试证据
status: recorded
parent: GOAL-002-s1-locale-core
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-001 · S1 多语种核心（2026-08-09）

## 事实

- **新增单元**（全部 shipped 代码，非测试专用实现）：
  - `src/i18n/locale.ts`：`SUPPORTED_LOCALES`、`normalizeLocaleCandidate`（BCP47 变体/大小写/下划线/语言前缀归一）、`resolveLocale`（冻结优先级纯函数：显式 → 系统默认非 auto → 浏览器偏好 → `en-US`）、`normalizePreference`。
  - `src/i18n/messages/en-US.json` + `zh-CN.json`：纯数据 catalog（en-US 为规范基线；S1 键：语种名/切换器/启动屏）。
  - `src/i18n/catalog.ts`：`translate`/`createTranslator`/`hasTranslation`/`lookupTranslation`/`interpolate`；缺失 key 判定（当前语种 + en-US 均无条目）→ `schema-ui:missing-translation` 事件（按 locale+key 去重）+ 回退链（当前 → en-US → key），不渲染空、不抛异常。
  - `src/i18n/format.ts`：`formatDate`/`formatNumber` 随有效 locale（`Intl.*`，无自定义模板）；无效时间戳/无效 IANA 时区安全降级。
  - `src/i18n/runtime.tsx`：`I18nProvider`/`useI18n`/`readStoredLocale`/`writeStoredLocale`/`applyLocaleToDocument`；localStorage 单通道（`schema-ui:locale`，auto=移除键），登出不清除；`<html lang>` 跟随。
  - `src/components/locale-switcher.tsx`：Shell 顶栏 + 匿名登录页可达的语种切换器（无需任何权限）。
- **接线**：`main.tsx` 以 `I18nProvider` 包裹 `AuthProvider`/`ManifestFailure`；`App.tsx` 顶栏与 `LoginPage.tsx` 头部加入 `LocaleSwitcher`；BootScreen 文案走 catalog（`app.boot.checkingSession`）。
- **既有测试适配**：App/LoginPage 相关 4 个测试文件渲染处包 `I18nProvider`（useI18n 必须位于 Provider 内）。
- **验证**：vitest **674/674**（35 文件）全绿，其中新增 i18n 测试 **45**（locale 优先级矩阵/归一化、catalog 回退链/缺失可观察/参数插值、format 双语种/时区/失败语义、runtime 持久化/切换/`lang`、切换器可达性与标签）；输出捕获 `{SCRATCH}/unit-s1-web.log`；`npm run build` 通过（tsc + vite，exit 0）。

## 产物

| 路径 | 说明 |
|------|------|
| `src/i18n/locale.ts` + `locale.test.ts` | 解析/归一化/优先级（C1） |
| `src/i18n/messages/*.json` + `catalog.ts` + `catalog.test.ts` | 资源装载/缺失可观察回退（C2/C3） |
| `src/i18n/format.ts` + `format.test.ts` | 日期/数字格式化（C5） |
| `src/i18n/runtime.tsx` + `runtime.test.tsx` | Provider/切换/持久化/`lang`（C4/C5） |
| `src/components/locale-switcher.tsx` + 测试 | 用户切换入口（C4） |

## 里程碑 checkpoint

- commit：`10f20ab`（2026-08-09，S1；owned paths = `apps/web/src/i18n/**`、`src/components/locale-switcher*`、`main.tsx`、`App.tsx`、`LoginPage.tsx`、受影响的 4 个测试文件 + workspace-007 治理文档，显式 `git add` 无 `-A`）。
