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

    // A-003: active filter chips appear for non-empty conditions.
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    const chips = container.querySelector('[data-filter-chips]');
    expect(chips).not.toBeNull();
    expect(chips!.textContent).toContain("ali");
    expect(chips!.textContent).toContain("Enabled");

    // A-003 pairing rule (user 2026-08-16): the search submit button is
    // adjacent to the keyword input — same grid cell, side by side.
    const searchButton = Array.from(container.querySelectorAll('button[type="submit"]')).find(
      (el) => el.textContent?.includes("Search"),
    );
    expect(searchButton).not.toBeUndefined();
    const cell = (searchButton as HTMLButtonElement).parentElement!;
    expect(cell.contains(container.querySelector('input') as Node)).toBe(true);
    // W13 T-03 (user 2026-08-16): the pair is ONE attached component — the
    // button overlaps the input's right border (-ml-px) and the inner
    // corners are squared off, so they read as a single control and can
    // never wrap onto separate rows.
    expect((searchButton as HTMLButtonElement).className).toContain("-ml-px");
    expect((searchButton as HTMLButtonElement).className).toContain("rounded-l-none");
    const searchInput = container.querySelector("input") as HTMLInputElement;
    expect(searchInput.className).toContain("rounded-r-none");

    // A-003: the reset button clears every condition and re-runs the search
    // (the request drops q and enabled).
    const resetButton = Array.from(container.querySelectorAll("button")).find((el) =>
      el.textContent?.includes("Reset"),
    );
    expect(resetButton).not.toBeUndefined();
    const before = calls.length;
    await act(async () => {
      (resetButton as HTMLButtonElement).click();
    });
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    const afterReset = calls[calls.length - 1];
    expect(calls.length).toBeGreaterThan(before);
    expect(afterReset).not.toContain("q=");
    expect(afterReset).not.toContain("enabled=");
    expect(container.querySelector('[data-filter-chips]')).toBeNull();
  });
});
