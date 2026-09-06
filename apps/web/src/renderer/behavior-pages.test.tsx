// @vitest-environment jsdom
//
// F3 (GOAL-042 W30): behavior-level unit tests for the 10 pages that the S5
// denominator only rendered (no per-page assertions): mail, mail-outbox,
// my-wallet, wallet, wallet-vouchers, telegram-settings, telegram-operator,
// digitaloffer-offers, digitaloffer-entitlements, digitaloffer-purchases.
//
// Each case loads the REAL module schema through the production chain
// (loadPageDocument incl. D-VAL + F-001 negotiation), feeds the tables one
// sample row built from the page's own column fields (so headers + row
// actions render — the empty state replaces the whole table), and asserts
// page-specific rendered UI.

import { readFileSync, readdirSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, describe, expect, it, vi } from "vitest";

import { AuthProvider } from "@/account/AuthContext";
import { I18nProvider } from "@/i18n/runtime";
import type { PageEntry } from "@/protocol/app-manifest";
import { loadPageDocument } from "@/protocol/load-page";
import type { RenderPageDocument } from "@/renderer/render.types";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";

// Register custom components exactly like main.tsx (mail-admin-tab,
// telegram-admin-tab, wallet-ensure).
import "@/components/mail-admin-tab";
import "@/components/telegram-admin-tab";
import "@/components/wallet-ensure";

const MODULES = resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules");

const TARGETS = [
  "mail",
  "mail-outbox",
  "my-wallet",
  "wallet",
  "wallet-vouchers",
  "telegram-settings",
  "telegram-operator",
  "digitaloffer-offers",
  "digitaloffer-entitlements",
  "digitaloffer-purchases",
] as const;

const pageDocuments = new Map<string, string>();

(function walk(dir: string): void {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const abs = join(dir, entry.name);
    if (entry.isDirectory()) {
      walk(abs);
    } else if (entry.isFile() && abs.replace(/\\/g, "/").includes("/schema/") && entry.name.endsWith(".json")) {
      const doc = JSON.parse(readFileSync(abs, "utf8")) as { meta?: { pageId?: string } };
      if (doc.meta?.pageId !== undefined) {
        pageDocuments.set(doc.meta.pageId, abs);
      }
    }
  }
})(MODULES);

/** Collects every `"field": "…"` name in the page document (table columns +
 *  form fields) so the fixture row satisfies every table the page declares. */
