// @vitest-environment jsdom
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { CronPreview } from "@/components/cron-preview";
import { I18nProvider } from "@/i18n/runtime";

const { authFetchMock } = vi.hoisted(() => ({
  authFetchMock: vi.fn(),
}));
vi.mock("@/account/AuthContext", () => ({
  useAuth: () => ({ authFetch: authFetchMock }),
}));

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  vi.useFakeTimers();
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  vi.useRealTimers();
  vi.clearAllMocks();
});

async function renderPreview(bindValue?: string) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <CronPreview
          node={{
            type: "custom",
            component: "cron-preview",
            ...(bindValue !== undefined ? { props: { bindValue } } : {}),
          }}
          context={{}}
        />
      </I18nProvider>,
    );
  });
  return container;
}

describe("CronPreview", () => {
  it("bound mode hides the standalone input and previews the field value", async () => {
    authFetchMock.mockResolvedValue(
      new Response(JSON.stringify({ description: "每天 02:00", nextRuns: ["2026-08-19T02:00:00Z"] }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }),
    );
    const container = await renderPreview("0 2 * * *");
    expect(container.querySelector("#cronPreviewInput")).toBeNull();
    expect(container.querySelector("[data-cron-bound='true']")).not.toBeNull();

    await act(async () => {
      await vi.advanceTimersByTimeAsync(400);
    });
    expect(authFetchMock).toHaveBeenCalledWith(
      "/api/scheduled-tasks/cron/preview",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ cron: "0 2 * * *" }),
      }),
    );
    expect(container.textContent).toContain("每天 02:00");
  });

  it("standalone mode keeps an independent input", async () => {
    const container = await renderPreview();
    expect(container.querySelector("#cronPreviewInput")).not.toBeNull();
    expect(container.querySelector("[data-cron-bound='true']")).toBeNull();
  });
});
