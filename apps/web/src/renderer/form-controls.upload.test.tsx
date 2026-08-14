// @vitest-environment jsdom

/**
 * W9 (GOAL-010): UploadField per-field remove button.
 * Single uploads with a committed value render a localized remove button
 * (I-008: the settings brand fields need a per-field clear affordance);
 * multiple uploads / readOnly / disabled fields never show it.
 */

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { FormControls } from "@/renderer/form-controls.tsx";
import type { FormControlField } from "@/renderer/form-controls";

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

async function renderControls(fields: FormControlField[], values: Record<string, unknown>, onChange: (id: string, value: unknown) => void, disabled = false) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider stored="zh-CN">
        <FormControls fields={fields} values={values} onChange={onChange} disabled={disabled} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("S3 upload field remove button (W9)", () => {
  it("renders a remove button for a committed single upload and clears the value", async () => {
    let cleared = false;
    const field: FormControlField = {
      id: "logoUrl",
      label: "Logo",
      type: "upload",
      actionRef: "uploadBrandingLogo",
    };
    const container = await renderControls(
      [field],
      { logoUrl: "/api/branding/assets/abcdef" },
      (id, value) => {
        if (id === "logoUrl" && value === "") {
          cleared = true;
        }
      },
    );
    expect(container.textContent).toContain("Value: /api/branding/assets/abcdef");
    const remove = [...container.querySelectorAll("button")].find((b) => b.textContent?.includes("移除图片"));
    expect(remove).not.toBeUndefined();
    await act(async () => {
      remove!.click();
    });
    expect(cleared).toBe(true);
  });

  it("does not render a remove button when the value is empty", async () => {
    const field: FormControlField = { id: "logoUrl", label: "Logo", type: "upload", actionRef: "uploadBrandingLogo" };
    const container = await renderControls([field], { logoUrl: "" }, () => undefined);
    expect([...container.querySelectorAll("button")].some((b) => b.textContent?.includes("移除图片"))).toBe(false);
  });

  it("does not render a remove button for multiple uploads", async () => {
    const field: FormControlField = { id: "files", label: "Files", type: "upload", actionRef: "up", multiple: true };
    const container = await renderControls([field], { files: ["/a", "/b"] }, () => undefined);
    expect([...container.querySelectorAll("button")].some((b) => b.textContent?.includes("移除图片"))).toBe(false);
  });

  it("does not render a remove button when the field is readOnly or the form is disabled", async () => {
    const readOnlyField: FormControlField = { id: "logoUrl", label: "Logo", type: "upload", actionRef: "up", readOnly: true };
    const container = await renderControls([readOnlyField], { logoUrl: "/x" }, () => undefined);
    expect([...container.querySelectorAll("button")].some((b) => b.textContent?.includes("移除图片"))).toBe(false);
    const disabled = await renderControls(
      [{ id: "logoUrl", label: "Logo", type: "upload", actionRef: "up" }],
      { logoUrl: "/x" },
      () => undefined,
      true,
    );
    expect([...disabled.querySelectorAll("button")].some((b) => b.textContent?.includes("移除图片"))).toBe(false);
  });
});
