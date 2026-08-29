// @vitest-environment jsdom

/**
 * W11 · U-01/U-02: dynamic option sources (optionsSource) for
 * select / radio / checkboxGroup. Options load from a same-origin path with
 * the auth-aware fetcher; static options remain the fallback while loading;
 * invalid sources and failed fetches fail closed to an empty option set.
 */

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { FormControls } from "@/renderer/form-controls.tsx";
import type { FormControlField } from "@/renderer/form-controls.types";

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

async function renderControls(
  fields: FormControlField[],
  values: Record<string, unknown>,
  fetcher?: typeof fetch,
) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider stored="en-US">
        <FormControls fields={fields} values={values} onChange={() => undefined} fetcher={fetcher} />
      </I18nProvider>,
    );
  });
  await act(async () => {
    await new Promise((resolve) => setTimeout(resolve, 0));
  });
  return container;
}

function jsonResponse(body: unknown): Response {
  return new Response(JSON.stringify(body), {
    status: 200,
    headers: { "Content-Type": "application/json" },
  });
}

describe("FormControls dynamic option sources (W11 · U-01/U-02)", () => {
  it("loads checkboxGroup options from {items:[...]} via optionsMapping", async () => {
    const fetcher = vi.fn().mockResolvedValue(
      jsonResponse({
        items: [
          { key: "admin", name: "Administrator" },
          { key: "editor", name: "Editor" },
          { key: "viewer", name: "Viewer" },
        ],
      }),
    ) as unknown as typeof fetch;
    const field: FormControlField = {
      id: "roles",
      label: "Roles",
      type: "checkboxGroup",
      optionsSource: {
        url: "/api/roles",
        params: { pageSize: 100 },
        valueField: "key",
        labelField: "name",
      },
    };
    const container = await renderControls([field], {}, fetcher);
    expect(fetcher).toHaveBeenCalledWith(expect.stringContaining("pageSize=100"), expect.anything());
    expect(container.textContent).toContain("Administrator");
    expect(container.textContent).toContain("Viewer");
  });

  it("falls back to static options while the source is loading", async () => {
    let resolveFetch: (r: Response) => void = () => undefined;
    const fetcher = vi.fn().mockReturnValue(
      new Promise<Response>((resolve) => {
        resolveFetch = resolve;
      }),
    ) as unknown as typeof fetch;
    const field: FormControlField = {
      id: "role",
      label: "Role",
      type: "select",
      options: [{ value: "fallback", label: "Fallback" }],
      optionsSource: { url: "/api/roles", valueField: "key", labelField: "name" },
    };
    const container = await renderControls([field], {}, fetcher);
    expect(container.textContent).toContain("Fallback");
    // Resolve the source → dynamic options replace the fallback.
    await act(async () => {
      resolveFetch(jsonResponse({ items: [{ key: "editor", name: "Editor" }] }));
    });
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(container.textContent).toContain("Editor");
    expect(container.textContent).not.toContain("Fallback");
  });

  it("fails closed to an empty option set on a failed fetch", async () => {
    const fetcher = vi.fn().mockRejectedValue(new Error("network")) as unknown as typeof fetch;
    const field: FormControlField = {
      id: "permissions",
      label: "Permissions",
      type: "checkboxGroup",
      optionsSource: { url: "/api/permissions", valueField: "key", labelField: "key" },
    };
    const container = await renderControls([field], {}, fetcher);
    // No options rendered (no checkbox labels beyond the fieldset legend).
    expect(container.textContent).toBe("Permissions");
  });

  it("rejects invalid (non single-slash) option sources without fetching", async () => {
    const fetcher = vi.fn() as unknown as typeof fetch;
    const field: FormControlField = {
      id: "role",
      label: "Role",
      type: "select",
      optionsSource: { url: "https://evil.example/api/roles", valueField: "key", labelField: "name" },
    };
    const container = await renderControls([field], {}, fetcher);
    expect(fetcher).not.toHaveBeenCalled();
    expect(container.textContent).toContain("Role");
  });

  // W10 F-007: WHATWG URL parsing normalizes "\" to "/" in special schemes,
  // so "/\host" would otherwise become the protocol-relative "//host" and
  // escape the origin. The validator must reject backslashes without fetching.
  it("rejects backslash option sources that would normalize to protocol-relative urls", async () => {
    const fetcher = vi.fn() as unknown as typeof fetch;
    for (const url of ["/\\evil.example/api/roles", "/api/roles\\", "/api\\evil"]) {
      const field: FormControlField = {
        id: "role",
        label: "Role",
        type: "select",
        optionsSource: { url, valueField: "key", labelField: "name" },
      };
      const container = await renderControls([field], {}, fetcher);
      expect(fetcher).not.toHaveBeenCalled();
      expect(container.textContent).toContain("Role");
    }
  });
});