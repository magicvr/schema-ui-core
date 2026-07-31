/**
 * request-construction fixture adapter (schema-ui-docs v2.7.0).
 *
 * Builds HTTP request / navigation / modal outcomes from declarative mappings.
 * Batch kinds are out of scope for MVP callers (Q1); this module still implements
 * non-batch kinds used by stage3 execution.
 */

import { encodeRFC3986, serializeQueryNumber } from "./query-serialize";

type Json = null | boolean | number | string | Json[] | { [k: string]: Json };
type JsonObject = Record<string, unknown>;

export type RequestConstructionResult =
  | {
      ok: true;
      request?: {
        method: string;
        url: string;
        body: unknown;
        headers?: Record<string, string>;
      };
      navigation?: { url: string };
      modalOpen?: { modalId: string };
      resolvedBase?: string;
    }
  | { ok: false; code: string; path: string };

const PROTOCOL_URL_RE = /^\/(?!\/)[^\s\\]*$/;
const UNSAFE_SEGMENTS = new Set(["__proto__", "constructor", "prototype"]);

function isObject(value: unknown): value is JsonObject {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function fail(code: string, path: string): RequestConstructionResult {
  return { ok: false, code, path };
}

function isProtocolRelativeUrl(url: string): boolean {
  return PROTOCOL_URL_RE.test(url);
}

function splitUrl(url: string): { path: string; query: Map<string, string> } {
  const qIndex = url.indexOf("?");
  if (qIndex < 0) {
    return { path: url, query: new Map() };
  }
  const path = url.slice(0, qIndex);
  const query = new Map<string, string>();
  const raw = url.slice(qIndex + 1);
  if (raw.length === 0) {
    return { path, query };
  }
  for (const part of raw.split("&")) {
    if (part === "") continue;
    const eq = part.indexOf("=");
    const key = eq < 0 ? part : part.slice(0, eq);
    const value = eq < 0 ? "" : part.slice(eq + 1);
    // Base query values in fixtures are already literal (not double-encoded).
    query.set(decodeURIComponent(key), decodeURIComponent(value));
  }
  return { path, query };
}

function serializeQueryValue(value: unknown): string {
  if (value === null) return "";
  if (typeof value === "boolean") return value ? "true" : "false";
  if (typeof value === "number") return serializeQueryNumber(value);
  if (typeof value === "string") return value;
  throw new Error("non-scalar query value");
}

function buildUrl(path: string, query: Map<string, string>): string {
  if (query.size === 0) return path;
  const keys = [...query.keys()].sort((a, b) => (a < b ? -1 : a > b ? 1 : 0));
  const parts = keys.map(
    (key) => `${encodeRFC3986(key)}=${encodeRFC3986(query.get(key) ?? "")}`,
  );
  return `${path}?${parts.join("&")}`;
}

function extractPathParams(path: string): string[] {
  const params: string[] = [];
  const re = /\{([A-Za-z_][A-Za-z0-9_]*)\}/g;
  let match: RegExpExecArray | null;
  while ((match = re.exec(path)) !== null) {
    params.push(match[1]);
  }
  return params;
}

function applyPathBindings(
  path: string,
  bindings: Record<string, string>,
  mappingPathPrefix: string,
): { ok: true; path: string } | { ok: false; code: string; path: string } {
  const needed = extractPathParams(path);
  const provided = Object.keys(bindings);
  for (const key of needed) {
    if (!(key in bindings)) {
      return { ok: false, code: "MISSING_PATH_BINDING", path: `${mappingPathPrefix}.${key}` };
    }
  }
  for (const key of provided) {
    if (!needed.includes(key)) {
      return { ok: false, code: "EXTRA_PATH_BINDING", path: `${mappingPathPrefix}.${key}` };
    }
  }
  let result = path;
  for (const key of needed) {
    result = result.replaceAll(`{${key}}`, encodeRFC3986(bindings[key]));
  }
  return { ok: true, path: result };
}

function resolveRowPath(
  expr: string,
  row: unknown,
  path: string,
): { ok: true; value: unknown } | { ok: false; code: string; path: string } {
  if (!expr.startsWith("$row.")) {
    return { ok: false, code: "INVALID_MAPPING_VALUE", path };
  }
  const segments = expr.slice("$row.".length).split(".");
  if (segments.some((s) => UNSAFE_SEGMENTS.has(s) || s === "")) {
    return { ok: false, code: "UNSAFE_ROW_PATH", path };
  }
  let current: unknown = row;
  for (const seg of segments) {
    if (!isObject(current) || !Object.prototype.hasOwnProperty.call(current, seg)) {
      return { ok: false, code: "UNRESOLVED_ROW_VALUE", path };
    }
    current = current[seg];
  }
  return { ok: true, value: current };
}

function resolveMappingValue(
  expr: unknown,
  row: unknown,
  path: string,
  opts: { allowNullAsTombstone?: boolean; pathSlot?: boolean },
):
  | { ok: true; value: unknown; tombstone?: boolean }
  | { ok: false; code: string; path: string } {
  if (typeof expr !== "string") {
    return { ok: true, value: expr };
  }
  if (!expr.startsWith("$")) {
    return { ok: true, value: expr };
  }
  if (expr.startsWith("$row.")) {
    const resolved = resolveRowPath(expr, row, path);
    if (!resolved.ok) return resolved;
    if (resolved.value === null) {
      if (opts.pathSlot) {
        return { ok: false, code: "NULL_PATH_VALUE", path };
      }
      if (opts.allowNullAsTombstone) {
        return { ok: true, value: null, tombstone: true };
      }
      return { ok: false, code: "NULL_PATH_VALUE", path };
    }
    if (opts.pathSlot) {
      if (
        typeof resolved.value !== "string" &&
        typeof resolved.value !== "number" &&
        typeof resolved.value !== "boolean"
      ) {
        return { ok: false, code: "INVALID_ROW_VALUE", path };
      }
      return { ok: true, value: String(resolved.value) };
    }
    // body/query scalars only (no nested objects/arrays for row body)
    if (isObject(resolved.value) || Array.isArray(resolved.value)) {
      return { ok: false, code: "INVALID_ROW_VALUE", path };
    }
    return { ok: true, value: resolved.value };
  }
  return { ok: false, code: "INVALID_MAPPING_VALUE", path };
}

function resolveRouteExpr(
  expr: string,
  route: { query?: JsonObject; params?: JsonObject },
  path: string,
): { ok: true; value: string } | { ok: false; code: string; path: string } {
  const queryPrefix = "$context.route.query.";
  const paramsPrefix = "$context.route.params.";
  if (expr.startsWith(queryPrefix)) {
    const key = expr.slice(queryPrefix.length);
    const value = route.query?.[key];
    if (value === undefined || value === null) {
      return { ok: false, code: "UNRESOLVED_ROUTE_VALUE", path };
    }
    return { ok: true, value: String(value) };
  }
  if (expr.startsWith(paramsPrefix)) {
    const key = expr.slice(paramsPrefix.length);
    const value = route.params?.[key];
    if (value === undefined || value === null) {
      return { ok: false, code: "UNRESOLVED_ROUTE_VALUE", path };
    }
    return { ok: true, value: String(value) };
  }
  return { ok: false, code: "INVALID_MAPPING_VALUE", path };
}

function joinBase(base: string | undefined, url: string): string {
  if (!base) return url;
  if (!url.startsWith("/")) return url;
  return `${base.replace(/\/$/, "")}${url}`;
}

function checkConfirm(input: JsonObject): RequestConstructionResult | null {
  if (input.confirm === undefined) return null;
  if (input.confirmAccepted === false) {
    return fail("CONFIRM_REJECTED", "confirm");
  }
  return null;
}

function applyIdempotency(
  action: JsonObject,
  invocationId: unknown,
  pathPrefix: string,
):
  | { ok: true; headers?: Record<string, string> }
  | { ok: false; code: string; path: string } {
  const policy = action.retryPolicy;
  if (policy === undefined || policy === "never") {
    return { ok: true };
  }
  if (policy !== "idempotent") {
    return { ok: false, code: "INVALID_RETRY_POLICY", path: `${pathPrefix}.retryPolicy` };
  }
  if (typeof invocationId !== "string" || invocationId.length === 0) {
    return { ok: false, code: "MISSING_INVOCATION_ID", path: "invocationId" };
  }
  return { ok: true, headers: { "Idempotency-Key": invocationId } };
}

function buildDataRef(input: JsonObject): RequestConstructionResult {
  const dataRef = input.dataRef as JsonObject;
  const url = dataRef.url as string;
  if (typeof url !== "string" || !isProtocolRelativeUrl(url)) {
    return fail("INVALID_PROTOCOL_URL", "dataRef.url");
  }
  const method = (dataRef.method as string | undefined) ?? "GET";
  if (method !== "GET") {
    return fail("DATA_REF_METHOD_NOT_READ_ONLY", "dataRef.method");
  }
  if (dataRef.requestInterceptor !== undefined) {
    return fail("INTERCEPTOR_VIOLATION", "requestInterceptor");
  }
  const { path, query } = splitUrl(url);
  const params = (dataRef.params as JsonObject | undefined) ?? {};
  for (const [key, value] of Object.entries(params)) {
    if (value === null) {
      query.delete(key);
      continue;
    }
    query.set(key, serializeQueryValue(value));
  }
  return {
    ok: true,
    request: { method: "GET", url: buildUrl(path, query), body: null },
  };
}

function isRelativeProtocolUrl(url: string): boolean {
  if (typeof url !== "string" || !url.startsWith("/") || url.startsWith("//")) {
    return false;
  }
  const pathOnly = url.includes("?") ? url.slice(0, url.indexOf("?")) : url;
  return isProtocolRelativeUrl(pathOnly) || isProtocolRelativeUrl(url);
}

function buildRowAction(input: JsonObject): RequestConstructionResult {
  const action = input.action as JsonObject;
  const method = action.method as string;
  const url = action.url as string;
  if (!isRelativeProtocolUrl(url)) {
    return fail("INVALID_PROTOCOL_URL", "action.url");
  }
  const mapping = (input.requestMapping as JsonObject | undefined) ?? {};
  const row = input.row;
  const { path: basePath, query } = splitUrl(url);

  const pathMap = (mapping.path as JsonObject | undefined) ?? {};
  const bindings: Record<string, string> = {};
  for (const [key, expr] of Object.entries(pathMap)) {
    const resolved = resolveMappingValue(expr, row, `requestMapping.path.${key}`, {
      pathSlot: true,
    });
    if (!resolved.ok) return resolved;
    bindings[key] = String(resolved.value);
  }
  const bound = applyPathBindings(basePath, bindings, "requestMapping.path");
  if (!bound.ok) return bound;

  const queryMap = (mapping.query as JsonObject | undefined) ?? {};
  for (const [key, expr] of Object.entries(queryMap)) {
    const resolved = resolveMappingValue(expr, row, `requestMapping.query.${key}`, {
      allowNullAsTombstone: true,
    });
    if (!resolved.ok) return resolved;
    if (resolved.tombstone) {
      query.delete(key);
      continue;
    }
    query.set(key, serializeQueryValue(resolved.value));
  }

  let body: unknown = null;
  const bodyMap = mapping.body as JsonObject | undefined;
  if (bodyMap !== undefined) {
    const out: JsonObject = {};
    for (const [key, expr] of Object.entries(bodyMap)) {
      const resolved = resolveMappingValue(expr, row, `requestMapping.body.${key}`, {});
      if (!resolved.ok) return resolved;
      out[key] = resolved.value as Json;
    }
    body = out;
  }

  const idem = applyIdempotency(action, input.invocationId, "action");
  if (!idem.ok) return idem;

  return {
    ok: true,
    request: {
      method,
      url: buildUrl(bound.path, query),
      body,
      ...(idem.headers ? { headers: idem.headers } : {}),
    },
  };
}

function formFieldIncluded(proj: JsonObject | undefined): boolean {
  if (!proj) return true;
  if (proj.visible === false) return false;
  if (proj.disabled === true) return false;
  if (proj.uploadStatus === "error") return false;
  return Object.prototype.hasOwnProperty.call(proj, "value");
}

function buildFormAction(input: JsonObject): RequestConstructionResult {
  const action = input.action as JsonObject;
  const method = action.method as string;
  const url = action.url as string;
  if (typeof url !== "string" || !isProtocolRelativeUrl(url)) {
    return fail("INVALID_PROTOCOL_URL", "action.url");
  }
  if (method === "GET") {
    return fail("FORM_GET_NOT_ALLOWED", "action.method");
  }

  const formValues = (input.formValues as JsonObject | undefined) ?? {};
  const formProjection = input.formProjection as Record<string, JsonObject> | undefined;
  const bodyMapping = action.bodyMapping as Record<string, string> | undefined;

  let sourceValues: JsonObject = {};
  if (formProjection) {
    for (const [field, proj] of Object.entries(formProjection)) {
      if (!formFieldIncluded(proj)) continue;
      sourceValues[field] = proj.value;
    }
  } else {
    sourceValues = { ...formValues };
  }

  let body: JsonObject = {};
  if (bodyMapping) {
    for (const [source, target] of Object.entries(bodyMapping)) {
      if (!Object.prototype.hasOwnProperty.call(sourceValues, source)) {
        // When projection omits field (no value key) or formValues missing
        return fail("UNRESOLVED_FORM_VALUE", `bodyMapping.${source}`);
      }
      body[target] = sourceValues[source];
    }
  } else {
    body = { ...sourceValues };
  }

  const idem = applyIdempotency(action, input.invocationId, "action");
  if (!idem.ok) return idem;

  const finalUrl = joinBase(input.baseURL as string | undefined, url);
  const result: RequestConstructionResult = {
    ok: true,
    request: {
      method,
      url: finalUrl,
      body,
      ...(idem.headers ? { headers: idem.headers } : {}),
    },
  };
  if (input.baseURL) {
    (result as { resolvedBase?: string }).resolvedBase = "api.baseURL";
  }
  return result;
}

function buildRowNavigate(input: JsonObject): RequestConstructionResult {
  const action = input.action as JsonObject;
  const url = action.url as string;
  const mapping = input.navigateMapping as JsonObject | undefined;
  if (!mapping || Object.keys(mapping).length === 0) {
    return fail("EMPTY_NAVIGATE_MAPPING", "navigateMapping");
  }
  if (mapping.body !== undefined) {
    return fail("NAVIGATE_BODY_NOT_ALLOWED", "navigateMapping.body");
  }
  const row = input.row;
  const { path: basePath, query } = splitUrl(url);

  const pathMap = (mapping.path as JsonObject | undefined) ?? {};
  const bindings: Record<string, string> = {};
  for (const [key, expr] of Object.entries(pathMap)) {
    const resolved = resolveMappingValue(expr, row, `navigateMapping.path.${key}`, {
      pathSlot: true,
    });
    if (!resolved.ok) return resolved;
    bindings[key] = String(resolved.value);
  }
  const bound = applyPathBindings(basePath, bindings, "navigateMapping.path");
  if (!bound.ok) return bound;

  const queryMap = (mapping.query as JsonObject | undefined) ?? {};
  for (const [key, expr] of Object.entries(queryMap)) {
    if (typeof expr === "string" && expr.startsWith("$") && !expr.startsWith("$row.")) {
      return fail("INVALID_MAPPING_VALUE", `navigateMapping.query.${key}`);
    }
    const resolved = resolveMappingValue(expr, row, `navigateMapping.query.${key}`, {
      allowNullAsTombstone: true,
    });
    if (!resolved.ok) return resolved;
    if (resolved.tombstone) {
      query.delete(key);
      continue;
    }
    query.set(key, serializeQueryValue(resolved.value));
  }

  return {
    ok: true,
    navigation: { url: buildUrl(bound.path, query) },
  };
}

function buildRecordSource(input: JsonObject): RequestConstructionResult {
  const rs = input.recordSource as JsonObject;
  if (rs.ref !== undefined) {
    return fail("RECORD_SOURCE_REF_NOT_ALLOWED", "recordSource");
  }
  if (rs.method === undefined) {
    return fail("MISSING_RECORD_SOURCE_METHOD", "recordSource.method");
  }
  if (rs.method !== "GET") {
    return fail("RECORD_SOURCE_METHOD_NOT_GET", "recordSource.method");
  }
  const responseMapping = rs.responseMapping as JsonObject | undefined;
  if (!responseMapping || Object.keys(responseMapping).length === 0) {
    return fail("EMPTY_RESPONSE_MAPPING", "recordSource.responseMapping");
  }
  const url = rs.url as string;
  if (typeof url !== "string" || (!isProtocolRelativeUrl(url) && !url.startsWith("/"))) {
    return fail("INVALID_PROTOCOL_URL", "recordSource.url");
  }
  const route = (input.route as { query?: JsonObject; params?: JsonObject }) ?? {};
  const { path: basePath, query } = splitUrl(url);

  const pathMap = (rs.path as JsonObject | undefined) ?? {};
  const bindings: Record<string, string> = {};
  for (const [key, expr] of Object.entries(pathMap)) {
    if (typeof expr !== "string") {
      return fail("INVALID_MAPPING_VALUE", `recordSource.path.${key}`);
    }
    const resolved = resolveRouteExpr(expr, route, `recordSource.path.${key}`);
    if (!resolved.ok) return resolved;
    bindings[key] = resolved.value;
  }
  const bound = applyPathBindings(basePath, bindings, "recordSource.path");
  if (!bound.ok) return bound;

  const queryMap = (rs.query as JsonObject | undefined) ?? {};
  for (const [key, value] of Object.entries(queryMap)) {
    query.set(key, serializeQueryValue(value));
  }

  const finalUrl = joinBase(input.baseURL as string | undefined, buildUrl(bound.path, query));
  const result: RequestConstructionResult = {
    ok: true,
    request: { method: "GET", url: finalUrl, body: null },
  };
  if (input.baseURL) {
    (result as { resolvedBase?: string }).resolvedBase = "api.baseURL";
  }
  return result;
}

function buildPageTriggerRequest(input: JsonObject): RequestConstructionResult {
  const confirm = checkConfirm(input);
  if (confirm) return confirm;
  const action = input.action as JsonObject;
  const method = action.method as string;
  const url = action.url as string;
  if (method === "GET") {
    return fail("PAGE_TRIGGER_METHOD_NOT_ALLOWED", "action.method");
  }
  if (typeof url !== "string" || !url.startsWith("/") || url.startsWith("//")) {
    return fail("INVALID_PROTOCOL_URL", "action.url");
  }
  if (extractPathParams(url).length > 0) {
    return fail("UNBOUND_URL_TEMPLATE", "action.url");
  }
  return {
    ok: true,
    request: { method, url, body: null },
  };
}

function buildPageTriggerNavigate(input: JsonObject): RequestConstructionResult {
  const confirm = checkConfirm(input);
  if (confirm) return confirm;
  const action = input.action as JsonObject;
  const url = action.url as string;
  if (typeof url !== "string" || !url.startsWith("/") || url.startsWith("//")) {
    return fail("INVALID_PROTOCOL_URL", "action.url");
  }
  if (extractPathParams(url).length > 0) {
    return fail("UNBOUND_URL_TEMPLATE", "action.url");
  }
  const finalUrl = joinBase(input.appRouteRoot as string | undefined, url);
  const result: RequestConstructionResult = {
    ok: true,
    navigation: { url: finalUrl },
  };
  if (input.appRouteRoot) {
    (result as { resolvedBase?: string }).resolvedBase = "app.routeRoot";
  }
  return result;
}

function buildPageTriggerModal(input: JsonObject): RequestConstructionResult {
  const confirm = checkConfirm(input);
  if (confirm) return confirm;
  const action = input.action as JsonObject;
  return {
    ok: true,
    modalOpen: { modalId: String(action.modalId) },
  };
}

function buildOutcomeNavigate(input: JsonObject): RequestConstructionResult {
  const url = input.url as string;
  const finalUrl = joinBase(input.appRouteRoot as string | undefined, url);
  return {
    ok: true,
    navigation: { url: finalUrl },
    resolvedBase: "app.routeRoot",
  };
}

/**
 * Run one request-construction fixture case.
 * Batch kinds return a structured error; stage3 excludes them via Q1.
 */
export function constructRequest(input: Record<string, unknown>): RequestConstructionResult {
  const kind = input.kind as string;
  switch (kind) {
    case "dataRef":
      return buildDataRef(input);
    case "rowAction":
      return buildRowAction(input);
    case "formAction":
      return buildFormAction(input);
    case "rowNavigate":
      return buildRowNavigate(input);
    case "recordSource":
      return buildRecordSource(input);
    case "pageTriggerRequest":
      return buildPageTriggerRequest(input);
    case "pageTriggerNavigate":
      return buildPageTriggerNavigate(input);
    case "pageTriggerModal":
      return buildPageTriggerModal(input);
    case "outcomeNavigate":
      return buildOutcomeNavigate(input);
    case "batchRequest":
      return fail("PAGE_TRIGGER_METHOD_NOT_ALLOWED", "batchRequest");
    default:
      return fail("INVALID_MAPPING_VALUE", "kind");
  }
}
