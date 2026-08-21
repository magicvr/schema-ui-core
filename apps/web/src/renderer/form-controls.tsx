import { useEffect, useState, type ReactNode } from "react";

import { ChevronDown, Search, X } from "lucide-react";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { resolveTextProp, type MessageParams } from "@/i18n/catalog";
import { useTranslate } from "@/i18n/runtime";
import { cn } from "@/lib/utils";
import type { UploadableFile } from "@/protocol/conformance/upload-orchestration";

import { getCustomComponent } from "@/renderer/custom-components";
import {
  coerceFieldValue,
  type FormControlField,
} from "@/renderer/form-controls";

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

/**
 * W11 · U-01/U-02: loads a form field's dynamic option source. Returns
 * null while no source is declared or the fetch is in flight (callers fall
 * back to static options), and an option list once resolved. Invalid sources
 * and failed fetches fail closed to an empty list.
 */
function useDynamicOptions(
  field: FormControlField,
  fetcher: typeof fetch | undefined,
): Array<{ value: string; label: string }> | null {
  const [items, setItems] = useState<Array<{ value: string; label: string }> | null>(null);

  useEffect(() => {
    const source = field.optionsSource;
    if (source === undefined) {
      setItems(null);
      return;
    }
    // Options-source urls are same-origin single-slash paths like table data
    // sources, but MAY carry a query string; fragments and scheme/host forms
    // are rejected (fail closed, no fetch).
    if (!/^\/(?!\/)[^\s\#]*$/.test(source.url)) {
      setItems([]);
      return;
    }
    let cancelled = false;
    const url = new URL(source.url, window.location.origin);
    if (source.params !== undefined) {
      for (const [key, value] of Object.entries(source.params)) {
        if (value === null || value === undefined) {
          url.searchParams.delete(key);
        } else {
          url.searchParams.set(key, String(value));
        }
      }
    }
    (fetcher ?? globalThis.fetch)(url.toString(), { headers: { Accept: "application/json" } })
      .then((response) => (response.ok ? response.json() : null))
      .then((body: unknown) => {
        if (cancelled) {
          return;
        }
        const raw = Array.isArray(body)
          ? body
          : (body as { items?: unknown[] } | null)?.items;
        const mapped: Array<{ value: string; label: string }> = [];
        for (const entry of raw ?? []) {
          if (typeof entry !== "object" || entry === null || Array.isArray(entry)) {
            continue;
          }
          const record = entry as Record<string, unknown>;
          const value = record[source.valueField];
          if (typeof value !== "string" || value === "") {
            continue;
          }
          const rawLabel = record[source.labelField];
          const label = typeof rawLabel === "string" && rawLabel !== "" ? rawLabel : value;
          mapped.push({ value, label });
        }
        setItems(mapped);
      })
      .catch(() => {
        if (!cancelled) {
          setItems([]);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [field.optionsSource, fetcher]);

  return items;
}

export function optionList(field: FormControlField, t: FieldTranslator): Array<{ value: string; label: string }> {
  return (field.options ?? []).map((option) => ({
    value: option.value,
    label: resolveTextProp(
      option as unknown as Record<string, unknown>,
      "labelKey",
      "label",
      t,
      String(option.value),
    ),
  }));
}

function isDateRangeValue(value: unknown): value is { start: string; end: string } {
  return (
    typeof value === "object" &&
    value !== null &&
    !Array.isArray(value) &&
    typeof (value as { start?: unknown }).start === "string" &&
    typeof (value as { end?: unknown }).end === "string"
  );
}

function displayValue(field: FormControlField, value: unknown): unknown {
  return coerceFieldValue(field, value);
}

const controlClass =
  "h-9 w-full rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20 disabled:cursor-not-allowed disabled:opacity-50";

function BaseInput({
  id,
  label,
  value,
  onChange,
  type,
  disabled,
  readOnly,
  placeholder,
  required,
  describedBy,
  searchMode,
  paired,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  type: "text" | "password";
  disabled?: boolean;
  readOnly?: boolean;
  placeholder?: string;
  required?: boolean;
  describedBy?: string;
  /** A-003: search keyword input — magnifier prefix + one-click clear. */
  searchMode?: boolean;
  /** W13 T-03: inside the [input+button] group — square right corner. */
  paired?: boolean;
}) {
  const t = useTranslate();
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs font-medium text-muted-foreground/80 select-none">
        {label}
      </Label>
      <div className="relative">
        {searchMode === true ? (
          <Search
            aria-hidden="true"
            className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground/50"
          />
        ) : null}
        <Input
          id={id}
          type={type}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          placeholder={placeholder}
          aria-required={required === true ? true : undefined}
          aria-describedby={describedBy}
          // W4 P2-2: password fields in schema-driven forms (change/reset) must
          // not be auto-filled from a saved login password — declare a new
          // password context so the browser suggests a fresh one.
          autoComplete={type === "password" ? "new-password" : undefined}
          onChange={(event) => onChange(event.target.value)}
          className={cn(
            searchMode === true ? "pl-9" : undefined,
            searchMode === true && value !== "" ? "pr-8" : undefined,
            paired === true ? "rounded-r-none" : undefined,
          )}
        />
        {searchMode === true && value !== "" ? (
          <button
            type="button"
            aria-label={t("feedback.clearSearch")}
            className="absolute right-2.5 top-1/2 inline-flex size-4.5 -translate-y-1/2 cursor-pointer items-center justify-center rounded-full text-muted-foreground/50 transition-colors hover:bg-muted hover:text-foreground active:scale-95"
            onClick={() => onChange("")}
          >
            <X aria-hidden="true" className="size-3" />
          </button>
        ) : null}
      </div>
    </div>
  );
}

function SelectField({
  id,
  label,
  field,
  value,
  onChange,
  disabled,
  readOnly,
  required,
  describedBy,
  t,
  fetcher,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
  required?: boolean;
  describedBy?: string;
  t: FieldTranslator;
  fetcher?: typeof fetch;
}) {
  // W11 · U-01/U-02: dynamic options win over static options once loaded.
  const dynamic = useDynamicOptions(field, fetcher);
  const options = dynamic !== null ? dynamic : optionList(field, t);
  if (field.mode === "multiple") {
    const selected = Array.isArray(value) ? value.map(String) : [];
    const toggle = (option: string) => {
      onChange(
        selected.includes(option)
          ? selected.filter((entry) => entry !== option)
          : [...selected, option],
      );
    };
    return (
      <fieldset className="space-y-1.5">
        <legend className="text-xs font-medium text-muted-foreground/80 select-none">{label}</legend>
        <div className="space-y-1.5">
          {options.map((option) => (
            <label key={option.value} className="flex cursor-pointer items-center gap-2 text-sm text-foreground/90 transition-colors hover:text-foreground">
              <input
                type="checkbox"
                className="size-4 cursor-pointer rounded border-input text-primary accent-primary transition-colors focus:ring-2 focus:ring-ring/20"
                checked={selected.includes(option.value)}
                disabled={disabled || readOnly}
                onChange={() => toggle(option.value)}
              />
              {option.label}
            </label>
          ))}
        </div>
      </fieldset>
    );
  }
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs font-medium text-muted-foreground/80 select-none">
        {label}
      </Label>
      {/* W8 P3 (GOAL-009): control-level color-scheme so the native option
          popup renders dark immediately in dark mode (root cascades too). */}
      <div className="relative">
        <select
          id={id}
          value={value === undefined || value === null ? "" : String(value)}
          disabled={disabled || readOnly}
          aria-required={required === true ? true : undefined}
          aria-describedby={describedBy}
          onChange={(event) => onChange(event.target.value)}
          className={cn(
            controlClass,
            "appearance-none pr-8 cursor-pointer scheme-light dark:scheme-dark",
            value === "" || value === undefined || value === null
              ? "text-muted-foreground"
              : "text-foreground font-normal",
          )}
        >
          {options.length === 0 || options[0]?.value !== "" ? (
            <option value="" className="bg-background text-muted-foreground">
              {t("feedback.selectPlaceholder")}
            </option>
          ) : null}
          {options.map((option) => (
            <option key={option.value} value={option.value} className="bg-background text-foreground">
              {option.label}
            </option>
          ))}
        </select>
        <ChevronDown
          aria-hidden="true"
          className="pointer-events-none absolute right-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground/60 transition-transform"
        />
      </div>
    </div>
  );
}

function RadioField({
  id,
  label,
  field,
  value,
  onChange,
  disabled,
  readOnly,
  t,
  fetcher,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
  t: FieldTranslator;
  fetcher?: typeof fetch;
}) {
  const dynamic = useDynamicOptions(field, fetcher);
  const options = dynamic !== null ? dynamic : optionList(field, t);
  const current = value === undefined || value === null ? "" : String(value);
  return (
    <fieldset className="space-y-1.5" id={id}>
      <legend className="text-xs font-medium text-muted-foreground/80 select-none">{label}</legend>
      <div className="space-y-1.5">
        {options.map((option) => (
          <label key={option.value} className="flex cursor-pointer items-center gap-2 text-sm text-foreground/90 transition-colors hover:text-foreground">
            <input
              type="radio"
              name={id}
              value={option.value}
              className="size-4 cursor-pointer accent-primary"
              checked={current === option.value}
              disabled={disabled || readOnly}
              onChange={() => onChange(option.value)}
            />
            {option.label}
          </label>
        ))}
      </div>
    </fieldset>
  );
}

function CheckboxGroupField({
  id,
  label,
  field,
  value,
  onChange,
  disabled,
  readOnly,
  t,
  fetcher,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
  t: FieldTranslator;
  fetcher?: typeof fetch;
}) {
  const dynamic = useDynamicOptions(field, fetcher);
  const options = dynamic !== null ? dynamic : optionList(field, t);
  const selected = Array.isArray(value) ? value.map(String) : [];
  const toggle = (option: string) => {
    onChange(
      selected.includes(option)
        ? selected.filter((entry) => entry !== option)
        : [...selected, option],
    );
  };
  return (
    <fieldset className="space-y-1.5" id={id}>
      <legend className="text-xs font-medium text-muted-foreground/80 select-none">{label}</legend>
      <div className="space-y-1.5">
        {options.map((option) => (
          <label key={option.value} className="flex cursor-pointer items-center gap-2 text-sm text-foreground/90 transition-colors hover:text-foreground">
            <input
              type="checkbox"
              checked={selected.includes(option.value)}
              disabled={disabled || readOnly}
              onChange={() => toggle(option.value)}
              className="size-4 cursor-pointer rounded border-input text-primary accent-primary transition-colors focus:ring-2 focus:ring-ring/20"
            />
            {option.label}
          </label>
        ))}
      </div>
    </fieldset>
  );
}

function BooleanField({
  id,
  label,
  field,
  value,
  onChange,
  disabled,
  readOnly,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
}) {
  const checked = value === true;
  const t = useTranslate();
  return (
    <div className="flex items-center gap-2.5 pt-1 text-sm">
      <input
        id={id}
        type="checkbox"
        checked={checked}
        disabled={disabled || readOnly}
        onChange={(event) => onChange(event.target.checked)}
        className="size-4 cursor-pointer rounded border-input text-primary accent-primary transition-colors focus:ring-2 focus:ring-ring/20"
      />
      <Label htmlFor={id} className="cursor-pointer font-normal text-foreground select-none">{label}</Label>
      {field.type === "switch" ? (
        <span className="rounded-full bg-muted/60 px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground">
          {t("feedback.switchMarker")}
        </span>
      ) : null}
    </div>
  );
}

function TextAreaField({
  id,
  label,
  value,
  onChange,
  disabled,
  readOnly,
  required,
  describedBy,
  placeholder,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
  readOnly?: boolean;
  required?: boolean;
  describedBy?: string;
  placeholder?: string;
}) {
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs font-medium text-muted-foreground/80 select-none">
        {label}
      </Label>
      <Textarea
        id={id}
        value={value}
        disabled={disabled}
        readOnly={readOnly}
        aria-required={required === true ? true : undefined}
        aria-describedby={describedBy}
        placeholder={placeholder}
        onChange={(event) => onChange(event.target.value)}
        rows={4}
      />
    </div>
  );
}

function NumberField({
  id,
  label,
  value,
  onChange,
  disabled,
  readOnly,
  required,
  describedBy,
  min,
  max,
  step,
}: {
  id: string;
  label: string;
  value: number | undefined;
  onChange: (value: number | undefined) => void;
  disabled?: boolean;
  readOnly?: boolean;
  required?: boolean;
  describedBy?: string;
  min?: number;
  max?: number;
  step?: number;
}) {
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs font-medium text-muted-foreground/80 select-none">
        {label}
      </Label>
      <Input
        id={id}
        type="number"
        // An empty input renders blank, not "0"; a cleared number is submitted
        // as undefined (field omitted) so the backend default applies (D7).
        value={Number.isFinite(value) ? value : ""}
        disabled={disabled}
        readOnly={readOnly}
        aria-required={required === true ? true : undefined}
        aria-describedby={describedBy}
        min={min}
        max={max}
        step={step}
        onChange={(event) => {
          const next = event.target.valueAsNumber;
          onChange(Number.isFinite(next) ? next : undefined);
        }}
      />
    </div>
  );
}

function DateField({
  id,
  label,
  value,
  onChange,
  disabled,
  readOnly,
  required,
  describedBy,
  min,
  max,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
  readOnly?: boolean;
  required?: boolean;
  describedBy?: string;
  min?: string;
  max?: string;
}) {
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs font-medium text-muted-foreground/80 select-none">
        {label}
      </Label>
      <Input
        id={id}
        type="date"
        value={value}
        disabled={disabled}
        readOnly={readOnly}
        aria-required={required === true ? true : undefined}
        aria-describedby={describedBy}
        min={min}
        max={max}
        className="cursor-pointer scheme-light dark:scheme-dark"
        onChange={(event) => onChange(event.target.value)}
      />
    </div>
  );
}

