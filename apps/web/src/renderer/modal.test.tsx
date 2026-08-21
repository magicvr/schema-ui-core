// @vitest-environment jsdom
//
// F-002 跨模块 UI 可访问性下限 · 可复跑焦点断言（S0 D-003 §8）：
// 模态焦点进入/约束（Tab 循环）/恢复与 Escape 关闭成立。
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { ModalHost } from "@/renderer/modal";

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

async function renderModal(onClose: () => void = () => {}) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <ModalHost title="Edit" onClose={onClose}>
        <input aria-label="Name" />
        <input aria-label="Email" />
        <button type="button">Save</button>
      </ModalHost>,
    );
  });
  return container;
}

function pressKey(key: string) {
  document.dispatchEvent(new KeyboardEvent("keydown", { key, bubbles: true }));
}

describe("ModalHost shared-host accessibility floor (F-002)", () => {
  it("moves focus into the dialog on open and restores it on close", async () => {
    const trigger = document.createElement("button");
    trigger.textContent = "Open";
    document.body.appendChild(trigger);
    trigger.focus();

    await renderModal();
    // Focus enters the dialog (first focusable = the close button in header order
    // is inside the dialog; assert activeElement is within the dialog container).
    const dialog = document.querySelector('[role="dialog"]') as HTMLElement;
    expect(dialog).not.toBeNull();
    expect(dialog.contains(document.activeElement)).toBe(true);

    // Unmount → focus restored to the trigger.
    const { root, container } = activeRoots[0];
    await act(async () => root.unmount());
    container.remove();
    activeRoots.pop();
    expect(document.activeElement).toBe(trigger);
    trigger.remove();
  });

  it("traps Tab focus within the dialog", async () => {
    await renderModal();
    const dialog = document.querySelector('[role="dialog"]') as HTMLElement;
    const focusables = Array.from(
      dialog.querySelectorAll<HTMLElement>(
        'a[href], button:not([disabled]), input:not([disabled]), [tabindex]:not([tabindex="-1"])',
      ),
    );
    expect(focusables.length).toBeGreaterThanOrEqual(1);
    const first = focusables[0];
    const last = focusables[focusables.length - 1];

    // Focus on the last element; Tab wraps to the first (no escape).
    last.focus();
    await act(async () => {
      pressKey("Tab");
    });
    expect(document.activeElement).toBe(first);

    // Shift+Tab from the first wraps to the last.
    first.focus();
    await act(async () => {
      document.dispatchEvent(new KeyboardEvent("keydown", { key: "Tab", shiftKey: true, bubbles: true }));
    });
    expect(document.activeElement).toBe(last);
  });

  it("closes on Escape", async () => {
    const onClose = vi.fn();
    await renderModal(onClose);
    await act(async () => {
      pressKey("Escape");
    });
    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
