/**
 * Host failure result and recovery (ADR-0036 / spec 10 §3, `host.failure-recovery`).
 *
 * Production Host constructs the normalized failure result from the closed
 * elevation predicates; backend `code` stays opaque. `validateFailure`
 * enforces the semantic invariants (kind/hostCode pairing, scope/kind
 * compatibility, retry rules, per-kind recovery filtering) and
 * `validateReturnIntent` enforces the auth-return allowlist.
 */

export type HostScope = "bootstrap" | "manifest" | "page" | "auth" | "route" | "runtime";

export type HostKind =
  | "maintenance"
  | "upgrade-required"
  | "authentication-required"
  | "reauth-required"
  | "account-locked"
  | "forbidden"
  | "not-found"
  | "rate-limited"
  | "timeout"
  | "offline"
  | "protocol-rejected"
  | "render-failed"
  | "unavailable";

export type HostCode =
  | "HOST_MAINTENANCE"
  | "HOST_UPGRADE_REQUIRED"
  | "HOST_AUTH_REQUIRED"
  | "HOST_REAUTH_REQUIRED"
  | "HOST_ACCOUNT_LOCKED"
  | "HOST_FORBIDDEN"
  | "HOST_ROUTE_NOT_FOUND"
  | "HOST_RATE_LIMITED"
  | "HOST_TIMEOUT"
  | "HOST_OFFLINE"
  | "HOST_PROTOCOL_REJECTED"
  | "HOST_RENDER_FAILED"
  | "HOST_UNAVAILABLE";

export interface HostFailure {
  failureVersion: "1.0";
  failureId: string;
  scope: HostScope;
  kind: HostKind;
  hostCode: HostCode;
  retry?: { mode: "none" | "manual" | "after"; afterSeconds?: number };
  message: { messageKey: string; params?: Record<string, string | number | boolean> };
  correlation?: { requestId?: string; traceId?: string };
  diagnostics?: {
    phase?: string;
    protocolCode?: string;
    hostVersion?: string;
    protocolVersion?: string;
    manifestSha256?: string;
  };
  recoveryActions?: Array<{ type: "retry" | "reauth" | "home" | "back" | "reload" | "support"; url?: string }>;
}

const KIND_HOST_CODE: Record<HostKind, HostCode> = {
  "maintenance": "HOST_MAINTENANCE",
  "upgrade-required": "HOST_UPGRADE_REQUIRED",
  "authentication-required": "HOST_AUTH_REQUIRED",
  "reauth-required": "HOST_REAUTH_REQUIRED",
  "account-locked": "HOST_ACCOUNT_LOCKED",
  "forbidden": "HOST_FORBIDDEN",
  "not-found": "HOST_ROUTE_NOT_FOUND",
  "rate-limited": "HOST_RATE_LIMITED",
  "timeout": "HOST_TIMEOUT",
  "offline": "HOST_OFFLINE",
  "protocol-rejected": "HOST_PROTOCOL_REJECTED",
  "render-failed": "HOST_RENDER_FAILED",
  "unavailable": "HOST_UNAVAILABLE",
};

const SCOPE_KINDS: Record<HostScope, HostKind[]> = {
  bootstrap: ["maintenance", "upgrade-required", "rate-limited", "timeout", "offline", "unavailable", "protocol-rejected"],
  manifest: ["protocol-rejected", "rate-limited", "timeout", "offline", "unavailable", "authentication-required", "reauth-required", "forbidden"],
  page: ["protocol-rejected", "rate-limited", "timeout", "offline", "unavailable", "authentication-required", "reauth-required", "forbidden"],
  auth: ["authentication-required", "reauth-required", "account-locked", "forbidden"],
  route: ["not-found"],
  runtime: ["render-failed"],
};

const PERMANENT_REJECTED_QUERY_KEYS = new Set([
  "token", "access_token", "id_token", "code", "state", "session", "redirect", "returnto",
]);
export const PROTOCOL_RETURN_INTENT_ALLOWLIST = ["tab", "view", "page", "pageSize", "sort"];
const RETURN_INTENT_KEY_PATTERN = /^[a-z][a-zA-Z0-9_]*$/;
const HTML_PATTERN = /<\/?[a-zA-Z]/;
const RFC3339_PATTERN = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z$/;

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

/**
 * Host-level fetch classification (spec 10 §2.7 priority 1–2): 403 wins and
 * uniquely maps to forbidden; then 401 → auth; then explicit HTTP classes
 * (426, 429); then transport (timeout, offline); then other 5xx → unavailable.
 * Everything else returns null (stays in existing node/Action semantics).
 */
