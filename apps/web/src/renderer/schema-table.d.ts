import { type ResourceItem } from "@/renderer/resource";
import type { RenderTableNode } from "@/renderer/render.types";
/**
 * Default schema-driven table surface (R1 · GOAL-004 / D-004) + S4 CRUD wiring
 * + A-002 generic adapter (GOAL-010 S3).
 *
 * Renders a whitelisted table node's `props.columns` over its `props.dataSource`.
 * Fails closed when the node declares no columns, no (valid) data source, or a
 * response that violates the frozen rowKey invariant (F-001 / F-002).
 *
 * F-001 (I-010-001 v0.2.0 §2): `dataSource` must be a single-slash same-origin
 * path; invalid/absent sources render an observable fail-closed state and never
 * reach the (auth) fetcher.
 *
 * F-002 (I-010-001 v0.2.0 §3): `props.rowKey` (default `"id"`) is the direct
 * field name of each row's unique, non-empty scalar key (string / finite
 * number). A missing / empty / non-scalar / duplicate key stops rendering the
 * table data and forbids row actions and selection.
 *
 * S4 (GOAL-007 · I-007-003 v0.2.2 §9): when rendered inside a `RenderPage`
 * (which always wraps content in the SchemaCrudProvider), the table surfaces
 * `props.toolbar` (create trigger) and `props.actions` (row edit/delete),
 * participates in row selection for the page's recordView / edit-form prefill,
 * reads its query from the shared per-table query state (so the search-form
 * binding and the post-write reload both reach it), and registers its injected
 * fetcher so modal form submits / row actions use the same transport. All of
 * this is fixture-driven: no record-specific logic lives here.
 */
export interface SchemaTableProps {
    node: RenderTableNode;
    /** Injectable fetch (defaults to `globalThis.fetch`). */
    fetcher?: typeof fetch;
}
export interface SchemaTableColumnSpec {
    field: string;
    label?: string;
    sortable?: boolean;
    /** W4 · GOAL-005: single-line truncate + title affordance for long values. */
    truncate?: boolean;
    /** Column width hint (px number or CSS length). */
    width?: number | string;
    /** Minimum column width (px number or CSS length). */
    minWidth?: number | string;
    /** W16-F04: render a cent-valued number as a localized currency string. */
    format?: "currency";
    /** W16-F09: render this cell as a colored badge using the row field value. */
    badgeStyleField?: string;
}
/** Structured row predicate used by Schema actions; malformed predicates fail closed. */
export declare function rowActionDisabled(action: unknown, row: ResourceItem): boolean;
/** Extracts the table node's column specs, fail-closed on malformed entries. */
export declare function schemaTableColumns(node: RenderTableNode): SchemaTableColumnSpec[];
/**
 * Resolves the table node's list endpoint (F-001). v2.9 prefers the node-level
 * DataRef (ADR-0039, source:api url); the legacy props.dataSource string stays
 * supported. Returns null when absent or not a single-slash same-origin path;
 * the table then fails closed and never fetches (the fixture fallback was
 * removed in GOAL-010 S3).
 */
export declare function schemaTableDataSource(node: RenderTableNode): string | null;
/** v2.9 node-level DataRef params (ADR-0039 route bindings), else empty. */
export declare function schemaTableDataParams(node: RenderTableNode): Record<string, unknown> | undefined;
/** F-002: the direct field name used as each row's unique key (default "id"). */
export declare function schemaTableRowKey(node: RenderTableNode): string;
/** One select option of a schema-driven table filter. */
export interface SchemaTableFilterOptionSpec {
    value: string;
    label?: string;
    labelKey?: string;
}
/** A schema-driven table filter (table node `props.filters`; select-only). */
export interface SchemaTableFilterSpec {
    field: string;
    type: "select";
    label?: string;
    labelKey?: string;
    options: SchemaTableFilterOptionSpec[];
}
/**
 * Extracts the table node's filter specs (fail-closed on malformed entries).
 * Only `select` filters are supported; other/unknown types are dropped so a
 * schema typo never renders a broken control.
 */
export declare function schemaTableFilters(node: RenderTableNode): SchemaTableFilterSpec[];
/**
 * Windowed pager page list: 1 … current±1 … total, with gap markers for
 * long page runs (stable shape for tests and a11y).
 */
export declare function pagerPages(current: number, total: number): Array<number | "gap">;
export declare function SchemaTable({ node, fetcher }: SchemaTableProps): import("react").JSX.Element;
