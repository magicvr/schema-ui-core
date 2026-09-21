// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { hasDirtyState, resetDirtyStateSources } from "@/renderer/dirty-state";
import { RenderPage } from "@/renderer/render.tsx";
import type { RenderPageDocument } from "@/renderer/render.types";

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
  resetDirtyStateSources();
  vi.restoreAllMocks();
});

function requestDocument(): RenderPageDocument {
  return {
    meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest", "form.controls.advanced"] },
    actions: {
      save: { type: "request", method: "POST", url: "/api/save" },
    },
    body: {
      type: "form",
      id: "profile-form",
      props: {
        fields: [{ id: "name", label: "Name", type: "input", defaultValue: "Alice" }],
        submitAction: "save",
        submitLabel: "Save",
      },
    },
  } as RenderPageDocument;
}

function resetDocument(): RenderPageDocument {
  return {
    meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest", "form.controls.advanced"] },
    body: {
      type: "form",
      id: "profile-form",
      props: {
        fields: [{ id: "name", label: "Name", type: "input", defaultValue: "Alice" }],
      },
    },
  } as RenderPageDocument;
}

function modalDocument(): RenderPageDocument {
  return {
    meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest", "form.controls.advanced"] },
    actions: {
      open: {
        type: "modal",
        content: {
          type: "form",
          id: "modal-form",
          props: {
            fields: [{ id: "name", label: "Name", type: "input", defaultValue: "Alice" }],
          },
        },
      },
    },
    body: {
      type: "actionButton",
      id: "open-form",
      props: { actionId: "open", label: "Open form" },
    },
  } as unknown as RenderPageDocument;
}

function modalSubmitDocument(): RenderPageDocument {
  return {
    meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest", "form.controls.advanced"] },
    actions: {
      open: {
        type: "modal",
        content: {
          type: "form",
          id: "modal-form",
          props: {
            fields: [{ id: "name", label: "Name", type: "input", defaultValue: "Alice" }],
            submitAction: "save",
            submitLabel: "Save",
          },
        },
      },
      save: { type: "request", method: "POST", url: "/api/save" },
    },
    body: {
      type: "actionButton",
      id: "open-form",
      props: { actionId: "open", label: "Open form" },
    },
  } as unknown as RenderPageDocument;
}

function validationDocument(): RenderPageDocument {
  return {
    meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest", "form.controls.advanced"] },
    actions: {
      save: { type: "request", method: "POST", url: "/api/save" },
    },
    body: {
      type: "form",
      id: "profile-form",
      props: {
        fields: [{ id: "name", label: "Name", type: "input", defaultValue: "Alice", required: true }],
        submitAction: "save",
        submitLabel: "Save",
      },
    },
  } as RenderPageDocument;
}

function searchDocument(): RenderPageDocument {
  return {
    meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
    body: {
      type: "form",
      id: "search-form",
      props: {
        mode: "search",
        targetTable: "users",
        fields: [{ id: "q", label: "Search", type: "input" }],
        submitLabel: "Search",
      },
    },
  } as RenderPageDocument;
}

async function renderDocument(
  document: RenderPageDocument,
  dataFetcher?: typeof fetch,
): Promise<HTMLDivElement> {
  const container = documentElement();
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <RenderPage document={document} context={{}} dataFetcher={dataFetcher} />
      </I18nProvider>,
    );
  });
  return container;
}

function documentElement(): HTMLDivElement {
  const container = document.createElement("div");
  document.body.appendChild(container);
  return container;
}

function setInputValue(input: HTMLInputElement, value: string): void {
  const setter = Object.getOwnPropertyDescriptor(HTMLInputElement.prototype, "value")?.set;
  setter?.call(input, value);
  input.dispatchEvent(new Event("input", { bubbles: true }));
}

