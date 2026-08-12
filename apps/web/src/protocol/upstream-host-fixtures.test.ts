import { createHash } from "node:crypto";
import { readFileSync } from "node:fs";

import { describe, expect, it } from "vitest";

import { evaluateBootstrap, type BootstrapEvaluation } from "@/host/bootstrap";
import {
  classifyHostFetch,
  mapBootstrapResult,
  validateFailure,
  validateReturnIntent,
  type HostFailure,
} from "@/host/failure";
import {
  canonicalize,
  claimDigest,
  validateClaim,
  validateRegistry,
} from "@/host/claim";

/**
 * Host/App 互操作 conformance（S4 生产实现消费上游候选 fixtures，零排除）。
 *
 * 与 R3/R5 的 fixture 消费同构：上游 `schema-ui-docs@453008d`（H2 机器契约）
 * 的 host-bootstrap / host-failure / host-conformance-claim 三个 suite 直接
 * 逐字段核对本仓生产 Host 模块（`src/host/*`）。这些模块是 main.tsx /
 * App shell 的生产代码路径，不是 fixture-only adapter；任何 skip/改写期望
 * 均不允许（上游消费规则）。
 */

const HOST_BOOTSTRAP_FIXTURE_SHA256 =
  "2ea0f81841f974596d76e041c67d83ae9544a37115d3b8fce40915e918ef03db";
const HOST_FAILURE_FIXTURE_SHA256 =
  "de950c074ba150c31ec02d07c07950b5c5004a24a5cfc17024407c3d69bdd43f";
const HOST_CLAIM_FIXTURE_SHA256 =
  "6d27589836ab05636c3f598fc64b376207d0b1a140ca661e089b6f4d2939400a";

interface FixtureCase {
  id: string;
  input: Record<string, unknown>;
  expected: Record<string, unknown> | null;
}

interface FixtureSuite {
  fixtureVersion: string;
  category: string;
  cases: FixtureCase[];
}

function readCanonicalJson<T>(relativePath: string): { bytes: Buffer; value: T } {
  const bytes = readFileSync(new URL(relativePath, import.meta.url));
  const canonical = Buffer.from(bytes.toString("utf8").replace(/\r\n/g, "\n"), "utf8");
  return { bytes: canonical, value: JSON.parse(canonical.toString("utf8")) as T };
}

function assertPin(bytes: Buffer, expectedSha256: string, label: string): void {
  expect(createHash("sha256").update(bytes).digest("hex"), label).toBe(expectedSha256);
}

const bootstrapArtifact = readCanonicalJson<FixtureSuite>(
  "./upstream/host-bootstrap.cases.json",
);
const failureArtifact = readCanonicalJson<FixtureSuite>(
  "./upstream/host-failure.cases.json",
);
const claimArtifact = readCanonicalJson<FixtureSuite>(
  "./upstream/host-conformance-claim.cases.json",
);

describe("upstream host-bootstrap fixtures (production host implementation, no exclusions)", () => {
  it("matches the pinned provenance sha256", () => {
    assertPin(bootstrapArtifact.bytes, HOST_BOOTSTRAP_FIXTURE_SHA256, "host-bootstrap fixture pin");
  });

  it.each(bootstrapArtifact.value.cases.map((fixture) => [fixture.id, fixture] as const))(
    "%s",
    (_id, fixture) => {
      const actual = evaluateBootstrap(fixture.input as never) as BootstrapEvaluation;
      expect(actual).toEqual(fixture.expected);
    },
  );
});

describe("upstream host-failure fixtures (production host implementation, no exclusions)", () => {
  it("matches the pinned provenance sha256", () => {
    assertPin(failureArtifact.bytes, HOST_FAILURE_FIXTURE_SHA256, "host-failure fixture pin");
  });

  it.each(failureArtifact.value.cases.map((fixture) => [fixture.id, fixture] as const))(
    "%s",
    (_id, fixture) => {
      const { input } = fixture;
      let actual: unknown;
      switch (input.operation) {
        case "classify":
          actual = classifyHostFetch(input as never);
          break;
        case "mapBootstrapResult":
          actual = mapBootstrapResult(input as never);
          break;
        case "validateFailure":
          actual = validateFailure(input.failure as HostFailure);
          break;
        case "validateReturnIntent":
          actual = validateReturnIntent(
            input.intent as never,
            input.options as { registeredKeys?: string[]; nowIso: string },
          );
          break;
        default:
          throw new Error(`Unknown operation: ${String(input.operation)}`);
      }
      expect(actual).toEqual(fixture.expected);
    },
  );
});

describe("upstream host-conformance-claim fixtures (production host implementation, no exclusions)", () => {
  it("matches the pinned provenance sha256", () => {
    assertPin(claimArtifact.bytes, HOST_CLAIM_FIXTURE_SHA256, "host-conformance-claim fixture pin");
  });

  it.each(claimArtifact.value.cases.map((fixture) => [fixture.id, fixture] as const))(
    "%s",
    async (_id, fixture) => {
      const { input } = fixture;
      let actual: unknown;
      switch (input.operation) {
        case "canonicalize":
          actual = { canonical: canonicalize(input.value) };
          break;
        case "claimDigest":
          actual = { digest: await claimDigest(input.claim) };
          break;
        case "validateRegistry":
          actual = validateRegistry(input.registry);
          break;
        case "validateClaim":
          actual = validateClaim(input.claim as never, input.options as never);
          break;
        default:
          throw new Error(`Unknown operation: ${String(input.operation)}`);
      }
      expect(actual).toEqual(fixture.expected);
    },
  );
});
