/**
 * Production host boot orchestration (ADR-0035 §2.3 stage order).
 *
 * Stage order is mandatory: availability-gate terminals must render WITHOUT
 * fetching the manifest; auth-resolution terminals likewise. Each decision
 * reuses the fixture-pinned `evaluateBootstrap`, so the vendored upstream
 * host-bootstrap suite covers every stage this orchestrator executes.
 */

import registryJson from "@schemas/capability-registry.json";

import type { AppManifest } from "@/protocol/app-manifest";
import {
  BOOTSTRAP_VERSION,
  discoverBootstrapDocument,
  evaluateBootstrap,
  sha256Hex,
  type BootstrapAuth,
  type BootstrapEvaluation,
  type HostSupport,
} from "@/host/bootstrap";
import { mapBootstrapResult, nextFailureId, type HostFailure } from "@/host/failure";
import { captureReturnIntent } from "@/host/return-intent";

/** Session adapter state (ADR-0035 D4): normalized by AuthContext. */
export type SessionAdapterState = "loading" | "authenticated" | "unauthenticated" | "reauth-required" | "locked";

/** Maps the session adapter state to the bootstrap normalized auth input (D4). */
export function adapterAuthFor(
  status: SessionAdapterState,
  user: { id: string; name?: string } | null,
): BootstrapAuth {
  if (status === "authenticated" && user !== null) {
    return {
      state: "authenticated",
      principal: { id: user.id, name: user.name ?? "", roles: [] },
      provenance: "host-session-adapter",
    };
  }
  if (status === "reauth-required") {
    // The session adapter has a credential that no longer authenticates:
    // reauth-required terminal — never anonymous, never a stale principal.
    return { state: "reauth-required" };
  }
  if (status === "locked") {
    // GOAL-004 S4-6: the account-lock terminal. Never anonymous, never a
    // stale principal.
    return { state: "locked" };
  }
  return { state: "anonymous" };
}

export interface HostBootState {
  evaluation: BootstrapEvaluation;
  failure: HostFailure | null;
  manifest: AppManifest | null;
}

export interface HostBootInput {
  documentResult: Awaited<ReturnType<typeof discoverBootstrapDocument>>;
  auth: BootstrapAuth;
  manifestLoader: () => Promise<{ manifest: AppManifest; bytes: Uint8Array }>;
  registry?: unknown;
}

const HOST_SUPPORT: HostSupport = {
  supportedBootstrapVersions: [BOOTSTRAP_VERSION],
  supportedCapabilities: [
    "app.manifest",
    "app.navigation",
    "host.bootstrap",
    "host.failure-recovery",
    "host.conformance-claim",
  ],
};

/** Builds the terminal HostFailure for a bootstrap evaluation result. */
function terminalFailure(
  evaluation: BootstrapEvaluation,
  document: { availability?: { mode?: string; retryAfterSeconds?: number; messageKey?: string } } | null,
): HostFailure {
  const mapped = mapBootstrapResult({
    result: evaluation.result,
    fetchClassification: evaluation.fetchClassification,
  });
  const scope = mapped?.scope ?? "bootstrap";
  const kind = mapped?.kind ?? "protocol-rejected";
  const hostCode = mapped?.hostCode ?? "HOST_PROTOCOL_REJECTED";
  const availability = document?.availability;
  // Maintenance may carry a countdown (after) — its retry always rebuilds the
  // whole application instance; absent retryAfterSeconds means manual only.
  const retry =
    kind === "maintenance" && availability?.retryAfterSeconds !== undefined
      ? { mode: "after" as const, afterSeconds: availability.retryAfterSeconds }
      : kind === "maintenance"
        ? { mode: "manual" as const }
        : kind === "rate-limited" || kind === "timeout" || kind === "offline" || kind === "unavailable"
          ? { mode: "manual" as const }
          : kind === "forbidden" || kind === "protocol-rejected"
            ? { mode: "none" as const }
            : undefined;
  const failure: HostFailure = {
    failureVersion: "1.0",
    failureId: nextFailureId(),
    scope,
    kind,
    hostCode,
    ...(retry === undefined ? {} : { retry }),
    message: {
      messageKey: availability?.messageKey ?? kindMessageKey(kind),
    },
    diagnostics: {
      phase: evaluation.phase,
      ...(evaluation.code !== "OK" ? { protocolCode: evaluation.code } : {}),
    },
    recoveryActions: recoveryActionsFor(kind, evaluation),
  };
  return failure;
}

