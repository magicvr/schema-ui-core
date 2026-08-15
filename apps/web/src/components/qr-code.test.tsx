// @vitest-environment jsdom
// W11 · M-01: QrCode renders a scannable SVG matrix for an otpauth URI.
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { QrCode } from "@/components/qr-code";

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

async function renderQr(props: Parameters<typeof QrCode>[0]) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(<QrCode {...props} />);
  });
  return container;
}

describe("QrCode", () => {
  it("renders a square SVG matrix for a valid otpauth URI", async () => {
    const container = await renderQr({
      value: "otpauth://totp/Schema%20UI%20Core:admin?secret=JBSWY3DPEHPK3PXP&issuer=Schema%20UI%20Core",
      size: 160,
    });
    const svg = container.querySelector<SVGSVGElement>("[data-qr-code]");
    expect(svg).not.toBeNull();
    expect(svg!.getAttribute("width")).toBe("160");
    expect(svg!.getAttribute("height")).toBe("160");
    expect(svg!.getAttribute("role")).toBe("img");
    // The matrix is square and contains dark modules (black rects).
    const viewBox = svg!.getAttribute("viewBox") ?? "";
    const [w, h] = viewBox.split(" ").slice(2).map(Number);
    expect(w).toBeGreaterThan(0);
    expect(w).toBe(h);
    const dark = svg!.querySelectorAll('rect[fill="black"]');
    expect(dark.length).toBeGreaterThan(0);
  });

  it("renders nothing for an empty value", async () => {
    const container = await renderQr({ value: "" });
    expect(container.querySelector("[data-qr-code]")).toBeNull();
  });

  it("produces a deterministic matrix for the same input", async () => {
    const uri =
      "otpauth://totp/Schema%20UI%20Core:admin?secret=JBSWY3DPEHPK3PXP&issuer=Schema%20UI%20Core";
    const a = await renderQr({ value: uri });
    const b = await renderQr({ value: uri });
    const darkA = a.querySelectorAll('rect[fill="black"]').length;
    const darkB = b.querySelectorAll('rect[fill="black"]').length;
    expect(darkA).toBe(darkB);
    expect(darkA).toBeGreaterThan(0);
  });
});
