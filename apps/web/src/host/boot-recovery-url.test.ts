// @vitest-environment jsdom

import { afterEach, describe, expect, it, vi } from "vitest";

import { executeBootRecovery, isSafeSupportUrl } from "@/host/boot";

/**
 * W13 F-015 (GOAL-013 A-001): the "support" recovery action used to assign
 * action.url to an anchor href unchecked — a javascript: value would execute
 * on click. Only http(s) targets (absolute or same-origin relative) may pass.
 */

describe("boot recovery support URL safety (W13 F-015)", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it.each(["https://support.example.com/help", "http://localhost:8080/help", "/relative-help"])(
    "navigates an http(s) target: %s",
    (url) => {
      expect(isSafeSupportUrl(url)).toBe(true);
    },
  );

  it.each([
    "javascript:alert(1)",
    "JAVASCRIPT:alert(1)",
    "data:text/html,<script>alert(1)</script>",
    "vbscript:msgbox",
    "file:///C:/Windows/system32/drivers/etc/hosts",
  ])("rejects a dangerous scheme: %s", (url) => {
    expect(isSafeSupportUrl(url)).toBe(false);
  });

  it("never builds an anchor for a javascript: url", () => {
    const click = vi.fn();
    const create = vi.spyOn(document, "createElement").mockReturnValue({
      href: "",
      rel: "",
      click,
    } as unknown as HTMLAnchorElement);
    executeBootRecovery({ type: "support", url: "javascript:alert(1)" });
    expect(create).not.toHaveBeenCalled();
    expect(click).not.toHaveBeenCalled();
  });
});
