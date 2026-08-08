/**
 * Shared pure state determination for async display regions (S4 · GOAL-004).
 *
 * statCard / chart / list-table each fetch a resource independently and used
 * to invent their own ad-hoc "Loading…" text placeholder. This module
 * centralizes the loading / error / empty / ready decision into one pure,
 * directly-testable function so every consumer renders the same sequence
 * (Skeleton while loading, a `role="alert"` message on error, a muted empty
 * message otherwise) instead of drifting independently.
 */

export type AsyncDisplayState = "loading" | "error" | "empty" | "ready";

export interface AsyncDisplayInput {
  /** True while the underlying fetch/request has not yet settled. */
  loading: boolean;
  /** Non-null when the fetch/request failed. */
  error: string | null;
  /** True when the fetch succeeded but produced no renderable rows/points. */
  isEmpty?: boolean;
}

/**
 * Resolves the single display state a region should show.
 *
 * Precedence: `error` wins over `loading` (a failed fetch is not "still
 * loading" even if a stale loading flag lingers), and `loading` wins over
 * `isEmpty` (emptiness is unknown until the fetch settles).
 */
export function resolveAsyncDisplayState({
  loading,
  error,
  isEmpty = false,
}: AsyncDisplayInput): AsyncDisplayState {
  if (error !== null) {
    return "error";
  }
  if (loading) {
    return "loading";
  }
  if (isEmpty) {
    return "empty";
  }
  return "ready";
}
