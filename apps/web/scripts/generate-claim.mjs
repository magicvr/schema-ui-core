#!/usr/bin/env node
/**
 * Build-time conformance claim generator (ADR-0037 / spec 10 §4.1).
 *
 * Emits into apps/web/public/protocol/:
 *   conformance-local-report.json   — evidence report (kind: local-report)
 *   conformance-claim.json          — the static claim (closed C0 shape)
 *   conformance-claim.json.sha256   — canonical digest (D1a) of the claim
 *
 * Bindings: schema-ui-protocol v2.8.0 (formal tag v2.8.0, commit 521cff8) —
 *   fixtureSha256          = release-check fixture-tree digest of the release
 *   protocolContentSha256  = protocol artifact contentDigest of the release
 * （2026-08-13 身份纠偏：此前 593f625/40690917… 为上游 H4 预备身份，见
 *   上游审计 0080 V379；正式 tag v2.8.0 @ 521cff8，content 4fae4605…。）
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

// Formal 2.8.0 release bindings (tag v2.8.0, commit 521cff8; upstream audit
// 0080 V379). 593f625/40690917… was the H4 preparatory identity, not the tag.
const UPSTREAM_SOURCE_COMMIT = "521cff8";
const UPSTREAM_FIXTURE_SHA256 =
  "7aacf1332ec66a16db8c79c5f3af37d241bd69b88103e503fe4d91984dd138a2";
const UPSTREAM_PROTOCOL_CONTENT_SHA256 =
  "4fae46058d01bb62d8ff5a17b35f57021a417302c9d8b932916e17ab8acf3c30";

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
    artifactVersion: "2.8.0",
    fixtureSha256: UPSTREAM_FIXTURE_SHA256,
    protocolContentSha256: UPSTREAM_PROTOCOL_CONTENT_SHA256,
  },
  residuals: [
    // GOAL-004 S4-4 纠错（2026-08-13）：页面协议 2.7 mandatory 的 multi-round
    // $deps reactions 已于 e18edce 实现（stage3 reactions 套件零排除）；本 claim
    // pageVersions ["2.7"] 条目不再视为候选绑定。
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
