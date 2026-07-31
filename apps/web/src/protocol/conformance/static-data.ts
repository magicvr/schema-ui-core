/**
 * static-data fixture adapter (schema-ui-docs v2.7.0).
 */

type Json = unknown;

export type StaticDataResult =
  | { ok: true; value: Json; network: false }
  | {
      ok: false;
      code:
        | "STATIC_DATA_SHAPE_MISMATCH"
        | "STATIC_DATA_REF_INVALID"
        | "STATIC_RESPONSE_MAPPING_NOT_ALLOWED";
    };

interface DataRef {
  source: "static" | "ref" | "api";
  value?: Json;
  ref?: string;
  responseMapping?: unknown;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function resolveRef(
  ref: string,
  datasources: Record<string, unknown> | undefined,
): { ok: true; value: Json } | { ok: false; code: "STATIC_DATA_REF_INVALID" } {
  if (!datasources || !Object.prototype.hasOwnProperty.call(datasources, ref)) {
    return { ok: false, code: "STATIC_DATA_REF_INVALID" };
  }
  const ds = datasources[ref];
  if (!isRecord(ds) || ds.source !== "static") {
    return { ok: false, code: "STATIC_DATA_REF_INVALID" };
  }
  return { ok: true, value: ds.value as Json };
}

function shapeOk(component: string, value: Json, valueField?: string): boolean {
  switch (component) {
    case "table":
    case "chart":
      return Array.isArray(value);
    case "text":
      return (
        value === null ||
        typeof value === "string" ||
        typeof value === "number" ||
        typeof value === "boolean"
      );
    case "statCard":
      if (valueField !== undefined) {
        // Already extracted field; null rejected.
        return value !== null && value !== undefined && !Array.isArray(value) && typeof value !== "object";
      }
      // Bare static value for statCard: object ok if we will read field; null rejected.
      if (value === null) {
        return false;
      }
      return true;
    default:
      return true;
  }
}

export function resolveStaticData(input: {
  component: string;
  data: DataRef;
  datasources?: Record<string, unknown>;
  props?: { valueField?: string };
}): StaticDataResult {
  const { component, data, datasources, props } = input;

  if (data.responseMapping !== undefined) {
    return { ok: false, code: "STATIC_RESPONSE_MAPPING_NOT_ALLOWED" };
  }

  let value: Json;
  if (data.source === "static") {
    value = data.value as Json;
  } else if (data.source === "ref") {
    if (typeof data.ref !== "string") {
      return { ok: false, code: "STATIC_DATA_REF_INVALID" };
    }
    const resolved = resolveRef(data.ref, datasources);
    if (!resolved.ok) {
      return resolved;
    }
    value = resolved.value;
    if (component === "statCard" && props?.valueField) {
      if (!isRecord(value) || !Object.prototype.hasOwnProperty.call(value, props.valueField)) {
        return { ok: false, code: "STATIC_DATA_SHAPE_MISMATCH" };
      }
      value = value[props.valueField];
      if (!shapeOk(component, value, props.valueField)) {
        return { ok: false, code: "STATIC_DATA_SHAPE_MISMATCH" };
      }
      return { ok: true, value, network: false };
    }
  } else {
    return { ok: false, code: "STATIC_DATA_SHAPE_MISMATCH" };
  }

  if (!shapeOk(component, value, props?.valueField)) {
    return { ok: false, code: "STATIC_DATA_SHAPE_MISMATCH" };
  }

  // text rejects arrays
  if (component === "text" && Array.isArray(value)) {
    return { ok: false, code: "STATIC_DATA_SHAPE_MISMATCH" };
  }
  if (component === "statCard" && value === null) {
    return { ok: false, code: "STATIC_DATA_SHAPE_MISMATCH" };
  }
  if (component === "table" && !Array.isArray(value)) {
    return { ok: false, code: "STATIC_DATA_SHAPE_MISMATCH" };
  }

  return { ok: true, value, network: false };
}
