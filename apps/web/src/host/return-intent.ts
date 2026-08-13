/**
 * Recoverable auth return intent (ADR-0036 D6 / spec 10 §3.7) — production
 * capture/consume wiring.
 *
 * The pure validator lives in `host/failure.ts` (`validateReturnIntent`, pinned
 * by the vendored upstream host-failure fixtures). This module owns the
 * Host-side lifecycle the protocol leaves to the Host:
 *
 * - capture: auth loss (boot reauth terminal or mid-session session loss) stores
 *   a single-use intent bound to the current in-app path/query;
 * - consume: a successful login validates it (allowlist narrowing, expiry,
 *   nonce single-use, login-path self-loop rejection) and yields the target
 *   path; the caller restores navigation.
 *
 * The boot-time auth terminal happens BEFORE the manifest is loaded, so capture
 * uses only the protocol allowlist (`validateReturnIntent` without
 * `registeredKeys`) — the manifest `returnIntentQueryKeys` extension is a
 * narrowing-only source per ADR-0036 D6, and omitting it here can only be
 * narrower, never wider.
 */

import { validateReturnIntent, type ReturnIntent } from "@/host/failure";

const STORAGE_KEY = "schema-ui.return-intent";
const INTENT_TTL_MS = 10 * 60 * 1000;
const LOGIN_PATHS = new Set(["/login", "/login/"]);

export interface ReturnIntentTarget {
  path: string;
  query: Record<string, string>;
}

function buildQueryString(query: Record<string, string>): string {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(query)) {
    params.set(key, value);
  }
  const encoded = params.toString();
  return encoded === "" ? "" : `?${encoded}`;
}

/** Serializes the current location into an intent (null when nothing to restore). */
function intentFromLocation(path: string, query: Record<string, string>, nowIso: string): ReturnIntent | null {
  if (!path.startsWith("/") || LOGIN_PATHS.has(path)) {
    return null;
  }
  const kept = Object.fromEntries(
    Object.entries(query).filter(([, value]) => typeof value === "string"),
  );
  return {
    path,
    query: kept,
    expiresAt: new Date(Date.parse(nowIso) + INTENT_TTL_MS).toISOString(),
    nonce: crypto.randomUUID(),
  };
}

/**
 * Captures a return intent for the current location (auth loss paths).
 * No-op when the location is the login surface itself (self-loop prevention).
 */
export function captureReturnIntent(options?: {
  path?: string;
  query?: Record<string, string>;
  nowIso?: string;
  storage?: Storage;
}): void {
  const storage = options?.storage ?? window.sessionStorage;
  const nowIso = options?.nowIso ?? new Date().toISOString();
  const intent = intentFromLocation(
    options?.path ?? window.location.pathname,
    options?.query ?? {},
    nowIso,
  );
  if (intent === null) {
    return;
  }
  storage.setItem(STORAGE_KEY, JSON.stringify(intent));
}

/**
 * Consumes the stored intent after a successful login: validates it (protocol
 * allowlist, expiry, nonce) and removes it from storage — nonce is single-use
 * and consumed on read, whether valid or not (no replay).
 */
export function takeReturnIntent(options?: {
  nowIso?: string;
  storage?: Storage;
}): ReturnIntentTarget | null {
  const storage = options?.storage ?? window.sessionStorage;
  const nowIso = options?.nowIso ?? new Date().toISOString();
  const raw = storage.getItem(STORAGE_KEY);
  if (raw === null) {
    return null;
  }
  storage.removeItem(STORAGE_KEY);
  let intent: unknown;
  try {
    intent = JSON.parse(raw) as unknown;
  } catch {
    return null;
  }
  const result = validateReturnIntent(intent as ReturnIntent, { nowIso });
  if (!result.valid || LOGIN_PATHS.has((intent as ReturnIntent).path)) {
    return null;
  }
  return { path: (intent as ReturnIntent).path, query: result.keptQuery };
}

/** Returns the storage key (test seam). */
export function returnIntentStorageKey(): string {
  return STORAGE_KEY;
}

export { buildQueryString };
