// @vitest-environment jsdom
//
// GOAL-044 W32 (D-001 §3) · the idle-poll decision, at the component boundary.
//
// The page-level test in renderer/jobs-result-center.test.tsx drives this through
// the real jobs document; this one drives it through a stub CRUD context so the
// branch that a page cannot easily reach — "the table has not published its rows
// yet" — is pinned too. It must refresh conservatively there: a control that
// stops polling because it cannot SEE the rows would silently stop updating a
// busy list (the failure mode is invisible), while refreshing one extra time is
// merely one extra request.

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { JobsAutoRefresh } from "@/components/jobs-auto-refresh";
import { SchemaCrudContext, type SchemaCrudValue } from "@/renderer/render.tsx";
import type { RenderCustomNode } from "@/renderer/render.types";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  vi.useFakeTimers();
});

afterEach(async () => {
  vi.useRealTimers();
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

function node(): RenderCustomNode {
  return {
    type: "custom",
    id: "jobs-auto-refresh",
    component: "jobs-auto-refresh",
    props: { targetTable: "jobs-table", statusField: "status", activeStatuses: ["queued", "running"] },
  } as unknown as RenderCustomNode;
}

async function renderWithRows(
  rows: ReadonlyArray<Record<string, unknown>> | undefined,
): Promise<{ refreshTable: ReturnType<typeof vi.fn>; container: HTMLDivElement }> {
  const refreshTable = vi.fn();
  // Only the two seams this component consumes are needed; the rest of the page
  // context is irrelevant here, so the stub is intentionally partial.
  const crud = {
    refreshTable,
    tableRows: () => rows,
    selection: () => undefined,
  } as unknown as SchemaCrudValue;

  const container = globalThis.document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider stored="en-US">
        <SchemaCrudContext.Provider value={crud}>
          <JobsAutoRefresh node={node()} context={{}} />
        </SchemaCrudContext.Provider>
      </I18nProvider>,
    );
  });
  const select = container.querySelector("[data-jobs-refresh-select]") as HTMLSelectElement;
  await act(async () => {
    select.value = "5000";
    select.dispatchEvent(new Event("change", { bubbles: true }));
  });
  return { refreshTable, container };
}

async function tick(ms: number): Promise<void> {
  await act(async () => {
    vi.advanceTimersByTime(ms);
  });
  await act(async () => {
    await Promise.resolve();
  });
}

describe("W32 · jobs auto-refresh idle decision", () => {
  it("refreshes while a row is in an active status", async () => {
    const { refreshTable } = await renderWithRows([{ status: "running" }]);
    await tick(5000);
    expect(refreshTable).toHaveBeenCalledTimes(1);
    expect(refreshTable).toHaveBeenCalledWith("jobs-table");
  });

  it("skips the tick once every row is terminal", async () => {
    const { refreshTable } = await renderWithRows([{ status: "succeeded" }, { status: "failed" }]);
    await tick(15000);
    expect(refreshTable).not.toHaveBeenCalled();
  });

  it("refreshes conservatively while the table has published no rows", async () => {
    const { refreshTable } = await renderWithRows(undefined);
    await tick(5000);
    expect(refreshTable).toHaveBeenCalledTimes(1);
  });

  it("stays idle while the control is off", async () => {
    const { refreshTable, container } = await renderWithRows([{ status: "running" }]);
    const select = container.querySelector("[data-jobs-refresh-select]") as HTMLSelectElement;
    await act(async () => {
      select.value = "0";
      select.dispatchEvent(new Event("change", { bubbles: true }));
    });
    await tick(30000);
    expect(refreshTable).not.toHaveBeenCalled();
  });
});