describe("R3 dirty form lifecycle", () => {
  it("marks a default form dirty, clears it when restored, and resets to baseline", async () => {
    const container = await renderDocument(resetDocument());
    const input = container.querySelector<HTMLInputElement>("#field-name");
    const form = container.querySelector<HTMLFormElement>("form");
    expect(input?.value).toBe("Alice");
    expect(form?.dataset.formDirty).toBe("false");
    expect(hasDirtyState()).toBe(false);

    await act(async () => setInputValue(input!, "Bob"));
    expect(form?.dataset.formDirty).toBe("true");
    expect(hasDirtyState()).toBe(true);

    await act(async () => {
      setInputValue(input!, "Alice");
    });
    expect(form?.dataset.formDirty).toBe("false");
    expect(hasDirtyState()).toBe(false);

    await act(async () => {
      setInputValue(input!, "Carol");
      form?.dispatchEvent(new Event("reset", { bubbles: true, cancelable: true }));
    });
    expect(input?.value).toBe("Alice");
    expect(form?.dataset.formDirty).toBe("false");
    expect(hasDirtyState()).toBe(false);
  });

  it("does not register search form query changes as business dirty state", async () => {
    const container = await renderDocument(searchDocument());
    const input = container.querySelector<HTMLInputElement>("#field-q");
    expect(hasDirtyState()).toBe(false);
    await act(async () => setInputValue(input!, "alice"));
    expect(hasDirtyState()).toBe(false);
    expect(container.querySelector<HTMLFormElement>("form")?.dataset.formDirty).toBe("false");
  });

  it("clears dirty after a successful submit and preserves it after a failed submit", async () => {
    let response: Response = new Response("{}", { status: 200 });
    const fetcher: typeof fetch = vi.fn(async () => response) as typeof fetch;
    const container = await renderDocument(requestDocument(), fetcher);
    const input = container.querySelector<HTMLInputElement>("#field-name");
    const form = container.querySelector<HTMLFormElement>("form");
    const submit = Array.from(container.querySelectorAll<HTMLButtonElement>("button")).find(
      (button) => button.textContent?.trim() === "Save",
    );

    await act(async () => setInputValue(input!, "Bob"));
    expect(hasDirtyState()).toBe(true);
    await act(async () => submit?.click());
    expect(hasDirtyState()).toBe(false);
    expect(form?.dataset.formDirty).toBe("false");

    response = new Response(JSON.stringify({ error: "SAVE_FAILED", message: "try again" }), {
      status: 400,
      headers: { "Content-Type": "application/json" },
    });
    await act(async () => setInputValue(input!, "Carol"));
    await act(async () => submit?.click());
    expect(hasDirtyState()).toBe(true);
    expect(form?.dataset.formDirty).toBe("true");
    expect(container.querySelector('[role="alert"]')?.textContent).toContain("SAVE_FAILED");
  });

  it("keeps dirty after client validation failure and a transport throw", async () => {
    const fetcher = vi.fn(async () => {
      throw new Error("offline");
    }) as typeof fetch;
    const container = await renderDocument(validationDocument(), fetcher);
    const input = container.querySelector<HTMLInputElement>("#field-name");
    const submit = Array.from(container.querySelectorAll<HTMLButtonElement>("button")).find(
      (button) => button.textContent?.trim() === "Save",
    );

    await act(async () => setInputValue(input!, ""));
    expect(hasDirtyState()).toBe(true);
    await act(async () => submit?.click());
    expect(fetcher).not.toHaveBeenCalled();
    expect(hasDirtyState()).toBe(true);

    await act(async () => setInputValue(input!, "Bob"));
    await act(async () => submit?.click());
    expect(fetcher).toHaveBeenCalledTimes(1);
    expect(hasDirtyState()).toBe(true);
    expect(container.querySelector('[role="alert"]')?.textContent).toContain("REQUEST_FAILED");
  });

  it("keeps a dirty modal open when canceling and closes only after confirmation", async () => {
    const container = await renderDocument(modalDocument());
    await act(async () => {
      container.querySelector<HTMLButtonElement>("button")?.click();
    });
    const input = container.querySelector<HTMLInputElement>("#field-name");
    await act(async () => setInputValue(input!, "Bob"));
    expect(container.querySelector('[role="dialog"]')).not.toBeNull();
    expect(hasDirtyState()).toBe(true);

    const confirm = vi.spyOn(window, "confirm").mockReturnValue(false);
    await act(async () => {
      container.querySelector<HTMLButtonElement>('button[aria-label="Close dialog"]')?.click();
    });
    expect(confirm).toHaveBeenCalledTimes(1);
    expect(container.querySelector('[role="dialog"]')).not.toBeNull();
    expect(hasDirtyState()).toBe(true);

    confirm.mockReturnValue(true);
    await act(async () => {
      container.querySelector<HTMLButtonElement>('button[aria-label="Close dialog"]')?.click();
    });
    expect(container.querySelector('[role="dialog"]')).toBeNull();
    expect(hasDirtyState()).toBe(false);
  });

  it("clears the modal dirty source when its submit succeeds", async () => {
    const fetcher: typeof fetch = vi.fn(async () => new Response("{}", { status: 200 }));
    const container = await renderDocument(modalSubmitDocument(), fetcher);
    await act(async () => {
      container.querySelector<HTMLButtonElement>("button")?.click();
    });
    const input = container.querySelector<HTMLInputElement>("#field-name");
    const submit = Array.from(container.querySelectorAll<HTMLButtonElement>("button")).find(
      (button) => button.textContent?.trim() === "Save",
    );
    await act(async () => setInputValue(input!, "Bob"));
    expect(hasDirtyState()).toBe(true);
    await act(async () => submit?.click());
    expect(container.querySelector('[role="dialog"]')).toBeNull();
    expect(hasDirtyState()).toBe(false);
  });
});
