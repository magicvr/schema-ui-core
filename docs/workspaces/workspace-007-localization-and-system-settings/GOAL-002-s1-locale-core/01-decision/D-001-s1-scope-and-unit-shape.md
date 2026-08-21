---
id: D-001
doc: decision
title: S1 · 实施范围与 API 形状（locale 解析/资源/切换/格式化单元边界）
status: accepted
parent: GOAL-002-s1-locale-core
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# D-001 · S1 实施范围与 API 形状

## 触发

S0 契约冻结完成（Root D-002 accepted，`I-L10N-001`/`I-L10N-002` verified）。S1 实施前冻结本阶段单元边界与 API 形状，避免实现期漂移。

## 决定

1. **纯逻辑与 I/O 分离**：locale 解析/匹配/回退、缺失 key 判定、日期/数字格式化决策做成无副作用纯单元（`src/i18n/`），直接可单测；DOM/事件副作用隔离在薄层。
2. **单元边界**：
   - `src/i18n/locale.ts`：`SUPPORTED_LOCALES = ["zh-CN","en-US"]`；`resolveLocale(input)` 纯函数（stored / systemDefault / browserLanguages → 有效 locale 或 `auto` 语义），优先级：显式 → 系统默认（非 auto）→ 浏览器偏好（首支持语种）→ `en-US`。
   - `src/i18n/messages/`：`zh-CN.json` / `en-US.json` 纯数据 catalog；`en-US` 为规范基线。
   - `src/i18n/catalog.ts`：装载 + `translate(key, params?, locale?)`；缺失 key 判定（当前语种与 en-US 均无条目）→ 可观察（事件）且回退链（当前 → en-US → key）。
   - `src/i18n/format.ts`：`formatDate`/`formatNumber` 按有效 locale（`Intl.*`），无自定义模板。
   - `src/i18n/runtime.tsx`：React provider + `useI18n()`；切换写 `localStorage["schema-ui:locale"]`（移除 = auto）；`document.documentElement.lang` 跟随；`schema-ui:missing-translation` 事件。
3. **切换器**：`src/components/locale-switcher.tsx`（Shell/用户菜单 + 匿名登录页可达；无需设置权限）。
4. **范围边界**：本阶段不翻译任何现有页面文案（S2）；不改动任何页面/schema 文件；不动 API。
5. **验证**：vitest 驱动 shipped 函数（不 mock 被测单元）；输出捕获 scratch；`npm run build` 通过。

## 未选方案

- 不引入第三方 i18n 库（如 i18next/react-i18next）：自研薄运行时满足范围且无依赖面扩张。
- 不做动态资源懒加载（两语种内置 bundling，装载失败降级为 en-US）。
- 不把 locale 写入账号资料（I-L10N-002 已冻结 localStorage 单通道）。

## 影响

- C1～C6 检查点全部由上述单元承载；S2 在 `src/i18n/` 之上做页面双语化。
