/**
 * R5 D-FORM form control surface (frozen §5 whitelist).
 *
 * Wire rules from schema-ui-docs@2.7.0 (fixed commit ca9e5fe…):
 *  - base: input → string, select (single) → string
 *  - 2.6 (capability form.controls.extended): textarea → string, switch → boolean,
 *    checkbox → boolean, radio → single string, select.mode=multiple → string[]
 *  - 2.7 (capability form.controls.advanced): cascader → path string[],
 *    checkboxGroup → value string[], richText → markdown string, password → string
 *  - 2.7 props.defaultValue must match the field's wire type
 *
 * The schema-driven page gate (meta.protocolVersion + requiredCapabilities) is
 * enforced by checkFormCapabilities; the Renderer-level page gate lands in 2c.
 */

export type FormControlType =
  | "input"
  | "select"
  | "textarea"
  | "switch"
  | "checkbox"
  | "radio"
  | "cascader"
  | "checkboxGroup"
  | "richText"
  | "password";

export const FORM_CONTROLS_EXTENDED_CAPABILITY = "form.controls.extended";
export const FORM_CONTROLS_ADVANCED_CAPABILITY = "form.controls.advanced";

/** Whitelisted base controls (no capability gate). */
const BASE_CONTROLS = new Set<FormControlType>(["input", "select"]);
/** 2.6 controls: require protocol >= 2.6 + form.controls.extended. */
const EXTENDED_CONTROLS = new Set<FormControlType>([
  "textarea",
  "switch",
  "checkbox",
  "radio",
]);
/** 2.7 controls: require protocol >= 2.7 + form.controls.advanced. */
const ADVANCED_CONTROLS = new Set<FormControlType>([
  "cascader",
  "checkboxGroup",
  "richText",
  "password",
]);

export type WireKind = "string" | "boolean" | "string-array";

export interface FormOption {
  value: string;
  label?: string;
}

export interface FormControlField {
  id: string;
  label?: string;
  type: FormControlType;
  /** select only: single (default) or multiple. */
  mode?: "single" | "multiple";
  options?: FormOption[];
  defaultValue?: unknown;
}

export interface FormControlMeta {
  protocolVersion: string;
  requiredCapabilities: string[];
}

