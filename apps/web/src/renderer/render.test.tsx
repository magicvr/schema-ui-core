// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { RenderPage } from "@/renderer/render.tsx";
import type { RenderMeta, RenderPageDocument } from "@/renderer/render";

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

async function renderDocument(
  pageDoc: RenderPageDocument,
  context: Record<string, unknown>,
  dataFetcher?: typeof fetch,
) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(<RenderPage document={pageDoc} context={context} dataFetcher={dataFetcher} />);
  });
  return container;
}

/** Fixture transport returning the pinned list envelope for any dataSource. */
function fixtureListFetcher(items: Array<Record<string, unknown>>): typeof fetch {
  return (async (input: RequestInfo | URL) => {
    const url = String(input);
    if (!url.startsWith("/api/")) {
      return new Response(JSON.stringify({ error: "NOT_FOUND", message: "not found" }), {
        status: 404,
        headers: { "Content-Type": "application/json" },
      });
    }
    return new Response(
      JSON.stringify({ items, total: items.length, page: 1, pageSize: 100 }),
      { status: 200, headers: { "Content-Type": "application/json" } },
    );
  }) as typeof fetch;
}

function displayDocument(body: RenderPageDocument["body"]): RenderPageDocument {
  return {
    meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
    body,
  };
}

describe("RenderPage full $deps reactions (I-PROTO-FULL-001 · D-EXPR integration)", () => {
  it("commits fulfill values and hides a field through the multi-round engine", async () => {
    const pageDoc: RenderPageDocument = {
      meta: {
        protocolVersion: "2.7",
        requiredCapabilities: ["app.manifest", "form.controls.extended", "form.controls.advanced"],
      },
      body: {
        type: "form",
        id: "reactive-form",
        props: {
          fields: [
            { id: "trigger", label: "Trigger", type: "switch" },
            {
              id: "status",
              label: "Status",
              type: "input",
              defaultValue: "baseline",
              reactions: [{ when: "$deps.trigger == true", fulfill: { value: "closed", disabled: true } }],
            },
          ],
        },
      },
    };
    const container = await renderDocument(pageDoc, {});
    // The defaultValue initializes status; the engine then commits "closed"
    // because trigger defaults to false… — flip the switch to drive a commit.
    const inputs = container.querySelectorAll("input");
    const switchInput = inputs[0]!;
    await act(async () => {
      switchInput.click();
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    const statusInput = container.querySelector("input[type='text']") as HTMLInputElement;
    expect(statusInput.value).toBe("closed");
    expect(statusInput.disabled).toBe(true);
  });

  it("reports a loop limit as a blocking reaction error", async () => {
    const pageDoc: RenderPageDocument = {
      meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
      body: {
        type: "form",
        id: "loop-form",
        props: {
          fields: [
            {
              id: "toggle",
              label: "Toggle",
              type: "input",
              defaultValue: "a",
              reactions: [
                {
                  when: "$deps.toggle == 'a'",
                  fulfill: { value: "b" },
                  otherwise: { value: "a" },
                },
              ],
            },
          ],
        },
      },
    };
    const container = await renderDocument(pageDoc, {});
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(container.textContent).toContain("REACTION_LOOP_LIMIT");
  });
});

describe("RenderPage display types (I-PROTO-FULL-001 · statCard/chart)", () => {
  it("renders a statCard value with a currency format from its dataSource", async () => {
    const pageDoc = displayDocument({
      type: "statCard",
      id: "revenue",
      props: { label: "Revenue", format: "currency", valueField: "amount", dataSource: "/api/orders" },
    });
    const container = await renderDocument(
      pageDoc,
      {},
      fixtureListFetcher([{ id: "o1", amount: 1250, month: "2026-08", count: 12 }]),
    );
    await act(async () => {
      // Let the async data fetch settle.
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(container.textContent).toContain("Revenue");
    expect(container.textContent).toContain("1250");
  });

  it("fails closed when statCard format rejects the value type", async () => {
    const pageDoc = displayDocument({
      type: "statCard",
      id: "revenue",
      props: { label: "Revenue", format: "currency", valueField: "name", dataSource: "/api/orders" },
    });
    const container = await renderDocument(
      pageDoc,
      {},
      fixtureListFetcher([{ id: "o1", name: "order-a", amount: 1250, month: "2026-08", count: 12 }]),
    );
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "statCard format",
    );
  });

  it("renders a chart from its dataSource and lists its points", async () => {
    const pageDoc = displayDocument({
      type: "chart",
      id: "orders-chart",
      props: { chartType: "bar", xField: "month", yField: "count", dataSource: "/api/orders" },
    });
    const container = await renderDocument(
      pageDoc,
      {},
      fixtureListFetcher([
        { id: "o1", amount: 1250, month: "2026-08", count: 12 },
        { id: "o2", amount: 800, month: "2026-09", count: 7 },
      ]),
    );
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(container.querySelector("svg[role='img']")).not.toBeNull();
    expect(container.textContent).toContain("2026-08");
    expect(container.textContent).toContain("12");
  });

  it("fails closed when a chart node lacks its required props", async () => {
    const pageDoc = displayDocument({
      type: "chart",
      id: "bad",
      props: { chartType: "bar", dataSource: "/api/orders" },
    });
    const container = await renderDocument(pageDoc, {});
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "chart node requires chartType",
    );
  });
});

function reactionFormDocument(reactions: unknown[]): RenderPageDocument {
  return {
    meta: {
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "form.controls.extended"],
    },
    body: {
      type: "form",
      id: "reactive-form",
      props: {
        fields: [
          { id: "name", label: "Name", type: "input" },
          { id: "notes", label: "Notes", type: "textarea" },
          { id: "approval", label: "Approval", type: "switch" },
        ],
        reactions,
      },
    },
  };
}

// GOAL-009 S1 (A-002 F-002-002): a default-mode form with a real submitAction,
// so gate/reaction errors must block the outgoing request, not just be shown.
function submitFormDocument(
  fields: unknown,
  reactions: unknown[],
  meta: RenderMeta = {
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "form.controls.extended"],
  },
): RenderPageDocument {
  return {
    meta,
    actions: {
      submit: { type: "request", method: "POST", url: "/api/users" },
    },
    body: {
      type: "form",
      id: "gated-submit-form",
      props: {
        fields,
        reactions,
        submitAction: "submit",
        submitLabel: "Submit record",
      },
    },
  } as RenderPageDocument;
}

