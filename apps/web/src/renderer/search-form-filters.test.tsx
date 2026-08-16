// @vitest-environment jsdom
//
// T-02 (GOAL-013 D-003): search-mode forms bind every non-q field into the
// target table's query filters (serialized to URL params). A select filter is
// sent as its own param; clearing it removes the param; table-level props
// filters are preserved.

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table.tsx";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", { configurable: true, value: true });
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

// Document: search form (q + enabled select) targeting a table.
function pageDocument(): unknown {
  return {
    meta: {
      pageId: "users",
      title: "Users",
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation", "form.controls.extended"],
    },
    body: {
      type: "section",
      children: [
        {
          type: "form",
          id: "users-search",
          props: {
            mode: "search",
            targetTable: "users-table",
            fields: [
              { id: "q", type: "input" },
              {
                id: "enabled",
                type: "select",
                options: [
                  { value: "", label: "All" },
                  { value: "true", label: "Enabled" },
                  { value: "false", label: "Disabled" },
                ],
              },
            ],
            submitLabel: "Search",
          },
        },
        {
          type: "table",
          id: "users-table",
          props: {
            columns: [
              { field: "id", label: "ID" },
              { field: "username", label: "Username" },
            ],
            dataSource: "/api/users",
          },
        },
      ],
    },
  };
}

function renderSurface(fetcher: typeof fetch): HTMLDivElement {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  act(() => {
    root.render(
      <I18nProvider>
        <RenderPage
          document={pageDocument() as never}
          context={{}}
          dataFetcher={fetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
        />
      </I18nProvider>,
    );
  });
  return container;
}

describe("T-02 search form filter binding (GOAL-013 D-003)", () => {
  it("sends select filters as query params and drops cleared ones", async () => {
    const calls: string[] = [];
    const fetcher = (async (input: RequestInfo | URL) => {
      calls.push(String(input));
      return new Response(
        JSON.stringify({ items: [{ id: "u1", username: "alice" }], total: 1, page: 1, pageSize: 10 }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }) as typeof fetch;

    const container = renderSurface(fetcher);
    // Wait for the initial table fetch.
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(calls.length).toBeGreaterThanOrEqual(1);

    const form = container.querySelector('form');
    expect(form).not.toBeNull();

    // Type a keyword and pick enabled=true from the select.
    const qInput = form!.querySelector('input');
    const select = form!.querySelector('select');
    await act(async () => {
      const setter = Object.getOwnPropertyDescriptor(
        window.HTMLInputElement.prototype,
        "value",
      )!.set!;
      setter.call(qInput, "ali");
      qInput!.dispatchEvent(new Event("input", { bubbles: true }));
      const selectSetter = Object.getOwnPropertyDescriptor(
        window.HTMLSelectElement.prototype,
        "value",
      )!.set!;
      selectSetter.call(select, "true");
      select!.dispatchEvent(new Event("change", { bubbles: true }));
    });
    await act(async () => {
      form!.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true }));
    });
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    const last = calls[calls.length - 1];
    expect(last).toContain("q=ali");
    expect(last).toContain("enabled=true");
  });
});
