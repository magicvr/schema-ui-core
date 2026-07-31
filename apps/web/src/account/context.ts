import type { NavigationContext } from "@/protocol/app-manifest";

export const DEFAULT_ACCOUNT_URL = "/api/accounts/me";

/**
 * Loads the R4 account session snapshot. Returns a NavigationContext built
 * from GET /api/accounts/me, or an empty context on failure. The empty
 * fallback is a fail-closed renderer posture (missing context denies
 * expressions); callers can distinguish failure via the error return.
 */
export async function loadAccountContext(
  fetcher: typeof fetch = fetch,
  url = DEFAULT_ACCOUNT_URL,
): Promise<{ context: NavigationContext; error: unknown | null }> {
  try {
    const response = await fetcher(url);
    if (!response.ok) {
      return { context: {}, error: new Error(`account fetch failed: HTTP ${response.status}`) };
    }
    const body = (await response.json()) as {
      user?: Record<string, unknown>;
      features?: Record<string, unknown>;
    };
    return {
      context: {
        ...(body.user === undefined ? {} : { user: body.user }),
        ...(body.features === undefined ? {} : { features: body.features }),
      },
      error: null,
    };
  } catch (error) {
    return { context: {}, error };
  }
}