function DateRangeField({
  id,
  label,
  value,
  onChange,
  disabled,
  readOnly,
  required,
  describedBy,
}: {
  id: string;
  label: string;
  value: { start: string; end: string };
  onChange: (value: { start: string; end: string }) => void;
  disabled?: boolean;
  readOnly?: boolean;
  required?: boolean;
  describedBy?: string;
}) {
  return (
    <fieldset className="space-y-1.5" id={id}>
      <legend className="text-xs font-medium text-muted-foreground/80 select-none">{label}</legend>
      <div className="flex items-center gap-2">
        <Input
          type="date"
          value={value.start}
          disabled={disabled}
          readOnly={readOnly}
          aria-required={required === true ? true : undefined}
          aria-describedby={describedBy}
          aria-label={`${label} start`}
          className="cursor-pointer scheme-light dark:scheme-dark"
          onChange={(event) => onChange({ ...value, start: event.target.value })}
        />
        <span className="text-xs font-medium text-muted-foreground/60 select-none">–</span>
        <Input
          type="date"
          value={value.end}
          disabled={disabled}
          readOnly={readOnly}
          aria-required={required === true ? true : undefined}
          aria-describedby={describedBy}
          aria-label={`${label} end`}
          className="cursor-pointer scheme-light dark:scheme-dark"
          onChange={(event) => onChange({ ...value, end: event.target.value })}
        />
      </div>
    </fieldset>
  );
}

