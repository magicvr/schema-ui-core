/**
 * D-COMP component-format fixture adapter (schema-ui-docs v2.7.0).
 * Validates format wire types without coercion.
 */

export type ComponentFormat =
  | "currency"
  | "percent"
  | "datetime"
  | string;

export type FormatResult =
  | { ok: true; value: unknown }
  | { ok: false; code: "COMPONENT_DATA_TYPE_MISMATCH" };

export function applyComponentFormat(
  format: ComponentFormat,
  value: unknown,
): FormatResult {
  switch (format) {
    case "currency":
    case "percent":
      if (typeof value === "number" && Number.isFinite(value)) {
        return { ok: true, value };
      }
      return { ok: false, code: "COMPONENT_DATA_TYPE_MISMATCH" };
    case "datetime":
      if (typeof value === "string") {
        return { ok: true, value };
      }
      return { ok: false, code: "COMPONENT_DATA_TYPE_MISMATCH" };
    default:
      return { ok: false, code: "COMPONENT_DATA_TYPE_MISMATCH" };
  }
}
