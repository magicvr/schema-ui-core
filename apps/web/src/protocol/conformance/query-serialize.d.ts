/**
 * ADR-0010 query-serialization adapter (schema-ui-docs v2.7.0 fixtures).
 *
 * Merges base URL query with ordered source layers; encodes keys/values per
 * RFC 3986 (unreserved A-Za-z0-9-._~); sorts final keys by Unicode code point.
 */
export type QueryScalar = string | number | boolean | null;
export type QueryPair = [string, unknown];
export type QuerySource = QueryPair[];
export type QuerySerializeResult = {
    ok: true;
    url: string;
} | {
    ok: false;
    code: "INVALID_BASE_URL_QUERY" | "INVALID_QUERY_KEY" | "INVALID_QUERY_VALUE";
};
/** RFC3986 encode: percent-encode everything except unreserved. */
export declare function encodeRFC3986(value: string): string;
/**
 * JCS-style number serialization used by query-serialization fixtures.
 * Matches JSON number text for typical finite values (1e+21, 1e-7, 0.000001).
 */
export declare function serializeQueryNumber(value: number): string;
export declare function serializeQuery(baseUrl: string, sources: QuerySource[]): QuerySerializeResult;
