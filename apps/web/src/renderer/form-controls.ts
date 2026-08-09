/**
 * D-FORM form control surface (frozen §5 whitelist + I-PROTO-FULL-001 full
 * registry surface).
 *
 * Wire rules from schema-ui-docs@2.7.0 (fixed commit ca9e5fe…):
 *  - base: input → string, select (single) → string, inputNumber → number,
 *    datePicker → ISO 8601 string, dateRangePicker → {start,end} pair bound
 *    to startField/endField (registry props; no single-field wire)
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
  | "inputNumber"
  | "datePicker"
  | "dateRangePicker"
  | "textarea"
  | "switch"
  | "checkbox"
  | "radio"
  | "cascader"
  | "checkboxGroup"
  | "richText"
  | "password"
  | "upload";

export const FORM_CONTROLS_EXTENDED_CAPABILITY = "form.controls.extended";
export const FORM_CONTROLS_ADVANCED_CAPABILITY = "form.controls.advanced";

/** Whitelisted base controls (no capability gate; registry has no `since`). */
const BASE_CONTROLS = new Set<FormControlType>([
  "input",
  "select",
  "inputNumber",
  "datePicker",
  "dateRangePicker",
  "upload",
]);
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

export type WireKind = "string" | "boolean" | "string-array" | "number" | "date-range";

export interface FormOption {
  value: string;
  label?: string;
  /** S2 (VP-007): i18n key resolved before `label` (upstream registry field). */
  labelKey?: string;
}

export interface DateRangeValue {
  start: string;
  end: string;
}

export interface FormControlField {
  id: string;
  label?: string;
  /** S2 (VP-007): i18n key resolved before `label` (missing-key observable). */
  labelKey?: string;
  placeholder?: string;
  /** S2 (VP-007): i18n key resolved before `placeholder` (local doc convention). */
  placeholderKey?: string;
  type: FormControlType;
  /** select only: single (default) or multiple. */
  mode?: "single" | "multiple";
  options?: FormOption[];
  defaultValue?: unknown;
  /** dateRangePicker only: the two bound output fields (registry props). */
  startField?: string;
  endField?: string;
  /** inputNumber constraints (registry props, since 0.2.1). */
  min?: number;
  max?: number;
  step?: number;
  precision?: number;
  /** datePicker display format (display-only; data stays ISO 8601). */
  format?: string;
  /** upload only: direct-URL mode (registry oneOf with actionRef). */
  action?: string;
  /** upload only: references a top-level type=upload action (requires actions.upload). */
  actionRef?: string;
  /** upload constraints (direct-URL mode only; actionRef mode reads the action). */
  accept?: string;
  maxSize?: number;
  multiple?: boolean;
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
    case "inputNumber":
      return "number";
    case "dateRangePicker":
      return "date-range";
    case "upload":
      // 07-actions-contract.md §7.2: multiple → array, single → string.
      return field.multiple === true ? "string-array" : "string";
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
    !(kind === "string-array" && Array.isArray(raw) && raw.length === 0) &&
    !(kind === "number" && raw === "") &&
    !(kind === "date-range" && isRecord(raw) && raw.start === "" && raw.end === "");

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
    case "number": {
      const numeric = typeof value === "number" ? value : Number(value);
      return Number.isFinite(numeric) ? numeric : 0;
    }
    case "date-range": {
      const record = isRecord(value) ? value : {};
      return {
        start: typeof record.start === "string" ? record.start : "",
        end: typeof record.end === "string" ? record.end : "",
      };
    }
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
  if (field.type === "dateRangePicker") {
    // Registry: dateRangePicker has no defaultValue prop (binds two fields).
    return {
      code: "DATE_RANGE_DEFAULT_VALUE_FORBIDDEN",
      path: `fields[${field.id}].defaultValue`,
      message: "dateRangePicker has no defaultValue prop (binds startField/endField)",
    };
  }
  const kind = wireKindOf(field);
  const matches =
    kind === "boolean"
      ? typeof field.defaultValue === "boolean"
      : kind === "string-array"
        ? Array.isArray(field.defaultValue) &&
          field.defaultValue.every((entry) => typeof entry === "string")
        : kind === "number"
          ? typeof field.defaultValue === "number" && Number.isFinite(field.defaultValue)
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
    if (field.type === "dateRangePicker") {
      // Registry: startField / endField are required non-empty props.
      if (typeof field.startField !== "string" || field.startField === "") {
        errors.push({
          code: "DATE_RANGE_START_FIELD_REQUIRED",
          path: `fields[${field.id}].startField`,
          message: "dateRangePicker requires a non-empty startField",
        });
      }
      if (typeof field.endField !== "string" || field.endField === "") {
        errors.push({
          code: "DATE_RANGE_END_FIELD_REQUIRED",
          path: `fields[${field.id}].endField`,
          message: "dateRangePicker requires a non-empty endField",
        });
      }
    }
    if (field.type === "upload") {
      // Registry oneOf: exactly one of action / actionRef (fail-closed).
      const hasAction = typeof field.action === "string" && field.action !== "";
      const hasActionRef = typeof field.actionRef === "string" && field.actionRef !== "";
      if (!hasAction && !hasActionRef) {
        errors.push({
          code: "UPLOAD_ACTION_REQUIRED",
          path: `fields[${field.id}].action`,
          message: "upload requires exactly one of action / actionRef",
        });
      }
      if (hasAction && hasActionRef) {
        errors.push({
          code: "UPLOAD_ACTION_CONFLICT",
          path: `fields[${field.id}].actionRef`,
          message: "upload must not declare both action and actionRef",
        });
      }
      if (hasActionRef && !capabilities.has("actions.upload")) {
        errors.push({
          code: "UPLOAD_CAPABILITY_REQUIRED",
          path: "meta.requiredCapabilities",
          message: "upload actionRef requires actions.upload",
        });
      }
      if (hasAction && !/^\/(?!\/)[^\s\\?#]*$/.test(field.action!)) {
        errors.push({
          code: "UPLOAD_ACTION_INVALID",
          path: `fields[${field.id}].action`,
          message: "upload action must be a single-slash same-origin path",
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
