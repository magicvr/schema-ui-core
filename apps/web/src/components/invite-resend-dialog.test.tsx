// @vitest-environment jsdom
// Resend dialog coverage (workspace-019, modal action): the dialog reads the
// triggering invite from context.modalRow, rotates the token via the admin
// API, discloses the one-time link inline, copies it on demand, and finishes
// by closing the modal plus reloading the list.
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { InviteResendDialog } from "@/components/invite-resend-dialog";
import { I18nProvider } from "@/i18n/runtime";
import { SchemaCrudContext } from "@/renderer/render.tsx";

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
  vi.unstubAllGlobals();
});

async function renderDialog(row: Record<string, unknown> | null, onClose?: () => void, onReload?: () => void) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  vi.stubGlobal("fetch", () =>
    Promise.resolve(
      new Response(JSON.stringify({ link: "/invite/accept?token=rotated-77" }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }),
    ),
  );
  await act(async () => {
    root.render(
      <I18nProvider>
        <SchemaCrudContext.Provider
          value={
            {
              fetcher: globalThis.fetch,
              closeModal: onClose ?? (() => undefined),
              reloadList: onReload ?? (() => undefined),
            } as never
          }
        >
          <InviteResendDialog
            node={{ type: "custom", component: "invite-resend-dialog" } as never}
            context={{ modalRow: row }}
            children={[]}
          />
        </SchemaCrudContext.Provider>
      </I18nProvider>,
    );
  });
  return container;
}

async function settled() {
  await act(async () => {
    await new Promise((resolve) => setTimeout(resolve, 0));
  });
}

describe("InviteResendDialog", () => {
  it("renders the invite summary and resends, disclosing the rotated link", async () => {
    const container = await renderDialog({ id: "inv-42", email: "a@example.com", roles: ["viewer"] });
    expect(container.querySelector("[data-invite-resend-dialog]")).not.toBeNull();
    expect(container.textContent).toContain("inv-42");

    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-resend-send]")!.click();
    });
    await settled();
    expect(container.querySelector("[data-resend-link]")?.textContent).toContain("/invite/accept?token=rotated-77");
  });

  it("copies the disclosed link to the clipboard", async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    vi.stubGlobal("navigator", { ...navigator, clipboard: { writeText } });
    const container = await renderDialog({ id: "inv-42", email: "", roles: [] });
    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-resend-send]")!.click();
    });
    await settled();
    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-resend-copy]")!.click();
    });
    await settled();
    expect(writeText).toHaveBeenCalledTimes(1);
    expect(String(writeText.mock.calls[0]?.[0] ?? "")).toBe("/invite/accept?token=rotated-77");
    expect(container.querySelector("[data-copied-hint]")).not.toBeNull();
  });

  it("closes the modal and reloads the list on done", async () => {
    const onClose = vi.fn();
    const onReload = vi.fn();
    const container = await renderDialog({ id: "inv-42", email: "", roles: [] }, onClose, onReload);
    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-resend-send]")!.click();
    });
    await settled();
    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-resend-done]")!.click();
    });
    await settled();
    expect(onClose).toHaveBeenCalledTimes(1);
    expect(onReload).toHaveBeenCalledTimes(1);
  });

  it("fails closed without row context", async () => {
    const container = await renderDialog(null);
    expect(container.querySelector<HTMLButtonElement>("[data-resend-send]")!.disabled).toBe(true);
    expect(container.textContent).toContain("No invitation row context.");
  });
});