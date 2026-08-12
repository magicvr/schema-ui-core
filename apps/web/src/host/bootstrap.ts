/**
 * Host bootstrap runtime (ADR-0035 / spec 10 §2, `host.bootstrap`).
 *
 * Production entry: `discoverBootstrapDocument()` fetches the optional public
 * bootstrap document from the real entry; `evaluateBootstrap()` implements
 * the deterministic lifecycle stages. This module is consumed by the
 * production boot path (`main.tsx`) — it is not a fixture adapter.
 *
 * Pinned upstream machine contracts: `src/protocol/upstream/provenance-v2.8-candidate.json`.
 */

const CAPABILITY_PATTERN = /^[a-z][a-z0-9-]*(?:\.[a-z][a-z0-9-]*)*$/;
const BOOTSTRAP_CONTENT_TYPES = new Set([
  "application/json",
  "application/schema-ui+json",
]);

export const DEFAULT_BOOTSTRAP_PATH = "/.well-known/schema-ui/host-bootstrap.json";
export const BOOTSTRAP_VERSION = "1.0" as const;

export type BootstrapAvailabilityMode = "normal" | "maintenance" | "upgrade-required" | "degraded";

export interface BootstrapAvailability {
  mode: BootstrapAvailabilityMode;
  messageKey?: string;
  retryAfterSeconds?: number;
  minimumHostVersion?: string;
  disabledCapabilities?: string[];
}

export interface BootstrapDocument {
  bootstrapVersion: string;
  requiredCapabilities: string[];
  manifest: { url: string; sha256?: string };
  availability: BootstrapAvailability;
}

export type BootstrapAuthState = "anonymous" | "authenticated" | "reauth-required" | "locked";

export interface BootstrapAuth {
  state: BootstrapAuthState;
  principal?: { id?: string; name?: string; roles?: string[] };
  expiresAt?: string;
  provenance?: string;
}

export interface HostSupport {
  supportedBootstrapVersions: string[];
  supportedCapabilities: string[];
}

export type BootstrapResultCode =
  | "OK"
  | "INVALID_HOST_SUPPORT"
  | "INVALID_REQUIRED_CAPABILITIES"
  | "UNSUPPORTED_BOOTSTRAP_VERSION"
  | "INVALID_BOOTSTRAP_DOCUMENT"
  | "MISSING_REQUIRED_CAPABILITY"
  | "BOOTSTRAP_DOCUMENT_FAILED"
  | "MANIFEST_CAPABILITY_REJECTED"
  | "MANIFEST_INTEGRITY_FAILED";

export type BootstrapResult =
  | "READY"
  | "READY_DEGRADED"
  | "MAINTENANCE"
  | "UPGRADE_REQUIRED"
  | "REAUTH_REQUIRED"
  | "ACCOUNT_LOCKED"
  | "BOOTSTRAP_DOCUMENT_FAILED"
  | "BOOTSTRAP_NEGOTIATION_REJECTED"
  | "MANIFEST_CAPABILITY_REJECTED"
  | "MANIFEST_INTEGRITY_FAILED";

export type BootstrapFetchClassification =
  | "rate-limited"
  | "timeout"
  | "offline"
  | "unavailable"
  | "protocol";

export interface BootstrapEvaluation {
  code: BootstrapResultCode;
  result: BootstrapResult;
  phase: string;
  fetchClassification: BootstrapFetchClassification | null;
  missingCapabilities: string[];
  effectiveCapabilities: string[] | null;
  context: { user: { id: string; name: string; roles: string[] } } | null;
}

export interface BootstrapDiscovery {
  status: "ok" | "not-provided" | "failed";
  document: BootstrapDocument | null;
  classification: BootstrapFetchClassification | null;
  bytesSha256: string | null;
}

const ANONYMOUS_USER = { id: "", name: "", roles: [] };