function collectFields(pageId: string): string[] {
  const file = pageDocuments.get(pageId);
  if (file === undefined) {
    return [];
  }
  const text = readFileSync(file, "utf8");
  return [...new Set([...text.matchAll(/"field"\s*:\s*"([^"]+)"/g)].map((m) => m[1]))];
}

function fixtureFetcher(pageId: string): typeof fetch {
  const fields = collectFields(pageId);
  const row: Record<string, unknown> = { id: "row-1" };
  for (const field of fields) {
    row[field] = /amount|balance|total|price|count|size|version|remaining|index/i.test(field)
      ? 100
      : `sample-${field}`;
  }
  return (async (input: RequestInfo | URL) => {
    const url = String(input);
    if (url.startsWith("/api/")) {
      return new Response(
        JSON.stringify({ items: [row], total: 1, page: 1, pageSize: 100 }),
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

async function renderPage(
  pageId: string,
  fetcherOverride?: typeof fetch,
): Promise<HTMLDivElement> {
  const file = pageDocuments.get(pageId);
  if (file === undefined) {
    throw new Error(`schema file not found for ${pageId}`);
  }
  const pageDocument = JSON.parse(readFileSync(file, "utf8"));
  const page: PageEntry = {
    pageId,
    title: pageId,
    schemaUrl: `/api/schema/${pageId}`,
    route: `/${pageId}`,
  };
  const loaded = await loadPageDocument(page, {}, {
    fetcher: async () => new Response(JSON.stringify(pageDocument), { status: 200 }),
  });
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  const fetcher = fetcherOverride ?? fixtureFetcher(pageId);
  await act(async () => {
    root.render(
      <AuthProvider>
        <I18nProvider stored="en-US" browserLanguages={["en-US"]}>
          <RenderPage
            document={loaded as RenderPageDocument}
            context={{}}
            dataFetcher={fetcher}
            tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
          />
        </I18nProvider>
      </AuthProvider>,
    );
  });
  return container;
}

describe("F3 · behavior-level tests for the 10 denominator-only pages", () => {
  for (const pageId of TARGETS) {
    it(`${pageId} renders its page-specific UI`, async () => {
      const container = await renderPage(pageId);
      const text = container.textContent ?? "";
      expect(text.trim().length, `${pageId} must render content`).toBeGreaterThan(0);
      expect(container.querySelector('[aria-labelledby="schema-error-title"]')).toBeNull();
      expect(text).not.toContain("unknown custom component");
    });
  }

  it("mail renders the outbound-mail console surface (mail-admin-tab)", async () => {
    // A-002 F-002 (independent · GOAL-042): assert the custom surface really
    // issues its data request (fetch spy) — text assertions alone could pass
    // if the custom component short-circuited.
    const fetcher = vi.fn(fixtureFetcher("mail"));
    const container = await renderPage("mail", fetcher);
    const text = container.textContent ?? "";
    // Either the config form loaded or the fail-closed alert shows — the
    // console heading itself must be present either way.
    expect(text).toMatch(/Outbound mail|Could not load the outbound-mail configuration/);
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(fetcher).toHaveBeenCalledWith(expect.stringContaining("/api/mail/config"));
  });

  it("mail-outbox renders the outbox table headers", async () => {
    const container = await renderPage("mail-outbox");
    const text = container.textContent ?? "";
    expect(text).toContain("Subject");
    expect(text).toContain("Delivery status");
  });

  it("my-wallet renders the wallet surface (ensure probe + balance intro)", async () => {
    const container = await renderPage("my-wallet");
    const text = container.textContent ?? "";
    expect(text).toContain("Your wallet balance and recent ledger entries");
    // Either the wallet exists (stat cards) or the lazy-open surface shows.
    expect(text).toMatch(/Open wallet|Total balance/);
  });

  it("wallet renders the account search/toolbar surface", async () => {
    const container = await renderPage("wallet");
    const text = container.textContent ?? "";
    expect(text).toContain("Owner type");
    expect(text).toContain("New account");
    expect(text).toContain("Reconcile");
  });

  it("wallet-vouchers renders the voucher table headers and generate action", async () => {
    const container = await renderPage("wallet-vouchers");
    const text = container.textContent ?? "";
    expect(text).toContain("Redeemed");
    expect(text).toContain("Generate vouchers");
  });

  it("telegram-settings renders the Telegram channel surface", async () => {
    const container = await renderPage("telegram-settings");
    const text = container.textContent ?? "";
    expect(text).toContain("Telegram channel");
  });

  it("telegram-operator renders the operator conversations surface", async () => {
    const container = await renderPage("telegram-operator");
    const text = container.textContent ?? "";
    expect(text).toContain("Operator conversations");
  });

  it("digitaloffer-offers renders the product table and create action", async () => {
    const container = await renderPage("digitaloffer-offers");
    const text = container.textContent ?? "";
    expect(text).toContain("Name");
    expect(text).toContain("New product");
  });

  it("digitaloffer-entitlements renders the entitlements table with void row action", async () => {
    const container = await renderPage("digitaloffer-entitlements");
    const text = container.textContent ?? "";
    expect(text).toContain("Remaining");
    expect(text).toContain("Void");
  });

  it("digitaloffer-purchases renders the orders table headers", async () => {
    const container = await renderPage("digitaloffer-purchases");
    const text = container.textContent ?? "";
    expect(text).toContain("Purchased at");
    expect(text).toContain("Request ID");
  });
});
