// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";
import type { RenderPageDocument } from "@/renderer/render";

/**
 * F-02 local extension (GOAL-004 D-002 §5): the protocol's CustomAction
 * extension point — a whitelisted "export.*" handler fetches the CSV with the
 * authed transport and triggers a browser download. Unknown handler names fail
 * closed (CUSTOM_HANDLER_NOT_FOUND).
 */

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

// jsdom lacks URL.createObjectURL/revokeObjectURL — install local stubs.
const originalCreateObjectURL = URL.createObjectURL;
const originalRevokeObjectURL = URL.revokeObjectURL;
const objectUrls: string[] = [];
const revokes: string[] = [];
URL.createObjectURL = (obj: Blob | MediaSource) => {
  objectUrls.push(obj instanceof Blob ? obj.type : "media");
  return "blob:mock-url";
};
URL.revokeObjectURL = (url: string) => {
  revokes.push(url);
};

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  objectUrls.length = 0;
  revokes.length = 0;
  // afterEach restores the originals, so re-install the jsdom stubs for every
  // test (jsdom does not implement createObjectURL/revokeObjectURL).
  URL.createObjectURL = (obj: Blob | MediaSource) => {
    objectUrls.push(obj instanceof Blob ? obj.type : "media");
    return "blob:mock-url";
  };
  URL.revokeObjectURL = (url: string) => {
    revokes.push(url);
  };
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  vi.restoreAllMocks();
  URL.createObjectURL = originalCreateObjectURL;
  URL.revokeObjectURL = originalRevokeObjectURL;
});

async function renderDocument(pageDoc: RenderPageDocument, dataFetcher?: typeof fetch) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <RenderPage
          document={pageDoc}
          context={{}}
          dataFetcher={dataFetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={dataFetcher ?? fetch} />}
        />
      </I18nProvider>,
    );
  });
  return container;
}

function downloadDocument(): RenderPageDocument {
  // The export action uses the protocol's CustomAction shape (action.schema.json).
  return {
    meta: {
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation", "actions.page.trigger"],
    },
    actions: {
      exportUsers: {
        type: "custom",
        handler: "export.users",
      },
      unknownHandler: {
        type: "custom",
        handler: "export.nope",
      },
      plainAction: {
        type: "request",
        method: "POST",
        url: "/api/export/users",
        onSuccess: { behavior: "reload" },
      },
    },
    body: {
      type: "section",
      id: "root",
      children: [
        {
          type: "actionButton",
          id: "btn-export",
          props: { label: "Export", actionId: "exportUsers" },
        },
        {
          type: "actionButton",
          id: "btn-unknown",
          props: { label: "Unknown", actionId: "unknownHandler" },
        },
        {
          type: "actionButton",
          id: "btn-plain",
          props: { label: "Plain", actionId: "plainAction" },
        },
      ],
    },
  } as unknown as RenderPageDocument;
}

function findButton(container: HTMLElement, text: string): HTMLButtonElement {
  const buttons = container.querySelectorAll("button");
  const found = Array.from(buttons).find((b) => b.textContent === text);
  if (!(found instanceof HTMLButtonElement)) {
    throw new Error(`button "${text}" not found`);
  }
  return found;
}

it("custom export handler fetches the CSV and triggers a download", async () => {
  const clickSpy = vi.fn();
  const originalCreateElement = document.createElement.bind(document);
  vi.spyOn(document, "createElement").mockImplementation((tag: string) => {
    const el = originalCreateElement(tag) as HTMLAnchorElement;
    if (tag === "a") {
      el.click = clickSpy as () => void;
    }
    return el;
  });
  const fetcher = vi.fn(async () => {
    return new Response("id,username\n1,admin", {
      status: 200,
      headers: { "Content-Type": "text/csv" },
    });
  });
  const container = await renderDocument(downloadDocument(), fetcher as typeof fetch);
  await act(async () => {
    findButton(container, "Export").click();
    await new Promise((resolve) => setTimeout(resolve, 50));
  });
  expect(fetcher).toHaveBeenCalledWith(
    "/api/export/users",
    expect.objectContaining({ method: "GET" }),
  );
  expect(clickSpy).toHaveBeenCalled();
  expect(objectUrls.length).toBeGreaterThan(0);
  expect(revokes.length).toBeGreaterThan(0);
});

