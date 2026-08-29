/**
 * D-COMP component-format fixture adapter (schema-ui-docs v2.7.0).
 * Validates format wire types without coercion.
 */
export type ComponentFormat = "currency" | "percent" | "datetime" | string;
export type FormatResult = {
    ok: true;
    value: unknown;
} | {
    ok: false;
    code: "COMPONENT_DATA_TYPE_MISMATCH";
};
export declare function applyComponentFormat(format: ComponentFormat, value: unknown): FormatResult;
