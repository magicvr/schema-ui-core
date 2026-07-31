/**
 * response-mapping fixture adapter (schema-ui-docs v2.7.0 / ADR-0005).
 */

type Json = unknown;

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

/** Own-property path walk; prototype properties do not count as present. */
function getPath(root: unknown, path: string): { found: boolean; value: Json } {
  if (path === "") {
    return { found: true, value: root };
  }
  const parts = path.split(".");
  let current: unknown = root;
  for (const part of parts) {
    if (!isRecord(current) || !Object.prototype.hasOwnProperty.call(current, part)) {
      return { found: false, value: undefined };
    }
    current = current[part];
  }
  return { found: true, value: current as Json };
}

function mappingIsEmptyObject(mapping: unknown): boolean {
  return isRecord(mapping) && Object.keys(mapping).length === 0;
}

function supportsMapping(component: string): boolean {
  return (
    component === "table" ||
    component === "chart" ||
    component === "formRecord" ||
    component === "recordView"
  );
}

function wireOk(value: unknown, wire: string | undefined): boolean {
  if (wire === undefined || value === null) {
    return true;
  }
  switch (wire) {
    case "string":
      return typeof value === "string";
    case "boolean":
      return typeof value === "boolean";
    case "array":
      return Array.isArray(value);
    case "number":
      return typeof value === "number" && Number.isFinite(value);
    default:
      return true;
  }
}

export function mapResponse(input: Record<string, unknown>): Record<string, unknown> {
  const component = input.component as string;

  // Explicit null mapping invalid for any component that receives localMapping.
  if (Object.prototype.hasOwnProperty.call(input, "localMapping") && input.localMapping === null) {
    return { ok: false, code: "INVALID_RESPONSE_MAPPING", path: "localMapping" };
  }

  if (!supportsMapping(component)) {
    if (
      Object.prototype.hasOwnProperty.call(input, "localMapping") ||
      Object.prototype.hasOwnProperty.call(input, "datasourceMapping")
    ) {
      // text/statCard reject mapping support (path reported as localMapping per fixtures).
      if (input.localMapping === null) {
        return { ok: false, code: "INVALID_RESPONSE_MAPPING", path: "localMapping" };
      }
      return { ok: false, code: "RESPONSE_MAPPING_NOT_SUPPORTED", path: "localMapping" };
    }
  }

  if (component === "formRecord" || component === "recordView") {
    const mapping = input.responseMapping;
    if (!isRecord(mapping) || mappingIsEmptyObject(mapping)) {
      return { ok: false, code: "INVALID_RESPONSE_MAPPING", path: "responseMapping" };
    }
    const response = input.response;
    const fieldWireTypes = (input.fieldWireTypes as Record<string, string>) ?? {};
    const values: Record<string, unknown> = {};
    const skipped: Record<string, string> = {};
    for (const [target, path] of Object.entries(mapping)) {
      if (typeof path !== "string") {
        return { ok: false, code: "INVALID_RESPONSE_MAPPING", path: "responseMapping" };
      }
      const resolved = getPath(response, path);
      if (!resolved.found) {
        values[target] = null;
        continue;
      }
      const wire = fieldWireTypes[target];
      if (!wireOk(resolved.value, wire)) {
        skipped[target] = "FIELD_WIRE_TYPE_MISMATCH";
        continue;
      }
      values[target] = resolved.value;
    }
    if (Object.keys(skipped).length > 0) {
      return { ok: true, values, skipped };
    }
    return { ok: true, values };
  }

  if (component === "chart") {
    const response = input.response;
    if (Array.isArray(response)) {
      return { ok: true, data: { list: response } };
    }
    // fall through to mapping if present
  }

  // table (and chart with mapping)
  let mapping: Record<string, string> | undefined;
  const hasLocal = Object.prototype.hasOwnProperty.call(input, "localMapping");
  const hasDs = Object.prototype.hasOwnProperty.call(input, "datasourceMapping");

  if (hasLocal) {
    if (!isRecord(input.localMapping)) {
      return { ok: false, code: "INVALID_RESPONSE_MAPPING", path: "localMapping" };
    }
    if (mappingIsEmptyObject(input.localMapping)) {
      return { ok: false, code: "INVALID_RESPONSE_MAPPING", path: "localMapping" };
    }
    mapping = input.localMapping as Record<string, string>;
  } else if (hasDs) {
    mapping = input.datasourceMapping as Record<string, string>;
  }

  const response = input.response;
  const paginationMode = input.paginationMode as string | undefined;

  if (!mapping) {
    // Default table mapping
    if (component === "table") {
      if (!isRecord(response)) {
        return { ok: false, code: "INVALID_RESPONSE_MAPPING", path: "response" };
      }
      const list = response.list;
      if (!Array.isArray(list)) {
        return { ok: false, code: "RESPONSE_MAPPING_PATH_MISSING", path: "list" };
      }
      if (paginationMode === "server") {
        if (typeof response.total !== "number") {
          return { ok: false, code: "RESPONSE_MAPPING_TYPE_MISMATCH", path: "total" };
        }
        return { ok: true, data: { list, total: response.total } };
      }
      return { ok: true, data: { list } };
    }
    if (component === "chart" && Array.isArray(response)) {
      return { ok: true, data: { list: response } };
    }
  }

  const listPath = mapping?.list;
  if (typeof listPath !== "string") {
    return { ok: false, code: "INVALID_RESPONSE_MAPPING", path: "localMapping" };
  }
  const listResolved = getPath(response, listPath);
  if (!listResolved.found) {
    return { ok: false, code: "RESPONSE_MAPPING_PATH_MISSING", path: listPath };
  }
  if (!Array.isArray(listResolved.value)) {
    return { ok: false, code: "RESPONSE_MAPPING_TYPE_MISMATCH", path: listPath };
  }

  if (paginationMode === "server") {
    const totalPath = mapping?.total;
    if (typeof totalPath !== "string") {
      return { ok: false, code: "INVALID_RESPONSE_MAPPING", path: "localMapping" };
    }
    const totalResolved = getPath(response, totalPath);
    if (!totalResolved.found) {
      return { ok: false, code: "RESPONSE_MAPPING_PATH_MISSING", path: totalPath };
    }
    if (typeof totalResolved.value !== "number") {
      return { ok: false, code: "RESPONSE_MAPPING_TYPE_MISMATCH", path: totalPath };
    }
    return {
      ok: true,
      data: { list: listResolved.value, total: totalResolved.value },
    };
  }

  return { ok: true, data: { list: listResolved.value } };
}
