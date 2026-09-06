// @vitest-environment node
//
// W25 · 防复发机制：模块 schema 中声明的所有渲染层自定义组件
// ({type:"custom", component:"…"}) 必须在 web 前端注册表中存在
// (getCustomComponent !== null)。作用域仅为渲染层 custom 节点——动作层的
// {type:"custom", handler:"…"}（如 users/roles export、file-library 下载/
// 预览/复制）走动作分发器，不属于本校验。
import { readdirSync, readFileSync, existsSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

import { getCustomComponent } from "@/renderer/custom-components";
// Side-effect imports mirror main.tsx: each component module self-registers
// into the custom-component registry (the test asserts schema ↔ registry
// agreement, so the registry must be populated the same way the app does).
import "@/components/mfa-manager";
import "@/components/account-session-toolbar";
import "@/components/cron-preview";
import "@/components/monitoring-auto-refresh";
import "@/components/import-template-download";
import "@/components/wallet-ensure";
import "@/components/notification-center";
import "@/components/data-permission-scopes";
import "@/components/activity-export";
import "@/components/mail-admin-tab";
// workspace-018 R3 shipped the component but missed this side-effect import,
// so the W25 guard failed at HEAD (account:email-identity "unregistered").
import "@/components/email-identity";
import "@/components/password-policy-tab";
import "@/components/invite-issue-card";
import "@/components/invite-resend-dialog";
// A-002 F-003 (GOAL-041 S2): telegram-admin-tab lives in a nested module
// (channel/telegram/schema) that the historical one-level walker never saw;
// the recursive walk now scans it, so the component must be imported here.
import "@/components/telegram-admin-tab";

const MODULES_DIR = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../../api/modules",
);

// Recursive schema walk (A-002 F-003, GOAL-041 S2): nested modules
// (channel/telegram/schema, dev/examples/schema) were invisible to the
// historical one-level walk, so telegram-admin-tab etc. escaped this guard.
function walkSchemaFiles(dir: string, out: string[]): void {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const abs = join(dir, entry.name);
    if (entry.isDirectory()) {
      walkSchemaFiles(abs, out);
    } else if (
      entry.isFile() &&
      abs.replace(/\\/g, "/").includes("/schema/") &&
      entry.name.endsWith(".json")
    ) {
      out.push(abs);
    }
  }
}

function collectSchemaFiles(): string[] {
  const out: string[] = [];
  for (const moduleName of readdirSync(MODULES_DIR)) {
    const moduleRoot = join(MODULES_DIR, moduleName);
    if (!existsSync(moduleRoot)) {
      continue;
    }
    walkSchemaFiles(moduleRoot, out);
  }
  return out;
}

interface ComponentRef {
  pageId: string;
  component: string;
}

function collectComponentRefs(value: unknown, pageId: string, out: ComponentRef[]): void {
  if (Array.isArray(value)) {
    for (const item of value) {
      collectComponentRefs(item, pageId, out);
    }
    return;
  }
  if (value !== null && typeof value === "object") {
    const obj = value as Record<string, unknown>;
    if (obj.type === "custom" && typeof obj.component === "string") {
      out.push({ pageId, component: obj.component });
    }
    for (const child of Object.values(obj)) {
      collectComponentRefs(child, pageId, out);
    }
  }
}

describe("module schema custom components (W25 防复发)", () => {
  it("every schema-declared renderer custom component is registered", () => {
    const refs: ComponentRef[] = [];
    for (const file of collectSchemaFiles()) {
      const raw = readFileSync(file, "utf8");
      const document = JSON.parse(raw) as { meta?: { pageId?: unknown } };
      const pageId =
        typeof document?.meta?.pageId === "string" ? document.meta.pageId : file;
      collectComponentRefs(document, pageId, refs);
    }
    expect(refs.length).toBeGreaterThan(0);
    const missing = refs
      .filter((ref) => getCustomComponent(ref.component) === null)
      .map((ref) => `${ref.pageId}:${ref.component}`);
    expect(missing).toEqual([]);
  });
});