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
export type FormControlType = "input" | "select" | "inputNumber" | "datePicker" | "dateRangePicker" | "textarea" | "switch" | "checkbox" | "radio" | "cascader" | "checkboxGroup" | "richText" | "password" | "upload";
export declare const FORM_CONTROLS_EXTENDED_CAPABILITY = "form.controls.extended";
export declare const FORM_CONTROLS_ADVANCED_CAPABILITY = "form.controls.advanced";
/** ADR-0021: `form.props.recordSource` prefill GET (registry since 2.1). */
export declare const FORM_RECORD_LOAD_CAPABILITY = "form.record.load";
/** ADR-0040 (since 2.9): `readOnly` field declaration (value still projects). */
export declare const FORM_CONTROLS_READONLY_CAPABILITY = "form.controls.readonly";
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
    /**
     * W11 · U-01/U-02 — dynamic option source, aligned with the upstream
     * registry shape (component-registry.json, since 0.2): an object with a
     * required single-slash same-origin url plus the response item fields used
     * for value/label. The response is {items:[...]} (or a bare array); while
     * the source loads, static options (if any) remain the fallback; an invalid
     * source or failed fetch fails closed to an empty option set.
     */
    optionsSource?: {
        url: string;
        /** Optional scalar query params appended to url (e.g. pageSize). */
        params?: Record<string, string | number | boolean | null>;
        labelField: string;
        valueField: string;
    };
    defaultValue?: unknown;
    /** dateRangePicker only: the two bound output fields (registry props). */
    startField?: string;
    endField?: string;
    /** inputNumber constraints (registry props, since 0.2.1). */
    min?: number;
    max?: number;
    step?: number;
    precision?: number;
    /** GOAL-014 D-002 §3: field-level validation constraints (optional). */
    required?: boolean;
    /** ADR-0040 (since 2.9): read-only field — user cannot edit, value still
     * participates in values and the submit projection (bodyMapping);
     * recordSource backfill and reactions keep writing. Requires protocol
     * >= 2.9 and form.controls.readonly. */
    readOnly?: boolean;
    /** Regex pattern for string-typed fields (submit-time validation). */
    pattern?: string;
    /** String length bounds for input/textarea. */
    minLength?: number;
    maxLength?: number;
    /** datePicker display format (display-only; data stays ISO 8601). */
    format?: string;
    /** upload only: direct-URL mode (registry oneOf with actionRef). */
    action?: string;
    /** upload only: references a top-level type=upload action (requires actions.upload). */
    actionRef?: string;
    /**
     * W17: optional Host-local addon rendered under this field. The value is a
     * registered custom-component key (e.g. cron-preview). Not a protocol
     * control type.
     */
    afterComponent?: string;
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
export declare function isWhitelistedFormControl(type: string): type is FormControlType;
export declare function wireKindOf(field: FormControlField): WireKind;
/** Coerces a raw control value to its wire kind; defaultValue applies when raw is empty. */
export declare function coerceFieldValue(field: FormControlField, raw: unknown): unknown;
/** Fails closed when defaultValue does not match the field's wire type. */
export declare function validateDefaultValue(field: FormControlField): FormControlGateError | null;
/** Gates a control set against page meta (P-005 / frozen capability rules). */
export declare function checkFormCapabilities(meta: FormControlMeta, fields: FormControlField[]): FormControlGateError[];
/** Validates an arbitrary page-meta value without assuming record types. */
export declare function checkFormCapabilitiesRaw(metaValue: unknown, fields: FormControlField[]): FormControlGateError[];
/** One submit-time field validation failure (GOAL-014 D-002 §3). */
export interface FieldValidationError {
    field: string;
    code: "REQUIRED" | "PATTERN" | "MIN_LENGTH" | "MAX_LENGTH" | "MIN_VALUE" | "MAX_VALUE";
    /** Stable i18n key (form.validation.*); message is the en fallback. */
    messageKey: string;
    message: string;
}
/**
 * Validates form values against each field's declared constraints (GOAL-014
 * D-002 §3.1). Pure function: returns all failures (not just the first) so
 * the host can inline every field error. Boolean fields (switch/checkbox)
 * never fail REQUIRED — their wire kind is boolean and the toggle has an
 * explicit state.
 */
export declare function validateFieldValues(fields: FormControlField[], values: Record<string, unknown>): FieldValidationError[];
