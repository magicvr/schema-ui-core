// @vitest-environment jsdom
// W15 F-005 regression lock (A-003 F-008 · recommended): the invitation token
// is a one-time bearer with a multi-day TTL — the page must scrub it from the
// URL the moment it mounts, keep unrelated query params, and stay idempotent
// (a second mount with an already-clean URL must not touch history again).
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { InviteAcceptPage } from "@/components/invite-accept";
import { I18nProvider } from "@/i18n/runtime";

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
  vi.restoreAllMocks();
});

async function mountWithUrl(url: string) {
  // Set the URL BEFORE spying so the setup rewrite is not counted.
  window.history.replaceState(window.history.state, "", url);
  const replaceState = vi.spyOn(window.history, "replaceState");
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <InviteAcceptPage />
      </I18nProvider>,
    );
  });
  return replaceState;
}

describe("W15 F-005 · invite token URL scrub", () => {
  it("removes token from the query while keeping unrelated params", async () => {
    const replaceState = await mountWithUrl("/invite/accept?token=one-time-7d&utm=mail");
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(replaceState).toHaveBeenCalled();
    const params = new URLSearchParams(window.location.search);
    expect(params.get("token")).toBeNull();
    expect(params.get("utm")).toBe("mail");
  });

  it("mounts without a token without touching history", async () => {
    const replaceState = await mountWithUrl("/invite/accept");
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(replaceState).not.toHaveBeenCalled();
  });

  it("is idempotent: a second mount on the cleaned URL performs no extra scrub", async () => {
    const replaceState = await mountWithUrl("/invite/accept?token=one-time-7d");
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    const callsAfterFirstMount = replaceState.mock.calls.length;
    expect(callsAfterFirstMount).toBeGreaterThan(0);

    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <I18nProvider>
          <InviteAcceptPage />
        </I18nProvider>,
      );
    });
    expect(replaceState.mock.calls.length).toBe(callsAfterFirstMount);
  });
});