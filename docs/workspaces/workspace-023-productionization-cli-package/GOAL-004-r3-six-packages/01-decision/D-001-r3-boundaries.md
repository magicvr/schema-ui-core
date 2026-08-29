# D-001 · 六包边界定案（2026-08-29）

| 包 | 来源 | 导出面（v0.1） | peer/external | 备注 |
|----|------|-----------------|----------------|------|
| @magicvr/schema-ui-protocol | src/protocol | 已有 0.2.0 | 自包含 | 不动 |
| @magicvr/schema-ui-lib | src/lib + src/i18n | cn/formatDisplayTime/fetch-timeout + i18n 运行时（I18nProvider/useI18n/locale/catalog） | 自包含（bundle） | i18n 并入 lib（语言运行时 = 通用面） |
| @magicvr/schema-ui-theme | src/theme | resolveTheme/applyThemeToElement/initTheme/setTheme + Theme 类型 | 自包含 | CSS 面（token 变量/品牌覆盖）留下游（README 指引） |
| @magicvr/schema-ui-ui | src/components/ui + data-table | 原子（button/card/input/label/badge 等）+ DataTable | react peer（external） | 业务示范组件（locale-switcher 等）留 shell |
| @magicvr/schema-ui-renderer | src/renderer | 已有 0.1.0（自包含 bundle） | react peer | 依赖图细化（ui 包消费）留 go 后 |
| @magicvr/schema-ui-shell | src/app + src/host | App/AuthGate/LoginPage/navigation/branding + host bootstrap | react peer | 自包含 v0.1 |

## 决策点

- 不重排源码（B 路径产物层六包）；renderer 依赖图 external 化与纯原子拆分 = go 后六包化专项（登记）。
- d.ts：尝试 dts-bundle-generator 单文件管线（I-023-005）。