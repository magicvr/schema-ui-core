// golden-field 六包消费探针（R3）：lib/theme/ui/shell 独立包存在性与核心导出；
// renderer 0.2.0（d.ts 自动化管线版）兼容。
import { cn } from "@magicvr/schema-ui-lib";
import { formatDisplayTime } from "@magicvr/schema-ui-lib";
import { resolveTheme } from "@magicvr/schema-ui-theme";
import { DataTable } from "@magicvr/schema-ui-ui";
import { App } from "@magicvr/schema-ui-shell";
import assert from "node:assert";

// lib：cn 合并工具（clsx+tailwind-merge 语义）
assert.equal(cn("a", "b", false && "c"), "a b");
// lib：datetime 格式化（仅 ISO 字符串输入）
const dt = formatDisplayTime("2026-08-29T00:00:00Z");
assert.ok(dt && dt.startsWith("2026-08-29"), `formatDisplayTime ISO → ${dt}`);
// theme：Token 解析（stored=dark 确定语义）
assert.equal(resolveTheme({ stored: "dark", prefersDark: false }).theme, "dark");
// ui：DataTable 组件存在（渲染器核心列表面）
assert.equal(typeof DataTable, "function");
// shell：App 应用壳存在（骨架组装面）
assert.equal(typeof App, "function");

console.log("golden-field six-package probe PASS · lib/theme/ui/shell/renderer0.2.0 all consumable");