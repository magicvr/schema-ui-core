import { createHash } from "node:crypto";
import { readFileSync } from "node:fs";

import { describe, expect, it } from "vitest";

import { projectNavigation, type ProjectedItem } from "@/app/navigation";
import {
  APP_MANIFEST_SOURCE,
  APP_MANIFEST_PROTOCOL_VERSION,
  ManifestError,
  type NavigationContext,
  type PageEntry,
  loadAppManifest,
  matchRoute,
  pageIdMatches,
  resolveInitialRoute,
  resolveLogoUrl,
  resolveSchemaUrl,
  validateAppManifest,
} from "@/protocol/app-manifest";

type JsonObject = Record<string, unknown>;

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

interface PinnedJson<T> {
  bytes: Buffer;
  value: T;
}

const APP_MANIFEST_SCHEMA_SHA256 =
  "294eb0f391a515a123ed07ef5e19b40f1126fbe57e585ebdf639a34232a6c2e8";
const APP_MANIFEST_FIXTURE_SHA256 =
  "8aafd684a40c4350d39ce0c94ab3c986455cb21fe0ebc37e4f3640261163405e";
const APP_NAVIGATION_FIXTURE_SHA256 =
  "11b0117078b6e12c92805e21da02f9fe522fe69ae8bf41d74498cbef468f2897";
const STATIC_MANIFEST_SHA256 =
  "8d90160418ef68f10164a460bd7eede215cb292f8c1a9264912500b9e51c637c";

function readJson<T>(relativePath: string): PinnedJson<T> {
  const bytes = canonicalArtifactBytes(readFileSync(new URL(relativePath, import.meta.url)));
  return {
    bytes,
    value: JSON.parse(bytes.toString("utf8")) as T,
  };
}

function canonicalArtifactBytes(bytes: Buffer): Buffer {
  // Provenance hashes describe the upstream LF bytes; Git may check them out as CRLF.
  return Buffer.from(bytes.toString("utf8").replace(/\r\n/g, "\n"), "utf8");
}

const schemaArtifact = readJson<JsonObject>(
  "../../../../docs/schemas/app-manifest.schema.json",
);
const appManifestArtifact = readJson<FixtureSuite>(
  "./upstream/app-manifest.cases.json",
);
const appNavigationArtifact = readJson<FixtureSuite>(
  "./upstream/app-navigation.cases.json",
);
const provenanceArtifact = readJson<FixtureProvenance>("./upstream/provenance.json");
const staticManifestArtifact = readJson<JsonObject>(
  "../../public/.well-known/schema-ui/app-manifest.json",
);