export function classifyHostFetch(input: {
  scope: "bootstrap" | "manifest" | "page";
  status: number | null;
  authState: "anonymous" | "authenticated";
  transport: "timeout" | "offline" | null;
}): { scope: HostScope; kind: HostKind; hostCode: HostCode } | null {
  const { scope, status, authState, transport } = input;
  if (scope !== "bootstrap" && scope !== "manifest" && scope !== "page") return null;
  if (status === 403) return { scope, kind: "forbidden", hostCode: "HOST_FORBIDDEN" };
  if (status === 401) {
    if (authState === "authenticated") {
      return { scope, kind: "reauth-required", hostCode: "HOST_REAUTH_REQUIRED" };
    }
    return { scope, kind: "authentication-required", hostCode: "HOST_AUTH_REQUIRED" };
  }
  if (status === 426) return { scope, kind: "upgrade-required", hostCode: "HOST_UPGRADE_REQUIRED" };
  if (status === 429) return { scope, kind: "rate-limited", hostCode: "HOST_RATE_LIMITED" };
  if (transport === "timeout") return { scope, kind: "timeout", hostCode: "HOST_TIMEOUT" };
  if (transport === "offline") return { scope, kind: "offline", hostCode: "HOST_OFFLINE" };
  if (typeof status === "number" && status >= 500 && status <= 599) {
    return { scope, kind: "unavailable", hostCode: "HOST_UNAVAILABLE" };
  }
  return null;
}

/** Bootstrap stable result → Host failure triple (ADR-0035 D7 / spec 10 §2.7). */
export function mapBootstrapResult(input: {
  result: string;
  fetchClassification?: string | null;
}): { scope: HostScope; kind: HostKind; hostCode: HostCode } | null {
  const { result, fetchClassification } = input;
  switch (result) {
    case "MAINTENANCE":
      return { scope: "bootstrap", kind: "maintenance", hostCode: "HOST_MAINTENANCE" };
    case "UPGRADE_REQUIRED":
      return { scope: "bootstrap", kind: "upgrade-required", hostCode: "HOST_UPGRADE_REQUIRED" };
    case "REAUTH_REQUIRED":
      return { scope: "auth", kind: "reauth-required", hostCode: "HOST_REAUTH_REQUIRED" };
    case "ACCOUNT_LOCKED":
      return { scope: "auth", kind: "account-locked", hostCode: "HOST_ACCOUNT_LOCKED" };
    case "BOOTSTRAP_NEGOTIATION_REJECTED":
      return { scope: "bootstrap", kind: "protocol-rejected", hostCode: "HOST_PROTOCOL_REJECTED" };
    case "MANIFEST_CAPABILITY_REJECTED":
      return { scope: "manifest", kind: "protocol-rejected", hostCode: "HOST_PROTOCOL_REJECTED" };
    case "MANIFEST_INTEGRITY_FAILED":
      return { scope: "manifest", kind: "protocol-rejected", hostCode: "HOST_PROTOCOL_REJECTED" };
    case "BOOTSTRAP_DOCUMENT_FAILED": {
      switch (fetchClassification) {
        case "rate-limited":
          return { scope: "bootstrap", kind: "rate-limited", hostCode: "HOST_RATE_LIMITED" };
        case "timeout":
          return { scope: "bootstrap", kind: "timeout", hostCode: "HOST_TIMEOUT" };
        case "offline":
          return { scope: "bootstrap", kind: "offline", hostCode: "HOST_OFFLINE" };
        case "unavailable":
          return { scope: "bootstrap", kind: "unavailable", hostCode: "HOST_UNAVAILABLE" };
        case "protocol":
          return { scope: "bootstrap", kind: "protocol-rejected", hostCode: "HOST_PROTOCOL_REJECTED" };
        default:
          return null;
      }
    }
    default:
      return null;
  }
}

/** Generates a unique printable failure ID per occurrence within an app instance. */
let failureSequence = 0;
export function nextFailureId(): string {
  failureSequence += 1;
  return `hf-${Date.now().toString(36)}-${failureSequence.toString(36)}`;
}