function kindMessageKey(kind: string): string {
  return `hostFailure.${kind.replace(/-([a-z])/g, (_m, letter: string) => letter.toUpperCase())}`;
}

function recoveryActionsFor(
  kind: string,
  evaluation: BootstrapEvaluation,
): HostFailure["recoveryActions"] {
  switch (kind) {
    case "maintenance":
      return [{ type: "retry" as const }];
    case "upgrade-required":
      return [];
    case "reauth-required":
      return [{ type: "reauth" as const }];
    case "account-locked":
      return [{ type: "home" as const }];
    case "rate-limited":
      return [{ type: "retry" as const }];
    case "timeout":
    case "offline":
      return [{ type: "retry" as const }];
    case "unavailable":
      return [{ type: "retry" as const }];
    case "protocol-rejected":
      return [];
    default:
      return evaluation.fetchClassification !== null ? [{ type: "retry" as const }] : [];
  }
}

/**
 * Executes the deterministic bootstrap lifecycle in production stage order.
 *
 * 1. discovery result (already resolved by the caller);
 * 2. bootstrap-validation + availability-gate + auth-resolution (no manifest
 *    fetch happens before terminals are decided);
 * 3. manifest load + integrity check (declared sha256);
 * 4. manifest capability narrowing (degraded);
 * 5. READY / READY_DEGRADED.
 */
export async function bootHost(input: HostBootInput): Promise<HostBootState> {
  const { documentResult, auth, manifestLoader } = input;
  const registry = input.registry ?? registryJson;

  if (documentResult.status === "failed") {
    const evaluation: BootstrapEvaluation = {
      code: "BOOTSTRAP_DOCUMENT_FAILED",
      result: "BOOTSTRAP_DOCUMENT_FAILED",
      phase: "bootstrap-discovery",
      fetchClassification: documentResult.classification,
      missingCapabilities: [],
      effectiveCapabilities: null,
      context: null,
    };
    return { evaluation, failure: terminalFailure(evaluation, null), manifest: null };
  }

  const fetchStatus = documentResult.status === "ok" ? "ok" : "not-provided";

  // Stages 2–4 with no manifest: any non-READY outcome is a pre-load terminal.
  // The integrity stage (6) is fed a synthetic pass here — its real check runs
  // below against the fetched manifest bytes; this evaluation only decides the
  // pre-load gates (validation, availability, auth).
  const declaredSha256 = documentResult.document?.manifest.sha256;
  const preLoad = evaluateBootstrap({
    document: documentResult.document,
    fetch: { status: fetchStatus },
    hostSupport: HOST_SUPPORT,
    auth,
    manifest: null,
    integrity:
      declaredSha256 !== undefined
        ? { declaredSha256, computedSha256: declaredSha256 }
        : null,
    capabilityRegistry: registry as never,
  });
  if (preLoad.result !== "READY" && preLoad.result !== "READY_DEGRADED") {
    return {
      evaluation: preLoad,
      failure: terminalFailure(preLoad, documentResult.document),
      manifest: null,
    };
  }

  // Stage 5–6: manifest load + integrity. Manifest load failures keep
  // ADR-0025/09 semantics (MANIFEST_LOAD_FAILED → existing ManifestFailure
  // surface); bootstrap must not rename them.
  let loaded: { manifest: AppManifest; bytes: Uint8Array };
  try {
    loaded = await manifestLoader();
  } catch (error) {
    throw error;
  }

  if (declaredSha256 !== undefined) {
    const computedSha256 = await sha256Hex(loaded.bytes);
    if (computedSha256 !== declaredSha256) {
      const evaluation: BootstrapEvaluation = {
        code: "MANIFEST_INTEGRITY_FAILED",
        result: "MANIFEST_INTEGRITY_FAILED",
        phase: "manifest-integrity",
        fetchClassification: null,
        missingCapabilities: [],
        effectiveCapabilities: null,
        context: null,
      };
      return { evaluation, failure: terminalFailure(evaluation, null), manifest: null };
    }
  }

  // Stages 7–9: capability narrowing + ready.
  const evaluation = evaluateBootstrap({
    document: documentResult.document,
    fetch: { status: fetchStatus },
    hostSupport: HOST_SUPPORT,
    auth,
    manifest: {
      protocolVersion: loaded.manifest.protocolVersion,
      requiredCapabilities: loaded.manifest.requiredCapabilities,
    },
    integrity:
      declaredSha256 !== undefined
        ? { declaredSha256, computedSha256: declaredSha256 }
        : null,
    capabilityRegistry: registry as never,
  });

  if (evaluation.result === "MANIFEST_CAPABILITY_REJECTED") {
    return { evaluation, failure: terminalFailure(evaluation, null), manifest: null };
  }

  return { evaluation, failure: null, manifest: loaded.manifest };
}

