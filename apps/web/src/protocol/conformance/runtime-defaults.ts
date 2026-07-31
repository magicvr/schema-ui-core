/**
 * runtime-defaults fixture adapter (schema-ui-docs v2.7.0).
 */

export type RuntimeDefaultsResult = Record<string, unknown>;

export function applyRuntimeDefaults(input: Record<string, unknown>): RuntimeDefaultsResult {
  const kind = input.kind;

  if (kind === "requestConfig") {
    const requiresNetwork = input.requiresNetwork === true;
    const baseURL = input.baseURL;
    if (requiresNetwork && (baseURL === undefined || baseURL === null || baseURL === "")) {
      return { ok: false, code: "MISSING_BASE_URL" };
    }
    return { ok: true };
  }

  if (kind === "defaults") {
    const target = input.target;
    const value = (input.value ?? {}) as Record<string, unknown>;
    if (target === "dataRef") {
      return {
        ok: true,
        value: {
          ...value,
          method: value.method ?? "GET",
        },
      };
    }
    if (target === "uploadAction") {
      return {
        ok: true,
        value: {
          ...value,
          method: value.method ?? "POST",
          retryPolicy: value.retryPolicy ?? "never",
          fieldName: value.fieldName ?? "file",
          multiple: value.multiple ?? false,
        },
      };
    }
    return { ok: false, code: "INVALID_DEFAULTS_TARGET" };
  }

  if (kind === "component") {
    const type = input.type as string;
    const installed = new Set((input.installedTypes as string[]) ?? []);
    if (!installed.has(type)) {
      return { ok: false, code: "UNKNOWN_COMPONENT_TYPE" };
    }
    const requiredProps = (input.requiredProps as string[]) ?? [];
    const props = (input.props as Record<string, unknown>) ?? {};
    for (const prop of requiredProps) {
      if (!Object.prototype.hasOwnProperty.call(props, prop)) {
        return { ok: false, code: "INVALID_COMPONENT", path: `props.${prop}` };
      }
    }
    return { ok: true };
  }

  if (kind === "formFieldInit") {
    const fields = (input.fields as Array<Record<string, unknown>>) ?? [];
    const recordValues = (input.recordValues as Record<string, unknown>) ?? {};
    const reactionWrites = (input.reactionWrites as Array<{ field: string; value: unknown }>) ?? [];
    const values: Record<string, unknown> = {};
    for (const field of fields) {
      const name = field.field as string;
      if (field.defaultValue !== undefined) {
        values[name] = field.defaultValue;
      }
    }
    for (const [key, value] of Object.entries(recordValues)) {
      values[key] = value;
    }
    for (const write of reactionWrites) {
      values[write.field] = write.value;
    }
    return { ok: true, values };
  }

  return { ok: false, code: "UNKNOWN_KIND" };
}
