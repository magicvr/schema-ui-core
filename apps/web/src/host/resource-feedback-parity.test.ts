/**
 * GOAL-043 · Host terminal failure vs ordinary resource feedback: direct
 * parity.
 *
 * WHY THIS EXISTS
 * ---------------
 * R4 (VP-037) unified resource feedback and fixed the transport/offline/timeout
 * classification, but the two surfaces that classify the SAME conditions were
 * only ever tested apart: `HostFailureScreen` maps a host `kind` to its copy,
 * and `feedback-policy` maps a status/code to a user-facing category. Its audit
 * left the direct comparison as a bounded recommended (GOAL-005 A-002 F-002 /
 * R5-I-005): redundant copy could drift, or one surface could lose explicit
 * handling for a shared condition, without any test noticing.
 *
 * This test compares the two tables against each other: for every condition both
 * surfaces can see, the host kind must exist with its own copy (not the generic
 * fallback), the resource policy must classify it to the same concept, and both
 * keys must exist in both locales.
 */

import { readFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

import { HOST_FAILURE_GENERIC_KEY, HOST_FAILURE_MESSAGE_KEYS } from "@/app/HostFailureScreen";
import { classifyFeedbackFailure, type FeedbackFailureCategory } from "@/renderer/feedback-policy";

const HERE = dirname(fileURLToPath(import.meta.url));
/** This test lives at src/host/, so the package root is two levels up. */
const WEB_ROOT = resolve(HERE, "../..");
const CATALOGS = {
  "zh-CN": join(WEB_ROOT, "src/i18n/messages/zh-CN.json"),
  "en-US": join(WEB_ROOT, "src/i18n/messages/en-US.json"),
} as const;

/**
 * Conditions both surfaces handle. `resource` is the signal an ordinary list or
 * action sees; `hostKind` is the terminal kind the host raises for it.
 */
const SHARED_CONDITIONS: Array<{
  condition: FeedbackFailureCategory;
  resource: { status?: number; code?: string };
  hostKind: string;
}> = [
  { condition: "maintenance", resource: { code: "SERVICE_MAINTENANCE" }, hostKind: "maintenance" },
  { condition: "timeout", resource: { code: "REQUEST_TIMEOUT" }, hostKind: "timeout" },
  { condition: "offline", resource: { code: "REQUEST_FAILED" }, hostKind: "offline" },
  { condition: "rate-limited", resource: { status: 429 }, hostKind: "rate-limited" },
  { condition: "authentication", resource: { status: 401 }, hostKind: "authentication-required" },
  { condition: "forbidden", resource: { status: 403 }, hostKind: "forbidden" },
  { condition: "not-found", resource: { status: 404 }, hostKind: "not-found" },
  { condition: "unavailable", resource: { status: 503 }, hostKind: "unavailable" },
];

const GENERIC_RESOURCE_KEY = "feedback.operationFailed";

describe("GOAL-043 · host failure and resource feedback parity", () => {
  it("classifies every shared condition the same way on both surfaces", () => {
    const mismatches = SHARED_CONDITIONS.map((row) => {
      const policy = classifyFeedbackFailure(row.resource);
      return policy.category === row.condition
        ? null
        : `${row.condition}: resource policy said ${policy.category} for ${JSON.stringify(row.resource)}`;
    }).filter((entry): entry is string => entry !== null);
    expect(mismatches, "resource feedback policy disagrees with the shared table").toEqual([]);

    const missingKinds = SHARED_CONDITIONS.filter(
      (row) => HOST_FAILURE_MESSAGE_KEYS[row.hostKind] === undefined,
    ).map((row) => row.hostKind);
    expect(missingKinds, "host surface lost explicit handling for a shared condition").toEqual([]);
  });

  it("gives each shared condition its own copy on both surfaces", () => {
    const hostGeneric = SHARED_CONDITIONS.filter(
      (row) => HOST_FAILURE_MESSAGE_KEYS[row.hostKind] === HOST_FAILURE_GENERIC_KEY,
    ).map((row) => row.hostKind);
    expect(hostGeneric, "host surface must not fall back to generic copy").toEqual([]);

    const resourceGeneric = SHARED_CONDITIONS.filter(
      (row) => classifyFeedbackFailure(row.resource).messageKey === GENERIC_RESOURCE_KEY,
    ).map((row) => row.condition);
    expect(resourceGeneric, "resource policy must not fall back to generic copy").toEqual([]);

    // A condition must not resolve to the same key on both surfaces by accident:
    // the two namespaces are distinct so each surface keeps its own copy.
    for (const row of SHARED_CONDITIONS) {
      const hostKey = HOST_FAILURE_MESSAGE_KEYS[row.hostKind];
      const resourceKey = classifyFeedbackFailure(row.resource).messageKey;
      expect(hostKey.startsWith("hostFailure."), `${hostKey} must live in the host namespace`).toBe(
        true,
      );
      expect(
        resourceKey.startsWith("feedback."),
        `${resourceKey} must live in the feedback namespace`,
      ).toBe(true);
      expect(hostKey).not.toBe(resourceKey);
    }
  });

  it("keeps every referenced key present in both locales", () => {
    const catalogs = Object.entries(CATALOGS).map(([locale, path]) => ({
      locale,
      keys: JSON.parse(readFileSync(path, "utf8")) as Record<string, string>,
    }));
    const missing: string[] = [];
    for (const { locale, keys } of catalogs) {
      for (const row of SHARED_CONDITIONS) {
        for (const key of [
          HOST_FAILURE_MESSAGE_KEYS[row.hostKind],
          classifyFeedbackFailure(row.resource).messageKey,
        ]) {
          if (typeof keys[key] !== "string" || keys[key] === "") {
            missing.push(`${locale}: ${key}`);
          }
        }
      }
    }
    expect(missing, "a shared condition is missing copy in a locale").toEqual([]);
  });
});
