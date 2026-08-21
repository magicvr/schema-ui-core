import { createHash } from "node:crypto";
import { readFileSync } from "node:fs";

import Ajv from "ajv";
import { describe, expect, it } from "vitest";

import claimSchema from "@schemas/host-conformance-claim.schema.json";
import {
  canonicalize,
  claimDigest,
  validateClaim,
} from "@/host/claim";

/**
 * Verifies the build-generated claim artifact (ADR-0037 D4/D5 binding):
 * the claim must be structurally valid (C0), pass the C1 check order against
 * the pinned registry and digests, bind evidence bytes by sha256, and match
 * its recorded canonical digest. Any mismatch fails the build gate.
 */

const CLAIM_PATH = new URL("../../public/protocol/conformance-claim.json", import.meta.url);
const REPORT_PATH = new URL("../../public/protocol/conformance-local-report.json", import.meta.url);
const DIGEST_PATH = new URL("../../public/protocol/conformance-claim.json.sha256", import.meta.url);

function readCanonicalJson<T>(path: URL): { bytes: Buffer; value: T } {
  const bytes = readFileSync(path);
  const canonical = Buffer.from(bytes.toString("utf8").replace(/\r\n/g, "\n"), "utf8");
  return { bytes: canonical, value: JSON.parse(canonical.toString("utf8")) as T };
}

interface ClaimFixture {
  claimVersion: string;
  host: { hostId: string; hostVersion: string; buildId: string };
  protocolArtifact: { artifactVersion: string; contentSha256: string };
  support: { pageVersions: string[]; manifestVersions: string[]; capabilities: string[] };
  conformance: {
    fixtureVersion: string;
    fixtureSha256: string;
    suites: Array<{ suiteId: string; suiteVersion: string; result: string }>;
  };
  evidence: Array<{ kind: string; subjectBuildId: string; uri: string; sha256: string }>;
}

interface ReportFixture {
  reportVersion: string;
  subjectBuildId: string;
  pinnedUpstream: { sourceCommit: string; fixtureSha256: string; protocolContentSha256: string };
}

describe("build-generated conformance claim artifact", () => {
  const claimArtifact = readCanonicalJson<ClaimFixture>(CLAIM_PATH);
  const reportArtifact = readCanonicalJson<ReportFixture>(REPORT_PATH);

  it("is structurally valid against the pinned C0 schema", () => {
    const validate = new Ajv({ allErrors: true, strict: false }).compile(claimSchema as object);
    expect(validate(claimArtifact.value), JSON.stringify(validate.errors)).toBe(true);
  });

  it("passes the C1 check order against the pinned registry and digests", () => {
    const report = reportArtifact.value;
    const result = validateClaim(claimArtifact.value as never, {
      artifactContentSha256: report.pinnedUpstream.protocolContentSha256,
      fixtureSha256: report.pinnedUpstream.fixtureSha256,
    });
    expect(result).toEqual({ code: "CLAIM_OK" });
  });

  it("binds the evidence report bytes by sha256", () => {
    const evidence = claimArtifact.value.evidence[0];
    expect(evidence.kind).toBe("local-report");
    expect(evidence.subjectBuildId).toBe(claimArtifact.value.host.buildId);
    const digest = createHash("sha256").update(reportArtifact.bytes).digest("hex");
    expect(evidence.sha256).toBe(digest);
  });

  it("matches the recorded canonical claim digest (D1a reproducibility)", async () => {
    const recorded = readFileSync(DIGEST_PATH, "utf8").trim();
    expect(await claimDigest(claimArtifact.value)).toBe(recorded);
    // Cross-check: the canonical serialization itself is pinned by the
    // upstream claim fixtures (canonicalize cases).
    expect(typeof canonicalize(claimArtifact.value)).toBe("string");
  });

  it("binds the same build id across report, evidence and claim host", () => {
    expect(reportArtifact.value.subjectBuildId).toBe(claimArtifact.value.host.buildId);
    expect(claimArtifact.value.evidence[0].subjectBuildId).toBe(claimArtifact.value.host.buildId);
    expect(claimArtifact.value.host.buildId).toMatch(/^git:[0-9a-f]{40}$/);
  });
});
