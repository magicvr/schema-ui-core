---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-web-package-consumption
version: 0.1.0
---

# D-002 · S3 渲染闭环实施方案（关键决策 · 待用户裁决）

## 侦察事实（技术基线）

1. **顶层渲染入口**：`RenderPage({ document, context, tableRenderer, dataFetcher, onAction, onNavigate, formComponent })`——document（schema 页面文档）+ 可注入 `dataFetcher`（`typeof fetch`）→ 可在无浏览器环境（SSR/`renderToString`）执行初始渲染。
2. **renderer 依赖面**：`@/components/ui/*`（async-state/card/skeleton/input/label/textarea）+ `@/components/data-table` + `@/i18n/{runtime,catalog,…}`（i18n 是主要内部面）；**不直接依赖 theme**。
3. **CSS 面**：全仓唯一样式 import = `main.tsx → index.css`（Token 体系）；组件只含 className → **包组件零 CSS**，样式/Tailwind/Token 完全留在应用侧（下游自带 index.css + brand 覆盖）——Token 覆盖定制天然成立。

## 方案（B 路径增量）

**粗粒度单包 `@schema-ui/renderer` v0.1.0**：

- entry：新建 `src/renderer/index.ts` → 导出 `RenderPage` + `I18nProvider`（re-export `@/i18n/runtime`）+ 关键类型（`RenderPageDocument`/`RenderMeta`/`SchemaCrudValue`）+ `registerCustomComponent`（扩展接缝）。
- 打包：Vite lib（同 protocol 链路）；alias `@` 解析 components/i18n/protocol 引用（内部 bundle 自包含）；**external: react / react-dom / react/jsx-runtime**（React 组件库硬性 peer，防双实例）；declaration 同链路。
- 产物：`dist-lib/@schema-ui/renderer/` + package.json（peerDeps: react ^19、react-dom ^19；指引 `@source` 指令供下游 Tailwind 4 扫描包内 className）。

**验证（判据 #2 核心）**：

1. **SSR 渲染**：golden-web 仅包依赖（react + @schema-ui/renderer + protocol）→ `renderToString(<I18nProvider><RenderPage document={真实主线 schema fixture}/></I18nProvider>)` → 结构断言（页面化节点/表格壳/标题文本存在、无异常）。
2. **schema 同源性**：fixture 直接取自主线模块 schema（`apps/api/modules/users/schema/users.json` 裁剪为最小可渲染文档）→ 「与主线同一的 schema 页面集」证据。
3. **Token 覆盖断言**：golden-web 自带 `index.css`（Token 变量集）+ `brand.css`（覆盖）→ 轻量一致性脚本（仅覆盖 index.css 已声明变量——复用主仓 `brand-example.test.ts` 模式）验证定制面成立。
4. **peer/Tailwind 指引**：README 段落（下游 Tailwind 4 `@source` 扫描 renderer 产物类名）→ peer 矩阵实测点。

**六包细化（ui/theme 独立包）**：保留为 go 后正式化（冻结面 Web 侧 v1.2 按六包列）；试点以粗粒度换闭环速度，产物级可逆。

## 不进本目标

- 不包 shell/登录流（App/LoginPage 依赖 API 会话，属 fork 应用壳——R4 演练再议）。
- 不做浏览器级渲染（e2e 主仓已有视觉回归基线；SSR 结构断言 + Token 文件断言为本试点证据面）。
- 不引入 golden-web 的 vitest/jsdom 重型依赖（SSR 零额外依赖）。