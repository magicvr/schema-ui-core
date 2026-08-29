// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

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
  vi.restoreAllMocks();
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

async function renderDocument(pageDoc: RenderPageDocument) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(<RenderPage document={pageDoc} context={{}} />);
  });
  return container;
}

function intentButtonDocument(metaProtocolVersion: string): RenderPageDocument {
  return {
    // The 2.2 variant intentionally trips PROTOCOL_VERSION_TOO_LOW (and
    // CAPABILITY_REQUIRED) — the invalid-structure fixture. The 2.3 variant is
    // the valid control: version floor met AND the permissions.inheritance
    // capability declared, so L2 is clean and the gate stays open.
    meta: {
      protocolVersion: metaProtocolVersion,
      requiredCapabilities:
        metaProtocolVersion === "2.3"
          ? ["app.manifest", "permissions.inheritance"]
          : ["app.manifest"],
    },
    body: {
      type: "actionButton",
      id: "btn-1",
      props: { label: "Go", permissionIntent: "edit", key: "go-btn" },
    },
  } as RenderPageDocument;
}

/**
 * W9 A-005 R-F-001 regression lock: the L2 permission validator runs in the
 * production render path — a malformed permission structure surfaces at load
 * (console.error) and fails CLOSED (registered targets denied), while a valid
 * page keeps its gate open.
 */
describe("permission L2 validation wired into render (W9 R-F-001)", () => {
  it("denies gated targets and surfaces console.error on an invalid structure", async () => {
    const errorSpy = vi.spyOn(console, "error").mockImplementation(() => {});
    // protocolVersion 2.2 < 2.3 with permission fields present →
    // PROTOCOL_VERSION_TOO_LOW.
    const container = await renderDocument(intentButtonDocument("2.2"));
    const button = container.querySelector("button") as HTMLButtonElement | null;
    expect(button).not.toBeNull();
    expect(button!.disabled).toBe(true);
    expect(
      errorSpy.mock.calls.some((call) =>
        String(call[0]).includes("permission L2 validation failed"),
      ),
    ).toBe(true);
  });

  it("keeps a valid structure's gate open with no validation noise", async () => {
    const errorSpy = vi.spyOn(console, "error").mockImplementation(() => {});
    const container = await renderDocument(intentButtonDocument("2.3"));
    const button = container.querySelector("button") as HTMLButtonElement | null;
    expect(button).not.toBeNull();
    expect(button!.disabled).toBe(false);
    expect(errorSpy).not.toHaveBeenCalled();
  });
});