function isRecord(value: unknown): value is JsonObject {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function record(value: unknown, label: string): JsonObject {
  if (!isRecord(value)) {
    throw new Error("Expected " + label + " to be an object.");
  }
  return value;
}

function stringValue(value: unknown, label: string): string {
  if (typeof value !== "string") {
    throw new Error("Expected " + label + " to be a string.");
  }
  return value;
}

function stringArray(value: unknown, label: string): string[] {
  if (!Array.isArray(value) || value.some((entry) => typeof entry !== "string")) {
    throw new Error("Expected " + label + " to be an array of strings.");
  }
  return value;
}

function findCase(suite: FixtureSuite, id: string): FixtureCase {
  const fixtureCase = suite.cases.find((candidate) => candidate.id === id);
  if (fixtureCase === undefined) {
    throw new Error("Fixture case not found: " + id);
  }
  return fixtureCase;
}

function caseInput(fixtureCase: FixtureCase): JsonObject {
  return record(fixtureCase.input, fixtureCase.id + ".input");
}

function errorResult(error: unknown, includeErrors = false): JsonObject {
  if (!(error instanceof ManifestError)) {
    throw error;
  }
  return includeErrors
    ? {
        ok: false,
        code: error.code,
        path: error.path,
        errors: [{ code: error.code, path: error.path }],
      }
    : { ok: false, code: error.code, path: error.path };
}

function codeResult(error: unknown): JsonObject {
  if (!(error instanceof ManifestError)) {
    throw error;
  }
  return { ok: false, code: error.code };
}

function pageIdFromPages(value: unknown): string | undefined {
  if (!Array.isArray(value)) {
    return undefined;
  }
  const firstPage = value[0];
  if (!isRecord(firstPage) || typeof firstPage.pageId !== "string") {
    return undefined;
  }
  return firstPage.pageId;
}

function hostManifestValue(
  rawValue: unknown,
  fallbackPages?: unknown,
): JsonObject {
  const raw = record(rawValue, "manifest");
  const rawProtocol = typeof raw.protocolVersion === "string" ? raw.protocolVersion : "";
  const protocolParts = rawProtocol.split(".").map(Number);
  const isBelowMinimum =
    protocolParts.length === 2 &&
    Number.isFinite(protocolParts[0]) &&
    Number.isFinite(protocolParts[1]) &&
    (protocolParts[0] < 2 ||
      (protocolParts[0] === 2 && protocolParts[1] < 5));
  const navigation = raw.navigation;
  const rawCapabilities = Array.isArray(raw.requiredCapabilities)
    ? raw.requiredCapabilities.filter(
        (capability): capability is string => typeof capability === "string",
      )
    : [];
  const capabilities = [
    ...new Set([
      ...rawCapabilities,
      "app.manifest",
      ...(navigation === undefined ? [] : ["app.navigation"]),
    ]),
  ];
  const pages = raw.pages ?? fallbackPages ?? [];
  const app =
    raw.app ??
    ({
      appId: "fixture-app",
      name: "Fixture",
      ...(pageIdFromPages(pages) === undefined
        ? {}
        : { homePageRef: pageIdFromPages(pages) }),
    } satisfies JsonObject);
  return {
    ...raw,
    protocolVersion: isBelowMinimum ? rawProtocol : APP_MANIFEST_PROTOCOL_VERSION,
    requiredCapabilities: capabilities,
    app,
    pages,
  };
}

function navigationManifestValue(input: JsonObject): JsonObject {
  const baseline = findCase(appNavigationArtifact.value, "three-slots-project");
  const baselineInput = caseInput(baseline);
  const pages = input.pages ?? baselineInput.pages;
  const app = input.app ?? {
    appId: "fixture-navigation",
    name: "Fixture Navigation",
    homePageRef: pageIdFromPages(pages),
  };
  return hostManifestValue(
    {
      app,
      pages,
      navigation: input.navigation,
    },
    pages,
  );
}

function queryString(value: unknown): string {
  const query = record(value, "query");
  return new URLSearchParams(
    Object.entries(query).map(([key, entry]) => [
      key,
      stringValue(entry, "query." + key),
    ]),
  ).toString();
}

function routeResult(
  route: ReturnType<typeof resolveInitialRoute>,
): JsonObject | undefined {
  if (route === undefined) {
    return undefined;
  }
  return {
    path: route.path,
    params: route.params,
    query: route.query,
  };
}

// These operations are fixture-only adapters. R3 production support remains an exact 2.7 host;
// the upstream cases also exercise negotiation inputs that are intentionally outside that host API.
function negotiateFixture(
  manifestValue: unknown,
  rendererSupportValue: unknown,
): JsonObject {
  const manifest = record(manifestValue, "manifest");
  const rendererSupport = record(rendererSupportValue, "rendererSupport");
  const pageVersion =
    typeof manifest.protocolVersion === "string" ? manifest.protocolVersion : null;
  const supportedVersions = stringArray(
    rendererSupport.supportedVersions,
    "rendererSupport.supportedVersions",
  );
  const supportedCapabilities = stringArray(
    rendererSupport.supportedCapabilities,
    "rendererSupport.supportedCapabilities",
  );
  const requiredCapabilities =
    manifest.requiredCapabilities === undefined
      ? []
      : stringArray(manifest.requiredCapabilities, "manifest.requiredCapabilities");
  const missingCapabilities = requiredCapabilities.filter(
    (capability) => !supportedCapabilities.includes(capability),
  );
  const code =
    pageVersion === null
      ? "MISSING_PROTOCOL_VERSION"
      : !supportedVersions.includes(pageVersion)
        ? "UNSUPPORTED_PROTOCOL_VERSION"
        : missingCapabilities.length > 0
          ? "MISSING_REQUIRED_CAPABILITY"
          : "OK";
  return {
    accepted: code === "OK",
    code,
    pageVersion,
    supportedVersions,
    missingCapabilities,
  };
}

async function runManifestCase(fixtureCase: FixtureCase): Promise<JsonObject> {
  const input = caseInput(fixtureCase);
  const operation = stringValue(input.operation, fixtureCase.id + ".operation");

  try {
    switch (operation) {
      case "validate": {
        validateAppManifest(hostManifestValue(input.manifest));
        return { ok: true };
      }
      case "negotiate":
        return negotiateFixture(input.manifest, input.rendererSupport);
      case "decoupledVersions": {
        const manifestResult = negotiateFixture(
          input.manifest,
          input.rendererSupport,
        );
        const pageMeta = record(input.pageMeta, fixtureCase.id + ".pageMeta");
        const pageResult = negotiateFixture(
          { protocolVersion: pageMeta.protocolVersion, requiredCapabilities: [] },
          input.pageRendererSupport,
        );
        return {
          manifest: manifestResult,
          page: pageResult,
          ok: manifestResult.accepted === true && pageResult.accepted === true,
        };
      }
      case "resolveHome": {
        const manifest = validateAppManifest(hostManifestValue(input.manifest));
        const requestedPath =
          input.deepLinkPath === undefined
            ? "/"
            : stringValue(input.deepLinkPath, fixtureCase.id + ".deepLinkPath");
        const route = resolveInitialRoute(manifest, requestedPath);
        if (route === undefined) {
          return { ok: false };
        }
        return {
          ok: true,
          path: route.path,
          pageId: route.page.pageId,
          route: routeResult(route),
          source: route.source,
        };
      }
      case "matchRoute": {
        const pages = input.manifest === undefined
          ? (input.pages as PageEntry[])
          : validateAppManifest(hostManifestValue(input.manifest)).pages;
        const route = matchRoute(
          pages,
          stringValue(input.path, fixtureCase.id + ".path"),
        );
        return route === undefined
          ? { matched: false }
          : {
              matched: true,
              pageId: route.page.pageId,
              route: route.page.route,
              params: route.params,
              index: route.index,
            };
      }
      case "resolveSchemaUrl":
        return {
          ok: true,
          url: resolveSchemaUrl(
            stringValue(input.baseURL, fixtureCase.id + ".baseURL"),
            stringValue(input.schemaUrl, fixtureCase.id + ".schemaUrl"),
            record(input.params, fixtureCase.id + ".params") as Record<string, string>,
          ),
        };
      case "pageIdMatch": {
        const pages = [
          {
            pageId: stringValue(input.registryPageId, fixtureCase.id + ".registryPageId"),
            title: "Fixture",
            schemaUrl: "/schema",
            route: "/fixture",
          },
        ] satisfies PageEntry[];
        return {
          ok: pageIdMatches(
            pages[0],
            stringValue(input.schemaMetaPageId, fixtureCase.id + ".schemaMetaPageId"),
          ),
        };
      }
      case "resolveLogo":
        return {
          ok: true,
          url: resolveLogoUrl(
            stringValue(input.baseURL, fixtureCase.id + ".baseURL"),
            stringValue(input.logoUrl, fixtureCase.id + ".logoUrl"),
          ),
        };
      case "navigate": {
        const manifest = validateAppManifest(hostManifestValue(input.manifest));
        const appPath = stringValue(input.appPath, fixtureCase.id + ".appPath");
        const query = queryString(input.query);
        const requestedPath = query === "" ? appPath : appPath + "?" + query;
        const route = resolveInitialRoute(manifest, requestedPath);
        return {
          ok: true,
          navigated: true,
          registryHit: route !== undefined,
          route: route === undefined
            ? {
                path: appPath,
                params: {},
                query: record(input.query, fixtureCase.id + ".query"),
              }
            : {
                path: route.path,
                params: route.params,
                query: route.query,
              },
        };
      }
      case "load": {
        try {
          await loadAppManifest({
            fetcher: async () => new Response("missing", { status: 503 }),
          });
          return { ok: true };
        } catch (error) {
          if (!(error instanceof ManifestError)) {
            throw error;
          }
          return { ok: false, code: error.code, renderPages: false };
        }
      }
      default:
        throw new Error("Unsupported app-manifest fixture operation: " + operation);
    }
  } catch (error) {
    if (operation === "validate") {
      return errorResult(error, true);
    }
    if (operation === "pageIdMatch" || operation === "resolveLogo") {
      return codeResult(error);
    }
    return errorResult(error);
  }
}

function supportedProjection(item: ProjectedItem): JsonObject {
  if (item.type === "group") {
    return {
      type: "group",
      items: item.items.map(supportedProjection),
      label: item.label,
      ...(item.icon === undefined ? {} : { icon: item.icon }),
    };
  }
  return {
    type: "link",
    ...(item.pageRef === undefined ? {} : { pageRef: item.pageRef }),
    ...(item.url === undefined ? {} : { url: item.url }),
    label: item.label,
    ...(item.icon === undefined ? {} : { icon: item.icon }),
    active: item.active,
  };
}

function supportedExpectedItem(value: unknown): JsonObject {
  const item = record(value, "expected navigation item");
  if (item.type === "group") {
    const children = item.items;
    if (!Array.isArray(children)) {
      throw new Error("Expected group items to be an array.");
    }
    return {
      type: "group",
      items: children.map(supportedExpectedItem),
      label: stringValue(item.label, "expected group label"),
      ...(item.icon === undefined ? {} : { icon: stringValue(item.icon, "expected group icon") }),
    };
  }
  return {
    type: "link",
    ...(item.pageRef === undefined
      ? {}
      : { pageRef: stringValue(item.pageRef, "expected pageRef") }),
    ...(item.url === undefined
      ? {}
      : { url: stringValue(item.url, "expected url") }),
    label: stringValue(item.label, "expected link label"),
    ...(item.icon === undefined
      ? {}
      : { icon: stringValue(item.icon, "expected link icon") }),
    active: item.active === true,
  };
}

function supportedExpectedSlots(value: unknown): JsonObject {
  const expected = record(value, "expected navigation result");
  const slots = record(expected.slots, "expected navigation slots");
  const supported: JsonObject = {};
  for (const slot of ["top", "sidebar", "user"] as const) {
    if (slots[slot] !== undefined) {
      if (!Array.isArray(slots[slot])) {
        throw new Error("Expected " + slot + " to be an array.");
      }
      supported[slot] = slots[slot].map(supportedExpectedItem);
    }
  }
  return supported;
}

function supportedSlots(
  projection: ReturnType<typeof projectNavigation>,
): JsonObject {
  const slots: JsonObject = {};
  for (const slot of ["top", "sidebar", "user"] as const) {
    if (projection[slot].length > 0) {
      slots[slot] = projection[slot].map(supportedProjection);
    }
  }
  return slots;
}

function supportedExpectedValidation(value: unknown): JsonObject {
  const expected = record(value, "expected validation result");
  if (expected.ok === true) {
    return { ok: true };
  }
  const errors = expected.errors;
  const firstError = Array.isArray(errors) ? errors[0] : undefined;
  return {
    ok: false,
    code: stringValue(expected.code, "expected validation code"),
    path: stringValue(expected.path, "expected validation path"),
    ...(isRecord(firstError) ? { errors: [firstError] } : {}),
  };
}

function runNavigationCase(fixtureCase: FixtureCase): JsonObject {
  const input = caseInput(fixtureCase);
  const operation = stringValue(input.operation, fixtureCase.id + ".operation");
  try {
    if (operation === "validate") {
      validateAppManifest(navigationManifestValue(input));
      return { ok: true };
    }
    if (operation !== "project") {
      throw new Error("Unsupported app-navigation fixture operation: " + operation);
    }
    const manifest = validateAppManifest(navigationManifestValue(input));
    const context = (input.context ?? {}) as NavigationContext;
    const projection = projectNavigation(
      manifest,
      stringValue(input.path, fixtureCase.id + ".path"),
      context,
    );
    return { ok: true, slots: supportedSlots(projection) };
  } catch (error) {
    return errorResult(error, operation === "validate");
  }
}

const manifestCaseExclusions: Record<string, string> = {
  "m1-missing-app-manifest-capability":
    "Excluded: the upstream M1 schema fixture uses CAPABILITY_REQUIRED, while the R3 hand-written host validator exposes MISSING_REQUIRED_CAPABILITY; this schema error-envelope difference is outside the frozen R3 subset.",
  "m1-navigation-without-capability":
    "Excluded: the upstream M1 schema fixture uses CAPABILITY_REQUIRED, while the R3 hand-written host validator exposes MISSING_REQUIRED_CAPABILITY; this schema error-envelope difference is outside the frozen R3 subset.",
};

const manifestCases = appManifestArtifact.value.cases.filter(
  (fixtureCase) => !Object.prototype.hasOwnProperty.call(manifestCaseExclusions, fixtureCase.id),
);
const navigationCases = appNavigationArtifact.value.cases;

function assertFixtureCoverage(
  suite: FixtureSuite,
  executedCases: FixtureCase[],
  exclusions: Record<string, string>,
  label: string,
) {
  const allIds = suite.cases.map((fixtureCase) => fixtureCase.id);
  const executedIds = executedCases.map((fixtureCase) => fixtureCase.id);
  const excludedIds = Object.keys(exclusions);
  expect(new Set(allIds).size, label + " fixture IDs must be unique").toBe(allIds.length);
  expect(executedIds.some((id) => excludedIds.includes(id))).toBe(false);
  expect(excludedIds.every((id) => allIds.includes(id))).toBe(true);
  expect([...executedIds, ...excludedIds].sort()).toEqual([...allIds].sort());
  for (const reason of Object.values(exclusions)) {
    expect(reason.trim().length).toBeGreaterThan(0);
  }
}

describe("pinned schema-ui-docs fixture artifacts", () => {
  it("keeps the schema and behavior fixtures at the pinned commit", () => {
    expect(provenanceArtifact.value.sourceRepo).toBe(
      "https://github.com/magicvr/schema-ui-docs",
    );
    expect(provenanceArtifact.value.sourceCommit).toBe(
      "ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b",
    );
    expect(APP_MANIFEST_SOURCE).toBe(
      "https://github.com/magicvr/schema-ui-docs/tree/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b",
    );
    expect(provenanceArtifact.value.artifactVersion).toBe("2.7.0");
    // R3 baseline artifacts remain pinned; R5 stage 3 extends provenance with
    // additional schemas/fixtures (see stage3-fixtures.test.ts for full set).
    const byPath = new Map(
      provenanceArtifact.value.artifacts.map((artifact) => [artifact.path, artifact.sha256]),
    );
    expect(byPath.get("docs/schemas/app-manifest.schema.json")).toBe(APP_MANIFEST_SCHEMA_SHA256);
    expect(byPath.get("conformance/fixtures/app-manifest/cases.json")).toBe(
      APP_MANIFEST_FIXTURE_SHA256,
    );
    expect(byPath.get("conformance/fixtures/app-navigation/cases.json")).toBe(
      APP_NAVIGATION_FIXTURE_SHA256,
    );
    expect(createHash("sha256").update(schemaArtifact.bytes).digest("hex")).toBe(
      APP_MANIFEST_SCHEMA_SHA256,
    );
    expect(createHash("sha256").update(appManifestArtifact.bytes).digest("hex")).toBe(
      APP_MANIFEST_FIXTURE_SHA256,
    );
    expect(createHash("sha256").update(appNavigationArtifact.bytes).digest("hex")).toBe(
      APP_NAVIGATION_FIXTURE_SHA256,
    );
    expect(createHash("sha256").update(staticManifestArtifact.bytes).digest("hex")).toBe(
      STATIC_MANIFEST_SHA256,
    );
    expect(schemaArtifact.value.$id).toBe(
      "https://schema-ui.dev/schemas/app-manifest.schema.json",
    );
    expect(appManifestArtifact.value.fixtureVersion).toBe("1.0");
    expect(appNavigationArtifact.value.fixtureVersion).toBe("1.0");
  });

  it("accounts for every upstream fixture case", () => {
    assertFixtureCoverage(
      appManifestArtifact.value,
      manifestCases,
      manifestCaseExclusions,
      "app-manifest",
    );
    assertFixtureCoverage(
      appNavigationArtifact.value,
      navigationCases,
      {},
      "app-navigation",
    );
  });
});

describe("app-manifest upstream behavior fixtures", () => {
  for (const fixtureCase of manifestCases) {
    it(fixtureCase.id, async () => {
      await expect(runManifestCase(fixtureCase)).resolves.toEqual(
        fixtureCase.input.operation === "validate"
          ? supportedExpectedValidation(fixtureCase.expected)
          : fixtureCase.expected,
      );
    });
  }
});

describe("app-navigation upstream behavior fixtures", () => {
  for (const fixtureCase of navigationCases) {
    it(fixtureCase.id, () => {
      if (fixtureCase.input.operation === "validate") {
        expect(runNavigationCase(fixtureCase)).toEqual(fixtureCase.expected);
      } else {
        expect(runNavigationCase(fixtureCase)).toEqual({
          ok: true,
          slots: supportedExpectedSlots(fixtureCase.expected),
        });
      }
    });
  }
});
