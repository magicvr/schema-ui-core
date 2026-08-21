// @vitest-environment jsdom
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { ImportTemplateDownload } from "@/components/import-template-download";
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
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  vi.clearAllMocks();
});

async function renderDownload() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <ImportTemplateDownload
          node={{ type: "custom", component: "import-template-download" }}
          context={{}}
        />
      </I18nProvider>,
    );
  });
  return container;
}

describe("ImportTemplateDownload", () => {
  it("shows an error when the template request fails", async () => {
    authFetchMock.mockResolvedValue(new Response("nope", { status: 500 }));
    const container = await renderDownload();
    const button = container.querySelector<HTMLButtonElement>("[data-import-template-download]");
    expect(button).not.toBeNull();
    await act(async () => {
      button!.click();
    });
    expect(container.querySelector("[data-import-template-error]")).not.toBeNull();
  });
});
