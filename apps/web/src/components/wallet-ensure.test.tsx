// @vitest-environment jsdom
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { WalletEnsure } from "@/components/wallet-ensure";
import { I18nProvider } from "@/i18n/runtime";

const { authFetchMock, reloadListMock } = vi.hoisted(() => ({
  authFetchMock: vi.fn(),
  reloadListMock: vi.fn(),
}));
vi.mock("@/account/AuthContext", () => ({
  useAuth: () => ({ authFetch: authFetchMock }),
}));
vi.mock("@/renderer/render.tsx", () => ({
  useSchemaCrud: () => ({ reloadList: reloadListMock }),
}));

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  vi.clearAllMocks();
});

async function renderEnsure() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <WalletEnsure node={{ type: "custom", component: "wallet-ensure" }} context={{}} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("WalletEnsure", () => {
  it("POSTs /api/wallet/me on mount and reloads lists", async () => {
    authFetchMock.mockResolvedValue(new Response("{}", { status: 200 }));
    await renderEnsure();
    expect(authFetchMock).toHaveBeenCalledWith(
      "/api/wallet/me",
      expect.objectContaining({ method: "POST" }),
    );
    expect(reloadListMock).toHaveBeenCalled();
  });

  it("shows a retry CTA when open fails", async () => {
    authFetchMock.mockResolvedValue(new Response("nope", { status: 500 }));
    const container = await renderEnsure();
    expect(container.querySelector("[role='alert']")).not.toBeNull();
    expect(reloadListMock).not.toHaveBeenCalled();
  });
});
