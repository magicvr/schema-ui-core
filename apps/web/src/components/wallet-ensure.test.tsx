// @vitest-environment jsdom
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { WalletEnsure } from "@/components/wallet-ensure";
import { I18nProvider } from "@/i18n/runtime";
import { ResourceApiError } from "@/renderer/resource";

const { authFetchMock, reloadListMock, fetchListMock } = vi.hoisted(() => ({
  authFetchMock: vi.fn(),
  reloadListMock: vi.fn(),
  fetchListMock: vi.fn(),
}));
vi.mock("@/account/AuthContext", () => ({
  useAuth: () => ({ authFetch: authFetchMock }),
}));
vi.mock("@/renderer/render.tsx", () => ({
  useSchemaCrud: () => ({ reloadList: reloadListMock, fetchList: fetchListMock }),
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
  it("probes via the shared fetch cache and neither POSTs nor reloads when the wallet exists", async () => {
    fetchListMock.mockResolvedValue({
      items: [{ balanceTotal: 1 }],
      total: 1,
      page: 1,
      pageSize: 100,
    });
    await renderEnsure();
    expect(fetchListMock).toHaveBeenCalledWith(
      "/api/wallet/me",
      expect.objectContaining({ page: 1, pageSize: 100 }),
      undefined,
      expect.any(Function),
    );
    expect(authFetchMock).not.toHaveBeenCalled();
    expect(reloadListMock).not.toHaveBeenCalled();
  });

  it("POSTs to create the wallet and reloads lists only when the probe reports a missing wallet", async () => {
    fetchListMock.mockRejectedValue(new ResourceApiError(404, "WALLET_NOT_FOUND", "wallet account not found"));
    authFetchMock.mockResolvedValue(new Response("{}", { status: 200 }));
    await renderEnsure();
    expect(authFetchMock).toHaveBeenCalledWith(
      "/api/wallet/me",
      expect.objectContaining({ method: "POST" }),
    );
    expect(reloadListMock).toHaveBeenCalled();
  });

  it("shows a retry CTA when creating fails", async () => {
    fetchListMock.mockRejectedValue(new ResourceApiError(404, "WALLET_NOT_FOUND", "wallet account not found"));
    authFetchMock.mockResolvedValue(new Response("nope", { status: 500 }));
    const container = await renderEnsure();
    expect(container.querySelector("[role='alert']")).not.toBeNull();
    expect(reloadListMock).not.toHaveBeenCalled();
  });

  it("shows a retry CTA when the probe fails with a non-wallet error", async () => {
    fetchListMock.mockRejectedValue(new Error("network down"));
    const container = await renderEnsure();
    expect(container.querySelector("[role='alert']")).not.toBeNull();
    expect(authFetchMock).not.toHaveBeenCalled();
    expect(reloadListMock).not.toHaveBeenCalled();
  });
});