it("unknown custom handler names fail closed", async () => {
  const clickSpy = vi.fn();
  const originalCreateElement = document.createElement.bind(document);
  vi.spyOn(document, "createElement").mockImplementation((tag: string) => {
    const el = originalCreateElement(tag) as HTMLAnchorElement;
    if (tag === "a") {
      el.click = clickSpy as () => void;
    }
    return el;
  });
  const fetcher = vi.fn(async () => new Response("ok", { status: 200 }));
  const container = await renderDocument(downloadDocument(), fetcher as typeof fetch);
  await act(async () => {
    findButton(container, "Unknown").click();
    await new Promise((resolve) => setTimeout(resolve, 50));
  });
  expect(fetcher).not.toHaveBeenCalled();
  expect(clickSpy).not.toHaveBeenCalled();
});

it("non-custom actions are unaffected", async () => {
  const clickSpy = vi.fn();
  const originalCreateElement = document.createElement.bind(document);
  vi.spyOn(document, "createElement").mockImplementation((tag: string) => {
    const el = originalCreateElement(tag) as HTMLAnchorElement;
    if (tag === "a") {
      el.click = clickSpy as () => void;
    }
    return el;
  });
  const fetcher = vi.fn(async () => new Response("ok", { status: 200 }));
  const container = await renderDocument(downloadDocument(), fetcher as typeof fetch);
  await act(async () => {
    findButton(container, "Plain").click();
    await new Promise((resolve) => setTimeout(resolve, 50));
  });
  expect(fetcher).toHaveBeenCalled();
  expect(clickSpy).not.toHaveBeenCalled();
});

// S-02 (GOAL-007 D-002 §5): the library.download custom handler resolves the
// {id} slot from the row context and names the blob after the row's stored
// name (client-side only). A templated handler without a row fails closed.
it("library.download resolves the row id and names the download from the row", async () => {
  const clickSpy = vi.fn();
  const downloadNames: string[] = [];
  const originalCreateElement = document.createElement.bind(document);
  vi.spyOn(document, "createElement").mockImplementation((tag: string) => {
    const el = originalCreateElement(tag) as HTMLAnchorElement;
    if (tag === "a") {
      const realClick = el.click.bind(el);
      el.click = () => {
        downloadNames.push(el.download);
        clickSpy();
      };
    }
    return el;
  });
  const fetchCalls: Array<{ url: string; init?: RequestInit }> = [];
  const fetcher = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = String(input);
    fetchCalls.push({ url, init });
    if (url.includes("/download")) {
      return new Response("file-bytes", { status: 200 });
    }
    // Table list fetch: one row.
    return new Response(
      JSON.stringify({ items: [{ id: "abc123", name: "report (Q1).pdf", type: "text/plain", size: 10, owner: "user-admin", created: 1 }], total: 1, page: 1, pageSize: 10 }),
      { status: 200, headers: { "Content-Type": "application/json" } },
    );
  });
  const doc = {
    meta: {
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation", "permissions.inheritance", "table.sort"],
    },
    actions: {
      downloadFile: { type: "custom", handler: "library.download" },
    },
    body: {
      type: "table",
      id: "files-table",
      props: {
        columns: [{ field: "name", label: "Name" }],
        dataSource: "/api/library/files",
        actions: [
          { key: "download", label: "Download", actionRef: "downloadFile", permissionIntent: "edit" },
        ],
      },
    },
  } as unknown as RenderPageDocument;
  const container = await renderDocument(doc, fetcher as typeof fetch);
  // The table list loads asynchronously; wait for the row action to appear.
  let button: HTMLButtonElement | undefined;
  for (let attempt = 0; attempt < 20 && button === undefined; attempt++) {
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 25));
    });
    try {
      button = findButton(container, "Download");
    } catch {
      button = undefined;
    }
  }
  if (button === undefined) throw new Error("Download row action not rendered");
  await act(async () => {
    button!.click();
    await new Promise((resolve) => setTimeout(resolve, 200));
  });
  const downloadCall = fetchCalls.find((call) => call.url.includes("/download"));
  expect(downloadCall).toBeTruthy();
  expect(downloadCall!.url).toBe("/api/library/files/abc123/download");
  expect(clickSpy).toHaveBeenCalled();
  // The anchor's download attribute carries the row name (quote-safe).
  expect(downloadNames).toContain("report (Q1).pdf");
  expect(objectUrls.length).toBeGreaterThan(0);
});
