# apps/web · schema-ui-core Web 骨架

R1 React 前端骨架（GOAL-004）。Tailwind + shadcn/ui 基线 + 浅/深色占位；**无业务路由树、无 Admin 导航壳**。

## 要求

- Node.js 20+（本仓 R1 实测 Node 22）
- 包管理：**npm** + `package-lock.json`

## 工具链

| 项 | 值 |
|----|-----|
| Vite | 6.x |
| React | 19 |
| TypeScript | 5.8 |
| Tailwind | 4.x（`@tailwindcss/vite`） |
| shadcn | `components.json` style **new-york**；示例 `src/components/ui/button.tsx` |

## 目录

```text
src/app/           # 应用壳 / 占位页
src/components/ui/ # shadcn 组件
src/host/          # 空分层占位（R1）
src/protocol/      # 空分层占位（R1）
src/renderer/      # 空分层占位（R1）
src/lib/utils.ts   # cn() 等
```

## 运行

```bash
npm install
npm run dev
# http://localhost:5173

npm run build
```

## shadcn 痕迹

- `components.json`（new-york、cssVariables）
- `src/components/ui/button.tsx`
- 文档：本 README；初始化采用与 shadcn 约定一致的 aliases / CSS 变量（非手写 div 冒充）

后续可继续：

```bash
npx shadcn@latest add <component>
```

## 非目标（R1）

- Admin 侧栏 / App manifest 导航（R3）
- 协议 Renderer 全量与业务范例页（R5）
- 平行仓 `mock-app` 业务演示
