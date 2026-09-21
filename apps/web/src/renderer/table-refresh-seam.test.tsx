// @vitest-environment jsdom
//
// GOAL-044 W32 (D-001 §2) · the three refresh seams and what each does to the
// page's table selections.
//
//   reloadList()            full reload wave, CLEARS every selection (ADR-0022 D2)
//   refreshList(dataSource) display dataSources only (W25) — untouched here
//   refreshTable(tableId)   THIS table, current query, KEEPS the selection
//
// The third is new and exists so a polling control can refresh a list without
// deleting the selection an operator is working with. This test pins all three
// behaviours against one real page so the distinction cannot quietly collapse:
// if `refreshTable` ever started clearing selections, case 2 fails; if
// `reloadList` stopped clearing them, case 3 fails (that one is the frozen
// ADR-0022 contract, not a preference).

import { act, useState } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import type { RenderPageDocument } from "@/renderer/render.types";
import { registerCustomComponent, resetCustomComponentsForTests } from "@/renderer/custom-components";
import { RenderPage, useSchemaCrud } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";

const TABLE_ID = "seam-table";

interface Harness {
  container: HTMLDivElement;
  requests: () => number;
  click: (label: string) => Promise<void>;
  text: () => string;
}

function SeamProbe() {
  const crud = useSchemaCrud();
  const selection = crud?.selection(TABLE_ID);
  // The row registry is ref-backed on purpose (publishing rows must not
  // re-render the page), so the probe reads it on demand — exactly how a
  // polling control consumes it: inside a callback, never during render.
  const [observedRows, setObservedRows] = useState(-1);
  return (
    <div data-seam-probe>
      <span data-seam-selection>{String(selection?.count ?? 0)}</span>
      <span data-seam-rows>{String(observedRows)}</span>
      <button type="button" onClick={() => crud?.setSelection(TABLE_ID, ["1", "2"])}>
        select
      </button>
      <button type="button" onClick={() => crud?.refreshTable(TABLE_ID)}>
        refresh-table
      </button>
      <button type="button" onClick={() => crud?.refreshTable("no-such-table")}>
        refresh-unknown
      </button>
      <button type="button" onClick={() => crud?.reloadList()}>
        reload
      </button>
      <button type="button" onClick={() => setObservedRows(crud?.tableRows(TABLE_ID)?.length ?? -1)}>
        read-rows
      </button>
    </div>
  );
}

function pageDocument(): RenderPageDocument {
  return {
    meta: {
      pageId: "seam",
      title: "Seam",
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "table.selection"],
    },
    body: {
      type: "section",
      id: "seam",
      children: [
        { type: "custom", id: "seam-probe", component: "seam-probe" },
        {
          type: "table",
          id: TABLE_ID,
          props: {
            columns: [{ field: "id", label: "ID" }],
            dataSource: "/api/things",
            selection: { mode: "multiple" },
          },
        },
      ],
    },
  } as unknown as RenderPageDocument;
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  resetCustomComponentsForTests();
  registerCustomComponent("seam-probe", SeamProbe);
});

afterEach(async () => {
  resetCustomComponentsForTests();
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

async function renderSeamPage(): Promise<Harness> {
  let requests = 0;
  const fetcher = (async (input: RequestInfo | URL) => {
    const url = String(input);
    if (url.startsWith("/api/things")) {
      requests += 1;
      return new Response(
        JSON.stringify({
          items: [
            { id: "1", label: "one" },
            { id: "2", label: "two" },
            { id: "3", label: "three" },
          ],
          total: 3,
          page: 1,
          pageSize: 20,
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    return new Response(JSON.stringify({ items: [], total: 0, page: 1, pageSize: 20 }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  }) as typeof fetch;

  const container = globalThis.document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider stored="en-US">
        <RenderPage
          document={pageDocument()}
          context={{ user: { permissions: ["things.read"] } }}
          dataFetcher={fetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
        />
      </I18nProvider>,
    );
  });
  await act(async () => {
    await Promise.resolve();
  });

  const click = async (label: string) => {
    const button = Array.from(container.querySelectorAll("button")).find(
      (candidate) => (candidate.textContent ?? "").trim() === label,
    );
    expect(button, `button ${label}`).toBeDefined();
    await act(async () => {
      button!.dispatchEvent(new MouseEvent("click", { bubbles: true }));
    });
    await act(async () => {
      await Promise.resolve();
    });
  };

  return {
    container,
    requests: () => requests,
    click,
    text: () => container.textContent ?? "",
  };
}

function selectionCount(container: HTMLElement): string {
  return container.querySelector("[data-seam-selection]")?.textContent ?? "";
}

describe("W32 · refresh seams vs table selection", () => {
  it("refreshTable refetches the table and KEEPS the selection", async () => {
    const harness = await renderSeamPage();
    await harness.click("select");
    expect(selectionCount(harness.container)).toBe("2");

    const before = harness.requests();
    await harness.click("refresh-table");
    expect(harness.requests()).toBeGreaterThan(before);
    expect(selectionCount(harness.container)).toBe("2");
  });

  it("reloadList still CLEARS every selection (ADR-0022 D2 unchanged)", async () => {
    const harness = await renderSeamPage();
    await harness.click("select");
    expect(selectionCount(harness.container)).toBe("2");

    const before = harness.requests();
    await harness.click("reload");
    expect(harness.requests()).toBeGreaterThan(before);
    expect(selectionCount(harness.container)).toBe("0");
  });

  it("publishes the table's rows for idle decisions and tolerates unknown tables", async () => {
    const harness = await renderSeamPage();
    await harness.click("read-rows");
    expect(harness.container.querySelector("[data-seam-rows]")?.textContent).toBe("3");

    // An unknown table id is a no-op, not a crash and not a page-wide refresh.
    const before = harness.requests();
    await harness.click("refresh-unknown");
    expect(harness.requests()).toBe(before);
    expect(selectionCount(harness.container)).toBe("0");
  });
});
