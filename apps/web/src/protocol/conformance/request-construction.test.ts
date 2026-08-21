/**
 * D-001 P0 regression: recordSource.url must reject protocol-relative (`//`)
 * and absolute URLs so authFetch never attaches the Bearer token to an
 * external origin. Same-origin single-slash paths (optionally with query)
 * remain legal. The vendored upstream cases.json is pinned by sha256 and is
 * not modified; these local cases guard the hardened rule directly.
 */
import { describe, expect, it } from "vitest";

import { constructRequest } from "./request-construction";

function recordSourceInput(url: string): Record<string, unknown> {
  return {
    kind: "recordSource",
    recordSource: {
      method: "GET",
      url,
      responseMapping: { siteTitle: "siteTitle" },
    },
    route: { query: {}, params: {} },
  };
}

describe("recordSource.url protocol hardening (D-001 P0)", () => {
  it("accepts a same-origin single-slash path", () => {
    const result = constructRequest(recordSourceInput("/api/settings/default"));
    expect(result.ok).toBe(true);
  });

  it("accepts a same-origin path with a query", () => {
    const result = constructRequest(recordSourceInput("/api/orders?verbose=true"));
    expect(result.ok).toBe(true);
  });

  it("rejects a protocol-relative //host URL (Bearer leak)", () => {
    const result = constructRequest(recordSourceInput("//evil.example/api/settings"));
    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.code).toBe("INVALID_PROTOCOL_URL");
      expect(result.path).toBe("recordSource.url");
    }
  });

  it("rejects a URL whose backslash would resolve to an external host", () => {
    const result = constructRequest(recordSourceInput("/\\evil.example/api/settings"));
    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.code).toBe("INVALID_PROTOCOL_URL");
    }
  });

  it("rejects an absolute https URL", () => {
    const result = constructRequest(recordSourceInput("https://evil.example/api/settings"));
    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.code).toBe("INVALID_PROTOCOL_URL");
    }
  });
});