/** Login/session gate: locked and reauth-required must never reach ready. */
export function isBootTerminal(state: HostBootState): boolean {
  return state.failure !== null;
}

/**
 * Reauth-required terminal for the post-boot session-loss path (ADR-0035 D7):
 * the same closed failure result the boot orchestrator produces when the
 * adapter reports reauth-required before manifest-load.
 */
export function reauthFailure(): HostFailure {
  const mapped = mapBootstrapResult({ result: "REAUTH_REQUIRED", fetchClassification: null });
  const scope = mapped?.scope ?? "auth";
  const kind = mapped?.kind ?? "reauth-required";
  const hostCode = mapped?.hostCode ?? "HOST_REAUTH_REQUIRED";
  return {
    failureVersion: "1.0",
    failureId: nextFailureId(),
    scope,
    kind,
    hostCode,
    message: { messageKey: kindMessageKey(kind) },
    diagnostics: { phase: "auth-resolution" },
    recoveryActions: recoveryActionsFor(kind, {
      code: "REAUTH_REQUIRED",
      result: "REAUTH_REQUIRED",
      phase: "auth-resolution",
      fetchClassification: null,
      missingCapabilities: [],
      effectiveCapabilities: null,
      context: null,
    }),
  };
}

/**
 * Account-lock terminal (GOAL-004 S4-6, ADR-0035 D7 / ADR-0036 D6): the
 * closed failure result for the locked adapter state. account-locked allows
 * home/support only — no reauth, no retry loop.
 */
export function lockedFailure(): HostFailure {
  const mapped = mapBootstrapResult({ result: "ACCOUNT_LOCKED", fetchClassification: null });
  const scope = mapped?.scope ?? "auth";
  const kind = mapped?.kind ?? "account-locked";
  const hostCode = mapped?.hostCode ?? "HOST_ACCOUNT_LOCKED";
  return {
    failureVersion: "1.0",
    failureId: nextFailureId(),
    scope,
    kind,
    hostCode,
    message: { messageKey: kindMessageKey(kind) },
    diagnostics: { phase: "auth-resolution" },
    recoveryActions: recoveryActionsFor(kind, {
      code: "ACCOUNT_LOCKED",
      result: "ACCOUNT_LOCKED",
      phase: "auth-resolution",
      fetchClassification: null,
      missingCapabilities: [],
      effectiveCapabilities: null,
      context: null,
    }),
  };
}

/** Executes a recovery action; bootstrap retry always rebuilds the instance. */
export function executeBootRecovery(action: { type: string; url?: string }): void {
  switch (action.type) {
    case "retry":
    case "reload":
      window.location.reload();
      break;
    case "reauth":
      // Recoverable auth return intent (ADR-0036 D6): capture the current
      // in-app location before leaving for the login surface so a successful
      // login restores it.
      captureReturnIntent();
      window.location.href = "/login";
      break;
    case "home":
      window.location.href = "/";
      break;
    case "back":
      window.history.back();
      break;
    case "support":
      if (action.url) {
        const link = document.createElement("a");
        link.href = action.url;
        link.rel = "noopener noreferrer";
        link.click();
      }
      break;
    default:
      break;
  }
}
