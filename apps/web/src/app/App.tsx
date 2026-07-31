import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";

/**
 * R1 应用壳占位：单页 + 主题切换。
 * 不含 Admin 导航壳 / 多业务路由（R3）与协议 Renderer 全量（R5）。
 */
export function App() {
  return (
    <div className="mx-auto flex min-h-screen max-w-2xl flex-col gap-6 p-8">
      <header className="flex items-center justify-between gap-4">
        <div>
          <p className="text-sm text-muted-foreground">schema-ui-core · R1</p>
          <h1 className="text-2xl font-semibold tracking-tight">
            Admin 前端骨架
          </h1>
        </div>
        <ThemeToggle />
      </header>

      <main className="rounded-lg border border-border bg-card p-6 text-card-foreground shadow-sm">
        <p className="text-sm leading-6 text-muted-foreground">
          Vite + React 19 + TypeScript + Tailwind + shadcn/ui 基线已接入。
          业务页、协议范例与完整外壳尚未实现。
        </p>
        <div className="mt-4 flex flex-wrap gap-2">
          <Button type="button">Primary</Button>
          <Button type="button" variant="secondary">
            Secondary
          </Button>
          <Button type="button" variant="outline">
            Outline
          </Button>
        </div>
      </main>
    </div>
  );
}
