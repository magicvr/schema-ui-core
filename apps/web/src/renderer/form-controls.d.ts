import { type ReactNode } from "react";
import { type MessageParams } from "@/i18n/catalog";
import type { UploadableFile } from "@/protocol/conformance/upload-orchestration";
import { type FormControlField } from "@/renderer/form-controls.types";
export interface FormControlsProps {
    fields: FormControlField[];
    values: Record<string, unknown>;
    onChange: (id: string, value: unknown) => void;
    disabled?: boolean;
    /** Per-field disabled override (R5 renderer reaction state). */
    fieldDisabled?: (id: string) => boolean;
    idPrefix?: string;
    /** Upload control transport (ADR-0012): validates + uploads + returns the field value. */
    onUpload?: (field: FormControlField, files: UploadableFile[]) => Promise<unknown>;
    /** GOAL-014 D-002 §4: inline field errors keyed by field id (submit-time
     * validation + server fieldErrors echo). */
    fieldErrors?: Record<string, string>;
    /** GOAL-014 D-002 §4: column count (default 1 = single-column layout).
     * >1 enables a responsive grid; the mobile layout stays single-column. */
    columns?: number;
    /** W11 · U-01/U-02: auth-aware transport for dynamic option sources
     * (optionsSource); defaults to globalThis.fetch. */
    fetcher?: typeof fetch;
    /**
     * A-003 (GOAL-013 audit response): search-mode presentation — compact
     * responsive auto-grid (1..5 columns), keyword input with search prefix
     * icon + clear affordance. Search schemas keep their exact JSON shape.
     */
    searchMode?: boolean;
    /** A-003: action cluster (Reset button) rendered inside the
     * search-mode grid so it aligns with the field row. */
    actionSlot?: ReactNode;
    /**
     * A-003 (user pairing rule): one search button rendered side-by-side with
     * EVERY keyword input field in search mode — as many buttons as there are
     * text inputs, each pair adjacent in the same grid cell.
     */
    searchButtonSlot?: ReactNode;
}
export type FieldTranslator = (key: string, params?: MessageParams, literalFallback?: string) => string;
export declare function optionList(field: FormControlField, t: FieldTranslator): Array<{
    value: string;
    label: string;
}>;
export declare function FormControls({ fields, values, onChange, disabled, fieldDisabled, idPrefix, onUpload, fieldErrors, columns, fetcher, searchMode: searchModeProp, actionSlot, searchButtonSlot, }: FormControlsProps): import("react").JSX.Element;
