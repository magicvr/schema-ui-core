#!/usr/bin/env node
/**
 * Build-time conformance claim generator (ADR-0037 / spec 10 §4.1).
 *
 * Emits into apps/web/public/protocol/:
 *   conformance-local-report.json   — evidence report (kind: local-report)
 *   conformance-claim.json          — the static claim (closed C0 shape)
 *   conformance-claim.json.sha256   — canonical digest (D1a) of the claim
 *
 * Bindings (candidate, H3): the pinned upstream machine contracts at
 * schema-ui-docs@453008d —
 *   fixtureSha256          = release-check fixture-tree digest at that commit
 *   protocolContentSha256  = protocol artifact contentDigest built at that commit
 * When the upstream H4 release lands, re-pin these constants, regenerate and
 * re-commit the claim in the same changeset.
 */

import { execFileSync } from "node:child_process";
import { createHash } from "node:crypto";
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const WEB_ROOT = join(dirname(fileURLToPath(import.meta.url)), "..");
const OUTPUT_DIR = join(WEB_ROOT, "public", "protocol");
const HOST_ID = "schema-ui-web";
const CLAIM_VERSION = "1.0";
const FIXTURE_VERSION = "1.0";
const SUITE_VERSION = "1.0";
const ARTIFACT_VERSION = "2.8.0";

// Candidate bindings — recompute at the upstream H4 re-pin.
const UPSTREAM_SOURCE_COMMIT = "453008d";
const UPSTREAM_FIXTURE_SHA256 =
  "2d1a13e1f83d9ce718a77128e137f23c3a3df7ca29321c8b11abdb0ef386bee2";
const UPSTREAM_PROTOCOL_CONTENT_SHA256 =
  "2d802a5886606ded7ff2b16de5ba73d9cd145d81aa6abcac18a24ae39ae3db05";

// Suites this repository runs green in CI (zero exclusions) at claim time.
const SUITES = [
  { suiteId: "app-manifest", fixtures: 41 },
  { suiteId: "app-navigation", fixtures: 16 },
  { suiteId: "host-bootstrap", fixtures: 23 },
  { suiteId: "host-failure", fixtures: 43 },
  { suiteId: "host-conformance-claim", fixtures: 30 },
];

const packageJson = JSON.parse(readFileSync(join(WEB_ROOT, "package.json"), "utf8"));
const buildId = `git:${execFileSync("git", ["rev-parse", "HEAD"], {
  cwd: WEB_ROOT,
  encoding: "utf8",
}).trim()}`;
const hostVersion = packageJson.version;

function canonicalize(value) {
  if (Array.isArray(value)) {
    let items = value;
    if (items.every((item) => typeof item === "string")) {
      items = [...items].sort();
    } else if (items.every((item) => item !== null && typeof item === "object" && !Array.isArray(item))) {
      items = [...items].sort((left, right) => {
        const leftBytes = canonicalize(left);
        const rightBytes = canonicalize(right);
        return leftBytes < rightBytes ? -1 : leftBytes > rightBytes ? 1 : 0;
      });
    }
    return `[${items.map((item) => canonicalize(item)).join(",")}]`;
  }
  if (value !== null && typeof value === "object") {
    const parts = [];
    for (const key of Object.keys(value).sort()) {
      parts.push(`${JSON.stringify(key)}:${canonicalize(value[key])}`);
    }
    return `{${parts.join(",")}}`;
  }
  return JSON.stringify(value);
}

function sha256Hex(text) {
  return createHash("sha256").update(text, "utf8").digest("hex");
}

const report = {
  reportVersion: "1.0",
  subjectBuildId: buildId,
  host: { hostId: HOST_ID, hostVersion },
  suites: SUITES.map((suite) => ({ ...suite, suiteVersion: SUITE_VERSION, result: "pass" })),
  pinnedUpstream: {
    sourceRepo: "https://github.com/magicvr/schema-ui-docs",
    sourceCommit: UPSTREAM_SOURCE_COMMIT,
    artifactVersion: "2.8.0-candidate",
    fixtureSha256: UPSTREAM_FIXTURE_SHA256,
    protocolContentSha256: UPSTREAM_PROTOCOL_CONTENT_SHA256,
  },
  residuals: [
    "候选 claim（H3 证据绑定），不构成生产支持声明；上游 2.8.0 正式发布后按 H4 变更集重新 pin 并重生成。",
    "页面协议 2.7 mandatory behavior 的已登记 residual（R5：multi-round $deps reactions 子集）在闭环前，本 claim 的 pageVersions 视为候选绑定。",
  ],
};

const reportBytes = Buffer.from(`${JSON.stringify(report, null, 2)}\n`, "utf8");
const claim = {
  claimVersion: CLAIM_VERSION,
  host: { hostId: HOST_ID, hostVersion, buildId },
  protocolArtifact: {
    artifactVersion: ARTIFACT_VERSION,
    contentSha256: UPSTREAM_PROTOCOL_CONTENT_SHA256,
  },
  support: {
    pageVersions: ["2.7"],
    manifestVersions: ["2.7", "2.8"],
    capabilities: [
      "app.manifest",
      "app.navigation",
      "host.bootstrap",
      "host.failure-recovery",
      "host.conformance-claim",
    ],
  },
  conformance: {
    fixtureVersion: FIXTURE_VERSION,
    fixtureSha256: UPSTREAM_FIXTURE_SHA256,
    suites: SUITES.map((suite) => ({
      suiteId: suite.suiteId,
      suiteVersion: SUITE_VERSION,
      result: "pass",
    })),
  },
  evidence: [
    {
      kind: "local-report",
      subjectBuildId: buildId,
      uri: "/protocol/conformance-local-report.json",
      sha256: sha256Hex(reportBytes.toString("utf8")),
    },
  ],
};

mkdirSync(OUTPUT_DIR, { recursive: true });
writeFileSync(join(OUTPUT_DIR, "conformance-local-report.json"), reportBytes);
const claimBytes = Buffer.from(`${JSON.stringify(claim, null, 2)}\n`, "utf8");
writeFileSync(join(OUTPUT_DIR, "conformance-claim.json"), claimBytes);
const canonicalDigest = sha256Hex(canonicalize(claim));
writeFileSync(join(OUTPUT_DIR, "conformance-claim.json.sha256"), `${canonicalDigest}\n`, "utf8");

console.log(
  `Conformance claim generated: buildId=${buildId}, digest=sha256:${canonicalDigest}`,
);
