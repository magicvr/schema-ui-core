/**
 * ADR-0010 query-serialization adapter (schema-ui-docs v2.7.0 fixtures).
 *
 * Merges base URL query with ordered source layers; encodes keys/values per
 * RFC 3986 (unreserved A-Za-z0-9-._~); sorts final keys by Unicode code point.
 */

export type QueryScalar = string | number | boolean | null;

export type QueryPair = [string, unknown];

export type QuerySource = QueryPair[];

export type QuerySerializeResult =
  | { ok: true; url: string }
  | {
      ok: false;
      code: "INVALID_BASE_URL_QUERY" | "INVALID_QUERY_KEY" | "INVALID_QUERY_VALUE";
    };

function isUndefinedMarker(value: unknown): boolean {
  return (
    typeof value === "object" &&
    value !== null &&
    !Array.isArray(value) &&
    Object.prototype.hasOwnProperty.call(value, "$undefined") &&
    (value as { $undefined: unknown }).$undefined === true
  );
}

function isComposite(value: unknown): boolean {
  return (
    (typeof value === "object" && value !== null && !isUndefinedMarker(value)) ||
    Array.isArray(value)
  );
}

/** RFC3986 encode: percent-encode everything except unreserved. */
export function encodeRFC3986(value: string): string {
  return encodeURIComponent(value).replace(/[!'()*]/g, (char) => {
    return "%" + char.charCodeAt(0).toString(16).toUpperCase().padStart(2, "0");
  });
}

/**
 * JCS-style number serialization used by query-serialization fixtures.
 * Matches JSON number text for typical finite values (1e+21, 1e-7, 0.000001).
 */
export function serializeQueryNumber(value: number): string {
  if (!Number.isFinite(value)) {
    throw new Error("non-finite number");
  }
  // Prefer the same text JSON would emit for these fixtures.
  if (Object.is(value, -0)) {
    return "0";
  }
  const asJson = JSON.stringify(value);
  return asJson;
}

function serializeScalar(value: QueryScalar): string {
  if (value === null) {
    return "";
  }
  if (typeof value === "boolean") {
    return value ? "true" : "false";
  }
  if (typeof value === "number") {
    return serializeQueryNumber(value);
  }
  return value;
}

function decodePercentUtf8(encoded: string): string | null {
  try {
    // Reject incomplete / non-hex percent sequences before decodeURIComponent.
    if (/%(?:$|[^0-9A-Fa-f]|[0-9A-Fa-f](?:$|[^0-9A-Fa-f]))/.test(encoded)) {
      return null;
    }
    if (/%[0-9A-Fa-f]{2}/.test(encoded)) {
      // Reject overlong / invalid UTF-8 by checking decodeURIComponent + URIError.
      const decoded = decodeURIComponent(encoded);
      // decodeURIComponent accepts %FF as a single Latin-1-ish replacement in some
      // engines? Node throws URIError for invalid UTF-8 sequences in modern V8.
      // Also reject lone high bytes that don't form valid UTF-8 when re-encoded.
      for (let i = 0; i < encoded.length; i++) {
        if (encoded[i] === "%") {
          const hex = encoded.slice(i + 1, i + 3);
          const byte = Number.parseInt(hex, 16);
          // After decodeURIComponent, if original had invalid UTF-8 multi-byte, it throws.
          // Single-byte 0xFF is invalid as UTF-8 leading in a string context when
          // percent-decoded in isolation — Node's decodeURIComponent throws on %FF.
          void byte;
        }
      }
      return decoded;
    }
    return encoded;
  } catch {
    return null;
  }
}

interface ParsedBase {
  path: string;
  pairs: Array<[string, string]>;
  fragment: string;
}

