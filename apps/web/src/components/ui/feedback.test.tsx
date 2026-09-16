// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { FeedbackNoticeView } from "@/components/ui/feedback";
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
  vi.useRealTimers();
  vi.clearAllMocks();
});

async function renderFeedback(
  feedback: Parameters<typeof FeedbackNoticeView>[0]["feedback"],
  props: Omit<Parameters<typeof FeedbackNoticeView>[0], "feedback"> = {},
) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <FeedbackNoticeView feedback={feedback} {...props} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("FeedbackNoticeView (R4)", () => {
  it("uses status for success, auto-dismisses it, and keeps errors persistent", async () => {
    vi.useFakeTimers();
    const success = await renderFeedback({ kind: "success", message: "Saved" });
    expect(success.querySelector('[role="status"]')).not.toBeNull();
    await act(async () => vi.advanceTimersByTime(4000));
    expect(success.querySelector('[role="status"]')).toBeNull();

    const error = await renderFeedback({ kind: "error", code: "R4_ERR", message: "Failed" });
    await act(async () => vi.advanceTimersByTime(8000));
    expect(error.querySelector('[role="alert"]')).not.toBeNull();
    expect(error.querySelector('[data-feedback-code="R4_ERR"]')).not.toBeNull();
  });

  it("offers a keyboard-reachable retry once and a dismiss action", async () => {
    let resolveRetry: (() => void) | undefined;
    const retry = vi.fn(
      () =>
        new Promise<void>((resolve) => {
          resolveRetry = resolve;
        }),
    );
    const container = await renderFeedback(
      { kind: "error", code: "READ_DOWN", message: "Unavailable", retry },
      { surface: "inline", dismissible: true, retryButtonDataAttribute: "data-r4-retry" },
    );
    const retryButton = container.querySelector<HTMLButtonElement>("[data-r4-retry]");
    expect(retryButton).not.toBeNull();
    expect(retryButton?.type).toBe("button");
    await act(async () => {
      retryButton?.click();
      retryButton?.click();
      await Promise.resolve();
    });
    expect(retry).toHaveBeenCalledTimes(1);
    expect(retryButton?.disabled).toBe(true);
    await act(async () => {
      resolveRetry?.();
      await Promise.resolve();
    });
    expect(retryButton?.disabled).toBe(false);

    const dismiss = container.querySelector<HTMLButtonElement>("[aria-label]");
    expect(dismiss?.getAttribute("aria-label")).toBe("Dismiss");
    await act(async () => dismiss?.click());
    expect(container.querySelector('[role="alert"]')).toBeNull();
  });

  it("keeps the catalog copy primary while exposing a form diagnostic code", async () => {
    const container = await renderFeedback(
      {
        kind: "error",
        code: "INVALID_SITE_TITLE",
        message: "raw detail",
        messageKey: "error.invalidSiteTitle",
      },
      { surface: "inline", dismissible: false, showDiagnosticCode: true },
    );
    expect(container.querySelector('[role="alert"]')?.textContent).toContain("INVALID_SITE_TITLE");
    expect(container.querySelector('[role="alert"]')?.textContent).not.toContain("raw detail");
  });
});
