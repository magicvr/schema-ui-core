import { useState } from "react";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { resolveTextProp, type MessageParams } from "@/i18n/catalog";
import { useTranslate } from "@/i18n/runtime";
import { cn } from "@/lib/utils";
import type { UploadableFile } from "@/protocol/conformance/upload-orchestration";
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
}

type FieldTranslator = (key: string, params?: MessageParams, literalFallback?: string) => string;

function optionList(field: FormControlField, t: FieldTranslator): Array<{ value: string; label: string }> {
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
  "h-9 w-full rounded-md border border-input bg-background px-3 text-sm outline-none transition-colors focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-50";

function BaseInput({
  id,
  label,
  value,
  onChange,
  type,
  disabled,
  readOnly,
  placeholder,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  type: "text" | "password";
  disabled?: boolean;
  readOnly?: boolean;
  placeholder?: string;
}) {
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs text-muted-foreground">
        {label}
      </Label>
      <Input
        id={id}
        type={type}
        value={value}
        disabled={disabled}
        readOnly={readOnly}
        placeholder={placeholder}
        // W4 P2-2: password fields in schema-driven forms (change/reset) must
        // not be auto-filled from a saved login password — declare a new
        // password context so the browser suggests a fresh one.
        autoComplete={type === "password" ? "new-password" : undefined}
        onChange={(event) => onChange(event.target.value)}
      />
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
  t,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
  t: FieldTranslator;
}) {
  const options = optionList(field, t);
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
        <legend className="text-xs font-medium text-muted-foreground">{label}</legend>
        <div className="space-y-1">
          {options.map((option) => (
            <label key={option.value} className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
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
      <Label htmlFor={id} className="text-xs text-muted-foreground">
        {label}
      </Label>
      <select
        id={id}
        value={value === undefined || value === null ? "" : String(value)}
        disabled={disabled || readOnly}
        onChange={(event) => onChange(event.target.value)}
        className={controlClass}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
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
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
  t: FieldTranslator;
}) {
  const current = value === undefined || value === null ? "" : String(value);
  return (
    <fieldset className="space-y-1.5" id={id}>
      <legend className="text-xs font-medium text-muted-foreground">{label}</legend>
      <div className="space-y-1">
        {optionList(field, t).map((option) => (
          <label key={option.value} className="flex items-center gap-2 text-sm">
            <input
              type="radio"
              name={id}
              value={option.value}
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
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
  t: FieldTranslator;
}) {
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
      <legend className="text-xs font-medium text-muted-foreground">{label}</legend>
      <div className="space-y-1">
        {optionList(field, t).map((option) => (
          <label key={option.value} className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
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
    <div className="flex items-center gap-2 text-sm">
      <input
        id={id}
        type="checkbox"
        checked={checked}
        disabled={disabled || readOnly}
        onChange={(event) => onChange(event.target.checked)}
      />
      <Label htmlFor={id}>{label}</Label>
      {field.type === "switch" ? (
        <span className="text-xs text-muted-foreground">{t("feedback.switchMarker")}</span>
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
  placeholder,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
  readOnly?: boolean;
  placeholder?: string;
}) {
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs text-muted-foreground">
        {label}
      </Label>
      <Textarea
        id={id}
        value={value}
        disabled={disabled}
        readOnly={readOnly}
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
  min?: number;
  max?: number;
  step?: number;
}) {
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs text-muted-foreground">
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
  min,
  max,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
  readOnly?: boolean;
  min?: string;
  max?: string;
}) {
  return (
    <div className="space-y-1.5">
      <Label htmlFor={id} className="text-xs text-muted-foreground">
        {label}
      </Label>
      <Input
        id={id}
        type="date"
        value={value}
        disabled={disabled}
        readOnly={readOnly}
        min={min}
        max={max}
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
}: {
  id: string;
  label: string;
  value: { start: string; end: string };
  onChange: (value: { start: string; end: string }) => void;
  disabled?: boolean;
  readOnly?: boolean;
}) {
  return (
    <fieldset className="space-y-1.5" id={id}>
      <legend className="text-xs font-medium text-muted-foreground">{label}</legend>
      <div className="flex items-center gap-2">
        <Input
          type="date"
          value={value.start}
          disabled={disabled}
          readOnly={readOnly}
          aria-label={`${label} start`}
          onChange={(event) => onChange({ ...value, start: event.target.value })}
        />
        <span className="text-xs text-muted-foreground">–</span>
        <Input
          type="date"
          value={value.end}
          disabled={disabled}
          readOnly={readOnly}
          aria-label={`${label} end`}
          onChange={(event) => onChange({ ...value, end: event.target.value })}
        />
      </div>
    </fieldset>
  );
}

function UploadField({
  id,
  label,
  field,
  value,
  onChange,
  disabled,
  readOnly,
  onUpload,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
  readOnly?: boolean;
  onUpload?: (field: FormControlField, files: UploadableFile[]) => Promise<unknown>;
}) {
  const [status, setStatus] = useState<{ kind: "idle" | "uploading" | "error"; message?: string }>({
    kind: "idle",
  });
  const display =
    value === undefined || value === null || value === ""
      ? ""
      : Array.isArray(value)
        ? value.join(", ")
        : String(value);
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
        <span className="text-xs text-muted-foreground">Value: {display}</span>
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
}) {
  const id = `${idPrefix ?? "field"}-${field.id}`;
  const label = resolveTextProp(
    field as unknown as Record<string, unknown>,
    "labelKey",
    "label",
    t,
    field.type,
  );
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
    case "input":
      return (
        <BaseInput
          id={id}
          label={label}
          value={String(value)}
          type="text"
          disabled={disabled}
          readOnly={readOnly}
          onChange={emitChange}
          placeholder={placeholder}
        />
      );
    case "password":
      return (
        <BaseInput
          id={id}
          label={label}
          value={String(value)}
          type="password"
          disabled={disabled}
          readOnly={readOnly}
          onChange={emitChange}
          placeholder={placeholder}
        />
      );
    case "inputNumber":
      return (
        <NumberField
          id={id}
          label={label}
          value={typeof value === "number" ? value : undefined}
          disabled={disabled}
          readOnly={readOnly}
          min={field.min}
          max={field.max}
          step={field.step}
          onChange={emitChange}
        />
      );
    case "datePicker":
      return (
        <DateField
          id={id}
          label={label}
          value={String(value)}
          disabled={disabled}
          readOnly={readOnly}
          onChange={emitChange}
        />
      );
    case "dateRangePicker": {
      const range =
        isDateRangeValue(value) ? value : { start: "", end: "" };
      return (
        <DateRangeField
          id={id}
          label={label}
          value={range}
          disabled={disabled}
          readOnly={readOnly}
          onChange={emitChange}
        />
      );
    }
    case "select":
      return (
        <SelectField
          id={id}
          label={label}
          field={field}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          t={t}
          onChange={emitChange}
        />
      );
    case "radio":
      return (
        <RadioField
          id={id}
          label={label}
          field={field}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          t={t}
          onChange={emitChange}
        />
      );
    case "checkboxGroup":
    case "cascader":
      return (
        <CheckboxGroupField
          id={id}
          label={label}
          field={field}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          t={t}
          onChange={emitChange}
        />
      );
    case "switch":
    case "checkbox":
      return (
        <BooleanField
          id={id}
          label={label}
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
          label={label}
          value={String(value)}
          disabled={disabled}
          readOnly={readOnly}
          placeholder={placeholder}
          onChange={emitChange}
        />
      );
    case "richText":
      return (
        <TextAreaField
          id={id}
          label={label}
          value={String(value)}
          disabled={disabled}
          placeholder={placeholder ?? "Markdown"}
          onChange={emitChange}
        />
      );
    case "upload":
      return (
        <UploadField
          id={id}
          label={label}
          field={field}
          value={value}
          disabled={disabled}
          readOnly={readOnly}
          onUpload={onUpload}
          onChange={emitChange}
        />
      );
  }
  return null;
};
  return (
    <div className="min-w-0">
      {renderControl()}
      {error !== undefined ? (
        <p role="alert" className="mt-1 text-sm text-destructive">
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
}: FormControlsProps) {
  const t = useTranslate();
  // GOAL-014 D-002 §4: single-column is the default (industry convention for
  // modal forms); schema may opt into a responsive multi-column grid. Mobile
  // stays single-column regardless.
  const cols =
    columns !== undefined && columns > 1 ? Math.min(Math.max(Math.floor(columns), 1), 4) : undefined;
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
      className={cn("grid gap-4 grid-cols-1", gridClass)}
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
        />
      ))}
    </div>
  );
}