function submitButton(container: HTMLElement): HTMLButtonElement {
  return Array.from(container.querySelectorAll("button")).find(
    (button) => button.textContent?.trim() === "Submit record",
  ) as HTMLButtonElement;
}

async function withFetchSpy(run: (fetchSpy: ReturnType<typeof vi.fn>) => Promise<void>) {
  const fetchSpy = vi.fn(async () => new Response("{}", { status: 200 }));
  vi.stubGlobal("fetch", fetchSpy);
  try {
    await run(fetchSpy);
  } finally {
    vi.unstubAllGlobals();
  }
}

describe("RenderPage form node with reactions", () => {
  it("renders every field by default", async () => {
    const container = await renderDocument(reactionFormDocument([]), {});
    expect(container.textContent).toContain("Name");
    expect(container.textContent).toContain("Notes");
    expect(container.textContent).toContain("Approval");
  });

  it("hides a field when the reaction holds", async () => {
    const pageDoc = reactionFormDocument([
      {
        id: "hide-notes",
        when: "$context.features.audit == true",
        apply: [{ fieldId: "notes", visible: false }],
      },
    ]);
    const container = await renderDocument(pageDoc, { features: { audit: true } });
    expect(container.textContent).toContain("Name");
    expect(container.textContent).not.toContain("Notes");
  });

  it("disables a field when the reaction holds", async () => {
    const pageDoc = reactionFormDocument([
      {
        id: "lock-approval",
        when: '$context.user.roles contains "viewer"',
        apply: [{ fieldId: "approval", disabled: true }],
      },
    ]);
    const container = await renderDocument(pageDoc, { user: { roles: ["viewer"] } });
    const checkbox = container.querySelector('input[type="checkbox"]');
    expect(checkbox).not.toBeNull();
    expect((checkbox as HTMLInputElement).disabled).toBe(true);
  });

  it("fails closed on an unknown node type", async () => {
    const pageDoc: RenderPageDocument = {
      meta: { protocolVersion: "2.7", requiredCapabilities: [] },
      body: { type: "slider", id: "x", props: {} } as unknown as RenderPageDocument["body"],    };
    const container = await renderDocument(pageDoc, {});
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "outside the registry renderer whitelist",
    );
  });

  it("rejects a field whose capability gate fails and reports the gate error", async () => {
    // textarea requires 2.6 + form.controls.extended; the meta only declares app.manifest.
    const pageDoc: RenderPageDocument = {
      meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
      body: {
        type: "form",
        id: "gated-form",
        props: {
          fields: [
            { id: "name", label: "Name", type: "input" },
            { id: "notes", label: "Notes", type: "textarea" },
          ],
        },
      },
    };
    const container = await renderDocument(pageDoc, {});
    // The compliant input still renders; the gated textarea is rejected.
    expect(container.textContent).toContain("Name");
    expect(container.textContent).not.toContain("Notes");
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "FORM_CAPABILITY_REQUIRED",
    );
  });

  it("rejects an unknown field type and reports it", async () => {
    const pageDoc: RenderPageDocument = {
      meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
      body: {
        type: "form",
        id: "typed-form",
        props: {
          fields: [{ id: "up", label: "Upload", type: "slider" as "upload" }],
        },
      },
    };
    const container = await renderDocument(pageDoc, {});
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "FORM_TYPE_NOT_WHITELISTED",
    );
  });

  it("gates an upload field missing action/actionRef (registry oneOf)", async () => {
    const pageDoc: RenderPageDocument = {
      meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest", "actions.upload"] },
      body: {
        type: "form",
        id: "upload-form",
        props: {
          fields: [{ id: "file", label: "File", type: "upload" }],
        },
      },
    };
    const container = await renderDocument(pageDoc, {});
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "UPLOAD_ACTION_REQUIRED",
    );
  });

  it("dispatches the remaining frozen §5 node types", async () => {
    const pageDoc: RenderPageDocument = {
      meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
      body: {
        type: "grid",
        id: "g1",
        props: { columns: 2 },
        children: [
          { type: "text", props: { text: "hello text" } },
          {
            type: "recordView",
            props: { record: { id: "rec-1", status: "active" } },
          },
          {
            type: "actionButton",
            props: { label: "Approve", actionId: "approve", visibleWhen: true },
          },
          {
            type: "tabs",
            props: {},
            children: [
              { type: "text", props: { text: "tab one" } },
              { type: "text", props: { text: "tab two" } },
            ],
          },
        ],
      },
    };
    const container = await renderDocument(pageDoc, {});
    expect(container.textContent).toContain("hello text");
    expect(container.textContent).toContain("rec-1");
    expect(container.textContent).toContain("Approve");
    expect(container.textContent).toContain("tab one");
  });
});

