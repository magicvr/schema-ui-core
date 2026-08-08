/**
 * R5 stage 3 — structure + behavior fixture execution against vendored
 * schema-ui-docs@2.7.0 artifacts (I-PROTO-004 = vendor).
 */
import { createHash } from "node:crypto";
import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

import { applyComponentFormat } from "./component-format";
import { serializeQuery } from "./query-serialize";
import { resolveStaticData } from "./static-data";
import { runRequestLifecycle } from "./request-lifecycle";
import { applyRuntimeDefaults } from "./runtime-defaults";
import { mapResponse } from "./response-mapping";
import { runSearchTable } from "./search-table";
import { runTableSort } from "./table-sort";
import { negotiateVersion } from "./version-negotiate";
import { runActionOutcome } from "./actions-outcome";
import { constructRequest } from "./request-construction";
import { runReactionEngine } from "@/renderer/reaction-engine";
import { runUploadBatch } from "./upload-orchestration";
import {
  sampleWhitelistedPage,
  validateAgainstSchema,
} from "./schema-validate";

type JsonObject = Record<string, unknown>;

interface FixtureCase {
  id: string;
  input: JsonObject;
  expected: JsonObject;
}

interface FixtureSuite {
  fixtureVersion?: string;
  category?: string;
  cases: FixtureCase[];
}

interface ProvenanceArtifact {
  path: string;
  sha256: string;
}

interface FixtureProvenance {
  sourceRepo: string;
  sourceCommit: string;
  artifactVersion: string;
  artifacts: ProvenanceArtifact[];
}

const __dirname = dirname(fileURLToPath(import.meta.url));
const UPSTREAM = join(__dirname, "../upstream");
const SCHEMAS = join(__dirname, "../../../../../docs/schemas");

function readBytes(path: string): Buffer {
  return readFileSync(path);
}

function readJsonFile<T>(path: string): { bytes: Buffer; value: T } {
  const bytes = readBytes(path);
  return { bytes, value: JSON.parse(bytes.toString("utf8")) as T };
}

function sha256(bytes: Buffer): string {
  return createHash("sha256").update(canonicalArtifactBytes(bytes)).digest("hex");
}

function canonicalArtifactBytes(bytes: Buffer): Buffer {
  // Provenance hashes describe the upstream LF bytes; Git may check them out as CRLF.
  return Buffer.from(bytes.toString("utf8").replace(/\r\n/g, "\n"), "utf8");
}

function loadSuite(name: string): { bytes: Buffer; value: FixtureSuite } {
  const path = join(UPSTREAM, `${name}.cases.json`);
  const bytes = readBytes(path);
  const parsed = JSON.parse(bytes.toString("utf8")) as FixtureSuite | FixtureCase[];
  if (Array.isArray(parsed)) {
    return {
      bytes,
      value: { fixtureVersion: "1.0", category: name, cases: parsed },
    };
  }
  return { bytes, value: parsed };
}

const provenance = readJsonFile<FixtureProvenance>(join(UPSTREAM, "provenance.json"));

const SOURCE_COMMIT = "ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b";

function assertCoverage(
  suite: FixtureSuite,
  executedIds: string[],
  exclusions: Record<string, string>,
  label: string,
) {
  const allIds = suite.cases.map((c) => c.id);
  const excludedIds = Object.keys(exclusions);
  expect(new Set(allIds).size, `${label} unique ids`).toBe(allIds.length);
  expect(executedIds.some((id) => excludedIds.includes(id))).toBe(false);
  expect(excludedIds.every((id) => allIds.includes(id))).toBe(true);
  expect([...executedIds, ...excludedIds].sort()).toEqual([...allIds].sort());
  for (const reason of Object.values(exclusions)) {
    expect(reason.trim().length).toBeGreaterThan(0);
  }
}

describe("stage 3 · pinned vendor artifacts (I-PROTO-004)", () => {
  it("pins schemas and included fixture suites at schema-ui-docs@2.7.0", () => {
    expect(provenance.value.sourceRepo).toBe("https://github.com/magicvr/schema-ui-docs");
    expect(provenance.value.sourceCommit).toBe(SOURCE_COMMIT);
    expect(provenance.value.artifactVersion).toBe("2.7.0");

    for (const artifact of provenance.value.artifacts) {
      const isSchema = artifact.path.startsWith("docs/schemas/");
      const localPath = isSchema
        ? join(SCHEMAS, artifact.path.replace("docs/schemas/", ""))
        : join(
            UPSTREAM,
            artifact.path
              .replace("conformance/fixtures/", "")
              .replace("/cases.json", ".cases.json"),
          );
      const bytes = readBytes(localPath);
      expect(sha256(bytes), artifact.path).toBe(artifact.sha256);
    }
  });
});

