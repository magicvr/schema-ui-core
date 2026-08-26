/**
 * Host conformance claim (ADR-0037 / spec 10 §4, `host.conformance-claim`).
 *
 * Canonical serialization and digest are byte-identical with the upstream
 * JS/Python references (D1a); `validateClaim` enforces the §4.8 check order
 * against the pinned `capability-registry.json`.
 */

import registryJson from "@schemas/capability-registry.json";

const VERSION_PATTERN = /^(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)$/;
const CAPABILITY_PATTERN = /^[a-z][a-z0-9-]*(?:\.[a-z][a-z0-9-]*)*$/;
const SUITE_ID_PATTERN = /^[a-z0-9-]+$/;

export interface ClaimHost {
  hostId: string;
  hostVersion: string;
  buildId: string;
}

export interface ClaimSuite {
  suiteId: string;
  suiteVersion: string;
  result: "pass";
}

export interface ClaimEvidence {
  kind: "ci-artifact" | "signed-attestation" | "local-report";
  subjectBuildId: string;
  uri: string;
  sha256: string;
}

export interface Claim {
  claimVersion: "1.0";
  host: ClaimHost;
  protocolArtifact: { artifactVersion: string; contentSha256: string };
  support: { pageVersions: string[]; manifestVersions: string[]; capabilities: string[] };
  conformance: { fixtureVersion: "1.0"; fixtureSha256: string; suites: ClaimSuite[] };
  evidence: ClaimEvidence[];
}

export type ClaimCode =
  | "CLAIM_OK"
  | "INVALID_CLAIM"
  | "UNKNOWN_CLAIM_VERSION"
  | "UNKNOWN_CLAIM_CAPABILITY"
  | "INCOMPLETE_CAPABILITY_DEPENDENCY"
  | "CLAIM_ARTIFACT_MISMATCH"
  | "CLAIM_FIXTURE_MISMATCH"
  | "CLAIM_SUITE_INCOMPLETE"
  | "CLAIM_EVIDENCE_BUILD_MISMATCH"
  | "CLAIM_EVIDENCE_UNVERIFIABLE";

export const CLAIM_VERSION = "1.0" as const;

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isStringArray(value: unknown): value is string[] {
  return Array.isArray(value) && value.every((item) => typeof item === "string");
}

function isObjectArray(value: unknown): value is Record<string, unknown>[] {
  return Array.isArray(value) && value.every((item) => isRecord(item));
}

function compareVersion(left: string, right: string): number {
  const [leftMajor, leftMinor] = left.split(".").map(Number);
  const [rightMajor, rightMinor] = right.split(".").map(Number);
  if (leftMajor !== rightMajor) return leftMajor < rightMajor ? -1 : 1;
  if (leftMinor !== rightMinor) return leftMinor < rightMinor ? -1 : 1;
  return 0;
}

/**
 * Canonical serialization (D1a / §4.3): object keys byte-ascending, string
 * arrays byte-ascending, object arrays sorted by canonical bytes, minimal
 * RFC 8259 escaping, no whitespace. Number support is restricted to
 * integers (valid claims contain no numeric fields).
 */
export function canonicalize(value: unknown): string {
  if (isRecord(value)) {
    const parts: string[] = [];
    for (const key of Object.keys(value).sort()) {
      parts.push(`${JSON.stringify(key)}:${canonicalize(value[key])}`);
    }
    return `{${parts.join(",")}}`;
  }
  if (Array.isArray(value)) {
    let items = value;
    if (isStringArray(items)) {
      items = [...items].sort();
    } else if (isObjectArray(items)) {
      items = [...items].sort((left, right) => {
        const leftBytes = canonicalize(left);
        const rightBytes = canonicalize(right);
        return leftBytes < rightBytes ? -1 : leftBytes > rightBytes ? 1 : 0;
      });
    }
    return `[${items.map((item) => canonicalize(item)).join(",")}]`;
  }
  return JSON.stringify(value);
}

/** SHA-256 (lowercase hex) over the canonical UTF-8 bytes. */
export async function claimDigest(claim: unknown): Promise<string> {
  const bytes = new TextEncoder().encode(canonicalize(claim));
  const digest = await crypto.subtle.digest("SHA-256", bytes as unknown as BufferSource);
  return [...new Uint8Array(digest)].map((byte) => byte.toString(16).padStart(2, "0")).join("");
}

export interface ClaimOptions {
  registry?: unknown;
  artifactContentSha256: string;
  fixtureSha256: string;
  evidenceStates?: Array<"verifiable" | "unverifiable">;
}

