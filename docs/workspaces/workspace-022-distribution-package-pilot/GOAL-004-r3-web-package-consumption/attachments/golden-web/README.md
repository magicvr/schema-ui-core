# golden-web · Web 包消费试点黄金仓（GOAL-004）

仅经 npm 包形态消费 `@schema-ui/*`（`file:` 指向 `apps/web/dist-lib/` 产物，模拟上游已发布包）：

| 探针 | 内容 | 结果 |
|------|------|------|
| `node probe.mjs` | protocol 包功能断言（manifest 协商/表达式/URL 解析） | PASS |
| `node probe-render.mjs` | **渲染闭环**：`@schema-ui/renderer` SSR 渲染真实形态 schema 页面文档（表单 + defaultValue + 能力门控） | PASS |
| `node token-check.mjs` | Token 覆盖纪律：brand.css ⊆ index.css 变量集（fork 定制面） | PASS |

## 下游 Tailwind 4 指引（peer 矩阵实测点）

包组件只含 className（零 CSS 产物）；样式由下游构建编译。Tailwind 4（CSS-first）扫描 node_modules 产物：

```css
@import "tailwindcss";
@source "../../node_modules/@schema-ui/renderer/dist"; /* 或实际安装路径 */
```

若未配置 `@source`，包内 className 不会被生成 → 渲染无样式（结构仍在）。此为「CSS 面归下游」契约的显式声明。

## 非目标（试点边界）

- 不做浏览器级 e2e（主仓已有视觉回归基线）；SSR 结构断言 + Token 文件断言为本试点证据面。
- 不包 shell/登录流（App/LoginPage 属 fork 应用壳；R4 演练再议）。