export interface FormControlGateError {
  code: string;
  path: string;
  message: string;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

export function isWhitelistedFormControl(type: string): type is FormControlType {
  return (
    BASE_CONTROLS.has(type as FormControlType) ||
    EXTENDED_CONTROLS.has(type as FormControlType) ||
    ADVANCED_CONTROLS.has(type as FormControlType)
  );
}

export function wireKindOf(field: FormControlField): WireKind {
  switch (field.type) {
    case "switch":
    case "checkbox":
      return "boolean";
    case "select":
      return field.mode === "multiple" ? "string-array" : "string";
    case "cascader":
    case "checkboxGroup":
      return "string-array";
    default:
      return "string";
  }
}

/** Coerces a raw control value to its wire kind; defaultValue applies when raw is empty. */
export function coerceFieldValue(field: FormControlField, raw: unknown): unknown {
  const kind = wireKindOf(field);
  const present =
    raw !== undefined &&
    raw !== null &&
    !(kind === "string" && raw === "") &&
    !(kind === "string-array" && Array.isArray(raw) && raw.length === 0);

  if (!present && field.defaultValue !== undefined) {
    return coerceToKind(kind, field.defaultValue);
  }
  return coerceToKind(kind, raw);
}

function coerceToKind(kind: WireKind, value: unknown): unknown {
  switch (kind) {
    case "boolean":
      return typeof value === "boolean" ? value : false;
    case "string-array":
      return Array.isArray(value)
        ? value.filter((entry): entry is string => typeof entry === "string")
        : typeof value === "string" && value.trim() !== ""
          ? value.split(",").map((part) => part.trim()).filter((part) => part !== "")
          : [];
    case "string":
      // Resource rows often store multi-value fields as string[] (e.g. roles).
      // Textareas and free-text inputs expect a single wire string.
      if (Array.isArray(value)) {
        return value
          .filter((entry): entry is string => typeof entry === "string")
          .join(", ");
      }
      return value === undefined || value === null ? "" : String(value);
  }
}

/** Fails closed when defaultValue does not match the field's wire type. */
export function validateDefaultValue(field: FormControlField): FormControlGateError | null {
  if (field.defaultValue === undefined) {
    return null;
  }
  const kind = wireKindOf(field);
  const matches =
    kind === "boolean"
      ? typeof field.defaultValue === "boolean"
      : kind === "string-array"
        ? Array.isArray(field.defaultValue) &&
          field.defaultValue.every((entry) => typeof entry === "string")
        : typeof field.defaultValue === "string";
  if (matches) {
    return null;
  }
  return {
    code: "DEFAULT_VALUE_TYPE_MISMATCH",
    path: `fields[${field.id}].defaultValue`,
    message: `defaultValue for ${field.type} must be a ${kind}`,
  };
}

function versionAtLeast(version: string, major: number, minor: number): boolean {
  const match = /^(\d+)\.(\d+)$/.exec(version);
  if (!match) {
    return false;
  }
  const gotMajor = Number(match[1]);
  const gotMinor = Number(match[2]);
  return gotMajor > major || (gotMajor === major && gotMinor >= minor);
}

/** Gates a control set against page meta (P-005 / frozen capability rules). */
export function checkFormCapabilities(
  meta: FormControlMeta,
  fields: FormControlField[],
): FormControlGateError[] {
  const errors: FormControlGateError[] = [];
  const capabilities = new Set(meta.requiredCapabilities ?? []);
  const used = new Set(fields.map((field) => field.type));

  for (const type of [...used]) {
    if (EXTENDED_CONTROLS.has(type)) {
      if (!versionAtLeast(meta.protocolVersion, 2, 6)) {
        errors.push({
          code: "FORM_VERSION_TOO_LOW",
          path: `meta.protocolVersion`,
          message: `${type} requires protocol >= 2.6`,
        });
      }
      if (!capabilities.has(FORM_CONTROLS_EXTENDED_CAPABILITY)) {
        errors.push({
          code: "FORM_CAPABILITY_REQUIRED",
          path: "meta.requiredCapabilities",
          message: `${type} requires form.controls.extended`,
        });
      }
    }
    if (ADVANCED_CONTROLS.has(type)) {
      if (!versionAtLeast(meta.protocolVersion, 2, 7)) {
        errors.push({
          code: "FORM_VERSION_TOO_LOW",
          path: "meta.protocolVersion",
          message: `${type} requires protocol >= 2.7`,
        });
      }
      if (!capabilities.has(FORM_CONTROLS_ADVANCED_CAPABILITY)) {
        errors.push({
          code: "FORM_CAPABILITY_REQUIRED",
          path: "meta.requiredCapabilities",
          message: `${type} requires form.controls.advanced`,
        });
      }
    }
    if (!isWhitelistedFormControl(type)) {
      errors.push({
        code: "FORM_TYPE_NOT_WHITELISTED",
        path: `fields[${type}].type`,
        message: `${type} is outside the frozen §5 whitelist`,
      });
    }
  }

  for (const field of fields) {
    if (field.type === "select" && field.mode === "multiple") {
      if (!versionAtLeast(meta.protocolVersion, 2, 6)) {
        errors.push({
          code: "FORM_VERSION_TOO_LOW",
          path: `fields[${field.id}].mode`,
          message: "select multiple requires protocol >= 2.6",
        });
      }
      if (!capabilities.has(FORM_CONTROLS_EXTENDED_CAPABILITY)) {
        errors.push({
          code: "FORM_CAPABILITY_REQUIRED",
          path: `fields[${field.id}].mode`,
          message: "select multiple requires form.controls.extended",
        });
      }
    }
    if (field.defaultValue !== undefined) {
      // Frozen 2.7 rule: any field defaultValue requires 2.7 + advanced.
      if (!versionAtLeast(meta.protocolVersion, 2, 7)) {
        errors.push({
          code: "FORM_VERSION_TOO_LOW",
          path: `fields[${field.id}].defaultValue`,
          message: "defaultValue requires protocol >= 2.7",
        });
      }
      if (!capabilities.has(FORM_CONTROLS_ADVANCED_CAPABILITY)) {
        errors.push({
          code: "FORM_CAPABILITY_REQUIRED",
          path: `fields[${field.id}].defaultValue`,
          message: "defaultValue requires form.controls.advanced",
        });
      }
    }
    const defaultValueError = validateDefaultValue(field);
    if (defaultValueError !== null) {
      errors.push(defaultValueError);
    }
  }
  return errors;
}

/** Validates an arbitrary page-meta value without assuming record types. */
export function checkFormCapabilitiesRaw(
  metaValue: unknown,
  fields: FormControlField[],
): FormControlGateError[] {
  if (!isRecord(metaValue)) {
    return [
      {
        code: "FORM_META_INVALID",
        path: "meta",
        message: "page meta must be an object",
      },
    ];
  }
  const protocolVersion =
    typeof metaValue.protocolVersion === "string" ? metaValue.protocolVersion : "";
  const requiredCapabilities = Array.isArray(metaValue.requiredCapabilities)
    ? metaValue.requiredCapabilities.filter(
        (capability): capability is string => typeof capability === "string",
      )
    : [];
  return checkFormCapabilities({ protocolVersion, requiredCapabilities }, fields);
}