function isUniqueStringList(value: unknown, allowEmpty: boolean): value is string[] {
  return Array.isArray(value)
    && (allowEmpty || value.length > 0)
    && value.every((item) => typeof item === "string")
    && new Set(value).size === value.length;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function buildResult(
  code: BootstrapResultCode,
  result: BootstrapResult,
  phase: string,
  extra: Partial<BootstrapEvaluation> = {},
): BootstrapEvaluation {
  return {
    code,
    result,
    phase,
    fetchClassification: null,
    missingCapabilities: [],
    effectiveCapabilities: null,
    context: null,
    ...extra,
  };
}

function isValidBootstrapDocument(value: unknown): value is BootstrapDocument {
  if (!isRecord(value)) return false;
  if (value.bootstrapVersion !== BOOTSTRAP_VERSION) return false;
  if (!isUniqueStringList(value.requiredCapabilities, false)) return false;
  if (!value.requiredCapabilities.every((item) => CAPABILITY_PATTERN.test(item))) return false;
  if (!value.requiredCapabilities.includes("host.bootstrap")) return false;
  if (!isRecord(value.manifest)) return false;
  if (typeof value.manifest.url !== "string" || value.manifest.url.length === 0) return false;
  if (value.manifest.sha256 !== undefined && typeof value.manifest.sha256 !== "string") return false;
  if (!isRecord(value.availability)) return false;
  const mode = value.availability.mode;
  if (mode !== "normal" && mode !== "maintenance" && mode !== "upgrade-required" && mode !== "degraded") {
    return false;
  }
  if (value.availability.disabledCapabilities !== undefined
    && !isUniqueStringList(value.availability.disabledCapabilities, false)) return false;
  return true;
}

/**
 * Discovery (stage 1): GET the default (or explicit) bootstrap URL with
 * `credentials: omit`. Only 200 succeeds; 404/410 on the default entry means
 * "not provided" (fallback to the ADR-0025 manifest entry); every other
 * status, redirect, wrong content type or parse failure is fail-closed.
 */
export async function discoverBootstrapDocument(options: {
  url?: string;
  fetcher?: typeof fetch;
} = {}): Promise<BootstrapDiscovery> {
  const url = options.url ?? DEFAULT_BOOTSTRAP_PATH;
  const fetcher = options.fetcher ?? globalThis.fetch;
  if (!fetcher) {
    return { status: "failed", document: null, classification: "offline", bytesSha256: null };
  }
  let response: Response;
  try {
    response = await fetcher(url, { credentials: "omit", redirect: "manual" });
  } catch {
    // Transport unreachable (offline) or request timeout are indistinguishable
    // at this layer; callers overlay their own timeout classification.
    return { status: "failed", document: null, classification: "offline", bytesSha256: null };
  }
  if (response.status === 404 || response.status === 410) {
    if (url === DEFAULT_BOOTSTRAP_PATH) {
      return { status: "not-provided", document: null, classification: null, bytesSha256: null };
    }
    return { status: "failed", document: null, classification: "protocol", bytesSha256: null };
  }
  if (response.status === 429) {
    return { status: "failed", document: null, classification: "rate-limited", bytesSha256: null };
  }
  if (response.status !== 200) {
    return { status: "failed", document: null, classification: "unavailable", bytesSha256: null };
  }
  const contentType = (response.headers.get("content-type") ?? "")
    .split(";")[0].trim().toLowerCase();
  if (!BOOTSTRAP_CONTENT_TYPES.has(contentType)) {
    return { status: "failed", document: null, classification: "protocol", bytesSha256: null };
  }
  try {
    const bytes = new Uint8Array(await response.arrayBuffer());
    const text = new TextDecoder("utf-8", { fatal: true }).decode(bytes);
    const payload: unknown = JSON.parse(text);
    if (!isValidBootstrapDocument(payload)) {
      return { status: "failed", document: null, classification: "protocol", bytesSha256: null };
    }
    return {
      status: "ok",
      document: payload,
      classification: null,
      bytesSha256: await sha256Hex(bytes),
    };
  } catch {
    return { status: "failed", document: null, classification: "protocol", bytesSha256: null };
  }
}

/**
 * Deterministic lifecycle evaluation (stages 2–9). Mirrors the upstream B1
 * reference (`host-bootstrap.js`) so the vendored conformance fixtures must
 * pass against this production implementation without exclusions.
 */
export function evaluateBootstrap(input: {
  document: BootstrapDocument | null;
  fetch: { status: "ok" | "not-provided" | "failed"; classification?: string | null };
  hostSupport: HostSupport;
  auth: BootstrapAuth;
  manifest: { protocolVersion: string; requiredCapabilities?: string[] } | null;
  integrity: { declaredSha256: string; computedSha256: string } | null;
  capabilityRegistry: { capabilities: Record<string, unknown> };
}): BootstrapEvaluation {
  const { document, fetch, hostSupport, auth, manifest, integrity, capabilityRegistry } = input;
  const registryCapabilities = isRecord(capabilityRegistry?.capabilities)
    ? capabilityRegistry.capabilities
    : {};

  if (document === null) {
    if (fetch.status === "not-provided") {
      return runFallbackPath({ hostSupport, auth, manifest });
    }
    const classification = fetch.classification === "rate-limited"
      || fetch.classification === "timeout"
      || fetch.classification === "offline"
      || fetch.classification === "protocol"
      ? fetch.classification
      : "unavailable";
    return buildResult(
      "BOOTSTRAP_DOCUMENT_FAILED",
      "BOOTSTRAP_DOCUMENT_FAILED",
      "bootstrap-discovery",
      { fetchClassification: classification },
    );
  }

  // Stage 2: bootstrap-validation.
  if (!isUniqueStringList(hostSupport.supportedBootstrapVersions, false)
    || hostSupport.supportedBootstrapVersions.some((item) => item.length === 0)) {
    return buildResult("INVALID_HOST_SUPPORT", "BOOTSTRAP_NEGOTIATION_REJECTED", "bootstrap-validation");
  }
  if (!hostSupport.supportedBootstrapVersions.includes(document.bootstrapVersion)) {
    return buildResult("UNSUPPORTED_BOOTSTRAP_VERSION", "BOOTSTRAP_NEGOTIATION_REJECTED", "bootstrap-validation");
  }
  if (!isUniqueStringList(document.requiredCapabilities, false)
    || document.requiredCapabilities.some((item) => !CAPABILITY_PATTERN.test(item))
    || !document.requiredCapabilities.includes("host.bootstrap")) {
    return buildResult("INVALID_REQUIRED_CAPABILITIES", "BOOTSTRAP_NEGOTIATION_REJECTED", "bootstrap-validation");
  }
  const supportedCapabilities = Array.isArray(hostSupport.supportedCapabilities)
    ? hostSupport.supportedCapabilities
    : [];
  const supportedSet = new Set(supportedCapabilities);
  const missingCapabilities = document.requiredCapabilities.filter(
    (capability) => !supportedSet.has(capability),
  );
  if (missingCapabilities.length > 0) {
    return buildResult(
      "MISSING_REQUIRED_CAPABILITY",
      "BOOTSTRAP_NEGOTIATION_REJECTED",
      "bootstrap-validation",
      { missingCapabilities },
    );
  }

  // Stage 3: availability-gate.
  const mode = document.availability.mode;
  if (mode === "maintenance") return buildResult("OK", "MAINTENANCE", "availability-gate");
  if (mode === "upgrade-required") return buildResult("OK", "UPGRADE_REQUIRED", "availability-gate");

  let disabledCapabilities: string[] = [];
  if (mode === "degraded") {
    disabledCapabilities = document.availability.disabledCapabilities ?? [];
    if (disabledCapabilities.some((capability) => !Object.hasOwn(registryCapabilities, capability))) {
      return buildResult("INVALID_BOOTSTRAP_DOCUMENT", "BOOTSTRAP_DOCUMENT_FAILED", "bootstrap-validation");
    }
  }

  // Stage 4: auth-resolution.
  if (auth.state === "locked") return buildResult("OK", "ACCOUNT_LOCKED", "auth-resolution");
  if (auth.state === "reauth-required") return buildResult("OK", "REAUTH_REQUIRED", "auth-resolution");

  // Stage 6: manifest-integrity.
  if (document.manifest.sha256 !== undefined) {
    if (integrity === null || integrity.computedSha256 !== document.manifest.sha256) {
      return buildResult("MANIFEST_INTEGRITY_FAILED", "MANIFEST_INTEGRITY_FAILED", "manifest-integrity");
    }
  }

  // Stage 7: manifest capability narrowing (D5).
  const effectiveCapabilities = supportedCapabilities.filter(
    (capability) => !disabledCapabilities.includes(capability),
  );
  const manifestRequired = manifest?.requiredCapabilities ?? [];
  const effectiveSet = new Set(effectiveCapabilities);
  const missingAfterNarrowing = manifestRequired.filter(
    (capability) => !effectiveSet.has(capability),
  );
  if (missingAfterNarrowing.length > 0) {
    return buildResult(
      "MISSING_REQUIRED_CAPABILITY",
      "MANIFEST_CAPABILITY_REJECTED",
      "manifest-validation",
      { missingCapabilities: missingAfterNarrowing, effectiveCapabilities },
    );
  }

  // Stage 8–9: context-resolution + ready.
  const context = resolveContext(auth);
  if (mode === "degraded") {
    return buildResult("OK", "READY_DEGRADED", "ready", { context, effectiveCapabilities });
  }
  return buildResult("OK", "READY", "ready", { context });
}

function runFallbackPath(input: {
  hostSupport: HostSupport;
  auth: BootstrapAuth;
  manifest: { protocolVersion: string; requiredCapabilities?: string[] } | null;
}): BootstrapEvaluation {
  const { hostSupport, auth, manifest } = input;
  if (auth.state === "locked") return buildResult("OK", "ACCOUNT_LOCKED", "auth-resolution");
  if (auth.state === "reauth-required") return buildResult("OK", "REAUTH_REQUIRED", "auth-resolution");
  const supportedCapabilities = hostSupport.supportedCapabilities;
  const supportedSet = new Set(supportedCapabilities);
  const manifestRequired = manifest?.requiredCapabilities ?? [];
  const missingCapabilities = manifestRequired.filter((capability) => !supportedSet.has(capability));
  if (missingCapabilities.length > 0) {
    return buildResult(
      "MISSING_REQUIRED_CAPABILITY",
      "MANIFEST_CAPABILITY_REJECTED",
      "manifest-validation",
      { missingCapabilities },
    );
  }
  return buildResult("OK", "READY", "ready", { context: resolveContext(auth) });
}

function resolveContext(auth: BootstrapAuth): { user: { id: string; name: string; roles: string[] } } {
  if (auth.state === "authenticated") {
    return {
      user: {
        id: typeof auth.principal?.id === "string" ? auth.principal.id : "",
        name: typeof auth.principal?.name === "string" ? auth.principal.name : "",
        roles: Array.isArray(auth.principal?.roles)
          ? auth.principal.roles.filter((role): role is string => typeof role === "string")
          : [],
      },
    };
  }
  // anonymous sentinel: satisfies $context.user minimal shape, is not an
  // identity and must never act as a server authorization subject.
  return { user: ANONYMOUS_USER };
}

/** SHA-256 of raw bytes as lowercase hex (Web Crypto, browser + node). */
export async function sha256Hex(bytes: Uint8Array): Promise<string> {
  const digest = await crypto.subtle.digest("SHA-256", bytes as unknown as BufferSource);
  return [...new Uint8Array(digest)].map((byte) => byte.toString(16).padStart(2, "0")).join("");
}
