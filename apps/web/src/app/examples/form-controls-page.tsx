import { useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { FormControls } from "@/renderer/form-controls.tsx";
import {
  checkFormCapabilitiesRaw,
  type FormControlField,
  type FormControlGateError,
} from "@/renderer/form-controls";

const ADVANCED_FIELDS: FormControlField[] = [
  { id: "name", label: "Name (input)", type: "input", defaultValue: "Acme Console" },
  { id: "kind", label: "Kind (select single)", type: "select", options: [
    { value: "standard", label: "Standard" },
    { value: "priority", label: "Priority" },
  ], defaultValue: "standard" },
  { id: "tags", label: "Tags (select multiple)", type: "select", mode: "multiple", options: [
    { value: "core", label: "Core" },
    { value: "beta", label: "Beta" },
    { value: "legacy", label: "Legacy" },
  ], defaultValue: ["core"] },
  { id: "notes", label: "Notes (textarea)", type: "textarea", defaultValue: "" },
  { id: "enabled", label: "Enabled (switch)", type: "switch", defaultValue: true },
  { id: "agree", label: "Agree (checkbox)", type: "checkbox", defaultValue: false },
  { id: "priority", label: "Priority (radio)", type: "radio", options: [
    { value: "low", label: "Low" },
    { value: "medium", label: "Medium" },
    { value: "high", label: "High" },
  ], defaultValue: "medium" },
  { id: "owners", label: "Owners (checkboxGroup)", type: "checkboxGroup", options: [
    { value: "alice", label: "Alice" },
    { value: "bob", label: "Bob" },
    { value: "carol", label: "Carol" },
  ], defaultValue: ["alice"] },
  { id: "region", label: "Region (cascader)", type: "cascader", options: [
    { value: "na", label: "North America" },
    { value: "eu", label: "Europe" },
    { value: "apac", label: "APAC" },
  ], defaultValue: ["eu"] },
  { id: "description", label: "Description (richText)", type: "richText", defaultValue: "## Summary\nA short markdown note." },
  { id: "secret", label: "Secret (password)", type: "password", defaultValue: "" },
];

const META = {
  protocolVersion: "2.7",
  requiredCapabilities: [
    "app.manifest",
    "app.navigation",
    "form.controls.extended",
    "form.controls.advanced",
  ],
};

export function FormControlsPage() {
  const [values, setValues] = useState<Record<string, unknown>>(() =>
    Object.fromEntries(
      ADVANCED_FIELDS.map((field) => [field.id, field.defaultValue]),
    ),
  );
  const gates: FormControlGateError[] = useMemo(
    () => checkFormCapabilitiesRaw(META, ADVANCED_FIELDS),
    [],
  );
  const [serialized, setSerialized] = useState<Record<string, unknown> | null>(null);

  return (
    <section className="space-y-6" aria-labelledby="form-controls-title">
      <div className="space-y-2">
        <h1 id="form-controls-title" className="text-3xl font-semibold tracking-tight">
          Form controls
        </h1>
        <p className="max-w-2xl text-sm leading-6 text-muted-foreground">
          The frozen §5 whitelist: base (input / select), 2.6 extended (textarea /
          switch / checkbox / radio / select-multiple), 2.7 advanced (cascader /
          checkboxGroup / richText / password) with defaultValue.
        </p>
      </div>

      {gates.length > 0 ? (
        <p role="alert" className="text-sm text-destructive">
          {gates.length} capability gate error(s): {gates.map((gate) => gate.code).join(", ")}
        </p>
      ) : (
        <p className="text-xs text-muted-foreground">
          Capability gate: protocol 2.7 + form.controls.extended/advanced — passes.
        </p>
      )}

      <FormControls
        fields={ADVANCED_FIELDS}
        values={values}
        onChange={(id, value) => setValues((prev) => ({ ...prev, [id]: value }))}
      />

      <div className="flex items-center gap-2">
        <Button
          type="button"
          size="sm"
          onClick={() => setSerialized(Object.fromEntries(Object.entries(values)))}
        >
          Serialize values
        </Button>
        {serialized !== null ? (
          <code className="rounded bg-muted px-2 py-1 text-xs">{JSON.stringify(serialized)}</code>
        ) : null}
      </div>
    </section>
  );
}