/**
 * Same-origin path or http(s) URL — the only value shapes that may render as
 * an image preview. Everything else (bare ids, javascript:/data: URLs) stays
 * textual; a URL that fails to load falls back to the text too (onError).
 */
function isPreviewableUrl(value: string): boolean {
  if (/^\/(?!\/)[^\s\\]*$/.test(value)) {
    return true;
  }
  return /^https?:\/\//.test(value);
}

function UploadField({
  id,
  label,
  field,
  value,
  onChange,
  disabled,
  readOnly,
  required,
  describedBy,
  onUpload,
  removeLabel,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
  required?: boolean;
  describedBy?: string;
  onUpload?: (field: FormControlField, files: UploadableFile[]) => Promise<unknown>;
  /** W9 (GOAL-010): localized label for the per-field "remove image" button
   * (single uploads only; clears the field value back to ""). */
  removeLabel?: string;
}) {
  const [status, setStatus] = useState<{ kind: "idle" | "uploading" | "error"; message?: string }>({
    kind: "idle",
  });
  const [previewFailed, setPreviewFailed] = useState(false);
  const display =
    value === undefined || value === null || value === ""
      ? ""
      : Array.isArray(value)
        ? value.join(", ")
        : String(value);
  // W9 follow-up (user 2026-08-15): a new committed value (re-upload) replaces
  // a previously broken preview — reset the failure flag when the value changes.
  useEffect(() => {
    setPreviewFailed(false);
  }, [display]);
  const handleFiles = async (fileList: FileList | null) => {
    if (fileList === null || fileList.length === 0 || onUpload === undefined) {
      return;
    }
    const files: UploadableFile[] = [...fileList].map((file) => ({
      name: file.name,
      type: file.type,
      size: file.size,
      contentId: file.name,
      blob: file,
    }));
    setStatus({ kind: "uploading" });
    try {
      const result = await onUpload(field, files);
      setStatus({ kind: "idle" });
      onChange(result);
    } catch (error) {
      setStatus({ kind: "error", message: error instanceof Error ? error.message : String(error) });
    }
  };
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs text-muted-foreground">
        {label}
      </Label>
      <Input
        id={id}
        type="file"
        multiple={field.multiple === true}
        disabled={disabled || readOnly || status.kind === "uploading"}
        aria-required={required === true ? true : undefined}
        aria-describedby={describedBy}
        onChange={(event) => {
          void handleFiles(event.target.files).finally(() => {
            // Reset the input after the upload settles so re-selecting the same
            // file re-fires change (D7); a stale file reference would make a
            // retry after an error silently do nothing.
            event.target.value = "";
          });
        }}
      />
      {display !== "" ? (
        <div className="flex items-center gap-2">
          {isPreviewableUrl(display) && !previewFailed ? (
            // User-facing image preview (W9 follow-up): brand assets and other
            // URL values render as a thumbnail; non-image URLs (CSV/pdf ids,
            // bare ids) fail to load and fall back to the value text.
            <img
              src={display}
              alt=""
              className="h-10 w-auto max-w-full shrink-0 rounded border border-border object-contain"
              onError={() => setPreviewFailed(true)}
            />
          ) : (
            <span className="min-w-0 truncate text-xs text-muted-foreground">Value: {display}</span>
          )}
          {field.multiple !== true && removeLabel !== undefined && !disabled && !readOnly ? (
            <button
              type="button"
              className="shrink-0 text-xs text-destructive underline underline-offset-2"
              onClick={() => onChange("")}
            >
              {removeLabel}
            </button>
          ) : null}
        </div>
      ) : null}
      {status.kind === "error" ? (
        <span role="alert" className="text-xs text-destructive">
          {status.message}
        </span>
      ) : null}
    </div>
  );
}