/**
 * Registry validation (§4.4): closed entries, precise MAJOR.MINOR or null
 * versions, removedIn strictly after deprecatedSince, existing dependency
 * targets and an acyclic dependency graph.
 */
export function validateRegistry(registry: unknown): { valid: boolean; errors: string[] } {
  const errors: string[] = [];
  if (!isRecord(registry) || !isRecord(registry.capabilities)) {
    return { valid: false, errors: ["registry.capabilities must be an object"] };
  }
  const capabilities = registry.capabilities;
  const ids = Object.keys(capabilities);

  for (const id of ids) {
    if (!CAPABILITY_PATTERN.test(id)) {
      errors.push(`invalid capabilityId: ${id}`);
      continue;
    }
    const entry = capabilities[id];
    if (!isRecord(entry)) {
      errors.push(`${id}: entry must be an object`);
      continue;
    }
    const allowedKeys = new Set([
      "sinceProtocolVersion", "dependsOn", "mandatorySuites", "deprecatedSince", "removedIn",
    ]);
    for (const key of Object.keys(entry)) {
      if (!allowedKeys.has(key)) errors.push(`${id}: unknown key ${key}`);
    }
    if (typeof entry.sinceProtocolVersion !== "string"
      || !VERSION_PATTERN.test(entry.sinceProtocolVersion)) {
      errors.push(`${id}: sinceProtocolVersion must be a precise MAJOR.MINOR`);
    }
    for (const versionKey of ["deprecatedSince", "removedIn"] as const) {
      const value = entry[versionKey];
      if (value !== null && (typeof value !== "string" || !VERSION_PATTERN.test(value))) {
        errors.push(`${id}: ${versionKey} must be MAJOR.MINOR or null`);
      }
    }
    if (typeof entry.removedIn === "string") {
      if (entry.deprecatedSince === null || entry.deprecatedSince === undefined) {
        errors.push(`${id}: removedIn requires deprecatedSince`);
      } else if (typeof entry.deprecatedSince === "string"
        && compareVersion(entry.removedIn, entry.deprecatedSince) <= 0) {
        errors.push(`${id}: removedIn must be strictly later than deprecatedSince`);
      }
    }
    if (!isStringArray(entry.dependsOn) || new Set(entry.dependsOn).size !== entry.dependsOn.length) {
      errors.push(`${id}: dependsOn must be a unique string array`);
    } else {
      for (const dependency of entry.dependsOn) {
        if (!Object.hasOwn(capabilities, dependency)) {
          errors.push(`${id}: dependsOn references unregistered capability ${dependency}`);
        }
      }
    }
    if (!isStringArray(entry.mandatorySuites) || entry.mandatorySuites.length === 0
      || entry.mandatorySuites.some((suite) => !SUITE_ID_PATTERN.test(suite))
      || new Set(entry.mandatorySuites).size !== entry.mandatorySuites.length) {
      errors.push(`${id}: mandatorySuites must be a non-empty unique suite-id array`);
    }
  }

  // Dependency graph must be acyclic.
  const visiting = new Set<string>();
  const visited = new Set<string>();
  const visit = (id: string, chain: string[]): void => {
    if (visited.has(id)) return;
    if (visiting.has(id)) {
      errors.push(`dependency cycle detected: ${[...chain, id].join(" -> ")}`);
      return;
    }
    visiting.add(id);
    const entry = capabilities[id];
    if (isRecord(entry) && Array.isArray(entry.dependsOn)) {
      for (const dependency of entry.dependsOn) {
        if (Object.hasOwn(capabilities, dependency)) {
          visit(dependency, [...chain, id]);
        }
      }
    }
    visiting.delete(id);
    visited.add(id);
  };
  for (const id of ids) visit(id, []);

  return { valid: errors.length === 0, errors };
}

/**
 * Claim validation in the §4.8 order. The registry defaults to the vendored
 * `capability-registry.json` (pinned upstream artifact).
 */
