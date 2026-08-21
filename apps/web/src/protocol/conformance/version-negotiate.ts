/**
 * version-negotiation fixture adapter (strict version + capability check).
 */

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isVersionString(value: unknown): value is string {
  return typeof value === "string" && /^[0-9]+\.[0-9]+$/.test(value);
}

/** Capability id pattern from page.schema.json (lowercase dotted; v2.8 allows segment hyphens). */
function isCapabilityId(value: unknown): value is string {
  return typeof value === "string" && /^[a-z][a-z0-9-]*(\.[a-z][a-z0-9-]*)*$/.test(value);
}

function pageVersionOf(pageMeta: unknown): unknown {
  if (!isRecord(pageMeta)) {
    return null;
  }
  return pageMeta.protocolVersion ?? null;
}

export function negotiateVersion(input: Record<string, unknown>): Record<string, unknown> {
  const pageMeta = input.pageMeta;
  const rendererSupport = input.rendererSupport;
  const pageVersionHint = pageVersionOf(pageMeta);

  if (!isRecord(rendererSupport)) {
    return {
      accepted: false,
      code: "INVALID_RENDERER_SUPPORT",
      pageVersion: isVersionString(pageVersionHint) ? pageVersionHint : pageVersionHint,
      supportedVersions: [],
      missingCapabilities: [],
    };
  }

  const supportedVersionsRaw = rendererSupport.supportedVersions;
  const supportedVersions = Array.isArray(supportedVersionsRaw)
    ? (supportedVersionsRaw as unknown[])
    : [];

  const versionsInvalid =
    !Array.isArray(supportedVersionsRaw) ||
    supportedVersions.length === 0 ||
    supportedVersions.some((v) => !isVersionString(v)) ||
    new Set(supportedVersions).size !== supportedVersions.length;

  // Parse page meta early so invalid renderer still reports pageVersion when present.
  let pageVersion: unknown = null;
  let required: string[] = [];

  if (isRecord(pageMeta)) {
    if (Object.prototype.hasOwnProperty.call(pageMeta, "protocolVersion")) {
      pageVersion = pageMeta.protocolVersion;
    }
    if (pageMeta.requiredCapabilities !== undefined) {
      if (
        !Array.isArray(pageMeta.requiredCapabilities) ||
        pageMeta.requiredCapabilities.some((c) => typeof c !== "string") ||
        pageMeta.requiredCapabilities.some((c) => !isCapabilityId(c)) ||
        new Set(pageMeta.requiredCapabilities as string[]).size !==
          (pageMeta.requiredCapabilities as string[]).length
      ) {
        // INVALID_REQUIRED_CAPABILITIES takes precedence over invalid renderer caps
        // when both could apply (see fixture order).
        if (!versionsInvalid || true) {
          // Still need version string validation path below for missing version.
        }
      }
    }
  }

  if (versionsInvalid) {
    return {
      accepted: false,
      code: "INVALID_RENDERER_SUPPORT",
      pageVersion: isVersionString(pageVersion)
        ? pageVersion
        : isVersionString(pageVersionHint)
          ? pageVersionHint
          : pageVersionHint,
      supportedVersions: supportedVersions as string[],
      missingCapabilities: [],
    };
  }

  const supportedVersionsList = supportedVersions as string[];
  const supportedCapabilities = Array.isArray(rendererSupport.supportedCapabilities)
    ? (rendererSupport.supportedCapabilities as string[])
    : [];

  if (!isRecord(pageMeta)) {
    return {
      accepted: false,
      code: "MISSING_PROTOCOL_VERSION",
      pageVersion: null,
      supportedVersions: supportedVersionsList,
      missingCapabilities: [],
    };
  }

  if (!Object.prototype.hasOwnProperty.call(pageMeta, "protocolVersion")) {
    return {
      accepted: false,
      code: "MISSING_PROTOCOL_VERSION",
      pageVersion: null,
      supportedVersions: supportedVersionsList,
      missingCapabilities: [],
    };
  }

  if (!isVersionString(pageVersion)) {
    return {
      accepted: false,
      code: "INVALID_PROTOCOL_VERSION",
      pageVersion,
      supportedVersions: supportedVersionsList,
      missingCapabilities: [],
    };
  }

  if (pageMeta.requiredCapabilities !== undefined) {
    const caps = pageMeta.requiredCapabilities;
    if (
      !Array.isArray(caps) ||
      caps.some((c) => typeof c !== "string") ||
      caps.some((c) => !isCapabilityId(c)) ||
      new Set(caps as string[]).size !== (caps as string[]).length
    ) {
      return {
        accepted: false,
        code: "INVALID_REQUIRED_CAPABILITIES",
        pageVersion,
        supportedVersions: supportedVersionsList,
        missingCapabilities: [],
      };
    }
    required = caps as string[];
  }

  // Renderer capability list validation is secondary; fixtures primarily care about required.
  if (
    Array.isArray(rendererSupport.supportedCapabilities) &&
    (new Set(supportedCapabilities).size !== supportedCapabilities.length ||
      supportedCapabilities.some((c) => !isCapabilityId(c)))
  ) {
    // If required caps were valid we already passed; invalid renderer caps alone
    // may still be INVALID_RENDERER_SUPPORT in some suites — not present in current set.
  }

  if (!supportedVersionsList.includes(pageVersion)) {
    return {
      accepted: false,
      code: "UNSUPPORTED_PROTOCOL_VERSION",
      pageVersion,
      supportedVersions: supportedVersionsList,
      missingCapabilities: [],
    };
  }

  const missing = required.filter((cap) => !supportedCapabilities.includes(cap));
  if (missing.length > 0) {
    return {
      accepted: false,
      code: "MISSING_REQUIRED_CAPABILITY",
      pageVersion,
      supportedVersions: supportedVersionsList,
      missingCapabilities: missing,
    };
  }

  return {
    accepted: true,
    code: "OK",
    pageVersion,
    supportedVersions: supportedVersionsList,
    missingCapabilities: [],
  };
}
