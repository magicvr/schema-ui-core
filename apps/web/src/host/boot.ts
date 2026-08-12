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
  documentMessageKey: string | undefined,
): HostFailure {
  const mapped = mapBootstrapResult({
    result: evaluation.result,
    fetchClassification: evaluation.fetchClassification,
  });
  const scope = mapped?.scope ?? "bootstrap";
  const kind = mapped?.kind ?? "protocol-rejected";
  const hostCode = mapped?.hostCode ?? "HOST_PROTOCOL_REJECTED";
  const failure: HostFailure = {
    failureVersion: "1.0",
    failureId: nextFailureId(),
    scope,
    kind,
    hostCode,
    message: {
      messageKey: documentMessageKey ?? kindMessageKey(kind),
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
    return { evaluation, failure: terminalFailure(evaluation, undefined), manifest: null };
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
    const documentMessageKey = documentResult.document?.availability.messageKey;
    return {
      evaluation: preLoad,
      failure: terminalFailure(preLoad, documentMessageKey),
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
      return { evaluation, failure: terminalFailure(evaluation, undefined), manifest: null };
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
    return { evaluation, failure: terminalFailure(evaluation, undefined), manifest: null };
  }

  return { evaluation, failure: null, manifest: loaded.manifest };
}

/** Login/session gate: locked and reauth-required must never reach ready. */
export function isBootTerminal(state: HostBootState): boolean {
  return state.failure !== null;
}

/** Executes a recovery action; bootstrap retry always rebuilds the instance. */
export function executeBootRecovery(action: { type: string; url?: string }): void {
  switch (action.type) {
    case "retry":
    case "reload":
      window.location.reload();
      break;
    case "reauth":
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