export function validateClaim(claim: Claim, options: ClaimOptions): { code: ClaimCode } {
  const registry = isRecord(options.registry) && isRecord(options.registry.capabilities)
    ? options.registry
    : isRecord(registryJson) && isRecord((registryJson as Record<string, unknown>).capabilities)
      ? (registryJson as unknown as Record<string, unknown>)
      : null;
  const invalid = (code: ClaimCode) => ({ code });

  if (!isRecord(claim)) return invalid("INVALID_CLAIM");
  if (!isRecord(claim.host) || !isRecord(claim.protocolArtifact)
    || !isRecord(claim.support) || !isRecord(claim.conformance)
    || !Array.isArray(claim.evidence)) {
    return invalid("INVALID_CLAIM");
  }
  const { pageVersions, manifestVersions, capabilities } = claim.support;
  for (const list of [pageVersions, manifestVersions]) {
    if (!isStringArray(list) || list.length === 0
      || list.some((version) => !VERSION_PATTERN.test(version))
      || new Set(list).size !== list.length) {
      return invalid("INVALID_CLAIM");
    }
  }
  if (!isStringArray(capabilities) || capabilities.length === 0
    || capabilities.some((id) => !CAPABILITY_PATTERN.test(id))
    || new Set(capabilities).size !== capabilities.length) {
    return invalid("INVALID_CLAIM");
  }
  const sha256Pattern = /^[0-9a-f]{64}$/;
  if (typeof claim.conformance.fixtureSha256 !== "string"
    || !sha256Pattern.test(claim.conformance.fixtureSha256)
    || typeof claim.protocolArtifact.contentSha256 !== "string"
    || !sha256Pattern.test(claim.protocolArtifact.contentSha256)) {
    return invalid("INVALID_CLAIM");
  }
  if (!Array.isArray(claim.conformance.suites) || claim.conformance.suites.length === 0) {
    return invalid("INVALID_CLAIM");
  }
  if (claim.evidence.length === 0) return invalid("INVALID_CLAIM");
  if (typeof claim.host.buildId !== "string" || claim.host.buildId.length === 0) {
    return invalid("INVALID_CLAIM");
  }

  if (claim.claimVersion !== CLAIM_VERSION) return invalid("UNKNOWN_CLAIM_VERSION");
  if (registry === null) return invalid("INVALID_CLAIM");

  const registryCapabilities = registry.capabilities as Record<string, Record<string, unknown>>;
  const listedVersions = [...pageVersions, ...manifestVersions];
  for (const capability of capabilities) {
    const entry = registryCapabilities[capability];
    if (entry === undefined) return invalid("UNKNOWN_CLAIM_CAPABILITY");
    if (typeof entry.removedIn === "string") {
      for (const version of listedVersions) {
        if (compareVersion(entry.removedIn, version) <= 0) {
          return invalid("UNKNOWN_CLAIM_CAPABILITY");
        }
      }
    }
  }

  const listed = new Set(capabilities);
  for (const capability of capabilities) {
    // W13 F-016 (GOAL-013 A-001): this dependency walk previously re-pushed
    // every dependency's own dependsOn with no visited set — a cycle in the
    // registry graph looped forever (and shared DAG deps blew up
    // exponentially). Each dependency is now expanded at most once per walk;
    // the INCOMPLETE check still fires on every first encounter.
    const pending = [...((registryCapabilities[capability].dependsOn as string[]) ?? [])];
    const seen = new Set<string>([capability]);
    while (pending.length > 0) {
      const dependency = pending.pop()!;
      if (!listed.has(dependency)) return invalid("INCOMPLETE_CAPABILITY_DEPENDENCY");
      if (seen.has(dependency)) continue;
      seen.add(dependency);
      const dependencyEntry = registryCapabilities[dependency];
      if (dependencyEntry !== undefined && Array.isArray(dependencyEntry.dependsOn)) {
        pending.push(...(dependencyEntry.dependsOn as string[]));
      }
    }
  }

  if (claim.protocolArtifact.contentSha256 !== options.artifactContentSha256) {
    return invalid("CLAIM_ARTIFACT_MISMATCH");
  }
  if (claim.conformance.fixtureVersion !== "1.0"
    || claim.conformance.fixtureSha256 !== options.fixtureSha256) {
    return invalid("CLAIM_FIXTURE_MISMATCH");
  }

  const suiteResults = new Map(claim.conformance.suites.map((suite) => [suite.suiteId, suite]));
  for (const capability of capabilities) {
    for (const suiteId of (registryCapabilities[capability].mandatorySuites as string[]) ?? []) {
      const suite = suiteResults.get(suiteId);
      if (suite === undefined || suite.result !== "pass") return invalid("CLAIM_SUITE_INCOMPLETE");
    }
  }

  for (const evidence of claim.evidence) {
    if (!isRecord(evidence) || evidence.subjectBuildId !== claim.host.buildId) {
      return invalid("CLAIM_EVIDENCE_BUILD_MISMATCH");
    }
  }

  const evidenceStates = options.evidenceStates ?? [];
  for (let index = 0; index < claim.evidence.length; index += 1) {
    const state = index < evidenceStates.length ? evidenceStates[index] : "verifiable";
    if (state === "unverifiable") return invalid("CLAIM_EVIDENCE_UNVERIFIABLE");
  }

  return invalid("CLAIM_OK");
}
