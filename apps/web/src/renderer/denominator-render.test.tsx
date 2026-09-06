// @vitest-environment jsdom
//
// S5 (GOAL-041 W29) · full page-denominator render: every one of the 35 module
// schema documents must survive the production chain
//   validatePageDocument (D-VAL) → loadPageDocument (incl. F-001 page-level
//   version+capability negotiation) → RenderPage
// and render a non-empty surface without the schema-error page or an
// unknown-custom placeholder. This closes the C-007 denominator: structural
// coverage (35/35 D-VAL) alone was not runtime render proof. The representative
// pages (GOAL-004) remain the behavior-touched subset; this test proves the
// whole denominator renders.

import { readFileSync, readdirSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, describe, expect, it } from "vitest";

import { AuthProvider } from "@/account/AuthContext";
import type { PageEntry } from "@/protocol/app-manifest";
import { loadPageDocument } from "@/protocol/load-page";
import type { RenderPageDocument } from "@/renderer/render.types";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";

// Side-effect imports mirror main.tsx so every custom node in the denominator
// resolves to its real component (never the unknown-custom placeholder).
import "@/components/account-session-toolbar";
import "@/components/activity-export";
import "@/components/cron-preview";
import "@/components/data-permission-scopes";
import "@/components/email-identity";
import "@/components/import-template-download";
import "@/components/invite-issue-card";
import "@/components/invite-resend-dialog";
import "@/components/mail-admin-tab";
import "@/components/mfa-manager";
import "@/components/monitoring-auto-refresh";
import "@/components/notification-center";
import "@/components/password-policy-tab";
import "@/components/telegram-admin-tab";
import "@/components/wallet-ensure";

const MODULES = resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules");

function walkSchemaFiles(dir: string, out: string[]): void {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const abs = join(dir, entry.name);
    if (entry.isDirectory()) {
      walkSchemaFiles(abs, out);
    } else if (entry.isFile() && abs.replace(/\\/g, "/").includes("/schema/") && entry.name.endsWith(".json")) {
      out.push(abs);
    }
  }
}

function collectSchemaFiles(): string[] {
  const out: string[] = [];
  for (const moduleName of readdirSync(MODULES)) {
    walkSchemaFiles(join(MODULES, moduleName), out);
  }
  return out;
}

interface SchemaRef {
  file: string;
  pageId: string;
  document: unknown;
}

const REFS: SchemaRef[] = collectSchemaFiles().map((file) => {
  const raw = readFileSync(file, "utf8");
  const document = JSON.parse(raw) as { meta?: { pageId?: string } };
  return {
    file,
    pageId: document.meta?.pageId ?? file,
    document,
  };
});

function fixtureFetcher(): typeof fetch {
  return (async (input: RequestInfo | URL) => {
    const url = String(input);
    if (url.startsWith("/api/")) {
      return new Response(
        JSON.stringify({ items: [], total: 0, page: 1, pageSize: 100 }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    return new Response(JSON.stringify({ error: "NOT_FOUND" }), { status: 404 });
  }) as typeof fetch;
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

async function renderDocument(pageDoc: RenderPageDocument): Promise<HTMLDivElement> {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  const fetcher = fixtureFetcher();
  await act(async () => {
    root.render(
      <AuthProvider>
        <RenderPage
          document={pageDoc}
          context={{}}
          dataFetcher={fetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
        />
      </AuthProvider>,
    );
  });
  return container;
}

describe("S5 · full page denominator renders (35/35)", () => {
  it("collects all 35 module schema documents", () => {
    expect(REFS.length).toBe(35);
  });

  for (const ref of REFS) {
    it(`${ref.pageId} passes D-VAL → loadPageDocument → RenderPage`, async () => {
      const page: PageEntry = {
        pageId: ref.pageId,
        title: ref.pageId,
        schemaUrl: `/api/schema/${ref.pageId}`,
        route: `/${ref.pageId}`,
      };
      // loadPageDocument covers D-VAL + F-001 page-level negotiation; a
      // mismatch (unsupported version / missing capability) throws here.
      const loaded = await loadPageDocument(
        page,
        {},
        { fetcher: async () => new Response(JSON.stringify(ref.document), { status: 200 }) },
      );
      expect((loaded as { meta?: { pageId?: unknown } }).meta?.pageId).toBe(ref.pageId);

      const container = await renderDocument(loaded as RenderPageDocument);
      const text = container.textContent ?? "";
      expect(text.trim().length, `${ref.pageId} must render non-empty content`).toBeGreaterThan(0);
      // Fail-closed surfaces must not appear.
      expect(container.querySelector('[aria-labelledby="schema-error-title"]')).toBeNull();
      expect(container.textContent).not.toContain("unknown custom component");
    });
  }
});
