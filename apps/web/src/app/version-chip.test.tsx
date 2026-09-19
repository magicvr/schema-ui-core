// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { hasMonitoringRead, QUICKSTART_UPGRADE_URL, VersionChip } from "@/app/version-chip";
import { I18nProvider } from "@/i18n/runtime";

const active: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", { configurable: true, value: true });
});

afterEach(async () => {
  for (const { root, container } of active.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

async function renderChip(props: Parameters<typeof VersionChip>[0]): Promise<HTMLDivElement> {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  active.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider stored="en-US">
        <VersionChip {...props} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("VersionChip", () => {
  it("does not fetch or render without monitoring.read", async () => {
    const fetcher = vi.fn();
    const container = await renderChip({ canReadMonitoring: false, fetcher: fetcher as unknown as typeof fetch });
    expect(fetcher).not.toHaveBeenCalled();
    expect(container.querySelector("[data-version-chip]")).toBeNull();
    expect(hasMonitoringRead(["users.read"])).toBe(false);
    expect(hasMonitoringRead(["monitoring.read"])).toBe(true);
  });

  it("shows version (not commit) and the QUICKSTART upgrade link", async () => {
    const fetcher = vi.fn(async () =>
      new Response(
        JSON.stringify({
          items: [{ version: "1.2.3", commit: "deadbeef", modules: ["admin.users"] }],
          total: 1,
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      ),
    );
    const container = await renderChip({
      canReadMonitoring: true,
      fetcher: fetcher as unknown as typeof fetch,
    });
    await act(async () => {
      await Promise.resolve();
    });
    expect(fetcher).toHaveBeenCalled();
    expect(String(fetcher.mock.calls[0]?.[0])).toContain("/api/system-monitoring/status");
    expect(container.querySelector("[data-version-chip-value]")?.textContent).toBe("1.2.3");
    expect(container.textContent).not.toContain("deadbeef");
    const upgrade = container.querySelector(`a[href="${QUICKSTART_UPGRADE_URL}"]`);
    expect(upgrade).not.toBeNull();
  });

  it("calls onOpenDiagnostics when the diagnostics control is used", async () => {
    const fetcher = vi.fn(async () =>
      new Response(JSON.stringify({ items: [{ version: "1.0.0" }], total: 1 }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }),
    );
    const onOpenDiagnostics = vi.fn();
    const container = await renderChip({
      canReadMonitoring: true,
      fetcher: fetcher as unknown as typeof fetch,
      onOpenDiagnostics,
    });
    await act(async () => {
      await Promise.resolve();
    });
    const button = Array.from(container.querySelectorAll("button")).find((node) =>
      node.textContent === "Diagnostics",
    );
    expect(button).toBeDefined();
    await act(async () => {
      button?.click();
    });
    expect(onOpenDiagnostics).toHaveBeenCalledTimes(1);
  });

  it("hides the chip when status fetch fails", async () => {
    const fetcher = vi.fn(async () => new Response("", { status: 403 }));
    const container = await renderChip({
      canReadMonitoring: true,
      fetcher: fetcher as unknown as typeof fetch,
    });
    await act(async () => {
      await Promise.resolve();
    });
    expect(container.querySelector("[data-version-chip]")).toBeNull();
  });
});