describe("RenderPage form submit gate (GOAL-009 S1 · F-002-002)", () => {
  it("disables submit and sends no request while a field gate error is shown", async () => {
    await withFetchSpy(async (fetchSpy) => {
      // textarea requires 2.6 + form.controls.extended; the meta declares neither.
      const pageDoc = submitFormDocument(
        [{ id: "notes", label: "Notes", type: "textarea" }],
        [],
        { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
      );
      const container = await renderDocument(pageDoc, {});
      expect(container.querySelector('[role="alert"]')?.textContent).toContain(
        "FORM_CAPABILITY_REQUIRED",
      );
      const button = submitButton(container);
      expect(button.disabled).toBe(true);
      await act(async () => button.click());
      expect(fetchSpy).not.toHaveBeenCalled();
    });
  });

  it("disables submit and sends no request while a reaction error is shown", async () => {
    await withFetchSpy(async (fetchSpy) => {
      // $deps.* is outside the frozen reaction grammar �?REACTION_EXPRESSION_INVALID.
      const pageDoc = submitFormDocument(
        [{ id: "name", label: "Name", type: "input" }],
        [
          {
            id: "bad",
            when: "$deps.admin == true",
            apply: [{ fieldId: "name", visible: true }],
          },
        ],
      );
      const container = await renderDocument(pageDoc, {});
      expect(container.querySelector('[role="alert"]')?.textContent).toContain(
        "REACTION_EXPRESSION_INVALID",
      );
      const button = submitButton(container);
      expect(button.disabled).toBe(true);
      await act(async () => button.click());
      expect(fetchSpy).not.toHaveBeenCalled();
    });
  });

  it("still submits when the form is valid (positive control)", async () => {
    await withFetchSpy(async (fetchSpy) => {
      const pageDoc = submitFormDocument(
        [{ id: "name", label: "Name", type: "input" }],
        [],
      );
      const container = await renderDocument(pageDoc, {});
      const button = submitButton(container);
      expect(button.disabled).toBe(false);
      await act(async () => button.click());
      expect(fetchSpy).toHaveBeenCalledTimes(1);
      expect(fetchSpy).toHaveBeenCalledWith(
        "/api/users",
        expect.objectContaining({ method: "POST" }),
      );
    });
  });
});

