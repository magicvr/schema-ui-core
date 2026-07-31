import { useState, type ComponentType, type ReactNode } from "react";

import { FormControls } from "@/renderer/form-controls.tsx";
import type { FormControlField } from "@/renderer/form-controls";
import {
  parseRenderNode,
  resolveFormReactions,
  type RenderFormNode,
  type RenderPageDocument,
  type RenderSectionNode,
  type RenderTableNode,
} from "@/renderer/render";

/**
 * R5 D-COMP minimal Renderer (resolve R4 F-002).
 *
 * Dispatch layer: parses a page document, applies the frozen $context
 * reaction engine to form field state, and renders whitelisted node types
 * through the components in this directory. Unknown node types fail closed.
 *
 * Scope: only form/section/table nodes are handled; the form control
 * whitelist itself is enforced by D-FORM (isWhitelistedFormControl /
 * checkFormCapabilities). Table data wiring stays with the example page.
 */

export interface RendererComponentProps {
  document: RenderPageDocument;
  context: Record<string, unknown>;
  /** Renders a table node; provided by the example page that owns data. */
  tableRenderer?: (node: RenderTableNode) => ReactNode;
  /** Overrides the default FormControls component (keeps field wiring local). */
  formComponent?: ComponentType<{
    fields: FormControlField[];
    values: Record<string, unknown>;
    onChange: (id: string, value: unknown) => void;
    fieldDisabled?: (id: string) => boolean;
  }>;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function FormView({
  node,
  context,
  formComponent,
}: {
  node: RenderFormNode;
  context: Record<string, unknown>;
  formComponent?: RendererComponentProps["formComponent"];
}) {
  const Component = formComponent ?? FormControls;
  const reaction = resolveFormReactions(node, context);
  const [values, setValues] = useState<Record<string, unknown>>({});

  const visibleFields = node.props.fields.filter((raw) => {
    if (!isRecord(raw) || typeof raw.id !== "string") {
      return false;
    }
    return reaction.state[raw.id]?.visible !== false;
  }) as unknown as FormControlField[];

  const fieldDisabled = (id: string) => reaction.state[id]?.disabled === true;

  return (
    <div className="space-y-3">
      <Component
        fields={visibleFields}
        values={values}
        onChange={(id, value) => setValues((prev) => ({ ...prev, [id]: value }))}
        fieldDisabled={fieldDisabled}
      />
      {reaction.errors.length > 0 ? (
        <ul role="alert" className="space-y-1 text-sm text-destructive">
          {reaction.errors.map((error, index) => (
            <li key={index}>
              {error.code}: {error.message}
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  );
}

function SectionView({
  node,
  context,
  tableRenderer,
  formComponent,
}: {
  node: RenderSectionNode;
  context: Record<string, unknown>;
  tableRenderer?: RendererComponentProps["tableRenderer"];
  formComponent?: RendererComponentProps["formComponent"];
}) {
  return (
    <div className="space-y-6">
      {node.children.map((child, index) => {
        const parsed = parseRenderNode(child, `body.children[${index}]`);
        if ("code" in parsed) {
          return (
            <p key={index} role="alert" className="text-sm text-destructive">
              {parsed.message}
            </p>
          );
        }
        if (parsed.type === "form") {
          return <FormView key={index} node={parsed} context={context} formComponent={formComponent} />;
        }
        if (parsed.type === "section") {
          return (
            <SectionView
              key={index}
              node={parsed}
              context={context}
              tableRenderer={tableRenderer}
              formComponent={formComponent}
            />
          );
        }
        return tableRenderer?.(parsed) ?? <p key={index} className="text-sm text-muted-foreground">table nodes are wired by the example page</p>;
      })}
    </div>
  );
}

export function RenderPage({
  document,
  context,
  tableRenderer,
  formComponent,
}: RendererComponentProps) {
  const body = parseRenderNode(document.body, "body");
  if ("code" in body) {
    return (
      <p role="alert" className="text-sm text-destructive">
        {body.message}
      </p>
    );
  }
  if (body.type === "form") {
    return <FormView node={body} context={context} formComponent={formComponent} />;
  }
  if (body.type === "section") {
    return (
      <SectionView node={body} context={context} tableRenderer={tableRenderer} formComponent={formComponent} />
    );
  }
  return tableRenderer?.(body) ?? <p className="text-sm text-muted-foreground">table nodes are wired by the example page</p>;
}
