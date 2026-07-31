// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { RenderPage } from "@/renderer/render.tsx";
import type { RenderPageDocument } from "@/renderer/render";

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

async function renderDocument(pageDoc: RenderPageDocument, context: Record<string, unknown>) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(<RenderPage document={pageDoc} context={context} />);
  });
  return container;
}

function reactionFormDocument(reactions: unknown[]): RenderPageDocument {
  return {
    meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
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
      body: { type: "chart", id: "x", props: {} } as unknown as RenderPageDocument["body"],
    };
    const container = await renderDocument(pageDoc, {});
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "outside the §5 renderer whitelist",
    );
  });
});