/** Validates a produced failure result against the semantic invariants. */
export function validateFailure(failure: HostFailure): { valid: boolean; errors: string[] } {
  const errors: string[] = [];
  if (!isRecord(failure)) return { valid: false, errors: ["failure must be an object"] };

  const pair = KIND_HOST_CODE[failure.kind as HostKind];
  if (pair === undefined) {
    errors.push("unknown kind");
  } else if (failure.hostCode !== pair) {
    errors.push(`hostCode ${failure.hostCode} does not match kind ${failure.kind} (expected ${pair})`);
  }

  const allowedKinds = SCOPE_KINDS[failure.scope as HostScope];
  if (allowedKinds === undefined) {
    errors.push(`unknown scope ${failure.scope}`);
  } else if (pair !== undefined && !allowedKinds.includes(failure.kind as HostKind)) {
    errors.push(`kind ${failure.kind} is not valid in scope ${failure.scope}`);
  }

  if (failure.retry !== undefined && failure.retry !== null) {
    if (failure.retry.mode === "after" && !Number.isInteger(failure.retry.afterSeconds)) {
      errors.push("retry mode after requires positive integer afterSeconds");
    }
    if (failure.kind === "render-failed" && failure.retry.mode === "after") {
      errors.push("render-failed must not auto-loop reload (retry mode after forbidden)");
    }
  }

  const message = failure.message;
  if (!isRecord(message) || typeof message.messageKey !== "string" || message.messageKey.length === 0) {
    errors.push("message.messageKey is required");
  } else if (isRecord(message.params)) {
    for (const value of Object.values(message.params)) {
      if (typeof value === "string" && HTML_PATTERN.test(value)) {
        errors.push("message.params must not contain HTML");
      }
    }
  }

  for (const action of failure.recoveryActions ?? []) {
    if (!isRecord(action)) continue;
    if (failure.kind === "forbidden" && action.type === "reauth") {
      errors.push("forbidden must not offer reauth");
    }
    if (failure.kind === "account-locked") {
      if (action.type === "reauth") errors.push("account-locked must not offer reauth");
      if (action.type !== "home" && action.type !== "support") {
        errors.push("account-locked only allows home/support recovery actions");
      }
    }
  }

  return { valid: errors.length === 0, errors };
}

export interface ReturnIntent {
  path: string;
  query?: Record<string, string>;
  expiresAt: string;
  nonce: string;
}

export interface ReturnIntentValidation {
  valid: boolean;
  keptQuery: Record<string, string>;
  droppedKeys: string[];
  rejectedKeys: string[];
  reason: string | null;
}

/**
 * Validates a recoverable auth return intent (ADR-0036 D6 / spec 10 §3.7).
 * The Host computes its effective allowlist by only narrowing the protocol
 * allowlist plus the current page's registered `returnIntentQueryKeys`;
 * sensitive keys are permanently rejected case-insensitively.
 */
export function validateReturnIntent(
  intent: ReturnIntent,
  options: { registeredKeys?: string[]; nowIso: string },
): ReturnIntentValidation {
  const errors: string[] = [];
  const rejectedKeys: string[] = [];
  const droppedKeys: string[] = [];
  const keptQuery: Record<string, string> = {};

  if (!isRecord(intent)) {
    return { valid: false, keptQuery, droppedKeys, rejectedKeys, reason: "intent must be an object" };
  }

  const path = intent.path;
  if (typeof path !== "string"
    || !path.startsWith("/")
    || path.includes("#")
    || path.includes("://")
    || path.includes("\\")) {
    errors.push("path must be an absolute in-app path without scheme, authority, fragment or backslash");
  }
  if (typeof intent.nonce !== "string" || intent.nonce.length === 0) {
    errors.push("nonce is required");
  }
  if (typeof intent.expiresAt !== "string" || !RFC3339_PATTERN.test(intent.expiresAt)) {
    errors.push("expiresAt must be RFC 3339 UTC");
  } else if (Date.parse(intent.expiresAt) <= Date.parse(options.nowIso)) {
    errors.push("intent has expired");
  }

  const hostExtensions = (options.registeredKeys ?? []).filter(
    (key) => RETURN_INTENT_KEY_PATTERN.test(key),
  );
  const allowlist = new Set([...PROTOCOL_RETURN_INTENT_ALLOWLIST, ...hostExtensions]);

  if (isRecord(intent.query)) {
    for (const [key, value] of Object.entries(intent.query)) {
      if (typeof value !== "string") {
        errors.push(`query value for ${key} must be a string`);
        continue;
      }
      const lowered = key.toLowerCase();
      if (PERMANENT_REJECTED_QUERY_KEYS.has(lowered)) {
        rejectedKeys.push(key);
        continue;
      }
      if (!allowlist.has(key)) {
        droppedKeys.push(key);
        continue;
      }
      keptQuery[key] = value;
    }
  }

  return {
    valid: errors.length === 0,
    keptQuery,
    droppedKeys,
    rejectedKeys,
    reason: errors.length > 0 ? errors.join("; ") : null,
  };
}
