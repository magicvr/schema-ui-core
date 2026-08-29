# Web 包边界设计 v0.1（I-002 草案 · 2026-08-29）

来源：GOAL-004 E-001（apps/web 结构扫描）。状态：草案，S1 决策后生效。

## 1. 包清单（六包）

| 包 | 来源目录 | 导出面（草案） | 依赖 | 说明 |
|----|----------|----------------|------|------|
| `@schema-ui/protocol` | src/protocol | load-page / app-manifest / conformance 校验 / 类型 | lib（fetch-timeout）+ ajv（peer deps） | 最轻首包：纯 TS |
| `@schema-ui/lib` | src/lib | fetch-timeout 等通用工具 | — | 随 protocol 拆出 |
| `@schema-ui/theme` | src/theme | Token 语义 + brand 覆盖机制（CSS + TS） | tailwind | 定制面即此包 |
| `@schema-ui/ui` | src/components | shadcn 风格原子组件 + 导出类型 | react / radix / cva / clsx / tailwind-merge / lucide | peer 矩阵最大 |
| `@schema-ui/renderer` | src/renderer | renderSchemaPage / schema-table / form-controls / reactions / permissions | protocol / ui / theme / i18n？ | 渲染核心 |
| `@schema-ui/shell` | src/app + host | App / LoginPage / AuthGate / navigation / host-bootstrap 壳 | renderer / ui / theme / protocol | 下游骨架组装面 |

## 2. peer 依赖耦合矩阵（草案）

| 依赖 | protocol | lib | theme | ui | renderer | shell |
|------|----------|-----|-------|----|----------|-------|
| react ^19 | — | — | — | peer | peer | peer |
| react-dom ^19 | — | — | — | peer | peer | peer |
| tailwind ^4 | — | — | dep | peer(css) | peer(tokens) | peer |
| ajv ^8 | dep/peer | — | — | — | — | — |
| @radix-ui/* , cva, clsx, tailwind-merge, lucide | — | — | — | dep/peer | — | — |

规则草案：UI 相关包 = React/Tailwind 双 peer（版本窗在 changelog 声明）；protocol = ajv peer；theme 输出 CSS 产物（Tailwind 编译）与 TS Token。

## 3. 实施路径评估（关键决策）

| 路径 | 内容 | 优点 | 缺点 |
|------|------|------|------|
| **A · 源码 monorepo 化** | apps/web 拆为 `packages/*`（pnpm workspace），源码迁移 + import 相对化 + 各自 tsc/Vite lib 构建 | 边界最干净、长期正式化形态 | 改动面大（89 文件迁移 + CI/vitest 路径调整 + Tailwind 双编译方案） |
| **B · Vite lib 产物打包（推荐试点）** | 保持 src 结构，新增 lib 构建入口（vite build --mode lib 多产物），产出可发布 ESM+d.ts；下游 consume 产物 | 不改源码、快速闭环（与 R2 同思路）、可先发 protocol 单包 | 包内 import 别名（@/）需在构建 resolve；产物边界粗（renderer 包含内部图） |
| **C · 仅 protocol 最小切片** | 只把 protocol+lib 拆为真包（源码复制入 packages/ 或 lib 打包），下游消费协议层 | 最小证据：npm 包分发/消费链路成立 | 不达判据 #2「渲染 schema 页面集」主体 |

推荐：**B 先行 + A 择机正式化**（B 证明分发闭环与包消费形态；A 是 go 后 monorepo 化专项）。C 作为 B 的第一增量（protocol 首包）。

## 4. Token 覆盖验证设计（判据 #2 定制面）

下游 app：import `@schema-ui/theme`（Token 语义）+ 自建 `brand.css` 覆盖（复制 brand.example.css 机制）→ 渲染页面颜色换肤可核对。此机制已在主仓成型（theme/brand.example.css），包形态下不变。

## 5. 与冻结面关系

- Web 侧冻结面（v1.1.0 §范围 注明 R3 另行落盘）= 本设计 v1.0.x 生效后升格：六包导出签名 + peer 矩阵入清单（semver 同纪律）。
- `@schema-ui/shell` 的 AuthGate/登录与 Go 侧 assembly 对应：下游组合根 = shell 包 + config + 业务 schema 文档。