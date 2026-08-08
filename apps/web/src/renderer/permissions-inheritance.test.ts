import { readFileSync } from "node:fs";

import { describe, expect, it } from "vitest";

import {
  type ExecutionRequest,
  type L2Error,
  type PermissionTarget,
  evaluatePermissionTargets,
  executeAction,
  validatePermissions,
} from "@/renderer/permissions";

type JsonObject = Record<string, unknown>;

function isRecord(value: unknown): value is JsonObject {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

interface FixtureCase {
  id: string;
  input: JsonObject;
  expected: JsonObject;
}

interface FixtureSuite {
  fixtureVersion: string;
  category: string;
  cases: FixtureCase[];
}

const CASES_SHA256 = "ac124fa1d831d0aa2544b7544b1e177c3498c8c3b36ee4d535e8c3f2f5b8849e";
const casesBytes = canonicalArtifactBytes(
  readFileSync(
    new URL("../protocol/upstream/permissions-inheritance.cases.json", import.meta.url),
  ),
);

const { createHash } = await import("node:crypto");

function canonicalArtifactBytes(bytes: Buffer): Buffer {
  // Governance hashes describe the recorded LF bytes; Git may check them out as CRLF.
  return Buffer.from(bytes.toString("utf8").replace(/\r\n/g, "\n"), "utf8");
}

describe("pinned permissions-inheritance fixture integrity", () => {
  it("matches the SHA-256 recorded in GOAL-006 D-004 and provenance.json", () => {
    expect(createHash("sha256").update(casesBytes).digest("hex")).toBe(CASES_SHA256);
  });
});

const suite = JSON.parse(casesBytes.toString("utf8")) as FixtureSuite;
expect(suite.fixtureVersion).toBe("1.0");
expect(suite.category).toBe("permissions-inheritance");
expect(suite.cases.length).toBe(17);

function validationResult(page: JsonObject): { ok: boolean; errors: L2Error[] } {
  const errors = validatePermissions(page);
  return { ok: errors.length === 0, errors };
}

function sortTargets(targets: PermissionTarget[]): PermissionTarget[] {
  return [...targets].sort((a, b) => a.targetId.localeCompare(b.targetId));
}

function normalizeTarget(target: PermissionTarget): JsonObject {
  return {
    targetId: target.targetId,
    kind: target.kind,
    key: target.key,
    cascadeApplied: target.cascadeApplied,
    cascadedBy: target.cascadedBy,
    effectivePermission: target.effectivePermission,
  };
}

function expectedTargets(expected: JsonObject): JsonObject[] {
  const targets = expected.targets;
  return Array.isArray(targets) ? (targets as JsonObject[]) : [];
}

function pageFor(fixtureCase: FixtureCase): JsonObject {
  const input = fixtureCase.input;
  const page = isRecord(input.page) ? input.page : {};
  // navigatedPage is a top-level fixture key (new permission root), not a
  // field of page; merge it so the evaluator sees the full document.
  if (input.navigatedPage !== undefined) {
    return { ...page, navigatedPage: input.navigatedPage };
  }
  return page;
}

function runValidation(fixtureCase: FixtureCase): JsonObject {
  const result = validationResult(pageFor(fixtureCase));
  return result.ok ? { valid: true, errors: [] } : { valid: false, errors: result.errors };
}

function runEvaluation(fixtureCase: FixtureCase): JsonObject {
  const input = fixtureCase.input;
  const context = (input.context ?? {}) as Record<string, unknown>;
  const targets = sortTargets(evaluatePermissionTargets(pageFor(fixtureCase), context));
  const expected = expectedTargets(fixtureCase.expected).map((entry) => entry.targetId);
  return {
    ok: true,
    targets: targets
      .filter((target) => expected.includes(target.targetId))
      .map(normalizeTarget),
  };
}

function runExecution(fixtureCase: FixtureCase): JsonObject {
  const input = fixtureCase.input;
  const context = (input.context ?? {}) as Record<string, unknown>;
  const request = (input.execution ?? {}) as ExecutionRequest;
  const result = executeAction(pageFor(fixtureCase), request, context);
  return {
    ok: true,
    execution: result,
  };
}

describe("permissions-inheritance upstream behavior fixtures", () => {
  for (const fixtureCase of suite.cases) {
    it(fixtureCase.id, () => {
      const expected = fixtureCase.expected;
      const expectedValidation = isRecord(expected.validation) ? expected.validation : {};
      const valid = expectedValidation.valid === true;
      const actual: JsonObject = {};

      actual.validation = runValidation(fixtureCase);

      if (valid) {
        actual.targets = runEvaluation(fixtureCase).targets;
        if (fixtureCase.input.execution !== undefined) {
          actual.execution = runExecution(fixtureCase).execution;
        }
      }

      if (valid) {
        expect(actual.validation).toEqual({ valid: true, errors: [] });
        expect(sortTargets(actual.targets as unknown as PermissionTarget[])).toEqual(
          sortTargets(expectedTargets(expected) as unknown as PermissionTarget[]),
        );
        if (fixtureCase.input.execution !== undefined) {
          expect(actual.execution).toEqual(expected.execution);
        }
      } else {
        expect(actual.validation).toEqual({
          valid: false,
          errors: expectedValidation.errors,
        });
      }
    });
  }
});
