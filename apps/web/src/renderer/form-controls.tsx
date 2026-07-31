import { cn } from "@/lib/utils";
import {
  coerceFieldValue,
  type FormControlField,
} from "@/renderer/form-controls";

export interface FormControlsProps {
  fields: FormControlField[];
  values: Record<string, unknown>;
  onChange: (id: string, value: unknown) => void;
  disabled?: boolean;
  idPrefix?: string;
}

function optionList(field: FormControlField): Array<{ value: string; label: string }> {
  return (field.options ?? []).map((option) => ({
    value: option.value,
    label: option.label ?? option.value,
  }));
}

function displayValue(field: FormControlField, value: unknown): unknown {
  return coerceFieldValue(field, value);
}

function BaseInput({
  id,
  label,
  value,
  onChange,
  type,
  disabled,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  type: "text" | "password";
  disabled?: boolean;
}) {
  return (
    <label className="block space-y-1" htmlFor={id}>
      <span className="text-xs font-medium text-muted-foreground">{label}</span>
      <input
        id={id}
        type={type}
        value={value}
        disabled={disabled}
        onChange={(event) => onChange(event.target.value)}
        className="h-9 w-full rounded-md border border-input bg-background px-3 text-sm outline-none transition-colors focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-50"
      />
    </label>
  );
}

function SelectField({
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
  const options = optionList(field);
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
      <fieldset className="space-y-1">
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
    <label className="block space-y-1" htmlFor={id}>
      <span className="text-xs font-medium text-muted-foreground">{label}</span>
      <select
        id={id}
        value={value === undefined || value === null ? "" : String(value)}
        disabled={disabled}
        onChange={(event) => onChange(event.target.value)}
        className="h-9 w-full rounded-md border border-input bg-background px-3 text-sm outline-none transition-colors focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-50"
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </label>
  );
}

function RadioField({
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
  const current = value === undefined || value === null ? "" : String(value);
  return (
    <fieldset className="space-y-1" id={id}>
      <legend className="text-xs font-medium text-muted-foreground">{label}</legend>
      <div className="space-y-1">
        {optionList(field).map((option) => (
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
}: {
  id: string;
  label: string;
  field: FormControlField;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
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
    <fieldset className="space-y-1" id={id}>
      <legend className="text-xs font-medium text-muted-foreground">{label}</legend>
      <div className="space-y-1">
        {optionList(field).map((option) => (
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
  return (
    <label className="flex items-center gap-2 text-sm" htmlFor={id}>
      <input
        id={id}
        type="checkbox"
        checked={checked}
        disabled={disabled}
        onChange={(event) => onChange(event.target.checked)}
      />
      {label}
      {field.type === "switch" ? <span className="text-xs text-muted-foreground">(switch)</span> : null}
    </label>
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
    <label className="block space-y-1" htmlFor={id}>
      <span className="text-xs font-medium text-muted-foreground">{label}</span>
      <textarea
        id={id}
        value={value}
        disabled={disabled}
        placeholder={placeholder}
        onChange={(event) => onChange(event.target.value)}
        rows={4}
        className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none transition-colors focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-50"
      />
    </label>
  );
}

function FieldControl({
  field,
  values,
  onChange,
  disabled,
}: {
  field: FormControlField;
  values: Record<string, unknown>;
  onChange: (id: string, value: unknown) => void;
  disabled?: boolean;
}) {
  const id = `field-${field.id}`;
  const label = field.label ?? field.type;
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
        />
      );
    case "select":
      return (
        <SelectField
          id={id}
          label={label}
          field={field}
          value={value}
          disabled={disabled}
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
          placeholder="Markdown"
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
}: FormControlsProps) {
  return (
    <div className={cn("grid gap-4", fields.length > 1 && "sm:grid-cols-2")}>
      {fields.map((field) => (
        <FieldControl
          key={field.id}
          field={field}
          values={values}
          onChange={onChange}
          disabled={disabled}
        />
      ))}
    </div>
  );
}
