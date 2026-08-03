// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { AuthError } from "@/account/auth-client";
import { LoginPage } from "@/app/LoginPage";

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

async function renderLogin(onLogin: (u: string, p: string) => Promise<void>) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(<LoginPage onLogin={onLogin} />);
  });
  return container;
}

// React 19 tracks the input value via a property descriptor; assign through the
// native setter and dispatch `input` so onChange fires in the test DOM.
function fill(container: HTMLDivElement, selector: string, value: string) {
  const el = container.querySelector<HTMLInputElement>(selector);
  expect(el).not.toBeNull();
  const setter = Object.getOwnPropertyDescriptor(
    window.HTMLInputElement.prototype,
    "value",
  )?.set;
  act(() => {
    setter?.call(el, value);
    el!.dispatchEvent(new Event("input", { bubbles: true }));
  });
}

describe("LoginPage", () => {
  it("renders the sign-in form and submits credentials", async () => {
    const onLogin = vi.fn().mockResolvedValue(undefined);
    const container = await renderLogin(onLogin);

    expect(container.textContent).toContain("Sign in");
    fill(container, "#username", "admin");
    fill(container, "#password", "admin");

    const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
    expect(button).not.toBeNull();
    await act(async () => button!.click());
    expect(onLogin).toHaveBeenCalledWith("admin", "admin");
  });

  it("surfaces the server error when login fails", async () => {
    const onLogin = vi.fn().mockRejectedValue(new AuthError("INVALID_CREDENTIALS", "invalid username or password"));
    const container = await renderLogin(onLogin);

    fill(container, "#username", "admin");
    fill(container, "#password", "wrong");
    const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
    await act(async () => button!.click());

    expect(container.textContent).toContain("invalid username or password");
    expect(container.querySelector('[role="alert"]')).not.toBeNull();
  });

  it("disables submit until both fields are filled", async () => {
    const container = await renderLogin(vi.fn());
    const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
    expect(button?.disabled).toBe(true);
    fill(container, "#username", "admin");
    expect(button?.disabled).toBe(true);
    fill(container, "#password", "admin");
    expect(button?.disabled).toBe(false);
  });

  it("shows the local seed hint in development builds (F-002-004)", async () => {
    const env = import.meta.env as { DEV: boolean };
    const dev = env.DEV;
    env.DEV = true;
    try {
      const container = await renderLogin(vi.fn());
      expect(container.textContent).toContain("admin / admin");
    } finally {
      env.DEV = dev;
    }
  });

  it("hides the local seed hint in production builds (F-002-004)", async () => {
    const env = import.meta.env as { DEV: boolean };
    const dev = env.DEV;
    env.DEV = false;
    try {
      const container = await renderLogin(vi.fn());
      expect(container.textContent).not.toContain("admin / admin");
      expect(container.textContent).not.toContain("Local development seed");
    } finally {
      env.DEV = dev;
    }
  });
});
