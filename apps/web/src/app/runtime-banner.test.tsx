// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { RuntimeBanner } from "@/app/runtime-banner";
import { I18nProvider } from "@/i18n/runtime";

const active: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", { configurable: true, value: true });
});

afterEach(async () => {
  for (const { root, container } of active.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

function renderBanner(mode: Parameters<typeof RuntimeBanner>[0]["runtimeMode"]): HTMLDivElement {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  active.push({ root, container });
  act(() => {
    root.render(
      <I18nProvider stored="en-US">
        <RuntimeBanner runtimeMode={mode} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("RuntimeBanner", () => {
  it("renders nothing for normal or missing mode", () => {
    expect(renderBanner("normal").querySelector("[data-runtime-banner]")).toBeNull();
    expect(renderBanner(undefined).querySelector("[data-runtime-banner]")).toBeNull();
  });

  it("splits copy by exact runtimeMode and never uses Host availability names", () => {
    const maintenance = renderBanner("maintenance").querySelector("[data-runtime-banner='maintenance']");
    const degraded = renderBanner("degraded").querySelector("[data-runtime-banner='degraded']");
    const readOnly = renderBanner("read-only").querySelector("[data-runtime-banner='read-only']");
    expect(maintenance?.textContent).toMatch(/maintenance/i);
    expect(degraded?.textContent).toMatch(/degraded/i);
    expect(readOnly?.textContent).toMatch(/read-only/i);
    expect(maintenance?.textContent).not.toBe(degraded?.textContent);
    expect(degraded?.textContent).not.toBe(readOnly?.textContent);
  });
});
