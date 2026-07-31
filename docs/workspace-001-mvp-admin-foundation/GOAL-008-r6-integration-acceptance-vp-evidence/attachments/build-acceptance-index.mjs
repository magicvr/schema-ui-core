// R6 stage-2 acceptance evidence index builder + validator (I-008-004).
// Builds the formal acceptance index (mode: acceptance) from the actually
// captured result artifacts, computes SHA-256 over them, and validates against
// evidence-index.schema.json. This is the stage-2 persisted evidence index;
// planning artifacts remain under evidence/planning/.
import { readFileSync, writeFileSync } from "node:fs";
import { createHash } from "node:crypto";
import { fileURLToPath } from "node:url";
import { dirname, resolve } from "node:path";
import { createRequire } from "node:module";

const require = createRequire(import.meta.url);
const { default: Ajv } = require("ajv/dist/2020.js");

const here = dirname(fileURLToPath(import.meta.url)); // attachments/
const schema = JSON.parse(readFileSync(resolve(here, "evidence-index.schema.json"), "utf8"));
const resultsDir = resolve(here, "evidence", "acceptance", "results");

function sha256(relPath) {
  return createHash("sha256").update(readFileSync(resolve(resultsDir, relPath))).digest("hex");
}

const artifacts = [
  { id: "web-test", claimRefs: ["C-001", "C-005", "C-007", "VP-001.exit-1", "VP-001.exit-2", "VP-001.exit-3"], command: "npm test", cwd: "apps/web", exitCode: 0, outcome: "pass", artifact: "evidence/acceptance/results/web-test.log" },
  { id: "web-build", claimRefs: ["C-001", "VP-001.exit-1"], command: "npm run build", cwd: "apps/web", exitCode: 0, outcome: "pass", artifact: "evidence/acceptance/results/web-build.log" },
  { id: "api-test", claimRefs: ["C-002", "C-005", "VP-001.exit-1"], command: "go test ./...", cwd: "apps/api", exitCode: 0, outcome: "pass", artifact: "evidence/acceptance/results/api-test.log" },
  { id: "api-build", claimRefs: ["C-002", "VP-001.exit-1"], command: "go build ./...", cwd: "apps/api", exitCode: 0, outcome: "pass", artifact: "evidence/acceptance/results/api-build.log" },
  { id: "runtime-probes", claimRefs: ["C-003", "C-005", "VP-001.exit-1", "VP-001.exit-3"], command: "dual-service startup + HTTP probes (see log)", cwd: "repo-root", exitCode: 0, outcome: "pass", artifact: "evidence/acceptance/results/runtime-probes.log" },
  { id: "browser-e2e", claimRefs: ["C-004", "C-005", "VP-001.exit-1", "VP-001.exit-3"], command: "npm run test:e2e (Playwright Chromium)", cwd: "apps/web", exitCode: 0, outcome: "pass", artifact: "evidence/acceptance/results/browser-e2e.json" },
  { id: "browser-screenshot", claimRefs: ["C-004", "VP-001.exit-1"], command: "Playwright screenshot (r6-overview.png)", cwd: "apps/web", exitCode: 0, outcome: "pass", artifact: "evidence/acceptance/results/r6-overview.png" },
];

const index = {
  schemaVersion: "0.1.0-draft",
  mode: "acceptance",
  goalId: "GOAL-008-r6-integration-acceptance-vp-evidence",
  workspaceId: "workspace-001-mvp-admin-foundation",
  vpId: "VP-001-mvp-admin-foundation",
  repositoryRevision: "a941bedb1fc2cd4859a408df50653e867da35ff2",
  worktreeState: "clean",
  startedAt: "2026-08-01T05:54:00+08:00",
  completedAt: "2026-08-01T05:58:00+08:00",
  environmentRef: "windows/amd64; node v22.17.0; npm 10.9.2; go1.26.0; vitest 3.2.7; playwright 1.62.1",
  rerunRules: [
    "Run each command from the recorded cwd against one repository revision.",
    "Record the environment and worktree state before execution.",
    "Persist result artifacts before calculating SHA-256; missing artifacts remain explicit.",
    "Do not combine results from different revisions or environments.",
  ],
  results: artifacts.map((a) => ({
    ...a,
    artifactStatus: "verified",
    sha256: sha256(a.artifact.replace("evidence/acceptance/results/", "")),
  })),
  exclusions: [
    {
      id: "reactions-multiround-excluded",
      scope: "D-EXPR upstream reactions suite (16/16)",
      reason: "Upstream multi-round $deps writer semantics are not in the MVP $context subset (Root D-008); MVP formal entry is reactions.test.ts + /form-with-reactions.",
      reviewTrigger: "Re-evaluate if MVP $context subset expands to multi-round $deps.",
    },
    {
      id: "request-construction-batch-excluded",
      scope: "D-DATA request-construction batch cases (11)",
      reason: "Batch request construction excluded from MVP (Root D-010 Q1=否); non-batch 64/64 executed within npm test stage3.",
      reviewTrigger: "Re-evaluate when batch semantics are scoped.",
    },
    {
      id: "uploads-domain-excluded",
      scope: "D-UPLOAD",
      reason: "Upload domain excluded from MVP coverage (I-PROTO-001 v0.1.3).",
      reviewTrigger: "Only if upload becomes an MVP scope.",
    },
    {
      id: "local-not-clean-install",
      scope: "I-008-002 local run",
      reason: "This local acceptance run reused the existing dependency tree; clean install (npm ci) is evidenced by the GitHub Actions run (30667596846) on Linux.",
      reviewTrigger: "Re-run local clean install if dependency tree changes materially.",
    },
    {
      id: "denial-path-browser-not-asserted",
      scope: "C-006 browser-level denial",
      reason: "The real manifest has no permission-gated navigation item; shell denial is asserted at renderer/component level (App.integration.test, D-PERM fixtures) per account-permission-oracle D-1..D-6, not in a browser E2E.",
      reviewTrigger: "If a permission-gated nav item is added to the real manifest.",
    },
  ],
  residuals: [],
  overallOutcome: "pass",
};

const ajv = new Ajv();
ajv.addFormat("date-time", {
  validate: (value) => typeof value === "string" && !Number.isNaN(Date.parse(value)),
});
const validate = ajv.compile(schema);
const ok = validate(index);
if (!ok) {
  console.error("ACCEPTANCE_INDEX_VALIDATION_FAILED");
  console.error(JSON.stringify(validate.errors, null, 2));
  process.exit(1);
}

const outPath = resolve(here, "evidence", "acceptance", "evidence-index.json");
writeFileSync(outPath, JSON.stringify(index, null, 2) + "\n");
console.log("ACCEPTANCE_INDEX_VALIDATION_OK");
console.log(`results=${index.results.length} hashed=${index.results.filter((r) => r.sha256).length} overallOutcome=${index.overallOutcome}`);
for (const r of index.results) {
  console.log(`${r.id}: ${r.outcome} sha256=${r.sha256.slice(0, 12)}…`);
}
console.log(`wrote ${outPath}`);
