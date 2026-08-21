// @vitest-environment jsdom

/**
 * W10 F-002: fetch timeout wrapper. A hung request must abort after the
 * bounded window; a caller-provided signal must compose with the timeout;
 * an already-aborted caller signal fails fast without any fetch call.
 */

import { describe, expect, it, vi } from "vitest";

import { DEFAULT_FETCH_TIMEOUT_MS, withTimeout } from "@/lib/fetch-timeout";

/** A fetch double that hangs until its signal aborts, then rejects. */
function hangingFetch(): { impl: typeof fetch; seen: () => AbortSignal | undefined } {
  let captured: AbortSignal | null | undefined;
  const impl = (async (_input: RequestInfo | URL, init?: RequestInit) => {
    captured = init?.signal;
    return new Promise<Response>((_resolve, reject) => {
      init?.signal?.addEventListener("abort", () => reject(new Error("aborted")), { once: true });
    });
  }) as unknown as typeof fetch;
  return { impl, seen: () => captured ?? undefined };
}

describe("withTimeout fetch wrapper (W10 F-002)", () => {
  it("passes through a successful response and clears the timer", async () => {
    const response = new Response("{}", { status: 200 });
    const impl = vi.fn().mockResolvedValue(response) as unknown as typeof fetch;
    const wrapped = withTimeout(impl);
    const result = await wrapped("/api/x");
    expect(result).toBe(response);
    expect(impl).toHaveBeenCalledWith("/api/x", expect.objectContaining({ signal: expect.any(AbortSignal) }));
  });

  it("aborts a hung request after the timeout window", async () => {
    const { impl } = hangingFetch();
    const wrapped = withTimeout(impl, 20);
    await expect(wrapped("/api/hang")).rejects.toThrow("aborted");
  });

  it("composes a caller-provided abort signal with the timeout signal", async () => {
    const { impl, seen } = hangingFetch();
    const caller = new AbortController();
    const wrapped = withTimeout(impl, 5_000);
    const promise = wrapped("/api/x", { signal: caller.signal });
    // Yield a microtask so the wrapper attaches its relay listener first.
    await Promise.resolve();
    caller.abort();
    await expect(promise).rejects.toThrow("aborted");
    expect(seen()?.aborted).toBe(true);
  });

  it("fails fast without fetching when the caller signal is already aborted", async () => {
    const impl = vi.fn().mockResolvedValue(new Response("{}", { status: 200 })) as unknown as typeof fetch;
    const wrapped = withTimeout(impl);
    const caller = new AbortController();
    caller.abort();
    await expect(wrapped("/api/x", { signal: caller.signal })).rejects.toThrow();
    expect(impl).not.toHaveBeenCalled();
  });

  it("defaults to a 30s ceiling", () => {
    expect(DEFAULT_FETCH_TIMEOUT_MS).toBe(30_000);
  });
});
