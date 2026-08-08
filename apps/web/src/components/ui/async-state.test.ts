import { describe, expect, it } from "vitest";

import { resolveAsyncDisplayState } from "@/components/ui/async-state";

describe("resolveAsyncDisplayState (S4 · GOAL-004)", () => {
  it("resolves to error whenever an error string is present, even if loading is stale-true", () => {
    expect(resolveAsyncDisplayState({ loading: true, error: "boom" })).toBe("error");
    expect(resolveAsyncDisplayState({ loading: false, error: "boom" })).toBe("error");
  });

  it("resolves to loading when no error and the fetch has not settled", () => {
    expect(resolveAsyncDisplayState({ loading: true, error: null })).toBe("loading");
    expect(resolveAsyncDisplayState({ loading: true, error: null, isEmpty: true })).toBe("loading");
  });

  it("resolves to empty only once loading has settled with no rows/points", () => {
    expect(resolveAsyncDisplayState({ loading: false, error: null, isEmpty: true })).toBe("empty");
  });

  it("resolves to ready once loading has settled with data and no error", () => {
    expect(resolveAsyncDisplayState({ loading: false, error: null })).toBe("ready");
    expect(resolveAsyncDisplayState({ loading: false, error: null, isEmpty: false })).toBe("ready");
  });
});
