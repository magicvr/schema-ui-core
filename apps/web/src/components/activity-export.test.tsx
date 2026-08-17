// @vitest-environment jsdom
//
// W14 F-03 (GOAL-016): the Activity page export button downloads the CSV from
// /api/operations/export with the current table query.

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { ActivityExport } from "@/components/activity-export";
import { I18nProvider } from "@/i18n/runtime";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", { configurable: true, value: true });
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  vi.unstubAllGlobals();
});

const NODE = {
  type: "custom" as const,
  id: "activity-export",
  component: "activity-export",
  props: { targetTable: "operations-table" },
};

function renderExport() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  act(() => {
    root.render(
      <I18nProvider>
        <ActivityExport node={NODE as never} context={{}} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("ActivityExport (W14 F-03)", () => {
  it("downloads the operations CSV through the authed fetcher", async () => {
    const createObjectURL = vi.fn(() => "blob:mock");
    const revokeObjectURL = vi.fn();
    const click = vi.fn();
    vi.stubGlobal("URL", { ...URL, createObjectURL, revokeObjectURL });
    const originalCreateElement = document.createElement.bind(document);
    vi.spyOn(document, "createElement").mockImplementation((tag: string) => {
      const element = originalCreateElement(tag);
      if (tag === "a") {
        return Object.assign(element, { click });
      }
      return element;
    });

    const fetchMock = vi.fn(async () => new Response(new Blob(["id,event\nop-1,auth.login"], { type: "text/csv" }), { status: 200 }));
    vi.stubGlobal("fetch", fetchMock);

    const container = renderExport();
    const button = container.querySelector("button") as HTMLButtonElement;
    await act(async () => {
      button.click();
      await new Promise((resolve) => setTimeout(resolve, 30));
    });

    expect(fetchMock).toHaveBeenCalledWith(
      "/api/operations/export?pageSize=10000",
      expect.objectContaining({ headers: { Accept: "text/csv" } }),
    );
    expect(click).toHaveBeenCalled();
    expect(createObjectURL).toHaveBeenCalled();
  });
});
