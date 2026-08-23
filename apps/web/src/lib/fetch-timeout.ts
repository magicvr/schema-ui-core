/**
 * W10 F-002: fetch timeout wrapper.
 *
 * Wraps a fetch implementation so a hung request (slow network, stalled
 * server) aborts after a bounded window instead of pending forever. The
 * returned function is drop-in compatible with `typeof fetch`; a caller-
 * provided AbortSignal in `init` is composed with the timeout signal (either
 * firing aborts the request). Timers are cleared once the request settles, so
 * no handles leak on the happy path.
 */

/** Default ceiling for a single request (30s). */
export const DEFAULT_FETCH_TIMEOUT_MS = 30_000;

export function withTimeout(
  fetchImpl?: typeof fetch,
  timeoutMs: number = DEFAULT_FETCH_TIMEOUT_MS,
): typeof fetch {
  return async function fetchWithTimeout(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
    // Resolved per call (not at wrap time) so a later global fetch stub —
    // test doubles, service workers — is always honored.
    const doFetch = fetchImpl ?? globalThis.fetch;
    // RequestInit.signal is `AbortSignal | null | undefined`; normalize to
    // undefined so the guards below stay simple.
    const outer = init?.signal ?? undefined;
    // An already-aborted caller signal must fail fast without touching the
    // network (no fetch call, no timer).
    if (outer?.aborted === true) {
      throw new DOMException("The operation was aborted.", "AbortError");
    }
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), timeoutMs);
    // Relay the caller signal onto our controller; the relay listener is
    // removed once the request settles so a long-lived shared signal does not
    // accumulate listeners across calls (A-003 recommended F-002).
    const relayAbort = () => controller.abort();
    if (outer !== undefined) {
      outer.addEventListener("abort", relayAbort, { once: true });
    }
    try {
      return await doFetch(input, { ...init, signal: controller.signal });
    } finally {
      clearTimeout(timer);
      outer?.removeEventListener("abort", relayAbort);
    }
  };
}
