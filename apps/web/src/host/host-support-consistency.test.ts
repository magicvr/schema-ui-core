// @vitest-environment node
//
// F2 (GOAL-042 W30 · D-001): mechanical claim ↔ host-support consistency.
// host-support.json is the single source of truth for the runtime support set
// (host-support.ts) and the build-time claim (generate-claim.mjs). This test
// locks the chain:
//   1. the generated claim's support.* equals host-support.json exactly;
//   2. every claimed capability is a valid upstream capability-registry id;
//   3. every claimed capability's mandatorySuites are covered by the claim's
//      conformance.suites with result pass (the C1 suite-coverage rule made
//      explicit and independent of claim.ts's implementation).

import { readFileSync } from "node:fs";
import { createRequire } from "node:module";
import { resolve } from "node:path";

import { describe, expect, it } from "vitest";

const require = createRequire(import.meta.url);
const WEB_ROOT = resolve(require.resolve("../../package.json"), "..");

function readJson<T>(relative: string): T {
  return JSON.parse(
    readFileSync(resolve(WEB_ROOT, relative), "utf8"),
  ) as T;
}

interface HostSupport {
  supportedPageVersions: string[];
  supportedCapabilities: string[];
}

interface CapabilityRegistry {
  capabilities: Record<
    string,
    { sinceProtocolVersion?: string; mandatorySuites?: string[] }
  >;
}

interface ClaimFixture {
  support: { pageVersions: string[]; manifestVersions: string[]; capabilities: string[] };
  conformance: { suites: Array<{ suiteId: string; result: string }> };
}

const hostSupport = readJson<HostSupport>("src/host/host-support.json");
const registry = readJson<CapabilityRegistry>("../../docs/schemas/capability-registry.json");
const claim = readJson<ClaimFixture>("public/protocol/conformance-claim.json");

function sameSet(left: string[], right: string[]): boolean {
  const a = [...new Set(left)].sort();
  const b = [...new Set(right)].sort();
  return a.length === b.length && a.every((v, i) => v === b[i]);
}

describe("F2 · claim ↔ host-support mechanical consistency (D-001)", () => {
  it("host-support.json declares a non-trivial support set", () => {
    expect(hostSupport.supportedPageVersions.length).toBeGreaterThanOrEqual(3);
    expect(hostSupport.supportedCapabilities.length).toBe(19);
  });

  it("claim support.capabilities equals host-support.json (single source)", () => {
    expect(sameSet(claim.support.capabilities, hostSupport.supportedCapabilities)).toBe(true);
  });

  it("claim support.pageVersions equals host-support.json (single source)", () => {
    expect(sameSet(claim.support.pageVersions, hostSupport.supportedPageVersions)).toBe(true);
  });

  it("every claimed capability is a valid upstream capability-registry id", () => {
    const missing = claim.support.capabilities.filter((c) => !(c in registry.capabilities));
    expect(missing).toEqual([]);
  });

  it("every claimed capability's mandatorySuites are covered by claim suites as pass", () => {
    const suiteResult = new Map(claim.conformance.suites.map((s) => [s.suiteId, s.result]));
    const uncovered: string[] = [];
    for (const capability of claim.support.capabilities) {
      for (const suiteId of registry.capabilities[capability]?.mandatorySuites ?? []) {
        if (suiteResult.get(suiteId) !== "pass") {
          uncovered.push(`${capability} -> ${suiteId}`);
        }
      }
    }
    expect(uncovered).toEqual([]);
  });
});
