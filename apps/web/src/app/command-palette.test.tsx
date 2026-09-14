// @vitest-environment jsdom

import { act, useState } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { CommandPalette } from "@/app/CommandPalette";
import { I18nProvider } from "@/i18n/runtime";
import {
  SEARCHABLE_PROVIDER_VERSION,
  type SearchableProvider,
  type SearchableProviderContext,
} from "@/app/searchable";
import { validateAppManifest } from "@/protocol/app-manifest";

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
});

function context(): SearchableProviderContext {
  const manifest = validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "palette", name: "Palette", homePageRef: "home" },
    pages: [{ pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/home" }],
  });
  return {
    manifest,
    navigationContext: { user: { id: "u1" }, features: {} },
    currentPath: "/home",
    t: (key, _params, fallback) => fallback ?? key,
    loadPage: async () => ({ body: { type: "section", children: [] } }),
  };
}

const provider: SearchableProvider = {
  id: "test",
  version: SEARCHABLE_PROVIDER_VERSION,
  getItems: () => ({
    items: [
      { id: "page:home", kind: "page", label: "Home", href: "/home", rank: 0 },
      { id: "action:create", kind: "action", label: "Create user", href: "/home", rank: 1 },
      { id: "page:settings", kind: "page", label: "Settings", keywords: ["配置"], href: "/home", rank: 2 },
    ],
  }),
};

function render(element: React.ReactElement): HTMLDivElement {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  act(() => root.render(element));
  return container;
}

async function flush(): Promise<void> {
  await act(async () => {
    await Promise.resolve();
  });
}

describe("CommandPalette", () => {
  it("loads providers, exposes combobox/listbox semantics, and selects with arrows/Enter", async () => {
    const onSelect = vi.fn();
    const container = render(
      <I18nProvider stored="en-US">
        <CommandPalette
          open
          providers={[provider]}
          context={context()}
          onClose={vi.fn()}
          onSelect={onSelect}
        />
      </I18nProvider>,
    );
    await flush();

    const input = container.querySelector<HTMLInputElement>('[role="combobox"]');
    const listbox = container.querySelector('[role="listbox"]');
    expect(input).not.toBeNull();
    expect(document.activeElement).toBe(input);
    expect(input?.getAttribute("aria-controls")).toBe(listbox?.id);
    expect(container.querySelectorAll('[role="option"]')).toHaveLength(3);

    await act(async () => {
      input?.dispatchEvent(new KeyboardEvent("keydown", { key: "ArrowDown", bubbles: true }));
    });
    await act(async () => {
      input?.dispatchEvent(new KeyboardEvent("keydown", { key: "Enter", bubbles: true }));
    });
    expect(onSelect).toHaveBeenCalledWith(expect.objectContaining({ id: "action:create" }));
  });

  it("matches localized keywords and closes on Escape while restoring the trigger focus", async () => {
    function Harness() {
      const [open, setOpen] = useState(false);
      return (
        <>
          <button type="button" data-trigger onClick={() => setOpen(true)}>
            Trigger
          </button>
          <CommandPalette
            open={open}
            providers={[provider]}
            context={context()}
            onClose={() => setOpen(false)}
            onSelect={vi.fn()}
          />
        </>
      );
    }

    const container = render(
      <I18nProvider stored="en-US">
        <Harness />
      </I18nProvider>,
    );
    const trigger = container.querySelector<HTMLButtonElement>("[data-trigger]")!;
    await act(async () => {
      trigger.focus();
      trigger.click();
    });
    await flush();
    const input = container.querySelector<HTMLInputElement>('[role="combobox"]')!;
    expect(document.activeElement).toBe(input);

    await act(async () => {
      Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set?.call(
        input,
        "pei",
      );
      input.dispatchEvent(new Event("input", { bubbles: true }));
    });
    expect(container.querySelectorAll('[role="option"]')).toHaveLength(0);

    await act(async () => {
      input.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape", bubbles: true }));
    });
    expect(container.querySelector('[role="dialog"]')).toBeNull();
    expect(document.activeElement).toBe(trigger);
  });
});