function parseBaseUrl(baseUrl: string): ParsedBase | { error: "INVALID_BASE_URL_QUERY" | "INVALID_QUERY_KEY" } {
  const hashIndex = baseUrl.indexOf("#");
  const withoutFragment = hashIndex >= 0 ? baseUrl.slice(0, hashIndex) : baseUrl;
  const fragment = hashIndex >= 0 ? baseUrl.slice(hashIndex) : "";
  const qIndex = withoutFragment.indexOf("?");
  const path = qIndex >= 0 ? withoutFragment.slice(0, qIndex) : withoutFragment;
  const query = qIndex >= 0 ? withoutFragment.slice(qIndex + 1) : "";

  if (query === "") {
    return { path, pairs: [], fragment };
  }

  const pairs: Array<[string, string]> = [];
  for (const part of query.split("&")) {
    if (part === "") {
      // empty segment like trailing & — treat as empty key?
      continue;
    }
    const eq = part.indexOf("=");
    const rawKey = eq >= 0 ? part.slice(0, eq) : part;
    const rawVal = eq >= 0 ? part.slice(eq + 1) : "";
    if (rawKey === "") {
      return { error: "INVALID_QUERY_KEY" };
    }
    // Plus is literal in base query (not space).
    const keyDecoded = decodePercentUtf8(rawKey.replace(/\+/g, "%2B"));
    const valDecoded = decodePercentUtf8(rawVal.replace(/\+/g, "%2B"));
    if (keyDecoded === null || valDecoded === null) {
      return { error: "INVALID_BASE_URL_QUERY" };
    }
    if (keyDecoded === "") {
      return { error: "INVALID_QUERY_KEY" };
    }
    pairs.push([keyDecoded, valDecoded]);
  }
  return { path, pairs, fragment };
}

export function serializeQuery(
  baseUrl: string,
  sources: QuerySource[],
): QuerySerializeResult {
  const parsed = parseBaseUrl(baseUrl);
  if ("error" in parsed) {
    return { ok: false, code: parsed.error };
  }

  // Map of key → value text (or deleted).
  const map = new Map<string, string>();
  // Apply base pairs (last wins for duplicates).
  for (const [key, value] of parsed.pairs) {
    map.set(key, value);
  }

  for (const source of sources) {
    // Within a source, last write wins; apply in order onto map.
    const local = new Map<string, "delete" | string>();
    for (const pair of source) {
      if (!Array.isArray(pair) || pair.length !== 2) {
        return { ok: false, code: "INVALID_QUERY_VALUE" };
      }
      const [key, value] = pair;
      if (typeof key !== "string") {
        return { ok: false, code: "INVALID_QUERY_KEY" };
      }
      if (key === "") {
        return { ok: false, code: "INVALID_QUERY_KEY" };
      }
      if (isUndefinedMarker(value) || value === null || value === undefined) {
        local.set(key, "delete");
        continue;
      }
      if (isComposite(value)) {
        return { ok: false, code: "INVALID_QUERY_VALUE" };
      }
      if (
        typeof value !== "string" &&
        typeof value !== "number" &&
        typeof value !== "boolean"
      ) {
        return { ok: false, code: "INVALID_QUERY_VALUE" };
      }
      if (typeof value === "number" && !Number.isFinite(value)) {
        return { ok: false, code: "INVALID_QUERY_VALUE" };
      }
      local.set(key, serializeScalar(value as QueryScalar));
    }
    for (const [key, value] of local) {
      if (value === "delete") {
        map.delete(key);
      } else {
        map.set(key, value);
      }
    }
  }

  const keys = [...map.keys()].sort((a, b) => {
    // Unicode code point order (UTF-16 code unit order matches for BMP + surrogate pairs as code points when comparing well-formed strings via localeCompare with 'kn' or by iterating code points).
    const aPoints = [...a];
    const bPoints = [...b];
    const n = Math.min(aPoints.length, bPoints.length);
    for (let i = 0; i < n; i++) {
      const ca = aPoints[i]!.codePointAt(0)!;
      const cb = bPoints[i]!.codePointAt(0)!;
      if (ca !== cb) {
        return ca - cb;
      }
    }
    return aPoints.length - bPoints.length;
  });

  if (keys.length === 0) {
    return { ok: true, url: parsed.path + parsed.fragment };
  }

  const query = keys
    .map((key) => `${encodeRFC3986(key)}=${encodeRFC3986(map.get(key)!)}`)
    .join("&");
  return { ok: true, url: `${parsed.path}?${query}${parsed.fragment}` };
}
