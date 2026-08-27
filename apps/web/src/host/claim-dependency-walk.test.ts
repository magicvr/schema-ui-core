import { describe, expect, it } from "vitest";

import { CLAIM_VERSION, validateClaim, type Claim } from "@/host/claim";

/**
 * W13 F-016 (GOAL-013 A-001) regression locks for the validateClaim
 * dependency walk: a cyclic registry dependsOn graph used to loop forever
 * (no visited set); each dependency must now expand at most once while the
 * INCOMPLETE_CAPABILITY_DEPENDENCY check stays exact.
 */

const sha256 = "a".repeat(64);

function claimWithCapabilities(capabilities: string[]): Claim {
  return {
    claimVersion: CLAIM_VERSION,
    host: { hostId: "host", hostVersion: "1.0", buildId: "build-1" },
    protocolArtifact: { artifactVersion: "1.0", contentSha256: sha256 },
    support: { pageVersions: ["1.0"], manifestVersions: ["1.0"], capabilities },
    conformance: {
      fixtureVersion: "1.0",
      fixtureSha256: sha256,
      suites: [{ suiteId: "core-suite", suiteVersion: "1.0", result: "pass" }],
    },
    evidence: [
      { kind: "local-report", subjectBuildId: "build-1", uri: "file:///report.json", sha256 },
    ],
  };
}

describe("validateClaim dependency walk (W13 F-016)", () => {
  it("terminates on a cyclic registry graph and validates the listed set", () => {
    const registry = {
      capabilities: {
        "cap.a": { dependsOn: ["cap.b"], mandatorySuites: [] },
        "cap.b": { dependsOn: ["cap.a"], mandatorySuites: [] },
      },
    };
    expect(
      validateClaim(claimWithCapabilities(["cap.a", "cap.b"]), {
        registry,
        artifactContentSha256: sha256,
        fixtureSha256: sha256,
      }),
    ).toEqual({ code: "CLAIM_OK" });
  });

  it("still rejects an unlisted dependency", () => {
    const registry = {
      capabilities: {
        "cap.a": { dependsOn: ["cap.missing"], mandatorySuites: [] },
      },
    };
    expect(
      validateClaim(claimWithCapabilities(["cap.a"]), {
        registry,
        artifactContentSha256: sha256,
        fixtureSha256: sha256,
      }),
    ).toEqual({ code: "INCOMPLETE_CAPABILITY_DEPENDENCY" });
  });
});
