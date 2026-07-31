// One-shot validator for the R6 evidence-index draft schema (I-008-004).
// Builds a planning-mode dry-run from the actually captured artifact logs,
// validates it against evidence-index.schema.json, and reports parseability
// plus file-hash integrity. This proves the candidate record shape is
// parseable and that SHA-256 can be computed over real artifacts; it does NOT
// freeze the schema as the R6 acceptance contract.
import { readFileSync, writeFileSync } from "node:fs";
import { createHash } from "node:crypto";
import { fileURLToPath } from "node:url";
import { dirname, resolve } from "node:path";

import { createRequire } from "node:module";
const require = createRequire(import.meta.url);
const { default: Ajv } = require("ajv/dist/2020.js");

const here = dirname(fileURLToPath(import.meta.url));
const base = here;
const schema = JSON.parse(readFileSync(resolve(base, "evidence-index.schema.json"), "utf8"));
const resultsDir = resolve(base, "evidence", "planning", "results");

function sha256(relPath) {
  const data = readFileSync(resolve(resultsDir, relPath));
  return createHash("sha256").update(data).digest("hex");
}

const artifacts = {
  "web-test": { file: "web-test.log", exitCode: 0, outcome: "pass" },
  "web-build": { file: "web-build.log", exitCode: 0, outcome: "pass" },
  "api-test": { file: "api-test.log", exitCode: 0, outcome: "pass" },
  "api-build": { file: "api-build.log", exitCode: 0, outcome: "pass" },
  "runtime-probes": { file: "runtime-probes.log", exitCode: 0, outcome: "pass" },
};

const dryRun = {
  schemaVersion: "0.1.0-draft",
  mode: "planning",
  goalId: "GOAL-008-r6-integration-acceptance-vp-evidence",
  workspaceId: "workspace-001-mvp-admin-foundation",
  vpId: "VP-001-mvp-admin-foundation",
  repositoryRevision: "f3e04f6bd5c1f4ba6b7b72444fd9a0a0ab52d4d5",
  worktreeState: "clean",
  startedAt: "2026-08-01T04:43:00+08:00",
  completedAt: "2026-08-01T04:47:45+08:00",
  environmentRef: "windows/amd64; node v22.17.0; npm 10.9.2; go1.26.0; vitest 3.2.7",
  rerunRules: [
    "Run each command from the recorded cwd against one repository revision.",
    "Record the environment and worktree state before execution.",
    "Persist result artifacts before calculating SHA-256; missing artifacts remain explicit.",
    "Do not combine results from different revisions or environments.",
  ],
  results: [
    {
      id: "web-test",
      claimRefs: ["I-008-002", "VP-001.exit-1", "C-001"],
      command: "npm test",
      cwd: "apps/web",
      exitCode: artifacts["web-test"].exitCode,
      outcome: artifacts["web-test"].outcome,
      artifact: "evidence/planning/results/web-test.log",
      artifactStatus: "verified",
      sha256: sha256("web-test.log"),
    },
    {
      id: "web-build",
      claimRefs: ["I-008-002", "VP-001.exit-1", "C-001"],
      command: "npm run build",
      cwd: "apps/web",
      exitCode: artifacts["web-build"].exitCode,
      outcome: artifacts["web-build"].outcome,
      artifact: "evidence/planning/results/web-build.log",
      artifactStatus: "verified",
      sha256: sha256("web-build.log"),
    },
    {
      id: "api-test",
      claimRefs: ["I-008-002", "VP-001.exit-1", "C-001"],
      command: "go test ./...",
      cwd: "apps/api",
      exitCode: artifacts["api-test"].exitCode,
      outcome: artifacts["api-test"].outcome,
      artifact: "evidence/planning/results/api-test.log",
      artifactStatus: "verified",
      sha256: sha256("api-test.log"),
    },
    {
      id: "api-build",
      claimRefs: ["I-008-002", "VP-001.exit-1", "C-001"],
      command: "go build ./...",
      cwd: "apps/api",
      exitCode: artifacts["api-build"].exitCode,
      outcome: artifacts["api-build"].outcome,
      artifact: "evidence/planning/results/api-build.log",
      artifactStatus: "verified",
      sha256: sha256("api-build.log"),
    },
    {
      id: "runtime-probes",
      claimRefs: ["I-008-002", "I-008-003", "VP-001.exit-1", "VP-001.exit-3", "C-002"],
      command: "dual-service startup + HTTP probes (see log)",
      cwd: "repo-root",
      exitCode: artifacts["runtime-probes"].exitCode,
      outcome: artifacts["runtime-probes"].outcome,
      artifact: "evidence/planning/results/runtime-probes.log",
      artifactStatus: "verified",
      sha256: sha256("runtime-probes.log"),
    },
  ],
  exclusions: [
    {
      id: "clean-install-not-executed-locally",
      scope: "I-008-002",
      reason: "Local run reused the existing dependency tree; clean npm ci / fresh checkout is covered by the CI workflow being added for I-008-005.",
      reviewTrigger: "CI run on a fresh Linux checkout must produce equivalent install/build/test results.",
    },
  ],
  residuals: [
    {
      id: "platform-matrix-undecided",
      scope: "I-008-005",
      reason: "Windows local evidence captured; Linux/CI equivalence and browser E2E are being added as a minimal CI + browser matrix per user decision.",
      reviewTrigger: "CI workflow + browser runner must execute before I-008-005 is closed.",
    },
    {
      id: "cross-layer-oracle-undecided",
      scope: "I-008-003",
      reason: "Runtime probe confirms the forward path (GET /api/accounts/me returns dev session with admin+editor roles); denial-path visibility oracle at the API-to-Web/Renderer layer is being defined.",
      reviewTrigger: "Freeze the oracle before stage 2 account-permission execution.",
    },
  ],
  overallOutcome: "blocked",
};

const ajv = new Ajv();
// Slim ajv build has no bundled formats; register the one the draft schema
// references so it validates (and rejects non-date-time strings).
ajv.addFormat("date-time", {
  validate: (value) => typeof value === "string" && !Number.isNaN(Date.parse(value)),
});
const validate = ajv.compile(schema);
const ok = validate(dryRun);

// Persist the validated dry-run so the planning artifact is reproducible and
// itself hashable. It is deliberately NOT the frozen R6 acceptance index.
if (ok) {
  writeFileSync(resolve(base, "evidence", "planning", "evidence-index.dry-run.json"), JSON.stringify(dryRun, null, 2) + "\n");
}

if (!ok) {
  console.error("SCHEMA_VALIDATION_FAILED");
  console.error(JSON.stringify(validate.errors, null, 2));
  process.exit(1);
}

console.log("SCHEMA_VALIDATION_OK");
console.log(`results=${dryRun.results.length} artifacts_hashed=${dryRun.results.filter((r) => r.sha256).length}`);
for (const r of dryRun.results) {
  console.log(`${r.id}: outcome=${r.outcome} artifactStatus=${r.artifactStatus} sha256=${r.sha256.slice(0, 12)}…`);
}