function FieldControl({
  field,
  values,
  onChange,
  disabled,
  readOnly,
  idPrefix,
  onUpload,
  t,
  error,
  fetcher,
  searchMode,
  searchButtonSlot,
}: {
  field: FormControlField;
  values: Record<string, unknown>;
  onChange: (id: string, value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
  idPrefix?: string;
  onUpload?: (field: FormControlField, files: UploadableFile[]) => Promise<unknown>;
  t: (key: string, params?: MessageParams, literalFallback?: string) => string;
  error?: string;
  fetcher?: typeof fetch;
  /** A-003: search presentation for the keyword input (prefix icon + clear). */
  searchMode?: boolean;
  /** A-003 pairing rule: submit button rendered adjacent to the input. */
  searchButtonSlot?: ReactNode;
  /** W13 T-03: input is part of the [input+button] search group — square inner corner. */
  paired?: boolean;
}) {
  const id = `${idPrefix ?? "field"}-${field.id}`;
  const label = resolveTextProp(
    field as unknown as Record<string, unknown>,
    "labelKey",
    "label",
    t,
    field.type,
  );
  const required = field.required === true;
  const requiredLabel = required ? label + " *" : label;
  const errorId = error !== undefined ? `${id}-error` : undefined;
  const placeholder = resolveTextProp(
    field as unknown as Record<string, unknown>,
    "placeholderKey",
    "placeholder",
    t,
    undefined,
  );
  const value = displayValue(field, values[field.id]);
  // ADR-0040: a readOnly field never accepts user edits — user-originated
  // change events are dropped (reactions / recordSource / Host seeding still
  // write through the form state, which is not routed via onChange).
  const emitChange = (next: unknown) => {
    if (field.readOnly === true) {
      return;
    }
    onChange(field.id, next);
  };

  const renderControl = (): ReactNode => {
  switch (field.type) {
    case "input": {
      const control = (
        <BaseInput
          id={id}
          label={requiredLabel}
          value={String(value)}
          type="text"
          disabled={disabled}
          readOnly={readOnly}
          required={required}
          describedBy={errorId}
          onChange={emitChange}
          placeholder={placeholder}
          searchMode={searchMode}
          paired={searchMode === true && searchButtonSlot !== undefined}
        />
      );
      // W13 T-03 (user 2026-08-16): the keyword input and its search button
      // are ONE semantic component — zero-gap flex-nowrap so the pair stays
      // stuck together (贴在一起) and can never be wrapped onto separate
      // lines at any page width. The button overlaps the input's right edge
      // (-ml-px) and the inner corners are squared off on both sides.
      return searchMode === true && searchButtonSlot !== undefined ? (
        <div className="flex flex-nowrap items-end">
          <div className="min-w-0 flex-1">{control}</div>
          {searchButtonSlot}
        </div>
      ) : (
        control
      );
    }
    case "password":
      return (
        <BaseInput
          id={id}
          label={requiredLabel}
          value={String(value)}
          type="password"
          disabled={disabled}
          readOnly={readOnly}
          required={required}
          describedBy={errorId}
          onChange={emitChange}
          placeholder={placeholder}
        />
      );
    case "inputNumber":
      return (
        <>
          <NumberField
            id={id}
            label={requiredLabel}
            value={typeof value === "number" ? value : undefined}
            disabled={disabled}
            readOnly={readOnly}
            required={required}
            describedBy={errorId}
            min={field.min}
            max={field.max}
            step={field.step}
            onChange={emitChange}
          />
          {/* W16-F04: wallet adjustment warning for negative or large deltas. */}
          {field.id === "amountDelta" && typeof value === "number" && (value < 0 || Math.abs(value) > 100000) ? (
            <p className="text-xs text-warning" data-wallet-adjust-warning>
              {t("schema.wallet.adjustWarning")}
            </p>
          ) : null}
        </>
      );
    case "datePicker":
      return (
        <DateField
          id={id}
          label={requiredLabel}
          value={String(value)}
          disabled={disabled}
          readOnly={readOnly}
          required={required}
          describedBy={errorId}
          onChange={emitChange}
        />
      );
    case "dateRangePicker": {
      const range =
        isDateRangeValue(value) ? value : { start: "", end: "" };
      return (
        <DateRangeField
          id={id}
          label={requiredLabel}
          value={range}
          disabled={disabled}
          readOnly={readOnly}
          required={required}
          describedBy={errorId}
          onChange={emitChange}
        />
      );
    }
    case "select":
      return (
        <SelectField
          id={id}
          label={requiredLabel}
          field={field}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          required={required}
          describedBy={errorId}
          t={t}
          fetcher={fetcher}
          onChange={emitChange}
        />
      );
    case "radio":
      return (
        <RadioField
          id={id}
          label={requiredLabel}
          field={field}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          t={t}
          fetcher={fetcher}
          onChange={emitChange}
        />
      );
    case "checkboxGroup":
    case "cascader":
      return (
        <CheckboxGroupField
          id={id}
          label={requiredLabel}
          field={field}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          t={t}
          fetcher={fetcher}
          onChange={emitChange}
        />
      );
    case "switch":
    case "checkbox":
      return (
        <BooleanField
          id={id}
          label={requiredLabel}
          field={field}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          onChange={emitChange}
        />
      );
    case "textarea":
      return (
        <TextAreaField
          id={id}
          label={requiredLabel}
          value={String(value)}
          disabled={disabled}
          readOnly={readOnly}
          required={required}
          describedBy={errorId}
          placeholder={placeholder}
          onChange={emitChange}
        />
      );
    case "richText":
      return (
        <TextAreaField
          id={id}
          label={requiredLabel}
          value={String(value)}
          disabled={disabled}
          required={required}
          describedBy={errorId}
          placeholder={placeholder ?? "Markdown"}
          onChange={emitChange}
        />
      );
    case "upload":
      return (
        <UploadField
          id={id}
          label={requiredLabel}
          field={field}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          required={required}
          describedBy={errorId}
          onUpload={onUpload}
          removeLabel={t("form.upload.remove")}
          onChange={emitChange}
        />
      );
  }
  return null;
};
  const AfterComponent =
    typeof field.afterComponent === "string" && field.afterComponent !== ""
      ? getCustomComponent(field.afterComponent)
      : null;
  return (
    <div className="min-w-0">
      {renderControl()}
      {AfterComponent !== null ? (
        <AfterComponent
          node={{
            type: "custom",
            component: field.afterComponent ?? "",
            props: { bindValue: String(value) },
          }}
          context={{}}
        />
      ) : null}
      {error !== undefined ? (
        <p id={errorId} role="alert" className="mt-1 text-sm text-destructive">
          {error}
        </p>
      ) : null}
    </div>
  );
}

export function FormControls({
  fields,
  values,
  onChange,
  disabled = false,
  fieldDisabled,
  idPrefix,
  onUpload,
  fieldErrors,
  columns,
  fetcher,
  searchMode: searchModeProp,
  actionSlot,
  searchButtonSlot,
}: FormControlsProps) {
  const t = useTranslate();
  // GOAL-014 D-002 §4: single-column is the default (industry convention for
  // modal forms); schema may opt into a responsive multi-column grid. Mobile
  // stays single-column regardless.
  const cols =
    columns !== undefined && columns > 1 ? Math.min(Math.max(Math.floor(columns), 1), 4) : undefined;
  const searchMode = searchModeProp === true;
  // F-005 (A-003): Tailwind JIT cannot extract dynamically concatenated
  // class names; use a static lookup so sm:grid-cols-2/3/4 are real utilities.
  const GRID_COL_CLASSES: Record<number, string> = {
    2: "sm:grid-cols-2",
    3: "sm:grid-cols-3",
    4: "sm:grid-cols-4",
  };
  const gridClass = cols !== undefined ? GRID_COL_CLASSES[cols] : undefined;
  return (
    <div
      data-form-controls="design-system"
      data-form-columns={cols !== undefined ? String(cols) : "1"}
      data-form-search-mode={searchMode ? "true" : undefined}
      className={
        searchMode
          ? "grid grid-cols-1 items-end gap-3 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5"
          : cn("grid gap-4 grid-cols-1", gridClass)
      }
    >
      {fields.map((field) => (
        <FieldControl
          key={field.id}
          field={field}
          values={values}
          onChange={onChange}
          disabled={disabled || (fieldDisabled?.(field.id) ?? false)}
          // ADR-0040: readOnly fields render non-editable; their values stay
          // in the form state and the submit projection.
          readOnly={field.readOnly === true}
          idPrefix={idPrefix}
          onUpload={onUpload}
          t={t}
          error={fieldErrors?.[field.id]}
          fetcher={fetcher}
          searchMode={searchMode}
          searchButtonSlot={searchButtonSlot}
          paired={searchMode === true && searchButtonSlot !== undefined}
        />
      ))}
      {actionSlot !== undefined ? actionSlot : null}
    </div>
  );
}