describe("stage 3 · structural schema validation", () => {
  it("accepts a §5-whitelist sample page against page + node schemas", () => {
    const page = sampleWhitelistedPage();
    const pageResult = validateAgainstSchema("page", page);
    expect(pageResult.ok, JSON.stringify(pageResult.errors)).toBe(true);
    const nodeResult = validateAgainstSchema("node", page.body);
    expect(nodeResult.ok, JSON.stringify(nodeResult.errors)).toBe(true);
  });

  it("rejects a node missing type", () => {
    const result = validateAgainstSchema("node", { props: {} });
    expect(result.ok).toBe(false);
  });

  it("rejects a page missing meta.protocolVersion", () => {
    const result = validateAgainstSchema("page", {
      meta: { pageId: "x", title: "T" },
      body: { type: "text" },
    });
    expect(result.ok).toBe(false);
  });

  it("accepts a minimal legal reaction document (dependencies required)", () => {
    const result = validateAgainstSchema("reaction", {
      dependencies: [],
      when: "$context.user.roles.contains('admin')",
      fulfill: { visible: false },
    });
    expect(result.ok, JSON.stringify(result.errors)).toBe(true);
  });

  it("rejects a reaction missing required dependencies", () => {
    const result = validateAgainstSchema("reaction", {
      when: "$context.user.roles.contains('admin')",
      fulfill: { visible: false },
    });
    expect(result.ok).toBe(false);
  });

  // Protocol shape example only: /api/records is a legal action URL for schema
  // validation — not a mounted product resource (Records was retired by
  // migration 0006; no /api/records handler exists).
  it("accepts a minimal legal request action", () => {
    const result = validateAgainstSchema("action", {
      type: "request",
      method: "GET",
      url: "/api/records",
    });
    expect(result.ok, JSON.stringify(result.errors)).toBe(true);
  });

  it("accepts a minimal legal navigate action", () => {
    const result = validateAgainstSchema("action", {
      type: "navigate",
      url: "/data-table",
    });
    expect(result.ok, JSON.stringify(result.errors)).toBe(true);
  });

  it("rejects an action missing type", () => {
    const result = validateAgainstSchema("action", {
      method: "GET",
      url: "/api/records",
    });
    expect(result.ok).toBe(false);
  });

  it("rejects a request action missing url", () => {
    const result = validateAgainstSchema("action", {
      type: "request",
      method: "GET",
    });
    expect(result.ok).toBe(false);
  });
});

describe("stage 3 · component-format fixtures", () => {
  const suite = loadSuite("component-format");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "component-format");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      const result = applyComponentFormat(
        fixtureCase.input.format as string,
        fixtureCase.input.value,
      );
      expect(result).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · query-serialization fixtures", () => {
  const suite = loadSuite("query-serialization");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "query-serialization");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      const result = serializeQuery(
        fixtureCase.input.baseUrl as string,
        fixtureCase.input.sources as Array<Array<[string, unknown]>>,
      );
      expect(result).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · static-data fixtures", () => {
  const suite = loadSuite("static-data");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "static-data");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      const result = resolveStaticData({
        component: fixtureCase.input.component as string,
        data: fixtureCase.input.data as never,
        datasources: fixtureCase.input.datasources as Record<string, unknown> | undefined,
        props: fixtureCase.input.props as { valueField?: string } | undefined,
      });
      expect(result).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · request-lifecycle fixtures", () => {
  const suite = loadSuite("request-lifecycle");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "request-lifecycle");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      const result = runRequestLifecycle(
        fixtureCase.input.initialState,
        fixtureCase.input.events as never,
      );
      expect(result).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · runtime-defaults fixtures", () => {
  const suite = loadSuite("runtime-defaults");
  // upload defaults materialization is pure defaults text; suite included for D-VER.
  // D-UPLOAD domain remains excluded from product features, but this case only
  // materializes default fields on a payload shape.
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "runtime-defaults");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      expect(applyRuntimeDefaults(fixtureCase.input)).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · response-mapping fixtures", () => {
  const suite = loadSuite("response-mapping");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "response-mapping");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      expect(mapResponse(fixtureCase.input)).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · search-table fixtures", () => {
  const suite = loadSuite("search-table");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "search-table");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      expect(runSearchTable(fixtureCase.input)).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · table-sort fixtures", () => {
  const suite = loadSuite("table-sort");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "table-sort");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      expect(runTableSort(fixtureCase.input)).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · version-negotiation fixtures", () => {
  const suite = loadSuite("version-negotiation");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "version-negotiation");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      expect(negotiateVersion(fixtureCase.input)).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · actions fixtures (non-batch transport outcomes)", () => {
  const suite = loadSuite("actions");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "actions");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      expect(runActionOutcome(fixtureCase.input)).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · reactions fixtures (full $deps engine, I-PROTO-FULL-001)", () => {
  const suite = loadSuite("reactions");
  const cases = suite.value.cases;
  // Full multi-round $deps reaction engine (reaction-engine.ts) executes every
  // upstream case; no exclusions remain (I-PROTO-FULL-001 D-EXPR include).
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "reactions");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      const result = runReactionEngine(fixtureCase.input as never);
      expect(result).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · uploads fixtures (I-PROTO-FULL-001 · D-UPLOAD include)", () => {
  const suite = loadSuite("uploads");
  const cases = suite.value.cases;
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "uploads");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      const result = runUploadBatch(fixtureCase.input as never);
      expect(result).toEqual(fixtureCase.expected);
    });
  }
});

describe("stage 3 · request-construction fixtures (incl. batch, I-PROTO-FULL-001)", () => {
  const suite = loadSuite("request-construction");
  const cases = suite.value.cases;
  // Full suite executes: 64 non-batch + 11 batchRequest (ADR-0022 include).
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "request-construction");

  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => {
      expect(constructRequest(fixtureCase.input)).toEqual(fixtureCase.expected);
    });
  }
});
