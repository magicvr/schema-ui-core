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
  placeholder,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  type: "text" | "password";
  disabled?: boolean;
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
        placeholder={placeholder}
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
  t,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
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
                disabled={disabled}
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
        disabled={disabled}
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
  t,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
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
              disabled={disabled}
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
  t,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
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
              disabled={disabled}
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
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
}) {
  const checked = value === true;
  const t = useTranslate();
  return (
    <div className="flex items-center gap-2 text-sm">
      <input
        id={id}
        type="checkbox"
        checked={checked}
        disabled={disabled}
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
  placeholder,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
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
  min,
  max,
  step,
}: {
  id: string;
  label: string;
  value: number;
  onChange: (value: number) => void;
  disabled?: boolean;
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
        value={Number.isFinite(value) ? value : 0}
        disabled={disabled}
        min={min}
        max={max}
        step={step}
        onChange={(event) => {
          const next = event.target.valueAsNumber;
          onChange(Number.isFinite(next) ? next : 0);
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
  min,
  max,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
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
}: {
  id: string;
  label: string;
  value: { start: string; end: string };
  onChange: (value: { start: string; end: string }) => void;
  disabled?: boolean;
}) {
  return (
    <fieldset className="space-y-1.5" id={id}>
      <legend className="text-xs font-medium text-muted-foreground">{label}</legend>
      <div className="flex items-center gap-2">
        <Input
          type="date"
          value={value.start}
          disabled={disabled}
          aria-label={`${label} start`}
          onChange={(event) => onChange({ ...value, start: event.target.value })}
        />
        <span className="text-xs text-muted-foreground">–</span>
        <Input
          type="date"
          value={value.end}
          disabled={disabled}
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
  onUpload,
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
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
        disabled={disabled || status.kind === "uploading"}
        onChange={(event) => void handleFiles(event.target.files)}
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
  idPrefix,
  onUpload,
  t,
}: {
  field: FormControlField;
  values: Record<string, unknown>;
  onChange: (id: string, value: unknown) => void;
  disabled?: boolean;
  idPrefix?: string;
  onUpload?: (field: FormControlField, files: UploadableFile[]) => Promise<unknown>;
  t: (key: string, params?: MessageParams, literalFallback?: string) => string;
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

  switch (field.type) {
    case "input":
      return (
        <BaseInput
          id={id}
          label={label}
          value={String(value)}
          type="text"
          disabled={disabled}
          onChange={(next) => onChange(field.id, next)}
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
          onChange={(next) => onChange(field.id, next)}
          placeholder={placeholder}
        />
      );
    case "inputNumber":
      return (
        <NumberField
          id={id}
          label={label}
          value={Number(value)}
          disabled={disabled}
          min={field.min}
          max={field.max}
          step={field.step}
          onChange={(next) => onChange(field.id, next)}
        />
      );
    case "datePicker":
      return (
        <DateField
          id={id}
          label={label}
          value={String(value)}
          disabled={disabled}
          onChange={(next) => onChange(field.id, next)}
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
          onChange={(next) => onChange(field.id, next)}
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
          t={t}
          onChange={(next) => onChange(field.id, next)}
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
          t={t}
          onChange={(next) => onChange(field.id, next)}
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
          t={t}
          onChange={(next) => onChange(field.id, next)}
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
          onChange={(next) => onChange(field.id, next)}
        />
      );
    case "textarea":
      return (
        <TextAreaField
          id={id}
          label={label}
          value={String(value)}
          disabled={disabled}
          placeholder={placeholder}
          onChange={(next) => onChange(field.id, next)}
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
          onChange={(next) => onChange(field.id, next)}
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
          onUpload={onUpload}
          onChange={(next) => onChange(field.id, next)}
        />
      );
  }
  return null;
}

export function FormControls({
  fields,
  values,
  onChange,
  disabled = false,
  fieldDisabled,
  idPrefix,
  onUpload,
}: FormControlsProps) {
  const t = useTranslate();
  return (
    <div
      data-form-controls="design-system"
      className={cn("grid gap-4", fields.length > 1 && "sm:grid-cols-2")}
    >
      {fields.map((field) => (
        <FieldControl
          key={field.id}
          field={field}
          values={values}
          onChange={onChange}
          disabled={disabled || (fieldDisabled?.(field.id) ?? false)}
          idPrefix={idPrefix}
          onUpload={onUpload}
          t={t}
        />
      ))}
    </div>
  );
}
